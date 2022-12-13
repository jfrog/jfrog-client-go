package services

import (
	"encoding/json"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
	"net/url"
)

type SyncStatusService struct {
	client *jfroghttpclient.JfrogHttpClient
	auth.ServiceDetails
}

func (ss *SyncStatusService) getHttpDetails() httputils.HttpClientDetails {
	httpDetails := ss.ServiceDetails.CreateHttpClientDetails()
	return httpDetails
}

func NewSyncStatusService(client *jfroghttpclient.JfrogHttpClient) *SyncStatusService {
	return &SyncStatusService{client: client}
}

// GetSyncPipelineResourceStatus fetches pipeline sync status
func (ss *SyncStatusService) GetSyncPipelineResourceStatus(branch string) ([]PipelineSyncStatus, error) {
	queryParams := make(map[string]string, 0)
	queryParams["pipelineSourceBranches"] = branch
	uriVal, err := ss.constructURL(pipelineSyncStatus, queryParams)
	if err != nil {
		return []PipelineSyncStatus{}, errorutils.CheckError(err)
	}
	httpDetails := ss.getHttpDetails()
	log.Info("fetching pipeline sync status ...")

	resp, body, _, err := ss.client.SendGet(uriVal, true, &httpDetails)
	if err != nil {
		return []PipelineSyncStatus{}, errorutils.CheckError(err)
	}

	// Response Analysis
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return []PipelineSyncStatus{}, errorutils.CheckError(err)
	}

	rsc := make([]PipelineSyncStatus, 0)
	jsonErr := json.Unmarshal(body, &rsc)
	if jsonErr != nil {
		return []PipelineSyncStatus{}, errorutils.CheckError(jsonErr)
	}

	return rsc, nil
}

// constructURL from server config and api for fetching run status for a given branch
// and prepares URL string
func (ss *SyncStatusService) constructURL(api string, qParams map[string]string) (string, error) {
	uri, err := url.Parse(ss.ServiceDetails.GetUrl() + api)
	if err != nil {
		log.Error("Failed to parse pipelines fetch run status url")
		return "", errorutils.CheckError(err)
	}
	queryString := uri.Query()
	for k, v := range qParams {
		queryString.Set(k, v)
	}
	uri.RawQuery = queryString.Encode()

	return uri.String(), nil
}
