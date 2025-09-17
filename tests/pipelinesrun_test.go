//go:build itest

package tests

import (
	"fmt"
	pipelinesServices "github.com/jfrog/jfrog-client-go/pipelines/services"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/stretchr/testify/assert"
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

const (
	// Define default wait time
	defaultMaxWaitMinutes     = 10 * time.Minute
	defaultSyncSleepInterval  = 5 * time.Second
	defaultPipelineSourceName = "PipeRunTest"
)

func testTriggerSync(t *testing.T) {
	if !assert.NotEmpty(t, *PipelinesVcsToken, "cannot run pipelines tests without vcs token configured") {
		return
	}
	deleteSourceIfAlreadyExists(t)
	// Create integration with provided token.
	integrationName := getUniqueIntegrationName(pipelinesServices.GithubName)
	integrationId, err := testsPipelinesIntegrationsService.CreateGithubIntegration(integrationName, *PipelinesVcsToken)
	assert.NoError(t, err)
	defer deleteIntegrationAndAssert(t, integrationId)

	// Create source with the above integration and assert.
	sourceId, err := testsPipelinesSourcesService.AddSource(integrationId, *PipelinesVcsRepoFullPath, *PipelinesVcsBranch, pipelinesServices.DefaultPipelinesFileFilter, defaultPipelineSourceName)
	if !assert.NoError(t, err) {
		return
	}
	defer deleteSourceAndAssert(t, sourceId)
	pollSyncPipelineSource(t)
}

func testGetSyncStatus(t *testing.T) {
	if !assert.NotEmpty(t, *PipelinesVcsToken, "cannot run pipelines tests without vcs token configured") {
		return
	}
	deleteSourceIfAlreadyExists(t)
	// Create integration with provided token.
	integrationName := getUniqueIntegrationName(pipelinesServices.GithubName)
	integrationId, err := testsPipelinesIntegrationsService.CreateGithubIntegration(integrationName, *PipelinesVcsToken)
	assert.NoError(t, err)
	defer deleteIntegrationAndAssert(t, integrationId)

	// Create source with the above integration and assert.
	sourceId, err := testsPipelinesSourcesService.AddSource(integrationId, *PipelinesVcsRepoFullPath, *PipelinesVcsBranch, pipelinesServices.DefaultPipelinesFileFilter, defaultPipelineSourceName)
	if !assert.NoError(t, err) {
		return
	}
	defer deleteSourceAndAssert(t, sourceId)
	pollSyncPipelineSource(t)
	pollForSyncResourceStatus(t)
}

func testGetRunStatus(t *testing.T) {
	if !assert.NotEmpty(t, *PipelinesVcsToken, "cannot run pipelines tests without vcs token configured") {
		return
	}
	deleteSourceIfAlreadyExists(t)
	// Create integration with provided token.
	integrationName := getUniqueIntegrationName(pipelinesServices.GithubName)
	integrationId, err := testsPipelinesIntegrationsService.CreateGithubIntegration(integrationName, *PipelinesVcsToken)
	assert.NoError(t, err)
	defer deleteIntegrationAndAssert(t, integrationId)

	// Create source with the above integration and assert.
	sourceId, err := testsPipelinesSourcesService.AddSource(integrationId, *PipelinesVcsRepoFullPath, *PipelinesVcsBranch, pipelinesServices.DefaultPipelinesFileFilter, defaultPipelineSourceName)
	if !assert.NoError(t, err) {
		return
	}
	defer deleteSourceAndAssert(t, sourceId)

	pollSyncPipelineSource(t)

	pollForSyncResourceStatus(t)
	res, resourceErr := pipelinesServices.GetPipelineResource(testPipelinesSyncService.GetHTTPClient(),
		testPipelinesSyncService.GetServiceURL(),
		*PipelinesVcsRepoFullPath,
		testPipelinesSyncService.GetHttpDetails())

	assert.NoError(t, resourceErr)
	pipelineName := "pipelines_run_int_test"
	trigErr := testPipelinesRunService.TriggerPipelineRun(*PipelinesVcsBranch, pipelineName, *res.IsMultiBranch)
	assert.NoError(t, trigErr)

	pollGetRunStatus(t, pipelineName)
}

func deleteSourceIfAlreadyExists(t *testing.T) {
	queryParams := map[string]string{
		"name":               defaultPipelineSourceName,
		"repositoryFullName": strings.TrimSpace(*PipelinesVcsRepoFullPath),
		"branch":             strings.TrimSpace(*PipelinesVcsBranch),
	}
	sources, err := testsPipelinesSourcesService.GetSourceByFilter(queryParams)
	if err != nil {
		return
	}
	for _, source := range sources {
		deleteSourceAndAssert(t, source.Id)
	}
}

func pollGetRunStatus(t *testing.T, pipelineName string) {
	pollingAction := func() (shouldStop bool, responseBody []byte, err error) {
		res, resourceErr := pipelinesServices.GetPipelineResource(testPipelinesSyncService.GetHTTPClient(),
			testPipelinesSyncService.ServiceDetails.GetUrl(),
			*PipelinesVcsRepoFullPath,
			testPipelinesSyncService.ServiceDetails.CreateHttpClientDetails())

		assert.NoError(t, resourceErr)
		pipRunResponse, syncErr := testPipelinesRunService.GetRunStatus(*PipelinesVcsBranch, pipelineName, *res.IsMultiBranch)
		assert.NoError(t, syncErr)

		// Got the full valid response.
		if pipRunResponse != nil && len(pipRunResponse.Pipelines) > 0 && pipRunResponse.Pipelines[0].Name == pipelineName {
			log.Info("pipelines status code", pipRunResponse.Pipelines[0].Run.StatusCode)
			if isCancellable(pipRunResponse.Pipelines[0].Run.StatusCode) {
				runStatusCode := pipRunResponse.Pipelines[0].Run.StatusCode
				assertRunStatus(t, runStatusCode)
				runID := pipRunResponse.Pipelines[0].Run.ID
				cancelErr := testPipelinesRunService.CancelRun(runID)
				assert.NoError(t, cancelErr)
			} else {
				return false, []byte{}, nil
			}
			return true, []byte{}, nil
		}
		return false, []byte{}, nil
	}
	pollingExecutor := &httputils.PollingExecutor{
		Timeout:         defaultMaxWaitMinutes,
		PollingInterval: defaultSyncSleepInterval,
		PollingAction:   pollingAction,
		MsgPrefix:       "Get pipeline run status...",
	}
	// Polling execution
	_, err := pollingExecutor.Execute()
	assert.NoError(t, err)
}

func pollForSyncResourceStatus(t *testing.T) {
	// Define polling action
	pollingAction := func() (shouldStop bool, responseBody []byte, err error) {
		pipResStatus, syncErr := testPipelinesSyncStatusService.GetSyncPipelineResourceStatus(*PipelinesVcsRepoFullPath, *PipelinesVcsBranch)
		assert.NoError(t, syncErr)

		if len(pipResStatus) > 0 && pipResStatus[0].LastSyncStatusCode == 4002 {
			// Got the full valid response.
			return true, []byte{}, nil
		} else if len(pipResStatus) > 0 && pipResStatus[0].LastSyncStatusCode == 4003 {
			// Sync Failed, it is not needed to retry.
			return true, []byte{}, fmt.Errorf("failed to sync pipelines source\n%s", pipResStatus[0].LastSyncLogs)
		}
		return false, []byte{}, nil
	}
	pollingExecutor := &httputils.PollingExecutor{
		Timeout:         defaultMaxWaitMinutes,
		PollingInterval: defaultSyncSleepInterval,
		PollingAction:   pollingAction,
		MsgPrefix:       "Get pipeline sync status...",
	}
	// Polling execution
	_, err := pollingExecutor.Execute()
	assert.NoError(t, err)
}

func isCancellable(statusCode int) bool {
	switch statusCode {
	case 4000, 4001, 4005, 4006, 4016, 4022:
		return true
	}
	return false
}

func assertRunStatus(t *testing.T, statusCode int) {
	assert.GreaterOrEqual(t, statusCode, 4000)
	assert.LessOrEqual(t, statusCode, 4022)
}

func pollSyncPipelineSource(t *testing.T) {
	// Define polling action
	pollingAction := func() (shouldStop bool, responseBody []byte, err error) {
		syncErr := testPipelinesSyncService.SyncPipelineSource(*PipelinesVcsBranch, *PipelinesVcsRepoFullPath)
		assert.NoError(t, syncErr)

		return syncErr == nil, nil, syncErr
	}

	pollingExecutor := &httputils.PollingExecutor{
		Timeout:         defaultMaxWaitMinutes,
		PollingInterval: defaultSyncSleepInterval,
		PollingAction:   pollingAction,
		MsgPrefix:       "Syncing Pipeline Resource...",
	}
	// Polling execution
	_, err := pollingExecutor.Execute()
	assert.NoError(t, err)
}
