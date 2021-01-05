package xray

import (
	"net/http"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/xray/services"
	"github.com/jfrog/jfrog-client-go/xray/services/utils"
)

// XrayServicesManager defines the http client and general configuration
type XrayServicesManager struct {
	client *jfroghttpclient.JfrogHttpClient
	config config.Config
}

// New creates a service manager to interact with Xray
func New(details *auth.ServiceDetails, config config.Config) (*XrayServicesManager, error) {
	err := (*details).InitSsh()
	if err != nil {
		return nil, err
	}
	manager := &XrayServicesManager{config: config}
	manager.client, err = jfroghttpclient.JfrogClientBuilder().
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
func (sm *XrayServicesManager) Client() *jfroghttpclient.JfrogHttpClient {
	return sm.client
}

// GetVersion will return the xray version
func (sm *XrayServicesManager) GetVersion() (string, error) {
	versionService := services.NewVersionService(sm.client)
	versionService.XrayDetails = sm.config.GetServiceDetails()
	return versionService.GetVersion()
}

// CreateWatch will create a new xray watch
func (sm *XrayServicesManager) CreateWatch(params utils.WatchParams) (*http.Response, error) {
	WatchService := services.NewWatchService(sm.client)
	WatchService.XrayDetails = sm.config.GetServiceDetails()
	return WatchService.Create(params)
}

// GetWatch retrieves the details about an Xray watch by name
// It will error if no watch can be found by that name.
func (sm *XrayServicesManager) GetWatch(watchName string) (*utils.WatchParams, *http.Response, error) {
	WatchService := services.NewWatchService(sm.client)
	WatchService.XrayDetails = sm.config.GetServiceDetails()
	return WatchService.Get(watchName)
}

// UpdateWatch will update an existing Xray watch by name
// It will error if no watch can be found by that name.
func (sm *XrayServicesManager) UpdateWatch(params utils.WatchParams) (*http.Response, error) {
	WatchService := services.NewWatchService(sm.client)
	WatchService.XrayDetails = sm.config.GetServiceDetails()
	return WatchService.Update(params)
}

// DeleteWatch will delete an existing watch by name
// It will error if no watch can be found by that name.
func (sm *XrayServicesManager) DeleteWatch(watchName string) (*http.Response, error) {
	WatchService := services.NewWatchService(sm.client)
	WatchService.XrayDetails = sm.config.GetServiceDetails()
	return WatchService.Delete(watchName)
}
