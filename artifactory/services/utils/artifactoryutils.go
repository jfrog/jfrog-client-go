package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"sync"

	rthttpclient "github.com/jfrog/jfrog-client-go/artifactory/httpclient"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/httpclient"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	clientio "github.com/jfrog/jfrog-client-go/utils/io"
	"github.com/jfrog/jfrog-client-go/utils/io/content"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

const (
	ARTIFACTORY_SYMLINK = "symlink.dest"
	SYMLINK_SHA1        = "symlink.destsha1"
	Latest              = "LATEST"
	LastRelease         = "LAST_RELEASE"
)

func UploadFile(localPath, url, logMsgPrefix string, artifactoryDetails *auth.ServiceDetails, details *fileutils.FileDetails,
	httpClientsDetails httputils.HttpClientDetails, client *rthttpclient.ArtifactoryHttpClient, retries int, progress clientio.Progress) (*http.Response, []byte, error) {
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

func filterAqlSearchResultsByBuild(specFile *ArtifactoryCommonParams, cr *content.ContentReader, flags CommonConf, itemsAlreadyContainProperties bool) (*content.ContentReader, error) {
	var addPropsErr, aqlSearchErr error
	var resultItemWithProps *content.ContentReader
	var buildArtifactsSha1 map[string]byte
	var wg sync.WaitGroup
	// If 'build-number' is missing in spec file, we fetch the laster from artifactory.
	buildName, buildNumber, err := getBuildNameAndNumberFromBuildIdentifier(specFile.Build, flags)
	if err != nil {
		return nil, err
	}

	wg.Add(1)
	// Get Sha1 for artifacts by build name and number
	go func() {
		buildArtifactsSha1, aqlSearchErr = fetchBuildArtifactsSha1(buildName, buildNumber, flags)
		wg.Done()
	}()

	if !itemsAlreadyContainProperties {
		wg.Add(1)
		// Add properties to the previously found artifacts (in case properties weren't already fetched from Artifactory)
		go func() {
			defer wg.Done()
			resultItemWithProps, addPropsErr = searchProps(specFile.Aql.ItemsFind, "build.name", buildName, flags)
			if addPropsErr != nil {
				return
			}
			cr, addPropsErr = loadMissingProperties(cr, resultItemWithProps)
		}()
	}

	wg.Wait()
	if aqlSearchErr != nil {
		return nil, aqlSearchErr
	}
	if addPropsErr != nil {
		return nil, addPropsErr
	}
	return filterBuildAqlSearchResults(cr, buildArtifactsSha1, buildName, buildNumber)
}

// cr - File of sorted result item without properties
// crWithProps - File of result item with properties
// Load all properties to the sorted result items. Save the new result items to a file.
// return a content reader which point to the result file.
func loadMissingProperties(cr *content.ContentReader, crWithProps *content.ContentReader) (*content.ContentReader, error) {
	// Key ->.Repo + .Path + Name + Actual_Sha1, value -> ResultItem
	// Contain limited amount of items from a file, to not overflow memory.
	buffer := make(map[string]*ResultItem)
	var err error
	// Create New file to write result output
	resultFile, err := content.NewContentWriter("results", true, false)
	if err != nil {
		return nil, err
	}
	bufferCounter := 0
	for resultItem := new(ResultItem); cr.NextRecord(resultItem) == nil; resultItem = new(ResultItem) {
		// save the item in a buffer.
		if bufferCounter < 50000 {
			buffer[getResultItemKey(*resultItem)] = resultItem
		} else {
			// Buffer was full, write all data to a file.
			err = updateProps(crWithProps, resultFile, buffer)
			if err != nil {
				return nil, err
			}
			// Init buffer.
			bufferCounter = 1
			buffer = make(map[string]*ResultItem)
			buffer[getResultItemKey(*resultItem)] = resultItem
		}
	}
	if err = cr.GetError(); err != nil {
		return nil, err
	}
	if err := updateProps(crWithProps, resultFile, buffer); err != nil {
		return nil, err
	}
	resultFile.Close()
	return content.NewContentReader(resultFile.GetFilePath(), "results"), nil
}

// buffer - hold limited amount of items (sorted)
// crWithProps - file containing all the results with proprties
// cw - file to write sorted result item with properties
func updateProps(crWithProps *content.ContentReader, cw *content.ContentWriter, buffer map[string]*ResultItem) error {
	if len(buffer) == 0 {
		return nil
	}
	// Load buffer items with their properties.
	for resultItem := new(ResultItem); crWithProps.NextRecord(resultItem) == nil; resultItem = new(ResultItem) {
		if _, ok := buffer[getResultItemKey(*resultItem)]; ok {
			buffer[getResultItemKey(*resultItem)].Properties = resultItem.Properties
		}
	}
	if err := crWithProps.GetError(); err != nil {
		return err
	}
	// Write the items to a file.
	for _, v := range buffer {
		cw.Write(*v)
	}
	return nil
}

// Run AQL to retrieve all artifacts associated with a specific build.
// Return a map of the artifacts SHA1.
func fetchBuildArtifactsSha1(buildName, buildNumber string, flags CommonConf) (map[string]byte, error) {
	buildQuery := createAqlQueryForBuild(buildName, buildNumber, buildIncludeQueryPart([]string{"name", "repo", "path", "actual_sha1"}))
	cr, err := aqlSearch(buildQuery, flags)
	if err != nil {
		return nil, err
	}
	buildArtifactsSha, err := extractSha1AndPropertyFromAqlResponse(cr)
	if err != nil {
		return nil, err
	}
	return buildArtifactsSha, nil
}

/*
* ## Refactore##
 * Find artifact properties by the AQL, add them to the result items.
 *
 * resultItems - Artifacts to add properties to.
 * aqlBody - AQL to execute together with property filter.
 * filterByPropName - Property name to filter.
 * filterByPropValue - Property value to filter.
 * flags - Command flags for AQL execution.
*/
func searchProps(aqlBody, filterByPropName, filterByPropValue string, flags CommonConf) (*content.ContentReader, error) {
	return ExecAqlSaveToFile(createPropsQuery(aqlBody, filterByPropName, filterByPropValue), flags)
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

// Return a map of build's sha1: Key -> sha1, Value -> priority level
// By default every sha1 initialize with lowest priority(0 higher, 1 medium,2 lowest)
func extractSha1AndPropertyFromAqlResponse(cr *content.ContentReader) (elementsMap map[string]byte, err error) {
	elementsMap = make(map[string]byte)
	for resultItem := new(ResultItem); cr.NextRecord(resultItem) == nil; resultItem = new(ResultItem) {
		if err != nil {
			return nil, err
		}
		elementsMap[resultItem.Actual_Sha1] = 2
	}
	if err := cr.GetError(); err != nil {
		return nil, err
	}
	cr.Reset()
	return elementsMap, err
}

/*
 * buildArtifactsSha - List of all the build-name's sha1
 * cr - reader of the aql result
 * Returns a filtered search result file.
 *
 * Map each search result into one of the tree priority files:
 * 1st priority: Match {Sha1, build name, build number}
 * 2nd priority: Match {Sha1, build name}
 * 3rd priority: Match {Sha1}
 *
 *As a result, any duplicated search result items will be split into a different priority list.
 *Then merge all the priority list into a single file, so each item is present once in the result file according to the priority list.
 *
 * Side note: for each priority level, a single SHA1 can match multi artifacts under different modules
 */
func filterBuildAqlSearchResults(cr *content.ContentReader, buildArtifactsSha map[string]byte, buildName, buildNumber string) (*content.ContentReader, error) {
	priorityArray, err := createPrioritiesFiles()
	if err != nil {
		return nil, err
	}
	resultCw, err := content.NewContentWriter("results", true, false)
	if err != nil {
		return nil, err
	}
	// Step 1 - Populate 3 priorities files.
	for resultItem := new(ResultItem); cr.NextRecord(resultItem) == nil; resultItem = new(ResultItem) {
		if _, ok := buildArtifactsSha[resultItem.Actual_Sha1]; !ok {
			continue
		}
		resultBuildName, resultBuildNumber := getBuildNameAndNumberFromProps(resultItem.Properties)
		isBuildNameMatched := resultBuildName == buildName
		if isBuildNameMatched && resultBuildNumber == buildNumber {
			priorityArray[0].Write(*resultItem)
			buildArtifactsSha[resultItem.Actual_Sha1] = 0
			continue
		}
		if isBuildNameMatched && buildArtifactsSha[resultItem.Actual_Sha1] != 0 {
			priorityArray[1].Write(*resultItem)
			buildArtifactsSha[resultItem.Actual_Sha1] = 1
			continue
		}
		if buildArtifactsSha[resultItem.Actual_Sha1] == 2 {
			priorityArray[2].Write(*resultItem)
		}
	}
	if err := cr.GetError(); err != nil {
		return nil, err
	}
	var priorityLevel byte = 0
	// Step 2 - Append the files to the final results file.
	// Scan each priority artifacts and apply them to the final result, starting for priority 0 to 2.
	for _, priority := range priorityArray {
		err := priority.Close()
		if err != nil {
			return nil, err
		}
		temp := content.NewContentReader(priority.GetFilePath(), "results")
		for resultItem := new(ResultItem); temp.NextRecord(resultItem) == nil; resultItem = new(ResultItem) {
			if err != nil {
				return nil, err
			}
			if buildArtifactsSha[resultItem.Actual_Sha1] == priorityLevel {
				resultCw.Write(*resultItem)
			}
		}
		if err = temp.GetError(); err != nil {
			return nil, err
		}
		if err != nil {
			return nil, err
		}
		priorityLevel++
	}
	if err := resultCw.Close(); err != nil {
		return nil, err
	}
	return content.NewContentReader(resultCw.GetFilePath(), "results"), nil
}

// Create writers to hold each result item according to its priority.
func createPrioritiesFiles() ([]*content.ContentWriter, error) {
	firstPriority, err := content.NewContentWriter("results", true, false)
	if err != nil {
		return nil, err
	}
	secondPriority, err := content.NewContentWriter("results", true, false)
	if err != nil {
		return nil, err
	}
	thirdPriority, err := content.NewContentWriter("results", true, false)
	if err != nil {
		return nil, err
	}
	return []*content.ContentWriter{firstPriority, secondPriority, thirdPriority}, nil
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
