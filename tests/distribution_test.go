package tests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"testing"
	"time"

	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/distribution/services"
	distributionServicesUtils "github.com/jfrog/jfrog-client-go/distribution/services/utils"
	"github.com/jfrog/jfrog-client-go/httpclient"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"github.com/stretchr/testify/assert"
)

type distributableDistributionStatus string
type receivedDistributionStatus string

const (
	open                           distributableDistributionStatus = "OPEN"
	readyForDistribution           distributableDistributionStatus = "READY_FOR_DISTRIBUTION"
	signed                         distributableDistributionStatus = "SIGNED"
	gpgKeyId                                                       = "234503"
	artifactoryGpgKeyCreatePattern                                 = `{"alias":"cli tests distribution key","public_key":"%s"}`
	bundleVersion                                                  = "10"
)

var httpClient *httpclient.HttpClient
var distHttpDetails httputils.HttpClientDetails

func TestDistribution(t *testing.T) {
	if *DistUrl == "" {
		t.Skip("Distribution is not being tested, skipping...")
	}
	initClients(t)
	sendGpgKeys(t)

	t.Run("createDelete", createDelete)
	t.Run("createUpdate", createUpdate)
	t.Run("createSignDistributeDelete", createSignDistributeDelete)
	t.Run("createSignSyncDistributeDelete", createSignSyncDistributeDelete)

	artifactoryCleanup(t)
	deleteGpgKeys(t)
}

func initClients(t *testing.T) {
	var err error
	distHttpDetails = GetDistDetails().CreateHttpClientDetails()
	httpClient, err = httpclient.ClientBuilder().Build()
	assert.NoError(t, err)
}

func initDistributionTest(t *testing.T, bundleName string) string {
	artifactoryCleanup(t)
	deleteTestBundle(t, bundleName)
	uploadDummyFile(t)
	return bundleName
}

func createDelete(t *testing.T) {
	bundleName := initDistributionTest(t, "client-test-bundle-1")

	// Create signed release bundle
	createBundleParams := services.NewCreateReleaseBundleParams(bundleName, bundleVersion)
	createBundleParams.SignImmediately = true
	createBundleParams.SpecFiles = []*utils.ArtifactoryCommonParams{{Pattern: RtTargetRepo + "b.in"}}
	err := testsBundleCreateService.CreateReleaseBundle(createBundleParams)
	assert.NoError(t, err)
	distributionResponse := getLocalBundle(t, bundleName, true)
	assert.Contains(t, [2]distributableDistributionStatus{readyForDistribution, signed}, distributionResponse.State)

	// Delete local release bundle
	deleteLocalBundleParams := services.NewDeleteReleaseBundleParams(bundleName, bundleVersion)
	err = testsBundleDeleteLocalService.DeleteDistribution(deleteLocalBundleParams)
	assert.NoError(t, err)
	distributionResponse = getLocalBundle(t, bundleName, false)
	assert.Nil(t, distributionResponse)
}

func createUpdate(t *testing.T) {
	bundleName := initDistributionTest(t, "client-test-bundle-2")

	// Create unsigned release bundle
	createBundleParams := services.NewCreateReleaseBundleParams(bundleName, bundleVersion)
	createBundleParams.Description = "Release bundle description 1"
	createBundleParams.ReleaseNotes = "Release notes 1"
	createBundleParams.SpecFiles = []*utils.ArtifactoryCommonParams{{Pattern: RtTargetRepo + "b.in"}}
	err := testsBundleCreateService.CreateReleaseBundle(createBundleParams)
	assert.NoError(t, err)
	distributionResponse := getLocalBundle(t, bundleName, true)
	assert.Equal(t, open, distributionResponse.State)
	assert.Equal(t, createBundleParams.Description, distributionResponse.Description)
	assert.Equal(t, createBundleParams.ReleaseNotes, distributionResponse.ReleaseNotes.Content)
	spec := distributionResponse.Spec

	// Update release bundle
	updateBundleParams := services.NewUpdateReleaseBundleParams(bundleName, bundleVersion)
	updateBundleParams.Description = "Release bundle description 2"
	updateBundleParams.ReleaseNotes = "Release notes 2"
	updateBundleParams.SpecFiles = []*utils.ArtifactoryCommonParams{{Pattern: RtTargetRepo + "test/a.in"}}
	err = testsBundleUpdateService.UpdateReleaseBundle(updateBundleParams)
	assert.NoError(t, err)
	distributionResponse = getLocalBundle(t, bundleName, true)
	assert.Equal(t, open, distributionResponse.State)
	assert.Equal(t, updateBundleParams.Description, distributionResponse.Description)
	assert.Equal(t, updateBundleParams.ReleaseNotes, distributionResponse.ReleaseNotes.Content)
	assert.NotEqual(t, spec, distributionResponse.Spec)

	// Delete local release bundle
	deleteLocalBundleParams := services.NewDeleteReleaseBundleParams(bundleName, bundleVersion)
	err = testsBundleDeleteLocalService.DeleteDistribution(deleteLocalBundleParams)
	assert.NoError(t, err)
	distributionResponse = getLocalBundle(t, bundleName, false)
	assert.Nil(t, distributionResponse)
}

