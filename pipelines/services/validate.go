package services

import (
	"encoding/json"
	"errors"
	"github.com/gookit/color"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
	"net/url"
	"strconv"
)

const (
	validatePipResourcePath = "api/v1/validateYaml"
)

type ValidateService struct {
	client *jfroghttpclient.JfrogHttpClient
	auth.ServiceDetails
}

func NewValidateService(client *jfroghttpclient.JfrogHttpClient) *ValidateService {
	return &ValidateService{client: client}
}

func (vs *ValidateService) getHttpDetails() httputils.HttpClientDetails {
	httpDetails := vs.ServiceDetails.CreateHttpClientDetails()
	return httpDetails
}

func (vs *ValidateService) ValidatePipeline(data []byte) error {
	var err error
	httpDetails := vs.getHttpDetails()

	// Query params
	m := make(map[string]string, 0)

	// URL Construction
	uri := vs.constructValidateAPIURL(m, validatePipResourcePath)

	// Headers
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	httpDetails.Headers = headers
	log.Debug(string(data))

	// Send post request
	resp, body, err := vs.client.SendPost(uri, data, &httpDetails)
	if err != nil {
		return err
	}
	if err := errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return err
	}

	// Process response
	if nil != resp && resp.StatusCode == http.StatusOK {
		log.Info("Processing resources yaml response ")
		log.Info(string(body))
		err = processValidatePipResourceResponse(body)
		if err != nil {
			log.Error("Failed to process resources validation response")
		}
	}
	return nil
}

// processValidatePipResourceResponse process validate pipeline resource response
func processValidatePipResourceResponse(resp []byte) error {
	validationResponse := make(map[string]ValidationResponse)
	err := json.Unmarshal(resp, &validationResponse)
	if err != nil {
		return err
	}
	if len(validationResponse) == 0 {
		return errors.New("pipelines not found")
	}
	for k, v := range validationResponse {
		if v.IsValid != nil && *v.IsValid {
			log.Info("Validation of pipeline resources completed successfully ")
			msg := color.Green.Sprintf("Validation completed")
			log.Info(msg)
		} else {
			fileName := color.Red.Sprintf("%s", k)
			log.Error(fileName)
			validationErrors := v.Errors
			for _, validationError := range validationErrors {
				lineNum := validationError.LineNumber
				errorMessage := color.Red.Sprintf("%s", validationError.Text+":"+strconv.Itoa(lineNum))
				log.Error(errorMessage)
			}
		}
	}
	return nil
}

// constructPipelinesURL creates URL with all required details to make api call
// like headers, queryParams, apiPath
func (vs *ValidateService) constructValidateAPIURL(qParams map[string]string, apiPath string) string {
	uri, err := url.Parse(vs.ServiceDetails.GetUrl() + apiPath)
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
