package services

import (
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
)

const VirtualRepositoryRepoType = "virtual"

type VirtualRepositoryService struct {
	RepositoryService
}

func NewVirtualRepositoryService(client *jfroghttpclient.JfrogHttpClient, isUpdate bool) *VirtualRepositoryService {
	return &VirtualRepositoryService{
		RepositoryService: RepositoryService{
			repoType: VirtualRepositoryRepoType,
			client:   client,
			isUpdate: isUpdate,
		},
	}
}

func (vrs *VirtualRepositoryService) GetJfrogHttpClient() *jfroghttpclient.JfrogHttpClient {
	return vrs.client
}

func (vrs *VirtualRepositoryService) performRequest(params interface{}, repoKey string) error {
	return vrs.RepositoryService.performRequest(params, repoKey)
}

func (vrs *VirtualRepositoryService) Alpine(params AlpineVirtualRepositoryParams) error {
	return vrs.performRequest(params, params.Key)
}

func (vrs *VirtualRepositoryService) Bower(params BowerVirtualRepositoryParams) error {
	return vrs.performRequest(params, params.Key)
}

func (vrs *VirtualRepositoryService) Chef(params ChefVirtualRepositoryParams) error {
	return vrs.performRequest(params, params.Key)
}

func (vrs *VirtualRepositoryService) Conan(params ConanVirtualRepositoryParams) error {
	return vrs.performRequest(params, params.Key)
}

func (vrs *VirtualRepositoryService) Conda(params CondaVirtualRepositoryParams) error {
	return vrs.performRequest(params, params.Key)
}

func (vrs *VirtualRepositoryService) Cran(params CranVirtualRepositoryParams) error {
	return vrs.performRequest(params, params.Key)
}

func (vrs *VirtualRepositoryService) Debian(params DebianVirtualRepositoryParams) error {
	return vrs.performRequest(params, params.Key)
}

func (vrs *VirtualRepositoryService) Docker(params DockerVirtualRepositoryParams) error {
	return vrs.performRequest(params, params.Key)
}

func (vrs *VirtualRepositoryService) Gems(params GemsVirtualRepositoryParams) error {
	return vrs.performRequest(params, params.Key)
}

func (vrs *VirtualRepositoryService) Generic(params GenericVirtualRepositoryParams) error {
	return vrs.performRequest(params, params.Key)
}

func (vrs *VirtualRepositoryService) Gitlfs(params GitlfsVirtualRepositoryParams) error {
	return vrs.performRequest(params, params.Key)
}

func (vrs *VirtualRepositoryService) Go(params GoVirtualRepositoryParams) error {
	return vrs.performRequest(params, params.Key)
}

func (vrs *VirtualRepositoryService) Gradle(params GradleVirtualRepositoryParams) error {
	return vrs.performRequest(params, params.Key)
}

func (vrs *VirtualRepositoryService) Helm(params HelmVirtualRepositoryParams) error {
	return vrs.performRequest(params, params.Key)
}

func (vrs *VirtualRepositoryService) Ivy(params IvyVirtualRepositoryParams) error {
	return vrs.performRequest(params, params.Key)
}

func (vrs *VirtualRepositoryService) Maven(params MavenVirtualRepositoryParams) error {
	return vrs.performRequest(params, params.Key)
}

func (vrs *VirtualRepositoryService) Npm(params NpmVirtualRepositoryParams) error {
	return vrs.performRequest(params, params.Key)
}

func (vrs *VirtualRepositoryService) Nuget(params NugetVirtualRepositoryParams) error {
	return vrs.performRequest(params, params.Key)
}

func (vrs *VirtualRepositoryService) P2(params P2VirtualRepositoryParams) error {
	return vrs.performRequest(params, params.Key)
}

func (vrs *VirtualRepositoryService) Puppet(params PuppetVirtualRepositoryParams) error {
	return vrs.performRequest(params, params.Key)
}

func (vrs *VirtualRepositoryService) Pypi(params PypiVirtualRepositoryParams) error {
	return vrs.performRequest(params, params.Key)
}

