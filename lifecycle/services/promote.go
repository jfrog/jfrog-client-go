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
	reqBody   RbPromotionBody
	rbDetails ReleaseBundleDetails
	params    CreateOrPromoteReleaseBundleParams
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

func (p *promoteOperation) getOperationParams() CreateOrPromoteReleaseBundleParams {
	return p.params
}

func (rbs *ReleaseBundlesService) Promote(rbDetails ReleaseBundleDetails, params CreateOrPromoteReleaseBundleParams, environment string, overwrite bool) (RbPromotionResp, error) {
	operation := promoteOperation{
		reqBody: RbPromotionBody{
			Environment: environment,
			Overwrite:   overwrite,
		},
		rbDetails: rbDetails,
		params:    params,
	}
	respBody, err := rbs.doOperation(&operation)
	if err != nil {
		return RbPromotionResp{}, err
	}
	var promotionResp RbPromotionResp
	err = json.Unmarshal(respBody, &promotionResp)
	return promotionResp, errorutils.CheckError(err)
}

type RbPromotionBody struct {
	Environment            string   `json:"environment,omitempty"`
	Overwrite              bool     `json:"overwrite_existing_artifacts,omitempty"`
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
