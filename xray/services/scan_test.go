package services

import (
	"fmt"
	"testing"
)

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