func (vrs *VirtualRepositoryService) Rpm(params RpmVirtualRepositoryParams) error {
	return vrs.performRequest(params, params.Key)
}

func (vrs *VirtualRepositoryService) Sbt(params SbtVirtualRepositoryParams) error {
	return vrs.performRequest(params, params.Key)
}

func (vrs *VirtualRepositoryService) Swift(params SwiftVirtualRepositoryParams) error {
	return vrs.performRequest(params, params.Key)
}

func (vrs *VirtualRepositoryService) Yum(params YumVirtualRepositoryParams) error {
	return vrs.performRequest(params, params.Key)
}

type VirtualRepositoryBaseParams struct {
	RepositoryBaseParams
	KeyPairRefsRepositoryParams
	Repositories                                  []string `json:"repositories,omitempty"`
	ArtifactoryRequestsCanRetrieveRemoteArtifacts *bool    `json:"artifactoryRequestsCanRetrieveRemoteArtifacts,omitempty"`
	DefaultDeploymentRepo                         string   `json:"defaultDeploymentRepo,omitempty"`
}

type CommonJavaVirtualRepositoryParams struct {
	PomRepositoryReferencesCleanupPolicy string `json:"pomRepositoryReferencesCleanupPolicy,omitempty"`
	KeyPair                              string `json:"keyPair,omitempty"`
}

type CommonCacheVirtualRepositoryParams struct {
	VirtualRetrievalCachePeriodSecs int `json:"virtualRetrievalCachePeriodSecs,omitempty"`
}

func NewVirtualRepositoryBaseParams() VirtualRepositoryBaseParams {
	return VirtualRepositoryBaseParams{RepositoryBaseParams: RepositoryBaseParams{Rclass: VirtualRepositoryRepoType}}
}

func NewVirtualRepositoryPackageParams(packageType string) VirtualRepositoryBaseParams {
	return VirtualRepositoryBaseParams{RepositoryBaseParams: RepositoryBaseParams{Rclass: VirtualRepositoryRepoType, PackageType: packageType}}
}

type AlpineVirtualRepositoryParams struct {
	VirtualRepositoryBaseParams
	CommonCacheVirtualRepositoryParams
}

func NewAlpineVirtualRepositoryParams() AlpineVirtualRepositoryParams {
	return AlpineVirtualRepositoryParams{VirtualRepositoryBaseParams: NewVirtualRepositoryPackageParams("alpine")}
}

type BowerVirtualRepositoryParams struct {
	VirtualRepositoryBaseParams
	ExternalDependenciesEnabled    *bool    `json:"externalDependenciesEnabled,omitempty"`
	ExternalDependenciesPatterns   []string `json:"externalDependenciesPatterns,omitempty"`
	ExternalDependenciesRemoteRepo string   `json:"externalDependenciesRemoteRepo,omitempty"`
}

func NewBowerVirtualRepositoryParams() BowerVirtualRepositoryParams {
	return BowerVirtualRepositoryParams{VirtualRepositoryBaseParams: NewVirtualRepositoryPackageParams("bower")}
}

type ChefVirtualRepositoryParams struct {
	VirtualRepositoryBaseParams
	CommonCacheVirtualRepositoryParams
}

func NewChefVirtualRepositoryParams() ChefVirtualRepositoryParams {
	return ChefVirtualRepositoryParams{VirtualRepositoryBaseParams: NewVirtualRepositoryPackageParams("chef")}
}

type ConanVirtualRepositoryParams struct {
	VirtualRepositoryBaseParams
	CommonCacheVirtualRepositoryParams
}

func NewConanVirtualRepositoryParams() ConanVirtualRepositoryParams {
	return ConanVirtualRepositoryParams{VirtualRepositoryBaseParams: NewVirtualRepositoryPackageParams("conan")}
}

type CondaVirtualRepositoryParams struct {
	VirtualRepositoryBaseParams
	CommonCacheVirtualRepositoryParams
}

func NewCondaVirtualRepositoryParams() CondaVirtualRepositoryParams {
	return CondaVirtualRepositoryParams{VirtualRepositoryBaseParams: NewVirtualRepositoryPackageParams("conda")}
}

