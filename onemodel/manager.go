package onemodel

import (
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/onemodel/services"
)

type Manager interface {
	GraphqlQuery(query []byte) ([]byte, error)
}

type onemodelManager struct {
	client *jfroghttpclient.JfrogHttpClient
	config config.Config
}

func NewManager(config config.Config) (Manager, error) {
	details := config.GetServiceDetails()
	var err error
	manager := &onemodelManager{config: config}
	manager.client, err = jfroghttpclient.JfrogClientBuilder().
		SetCertificatesPath(config.GetCertificatesPath()).
		SetInsecureTls(config.IsInsecureTls()).
		SetClientCertPath(details.GetClientCertPath()).
		SetClientCertKeyPath(details.GetClientCertKeyPath()).
		AppendPreRequestInterceptor(details.RunPreRequestFunctions).
		SetContext(config.GetContext()).
		SetDialTimeout(config.GetDialTimeout()).
		SetOverallRequestTimeout(config.GetOverallRequestTimeout()).
		SetRetries(config.GetHttpRetries()).
		SetRetryWaitMilliSecs(config.GetHttpRetryWaitMilliSecs()).
		Build()

	return manager, err
}

func (omm *onemodelManager) Client() *jfroghttpclient.JfrogHttpClient {
	return omm.client
}

func (omm *onemodelManager) GraphqlQuery(query []byte) ([]byte, error) {
	onemodelService := services.NewOnemodelService(omm.config.GetServiceDetails(), omm.client)
	return onemodelService.Query(query)
}
