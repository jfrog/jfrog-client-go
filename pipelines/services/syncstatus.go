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
	"strconv"
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
func (ss *SyncStatusService) GetSyncPipelineResourceStatus(repoName, branch string) ([]PipelineSyncStatus, []byte, error) {
	// fetch resource ID
	resID, isMultiBranch, resourceErr := ss.getPipelineResourceID(repoName)
	if resourceErr != nil {
		log.Error("unable to fetch resourceID for: ", repoName)
		return []PipelineSyncStatus{}, []byte{}, errorutils.CheckError(resourceErr)
	}
	queryParams := make(map[string]string, 0)
	if isMultiBranch {
		queryParams["pipelineSourceBranches"] = branch
		queryParams["pipelineSourceIds"] = strconv.Itoa(resID)
	} else {
		queryParams["pipelineSourceIds"] = strconv.Itoa(resID)
	}
	uriVal, err := ss.constructURL(pipelineSyncStatus, queryParams)
	if err != nil {
		return []PipelineSyncStatus{}, []byte{}, errorutils.CheckError(err)
	}
	httpDetails := ss.getHttpDetails()
	log.Info("fetching pipeline sync status ...")

	resp, body, _, err := ss.client.SendGet(uriVal, true, &httpDetails)
	if err != nil {
		return []PipelineSyncStatus{}, body, errorutils.CheckError(err)
	}

	// Response Analysis
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return []PipelineSyncStatus{}, body, errorutils.CheckError(err)
	}
	rsc := make([]PipelineSyncStatus, 0)
	jsonErr := json.Unmarshal(body, &rsc)
	if jsonErr != nil {
		return []PipelineSyncStatus{}, body, errorutils.CheckError(jsonErr)
	}

	return rsc, body, nil
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

// getPipelineResourceID fetches resource ID for given full repository name
func (ss *SyncStatusService) getPipelineResourceID(repoName string) (int, bool, error) {
	httpDetails := ss.getHttpDetails()
	queryParams := make(map[string]string, 0)

	uriVal, errURL := ss.constructURL(pipelineResources, queryParams)
	if errURL != nil {
		return 0, false, errorutils.CheckError(errURL)
	}

	resp, body, _, err := ss.client.SendGet(uriVal, true, &httpDetails)
	if err != nil {
		return 0, false, errorutils.CheckError(err)
	}
	// Response Analysis
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return 0, false, err
	}
	if resp.StatusCode == http.StatusOK {
		log.Debug("received resource id")
	}
	p := make([]PipelineResources, 0)
	err = json.Unmarshal(body, &p)
	if err != nil {
		log.Error("Failed to unmarshal json response")
		return 0, false, errorutils.CheckError(err)
	}
	for _, res := range p {
		if res.RepositoryFullName == repoName && res.IsMultiBranch {
			return res.ID, res.IsMultiBranch, nil
		} else if res.RepositoryFullName == repoName && !res.IsMultiBranch {
			log.Debug("received repository name ", repoName, "is multi branch ", res.IsMultiBranch)
			return res.ID, res.IsMultiBranch, nil
		}
	}
	return 0, false, nil
}
