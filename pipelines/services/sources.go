package services

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type SourcesService struct {
	client *jfroghttpclient.JfrogHttpClient
	auth.ServiceDetails
}

func NewSourcesService(client *jfroghttpclient.JfrogHttpClient) *SourcesService {
	return &SourcesService{client: client}
}

const (
	DefaultPipelinesFileFilter        = "pipelines.yml"
	SourcesRestApi                    = "api/v1/pipelinesources/"
	sourceAlreadyExistsResponseString = "source already exists"
)

func (ss *SourcesService) AddSource(projectIntegrationId int, repositoryFullName, branch, fileFilter string) (id int, err error) {
	source := Source{
		ProjectId:            defaultProjectId,
		ProjectIntegrationId: projectIntegrationId,
		RepositoryFullName:   repositoryFullName,
		Branch:               branch,
		FileFilter:           fileFilter,
	}
	return ss.doAddSource(source)
}

func (ss *SourcesService) doAddSource(source Source) (id int, err error) {
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
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		if resp.StatusCode == http.StatusNotFound && strings.Contains(string(body), sourceAlreadyExistsResponseString) {
			return -1, &SourceAlreadyExistsError{InnerError: err}
		}
		return -1, err
	}

	created := &Source{}
	err = json.Unmarshal(body, created)
	return created.Id, errorutils.CheckError(err)
}

func (ss *SourcesService) GetSource(sourceId int) (*Source, error) {
	httpDetails := ss.ServiceDetails.CreateHttpClientDetails()
	url := ss.ServiceDetails.GetUrl() + SourcesRestApi + strconv.Itoa(sourceId)
	resp, body, _, err := ss.client.SendGet(url, true, &httpDetails)
	if err != nil {
		return nil, err
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return nil, err
	}
	source := &Source{}
	err = json.Unmarshal(body, source)
	return source, errorutils.CheckError(err)
}

func (ss *SourcesService) DeleteSource(sourceId int) error {
	httpDetails := ss.ServiceDetails.CreateHttpClientDetails()
	resp, body, err := ss.client.SendDelete(ss.ServiceDetails.GetUrl()+SourcesRestApi+strconv.Itoa(sourceId), nil, &httpDetails)
	if err != nil {
		return err
	}
	return errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK)
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

type SourceAlreadyExistsError struct {
	InnerError error
}

func (*SourceAlreadyExistsError) Error() string {
	return "Pipelines: Pipeline Source already exists."
}
