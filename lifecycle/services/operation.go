package services

import (
	"encoding/json"
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
)

const async = "async"

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
	queryParams := distribution.GetProjectQueryParam(operation.getOperationParams().ProjectKey)
	queryParams[async] = strconv.FormatBool(operation.getOperationParams().Async)
	requestFullUrl, err = utils.BuildUrl(rbs.GetLifecycleDetails().GetUrl(), operation.getOperationRestApi(), queryParams)
	if err != nil {
		return
	}
	httpClientDetails = rbs.GetLifecycleDetails().CreateHttpClientDetails()
	rtUtils.AddSigningKeyNameHeader(operation.getSigningKeyName(), &httpClientDetails.Headers)
	rtUtils.SetContentType("application/json", &httpClientDetails.Headers)
	return
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

type ReleaseBundleDetails struct {
	ReleaseBundleName    string `json:"release_bundle_name,omitempty"`
	ReleaseBundleVersion string `json:"release_bundle_version,omitempty"`
}

type CommonOptionalQueryParams struct {
	ProjectKey string
	Async      bool
}
