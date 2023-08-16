package manager

import (
	"encoding/json"
	"fmt"
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/xray/scan"
	"github.com/jfrog/jfrog-client-go/xray/services"
	"github.com/jfrog/jfrog-client-go/xray/services/utils"
	"net/http"
)

// XrayServicesManager defines the http client and general configuration
type XrayServicesManager struct {
	client *jfroghttpclient.JfrogHttpClient
	config config.Config
}

// New creates a service manager to interact with Xray
func New(config config.Config) (manager SecurityServiceManager, err error) {
	details := config.GetServiceDetails()
	if details.GetXscUrl() != "" {
		manager = &XscServicesManger{XrayServicesManager{config: config}}
	} else {
		manager = &XrayServicesManager{config: config}
	}

	client, err := jfroghttpclient.JfrogClientBuilder().
		SetCertificatesPath(config.GetCertificatesPath()).
		SetInsecureTls(config.IsInsecureTls()).
		SetContext(config.GetContext()).
		SetTimeout(config.GetHttpTimeout()).
		SetClientCertPath(details.GetClientCertPath()).
		SetClientCertKeyPath(details.GetClientCertKeyPath()).
		AppendPreRequestInterceptor(details.RunPreRequestFunctions).
		SetRetries(config.GetHttpRetries()).
		SetRetryWaitMilliSecs(config.GetHttpRetryWaitMilliSecs()).
		Build()

	if err != nil {
		return
	}
	manager.SetClient(client)
	return manager, err
}

// Client will return the http client
func (sm *XrayServicesManager) Client() *jfroghttpclient.JfrogHttpClient {
	return sm.client
}

func (sm *XrayServicesManager) SetClient(client *jfroghttpclient.JfrogHttpClient) {
	sm.client = client
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
func (sm *XrayServicesManager) ScanGraph(params scan.XrayGraphScanParams) (scanId string, err error) {
	scanService := scan.NewScanService(sm.client)
	scanService.XrayDetails = sm.config.GetServiceDetails()
	return scanService.ScanGraph(params)
}

// GetScanGraphResults returns an Xray scan output of the requested graph scan.
// The scanId input should be received from ScanGraph request.
func (sm *XrayServicesManager) GetScanGraphResults(scanID string, includeVulnerabilities, includeLicenses bool) (*scan.ScanResponse, error) {
	scanService := scan.NewScanService(sm.client)
	scanService.XrayDetails = sm.config.GetServiceDetails()
	return scanService.GetScanGraphResults(scanID, includeVulnerabilities, includeLicenses)
}

// BuildScan scans a published build-info with Xray.
// 'scanResponse' - Xray scan output of the requested build scan.
// 'noFailBuildPolicy' - Indicates that the Xray API returned a "No Xray Fail build...." error
func (sm *XrayServicesManager) BuildScan(params services.XrayBuildParams, includeVulnerabilities bool) (scanResponse *services.BuildScanResponse, noFailBuildPolicy bool, err error) {
	buildScanService := services.NewBuildScanService(sm.client)
	buildScanService.XrayDetails = sm.config.GetServiceDetails()
	return buildScanService.ScanBuild(params, includeVulnerabilities)
}

// GenerateVulnerabilitiesReport returns a Xray report response of the requested report
func (sm *XrayServicesManager) GenerateVulnerabilitiesReport(params services.ReportRequestParams) (resp *services.ReportResponse, err error) {
	reportService := services.NewReportService(sm.client)
	reportService.XrayDetails = sm.config.GetServiceDetails()
	return reportService.Vulnerabilities(params)
}

// ReportDetails returns a Xray details response for the requested report
func (sm *XrayServicesManager) ReportDetails(reportId string) (details *services.ReportDetails, err error) {
	reportService := services.NewReportService(sm.client)
	reportService.XrayDetails = sm.config.GetServiceDetails()
	return reportService.Details(reportId)
}

// ReportContent returns a Xray report content response for the requested report
func (sm *XrayServicesManager) ReportContent(params services.ReportContentRequestParams) (content *services.ReportContent, err error) {
	reportService := services.NewReportService(sm.client)
	reportService.XrayDetails = sm.config.GetServiceDetails()
	return reportService.Content(params)
}

// DeleteReport deletes a Xray report
func (sm *XrayServicesManager) DeleteReport(reportId string) error {
	reportService := services.NewReportService(sm.client)
	reportService.XrayDetails = sm.config.GetServiceDetails()
	return reportService.Delete(reportId)
}

// ArtifactSummary returns Xray artifact summaries for the requested checksums and/or paths
func (sm *XrayServicesManager) ArtifactSummary(params services.ArtifactSummaryParams) (*services.ArtifactSummaryResponse, error) {
	summaryService := services.NewSummaryService(sm.client)
	summaryService.XrayDetails = sm.config.GetServiceDetails()
	return summaryService.GetArtifactSummary(params)
}

// IsEntitled returns true if the user is entitled for the requested feature ID
func (sm *XrayServicesManager) IsEntitled(featureId string) (bool, error) {
	entitlementsService := services.NewEntitlementsService(sm.client)
	entitlementsService.XrayDetails = sm.config.GetServiceDetails()
	return entitlementsService.IsEntitled(featureId)
}

// IsXscEnabled Trying to get XSC version, if route is not available, user is not entitled for XSC.
func (sm *XrayServicesManager) IsXscEnabled() (xscEntitled bool, xsxVersion string) {
	httpDetails := sm.config.GetServiceDetails().CreateHttpClientDetails()
	serverDetails := sm.config.GetServiceDetails()

	resp, body, _, err := sm.client.SendGet(serverDetails.GetXscUrl()+scan.XscVersionAPI, true, &httpDetails)
	if err != nil {
		return
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		err = fmt.Errorf("failed while attempting to get XSC entitlements response with error:%s", err.Error())
		return
	}
	versionResponse := scan.XscVersionResponse{}
	if err = json.Unmarshal(body, &versionResponse); err != nil {
		err = errorutils.CheckErrorf("couldn't parse Xray server response: " + err.Error())
		return
	}
	if versionResponse.Version != "" {
		return true, versionResponse.Version
	}
	return
}
