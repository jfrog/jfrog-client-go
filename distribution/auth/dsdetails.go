package auth

import (
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/jfrog/jfrog-client-go/distribution"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

func NewDistributionDetails() *distributionDetails {
	return &distributionDetails{}
}

type distributionDetails struct {
	auth.CommonConfigFields
}

func (ds *distributionDetails) GetVersion() (string, error) {
	var err error
	if ds.Version == "" {
		ds.Version, err = ds.getDistributionVersion()
		if err != nil {
			return "", err
		}
		log.Debug("The Distribution version is:", ds.Version)
	}
	return ds.Version, nil
}

func (ds *distributionDetails) getDistributionVersion() (string, error) {
	cd := auth.CommonDetails(ds)
	serviceConfig, err := config.NewConfigBuilder().
		SetArtDetails(cd).
		SetCertificatesPath(cd.GetClientCertPath()).
		SetThreads(0).
		Build()
	sm, err := distribution.New(&cd, serviceConfig)
	if err != nil {
		return "", err
	}
	return sm.GetDistributionVersion()
}
