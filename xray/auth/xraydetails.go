package auth

import (
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/jfrog/jfrog-client-go/xray/manager"
)

// NewXrayDetails creates a struct of the Xray details
func NewXrayDetails() *XrayDetails {
	return &XrayDetails{}
}

type XrayDetails struct {
	auth.CommonConfigFields
}

type XscDetails struct {
	auth.CommonConfigFields
	XscUrl string
}

func (ds *XrayDetails) GetVersion() (string, error) {
	var err error
	if ds.Version == "" {
		ds.Version, err = ds.getXrayVersion()
		if err != nil {
			return "", err
		}
		log.Debug("JFrog Xray version is:", ds.Version)
	}
	return ds.Version, nil
}

func (ds *XrayDetails) getXrayVersion() (string, error) {
	cd := auth.ServiceDetails(ds)
	serviceConfig, err := config.NewConfigBuilder().
		SetServiceDetails(cd).
		SetCertificatesPath(cd.GetClientCertPath()).
		Build()
	if err != nil {
		return "", err
	}
	sm, err := manager.New(serviceConfig)
	if err != nil {
		return "", err
	}
	return sm.GetVersion()
}

func (ds *XrayDetails) GetXscUrl() string {
	return ds.XscUrl
}

func (ds *XrayDetails) SetXscUrl(url string) {
	ds.XscUrl = url
}

func (ds *XrayDetails) GetPlatformUrl() string {
	return ds.PlatformUrl
}
