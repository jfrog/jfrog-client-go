package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
	"net/url"
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
	runStatus            = "api/v1/search/pipelines/"
	triggerpipeline      = "api/v1/pipelines/trigger"
	pipelineSyncStatus   = "api/v1/pipelineSyncStatuses"
	pipelineResources    = "api/v1/pipelineSources"
	cancelRunPath        = "api/v1/runs/:runId/cancel"
	syncPipelineResource = "api/v1/pipelineSources"
	resourceVersions     = "api/v1/resourceVersions"
)

func (rs *RunService) GetRunStatus(branch, pipeName string, isMultiBranch bool) (*PipelineRunStatusResponse, error) {
	httpDetails := rs.getHttpDetails()

	// query params
	queryParams := make(map[string]string, 0)
	if isMultiBranch { // add this query param only when pipeline source is multi-branch
		queryParams["pipelineSourceBranch"] = branch
	}
	if pipeName != "" {
		queryParams["name"] = pipeName
	}

	// URL Construction
	uri, pipeURLErr := rs.constructPipelinesURL(queryParams, runStatus)
	if pipeURLErr != nil {
		return nil, pipeURLErr
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
	fmt.Printf("pipeline response %+v\n", string(body))
	return &pipRunResp, errorutils.CheckError(err)
}

func (rs *RunService) getHttpDetails() httputils.HttpClientDetails {
	return rs.ServiceDetails.CreateHttpClientDetails()
}

// constructPipelinesURL creates URL with all required details to make api call
// like headers, queryParams, apiPath
func (rs *RunService) constructPipelinesURL(qParams map[string]string, apiPath string) (string, error) {
	uri, err := url.Parse(rs.ServiceDetails.GetUrl() + apiPath)
	if err != nil {
		log.Error("Failed to parse pipelines fetch run status url")
		return "", errorutils.CheckError(err)
	}
	queryString := uri.Query()
	for k, v := range qParams {
		queryString.Set(k, v)
	}
	uri.RawQuery = queryString.Encode()

	return uri.String(), nil
}

func (rs *RunService) TriggerPipelineRun(branch, pipeline string, isMultiBranch bool) error {
	httpDetails := rs.getHttpDetails()
	queryParams := make(map[string]string, 0)

	var payload *strings.Reader
	if isMultiBranch { // add this query param only when pipeline source is multi-branch
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
		log.Error("Failed to read stream to send payload to trigger pipelines")
		return errorutils.CheckError(err)
	}

	// URL Construction
	utils.AddHeader("Content-Type", "application/json", &httpDetails.Headers)
	uri, pipeURLErr := rs.constructPipelinesURL(queryParams, triggerpipeline)
	if pipeURLErr != nil {
		return pipeURLErr
	}

	// Prepare Request
	resp, body, err := rs.client.SendPost(uri, buf.Bytes(), &httpDetails)
	if err != nil {
		return err
	}

	// Response Analysis
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return err
	}
	log.Info(fmt.Sprintf("triggered successfully\n%s %s \n%14s %s", "PipelineName :", pipeline, "Branch :", branch))

	return nil
}

func (rs *RunService) CancelRun(runID int) error {
	log.Info("cancelling the run ", runID)
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
		log.Error("Failed to read stream to send payload to trigger pipelines")
		return errorutils.CheckError(err)
	}

	// URL Construction
	utils.AddHeader("Content-Type", "application/json", &httpDetails.Headers)
	uri, pipeURLErr := rs.constructPipelinesURL(queryParams, cancelRun)
	if pipeURLErr != nil {
		return pipeURLErr
	}

	// Prepare Request
	resp, body, err := rs.client.SendPost(uri, buf.Bytes(), &httpDetails)
	if err != nil {
		return errorutils.CheckError(err)
	}

	// Response Analysis
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return err
	}
	if resp.StatusCode == http.StatusOK {
		log.Info(fmt.Sprintf("cancelled run %s successfully", runValue))
		return nil
	}
	log.Error("unable to find run id")
	return errors.New(fmt.Sprintf("Unable to find run ID: %d", runID))
}
