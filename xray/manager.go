package xray

import (
	rthttpclient "github.com/jfrog/jfrog-client-go/artifactory/httpclient"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/jfrog/jfrog-client-go/xray/services"
)

type XrayServicesManager struct {
	client *rthttpclient.ArtifactoryHttpClient
	config config.Config
}

func New(details *auth.ServiceDetails, config config.Config) (*XrayServicesManager, error) {
	err := (*details).InitSsh()
	if err != nil {
		return nil, err
	}
	manager := &XrayServicesManager{config: config}
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

func (sm *XrayServicesManager) Client() *rthttpclient.ArtifactoryHttpClient {
	return sm.client
}

func (sm *XrayServicesManager) GetXrayVersion() (string, error) {
	versionService := services.NewVersionService(sm.client)
	versionService.XrayDetails = sm.config.GetServiceDetails()
	return versionService.GetXrayVersion()
}
