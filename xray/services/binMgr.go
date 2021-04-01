package services

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"

	artUtils "github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

const (
	binMgrAPIURL = "api/v1/binMgr"
)

// BinMgrService defines the http client and Xray details
type BinMgrService struct {
	client      *jfroghttpclient.JfrogHttpClient
	XrayDetails auth.ServiceDetails
}

// NewBinMgrService creates a new Xray Binary Manager Service
func NewBinMgrService(client *jfroghttpclient.JfrogHttpClient) *BinMgrService {
	return &BinMgrService{client: client}
}

// GetXrayDetails returns the Xray details
func (xbms *BinMgrService) GetXrayDetails() auth.ServiceDetails {
	return xbms.XrayDetails
}

// GetJfrogHttpClient returns the http client
func (xbms *BinMgrService) GetJfrogHttpClient() *jfroghttpclient.JfrogHttpClient {
	return xbms.client
}

// The getBinMgrURL does not end with a slash
// So, calling functions will need to add it
func (xbms *BinMgrService) getBinMgrURL() string {
	return clientutils.AddTrailingSlashIfNeeded(xbms.XrayDetails.GetUrl()) + binMgrAPIURL
}

// AddBuildsToIndexing will add builds to indexing configuration
func (xbms *BinMgrService) AddBuildsToIndexing(buildNames []string) error {
	payloadBody := addBuildsToIndexBody{BuildNames: buildNames}

	content, err := json.Marshal(payloadBody)
	if err != nil {
		return errorutils.CheckError(err)
	}

	httpClientsDetails := xbms.XrayDetails.CreateHttpClientDetails()
	artUtils.SetContentType("application/json", &httpClientsDetails.Headers)
	var url = xbms.getBinMgrURL() + "/builds"
	var resp *http.Response
	var respBody []byte

	log.Info("Configuring Xray to index the build...")
	resp, respBody, err = xbms.client.SendPost(url, content, &httpClientsDetails)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return errorutils.CheckError(errors.New("Xray response: " + resp.Status + "\n" + clientutils.IndentJson(respBody)))
	}
	log.Debug("Xray response:", resp.Status)
	log.Debug("Done adding builds to indexing configuration.")
	return nil
}

type addBuildsToIndexBody struct {
	BuildNames []string `json:"names"`
}
