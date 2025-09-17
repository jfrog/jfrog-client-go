package apptrust

import (
	"github.com/jfrog/jfrog-client-go/apptrust/services"
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
)

type ApptrustServicesManager struct {
	client *jfroghttpclient.JfrogHttpClient
	config config.Config
}

func New(config config.Config) (*ApptrustServicesManager, error) {
	details := config.GetServiceDetails()
	var err error
	manager := &ApptrustServicesManager{config: config}
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

func (asm *ApptrustServicesManager) Client() *jfroghttpclient.JfrogHttpClient {
	return asm.client
}

func (asm *ApptrustServicesManager) GetApplicationDetails(applicationKey string) (*services.Application, error) {
	appService := services.NewApplicationService(asm.config.GetServiceDetails(), asm.client)
	return appService.GetApplicationDetails(applicationKey)
}

func (asm *ApptrustServicesManager) GetApplicationVersionPromotions(applicationKey, applicationVersion string, queryParams map[string]string) (*services.ApplicationPromotionsResponse, error) {
	appService := services.NewApplicationService(asm.config.GetServiceDetails(), asm.client)
	return appService.GetApplicationVersionPromotions(applicationKey, applicationVersion, queryParams)
}
