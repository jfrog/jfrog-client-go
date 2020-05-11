package mavensync

import (
	"encoding/json"
	"errors"
	"net/http"
	"path"

	"github.com/jfrog/jfrog-client-go/bintray/auth"
	"github.com/jfrog/jfrog-client-go/bintray/services/versions"
	"github.com/jfrog/jfrog-client-go/httpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

func NewService(client *httpclient.HttpClient) *MavenCentralSyncService {
	return &MavenCentralSyncService{client: client}
}

type MavenCentralSyncService struct {
	client         *httpclient.HttpClient
	BintrayDetails auth.BintrayDetails
}

type Params struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Close    string `json:"close,omitempty"`
}

func NewParams(username string, password string, dontClose bool) *Params {
	req := &Params{
		Username: username,
		Password: password,
		Close:    "1",
	}

	if dontClose {
		req.Close = "0"
	}

	return req
}

func (mcss *MavenCentralSyncService) Sync(p *Params, path *versions.Path) error {
	url, err := buildSyncURL(mcss.BintrayDetails, path)
	if err != nil {
		return err
	}

	if mcss.BintrayDetails.GetUser() == "" {
		mcss.BintrayDetails.SetUser(path.Subject)
	}

	log.Info("Requesting content sync...")
	client, err := httpclient.ClientBuilder().Build()
	if err != nil {
		return err
	}

	requestContent, err := json.Marshal(p)
	if err != nil {
		return errorutils.CheckError(err)
	}

	resp, body, err := client.SendPost(url, requestContent, mcss.BintrayDetails.CreateHttpClientDetails())
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errorutils.CheckError(errors.New("bintray response: " + resp.Status + "\n" + clientutils.IndentJson(body)))
	}

	log.Debug("Bintray response:", resp.Status)
	log.Output(clientutils.IndentJson(body))
	return nil
}

func buildSyncURL(bt auth.BintrayDetails, p *versions.Path) (string, error) {
	if anyEmpty(p.Package, p.Repo, p.Subject, p.Version) {
		return "", errorutils.CheckError(errors.New("invalid path input"))
	}
	return bt.GetApiUrl() + path.Join("maven_central_sync", p.Subject, p.Repo, p.Package, "versions", p.Version), nil
}

func anyEmpty(strs ...string) bool {
	for _, str := range strs {
		if str == "" {
			return true
		}
	}
	return false
}
