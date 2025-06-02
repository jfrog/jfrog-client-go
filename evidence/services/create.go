package services

import "path"

const (
	evidenceCreateAPI    = "api/v1/subject"
	providerIdQueryParam = "?providerId="
)

type createEvidenceOperation struct {
	evidenceBody EvidenceCreationBody
}

func (ce *createEvidenceOperation) getOperationRestApi() string {
	apiUrl := path.Join(evidenceCreateAPI, ce.evidenceBody.SubjectUri)
	if ce.evidenceBody.ProviderId != "" {
		apiUrl = path.Join(apiUrl, providerIdQueryParam, ce.evidenceBody.ProviderId)
	}

	return apiUrl
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
