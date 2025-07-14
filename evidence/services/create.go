package services

import (
	"path"

	"github.com/jfrog/jfrog-client-go/utils/log"
)

const (
	evidenceCreateAPI = "api/v1/subject"
)

type createEvidenceOperation struct {
	evidenceBody EvidenceCreationBody
}

func (ce *createEvidenceOperation) getOperationRestApi() string {
	return path.Join(evidenceCreateAPI, ce.evidenceBody.SubjectUri)
}

func (ce *createEvidenceOperation) getRequestBody() []byte {
	return ce.evidenceBody.DSSEFileRaw
}

func (ce *createEvidenceOperation) getProviderId() string {
	return ce.evidenceBody.ProviderId
}

func (es *EvidenceService) UploadEvidence(evidenceDetails EvidenceDetails) ([]byte, error) {
	operation := createEvidenceOperation{
		evidenceBody: EvidenceCreationBody{
			EvidenceDetails: evidenceDetails,
		},
	}
	if !es.IsEvidenceSupportsProviderId() {
		// If the evidence version does not support providerId, we will not set it in the request
		log.Warn("Evidence version does not support providerId. The providerId will not be set in the request.")
		operation.evidenceBody.ProviderId = ""
	}
	body, err := es.doOperation(&operation)
	return body, err
}

type EvidenceCreationBody struct {
	EvidenceDetails
}
