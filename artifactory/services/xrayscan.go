package services

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/httpclient"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
)

const apiUri = "api/xray/scanBuild"

// Retrying to resume the scan 10 times after a stable connection
const consecutiveRetries = 10

// Expecting \r\n every 30 seconds
const connectionTimeout = 90 * time.Second

// 15 seconds sleep between retry
const sleepBetweenRetries = 15 * time.Second
const stableConnectionWindow = 100 * time.Second
const fatalFailureStatus = -1

type XrayScanService struct {
	client     *jfroghttpclient.JfrogHttpClient
	ArtDetails auth.ServiceDetails
}

func NewXrayScanService(client *jfroghttpclient.JfrogHttpClient) *XrayScanService {
	return &XrayScanService{client: client}
}

// Deprecated legacy scan build. The new build scan command is in "/xray/commands/scan/buildscan"
func (ps *XrayScanService) ScanBuild(scanParams XrayScanParams) ([]byte, error) {
	url := ps.ArtDetails.GetUrl()
	requestFullUrl, err := utils.BuildArtifactoryUrl(url, apiUri, make(map[string]string))
	if err != nil {
		return []byte{}, err
	}
	data := XrayScanBody{
		BuildName:   scanParams.GetBuildName(),
		BuildNumber: scanParams.GetBuildNumber(),
		Project:     scanParams.GetProjectKey(),
		Context:     clientutils.GetUserAgent(),
	}

	requestContent, err := json.Marshal(data)
	if err != nil {
		return []byte{}, errorutils.CheckError(err)
	}

	connection := httpclient.RetryableConnection{
		ReadTimeout:            connectionTimeout,
		RetriesNum:             consecutiveRetries,
		StableConnectionWindow: stableConnectionWindow,
		SleepBetweenRetries:    sleepBetweenRetries,
		ConnectHandler: func() (*http.Response, error) {
			return ps.execScanRequest(requestFullUrl, requestContent)
		},
		ErrorHandler: func(content []byte) error {
			return checkForXrayResponseError(content, true)
		},
	}
	result, err := connection.Do()
	if err != nil {
		return []byte{}, err
	}

	return result, nil
}

func isFatalScanError(errResp *errorResponse) bool {
	if errResp == nil {
		return false
	}
	for _, v := range errResp.Errors {
		if v.Status == fatalFailureStatus {
			return true
		}
	}
	return false
}

func checkForXrayResponseError(content []byte, ignoreFatalError bool) error {
	respErrors := &errorResponse{}
	err := json.Unmarshal(content, respErrors)
	if errorutils.CheckError(err) != nil {
		return err
	}

	if respErrors.Errors == nil {
		return nil
	}

	if ignoreFatalError && isFatalScanError(respErrors) {
		// fatal error should be interpreted as no errors so no more retries will accrue
		return nil
	}
	return errorutils.CheckErrorf("Server response: " + string(content))
}

func (ps *XrayScanService) execScanRequest(url string, content []byte) (*http.Response, error) {
	httpClientsDetails := ps.ArtDetails.CreateHttpClientDetails()
	utils.SetContentType("application/json", &httpClientsDetails.Headers)

	// The scan build operation can take a long time to finish.
	// To keep the connection open, when Xray starts scanning the build, it starts sending new-lines
	// on the open channel. This tells the client that the operation is still in progress and the
	// connection does not get timed out.
	// We need make sure the new-lines are not buffered on the nginx and are flushed
	// as soon as Xray sends them.
	utils.DisableAccelBuffering(&httpClientsDetails.Headers)

	resp, body, _, err := ps.client.Send("POST", url, content, true, false, &httpClientsDetails, "")
	if err != nil {
		return resp, err
	}
	return resp, errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK)
}

type errorResponse struct {
	Errors []errorsStatusResponse `json:"errors,omitempty"`
}

type errorsStatusResponse struct {
	Status int `json:"status,omitempty"`
}

type XrayScanBody struct {
	BuildName   string `json:"buildName,omitempty"`
	BuildNumber string `json:"buildNumber,omitempty"`
	Project     string `json:"project,omitempty"`
	Context     string `json:"context,omitempty"`
}

type XrayScanParams struct {
	BuildName   string
	BuildNumber string
	ProjectKey  string
}

func (bp *XrayScanParams) GetBuildName() string {
	return bp.BuildName
}

func (bp *XrayScanParams) GetBuildNumber() string {
	return bp.BuildNumber
}

func (bp *XrayScanParams) GetProjectKey() string {
	return bp.ProjectKey
}

func NewXrayScanParams() XrayScanParams {
	return XrayScanParams{}
}
