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
		log.Debug("JFrog Distribution version is:", ds.Version)
	}
	return ds.Version, nil
}

func (ds *distributionDetails) getDistributionVersion() (string, error) {
	cd := auth.ServiceDetails(ds)
	serviceConfig, err := config.NewConfigBuilder().
		SetServiceDetails(cd).
		SetCertificatesPath(cd.GetClientCertPath()).
		Build()
	if err != nil {
		return "", err
	}
	sm, err := distribution.New(serviceConfig)
	if err != nil {
		return "", err
	}
	return sm.GetDistributionVersion()
}
