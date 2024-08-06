package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
	"strconv"
	"strings"
)

type RunService struct {
	client *jfroghttpclient.JfrogHttpClient
	auth.ServiceDetails
}

func NewRunService(client *jfroghttpclient.JfrogHttpClient) *RunService {
	return &RunService{client: client}
}

const (
	runStatus          = "api/v1/search/pipelines/"
	triggerpipeline    = "api/v1/pipelines/trigger"
	pipelineSyncStatus = "api/v1/pipelineSyncStatuses"
	pipelineResources  = "api/v1/pipelineSources"
	cancelRunPath      = "api/v1/runs/:runId/cancel"
)

func (rs *RunService) GetRunStatus(branch, pipeName string, isMultiBranch bool) (*PipelineRunStatusResponse, error) {
	httpDetails := rs.getHttpDetails()

	// Query params
	queryParams := make(map[string]string)
	if isMultiBranch {
		// Add this query param only when pipeline source is multi-branch
		queryParams["pipelineSourceBranch"] = branch
	}
	if pipeName != "" {
		queryParams["name"] = pipeName
	}

	// URL Construction
	uri, err := constructPipelinesURL(queryParams, rs.ServiceDetails.GetUrl(), runStatus)
	if err != nil {
		return nil, err
	}

	// Prepare Request
	resp, body, _, err := rs.client.SendGet(uri, true, &httpDetails)
	if err != nil {
		return nil, err
	}

	// Response Analysis
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return nil, err
	}
	pipRunResp := PipelineRunStatusResponse{}
	err = json.Unmarshal(body, &pipRunResp)
	return &pipRunResp, errorutils.CheckError(err)
}

func (rs *RunService) getHttpDetails() httputils.HttpClientDetails {
	return rs.ServiceDetails.CreateHttpClientDetails()
}

func (rs *RunService) TriggerPipelineRun(branch, pipeline string, isMultiBranch bool) error {
	httpDetails := rs.getHttpDetails()
	queryParams := make(map[string]string, 0)

	var payload *strings.Reader
	if isMultiBranch {
		// Add this query param only when pipeline source is multi-branch
		payload = strings.NewReader(`{
	    "branchName": "` + strings.TrimSpace(branch) + `",
		"pipelineName": "` + strings.TrimSpace(pipeline) + `"
		}`)
	} else {
		payload = strings.NewReader(`{
		"pipelineName": "` + strings.TrimSpace(pipeline) + `"
		}`)
	}
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(payload)
	if err != nil {
		return errorutils.CheckError(err)
	}

	// URL Construction
	httpDetails.SetContentTypeApplicationJson()
	uri, err := constructPipelinesURL(queryParams, rs.ServiceDetails.GetUrl(), triggerpipeline)
	if err != nil {
		return err
	}

	// Prepare Request
	resp, body, err := rs.client.SendPost(uri, buf.Bytes(), &httpDetails)
	if err != nil {
		return err
	}

	// Response Analysis
	if err := errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return err
	}
	log.Info(fmt.Sprintf("Triggered successfully\n%s %s \n%14s %s", "PipelineName :", pipeline, "Branch :", branch))

	return nil
}

func (rs *RunService) CancelRun(runID int) error {
	log.Info("Cancelling the run", runID)
	runValue := strconv.Itoa(runID)
	cancelRun := cancelRunPath
	cancelRun = strings.Replace(cancelRun, ":runId", runValue, 1)
	httpDetails := rs.getHttpDetails()
	queryParams := make(map[string]string, 0)
	payload := strings.NewReader(`{
		}`)
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(payload)
	if err != nil {
		return errorutils.CheckError(err)
	}

	// URL Construction
	httpDetails.SetContentTypeApplicationJson()
	uri, err := constructPipelinesURL(queryParams, rs.ServiceDetails.GetUrl(), cancelRun)
	if err != nil {
		return err
	}

	// Prepare Request
	resp, body, err := rs.client.SendPost(uri, buf.Bytes(), &httpDetails)
	if err != nil {
		return errorutils.CheckError(err)
	}

	// Response Analysis
	if err := errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return err
	}
	if resp.StatusCode == http.StatusOK {
		log.Info(fmt.Sprintf("Cancelled run %s successfully", runValue))
		return nil
	}
	return fmt.Errorf("unable to find run ID: %d", runID)
}
