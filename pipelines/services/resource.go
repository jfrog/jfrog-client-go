package services

import (
	"encoding/json"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
)

// GetPipelineResource fetches pipeline resource information for the given full repository name
func GetPipelineResource(client *jfroghttpclient.JfrogHttpClient, apiURL, repoName string, httpDetails httputils.HttpClientDetails) (*PipelineResources, error) {
	// Query params
	queryParams := make(map[string]string, 0)
	uriVal, errURL := constructPipelinesURL(queryParams, apiURL, pipelineResources)
	if errURL != nil {
		return nil, errURL
	}
	resp, body, _, err := client.SendGet(uriVal, true, &httpDetails)
	if err != nil {
		return nil, errorutils.CheckError(err)
	}

	// Response Analysis
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return nil, err
	}
	p := make([]PipelineResources, 0)
	err = json.Unmarshal(body, &p)
	if err != nil {
		return nil, errorutils.CheckError(err)
	}
	for _, res := range p {
		if res.RepositoryFullName == repoName {
			log.Debug("Received repository name ", repoName, "is multi branch ", *res.IsMultiBranch)
			return &res, nil
		}
	}
	return nil, nil
}
