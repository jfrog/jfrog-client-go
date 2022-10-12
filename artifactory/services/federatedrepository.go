package services

import (
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
)

const FederatedRepositoryRepoType = "federated"

type FederatedRepositoryService struct {
	RepositoryService
}

func NewFederatedRepositoryService(client *jfroghttpclient.JfrogHttpClient, isUpdate bool) *FederatedRepositoryService {
	return &FederatedRepositoryService{
		RepositoryService: RepositoryService{
			repoType: FederatedRepositoryRepoType,
			client:   client,
			isUpdate: isUpdate,
		},
	}
}

func (frs *FederatedRepositoryService) GetJfrogHttpClient() *jfroghttpclient.JfrogHttpClient {
	return frs.client
}

func (frs *FederatedRepositoryService) performRequest(params interface{}, repoKey string) error {
	return frs.RepositoryService.performRequest(params, repoKey)
}

func (frs *FederatedRepositoryService) Alpine(params AlpineFederatedRepositoryParams) error {
	return frs.performRequest(params, params.Key)
}

func (frs *FederatedRepositoryService) Bower(params BowerFederatedRepositoryParams) error {
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

func (frs *FederatedRepositoryService) Cran(params CranFederatedRepositoryParams) error {
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

func (frs *FederatedRepositoryService) Swift(params SwiftFederatedRepositoryParams) error {
	return frs.performRequest(params, params.Key)
}

func (frs *FederatedRepositoryService) Vagrant(params VagrantFederatedRepositoryParams) error {
	return frs.performRequest(params, params.Key)
}

func (frs *FederatedRepositoryService) Yum(params YumFederatedRepositoryParams) error {
	return frs.performRequest(params, params.Key)
}

type FederatedRepositoryMember struct {
	Url     string `json:"url"`
	Enabled *bool  `json:"enabled,omitempty"`
}

type FederatedRepositoryBaseParams struct {
	RepositoryBaseParams
	AdditionalRepositoryBaseParams
	KeyPairRefsRepositoryParams
	ArchiveBrowsingEnabled *bool                       `json:"archiveBrowsingEnabled,omitempty"`
	Members                []FederatedRepositoryMember `json:"members,omitempty"`
}

func NewFederatedRepositoryBaseParams() FederatedRepositoryBaseParams {
	return FederatedRepositoryBaseParams{RepositoryBaseParams: RepositoryBaseParams{Rclass: FederatedRepositoryRepoType}}
}

func NewFederatedRepositoryPackageParams(packageType string) FederatedRepositoryBaseParams {
	return FederatedRepositoryBaseParams{RepositoryBaseParams: RepositoryBaseParams{Rclass: FederatedRepositoryRepoType, PackageType: packageType}}
}

type AlpineFederatedRepositoryParams struct {
	FederatedRepositoryBaseParams
}

func NewAlpineFederatedRepositoryParams() AlpineFederatedRepositoryParams {
	return AlpineFederatedRepositoryParams{FederatedRepositoryBaseParams: NewFederatedRepositoryPackageParams("alpine")}
}

type BowerFederatedRepositoryParams struct {
	FederatedRepositoryBaseParams
}

func NewBowerFederatedRepositoryParams() BowerFederatedRepositoryParams {
	return BowerFederatedRepositoryParams{FederatedRepositoryBaseParams: NewFederatedRepositoryPackageParams("bower")}
}

type CargoFederatedRepositoryParams struct {
	FederatedRepositoryBaseParams
	CargoRepositoryParams
}

func NewCargoFederatedRepositoryParams() CargoFederatedRepositoryParams {
	return CargoFederatedRepositoryParams{FederatedRepositoryBaseParams: NewFederatedRepositoryPackageParams("cargo")}
}

type ChefFederatedRepositoryParams struct {
	FederatedRepositoryBaseParams
}

func NewChefFederatedRepositoryParams() ChefFederatedRepositoryParams {
	return ChefFederatedRepositoryParams{FederatedRepositoryBaseParams: NewFederatedRepositoryPackageParams("chef")}
}

type CocoapodsFederatedRepositoryParams struct {
	FederatedRepositoryBaseParams
}

func NewCocoapodsFederatedRepositoryParams() CocoapodsFederatedRepositoryParams {
	return CocoapodsFederatedRepositoryParams{FederatedRepositoryBaseParams: NewFederatedRepositoryPackageParams("cocoapods")}
}

type ComposerFederatedRepositoryParams struct {
	FederatedRepositoryBaseParams
}

func NewComposerFederatedRepositoryParams() ComposerFederatedRepositoryParams {
	return ComposerFederatedRepositoryParams{FederatedRepositoryBaseParams: NewFederatedRepositoryPackageParams("composer")}
}

type ConanFederatedRepositoryParams struct {
	FederatedRepositoryBaseParams
}

func NewConanFederatedRepositoryParams() ConanFederatedRepositoryParams {
	return ConanFederatedRepositoryParams{FederatedRepositoryBaseParams: NewFederatedRepositoryPackageParams("conan")}
}

type CondaFederatedRepositoryParams struct {
	FederatedRepositoryBaseParams
}

func NewCondaFederatedRepositoryParams() CondaFederatedRepositoryParams {
	return CondaFederatedRepositoryParams{FederatedRepositoryBaseParams: NewFederatedRepositoryPackageParams("conda")}
}

type CranFederatedRepositoryParams struct {
	FederatedRepositoryBaseParams
}

func NewCranFederatedRepositoryParams() CranFederatedRepositoryParams {
	return CranFederatedRepositoryParams{FederatedRepositoryBaseParams: NewFederatedRepositoryPackageParams("cran")}
}

type DebianFederatedRepositoryParams struct {
	FederatedRepositoryBaseParams
	DebianRepositoryParams
}

func NewDebianFederatedRepositoryParams() DebianFederatedRepositoryParams {
	return DebianFederatedRepositoryParams{FederatedRepositoryBaseParams: NewFederatedRepositoryPackageParams("debian")}
}

type DockerFederatedRepositoryParams struct {
	FederatedRepositoryBaseParams
	DockerRepositoryParams
}

func NewDockerFederatedRepositoryParams() DockerFederatedRepositoryParams {
	return DockerFederatedRepositoryParams{FederatedRepositoryBaseParams: NewFederatedRepositoryPackageParams("docker")}
}

type GemsFederatedRepositoryParams struct {
	FederatedRepositoryBaseParams
}

func NewGemsFederatedRepositoryParams() GemsFederatedRepositoryParams {
	return GemsFederatedRepositoryParams{FederatedRepositoryBaseParams: NewFederatedRepositoryPackageParams("gems")}
}

type GenericFederatedRepositoryParams struct {
	FederatedRepositoryBaseParams
}

func NewGenericFederatedRepositoryParams() GenericFederatedRepositoryParams {
	return GenericFederatedRepositoryParams{FederatedRepositoryBaseParams: NewFederatedRepositoryPackageParams("generic")}
}

type GitlfsFederatedRepositoryParams struct {
	FederatedRepositoryBaseParams
}

func NewGitlfsFederatedRepositoryParams() GitlfsFederatedRepositoryParams {
	return GitlfsFederatedRepositoryParams{FederatedRepositoryBaseParams: NewFederatedRepositoryPackageParams("gitlfs")}
}

type GoFederatedRepositoryParams struct {
	FederatedRepositoryBaseParams
}

func NewGoFederatedRepositoryParams() GoFederatedRepositoryParams {
	return GoFederatedRepositoryParams{FederatedRepositoryBaseParams: NewFederatedRepositoryPackageParams("go")}
}

type GradleFederatedRepositoryParams struct {
	FederatedRepositoryBaseParams
	JavaPackageManagersRepositoryParams
}

func NewGradleFederatedRepositoryParams() GradleFederatedRepositoryParams {
	return GradleFederatedRepositoryParams{FederatedRepositoryBaseParams: NewFederatedRepositoryPackageParams("gradle")}
}

type HelmFederatedRepositoryParams struct {
	FederatedRepositoryBaseParams
}

func NewHelmFederatedRepositoryParams() HelmFederatedRepositoryParams {
	return HelmFederatedRepositoryParams{FederatedRepositoryBaseParams: NewFederatedRepositoryPackageParams("helm")}
}

type IvyFederatedRepositoryParams struct {
	FederatedRepositoryBaseParams
	JavaPackageManagersRepositoryParams
}

func NewIvyFederatedRepositoryParams() IvyFederatedRepositoryParams {
	return IvyFederatedRepositoryParams{FederatedRepositoryBaseParams: NewFederatedRepositoryPackageParams("ivy")}
}

type MavenFederatedRepositoryParams struct {
	FederatedRepositoryBaseParams
	JavaPackageManagersRepositoryParams
}

func NewMavenFederatedRepositoryParams() MavenFederatedRepositoryParams {
	return MavenFederatedRepositoryParams{FederatedRepositoryBaseParams: NewFederatedRepositoryPackageParams("maven")}
}

type NpmFederatedRepositoryParams struct {
	FederatedRepositoryBaseParams
}

func NewNpmFederatedRepositoryParams() NpmFederatedRepositoryParams {
	return NpmFederatedRepositoryParams{FederatedRepositoryBaseParams: NewFederatedRepositoryPackageParams("npm")}
}

type NugetFederatedRepositoryParams struct {
	FederatedRepositoryBaseParams
	NugetRepositoryParams
}

func NewNugetFederatedRepositoryParams() NugetFederatedRepositoryParams {
	return NugetFederatedRepositoryParams{FederatedRepositoryBaseParams: NewFederatedRepositoryPackageParams("nuget")}
}

type OpkgFederatedRepositoryParams struct {
	FederatedRepositoryBaseParams
}

func NewOpkgFederatedRepositoryParams() OpkgFederatedRepositoryParams {
	return OpkgFederatedRepositoryParams{FederatedRepositoryBaseParams: NewFederatedRepositoryPackageParams("opkg")}
}

type PuppetFederatedRepositoryParams struct {
	FederatedRepositoryBaseParams
}

func NewPuppetFederatedRepositoryParams() PuppetFederatedRepositoryParams {
	return PuppetFederatedRepositoryParams{FederatedRepositoryBaseParams: NewFederatedRepositoryPackageParams("puppet")}
}

type PypiFederatedRepositoryParams struct {
	FederatedRepositoryBaseParams
}

func NewPypiFederatedRepositoryParams() PypiFederatedRepositoryParams {
	return PypiFederatedRepositoryParams{FederatedRepositoryBaseParams: NewFederatedRepositoryPackageParams("pypi")}
}

type RpmFederatedRepositoryParams struct {
	FederatedRepositoryBaseParams
	RpmRepositoryParams
}

func NewRpmFederatedRepositoryParams() RpmFederatedRepositoryParams {
	return RpmFederatedRepositoryParams{FederatedRepositoryBaseParams: NewFederatedRepositoryPackageParams("rpm")}
}

type SbtFederatedRepositoryParams struct {
	FederatedRepositoryBaseParams
	JavaPackageManagersRepositoryParams
}

func NewSbtFederatedRepositoryParams() SbtFederatedRepositoryParams {
	return SbtFederatedRepositoryParams{FederatedRepositoryBaseParams: NewFederatedRepositoryPackageParams("sbt")}
}

type SwiftFederatedRepositoryParams struct {
	FederatedRepositoryBaseParams
}

func NewSwiftFederatedRepositoryParams() SwiftFederatedRepositoryParams {
	return SwiftFederatedRepositoryParams{FederatedRepositoryBaseParams: NewFederatedRepositoryPackageParams("swift")}
}

type VagrantFederatedRepositoryParams struct {
	FederatedRepositoryBaseParams
}

func NewVagrantFederatedRepositoryParams() VagrantFederatedRepositoryParams {
	return VagrantFederatedRepositoryParams{FederatedRepositoryBaseParams: NewFederatedRepositoryPackageParams("vagrant")}
}

type YumFederatedRepositoryParams struct {
	FederatedRepositoryBaseParams
	RpmRepositoryParams
}

func NewYumFederatedRepositoryParams() YumFederatedRepositoryParams {
	return YumFederatedRepositoryParams{FederatedRepositoryBaseParams: NewFederatedRepositoryPackageParams("yum")}
}
