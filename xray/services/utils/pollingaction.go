package utils

import (
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"net/http"
)

func PollingAction(client *jfroghttpclient.JfrogHttpClient, endPoint string, httpClientDetails httputils.HttpClientDetails) (action func() (shouldStop bool, responseBody []byte, err error)) {
	pollingAction := func() (shouldStop bool, responseBody []byte, err error) {
		resp, body, _, err := client.SendGet(endPoint, true, &httpClientDetails)
		if err != nil {
			return true, nil, err
		}
		if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK, http.StatusAccepted); err != nil {
			return true, nil, err
		}
		// Got the full valid response.
		if resp.StatusCode == http.StatusOK {
			return true, body, nil
		}
		return false, nil, nil
	}
	return pollingAction
}
