package services

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type FederatedRepositoryService struct {
	isUpdate   bool
	client     *jfroghttpclient.JfrogHttpClient
	ArtDetails auth.ServiceDetails
}

func NewFederatedRepositoryService(client *jfroghttpclient.JfrogHttpClient, isUpdate bool) *FederatedRepositoryService {
	return &FederatedRepositoryService{client: client, isUpdate: isUpdate}
}

func (frs *FederatedRepositoryService) GetJfrogHttpClient() *jfroghttpclient.JfrogHttpClient {
	return frs.client
}

func (frs *FederatedRepositoryService) performRequest(params interface{}, repoKey string) error {
	content, err := json.Marshal(params)
	if errorutils.CheckError(err) != nil {
		return err
	}
	httpClientsDetails := frs.ArtDetails.CreateHttpClientDetails()
	utils.SetContentType("application/vnd.org.jfrog.artifactory.repositories.FederatedRepositoryConfiguration+json", &httpClientsDetails.Headers)
	var url = frs.ArtDetails.GetUrl() + "api/repositories/" + repoKey
	var operationString string
	var resp *http.Response
	var body []byte
	if frs.isUpdate {
		log.Info("Updating federated repository...")
		operationString = "updating"
		resp, body, err = frs.client.SendPost(url, content, &httpClientsDetails)
	} else {
		log.Info("Creating federated repository...")
		operationString = "creating"
		resp, body, err = frs.client.SendPut(url, content, &httpClientsDetails)
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

func (frs *FederatedRepositoryService) Alpine(params AlpineFederatedRepositoryParams) error {
	return frs.performRequest(params, params.Key)
}

func (frs *FederatedRepositoryService) Bower(params BowerFederatedRepositoryParams) error {
	return frs.performRequest(params, params.Key)
}

func (frs *FederatedRepositoryService) Cran(params CranFederatedRepositoryParams) error {
	return frs.performRequest(params, params.Key)
}

func (frs *FederatedRepositoryService) Cargo(params CargoFederatedRepositoryParams) error {
	return frs.performRequest(params, params.Key)
}

func (frs *FederatedRepositoryService) Chef(params ChefFederatedRepositoryParams) error {
	return frs.performRequest(params, params.Key)
}

func (frs *FederatedRepositoryService) Cocoapods(params CocoapodsFederatedRepositoryParams) error {
	return frs.performRequest(params, params.Key)
}

func (frs *FederatedRepositoryService) Composer(params ComposerFederatedRepositoryParams) error {
	return frs.performRequest(params, params.Key)
}

func (frs *FederatedRepositoryService) Conan(params ConanFederatedRepositoryParams) error {
	return frs.performRequest(params, params.Key)
}

func (frs *FederatedRepositoryService) Conda(params CondaFederatedRepositoryParams) error {
	return frs.performRequest(params, params.Key)
}

func (frs *FederatedRepositoryService) Debian(params DebianFederatedRepositoryParams) error {
	return frs.performRequest(params, params.Key)
}

func (frs *FederatedRepositoryService) Docker(params DockerFederatedRepositoryParams) error {
	return frs.performRequest(params, params.Key)
}

func (frs *FederatedRepositoryService) Gems(params GemsFederatedRepositoryParams) error {
	return frs.performRequest(params, params.Key)
}

func (frs *FederatedRepositoryService) Generic(params GenericFederatedRepositoryParams) error {
	return frs.performRequest(params, params.Key)
}

func (frs *FederatedRepositoryService) Gitlfs(params GitlfsFederatedRepositoryParams) error {
	return frs.performRequest(params, params.Key)
}

func (frs *FederatedRepositoryService) Go(params GoFederatedRepositoryParams) error {
	return frs.performRequest(params, params.Key)
}

func (frs *FederatedRepositoryService) Gradle(params GradleFederatedRepositoryParams) error {
	return frs.performRequest(params, params.Key)
}

func (frs *FederatedRepositoryService) Helm(params HelmFederatedRepositoryParams) error {
	return frs.performRequest(params, params.Key)
}

func (frs *FederatedRepositoryService) Ivy(params IvyFederatedRepositoryParams) error {
	return frs.performRequest(params, params.Key)
}

func (frs *FederatedRepositoryService) Maven(params MavenFederatedRepositoryParams) error {
	return frs.performRequest(params, params.Key)
}

func (frs *FederatedRepositoryService) Npm(params NpmFederatedRepositoryParams) error {
	return frs.performRequest(params, params.Key)
}

func (frs *FederatedRepositoryService) Nuget(params NugetFederatedRepositoryParams) error {
	return frs.performRequest(params, params.Key)
}

func (frs *FederatedRepositoryService) Opkg(params OpkgFederatedRepositoryParams) error {
	return frs.performRequest(params, params.Key)
}

func (frs *FederatedRepositoryService) Puppet(params PuppetFederatedRepositoryParams) error {
	return frs.performRequest(params, params.Key)
}

func (frs *FederatedRepositoryService) Pypi(params PypiFederatedRepositoryParams) error {
	return frs.performRequest(params, params.Key)
}

func (frs *FederatedRepositoryService) Rpm(params RpmFederatedRepositoryParams) error {
	return frs.performRequest(params, params.Key)
}

func (frs *FederatedRepositoryService) Sbt(params SbtFederatedRepositoryParams) error {
	return frs.performRequest(params, params.Key)
}

func (frs *FederatedRepositoryService) Vagrant(params VagrantFederatedRepositoryParams) error {
	return frs.performRequest(params, params.Key)
}

func (frs *FederatedRepositoryService) Yum(params YumFederatedRepositoryParams) error {
	return frs.performRequest(params, params.Key)
}

type FederatedRepositoryMemberParams struct {
	Url     string `json:"url"`
	Enabled *bool  `json:"enabled,omitempty"`
}

type FederatedRepositoryBaseParams struct {
	Key                    string                            `json:"key,omitempty"`
	Rclass                 string                            `json:"rclass"`
	PackageType            string                            `json:"packageType,omitempty"`
	Description            string                            `json:"description,omitempty"`
	Notes                  string                            `json:"notes,omitempty"`
	IncludesPattern        string                            `json:"includesPattern,omitempty"`
	ExcludesPattern        string                            `json:"excludesPattern,omitempty"`
	RepoLayoutRef          string                            `json:"repoLayoutRef,omitempty"`
	BlackedOut             *bool                             `json:"blackedOut,omitempty"`
	XrayIndex              *bool                             `json:"xrayIndex,omitempty"`
	PropertySets           []string                          `json:"propertySets,omitempty"`
	ArchiveBrowsingEnabled *bool                             `json:"archiveBrowsingEnabled,omitempty"`
	DownloadRedirect       *bool                             `json:"downloadRedirect,omitempty"`
	PriorityResolution     *bool                             `json:"priorityResolution,omitempty"`
	Members                []FederatedRepositoryMemberParams `json:"members,omitempty"`
}

func NewFederatedRepositoryBaseParams() FederatedRepositoryBaseParams {
	return FederatedRepositoryBaseParams{Rclass: "federated"}
}

type AlpineFederatedRepositoryParams struct {
	FederatedRepositoryBaseParams
}

func NewAlpineFederatedRepositoryParams() AlpineFederatedRepositoryParams {
	return AlpineFederatedRepositoryParams{FederatedRepositoryBaseParams{Rclass: "federated", PackageType: "alpine"}}
}

type BowerFederatedRepositoryParams struct {
	FederatedRepositoryBaseParams
}

func NewBowerFederatedRepositoryParams() BowerFederatedRepositoryParams {
	return BowerFederatedRepositoryParams{FederatedRepositoryBaseParams{Rclass: "federated", PackageType: "bower"}}
}

type CranFederatedRepositoryParams struct {
	FederatedRepositoryBaseParams
}

func NewCranFederatedRepositoryParams() CranFederatedRepositoryParams {
	return CranFederatedRepositoryParams{FederatedRepositoryBaseParams{Rclass: "federated", PackageType: "cran"}}
}

type CargoFederatedRepositoryParams struct {
	FederatedRepositoryBaseParams
	CargoAnonymousAccess *bool `json:"cargoAnonymousAccess,omitempty"`
}

func NewCargoFederatedRepositoryParams() CargoFederatedRepositoryParams {
	return CargoFederatedRepositoryParams{FederatedRepositoryBaseParams: FederatedRepositoryBaseParams{Rclass: "federated", PackageType: "cargo"}}
}

type ChefFederatedRepositoryParams struct {
	FederatedRepositoryBaseParams
}

func NewChefFederatedRepositoryParams() ChefFederatedRepositoryParams {
	return ChefFederatedRepositoryParams{FederatedRepositoryBaseParams{Rclass: "federated", PackageType: "chef"}}
}

type CocoapodsFederatedRepositoryParams struct {
	FederatedRepositoryBaseParams
}

func NewCocoapodsFederatedRepositoryParams() CocoapodsFederatedRepositoryParams {
	return CocoapodsFederatedRepositoryParams{FederatedRepositoryBaseParams{Rclass: "federated", PackageType: "cocoapods"}}
}

type ComposerFederatedRepositoryParams struct {
	FederatedRepositoryBaseParams
}

func NewComposerFederatedRepositoryParams() ComposerFederatedRepositoryParams {
	return ComposerFederatedRepositoryParams{FederatedRepositoryBaseParams{Rclass: "federated", PackageType: "composer"}}
}

type ConanFederatedRepositoryParams struct {
	FederatedRepositoryBaseParams
}

func NewConanFederatedRepositoryParams() ConanFederatedRepositoryParams {
	return ConanFederatedRepositoryParams{FederatedRepositoryBaseParams{Rclass: "federated", PackageType: "conan"}}
}

type CondaFederatedRepositoryParams struct {
	FederatedRepositoryBaseParams
}

func NewCondaFederatedRepositoryParams() CondaFederatedRepositoryParams {
	return CondaFederatedRepositoryParams{FederatedRepositoryBaseParams{Rclass: "federated", PackageType: "conda"}}
}

type DebianFederatedRepositoryParams struct {
	FederatedRepositoryBaseParams
	DebianTrivialLayout             *bool    `json:"debianTrivialLayout,omitempty"`
	OptionalIndexCompressionFormats []string `json:"optionalIndexCompressionFormats,omitempty"`
}

func NewDebianFederatedRepositoryParams() DebianFederatedRepositoryParams {
	return DebianFederatedRepositoryParams{FederatedRepositoryBaseParams: FederatedRepositoryBaseParams{Rclass: "federated", PackageType: "debian"}}
}

type DockerFederatedRepositoryParams struct {
	FederatedRepositoryBaseParams
	MaxUniqueTags       int    `json:"maxUniqueTags,omitempty"`
	DockerApiVersion    string `json:"dockerApiVersion,omitempty"`
	BlockPushingSchema1 *bool  `json:"blockPushingSchema1,omitempty"`
	DockerTagRetention  int    `json:"dockerTagRetention,omitempty"`
}

func NewDockerFederatedRepositoryParams() DockerFederatedRepositoryParams {
	return DockerFederatedRepositoryParams{FederatedRepositoryBaseParams: FederatedRepositoryBaseParams{Rclass: "federated", PackageType: "docker"}}
}

type GemsFederatedRepositoryParams struct {
	FederatedRepositoryBaseParams
}

func NewGemsFederatedRepositoryParams() GemsFederatedRepositoryParams {
	return GemsFederatedRepositoryParams{FederatedRepositoryBaseParams{Rclass: "federated", PackageType: "gems"}}
}

type GenericFederatedRepositoryParams struct {
	FederatedRepositoryBaseParams
}

func NewGenericFederatedRepositoryParams() GenericFederatedRepositoryParams {
	return GenericFederatedRepositoryParams{FederatedRepositoryBaseParams{Rclass: "federated", PackageType: "generic"}}
}

type GitlfsFederatedRepositoryParams struct {
	FederatedRepositoryBaseParams
}

func NewGitlfsFederatedRepositoryParams() GitlfsFederatedRepositoryParams {
	return GitlfsFederatedRepositoryParams{FederatedRepositoryBaseParams{Rclass: "federated", PackageType: "gitlfs"}}
}

type GoFederatedRepositoryParams struct {
	FederatedRepositoryBaseParams
}

func NewGoFederatedRepositoryParams() GoFederatedRepositoryParams {
	return GoFederatedRepositoryParams{FederatedRepositoryBaseParams{Rclass: "federated", PackageType: "go"}}
}

type GradleFederatedRepositoryParams struct {
	FederatedRepositoryBaseParams
	CommonMavenGradleFederatedRepositoryParams
}

func NewGradleFederatedRepositoryParams() GradleFederatedRepositoryParams {
	return GradleFederatedRepositoryParams{FederatedRepositoryBaseParams: FederatedRepositoryBaseParams{Rclass: "federated", PackageType: "gradle"}}
}

type HelmFederatedRepositoryParams struct {
	FederatedRepositoryBaseParams
}

func NewHelmFederatedRepositoryParams() HelmFederatedRepositoryParams {
	return HelmFederatedRepositoryParams{FederatedRepositoryBaseParams{Rclass: "federated", PackageType: "helm"}}
}

type IvyFederatedRepositoryParams struct {
	FederatedRepositoryBaseParams
	CommonMavenGradleFederatedRepositoryParams
}

func NewIvyFederatedRepositoryParams() IvyFederatedRepositoryParams {
	return IvyFederatedRepositoryParams{FederatedRepositoryBaseParams: FederatedRepositoryBaseParams{Rclass: "federated", PackageType: "ivy"}}
}

type CommonMavenGradleFederatedRepositoryParams struct {
	MaxUniqueSnapshots           int    `json:"maxUniqueSnapshots,omitempty"`
	HandleReleases               *bool  `json:"handleReleases,omitempty"`
	HandleSnapshots              *bool  `json:"handleSnapshots,omitempty"`
	SuppressPomConsistencyChecks *bool  `json:"suppressPomConsistencyChecks,omitempty"`
	SnapshotVersionBehavior      string `json:"snapshotVersionBehavior,omitempty"`
	ChecksumPolicyType           string `json:"checksumPolicyType,omitempty"`
}

type MavenFederatedRepositoryParams struct {
	FederatedRepositoryBaseParams
	CommonMavenGradleFederatedRepositoryParams
}

func NewMavenFederatedRepositoryParams() MavenFederatedRepositoryParams {
	return MavenFederatedRepositoryParams{FederatedRepositoryBaseParams: FederatedRepositoryBaseParams{Rclass: "federated", PackageType: "maven"}}
}

type NpmFederatedRepositoryParams struct {
	FederatedRepositoryBaseParams
}

func NewNpmFederatedRepositoryParams() NpmFederatedRepositoryParams {
	return NpmFederatedRepositoryParams{FederatedRepositoryBaseParams{Rclass: "federated", PackageType: "npm"}}
}

type NugetFederatedRepositoryParams struct {
	FederatedRepositoryBaseParams
	MaxUniqueSnapshots       int   `json:"maxUniqueSnapshots,omitempty"`
	ForceNugetAuthentication *bool `json:"forceNugetAuthentication,omitempty"`
}

func NewNugetFederatedRepositoryParams() NugetFederatedRepositoryParams {
	return NugetFederatedRepositoryParams{FederatedRepositoryBaseParams: FederatedRepositoryBaseParams{Rclass: "federated", PackageType: "nuget"}}
}

type OpkgFederatedRepositoryParams struct {
	FederatedRepositoryBaseParams
}

func NewOpkgFederatedRepositoryParams() OpkgFederatedRepositoryParams {
	return OpkgFederatedRepositoryParams{FederatedRepositoryBaseParams{Rclass: "federated", PackageType: "opkg"}}
}

type PuppetFederatedRepositoryParams struct {
	FederatedRepositoryBaseParams
}

func NewPuppetFederatedRepositoryParams() PuppetFederatedRepositoryParams {
	return PuppetFederatedRepositoryParams{FederatedRepositoryBaseParams{Rclass: "federated", PackageType: "puppet"}}
}

type PypiFederatedRepositoryParams struct {
	FederatedRepositoryBaseParams
}

func NewPypiFederatedRepositoryParams() PypiFederatedRepositoryParams {
	return PypiFederatedRepositoryParams{FederatedRepositoryBaseParams{Rclass: "federated", PackageType: "pypi"}}
}

type RpmFederatedRepositoryParams struct {
	FederatedRepositoryBaseParams
	YumRootDepth            int   `json:"yumRootDepth,omitempty"`
	CalculateYumMetadata    *bool `json:"calculateYumMetadata,omitempty"`
	EnableFileListsIndexing *bool `json:"enableFileListsIndexing,omitempty"`
}

func NewRpmFederatedRepositoryParams() RpmFederatedRepositoryParams {
	return RpmFederatedRepositoryParams{FederatedRepositoryBaseParams: FederatedRepositoryBaseParams{Rclass: "federated", PackageType: "rpm"}}
}

type SbtFederatedRepositoryParams struct {
	FederatedRepositoryBaseParams
	CommonMavenGradleFederatedRepositoryParams
}

func NewSbtFederatedRepositoryParams() SbtFederatedRepositoryParams {
	return SbtFederatedRepositoryParams{FederatedRepositoryBaseParams: FederatedRepositoryBaseParams{Rclass: "federated", PackageType: "sbt"}}
}

type VagrantFederatedRepositoryParams struct {
	FederatedRepositoryBaseParams
}

func NewVagrantFederatedRepositoryParams() VagrantFederatedRepositoryParams {
	return VagrantFederatedRepositoryParams{FederatedRepositoryBaseParams{Rclass: "federated", PackageType: "vagrant"}}
}

type YumFederatedRepositoryParams struct {
	FederatedRepositoryBaseParams
}

func NewYumFederatedRepositoryParams() YumFederatedRepositoryParams {
	return YumFederatedRepositoryParams{FederatedRepositoryBaseParams{Rclass: "federated", PackageType: "yum"}}
}
