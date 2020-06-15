package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/jfrog/jfrog-client-go/artifactory/buildinfo"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/content"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type RequiredArtifactProps int

// This enum defines which properties are required in the result of the aql.
// For example, when performing a copy/move command - the props are not needed, so we set RequiredArtifactProps to NONE.
const (
	MAX_BUFFER_SIZE                       = 50000
	ALL             RequiredArtifactProps = iota
	SYMLINK
	NONE
)

// Use this function when searching by build without pattern or aql.
// This will prevent unnecessary search upon all Artifactory.
func SearchBySpecWithBuild(specFile *ArtifactoryCommonParams, flags CommonConf) ([]ResultItem, error) {
	buildName, buildNumber, err := getBuildNameAndNumberFromBuildIdentifier(specFile.Build, flags)
	if err != nil {
		return nil, err
	}
	specFile.Aql = Aql{ItemsFind: createAqlBodyForBuild(buildName, buildNumber)}

	executionQuery := BuildQueryFromSpecFile(specFile, ALL)
	results, err := aqlSearch(executionQuery, flags)
	if err != nil {
		return nil, err
	}

	// If artifacts' properties weren't fetched in previous aql, fetch now and add to results.
	if !includePropertiesInAqlForSpec(specFile) {
		err = searchAndAddPropsToAqlResult(results, specFile.Aql.ItemsFind, "build.name", buildName, flags)
		if err != nil {
			return nil, err
		}
	}

	// Extract artifacts sha1 for filtering.
	buildArtifactsSha1, err := extractSha1FromAqlResponse(results)
	// Filter artifacts by priorities.
	return filterBuildAqlSearchResults(&results, &buildArtifactsSha1, buildName, buildNumber), err
}

// Search with builds returns many results, some are not part of the build and others may be duplicated of the same artifact.
// To shrink the results:
// 1. Save build-name's sha1s(all the artifact's sha1 that is bound to build-name build)
// 2. Remove artifacts that not are present on the sha1 list
// 3. If we have more than one artifacts with the same sha1:
// 	3.1 Compare the build-name & build-number among all the artifact with the same sha1.
func SearchBySpecWithBuildSaveToFile(specFile *ArtifactoryCommonParams, flags CommonConf) (*content.ContentReader, error) {
	buildName, buildNumber, err := getBuildNameAndNumberFromBuildIdentifier(specFile.Build, flags)
	if err != nil {
		return nil, err
	}
	specFile.Aql = Aql{ItemsFind: createAqlBodyForBuild(buildName, buildNumber)}
	executionQuery := BuildQueryFromSpecFile(specFile, ALL)
	cr, err := aqlSearchSaveToFile(executionQuery, flags)
	if err != nil {
		return nil, err
	}

	// If artifacts' properties weren't fetched in previous aql, fetch now and add to results.
	if !includePropertiesInAqlForSpec(specFile) {
		crWithProps, err := searchPropsSaveToFile(specFile.Aql.ItemsFind, "build.name", buildName, flags)
		if err != nil {
			return nil, err
		}
		cr, err = loadMissingProperties(cr, crWithProps)
		if err != nil {
			return nil, err
		}
	}

	buildArtifactsSha1, err := extractSha1AndPropertyFromAqlResponseSaveToFile(cr)
	return filterBuildAqlSearchResultsSaveToFile(cr, buildArtifactsSha1, buildName, buildNumber)
}

// Perform search by pattern.
func SearchBySpecWithPattern(specFile *ArtifactoryCommonParams, flags CommonConf, requiredArtifactProps RequiredArtifactProps) ([]ResultItem, error) {
	// Create AQL according to spec fields.
	query, err := CreateAqlBodyForSpecWithPattern(specFile)
	if err != nil {
		return nil, err
	}
	specFile.Aql = Aql{ItemsFind: query}
	return SearchBySpecWithAql(specFile, flags, requiredArtifactProps)
}

// Perform search by pattern.
func SearchBySpecWithPatternSaveToFile(specFile *ArtifactoryCommonParams, flags CommonConf, requiredArtifactProps RequiredArtifactProps) (*content.ContentReader, error) {
	// Create AQL according to spec fields.
	query, err := CreateAqlBodyForSpecWithPattern(specFile)
	if err != nil {
		return nil, err
	}
	specFile.Aql = Aql{ItemsFind: query}
	return SearchBySpecWithAqlSaveToFile(specFile, flags, requiredArtifactProps)
}

