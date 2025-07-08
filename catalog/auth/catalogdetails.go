package auth

import (
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/catalog"
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

// NewCatalogDetails creates a struct of the Catalog details
func NewCatalogDetails() *catalogDetails {
	return &catalogDetails{}
}

type catalogDetails struct {
	auth.CommonConfigFields
}

func (cs *catalogDetails) GetVersion() (string, error) {
	var err error
	if cs.Version == "" {
		cs.Version, err = cs.getCatalogVersion()
		if err != nil {
			return "", err
		}
		log.Debug("JFrog Catalog version is:", cs.Version)
	}
	return cs.Version, nil
}

func (ds *catalogDetails) getCatalogVersion() (string, error) {
	cd := auth.ServiceDetails(ds)
	serviceConfig, err := config.NewConfigBuilder().
		SetServiceDetails(cd).
		SetCertificatesPath(cd.GetClientCertPath()).
		Build()
	if err != nil {
		return "", err
	}
	cm, err := catalog.New(serviceConfig)
	if cm != nil {
		return "", err
	}
	return cm.GetVersion()
}
