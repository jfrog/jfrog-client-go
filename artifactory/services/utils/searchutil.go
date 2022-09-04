package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
	"sync"

	buildinfo "github.com/jfrog/build-info-go/entities"
	"github.com/jfrog/gofrog/version"

	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/content"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
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
// Collect build artifacts and build dependencies separately, then merge the results into one reader.
func SearchBySpecWithBuild(specFile *CommonParams, flags CommonConf) (*content.ContentReader, error) {
	buildName, buildNumber, err := getBuildNameAndNumberFromBuildIdentifier(specFile.Build, specFile.Project, flags)
	if err != nil {
		return nil, err
	}
	aggregatedBuilds, err := getAggregatedBuilds(buildName, buildNumber, specFile.Project, flags)
	if err != nil {
		return nil, err
	}

	// The specified build does not exist, so an empty reader is returned.
	if len(aggregatedBuilds) == 0 {
		return content.NewEmptyContentReader(content.DefaultKey), nil
	}

	var wg sync.WaitGroup
	wg.Add(2)

	// Get build artifacts.
	var artifactsReader *content.ContentReader
	var artErr error
	go func() {
		defer wg.Done()
		if !specFile.ExcludeArtifacts {
			artifactsReader, artErr = getBuildArtifactsForBuildSearch(*specFile, flags, aggregatedBuilds)
		}
	}()

	// Get build dependencies.
	var dependenciesReader *content.ContentReader
	var depErr error
	go func() {
		defer wg.Done()
		if specFile.IncludeDeps {
			dependenciesReader, depErr = getBuildDependenciesForBuildSearch(*specFile, flags, aggregatedBuilds)
		}
	}()

	wg.Wait()
	if artifactsReader != nil {
		defer artifactsReader.Close()
	}
	if dependenciesReader != nil {
		defer dependenciesReader.Close()
	}
	if artErr != nil {
		return nil, artErr
	}
	if depErr != nil {
		return nil, depErr
	}

	return filterBuildArtifactsAndDependencies(artifactsReader, dependenciesReader, specFile, flags, aggregatedBuilds)
}

func getBuildDependenciesForBuildSearch(specFile CommonParams, flags CommonConf, builds []Build) (*content.ContentReader, error) {
	specFile.Aql = Aql{ItemsFind: createAqlBodyForBuildDependencies(builds)}
	executionQuery := BuildQueryFromSpecFile(&specFile, ALL)
	return aqlSearch(executionQuery, flags)
}

func getBuildArtifactsForBuildSearch(specFile CommonParams, flags CommonConf, builds []Build) (*content.ContentReader, error) {
	specFile.Aql = Aql{ItemsFind: createAqlBodyForBuildArtifacts(builds)}
	executionQuery := BuildQueryFromSpecFile(&specFile, ALL)
	return aqlSearch(executionQuery, flags)
}

// Search with builds may return duplicated items, as the search is performed by checksums.
// Some are not part of the build and others may be duplicated of the same artifact.
// 1. Save SHA1 values received for build-name.
// 2. Remove artifacts that not are present on the sha1 list
// 3. If we have more than one artifact with the same sha1:
// 	3.1 Compare the build-name & build-number among all the artifact with the same sha1.
// This will prevent unnecessary search upon all Artifactory:
func filterBuildArtifactsAndDependencies(artifactsReader, dependenciesReader *content.ContentReader, specFile *CommonParams, flags CommonConf, builds []Build) (*content.ContentReader, error) {
	if includePropertiesInAqlForSpec(specFile) {
		// Don't fetch artifacts' properties from Artifactory.
		mergedReader, err := mergeArtifactsAndDependenciesReaders(artifactsReader, dependenciesReader)
		if err != nil {
			return nil, err
		}
		defer mergedReader.Close()
		buildArtifactsSha1, err := extractSha1FromAqlResponse(mergedReader)
		if err != nil {
			return nil, err
		}
		return filterBuildAqlSearchResults(mergedReader, buildArtifactsSha1, builds)
	}

	// Artifacts' properties weren't fetched in previous aql, fetch now and add to results.
	var buildNames []string
	for _, build := range builds {
		buildNames = append(buildNames, build.BuildName)
	}
	readerWithProps, err := searchProps(createAqlBodyForBuildArtifacts(builds), "build.name", buildNames, flags)
	if err != nil {
		return nil, err
	}
	defer readerWithProps.Close()
	artifactsSortedReaderWithProps, err := loadMissingProperties(artifactsReader, readerWithProps)
	if err != nil {
		return nil, err
	}
	defer artifactsSortedReaderWithProps.Close()
	mergedReader, err := mergeArtifactsAndDependenciesReaders(artifactsSortedReaderWithProps, dependenciesReader)
	if err != nil {
		return nil, err
	}
	defer mergedReader.Close()
	buildArtifactsSha1, err := extractSha1FromAqlResponse(mergedReader)
	if err != nil {
		return nil, err
	}
	return filterBuildAqlSearchResults(mergedReader, buildArtifactsSha1, builds)
}

