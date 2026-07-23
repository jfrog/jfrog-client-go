package services

import (
	"encoding/json"
	"net/http"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	xscutils "github.com/jfrog/jfrog-client-go/xsc/services/utils"
)

const (
	scanResultsRouteAPIUrl      = "gitinfo/scanResultsUiRoute"
	GetUIRouteAPIMinXrayVersion = "3.76.0"
)

type ScanResultsRouteService struct {
	client          *jfroghttpclient.JfrogHttpClient
	XrayDetails     auth.ServiceDetails
	ScopeProjectKey string
}

func NewScanResultsRouteService(client *jfroghttpclient.JfrogHttpClient) *ScanResultsRouteService {
	return &ScanResultsRouteService{client: client}
}

func (s *ScanResultsRouteService) getScanResultsUIRouteURL() string {
	return utils.AppendScopedProjectKeyParam(utils.AddTrailingSlashIfNeeded(s.XrayDetails.GetUrl())+xscutils.XscInXraySuffix+scanResultsRouteAPIUrl, s.ScopeProjectKey)
}

func (s *ScanResultsRouteService) GetScanResultsUIRoute(gitInfo *XscGitInfoContext) (*ScanResultsUIRouteResponse, error) {
	httpClientsDetails := s.XrayDetails.CreateHttpClientDetails()
	httpClientsDetails.SetContentTypeApplicationJson()

	request := ScanResultsUIRouteRequest{
		GitInfoUrl: gitInfo.Source.GitRepoHttpsCloneUrl,
		BranchName: gitInfo.Source.BranchName,
		CommitHash: gitInfo.Source.CommitHash,
	}
	if gitInfo.PullRequest != nil {
		request.PullRequestId = int64(gitInfo.PullRequest.PullRequestId)
	}
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, errorutils.CheckError(err)
	}
	resp, body, err := s.client.SendPost(s.getScanResultsUIRouteURL(), requestBody, &httpClientsDetails)
	if err != nil {
		return nil, errorutils.CheckError(err)
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return nil, err
	}
	var response ScanResultsUIRouteResponse
	err = json.Unmarshal(body, &response)
	return &response, errorutils.CheckError(err)
}

type ScanResultsUIRouteRequest struct {
	GitInfoUrl    string `json:"git_info_url"`
	BranchName    string `json:"branch_name,omitempty"`
	CommitHash    string `json:"commit_hash,omitempty"`
	PullRequestId int64  `json:"pull_request_id,omitempty"`
}

type ScanResultsUIRouteResponse struct {
	Url string `json:"url"`
	// Metadata
	RepoId   int64  `json:"repo_id"`
	BranchId int64  `json:"branch_id"`
	Path     string `json:"path"`
}
