package metadata

import (
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/metadata/services"
)

type Manager interface {
	GraphqlQuery(query []byte) ([]byte, error)
}

type metadataManager struct {
	client *jfroghttpclient.JfrogHttpClient
	config config.Config
}

func NewManager(config config.Config) (Manager, error) {
	details := config.GetServiceDetails()
	var err error
	manager := &metadataManager{config: config}
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

func (mm *metadataManager) Client() *jfroghttpclient.JfrogHttpClient {
	return mm.client
}

func (mm *metadataManager) GraphqlQuery(query []byte) ([]byte, error) {
	metadataService := services.NewMetadataService(mm.config.GetServiceDetails(), mm.client)
	return metadataService.Query(query)
}
