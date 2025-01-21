package application

import (
	"github.com/jfrog/jfrog-client-go/application/services"
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
)

// ApplicationServicesManager defines the http client and general configuration
type ApplicationServicesManager struct {
	client *jfroghttpclient.JfrogHttpClient
	config config.Config
}

// New creates a service manager to interact with Application
func New(config config.Config) (*ApplicationServicesManager, error) {
	details := config.GetServiceDetails()
	var err error
	manager := &ApplicationServicesManager{config: config}
	manager.client, err = jfroghttpclient.JfrogClientBuilder().
		SetCertificatesPath(config.GetCertificatesPath()).
		SetInsecureTls(config.IsInsecureTls()).
		SetContext(config.GetContext()).
		SetDialTimeout(config.GetDialTimeout()).
		SetOverallRequestTimeout(config.GetOverallRequestTimeout()).
		SetClientCertPath(details.GetClientCertPath()).
		SetClientCertKeyPath(details.GetClientCertKeyPath()).
		AppendPreRequestInterceptor(details.RunPreRequestFunctions).
		SetRetries(config.GetHttpRetries()).
		SetRetryWaitMilliSecs(config.GetHttpRetryWaitMilliSecs()).
		Build()
	return manager, err
}

func (ap *ApplicationServicesManager) Client() *jfroghttpclient.JfrogHttpClient {
	return ap.client
}

func (ap *ApplicationServicesManager) Config() config.Config {
	return ap.config
}

func (ap *ApplicationServicesManager) AddCommitInfo(commitInfo services.CreateApplicationCommitInfo, applicationKey string) error {
	commitInfoService := services.NewCommitInfoService(ap.client, ap.config.GetServiceDetails())
	return commitInfoService.AddCommitInfo(applicationKey, commitInfo)
}

func (ap *ApplicationServicesManager) GetVersion() (string, error) {
	versionService := services.NewVersionService(ap.client)
	versionService.ApplicationDetails = ap.config.GetServiceDetails()
	return versionService.GetVersion()
}