// Use this function when running Aql with pattern
func SearchBySpecWithAqlSaveToFile(specFile *ArtifactoryCommonParams, flags CommonConf, requiredArtifactProps RequiredArtifactProps) (*content.ContentReader, error) {
	// Execute the search according to provided aql in specFile.
	var crWithProps *content.ContentReader
	query := BuildQueryFromSpecFile(specFile, requiredArtifactProps)
	cr, err := aqlSearchSaveToFile(query, flags)
	if err != nil {
		return nil, err
	}
	isEmpty, err := cr.IsEmpty()
	if err != nil {
		return nil, err
	}
	// Filter results by build.
	if specFile.Build != "" && !isEmpty {
		// If requiredArtifactProps is not NONE and 'includePropertiesInAqlForSpec' for specFile returned true, results contains properties for artifacts.
		resultsArtifactsIncludeProperties := requiredArtifactProps != NONE && includePropertiesInAqlForSpec(specFile)
		cr, err = filterAqlSearchResultsByBuildSaveToFile(specFile, cr, flags, resultsArtifactsIncludeProperties)
		if err != nil {
			return nil, err
		}
	}

	// If:
	// 1. Properties weren't included in 'results'.
	// AND
	// 2. Properties weren't fetched during 'build' filtering
	// Then: we should fetch them now.
	if !includePropertiesInAqlForSpec(specFile) && specFile.Build == "" {
		switch requiredArtifactProps {
		case ALL:
			crWithProps, err = searchPropsSaveToFile(specFile.Aql.ItemsFind, "*", "*", flags)
			break
		case SYMLINK:
			crWithProps, err = searchPropsSaveToFile(specFile.Aql.ItemsFind, "symlink.dest", "*", flags)
			break
		}
		if err != nil {
			return nil, err
		}
		cr, err = loadMissingProperties(cr, crWithProps)
		if err != nil {
			return nil, err
		}
	}
	cr.Reset()
	return cr, err
}

// Use this function when running Aql with pattern
func SearchBySpecWithAql(specFile *ArtifactoryCommonParams, flags CommonConf, requiredArtifactProps RequiredArtifactProps) ([]ResultItem, error) {
	// Execute the search according to provided aql in specFile.
	query := BuildQueryFromSpecFile(specFile, requiredArtifactProps)
	results, err := aqlSearch(query, flags)
	if err != nil {
		return nil, err
	}

	// Filter results by build.
	if specFile.Build != "" && len(results) > 0 {
		// If requiredArtifactProps is not NONE and 'includePropertiesInAqlForSpec' for specFile returned true, results contains properties for artifacts.
		resultsArtifactsIncludeProperties := requiredArtifactProps != NONE && includePropertiesInAqlForSpec(specFile)
		results, err = filterAqlSearchResultsByBuild(specFile, results, flags, resultsArtifactsIncludeProperties)
		if err != nil {
			return nil, err
		}
	}

	// If:
	// 1. Properties weren't included in 'results'.
	// AND
	// 2. Properties weren't fetched during 'build' filtering
	// Then: we should fetch them now.
	if !includePropertiesInAqlForSpec(specFile) && specFile.Build == "" {
		switch requiredArtifactProps {
		case ALL:
			err = searchAndAddPropsToAqlResult(results, specFile.Aql.ItemsFind, "*", "*", flags)
			break
		case SYMLINK:
			err = searchAndAddPropsToAqlResult(results, specFile.Aql.ItemsFind, "symlink.dest", "*", flags)
			break
		}
		if err != nil {
			return nil, err
		}
	}
	return results, err
}

func aqlSearch(aqlQuery string, flags CommonConf) ([]ResultItem, error) {
	json, err := ExecAql(aqlQuery, flags)
	if err != nil {
		return nil, err
	}

	resultItems, err := parseAqlSearchResponse(json)
	return resultItems, err
}

func aqlSearchSaveToFile(aqlQuery string, flags CommonConf) (*content.ContentReader, error) {
	return ExecAqlSaveToFile(aqlQuery, flags)
}

func ExecAql(aqlQuery string, flags CommonConf) ([]byte, error) {
	client, err := flags.GetJfrogHttpClient()
	if err != nil {
		return nil, err
	}
	aqlUrl := flags.GetArtifactoryDetails().GetUrl() + "api/search/aql"
	log.Debug("Searching Artifactory using AQL query:\n", aqlQuery)

	httpClientsDetails := flags.GetArtifactoryDetails().CreateHttpClientDetails()
	resp, body, err := client.SendPost(aqlUrl, []byte(aqlQuery), &httpClientsDetails)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errorutils.CheckError(errors.New("Artifactory response: " + resp.Status + "\n" + utils.IndentJson(body)))
	}

	log.Debug("Artifactory response: ", resp.Status)
	return body, err
}

