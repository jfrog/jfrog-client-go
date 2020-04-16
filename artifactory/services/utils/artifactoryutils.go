package utils

import (
	"encoding/json"
	"errors"
	rthttpclient "github.com/jfrog/jfrog-client-go/artifactory/httpclient"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/httpclient"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const (
	ARTIFACTORY_SYMLINK = "symlink.dest"
	SYMLINK_SHA1        = "symlink.destsha1"
	Latest              = "LATEST"
	LastRelease         = "LAST_RELEASE"
)

func UploadFile(localPath, url, logMsgPrefix string, artifactoryDetails *auth.ServiceDetails, details *fileutils.FileDetails,
	httpClientsDetails httputils.HttpClientDetails, client *rthttpclient.ArtifactoryHttpClient, retries int, progress io.Progress) (*http.Response, []byte, error) {
	var err error
	if details == nil {
		details, err = fileutils.GetFileDetails(localPath)
	}
	if err != nil {
		return nil, nil, err
	}

	requestClientDetails := httpClientsDetails.Clone()
	AddChecksumHeaders(requestClientDetails.Headers, details)
	AddAuthHeaders(requestClientDetails.Headers, *artifactoryDetails)

	return client.UploadFile(localPath, url, logMsgPrefix, requestClientDetails, retries, progress)
}

func AddChecksumHeaders(headers map[string]string, fileDetails *fileutils.FileDetails) {
	AddHeader("X-Checksum-Sha1", fileDetails.Checksum.Sha1, &headers)
	AddHeader("X-Checksum-Md5", fileDetails.Checksum.Md5, &headers)
	if len(fileDetails.Checksum.Sha256) > 0 {
		AddHeader("X-Checksum", fileDetails.Checksum.Sha256, &headers)
	}
}

func AddAuthHeaders(headers map[string]string, artifactoryDetails auth.ServiceDetails) {
	if headers == nil {
		headers = make(map[string]string)
	}
	if artifactoryDetails.GetSshAuthHeaders() != nil {
		utils.MergeMaps(artifactoryDetails.GetSshAuthHeaders(), headers)
	}
}

func SetContentType(contentType string, headers *map[string]string) {
	AddHeader("Content-Type", contentType, headers)
}

func DisableAccelBuffering(headers *map[string]string) {
	AddHeader("X-Accel-Buffering", "no", headers)
}

func AddHeader(headerName, headerValue string, headers *map[string]string) {
	if *headers == nil {
		*headers = make(map[string]string)
	}
	(*headers)[headerName] = headerValue
}

func BuildArtifactoryUrl(baseUrl, path string, params map[string]string) (string, error) {
	u := url.URL{Path: path}
	escapedUrl, err := url.Parse(baseUrl + u.String())
	err = errorutils.CheckError(err)
	if err != nil {
		return "", err
	}
	q := escapedUrl.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	escapedUrl.RawQuery = q.Encode()
	return escapedUrl.String(), nil
}

func IsWildcardPattern(pattern string) bool {
	return strings.Contains(pattern, "*") || strings.HasSuffix(pattern, "/") || !strings.Contains(pattern, "/")
}

// @paths - sorted array
// @index - index of the current path which we want to check if it a prefix of any of the other previous paths
// @separator - file separator
// returns true paths[index] is a prefix of any of the paths[i] where i<index , otherwise returns false
func IsSubPath(paths []string, index int, separator string) bool {
	currentPath := paths[index]
	if !strings.HasSuffix(currentPath, separator) {
		currentPath += separator
	}
	for i := index - 1; i >= 0; i-- {
		if strings.HasPrefix(paths[i], currentPath) {
			return true
		}
	}
	return false
}

