package versions

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path"

	"github.com/jfrog/jfrog-client-go/bintray/auth"
	"github.com/jfrog/jfrog-client-go/httpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
	"path"
)

func NewService(client *httpclient.HttpClient) *VersionService {
	us := &VersionService{client: client}
	return us
}

func NewVersionParams() *Params {
	return &Params{Path: &Path{}}
}

type VersionService struct {
	client         *httpclient.HttpClient
	BintrayDetails auth.BintrayDetails
}

type Path struct {
	Subject string `yaml:"subject,omitempty"`
	Repo    string `yaml:"repo,omitempty"`
	Package string `yaml:"package,omitempty"`
	Version string `yaml:"version"`
}

type Params struct {
	*Path                    `yaml:"path"`
	Desc                     string `yaml:"desc"`
	VcsTag                   string `yaml:"vcstag"`
	Released                 string `yaml:"released"`
	GithubReleaseNotesFile   string `yaml:"githubreleasenotesfile"`
	GithubUseTagReleaseNotes bool   `yaml:"githubusetagreleasenotes"`
}

func (vs *VersionService) Create(params *Params) error {
	log.Info("Creating version...")
	resp, body, err := vs.doCreateVersion(params)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusCreated {
		return errorutils.CheckError(errors.New("Bintray response: " + resp.Status + "\n" + clientutils.IndentJson(body)))
	}
	log.Debug("Bintray response:", resp.Status)
	log.Output(clientutils.IndentJson(body))
	return nil
}

func (vs *VersionService) Update(params *Params) error {
	if vs.BintrayDetails.GetUser() == "" {
		vs.BintrayDetails.SetUser(params.Subject)
	}

	content, err := createVersionContent(params)
	if err != nil {
		return err
	}

	url := vs.BintrayDetails.GetApiUrl() + path.Join("packages/", params.Subject, params.Repo, params.Package, "versions", params.Version)

	log.Info("Updating version...")
	httpClientsDetails := vs.BintrayDetails.CreateHttpClientDetails()
	client, err := httpclient.ClientBuilder().Build()
	if err != nil {
		return err
	}
	resp, body, err := client.SendPatch(url, content, httpClientsDetails)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errorutils.CheckError(errors.New("Bintray response: " + resp.Status + "\n" + clientutils.IndentJson(body)))
	}

	log.Debug("Bintray response:", resp.Status)
	log.Info("Updated version", params.Version+".")
	return nil
}

func (vs *VersionService) Publish(versionPath *Path) error {
	if vs.BintrayDetails.GetUser() == "" {
		vs.BintrayDetails.SetUser(versionPath.Subject)
	}
	url := vs.BintrayDetails.GetApiUrl() + path.Join("content", versionPath.Subject, versionPath.Repo, versionPath.Package, versionPath.Version, "publish")

	log.Info("Publishing version...")
	httpClientsDetails := vs.BintrayDetails.CreateHttpClientDetails()
	client, err := httpclient.ClientBuilder().Build()
	if err != nil {
		return err
	}
	resp, body, err := client.SendPost(url, nil, httpClientsDetails)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errorutils.CheckError(errors.New("Bintray response: " + resp.Status + "\n" + clientutils.IndentJson(body)))
	}

	log.Debug("Bintray response:", resp.Status)
	log.Output(clientutils.IndentJson(body))
	return nil
}

func (vs *VersionService) Delete(versionPath *Path) error {
	if vs.BintrayDetails.GetUser() == "" {
		vs.BintrayDetails.SetUser(versionPath.Subject)
	}
	url := vs.BintrayDetails.GetApiUrl() + path.Join("packages", versionPath.Subject, versionPath.Repo, versionPath.Package, "versions", versionPath.Version)

	log.Info("Deleting version...")
	httpClientsDetails := vs.BintrayDetails.CreateHttpClientDetails()
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
	log.Info("Deleted version", versionPath.Version+".")
	return nil
}

