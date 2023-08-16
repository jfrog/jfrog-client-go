package manager

import (
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/xray/scan"
	"github.com/jfrog/jfrog-client-go/xray/services"
	"github.com/jfrog/jfrog-client-go/xray/services/utils"
)

// SecurityServiceManager holds operations to Xray ( regrading if for Xray backend or XSC )
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
	ScanGraph(params scan.XrayGraphScanParams) (scanId string, err error)
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
}
