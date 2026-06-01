package jpd

import (
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
)

type JPDServicesManager struct {
	client *jfroghttpclient.JfrogHttpClient
	config config.Config
}

func NewJPDServicesManager(config config.Config) (*JPDServicesManager, error) {
	details := config.GetServiceDetails()
	var err error
	manager := &JPDServicesManager{config: config}
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

func (jpdsm *JPDServicesManager) Client() *jfroghttpclient.JfrogHttpClient {
	return jpdsm.client
}

func (jpdsm *JPDServicesManager) GetJPDsStats(serverUrl string) ([]byte, error) {
	jpdStatsService := NewJPDsStatsService(jpdsm.config.GetServiceDetails(), jpdsm.client)
	return jpdStatsService.GetAllJPDs(serverUrl)
}
