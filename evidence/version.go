package evidence

import "github.com/jfrog/jfrog-client-go/evidence/services"

func (esm *EvidenceServicesManager) GetVersion() (string, error) {
	evidenceService := services.NewEvidenceService(esm.config.GetServiceDetails(), esm.client)
	return evidenceService.GetVersion()
}