// This method parses buildIdentifier. buildIdentifier should be from the format "buildName/buildNumber".
// If no buildNumber provided LATEST wil be downloaded.
// If buildName or buildNumber contains "/" (slash) it should be escaped by "\" (backslash).
// Result examples of parsing: "aaa/123" > "aaa"-"123", "aaa" > "aaa"-"LATEST", "aaa\\/aaa" > "aaa/aaa"-"LATEST",  "aaa/12\\/3" > "aaa"-"12/3".
func getBuildNameAndNumberFromBuildIdentifier(buildIdentifier string, flags CommonConf) (string, string, error) {
	buildName, buildNumber, err := parseNameAndVersion(buildIdentifier, true)
	if err != nil {
		return "", "", err
	}
	return GetBuildNameAndNumberFromArtifactory(buildName, buildNumber, flags)
}

func GetBuildNameAndNumberFromArtifactory(buildName, buildNumber string, flags CommonConf) (string, string, error) {
	if buildNumber == Latest || buildNumber == LastRelease {
		return getLatestBuildNumberFromArtifactory(buildName, buildNumber, flags)
	}
	return buildName, buildNumber, nil
}

func getBuildNameAndNumberFromProps(properties []Property) (buildName string, buildNumber string) {
	for _, property := range properties {
		if property.Key == "build.name" {
			buildName = property.Value
		} else if property.Key == "build.number" {
			buildNumber = property.Value
		}
		if len(buildName) > 0 && len(buildNumber) > 0 {
			return buildName, buildNumber
		}
	}
	return
}

// For builds (useLatestPolicy = true) - Parse build name and number. The build number can be LATEST if absent.
// For release bundles - Parse bundle name and version.
func parseNameAndVersion(identifier string, useLatestPolicy bool) (string, string, error) {
	const Delimiter = "/"
	const EscapeChar = "\\"

	if identifier == "" {
		return "", "", nil
	}
	if !strings.Contains(identifier, Delimiter) {
		if useLatestPolicy {
			log.Debug("No '" + Delimiter + "' is found in the build, build number is set to " + Latest)
			return identifier, Latest, nil
		} else {
			return "", "", errorutils.CheckError(errors.New("No '" + Delimiter + "' is found in the bundle"))
		}
	}
	name, version := "", ""
	versionsArray := []string{}
	identifiers := strings.Split(identifier, Delimiter)
	// The delimiter must not be prefixed with escapeChar (if it is, it should be part of the version)
	// the code below gets substring from before the last delimiter.
	// If the new string ends with escape char it means the last delimiter was part of the version and we need
	// to go back to the previous delimiter.
	// If no proper delimiter was found the full string will be the name.
	for i := len(identifiers) - 1; i >= 1; i-- {
		versionsArray = append([]string{identifiers[i]}, versionsArray...)
		if !strings.HasSuffix(identifiers[i-1], EscapeChar) {
			name = strings.Join(identifiers[:i], Delimiter)
			version = strings.Join(versionsArray, Delimiter)
			break
		}
	}
	if name == "" {
		if useLatestPolicy {
			log.Debug("No delimiter char (" + Delimiter + ") without escaping char was found in the build, build number is set to " + Latest)
			name = identifier
			version = Latest
		} else {
			return "", "", errorutils.CheckError(errors.New("No delimiter char (" + Delimiter + ") without escaping char was found in the bundle"))
		}
	}
	// Remove escape chars
	name = strings.Replace(name, "\\/", "/", -1)
	version = strings.Replace(version, "\\/", "/", -1)
	return name, version, nil
}

type build struct {
	BuildName   string `json:"buildName"`
	BuildNumber string `json:"buildNumber"`
}

