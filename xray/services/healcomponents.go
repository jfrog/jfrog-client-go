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
		return ComponentResolutionResponse{Lockfile: fakeMavenHealedPom, Changes: []Change{
			{
				Package:         "org.springframework:spring-core:5.3.39-0.cgr.4",
				BeforeIntegrity: "5.3.20",
				AfterIntegrity:  "5.3.39-0.cgr.4",
			},
		}}, nil
	}
	if req.BuildTool == "npm" {
		return ComponentResolutionResponse{Lockfile: fakeNpmLockfile, Changes: []Change{
			{
				Package:         "lodash",
				BeforeIntegrity: "sha512-v2kDEe57lecTulaDIuNTPy3Ry4gLGJ6Z1O3vE1krgXZNrsQ+LFTGHVxVjcXPs17LhbZVGedAJv8XZ1tvj5FvSg==",
				AfterIntegrity:  "sha512-Hpyrx+puvIK8/81t1qrv51FEytpvZ78WB88A/NQYFttbH0nQGXQETXvVBK7R+meEGcdC2E8US4QLY6tMAIE2Vw==",
			},
		}}, nil
	}
	return ComponentResolutionResponse{}, fmt.Errorf("unsupported build tool: %s", req.BuildTool)
}

// fakeMavenHealedPom is a sample healed pom.xml returned by the maven heal stub.
const fakeMavenHealedPom = `<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0"
         xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
         xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd">
  <modelVersion>4.0.0</modelVersion>
  <groupId>com.example</groupId>
  <artifactId>demo-app</artifactId>
  <version>1.0.0</version>
  <dependencyManagement>
    <dependencies>
      <dependency>
        <groupId>org.springframework</groupId>
        <artifactId>spring-core</artifactId>
        <version>5.3.39-0.cgr.4</version>
      </dependency>
    </dependencies>
  </dependencyManagement>
</project>`

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
	BuildTool string `json:"build-tool"`
	Repo      string `json:"repo"`
	Lockfile  string `json:"lockfile"`
}

type ComponentResolutionResponse struct {
	Lockfile string   `json:"lockfile"`
	Changes  []Change `json:"changes,omitempty"`
}
