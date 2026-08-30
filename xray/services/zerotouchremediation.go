package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
)

const componentResolutionApi = "api/v1/ztr/lockfile/remediate"

type ZeroTouchRemediationService struct {
	client          *jfroghttpclient.JfrogHttpClient
	XrayDetails     auth.ServiceDetails
	ScopeProjectKey string
}

func NewZeroTouchRemediationService(client *jfroghttpclient.JfrogHttpClient) *ZeroTouchRemediationService {
	return &ZeroTouchRemediationService{client: client}
}

func (ztr *ZeroTouchRemediationService) getUrl() string {
	return utils.AppendScopedProjectKeyParam(utils.AddTrailingSlashIfNeeded(ztr.XrayDetails.GetUrl())+componentResolutionApi, ztr.ScopeProjectKey)
}

func (ztr *ZeroTouchRemediationService) Remediate(req ComponentResolutionRequest) (*ComponentResolutionResponse, bool, error) {
	httpDetails := ztr.XrayDetails.CreateHttpClientDetails()
	// SendPost retries on 5xx, disable retries for this request
	httpDetails.AddPreRetryInterceptor(func() bool { return false })
	body, err := json.Marshal(req)
	if err != nil {
		return nil, false, err
	}
	resp, body, err := ztr.client.SendPost(ztr.getUrl(), body, &httpDetails)
	if err != nil {
		return nil, false, fmt.Errorf("failed while attempting to resolve component: %w", err)
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK, http.StatusServiceUnavailable); err != nil {
		return nil, false, fmt.Errorf("got unexpected server response while attempting to resolve component: %w", err)
	}
	if resp.StatusCode == http.StatusServiceUnavailable {
		return &ComponentResolutionResponse{Lockfile: req.Lockfile}, true, nil
	}
	var response ComponentResolutionResponse
	if err = json.Unmarshal(body, &response); err != nil {
		return nil, false, fmt.Errorf("failed to decode component resolution response: %w", err)
	}
	return &response, false, nil
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