func getLatestBuildNumberFromArtifactory(buildName, buildNumber string, flags CommonConf) (string, string, error) {
	restUrl := flags.GetArtifactoryDetails().GetUrl() + "api/build/patternArtifacts"
	body, err := createBodyForLatestBuildRequest(buildName, buildNumber)
	if err != nil {
		return "", "", err
	}
	log.Debug("Getting build name and number from Artifactory: " + buildName + ", " + buildNumber)
	httpClientsDetails := flags.GetArtifactoryDetails().CreateHttpClientDetails()
	SetContentType("application/json", &httpClientsDetails.Headers)
	log.Debug("Sending post request to: " + restUrl + ", with the following body: " + string(body))
	client, err := httpclient.ClientBuilder().Build()
	if err != nil {
		return "", "", err
	}
	resp, body, err := client.SendPost(restUrl, body, httpClientsDetails)
	if err != nil {
		return "", "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", "", errorutils.CheckError(errors.New("Artifactory response: " + resp.Status + "\n" + utils.IndentJson(body)))
	}
	log.Debug("Artifactory response: ", resp.Status)
	var responseBuild []build
	err = json.Unmarshal(body, &responseBuild)
	if errorutils.CheckError(err) != nil {
		return "", "", err
	}
	if responseBuild[0].BuildNumber != "" {
		log.Debug("Found build number: " + responseBuild[0].BuildNumber)
	} else {
		log.Debug("The build could not be found in Artifactory")
	}

	return buildName, responseBuild[0].BuildNumber, nil
}

func createBodyForLatestBuildRequest(buildName, buildNumber string) (body []byte, err error) {
	buildJsonArray := []build{{buildName, buildNumber}}
	body, err = json.Marshal(buildJsonArray)
	err = errorutils.CheckError(err)
	return
}

func filterAqlSearchResultsByBuild(specFile *ArtifactoryCommonParams, itemsToFilter []ResultItem, flags CommonConf, itemsAlreadyContainProperties bool) ([]ResultItem, error) {
	var addPropsErr error
	var aqlSearchErr error
	var buildArtifactsSha1 map[string]bool
	var wg sync.WaitGroup
	buildName, buildNumber, err := getBuildNameAndNumberFromBuildIdentifier(specFile.Build, flags)
	if err != nil {
		return nil, err
	}

	wg.Add(2)
	// Get Sha1 for artifacts by build name and number
	go func() {
		buildArtifactsSha1, aqlSearchErr = fetchBuildArtifactsSha1(buildName, buildNumber, flags)
		wg.Done()
	}()

	if !itemsAlreadyContainProperties {
		// Add properties to the previously found artifacts (in case properties weren't already fetched from Artifactory)
		go func() {
			addPropsErr = searchAndAddPropsToAqlResult(itemsToFilter, specFile.Aql.ItemsFind, "build.name", buildName, flags)
			wg.Done()
		}()
	} else {
		wg.Done()
	}

	wg.Wait()
	if aqlSearchErr != nil {
		return nil, aqlSearchErr
	}
	if addPropsErr != nil {
		return nil, addPropsErr
	}

	return filterBuildAqlSearchResults(&itemsToFilter, &buildArtifactsSha1, buildName, buildNumber), err
}

// Run AQL to retrieve all artifacts associated with a specific build.
// Return a map of the artifacts SHA1.
func fetchBuildArtifactsSha1(buildName, buildNumber string, flags CommonConf) (map[string]bool, error) {
	buildQuery := createAqlQueryForBuild(buildName, buildNumber, buildIncludeQueryPart([]string{"name", "repo", "path", "actual_sha1"}))

	parsedBuildAqlResponse, err := aqlSearch(buildQuery, flags)
	if err != nil {
		return nil, err
	}
	buildArtifactsSha, err := extractSha1FromAqlResponse(parsedBuildAqlResponse)
	if err != nil {
		return nil, err
	}

	return buildArtifactsSha, nil
}

/*
 * Find artifact properties by the AQL, add them to the result items.
 *
 * resultItems - Artifacts to add properties to.
 * aqlBody - AQL to execute together with property filter.
 * filterByPropName - Property name to filter.
 * filterByPropValue - Property value to filter.
 * flags - Command flags for AQL execution.
 */
func searchAndAddPropsToAqlResult(resultItems []ResultItem, aqlBody, filterByPropName, filterByPropValue string, flags CommonConf) error {
	propsAqlResponseJson, err := ExecAql(createPropsQuery(aqlBody, filterByPropName, filterByPropValue), flags)
	if err != nil {
		return err
	}
	propsAqlResponse, err := parseAqlSearchResponse(propsAqlResponseJson)
	if err != nil {
		return err
	}
	addPropsToAqlResult(resultItems, propsAqlResponse)
	return nil
}

