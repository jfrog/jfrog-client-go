package tests

import (
	"strings"
	"testing"

	"github.com/jfrog/jfrog-client-go/pipelines/services"
	"github.com/stretchr/testify/assert"
)

const (
	integrationNamesPrefix = "jfrog_client_pipelines_integrations_tests"
	testsDummyRtUrl        = "https://pipelines.integration.com/artifactory/"
	testsDummyVcsUrl       = "https://non-existing-vcs.com/"
	testsDummyUser         = "nonexistinguser"
	testsDummyToken        = "nonexistingtoken"
	testsDummyApiKey       = "nonexistingkey"
)

func TestPipelinesIntegrations(t *testing.T) {
	initPipelinesTest(t)
	t.Run(services.GithubName, testCreateGithubIntegrationAndGetByName)
	t.Run(services.GithubEnterpriseName, testCreateGithubEnterpriseIntegration)
	t.Run(services.BitbucketName, testCreateBitbucketIntegration)
	t.Run(services.BitbucketServerName, testCreateBitbucketServerIntegration)
	t.Run(services.GitlabName, testCreateGitlabIntegration)
	t.Run(services.ArtifactoryName, testCreateArtifactoryIntegration)
}

func testCreateGithubIntegrationAndGetByName(t *testing.T) {
	name := getUniqueIntegrationName(services.GithubName)
	id, err := testsPipelinesIntegrationsService.CreateGithubIntegration(name, testsDummyToken)
	if err != nil {
		assert.NoError(t, err)
		return
	}
	defer deleteIntegrationAndAssert(t, id)
	getIntegrationAndAssert(t, id, name, services.GithubName)

	// Test get by name.
	integration, err := testsPipelinesIntegrationsService.GetIntegrationByName(name)
	if err != nil {
		assert.NoError(t, err)
		return
	}
	assert.Equal(t, name, integration.Name)
	assert.Equal(t, id, integration.Id)
}

func testCreateGithubEnterpriseIntegration(t *testing.T) {
	name := getUniqueIntegrationName(services.GithubEnterpriseName)
	id, err := testsPipelinesIntegrationsService.CreateGithubEnterpriseIntegration(name, testsDummyVcsUrl, testsDummyToken)
	if err != nil {
		assert.NoError(t, err)
		return
	}
	defer deleteIntegrationAndAssert(t, id)
	getIntegrationAndAssert(t, id, name, services.GithubEnterpriseName)
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

func testCreateBitbucketServerIntegration(t *testing.T) {
	name := getUniqueIntegrationName(services.BitbucketServerName)
	id, err := testsPipelinesIntegrationsService.CreateBitbucketServerIntegration(name, testsDummyVcsUrl, testsDummyUser, testsDummyToken)
	if err != nil {
		assert.NoError(t, err)
		return
	}
	defer deleteIntegrationAndAssert(t, id)
	getIntegrationAndAssert(t, id, name, services.BitbucketServerName)
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
	integration, err := testsPipelinesIntegrationsService.GetIntegrationById(id)
	if err != nil {
		assert.NoError(t, err)
		return
	}
	assert.NotNil(t, integration)
	assert.Equal(t, name, integration.Name)
	assert.Equal(t, integrationType, integration.MasterIntegrationName)
}

func getUniqueIntegrationName(integrationType string) string {
	return strings.Join([]string{integrationNamesPrefix, integrationType, getCustomRunId('_')}, "_")
}

func deleteIntegrationAndAssert(t *testing.T, id int) {
	err := testsPipelinesIntegrationsService.DeleteIntegration(id)
	assert.NoError(t, err)
}
