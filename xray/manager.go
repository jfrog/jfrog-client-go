package xray

import (
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

// GetVersion will return the Xray version
func (sm *XrayServicesManager) GetVersion() (string, error) {
	versionService := services.NewVersionService(sm.client)
	versionService.XrayDetails = sm.config.GetServiceDetails()
	return versionService.GetVersion()
}

// CreateWatch will create a new Xray watch
func (sm *XrayServicesManager) CreateWatch(params utils.WatchParams) error {
	watchService := services.NewWatchService(sm.client)
	watchService.XrayDetails = sm.config.GetServiceDetails()
	return watchService.Create(params)
}

// GetWatch retrieves the details about an Xray watch by name
// It will error if no watch can be found by that name.
func (sm *XrayServicesManager) GetWatch(watchName string) (*utils.WatchParams, error) {
	watchService := services.NewWatchService(sm.client)
	watchService.XrayDetails = sm.config.GetServiceDetails()
	return watchService.Get(watchName)
}

// UpdateWatch will update an existing Xray watch by name
// It will error if no watch can be found by that name.
func (sm *XrayServicesManager) UpdateWatch(params utils.WatchParams) error {
	watchService := services.NewWatchService(sm.client)
	watchService.XrayDetails = sm.config.GetServiceDetails()
	return watchService.Update(params)
}

// DeleteWatch will delete an existing watch by name
// It will error if no watch can be found by that name.
func (sm *XrayServicesManager) DeleteWatch(watchName string) error {
	watchService := services.NewWatchService(sm.client)
	watchService.XrayDetails = sm.config.GetServiceDetails()
	return watchService.Delete(watchName)
}

// CreatePolicy will create a new Xray policy
func (sm *XrayServicesManager) CreatePolicy(params utils.PolicyParams) error {
	policyService := services.NewPolicyService(sm.client)
	policyService.XrayDetails = sm.config.GetServiceDetails()
	return policyService.Create(params)
}

// GetPolicy retrieves the details about an Xray policy by name
// It will error if no policy can be found by that name.
func (sm *XrayServicesManager) GetPolicy(policyName string) (*utils.PolicyParams, error) {
	policyService := services.NewPolicyService(sm.client)
	policyService.XrayDetails = sm.config.GetServiceDetails()
	return policyService.Get(policyName)
}

// UpdatePolicy will update an existing Xray policy by name
// It will error if no policy can be found by that name.
func (sm *XrayServicesManager) UpdatePolicy(params utils.PolicyParams) error {
	policyService := services.NewPolicyService(sm.client)
	policyService.XrayDetails = sm.config.GetServiceDetails()
	return policyService.Update(params)
}

// DeletePolicy will delete an existing policy by name
// It will error if no policy can be found by that name.
func (sm *XrayServicesManager) DeletePolicy(policyName string) error {
	policyService := services.NewPolicyService(sm.client)
	policyService.XrayDetails = sm.config.GetServiceDetails()
	return policyService.Delete(policyName)
}

// AddBuildsToIndexing will add builds to Xray indexing configuration
func (sm *XrayServicesManager) AddBuildsToIndexing(buildNames []string) error {
	binMgrService := services.NewBinMgrService(sm.client)
	binMgrService.XrayDetails = sm.config.GetServiceDetails()
	return binMgrService.AddBuildsToIndexing(buildNames)
}
