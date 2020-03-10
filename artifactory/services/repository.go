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

type RepositoryService struct {
	client     *rthttpclient.ArtifactoryHttpClient
	ArtDetails auth.ArtifactoryDetails
}

func NewRepositoryService(client *rthttpclient.ArtifactoryHttpClient) *RepositoryService {
	return &RepositoryService{client: client}
}

func (rs *RepositoryService) GetJfrogHttpClient() *rthttpclient.ArtifactoryHttpClient {
	return rs.client
}

func (rs *RepositoryService) performRequest(params interface{}, repoKey string) error {
	content, err := json.Marshal(params)
	if errorutils.CheckError(err) != nil {
		return err
	}
	httpClientsDetails := rs.ArtDetails.CreateHttpClientDetails()
	utils.SetContentType("application/vnd.org.jfrog.artifactory.repositories.LocalRepositoryConfiguration+json", &httpClientsDetails.Headers)
	log.Info("Creating local repository......")
	resp, body, err := rs.client.SendPut(rs.ArtDetails.GetUrl()+"api/repositories/"+repoKey(), content, &httpClientsDetails)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errorutils.CheckError(errors.New("Artifactory response: " + resp.Status + "\n" + clientutils.IndentJson(body)))
	}

	log.Debug("Artifactory response:", resp.Status)
	log.Info("Done creating repository.")
	return nil
}

//func (rs *RepositoryService) UpdateRepository(propsParams PropsParams) (int, error) {
//	log.Info("Updating repository...")
//	totalSuccess, err := ps.performRequest(propsParams, false)
//	if err != nil {
//		return totalSuccess, err
//	}
//	log.Info("Done updating repository.")
//	return totalSuccess, nil
//}

//func (rs *RepositoryService) DeleteRepository(propsParams PropsParams) (int, error) {
//	log.Info("Deleting repository...")
//	totalSuccess, err := ps.performRequest(propsParams, true)
//	if err != nil {
//		return totalSuccess, err
//	}
//	log.Info("Done deleting repository.")
//	return totalSuccess, nil
//}

func (rs *RepositoryService) Maven(params MavenGradleLocalRepositoryParams) error {
	return rs.performRequest(params, params.Key)
}

func (rs *RepositoryService) Gradle(params MavenGradleLocalRepositoryParams) error {
	return rs.performRequest(params, params.Key)
}

func (rs *RepositoryService) Rpm(params RpmLocalRepositoryParam) error {
	return rs.performRequest(params, params.Key)
}

func (rs *RepositoryService) Nuget(params NugetLocalRepositoryParam) error {
	return rs.performRequest(params, params.Key)
}

func (rs *RepositoryService) Debian(params DebianLocalRepositoryParam) error {
	return rs.performRequest(params, params.Key)
}

func (rs *RepositoryService) Docker(params DockerLocalRepositoryParam) error {
	return rs.performRequest(params, params.Key)
}

func (rs *RepositoryService) Ivy(params IvyLocalRepositoryParam) error {
	return rs.performRequest(params, params.Key)
}

func (rs *RepositoryService) Sbt(params SbtLocalRepositoryParam) error {
	return rs.performRequest(params, params.Key)
}

func (rs *RepositoryService) Helm(params HelmLocalRepositoryParam) error {
	return rs.performRequest(params, params.Key)
}

func (rs *RepositoryService) Cocapods(params CocapodsLocalRepositoryParam) error {
	return rs.performRequest(params, params.Key)
}

func (rs *RepositoryService) Opkg(params OpkgLocalRepositoryParam) error {
	return rs.performRequest(params, params.Key)
}

func (rs *RepositoryService) Cran(params CranLocalRepositoryParam) error {
	return rs.performRequest(params, params.Key)
}

func (rs *RepositoryService) Gems(params GemsLocalRepositoryParam) error {
	return rs.performRequest(params, params.Key)
}

func (rs *RepositoryService) Npm(params NpmLocalRepositoryParam) error {
	return rs.performRequest(params, params.Key)
}

func (rs *RepositoryService) Bower(params BowerLocalRepositoryParam) error {
	return rs.performRequest(params, params.Key)
}

func (rs *RepositoryService) Composer(params ComposerLocalRepositoryParam) error {
	return rs.performRequest(params, params.Key)
}

func (rs *RepositoryService) Pypi(params PypiLocalRepositoryParam) error {
	return rs.performRequest(params, params.Key)
}

func (rs *RepositoryService) Vagrant(params VagrantLocalRepositoryParam) error {
	return rs.performRequest(params, params.Key)
}

func (rs *RepositoryService) Gitlfs(params GitlfsLocalRepositoryParam) error {
	return rs.performRequest(params, params.Key)
}

func (rs *RepositoryService) Go(params GoLocalRepositoryParam) error {
	return rs.performRequest(params, params.Key)
}

func (rs *RepositoryService) Yum(params YumLocalRepositoryParam) error {
	return rs.performRequest(params, params.Key)
}

