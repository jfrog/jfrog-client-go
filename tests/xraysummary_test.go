package tests

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"strconv"
	"strings"
	"testing"

	"github.com/jfrog/jfrog-client-go/artifactory/services/utils/tests/xray"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/xray/services"
)

var testsXraySummaryService *services.SummaryService

func TestNewXraySummaryService(t *testing.T) {
	initXrayTest(t)
	xrayServerPort := xray.StartXrayMockServer(t)
	xrayDetails := GetXrayDetails()
	client, err := jfroghttpclient.JfrogClientBuilder().
		SetClientCertPath(xrayDetails.GetClientCertPath()).
		SetClientCertKeyPath(xrayDetails.GetClientCertKeyPath()).
		AppendPreRequestInterceptor(xrayDetails.RunPreRequestFunctions).
		Build()
	if err != nil {
		t.Error(err)
	}

	testsXraySummaryService = services.NewSummaryService(client)
	testsXraySummaryService.XrayDetails = xrayDetails
	testsXraySummaryService.XrayDetails.SetUrl("http://localhost:" + strconv.Itoa(xrayServerPort) + "/")

	// Run tests
	tests := []struct {
		name      string
		checksums []string
		paths     []string
		expected  string
	}{
		{name: "getVulnerableArtifactSummary", checksums: []string{"0000"}, paths: nil, expected: xray.VulnerableXraySummaryArtifactResponse},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			artifactSummary(t, test.checksums, test.paths, test.expected)
		})
	}
}

func artifactSummary(t *testing.T, checksums []string, paths []string, expected string) {
	params := services.ArtifactSummaryParams{
		Checksums: checksums,
		Paths:     paths,
	}
	result, err := testsXraySummaryService.GetArtifactSummary(params)
	if err != nil {
		t.Error(err)
	}

	resultString, err := json.Marshal(result)
	if err != nil {
		t.Error(err)
	}

	buf := bytes.NewBuffer([]byte{})
	err = json.Compact(buf, []byte(expected))
	assert.NoError(t, err)
	expected = buf.String()

	expected = strings.ReplaceAll(expected, "\n", "")
	actual := strings.ReplaceAll(string(resultString), "\n", "")
	if actual != expected {
		t.Error("\nExpected:\n", expected, "\n\nGot:\n", actual)
	}
}
