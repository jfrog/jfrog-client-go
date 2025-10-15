package services

import (
	"encoding/json"
	"fmt"
	rtUtils "github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/distribution"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
	"strconv"
	"time"
)

const async = "async"
const operation = "operation"

type ReleaseBundlesService struct {
	client    *jfroghttpclient.JfrogHttpClient
	lcDetails *auth.ServiceDetails
}

func NewReleaseBundlesService(lcDetails auth.ServiceDetails, client *jfroghttpclient.JfrogHttpClient) *ReleaseBundlesService {
	return &ReleaseBundlesService{lcDetails: &lcDetails, client: client}
}

func (rbs *ReleaseBundlesService) GetLifecycleDetails() auth.ServiceDetails {
	return *rbs.lcDetails
}

type ReleaseBundleOperation interface {
	getOperationRestApi() string
	getRequestBody() any
	getOperationSuccessfulMsg() string
	getOperationParams() CommonOptionalQueryParams
	getSigningKeyName() string
}

func (rbs *ReleaseBundlesService) doPostOperation(operation ReleaseBundleOperation) (response []byte, err error) {
	requestFullUrl, httpClientDetails, err := prepareRequest(operation, rbs)
	if err != nil {
		return
	}

	content, err := json.Marshal(operation.getRequestBody())
	if err != nil {
		err = errorutils.CheckError(err)
		return
	}

	resp, body, err := rbs.client.SendPost(requestFullUrl, content, &httpClientDetails)
	if err != nil {
		return
	}
	return handleResponse(operation, resp, body)
}

func (rbs *ReleaseBundlesService) doGetOperation(operation ReleaseBundleOperation) (response []byte, err error) {
	requestFullUrl, httpClientDetails, err := prepareRequest(operation, rbs)
	if err != nil {
		return
	}

	resp, body, _, err := rbs.client.SendGet(requestFullUrl, false, &httpClientDetails)
	if err != nil {
		return
	}
	return handleResponse(operation, resp, body)
}

func prepareRequest(operation ReleaseBundleOperation, rbs *ReleaseBundlesService) (requestFullUrl string, httpClientDetails httputils.HttpClientDetails, err error) {
	queryParams := buildQueryParams(operation.getOperationParams())
	requestFullUrl, err = utils.BuildUrl(rbs.GetLifecycleDetails().GetUrl(), operation.getOperationRestApi(), queryParams)
	if err != nil {
		return
	}
	httpClientDetails = rbs.GetLifecycleDetails().CreateHttpClientDetails()
	rtUtils.AddSigningKeyNameHeader(operation.getSigningKeyName(), &httpClientDetails.Headers)
	httpClientDetails.SetContentTypeApplicationJson()
	return
}

func buildQueryParams(commonOptionalQueryParams CommonOptionalQueryParams) map[string]string {
	params := distribution.GetProjectQueryParam(commonOptionalQueryParams.ProjectKey)
	params[async] = strconv.FormatBool(commonOptionalQueryParams.Async)
	if commonOptionalQueryParams.PromotionType != "" {
		params[operation] = commonOptionalQueryParams.PromotionType
	}
	return params
}

func handleResponse(operation ReleaseBundleOperation, resp *http.Response, body []byte) (response []byte, err error) {
	if !operation.getOperationParams().Async {
		if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusCreated); err != nil {
			return
		}
		log.Info(operation.getOperationSuccessfulMsg())
		return body, nil
	}

	return body, errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK, http.StatusAccepted)
}

func (rbs *ReleaseBundlesService) doHttpRequestWithRetry(
	requestFullUrl string,
	httpClientsDetails *httputils.HttpClientDetails,
	responsePtr interface{}) error {

	for i := 0; i < maxAttempts; i++ {
		resp, body, _, sendErr := rbs.client.SendGet(requestFullUrl, true, httpClientsDetails)
		if sendErr != nil {
			log.Debug(fmt.Sprintf("Attempt %d/%d: Network error during SendGet to %s: %v", i+1, maxAttempts, requestFullUrl, sendErr))
			if i < maxAttempts-1 {
				time.Sleep(utils.CalculateBackoff(i, initialBackoff, maxBackoff))
				continue
			}
			return sendErr
		}
		log.Debug("Artifactory response status:", resp.Status)
		switch resp.StatusCode {
		case http.StatusOK:
			return errorutils.CheckError(json.Unmarshal(body, responsePtr))
		case http.StatusUnauthorized:
			return ErrAuth
		case http.StatusForbidden:
			return ErrPermission
		case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
			if i < maxAttempts-1 {
				log.Debug(fmt.Sprintf("Attempt %d/%d: Server error (%d) from %s. Retrying...", i+1, maxAttempts, resp.StatusCode, requestFullUrl))
				time.Sleep(utils.CalculateBackoff(i, initialBackoff, maxBackoff))
				continue
			}
			fallthrough
		default:
			return errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK)
		}
	}
	return fmt.Errorf("failed to get response from %s after %d attempts", requestFullUrl, maxAttempts)
}

type ReleaseBundleDetails struct {
	ReleaseBundleName    string `json:"release_bundle_name,omitempty"`
	ReleaseBundleVersion string `json:"release_bundle_version,omitempty"`
}

type CommonOptionalQueryParams struct {
	ProjectKey    string
	Async         bool
	PromotionType string
}
