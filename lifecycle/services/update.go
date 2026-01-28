package services

type updateOperation struct {
	reqBody        RbUpdateBody
	params         CommonOptionalQueryParams
	signingKeyName string
}

func (u *updateOperation) getOperationRestApi() string {
	return releaseBundleBaseApi
}

func (u *updateOperation) getRequestBody() any {
	return u.reqBody
}

func (u *updateOperation) getOperationSuccessfulMsg() string {
	return "Release Bundle successfully updated"
}

func (u *updateOperation) getOperationParams() CommonOptionalQueryParams {
	return u.params
}

func (u *updateOperation) getSigningKeyName() string {
	return u.signingKeyName
}

// UpdateReleaseBundleFromMultipleSources updates an existing draft release bundle by adding sources
func (rbs *ReleaseBundlesService) UpdateReleaseBundleFromMultipleSources(rbDetails ReleaseBundleDetails, params CommonOptionalQueryParams,
	signingKeyName string, addSources []RbSource) (response []byte, err error) {
	operation := updateOperation{
		reqBody: RbUpdateBody{
			ReleaseBundleDetails: rbDetails,
			AddSources:           addSources,
		},
		params:         params,
		signingKeyName: signingKeyName,
	}
	return rbs.doPatchOperation(&operation)
}

type RbUpdateBody struct {
	ReleaseBundleDetails
	AddSources []RbSource `json:"add_sources,omitempty"`
}
