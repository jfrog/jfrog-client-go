package artifactory

import (
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type Config interface {
	GetCertificatesPath() string
	GetThreads() int
	IsDryRun() bool
	GetArtDetails() auth.CommonDetails
	GetLogger() log.Log
	IsInsecureTls() bool
}

type artifactoryServicesConfig struct {
	auth.CommonDetails
	certificatesPath string
	dryRun           bool
	threads          int
	logger           log.Log
	insecureTls      bool
}

func (config *artifactoryServicesConfig) IsDryRun() bool {
	return config.dryRun
}

func (config *artifactoryServicesConfig) GetCertificatesPath() string {
	return config.certificatesPath
}

func (config *artifactoryServicesConfig) GetThreads() int {
	return config.threads
}

func (config *artifactoryServicesConfig) GetArtDetails() auth.CommonDetails {
	return config.CommonDetails
}

func (config *artifactoryServicesConfig) GetLogger() log.Log {
	return config.logger
}

func (config *artifactoryServicesConfig) IsInsecureTls() bool {
	return config.insecureTls
}
