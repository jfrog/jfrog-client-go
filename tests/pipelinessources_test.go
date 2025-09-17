//go:build itest

package tests

import (
	"github.com/jfrog/jfrog-client-go/pipelines/services"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/stretchr/testify/assert"
	"testing"
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
	if !assert.NoError(t, err) {
		return
	}
	defer deleteIntegrationAndAssert(t, integrationId)

	// Create source with the above integration and assert.
	sourceId, err := testsPipelinesSourcesService.AddSource(integrationId, *PipelinesVcsRepoFullPath, *PipelinesVcsBranch, services.DefaultPipelinesFileFilter, "")
	if !assert.NoError(t, err) {
		return
	}
	defer deleteSourceAndAssert(t, sourceId)
	getSourceAndAssert(t, sourceId, integrationId)
}

func getSourceAndAssert(t *testing.T, sourceId, intId int) {
	source, err := testsPipelinesSourcesService.GetSource(sourceId)
	if !assert.NoError(t, err) {
		return
	}
	assert.NotNil(t, source)
	assert.Equal(t, intId, source.ProjectIntegrationId)
	assert.Equal(t, *PipelinesVcsRepoFullPath, source.RepositoryFullName)
	assert.Equal(t, *PipelinesVcsBranch, source.Branch)
	assert.Equal(t, services.DefaultPipelinesFileFilter, source.FileFilter)
}

func deleteSourceAndAssert(t *testing.T, id int) {
	pollingExecutor := &utils.RetryExecutor{
		MaxRetries:               10,
		RetriesIntervalMilliSecs: 3000,
		ErrorMessage:             "Failed deleting source. Trying again after sleep...",
		ExecutionHandler: func() (shouldRetry bool, err error) {
			err = testsPipelinesSourcesService.DeleteSource(id)
			return err != nil, err
		},
	}
	assert.NoError(t, pollingExecutor.Execute())
}
