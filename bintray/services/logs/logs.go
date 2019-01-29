package logs

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

func NewService(client *httpclient.HttpClient) *LogsService {
	us := &LogsService{client: client}
	return us
}

type LogsService struct {
	client         *httpclient.HttpClient
	BintrayDetails auth.BintrayDetails
}

func (ls *LogsService) List(versionPath *versions.Path) error {
	if ls.BintrayDetails.GetUser() == "" {
		ls.BintrayDetails.SetUser(versionPath.Subject)
	}
	listUrl := ls.BintrayDetails.GetApiUrl() + path.Join("packages", versionPath.Subject, versionPath.Repo, versionPath.Package, "logs")
	httpClientsDetails := ls.BintrayDetails.CreateHttpClientDetails()
	log.Info("Getting logs...")
	client, err := httpclient.ClientBuilder().Build()
	if err != nil {
		return err
	}
	resp, body, _, _ := client.SendGet(listUrl, true, httpClientsDetails)

	if resp.StatusCode != http.StatusOK {
		return errorutils.CheckError(errors.New("Bintray response: " + resp.Status + "\n" + clientutils.IndentJson(body)))
	}

	log.Debug("Bintray response:", resp.Status)
	log.Output(clientutils.IndentJson(body))
	return nil
}

func (ls *LogsService) Download(versionPath *versions.Path, logName string) error {
	if ls.BintrayDetails.GetUser() == "" {
		ls.BintrayDetails.SetUser(versionPath.Subject)
	}
	downloadUrl := ls.BintrayDetails.GetApiUrl() + path.Join("packages", versionPath.Subject, versionPath.Repo, versionPath.Package, "logs")

	httpClientsDetails := ls.BintrayDetails.CreateHttpClientDetails()
	log.Info("Downloading logs...")
	client, err := httpclient.ClientBuilder().Build()
	if err != nil {
		return err
	}
	details := &httpclient.DownloadFileDetails{
		FileName:      logName,
		DownloadPath:  downloadUrl,
		LocalPath:     "",
		LocalFileName: logName}
	resp, err := client.DownloadFile(details, "", httpClientsDetails, utils.BintrayDownloadRetries, false)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errorutils.CheckError(errors.New("Bintray response: " + resp.Status))
	}
	log.Debug("Bintray response:", resp.Status)
	log.Info("Downloaded log.")
	return nil
}
