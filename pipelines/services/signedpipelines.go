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

type ArtifactType int

const (
	Artifact ArtifactType = iota
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

func (a ArtifactType) String() string {
	switch a {
	case Artifact:
		return "artifact"
	case BuildInfo:
		return "buildInfo"
	case ReleaseBundle:
		return "releaseBundle"
	}
	return ""
}

func (sp *SignedPipelinesService) getHttpDetails() httputils.HttpClientDetails {
	return sp.ServiceDetails.CreateHttpClientDetails()
}

func NewSignedPipelinesService(client *jfroghttpclient.JfrogHttpClient) *SignedPipelinesService {
	return &SignedPipelinesService{client: client}
}

func (sp *SignedPipelinesService) ValidateSignedPipelines(artifactTypeInfo ArtifactTypeInfo, artifactType ArtifactType) error {
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

func (sp *SignedPipelinesService) constructQueryParamsBasedOnArtifactType(artifactTypeInfo ArtifactTypeInfo, artifactType ArtifactType) map[string]string {
	queryParams := map[string]string{}
	switch artifactType {
	case BuildInfo:
		queryParams = map[string]string{
			"buildName":                 artifactTypeInfo.BuildName,
			"buildNumber":               artifactTypeInfo.BuildNumber,
			"projectKey":                artifactTypeInfo.ProjectKey,
			signedPipelinesArtifactType: BuildInfo.String(),
		}
	case Artifact:
		queryParams = map[string]string{
			"artifactPath":              artifactTypeInfo.ArtifactPath,
			signedPipelinesArtifactType: Artifact.String(),
		}
	case ReleaseBundle:
		queryParams = map[string]string{
			"rbName":                    artifactTypeInfo.RbName,
			"rbVersion":                 artifactTypeInfo.RbVersion,
			signedPipelinesArtifactType: ReleaseBundle.String(),
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
