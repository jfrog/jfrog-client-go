package services

import (
	"encoding/json"
	"fmt"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	utils2 "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
	"strings"
)

const (
	conflictErrorMessage                    = "Bundle already exists"
	ReleaseBundleImportRestApiEndpoint      = "api/release/import/"
	octetStream                             = "application/octet-stream"
	ReleaseBundleExistInRbV2RestApiEndpoint = "lifecycle/api/v2/release_bundle/existence"
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

type isReleaseBundleExistResponse struct {
	Exists bool `json:"exists"`
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

func (rs *releaseService) IsReleaseBundleExistInRbV2(project, bundleNameAndVersion string) (bool, error) {
	httpClientsDetails := rs.ArtDetails.CreateHttpClientDetails()
	if project != "" {
		project = fmt.Sprintf("project=%s&", project)
	} else {
		project = fmt.Sprintf("project=default")
	}

	rtUrl := strings.Replace(rs.ArtDetails.GetUrl(), "/artifactory", "", 1)
	url := fmt.Sprintf("%s%s/%s/?%s", rtUrl, ReleaseBundleExistInRbV2RestApiEndpoint, bundleNameAndVersion, project)
	resp, body, _, err := rs.client.SendGet(url, true, &httpClientsDetails)
	if err != nil {
		return false, err
	}

	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return false, err
	}

	response := &isReleaseBundleExistResponse{}
	if err := json.Unmarshal(body, response); err != nil {
		return false, err
	}

	return response.Exists, nil
}
