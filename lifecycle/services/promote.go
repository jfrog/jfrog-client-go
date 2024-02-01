package services

import (
	"encoding/json"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"path"
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
	respBody, err := rbs.doOperation(&operation)
	if err != nil {
		return RbPromotionResp{}, err
	}
	var promotionResp RbPromotionResp
	err = json.Unmarshal(respBody, &promotionResp)
	return promotionResp, errorutils.CheckError(err)
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
