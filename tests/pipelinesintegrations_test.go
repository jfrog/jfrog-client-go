package tests

import (
	"github.com/jfrog/jfrog-client-go/pipelines/services"
	"github.com/stretchr/testify/assert"
	"strconv"
	"strings"
	"testing"
	"time"
)

const (
	integrationNamesPrefix = "jfrog_client_pipelines_integrations_tests"
	testsDummyRtUrl        = "https://pipelines.integration.com/artifactory/"
	testsDummyVcsUrl       = "https://non-existing-vcs.com/"
	testsDummyUser         = "nonexistinguser"
	testsDummyToken        = "nonexistingtoken"
	testsDummyApiKey       = "nonexistingkey"
)

func TestIntegrations(t *testing.T) {
	t.Run(services.GithubName, testCreateGithubIntegration)
	t.Run(services.BitbucketName, testCreateBitbucketIntegration)
	t.Run(services.GitlabName, testCreateGitlabIntegration)
	t.Run(services.ArtifactoryName, testCreateArtifactoryIntegration)
}

func testCreateGithubIntegration(t *testing.T) {
	name := getUniqueIntegrationName(services.GithubName)
	id, err := testsPipelinesIntegrationsService.CreateGithubIntegration(name, testsDummyToken)
	if err != nil {
		assert.NoError(t, err)
		return
	}
	defer deleteIntegrationAndAssert(t, id)
	getIntegrationAndAssert(t, id, name, services.GithubName)
}

func testCreateBitbucketIntegration(t *testing.T) {
	name := getUniqueIntegrationName(services.BitbucketName)
	id, err := testsPipelinesIntegrationsService.CreateBitbucketIntegration(name, testsDummyUser, testsDummyToken)
	if err != nil {
		assert.NoError(t, err)
		return
	}
	defer deleteIntegrationAndAssert(t, id)
	getIntegrationAndAssert(t, id, name, services.BitbucketName)
}

func testCreateGitlabIntegration(t *testing.T) {
	name := getUniqueIntegrationName(services.GitlabName)
	id, err := testsPipelinesIntegrationsService.CreateGitlabIntegration(name, testsDummyVcsUrl, testsDummyToken)
	if err != nil {
		assert.NoError(t, err)
		return
	}
	defer deleteIntegrationAndAssert(t, id)
	getIntegrationAndAssert(t, id, name, services.GitlabName)
}

func testCreateArtifactoryIntegration(t *testing.T) {
	name := getUniqueIntegrationName(services.ArtifactoryName)
	id, err := testsPipelinesIntegrationsService.CreateArtifactoryIntegration(name, testsDummyRtUrl, testsDummyUser, testsDummyApiKey)
	if err != nil {
		assert.NoError(t, err)
		return
	}
	defer deleteIntegrationAndAssert(t, id)
	getIntegrationAndAssert(t, id, name, services.ArtifactoryName)
}

func getIntegrationAndAssert(t *testing.T, id int, name, integrationType string) {
	integration, err := testsPipelinesIntegrationsService.GetIntegration(id)
	if err != nil {
		assert.NoError(t, err)
		return
	}
	assert.NotNil(t, integration)
	assert.Equal(t, name, integration.Name)
	assert.Equal(t, integrationType, integration.MasterIntegrationName)
}

func getUniqueIntegrationName(integrationType string) string {
	return strings.Join([]string{integrationNamesPrefix, integrationType, strconv.FormatInt(time.Now().Unix(), 10)}, "_")
}

func deleteIntegrationAndAssert(t *testing.T, id int) {
	err := testsPipelinesIntegrationsService.DeleteIntegration(id)
	assert.NoError(t, err)
}
