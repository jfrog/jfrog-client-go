package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	artifactoryServices "github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/distribution/services"
	distributionServicesUtils "github.com/jfrog/jfrog-client-go/distribution/services/utils"
	"github.com/jfrog/jfrog-client-go/http/httpclient"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"github.com/stretchr/testify/assert"
)

type distributableDistributionStatus string

const (
	// Release bundle created and open for changes:
	open distributableDistributionStatus = "OPEN"
	// Release bundle is signed, but not stored:
	signed distributableDistributionStatus = "SIGNED"
	// Release bundle is signed and stored, but not scanned by Xray:
	stored distributableDistributionStatus = "STORED"
	// Release bundle is signed, stored and scanned by Xray:
	readyForDistribution distributableDistributionStatus = "READY_FOR_DISTRIBUTION"

	gpgKeyAlias                    = "client tests distribution key"
	artifactoryGpgKeyCreatePattern = `{"alias":"` + gpgKeyAlias + `","public_key":"%s"}`
	bundleVersion                  = "10"
)

var httpClient *httpclient.HttpClient
var distHttpDetails httputils.HttpClientDetails

func TestDistributionServices(t *testing.T) {
	initDistributionTest(t)
	initClients(t)
	sendGpgKeys(t)

	// Local release bundle tests
	t.Run("createDelete", createDelete)
	t.Run("createUpdate", createUpdate)
	t.Run("createWithProps", createWithProps)

	// Remote release bundle tests
	t.Run("createSignDistributeDelete", createSignDistributeDelete)
	t.Run("createSignSyncDistributeDelete", createSignSyncDistributeDelete)
	t.Run("createDistributeMapping", createDistributeMapping)
	t.Run("createDistributeMappingPlaceholder", createDistributeMappingPlaceholder)

	artifactoryCleanup(t)
	deleteGpgKeys(t)
}

func initDistributionTest(t *testing.T) {
	if !*TestDistribution {
		t.Skip("Skipping distribution test. To run distribution test add the '-test.distribution=true' option.")
	}
}

func initClients(t *testing.T) {
	var err error
	distHttpDetails = GetDistDetails().CreateHttpClientDetails()
	httpClient, err = httpclient.ClientBuilder().Build()
	assert.NoError(t, err)
}

func setupDistributionTest(t *testing.T, bundleName string) string {
	artifactoryCleanup(t)
	uploadDummyFile(t)
	return bundleName
}

func initLocalDistributionTest(t *testing.T, bundleName string) string {
	return setupDistributionTest(t, bundleName)
}

func initRemoteDistributionTest(t *testing.T, bundleName string) string {
	testsBundleDistributeService.Sync = false
	return setupDistributionTest(t, bundleName)
}

func createDelete(t *testing.T) {
	bundleName := initLocalDistributionTest(t, "client-test-bundle-"+getRunId())
	defer deleteLocalBundle(t, bundleName, true)

	// Create signed release bundle
	createBundleParams := services.NewCreateReleaseBundleParams(bundleName, bundleVersion)
	createBundleParams.SignImmediately = true
	createBundleParams.SpecFiles = []*utils.CommonParams{{Pattern: getRtTargetRepo() + "b.in"}}
	summary, err := testsBundleCreateService.CreateReleaseBundle(createBundleParams)
	assert.NotNil(t, summary)
	verifyValidSha256(t, summary.GetSha256())
	assert.NoError(t, err)
	distributionResponse := getLocalBundle(t, bundleName, true)
	if assert.NotNil(t, distributionResponse) {
		assertReleaseBundleSigned(t, distributionResponse.State)
	}
}

