package services

import (
	"encoding/json"
	"errors"
	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	rthttpclient "github.com/jfrog/jfrog-client-go/artifactory/httpclient"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
)

type VirtualRepositoryService struct {
	isUpdate   bool
	client     *rthttpclient.ArtifactoryHttpClient
	ArtDetails auth.ArtifactoryDetails
}

func NewVirtualRepositoryService(client *rthttpclient.ArtifactoryHttpClient, isUpdate bool) *VirtualRepositoryService {
	return &VirtualRepositoryService{client: client, isUpdate: isUpdate}
}

func (vrs *VirtualRepositoryService) GetJfrogHttpClient() *rthttpclient.ArtifactoryHttpClient {
	return vrs.client
}

func (vrs *VirtualRepositoryService) performRequest(params interface{}, repoKey string) error {
	content, err := json.Marshal(params)
	if errorutils.CheckError(err) != nil {
		return err
	}
	httpClientsDetails := vrs.ArtDetails.CreateHttpClientDetails()
	utils.SetContentType("application/vnd.org.jfrog.artifactory.repositories.VirtualRepositoryConfiguration+json", &httpClientsDetails.Headers)
	var url = vrs.ArtDetails.GetUrl() + "api/repositories/" + repoKey
	var operationString string
	var resp *http.Response
	var body []byte
	if vrs.isUpdate {
		log.Info("Updating virtual repository......")
		operationString = "updating"
		resp, body, err = vrs.client.SendPost(url, content, &httpClientsDetails)
	} else {
		log.Info("Creating virtual repository......")
		operationString = "creating"
		resp, body, err = vrs.client.SendPut(url, content, &httpClientsDetails)
	}
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errorutils.CheckError(errors.New("Artifactory response: " + resp.Status + "\n" + clientutils.IndentJson(body)))
	}

	log.Debug("Artifactory response:", resp.Status)
	log.Info("Done " + operationString + " repository.")
	return nil
}

func (vrs *VirtualRepositoryService) Maven(params MavenGradleVirtualRepositoryParams) error {
	return vrs.performRequest(params, params.Key)
}

func (vrs *VirtualRepositoryService) Gradle(params MavenGradleVirtualRepositoryParams) error {
	return vrs.performRequest(params, params.Key)
}

func (vrs *VirtualRepositoryService) Ivy(params IvyVirtualRepositoryParam) error {
	return vrs.performRequest(params, params.Key)
}

func (vrs *VirtualRepositoryService) Sbt(params SbtVirtualRepositoryParam) error {
	return vrs.performRequest(params, params.Key)
}

func (vrs *VirtualRepositoryService) Helm(params HelmVirtualRepositoryParam) error {
	return vrs.performRequest(params, params.Key)
}

func (vrs *VirtualRepositoryService) Rpm(params RpmVirtualRepositoryParam) error {
	return vrs.performRequest(params, params.Key)
}

func (vrs *VirtualRepositoryService) Nuget(params NugetVirtualRepositoryParam) error {
	return vrs.performRequest(params, params.Key)
}

func (vrs *VirtualRepositoryService) Cran(params CranVirtualRepositoryParam) error {
	return vrs.performRequest(params, params.Key)
}

func (vrs *VirtualRepositoryService) Gems(params GemsVirtualRepositoryParam) error {
	return vrs.performRequest(params, params.Key)
}

func (vrs *VirtualRepositoryService) Npm(params NpmVirtualRepositoryParam) error {
	return vrs.performRequest(params, params.Key)
}

func (vrs *VirtualRepositoryService) Bower(params BowerVirtualRepositoryParam) error {
	return vrs.performRequest(params, params.Key)
}

func (vrs *VirtualRepositoryService) Debian(params DebianVirtualRepositoryParam) error {
	return vrs.performRequest(params, params.Key)
}

func (vrs *VirtualRepositoryService) Pypi(params PypiVirtualRepositoryParam) error {
	return vrs.performRequest(params, params.Key)
}

func (vrs *VirtualRepositoryService) Docker(params DockerVirtualRepositoryParam) error {
	return vrs.performRequest(params, params.Key)
}

func (vrs *VirtualRepositoryService) Yum(params YumVirtualRepositoryParam) error {
	return vrs.performRequest(params, params.Key)
}

func (vrs *VirtualRepositoryService) Go(params GoVirtualRepositoryParam) error {
	return vrs.performRequest(params, params.Key)
}

func (vrs *VirtualRepositoryService) P2(params P2VirtualRepositoryParam) error {
	return vrs.performRequest(params, params.Key)
}

func (vrs *VirtualRepositoryService) Chef(params ChefVirtualRepositoryParam) error {
	return vrs.performRequest(params, params.Key)
}

func (vrs *VirtualRepositoryService) Puppet(params PuppetVirtualRepositoryParam) error {
	return vrs.performRequest(params, params.Key)
}

func (vrs *VirtualRepositoryService) Conda(params CondaVirtualRepositoryParam) error {
	return vrs.performRequest(params, params.Key)
}

func (vrs *VirtualRepositoryService) Conan(params ConanVirtualRepositoryParam) error {
	return vrs.performRequest(params, params.Key)
}

