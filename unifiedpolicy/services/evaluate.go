package services

import (
	"encoding/json"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"net/http"
)

const (
	evaluateEndpoint = "unifiedpolicy/api/v1/evaluate"
)

type EvaluateService struct {
	client         *jfroghttpclient.JfrogHttpClient
	serviceDetails auth.ServiceDetails
	config         config.Config
}

func NewEvaluateService(client *jfroghttpclient.JfrogHttpClient, serviceDetails auth.ServiceDetails) *EvaluateService {
	return &EvaluateService{client: client, serviceDetails: serviceDetails}
}

func (c *EvaluateService) sendPostRequest(requestContent []byte) (resp *http.Response, body []byte, err error) {
	commitInfoUrl := c.serviceDetails.GetUrl() + evaluateEndpoint
	clientDetails := c.serviceDetails.CreateHttpClientDetails()
	resp, body, err = c.client.SendPost(commitInfoUrl, requestContent, &clientDetails)
	return
}

func (c *EvaluateService) Evaluate(evaluateRequest *EvaluateRequest) (*EvaluateResponse, error) {
	requestContent, err := json.Marshal(evaluateRequest)
	if err != nil {
		return nil, errorutils.CheckError(err)
	}
	resp, body, err := c.sendPostRequest(requestContent)
	if err != nil {
		return nil, err
	}
	if err = errorutils.CheckResponseStatus(resp, http.StatusOK); err != nil {
		return nil, errorutils.CheckError(errorutils.GenerateResponseError(resp.Status, utils.IndentJson(body)))
	}
	var response *EvaluateResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, errorutils.CheckErrorf("failed to unmarshal evaluate response: " + err.Error())
	}
	return response, nil
}

type EvaluateResponse struct {
	Decision     string `json:"stage"`
	Explanations string `json:"explanation"`
}

type Context struct {
	Stage string `json:"stage"`
}

type Resource struct {
	ApplicationKey string `json:"application_key"`
	Type           string `json:"type"`
	MultiScanId    string `json:"multi_scan_id"`
	GitRepoUrl     string `json:"git_repo_url"`
	PullRequestId  string `json:"pr_id"`
}

type EvaluateRequest struct {
	Action   string   `json:"action"`
	Context  Context  `json:"context"`
	Resource Resource `json:"resource"`
}
