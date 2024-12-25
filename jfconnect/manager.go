package jfconnect

import (
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/jfconnect/services"
)

type Manager interface {
	PostVisibilityMetric(services.VisibilityMetric) error
}

type jfConnectManager struct {
	client *jfroghttpclient.JfrogHttpClient
	config config.Config
}

func NewManager(config config.Config) (Manager, error) {
	details := config.GetServiceDetails()
	var err error
	manager := &jfConnectManager{config: config}
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

func (jm *jfConnectManager) Client() *jfroghttpclient.JfrogHttpClient {
	return jm.client
}

func (jm *jfConnectManager) PostVisibilityMetric(metric services.VisibilityMetric) error {
	jfConnectService := services.NewJfConnectService(jm.config.GetServiceDetails(), jm.client)
	return jfConnectService.PostVisibilityMetric(metric)
}
