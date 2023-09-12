package services

import (
	"fmt"
	"testing"
)

func TestCreateScanGraphQueryParams(t *testing.T) {
	tests := []struct {
		testName      string
		projectKey    string
		repoPath      string
		watches       []string
		scanType      ScanType
		expectedQuery string
	}{
		{"with_project_key", "p1", "", nil, Binary,
			fmt.Sprintf("?%s%s&%s%s", projectQueryParam, "p1", scanTypeQueryParam, Binary)},

		{"with_repo_path", "", "r1", nil, Binary,
			fmt.Sprintf("?%s%s&%s%s", repoPathQueryParam, "r1", scanTypeQueryParam, Binary)},

		{"with_watches", "", "", []string{"w1", "w2"}, Binary,
			fmt.Sprintf("?%s%s&%s%s&%s%s", watchesQueryParam, "w1", watchesQueryParam, "w2", scanTypeQueryParam, Binary)},

		{"with_empty_watch_string", "", "", []string{""}, "",
			""},

		{"without_context", "", "", nil, Dependency,
			fmt.Sprintf("?%s%s", scanTypeQueryParam, Dependency)},

		{"without_scan_type", "", "", []string{"w1", "w2"}, "",
			fmt.Sprintf("?%s%s&%s%s", watchesQueryParam, "w1", watchesQueryParam, "w2")},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			params := XrayGraphScanParams{
				RepoPath:   test.repoPath,
				Watches:    test.watches,
				ProjectKey: test.projectKey,
				ScanType:   test.scanType,
			}
			actualQuery := createScanGraphQueryParams(params)
			if actualQuery != test.expectedQuery {
				t.Error(test.testName, "Expecting:", test.expectedQuery, "Got:", actualQuery)
			}
		})
	}
}
