package services

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

const versionApi = "api/v1/system/version"

func (es *EvidenceService) IsEvidenceSupportsProviderId() bool {
	// providerId is supported from evidence version XXX
	// get evidence version API was added afterwards so we will check that the API returns 200 OK
	// and not 404 Not Found

	requestFullUrl, err := url.Parse(es.GetEvidenceDetails().GetUrl() + versionApi)
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

func (es *EvidenceService) IsEvidenceVersionSupportsFeature(version string) (bool, error) {
	if strings.TrimSpace(version) == "" {
		return false, errorutils.CheckErrorf("supported version cannot be empty")
	}

	requestFullUrl, err := url.Parse(es.GetEvidenceDetails().GetUrl() + versionApi)
	if err != nil {
		return false, err
	}

	httpClientDetails := es.GetEvidenceDetails().CreateHttpClientDetails()
	httpClientDetails.SetContentTypeApplicationJson()

	log.Debug("Checking evidence version: ")
	resp, body, _, err := es.client.SendGet(requestFullUrl.String(), true, &httpClientDetails)
	if err != nil {
		return false, err
	}

	if resp.StatusCode != http.StatusOK {
		return false, errorutils.CheckErrorf("failed to get evidence version, status code: %d", resp.StatusCode)
	}

	// Try to parse as JSON first (as per existing tests)
	var versionResponse map[string]interface{}
	var respVersion string

	if err := json.Unmarshal(body, &versionResponse); err == nil {
		if respVersionVal, ok := versionResponse["version"].(string); ok {
			respVersion = respVersionVal
		} else {
			return false, errorutils.CheckErrorf("version field not found or not a string in response")
		}
	} else {
		// Fallback to raw string if not JSON
		respVersion = strings.TrimSpace(string(body))
	}

	return CompareVersions(respVersion, version)
}

const versionParts = 3 // version format is x.y.z 3 parts major, minor, patch

func CompareVersions(version, minVersion string) (bool, error) {
	version = strings.TrimSpace(version)
	minVersion = strings.TrimSpace(minVersion)

	if version == "" || minVersion == "" {
		return false, errorutils.CheckErrorf("version strings cannot be empty")
	}

	if strings.Compare(version, minVersion) == 0 {
		return true, nil
	}

	curParts := strings.Split(version, ".")
	minParts := strings.Split(minVersion, ".")

	if len(curParts) != versionParts || len(minParts) != versionParts {
		return false, errorutils.CheckErrorf("invalid version format: %s or %s. Expected format: x.y.z", version, minVersion)
	}

	for i := 0; i < versionParts; i++ {
		curNum, err := strconv.Atoi(curParts[i])
		if err != nil {
			return false, errorutils.CheckErrorf("Invalid version number: %s", curParts[i])
		}

		minNum, err := strconv.Atoi(minParts[i])
		if err != nil {
			return false, errorutils.CheckErrorf("Invalid minimum version number: %s", minParts[i])
		}

		if curNum < minNum {
			return false, nil
		} else if curNum > minNum {
			return true, nil
		}
	}

	return true, nil
}
