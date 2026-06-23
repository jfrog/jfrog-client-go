package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

const componentResolutionApi = "api/v1/lockfile/heal"

type ComponentsHealService struct {
	client          *jfroghttpclient.JfrogHttpClient
	XrayDetails     auth.ServiceDetails
	ScopeProjectKey string
}

func NewComponentsHealService(client *jfroghttpclient.JfrogHttpClient) *ComponentsHealService {
	return &ComponentsHealService{client: client}
}

func (chs *ComponentsHealService) getUrl() string {
	return utils.AppendScopedProjectKeyParam(utils.AddTrailingSlashIfNeeded(chs.XrayDetails.GetUrl())+componentResolutionApi, chs.ScopeProjectKey)
}

func (chs *ComponentsHealService) Heal(req ComponentResolutionRequest) (*ComponentResolutionResponse, error) {
	httpDetails := chs.XrayDetails.CreateHttpClientDetails()
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, body, err := chs.client.SendPost(chs.getUrl(), body, &httpDetails)
	if err != nil {
		return nil, fmt.Errorf("failed while attempting to resolve component: %w", err)
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK, http.StatusServiceUnavailable); err != nil {
		return nil, fmt.Errorf("got unexpected server response while attempting to resolve component: %w", err)
	}
	if resp.StatusCode == http.StatusServiceUnavailable {
		log.Warn("Self-heal is disabled on JFrog Xray server. Skipping component resolution.")
		return &ComponentResolutionResponse{Lockfile: req.Lockfile}, nil
	}
	var response ComponentResolutionResponse
	if err = json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to decode component resolution response: %w", err)
	}
	return &response, nil
}

type Change struct {
	Package         string `json:"package"`
	BeforeIntegrity string `json:"before_integrity"`
	AfterIntegrity  string `json:"after_integrity"`
}

type ComponentResolutionRequest struct {
	BuildTool string `json:"build-tool"`
	Repo      string `json:"repo"`
	Lockfile  string `json:"lockfile"`
}

type ComponentResolutionResponse struct {
	Lockfile string   `json:"lockfile"`
	Changes  []Change `json:"changes,omitempty"`
}
