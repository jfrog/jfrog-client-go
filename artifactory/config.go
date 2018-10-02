package artifactory

import (
	"github.com/jfrog/jfrog-cli-go/jfrog-client/artifactory/auth"
	"github.com/jfrog/jfrog-cli-go/jfrog-client/utils/log"
)

type artifactoryServicesConfig struct {
	auth.ArtifactoryDetails
	certifactesPath   string
	dryRun            bool
	threads           int
	minSplitSize      int64
	splitCount        int
	minChecksumDeploy int64
	logger            log.Log
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
