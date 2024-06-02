package services

import "path"

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

func (es *EvidenceService) UploadEvidence(evidenceDetails EvidenceDetails) ([]byte, error) {
	operation := createEvidenceOperation{
		evidenceBody: EvidenceCreationBody{
			EvidenceDetails: evidenceDetails,
		},
	}
	body, err := es.doOperation(&operation)
	return body, err
}

type EvidenceCreationBody struct {
	EvidenceDetails
}
