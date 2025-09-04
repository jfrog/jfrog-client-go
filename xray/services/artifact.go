package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
)

const (
	statusAPI = "api/v1/artifact/status"
)

// ArtifactStatus represents the status of an artifact in Xray
const (
	ArtifactStatusNotSupported ArtifactStatus = "NOT_SUPPORTED"
	ArtifactStatusNotScanned   ArtifactStatus = "NOT_SCANNED"
	ArtifactStatusPending      ArtifactStatus = "PENDING"
	ArtifactStatusScanning     ArtifactStatus = "SCANNING"
	ArtifactStatusDone         ArtifactStatus = "DONE"
	ArtifactStatusPartial      ArtifactStatus = "PARTIAL"
	ArtifactStatusFailed       ArtifactStatus = "FAILED"
)

type ArtifactStatus string

type ArtifactService struct {
	client          *jfroghttpclient.JfrogHttpClient
	XrayDetails     auth.ServiceDetails
	ScopeProjectKey string
}

// NewArtifactService creates a new service for interacting with artifacts in Xray.
func NewArtifactService(client *jfroghttpclient.JfrogHttpClient) *ArtifactService {
	return &ArtifactService{client: client}
}

func (as *ArtifactService) GetStatus(repo, path string) (response *ArtifactStatusResponse, err error) {
	httpClientsDetails := as.XrayDetails.CreateHttpClientDetails()
	httpClientsDetails.SetContentTypeApplicationJson()

	requestBody, err := json.Marshal(ArtifactStatusRequest{Repository: repo, Path: path})
	if errorutils.CheckError(err) != nil {
		return
	}

	resp, body, err := as.client.SendPost(clientutils.AppendScopedProjectKeyParam(as.XrayDetails.GetUrl()+statusAPI, as.ScopeProjectKey), requestBody, &httpClientsDetails)
	if err != nil {
		return
	}

	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		err = fmt.Errorf("got unexpected server response while attempting to get artifact status for %s/%s:\n%s", repo, path, err.Error())
		return
	}

	response = &ArtifactStatusResponse{}
	if err = json.Unmarshal(body, response); err != nil {
		err = errorutils.CheckErrorf("couldn't parse JFrog Xray artifact status response: %s", err.Error())
	}
	return
}

type ArtifactStatusRequest struct {
	Repository string `json:"repo"`
	Path       string `json:"path"`
}

type ArtifactStatusResponse struct {
	Overall ArtifactScanStatus     `json:"overall"`
	Details ArtifactDetailedStatus `json:"details"`
}

type ArtifactScanStatus struct {
	Status ArtifactStatus `json:"status"`
	// Timestamp in RFC 3339 format
	Timestamp string `json:"time"`
}

type ArtifactDetailedStatus struct {
	Sca                ArtifactScanStatus `json:"sca"`
	ContextualAnalysis ArtifactScanStatus `json:"contextual_analysis"`
	Exposures          ArtifactScanStatus `json:"exposures"`
	Violations         ArtifactScanStatus `json:"violations"`
}

type ArtifactExposureScanStatus struct {
	// Overall status of the exposure scans
	ArtifactScanStatus
	Categories ArtifactExposureCategoriesStatus `json:"categories"`
}

type ArtifactExposureCategoriesStatus struct {
	IaC          ArtifactScanStatus `json:"iac"`
	Secrets      ArtifactScanStatus `json:"secrets"`
	Services     ArtifactScanStatus `json:"services"`
	Applications ArtifactScanStatus `json:"applications"`
}
