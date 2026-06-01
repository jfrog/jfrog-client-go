package services

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/httpclient"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
)

const (
	downloadIndexerAPI = "api/v1/indexer-resources/download"
)

type IndexerService struct {
	client          *jfroghttpclient.JfrogHttpClient
	XrayDetails     auth.ServiceDetails
	ScopeProjectKey string
}

func NewIndexerService(client *jfroghttpclient.JfrogHttpClient) *IndexerService {
	return &IndexerService{client: client}
}

func (is *IndexerService) Download(localDirPath, localBinaryName string) (string, error) {
	httpClientDetails := is.XrayDetails.CreateHttpClientDetails()
	url := is.getUrlForDownloadApi()
	// Download the indexer from Xray to the provided directory
	downloadFileDetails := &httpclient.DownloadFileDetails{DownloadPath: url, LocalPath: localDirPath, LocalFileName: localBinaryName}
	resp, err := is.client.DownloadFile(downloadFileDetails, "", &httpClientDetails, false, false)
	if err != nil {
		return "", fmt.Errorf("failed while attempting to download %q: %w", url, err)
	}
	if err = errorutils.CheckResponseStatus(resp, http.StatusOK); err != nil {
		if resp.StatusCode == http.StatusUnauthorized {
			err = fmt.Errorf("%s\nHint: It appears that the credentials provided do not have sufficient permissions for JFrog Xray. This could be due to either incorrect credentials or limited permissions restricted to Artifactory only", err.Error())
		}
		return "", fmt.Errorf("failed to download %q: %w", url, err)
	}
	// Add execution permissions to the indexer binary
	downloadedFilePath := filepath.Join(localDirPath, localBinaryName)
	if err = os.Chmod(downloadedFilePath, 0o755); err != nil {
		return "", errorutils.CheckError(err)
	}
	return downloadedFilePath, errorutils.CheckError(err)
}

func (is *IndexerService) getUrlForDownloadApi() string {
	return utils.AppendScopedProjectKeyParam(fmt.Sprintf("%s%s/%s/%s", is.XrayDetails.GetUrl(), downloadIndexerAPI, runtime.GOOS, runtime.GOARCH), is.ScopeProjectKey)
}
