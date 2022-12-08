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
	"net/url"
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
	syncPipelineResource = "api/v1/pipelineSources"
	resourceVersions     = "api/v1/resourceVersions"
)

func (rs *RunService) GetRunStatus(branch, pipeName string) (*PipelineRunStatusResponse, error) {
	httpDetails := rs.getHttpDetails()

	// query params
	queryParams := make(map[string]string, 0)
	queryParams["pipelineSourceBranch"] = branch
	if pipeName != "" {
		queryParams["name"] = pipeName
	}

	// URL Construction
	uri, pipeURLErr := rs.constructPipelinesURL(queryParams, runStatus)
	if pipeURLErr != nil {
		return nil, errorutils.CheckError(pipeURLErr)
	}

	// Prepare Request
	resp, body, _, err := rs.client.SendGet(uri, true, &httpDetails)
	if err != nil {
		return nil, errorutils.CheckError(err)
	}

	// Response Analysis
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return nil, errorutils.CheckError(err)
	}
	p := PipelineRunStatusResponse{}
	err = json.Unmarshal(body, &p)
	if err != nil {
		log.Error("Failed to unmarshal json response")
		return &PipelineRunStatusResponse{}, errorutils.CheckError(err)
	}
	return &p, nil
}

func (rs *RunService) getHttpDetails() httputils.HttpClientDetails {
	httpDetails := rs.ServiceDetails.CreateHttpClientDetails()
	return httpDetails
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

func (rs *RunService) TriggerPipelineRun(branch, pipeline string) (string, error) {
	httpDetails := rs.getHttpDetails()
	m := make(map[string]string, 0)

	payload := strings.NewReader(`{
	    "branchName": "` + strings.TrimSpace(branch) + `",
		"pipelineName": "` + strings.TrimSpace(pipeline) + `"
		}`)
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(payload)
	if err != nil {
		log.Error("Failed to read stream to send payload to trigger pipelines")
		return "", errorutils.CheckError(err)
	}

	// URL Construction
	headers := make(map[string]string, 0)
	headers["Content-Type"] = "application/json"
	httpDetails.Headers = headers
	uri, pipeURLErr := rs.constructPipelinesURL(m, triggerpipeline)
	if pipeURLErr != nil {
		return "", pipeURLErr
	}

	// Prepare Request
	resp, body, err := rs.client.SendPost(uri, buf.Bytes(), &httpDetails)
	if err != nil {
		return "", errorutils.CheckError(err)
	}

	// Response Analysis
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return "", err
	}
	if resp.StatusCode == http.StatusOK {
		s := fmt.Sprintf("triggered successfully\n%s %s \n%14s %s", "PipelineName :", pipeline, "Branch :", branch)
		return s, nil
	}

	return "", nil
}
