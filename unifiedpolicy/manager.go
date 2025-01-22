package unifiedpolicy

import (
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/unifiedpolicy/services"
)

type UnifiedPolicyServicesManager struct {
	client *jfroghttpclient.JfrogHttpClient
	config config.Config
}

// New creates a service manager to interact with Application
func New(config config.Config) (*UnifiedPolicyServicesManager, error) {
	details := config.GetServiceDetails()
	var err error
	manager := &UnifiedPolicyServicesManager{config: config}
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

func (up *UnifiedPolicyServicesManager) Evaluate(evaluateRequest *services.EvaluateRequest) (*services.EvaluateResponse, error) {
	evaluateService := services.NewEvaluateService(up.client, up.config.GetServiceDetails())
	return evaluateService.Evaluate(evaluateRequest)
}

func (up *UnifiedPolicyServicesManager) GetVersion() (string, error) {
	versionService := services.NewVersionService(up.client)
	versionService.UnifiedPolicyDetails = up.config.GetServiceDetails()
	return versionService.GetVersion()
}
