package services

import (
	"path"

	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

const (
	evidenceCreateAPI                   = "api/v1/subject"
	minEvidenceVersionForAttachmentsAPI = "7.646.1"
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
	if err := es.resolveUploadCompatibility(&evidenceDetails); err != nil {
		return nil, err
	}
	operation := createEvidenceOperation{
		evidenceBody: EvidenceCreationBody{
			EvidenceDetails: evidenceDetails,
		},
	}
	body, opErr := es.doOperation(&operation)
	return body, opErr
}

type EvidenceCreationBody struct {
	EvidenceDetails
}

func (es *EvidenceService) resolveUploadCompatibility(evidenceDetails *EvidenceDetails) error {
	if evidenceDetails == nil {
		return nil
	}
	requiresVersionLookup := len(evidenceDetails.Attachments) > 0 || evidenceDetails.ProviderId != ""
	if !requiresVersionLookup {
		return nil
	}

	currentVersion, err := es.GetVersion()
	if err != nil {
		if len(evidenceDetails.Attachments) > 0 {
			return err
		}
		// ProviderId support is inferred from version endpoint availability.
		log.Warn("Evidence version endpoint is unavailable. The providerId will not be set in the request.")
		evidenceDetails.ProviderId = ""
		return nil
	}
	if len(evidenceDetails.Attachments) > 0 {
		if err = clientutils.ValidateMinimumVersion("JFrog Evidence", currentVersion, minEvidenceVersionForAttachmentsAPI); err != nil {
			return err
		}
	}
	return nil
}
