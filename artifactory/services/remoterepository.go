package services

import (
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
)

const RemoteRepositoryRepoType = "remote"

type RemoteRepositoryService struct {
	RepositoryService
}

func NewRemoteRepositoryService(client *jfroghttpclient.JfrogHttpClient, isUpdate bool) *RemoteRepositoryService {
	return &RemoteRepositoryService{
		RepositoryService: RepositoryService{
			repoType: RemoteRepositoryRepoType,
			client:   client,
			isUpdate: isUpdate,
		},
	}
}

func (rrs *RemoteRepositoryService) GetJfrogHttpClient() *jfroghttpclient.JfrogHttpClient {
	return rrs.client
}

func (rrs *RemoteRepositoryService) performRequest(params interface{}, repoKey string) error {
	return rrs.RepositoryService.performRequest(params, repoKey)
}

func (rrs *RemoteRepositoryService) Alpine(params AlpineRemoteRepositoryParams) error {
	return rrs.performRequest(params, params.Key)
}

func (rrs *RemoteRepositoryService) Bower(params BowerRemoteRepositoryParams) error {
	return rrs.performRequest(params, params.Key)
}

func (rrs *RemoteRepositoryService) Cargo(params CargoRemoteRepositoryParams) error {
	return rrs.performRequest(params, params.Key)
}

func (rrs *RemoteRepositoryService) Chef(params ChefRemoteRepositoryParams) error {
	return rrs.performRequest(params, params.Key)
}

func (rrs *RemoteRepositoryService) Cocoapods(params CocoapodsRemoteRepositoryParams) error {
	return rrs.performRequest(params, params.Key)
}

func (rrs *RemoteRepositoryService) Composer(params ComposerRemoteRepositoryParams) error {
	return rrs.performRequest(params, params.Key)
}

func (rrs *RemoteRepositoryService) Conan(params ConanRemoteRepositoryParams) error {
	return rrs.performRequest(params, params.Key)
}

func (rrs *RemoteRepositoryService) Conda(params CondaRemoteRepositoryParams) error {
	return rrs.performRequest(params, params.Key)
}

func (rrs *RemoteRepositoryService) Cran(params CranRemoteRepositoryParams) error {
	return rrs.performRequest(params, params.Key)
}

func (rrs *RemoteRepositoryService) Debian(params DebianRemoteRepositoryParams) error {
	return rrs.performRequest(params, params.Key)
}

func (rrs *RemoteRepositoryService) Docker(params DockerRemoteRepositoryParams) error {
	return rrs.performRequest(params, params.Key)
}

func (rrs *RemoteRepositoryService) Gems(params GemsRemoteRepositoryParams) error {
	return rrs.performRequest(params, params.Key)
}

func (rrs *RemoteRepositoryService) Generic(params GenericRemoteRepositoryParams) error {
	return rrs.performRequest(params, params.Key)
}

func (rrs *RemoteRepositoryService) Gitlfs(params GitlfsRemoteRepositoryParams) error {
	return rrs.performRequest(params, params.Key)
}

func (rrs *RemoteRepositoryService) Go(params GoRemoteRepositoryParams) error {
	return rrs.performRequest(params, params.Key)
}

func (rrs *RemoteRepositoryService) Gradle(params GradleRemoteRepositoryParams) error {
	return rrs.performRequest(params, params.Key)
}

func (rrs *RemoteRepositoryService) Helm(params HelmRemoteRepositoryParams) error {
	return rrs.performRequest(params, params.Key)
}

func (rrs *RemoteRepositoryService) Ivy(params IvyRemoteRepositoryParams) error {
	return rrs.performRequest(params, params.Key)
}

func (rrs *RemoteRepositoryService) Maven(params MavenRemoteRepositoryParams) error {
	return rrs.performRequest(params, params.Key)
}

func (rrs *RemoteRepositoryService) Npm(params NpmRemoteRepositoryParams) error {
	return rrs.performRequest(params, params.Key)
}

func (rrs *RemoteRepositoryService) Nuget(params NugetRemoteRepositoryParams) error {
	return rrs.performRequest(params, params.Key)
}

func (rrs *RemoteRepositoryService) Opkg(params OpkgRemoteRepositoryParams) error {
	return rrs.performRequest(params, params.Key)
}

func (rrs *RemoteRepositoryService) P2(params P2RemoteRepositoryParams) error {
	return rrs.performRequest(params, params.Key)
}

