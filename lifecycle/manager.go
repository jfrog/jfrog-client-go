package lifecycle

import (
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	lifecycle "github.com/jfrog/jfrog-client-go/lifecycle/services"
	"github.com/jfrog/jfrog-client-go/utils/distribution"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
)

type LifecycleServicesManager struct {
	client *jfroghttpclient.JfrogHttpClient
	config config.Config
}

func New(config config.Config) (*LifecycleServicesManager, error) {
	details := config.GetServiceDetails()
	var err error
	manager := &LifecycleServicesManager{config: config}
	manager.client, err = jfroghttpclient.JfrogClientBuilder().
		SetCertificatesPath(config.GetCertificatesPath()).
		SetInsecureTls(config.IsInsecureTls()).
		SetClientCertPath(details.GetClientCertPath()).
		SetClientCertKeyPath(details.GetClientCertKeyPath()).
		AppendPreRequestInterceptor(details.RunPreRequestFunctions).
		SetContext(config.GetContext()).
		SetDialTimeout(config.GetDialTimeout()).
		SetOverallRequestTimeout(config.GetOverallRequestTimeout()).
		SetRetries(config.GetHttpRetries()).
		SetRetryWaitMilliSecs(config.GetHttpRetryWaitMilliSecs()).
		Build()

	return manager, err
}

func (lcs *LifecycleServicesManager) Client() *jfroghttpclient.JfrogHttpClient {
	return lcs.client
}

func (lcs *LifecycleServicesManager) IsDryRun() bool {
	return lcs.config.IsDryRun()
}

// CreateReleaseBundleFromAql creates a release bundle from AQL query (backward compatible, draft defaults to false)
func (lcs *LifecycleServicesManager) CreateReleaseBundleFromAql(rbDetails lifecycle.ReleaseBundleDetails,
	queryParams lifecycle.CommonOptionalQueryParams, signingKeyName string, aqlQuery string) error {
	return lcs.CreateReleaseBundleFromAqlDraft(rbDetails, queryParams, signingKeyName, aqlQuery, false)
}

// CreateReleaseBundleFromAqlDraft creates a release bundle from AQL query with draft option
func (lcs *LifecycleServicesManager) CreateReleaseBundleFromAqlDraft(rbDetails lifecycle.ReleaseBundleDetails,
	queryParams lifecycle.CommonOptionalQueryParams, signingKeyName string, aqlQuery string, draft bool) error {
	rbService := lifecycle.NewReleaseBundlesService(lcs.config.GetServiceDetails(), lcs.client)
	return rbService.CreateFromAqlDraft(rbDetails, queryParams, signingKeyName, aqlQuery, draft)
}

// CreateReleaseBundleFromArtifacts creates a release bundle from artifacts (backward compatible, draft defaults to false)
func (lcs *LifecycleServicesManager) CreateReleaseBundleFromArtifacts(rbDetails lifecycle.ReleaseBundleDetails,
	queryParams lifecycle.CommonOptionalQueryParams, signingKeyName string, sourceArtifacts lifecycle.CreateFromArtifacts) error {
	return lcs.CreateReleaseBundleFromArtifactsDraft(rbDetails, queryParams, signingKeyName, sourceArtifacts, false)
}

// CreateReleaseBundleFromArtifactsDraft creates a release bundle from artifacts with draft option
func (lcs *LifecycleServicesManager) CreateReleaseBundleFromArtifactsDraft(rbDetails lifecycle.ReleaseBundleDetails,
	queryParams lifecycle.CommonOptionalQueryParams, signingKeyName string, sourceArtifacts lifecycle.CreateFromArtifacts, draft bool) error {
	rbService := lifecycle.NewReleaseBundlesService(lcs.config.GetServiceDetails(), lcs.client)
	return rbService.CreateFromArtifactsDraft(rbDetails, queryParams, signingKeyName, sourceArtifacts, draft)
}

// CreateReleaseBundleFromBuilds creates a release bundle from builds (backward compatible, draft defaults to false)
func (lcs *LifecycleServicesManager) CreateReleaseBundleFromBuilds(rbDetails lifecycle.ReleaseBundleDetails,
	queryParams lifecycle.CommonOptionalQueryParams, signingKeyName string, sourceBuilds lifecycle.CreateFromBuildsSource) error {
	return lcs.CreateReleaseBundleFromBuildsDraft(rbDetails, queryParams, signingKeyName, sourceBuilds, false)
}

