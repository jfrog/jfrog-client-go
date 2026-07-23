//go:build itest

package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/xsc/services"
	"github.com/jfrog/jfrog-client-go/xsc/services/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	scanResultsUIRouteAPIPath = "/xray/api/v1/xsc/gitinfo/scanResultsUiRoute"
	testGitRepoURL            = "https://github.com/jfrog/jfrog-client-go.git"
	testBranchName            = "feature/XRAY-132336-add-git-route-api"
	testCommitHash            = "abc123def456"
	testPullRequestId         = 42
	testScanResultsURL        = "https://example.jfrog.io/ui/scans-list/repositories/1/scan-desc"
	testRepoId                = int64(101)
	testBranchId              = int64(202)
	testScanResultsPath       = "/ui/scans-list/repositories/1/scan-desc"
)

func TestXscGetScanResultsUIRoute(t *testing.T) {
	initXscTest(t, "", utils.MinXrayVersionXscTransitionToXray)

	baseGitInfo := func() *services.XscGitInfoContext {
		return &services.XscGitInfoContext{
			Source: services.CommitContext{
				GitRepoHttpsCloneUrl: testGitRepoURL,
				BranchName:           testBranchName,
				CommitHash:           testCommitHash,
			},
		}
	}

	testCases := []struct {
		name         string
		gitInfo      *services.XscGitInfoContext
		serverStatus int
		expectError  bool
		expectPRId   bool
	}{
		{
			name:         "branch only success",
			gitInfo:      baseGitInfo(),
			serverStatus: http.StatusOK,
			expectError:  false,
			expectPRId:   false,
		},
		{
			name: "with pull request success",
			gitInfo: func() *services.XscGitInfoContext {
				gitInfo := baseGitInfo()
				gitInfo.PullRequest = &services.PullRequestContext{PullRequestId: testPullRequestId}
				return gitInfo
			}(),
			serverStatus: http.StatusOK,
			expectError:  false,
			expectPRId:   true,
		},
		{
			name:         "non-200 response",
			gitInfo:      baseGitInfo(),
			serverStatus: http.StatusInternalServerError,
			expectError:  true,
			expectPRId:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockServer, routeService := createXscMockServerForScanResultsUIRoute(t, tc.serverStatus, tc.expectPRId)
			defer mockServer.Close()

			response, err := routeService.GetScanResultsUIRoute(tc.gitInfo)
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, response)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, response)
			assert.Equal(t, testScanResultsURL, response.Url)
			assert.Equal(t, testRepoId, response.RepoId)
			assert.Equal(t, testBranchId, response.BranchId)
			assert.Equal(t, testScanResultsPath, response.Path)
		})
	}
}

func createXscMockServerForScanResultsUIRoute(t *testing.T, statusCode int, expectPRId bool) (mockServer *httptest.Server, routeService *services.ScanResultsRouteService) {
	mockServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI != scanResultsUIRouteAPIPath || r.Method != http.MethodPost {
			assert.Fail(t, "received an unexpected request: "+r.Method+" "+r.RequestURI)
			return
		}

		var reqBody services.ScanResultsUIRouteRequest
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		assert.NoError(t, err, "Invalid JSON request body")
		if err != nil {
			return
		}

		assert.Equal(t, testGitRepoURL, reqBody.GitInfoUrl)
		assert.Equal(t, testBranchName, reqBody.BranchName)
		assert.Equal(t, testCommitHash, reqBody.CommitHash)
		if expectPRId {
			assert.Equal(t, int64(testPullRequestId), reqBody.PullRequestId)
		} else {
			assert.Zero(t, reqBody.PullRequestId)
		}

		w.WriteHeader(statusCode)
		if statusCode != http.StatusOK {
			_, err = w.Write([]byte(`{"error":"internal server error"}`))
			assert.NoError(t, err)
			return
		}

		response := services.ScanResultsUIRouteResponse{
			Url:      testScanResultsURL,
			RepoId:   testRepoId,
			BranchId: testBranchId,
			Path:     testScanResultsPath,
		}
		err = json.NewEncoder(w).Encode(response)
		assert.NoError(t, err)
	}))

	xrayDetails := GetXrayDetails()
	xrayDetails.SetUrl(mockServer.URL + "/xray")
	xrayDetails.SetAccessToken("")

	client, err := jfroghttpclient.JfrogClientBuilder().Build()
	assert.NoError(t, err)

	routeService = services.NewScanResultsRouteService(client)
	routeService.XrayDetails = xrayDetails
	return
}