func (rs *RepositoryService) Conan(params ConanLocalRepositoryParam) error {
	return rs.performRequest(params, params.Key)
}

func (rs *RepositoryService) Chef(params ChefLocalRepositoryParam) error {
	return rs.performRequest(params, params.Key)
}

func (rs *RepositoryService) Puppet(params PuppetLocalRepositoryParam) error {
	return rs.performRequest(params, params.Key)
}

func (rs *RepositoryService) Generic(params GenericRepositoryParam) error {
	return rs.performRequest(params, params.Key)
}

type LocalRepositoryBaseParams struct {
	Key                             string   `json:"key,omitempty"`
	Rclass                          string   `json:"rclass"`
	PackageType                     string   `json:"packageTyoe,omitempty"`
	Description                     string   `json:"description,omitempty"`
	Notes                           string   `json:"notes,omitempty"`
	IncludesPattern                 string   `json:"includesPattern,omitempty"`
	ExcludesPattern                 string   `json:"excludesPattern,omitempty"`
	RepoLayoutRef                   string   `json:"repoLayoutRef, omitempty"`
	BlackedOut                      bool     `json:"blackedOut, omitempty"`
	XrayIndex                       bool     `json:"xrayIndex, omitempty"`
	PropertySet                     []string `json:"propertySet, omitempty"`
	ArchiveBrowsingEnabled          bool     `json:"archiveBrowsingEnabled, omitempty"`
	OptionalIndexCompressionFormats []string `json:"optionalIndexCompressionFormats, omitempty"`
	DownloadRedirect                bool     `json:"downloadRedirect, omitempty"`
	BlockPushingSchema1             bool     `json:"blockPushingSchema1, omitempty"`
}

type MavenGradleLocalRepositoryParams struct {
	LocalRepositoryBaseParams
	maxUniqueSnapshots           int    `json:"maxUniqueSnapshots,omitempty"`
	handleReleases               bool   `json:"handleReleases,omitempty"`
	handleSnapshot               bool   `json:"handleSnapshot,omitempty"`
	suppressPomConsistencyChecks bool   `json:"suppressPomConsistencyChecks,omitempty"`
	snapshotVersionBehavior      string `json:"snapshotVersionBehavior,omitempty"`
	checksumPolicyType           string `json:"checksumPolicyType,omitempty"`
}

type RpmLocalRepositoryParam struct {
	LocalRepositoryBaseParams
	YumRootDepth            int  `json:"yumRootDepth,omitempty"`
	CalculateYumMetadata    bool `json:"calculateYumMetadata,omitempty"`
	EnableFileListsIndexing bool `json:"enableFileListsIndexing ,omitempty"`
}

type NugetLocalRepositoryParam struct {
	LocalRepositoryBaseParams
	MaxUniqueSnapshots       int  `json:"maxUniqueSnapshots,omitempty"`
	ForceNugetAuthentication bool `json:"forceNugetAuthentication ,omitempty"`
}

type DebianLocalRepositoryParam struct {
	LocalRepositoryBaseParams
	DebianTrivialLayout bool `json:"debianTrivialLayout ,omitempty"`
}

type DockerLocalRepositoryParam struct {
	LocalRepositoryBaseParams
	MaxUniqueTags    int  `json:"maxUniqueTags,omitempty"`
	DockerApiVersion bool `json:"dockerApiVersion ,omitempty"`
}

type IvyLocalRepositoryParam struct {
	LocalRepositoryBaseParams
}

type SbtLocalRepositoryParam struct {
	LocalRepositoryBaseParams
}

type HelmLocalRepositoryParam struct {
	LocalRepositoryBaseParams
}

type CocapodsLocalRepositoryParam struct {
	LocalRepositoryBaseParams
}

type OpkgLocalRepositoryParam struct {
	LocalRepositoryBaseParams
}

type CranLocalRepositoryParam struct {
	LocalRepositoryBaseParams
}

type GemsLocalRepositoryParam struct {
	LocalRepositoryBaseParams
}

type NpmLocalRepositoryParam struct {
	LocalRepositoryBaseParams
}

type BowerLocalRepositoryParam struct {
	LocalRepositoryBaseParams
}

type ComposerLocalRepositoryParam struct {
	LocalRepositoryBaseParams
}

type PypiLocalRepositoryParam struct {
	LocalRepositoryBaseParams
}

type VagrantLocalRepositoryParam struct {
	LocalRepositoryBaseParams
}

type GitlfsLocalRepositoryParam struct {
	LocalRepositoryBaseParams
}

type GoLocalRepositoryParam struct {
	LocalRepositoryBaseParams
}

type YumLocalRepositoryParam struct {
	LocalRepositoryBaseParams
}

type ConanLocalRepositoryParam struct {
	LocalRepositoryBaseParams
}

type ChefLocalRepositoryParam struct {
	LocalRepositoryBaseParams
}

type PuppetLocalRepositoryParam struct {
	LocalRepositoryBaseParams
}

type GenericRepositoryParam struct {
	LocalRepositoryBaseParams
}