// CreateReleaseBundleFromBuildsDraft creates a release bundle from builds with draft option
func (lcs *LifecycleServicesManager) CreateReleaseBundleFromBuildsDraft(rbDetails lifecycle.ReleaseBundleDetails,
	queryParams lifecycle.CommonOptionalQueryParams, signingKeyName string, sourceBuilds lifecycle.CreateFromBuildsSource, draft bool) error {
	rbService := lifecycle.NewReleaseBundlesService(lcs.config.GetServiceDetails(), lcs.client)
	return rbService.CreateFromBuildsDraft(rbDetails, queryParams, signingKeyName, sourceBuilds, draft)
}

// CreateReleaseBundleFromBundles creates a release bundle from other release bundles (backward compatible, draft defaults to false)
func (lcs *LifecycleServicesManager) CreateReleaseBundleFromBundles(rbDetails lifecycle.ReleaseBundleDetails,
	queryParams lifecycle.CommonOptionalQueryParams, signingKeyName string, sourceReleaseBundles lifecycle.CreateFromReleaseBundlesSource) error {
	return lcs.CreateReleaseBundleFromBundlesDraft(rbDetails, queryParams, signingKeyName, sourceReleaseBundles, false)
}

// CreateReleaseBundleFromBundlesDraft creates a release bundle from other release bundles with draft option
func (lcs *LifecycleServicesManager) CreateReleaseBundleFromBundlesDraft(rbDetails lifecycle.ReleaseBundleDetails,
	queryParams lifecycle.CommonOptionalQueryParams, signingKeyName string, sourceReleaseBundles lifecycle.CreateFromReleaseBundlesSource, draft bool) error {
	rbService := lifecycle.NewReleaseBundlesService(lcs.config.GetServiceDetails(), lcs.client)
	return rbService.CreateFromBundlesDraft(rbDetails, queryParams, signingKeyName, sourceReleaseBundles, draft)
}

// CreateReleaseBundleFromPackages creates a release bundle from packages (backward compatible, draft defaults to false)
func (lcs *LifecycleServicesManager) CreateReleaseBundleFromPackages(rbDetails lifecycle.ReleaseBundleDetails,
	queryParams lifecycle.CommonOptionalQueryParams, signingKeyName string, packageSource lifecycle.CreateFromPackagesSource) error {
	return lcs.CreateReleaseBundleFromPackagesDraft(rbDetails, queryParams, signingKeyName, packageSource, false)
}

// CreateReleaseBundleFromPackagesDraft creates a release bundle from packages with draft option
func (lcs *LifecycleServicesManager) CreateReleaseBundleFromPackagesDraft(rbDetails lifecycle.ReleaseBundleDetails,
	queryParams lifecycle.CommonOptionalQueryParams, signingKeyName string, packageSource lifecycle.CreateFromPackagesSource, draft bool) error {
	rbService := lifecycle.NewReleaseBundlesService(lcs.config.GetServiceDetails(), lcs.client)
	return rbService.CreateFromPackagesDraft(rbDetails, queryParams, signingKeyName, packageSource, draft)
}

// CreateReleaseBundlesFromMultipleSources creates a release bundle from multiple sources (backward compatible, draft defaults to false)
func (lcs *LifecycleServicesManager) CreateReleaseBundlesFromMultipleSources(rbDetails lifecycle.ReleaseBundleDetails, queryParams lifecycle.CommonOptionalQueryParams, signingKeyName string, sources []lifecycle.RbSource) (response []byte, err error) {
	return lcs.CreateReleaseBundlesFromMultipleSourcesDraft(rbDetails, queryParams, signingKeyName, sources, false)
}

// CreateReleaseBundlesFromMultipleSourcesDraft creates a release bundle from multiple sources with draft option
func (lcs *LifecycleServicesManager) CreateReleaseBundlesFromMultipleSourcesDraft(rbDetails lifecycle.ReleaseBundleDetails, queryParams lifecycle.CommonOptionalQueryParams, signingKeyName string, sources []lifecycle.RbSource, draft bool) (response []byte, err error) {
	rbService := lifecycle.NewReleaseBundlesService(lcs.config.GetServiceDetails(), lcs.client)
	resp, err := rbService.CreateReleaseBundleFromMultipleSourcesDraft(rbDetails, queryParams, signingKeyName, sources, draft)
	return resp, errorutils.CheckError(err)
}

