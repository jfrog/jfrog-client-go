package auth

import (
	"github.com/jfrog/jfrog-client-go/apptrust"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

func NewXscDetails() *ApptrustDetails {
	return &ApptrustDetails{}
}

type ApptrustDetails struct {
	auth.CommonConfigFields
}

func (ds *ApptrustDetails) GetVersion() (string, error) {
	var err error
	if ds.Version == "" {
		ds.Version, err = ds.getApptrustVersion()
		if err != nil {
			return "", err
		}
		log.Debug("JFrog AppTrust version is:", ds.Version)
	}
	return ds.Version, nil
}

func (ds *ApptrustDetails) getApptrustVersion() (string, error) {
	cd := auth.ServiceDetails(ds)
	serviceConfig, err := config.NewConfigBuilder().
		SetServiceDetails(cd).
		SetCertificatesPath(cd.GetClientCertPath()).
		Build()
	if err != nil {
		return "", err
	}
	sm, err := apptrust.New(serviceConfig)
	if err != nil {
		return "", err
	}
	return sm.GetVersion()
}
