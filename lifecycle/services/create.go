package services

const (
	releaseBundleBaseApi = "api/v2/release_bundle"
)

type SourceType string

const (
	Aql            SourceType = "aql"
	Artifacts      SourceType = "artifacts"
	Builds         SourceType = "builds"
	ReleaseBundles SourceType = "release_bundles"
	Packages       SourceType = "packages"
)

type createOperation struct {
	reqBody        RbCreationBody
	params         CommonOptionalQueryParams
	signingKeyName string
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

func (c *createOperation) getOperationParams() CommonOptionalQueryParams {
	return c.params
}

func (c *createOperation) getSigningKeyName() string {
	return c.signingKeyName
}

// CreateFromAql creates a release bundle from AQL query (backward compatible, draft defaults to false)
func (rbs *ReleaseBundlesService) CreateFromAql(rbDetails ReleaseBundleDetails, params CommonOptionalQueryParams, signingKeyName string, aqlQuery string) error {
	return rbs.CreateFromAqlDraft(rbDetails, params, signingKeyName, aqlQuery, false)
}

// CreateFromAqlDraft creates a release bundle from AQL query with draft option
func (rbs *ReleaseBundlesService) CreateFromAqlDraft(rbDetails ReleaseBundleDetails, params CommonOptionalQueryParams, signingKeyName string, aqlQuery string, draft bool) error {
	return rbs.CreateReleaseBundleDraft(rbDetails, params, signingKeyName, Aql, CreateFromAqlSource{Aql: aqlQuery}, draft)
}

// CreateFromArtifacts creates a release bundle from artifacts (backward compatible, draft defaults to false)
func (rbs *ReleaseBundlesService) CreateFromArtifacts(rbDetails ReleaseBundleDetails, params CommonOptionalQueryParams, signingKeyName string, sourceArtifacts CreateFromArtifacts) error {
	return rbs.CreateFromArtifactsDraft(rbDetails, params, signingKeyName, sourceArtifacts, false)
}

// CreateFromArtifactsDraft creates a release bundle from artifacts with draft option
func (rbs *ReleaseBundlesService) CreateFromArtifactsDraft(rbDetails ReleaseBundleDetails, params CommonOptionalQueryParams, signingKeyName string, sourceArtifacts CreateFromArtifacts, draft bool) error {
	return rbs.CreateReleaseBundleDraft(rbDetails, params, signingKeyName, Artifacts, sourceArtifacts, draft)
}

// CreateFromBuilds creates a release bundle from builds (backward compatible, draft defaults to false)
func (rbs *ReleaseBundlesService) CreateFromBuilds(rbDetails ReleaseBundleDetails, params CommonOptionalQueryParams, signingKeyName string, sourceBuilds CreateFromBuildsSource) error {
	return rbs.CreateFromBuildsDraft(rbDetails, params, signingKeyName, sourceBuilds, false)
}

// CreateFromBuildsDraft creates a release bundle from builds with draft option
func (rbs *ReleaseBundlesService) CreateFromBuildsDraft(rbDetails ReleaseBundleDetails, params CommonOptionalQueryParams, signingKeyName string, sourceBuilds CreateFromBuildsSource, draft bool) error {
	return rbs.CreateReleaseBundleDraft(rbDetails, params, signingKeyName, Builds, sourceBuilds, draft)
}

// CreateFromBundles creates a release bundle from other release bundles (backward compatible, draft defaults to false)
func (rbs *ReleaseBundlesService) CreateFromBundles(rbDetails ReleaseBundleDetails, params CommonOptionalQueryParams, signingKeyName string, sourceReleaseBundles CreateFromReleaseBundlesSource) error {
	return rbs.CreateFromBundlesDraft(rbDetails, params, signingKeyName, sourceReleaseBundles, false)
}

// CreateFromBundlesDraft creates a release bundle from other release bundles with draft option
func (rbs *ReleaseBundlesService) CreateFromBundlesDraft(rbDetails ReleaseBundleDetails, params CommonOptionalQueryParams, signingKeyName string, sourceReleaseBundles CreateFromReleaseBundlesSource, draft bool) error {
	return rbs.CreateReleaseBundleDraft(rbDetails, params, signingKeyName, ReleaseBundles, sourceReleaseBundles, draft)
}

// CreateFromPackages creates a release bundle from packages (backward compatible, draft defaults to false)
func (rbs *ReleaseBundlesService) CreateFromPackages(rbDetails ReleaseBundleDetails, params CommonOptionalQueryParams, signingKeyName string, sourcePackages CreateFromPackagesSource) error {
	return rbs.CreateFromPackagesDraft(rbDetails, params, signingKeyName, sourcePackages, false)
}

// CreateFromPackagesDraft creates a release bundle from packages with draft option
func (rbs *ReleaseBundlesService) CreateFromPackagesDraft(rbDetails ReleaseBundleDetails, params CommonOptionalQueryParams, signingKeyName string, sourcePackages CreateFromPackagesSource, draft bool) error {
	return rbs.CreateReleaseBundleDraft(rbDetails, params, signingKeyName, Packages, sourcePackages, draft)
}

// CreateReleaseBundleFromMultipleSources creates a release bundle from multiple sources (backward compatible, draft defaults to false)
func (rbs *ReleaseBundlesService) CreateReleaseBundleFromMultipleSources(rbDetails ReleaseBundleDetails, params CommonOptionalQueryParams,
	signingKeyName string, sources []RbSource) (response []byte, err error) {
	return rbs.CreateReleaseBundleFromMultipleSourcesDraft(rbDetails, params, signingKeyName, sources, false)
}

// CreateReleaseBundleFromMultipleSourcesDraft creates a release bundle from multiple sources with draft option
func (rbs *ReleaseBundlesService) CreateReleaseBundleFromMultipleSourcesDraft(rbDetails ReleaseBundleDetails, params CommonOptionalQueryParams,
	signingKeyName string, sources []RbSource, draft bool) (response []byte, err error) {
	operation := createOperation{
		reqBody: RbCreationBody{
			ReleaseBundleDetails: rbDetails,
			Sources:              sources,
			Draft:                draft,
		},
		params:         params,
		signingKeyName: signingKeyName,
	}
	response, err = rbs.doPostOperation(&operation)
	return response, err
}

// CreateReleaseBundle creates a release bundle (backward compatible, draft defaults to false)
func (rbs *ReleaseBundlesService) CreateReleaseBundle(rbDetails ReleaseBundleDetails, params CommonOptionalQueryParams,
	signingKeyName string, rbSourceType SourceType, source interface{}) error {
	return rbs.CreateReleaseBundleDraft(rbDetails, params, signingKeyName, rbSourceType, source, false)
}

// CreateReleaseBundleDraft creates a release bundle with draft option
func (rbs *ReleaseBundlesService) CreateReleaseBundleDraft(rbDetails ReleaseBundleDetails, params CommonOptionalQueryParams,
	signingKeyName string, rbSourceType SourceType, source interface{}, draft bool) error {
	operation := createOperation{
		reqBody: RbCreationBody{
			ReleaseBundleDetails: rbDetails,
			SourceType:           rbSourceType,
			Source:               source,
			Draft:                draft,
		},
		params:         params,
		signingKeyName: signingKeyName,
	}
	_, err := rbs.doPostOperation(&operation)
	return err
}

type CreateFromAqlSource struct {
	Aql string `json:"aql,omitempty"`
}

type SourceBuildDetails struct {
	BuildName   string
	BuildNumber string
	ProjectKey  string
}

type CreateFromArtifacts struct {
	Artifacts []ArtifactSource `json:"artifacts,omitempty"`
}

type CreateFromBuildsSource struct {
	Builds []BuildSource `json:"builds,omitempty"`
}

type CreateFromPackagesSource struct {
	Packages []PackageSource `json:"packages,omitempty"`
}

type ArtifactSource struct {
	Path   string `json:"path,omitempty"`
	Sha256 string `json:"sha256,omitempty"`
}

type BuildSource struct {
	BuildName           string `json:"build_name,omitempty"`
	BuildNumber         string `json:"build_number,omitempty"`
	BuildRepository     string `json:"build_repository,omitempty"`
	IncludeDependencies bool   `json:"include_dependencies,omitempty"`
}

type PackageSource struct {
	PackageName    string `json:"package_name,omitempty"`
	PackageVersion string `json:"package_version,omitempty"`
	PackageType    string `json:"package_type,omitempty"`
	RepositoryKey  string `json:"repository_key,omitempty"`
}

type CreateFromReleaseBundlesSource struct {
	ReleaseBundles []ReleaseBundleSource `json:"release_bundles,omitempty"`
}

type ReleaseBundleSource struct {
	ReleaseBundleName    string `json:"release_bundle_name,omitempty"`
	ReleaseBundleVersion string `json:"release_bundle_version,omitempty"`
	ProjectKey           string `json:"project,omitempty"`
	RepositoryKey        string `json:"repository_key,omitempty"`
}

type RbSource struct {
	SourceType     SourceType            `json:"source_type"`
	Builds         []BuildSource         `json:"builds,omitempty"`
	ReleaseBundles []ReleaseBundleSource `json:"release_bundles,omitempty"`
	Artifacts      []ArtifactSource      `json:"artifacts,omitempty"`
	Packages       []PackageSource       `json:"packages,omitempty"`
	Aql            string                `json:"aql,omitempty"`
}
type RbCreationBody struct {
	ReleaseBundleDetails
	SourceType SourceType  `json:"source_type,omitempty"`
	Source     interface{} `json:"source,omitempty"`
	Sources    []RbSource  `json:"sources,omitempty"`
	Draft      bool        `json:"draft,omitempty"`
}
