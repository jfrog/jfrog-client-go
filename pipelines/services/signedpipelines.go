package services

import (
	"encoding/json"
	"errors"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
)

type SignedPipelinesService struct {
	client *jfroghttpclient.JfrogHttpClient
	auth.ServiceDetails
}

const (
	validatePipelines           = "api/v1/pipeinfo/verify"
	signedPipelinesArtifactType = "artifactType"
)

type SignedPipelinesValidation struct {
	Result   bool     `json:"result"`
	Messages []string `json:"messages"`
	Message  string   `json:"message"`
}

func (sp *SignedPipelinesService) getHttpDetails() httputils.HttpClientDetails {
	return sp.ServiceDetails.CreateHttpClientDetails()
}

func NewSignedPipelinesService(client *jfroghttpclient.JfrogHttpClient) *SignedPipelinesService {
	return &SignedPipelinesService{client: client}
}

func (sp *SignedPipelinesService) ValidateSignedPipelines(artifactType, buildName, buildNumber, projectKey, artifactPath, rbName, rbVersion string) error {
	// Fetch pipeline resource to retrieve resource ID
	log.Info("Validating signed pipelines for ", artifactType)
	httpDetails := sp.getHttpDetails()
	queryParams := sp.constructQueryParamsBasedOnArtifactType(artifactType, buildName, buildNumber, projectKey, artifactPath, rbName, rbVersion)
	uriVal, err := constructPipelinesURL(queryParams, sp.ServiceDetails.GetUrl(), validatePipelines)
	if err != nil {
		return err
	}
	resp, body, _, httpErr := sp.client.SendGet(uriVal, true, &httpDetails)
	if httpErr != nil {
		return errorutils.CheckError(httpErr)
	}
	if err := errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return errorutils.CheckError(err)
	}
	return parseValidateSignedPipelinesResponse(body)
}

func (sp *SignedPipelinesService) constructQueryParamsBasedOnArtifactType(artifactType, buildName, buildNumber, projectKey, artifactPath, rbName, rbVersion string) map[string]string {
	queryParams := map[string]string{}
	switch artifactType {
	case "buildInfo":
		queryParams = map[string]string{
			signedPipelinesArtifactType: artifactType,
			"buildName":                 buildName,
			"buildNumber":               buildNumber,
			"projectKey":                projectKey,
		}
	case "artifact":
		queryParams = map[string]string{
			signedPipelinesArtifactType: artifactType,
			"artifactPath":              artifactPath,
		}
	case "releaseBundle":
		queryParams = map[string]string{
			signedPipelinesArtifactType: artifactType,
			"rbName":                    rbName,
			"rbVersion":                 rbVersion,
		}
	}
	return queryParams
}

func parseValidateSignedPipelinesResponse(body []byte) error {
	signedPipelinesValidationResponse := SignedPipelinesValidation{}
	jsonErr := json.Unmarshal(body, &signedPipelinesValidationResponse)
	if jsonErr != nil {
		return errorutils.CheckError(jsonErr)
	}
	if !signedPipelinesValidationResponse.Result {
		log.Output("Validation failed with below message/messages")
		for _, message := range signedPipelinesValidationResponse.Messages {
			log.Output(message)
		}
		log.Output(signedPipelinesValidationResponse.Message)
		return errorutils.CheckError(errors.New("Signed pipelines validation failed"))
	}
	log.Output("Validation is completed successfully")
	return nil
}
