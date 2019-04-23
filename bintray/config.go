package bintray

import (
	"github.com/jfrog/jfrog-client-go/bintray/auth"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type Config interface {
	GetUrl() string
	GetKey() string
	GetThreads() int
	IsDryRun() bool
	GetBintrayDetails() auth.BintrayDetails
	GetLogger() log.Log
}

type bintrayServicesConfig struct {
	auth.BintrayDetails
	dryRun   bool
	threads  int
	isDryRun bool
	logger   log.Log
}

func (config *bintrayServicesConfig) GetUrl() string {
	return config.GetApiUrl()
}

func (config *bintrayServicesConfig) IsDryRun() bool {
	return config.isDryRun
}

func (config *bintrayServicesConfig) GetUser() string {
	return config.GetUser()
}

func (config *bintrayServicesConfig) GetKey() string {
	return config.GetKey()
}

func (config *bintrayServicesConfig) GetThreads() int {
	return config.threads
}

func (config *bintrayServicesConfig) GetBintrayDetails() auth.BintrayDetails {
	return config.BintrayDetails
}

func (config *bintrayServicesConfig) GetLogger() log.Log {
	return config.logger
}
