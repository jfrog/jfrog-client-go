package auth

import (
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/jfrog/jfrog-client-go/xsc"
)

// NewXscDetails creates a struct of the Xsc details
func NewXscDetails() *XscDetails {
	return &XscDetails{}
}

type XscDetails struct {
	auth.CommonConfigFields
}

func (ds *XscDetails) GetVersion() (string, error) {
	var err error
	if ds.Version == "" {
		ds.Version, err = ds.getXscVersion()
		if err != nil {
			return "", err
		}
		log.Debug("The Xsc version is:", ds.Version)
	}
	return ds.Version, nil
}

func (ds *XscDetails) getXscVersion() (string, error) {
	cd := auth.ServiceDetails(ds)
	serviceConfig, err := config.NewConfigBuilder().
		SetServiceDetails(cd).
		SetCertificatesPath(cd.GetClientCertPath()).
		Build()
	if err != nil {
		return "", err
	}
	sm, err := xsc.New(serviceConfig)
	if err != nil {
		return "", err
	}
	return sm.GetVersion()
}
