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

type LocalRepositoryService struct {
	client     *rthttpclient.ArtifactoryHttpClient
	ArtDetails auth.ArtifactoryDetails
}

func NewLocalRepositoryService(client *rthttpclient.ArtifactoryHttpClient) *LocalRepositoryService {
	return &LocalRepositoryService{client: client}
}

func (lrs *LocalRepositoryService) GetJfrogHttpClient() *rthttpclient.ArtifactoryHttpClient {
	return lrs.client
}

func (lrs *LocalRepositoryService) performRequest(params interface{}, repoKey string) error {
	content, err := json.Marshal(params)
	if errorutils.CheckError(err) != nil {
		return err
	}
	httpClientsDetails := lrs.ArtDetails.CreateHttpClientDetails()
	utils.SetContentType("application/vnd.org.jfrog.artifactory.repositories.LocalRepositoryConfiguration+json", &httpClientsDetails.Headers)
	log.Info("Creating local repository......")
	resp, body, err := lrs.client.SendPut(lrs.ArtDetails.GetUrl()+"api/repositories/"+repoKey, content, &httpClientsDetails)
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

//func (rs *LocalRepositoryService) UpdateRepository(propsParams PropsParams) (int, error) {
//	log.Info("Updating repository...")
//	totalSuccess, err := ps.performRequest(propsParams, false)
//	if err != nil {
//		return totalSuccess, err
//	}
//	log.Info("Done updating repository.")
//	return totalSuccess, nil
//}

//func (rs *LocalRepositoryService) DeleteRepository(propsParams PropsParams) (int, error) {
//	log.Info("Deleting repository...")
//	totalSuccess, err := ps.performRequest(propsParams, true)
//	if err != nil {
//		return totalSuccess, err
//	}
//	log.Info("Done deleting repository.")
//	return totalSuccess, nil
//}

func (lrs *LocalRepositoryService) Maven(params MavenGradleLocalRepositoryParams) error {
	return lrs.performRequest(params, params.Key)
}

func (lrs *LocalRepositoryService) Gradle(params MavenGradleLocalRepositoryParams) error {
	return lrs.performRequest(params, params.Key)
}

func (lrs *LocalRepositoryService) Rpm(params RpmLocalRepositoryParam) error {
	return lrs.performRequest(params, params.Key)
}

func (lrs *LocalRepositoryService) Nuget(params NugetLocalRepositoryParam) error {
	return lrs.performRequest(params, params.Key)
}

func (lrs *LocalRepositoryService) Debian(params DebianLocalRepositoryParam) error {
	return lrs.performRequest(params, params.Key)
}

func (lrs *LocalRepositoryService) Docker(params DockerLocalRepositoryParam) error {
	return lrs.performRequest(params, params.Key)
}

func (lrs *LocalRepositoryService) Ivy(params IvyLocalRepositoryParam) error {
	return lrs.performRequest(params, params.Key)
}

func (lrs *LocalRepositoryService) Sbt(params SbtLocalRepositoryParam) error {
	return lrs.performRequest(params, params.Key)
}

func (lrs *LocalRepositoryService) Helm(params HelmLocalRepositoryParam) error {
	return lrs.performRequest(params, params.Key)
}

func (lrs *LocalRepositoryService) Cocapods(params CocapodsLocalRepositoryParam) error {
	return lrs.performRequest(params, params.Key)
}

func (lrs *LocalRepositoryService) Opkg(params OpkgLocalRepositoryParam) error {
	return lrs.performRequest(params, params.Key)
}

func (lrs *LocalRepositoryService) Cran(params CranLocalRepositoryParam) error {
	return lrs.performRequest(params, params.Key)
}

func (lrs *LocalRepositoryService) Gems(params GemsLocalRepositoryParam) error {
	return lrs.performRequest(params, params.Key)
}

func (lrs *LocalRepositoryService) Npm(params NpmLocalRepositoryParam) error {
	return lrs.performRequest(params, params.Key)
}

func (lrs *LocalRepositoryService) Bower(params BowerLocalRepositoryParam) error {
	return lrs.performRequest(params, params.Key)
}

func (lrs *LocalRepositoryService) Composer(params ComposerLocalRepositoryParam) error {
	return lrs.performRequest(params, params.Key)
}

func (lrs *LocalRepositoryService) Pypi(params PypiLocalRepositoryParam) error {
	return lrs.performRequest(params, params.Key)
}

func (lrs *LocalRepositoryService) Vagrant(params VagrantLocalRepositoryParam) error {
	return lrs.performRequest(params, params.Key)
}

func (lrs *LocalRepositoryService) Gitlfs(params GitlfsLocalRepositoryParam) error {
	return lrs.performRequest(params, params.Key)
}

func (lrs *LocalRepositoryService) Go(params GoLocalRepositoryParam) error {
	return lrs.performRequest(params, params.Key)
}

func (lrs *LocalRepositoryService) Yum(params YumLocalRepositoryParam) error {
	return lrs.performRequest(params, params.Key)
}

func (lrs *LocalRepositoryService) Conan(params ConanLocalRepositoryParam) error {
	return lrs.performRequest(params, params.Key)
}

func (lrs *LocalRepositoryService) Chef(params ChefLocalRepositoryParam) error {
	return lrs.performRequest(params, params.Key)
}

func (lrs *LocalRepositoryService) Puppet(params PuppetLocalRepositoryParam) error {
	return lrs.performRequest(params, params.Key)
}

func (lrs *LocalRepositoryService) Generic(params GenericLocalRepositoryParam) error {
	return lrs.performRequest(params, params.Key)
}

type LocalRepositoryBaseParams struct {
	Key                             string   `json:"key,omitempty"`
	Rclass                          string   `json:"rclass"`
	PackageType                     string   `json:"packageType,omitempty"`
	Description                     string   `json:"description,omitempty"`
	Notes                           string   `json:"notes,omitempty"`
	IncludesPattern                 string   `json:"includesPattern,omitempty"`
	ExcludesPattern                 string   `json:"excludesPattern,omitempty"`
	RepoLayoutRef                   string   `json:"repoLayoutRef, omitempty"`
	BlackedOut                      bool     `json:"blackedOut, omitempty"`
	XrayIndex                       bool     `json:"xrayIndex, omitempty"`
	PropertySets                    []string `json:"propertySets, omitempty"`
	ArchiveBrowsingEnabled          bool     `json:"archiveBrowsingEnabled, omitempty"`
	OptionalIndexCompressionFormats []string `json:"optionalIndexCompressionFormats, omitempty"`
	DownloadRedirect                bool     `json:"downloadRedirect, omitempty"`
	BlockPushingSchema1             bool     `json:"blockPushingSchema1, omitempty"`
}

type MavenGradleLocalRepositoryParams struct {
	LocalRepositoryBaseParams
	MaxUniqueSnapshots           int    `json:"maxUniqueSnapshots,omitempty"`
	HandleReleases               bool   `json:"handleReleases,omitempty"`
	HandleSnapshot               bool   `json:"handleSnapshot,omitempty"`
	SuppressPomConsistencyChecks bool   `json:"suppressPomConsistencyChecks,omitempty"`
	SnapshotVersionBehavior      string `json:"snapshotVersionBehavior,omitempty"`
	ChecksumPolicyType           string `json:"checksumPolicyType,omitempty"`
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

type GenericLocalRepositoryParam struct {
	LocalRepositoryBaseParams
}
