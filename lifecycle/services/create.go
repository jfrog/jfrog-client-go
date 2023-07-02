package services

const (
	releaseBundleBaseApi = "api/v2/release_bundle"
)

type sourceType string

const (
	builds         sourceType = "builds"
	releaseBundles sourceType = "release_bundles"
)

type createOperation struct {
	reqBody RbCreationBody
	params  CreateOrPromoteReleaseBundleParams
}

func (c *createOperation) getOperationRestApi() string {
	return releaseBundleBaseApi
}

func (c *createOperation) getRequestBody() any {
	return c.reqBody
}

func (c *createOperation) getOperationSuccessfulMsg() string {
	return "Release Bundle successfully created"
}

func (c *createOperation) getOperationParams() CreateOrPromoteReleaseBundleParams {
	return c.params
}

func (rbs *ReleaseBundlesService) CreateFromBuilds(rbDetails ReleaseBundleDetails, params CreateOrPromoteReleaseBundleParams, sourceBuilds CreateFromBuildsSource) error {
	operation := createOperation{
		reqBody: RbCreationBody{
			ReleaseBundleDetails: rbDetails,
			SourceType:           builds,
			Source:               sourceBuilds},
		params: params,
	}
	_, err := rbs.doOperation(&operation)
	return err
}

func (rbs *ReleaseBundlesService) CreateFromBundles(rbDetails ReleaseBundleDetails, params CreateOrPromoteReleaseBundleParams, sourceReleaseBundles CreateFromReleaseBundlesSource) error {
	operation := createOperation{
		reqBody: RbCreationBody{
			ReleaseBundleDetails: rbDetails,
			SourceType:           releaseBundles,
			Source:               sourceReleaseBundles},
		params: params,
	}
	_, err := rbs.doOperation(&operation)
	return err
}

type SourceBuildDetails struct {
	BuildName   string
	BuildNumber string
	ProjectKey  string
}

type CreateFromBuildsSource struct {
	Builds []BuildSource `json:"builds,omitempty"`
}

type BuildSource struct {
	BuildName       string `json:"build_name,omitempty"`
	BuildNumber     string `json:"build_number,omitempty"`
	BuildRepository string `json:"build_repository,omitempty"`
}

type CreateFromReleaseBundlesSource struct {
	ReleaseBundles []ReleaseBundleSource `json:"release_bundles,omitempty"`
}

type ReleaseBundleSource struct {
	ReleaseBundleName    string `json:"release_bundle_name,omitempty"`
	ReleaseBundleVersion string `json:"release_bundle_version,omitempty"`
	ProjectKey           string `json:"project_key,omitempty"`
}

type RbCreationBody struct {
	ReleaseBundleDetails
	SourceType sourceType  `json:"source_type,omitempty"`
	Source     interface{} `json:"source,omitempty"`
}