func (rrs *RemoteRepositoryService) Puppet(params PuppetRemoteRepositoryParams) error {
	return rrs.performRequest(params, params.Key)
}

func (rrs *RemoteRepositoryService) Pypi(params PypiRemoteRepositoryParams) error {
	return rrs.performRequest(params, params.Key)
}

func (rrs *RemoteRepositoryService) Rpm(params RpmRemoteRepositoryParams) error {
	return rrs.performRequest(params, params.Key)
}

func (rrs *RemoteRepositoryService) Sbt(params SbtRemoteRepositoryParams) error {
	return rrs.performRequest(params, params.Key)
}

func (rrs *RemoteRepositoryService) Swift(params SwiftRemoteRepositoryParams) error {
	return rrs.performRequest(params, params.Key)
}

func (rrs *RemoteRepositoryService) Vcs(params VcsRemoteRepositoryParams) error {
	return rrs.performRequest(params, params.Key)
}

func (rrs *RemoteRepositoryService) Yum(params YumRemoteRepositoryParams) error {
	return rrs.performRequest(params, params.Key)
}

type ContentSynchronisationStatistics struct {
	Enabled *bool `json:"enabled,omitempty"`
}

type ContentSynchronisationProperties struct {
	Enabled *bool `json:"enabled,omitempty"`
}

type ContentSynchronisationSource struct {
	OriginAbsenceDetection *bool `json:"originAbsenceDetection,omitempty"`
}

type ContentSynchronisation struct {
	Enabled    *bool                             `json:"enabled,omitempty"`
	Statistics *ContentSynchronisationStatistics `json:"statistics,omitempty"`
	Properties *ContentSynchronisationProperties `json:"properties,omitempty"`
	Source     *ContentSynchronisationSource     `json:"source,omitempty"`
}

type RemoteRepositoryBaseParams struct {
	RepositoryBaseParams
	AdditionalRepositoryBaseParams
	Url                               string                  `json:"url"`
	Username                          string                  `json:"username,omitempty"`
	Password                          string                  `json:"password,omitempty"`
	Proxy                             string                  `json:"proxy,omitempty"`
	HardFail                          *bool                   `json:"hardFail,omitempty"`
	Offline                           *bool                   `json:"offline,omitempty"`
	StoreArtifactsLocally             *bool                   `json:"storeArtifactsLocally,omitempty"`
	SocketTimeoutMillis               int                     `json:"socketTimeoutMillis,omitempty"`
	LocalAddress                      string                  `json:"localAddress,omitempty"`
	RetrievalCachePeriodSecs          int                     `json:"retrievalCachePeriodSecs,omitempty"`
	MetadataRetrievalTimeoutSecs      int                     `json:"metadataRetrievalTimeoutSecs,omitempty"`
	MissedRetrievalCachePeriodSecs    int                     `json:"missedRetrievalCachePeriodSecs,omitempty"`
	UnusedArtifactsCleanupPeriodHours int                     `json:"unusedArtifactsCleanupPeriodHours,omitempty"`
	AssumedOfflinePeriodSecs          int                     `json:"assumedOfflinePeriodSecs,omitempty"`
	ShareConfiguration                *bool                   `json:"shareConfiguration,omitempty"`
	SynchronizeProperties             *bool                   `json:"synchronizeProperties,omitempty"`
	BlockMismatchingMimeTypes         *bool                   `json:"blockMismatchingMimeTypes,omitempty"`
	MismatchingMimeTypesOverrideList  string                  `json:"mismatchingMimeTypesOverrideList,omitempty"`
	AllowAnyHostAuth                  *bool                   `json:"allowAnyHostAuth,omitempty"`
	EnableCookieManagement            *bool                   `json:"enableCookieManagement,omitempty"`
	BypassHeadRequests                *bool                   `json:"bypassHeadRequests,omitempty"`
	ClientTlsCertificate              string                  `json:"clientTlsCertificate,omitempty"`
	ContentSynchronisation            *ContentSynchronisation `json:"contentSynchronisation,omitempty"`
}

func NewRemoteRepositoryBaseParams() RemoteRepositoryBaseParams {
	return RemoteRepositoryBaseParams{RepositoryBaseParams: RepositoryBaseParams{Rclass: RemoteRepositoryRepoType}}
}

func NewRemoteRepositoryPackageParams(packageType string) RemoteRepositoryBaseParams {
	return RemoteRepositoryBaseParams{RepositoryBaseParams: RepositoryBaseParams{Rclass: RemoteRepositoryRepoType, PackageType: packageType}}
}

