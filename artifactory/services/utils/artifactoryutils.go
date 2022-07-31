package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"sync"

	buildinfo "github.com/jfrog/build-info-go/entities"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	clientio "github.com/jfrog/jfrog-client-go/utils/io"
	"github.com/jfrog/jfrog-client-go/utils/io/content"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

const (
	ArtifactorySymlink           = "symlink.dest"
	SymlinkSha1                  = "symlink.destsha1"
	LatestBuildNumberKey         = "LATEST"
	lastRelease                  = "LAST_RELEASE"
	buildRepositoriesSuffix      = "-build-info"
	defaultBuildRepositoriesName = "artifactory"
)

func UploadFile(localPath, url, logMsgPrefix string, artifactoryDetails *auth.ServiceDetails, details *fileutils.FileDetails,
	httpClientsDetails httputils.HttpClientDetails, client *jfroghttpclient.JfrogHttpClient, includeChecksums bool,
	progress clientio.ProgressMgr) (*http.Response, []byte, error) {
	var err error
	requestClientDetails := httpClientsDetails.Clone()
	if includeChecksums {
		if details == nil {
			details, err = fileutils.GetFileDetails(localPath, includeChecksums)
		}
		if err != nil {
			return nil, nil, err
		}
		AddChecksumHeaders(requestClientDetails.Headers, details)
	}
	AddAuthHeaders(requestClientDetails.Headers, *artifactoryDetails)
	return client.UploadFile(localPath, url, logMsgPrefix, requestClientDetails, progress)
}