func createUpdate(t *testing.T) {
	bundleName := initLocalDistributionTest(t, "client-test-bundle-"+getRunId())
	defer deleteLocalBundle(t, bundleName, true)

	// Create release bundle params
	createBundleParams := services.NewCreateReleaseBundleParams(bundleName, bundleVersion)
	createBundleParams.Description = "Release bundle description 1"
	createBundleParams.ReleaseNotes = "Release notes 1"
	createBundleParams.SpecFiles = []*utils.CommonParams{{Pattern: getRtTargetRepo() + "b.in"}}

	// Test DryRun first
	err := createDryRun(createBundleParams)
	if err != nil {
		assert.NoError(t, err)
		return
	}
	// Verify was not created.
	getLocalBundle(t, bundleName, false)

	// Create unsigned release bundle
	summary, err := testsBundleCreateService.CreateReleaseBundle(createBundleParams)
	if err != nil {
		assert.NoError(t, err)
		return
	}
	assert.Nil(t, summary)
	distributionResponse := assertCreatedLocalBundle(t, bundleName, createBundleParams)
	spec := distributionResponse.BundleSpec

	// Create update release bundle params
	updateBundleParams := services.NewUpdateReleaseBundleParams(bundleName, bundleVersion)
	updateBundleParams.Description = "Release bundle description 2"
	updateBundleParams.ReleaseNotes = "Release notes 2"
	updateBundleParams.SpecFiles = []*utils.CommonParams{{Pattern: getRtTargetRepo() + "test/a.in"}}
	updateBundleParams.SignImmediately = false
	// Test DryRun first
	err = updateDryRun(updateBundleParams)
	if err != nil {
		assert.NoError(t, err)
		return
	}
	// Verify the release bundle was not updated.
	assertCreatedLocalBundle(t, bundleName, createBundleParams)

	summary, err = testsBundleUpdateService.UpdateReleaseBundle(updateBundleParams)
	if err != nil {
		assert.NoError(t, err)
		return
	}
	assert.Nil(t, summary)
	distributionResponse = getLocalBundle(t, bundleName, true)
	assert.Equal(t, open, distributionResponse.State)
	assert.Equal(t, updateBundleParams.Description, distributionResponse.Description)
	assert.Equal(t, updateBundleParams.ReleaseNotes, distributionResponse.ReleaseNotes.Content)
	assert.NotEqual(t, spec, distributionResponse.BundleSpec)
}

func assertCreatedLocalBundle(t *testing.T, bundleName string, createBundleParams services.CreateReleaseBundleParams) *distributableResponse {
	distributionResponse := getLocalBundle(t, bundleName, true)
	assert.Equal(t, open, distributionResponse.State)
	assert.Equal(t, createBundleParams.Description, distributionResponse.Description)
	assert.Equal(t, createBundleParams.ReleaseNotes, distributionResponse.ReleaseNotes.Content)
	return distributionResponse
}

func createDryRun(createBundleParams services.CreateReleaseBundleParams) error {
	defer setServicesToDryRunFalse()
	testsBundleCreateService.DryRun = true
	_, err := testsBundleCreateService.CreateReleaseBundle(createBundleParams)
	return err
}

func updateDryRun(updateBundleParams services.UpdateReleaseBundleParams) error {
	defer setServicesToDryRunFalse()
	testsBundleUpdateService.DryRun = true
	_, err := testsBundleUpdateService.UpdateReleaseBundle(updateBundleParams)
	return err
}

func distributeDryRun(distributionParams services.DistributionParams) error {
	defer setServicesToDryRunFalse()
	testsBundleDistributeService.DryRun = true
	testsBundleDistributeService.AutoCreateRepo = true
	return testsBundleDistributeService.Distribute(distributionParams)
}

func setServicesToDryRunFalse() {
	testsBundleCreateService.DryRun = false
	testsBundleUpdateService.DryRun = false
	testsBundleDistributeService.DryRun = false
}

func createWithProps(t *testing.T) {
	bundleName := initLocalDistributionTest(t, "client-test-bundle-"+getRunId())
	defer deleteLocalBundle(t, bundleName, true)

	// Create release bundle with properties
	targetProps, err := utils.ParseProperties("key1=value1;key2=value2,value3")
	assert.NoError(t, err)
	createBundleParams := services.NewCreateReleaseBundleParams(bundleName, bundleVersion)
	createBundleParams.SpecFiles = []*utils.CommonParams{{
		Pattern:     getRtTargetRepo() + "b.in",
		TargetProps: targetProps,
	}}
	summary, err := testsBundleCreateService.CreateReleaseBundle(createBundleParams)
	assert.Nil(t, summary)
	assert.NoError(t, err)

	// Check results
	distributionResponse := getLocalBundle(t, bundleName, true)
	addedProps := distributionResponse.BundleSpec.Queries[0].AddedProps
	assert.Len(t, addedProps, 2)

	// Populate prop1Values and prop2Values
	var prop1Values []string
	var prop2Values []string
	if addedProps[0].Key == "key1" {
		assert.Equal(t, "key2", addedProps[1].Key)
		prop1Values = addedProps[0].Values
		prop2Values = addedProps[1].Values
	} else if addedProps[0].Key == "key2" {
		assert.Equal(t, "key1", addedProps[1].Key)
		prop1Values = addedProps[1].Values
		prop2Values = addedProps[0].Values
	} else {
		assert.Fail(t, "Unexpected key", addedProps[0].Key)
	}

	// Check prop1Values and prop2Values
	assert.Len(t, prop1Values, 1)
	assert.Len(t, prop2Values, 2)
	if len(prop1Values) == 1 {
		assert.Equal(t, "value1", prop1Values[0])
	}
	if len(prop2Values) == 2 {
		assert.Equal(t, "value2", prop2Values[0])
		assert.Equal(t, "value3", prop2Values[1])
	}
}

