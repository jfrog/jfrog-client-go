package services

import (
	"net/http"
	"net/url"
	"path"

	"github.com/jfrog/jfrog-client-go/utils/log"
)

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
	if !es.isEvidenceSupportsProviderId() {
		// If the evidence version does not support providerId, we will not set it in the request
		log.Warn("Evidence version does not support providerId. The providerId will not be set in the request.")
		operation.evidenceBody.ProviderId = ""
	}
	body, err := es.doOperation(&operation)
	return body, err
}

func (es *EvidenceService) isEvidenceSupportsProviderId() bool {
	// providerId is supported from evidence version XXX
	// get evidence version API was added afterwards so we will check that the API returns 200 OK
	// and not 404 Not Found
	requestFullUrl, err := url.Parse(es.GetEvidenceDetails().GetUrl() + "api/v1/system/version")
	if err != nil {
		return false
	}

	httpClientDetails := es.GetEvidenceDetails().CreateHttpClientDetails()
	httpClientDetails.SetContentTypeApplicationJson()

	log.Debug("Check evidence version. Sending request to: ", requestFullUrl.String())
	resp, _, _, err := es.client.SendGet(requestFullUrl.String(), true, &httpClientDetails)
	if err != nil {
		return false
	}

	return resp.StatusCode == http.StatusOK
}

type EvidenceCreationBody struct {
	EvidenceDetails
}
