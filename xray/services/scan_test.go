package services

import (
	"fmt"
	"testing"

	xscUtils "github.com/jfrog/jfrog-client-go/xsc/services/utils"
)

func TestCreateScanGraphQueryParams(t *testing.T) {
	tests := []struct {
		testName      string
		projectKey    string
		repoPath      string
		gitRepoUrl    string
		watches       []string
		scanType      ScanType
		xrayVersion   string
		expectedQuery string
	}{
		{"with_project_key", "p1", "", "", nil, Binary, "0.0.0", fmt.Sprintf("?%s%s&%s%s", projectQueryParam, "p1", scanTypeQueryParam, Binary)},

		{"with_repo_path", "", "r1", "", nil, Binary, "0.0.0", fmt.Sprintf("?%s%s&%s%s", repoPathQueryParam, "r1", scanTypeQueryParam, Binary)},

		{"with_watches", "", "", "", []string{"w1", "w2"}, Binary, "0.0.0", fmt.Sprintf("?%s%s&%s%s&%s%s", watchesQueryParam, "w1", watchesQueryParam, "w2", scanTypeQueryParam, Binary)},

		{"with_empty_watch_string", "", "", "", []string{""}, "", xscUtils.MinXrayVersionXscTransitionToXray, ""},

		{"without_context", "", "", "", nil, Dependency, xscUtils.MinXrayVersionXscTransitionToXray, fmt.Sprintf("?%s%s", scanTypeQueryParam, Dependency)},

		{"without_scan_type", "", "", "", []string{"w1", "w2"}, "", "0.0.0", fmt.Sprintf("?%s%s&%s%s", watchesQueryParam, "w1", watchesQueryParam, "w2")},

		{"with_git_repo_url", "", "", "http://some-url", nil, Dependency, xscUtils.MinXrayVersionXscTransitionToXray, fmt.Sprintf("?%s%s&%s%s", scanTypeQueryParam, Dependency, gitRepoKeyQueryParam, "some-url.git")},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			params := XrayGraphScanParams{
				RepoPath:             test.repoPath,
				Watches:              test.watches,
				ProjectKey:           test.projectKey,
				ScanType:             test.scanType,
				XrayVersion:          test.xrayVersion,
				GitRepoHttpsCloneUrl: test.gitRepoUrl,
			}
			actualQuery := createScanGraphQueryParams(params)
			if actualQuery != test.expectedQuery {
				t.Error(test.testName, "Expecting:", test.expectedQuery, "Got:", actualQuery)
			}
		})
	}
}
