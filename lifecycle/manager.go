package lifecycle

import (
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

func (lcs *LifecycleServicesManager) DeleteReleaseBundle(rbDetails lifecycle.ReleaseBundleDetails, queryParams lifecycle.CommonOptionalQueryParams) error {
	rbService := lifecycle.NewReleaseBundlesService(lcs.config.GetServiceDetails(), lcs.client)
	return rbService.DeleteReleaseBundle(rbDetails, queryParams)
}

func (lcs *LifecycleServicesManager) DistributeReleaseBundle(params distribution.DistributionParams, autoCreateRepo bool, pathMapping lifecycle.PathMapping) error {
	distributeBundleService := lifecycle.NewDistributeReleaseBundleService(lcs.client)
	distributeBundleService.LcDetails = lcs.config.GetServiceDetails()
	distributeBundleService.DryRun = lcs.config.IsDryRun()
	distributeBundleService.AutoCreateRepo = autoCreateRepo
	distributeBundleService.DistributeParams = params
	distributeBundleService.PathMapping = pathMapping
	return distributeBundleService.Distribute()
}

func (lcs *LifecycleServicesManager) RemoteDeleteReleaseBundle(params distribution.DistributionParams, dryRun bool) error {
	rbService := lifecycle.NewReleaseBundlesService(lcs.config.GetServiceDetails(), lcs.client)
	return rbService.RemoteDeleteReleaseBundle(params, dryRun)
}
