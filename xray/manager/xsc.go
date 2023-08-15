package manager

import "github.com/jfrog/jfrog-client-go/xray/services"

type XscServicesManger struct {
	XrayServicesManager
}

// ScanGraph will send Xray the given graph for scan
// Returns a string represents the scan ID.
func (xsc *XscServicesManger) ScanGraph(params services.XrayGraphScanParams) (scanId string, err error) {
	scanService := services.NewScanService(xsc.client)
	scanService.XrayDetails = xsc.config.GetServiceDetails()
	return scanService.ScanGraph(params)
}

// GetScanGraphResults returns an Xray scan output of the requested graph scan.
// The scanId input should be received from ScanGraph request.
func (xsc *XscServicesManger) GetScanGraphResults(scanID string, includeVulnerabilities, includeLicenses bool) (*services.ScanResponse, error) {
	scanService := services.NewScanService(xsc.client)
	scanService.XrayDetails = xsc.config.GetServiceDetails()
	return scanService.GetScanGraphResults(scanID, includeVulnerabilities, includeLicenses)
}
