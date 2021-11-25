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
		expectedQuery string
	}{
		{"with_project_key", "p1", "", nil,
			fmt.Sprintf("?%s%s", projectQueryParam, "p1")},

		{"with_repo_path", "", "r1", nil,
			fmt.Sprintf("?%s%s", repoPathQueryParam, "r1")},

		{"with_watches", "", "", []string{"w1", "w2"},
			fmt.Sprintf("?%s%s&%s%s", watchesQueryParam, "w1", watchesQueryParam, "w2")},

		{"with_empty_watch_string", "", "", []string{""},
			""},

		{"without_scan_type", "", "", []string{"w1", "w2"},
			fmt.Sprintf("?%s%s&%s%s", watchesQueryParam, "w1", watchesQueryParam, "w2")},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			params := XrayGraphScanParams{
				RepoPath:   test.repoPath,
				Watches:    test.watches,
				ProjectKey: test.projectKey,
			}
			actualQuery := createScanGraphQueryParams(params)
			if actualQuery != test.expectedQuery {
				t.Error(test.testName, "Expecting:", test.expectedQuery, "Got:", actualQuery)
			}
		})
	}
}
