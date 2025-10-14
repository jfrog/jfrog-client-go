package services

import (
	"errors"
	"github.com/jfrog/jfrog-client-go/utils"
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

func (rbs *ReleaseBundlesService) ReleaseBundlesSearchGroups(optionalQueryParams GetSearchOptionalQueryParams) (ReleaseBundlesGroupResponse, error) {
	restApi := GetReleaseBundleSearchGroupApi()
	requestFullUrl, err := utils.BuildUrl(rbs.GetLifecycleDetails().GetUrl(), restApi, buildGetSearchQueryParams(optionalQueryParams))
	if err != nil {
		return ReleaseBundlesGroupResponse{}, err
	}
	httpClientsDetails := rbs.GetLifecycleDetails().CreateHttpClientDetails()
	var response ReleaseBundlesGroupResponse
	err = rbs.doHttpRequestWithRetry(requestFullUrl, &httpClientsDetails, &response)
	return response, err
}

func (rbs *ReleaseBundlesService) ReleaseBundlesSearchVersions(releaseBundleName string, optionalQueryParams GetSearchOptionalQueryParams) (ReleaseBundleVersionsResponse, error) {
	restApi := GetReleaseBundleSearchVersionsApi(releaseBundleName)
	requestFullUrl, err := utils.BuildUrl(rbs.GetLifecycleDetails().GetUrl(), restApi, buildGetSearchQueryParams(optionalQueryParams))
	if err != nil {
		return ReleaseBundleVersionsResponse{}, err
	}
	httpClientsDetails := rbs.GetLifecycleDetails().CreateHttpClientDetails()
	var response ReleaseBundleVersionsResponse
	err = rbs.doHttpRequestWithRetry(requestFullUrl, &httpClientsDetails, &response)
	return response, err
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
