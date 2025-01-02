package lifecycle

import (
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	lifecycle "github.com/jfrog/jfrog-client-go/lifecycle/services"
	"github.com/jfrog/jfrog-client-go/utils/distribution"
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

func (lcs *LifecycleServicesManager) CreateReleaseBundleFromAql(rbDetails lifecycle.ReleaseBundleDetails,
	queryParams lifecycle.CommonOptionalQueryParams, signingKeyName string, aqlQuery string) error {
	rbService := lifecycle.NewReleaseBundlesService(lcs.config.GetServiceDetails(), lcs.client)
	return rbService.CreateFromAql(rbDetails, queryParams, signingKeyName, aqlQuery)
}

func (lcs *LifecycleServicesManager) CreateReleaseBundleFromArtifacts(rbDetails lifecycle.ReleaseBundleDetails,
	queryParams lifecycle.CommonOptionalQueryParams, signingKeyName string, sourceArtifacts lifecycle.CreateFromArtifacts) error {
	rbService := lifecycle.NewReleaseBundlesService(lcs.config.GetServiceDetails(), lcs.client)
	return rbService.CreateFromArtifacts(rbDetails, queryParams, signingKeyName, sourceArtifacts)
}

func (lcs *LifecycleServicesManager) CreateReleaseBundleFromBuilds(rbDetails lifecycle.ReleaseBundleDetails,
	queryParams lifecycle.CommonOptionalQueryParams, signingKeyName string, sourceBuilds lifecycle.CreateFromBuildsSource) error {
	rbService := lifecycle.NewReleaseBundlesService(lcs.config.GetServiceDetails(), lcs.client)
	return rbService.CreateFromBuilds(rbDetails, queryParams, signingKeyName, sourceBuilds)
}

func (lcs *LifecycleServicesManager) CreateReleaseBundleFromBundles(rbDetails lifecycle.ReleaseBundleDetails,
	queryParams lifecycle.CommonOptionalQueryParams, signingKeyName string, sourceReleaseBundles lifecycle.CreateFromReleaseBundlesSource) error {
	rbService := lifecycle.NewReleaseBundlesService(lcs.config.GetServiceDetails(), lcs.client)
	return rbService.CreateFromBundles(rbDetails, queryParams, signingKeyName, sourceReleaseBundles)
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

	mappings := &distributeBundleService.Modifications.PathMappings
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

func (lcs *LifecycleServicesManager) IsReleaseBundleExist(projectKey, releaseBundleNameAndVersion string) (bool, error) {
	rbService := lifecycle.NewReleaseBundlesService(lcs.config.GetServiceDetails(), lcs.client)
	return rbService.IsExists(projectKey, releaseBundleNameAndVersion)
}