func mergeArtifactsAndDependenciesReaders(artifactsReader, dependenciesReader *content.ContentReader) (*content.ContentReader, error) {
	var readers []*content.ContentReader
	if artifactsReader != nil {
		readers = append(readers, artifactsReader)
	}
	if dependenciesReader != nil {
		readers = append(readers, dependenciesReader)
	}
	return content.MergeReaders(readers, content.DefaultKey)
}

// Perform search by pattern.
func SearchBySpecWithPattern(specFile *CommonParams, flags CommonConf, requiredArtifactProps RequiredArtifactProps) (*content.ContentReader, error) {
	// Create AQL according to spec fields.
	query, err := CreateAqlBodyForSpecWithPattern(specFile)
	if err != nil {
		return nil, err
	}
	specFile.Aql = Aql{ItemsFind: query}
	return SearchBySpecWithAql(specFile, flags, requiredArtifactProps)
}

// Use this function when running Aql with pattern
func SearchBySpecWithAql(specFile *CommonParams, flags CommonConf, requiredArtifactProps RequiredArtifactProps) (*content.ContentReader, error) {
	// Execute the search according to provided aql in specFile.
	var fetchedProps *content.ContentReader
	query := BuildQueryFromSpecFile(specFile, requiredArtifactProps)
	reader, err := aqlSearch(query, flags)
	if err != nil {
		return nil, err
	}
	filteredReader, err := FilterResultsByBuild(specFile, flags, requiredArtifactProps, reader)
	if err != nil {
		return nil, err
	}
	if filteredReader != nil {
		// This one will close the original reader that was used
		// to create the filteredReader (a new pointer will be created by the defer mechanism).
		defer reader.Close()
		// The new reader assignment will not affect the defer statement.
		reader = filteredReader
	}
	fetchedProps, err = fetchProps(specFile, flags, requiredArtifactProps, reader)
	if fetchedProps != nil {
		// Before returning the new reader, we close the one we used to creat it.
		defer reader.Close()
		return fetchedProps, err
	}
	// Returns the open filteredReader or the original reader that returned from the AQL search.
	return reader, err
}

// Filter the results by build, if no build found or items to filter, nil will be returned.
func FilterResultsByBuild(specFile *CommonParams, flags CommonConf, requiredArtifactProps RequiredArtifactProps, reader *content.ContentReader) (*content.ContentReader, error) {
	length, err := reader.Length()
	if err != nil {
		return nil, err
	}
	if specFile.Build != "" && length > 0 {
		// If requiredArtifactProps is not NONE and 'includePropertiesInAqlForSpec' for specFile returned true, results contains properties for artifacts.
		resultsArtifactsIncludeProperties := requiredArtifactProps != NONE && includePropertiesInAqlForSpec(specFile)
		return filterAqlSearchResultsByBuild(specFile, reader, flags, resultsArtifactsIncludeProperties)
	}
	return nil, nil
}

// Fetch properties only if:
// 1. Properties weren't included in 'results'.
// AND
// 2. Properties weren't fetched during 'build' filtering
// Otherwise, nil will be returned
func fetchProps(specFile *CommonParams, flags CommonConf, requiredArtifactProps RequiredArtifactProps, reader *content.ContentReader) (*content.ContentReader, error) {
	if !includePropertiesInAqlForSpec(specFile) && specFile.Build == "" && requiredArtifactProps != NONE {
		var readerWithProps *content.ContentReader
		var err error
		switch requiredArtifactProps {
		case ALL:
			readerWithProps, err = searchProps(specFile.Aql.ItemsFind, "*", []string{"*"}, flags)
		case SYMLINK:
			readerWithProps, err = searchProps(specFile.Aql.ItemsFind, "symlink.dest", []string{"*"}, flags)
		}
		if err != nil {
			return nil, err
		}
		defer readerWithProps.Close()
		return loadMissingProperties(reader, readerWithProps)
	}
	return nil, nil
}

func aqlSearch(aqlQuery string, flags CommonConf) (*content.ContentReader, error) {
	return ExecAqlSaveToFile(aqlQuery, flags)
}

func ExecAql(aqlQuery string, flags CommonConf) (io.ReadCloser, error) {
	client := flags.GetJfrogHttpClient()
	aqlUrl := flags.GetArtifactoryDetails().GetUrl() + "api/search/aql"
	log.Debug("Searching Artifactory using AQL query:\n", aqlQuery)
	httpClientsDetails := flags.GetArtifactoryDetails().CreateHttpClientDetails()
	resp, err := client.SendPostLeaveBodyOpen(aqlUrl, []byte(aqlQuery), &httpClientsDetails)
	if err != nil {
		return nil, err
	}
	if err = errorutils.CheckResponseStatus(resp, http.StatusOK); err != nil {
		return nil, err
	}
	log.Debug("Artifactory response: ", resp.Status)
	return resp.Body, err
}

