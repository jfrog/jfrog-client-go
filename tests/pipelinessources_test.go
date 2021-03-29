package tests

import (
	"github.com/jfrog/jfrog-client-go/pipelines/services"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	testsRepo   = "ecoswamp/repo-for-pipelines-tests"
	testsBranch = "main"
)

func TestPipelinesSources(t *testing.T) {
	initPipelinesTest(t)
	t.Run("addPipelineSource", testAddPipelineSource)
}

func testAddPipelineSource(t *testing.T) {
	if *PipelinesVcsToken == "" {
		assert.NotEmpty(t, *PipelinesVcsToken, "cannot run pipelines tests without vcs token configured")
		return
	}
	// Create integration with provided token.
	integrationName := getUniqueIntegrationName(services.GithubName)
	integrationId, err := testsPipelinesIntegrationsService.CreateGithubIntegration(integrationName, *PipelinesVcsToken)
	if err != nil {
		assert.NoError(t, err)
		return
	}
	defer deleteIntegrationAndAssert(t, integrationId)

	// Create source with the above integration and assert.
	sourceId, err := testsPipelinesSourcesService.AddPipelineSource(integrationId, testsRepo, testsBranch, services.DefaultPipelinesFileFilter)
	if err != nil {
		assert.NoError(t, err)
		return
	}
	defer deleteSourceAndAssert(t, sourceId)
	getSourceAndAssert(t, sourceId, integrationId)
}

func getSourceAndAssert(t *testing.T, sourceId, intId int) {
	source, err := testsPipelinesSourcesService.GetSource(sourceId)
	if err != nil {
		assert.NoError(t, err)
		return
	}
	assert.NotNil(t, source)
	assert.Equal(t, intId, source.ProjectIntegrationId)
	assert.Equal(t, testsRepo, source.RepositoryFullName)
	assert.Equal(t, testsBranch, source.Branch)
	assert.Equal(t, services.DefaultPipelinesFileFilter, source.FileFilter)
}

func deleteSourceAndAssert(t *testing.T, id int) {
	err := testsPipelinesSourcesService.DeleteSource(id)
	assert.NoError(t, err)
}
