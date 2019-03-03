package artifactory

import (
	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type Config interface {
	GetUrl() string
	GetPassword() string
	GetApiKey() string
	GetCertifactesPath() string
	GetThreads() int
	GetMinSplitSize() int64
	GetSplitCount() int
	GetMinChecksumDeploy() int64
	IsDryRun() bool
	GetArtDetails() auth.ArtifactoryDetails
	GetLogger() log.Log
	IsInsecureTls() bool
}

type ArtifactoryServicesSetter interface {
	SetThread(threads int)
	SetArtDetails(artDetails auth.ArtifactoryDetails)
	SetDryRun(isDryRun bool)
}

type artifactoryServicesConfig struct {
	auth.ArtifactoryDetails
	certifactesPath   string
	dryRun            bool
	threads           int
	minSplitSize      int64
	splitCount        int
	minChecksumDeploy int64
	logger            log.Log
	insecureTls       bool
}

func (config *artifactoryServicesConfig) GetUrl() string {
	return config.GetUrl()
}

func (config *artifactoryServicesConfig) IsDryRun() bool {
	return config.dryRun
}

func (config *artifactoryServicesConfig) GetPassword() string {
	return config.GetPassword()
}

func (config *artifactoryServicesConfig) GetApiKey() string {
	return config.GetApiKey()
}

func (config *artifactoryServicesConfig) GetCertifactesPath() string {
	return config.certifactesPath
}

func (config *artifactoryServicesConfig) GetThreads() int {
	return config.threads
}

func (config *artifactoryServicesConfig) GetMinSplitSize() int64 {
	return config.minSplitSize
}

func (config *artifactoryServicesConfig) GetSplitCount() int {
	return config.splitCount
}
func (config *artifactoryServicesConfig) GetMinChecksumDeploy() int64 {
	return config.minChecksumDeploy
}

func (config *artifactoryServicesConfig) GetArtDetails() auth.ArtifactoryDetails {
	return config.ArtifactoryDetails
}

func (config *artifactoryServicesConfig) GetLogger() log.Log {
	return config.logger
}

func (config *artifactoryServicesConfig) IsInsecureTls() bool {
	return config.insecureTls
}
