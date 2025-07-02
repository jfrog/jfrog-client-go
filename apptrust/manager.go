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

func (sm *ApptrustServicesManager) Client() *jfroghttpclient.JfrogHttpClient {
	return sm.client
}

func (sm *ApptrustServicesManager) Config() config.Config {
	return sm.config
}

// Returns AppTrust server version
func (sm *ApptrustServicesManager) GetVersion() (string, error) {
	// TODO eran complete this service
	return "1.0.0", nil
	/*
		versionService := services.NewVersionService(sm.client)
		versionService.XrayDetails = sm.config.GetServiceDetails()
		return versionService.GetVersion()

	*/
}

// Returns evaluation response from Unified Policy server indicating on evaluation result OR missing scans that are required for performing the evaluation.
func (sm *ApptrustServicesManager) Evaluate(params services.EvaluateRequest) (response services.EvaluateResponse, err error) {
	evaluationService := services.NewEvaluationService(sm.client)
	evaluationService.ApptrustDetails = sm.config.GetServiceDetails()
	return evaluationService.Evaluate(params)
}
