package services

import (
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
)

const LocalRepositoryRepoType = "local"

type LocalRepositoryService struct {
	RepositoryService
}

func NewLocalRepositoryService(client *jfroghttpclient.JfrogHttpClient, isUpdate bool) *LocalRepositoryService {
	return &LocalRepositoryService{
		RepositoryService: RepositoryService{
			repoType: LocalRepositoryRepoType,
			client:   client,
			isUpdate: isUpdate,
		},
	}
}

func (lrs *LocalRepositoryService) GetJfrogHttpClient() *jfroghttpclient.JfrogHttpClient {
	return lrs.client
}

func (lrs *LocalRepositoryService) performRequest(params interface{}, repoKey string) error {
	return lrs.RepositoryService.performRequest(params, repoKey)
}

func (lrs *LocalRepositoryService) Alpine(params AlpineLocalRepositoryParams) error {
	return lrs.performRequest(params, params.Key)
}

func (lrs *LocalRepositoryService) Bower(params BowerLocalRepositoryParams) error {
	return lrs.performRequest(params, params.Key)
}

func (lrs *LocalRepositoryService) Cargo(params CargoLocalRepositoryParams) error {
	return lrs.performRequest(params, params.Key)
}

func (lrs *LocalRepositoryService) Chef(params ChefLocalRepositoryParams) error {
	return lrs.performRequest(params, params.Key)
}

func (lrs *LocalRepositoryService) Cocoapods(params CocoapodsLocalRepositoryParams) error {
	return lrs.performRequest(params, params.Key)
}

func (lrs *LocalRepositoryService) Composer(params ComposerLocalRepositoryParams) error {
	return lrs.performRequest(params, params.Key)
}

func (lrs *LocalRepositoryService) Conan(params ConanLocalRepositoryParams) error {
	return lrs.performRequest(params, params.Key)
}

func (lrs *LocalRepositoryService) Conda(params CondaLocalRepositoryParams) error {
	return lrs.performRequest(params, params.Key)
}

func (lrs *LocalRepositoryService) Cran(params CranLocalRepositoryParams) error {
	return lrs.performRequest(params, params.Key)
}

func (lrs *LocalRepositoryService) Debian(params DebianLocalRepositoryParams) error {
	return lrs.performRequest(params, params.Key)
}

func (lrs *LocalRepositoryService) Docker(params DockerLocalRepositoryParams) error {
	return lrs.performRequest(params, params.Key)
}

func (lrs *LocalRepositoryService) Gems(params GemsLocalRepositoryParams) error {
	return lrs.performRequest(params, params.Key)
}

func (lrs *LocalRepositoryService) Generic(params GenericLocalRepositoryParams) error {
	return lrs.performRequest(params, params.Key)
}

func (lrs *LocalRepositoryService) Gitlfs(params GitlfsLocalRepositoryParams) error {
	return lrs.performRequest(params, params.Key)
}

func (lrs *LocalRepositoryService) Go(params GoLocalRepositoryParams) error {
	return lrs.performRequest(params, params.Key)
}

func (lrs *LocalRepositoryService) Gradle(params GradleLocalRepositoryParams) error {
	return lrs.performRequest(params, params.Key)
}

func (lrs *LocalRepositoryService) Helm(params HelmLocalRepositoryParams) error {
	return lrs.performRequest(params, params.Key)
}

func (lrs *LocalRepositoryService) Ivy(params IvyLocalRepositoryParams) error {
	return lrs.performRequest(params, params.Key)
}

func (lrs *LocalRepositoryService) Maven(params MavenLocalRepositoryParams) error {
	return lrs.performRequest(params, params.Key)
}

func (lrs *LocalRepositoryService) Npm(params NpmLocalRepositoryParams) error {
	return lrs.performRequest(params, params.Key)
}

func (lrs *LocalRepositoryService) Nuget(params NugetLocalRepositoryParams) error {
	return lrs.performRequest(params, params.Key)
}

func (lrs *LocalRepositoryService) Opkg(params OpkgLocalRepositoryParams) error {
	return lrs.performRequest(params, params.Key)
}

func (lrs *LocalRepositoryService) Puppet(params PuppetLocalRepositoryParams) error {
	return lrs.performRequest(params, params.Key)
}

func (lrs *LocalRepositoryService) Pypi(params PypiLocalRepositoryParams) error {
	return lrs.performRequest(params, params.Key)
}

func (lrs *LocalRepositoryService) Rpm(params RpmLocalRepositoryParams) error {
	return lrs.performRequest(params, params.Key)
}

func (lrs *LocalRepositoryService) Sbt(params SbtLocalRepositoryParams) error {
	return lrs.performRequest(params, params.Key)
}

func (lrs *LocalRepositoryService) Swift(params SwiftLocalRepositoryParams) error {
	return lrs.performRequest(params, params.Key)
}

func (lrs *LocalRepositoryService) Vagrant(params VagrantLocalRepositoryParams) error {
	return lrs.performRequest(params, params.Key)
}

