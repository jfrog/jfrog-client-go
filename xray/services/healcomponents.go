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
	log.Debug("Sending request to Xray component-resolution API: %s", chs.getUrl())
	resp, body, err := chs.client.SendPost(chs.getUrl(), body, &httpDetails)
	if err != nil {
		return nil, fmt.Errorf("failed while attempting to resolve component: %w", err)
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return nil, fmt.Errorf("got unexpected server response while attempting to resolve component: %w", err)
	}
	var response ComponentResolutionResponse
	if err = json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to decode component resolution response: %w", err)
	}
	return &response, nil
}

func (chs *ComponentsHealService) mockResponse(req ComponentResolutionRequest) (ComponentResolutionResponse, error) {
	if req.BuildTool == "maven" {
		return ComponentResolutionResponse{Content: json.RawMessage(fakeMavenLockfile), Changes: []Change{
			{
				Package:         "org.apache.commons:commons-lang3:3.14.0",
				BeforeIntegrity: "sha256-orig",
				AfterIntegrity:  "sha256-chainguard",
			},
		}}, nil
	}
	if req.BuildTool == "npm" {
		return ComponentResolutionResponse{Content: json.RawMessage(fakeNpmLockfile), Changes: []Change{
			{
				Package:         "lodash",
				BeforeIntegrity: "sha512-v2kDEe57lecTulaDIuNTPy3Ry4gLGJ6Z1O3vE1krgXZNrsQ+LFTGHVxVjcXPs17LhbZVGedAJv8XZ1tvj5FvSg==",
				AfterIntegrity:  "sha512-Hpyrx+puvIK8/81t1qrv51FEytpvZ78WB88A/NQYFttbH0nQGXQETXvVBK7R+meEGcdC2E8US4QLY6tMAIE2Vw==",
			},
		}}, nil
	}
	return ComponentResolutionResponse{}, fmt.Errorf("unsupported build tool: %s", req.BuildTool)
}

const fakeMavenLockfile = `{
  "groupId": "com.example",
  "artifactId": "demo-app",
  "version": "1.0.0",
  "lockFileVersion": 1,
  "dependencies": [{
    "groupId": "org.apache.commons",
    "artifactId": "commons-lang3",
    "version": "3.14.0",
    "checksumAlgorithm": "SHA-256",
    "checksum": "sha256-chainguard",
    "resolved": "https://example.jfrog.io/artifactory/maven-virtual-chainguard/org/apache/commons/commons-lang3/3.14.0/commons-lang3-3.14.0.jar",
    "id": "org.apache.commons:commons-lang3:3.14.0"
  }]
}`

const (
	fakeNpmLockfile = `{
  "name": "xray-simple-npm-app",
  "version": "1.0.0",
  "lockfileVersion": 3,
  "requires": true,
  "packages": {
    "": {
      "name": "xray-simple-npm-app",
      "version": "1.0.0",
      "dependencies": {
        "lodash": "4.17.21"
      }
    },
    "node_modules/lodash": {
      "version": "4.17.21",
      "resolved": "https://z0xraylnp2.jfrogdev.org/artifactory/api/npm/npm-virtual-chainguard/lodash/-/lodash-4.17.21.tgz",
      "integrity": "sha512-Hpyrx+puvIK8/81t1qrv51FEytpvZ78WB88A/NQYFttbH0nQGXQETXvVBK7R+meEGcdC2E8US4QLY6tMAIE2Vw=="
    }
  }
}`
)

type Change struct {
	Package         string `json:"package"`
	BeforeIntegrity string `json:"before_integrity"`
	AfterIntegrity  string `json:"after_integrity"`
}

type ComponentResolutionRequest struct {
	BuildTool string          `json:"build-tool"`
	Repo      string          `json:"repo"`
	Lockfile  json.RawMessage `json:"lockfile"`
}

type ComponentResolutionResponse struct {
	Content json.RawMessage `json:"lockfile"`
	Changes []Change        `json:"changes,omitempty"`
}