func ExecAqlSaveToFile(aqlQuery string, flags CommonConf) (*content.ContentReader, error) {
	client, err := flags.GetJfrogHttpClient()
	if err != nil {
		return nil, err
	}
	aqlUrl := flags.GetArtifactoryDetails().GetUrl() + "api/search/aql"
	log.Debug("Searching Artifactory using AQL query:\n", aqlQuery)
	httpClientsDetails := flags.GetArtifactoryDetails().CreateHttpClientDetails()
	resp, cr, err := client.SendPostBodyToFile(aqlUrl, []byte(aqlQuery), &httpClientsDetails)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errorutils.CheckError(errors.New("Artifactory response: " + resp.Status + "\n"))
	}
	log.Debug("Artifactory response: ", resp.Status)
	return cr, err
}

func LogSearchResults(numOfArtifacts int) {
	var msgSuffix = "artifacts."
	if numOfArtifacts == 1 {
		msgSuffix = "artifact."
	}
	log.Info("Found", strconv.Itoa(numOfArtifacts), msgSuffix)
}

func parseAqlSearchResponse(resp []byte) ([]ResultItem, error) {
	var result AqlSearchResult
	err := json.Unmarshal(resp, &result)
	if errorutils.CheckError(err) != nil {
		return nil, err
	}
	return result.Results, nil
}

type AqlSearchResult struct {
	Results []ResultItem
}

type ResultItem struct {
	Repo        string     `json:"repo,omitempty"`
	Path        string     `json:"path,omitempty"`
	Name        string     `json:"name,omitempty"`
	Actual_Md5  string     `json:"actual_md5,omitempty"`
	Actual_Sha1 string     `json:"actual_sha1,omitempty"`
	Size        int64      `json:"size,omitempty"`
	Created     string     `json:"created,omitempty"`
	Modified    string     `json:"modified,omitempty"`
	Properties  []Property `json:"properties,omitempty"`
	Type        string     `json:"type,omitempty"`
}

func (item ResultItem) GetItemRelativePath() string {
	if item.Path == "." {
		return path.Join(item.Repo, item.Name)
	}

	url := item.Repo
	url = addSeparator(url, "/", item.Path)
	url = addSeparator(url, "/", item.Name)
	if item.Type == "folder" && !strings.HasSuffix(url, "/") {
		url = url + "/"
	}
	return url
}

func addSeparator(str1, separator, str2 string) string {
	if str2 == "" {
		return str1
	}
	if str1 == "" {
		return str2
	}

	return str1 + separator + str2
}

func (item *ResultItem) ToArtifact() buildinfo.Artifact {
	return buildinfo.Artifact{Name: item.Name, Checksum: &buildinfo.Checksum{Sha1: item.Actual_Sha1, Md5: item.Actual_Md5}, Path: path.Join(item.Repo, item.Path, item.Name)}
}

func (item *ResultItem) ToDependency() buildinfo.Dependency {
	return buildinfo.Dependency{Id: item.Name, Checksum: &buildinfo.Checksum{Sha1: item.Actual_Sha1, Md5: item.Actual_Md5}}
}

type AqlSearchResultItemFilter func(map[string]ResultItem, []string) []ResultItem
type AqlSearchResultItemFilterSaveToFile func(*content.ContentReader) (*content.ContentReader, error)

func FilterBottomChainResults(cr *content.ContentReader) (*content.ContentReader, error) {
	cw, err := content.NewContentWriter("results", true, false)
	if err != nil {
		return nil, err
	}
	var temp string
	for resultItem := new(ResultItem); cr.NextRecord(resultItem) == nil; resultItem = new(ResultItem) {
		rPath := resultItem.GetItemRelativePath()
		if resultItem.Type == "folder" && !strings.HasSuffix(rPath, "/") {
			rPath += "/"
		}
		if temp == "" || !strings.HasPrefix(temp, rPath) {
			cw.Write(*resultItem)
			temp = rPath
		}
	}
	if err := cr.GetError(); err != nil {
		return nil, err
	}
	if err := cw.Close(); err != nil {
		return nil, err
	}
	return content.NewContentReader(cw.GetFilePath(), cw.GetArrayKey()), nil
}

// Reduce the amount of items by saves only the shortest item path for each uniq path e.g.:
// a | a/b | c | e/f -> a | c | e/f
func FilterTopChainResults(cr *content.ContentReader) (*content.ContentReader, error) {
	cw, err := content.NewContentWriter("results", true, false)
	if err != nil {
		return nil, err
	}
	var temp string
	for resultItem := new(ResultItem); cr.NextRecord(resultItem) == nil; resultItem = new(ResultItem) {
		rPath := resultItem.GetItemRelativePath()
		if resultItem.Type == "folder" && !strings.HasSuffix(rPath, "/") {
			rPath += "/"
		}
		if temp == "" || !strings.HasPrefix(rPath, temp) {
			cw.Write(*resultItem)
			temp = rPath
		}
	}
	if err := cr.GetError(); err != nil {
		return nil, err
	}
	if err := cw.Close(); err != nil {
		return nil, err
	}
	return content.NewContentReader(cw.GetFilePath(), cw.GetArrayKey()), nil
}

