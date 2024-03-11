package services

import (
	"encoding/json"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	ioutils "github.com/jfrog/jfrog-client-go/utils/io"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
	"strings"
)

const (
	conflictErrorMessage  = "Bundle already exists"
	importRestApiEndpoint = "api/release/import"
	// TODO change this when the bug is fixed
	// https://jfrog-int.atlassian.net/browse/JR-8542
	tempImportRestApiEndpoint = "ui/api/v1/ui/release/import"
	octetStream               = "application/octet-stream"
)

type ReleaseService struct {
	client     *jfroghttpclient.JfrogHttpClient
	ArtDetails auth.ServiceDetails
	Progress   ioutils.ProgressMgr
}

type ErrorResponseWithMessage struct {
	Errors []ErrorDetail `json:"errors"`
}

type ErrorDetail struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func NewReleaseService(artDetails auth.ServiceDetails, client *jfroghttpclient.JfrogHttpClient) *ReleaseService {
	return &ReleaseService{client: client, ArtDetails: artDetails}
}

func (rs *ReleaseService) GetJfrogHttpClient() *jfroghttpclient.JfrogHttpClient {
	return rs.client
}

func (rs *ReleaseService) ImportReleaseBundle(filePath string) (err error) {
	// Load desired file
	content, err := fileutils.ReadFile(filePath)
	if err != nil {
		return
	}
	// Upload file
	httpClientsDetails := rs.ArtDetails.CreateHttpClientDetails()

	// TODO replace URL when artifactory bug is fixed
	// url := rs.ArtDetails.GetUrl() + importRestApiEndpoint
	tempUrl := strings.TrimSuffix(rs.ArtDetails.GetUrl(), "/artifactory/") + "/" + tempImportRestApiEndpoint

	utils.SetContentType(octetStream, &httpClientsDetails.Headers)
	var resp *http.Response
	var body []byte
	log.Info("Uploading archive...")
	if resp, body, err = rs.client.SendPost(tempUrl, content, &httpClientsDetails); err != nil {
		return
	}
	// When a release bundle already exists, don't return an error message of failure.
	if resp.StatusCode == http.StatusBadRequest {
		response := ErrorResponseWithMessage{}
		if err = json.Unmarshal(body, &response); err != nil {
			return
		}
		if response.Errors[0].Message == conflictErrorMessage {
			log.Warn("Bundle already exists, did not upload a new bundle")
			return
		}
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusAccepted); err != nil {
		return
	}
	log.Info("Upload Successful")
	return
}
