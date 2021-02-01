package entitlements

import (
	"encoding/json"
	"errors"
	"github.com/jfrog/jfrog-client-go/bintray/auth"
	"github.com/jfrog/jfrog-client-go/bintray/services/versions"
	"github.com/jfrog/jfrog-client-go/http/httpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
	"strings"
)

func NewService(client *httpclient.HttpClient) *EntitlementsService {
	us := &EntitlementsService{client: client}
	return us
}

func NewEntitlementsParams() *Params {
	return &Params{VersionPath: &versions.Path{}}
}

type EntitlementsService struct {
	client         *httpclient.HttpClient
	BintrayDetails auth.BintrayDetails
}

type Params struct {
	VersionPath *versions.Path
	Id          string
	Path        string
	Access      string
	Keys        string
}

func (es *EntitlementsService) ShowAll(path *versions.Path) error {
	url := BuildEntitlementsUrl(es.BintrayDetails, path)
	if es.BintrayDetails.GetUser() == "" {
		es.BintrayDetails.SetUser(path.Subject)
	}
	httpClientsDetails := es.BintrayDetails.CreateHttpClientDetails()
	log.Info("Getting entitlements...")
	client, err := httpclient.ClientBuilder().Build()
	if err != nil {
		return err
	}
	resp, body, _, err := client.SendGet(url, true, httpClientsDetails)
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

func (es *EntitlementsService) Show(id string, path *versions.Path) error {
	url := BuildEntitlementUrl(es.BintrayDetails, path, id)
	if es.BintrayDetails.GetUser() == "" {
		es.BintrayDetails.SetUser(path.Subject)
	}
	httpClientsDetails := es.BintrayDetails.CreateHttpClientDetails()
	log.Info("Getting entitlement...")
	client, err := httpclient.ClientBuilder().Build()
	if err != nil {
		return err
	}
	resp, body, _, err := client.SendGet(url, true, httpClientsDetails)
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

func (es *EntitlementsService) Create(params *Params) error {
	var path = BuildEntitlementsUrl(es.BintrayDetails, params.VersionPath)

	if es.BintrayDetails.GetUser() == "" {
		es.BintrayDetails.SetUser(params.VersionPath.Subject)
	}
	content, err := createEntitlementContent(params)
	if err != nil {
		return err
	}

	httpClientsDetails := es.BintrayDetails.CreateHttpClientDetails()
	log.Info("Creating entitlement...")
	client, err := httpclient.ClientBuilder().Build()
	if err != nil {
		return err
	}
	resp, body, err := client.SendPost(path, content, httpClientsDetails)
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

func (es *EntitlementsService) Delete(id string, path *versions.Path) error {
	url := BuildEntitlementUrl(es.BintrayDetails, path, id)
	if es.BintrayDetails.GetUser() == "" {
		es.BintrayDetails.SetUser(path.Subject)
	}

	httpClientsDetails := es.BintrayDetails.CreateHttpClientDetails()
	log.Info("Deleting entitlement...")
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
	log.Output(clientutils.IndentJson(body))
	return nil
}

func (es *EntitlementsService) Update(params *Params) error {
	path := BuildEntitlementUrl(es.BintrayDetails, params.VersionPath, params.Id)
	if es.BintrayDetails.GetUser() == "" {
		es.BintrayDetails.SetUser(params.VersionPath.Subject)
	}
	content, err := createEntitlementContent(params)
	if err != nil {
		return err
	}

	httpClientsDetails := es.BintrayDetails.CreateHttpClientDetails()
	log.Info("Updating entitlement...")
	client, err := httpclient.ClientBuilder().Build()
	if err != nil {
		return err
	}
	resp, body, err := client.SendPatch(path, content, httpClientsDetails)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errorutils.CheckError(errors.New("Bintray response: " + resp.Status + "\n" + clientutils.IndentJson(body)))
	}

	log.Debug("Bintray response:", resp.Status)
	log.Output(clientutils.IndentJson(body))
	return err
}

func createEntitlementContent(params *Params) ([]byte, error) {
	var downloadKeys []string
	if params.Keys != "" {
		downloadKeys = strings.Split(params.Keys, ",")
	}
	Config := contentConfig{
		Access:       params.Access,
		DownloadKeys: downloadKeys,
		Path:         params.Path,
	}
	requestContent, err := json.Marshal(Config)
	if err != nil {
		return nil, errorutils.CheckError(errors.New("Failed to execute request."))
	}
	return requestContent, nil
}

type contentConfig struct {
	Access       string   `json:"access,omitempty"`
	DownloadKeys []string `json:"download_keys,omitempty"`
	Path         string   `json:"path,omitempty"`
}