func createSignDistributeDelete(t *testing.T) {
	bundleName := initRemoteDistributionTest(t, "client-test-bundle-"+getRunId())
	defer deleteRemoteAndLocalBundle(t, bundleName, true)

	// Create unsigned release bundle
	createBundleParams := services.NewCreateReleaseBundleParams(bundleName, bundleVersion)
	createBundleParams.SpecFiles = []*utils.CommonParams{{Pattern: getRtTargetRepo() + "b.in"}}
	summary, err := testsBundleCreateService.CreateReleaseBundle(createBundleParams)
	assert.NoError(t, err)
	assert.Nil(t, summary)
	distributionResponse := getLocalBundle(t, bundleName, true)
	assert.Equal(t, open, distributionResponse.State)

	// Sign release bundle
	signBundleParams := services.NewSignBundleParams(bundleName, bundleVersion)
	summary, err = testsBundleSignService.SignReleaseBundle(signBundleParams)
	assert.NoError(t, err)
	assert.NotNil(t, summary)
	verifyValidSha256(t, summary.GetSha256())
	distributionResponse = getLocalBundle(t, bundleName, true)
	assertReleaseBundleSigned(t, distributionResponse.State)

	// Create distribute params.
	distributeBundleParams := services.NewDistributeReleaseBundleParams(bundleName, bundleVersion)
	distributeBundleParams.DistributionRules = []*distributionServicesUtils.DistributionCommonParams{{SiteName: "*"}}

	// Create response params.
	distributionStatusParams := services.DistributionStatusParams{
		Name:    bundleName,
		Version: bundleVersion,
	}

	// Test DryRun first.
	err = distributeDryRun(distributeBundleParams)
	if err != nil {
		assert.NoError(t, err)
		return
	}
	// Assert release bundle not in distribution yet.
	response, err := testsBundleDistributionStatusService.GetStatus(distributionStatusParams)
	assert.NoError(t, err)
	assert.Len(t, *response, 0)

	// Distribute release bundle
	testsBundleDistributeService.AutoCreateRepo = true
	err = testsBundleDistributeService.Distribute(distributeBundleParams)
	assert.NoError(t, err)
	waitForDistribution(t, bundleName)

	// Assert release bundle in "completed" status
	response, err = testsBundleDistributionStatusService.GetStatus(distributionStatusParams)
	assert.NoError(t, err)
	assert.Equal(t, services.Completed, (*response)[0].Status)
}

func createSignSyncDistributeDelete(t *testing.T) {
	bundleName := initRemoteDistributionTest(t, "client-test-bundle-"+getRunId())
	defer deleteRemoteAndLocalBundle(t, bundleName, true)

	// Create unsigned release bundle
	createBundleParams := services.NewCreateReleaseBundleParams(bundleName, bundleVersion)
	createBundleParams.SpecFiles = []*utils.CommonParams{{Pattern: getRtTargetRepo() + "b.in"}}
	summary, err := testsBundleCreateService.CreateReleaseBundle(createBundleParams)
	assert.NoError(t, err)
	assert.Nil(t, summary)
	distributionResponse := getLocalBundle(t, bundleName, true)
	assert.Equal(t, open, distributionResponse.State)

	// Sign release bundle
	signBundleParams := services.NewSignBundleParams(bundleName, bundleVersion)
	summary, err = testsBundleSignService.SignReleaseBundle(signBundleParams)
	assert.NoError(t, err)
	assert.NotNil(t, summary)
	verifyValidSha256(t, summary.GetSha256())
	distributionResponse = getLocalBundle(t, bundleName, true)
	assertReleaseBundleSigned(t, distributionResponse.State)

	// Distribute release bundle
	distributeBundleParams := services.NewDistributeReleaseBundleParams(bundleName, bundleVersion)
	distributeBundleParams.DistributionRules = []*distributionServicesUtils.DistributionCommonParams{{SiteName: "*"}}
	testsBundleDistributeService.Sync = true
	testsBundleDistributeService.AutoCreateRepo = true
	err = testsBundleDistributeService.Distribute(distributeBundleParams)
	assert.NoError(t, err)

	// Assert release bundle in "completed" status
	distributionStatusParams := services.DistributionStatusParams{
		Name:    bundleName,
		Version: bundleVersion,
	}
	response, err := testsBundleDistributionStatusService.GetStatus(distributionStatusParams)
	assert.NoError(t, err)
	assert.Equal(t, services.Completed, (*response)[0].Status)
}

