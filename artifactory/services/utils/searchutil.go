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
	ALL RequiredArtifactProps = iota
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
	resp, cr, err := client.SendPostResponseToFile(aqlUrl, []byte(aqlQuery), &httpClientsDetails)
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
	Repo        string
	Path        string
	Name        string
	Actual_Md5  string
	Actual_Sha1 string
	Size        int64
	Created     string
	Modified    string
	Properties  []Property
	Type        string
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

func FilterBottomChainResults(paths map[string]ResultItem, pathsKeys []string) []ResultItem {
	var result []ResultItem
	sort.Sort(sort.Reverse(sort.StringSlice(pathsKeys)))
	for i, k := range pathsKeys {
		if i == 0 || !IsSubPath(pathsKeys, i, "/") {
			result = append(result, paths[k])
		}
	}

	return result
}

func FilterTopChainResults(paths map[string]ResultItem, pathsKeys []string) []ResultItem {
	sort.Strings(pathsKeys)
	for _, k := range pathsKeys {
		for _, k2 := range pathsKeys {
			prefix := k2
			if paths[k2].Type == "folder" && !strings.HasSuffix(k2, "/") {
				prefix += "/"
			}

			if k != k2 && strings.HasPrefix(k, prefix) {
				delete(paths, k)
				continue
			}
		}
	}

	var result []ResultItem
	for _, v := range paths {
		result = append(result, v)
	}

	return result
}

// Reduce the amount of items by saves only the shortest item path for each uniq path e.g.:
// a | a/b | c | e/f -> a | c | e/f
func FilterTopChainResultsSaveToFile(cr *content.ContentReader) (*content.ContentReader, error) {
	cw, err := content.NewContentWriter("results", true, false)
	if err != nil {
		return nil, err
	}
	dupCr, err := cr.Duplicate()
	if err != nil {
		return nil, err
	}
	for resultItem := new(ResultItem); cr.NextRecord(resultItem) == nil; resultItem = new(ResultItem) {
		skip := false
		for temp := new(ResultItem); dupCr.NextRecord(temp) == nil; temp = new(ResultItem) {
			prefix := temp.GetItemRelativePath()
			if temp.Type == "folder" && !strings.HasSuffix(temp.GetItemRelativePath(), "/") {
				prefix += "/"
			}
			if resultItem.GetItemRelativePath() != temp.GetItemRelativePath() && strings.HasPrefix(resultItem.GetItemRelativePath(), prefix) {
				skip = true
			}
		}
		if err := dupCr.GetError(); err != nil {
			return nil, err
		}
		if !skip {
			cw.Write(*resultItem)
		}
		dupCr.Reset()
	}
	if err := cr.GetError(); err != nil {
		return nil, err
	}
	cw.Close()
	// TODO: Remove this
	// dupCr.Close()
	// cr.Close()
	return content.NewContentReader(cw.GetFilePath(), cw.GetArrayKey()), nil
}

// Reduce Dir results by using the resultsFilter
func ReduceDirResult(searchResults []ResultItem, resultsFilter AqlSearchResultItemFilter) []ResultItem {
	paths := make(map[string]ResultItem)
	pathsKeys := make([]string, 0, len(searchResults))
	for _, file := range searchResults {
		if file.Name == "." {
			continue
		}

		url := file.GetItemRelativePath()
		paths[url] = file
		pathsKeys = append(pathsKeys, url)
	}
	return resultsFilter(paths, pathsKeys)
}

func ReduceDirResultSaveToFile(searchResults *content.ContentReader, resultsFilter AqlSearchResultItemFilterSaveToFile) (*content.ContentReader, error) {
	cw, err := content.NewContentWriter(searchResults.GetArrayKey(), true, false)
	if err != nil {
		return nil, err
	}
	for resultItem := new(ResultItem); searchResults.NextRecord(resultItem) == nil; resultItem = new(ResultItem) {
		if resultItem.Name == "." {
			continue
		}
		cw.Write(*resultItem)
	}
	if err := searchResults.GetError(); err != nil {
		return nil, err
	}
	cw.Close()
	searchResults.SetFilePath(cw.GetFilePath())
	return resultsFilter(searchResults)
}
