package artifactory

import (
	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type Config interface {
	GetUrl() string
	GetPassword() string
	GetApiKey() string
	GetCertificatesPath() string
	GetThreads() int
	IsDryRun() bool
	GetArtDetails() auth.ArtifactoryDetails
	GetLogger() log.Log
	IsInsecureTls() bool
}

type artifactoryServicesConfig struct {
	auth.ArtifactoryDetails
	certificatesPath string
	dryRun           bool
	threads          int
	logger           log.Log
	insecureTls      bool
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

func (config *artifactoryServicesConfig) GetCertificatesPath() string {
	return config.certificatesPath
}

func (config *artifactoryServicesConfig) GetThreads() int {
	return config.threads
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
