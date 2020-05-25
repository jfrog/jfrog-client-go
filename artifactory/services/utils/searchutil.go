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
