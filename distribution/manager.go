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

func (sm *DistributionServicesManager) CreateReleaseBundle(params services.CreateBundleParams) error {
	createBundleService := services.NewCreateBundleService(sm.client)
	createBundleService.DistDetails = sm.config.GetCommonDetails()
	createBundleService.DryRun = sm.config.IsDryRun()
	return createBundleService.CreateReleaseBundle(params)
}

func (sm *DistributionServicesManager) DistributeReleaseBundle(params services.DistributionParams) error {
	distributeBundleService := services.NewDistributeService(sm.client)
	distributeBundleService.DistDetails = sm.config.GetCommonDetails()
	distributeBundleService.DryRun = sm.config.IsDryRun()
	return distributeBundleService.Distribute(params)
}

func (sm *DistributionServicesManager) Client() *rthttpclient.ArtifactoryHttpClient {
	return sm.client
}
