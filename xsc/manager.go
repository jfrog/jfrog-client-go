package xsc

import (
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/xsc/services"
)

// XscServicesManager defines the http client and general configuration
type XscServicesManager struct {
	client *jfroghttpclient.JfrogHttpClient
	config config.Config
}

// New creates a service manager to interact with Xsc
func New(config config.Config) (*XscServicesManager, error) {
	details := config.GetServiceDetails()
	var err error
	manager := &XscServicesManager{config: config}
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
func (sm *XscServicesManager) Client() *jfroghttpclient.JfrogHttpClient {
	return sm.client
}

func (sm *XscServicesManager) Config() config.Config {
	return sm.config
}

// GetVersion will return the Xsc version
func (sm *XscServicesManager) GetVersion() (string, error) {
	versionService := services.NewVersionService(sm.client)
	versionService.XscDetails = sm.config.GetServiceDetails()
	return versionService.GetVersion()
}

// TODO: NEED TO IMPLEMENT
//// ScanGraph will send Xsc the given graph for scan
//// Returns a string represents the scan ID.
//func (sm *XscServicesManager) ScanGraph(params services.XscGraphScanParams) (scanId string, err error) {
//	scanService := services.NewScanService(sm.client)
//	scanService.XscDetails = sm.config.GetServiceDetails()
//	return scanService.ScanGraph(params)
//}
//
//// GetScanGraphResults returns an Xsc scan output of the requested graph scan.
//// The scanId input should be received from ScanGraph request.
//func (sm *XscServicesManager) GetScanGraphResults(scanID string, includeVulnerabilities, includeLicenses bool) (*services.ScanResponse, error) {
//	scanService := services.NewScanService(sm.client)
//	scanService.XscDetails = sm.config.GetServiceDetails()
//	return scanService.GetScanGraphResults(scanID, includeVulnerabilities, includeLicenses)
//}
//
//// BuildScan scans a published build-info with Xsc.
//// Returns a string represents the scan ID.
//func (sm *XscServicesManager) BuildScan(params services.XscBuildParams) (*services.BuildScanResponse, error) {
//	buildScanService := services.NewBuildScanService(sm.client)
//	buildScanService.XscDetails = sm.config.GetServiceDetails()
//	err := buildScanService.Scan(params)
//	if err != nil {
//		return nil, err
//	}
//	return buildScanService.GetBuildScanResults(params)
//}
//
//// BuildSummary returns the summary of build scan which had been previously performed.
//func (sm *XscServicesManager) BuildSummary(params services.XscBuildParams) (*services.SummaryResponse, error) {
//	buildSummary := services.NewSummaryService(sm.client)
//	buildSummary.XscDetails = sm.config.GetServiceDetails()
//	return buildSummary.GetBuildSummary(params)
//}
//
//// GenerateVulnerabilitiesReport returns a Xsc report response of the requested report
//func (sm *XscServicesManager) GenerateVulnerabilitiesReport(params services.ReportRequestParams) (resp *services.ReportResponse, err error) {
//	reportService := services.NewReportService(sm.client)
//	reportService.XscDetails = sm.config.GetServiceDetails()
//	return reportService.Vulnerabilities(params)
//}
//
//// ReportDetails returns a Xsc details response for the requested report
//func (sm *XscServicesManager) ReportDetails(reportId string) (details *services.ReportDetails, err error) {
//	reportService := services.NewReportService(sm.client)
//	reportService.XscDetails = sm.config.GetServiceDetails()
//	return reportService.Details(reportId)
//}
//
//// ReportContent returns a Xsc report content response for the requested report
//func (sm *XscServicesManager) ReportContent(params services.ReportContentRequestParams) (content *services.ReportContent, err error) {
//	reportService := services.NewReportService(sm.client)
//	reportService.XscDetails = sm.config.GetServiceDetails()
//	return reportService.Content(params)
//}
//
//// DeleteReport deletes a Xsc report
//func (sm *XscServicesManager) DeleteReport(reportId string) error {
//	reportService := services.NewReportService(sm.client)
//	reportService.XscDetails = sm.config.GetServiceDetails()
//	return reportService.Delete(reportId)
//}
