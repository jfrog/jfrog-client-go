package services

import (
	"encoding/json"
	"errors"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
	"strconv"
)

type SourcesService struct {
	client *jfroghttpclient.JfrogHttpClient
	auth.ServiceDetails
}

func NewSourcesService(client *jfroghttpclient.JfrogHttpClient) *SourcesService {
	return &SourcesService{client: client}
}

const (
	DefaultPipelinesFileFilter = "pipelines.yml"
	SourcesRestApi             = "api/v1/pipelinesources/"
)

func (ss *SourcesService) AddPipelineSource(projectIntegrationId int, repositoryFullName, branch, fileFilter string) (id int, err error) {
	source := Source{
		ProjectId:            defaultProjectId,
		ProjectIntegrationId: projectIntegrationId,
		RepositoryFullName:   repositoryFullName,
		Branch:               branch,
		FileFilter:           fileFilter,
	}
	return ss.addSource(source)
}

func (ss *SourcesService) addSource(source Source) (id int, err error) {
	log.Debug("Adding Pipeline Source...")
	content, err := json.Marshal(source)
	if err != nil {
		return -1, errorutils.CheckError(err)
	}
	httpDetails := ss.ServiceDetails.CreateHttpClientDetails()
	headers := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}
	utils.MergeMaps(httpDetails.Headers, headers)
	httpDetails.Headers = headers

	resp, body, err := ss.client.SendPost(ss.ServiceDetails.GetUrl()+SourcesRestApi, content, &httpDetails)
	if err != nil {
		return -1, err
	}
	if resp.StatusCode != http.StatusOK {
		return -1, errorutils.CheckError(errors.New("Pipelines response: " + resp.Status + "\n" + utils.IndentJson(body)))
	}

	created := &Source{}
	err = json.Unmarshal(body, created)
	if err != nil {
		return -1, err
	}
	return created.Id, nil
}

func (ss *SourcesService) GetSource(sourceId int) (*Source, error) {
	httpDetails := ss.ServiceDetails.CreateHttpClientDetails()
	url := ss.ServiceDetails.GetUrl() + SourcesRestApi + strconv.Itoa(sourceId)
	resp, body, _, err := ss.client.SendGet(url, true, &httpDetails)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errorutils.CheckError(errors.New("Pipelines response: " + resp.Status + "\n" + utils.IndentJson(body)))
	}
	source := &Source{}
	err = json.Unmarshal(body, source)
	if err != nil {
		return nil, err
	}
	return source, nil
}

func (ss *SourcesService) DeleteSource(sourceId int) error {
	httpDetails := ss.ServiceDetails.CreateHttpClientDetails()
	resp, body, err := ss.client.SendDelete(ss.ServiceDetails.GetUrl()+SourcesRestApi+strconv.Itoa(sourceId), nil, &httpDetails)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errorutils.CheckError(errors.New("Pipelines response: " + resp.Status + "\n" + utils.IndentJson(body)))
	}
	return nil
}

type Source struct {
	ProjectId            int    `json:"projectId,omitempty"`
	ProjectIntegrationId int    `json:"projectIntegrationId,omitempty"`
	RepositoryFullName   string `json:"repositoryFullName,omitempty"`
	Branch               string `json:"branch,omitempty"`
	FileFilter           string `json:"fileFilter,omitempty"`
	// For multibranch pipelines only:
	IsMultiBranch        bool   `json:"isMultiBranch,omitempty"`
	BranchExcludePattern string `json:"branchExcludePattern,omitempty"`
	BranchIncludePattern string `json:"branchIncludePattern,omitempty"`

	Id int `json:"id,omitempty"`
}