type JavaPackageManagersRemoteRepositoryParams struct {
	RemoteRepoChecksumPolicyType string `json:"remoteRepoChecksumPolicyType,omitempty"`
	MaxUniqueSnapshots           int    `json:"maxUniqueSnapshots,omitempty"`
	FetchJarsEagerly             *bool  `json:"fetchJarsEagerly,omitempty"`
	SuppressPomConsistencyChecks *bool  `json:"suppressPomConsistencyChecks,omitempty"`
	FetchSourcesEagerly          *bool  `json:"fetchSourcesEagerly,omitempty"`
	HandleReleases               *bool  `json:"handleReleases,omitempty"`
	HandleSnapshots              *bool  `json:"handleSnapshots,omitempty"`
	RejectInvalidJars            *bool  `json:"rejectInvalidJars,omitempty"`
}

type VcsGitRemoteRepositoryParams struct {
	VcsType           string `json:"vcsType,omitempty"`
	VcsGitProvider    string `json:"vcsGitProvider,omitempty"`
	VcsGitDownloadUrl string `json:"vcsGitDownloadUrl,omitempty"`
}

type AlpineRemoteRepositoryParams struct {
	RemoteRepositoryBaseParams
}

func NewAlpineRemoteRepositoryParams() AlpineRemoteRepositoryParams {
	return AlpineRemoteRepositoryParams{RemoteRepositoryBaseParams: NewRemoteRepositoryPackageParams("alpine")}
}

type BowerRemoteRepositoryParams struct {
	RemoteRepositoryBaseParams
	VcsGitRemoteRepositoryParams
	BowerRegistryUrl string `json:"bowerRegistryUrl,omitempty"`
}

func NewBowerRemoteRepositoryParams() BowerRemoteRepositoryParams {
	return BowerRemoteRepositoryParams{RemoteRepositoryBaseParams: NewRemoteRepositoryPackageParams("bower")}
}

type CargoRemoteRepositoryParams struct {
	RemoteRepositoryBaseParams
	CargoRepositoryParams
	GitRegistryUrl string `json:"gitRegistryUrl,omitempty"`
}

func NewCargoRemoteRepositoryParams() CargoRemoteRepositoryParams {
	return CargoRemoteRepositoryParams{RemoteRepositoryBaseParams: NewRemoteRepositoryPackageParams("cargo")}
}

type ChefRemoteRepositoryParams struct {
	RemoteRepositoryBaseParams
}

func NewChefRemoteRepositoryParams() ChefRemoteRepositoryParams {
	return ChefRemoteRepositoryParams{RemoteRepositoryBaseParams: NewRemoteRepositoryPackageParams("chef")}
}

type CocoapodsRemoteRepositoryParams struct {
	RemoteRepositoryBaseParams
	VcsGitRemoteRepositoryParams
	PodsSpecsRepoUrl string `json:"podsSpecsRepoUrl,omitempty"`
}

func NewCocoapodsRemoteRepositoryParams() CocoapodsRemoteRepositoryParams {
	return CocoapodsRemoteRepositoryParams{RemoteRepositoryBaseParams: NewRemoteRepositoryPackageParams("cocoapods")}
}

type ComposerRemoteRepositoryParams struct {
	RemoteRepositoryBaseParams
	VcsGitRemoteRepositoryParams
	ComposerRegistryUrl string `json:"composerRegistryUrl,omitempty"`
}

func NewComposerRemoteRepositoryParams() ComposerRemoteRepositoryParams {
	return ComposerRemoteRepositoryParams{RemoteRepositoryBaseParams: NewRemoteRepositoryPackageParams("composer")}
}

type ConanRemoteRepositoryParams struct {
	RemoteRepositoryBaseParams
}

func NewConanRemoteRepositoryParams() ConanRemoteRepositoryParams {
	return ConanRemoteRepositoryParams{RemoteRepositoryBaseParams: NewRemoteRepositoryPackageParams("conan")}
}

type CondaRemoteRepositoryParams struct {
	RemoteRepositoryBaseParams
}

func NewCondaRemoteRepositoryParams() CondaRemoteRepositoryParams {
	return CondaRemoteRepositoryParams{RemoteRepositoryBaseParams: NewRemoteRepositoryPackageParams("conda")}
}

type CranRemoteRepositoryParams struct {
	RemoteRepositoryBaseParams
}