// UpdateReleaseBundleFromMultipleSources updates an existing draft release bundle by adding sources
func (lcs *LifecycleServicesManager) UpdateReleaseBundleFromMultipleSources(rbDetails lifecycle.ReleaseBundleDetails, queryParams lifecycle.CommonOptionalQueryParams, signingKeyName string, addSources []lifecycle.RbSource) (response []byte, err error) {
	rbService := lifecycle.NewReleaseBundlesService(lcs.config.GetServiceDetails(), lcs.client)
	resp, err := rbService.UpdateReleaseBundleFromMultipleSources(rbDetails, queryParams, signingKeyName, addSources)
	return resp, errorutils.CheckError(err)
}

// FinalizeReleaseBundle finalizes a draft release bundle
func (lcs *LifecycleServicesManager) FinalizeReleaseBundle(rbDetails lifecycle.ReleaseBundleDetails, queryParams lifecycle.CommonOptionalQueryParams, signingKeyName string) ([]byte, error) {
	rbService := lifecycle.NewReleaseBundlesService(lcs.config.GetServiceDetails(), lcs.client)
	resp, err := rbService.FinalizeReleaseBundle(rbDetails, queryParams, signingKeyName)
	return resp, errorutils.CheckError(err)
}

func (lcs *LifecycleServicesManager) GetReleaseBundleSpecification(rbDetails lifecycle.ReleaseBundleDetails) (lifecycle.ReleaseBundleSpecResponse, error) {
	rbService := lifecycle.NewReleaseBundlesService(lcs.config.GetServiceDetails(), lcs.client)
	return rbService.GetReleaseBundleSpecification(rbDetails)
}

func (lcs *LifecycleServicesManager) PromoteReleaseBundle(rbDetails lifecycle.ReleaseBundleDetails, queryParams lifecycle.CommonOptionalQueryParams, signingKeyName string, promotionParams lifecycle.RbPromotionParams) (lifecycle.RbPromotionResp, error) {
	rbService := lifecycle.NewReleaseBundlesService(lcs.config.GetServiceDetails(), lcs.client)
	return rbService.Promote(rbDetails, queryParams, signingKeyName, promotionParams)
}

func (lcs *LifecycleServicesManager) GetReleaseBundleCreationStatus(rbDetails lifecycle.ReleaseBundleDetails, projectKey string, sync bool) (lifecycle.ReleaseBundleStatusResponse, error) {
	rbService := lifecycle.NewReleaseBundlesService(lcs.config.GetServiceDetails(), lcs.client)
	return rbService.GetReleaseBundleCreationStatus(rbDetails, projectKey, sync)
}

func (lcs *LifecycleServicesManager) GetReleaseBundlePromotionStatus(rbDetails lifecycle.ReleaseBundleDetails, projectKey, createdMillis string, sync bool) (lifecycle.ReleaseBundleStatusResponse, error) {
	rbService := lifecycle.NewReleaseBundlesService(lcs.config.GetServiceDetails(), lcs.client)
	return rbService.GetReleaseBundlePromotionStatus(rbDetails, projectKey, createdMillis, sync)
}

func (lcs *LifecycleServicesManager) GetReleaseBundleVersionPromotions(rbDetails lifecycle.ReleaseBundleDetails, optionalQueryParams lifecycle.GetPromotionsOptionalQueryParams) (lifecycle.RbPromotionsResponse, error) {
	rbService := lifecycle.NewReleaseBundlesService(lcs.config.GetServiceDetails(), lcs.client)
	return rbService.GetReleaseBundleVersionPromotions(rbDetails, optionalQueryParams)
}

func (lcs *LifecycleServicesManager) DeleteReleaseBundleVersion(rbDetails lifecycle.ReleaseBundleDetails, queryParams lifecycle.CommonOptionalQueryParams) error {
	rbService := lifecycle.NewReleaseBundlesService(lcs.config.GetServiceDetails(), lcs.client)
	return rbService.DeleteReleaseBundleVersion(rbDetails, queryParams)
}

