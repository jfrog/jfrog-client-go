package gpg

import (
	"errors"
	"github.com/jfrog/jfrog-client-go/bintray/auth"
	"github.com/jfrog/jfrog-client-go/bintray/services/utils"
	"github.com/jfrog/jfrog-client-go/bintray/services/versions"
	"github.com/jfrog/jfrog-client-go/httpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
	"path"
)

func NewService(client *httpclient.HttpClient) *GpgService {
	us := &GpgService{client: client}
	return us
}

type GpgService struct {
	client         *httpclient.HttpClient
	BintrayDetails auth.BintrayDetails
}

type Path struct {
	Subject string
	Repo    string
	Package string
	Version string
}

func (gs *GpgService) SignFile(pathDetails *utils.PathDetails, passphrase string) error {
	if gs.BintrayDetails.GetUser() == "" {
		gs.BintrayDetails.SetUser(pathDetails.Subject)
	}
	url := gs.BintrayDetails.GetApiUrl() + path.Join("gpg", pathDetails.Subject, pathDetails.Repo, pathDetails.Path)

	var data string
	if passphrase != "" {
		data = "{ \"passphrase\": \"" + passphrase + "\" }"
	}

	log.Info("GPG signing file...")
	httpClientsDetails := gs.BintrayDetails.CreateHttpClientDetails()
	client, err := httpclient.ClientBuilder().Build()
	if err != nil {
		return err
	}
	resp, body, err := client.SendPost(url, []byte(data), httpClientsDetails)
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

func (gs *GpgService) SignVersion(versionPath *versions.Path, passphrase string) error {
	if gs.BintrayDetails.GetUser() == "" {
		gs.BintrayDetails.SetUser(versionPath.Subject)
	}
	url := gs.BintrayDetails.GetApiUrl() + path.Join("gpg", versionPath.Subject, versionPath.Repo, versionPath.Package, "versions", versionPath.Version)

	var data string
	if passphrase != "" {
		data = "{ \"passphrase\": \"" + passphrase + "\" }"
	}

	log.Info("GPG signing version...")
	httpClientsDetails := gs.BintrayDetails.CreateHttpClientDetails()
	client, err := httpclient.ClientBuilder().Build()
	if err != nil {
		return err
	}
	resp, body, err := client.SendPost(url, []byte(data), httpClientsDetails)
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
