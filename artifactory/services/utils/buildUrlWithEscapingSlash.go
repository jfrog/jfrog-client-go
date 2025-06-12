package utils

import (
	"net/url"
	"path"
	"strings"

	"github.com/jfrog/jfrog-client-go/utils/errorutils"
)

func BuildUrlWithEscapingSlash(baseUrl, restApi, buildName, buildNumber string, queryParams map[string]string) (string, error) {
	u := url.URL{Path: restApi}
	parsedUrl, err := url.Parse(baseUrl + u.String())
	if err = errorutils.CheckError(err); err != nil {
		return "", err
	}
	q := parsedUrl.Query()
	for k, v := range queryParams {
		q.Set(k, v)
	}
	parsedUrl.RawQuery = q.Encode()

	// Semicolons are reserved as separators in some Artifactory APIs, so they'd better be encoded when used for other purposes
	encodedUrl := strings.ReplaceAll(parsedUrl.String(), ";", url.QueryEscape(";"))

	escapedBuildPath := path.Join(url.QueryEscape(buildName), url.QueryEscape(buildNumber))
	urlParts := strings.Split(encodedUrl, "?")
	resultUrl := urlParts[0] + "/" + escapedBuildPath
	if len(urlParts) > 1 && urlParts[1] != "" {
		resultUrl += "?" + urlParts[1]
	}

	return resultUrl, nil
}