func UploadFileFromReader(reader io.Reader, url string, artifactoryDetails *auth.ServiceDetails, details *fileutils.FileDetails,
	httpClientsDetails httputils.HttpClientDetails, client *jfroghttpclient.JfrogHttpClient) (*http.Response, []byte, error) {
	requestClientDetails := httpClientsDetails.Clone()
	AddChecksumHeaders(requestClientDetails.Headers, details)
	AddAuthHeaders(requestClientDetails.Headers, *artifactoryDetails)

	return client.UploadFileFromReader(reader, url, requestClientDetails, details.Size)
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

// Builds a URL for Artifactory requests.
// Pay attention: semicolons are escaped!
func BuildArtifactoryUrl(baseUrl, path string, params map[string]string) (string, error) {
	u := url.URL{Path: path}
	parsedUrl, err := url.Parse(baseUrl + u.String())
	err = errorutils.CheckError(err)
	if err != nil {
		return "", err
	}
	q := parsedUrl.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	parsedUrl.RawQuery = q.Encode()

	// Semicolons are reserved as separators in some Artifactory APIs, so they'd better be encoded when used for other purposes
	encodedUrl := strings.Replace(parsedUrl.String(), ";", url.QueryEscape(";"), -1)
	return encodedUrl, nil
}

func IsWildcardPattern(pattern string) bool {
	return strings.Contains(pattern, "*") || strings.HasSuffix(pattern, "/") || !strings.Contains(pattern, "/")
}

func GetProjectQueryParam(projectKey string) string {
	if projectKey == "" {
		return ""
	}
	return "?project=" + projectKey
}

// paths - Sorted array.
// index - Index of the current path which we want to check if it a prefix of any of the other previous paths.
// separator - File separator.
// Returns true paths[index] is a prefix of any of the paths[i] where i<index, otherwise returns false.
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
// If no buildNumber provided LATEST will be downloaded.
// If buildName or buildNumber contains "/" (slash) it should be escaped by "\" (backslash).
// Result examples of parsing: "aaa/123" > "aaa"-"123", "aaa" > "aaa"-"LATEST", "aaa\\/aaa" > "aaa/aaa"-"LATEST",  "aaa/12\\/3" > "aaa"-"12/3".
func getBuildNameAndNumberFromBuildIdentifier(buildIdentifier, projectKey string, flags CommonConf) (string, string, error) {
	buildName, buildNumber, err := ParseNameAndVersion(buildIdentifier, true)
	if err != nil {
		return "", "", err
	}
	return GetBuildNameAndNumberFromArtifactory(buildName, buildNumber, projectKey, flags)
}

func GetBuildNameAndNumberFromArtifactory(buildName, buildNumber, projectKey string, flags CommonConf) (string, string, error) {
	if buildNumber == LatestBuildNumberKey || buildNumber == lastRelease {
		return getLatestBuildNumberFromArtifactory(buildName, buildNumber, projectKey, flags)
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
// For module - Parse module name and number.
func ParseNameAndVersion(identifier string, useLatestPolicy bool) (string, string, error) {
	const Delimiter = "/"
	const EscapeChar = "\\"

	if identifier == "" {
		return "", "", nil
	}
	if !strings.Contains(identifier, Delimiter) {
		if useLatestPolicy {
			log.Debug("No '" + Delimiter + "' is found in the build, build number is set to " + LatestBuildNumberKey)
			return identifier, LatestBuildNumberKey, nil
		} else {
			return "", "", errorutils.CheckErrorf("No '" + Delimiter + "' is found in '" + identifier + "'")
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
			log.Debug("No delimiter char (" + Delimiter + ") without escaping char was found in the build, build number is set to " + LatestBuildNumberKey)
			name = identifier
			version = LatestBuildNumberKey
		} else {
			return "", "", errorutils.CheckErrorf("No delimiter char (" + Delimiter + ") without escaping char was found in '" + identifier + "'")
		}
	}
	// Remove escape chars.
	name = strings.Replace(name, "\\/", "/", -1)
	version = strings.Replace(version, "\\/", "/", -1)
	return name, version, nil
}

type Build struct {
	BuildName   string `json:"buildName"`
	BuildNumber string `json:"buildNumber"`
}

func getLatestBuildNumberFromArtifactory(buildName, buildNumber, projectKey string, flags CommonConf) (string, string, error) {
	buildRepo := defaultBuildRepositoriesName
	if projectKey != "" {
		buildRepo = projectKey
	}
	buildRepo += buildRepositoriesSuffix
	aqlBody := CreateAqlQueryForLatestCreated(buildRepo, buildName)
	reader, err := aqlSearch(aqlBody, flags)
	if err != nil {
		return "", "", err
	}
	defer reader.Close()
	for resultItem := new(ResultItem); reader.NextRecord(resultItem) == nil; resultItem = new(ResultItem) {
		if i := strings.LastIndex(resultItem.Name, "-"); i != -1 {
			// Remove the timestamp and .json to get the build number
			buildNumber = resultItem.Name[:i]
			return buildName, buildNumber, nil
		}
	}
	log.Debug(fmt.Sprintf("A build-name: <%s> with a build-number: <%s> could not be found in Artifactory.", buildName, buildNumber))
	return "", "", nil
}

func filterAqlSearchResultsByBuild(specFile *CommonParams, reader *content.ContentReader, flags CommonConf, itemsAlreadyContainProperties bool) (*content.ContentReader, error) {
	var artifactsAqlSearchErr, dependenciesAqlSearchErr error
	var readerWithProps *content.ContentReader
	buildArtifactsSha1 := make(map[string]int)
	buildDependenciesSha1 := make(map[string]int)
	var wg sync.WaitGroup
	wg.Add(2)
	// If 'build-number' is missing in spec file, we fetch the latest from artifactory.
	buildName, buildNumber, err := getBuildNameAndNumberFromBuildIdentifier(specFile.Build, specFile.Project, flags)
	if err != nil {
		return nil, err
	}

	aggregatedBuilds, err := getAggregatedBuilds(buildName, buildNumber, specFile.Project, flags)
	if err != nil {
		return nil, err
	}
	go func() {
		// Get Sha1 for artifacts.
		defer wg.Done()
		if !specFile.ExcludeArtifacts {
			buildArtifactsSha1, artifactsAqlSearchErr = fetchBuildArtifactsOrDependenciesSha1(flags, true, aggregatedBuilds)
		}
	}()

	go func() {
		// Get Sha1 for dependencies.
		defer wg.Done()
		if specFile.IncludeDeps {
			buildDependenciesSha1, dependenciesAqlSearchErr = fetchBuildArtifactsOrDependenciesSha1(flags, false, aggregatedBuilds)
		}
	}()

	if specFile.ExcludeArtifacts || itemsAlreadyContainProperties {
		// No need to add properties to the search results.
		wg.Wait()
		for k, v := range buildDependenciesSha1 {
			buildArtifactsSha1[k] = v
		}
		if artifactsAqlSearchErr != nil {
			return nil, artifactsAqlSearchErr
		}
		if dependenciesAqlSearchErr != nil {
			return nil, dependenciesAqlSearchErr
		}
		return filterBuildAqlSearchResults(reader, buildArtifactsSha1, aggregatedBuilds)
	}

	// Add properties to the previously found artifacts.
	var buildNames []string
	for _, build := range aggregatedBuilds {
		buildNames = append(buildNames, build.BuildName)
	}
	readerWithProps, err = searchProps(specFile.Aql.ItemsFind, "build.name", buildNames, flags)
	if err != nil {
		return nil, err
	}
	defer readerWithProps.Close()
	tempReader, err := loadMissingProperties(reader, readerWithProps)
	if err != nil {
		return nil, err
	}
	defer tempReader.Close()

	wg.Wait()
	// Merge artifacts and dependencies Sha1 maps.
	for k, v := range buildDependenciesSha1 {
		buildArtifactsSha1[k] = v
	}
	if artifactsAqlSearchErr != nil {
		return nil, artifactsAqlSearchErr
	}
	if dependenciesAqlSearchErr != nil {
		return nil, dependenciesAqlSearchErr
	}
	return filterBuildAqlSearchResults(tempReader, buildArtifactsSha1, aggregatedBuilds)
}

// Load all properties to the sorted result items. Save the new result items to a file.
// cr - Sorted result without properties
// crWithProps - Result item with properties
// Return a content reader which points to the result file.
func loadMissingProperties(reader *content.ContentReader, readerWithProps *content.ContentReader) (*content.ContentReader, error) {
	// Key -> Relative path, value -> *ResultItem
	// Contains a limited amount of items from a file, to not overflow memory.
	buffer := make(map[string]*ResultItem)
	var writeOrder []*ResultItem
	var err error
	// Create new file to write result output
	resultFile, err := content.NewContentWriter(content.DefaultKey, true, false)
	if err != nil {
		return nil, err
	}
	defer resultFile.Close()
	for resultItem := new(ResultItem); reader.NextRecord(resultItem) == nil; resultItem = new(ResultItem) {
		buffer[resultItem.GetItemRelativePath()] = resultItem
		// Since maps are an unordered collection, we use slice to save the order of the items
		writeOrder = append(writeOrder, resultItem)
		if len(buffer) == utils.MaxBufferSize {
			// Buffer was full, write all data to a file.
			err = updateProps(readerWithProps, resultFile, buffer, writeOrder)
			if err != nil {
				return nil, err
			}
			buffer = make(map[string]*ResultItem)
			writeOrder = make([]*ResultItem, 0)
		}
	}
	if reader.GetError() != nil {
		return nil, err
	}
	reader.Reset()
	if err := updateProps(readerWithProps, resultFile, buffer, writeOrder); err != nil {
		return nil, err
	}
	return content.NewContentReader(resultFile.GetFilePath(), content.DefaultKey), nil
}

// Load the properties from readerWithProps into buffer's ResultItem and write its values into the resultWriter.
// buffer - Search result buffer Key -> relative path, value -> ResultItem. We use this to load the props into the item by matching the uniqueness of relevant path.
// crWithProps - File containing all the results with proprties.
// writeOrder - List of sorted buffer's searchResults(Map is an unordered collection).
// resultWriter - Search results (sorted) with props.
func updateProps(readerWithProps *content.ContentReader, resultWriter *content.ContentWriter, buffer map[string]*ResultItem, writeOrder []*ResultItem) error {
	if len(buffer) == 0 {
		return nil
	}
	// Load buffer items with their properties.
	for resultItem := new(ResultItem); readerWithProps.NextRecord(resultItem) == nil; resultItem = new(ResultItem) {
		if value, ok := buffer[resultItem.GetItemRelativePath()]; ok {
			value.Properties = resultItem.Properties
		}
	}
	if err := readerWithProps.GetError(); err != nil {
		return err
	}
	readerWithProps.Reset()
	// Write the items to a file with the same search result order.
	for _, itemToWrite := range writeOrder {
		resultWriter.Write(*itemToWrite)
	}
	return nil
}

// Run AQL to retrieve artifacts or dependencies which are associated with a specific build.
// Return a map of the items' SHA1.
func fetchBuildArtifactsOrDependenciesSha1(flags CommonConf, artifacts bool, builds []Build) (map[string]int, error) {
	buildQuery := createAqlQueryForBuild(buildIncludeQueryPart([]string{"name", "repo", "path", "actual_sha1"}), artifacts, builds)
	reader, err := aqlSearch(buildQuery, flags)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return extractSha1FromAqlResponse(reader)
}

// Find artifacts with a specific property.
// aqlBody - AQL to execute together with property filter.
// filterByPropName - Property name to filter.
// filterByPropValue - Property value to filter.
// flags - Command flags for AQL execution.
func searchProps(aqlBody, filterByPropName string, filterByPropValues []string, flags CommonConf) (*content.ContentReader, error) {
	return ExecAqlSaveToFile(createPropsQuery(aqlBody, filterByPropName, filterByPropValues), flags)
}

// Gets a reader of AQL results, and return map with all the SHA1's as keys.
// The values for all the keys in the map is 2
func extractSha1FromAqlResponse(reader *content.ContentReader) (elementsMap map[string]int, err error) {
	elementsMap = make(map[string]int)
	for resultItem := new(ResultItem); reader.NextRecord(resultItem) == nil; resultItem = new(ResultItem) {
		elementsMap[resultItem.Actual_Sha1] = 2
	}
	if err = reader.GetError(); err != nil {
		return
	}
	reader.Reset()
	return
}

// Returns a filtered search result file.
// Map each search result in one of three priority files:
// 1st priority: Match {Sha1, build name, build number}
// 2nd priority: Match {Sha1, build name}
// 3rd priority: Match {Sha1}
// As a result, any duplicated search result item will be split into a different priority list.
// Then merge all the priority list into a single file, so each item is present once in the result file according to the priority list.
// Side note: For each priority level, a single SHA1 can match multi artifacts under different modules.
// reader - Reader of the aql result.
// buildArtifactsSha - Map of all the build-name's sha1 as keys and int as its values. The int value represents priority wheres 0 is a high priority and 2 is lowest.
func filterBuildAqlSearchResults(reader *content.ContentReader, buildArtifactsSha map[string]int, builds []Build) (*content.ContentReader, error) {
	priorityArray, err := createPrioritiesFiles()
	if err != nil {
		return nil, err
	}
	resultCw, err := content.NewContentWriter(content.DefaultKey, true, false)
	if err != nil {
		return nil, err
	}
	defer resultCw.Close()
	// Step 1 - Fill the priority files with search results.
	for resultItem := new(ResultItem); reader.NextRecord(resultItem) == nil; resultItem = new(ResultItem) {
		if _, ok := buildArtifactsSha[resultItem.Actual_Sha1]; !ok {
			continue
		}
		resultBuildName, resultBuildNumber := getBuildNameAndNumberFromProps(resultItem.Properties)
		if isBuildContained(resultBuildName, resultBuildNumber, builds) {
			priorityArray[0].Write(*resultItem)
			buildArtifactsSha[resultItem.Actual_Sha1] = 0
			continue
		}
		if isBuildNameContained(resultBuildName, builds) && buildArtifactsSha[resultItem.Actual_Sha1] != 0 {
			priorityArray[1].Write(*resultItem)
			buildArtifactsSha[resultItem.Actual_Sha1] = 1
			continue
		}
		if buildArtifactsSha[resultItem.Actual_Sha1] == 2 {
			priorityArray[2].Write(*resultItem)
		}
	}
	if err = reader.GetError(); err != nil {
		return nil, err
	}
	reader.Reset()
	var priorityLevel = 0
	// Step 2 - Append the files to the final results file.
	// Scan each priority artifacts and apply them to the final result, skip results that have been already written, by higher priority.
	for _, priority := range priorityArray {
		if err = priority.Close(); err != nil {
			return nil, err
		}
		temp := content.NewContentReader(priority.GetFilePath(), content.DefaultKey)
		for resultItem := new(ResultItem); temp.NextRecord(resultItem) == nil; resultItem = new(ResultItem) {
			if buildArtifactsSha[resultItem.Actual_Sha1] == priorityLevel {
				resultCw.Write(*resultItem)
				// Remove item from map to avoid duplicates.
				delete(buildArtifactsSha, resultItem.Actual_Sha1)
			}
		}
		if err = temp.GetError(); err != nil {
			return nil, err
		}
		if err = temp.Close(); err != nil {
			return nil, err
		}
		priorityLevel++
	}
	return content.NewContentReader(resultCw.GetFilePath(), content.DefaultKey), nil
}

// Return true if the input buildName and buildNumber contained in the builds array.
func isBuildContained(buildName, buildNumber string, builds []Build) bool {
	for _, build := range builds {
		if build.BuildName == buildName && build.BuildNumber == buildNumber {
			return true
		}
	}
	return false
}

// Return true if the input buildName contained in the builds array.
func isBuildNameContained(buildName string, builds []Build) bool {
	for _, build := range builds {
		if build.BuildName == buildName {
			return true
		}
	}
	return false
}

// Create priority files.
func createPrioritiesFiles() ([]*content.ContentWriter, error) {
	firstPriority, err := content.NewContentWriter(content.DefaultKey, true, false)
	if err != nil {
		return nil, err
	}
	secondPriority, err := content.NewContentWriter(content.DefaultKey, true, false)
	if err != nil {
		return nil, err
	}
	thirdPriority, err := content.NewContentWriter(content.DefaultKey, true, false)
	if err != nil {
		return nil, err
	}
	return []*content.ContentWriter{firstPriority, secondPriority, thirdPriority}, nil
}

func GetBuildInfo(buildName, buildNumber, projectKey string, flags CommonConf) (pbi *buildinfo.PublishedBuildInfo, found bool, err error) {
	// Resolve LATEST build number from Artifactory if required.
	name, number, err := GetBuildNameAndNumberFromArtifactory(buildName, buildNumber, projectKey, flags)
	if err != nil {
		return nil, false, err
	}

	// Get build-info json from Artifactory.
	httpClientsDetails := flags.GetArtifactoryDetails().CreateHttpClientDetails()
	restApi := path.Join("api/build/", name, number)

	queryParams := make(map[string]string)
	if projectKey != "" {
		queryParams["project"] = projectKey
	}

	requestFullUrl, err := BuildArtifactoryUrl(flags.GetArtifactoryDetails().GetUrl(), restApi, queryParams)
	if err != nil {
		return nil, false, err
	}

	httpClient := flags.GetJfrogHttpClient()
	log.Debug("Getting build-info from: ", requestFullUrl)
	resp, body, _, err := httpClient.SendGet(requestFullUrl, true, &httpClientsDetails)
	if err != nil {
		return nil, false, err
	}
	if resp.StatusCode == http.StatusNotFound {
		log.Debug("Artifactory response: " + resp.Status + "\n" + utils.IndentJson(body))
		return nil, false, nil
	}
	if err = errorutils.CheckResponseStatus(resp, body, http.StatusOK); err != nil {
		return nil, false, err
	}

	// Build BuildInfo struct from json.
	publishedBuildInfo := &buildinfo.PublishedBuildInfo{}
	if err := json.Unmarshal(body, publishedBuildInfo); err != nil {
		return nil, true, err
	}

	return publishedBuildInfo, true, nil
}

// Recursively, aggregate all transitive builds of the input buildName and buildNumber.
// Build B is considered transitive of build A if the 2 following conditions are met:
// 1. B is a submodule of A or another build that is a transitive of A (direct or indirect descendant).
// 2. B is a module with "Build" type.
func getAggregatedBuilds(buildName, buildNumber, projectKey string, flags CommonConf) ([]Build, error) {
	buildInfo, _, err := GetBuildInfo(buildName, buildNumber, projectKey, flags)
	if err != nil || buildInfo == nil {
		return []Build{}, err
	}
	aggregatedBuilds := []Build{{
		BuildName:   buildName,
		BuildNumber: buildNumber,
	}}
	for _, module := range buildInfo.BuildInfo.Modules {
		if module.Type == buildinfo.Build {
			name, version, err := ParseNameAndVersion(module.Id, false)
			if err != nil {
				return []Build{}, err
			}
			childAggregatedBuilds, err := getAggregatedBuilds(name, version, projectKey, flags)
			if err != nil {
				return []Build{}, err
			}
			aggregatedBuilds = append(aggregatedBuilds, childAggregatedBuilds...)
		}
	}
	return aggregatedBuilds, nil
}

type CommonConf interface {
	GetArtifactoryDetails() auth.ServiceDetails
	GetJfrogHttpClient() *jfroghttpclient.JfrogHttpClient
}

type CommonConfImpl struct {
	client     *jfroghttpclient.JfrogHttpClient
	artDetails *auth.ServiceDetails
}

func NewCommonConfImpl(artDetails auth.ServiceDetails) (CommonConf, error) {
	client, err := jfroghttpclient.JfrogClientBuilder().
		SetClientCertPath(artDetails.GetClientCertPath()).
		SetClientCertKeyPath(artDetails.GetClientCertKeyPath()).
		AppendPreRequestInterceptor(artDetails.RunPreRequestFunctions).
		Build()
	if err != nil {
		return nil, err
	}
	return &CommonConfImpl{artDetails: &artDetails, client: client}, nil
}

func (flags *CommonConfImpl) GetArtifactoryDetails() auth.ServiceDetails {
	return *flags.artDetails
}

func (flags *CommonConfImpl) GetJfrogHttpClient() *jfroghttpclient.JfrogHttpClient {
	return flags.client
}
