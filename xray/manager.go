package xray

import (
	rthttpclient "github.com/jfrog/jfrog-client-go/artifactory/httpclient"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/jfrog/jfrog-client-go/xray/services"
	"github.com/jfrog/jfrog-client-go/xray/services/utils"
)

// XrayServicesManager defines the http client and general configuration
type XrayServicesManager struct {
	client *rthttpclient.ArtifactoryHttpClient
	config config.Config
}

// New creates a service manager to interact with Xray
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

// Client will return the http client
func (sm *XrayServicesManager) Client() *rthttpclient.ArtifactoryHttpClient {
	return sm.client
}

// GetXrayVersion will return the xray version
func (sm *XrayServicesManager) GetXrayVersion() (string, error) {
	versionService := services.NewVersionService(sm.client)
	versionService.XrayDetails = sm.config.GetServiceDetails()
	return versionService.GetXrayVersion()
}

// CreateXrayWatch will create a new xray watch
func (sm *XrayServicesManager) CreateXrayWatch(params utils.XrayWatchParams) error {
	XrayWatchService := services.NewXrayWatchService(sm.client)
	XrayWatchService.XrayDetails = sm.config.GetServiceDetails()
	return XrayWatchService.Create(params)
}

// GetXrayWatch retrieves the details about an Xray watch by name
// It will error if no watch can be found by that name.
func (sm *XrayServicesManager) GetXrayWatch(watchName string) (*utils.XrayWatchParams, error) {
	XrayWatchService := services.NewXrayWatchService(sm.client)
	XrayWatchService.XrayDetails = sm.config.GetServiceDetails()
	return XrayWatchService.Get(watchName)
}

// UpdateXrayWatch will update an existing Xray watch by name
// It will error if no watch can be found by that name.
func (sm *XrayServicesManager) UpdateXrayWatch(params utils.XrayWatchParams) error {
	XrayWatchService := services.NewXrayWatchService(sm.client)
	XrayWatchService.XrayDetails = sm.config.GetServiceDetails()
	return XrayWatchService.Update(params)
}

// DeleteXrayWatch will delete an existing watch by name
// It will error if no watch can be found by that name.
func (sm *XrayServicesManager) DeleteXrayWatch(watchName string) error {
	XrayWatchService := services.NewXrayWatchService(sm.client)
	XrayWatchService.XrayDetails = sm.config.GetServiceDetails()
	return XrayWatchService.Delete(watchName)
}
