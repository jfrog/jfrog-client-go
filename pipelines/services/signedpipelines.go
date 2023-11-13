package services

import (
	"encoding/json"
	"net/http"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type SignedPipelinesService struct {
	client *jfroghttpclient.JfrogHttpClient
	auth.ServiceDetails
}

const (
	Artifact = iota
	BuildInfo
	ReleaseBundle
)

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

func (sp *SignedPipelinesService) ValidateSignedPipelines(artifactTypeInfo ArtifactTypeInfo, artifactType int) error {
	// Fetch pipeline resource to retrieve resource ID
	log.Info("Validating signed pipelines for", artifactType)
	httpDetails := sp.getHttpDetails()
	queryParams := sp.constructQueryParamsBasedOnArtifactType(artifactTypeInfo, artifactType)
	uriVal, err := constructPipelinesURL(queryParams, sp.GetUrl(), validatePipelines)
	if err != nil {
		return err
	}
	resp, body, _, httpErr := sp.client.SendGet(uriVal, true, &httpDetails)
	if httpErr != nil {
		return httpErr
	}
	if err := errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return err
	}
	return parseValidateSignedPipelinesResponse(body)
}

func (sp *SignedPipelinesService) constructQueryParamsBasedOnArtifactType(artifactTypeInfo ArtifactTypeInfo, artifactType int) map[string]string {
	queryParams := map[string]string{}
	switch artifactType {
	case BuildInfo:
		queryParams = map[string]string{
			signedPipelinesArtifactType: "buildInfo",
			"buildName":                 artifactTypeInfo.BuildName,
			"buildNumber":               artifactTypeInfo.BuildNumber,
			"projectKey":                artifactTypeInfo.ProjectKey,
		}
	case Artifact:
		queryParams = map[string]string{
			signedPipelinesArtifactType: "artifact",
			"artifactPath":              artifactTypeInfo.ArtifactPath,
		}
	case ReleaseBundle:
		queryParams = map[string]string{
			signedPipelinesArtifactType: "releaseBundle",
			"rbName":                    artifactTypeInfo.RbName,
			"rbVersion":                 artifactTypeInfo.RbVersion,
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
		log.Output("Signed Pipelines validation failed with below message/messages")
		for _, message := range signedPipelinesValidationResponse.Messages {
			log.Output(message)
		}
		log.Output(signedPipelinesValidationResponse.Message)
		return errorutils.CheckErrorf("signed pipelines validation failed")
	}
	log.Output("Signed Pipelines validation is completed successfully")
	return nil
}
