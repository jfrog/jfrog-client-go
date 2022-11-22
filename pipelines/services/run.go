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

const (
	RunStatus = "api/v1/search/pipelines/?pipelineSourceBranch="
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

func (rs *RunService) GetRunStatus(branch, pipeName string) (*PipResponse, error) {
	httpDetails := rs.getHttpDetails()

	/* query params */
	m := make(map[string]string, 0)
	m["pipelineSourceBranch"] = branch
	if pipeName != "" {
		m["name"] = pipeName
	}
	/* query params */

	/* URL Construction */
	uri := rs.constructPipelinesURL(m, runStatus)
	/* URL Construction */

	/* Prepare Request */
	resp, body, _, err := rs.client.SendGet(uri, true, &httpDetails)
	if err != nil {
		return nil, err
	}
	/* Prepare Request */

	/* Response Analysis */
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return nil, err
	}
	p := GetPipelineResponse()
	err = json.Unmarshal(body, &p)
	if err != nil {
		log.Error("Failed to unmarshal json response")
		return &PipResponse{}, err
	}
	/* Response Analysis */
	return &p, nil
}

func (rs *RunService) getHttpDetails() httputils.HttpClientDetails {
	httpDetails := rs.ServiceDetails.CreateHttpClientDetails()
	return httpDetails
}

func GetPipelineResponse() PipResponse {
	r := PipResponse{}
	return r
}

/*
constructPipelinesURL creates URL with all required details to make api call
like headers, queryParams, apiPath
*/
func (rs *RunService) constructPipelinesURL(qParams map[string]string, apiPath string) string {
	uri, err := url.Parse(rs.ServiceDetails.GetUrl() + apiPath)
	if err != nil {
		log.Error("Failed to parse pipelines fetch run status url")
	}
	queryString := uri.Query()
	for k, v := range qParams {
		queryString.Set(k, v)
	}
	uri.RawQuery = queryString.Encode()

	return uri.String()
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
		return "", err
	}

	/* URL Construction */
	headers := make(map[string]string, 0)
	headers["Content-Type"] = "application/json"
	httpDetails.Headers = headers
	uri := rs.constructPipelinesURL(m, triggerpipeline)
	/* URL Construction */

	/* Prepare Request */
	resp, body, err := rs.client.SendPost(uri, buf.Bytes(), &httpDetails)
	if err != nil {
		return "", err
	}
	/* Prepare Request */

	/* Response Analysis */
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return "", err
	}
	if resp.StatusCode == http.StatusOK {
		s := fmt.Sprintf("triggered successfully\n%s %s \n%14s %s", "PipelineName :", pipeline, "Branch :", branch)
		return s, nil
	}
	/* Response Analysis */

	return "", nil
}
