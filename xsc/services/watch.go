package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"

	api "github.com/jfrog/jfrog-client-go/xray/services/utils"
	xscutils "github.com/jfrog/jfrog-client-go/xsc/services/utils"
)

const (
	watchResourceAPIUrl   = "watches/resource"
	gitRepoResourceUrlKey = "git_repository"
	projectResourceUrlKey = "project"
)

// WatchService defines the http client and Xray details
type WatchService struct {
	client      *jfroghttpclient.JfrogHttpClient
	XrayDetails auth.ServiceDetails
}

// NewWatchService creates a new Xray Watch Service
func NewWatchService(client *jfroghttpclient.JfrogHttpClient) *WatchService {
	return &WatchService{client: client}
}

// GetResourceWatches retrieves the active watches that are associated with a specific git repository and project
func (xws *WatchService) GetResourceWatches(gitRepo, project string) (watches *api.ResourcesWatchesBody, err error) {
	if gitRepo == "" && project == "" {
		return nil, errors.New("no resources provided")
	}
	httpClientsDetails := xws.XrayDetails.CreateHttpClientDetails()
	log.Info(fmt.Sprintf("Getting resources (%s) active watches...", getResourcesString(gitRepo, project)))
	resp, body, _, err := xws.client.SendGet(xws.getWatchURL(gitRepo, project), true, &httpClientsDetails)
	if err != nil {
		return
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return
	}
	log.Debug(fmt.Sprintf("Xray response (status %s): %s", resp.Status, body))
	watches = &api.ResourcesWatchesBody{}
	err = json.Unmarshal(body, &watches)
	if err != nil {
		return nil, errors.New("failed un-marshalling resources watches body")
	}
	log.Info(fmt.Sprintf("Found %d active watches", len(watches.GitRepositoryWatches)+len(watches.ProjectWatches)))
	return
}

func getResourcesString(gitRepo, project string) string {
	providedResources := []string{}
	if gitRepo != "" {
		providedResources = append(providedResources, fmt.Sprintf("git repository: %s", gitRepo))
	}
	if project != "" {
		providedResources = append(providedResources, fmt.Sprintf("project: %s", project))
	}
	return strings.Join(providedResources, ", ")
}

func (xws *WatchService) getWatchURL(gitRepo, project string) string {
	url := utils.AddTrailingSlashIfNeeded(xws.XrayDetails.GetUrl()) + xscutils.XscInXraySuffix + watchResourceAPIUrl
	params := []string{}
	if gitRepo != "" {
		params = append(params, fmt.Sprintf("%s=%s", gitRepoResourceUrlKey, gitRepo))
	}
	if project != "" {
		params = append(params, fmt.Sprintf("%s=%s", projectResourceUrlKey, project))
	}
	return url + "?" + strings.Join(params, "&")
}
