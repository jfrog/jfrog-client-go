package services

import (
	"encoding/json"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	workspaces         = "api/v1/workspaces"
	deleteWorkspace    = "api/v1/workspaces/:id"
	validateWorkspace  = "api/v1/validateWorkspace"
	workspaceSync      = "api/v1/syncWorkspace"
	workspacePipelines = "api/v1/pipelines"
	workspaceRuns      = "api/v1/runs"
	workspaceSteps     = "api/v1/steps"
	stepConsoles       = "api/v1/steplets/:stepID/consoles"
)

type WorkspaceService struct {
	client *jfroghttpclient.JfrogHttpClient
	auth.ServiceDetails
}

func NewWorkspaceService(client *jfroghttpclient.JfrogHttpClient) *WorkspaceService {
	return &WorkspaceService{client: client}
}

func (ws *WorkspaceService) getHttpDetails() httputils.HttpClientDetails {
	return ws.ServiceDetails.CreateHttpClientDetails()
}

func (ws *WorkspaceService) GetWorkspace() ([]WorkspacesResponse, error) {
	httpDetails := ws.getHttpDetails()

	// Query params
	queryParams := make(map[string]string, 0)

	// URL construction
	uri, err := constructPipelinesURL(queryParams, ws.ServiceDetails.GetUrl(), workspaces)
	if err != nil {
		return nil, err
	}

	// Prepare request
	resp, body, _, err := ws.client.SendGet(uri, true, &httpDetails)
	if err != nil {
		return nil, err
	}

	// Response Analysis
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return nil, err
	}
	wsStatusResp := make([]WorkspacesResponse, 0)
	err = json.Unmarshal(body, &wsStatusResp)
	return wsStatusResp, nil
}

func (ws *WorkspaceService) DeleteWorkspace(wsID string) error {
	httpDetails := ws.getHttpDetails()
	deleteWorkspaceAPI := strings.Replace(deleteWorkspace, ":id", wsID, 1)

	// Query params
	queryParams := make(map[string]string, 0)

	// URL construction
	uri, err := constructPipelinesURL(queryParams, ws.ServiceDetails.GetUrl(), deleteWorkspaceAPI)
	if err != nil {
		return err
	}

	// Prepare request
	resp, body, err := ws.client.SendDelete(uri, nil, &httpDetails)
	if err != nil {
		return err
	}

	// Response analysis
	err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK)
	return err
}

func (ws *WorkspaceService) ValidateWorkspace(data []byte) error {
	httpDetails := ws.getHttpDetails()

	// Query params
	queryParams := make(map[string]string, 0)

	// URL construction
	uri, err := constructPipelinesURL(queryParams, ws.ServiceDetails.GetUrl(), validateWorkspace)
	if err != nil {
		return err
	}

	// Headers
	headers := make(map[string]string, 0)
	headers["Content-Type"] = "application/json"
	httpDetails.Headers = headers

	// Prepare request
	resp, body, err := ws.client.SendPost(uri, data, &httpDetails)
	if err != nil {
		return err
	}

	// Response analysis
	return errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK)
}

func (ws *WorkspaceService) WorkspaceSync() error {
	httpDetails := ws.getHttpDetails()

	// Query params
	queryParams := make(map[string]string, 0)

	// URL construction
	uri, err := constructPipelinesURL(queryParams, ws.ServiceDetails.GetUrl(), workspaceSync)
	if err != nil {
		return err
	}

	// Prepare request
	resp, body, _, err := ws.client.SendGet(uri, true, &httpDetails)
	if err != nil {
		return err
	}

	// Response analysis
	return errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK)
}

func (ws *WorkspaceService) WorkspaceRunIDs(pipelines []string) ([]PipelinesRunID, error) {
	httpDetails := ws.getHttpDetails()
	pipelineFilter := strings.Join(pipelines, ",")

	// Query params
	// TODO ADD include in query param if needed
	queryParams := map[string]string{
		"names":   pipelineFilter,
		"limit":   "1",
		"include": "latestRunId,name",
	}

	// URL construction
	uri, err := constructPipelinesURL(queryParams, ws.ServiceDetails.GetUrl(), workspacePipelines)
	if err != nil {
		return nil, err
	}

	// Prepare request
	resp, body, _, err := ws.client.SendGet(uri, true, &httpDetails)
	if err != nil {
		return nil, err
	}
	// Response analysis
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return nil, err
	}
	pipeRunIDs := make([]PipelinesRunID, 0)
	err = json.Unmarshal(body, &pipeRunIDs)
	return pipeRunIDs, err
}

