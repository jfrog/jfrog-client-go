package auth

import (
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/jfrog/jfrog-client-go/xray"
)

// NewXrayDetails creates a struct of the Xray details
func NewXrayDetails() *xrayDetails {
	return &xrayDetails{}
}

type xrayDetails struct {
	auth.CommonConfigFields
}

func (ds *xrayDetails) GetVersion() (string, error) {
	var err error
	if ds.Version == "" {
		ds.Version, err = ds.getXrayVersion()
		if err != nil {
			return "", err
		}
		log.Debug("The Xray version is:", ds.Version)
	}
	return ds.Version, nil
}

func (ds *xrayDetails) getXrayVersion() (string, error) {
	cd := auth.ServiceDetails(ds)
	serviceConfig, err := config.NewConfigBuilder().
		SetServiceDetails(cd).
		SetCertificatesPath(cd.GetClientCertPath()).
		Build()
	sm, err := xray.New(serviceConfig)
	if err != nil {
		return "", err
	}
	return sm.GetVersion()
}