func (lrs *LocalRepositoryService) Yum(params YumLocalRepositoryParams) error {
	return lrs.performRequest(params, params.Key)
}

type LocalRepositoryBaseParams struct {
	RepositoryBaseParams
	AdditionalRepositoryBaseParams
	KeyPairRefsRepositoryParams
	ArchiveBrowsingEnabled *bool `json:"archiveBrowsingEnabled,omitempty"`
}

func NewLocalRepositoryBaseParams() LocalRepositoryBaseParams {
	return LocalRepositoryBaseParams{RepositoryBaseParams: RepositoryBaseParams{Rclass: LocalRepositoryRepoType}}
}

func NewLocalRepositoryPackageParams(packageType string) LocalRepositoryBaseParams {
	return LocalRepositoryBaseParams{RepositoryBaseParams: RepositoryBaseParams{Rclass: LocalRepositoryRepoType, PackageType: packageType}}
}

type AlpineLocalRepositoryParams struct {
	LocalRepositoryBaseParams
}

func NewAlpineLocalRepositoryParams() AlpineLocalRepositoryParams {
	return AlpineLocalRepositoryParams{LocalRepositoryBaseParams: NewLocalRepositoryPackageParams("alpine")}
}

type BowerLocalRepositoryParams struct {
	LocalRepositoryBaseParams
}

func NewBowerLocalRepositoryParams() BowerLocalRepositoryParams {
	return BowerLocalRepositoryParams{LocalRepositoryBaseParams: NewLocalRepositoryPackageParams("bower")}
}

type CargoLocalRepositoryParams struct {
	LocalRepositoryBaseParams
	CargoRepositoryParams
}

func NewCargoLocalRepositoryParams() CargoLocalRepositoryParams {
	return CargoLocalRepositoryParams{LocalRepositoryBaseParams: NewLocalRepositoryPackageParams("cargo")}
}

type ChefLocalRepositoryParams struct {
	LocalRepositoryBaseParams
}

func NewChefLocalRepositoryParams() ChefLocalRepositoryParams {
	return ChefLocalRepositoryParams{LocalRepositoryBaseParams: NewLocalRepositoryPackageParams("chef")}
}

type CocoapodsLocalRepositoryParams struct {
	LocalRepositoryBaseParams
}

func NewCocoapodsLocalRepositoryParams() CocoapodsLocalRepositoryParams {
	return CocoapodsLocalRepositoryParams{LocalRepositoryBaseParams: NewLocalRepositoryPackageParams("cocoapods")}
}

type ComposerLocalRepositoryParams struct {
	LocalRepositoryBaseParams
}

func NewComposerLocalRepositoryParams() ComposerLocalRepositoryParams {
	return ComposerLocalRepositoryParams{LocalRepositoryBaseParams: NewLocalRepositoryPackageParams("composer")}
}

type ConanLocalRepositoryParams struct {
	LocalRepositoryBaseParams
}

func NewConanLocalRepositoryParams() ConanLocalRepositoryParams {
	return ConanLocalRepositoryParams{LocalRepositoryBaseParams: NewLocalRepositoryPackageParams("conan")}
}

type CondaLocalRepositoryParams struct {
	LocalRepositoryBaseParams
}

func NewCondaLocalRepositoryParams() CondaLocalRepositoryParams {
	return CondaLocalRepositoryParams{LocalRepositoryBaseParams: NewLocalRepositoryPackageParams("conda")}
}

type CranLocalRepositoryParams struct {
	LocalRepositoryBaseParams
}

func NewCranLocalRepositoryParams() CranLocalRepositoryParams {
	return CranLocalRepositoryParams{LocalRepositoryBaseParams: NewLocalRepositoryPackageParams("cran")}
}

type DebianLocalRepositoryParams struct {
	LocalRepositoryBaseParams
	DebianRepositoryParams
}

func NewDebianLocalRepositoryParams() DebianLocalRepositoryParams {
	return DebianLocalRepositoryParams{LocalRepositoryBaseParams: NewLocalRepositoryPackageParams("debian")}
}

type DockerLocalRepositoryParams struct {
	LocalRepositoryBaseParams
	DockerRepositoryParams
}

func NewDockerLocalRepositoryParams() DockerLocalRepositoryParams {
	return DockerLocalRepositoryParams{LocalRepositoryBaseParams: NewLocalRepositoryPackageParams("docker")}
}

type GemsLocalRepositoryParams struct {
	LocalRepositoryBaseParams
}

func NewGemsLocalRepositoryParams() GemsLocalRepositoryParams {
	return GemsLocalRepositoryParams{LocalRepositoryBaseParams: NewLocalRepositoryPackageParams("gems")}
}

type GenericLocalRepositoryParams struct {
	LocalRepositoryBaseParams
}

func NewGenericLocalRepositoryParams() GenericLocalRepositoryParams {
	return GenericLocalRepositoryParams{LocalRepositoryBaseParams: NewLocalRepositoryPackageParams("generic")}
}

type GitlfsLocalRepositoryParams struct {
	LocalRepositoryBaseParams
}