func createSignDistributeDelete(t *testing.T) {
	bundleName := initDistributionTest(t, "client-test-bundle-3")

	// Create unsigned release bundle
	createBundleParams := services.NewCreateReleaseBundleParams(bundleName, bundleVersion)
	createBundleParams.SpecFiles = []*utils.ArtifactoryCommonParams{{Pattern: RtTargetRepo + "b.in"}}
	err := testsBundleCreateService.CreateReleaseBundle(createBundleParams)
	assert.NoError(t, err)
	distributionResponse := getLocalBundle(t, bundleName, true)
	assert.Equal(t, open, distributionResponse.State)

	// Sign release bundle
	signBundleParams := services.NewSignBundleParams(bundleName, bundleVersion)
	err = testsBundleSignService.SignReleaseBundle(signBundleParams)
	assert.NoError(t, err)
	distributionResponse = getLocalBundle(t, bundleName, true)
	assert.Contains(t, [2]distributableDistributionStatus{readyForDistribution, signed}, distributionResponse.State)

	// Distribute release bundle
	distributeBundleParams := services.NewDistributeReleaseBundleParams(bundleName, bundleVersion)
	distributeBundleParams.DistributionRules = []*distributionServicesUtils.DistributionCommonParams{{SiteName: "*"}}
	err = testsBundleDistributeService.Distribute(distributeBundleParams)
	assert.NoError(t, err)
	waitForDistribution(t, bundleName)

	// Assert release bundle in "completed" status
	distributionStatusParams := services.DistributionStatusParams{
		Name:    bundleName,
		Version: bundleVersion,
	}
	response, err := testsBundleDistributionStatusService.GetStatus(distributionStatusParams)
	assert.NoError(t, err)
	assert.Equal(t, services.Completed, (*response)[0].Status)

	// Delete release bundle
	err = deleteTestBundle(t, bundleName)
	assert.NoError(t, err)
	waitForDeletion(t, bundleName)
}

func createSignSyncDistributeDelete(t *testing.T) {
	bundleName := initDistributionTest(t, "client-test-bundle-4")

	// Create unsigned release bundle
	createBundleParams := services.NewCreateReleaseBundleParams(bundleName, bundleVersion)
	createBundleParams.SpecFiles = []*utils.ArtifactoryCommonParams{{Pattern: RtTargetRepo + "b.in"}}
	err := testsBundleCreateService.CreateReleaseBundle(createBundleParams)
	assert.NoError(t, err)
	distributionResponse := getLocalBundle(t, bundleName, true)
	assert.Equal(t, open, distributionResponse.State)

	// Sign release bundle
	signBundleParams := services.NewSignBundleParams(bundleName, bundleVersion)
	err = testsBundleSignService.SignReleaseBundle(signBundleParams)
	assert.NoError(t, err)
	distributionResponse = getLocalBundle(t, bundleName, true)
	assert.Contains(t, [2]distributableDistributionStatus{readyForDistribution, signed}, distributionResponse.State)

	// Distribute release bundle
	distributeBundleParams := services.NewDistributeReleaseBundleParams(bundleName, bundleVersion)
	distributeBundleParams.DistributionRules = []*distributionServicesUtils.DistributionCommonParams{{SiteName: "*"}}
	testsBundleDistributeService.Sync = true
	err = testsBundleDistributeService.Distribute(distributeBundleParams)
	testsBundleDistributeService.Sync = false
	assert.NoError(t, err)

	// Assert release bundle in "completed" status
	distributionStatusParams := services.DistributionStatusParams{
		Name:    bundleName,
		Version: bundleVersion,
	}
	response, err := testsBundleDistributionStatusService.GetStatus(distributionStatusParams)
	assert.NoError(t, err)
	assert.Equal(t, services.Completed, (*response)[0].Status)

	// Delete release bundle
	err = deleteTestBundle(t, bundleName)
	assert.NoError(t, err)
	waitForDeletion(t, bundleName)
}

