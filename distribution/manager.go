package distribution

import (
	rthttpclient "github.com/jfrog/jfrog-client-go/artifactory/httpclient"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/jfrog/jfrog-client-go/distribution/services"
)

type DistributionServicesManager struct {
	client *rthttpclient.ArtifactoryHttpClient
	config config.Config
}

func New(details *auth.ServiceDetails, config config.Config) (*DistributionServicesManager, error) {
	err := (*details).InitSsh()
	if err != nil {
		return nil, err
	}
	manager := &DistributionServicesManager{config: config}
	manager.client, err = rthttpclient.ArtifactoryClientBuilder().
		SetCertificatesPath(config.GetCertificatesPath()).
		SetInsecureTls(config.IsInsecureTls()).
		SetServiceDetails(details).
		Build()
	if err != nil {
		return nil, err
	}
	return manager, err
}

func (sm *DistributionServicesManager) SetSigningKey(params services.SetSigningKeyParams) error {
	setSigningKeyService := services.NewSetSigningKeyService(sm.client)
	setSigningKeyService.DistDetails = sm.config.GetServiceDetails()
	return setSigningKeyService.SetSigningKey(params)
}

func (sm *DistributionServicesManager) CreateReleaseBundle(params services.CreateReleaseBundleParams) error {
	createBundleService := services.NewCreateReleseBundleService(sm.client)
	createBundleService.DistDetails = sm.config.GetServiceDetails()
	createBundleService.DryRun = sm.config.IsDryRun()
	return createBundleService.CreateReleaseBundle(params)
}

func (sm *DistributionServicesManager) UpdateReleaseBundle(params services.UpdateReleaseBundleParams) error {
	createBundleService := services.NewUpdateReleaseBundleService(sm.client)
	createBundleService.DistDetails = sm.config.GetServiceDetails()
	createBundleService.DryRun = sm.config.IsDryRun()
	return createBundleService.UpdateReleaseBundle(params)
}

func (sm *DistributionServicesManager) SignReleaseBundle(params services.SignBundleParams) error {
	signBundleService := services.NewSignBundleService(sm.client)
	signBundleService.DistDetails = sm.config.GetServiceDetails()
	return signBundleService.SignReleaseBundle(params)
}

func (sm *DistributionServicesManager) DistributeReleaseBundle(params services.DistributionParams) error {
	distributeBundleService := services.NewDistributeReleaseBundleService(sm.client)
	distributeBundleService.DistDetails = sm.config.GetServiceDetails()
	distributeBundleService.DryRun = sm.config.IsDryRun()
	return distributeBundleService.Distribute(params)
}

func (sm *DistributionServicesManager) DeleteReleaseBundle(params services.DeleteDistributionParams) error {
	deleteBundleService := services.NewDeleteReleaseBundleService(sm.client)
	deleteBundleService.DistDetails = sm.config.GetServiceDetails()
	deleteBundleService.DryRun = sm.config.IsDryRun()
	return deleteBundleService.DeleteDistribution(params)
}

func (sm *DistributionServicesManager) DeleteLocalReleaseBundle(params services.DeleteDistributionParams) error {
	deleteLocalBundleService := services.NewDeleteLocalDistributionService(sm.client)
	deleteLocalBundleService.DistDetails = sm.config.GetServiceDetails()
	deleteLocalBundleService.DryRun = sm.config.IsDryRun()
	return deleteLocalBundleService.DeleteDistribution(params)
}

func (sm *DistributionServicesManager) Client() *rthttpclient.ArtifactoryHttpClient {
	return sm.client
}

func (sm *DistributionServicesManager) GetDistributionVersion() (string, error) {
	versionService := services.NewVersionService(sm.client)
	versionService.DistDetails = sm.config.GetServiceDetails()
	return versionService.GetDistributionVersion()
}