func NewCranRemoteRepositoryParams() CranRemoteRepositoryParams {
	return CranRemoteRepositoryParams{RemoteRepositoryBaseParams: NewRemoteRepositoryPackageParams("cran")}
}

type DebianRemoteRepositoryParams struct {
	RemoteRepositoryBaseParams
	ListRemoteFolderItems *bool `json:"listRemoteFolderItems,omitempty"`
}

func NewDebianRemoteRepositoryParams() DebianRemoteRepositoryParams {
	return DebianRemoteRepositoryParams{RemoteRepositoryBaseParams: NewRemoteRepositoryPackageParams("debian")}
}

type DockerRemoteRepositoryParams struct {
	RemoteRepositoryBaseParams
	ExternalDependenciesEnabled  *bool    `json:"externalDependenciesEnabled,omitempty"`
	ExternalDependenciesPatterns []string `json:"externalDependenciesPatterns,omitempty"`
	EnableTokenAuthentication    *bool    `json:"enableTokenAuthentication,omitempty"`
	BlockPullingSchema1          *bool    `json:"blockPushingSchema1,omitempty"`
}

func NewDockerRemoteRepositoryParams() DockerRemoteRepositoryParams {
	return DockerRemoteRepositoryParams{RemoteRepositoryBaseParams: NewRemoteRepositoryPackageParams("docker")}
}

type GemsRemoteRepositoryParams struct {
	RemoteRepositoryBaseParams
	ListRemoteFolderItems *bool `json:"listRemoteFolderItems,omitempty"`
}

func NewGemsRemoteRepositoryParams() GemsRemoteRepositoryParams {
	return GemsRemoteRepositoryParams{RemoteRepositoryBaseParams: NewRemoteRepositoryPackageParams("gems")}
}

type GenericRemoteRepositoryParams struct {
	RemoteRepositoryBaseParams
	ListRemoteFolderItems *bool `json:"listRemoteFolderItems,omitempty"`
}

func NewGenericRemoteRepositoryParams() GenericRemoteRepositoryParams {
	return GenericRemoteRepositoryParams{RemoteRepositoryBaseParams: NewRemoteRepositoryPackageParams("generic")}
}

type GitlfsRemoteRepositoryParams struct {
	RemoteRepositoryBaseParams
}

func NewGitlfsRemoteRepositoryParams() GitlfsRemoteRepositoryParams {
	return GitlfsRemoteRepositoryParams{RemoteRepositoryBaseParams: NewRemoteRepositoryPackageParams("gitlfs")}
}

type GoRemoteRepositoryParams struct {
	RemoteRepositoryBaseParams
	VcsGitProvider string `json:"vcsGitProvider,omitempty"`
}

func NewGoRemoteRepositoryParams() GoRemoteRepositoryParams {
	return GoRemoteRepositoryParams{RemoteRepositoryBaseParams: NewRemoteRepositoryPackageParams("go")}
}

type GradleRemoteRepositoryParams struct {
	RemoteRepositoryBaseParams
	JavaPackageManagersRemoteRepositoryParams
}

func NewGradleRemoteRepositoryParams() GradleRemoteRepositoryParams {
	return GradleRemoteRepositoryParams{RemoteRepositoryBaseParams: NewRemoteRepositoryPackageParams("gradle")}
}

type HelmRemoteRepositoryParams struct {
	RemoteRepositoryBaseParams
	ChartsBaseUrl string `json:"chartsBaseUrl,omitempty"`
}

func NewHelmRemoteRepositoryParams() HelmRemoteRepositoryParams {
	return HelmRemoteRepositoryParams{RemoteRepositoryBaseParams: NewRemoteRepositoryPackageParams("helm")}
}

type IvyRemoteRepositoryParams struct {
	RemoteRepositoryBaseParams
	JavaPackageManagersRemoteRepositoryParams
}

func NewIvyRemoteRepositoryParams() IvyRemoteRepositoryParams {
	return IvyRemoteRepositoryParams{RemoteRepositoryBaseParams: NewRemoteRepositoryPackageParams("ivy")}
}

type MavenRemoteRepositoryParams struct {
	RemoteRepositoryBaseParams
	JavaPackageManagersRemoteRepositoryParams
}

func NewMavenRemoteRepositoryParams() MavenRemoteRepositoryParams {
	return MavenRemoteRepositoryParams{RemoteRepositoryBaseParams: NewRemoteRepositoryPackageParams("maven")}
}

type NpmRemoteRepositoryParams struct {
	RemoteRepositoryBaseParams
}

