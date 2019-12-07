package repositories

import (
	"errors"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/jfrog/jfrog-client-go/bintray/auth"
	"github.com/jfrog/jfrog-client-go/httpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

// NewService -> used for repositories package funcs
func NewService(client *httpclient.HttpClient) *RepositoryService {
	us := &RepositoryService{client: client}
	return us
}

// RepositoryService holds HTTP client and bintray auth details
type RepositoryService struct {
	client         *httpclient.HttpClient
	BintrayDetails auth.BintrayDetails
}

// Path is the URL path of repo
type Path struct {
	Subject string
	Repo    string
}

// Config is the equivalent of repo config json
type Config struct {
	Type               string   `yaml:"type,omitempty"`
	IsPrivate          bool     `yaml:"isprivate,omitempty"`
	Desc               string   `yaml:"desc,omitempty"`
	Labels             []string `yaml:"labels,omitempty"`
	GpgSignFiles       bool     `yaml:"gpgsignfiles,omitempty"`
	GpgSignMetadata    bool     `yaml:"gpgsignmetadata,omitempty"`
	GpgUseOwnerKey     bool     `yaml:"gpguseownerkey,omitempty"`
	RepoConfigFilePath string   `yaml:"repoconfigfilepath,omitempty"`
	YumMetadataDepth   int      `yaml:"yum_metadata_depth,omitempty"`
	YumGroupsFile      string   `yaml:"yum_groups_file,omitempty"`
}

// IsRepoExists -> to Check if Repo exists under owner
func (rs *RepositoryService) IsRepoExists(repositoryPath *Path) (bool, error) {
	url := rs.BintrayDetails.GetApiUrl() + path.Join("repos", repositoryPath.Subject, repositoryPath.Repo)
	httpClientsDetails := rs.BintrayDetails.CreateHttpClientDetails()

	client, err := httpclient.ClientBuilder().Build()
	if err != nil {
		return false, err
	}
	resp, _, err := client.SendHead(url, httpClientsDetails)
	if err != nil {
		return false, err
	}
	if resp.StatusCode == http.StatusOK {
		return true, nil
	}
	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	return false, errorutils.CheckError(errors.New("Bintray response: " + resp.Status))
}

// CreateReposIfNeeded -> makes a call to IsRepoExists(), then creates one if needed
func (rs *RepositoryService) CreateReposIfNeeded(repositoryPath *Path, repositoryParams *Config, configPath string) (bool, error) {
	var err error
	var existsOk bool

	existsOk, _ = rs.IsRepoExists(repositoryPath)
	if !existsOk {
		existsOk, err = rs.ExecCreateRepoRest(repositoryPath, repositoryParams, configPath)
		if err != nil {
			log.Error(err)
		}
	}
	return existsOk, err
}

// ExecCreateRepoRest -> creates the repo under owner (subject)
func (rs *RepositoryService) ExecCreateRepoRest(repositoryPath *Path, repositoryParams *Config, repoConfig string) (bool, error) {
	repoName := path.Join("repos", repositoryPath.Subject, repositoryPath.Repo)
	content, err := ioutil.ReadFile(getRepoConfigPath(repoConfig))
	if err != nil {
		return false, err
	}
	url := rs.BintrayDetails.GetApiUrl() + repoName
	httpClientsDetails := rs.BintrayDetails.CreateHttpClientDetails()
	httpClientsDetails.Headers = map[string]string{"Content-Type": "application/json"}
	client, err := httpclient.ClientBuilder().Build()
	if err != nil {
		return false, err
	}
	resp, _, err := client.SendPost(url, content, httpClientsDetails)
	if err != nil {
		return false, err
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return false, errors.New("Fail to create repository. Reason local repository with key: " + repoName + " already exist\n")
	}
	log.Info("Repository", repoName, "created.")
	return true, nil
}

// ExecDeleteRepoRest -> Deletes the repo under owner (subject)
func (rs *RepositoryService) ExecDeleteRepoRest(repositoryPath *Path) error {
	repoName := path.Join("repos", repositoryPath.Subject, repositoryPath.Repo)
	url := rs.BintrayDetails.GetApiUrl() + repoName

	log.Info("Deleting Repo...")
	httpClientsDetails := rs.BintrayDetails.CreateHttpClientDetails()
	client, err := httpclient.ClientBuilder().Build()
	if err != nil {
		return err
	}
	resp, body, err := client.SendDelete(url, nil, httpClientsDetails)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errorutils.CheckError(errors.New("Bintray response: " + resp.Status + "\n" + clientutils.IndentJson(body)))
	}

	log.Debug("Bintray response:", resp.Status)
	log.Info("Deleted Repo", repositoryPath.Repo+".")
	return nil
}
