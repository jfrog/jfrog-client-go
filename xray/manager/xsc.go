package manager

import (
	"github.com/jfrog/jfrog-client-go/xray/scan"
)

type XscServicesManger struct {
	XrayServicesManager
}

// ScanGraph will send XSC the given graph for scan
// Returns a string represents the scan ID.
func (xsc *XscServicesManger) ScanGraph(params scan.XrayGraphScanParams) (scanId string, err error) {
	scanService := scan.NewScanService(xsc.client, xsc.config.GetServiceDetails())
	return scanService.ScanGraph(params)
}

// GetScanGraphResults returns an XSC scan output of the requested graph scan.
// The scanId input should be received from ScanGraph request.
func (xsc *XscServicesManger) GetScanGraphResults(scanID string, includeVulnerabilities, includeLicenses bool) (*scan.ScanResponse, error) {
	scanService := scan.NewScanService(xsc.client, xsc.config.GetServiceDetails())
	return scanService.GetScanGraphResults(scanID, includeVulnerabilities, includeLicenses)
}
