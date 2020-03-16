package config

import (
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type Config interface {
	GetCertificatesPath() string
	GetThreads() int
	IsDryRun() bool
	GetCommonDetails() auth.CommonDetails
	GetLogger() log.Log
	IsInsecureTls() bool
}

type servicesConfig struct {
	auth.CommonDetails
	certificatesPath string
	dryRun           bool
	threads          int
	logger           log.Log
	insecureTls      bool
}

func (config *servicesConfig) IsDryRun() bool {
	return config.dryRun
}

func (config *servicesConfig) GetCertificatesPath() string {
	return config.certificatesPath
}

func (config *servicesConfig) GetThreads() int {
	return config.threads
}

func (config *servicesConfig) GetCommonDetails() auth.CommonDetails {
	return config.CommonDetails
}

func (config *servicesConfig) GetLogger() log.Log {
	return config.logger
}

func (config *servicesConfig) IsInsecureTls() bool {
	return config.insecureTls
}
