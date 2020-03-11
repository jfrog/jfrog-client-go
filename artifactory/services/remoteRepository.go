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

type RemoteRepositoryService struct {
	isUpdate   bool
	client     *rthttpclient.ArtifactoryHttpClient
	ArtDetails auth.ArtifactoryDetails
}

func NewRemoteRepositoryService(client *rthttpclient.ArtifactoryHttpClient, isUpdate bool) *RemoteRepositoryService {
	return &RemoteRepositoryService{client: client, isUpdate: isUpdate}
}

func (rrs *RemoteRepositoryService) GetJfrogHttpClient() *rthttpclient.ArtifactoryHttpClient {
	return rrs.client
}

func (rrs *RemoteRepositoryService) performRequest(params interface{}, repoKey string) error {
	content, err := json.Marshal(params)
	if errorutils.CheckError(err) != nil {
		return err
	}
	httpClientsDetails := rrs.ArtDetails.CreateHttpClientDetails()
	utils.SetContentType("application/vnd.org.jfrog.artifactory.repositories.RemoteRepositoryConfiguration+json", &httpClientsDetails.Headers)
	var url = rrs.ArtDetails.GetUrl() + "api/repositories/" + repoKey
	var operationString string
	var resp *http.Response
	var body []byte
	if rrs.isUpdate {
		log.Info("Updating remote repository......")
		operationString = "updating"
		resp, body, err = rrs.client.SendPost(url, content, &httpClientsDetails)
	} else {
		log.Info("Creating remote repository......")
		operationString = "creating"
		resp, body, err = rrs.client.SendPut(url, content, &httpClientsDetails)
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

func (rrs *RemoteRepositoryService) Maven(params MavenGradleRemoteRepositoryParams) error {
	return rrs.performRequest(params, params.Key)
}

func (rrs *RemoteRepositoryService) Gradle(params MavenGradleRemoteRepositoryParams) error {
	return rrs.performRequest(params, params.Key)
}

func (rrs *RemoteRepositoryService) Ivy(params IvyRemoteRepositoryParam) error {
	return rrs.performRequest(params, params.Key)
}

func (rrs *RemoteRepositoryService) Sbt(params SbtRemoteRepositoryParam) error {
	return rrs.performRequest(params, params.Key)
}

func (rrs *RemoteRepositoryService) Helm(params HelmRemoteRepositoryParam) error {
	return rrs.performRequest(params, params.Key)
}

func (rrs *RemoteRepositoryService) Cocoapods(params CocoapodsRemoteRepositoryParam) error {
	return rrs.performRequest(params, params.Key)
}

func (rrs *RemoteRepositoryService) Opkg(params OpkgRemoteRepositoryParam) error {
	return rrs.performRequest(params, params.Key)
}

func (rrs *RemoteRepositoryService) Rpm(params RpmRemoteRepositoryParam) error {
	return rrs.performRequest(params, params.Key)
}

func (rrs *RemoteRepositoryService) Nuget(params NugetRemoteRepositoryParam) error {
	return rrs.performRequest(params, params.Key)
}

func (rrs *RemoteRepositoryService) Cran(params CranRemoteRepositoryParam) error {
	return rrs.performRequest(params, params.Key)
}

func (rrs *RemoteRepositoryService) Gems(params GemsRemoteRepositoryParam) error {
	return rrs.performRequest(params, params.Key)
}

func (rrs *RemoteRepositoryService) Npm(params NpmRemoteRepositoryParam) error {
	return rrs.performRequest(params, params.Key)
}

func (rrs *RemoteRepositoryService) Bower(params BowerRemoteRepositoryParam) error {
	return rrs.performRequest(params, params.Key)
}

func (rrs *RemoteRepositoryService) Debian(params DebianRemoteRepositoryParam) error {
	return rrs.performRequest(params, params.Key)
}

func (rrs *RemoteRepositoryService) Pypi(params PypiRemoteRepositoryParam) error {
	return rrs.performRequest(params, params.Key)
}

func (rrs *RemoteRepositoryService) Docker(params DockerRemoteRepositoryParam) error {
	return rrs.performRequest(params, params.Key)
}

func (rrs *RemoteRepositoryService) Yum(params YumRemoteRepositoryParam) error {
	return rrs.performRequest(params, params.Key)
}

func (rrs *RemoteRepositoryService) Vcs(params VcsRemoteRepositoryParam) error {
	return rrs.performRequest(params, params.Key)
}

func (rrs *RemoteRepositoryService) Composer(params ComposerRemoteRepositoryParam) error {
	return rrs.performRequest(params, params.Key)
}

func (rrs *RemoteRepositoryService) Go(params GoRemoteRepositoryParam) error {
	return rrs.performRequest(params, params.Key)
}

func (rrs *RemoteRepositoryService) P2(params P2RemoteRepositoryParam) error {
	return rrs.performRequest(params, params.Key)
}

func (rrs *RemoteRepositoryService) Chef(params ChefRemoteRepositoryParam) error {
	return rrs.performRequest(params, params.Key)
}

func (rrs *RemoteRepositoryService) Puppet(params PuppetRemoteRepositoryParam) error {
	return rrs.performRequest(params, params.Key)
}

func (rrs *RemoteRepositoryService) Conda(params CondaRemoteRepositoryParam) error {
	return rrs.performRequest(params, params.Key)
}

func (rrs *RemoteRepositoryService) Conan(params ConanRemoteRepositoryParam) error {
	return rrs.performRequest(params, params.Key)
}

func (rrs *RemoteRepositoryService) Gitlfs(params GitlfsRemoteRepositoryParam) error {
	return rrs.performRequest(params, params.Key)
}

func (rrs *RemoteRepositoryService) Generic(params GenericRemoteRepositoryParam) error {
	return rrs.performRequest(params, params.Key)
}

type ContentSynchronisation struct {
	Enabled    bool `json:"enables,omitempty"`
	Statistics struct {
		Enabled bool `json:"enables,omitempty"`
	} `json:"statistics,omitempty"`
	Properties struct {
		Enabled bool `json:"enables,omitempty"`
	} `json:"properties,omitempty"`
	Source struct {
		OriginAbsenceDetection bool `json:"originAbsenceDetection,omitempty"`
	} `json:"source,omitempty"`
}

type RemoteRepositoryBaseParams struct {
	Key                               string                 `json:"key,omitempty"`
	Rclass                            string                 `json:"rclass"`
	PackageType                       string                 `json:"packageType,omitempty"`
	Url                               string                 `json:"url"`
	Username                          string                 `json:"username,omitempty"`
	Password                          string                 `json:"password,omitempty"`
	Proxy                             string                 `json:"proxy,omitempty"`
	Description                       string                 `json:"description,omitempty"`
	Notes                             string                 `json:"notes,omitempty"`
	IncludesPattern                   string                 `json:"includesPattern,omitempty"`
	ExcludesPattern                   string                 `json:"excludesPattern,omitempty"`
	RepoLayoutRef                     string                 `json:"repoLayoutRef, omitempty"`
	HardFail                          bool                   `json:"hardFail, omitempty"`
	Offline                           bool                   `json:"offline, omitempty"`
	BlackedOut                        bool                   `json:"blackedOut, omitempty"`
	StoreArtifactsLocally             bool                   `json:"storeArtifactsLocally, omitempty"`
	SocketTimeoutMillis               int                    `json:"socketTimeoutMillis, omitempty"`
	LocalAddress                      string                 `json:"localAddress, omitempty"`
	RetrievalCachePeriodSecs          int                    `json:"retrievalCachePeriodSecs, omitempty"`
	FailedRetrievalCachePeriodSecs    int                    `json:"failedRetrievalCachePeriodSecs, omitempty"`
	MissedRetrievalCachePeriodSecs    int                    `json:"missedRetrievalCachePeriodSecs, omitempty"`
	UnusedArtifactsCleanupEnabled     bool                   `json:"unusedArtifactsCleanupEnabled, omitempty"`
	UnusedArtifactsCleanupPeriodHours int                    `json:"unusedArtifactsCleanupPeriodHours, omitempty"`
	AssumedOfflinePeriodSecs          int                    `json:"assumedOfflinePeriodSecs, omitempty"`
	ShareConfiguration                bool                   `json:"shareConfiguration, omitempty"`
	SynchronizeProperties             bool                   `json:"synchronizeProperties, omitempty"`
	BlockMismatchingMimeTypes         bool                   `json:"blockMismatchingMimeTypes, omitempty"`
	PropertySets                      []string               `json:"propertySets, omitempty"`
	AllowAnyHostAuth                  bool                   `json:""allowAnyHostAuth", omitempty"`
	EnableCookieManagement            bool                   `json:"enableCookieManagement, omitempty"`
	BypassHeadRequests                bool                   `json:"bypassHeadRequests, omitempty"`
	ClientTlsCertificate              string                 `json:"clientTlsCertificate, omitempty"`
	BlockPushingSchema1               bool                   `json:"blockPushingSchema1, omitempty"`
	contentSynchronisation            ContentSynchronisation `json:"contentSynchronisation, omitempty"`
}

type MavenGradleRemoteRepositoryParams struct {
	RemoteRepositoryBaseParams
	FetchJarsEagerly             bool   `json:"fetchJarsEagerly, omitempty"`
	FetchSourcesEagerly          bool   `json:"fetchSourcesEagerly, omitempty"`
	RemoteRepoChecksumPolicyType string `json:"remoteRepoChecksumPolicyType, omitempty"`
	ListRemoteFolderItems        bool   `json:"listRemoteFolderItems, omitempty"`
	HandleReleases               bool   `json:"handleReleases, omitempty"`
	HandleSnapshot               bool   `json:"handleSnapshot, omitempty"`
	SuppressPomConsistencyChecks bool   `json:"suppressPomConsistencyChecks, omitempty"`
	RejectInvalidJars            bool   `json:"rejectInvalidJars, omitempty"`
}

type CocoapodsRemoteRepositoryParam struct {
	RemoteRepositoryBaseParams
	PodsSpecsRepoUrl string `json:"podsSpecsRepoUrl, omitempty"`
}

type OpkgRemoteRepositoryParam struct {
	RemoteRepositoryBaseParams
	ListRemoteFolderItems bool `json:"listRemoteFolderItems, omitempty"`
}

type RpmRemoteRepositoryParam struct {
	RemoteRepositoryBaseParams
	ListRemoteFolderItems bool `json:"listRemoteFolderItems, omitempty"`
}

type NugetRemoteRepositoryParam struct {
	RemoteRepositoryBaseParams
	FeedContextPath          string `json:"feedContextPath, omitempty"`
	DownloadContextPath      string `json:"downloadContextPath, omitempty"`
	V3FeedUrl                string `json:"v3FeedUrl, omitempty"`
	ForceNugetAuthentication bool   `json:"forceNugetAuthentication ,omitempty"`
}

type GemsRemoteRepositoryParam struct {
	RemoteRepositoryBaseParams
	ListRemoteFolderItems bool `json:"listRemoteFolderItems, omitempty"`
}

type NpmRemoteRepositoryParam struct {
	RemoteRepositoryBaseParams
	ListRemoteFolderItems bool `json:"listRemoteFolderItems, omitempty"`
}

type BowerRemoteRepositoryParam struct {
	RemoteRepositoryBaseParams
	BowerRegistryUrl string `json:"bowerRegistryUrl, omitempty"`
}

type DebianRemoteRepositoryParam struct {
	RemoteRepositoryBaseParams
	ListRemoteFolderItems bool `json:"listRemoteFolderItems, omitempty"`
}

type ComposerRemoteRepositoryParam struct {
	RemoteRepositoryBaseParams
	composerRegistryUrl string `json:"composerRegistryUrl, omitempty"`
}

type PypiRemoteRepositoryParam struct {
	RemoteRepositoryBaseParams
	ListRemoteFolderItems bool   `json:"listRemoteFolderItems, omitempty"`
	PypiRegistryUrl       string `json:"pypiRegistryUrl, omitempty"`
}

type DockerRemoteRepositoryParam struct {
	RemoteRepositoryBaseParams
	ExternalDependenciesEnabled  bool     `json:"externalDependenciesEnabled, omitempty"`
	ExternalDependenciesPatterns []string `json:"externalDependenciesPatterns, omitempty"`
	EnableTokenAuthentication    bool     `json:"enableTokenAuthentication, omitempty"`
}

type GitlfsRemoteRepositoryParam struct {
	RemoteRepositoryBaseParams
	ListRemoteFolderItems bool `json:"listRemoteFolderItems, omitempty"`
}

type VcsRemoteRepositoryParam struct {
	RemoteRepositoryBaseParams
	VcsGitProvider        string `json:"vcsGitProvider, omitempty"`
	VcsType               string `json:"vcsType, omitempty"`
	MaxUniqueSnapshots    int    `json:"maxUniqueSnapshots, omitempty"`
	VcsGitDownloadUrl     string `json:"vcsGitDownloadUrl, omitempty"`
	ListRemoteFolderItems bool   `json:"listRemoteFolderItems, omitempty"`
}

type GenericRemoteRepositoryParam struct {
	RemoteRepositoryBaseParams
	ListRemoteFolderItems bool `json:"listRemoteFolderItems, omitempty"`
}

type IvyRemoteRepositoryParam struct {
	RemoteRepositoryBaseParams
}

type SbtRemoteRepositoryParam struct {
	RemoteRepositoryBaseParams
}

type HelmRemoteRepositoryParam struct {
	RemoteRepositoryBaseParams
}

type CranRemoteRepositoryParam struct {
	RemoteRepositoryBaseParams
}

type GoRemoteRepositoryParam struct {
	RemoteRepositoryBaseParams
}

type YumRemoteRepositoryParam struct {
	RemoteRepositoryBaseParams
}

type P2RemoteRepositoryParam struct {
	RemoteRepositoryBaseParams
}

type ChefRemoteRepositoryParam struct {
	RemoteRepositoryBaseParams
}

type PuppetRemoteRepositoryParam struct {
	RemoteRepositoryBaseParams
}

type CondaRemoteRepositoryParam struct {
	RemoteRepositoryBaseParams
}

type ConanRemoteRepositoryParam struct {
	RemoteRepositoryBaseParams
}
