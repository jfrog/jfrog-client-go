package services

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScanResponseCVSSRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		cve  string
	}{
		{
			name: "all CVSS versions",
			cve:  `{"cve":"CVE-2024-1234","cvss_v2_score":"7.5","cvss_v2_vector":"AV:N/AC:L/Au:N/C:P/I:P/A:P","cvss_v3_score":"9.8","cvss_v3_vector":"CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:H/I:H/A:H","cvss_v4_score":"9.3","cvss_v4_vector":"CVSS:4.0/AV:N/AC:L/AT:N/PR:N/UI:N/VC:H/VI:H/VA:H/SC:N/SI:N/SA:N"}`,
		},
		{
			name: "CVSS v4 only",
			cve:  `{"cve":"CVE-2024-1234","cvss_v4_score":"9.3","cvss_v4_vector":"CVSS:4.0/AV:N/AC:L/AT:N/PR:N/UI:N/VC:H/VI:H/VA:H/SC:N/SI:N/SA:N"}`,
		},
		{
			name: "legacy response without CVSS v4",
			cve:  `{"cve":"CVE-2024-1234","cvss_v2_score":"7.5","cvss_v3_score":"9.8"}`,
		},
		{
			name: "no CVSS data",
			cve:  `{"cve":"CVE-2024-1234"}`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			payload := fmt.Sprintf(`{"violations":[{"cves":[%s]}],"vulnerabilities":[{"cves":[%s]}]}`, test.cve, test.cve)
			for _, response := range []struct {
				name  string
				value interface{}
			}{
				{name: "build scan", value: &BuildScanResponse{}},
				{name: "graph scan", value: &ScanResponse{}},
			} {
				t.Run(response.name, func(t *testing.T) {
					require.NoError(t, json.Unmarshal([]byte(payload), response.value))
					output, err := json.Marshal(response.value)
					require.NoError(t, err)
					assert.JSONEq(t, payload, string(output))
				})
			}
		})
	}
}

func TestCreateScanGraphQueryParams(t *testing.T) {
	tests := []struct {
		testName      string
		params        XrayGraphScanParams
		expectedQuery string
	}{
		{
			testName:      "with_project_key",
			params:        XrayGraphScanParams{ProjectKey: "p1", ScanType: Binary},
			expectedQuery: fmt.Sprintf("?%s%s&%s%s", projectQueryParam, "p1", scanTypeQueryParam, Binary),
		},
		{
			testName:      "with_repo_path",
			params:        XrayGraphScanParams{RepoPath: "r1", ScanType: Binary},
			expectedQuery: fmt.Sprintf("?%s%s&%s%s", repoPathQueryParam, "r1", scanTypeQueryParam, Binary),
		},
		{
			testName:      "with_watches",
			params:        XrayGraphScanParams{Watches: []string{"w1", "w2"}, ScanType: Binary},
			expectedQuery: fmt.Sprintf("?%s%s&%s%s&%s%s", watchesQueryParam, "w1", watchesQueryParam, "w2", scanTypeQueryParam, Binary),
		},
		{
			testName:      "with_empty_watch_string",
			params:        XrayGraphScanParams{Watches: []string{""}},
			expectedQuery: "",
		},
		{
			testName:      "without_context",
			params:        XrayGraphScanParams{ScanType: Dependency, XrayVersion: MinXrayVersionGitRepoKey},
			expectedQuery: fmt.Sprintf("?%s%s", scanTypeQueryParam, Dependency),
		},
		{
			testName:      "without_scan_type",
			params:        XrayGraphScanParams{Watches: []string{"w1", "w2"}},
			expectedQuery: fmt.Sprintf("?%s%s&%s%s", watchesQueryParam, "w1", watchesQueryParam, "w2"),
		},
		{
			testName:      "with_git_repo_url",
			params:        XrayGraphScanParams{GitRepoHttpsCloneUrl: "http://some-url", ScanType: Dependency, XrayVersion: MinXrayVersionGitRepoKey},
			expectedQuery: fmt.Sprintf("?%s%s&%s%s", scanTypeQueryParam, Dependency, gitRepoKeyQueryParam, "some-url.git"),
		},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			actualQuery := createScanGraphQueryParams(test.params)
			if actualQuery != test.expectedQuery {
				t.Error(test.testName, "Expecting:", test.expectedQuery, "Got:", actualQuery)
			}
		})
	}
}
