package services

import (
	"encoding/json"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
	"strconv"
)

type SyncStatusService struct {
	client *jfroghttpclient.JfrogHttpClient
	auth.ServiceDetails
}

func (ss *SyncStatusService) getHttpDetails() httputils.HttpClientDetails {
	return ss.ServiceDetails.CreateHttpClientDetails()
}

func NewSyncStatusService(client *jfroghttpclient.JfrogHttpClient) *SyncStatusService {
	return &SyncStatusService{client: client}
}

// GetSyncPipelineResourceStatus fetches pipeline sync status
func (ss *SyncStatusService) GetSyncPipelineResourceStatus(repoName, branch string) ([]PipelineSyncStatus, error) {
	// fetch resource ID
	resID, isMultiBranch, resourceErr := GetPipelineResourceID(ss.client, ss.GetUrl(), repoName, ss.getHttpDetails())
	if resourceErr != nil {
		log.Error("Unable to fetch resourceID for: ", repoName)
		return []PipelineSyncStatus{}, resourceErr
	}
	queryParams := make(map[string]string, 0)
	if isMultiBranch {
		queryParams["pipelineSourceBranches"] = branch
	}
	queryParams["pipelineSourceIds"] = strconv.Itoa(resID)

	uriVal, errURL := constructPipelinesURL(queryParams, ss.ServiceDetails.GetUrl(), pipelineSyncStatus)
	if errURL != nil {
		return []PipelineSyncStatus{}, errURL
	}
	httpDetails := ss.getHttpDetails()
	log.Info("fetching pipeline sync status ...")

	resp, body, _, err := ss.client.SendGet(uriVal, true, &httpDetails)
	if err != nil {
		return []PipelineSyncStatus{}, errorutils.CheckError(err)
	}

	// Response Analysis
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return []PipelineSyncStatus{}, err
	}
	rsc := make([]PipelineSyncStatus, 0)
	jsonErr := json.Unmarshal(body, &rsc)
	if jsonErr != nil {
		return []PipelineSyncStatus{}, errorutils.CheckError(jsonErr)
	}

	return rsc, nil
}
