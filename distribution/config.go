package distribution

import (
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type Config interface {
	GetDistributionDetails() auth.DistributionDetails
	IsDryRun() bool
	GetLogger() log.Log
}

type distributionServicesConfig struct {
	auth.DistributionDetails
	dryRun   bool
	logger   log.Log
}

func (config *distributionServicesConfig) IsDryRun() bool {
	return config.dryRun
}

func (config *distributionServicesConfig) GetDistributionDetails() auth.DistributionDetails {
	return config.DistributionDetails
}

func (config *distributionServicesConfig) GetLogger() log.Log {
	return config.logger
}