func (vrs *VirtualRepositoryService) Gitlfs(params GitlfsVirtualRepositoryParam) error {
	return vrs.performRequest(params, params.Key)
}

func (vrs *VirtualRepositoryService) Generic(params GenericVirtualRepositoryParam) error {
	return vrs.performRequest(params, params.Key)
}

type VirtualRepositoryBaseParams struct {
	Key                                           string   `json:"key,omitempty"`
	Rclass                                        string   `json:"rclass"`
	PackageType                                   string   `json:"packageType,omitempty"`
	Description                                   string   `json:"description,omitempty"`
	Notes                                         string   `json:"notes,omitempty"`
	IncludesPattern                               string   `json:"includesPattern,omitempty"`
	ExcludesPattern                               string   `json:"excludesPattern,omitempty"`
	RepoLayoutRef                                 string   `json:"repoLayoutRef, omitempty"`
	Repositories                                  []string `json:"repositories, omitempty"`
	ArtifactoryRequestsCanRetrieveRemoteArtifacts bool     `json:"artifactoryRequestsCanRetrieveRemoteArtifacts, omitempty"`
	DefaultDeploymentRepo                         string   `json:"defaultDeploymentRepo, omitempty"`
}

type MavenGradleVirtualRepositoryParams struct {
	VirtualRepositoryBaseParams
	ForceMavenAuthentication             bool   `json:"forceMavenAuthentication, omitempty"`
	PomRepositoryReferencesCleanupPolicy string `json:"pomRepositoryReferencesCleanupPolicy, omitempty"`
	KeyPair                              string `json:"keyPair, omitempty"`
}

type NugetVirtualRepositoryParam struct {
	VirtualRepositoryBaseParams
	ForceNugetAuthentication bool `json:"forceNugetAuthentication ,omitempty"`
}

type NpmVirtualRepositoryParam struct {
	RemoteRepositoryBaseParams
	ExternalDependenciesEnabled     bool     `json:"externalDependenciesEnabled, omitempty"`
	ExternalDependenciesPatterns    []string `json:"externalDependenciesPatterns, omitempty"`
	ExternalDependenciesRemoteRepo  string   `json:"externalDependenciesRemoteRepo, omitempty"`
	VirtualRetrievalCachePeriodSecs int      `json:"virtualRetrievalCachePeriodSecs, omitempty"`
}

type BowerVirtualRepositoryParam struct {
	RemoteRepositoryBaseParams
	ExternalDependenciesEnabled    bool     `json:"externalDependenciesEnabled, omitempty"`
	ExternalDependenciesPatterns   []string `json:"externalDependenciesPatterns, omitempty"`
	ExternalDependenciesRemoteRepo string   `json:"externalDependenciesRemoteRepo, omitempty"`
}

type DebianVirtualRepositoryParam struct {
	RemoteRepositoryBaseParams
	DebianTrivialLayout bool `json:"debianTrivialLayout, omitempty"`
}

type GoVirtualRepositoryParam struct {
	RemoteRepositoryBaseParams
	ExternalDependenciesEnabled  bool     `json:"externalDependenciesEnabled, omitempty"`
	ExternalDependenciesPatterns []string `json:"externalDependenciesPatterns, omitempty"`
}

type ConanVirtualRepositoryParam struct {
	VirtualRepositoryBaseParams
	VirtualRetrievalCachePeriodSecs int `json:"virtualRetrievalCachePeriodSecs, omitempty"`
}

type HelmVirtualRepositoryParam struct {
	VirtualRepositoryBaseParams
	VirtualRetrievalCachePeriodSecs int `json:"virtualRetrievalCachePeriodSecs, omitempty"`
}

type RpmVirtualRepositoryParam struct {
	VirtualRepositoryBaseParams
	VirtualRetrievalCachePeriodSecs int `json:"virtualRetrievalCachePeriodSecs, omitempty"`
}

type CranVirtualRepositoryParam struct {
	VirtualRepositoryBaseParams
	VirtualRetrievalCachePeriodSecs int `json:"virtualRetrievalCachePeriodSecs, omitempty"`
}

type ChefVirtualRepositoryParam struct {
	VirtualRepositoryBaseParams
	VirtualRetrievalCachePeriodSecs int `json:"virtualRetrievalCachePeriodSecs, omitempty"`
}

type CondaVirtualRepositoryParam struct {
	VirtualRepositoryBaseParams
}

type GitlfsVirtualRepositoryParam struct {
	VirtualRepositoryBaseParams
}

type P2VirtualRepositoryParam struct {
	VirtualRepositoryBaseParams
}

type GemsVirtualRepositoryParam struct {
	VirtualRepositoryBaseParams
}

type PypiVirtualRepositoryParam struct {
	VirtualRepositoryBaseParams
}

type PuppetVirtualRepositoryParam struct {
	VirtualRepositoryBaseParams
}

type IvyVirtualRepositoryParam struct {
	VirtualRepositoryBaseParams
}

type SbtVirtualRepositoryParam struct {
	VirtualRepositoryBaseParams
}

type DockerVirtualRepositoryParam struct {
	VirtualRepositoryBaseParams
}

type YumVirtualRepositoryParam struct {
	VirtualRepositoryBaseParams
}

type GenericVirtualRepositoryParam struct {
	VirtualRepositoryBaseParams
}
