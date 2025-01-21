package auth

import (
	"github.com/jfrog/jfrog-client-go/application"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

func NewApplicationDetails() auth.ServiceDetails {
	return &applicationDetails{}
}

type applicationDetails struct {
	auth.CommonConfigFields
}

func (ap *applicationDetails) GetVersion() (string, error) {
	var err error
	if ap.Version == "" {
		ap.Version, err = ap.getApplicationVersion()
		if err != nil {
			return "", err
		}
		log.Debug("JFrog Application version is:", ap.Version)
	}
	return ap.Version, nil
}

func (ap *applicationDetails) getApplicationVersion() (string, error) {
	cd := auth.ServiceDetails(ap)
	serviceConfig, err := config.NewConfigBuilder().
		SetServiceDetails(cd).
		SetCertificatesPath(cd.GetClientCertPath()).
		Build()
	if err != nil {
		return "", err
	}
	sm, err := application.New(serviceConfig)
	if err != nil {
		return "", err
	}
	return sm.GetVersion()
}
