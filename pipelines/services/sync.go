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
	"path"
	"strconv"
)

type SyncService struct {
	client *jfroghttpclient.JfrogHttpClient
	auth.ServiceDetails
}

func (ss *SyncService) getHttpDetails() httputils.HttpClientDetails {
	return ss.ServiceDetails.CreateHttpClientDetails()
}

func NewSyncService(client *jfroghttpclient.JfrogHttpClient) *SyncService {
	return &SyncService{client: client}
}

// SyncPipelineSource trigger sync for pipeline resource
func (ss *SyncService) SyncPipelineSource(branch string, repoName string) (int, []byte, error) {
	// fetch resource ID
	resID, _, resourceErr := ss.GetPipelineResourceID(repoName)
	if resourceErr != nil {
		log.Error("unable to fetch resourceID for: ", repoName)
		return 0, []byte{}, resourceErr
	}
	log.Info("Triggering pipeline source sync ...")

	// trigger sync
	httpDetails := ss.getHttpDetails()
	queryParams := map[string]string{
		"sync":   "true",
		"branch": branch,
	}

	apiPath := path.Join(pipelineResources, strconv.Itoa(resID))
	uriVal, errURL := ss.constructURL(apiPath, queryParams)
	if errURL != nil {
		return 0, []byte{}, errURL
	}
	resp, body, _, httpErr := ss.client.SendGet(uriVal, true, &httpDetails)
	if httpErr != nil {
		return 0, body, errorutils.CheckError(httpErr)
	}
	if err := errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return resp.StatusCode, body, errorutils.CheckError(err)
	}
	log.Info("Triggered pipeline sync successfully")
	return resp.StatusCode, body, nil
}

// GetPipelineResourceID fetches resource ID for given full repository name
func (ss *SyncService) GetPipelineResourceID(repoName string) (int, bool, error) {
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
	// Response analysis
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return resp.StatusCode, false, err
	}
	log.Debug("Received resource id")
	p := make([]PipelineResources, 0)
	err = json.Unmarshal(body, &p)
	if err != nil {
		log.Error("Failed to unmarshal json response")
		return 0, false, errorutils.CheckError(err)
	}
	for _, res := range p {
		if res.RepositoryFullName == repoName {
			log.Debug("received repository name ", repoName, "is multi branch ", res.IsMultiBranch)
			return res.ID, *res.IsMultiBranch, nil
		}
	}
	return 0, false, nil
}

// constructURL from server config and api for fetching run status for a given branch
// and prepares URL string
func (ss *SyncService) constructURL(api string, qParams map[string]string) (string, error) {
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