type CranVirtualRepositoryParams struct {
	VirtualRepositoryBaseParams
	CommonCacheVirtualRepositoryParams
}

func NewCranVirtualRepositoryParams() CranVirtualRepositoryParams {
	return CranVirtualRepositoryParams{VirtualRepositoryBaseParams: NewVirtualRepositoryPackageParams("cran")}
}

type DebianVirtualRepositoryParams struct {
	VirtualRepositoryBaseParams
	CommonCacheVirtualRepositoryParams
	DebianDefaultArchitectures      string   `json:"debianDefaultArchitectures,omitempty"`
	OptionalIndexCompressionFormats []string `json:"optionalIndexCompressionFormats,omitempty"`
}

func NewDebianVirtualRepositoryParams() DebianVirtualRepositoryParams {
	return DebianVirtualRepositoryParams{VirtualRepositoryBaseParams: NewVirtualRepositoryPackageParams("debian")}
}

type DockerVirtualRepositoryParams struct {
	VirtualRepositoryBaseParams
	ResolveDockerTagsByTimestamp *bool `json:"resolveDockerTagsByTimestamp,omitempty"`
}

func NewDockerVirtualRepositoryParams() DockerVirtualRepositoryParams {
	return DockerVirtualRepositoryParams{VirtualRepositoryBaseParams: NewVirtualRepositoryPackageParams("docker")}
}

type GemsVirtualRepositoryParams struct {
	VirtualRepositoryBaseParams
}

func NewGemsVirtualRepositoryParams() GemsVirtualRepositoryParams {
	return GemsVirtualRepositoryParams{VirtualRepositoryBaseParams: NewVirtualRepositoryPackageParams("gems")}
}

type GenericVirtualRepositoryParams struct {
	VirtualRepositoryBaseParams
}

func NewGenericVirtualRepositoryParams() GenericVirtualRepositoryParams {
	return GenericVirtualRepositoryParams{VirtualRepositoryBaseParams: NewVirtualRepositoryPackageParams("generic")}
}

type GitlfsVirtualRepositoryParams struct {
	VirtualRepositoryBaseParams
}

func NewGitlfsVirtualRepositoryParams() GitlfsVirtualRepositoryParams {
	return GitlfsVirtualRepositoryParams{VirtualRepositoryBaseParams: NewVirtualRepositoryPackageParams("gitlfs")}
}

type GoVirtualRepositoryParams struct {
	VirtualRepositoryBaseParams
	ExternalDependenciesEnabled  *bool    `json:"externalDependenciesEnabled,omitempty"`
	ExternalDependenciesPatterns []string `json:"externalDependenciesPatterns,omitempty"`
}

func NewGoVirtualRepositoryParams() GoVirtualRepositoryParams {
	return GoVirtualRepositoryParams{VirtualRepositoryBaseParams: NewVirtualRepositoryPackageParams("go")}
}

type GradleVirtualRepositoryParams struct {
	VirtualRepositoryBaseParams
	CommonJavaVirtualRepositoryParams
}

func NewGradleVirtualRepositoryParams() GradleVirtualRepositoryParams {
	return GradleVirtualRepositoryParams{VirtualRepositoryBaseParams: NewVirtualRepositoryPackageParams("gradle")}
}

type HelmVirtualRepositoryParams struct {
	VirtualRepositoryBaseParams
	CommonCacheVirtualRepositoryParams
}

func NewHelmVirtualRepositoryParams() HelmVirtualRepositoryParams {
	return HelmVirtualRepositoryParams{VirtualRepositoryBaseParams: NewVirtualRepositoryPackageParams("helm")}
}

type IvyVirtualRepositoryParams struct {
	VirtualRepositoryBaseParams
	CommonJavaVirtualRepositoryParams
}

func NewIvyVirtualRepositoryParams() IvyVirtualRepositoryParams {
	return IvyVirtualRepositoryParams{VirtualRepositoryBaseParams: NewVirtualRepositoryPackageParams("ivy")}
}