func NewGitlfsLocalRepositoryParams() GitlfsLocalRepositoryParams {
	return GitlfsLocalRepositoryParams{LocalRepositoryBaseParams: NewLocalRepositoryPackageParams("gitlfs")}
}

type GoLocalRepositoryParams struct {
	LocalRepositoryBaseParams
}

func NewGoLocalRepositoryParams() GoLocalRepositoryParams {
	return GoLocalRepositoryParams{LocalRepositoryBaseParams: NewLocalRepositoryPackageParams("go")}
}

type GradleLocalRepositoryParams struct {
	LocalRepositoryBaseParams
	JavaPackageManagersRepositoryParams
}

func NewGradleLocalRepositoryParams() GradleLocalRepositoryParams {
	return GradleLocalRepositoryParams{LocalRepositoryBaseParams: NewLocalRepositoryPackageParams("gradle")}
}

type HelmLocalRepositoryParams struct {
	LocalRepositoryBaseParams
}

func NewHelmLocalRepositoryParams() HelmLocalRepositoryParams {
	return HelmLocalRepositoryParams{LocalRepositoryBaseParams: NewLocalRepositoryPackageParams("helm")}
}

type IvyLocalRepositoryParams struct {
	LocalRepositoryBaseParams
	JavaPackageManagersRepositoryParams
}

func NewIvyLocalRepositoryParams() IvyLocalRepositoryParams {
	return IvyLocalRepositoryParams{LocalRepositoryBaseParams: NewLocalRepositoryPackageParams("ivy")}
}

type MavenLocalRepositoryParams struct {
	LocalRepositoryBaseParams
	JavaPackageManagersRepositoryParams
}

func NewMavenLocalRepositoryParams() MavenLocalRepositoryParams {
	return MavenLocalRepositoryParams{LocalRepositoryBaseParams: NewLocalRepositoryPackageParams("maven")}
}

type NpmLocalRepositoryParams struct {
	LocalRepositoryBaseParams
}

func NewNpmLocalRepositoryParams() NpmLocalRepositoryParams {
	return NpmLocalRepositoryParams{LocalRepositoryBaseParams: NewLocalRepositoryPackageParams("npm")}
}

type NugetLocalRepositoryParams struct {
	LocalRepositoryBaseParams
	NugetRepositoryParams
}

func NewNugetLocalRepositoryParams() NugetLocalRepositoryParams {
	return NugetLocalRepositoryParams{LocalRepositoryBaseParams: NewLocalRepositoryPackageParams("nuget")}
}

type OpkgLocalRepositoryParams struct {
	LocalRepositoryBaseParams
}

func NewOpkgLocalRepositoryParams() OpkgLocalRepositoryParams {
	return OpkgLocalRepositoryParams{LocalRepositoryBaseParams: NewLocalRepositoryPackageParams("opkg")}
}

type PuppetLocalRepositoryParams struct {
	LocalRepositoryBaseParams
}

func NewPuppetLocalRepositoryParams() PuppetLocalRepositoryParams {
	return PuppetLocalRepositoryParams{LocalRepositoryBaseParams: NewLocalRepositoryPackageParams("puppet")}
}

type PypiLocalRepositoryParams struct {
	LocalRepositoryBaseParams
}

func NewPypiLocalRepositoryParams() PypiLocalRepositoryParams {
	return PypiLocalRepositoryParams{LocalRepositoryBaseParams: NewLocalRepositoryPackageParams("pypi")}
}

type RpmLocalRepositoryParams struct {
	LocalRepositoryBaseParams
	RpmRepositoryParams
}

func NewRpmLocalRepositoryParams() RpmLocalRepositoryParams {
	return RpmLocalRepositoryParams{LocalRepositoryBaseParams: NewLocalRepositoryPackageParams("rpm")}
}

type SbtLocalRepositoryParams struct {
	LocalRepositoryBaseParams
	JavaPackageManagersRepositoryParams
}

func NewSbtLocalRepositoryParams() SbtLocalRepositoryParams {
	return SbtLocalRepositoryParams{LocalRepositoryBaseParams: NewLocalRepositoryPackageParams("sbt")}
}

type SwiftLocalRepositoryParams struct {
	LocalRepositoryBaseParams
}

func NewSwiftLocalRepositoryParams() SwiftLocalRepositoryParams {
	return SwiftLocalRepositoryParams{LocalRepositoryBaseParams: NewLocalRepositoryPackageParams("swift")}
}

type VagrantLocalRepositoryParams struct {
	LocalRepositoryBaseParams
}

func NewVagrantLocalRepositoryParams() VagrantLocalRepositoryParams {
	return VagrantLocalRepositoryParams{LocalRepositoryBaseParams: NewLocalRepositoryPackageParams("vagrant")}
}

type YumLocalRepositoryParams struct {
	LocalRepositoryBaseParams
	RpmRepositoryParams
}

func NewYumLocalRepositoryParams() YumLocalRepositoryParams {
	return YumLocalRepositoryParams{LocalRepositoryBaseParams: NewLocalRepositoryPackageParams("yum")}
}