func NewNpmRemoteRepositoryParams() NpmRemoteRepositoryParams {
	return NpmRemoteRepositoryParams{RemoteRepositoryBaseParams: NewRemoteRepositoryPackageParams("npm")}
}

type NugetRemoteRepositoryParams struct {
	RemoteRepositoryBaseParams
	FeedContextPath          string `json:"feedContextPath,omitempty"`
	DownloadContextPath      string `json:"downloadContextPath,omitempty"`
	V3FeedUrl                string `json:"v3FeedUrl,omitempty"`
	ForceNugetAuthentication *bool  `json:"forceNugetAuthentication,omitempty"`
}

func NewNugetRemoteRepositoryParams() NugetRemoteRepositoryParams {
	return NugetRemoteRepositoryParams{RemoteRepositoryBaseParams: NewRemoteRepositoryPackageParams("nuget")}
}

type OpkgRemoteRepositoryParams struct {
	RemoteRepositoryBaseParams
}

func NewOpkgRemoteRepositoryParams() OpkgRemoteRepositoryParams {
	return OpkgRemoteRepositoryParams{RemoteRepositoryBaseParams: NewRemoteRepositoryPackageParams("opkg")}
}

type P2RemoteRepositoryParams struct {
	RemoteRepositoryBaseParams
	ListRemoteFolderItems *bool `json:"listRemoteFolderItems,omitempty"`
}

func NewP2RemoteRepositoryParams() P2RemoteRepositoryParams {
	return P2RemoteRepositoryParams{RemoteRepositoryBaseParams: NewRemoteRepositoryPackageParams("p2")}
}

type PuppetRemoteRepositoryParams struct {
	RemoteRepositoryBaseParams
}

func NewPuppetRemoteRepositoryParams() PuppetRemoteRepositoryParams {
	return PuppetRemoteRepositoryParams{RemoteRepositoryBaseParams: NewRemoteRepositoryPackageParams("puppet")}
}

type PypiRemoteRepositoryParams struct {
	RemoteRepositoryBaseParams
	PypiRegistryUrl      string `json:"pyPIRegistryUrl,omitempty"`
	PypiRepositorySuffix string `json:"pyPIRepositorySuffix,omitempty"`
}

func NewPypiRemoteRepositoryParams() PypiRemoteRepositoryParams {
	return PypiRemoteRepositoryParams{RemoteRepositoryBaseParams: NewRemoteRepositoryPackageParams("pypi")}
}

type RpmRemoteRepositoryParams struct {
	RemoteRepositoryBaseParams
	ListRemoteFolderItems *bool `json:"listRemoteFolderItems,omitempty"`
}

func NewRpmRemoteRepositoryParams() RpmRemoteRepositoryParams {
	return RpmRemoteRepositoryParams{RemoteRepositoryBaseParams: NewRemoteRepositoryPackageParams("rpm")}
}

type SbtRemoteRepositoryParams struct {
	RemoteRepositoryBaseParams
	JavaPackageManagersRemoteRepositoryParams
}

func NewSbtRemoteRepositoryParams() SbtRemoteRepositoryParams {
	return SbtRemoteRepositoryParams{RemoteRepositoryBaseParams: NewRemoteRepositoryPackageParams("sbt")}
}

type SwiftRemoteRepositoryParams struct {
	RemoteRepositoryBaseParams
}

func NewSwiftRemoteRepositoryParams() SwiftRemoteRepositoryParams {
	return SwiftRemoteRepositoryParams{RemoteRepositoryBaseParams: NewRemoteRepositoryPackageParams("swift")}
}

type VcsRemoteRepositoryParams struct {
	RemoteRepositoryBaseParams
	VcsGitRemoteRepositoryParams
	MaxUniqueSnapshots int `json:"maxUniqueSnapshots,omitempty"`
}

func NewVcsRemoteRepositoryParams() VcsRemoteRepositoryParams {
	return VcsRemoteRepositoryParams{RemoteRepositoryBaseParams: NewRemoteRepositoryPackageParams("vcs")}
}

type YumRemoteRepositoryParams struct {
	RemoteRepositoryBaseParams
	ListRemoteFolderItems *bool `json:"listRemoteFolderItems,omitempty"`
}

func NewYumRemoteRepositoryParams() YumRemoteRepositoryParams {
	return YumRemoteRepositoryParams{RemoteRepositoryBaseParams: NewRemoteRepositoryPackageParams("yum")}
}