// Send GPG keys to Distribution and Artifactory to allow signing of release bundles
func sendGpgKeys(t *testing.T) {
	// Read gpg public and private key
	publicKey, err := ioutil.ReadFile(filepath.Join(getTestDataPath(), "public.key"))
	assert.NoError(t, err)
	privateKey, err := ioutil.ReadFile(filepath.Join(getTestDataPath(), "private.key"))
	assert.NoError(t, err)

	err = testsBundleSetSigningKeyService.SetSigningKey(services.NewSetSigningKeyParams(string(publicKey), string(privateKey)))
	assert.NoError(t, err)

	// Send public key to Artifactory
	content := fmt.Sprintf(artifactoryGpgKeyCreatePattern, publicKey)
	resp, body, err := httpClient.SendPost(GetRtDetails().GetUrl()+"api/security/keys/trusted", []byte(content), distHttpDetails)
	assert.NoError(t, err)
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusConflict {
		t.Error(resp.Status)
		t.Error(string(body))
	}
}

// Delete GPG key from Artifactory to clean up the test environment
func deleteGpgKeys(t *testing.T) {
	// Delete public key from Artifactory
	resp, body, err := httpClient.SendDelete(GetRtDetails().GetUrl()+"api/security/keys/trusted/"+gpgKeyId, nil, distHttpDetails)
	assert.NoError(t, err)
	if resp.StatusCode != http.StatusNoContent {
		t.Error(resp.Status)
		t.Error(string(body))
	}
}

func deleteTestBundle(t *testing.T, bundleName string) error {
	deleteBundleParams := services.NewDeleteReleaseBundleParams(bundleName, bundleVersion)
	deleteBundleParams.DeleteFromDistribution = true
	deleteBundleParams.DistributionRules = []*distributionServicesUtils.DistributionCommonParams{{SiteName: "*"}}
	return testsBundleDeleteRemoteService.DeleteDistribution(deleteBundleParams)
}

type receivedResponse struct {
	Id     string                     `json:"id,omitempty"`
	Status receivedDistributionStatus `json:"status,omitempty"`
}

type receivedResponses struct {
	receivedResponses []receivedResponse
}

// Wait for distribution of a release bundle
func waitForDistribution(t *testing.T, bundleName string) {
	distributionStatusParams := services.DistributionStatusParams{
		Name:    bundleName,
		Version: bundleVersion,
	}
	for i := 0; i < 120; i++ {
		response, err := testsBundleDistributionStatusService.GetStatus(distributionStatusParams)
		assert.NoError(t, err)
		assert.Len(t, *response, 1)

		switch (*response)[0].Status {
		case services.Completed:
			return
		case services.Failed:
			t.Error("Distribution failed for " + bundleName + "/" + bundleVersion)
			return
		case services.InProgress, services.NotDistributed:
			// Wait
		}
		t.Log("Waiting for " + bundleName + "/" + bundleVersion + "...")
		time.Sleep(time.Second)
	}
	t.Error("Timeout for release bundle distribution " + bundleName + "/" + bundleVersion)
}

// Wait for deletion of a release bundle
func waitForDeletion(t *testing.T, bundleName string) {
	for i := 0; i < 120; i++ {
		resp, body, _, err := httpClient.SendGet(GetDistDetails().GetUrl()+"api/v1/release_bundle/"+bundleName+"/"+bundleVersion+"/distribution", true, distHttpDetails)
		assert.NoError(t, err)
		if resp.StatusCode == http.StatusNotFound {
			return
		}
		if resp.StatusCode != http.StatusOK {
			t.Error(resp.Status)
			t.Error(string(body))
			return
		}
		t.Log("Waiting for distribution deletion " + bundleName + "/" + bundleVersion + "...")
		time.Sleep(time.Second)
	}
	t.Error("Timeout for release bundle deletion " + bundleName + "/" + bundleVersion)
}

type distributableResponse struct {
	Name         string                          `json:"name,omitempty"`
	Version      string                          `json:"version,omitempty"`
	State        distributableDistributionStatus `json:"state,omitempty"`
	Description  string                          `json:"description,omitempty"`
	ReleaseNotes releaseNotesResponse            `json:"release_notes,omitempty"`
	Spec         interface{}                     `json:"spec,omitempty"`
}

type releaseNotesResponse struct {
	Content string `json:"content,omitempty"`
	Syntax  string `json:"syntax,omitempty"`
}

func getLocalBundle(t *testing.T, bundleName string, expectExist bool) *distributableResponse {
	resp, body, _, err := httpClient.SendGet(GetDistDetails().GetUrl()+"api/v1/release_bundle/"+bundleName+"/"+bundleVersion, true, distHttpDetails)
	assert.NoError(t, err)
	if !expectExist && resp.StatusCode == http.StatusNotFound {
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		t.Error(resp.Status)
		t.Error(string(body))
		return nil
	}
	response := &distributableResponse{}
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)
	return response
}
