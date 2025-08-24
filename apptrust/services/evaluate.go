package services

import (
	"encoding/json"
	"fmt"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"net/http"
)

const (
	apptrustDomain   = "unifiedpolicy/"
	evaluateEndpoint = "evaluate"
	evaluateApi      = apptrustDomain + "api/v1/" + evaluateEndpoint
)

type EvaluationService struct {
	client          *jfroghttpclient.JfrogHttpClient
	ApptrustDetails auth.ServiceDetails
}

func NewEvaluationService(client *jfroghttpclient.JfrogHttpClient) *EvaluationService {
	return &EvaluationService{client: client}
}

type EvaluateRequest struct {
	Action   string           `json:"action"`
	Context  EvaluateContext  `json:"context"`
	Resource EvaluateResource `json:"resource"`
}

type EvaluateContext struct {
	Stage string `json:"stage"`
}

type EvaluateResource struct {
	ApplicationKey string `json:"application_key"`
	Type           string `json:"type"`
	MultiScanId    string `json:"multi_scan_id"`
	GitRepoUrl     string `json:"git_repo_url"`
}

type EvaluateResponse struct {
	Id          string   `json:"id"`
	Decision    string   `json:"decision"`
	Explanation string   `json:"explanation"`
	MissingData []string `json:"missing_data"`
}

func (es *EvaluationService) Evaluate(params EvaluateRequest) (EvaluateResponse, error) {
	// TODO eran delete this section - this is a temporary solution that returns the expected response for testing and development until API is implemented
	return EvaluateResponse{
		Id:          "mock-evaluate-id",
		Decision:    "Error",
		Explanation: "",
		MissingData: []string{"sast", "secrets"},
	}, nil
	// TODO eran delete up to here

	httpDetail := es.ApptrustDetails.CreateHttpClientDetails()
	url := utils.AddTrailingSlashIfNeeded(es.ApptrustDetails.GetUrl()) + evaluateApi
	requestBody, err := json.Marshal(params)
	if err != nil {
		return EvaluateResponse{}, fmt.Errorf("failed to marshal request body: %w", err)
	}

	resp, body, err := es.client.SendPost(url, requestBody, &httpDetail)
	if err != nil {
		return EvaluateResponse{}, fmt.Errorf("failed to send POST request to '%s': %w", url, err)
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return EvaluateResponse{}, err
	}

	var response EvaluateResponse
	err = errorutils.CheckError(json.Unmarshal(body, &response))
	return response, err
}
