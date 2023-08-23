package manager

import (
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/xray/scan"
	"github.com/jfrog/jfrog-client-go/xray/services"
	"github.com/jfrog/jfrog-client-go/xray/services/utils"
)

// SecurityServiceManager manages JFrog Xray service operations (Xray backend or XSC).
type SecurityServiceManager interface {
	// Attributes
	Client() *jfroghttpclient.JfrogHttpClient
	SetClient(client *jfroghttpclient.JfrogHttpClient)
	Config() config.Config
	GetVersion() (string, error)

	// Watches
	CreateWatch(params utils.WatchParams) error
	GetWatch(watchName string) (*utils.WatchParams, error)
	UpdateWatch(params utils.WatchParams) error
	DeleteWatch(watchName string) error
	// Policies
	CreatePolicy(params utils.PolicyParams) error
	GetPolicy(policyName string) (*utils.PolicyParams, error)
	UpdatePolicy(params utils.PolicyParams) error
	DeletePolicy(policyName string) error
	// Scan
	ScanGraph(params *scan.XrayGraphScanParams) (scanId string, err error)
	GetScanGraphResults(scanID string, includeVulnerabilities, includeLicenses bool) (*scan.ScanResponse, error)
	BuildScan(params services.XrayBuildParams, includeVulnerabilities bool) (scanResponse *services.BuildScanResponse, noFailBuildPolicy bool, err error)
	// Report
	GenerateVulnerabilitiesReport(params services.ReportRequestParams) (resp *services.ReportResponse, err error)
	ReportDetails(reportId string) (details *services.ReportDetails, err error)
	ReportContent(params services.ReportContentRequestParams) (content *services.ReportContent, err error)
	DeleteReport(reportId string) error
	// Utilities
	AddBuildsToIndexing(buildNames []string) error
	ArtifactSummary(params services.ArtifactSummaryParams) (*services.ArtifactSummaryResponse, error)
	IsEntitled(featureId string) (bool, error)
	IsXscEnabled() (bool, string, error)
}

// New creates a service manager to interact with Xray
// When XSC is enabled returns XscServicesManger.
func New(config config.Config) (manager SecurityServiceManager, err error) {
	details := config.GetServiceDetails()
	if details.GetXscVersion() != "" {
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
