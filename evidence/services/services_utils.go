package services

import (
	"net/http"
	"net/url"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

func IsEvidenceSupportsProviderId(evidenceDetails auth.ServiceDetails, client *jfroghttpclient.JfrogHttpClient) bool {
	// providerId is supported from evidence version XXX
	// get evidence version API was added afterwards so we will check that the API returns 200 OK
	// and not 404 Not Found
	es := NewEvidenceService(evidenceDetails, client)

	requestFullUrl, err := url.Parse(es.GetEvidenceDetails().GetUrl() + "api/v1/system/version")
	if err != nil {
		return false
	}

	httpClientDetails := es.GetEvidenceDetails().CreateHttpClientDetails()
	httpClientDetails.SetContentTypeApplicationJson()

	log.Debug("Checking evidence version: ")
	resp, _, _, err := client.SendGet(requestFullUrl.String(), true, &httpClientDetails)
	if err != nil {
		return false
	}

	return resp.StatusCode == http.StatusOK
}
