package services

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	xscutils "github.com/jfrog/jfrog-client-go/xsc/services/utils"
)

const gitScanUIRouteApi = "gitinfo/scan-ui-route"

type GitScanService struct {
	client      *jfroghttpclient.JfrogHttpClient
	XrayDetails auth.ServiceDetails
}

func NewGitScanService(client *jfroghttpclient.JfrogHttpClient) *GitScanService {
	return &GitScanService{client: client}
}

type GetGitScanUIRouteRequest struct {
	CommitHash    string `json:"commit_hash"`
	RepoName      string `json:"repo_name"`
	BranchName    string `json:"branch_name"`
	CdxPath       string `json:"cdx_path"`
	PackageId     string `json:"package_id"`
	PullRequestId int    `json:"pull_request_id,omitempty"`
}

type GetGitScanUIRouteResponse struct {
	Url string `json:"url"`
}

func (gs *GitScanService) getGitScanUIRouteEndpoint() string {
	return utils.AddTrailingSlashIfNeeded(gs.XrayDetails.GetUrl()) + xscutils.XscInXraySuffix + gitScanUIRouteApi
}

// GetGitScanUIRoute returns a full UI URL for viewing git scan results. Only supported for static SCA scans.
func (gs *GitScanService) GetGitScanUIRoute(request GetGitScanUIRouteRequest) (string, error) {
	requestBody, err := json.Marshal(request)
	if err != nil {
		return "", errorutils.CheckError(err)
	}

	httpClientDetails := gs.XrayDetails.CreateHttpClientDetails()
	resp, body, err := gs.client.SendPost(gs.getGitScanUIRouteEndpoint(), requestBody, &httpClientDetails)
	if err != nil {
		return "", err
	}

	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return "", err
	}

	var response GetGitScanUIRouteResponse
	err = errorutils.CheckError(json.Unmarshal(body, &response))
	if err != nil {
		return "", err
	}

	baseUrl := strings.TrimSuffix(gs.XrayDetails.GetUrl(), "/")
	return baseUrl + response.Url, nil
}