type MavenVirtualRepositoryParams struct {
	VirtualRepositoryBaseParams
	CommonJavaVirtualRepositoryParams
	ForceMavenAuthentication *bool `json:"forceMavenAuthentication,omitempty"`
}

func NewMavenVirtualRepositoryParams() MavenVirtualRepositoryParams {
	return MavenVirtualRepositoryParams{VirtualRepositoryBaseParams: NewVirtualRepositoryPackageParams("maven")}
}

type NpmVirtualRepositoryParams struct {
	VirtualRepositoryBaseParams
	CommonCacheVirtualRepositoryParams
	ExternalDependenciesEnabled    *bool    `json:"externalDependenciesEnabled,omitempty"`
	ExternalDependenciesPatterns   []string `json:"externalDependenciesPatterns,omitempty"`
	ExternalDependenciesRemoteRepo string   `json:"externalDependenciesRemoteRepo,omitempty"`
}

func NewNpmVirtualRepositoryParams() NpmVirtualRepositoryParams {
	return NpmVirtualRepositoryParams{VirtualRepositoryBaseParams: NewVirtualRepositoryPackageParams("npm")}
}

type NugetVirtualRepositoryParams struct {
	VirtualRepositoryBaseParams
	ForceNugetAuthentication *bool `json:"forceNugetAuthentication,omitempty"`
}

func NewNugetVirtualRepositoryParams() NugetVirtualRepositoryParams {
	return NugetVirtualRepositoryParams{VirtualRepositoryBaseParams: NewVirtualRepositoryPackageParams("nuget")}
}

type P2VirtualRepositoryParams struct {
	VirtualRepositoryBaseParams
	P2Urls []string `json:"p2Urls,omitempty"`
}

func NewP2VirtualRepositoryParams() P2VirtualRepositoryParams {
	return P2VirtualRepositoryParams{VirtualRepositoryBaseParams: NewVirtualRepositoryPackageParams("p2")}
}

type PuppetVirtualRepositoryParams struct {
	VirtualRepositoryBaseParams
}

func NewPuppetVirtualRepositoryParams() PuppetVirtualRepositoryParams {
	return PuppetVirtualRepositoryParams{VirtualRepositoryBaseParams: NewVirtualRepositoryPackageParams("puppet")}
}

type PypiVirtualRepositoryParams struct {
	VirtualRepositoryBaseParams
}

func NewPypiVirtualRepositoryParams() PypiVirtualRepositoryParams {
	return PypiVirtualRepositoryParams{VirtualRepositoryBaseParams: NewVirtualRepositoryPackageParams("pypi")}
}

type RpmVirtualRepositoryParams struct {
	VirtualRepositoryBaseParams
	CommonCacheVirtualRepositoryParams
}

func NewRpmVirtualRepositoryParams() RpmVirtualRepositoryParams {
	return RpmVirtualRepositoryParams{VirtualRepositoryBaseParams: NewVirtualRepositoryPackageParams("rpm")}
}

type SbtVirtualRepositoryParams struct {
	VirtualRepositoryBaseParams
	CommonJavaVirtualRepositoryParams
}

func NewSbtVirtualRepositoryParams() SbtVirtualRepositoryParams {
	return SbtVirtualRepositoryParams{VirtualRepositoryBaseParams: NewVirtualRepositoryPackageParams("sbt")}
}

type SwiftVirtualRepositoryParams struct {
	VirtualRepositoryBaseParams
}

func NewSwiftVirtualRepositoryParams() SwiftVirtualRepositoryParams {
	return SwiftVirtualRepositoryParams{VirtualRepositoryBaseParams: NewVirtualRepositoryPackageParams("swift")}
}

type YumVirtualRepositoryParams struct {
	VirtualRepositoryBaseParams
	CommonCacheVirtualRepositoryParams
}

func NewYumVirtualRepositoryParams() YumVirtualRepositoryParams {
	return YumVirtualRepositoryParams{VirtualRepositoryBaseParams: NewVirtualRepositoryPackageParams("yum")}
}
