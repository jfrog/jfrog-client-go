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

func New(commonDetails *auth.CommonDetails, config config.Config) (*DistributionServicesManager, error) {
	var err error
	manager := &DistributionServicesManager{config: config}
	manager.client, err = rthttpclient.ArtifactoryClientBuilder().
		SetCertificatesPath(config.GetCertificatesPath()).
		SetInsecureTls(config.IsInsecureTls()).
		SetCommonDetails(commonDetails).
		Build()
	if err != nil {
		return nil, err
	}
	return manager, err
}

func (sm *DistributionServicesManager) SetSigningKey(params services.SetSigningKeyParams) error {
	setSigningKeyService := services.NewSetSigningKeyService(sm.client)
	setSigningKeyService.DistDetails = sm.config.GetCommonDetails()
	return setSigningKeyService.SetSigningKey(params)
}

func (sm *DistributionServicesManager) CreateReleaseBundle(params services.CreateReleaseBundleParams) error {
	createBundleService := services.NewCreateReleseBundleService(sm.client)
	createBundleService.DistDetails = sm.config.GetCommonDetails()
	createBundleService.DryRun = sm.config.IsDryRun()
	return createBundleService.CreateReleaseBundle(params)
}

func (sm *DistributionServicesManager) UpdateReleaseBundle(params services.UpdateReleaseBundleParams) error {
	createBundleService := services.NewUpdateReleaseBundleService(sm.client)
	createBundleService.DistDetails = sm.config.GetCommonDetails()
	createBundleService.DryRun = sm.config.IsDryRun()
	return createBundleService.UpdateReleaseBundle(params)
}

func (sm *DistributionServicesManager) SignReleaseBundle(params services.SignBundleParams) error {
	signBundleService := services.NewSignBundleService(sm.client)
	signBundleService.DistDetails = sm.config.GetCommonDetails()
	return signBundleService.SignReleaseBundle(params)
}

func (sm *DistributionServicesManager) DistributeReleaseBundle(params services.DistributionParams) error {
	distributeBundleService := services.NewDistributeReleaseBundleService(sm.client)
	distributeBundleService.DistDetails = sm.config.GetCommonDetails()
	distributeBundleService.DryRun = sm.config.IsDryRun()
	return distributeBundleService.Distribute(params)
}

func (sm *DistributionServicesManager) DeleteReleaseBundle(params services.DeleteDistributionParams) error {
	deleteBundleService := services.NewDeleteReleaseBundleService(sm.client)
	deleteBundleService.DistDetails = sm.config.GetCommonDetails()
	deleteBundleService.DryRun = sm.config.IsDryRun()
	return deleteBundleService.DeleteDistribution(params)
}

func (sm *DistributionServicesManager) DeleteLocalReleaseBundle(params services.DeleteDistributionParams) error {
	deleteLocalBundleService := services.NewDeleteLocalDistributionService(sm.client)
	deleteLocalBundleService.DistDetails = sm.config.GetCommonDetails()
	deleteLocalBundleService.DryRun = sm.config.IsDryRun()
	return deleteLocalBundleService.DeleteDistribution(params)
}

func (sm *DistributionServicesManager) Client() *rthttpclient.ArtifactoryHttpClient {
	return sm.client
}
