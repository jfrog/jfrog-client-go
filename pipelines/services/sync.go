package services

import (
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
	"path"
	"strconv"
)

type SyncService struct {
	client *jfroghttpclient.JfrogHttpClient
	auth.ServiceDetails
}

func (ss *SyncService) GetHttpDetails() httputils.HttpClientDetails {
	return ss.ServiceDetails.CreateHttpClientDetails()
}

func NewSyncService(client *jfroghttpclient.JfrogHttpClient) *SyncService {
	return &SyncService{client: client}
}

func (ss *SyncService) GetHTTPClient() *jfroghttpclient.JfrogHttpClient {
	return ss.client
}

func (ss *SyncService) GetServiceURL() string {
	return ss.GetUrl()
}

// SyncPipelineSource trigger sync for pipeline resource
func (ss *SyncService) SyncPipelineSource(branch string, repoName string) error {
	// fetch resource ID
	resID, _, resourceErr := GetPipelineResourceID(ss.client, ss.GetUrl(), repoName, ss.GetHttpDetails())
	if resourceErr != nil {
		log.Error("unable to fetch resourceID for: ", repoName)
		return resourceErr
	}
	log.Info("Triggering pipeline source sync ...")

	// trigger sync
	httpDetails := ss.GetHttpDetails()
	queryParams := map[string]string{
		"sync":   "true",
		"branch": branch,
	}

	apiPath := path.Join(pipelineResources, strconv.Itoa(resID))
	uriVal, errURL := constructPipelinesURL(queryParams, ss.ServiceDetails.GetUrl(), apiPath)
	if errURL != nil {
		return errURL
	}
	resp, body, _, httpErr := ss.client.SendGet(uriVal, true, &httpDetails)
	if httpErr != nil {
		return errorutils.CheckError(httpErr)
	}
	if err := errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return errorutils.CheckError(err)
	}
	log.Info("Triggered pipeline sync successfully")
	return nil
}
