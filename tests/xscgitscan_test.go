//go:build itest

package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/xsc/services"
	xscutils "github.com/jfrog/jfrog-client-go/xsc/services/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetGitScanUIRoute(t *testing.T) {
	initXscTest(t, "", xscutils.MinXrayVersionXscTransitionToXray)

	mockServer, gitScanService := createXscMockServerForGitScan(t)
	defer mockServer.Close()

	request := services.GetGitScanUIRouteRequest{
		CommitHash:    "abc123",
		RepoName:      "test-repo",
		BranchName:    "main",
		CdxPath:       "frogbot/github.com/test/repo/main/commits/source_code.cdx.json",
		PackageId:     "generic://sha256:abc123/source_code.cdx.json",
		PullRequestId: 0,
	}

	url, err := gitScanService.GetGitScanUIRoute(request)
	assert.NoError(t, err)
	assert.Contains(t, url, "/ui/scans-list/git-repos-scans/test-repo/scan-descendants/main")
}

func TestGetGitScanUIRoute_PullRequest(t *testing.T) {
	initXscTest(t, "", xscutils.MinXrayVersionXscTransitionToXray)

	mockServer, gitScanService := createXscMockServerForGitScan(t)
	defer mockServer.Close()

	request := services.GetGitScanUIRouteRequest{
		CommitHash:    "abc123",
		RepoName:      "test-repo",
		BranchName:    "feature-branch",
		CdxPath:       "frogbot/github.com/test/repo/feature-branch/PR/source_code.cdx.json",
		PackageId:     "generic://sha256:abc123/source_code.cdx.json",
		PullRequestId: 42,
	}

	url, err := gitScanService.GetGitScanUIRoute(request)
	assert.NoError(t, err)
	assert.Contains(t, url, "/ui/scans-list/git-repos-scans/test-repo/scan-descendants/main")
	assert.Contains(t, url, "isPullRequest=true")
}

func createXscMockServerForGitScan(t *testing.T) (mockServer *httptest.Server, gitScanService *services.GitScanService) {
	mockServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.RequestURI, "gitinfo/scan-ui-route") && r.Method == http.MethodPost {
			w.WriteHeader(http.StatusOK)
			// Return a mock response with the UI route
			_, err := w.Write([]byte(`{"url":"/ui/scans-list/git-repos-scans/test-repo/scan-descendants/main?repoId=1&repoName=test-repo&branchId=1&isPullRequest=true"}`))
			assert.NoError(t, err)
		} else {
			assert.Fail(t, "received an unexpected request: "+r.RequestURI)
		}
	}))

	xrayDetails := GetXrayDetails()
	xrayDetails.SetUrl(mockServer.URL + "/xray")
	xrayDetails.SetAccessToken("")

	client, err := jfroghttpclient.JfrogClientBuilder().Build()
	assert.NoError(t, err)

	gitScanService = services.NewGitScanService(client)
	gitScanService.XrayDetails = xrayDetails
	return
}
