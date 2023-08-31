package manager

import (
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/jfrog/jfrog-client-go/xray/scan"
)

type XscServicesManger struct {
	XrayServicesManager
}

func (xsc *XscServicesManger) IsXscEnabled() (bool, string, error) {
	return xsc.XrayServicesManager.IsXscEnabled()
}

func (xsc *XscServicesManger) SetClient(client *jfroghttpclient.JfrogHttpClient) {
	xsc.XrayServicesManager.SetClient(client)
}

// ScanGraph scans dependency graph with XscGitInfoContext.
// XscGitInfoContext allows linking of scans and other data to the corresponding git repository.
// By passing multi-scan-id in the api calls.
// Returns a string represents the scan ID.
func (xsc *XscServicesManger) ScanGraph(params *scan.XrayGraphScanParams) (scanId string, err error) {
	log.Debug("Scanning graph using XSC service...")
	scanService := scan.NewXscScanService(xsc.client, xsc.config.GetServiceDetails())
	multiScanId, err := scanService.SendScanContext(params.XscGitInfoContext)
	if err != nil {
		// Don't fail the entire scan when failed to send XscGitInfoContext
		log.Warn("failed to send xsc git info context with the following error: ", err.Error())
	}
	if multiScanId != "" {
		params.MultiScanId = multiScanId
	}
	return scanService.ScanGraph(params)
}

// GetScanGraphResults returns an XSC scan output of the requested graph scan.
// The scanId input should be received from ScanGraph request.
func (xsc *XscServicesManger) GetScanGraphResults(scanID string, includeVulnerabilities, includeLicenses bool) (*scan.ScanResponse, error) {
	scanService := scan.NewXscScanService(xsc.client, xsc.config.GetServiceDetails())
	return scanService.GetScanGraphResults(scanID, includeVulnerabilities, includeLicenses)
}
