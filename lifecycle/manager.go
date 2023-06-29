package lifecycle

import (
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	lifecycle "github.com/jfrog/jfrog-client-go/lifecycle/services"
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
		SetRetries(config.GetHttpRetries()).
		SetRetryWaitMilliSecs(config.GetHttpRetryWaitMilliSecs()).
		Build()

	return manager, err
}

func (lcs *LifecycleServicesManager) Client() *jfroghttpclient.JfrogHttpClient {
	return lcs.client
}

func (lcs *LifecycleServicesManager) CreateReleaseBundleFromBuilds(rbDetails lifecycle.ReleaseBundleDetails,
	params lifecycle.CreateOrPromoteReleaseBundleParams, sourceBuilds lifecycle.CreateFromBuildsSource) error {
	rbService := lifecycle.NewReleaseBundlesService(lcs.config.GetServiceDetails(), lcs.client)
	return rbService.CreateFromBuilds(rbDetails, params, sourceBuilds)
}

func (lcs *LifecycleServicesManager) CreateReleaseBundleFromBundles(rbDetails lifecycle.ReleaseBundleDetails,
	params lifecycle.CreateOrPromoteReleaseBundleParams, sourceReleaseBundles lifecycle.CreateFromReleaseBundlesSource) error {
	rbService := lifecycle.NewReleaseBundlesService(lcs.config.GetServiceDetails(), lcs.client)
	return rbService.CreateFromBundles(rbDetails, params, sourceReleaseBundles)
}

func (lcs *LifecycleServicesManager) PromoteReleaseBundle(rbDetails lifecycle.ReleaseBundleDetails, params lifecycle.CreateOrPromoteReleaseBundleParams, environment string, overwrite bool) (lifecycle.RbPromotionResp, error) {
	rbService := lifecycle.NewReleaseBundlesService(lcs.config.GetServiceDetails(), lcs.client)
	return rbService.Promote(rbDetails, params, environment, overwrite)
}

func (lcs *LifecycleServicesManager) GetReleaseBundleCreationStatus(rbDetails lifecycle.ReleaseBundleDetails, projectKey string, sync bool) (lifecycle.ReleaseBundleStatusResponse, error) {
	rbService := lifecycle.NewReleaseBundlesService(lcs.config.GetServiceDetails(), lcs.client)
	return rbService.GetReleaseBundleCreationStatus(rbDetails, projectKey, sync)
}

func (lcs *LifecycleServicesManager) GetReleaseBundlePromotionStatus(rbDetails lifecycle.ReleaseBundleDetails, projectKey, createdMillis string, sync bool) (lifecycle.ReleaseBundleStatusResponse, error) {
	rbService := lifecycle.NewReleaseBundlesService(lcs.config.GetServiceDetails(), lcs.client)
	return rbService.GetReleaseBundlePromotionStatus(rbDetails, projectKey, createdMillis, sync)
}

func (lcs *LifecycleServicesManager) DeleteReleaseBundle(rbDetails lifecycle.ReleaseBundleDetails, params lifecycle.ReleaseBundleQueryParams) error {
	rbService := lifecycle.NewReleaseBundlesService(lcs.config.GetServiceDetails(), lcs.client)
	return rbService.DeleteReleaseBundle(rbDetails, params)
}
