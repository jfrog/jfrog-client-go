package manager

import (
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/xray/scan"
)

type XscServicesManger struct {
	XrayServicesManager
}

func (xsc *XscServicesManger) SetClient(client *jfroghttpclient.JfrogHttpClient) {
	xsc.XrayServicesManager.SetClient(client)
}

// ScanGraph will send XSC the given graph for scan
// Sends ContextDetails before scanning in order to show relevant information about the scan in the platform,
// getting multi-scan-id to pass in the calls.
// Returns a string represents the scan ID.
func (xsc *XscServicesManger) ScanGraph(params scan.XrayGraphScanParams) (scanId string, err error) {
	var multiScanId string
	scanService := scan.NewXscScanService(xsc.client, xsc.config.GetServiceDetails())
	if multiScanId, err = scanService.SendScanContext(params.ContextDetails); err != nil {
		return
	}
	return scanService.ScanGraph(params, multiScanId)
}

// GetScanGraphResults returns an XSC scan output of the requested graph scan.
// The scanId input should be received from ScanGraph request.
func (xsc *XscServicesManger) GetScanGraphResults(scanID string, includeVulnerabilities, includeLicenses bool) (*scan.ScanResponse, error) {
	scanService := scan.NewXscScanService(xsc.client, xsc.config.GetServiceDetails())
	return scanService.GetScanGraphResults(scanID, includeVulnerabilities, includeLicenses)
}
