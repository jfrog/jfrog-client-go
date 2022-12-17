package tests

import (
	"fmt"
	"github.com/jfrog/jfrog-client-go/pipelines/services"
	"github.com/stretchr/testify/assert"
	"strconv"
	"strings"
	"testing"
	"time"
)

const (
	defaultMaxWaitMinutes    = 45 * time.Minute // 45 minutes
	defaultSyncSleepInterval = 5 * time.Second  // 5 seconds
)

func TestPipelinesRunService(t *testing.T) {
	initPipelinesTest(t)
	t.Run("test trigger pipeline resource sync", testTriggerSync)
	t.Run("test get sync status", testGetSyncStatus)
	t.Run("test get run status", testGetRunStatus)
}

func testTriggerSync(t *testing.T) {
	if *PipelinesVcsToken == "" {
		assert.NotEmpty(t, *PipelinesVcsToken, "cannot run pipelines tests without vcs token configured")
		return
	}
	// Create integration with provided token.
	unixTime := time.Now().Unix()
	timeString := strconv.Itoa(int(unixTime))
	integrationName := strings.Join([]string{"github", "integration_test", timeString}, "_")
	integrationId, err := testsPipelinesIntegrationsService.CreateGithubIntegration(integrationName, *PipelinesVcsToken)
	assert.NoError(t, err)
	defer deleteIntegrationAndAssert(t, integrationId)

	// Create source with the above integration and assert.
	sourceId, srcErr := testsPipelinesSourcesService.AddSource(integrationId, *PipelinesVcsRepoFullPath, *PipelinesVcsBranch, services.DefaultPipelinesFileFilter)
	assert.NoError(t, srcErr)
	defer deleteSourceAndAssert(t, sourceId)
	syncErr := testPipelinesSyncService.SyncPipelineSource(*PipelinesVcsBranch, *PipelinesVcsRepoFullPath)
	assert.NoError(t, syncErr)
}

func testGetSyncStatus(t *testing.T) {
	if *PipelinesVcsToken == "" {
		assert.NotEmpty(t, *PipelinesVcsToken, "cannot run pipelines tests without vcs token configured")
		return
	}
	// Create integration with provided token.'
	unixTime := time.Now().Unix()
	timeString := strconv.Itoa(int(unixTime))
	integrationName := strings.Join([]string{"github", "integration_test", timeString}, "_")
	integrationId, err := testsPipelinesIntegrationsService.CreateGithubIntegration(integrationName, *PipelinesVcsToken)
	assert.NoError(t, err)
	defer deleteIntegrationAndAssert(t, integrationId)

	// Create source with the above integration and assert.
	sourceId, srcErr := testsPipelinesSourcesService.AddSource(integrationId, *PipelinesVcsRepoFullPath, *PipelinesVcsBranch, services.DefaultPipelinesFileFilter)
	assert.NoError(t, srcErr)

	defer deleteSourceAndAssert(t, sourceId)
	syncErr := testPipelinesSyncService.SyncPipelineSource(*PipelinesVcsBranch, *PipelinesVcsRepoFullPath)
	assert.NoError(t, syncErr)
	time.Sleep(15 * time.Second)
	resourceStatus, syncStatusErr := testPipelinesSyncStatusService.GetSyncPipelineResourceStatus(*PipelinesVcsBranch)
	assert.NoError(t, syncStatusErr)
	if resourceStatus[0].LastSyncStatusCode != 4002 {
		time.Sleep(15 * time.Second)
		_, syncStatusErr := testPipelinesSyncStatusService.GetSyncPipelineResourceStatus(*PipelinesVcsBranch)
		assert.NoError(t, syncStatusErr)
	}
}

func testGetRunStatus(t *testing.T) {
	assert.NotEmpty(t, *PipelinesVcsToken, "cannot run pipelines tests without vcs token configured")
	// Create integration with provided token.
	unixTime := time.Now().Unix()
	timeString := strconv.Itoa(int(unixTime))
	integrationName := strings.Join([]string{"github", "integration_test", timeString}, "_")
	integrationId, err := testsPipelinesIntegrationsService.CreateGithubIntegration(integrationName, *PipelinesVcsToken)
	assert.NoError(t, err)

	defer deleteIntegrationAndAssert(t, integrationId)

	// Create source with the above integration and assert.
	sourceId, sourceErr := testsPipelinesSourcesService.AddSource(integrationId, *PipelinesVcsRepoFullPath, *PipelinesVcsBranch, services.DefaultPipelinesFileFilter)
	assert.NoError(t, sourceErr)
	defer deleteSourceAndAssert(t, sourceId)

	syncErr := testPipelinesSyncService.SyncPipelineSource(*PipelinesVcsBranch, *PipelinesVcsRepoFullPath)
	assert.NoError(t, syncErr)

	time.Sleep(15 * time.Second)
	resourceStatus, syncStatusErr := testPipelinesSyncStatusService.GetSyncPipelineResourceStatus(*PipelinesVcsBranch)
	assert.NoError(t, syncStatusErr)
	if resourceStatus[0].LastSyncStatusCode != 4002 {
		time.Sleep(15 * time.Second)
		_, syncStatusErr := testPipelinesSyncStatusService.GetSyncPipelineResourceStatus(*PipelinesVcsBranch)
		if syncStatusErr != nil {
			assert.NoError(t, syncStatusErr)
			return
		}
	}
	pipelineName := "pipelines_run_int_test"
	status, trigErr := testPipelinesRunService.TriggerPipelineRun(*PipelinesVcsBranch, pipelineName, false)
	assert.NoError(t, trigErr)
	assertTriggerRun(t, pipelineName, *PipelinesVcsBranch, status)
	time.Sleep(60 * time.Second)

	runStatus, runErr := testPipelinesRunService.GetRunStatus(*PipelinesVcsBranch, pipelineName, false)
	assert.NoError(t, runErr)
	if runStatus != nil && len(runStatus.Pipelines) == 0 {
		_, runErr := testPipelinesRunService.GetRunStatus(*PipelinesVcsBranch, pipelineName, true)
		assert.NoError(t, runErr)
	}
	if runStatus != nil && len(runStatus.Pipelines) > 0 && isCancellable(runStatus.Pipelines[0].Run.StatusCode) {

		runStatusCode := runStatus.Pipelines[0].Run.StatusCode
		assertRunStatus(t, runStatusCode)
		runID := runStatus.Pipelines[0].Run.ID
		run, cancelErr := testPipelinesRunService.CancelTheRun(runID)
		assert.NoError(t, cancelErr)
		assert.Equal(t, "cancelled run "+strconv.Itoa(runID)+" successfully", run)
	}

}

func isCancellable(statusCode int) bool {
	switch statusCode {
	case 4000:
		fallthrough
	case 4001:
		fallthrough
	case 4005:
		fallthrough
	case 4006:
		fallthrough
	case 4016:
		fallthrough
	case 4022:
		return true

	}
	return false
}

func assertRunStatus(t *testing.T, statusCode int) {
	assert.True(t, statusCode >= 4000 && statusCode <= 4022)
}

func assertTriggerRun(t *testing.T, pipeline string, branch string, result string) {
	expected := fmt.Sprintf("triggered successfully\n%s %s \n%14s %s", "PipelineName :", pipeline, "Branch :", branch)
	assert.Equal(t, expected, result)
}
