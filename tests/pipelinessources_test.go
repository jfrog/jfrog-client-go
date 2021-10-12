package tests

import (
	"testing"

	"github.com/jfrog/jfrog-client-go/pipelines/services"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/stretchr/testify/assert"
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
	if !utils.IsWindows() {
		// The following test modify a single Pipelines resource. Therefore, we limit the parallelism by running on windows (tests run on all OS parallel).
		t.Skip("Not running on Windows, skipping...")
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
	sourceId, err := testsPipelinesSourcesService.AddSource(integrationId, *PipelinesVcsRepoFullPath, *PipelinesVcsBranch, services.DefaultPipelinesFileFilter)
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
	assert.Equal(t, *PipelinesVcsRepoFullPath, source.RepositoryFullName)
	assert.Equal(t, *PipelinesVcsBranch, source.Branch)
	assert.Equal(t, services.DefaultPipelinesFileFilter, source.FileFilter)
}

func deleteSourceAndAssert(t *testing.T, id int) {
	err := testsPipelinesSourcesService.DeleteSource(id)
	assert.NoError(t, err)
}
