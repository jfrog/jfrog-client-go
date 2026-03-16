package services

import (
	"net/http"
	"strings"

	"github.com/jfrog/jfrog-client-go/utils/errorutils"
)

const evidenceVersionAPI = "api/v1/system/version"

func (es *EvidenceService) GetVersion() (string, error) {
	httpClientDetails := es.GetEvidenceDetails().CreateHttpClientDetails()
	resp, body, _, err := es.client.SendGet(es.GetEvidenceDetails().GetUrl()+evidenceVersionAPI, true, &httpClientDetails)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return "", err
	}
	return strings.TrimSpace(string(body)), nil
}
