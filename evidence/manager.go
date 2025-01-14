package evidence

import (
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/jfrog/jfrog-client-go/evidence/services"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
)

type EvidenceServicesManager struct {
	client *jfroghttpclient.JfrogHttpClient
	config config.Config
}

func New(config config.Config) (*EvidenceServicesManager, error) {
	details := config.GetServiceDetails()
	var err error
	manager := &EvidenceServicesManager{config: config}
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

func (esm *EvidenceServicesManager) Client() *jfroghttpclient.JfrogHttpClient {
	return esm.client
}

func (esm *EvidenceServicesManager) UploadEvidence(evidenceDetails services.EvidenceDetails) ([]byte, error) {
	evidenceService := services.NewEvidenceService(esm.config.GetServiceDetails(), esm.client)
	return evidenceService.UploadEvidence(evidenceDetails)
}