func addPropsToAqlResult(items []ResultItem, props []ResultItem) {
	propsMap := createPropsMap(props)
	for i := range items {
		props, propsExists := propsMap[getResultItemKey(items[i])]
		if propsExists {
			items[i].Properties = props
		}
	}
}

func createPropsMap(items []ResultItem) (propsMap map[string][]Property) {
	propsMap = make(map[string][]Property)
	for _, item := range items {
		propsMap[getResultItemKey(item)] = item.Properties
	}
	return
}

func getResultItemKey(item ResultItem) string {
	return item.Repo + item.Path + item.Name + item.Actual_Sha1
}

func extractSha1FromAqlResponse(elements []ResultItem) (map[string]bool, error) {
	elementsMap := make(map[string]bool)
	for _, element := range elements {
		elementsMap[element.Actual_Sha1] = true
	}
	return elementsMap, nil
}

/*
 * Filter search results by the following priorities:
 * 1st priority: Match {Sha1, build name, build number}
 * 2nd priority: Match {Sha1, build name}
 * 3rd priority: Match {Sha1}
 */
func filterBuildAqlSearchResults(itemsToFilter *[]ResultItem, buildArtifactsSha *map[string]bool, buildName, buildNumber string) []ResultItem {
	filteredResults := []ResultItem{}
	firstPriority := map[string][]ResultItem{}
	secondPriority := map[string][]ResultItem{}
	thirdPriority := map[string][]ResultItem{}

	// Step 1 - Populate 3 priorities mappings.
	for _, item := range *itemsToFilter {
		if _, ok := (*buildArtifactsSha)[item.Actual_Sha1]; !ok {
			continue
		}
		resultBuildName, resultBuildNumber := getBuildNameAndNumberFromProps(item.Properties)
		isBuildNameMatched := resultBuildName == buildName
		if isBuildNameMatched && resultBuildNumber == buildNumber {
			firstPriority[item.Actual_Sha1] = append(firstPriority[item.Actual_Sha1], item)
			continue
		}
		if isBuildNameMatched {
			secondPriority[item.Actual_Sha1] = append(secondPriority[item.Actual_Sha1], item)
			continue
		}
		thirdPriority[item.Actual_Sha1] = append(thirdPriority[item.Actual_Sha1], item)
	}

	// Step 2 - Append mappings to the final results, respectively.
	for shaToMatch := range *buildArtifactsSha {
		if _, ok := firstPriority[shaToMatch]; ok {
			filteredResults = append(filteredResults, firstPriority[shaToMatch]...)
		} else if _, ok := secondPriority[shaToMatch]; ok {
			filteredResults = append(filteredResults, secondPriority[shaToMatch]...)
		} else if _, ok := thirdPriority[shaToMatch]; ok {
			filteredResults = append(filteredResults, thirdPriority[shaToMatch]...)
		}
	}

	return filteredResults
}

type CommonConf interface {
	GetArtifactoryDetails() auth.ServiceDetails
	SetArtifactoryDetails(rt auth.ServiceDetails)
	GetJfrogHttpClient() (*rthttpclient.ArtifactoryHttpClient, error)
	IsDryRun() bool
}

type CommonConfImpl struct {
	artDetails auth.ServiceDetails
	DryRun     bool
}

func (flags *CommonConfImpl) GetArtifactoryDetails() auth.ServiceDetails {
	return flags.artDetails
}

func (flags *CommonConfImpl) SetArtifactoryDetails(rt auth.ServiceDetails) {
	flags.artDetails = rt
}

func (flags *CommonConfImpl) IsDryRun() bool {
	return flags.DryRun
}

func (flags *CommonConfImpl) GetJfrogHttpClient() (*rthttpclient.ArtifactoryHttpClient, error) {
	return rthttpclient.ArtifactoryClientBuilder().SetServiceDetails(&flags.artDetails).Build()
}
