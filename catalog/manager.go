package catalog

import (
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/catalog/services"
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
)

// CatalogServicesManager defines the http client and general configuration
type CatalogServicesManager struct {
	client *jfroghttpclient.JfrogHttpClient
	config config.Config
	// Global reference to the provided project key, used for API endpoints that require it for authentication
	scopeProjectKey string
}

// New creates a service manager to interact with Catalog
func New(config config.Config) (*CatalogServicesManager, error) {
	details := config.GetServiceDetails()
	var err error
	manager := &CatalogServicesManager{config: config}
	manager.client, err = buildJFrogHttpClient(config, details)
	return manager, err
}

func buildJFrogHttpClient(config config.Config, authDetails auth.ServiceDetails) (*jfroghttpclient.JfrogHttpClient, error) {
	return jfroghttpclient.JfrogClientBuilder().
		SetCertificatesPath(config.GetCertificatesPath()).
		SetInsecureTls(config.IsInsecureTls()).
		SetContext(config.GetContext()).
		SetDialTimeout(config.GetDialTimeout()).
		SetOverallRequestTimeout(config.GetOverallRequestTimeout()).
		SetClientCertPath(authDetails.GetClientCertPath()).
		SetClientCertKeyPath(authDetails.GetClientCertKeyPath()).
		AppendPreRequestInterceptor(authDetails.RunPreRequestFunctions).
		SetRetries(config.GetHttpRetries()).
		SetRetryWaitMilliSecs(config.GetHttpRetryWaitMilliSecs()).
		Build()
}

// GetVersion will return the Catalog version
func (cm *CatalogServicesManager) GetVersion() (string, error) {
	versionService := services.NewVersionService(cm.client)
	versionService.CatalogDetails = cm.config.GetServiceDetails()
	return versionService.GetVersion()
}
