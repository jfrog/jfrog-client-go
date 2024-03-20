package auth

import (
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/jfrog/jfrog-client-go/pipelines"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

func NewPipelinesDetails() auth.ServiceDetails {
	return &pipelinesDetails{}
}

type pipelinesDetails struct {
	auth.CommonConfigFields
}

func (pd *pipelinesDetails) GetVersion() (string, error) {
	var err error
	if pd.Version == "" {
		pd.Version, err = pd.getPipelinesVersion()
		if err != nil {
			return "", err
		}
		log.Debug("JFrog Pipelines version is:", pd.Version)
	}
	return pd.Version, nil
}

func (pd *pipelinesDetails) getPipelinesVersion() (string, error) {
	cd := auth.ServiceDetails(pd)
	serviceConfig, err := config.NewConfigBuilder().
		SetServiceDetails(cd).
		SetCertificatesPath(cd.GetClientCertPath()).
		Build()
	if err != nil {
		return "", err
	}
	sm, err := pipelines.New(serviceConfig)
	if err != nil {
		return "", err
	}
	sys, err := sm.GetSystemInfo()
	if err != nil {
		return "", err
	}
	return sys.Version, nil
}
