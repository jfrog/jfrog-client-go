package services

import (
	"encoding/json"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	utils2 "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
)

const (
	conflictErrorMessage               = "Bundle already exists"
	ReleaseBundleImportRestApiEndpoint = "api/release/import/"
	octetStream                        = "application/octet-stream"
)

type releaseService struct {
	client     *jfroghttpclient.JfrogHttpClient
	ArtDetails auth.ServiceDetails
}

type ErrorResponseWithMessage struct {
	Errors []ErrorDetail `json:"errors"`
}

type ErrorDetail struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func NewReleaseService(artDetails auth.ServiceDetails, client *jfroghttpclient.JfrogHttpClient) *releaseService {
	return &releaseService{client: client, ArtDetails: artDetails}
}

func (rs *releaseService) GetJfrogHttpClient() *jfroghttpclient.JfrogHttpClient {
	return rs.client
}

func (rs *releaseService) ImportReleaseBundle(filePath string) (err error) {
	// Load desired file
	content, err := fileutils.ReadFile(filePath)
	if err != nil {
		return
	}
	// Upload file
	httpClientsDetails := rs.ArtDetails.CreateHttpClientDetails()

	url := utils2.AddTrailingSlashIfNeeded(rs.ArtDetails.GetUrl() + ReleaseBundleImportRestApiEndpoint)

	utils.SetContentType(octetStream, &httpClientsDetails.Headers)
	var resp *http.Response
	var body []byte
	log.Info("Uploading archive...")
	if resp, body, err = rs.client.SendPost(url, content, &httpClientsDetails); err != nil {
		return
	}
	// When a release bundle already exists, the API returns 400.
	// Check the error message, and if it's a conflict, don't fail the operation.
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
	log.Info("Release Bundle Imported Successfully")
	return
}
