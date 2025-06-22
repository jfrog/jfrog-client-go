package services

import (
	"net/http"
	"net/url"

	"github.com/jfrog/jfrog-client-go/utils/log"
)

func (es *EvidenceService) IsEvidenceSupportsProviderId() bool {
	// providerId is supported from evidence version XXX
	// get evidence version API was added afterwards so we will check that the API returns 200 OK
	// and not 404 Not Found

	requestFullUrl, err := url.Parse(es.GetEvidenceDetails().GetUrl() + "api/v1/system/version")
	if err != nil {
		return false
	}

	httpClientDetails := es.GetEvidenceDetails().CreateHttpClientDetails()
	httpClientDetails.SetContentTypeApplicationJson()

	log.Debug("Checking evidence version: ")
	resp, _, _, err := es.client.SendGet(requestFullUrl.String(), true, &httpClientDetails)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}