func createDistributeMapping(t *testing.T) {
	bundleName := initRemoteDistributionTest(t, "client-test-bundle-"+getRunId())
	defer deleteRemoteAndLocalBundle(t, bundleName, true)

	// Create release bundle with path mapping from <RtTargetRepo>/b.in to <RtTargetRepo>/b.out
	createBundleParams := services.NewCreateReleaseBundleParams(bundleName, bundleVersion)
	createBundleParams.SpecFiles = []*utils.CommonParams{{Pattern: getRtTargetRepo() + "b.in", Target: getRtTargetRepo() + "b.out"}}
	createBundleParams.SignImmediately = true
	summary, err := testsBundleCreateService.CreateReleaseBundle(createBundleParams)
	assert.NoError(t, err)
	assert.NotNil(t, summary)
	verifyValidSha256(t, summary.GetSha256())

	// Distribute release bundle
	distributeBundleParams := services.NewDistributeReleaseBundleParams(bundleName, bundleVersion)
	distributeBundleParams.DistributionRules = []*distributionServicesUtils.DistributionCommonParams{{SiteName: "*"}}
	testsBundleDistributeService.Sync = true
	// On distribution with path mapping, the target repository cannot be auto-created
	testsBundleDistributeService.AutoCreateRepo = false
	err = testsBundleDistributeService.Distribute(distributeBundleParams)
	assert.NoError(t, err)

	// Make sure <RtTargetRepo>/b.out does exist in Artifactory
	searchParams := artifactoryServices.NewSearchParams()
	searchParams.Pattern = getRtTargetRepo() + "b.out"
	reader, err := testsSearchService.Search(searchParams)
	assert.NoError(t, err)
	readerCloseAndAssert(t, reader)
	length, err := reader.Length()
	assert.NoError(t, err)
	assert.Equal(t, 1, length)
}

func createDistributeMappingPlaceholder(t *testing.T) {
	bundleName := initRemoteDistributionTest(t, "client-test-bundle-"+getRunId())
	defer deleteRemoteAndLocalBundle(t, bundleName, true)

	// Create release bundle with path mapping from <RtTargetRepo>/b.in to <RtTargetRepo>/b.out
	createBundleParams := services.NewCreateReleaseBundleParams(bundleName, bundleVersion)
	createBundleParams.SpecFiles = []*utils.CommonParams{{Pattern: "(" + getRtTargetRepo() + ")" + "(*).in", Target: "{1}{2}.out"}}
	createBundleParams.SignImmediately = true
	summary, err := testsBundleCreateService.CreateReleaseBundle(createBundleParams)
	assert.NoError(t, err)
	assert.NotNil(t, summary)
	verifyValidSha256(t, summary.GetSha256())

	// Distribute release bundle
	distributeBundleParams := services.NewDistributeReleaseBundleParams(bundleName, bundleVersion)
	distributeBundleParams.DistributionRules = []*distributionServicesUtils.DistributionCommonParams{{SiteName: "*"}}
	testsBundleDistributeService.Sync = true
	// On distribution with path mapping, the target repository cannot be auto-created
	testsBundleDistributeService.AutoCreateRepo = false
	err = testsBundleDistributeService.Distribute(distributeBundleParams)
	assert.NoError(t, err)

	// Make sure <RtTargetRepo>/b.out does exist in Artifactory
	searchParams := artifactoryServices.NewSearchParams()
	searchParams.Pattern = getRtTargetRepo() + "b.out"
	reader, err := testsSearchService.Search(searchParams)
	assert.NoError(t, err)
	readerCloseAndAssert(t, reader)
	length, err := reader.Length()
	assert.NoError(t, err)
	assert.Equal(t, 1, length)
}

// Send GPG keys to Distribution and Artifactory to allow signing of release bundles
func sendGpgKeys(t *testing.T) {
	// Read gpg public and private key
	publicKey, err := os.ReadFile(filepath.Join(getTestDataPath(), "public.key"))
	assert.NoError(t, err)
	privateKey, err := os.ReadFile(filepath.Join(getTestDataPath(), "private.key"))
	assert.NoError(t, err)

	err = testsBundleSetSigningKeyService.SetSigningKey(services.NewSetSigningKeyParams(string(publicKey), string(privateKey)))
	assert.NoError(t, err)

	// Send public key to Artifactory
	content := fmt.Sprintf(artifactoryGpgKeyCreatePattern, publicKey)
	resp, body, err := httpClient.SendPost(GetRtDetails().GetUrl()+"api/security/keys/trusted", []byte(content), distHttpDetails, "")
	assert.NoError(t, err)
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusConflict {
		t.Error(resp.Status)
		t.Error(string(body))
	}
}