func ExecAqlSaveToFile(aqlQuery string, flags CommonConf) (reader *content.ContentReader, err error) {
	var body io.ReadCloser
	body, err = ExecAql(aqlQuery, flags)
	if err != nil {
		return
	}
	defer func() {
		if body != nil {
			e := body.Close()
			if err == nil {
				err = errorutils.CheckError(e)
			}
		}
	}()
	log.Debug("Streaming data to file...")
	var filePath string
	filePath, err = streamToFile(body)
	if err != nil {
		return
	}
	log.Debug("Finished streaming data successfully.")
	reader = content.NewContentReader(filePath, content.DefaultKey)
	return
}

// Save the reader output into a temp file.
// return the file path.
func streamToFile(reader io.Reader) (filePath string, err error) {
	var fd *os.File
	bufioReader := bufio.NewReaderSize(reader, 65536)
	fd, err = fileutils.CreateTempFile()
	if err != nil {
		return "", err
	}
	defer func() {
		e := fd.Close()
		if err == nil {
			err = errorutils.CheckError(e)
		}
	}()
	_, err = io.Copy(fd, bufioReader)
	return fd.Name(), errorutils.CheckError(err)
}

func LogSearchResults(numOfArtifacts int) {
	var msgSuffix = "artifacts."
	if numOfArtifacts == 1 {
		msgSuffix = "artifact."
	}
	log.Info("Found", strconv.Itoa(numOfArtifacts), msgSuffix)
}

type AqlSearchResult struct {
	Results []ResultItem
}

// Implement this interface to allow creating 'content.ContentReader' items which can be used with 'searchutils' functions.
type SearchBasedContentItem interface {
	content.SortableContentItem
	GetItemRelativePath() string
	GetName() string
	GetType() string
}

type ResultItem struct {
	Repo        string     `json:"repo,omitempty"`
	Path        string     `json:"path,omitempty"`
	Name        string     `json:"name,omitempty"`
	Created     string     `json:"created,omitempty"`
	Modified    string     `json:"modified,omitempty"`
	Updated     string     `json:"updated,omitempty"`
	CreatedBy   string     `json:"created_by,omitempty"`
	ModifiedBy  string     `json:"modified_by,omitempty"`
	Type        string     `json:"type,omitempty"`
	Actual_Md5  string     `json:"actual_md5,omitempty"`
	Actual_Sha1 string     `json:"actual_sha1,omitempty"`
	Sha256      string     `json:"sha256,omitempty"`
	Size        int64      `json:"size,omitempty"`
	Properties  []Property `json:"properties,omitempty"`
	Stats       []Stat     `json:"stats,omitempty"`
}

type Stat struct {
	Downloaded      string      `json:"downloaded,omitempty"`
	Downloads       json.Number `json:"downloads,omitempty"`
	DownloadedBy    string      `json:"downloaded_by,omitempty"`
	RemoteDownloads json.Number `json:"remote_downloads,omitempty"`
}

func (item ResultItem) GetSortKey() string {
	if item.Type == "folder" {
		return appendFolderSuffix(item.GetItemRelativePath())
	}
	return item.GetItemRelativePath()
}

func (item ResultItem) GetName() string {
	return item.Name
}

func (item ResultItem) GetType() string {
	return item.Type
}

func (item ResultItem) GetItemRelativePath() string {
	if item.Path == "." {
		return path.Join(item.Repo, item.Name)
	}

	url := item.Repo
	url = path.Join(url, item.Path, item.Name)
	if item.Type == "folder" {
		url = appendFolderSuffix(url)
	}
	return url
}

// Returns "item.Repo/item.Path/".
func (item ResultItem) GetItemRelativeLocation() string {
	return path.Join(item.Repo, item.Path) + "/"
}

func (item *ResultItem) ToArtifact() buildinfo.Artifact {
	return buildinfo.Artifact{
		Name: item.Name,
		Checksum: buildinfo.Checksum{
			Sha1:   item.Actual_Sha1,
			Md5:    item.Actual_Md5,
			Sha256: item.Sha256,
		},
		Path: path.Join(item.Path, item.Name),
	}
}

func (item *ResultItem) ToDependency() buildinfo.Dependency {
	return buildinfo.Dependency{
		Id: item.Name,
		Checksum: buildinfo.Checksum{
			Sha1:   item.Actual_Sha1,
			Md5:    item.Actual_Md5,
			Sha256: item.Sha256,
		},
	}
}

