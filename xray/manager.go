package xray

import (
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
func New(config config.Config) (*XrayServicesManager, error) {
	details := config.GetServiceDetails()
	var err error
	manager := &XrayServicesManager{config: config}
	manager.client, err = jfroghttpclient.JfrogClientBuilder().
		SetCertificatesPath(config.GetCertificatesPath()).
		SetInsecureTls(config.IsInsecureTls()).
		SetContext(config.GetContext()).
		SetTimeout(config.GetHttpTimeout()).
		SetClientCertPath(details.GetClientCertPath()).
		SetClientCertKeyPath(details.GetClientCertKeyPath()).
		AppendPreRequestInterceptor(details.RunPreRequestFunctions).
		SetRetries(config.GetHttpRetries()).
		Build()
	return manager, err
}

// Client will return the http client
func (sm *XrayServicesManager) Client() *jfroghttpclient.JfrogHttpClient {
	return sm.client
}

func (sm *XrayServicesManager) Config() config.Config {
	return sm.config
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

// ScanGraph will send Xray the given graph for scan
// Returns a string represents the scan ID.
func (sm *XrayServicesManager) ScanGraph(params services.XrayGraphScanParams) (scanId string, err error) {
	scanService := services.NewScanService(sm.client)
	scanService.XrayDetails = sm.config.GetServiceDetails()
	return scanService.ScanGraph(params)
}

// GetScanGraphResults returns an Xray scan output of the requested graph scan.
// The scanId input should be received from ScanGraph request.
func (sm *XrayServicesManager) GetScanGraphResults(scanID string) (*services.ScanResponse, error) {
	scanService := services.NewScanService(sm.client)
	scanService.XrayDetails = sm.config.GetServiceDetails()
	return scanService.GetScanGraphResults(scanID)
}
