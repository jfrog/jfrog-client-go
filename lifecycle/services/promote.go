package services

import (
	"encoding/json"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
	"path"
	"strconv"
)

const (
	promotionBaseApi = "api/v2/promotion"
	records          = "records"
)

type promoteOperation struct {
	reqBody        RbPromotionBody
	rbDetails      ReleaseBundleDetails
	queryParams    CommonOptionalQueryParams
	signingKeyName string
}

func (p *promoteOperation) getOperationRestApi() string {
	return path.Join(promotionBaseApi, records, p.rbDetails.ReleaseBundleName, p.rbDetails.ReleaseBundleVersion)
}

func (p *promoteOperation) getRequestBody() any {
	return p.reqBody
}

func (p *promoteOperation) getOperationSuccessfulMsg() string {
	return "Release Bundle successfully promoted"
}

func (p *promoteOperation) getOperationParams() CommonOptionalQueryParams {
	return p.queryParams
}

func (p *promoteOperation) getSigningKeyName() string {
	return p.signingKeyName
}

func (rbs *ReleaseBundlesService) Promote(rbDetails ReleaseBundleDetails, queryParams CommonOptionalQueryParams, signingKeyName string, promotionParams RbPromotionParams) (RbPromotionResp, error) {
	operation := promoteOperation{
		reqBody:        RbPromotionBody(promotionParams),
		rbDetails:      rbDetails,
		queryParams:    queryParams,
		signingKeyName: signingKeyName,
	}
	respBody, err := rbs.doPostOperation(&operation)
	if err != nil {
		return RbPromotionResp{}, err
	}
	var promotionResp RbPromotionResp
	err = json.Unmarshal(respBody, &promotionResp)
	return promotionResp, errorutils.CheckError(err)
}

type RbPromotionsResponse struct {
	Promotions []RbPromotion `json:"promotions,omitempty"`
}

type RbPromotion struct {
	Status               RbStatus    `json:"status,omitempty"`
	RepositoryKey        string      `json:"repository_key,omitempty"`
	ReleaseBundleName    string      `json:"release_bundle_name,omitempty"`
	ReleaseBundleVersion string      `json:"release_bundle_version,omitempty"`
	Environment          string      `json:"environment,omitempty"`
	ServiceId            string      `json:"service_id,omitempty"`
	CreatedBy            string      `json:"created_by,omitempty"`
	Created              string      `json:"created,omitempty"`
	CreatedMillis        json.Number `json:"created_millis,omitempty"`
	Messages             []Message   `json:"messages,omitempty"`
}

type GetPromotionsOptionalQueryParams struct {
	Include    string
	Offset     int
	Limit      int
	FilterBy   string
	OrderBy    string
	OrderAsc   bool
	ProjectKey string
}

func buildGetPromotionsQueryParams(optionalQueryParams GetPromotionsOptionalQueryParams) map[string]string {
	params := make(map[string]string)
	if optionalQueryParams.ProjectKey != "" {
		params["project"] = optionalQueryParams.ProjectKey
	}
	if optionalQueryParams.Include != "" {
		params["include"] = optionalQueryParams.Include
	}
	if optionalQueryParams.Offset > 0 {
		params["offset"] = strconv.Itoa(optionalQueryParams.Offset)
	}
	if optionalQueryParams.Limit > 0 {
		params["limit"] = strconv.Itoa(optionalQueryParams.Limit)
	}
	if optionalQueryParams.FilterBy != "" {
		params["filter_by"] = optionalQueryParams.FilterBy
	}
	if optionalQueryParams.OrderBy != "" {
		params["order_by"] = optionalQueryParams.OrderBy
	}
	if optionalQueryParams.OrderAsc {
		params["order_asc"] = strconv.FormatBool(optionalQueryParams.OrderAsc)
	}
	return params
}

func (rbs *ReleaseBundlesService) GetReleaseBundleVersionPromotions(rbDetails ReleaseBundleDetails, optionalQueryParams GetPromotionsOptionalQueryParams) (response RbPromotionsResponse, err error) {
	restApi := GetGetReleaseBundleVersionPromotionsApi(rbDetails)
	requestFullUrl, err := utils.BuildUrl(rbs.GetLifecycleDetails().GetUrl(), restApi, buildGetPromotionsQueryParams(optionalQueryParams))
	if err != nil {
		return
	}
	httpClientsDetails := rbs.GetLifecycleDetails().CreateHttpClientDetails()
	resp, body, _, err := rbs.client.SendGet(requestFullUrl, true, &httpClientsDetails)
	if err != nil {
		return
	}
	log.Debug("Artifactory response:", resp.Status)
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return
	}
	err = errorutils.CheckError(json.Unmarshal(body, &response))
	return
}

func GetGetReleaseBundleVersionPromotionsApi(rbDetails ReleaseBundleDetails) string {
	return path.Join(promotionBaseApi, records, rbDetails.ReleaseBundleName, rbDetails.ReleaseBundleVersion)
}

type RbPromotionParams struct {
	Environment            string
	IncludedRepositoryKeys []string
	ExcludedRepositoryKeys []string
}

type RbPromotionBody struct {
	Environment            string   `json:"environment,omitempty"`
	IncludedRepositoryKeys []string `json:"included_repository_keys,omitempty"`
	ExcludedRepositoryKeys []string `json:"excluded_repository_keys,omitempty"`
}

type RbPromotionResp struct {
	RepositoryKey string `json:"repository_key,omitempty"`
	ReleaseBundleDetails
	RbPromotionBody
	Created       string      `json:"created,omitempty"`
	CreatedMillis json.Number `json:"created_millis,omitempty"`
}