func (lcs *LifecycleServicesManager) DeleteReleaseBundleVersionPromotion(rbDetails lifecycle.ReleaseBundleDetails, queryParams lifecycle.CommonOptionalQueryParams, created string) error {
	rbService := lifecycle.NewReleaseBundlesService(lcs.config.GetServiceDetails(), lcs.client)
	return rbService.DeleteReleaseBundleVersionPromotion(rbDetails, queryParams, created)
}

func (lcs *LifecycleServicesManager) DistributeReleaseBundle(rbDetails lifecycle.ReleaseBundleDetails, distributeParams lifecycle.DistributeReleaseBundleParams) error {
	distributeBundleService := lifecycle.NewDistributeReleaseBundleService(lcs.client)
	distributeBundleService.LcDetails = lcs.config.GetServiceDetails()
	distributeBundleService.DryRun = lcs.config.IsDryRun()

	distributeBundleService.DistributeParams = distribution.DistributionParams{
		Name:              rbDetails.ReleaseBundleName,
		Version:           rbDetails.ReleaseBundleVersion,
		DistributionRules: distributeParams.DistributionRules,
	}
	distributeBundleService.AutoCreateRepo = distributeParams.AutoCreateRepo
	distributeBundleService.Sync = distributeParams.Sync
	distributeBundleService.MaxWaitMinutes = distributeParams.MaxWaitMinutes
	distributeBundleService.ProjectKey = distributeParams.ProjectKey

	mappings := &distributeBundleService.PathMappings
	*mappings = []utils.PathMapping{}
	for _, pathMapping := range distributeParams.PathMappings {
		*mappings = append(*mappings,
			distribution.CreatePathMappingsFromPatternAndTarget(pathMapping.Pattern, pathMapping.Target)...)
	}
	return distributeBundleService.Distribute()
}

func (lcs *LifecycleServicesManager) RemoteDeleteReleaseBundle(rbDetails lifecycle.ReleaseBundleDetails, params lifecycle.ReleaseBundleRemoteDeleteParams) error {
	rbService := lifecycle.NewReleaseBundlesService(lcs.config.GetServiceDetails(), lcs.client)
	return rbService.RemoteDeleteReleaseBundle(rbDetails, params)
}

func (lcs *LifecycleServicesManager) ExportReleaseBundle(rbDetails lifecycle.ReleaseBundleDetails, modifications lifecycle.Modifications, queryParams lifecycle.CommonOptionalQueryParams) (exportResponse lifecycle.ReleaseBundleExportedStatusResponse, err error) {
	rbService := lifecycle.NewReleaseBundlesService(lcs.config.GetServiceDetails(), lcs.client)
	return rbService.ExportReleaseBundle(rbDetails, modifications, queryParams)
}

func (lcs *LifecycleServicesManager) IsReleaseBundleExist(rbName, rbVersion, projectKey string) (bool, error) {
	rbService := lifecycle.NewReleaseBundlesService(lcs.config.GetServiceDetails(), lcs.client)
	return rbService.ReleaseBundleExists(rbName, rbVersion, projectKey)
}

func (lcs *LifecycleServicesManager) AnnotateReleaseBundle(params lifecycle.AnnotateOperationParams) error {
	rbService := lifecycle.NewReleaseBundlesService(lcs.config.GetServiceDetails(), lcs.client)
	return rbService.AnnotateReleaseBundle(params)
}

func (lcs *LifecycleServicesManager) GetReleaseBundlesStats(serverUrl string) ([]byte, error) {
	rbService := lifecycle.NewReleaseBundlesStatsService(lcs.config.GetServiceDetails(), lcs.client)
	return rbService.GetReleaseBundlesStats(serverUrl)
}

func (lcs *LifecycleServicesManager) ReleaseBundlesSearchGroup(params lifecycle.GetSearchOptionalQueryParams) (lifecycle.ReleaseBundlesGroupResponse, error) {
	rbService := lifecycle.NewReleaseBundlesService(lcs.config.GetServiceDetails(), lcs.client)
	return rbService.ReleaseBundlesSearchGroups(params)
}

func (lcs *LifecycleServicesManager) ReleaseBundlesSearchVersions(releaseBundleName string, params lifecycle.GetSearchOptionalQueryParams) (lifecycle.ReleaseBundleVersionsResponse, error) {
	rbService := lifecycle.NewReleaseBundlesService(lcs.config.GetServiceDetails(), lcs.client)
	return rbService.ReleaseBundlesSearchVersions(releaseBundleName, params)
}
