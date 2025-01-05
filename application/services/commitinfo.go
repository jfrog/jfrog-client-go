package services

import (
	"encoding/json"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"net/http"
)

const (
	commitInfoEndpoint = "application/api/v1/commits"
)

type CommitInfoService struct {
	client         *jfroghttpclient.JfrogHttpClient
	serviceDetails auth.ServiceDetails
	config         config.Config
}

func NewCommitInfoService(client *jfroghttpclient.JfrogHttpClient, serviceDetails auth.ServiceDetails) *CommitInfoService {
	return &CommitInfoService{client: client, serviceDetails: serviceDetails}
}

func (c *CommitInfoService) sendPostRequest(requestContent []byte) (resp *http.Response, body []byte, err error) {
	commitInfoUrl := c.serviceDetails.GetUrl() + commitInfoEndpoint
	clientDetails := c.serviceDetails.CreateHttpClientDetails()
	resp, body, err = c.client.SendPost(commitInfoUrl, requestContent, &clientDetails)
	return
}

func (c *CommitInfoService) AddCommitInfo(commitInfo CreateApplicationCommitInfo) error {
	requestContent, err := json.Marshal(commitInfo)
	if err != nil {
		return errorutils.CheckError(err)
	}
	resp, body, err := c.sendPostRequest(requestContent)
	if err != nil {
		return err
	}
	// If the commit info already exists, the response code is 409 (conflict), this scenario is not considered an error.
	if err = errorutils.CheckResponseStatus(resp, http.StatusCreated, http.StatusConflict); err != nil {
		return errorutils.CheckError(errorutils.GenerateResponseError(resp.Status, utils.IndentJson(body)))
	}
	return nil
}

type CreateApplicationCommitInfo struct {
	GitRepoUrl     string `json:"vcs_url"`
	CommitHash     string `json:"commit_hash"`
	ParentHash     string `json:"parent_hash"`
	Branch         string `json:"branch"`
	AuthorEmail    string `json:"author_email"`
	AuthorName     string `json:"author_name"`
	AuthorDate     int64  `json:"author_date"`
	CommitterEmail string `json:"committer_email"`
	CommitterName  string `json:"committer_name"`
	CommitterDate  int64  `json:"committer_date"`
	MessageSubject string `json:"message_subject"`
	MessageBody    string `json:"message_body"`
	ChangedFiles   []byte `json:"changed_files"`
}