func ReduceTopChainDirResult(searchResults *content.ContentReader) (*content.ContentReader, error) {
	return ReduceDirResult(searchResults, true, FilterTopChainResults)
}

func ReduceBottomChainDirResult(searchResults *content.ContentReader) (*content.ContentReader, error) {
	return ReduceDirResult(searchResults, false, FilterBottomChainResults)
}

// Reduce Dir results by using the resultsFilter
func ReduceDirResult(searchResults *content.ContentReader, sortIncreasingOrder bool, resultsFilter AqlSearchResultItemFilterSaveToFile) (*content.ContentReader, error) {
	if searchResults == nil || searchResults.Length() == 0 {
		return searchResults, nil
	}
	paths := make(map[string]ResultItem)
	pathsKeys := make([]string, 0, MAX_BUFFER_SIZE)
	sortedFiles := []*content.ContentReader{}
	for resultItem := new(ResultItem); searchResults.NextRecord(resultItem) == nil; resultItem = new(ResultItem) {
		if resultItem.Name == "." {
			continue
		}
		rPath := resultItem.GetItemRelativePath()
		paths[rPath] = *resultItem
		pathsKeys = append(pathsKeys, rPath)
		if len(pathsKeys) == MAX_BUFFER_SIZE {
			sortedFile, err := saveBufferToFile(paths, pathsKeys, sortIncreasingOrder)
			if err != nil {
				return nil, err
			}
			sortedFiles = append(sortedFiles, sortedFile)
			paths = make(map[string]ResultItem)
			pathsKeys = make([]string, MAX_BUFFER_SIZE)
		}
	}
	if err := searchResults.GetError(); err != nil {
		return nil, err
	}
	var sortedFile *content.ContentReader
	if len(pathsKeys) > 0 {
		sortedFile, err := saveBufferToFile(paths, pathsKeys, sortIncreasingOrder)
		if err != nil {
			return nil, err
		}
		sortedFiles = append(sortedFiles, sortedFile)
	}
	sortedFile, err := mergeSortedFiles(sortedFiles)
	if err != nil {
		return nil, err
	}
	return resultsFilter(sortedFile)
}

func saveBufferToFile(paths map[string]ResultItem, pathsKeys []string, sortIncreasingOrder bool) (*content.ContentReader, error) {
	if len(pathsKeys) == 0 {
		return nil, nil
	}
	cw, err := content.NewContentWriter("results", true, false)
	if err != nil {
		return nil, err
	}
	if sortIncreasingOrder {
		sort.Strings(pathsKeys)
	} else {
		sort.Sort(sort.Reverse(sort.StringSlice(pathsKeys)))
	}
	for _, v := range pathsKeys {
		cw.Write(paths[v])
	}
	if err := cw.Close(); err != nil {
		return nil, err
	}
	return content.NewContentReader(cw.GetFilePath(), cw.GetArrayKey()), nil
}

func mergeSortedFiles(sortedFiles []*content.ContentReader) (*content.ContentReader, error) {
	if len(sortedFiles) == 0 {
		cw, err := content.NewEmptyContentWriter("results", true, false)
		return content.NewContentReader(cw.GetFilePath(), cw.GetArrayKey()), err
	}
	if len(sortedFiles) == 1 {
		return sortedFiles[0], nil
	}
	resultWriter, err := content.NewContentWriter("results", true, false)
	if err != nil {
		return nil, err
	}
	arr := make([]*ResultItem, len(sortedFiles))
	for {
		var smallest *ResultItem
		smallestIndex := 0
		for i := 0; i < len(sortedFiles); i++ {
			if arr[i] == nil && sortedFiles[i] != nil {
				arr[i] = new(ResultItem)
				if err := sortedFiles[i].NextRecord(arr[i]); nil != err {
					sortedFiles[i] = nil
					continue
				}
			}
			if smallest == nil || smallest.GetItemRelativePath() > arr[i].GetItemRelativePath() {
				smallest = arr[i]
				smallestIndex = i
			}
		}
		if smallest == nil {
			break
		}
		resultWriter.Write(*smallest)
		arr[smallestIndex] = nil
	}
	if err := resultWriter.Close(); err != nil {
		return nil, err
	}
	return content.NewContentReader(resultWriter.GetFilePath(), resultWriter.GetArrayKey()), nil
}
