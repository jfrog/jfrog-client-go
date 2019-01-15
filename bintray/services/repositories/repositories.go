package repositories

import (
	"errors"
	"github.com/jfrog/jfrog-client-go/bintray/auth"
	"github.com/jfrog/jfrog-client-go/httpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"net/http"
	"path"
)

func NewService(client *httpclient.HttpClient) *RepositoryService {
	us := &RepositoryService{client: client}
	return us
}

type RepositoryService struct {
	client         *httpclient.HttpClient
	BintrayDetails auth.BintrayDetails
}

type Path struct {
	Subject string
	Repo    string
}

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
