package tests

import (
	"bytes"
	"encoding/json"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils/tests/xray"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/xray/services"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

var testXrayReportService *services.ReportService

func TestXrayReport(t *testing.T) {
	initXrayTest(t)
	xrayServerPort := xray.StartXrayMockServer()
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

func reportAll(t *testing.T) {
	request := services.ReportRequestParams{
		Name: "test-report",
		Filters: services.Filter{
			HasRemediation: true,
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
	report, err := testXrayReportService.Vulnerabilities(request)
	assert.NoError(t, err)
	validateResponse(t, xray.VulnerabilityRequestResponse, report)

	var reportId = strconv.Itoa(report.ReportId)
	details, err := testXrayReportService.Details(reportId)
	assert.NoError(t, err)
	validateResponse(t, xray.VulnerabilityReportStatusResponse, details)

	reportReqCont := services.ReportContentRequestParams{
		ReportId:  reportId,
		Direction: "asc",
		PageNum:   0,
		NumRows:   7,
		OrderBy:   "severity",
	}
	content, err := testXrayReportService.Content(reportReqCont)
	assert.NoError(t, err)
	validateResponse(t, xray.VulnerabilityReportDetailsResponse, content)

	err = testXrayReportService.Delete(reportId)
	assert.NoError(t, err)
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