// Delete GPG key from Artifactory to clean up the test environment
func deleteGpgKeys(t *testing.T) {
	// Delete public key from Artifactory
	gpgKeyId := getGpgKeyId(t)
	if gpgKeyId == "" {
		return
	}
	resp, body, err := httpClient.SendDelete(GetRtDetails().GetUrl()+"api/security/keys/trusted/"+gpgKeyId, nil, distHttpDetails, "")
	assert.NoError(t, err)
	if resp.StatusCode != http.StatusNoContent {
		t.Error(resp.Status)
		t.Error(string(body))
	}
}

// Get GPG key ID created in the tests
func getGpgKeyId(t *testing.T) string {
	resp, body, _, err := httpClient.SendGet(GetRtDetails().GetUrl()+"api/security/keys/trusted", true, distHttpDetails, "")
	assert.NoError(t, err)
	if resp.StatusCode != http.StatusOK {
		t.Error(resp.Status)
		t.Error(string(body))
		return ""
	}
	responses := &gpgKeysResponse{}
	err = json.Unmarshal(body, &responses)
	assert.NoError(t, err)
	for _, gpgKeyResponse := range responses.Keys {
		if gpgKeyResponse.Alias == gpgKeyAlias {
			return gpgKeyResponse.Kid
		}
	}
	return ""
}

func assertReleaseBundleSigned(t *testing.T, status distributableDistributionStatus) {
	assert.Contains(t, []distributableDistributionStatus{signed, stored, readyForDistribution}, status)
}

// Wait for distribution of a release bundle
func waitForDistribution(t *testing.T, bundleName string) {
	distributionStatusParams := services.DistributionStatusParams{
		Name:    bundleName,
		Version: bundleVersion,
	}
	for i := 0; i < 120; i++ {
		response, err := testsBundleDistributionStatusService.GetStatus(distributionStatusParams)
		if assert.NoError(t, err) {
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
	}
	t.Error("Timeout for release bundle distribution " + bundleName + "/" + bundleVersion)
}

type distributableResponse struct {
	distributionServicesUtils.ReleaseBundleBody
	Name    string                          `json:"name,omitempty"`
	Version string                          `json:"version,omitempty"`
	State   distributableDistributionStatus `json:"state,omitempty"`
}

type gpgKeysResponse struct {
	Keys []gpgKeyResponse `json:"keys,omitempty"`
}

type gpgKeyResponse struct {
	Kid   string `json:"kid,omitempty"`
	Alias string `json:"alias,omitempty"`
}

func getLocalBundle(t *testing.T, bundleName string, expectExist bool) *distributableResponse {
	resp, body, _, err := httpClient.SendGet(GetDistDetails().GetUrl()+"api/v1/release_bundle/"+bundleName+"/"+bundleVersion, true, distHttpDetails, "")
	assert.NoError(t, err)
	if !expectExist {
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
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

func deleteLocalBundle(t *testing.T, bundleName string, assertDeletion bool) {
	deleteLocalBundleParams := services.NewDeleteReleaseBundleParams(bundleName, bundleVersion)
	testsBundleDeleteLocalService.Sync = true
	err := testsBundleDeleteLocalService.DeleteDistribution(deleteLocalBundleParams)
	if !assertDeletion {
		return
	}
	assert.NoError(t, err)
	distributionResponse := getLocalBundle(t, bundleName, false)
	assert.Nil(t, distributionResponse)
}

func deleteRemoteAndLocalBundle(t *testing.T, bundleName string, assertDeletion bool) {
	deleteBundleParams := services.NewDeleteReleaseBundleParams(bundleName, bundleVersion)
	// Delete also local release bundle
	deleteBundleParams.DeleteFromDistribution = true
	deleteBundleParams.DistributionRules = []*distributionServicesUtils.DistributionCommonParams{{SiteName: "*"}}
	deleteBundleParams.Sync = true
	err := testsBundleDeleteRemoteService.DeleteDistribution(deleteBundleParams)
	artifactoryCleanup(t)
	if !assertDeletion {
		return
	}
	assert.NoError(t, err)
}
