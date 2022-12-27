package services

import (
	"encoding/json"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
)

// GetPipelineResourceID fetches resource ID for given full repository name
func GetPipelineResourceID(client *jfroghttpclient.JfrogHttpClient, apiURL, repoName string, httpDetails httputils.HttpClientDetails) (int, bool, error) {
	queryParams := make(map[string]string, 0)

	uriVal, errURL := constructPipelinesURL(queryParams, apiURL, pipelineResources)
	if errURL != nil {
		return 0, false, errURL
	}

	resp, body, _, err := client.SendGet(uriVal, true, &httpDetails)
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
		if res.RepositoryFullName == repoName {
			log.Debug("received repository name ", repoName, "is multi branch ", res.IsMultiBranch)
			return res.ID, *res.IsMultiBranch, nil
		}
	}
	return 0, false, nil
}
