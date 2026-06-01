package services

import (
	"fmt"
	"net/http"

	"strconv"

	rtUtils "github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/distribution"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type finalizeOperation struct {
	rbDetails      ReleaseBundleDetails
	params         CommonOptionalQueryParams
	signingKeyName string
}

func (f *finalizeOperation) getOperationRestApi() string {
	return fmt.Sprintf("%s/%s/%s/finalize", releaseBundleBaseApi, f.rbDetails.ReleaseBundleName, f.rbDetails.ReleaseBundleVersion)
}

func (f *finalizeOperation) getRequestBody() any {
	return nil
}

func (f *finalizeOperation) getOperationSuccessfulMsg() string {
	return "Release Bundle successfully finalized"
}

func (f *finalizeOperation) getOperationParams() CommonOptionalQueryParams {
	return f.params
}

func (f *finalizeOperation) getSigningKeyName() string {
	return f.signingKeyName
}

// FinalizeReleaseBundle finalizes a draft release bundle
func (rbs *ReleaseBundlesService) FinalizeReleaseBundle(rbDetails ReleaseBundleDetails, params CommonOptionalQueryParams, signingKeyName string) ([]byte, error) {
	operation := &finalizeOperation{
		rbDetails:      rbDetails,
		params:         params,
		signingKeyName: signingKeyName,
	}
	return rbs.doFinalizeOperation(operation)
}

// doFinalizeOperation performs a POST request without a body for finalize operation
func (rbs *ReleaseBundlesService) doFinalizeOperation(operation *finalizeOperation) (response []byte, err error) {
	queryParams := buildFinalizeQueryParams(operation.params)
	requestFullUrl, err := utils.BuildUrl(rbs.GetLifecycleDetails().GetUrl(), operation.getOperationRestApi(), queryParams)
	if err != nil {
		return nil, err
	}

	httpClientDetails := rbs.GetLifecycleDetails().CreateHttpClientDetails()
	rtUtils.AddSigningKeyNameHeader(operation.signingKeyName, &httpClientDetails.Headers)
	httpClientDetails.SetContentTypeApplicationJson()

	resp, body, err := rbs.client.SendPost(requestFullUrl, nil, &httpClientDetails)
	if err != nil {
		return nil, err
	}

	return handleFinalizeResponse(operation, resp, body)
}

func buildFinalizeQueryParams(params CommonOptionalQueryParams) map[string]string {
	queryParams := distribution.GetProjectQueryParam(params.ProjectKey)
	queryParams[async] = strconv.FormatBool(params.Async)
	return queryParams
}

func handleFinalizeResponse(operation *finalizeOperation, resp *http.Response, body []byte) ([]byte, error) {
	// Finalize API returns 200 OK for both sync and async modes
	if err := errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return nil, err
	}
	log.Info(operation.getOperationSuccessfulMsg())
	return body, nil
}