func (vs *VersionService) Show(versionPath *Path) error {
	if vs.BintrayDetails.GetUser() == "" {
		vs.BintrayDetails.SetUser(versionPath.Subject)
	}
	if versionPath.Version == "" {
		versionPath.Version = "_latest"
	}

	url := vs.BintrayDetails.GetApiUrl() + path.Join("packages", versionPath.Subject, versionPath.Repo, versionPath.Package, "versions", versionPath.Version)

	log.Info("Getting version details...")
	httpClientsDetails := vs.BintrayDetails.CreateHttpClientDetails()
	client, err := httpclient.ClientBuilder().Build()
	if err != nil {
		return err
	}
	resp, body, _, _ := client.SendGet(url, true, httpClientsDetails)

	if resp.StatusCode != http.StatusOK {
		return errorutils.CheckError(errors.New("Bintray response: " + resp.Status + "\n" + clientutils.IndentJson(body)))
	}

	log.Debug("Bintray response:", resp.Status)
	log.Output(clientutils.IndentJson(body))
	return nil
}

func (vs *VersionService) IsVersionExists(versionPath *Path) (bool, error) {
	url := vs.BintrayDetails.GetApiUrl() + path.Join("packages", versionPath.Subject, versionPath.Repo, versionPath.Package, "versions", versionPath.Version)
	httpClientsDetails := vs.BintrayDetails.CreateHttpClientDetails()

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

// CalcMetadata -> to schedule metadata calculation https://bintray.com/docs/api/#calc_metadata
func (vs *VersionService) CalcMetadata(versionPath *Path) (bool, error) {
	metaPath := path.Join("calc_metadata", versionPath.Subject, versionPath.Repo, versionPath.Package, versionPath.Version)
	url := vs.BintrayDetails.GetApiUrl() + metaPath
	httpClientsDetails := vs.BintrayDetails.CreateHttpClientDetails()

	client, err := httpclient.ClientBuilder().Build()
	if err != nil {
		return false, err
	}
	resp, body, err := client.SendPost(url, nil, httpClientsDetails)
	if err != nil {
		return false, err
	}
	if resp.StatusCode != http.StatusAccepted {
		return false, errors.New("failed to schedule metadata calculation")
	}
	fmt.Println(clientutils.IndentJson(body))
	log.Info("Metadata calculation scheduled")
	return true, nil
}

func (vs *VersionService) doCreateVersion(params *Params) (*http.Response, []byte, error) {
	if vs.BintrayDetails.GetUser() == "" {
		vs.BintrayDetails.SetUser(params.Subject)
	}

	content, err := createVersionContent(params)
	if err != nil {
		return nil, []byte{}, err
	}
	url := vs.BintrayDetails.GetApiUrl() + path.Join("packages", params.Subject, params.Repo, params.Package, "versions")
	httpClientsDetails := vs.BintrayDetails.CreateHttpClientDetails()
	client, err := httpclient.ClientBuilder().Build()
	if err != nil {
		return nil, []byte{}, err
	}
	return client.SendPost(url, content, httpClientsDetails)
}

func createVersionContent(params *Params) ([]byte, error) {
	Config := contentConfig{
		Name:                     params.Version,
		Desc:                     params.Desc,
		VcsTag:                   params.VcsTag,
		Released:                 params.Released,
		GithubReleaseNotesFile:   params.GithubReleaseNotesFile,
		GithubUseTagReleaseNotes: params.GithubUseTagReleaseNotes,
	}
	requestContent, err := json.Marshal(Config)
	if err != nil {
		return nil, errorutils.CheckError(errors.New("failed to execute request"))
	}
	return requestContent, nil
}

type contentConfig struct {
	Name                     string `json:"name,omitempty"`
	Desc                     string `json:"desc,omitempty"`
	VcsTag                   string `json:"vcs_tag,omitempty"`
	Released                 string `json:"released,omitempty"`
	GithubReleaseNotesFile   string `json:"github_release_notes_file,omitempty"`
	GithubUseTagReleaseNotes bool   `json:"github_use_tag_release_notes,omitempty"`
}
