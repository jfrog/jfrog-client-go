package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
	"path"
	"strconv"
	"time"
)

const (
	groupApi       = "groups"
	maxAttempts    = 3
	initialBackoff = 100 * time.Millisecond
	maxBackoff     = 2 * time.Second
)

var (
	ErrAuth       = errors.New("authentication failed")
	ErrPermission = errors.New("permission denied")
)

func buildGetSearchQueryParams(optionalQueryParams GetSearchOptionalQueryParams) map[string]string {
	params := make(map[string]string)
	if optionalQueryParams.Includes != "" {
		params["include"] = optionalQueryParams.Includes
	}
	if optionalQueryParams.Offset > 0 {
		params["offset"] = strconv.Itoa(optionalQueryParams.Offset)
	}
	if optionalQueryParams.Limit > 0 {
		params["limit"] = strconv.Itoa(optionalQueryParams.Limit)
	}
	if optionalQueryParams.FilterBy != "" {
		params["filter_by"] = optionalQueryParams.FilterBy
	}
	if optionalQueryParams.OrderBy != "" {
		params["order_by"] = optionalQueryParams.OrderBy
	}
	if optionalQueryParams.OrderAsc {
		params["order_asc"] = strconv.FormatBool(optionalQueryParams.OrderAsc)
	}
	return params
}

func (rbs *ReleaseBundlesService) ReleaseBundlesSearchGroups(optionalQueryParams GetSearchOptionalQueryParams) (response ReleaseBundlesGroupResponse, err error) {
	restApi := GetReleaseBundleSearchGroupApi()
	requestFullUrl, err := utils.BuildUrl(rbs.GetLifecycleDetails().GetUrl(), restApi, buildGetSearchQueryParams(optionalQueryParams))
	if err != nil {
		return
	}
	httpClientsDetails := rbs.GetLifecycleDetails().CreateHttpClientDetails()

	for i := 0; i < maxAttempts; i++ {
		resp, body, _, sendErr := rbs.client.SendGet(requestFullUrl, true, &httpClientsDetails)
		if sendErr != nil {
			log.Debug(fmt.Sprintf("Attempt %d/%d: Network error during SendGet to %s: %v", i+1, maxAttempts, requestFullUrl, sendErr))
			if i < maxAttempts-1 {
				time.Sleep(utils.CalculateBackoff(i, initialBackoff, maxBackoff))
				continue
			}
			return response, sendErr
		}
		log.Debug("Artifactory response status:", resp.Status)
		switch resp.StatusCode {
		case http.StatusOK:
			err = errorutils.CheckError(json.Unmarshal(body, &response))
			return response, err
		case http.StatusUnauthorized:
			return response, ErrAuth
		case http.StatusForbidden:
			return response, ErrPermission
		case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
			if i < maxAttempts-1 {
				log.Debug(fmt.Sprintf("Attempt %d/%d: Server error (%d) from %s. Retrying...", i+1, maxAttempts, resp.StatusCode, requestFullUrl))
				time.Sleep(utils.CalculateBackoff(i, initialBackoff, maxBackoff))
				continue
			}
		}
		err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK)
		return response, err
	}
	return response, fmt.Errorf("failed to get response from %s after %d attempts", requestFullUrl, maxAttempts)
}

func (rbs *ReleaseBundlesService) ReleaseBundlesSearchVersions(releaseBundleName string, optionalQueryParams GetSearchOptionalQueryParams) (response ReleaseBundleVersionsResponse, err error) {
	restApi := GetReleaseBundleSearchVersionsApi(releaseBundleName) // This is the main difference in API path
	requestFullUrl, err := utils.BuildUrl(rbs.GetLifecycleDetails().GetUrl(), restApi, buildGetSearchQueryParams(optionalQueryParams))
	if err != nil {
		return
	}
	httpClientsDetails := rbs.GetLifecycleDetails().CreateHttpClientDetails()
	for i := 0; i < maxAttempts; i++ {
		resp, body, _, sendErr := rbs.client.SendGet(requestFullUrl, true, &httpClientsDetails)
		if sendErr != nil {
			log.Debug(fmt.Sprintf("Attempt %d/%d: Network error during SendGet to %s: %v", i+1, maxAttempts, requestFullUrl, sendErr))
			if i < maxAttempts-1 {
				time.Sleep(utils.CalculateBackoff(i, initialBackoff, maxBackoff))
				continue
			}
			return response, sendErr
		}
		log.Debug("Artifactory response status:", resp.Status)
		switch resp.StatusCode {
		case http.StatusOK:
			err = errorutils.CheckError(json.Unmarshal(body, &response))
			return response, err
		case http.StatusUnauthorized:
			return response, ErrAuth
		case http.StatusForbidden:
			return response, ErrPermission
		case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
			if i < maxAttempts-1 {
				log.Debug(fmt.Sprintf("Attempt %d/%d: Server error (%d) from %s. Retrying...", i+1, maxAttempts, resp.StatusCode, requestFullUrl))
				time.Sleep(utils.CalculateBackoff(i, initialBackoff, maxBackoff))
				continue
			}
		}
		err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK)
		return response, err
	}
	return response, fmt.Errorf("failed to get response from %s after %d attempts", requestFullUrl, maxAttempts)
}

func GetReleaseBundleSearchGroupApi() string {
	return path.Join(releaseBundleNewApi, groupApi)
}

func GetReleaseBundleSearchVersionsApi(releaseBundleName string) string {
	return path.Join(releaseBundleNewApi, records, releaseBundleName)
}

type ReleaseBundleSearchGroup struct {
	RepositoryKey              string    `json:"repository_key"`
	ProjectKey                 string    `json:"project_key"`
	ProjectName                string    `json:"project_name"`
	ServiceID                  string    `json:"service_id"`
	Created                    time.Time `json:"created"`
	ReleaseBundleName          string    `json:"release_bundle_name"`
	ReleaseBundleVersionLatest string    `json:"release_bundle_version_latest"`
	ReleaseBundleVersionsCount int       `json:"release_bundle_versions_count"`
}

// ReleaseBundlesGroupResponse represents the entire JSON response structure
type ReleaseBundlesGroupResponse struct {
	ReleaseBundleSearchGroup []ReleaseBundleSearchGroup `json:"release_bundles"`
	Total                    int                        `json:"total"`
	Limit                    int                        `json:"limit"`
	Offset                   int                        `json:"offset"`
}

type ReleaseBundleVersion struct {
	Status               string    `json:"status"`
	RepositoryKey        string    `json:"repository_key"`
	ReleaseBundleName    string    `json:"release_bundle_name"`
	ReleaseBundleVersion string    `json:"release_bundle_version"`
	ServiceID            string    `json:"service_id"`
	CreatedBy            string    `json:"created_by"`
	Created              time.Time `json:"created"`
	ReleaseStatus        string    `json:"release_status"`
}

// ReleaseBundleVersionsResponse represents the entire JSON response structure for versions
type ReleaseBundleVersionsResponse struct {
	ReleaseBundles []ReleaseBundleVersion `json:"release_bundles"`
	Total          int                    `json:"total"`
	Limit          int                    `json:"limit"`
	Offset         int                    `json:"offset"`
}

type GetSearchOptionalQueryParams struct {
	Offset   int
	Limit    int
	FilterBy string
	OrderBy  string
	OrderAsc bool
	Includes string
}
