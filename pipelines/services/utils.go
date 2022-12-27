package services

import (
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/url"
)

// constructPipelinesURL creates URL with all required details to make api call
// like headers, queryParams, apiPath
func constructPipelinesURL(qParams map[string]string, apiURL, apiPath string) (string, error) {
	uri, err := url.Parse(apiURL + apiPath)
	if err != nil {
		log.Error("Failed to parse pipelines fetch run status url")
		return "", errorutils.CheckError(err)
	}
	queryString := uri.Query()
	for k, v := range qParams {
		queryString.Set(k, v)
	}
	uri.RawQuery = queryString.Encode()

	return uri.String(), nil
}
