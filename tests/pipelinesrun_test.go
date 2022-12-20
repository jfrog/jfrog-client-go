package tests

import (
	"fmt"
	"github.com/jfrog/jfrog-client-go/pipelines/services"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
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
	pollSyncPipelineSource(t)
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
	pollSyncPipelineSource(t)
	pollForSyncResourceStatus(t)
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

	pollSyncPipelineSource(t)

	pollForSyncResourceStatus(t)
	pipelineName := "pipelines_run_int_test"
	status, trigErr := testPipelinesRunService.TriggerPipelineRun(*PipelinesVcsBranch, pipelineName, false)
	assert.NoError(t, trigErr)
	assertTriggerRun(t, pipelineName, *PipelinesVcsBranch, status)

	pollGetRunStatus(t, pipelineName)
}

func pollGetRunStatus(t *testing.T, pipelineName string) {
	pollingAction := func() (shouldStop bool, responseBody []byte, err error) {
		pipRunResponse, syncErr := testPipelinesRunService.GetRunStatus(*PipelinesVcsBranch, pipelineName, false)
		assert.NoError(t, syncErr)

		// Got the full valid response.
		if pipRunResponse != nil && len(pipRunResponse.Pipelines) > 0 && pipRunResponse.Pipelines[0].Name == pipelineName {
			log.Info("pipelines status code ", pipRunResponse.Pipelines[0].Run.StatusCode)
			if isCancellable(pipRunResponse.Pipelines[0].Run.StatusCode) {

				runStatusCode := pipRunResponse.Pipelines[0].Run.StatusCode
				assertRunStatus(t, runStatusCode)
				runID := pipRunResponse.Pipelines[0].Run.ID
				run, cancelErr := testPipelinesRunService.CancelTheRun(runID)
				assert.NoError(t, cancelErr)
				assert.Equal(t, "cancelled run "+strconv.Itoa(runID)+" successfully", run)
			}
			return true, []byte{}, nil
		}
		return false, []byte{}, nil
	}
	// define default wait time
	defaultMaxWaitMinutes := 10 * time.Minute
	defaultSyncSleepInterval := 5 * time.Second // 5 seconds
	pollingExecutor := &httputils.PollingExecutor{
		Timeout:         defaultMaxWaitMinutes,
		PollingInterval: defaultSyncSleepInterval,
		PollingAction:   pollingAction,
		MsgPrefix:       "Get pipeline run status...",
	}
	// polling execution
	_, err := pollingExecutor.Execute()
	assert.NoError(t, err)
}

func pollForSyncResourceStatus(t *testing.T) {
	//define polling action
	pollingAction := func() (shouldStop bool, responseBody []byte, err error) {
		pipResStatus, body, syncErr := testPipelinesSyncStatusService.GetSyncPipelineResourceStatus(*PipelinesVcsRepoFullPath, *PipelinesVcsBranch)
		assert.NoError(t, syncErr)

		// Got the full valid response.
		if len(pipResStatus) > 0 && pipResStatus[0].LastSyncStatusCode == 4002 {
			return true, body, nil
		}
		return false, body, nil
	}
	// define default wait time
	defaultMaxWaitMinutes := 10 * time.Minute
	defaultSyncSleepInterval := 5 * time.Second // 5 seconds
	pollingExecutor := &httputils.PollingExecutor{
		Timeout:         defaultMaxWaitMinutes,
		PollingInterval: defaultSyncSleepInterval,
		PollingAction:   pollingAction,
		MsgPrefix:       "Get pipeline sync status...",
	}
	// polling execution
	_, err := pollingExecutor.Execute()
	assert.NoError(t, err)
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

func pollSyncPipelineSource(t *testing.T) {
	//define polling action
	pollingAction := func() (shouldStop bool, responseBody []byte, err error) {
		statusCode, body, syncErr := testPipelinesSyncService.SyncPipelineSource(*PipelinesVcsBranch, *PipelinesVcsRepoFullPath)
		assert.NoError(t, syncErr)

		// Got the full valid response.
		if statusCode == http.StatusOK {
			return true, body, nil
		}
		return false, body, nil
	}
	// define default wait time
	defaultMaxWaitMinutes := 10 * time.Minute
	defaultSyncSleepInterval := 5 * time.Second // 5 seconds
	pollingExecutor := &httputils.PollingExecutor{
		Timeout:         defaultMaxWaitMinutes,
		PollingInterval: defaultSyncSleepInterval,
		PollingAction:   pollingAction,
		MsgPrefix:       "Syncing Pipeline Resource...",
	}
	// polling execution
	_, err := pollingExecutor.Execute()
	assert.NoError(t, err)
}