func (ws *WorkspaceService) WorkspaceRunStatus(pipelinesRunID int) ([]byte, error) {
	httpDetails := ws.getHttpDetails()

	// Query params
	// TODO ADD include in query param if needed
	queryParams := map[string]string{}

	// URL construction
	uri, err := constructPipelinesURL(queryParams, ws.ServiceDetails.GetUrl(), workspaceRuns+"/"+strconv.Itoa(pipelinesRunID))
	if err != nil {
		return nil, err
	}

	// Prepare request
	resp, body, _, err := ws.client.SendGet(uri, true, &httpDetails)
	if err != nil {
		return nil, err
	}
	// Response analysis
	err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK)
	return body, err
}

func (ws *WorkspaceService) WorkspaceStepStatus(pipelinesRunID int) ([]byte, error) {
	httpDetails := ws.getHttpDetails()

	// Query params
	queryParams := map[string]string{
		"runIds": strconv.Itoa(pipelinesRunID),
		"limit":  "15",
		//"include": "name,statusCode,triggeredAt,externalBuildUrl",
	}

	// URL construction
	uri, err := constructPipelinesURL(queryParams, ws.ServiceDetails.GetUrl(), workspaceSteps)
	if err != nil {
		return nil, err
	}

	// Prepare request
	resp, body, _, err := ws.client.SendGet(uri, true, &httpDetails)
	if err != nil {
		return nil, err
	}

	// Response analysis
	err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK)
	return body, err
}

func (ws *WorkspaceService) GetWorkspacePipelines(workspaces []WorkspacesResponse) (map[string]string, error) {
	pipelineNames := make(map[string]string, 1)

	log.Info("Collecting pipeline names configured")
	// Validate and return pipeline names as slice
	if len(workspaces) > 0 && !(*workspaces[0].IsSyncing) {
		pipelines := workspaces[0].PipelinesYmlPropertyBag.Pipelines
		for _, pi := range pipelines {
			pipelineNames[pi.Name] = pi.PipelineSourceBranch
		}
	}
	return pipelineNames, nil
}

func (ws *WorkspaceService) WorkspacePollSyncStatus() ([]WorkspacesResponse, error) {
	httpDetails := ws.getHttpDetails()

	// Query params
	queryParams := make(map[string]string, 0)

	// URL construction
	uri, err := constructPipelinesURL(queryParams, ws.ServiceDetails.GetUrl(), workspaces)
	if err != nil {
		return nil, err
	}

	pollingAction := func() (shouldStop bool, responseBody []byte, err error) {
		log.Info("Polling for pipeline resource sync")
		// Prepare request
		resp, body, _, err := ws.client.SendGet(uri, true, &httpDetails)
		if err != nil {
			return false, body, err
		}

		// Response Analysis
		if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
			return false, body, err
		}
		log.Debug("response received ", string(body))
		wsStatusResp := make([]WorkspacesResponse, 0)
		err = json.Unmarshal(body, &wsStatusResp)
		if err != nil {
			log.Error("failed to unmarshal validation response")
			return true, body, err
		}
		if len(wsStatusResp) > 0 && *wsStatusResp[0].IsSyncing {
			return false, body, err
		} else if len(wsStatusResp) > 0 && !(*wsStatusResp[0].IsSyncing) {
			log.Info("Sync is completed")
			return true, body, err
		}
		return true, body, err
	}
	pollingExecutor := &httputils.PollingExecutor{
		Timeout:         10 * time.Minute,
		PollingInterval: 5 * time.Second,
		PollingAction:   pollingAction,
		MsgPrefix:       "Get pipeline workspace sync status...",
	}
	// Polling execution
	body, err := pollingExecutor.Execute()
	if err != nil {
		return nil, err
	}
	workspaceStatusResponse := make([]WorkspacesResponse, 0)
	err = json.Unmarshal(body, &workspaceStatusResponse)
	return workspaceStatusResponse, err
}

// GetStepLogsUsingStepID retrieve steps logs using step id
func (ws *WorkspaceService) GetStepLogsUsingStepID(stepID string) (map[string][]Console, error) {
	httpDetails := ws.getHttpDetails()
	stepConsolesAPI := strings.Replace(stepConsoles, ":stepID", stepID, 1)

	// Query params
	queryParams := make(map[string]string, 0)

	// URL construction
	uri, err := constructPipelinesURL(queryParams, ws.ServiceDetails.GetUrl(), stepConsolesAPI)
	if err != nil {
		return nil, err
	}

	// Prepare request
	resp, body, _, err := ws.client.SendGet(uri, true, &httpDetails)
	if err != nil {
		return nil, err
	}

	// Response Analysis
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return nil, err
	}
	consoles := make(map[string][]Console)
	err = json.Unmarshal(body, &consoles)
	return consoles, err
}