type AqlSearchResultItemFilter func(SearchBasedContentItem, *content.ContentReader) (*content.ContentReader, error)

func (item *ResultItem) GetProperty(key string) string {
	for _, prop := range item.Properties {
		if prop.Key == key {
			return prop.Value
		}
	}
	return ""
}

func FilterBottomChainResults(readerRecord SearchBasedContentItem, reader *content.ContentReader) (*content.ContentReader, error) {
	writer, err := content.NewContentWriter(content.DefaultKey, true, false)
	if err != nil {
		return nil, err
	}
	defer writer.Close()

	// Get the expected record type from the reader.
	recordType := reflect.ValueOf(readerRecord).Type()

	var temp string
	for newRecord := (reflect.New(recordType)).Interface(); reader.NextRecord(newRecord) == nil; newRecord = (reflect.New(recordType)).Interface() {
		resultItem, ok := newRecord.(SearchBasedContentItem)
		if !ok {
			return nil, errorutils.CheckErrorf("Reader record is not search-based.")
		}

		if resultItem.GetName() == "." {
			continue
		}
		rPath := resultItem.GetItemRelativePath()
		if !strings.HasSuffix(rPath, "/") {
			rPath += "/"
		}
		if temp == "" || !strings.HasPrefix(temp, rPath) {
			writer.Write(resultItem)
			temp = rPath
		}
	}
	if err := reader.GetError(); err != nil {
		return nil, err
	}
	reader.Reset()
	return content.NewContentReader(writer.GetFilePath(), writer.GetArrayKey()), nil
}

// Reduce the amount of items by saving only the shortest item path for each unique path e.g.:
// a | a/b | c | e/f -> a | c | e/f
func FilterTopChainResults(readerRecord SearchBasedContentItem, reader *content.ContentReader) (*content.ContentReader, error) {
	writer, err := content.NewContentWriter(content.DefaultKey, true, false)
	if err != nil {
		return nil, err
	}
	defer writer.Close()

	// Get the expected record type from the reader.
	recordType := reflect.ValueOf(readerRecord).Type()

	var prevFolder string
	for newRecord := (reflect.New(recordType)).Interface(); reader.NextRecord(newRecord) == nil; newRecord = (reflect.New(recordType)).Interface() {
		resultItem, ok := newRecord.(SearchBasedContentItem)
		if !ok {
			return nil, errorutils.CheckErrorf("Reader record is not search-based.")
		}

		if resultItem.GetName() == "." {
			continue
		}
		rPath := resultItem.GetItemRelativePath()
		if resultItem.GetType() == "folder" && !strings.HasSuffix(rPath, "/") {
			rPath += "/"
		}
		if prevFolder == "" || !strings.HasPrefix(rPath, prevFolder) {
			writer.Write(resultItem)
			if resultItem.GetType() == "folder" {
				prevFolder = rPath
			}
		}
	}
	if err := reader.GetError(); err != nil {
		return nil, err
	}
	reader.Reset()
	return content.NewContentReader(writer.GetFilePath(), writer.GetArrayKey()), nil
}

func ReduceTopChainDirResult(readerRecord SearchBasedContentItem, searchResults *content.ContentReader) (*content.ContentReader, error) {
	return ReduceDirResult(readerRecord, searchResults, true, FilterTopChainResults)
}

func ReduceBottomChainDirResult(readerRecord SearchBasedContentItem, searchResults *content.ContentReader) (*content.ContentReader, error) {
	return ReduceDirResult(readerRecord, searchResults, false, FilterBottomChainResults)
}

// Reduce Dir results by using the resultsFilter.
func ReduceDirResult(readerRecord SearchBasedContentItem, searchResults *content.ContentReader, ascendingOrder bool, resultsFilter AqlSearchResultItemFilter) (*content.ContentReader, error) {
	sortedFile, err := content.SortContentReader(readerRecord, searchResults, ascendingOrder)
	if err != nil {
		return nil, err
	}
	defer sortedFile.Close()
	return resultsFilter(readerRecord, sortedFile)
}

func DisableTransitiveSearchIfNotAllowed(params *CommonParams, artifactoryVersion *version.Version) {
	transitiveSearchMinVersion := "7.17.0"
	if params.Transitive && !artifactoryVersion.AtLeast(transitiveSearchMinVersion) {
		log.Info(fmt.Sprintf("Transitive search is available on Artifactory version %s or higher. Installed Artifactory version: %s. Transitive option is ignored.",
			transitiveSearchMinVersion, artifactoryVersion.GetVersion()))
		params.Transitive = false
	}
}

func appendFolderSuffix(folderPath string) string {
	if !strings.HasSuffix(folderPath, "/") {
		folderPath = folderPath + "/"
	}
	return folderPath
}
