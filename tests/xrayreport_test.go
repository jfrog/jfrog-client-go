//go:build itest

package tests

import (
	"bytes"
	"encoding/json"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils/tests/xray"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/xray/services"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

var testXrayReportService *services.ReportService

func TestXrayReport(t *testing.T) {
	initXrayTest(t)
	xrayServerPort := xray.StartXrayMockServer(t)
	xrayDetails := GetXrayDetails()
	client, err := jfroghttpclient.JfrogClientBuilder().
		SetClientCertPath(xrayDetails.GetClientCertPath()).
		SetClientCertKeyPath(xrayDetails.GetClientCertKeyPath()).
		AppendPreRequestInterceptor(xrayDetails.RunPreRequestFunctions).
		Build()

	assert.NoError(t, err)

	testXrayReportService = services.NewReportService(client)
	testXrayReportService.XrayDetails = xrayDetails
	testXrayReportService.XrayDetails.SetUrl("http://localhost:" + strconv.Itoa(xrayServerPort) + "/")

	t.Run("reportAll", reportAll)
}

var vulnerabilitiesReportRequestParams = services.VulnerabilitiesReportRequestParams{
	Name: "test-report",
	Filters: services.VulnerabilitiesFilter{
		HasRemediation: utils.Pointer(true),
		Severity:       []string{"high"},
	},
	Resources: services.Resource{
		Repositories: []services.Repository{
			{
				Name: "dummy-repo",
			},
		},
	},
}
var licensesReportRequestParams = services.LicensesReportRequestParams{
	Name: "test-report",
	Filters: services.LicensesFilter{
		LicensePatterns: []string{"*"},
	},
	Resources: services.Resource{
		Repositories: []services.Repository{
			{
				Name: "dummy-repo",
			},
		},
	},
}
var reportTypes = []string{
	xray.VulnerabilitiesEndpoint,
	xray.LicensesEndpoint,
}

func reportAll(t *testing.T) {
	for _, ep := range reportTypes {
		var report *services.ReportResponse
		var err error
		if ep == xray.VulnerabilitiesEndpoint {
			report, err = testXrayReportService.Vulnerabilities(vulnerabilitiesReportRequestParams)
		} else if ep == xray.LicensesEndpoint {
			report, err = testXrayReportService.Licenses(licensesReportRequestParams)
		}
		assert.NoError(t, err)
		validateResponse(t, xray.MapResponse[xray.MapReportIdEndpoint[report.ReportId]]["XrayReportRequest"], report)

		var reportId = strconv.Itoa(report.ReportId)
		details, err := testXrayReportService.Details(reportId)
		assert.NoError(t, err)
		validateResponse(t, xray.MapResponse[xray.MapReportIdEndpoint[report.ReportId]]["ReportStatus"], details)

		reportReqCont := services.ReportContentRequestParams{
			ReportType: ep,
			ReportId:   reportId,
			Direction:  "asc",
			PageNum:    0,
			NumRows:    7,
		}
		if ep == xray.VulnerabilitiesEndpoint {
			reportReqCont.OrderBy = "severity"
		} else if ep == xray.LicensesEndpoint {
			reportReqCont.OrderBy = "license"
		}
		content, err := testXrayReportService.Content(reportReqCont)
		assert.NoError(t, err)
		validateResponse(t, xray.MapResponse[ep]["ReportDetails"], content)

		err = testXrayReportService.Delete(reportId)
		assert.NoError(t, err)
	}
}

func validateResponse(t *testing.T, expects string, payload interface{}) {
	compactExpects := new(bytes.Buffer)
	err := json.Compact(compactExpects, []byte(expects))
	if err != nil {
		t.Error(err)
	}

	actualPayloadJsonBytes, err := json.Marshal(payload)
	if err != nil {
		t.Error(err)
	}

	compactActualPayload := new(bytes.Buffer)
	err = json.Compact(compactActualPayload, actualPayloadJsonBytes)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, compactExpects.Len(), compactActualPayload.Len())
	assert.Equal(t, compactActualPayload.String(), string(actualPayloadJsonBytes))
}
