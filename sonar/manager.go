package sonar

import (
	"fmt"

	"github.com/jfrog/jfrog-client-go/config"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/sonar/services"
)

type Manager interface {
	GetQualityGateAnalysis(analysisID string) (*services.QualityGatesAnalysis, error)
	GetTaskDetails(ceTaskID string) (*services.TaskDetails, error)
	GetSonarIntotoStatementRaw(ceTaskID string) ([]byte, error)
}

type sonarManager struct {
	client *jfroghttpclient.JfrogHttpClient
	config config.Config
}

func NewManager(config config.Config) (Manager, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	details := config.GetServiceDetails()
	var err error
	manager := &sonarManager{config: config}
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

func (sm *sonarManager) Client() *jfroghttpclient.JfrogHttpClient {
	return sm.client
}

func (sm *sonarManager) GetSonarIntotoStatementRaw(ceTaskID string) ([]byte, error) {
	sonarService := services.NewSonarService(sm.config.GetServiceDetails(), sm.client)
	return sonarService.GetSonarIntotoStatementRaw(ceTaskID)
}

func (sm *sonarManager) GetQualityGateAnalysis(analysisID string) (*services.QualityGatesAnalysis, error) {
	sonarService := services.NewSonarService(sm.config.GetServiceDetails(), sm.client)
	return sonarService.GetQualityGateAnalysis(analysisID)
}

func (sm *sonarManager) GetTaskDetails(ceTaskID string) (*services.TaskDetails, error) {
	sonarService := services.NewSonarService(sm.config.GetServiceDetails(), sm.client)
	return sonarService.GetTaskDetails(ceTaskID)
}
