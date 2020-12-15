package utils

import (
	"bufio"
	"errors"
	"io"
	"net/http"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/jfrog/jfrog-client-go/artifactory/buildinfo"
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
func SearchBySpecWithBuild(specFile *ArtifactoryCommonParams, flags CommonConf) (*content.ContentReader, error) {
	buildName, buildNumber, err := getBuildNameAndNumberFromBuildIdentifier(specFile.Build, flags)
	if err != nil {
		return nil, err
	}
	var wg sync.WaitGroup
	wg.Add(2)

	// Get build artifacts.
	var artifactsReader *content.ContentReader
	var artErr error
	go func() {
		defer wg.Done()
		if !specFile.ExcludeArtifacts {
			artifactsReader, artErr = getBuildArtifactsForBuildSearch(*specFile, flags, buildName, buildNumber)
		}
	}()

	// Get build dependencies.
	var dependenciesReader *content.ContentReader
	var depErr error
	go func() {
		defer wg.Done()
		if specFile.IncludeDeps {
			dependenciesReader, depErr = getBuildDependenciesForBuildSearch(*specFile, flags, buildName, buildNumber)
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
		return nil, err
	}
	if depErr != nil {
		return nil, err
	}

	return filterBuildArtifactsAndDependencies(artifactsReader, dependenciesReader, specFile, flags, buildName, buildNumber)
}

func getBuildDependenciesForBuildSearch(specFile ArtifactoryCommonParams, flags CommonConf, buildName, buildNumber string) (*content.ContentReader, error) {
	specFile.Aql = Aql{ItemsFind: createAqlBodyForBuildDependencies(buildName, buildNumber)}
	executionQuery := BuildQueryFromSpecFile(&specFile, ALL)
	return aqlSearch(executionQuery, flags)
}

func getBuildArtifactsForBuildSearch(specFile ArtifactoryCommonParams, flags CommonConf, buildName, buildNumber string) (*content.ContentReader, error) {
	specFile.Aql = Aql{ItemsFind: createAqlBodyForBuildArtifacts(buildName, buildNumber)}
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
func filterBuildArtifactsAndDependencies(artifactsReader, dependenciesReader *content.ContentReader, specFile *ArtifactoryCommonParams, flags CommonConf, buildName, buildNumber string) (*content.ContentReader, error) {
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
		return filterBuildAqlSearchResults(mergedReader, buildArtifactsSha1, buildName, buildNumber)
	}

	// Artifacts' properties weren't fetched in previous aql, fetch now and add to results.
	readerWithProps, err := searchProps(createAqlBodyForBuildArtifacts(buildName, buildNumber), "build.name", buildName, flags)
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
	return filterBuildAqlSearchResults(mergedReader, buildArtifactsSha1, buildName, buildNumber)
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
func SearchBySpecWithPattern(specFile *ArtifactoryCommonParams, flags CommonConf, requiredArtifactProps RequiredArtifactProps) (*content.ContentReader, error) {
	// Create AQL according to spec fields.
	query, err := CreateAqlBodyForSpecWithPattern(specFile)
	if err != nil {
		return nil, err
	}
	specFile.Aql = Aql{ItemsFind: query}
	return SearchBySpecWithAql(specFile, flags, requiredArtifactProps)
}

// Use this function when running Aql with pattern
func SearchBySpecWithAql(specFile *ArtifactoryCommonParams, flags CommonConf, requiredArtifactProps RequiredArtifactProps) (*content.ContentReader, error) {
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
func FilterResultsByBuild(specFile *ArtifactoryCommonParams, flags CommonConf, requiredArtifactProps RequiredArtifactProps, reader *content.ContentReader) (*content.ContentReader, error) {
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
func fetchProps(specFile *ArtifactoryCommonParams, flags CommonConf, requiredArtifactProps RequiredArtifactProps, reader *content.ContentReader) (*content.ContentReader, error) {
	if !includePropertiesInAqlForSpec(specFile) && specFile.Build == "" && requiredArtifactProps != NONE {
		var readerWithProps *content.ContentReader
		var err error
		switch requiredArtifactProps {
		case ALL:
			readerWithProps, err = searchProps(specFile.Aql.ItemsFind, "*", "*", flags)
		case SYMLINK:
			readerWithProps, err = searchProps(specFile.Aql.ItemsFind, "symlink.dest", "*", flags)
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
	client, err := flags.GetJfrogHttpClient()
	if err != nil {
		return nil, err
	}
	aqlUrl := flags.GetArtifactoryDetails().GetUrl() + "api/search/aql"
	log.Debug("Searching Artifactory using AQL query:\n", aqlQuery)
	httpClientsDetails := flags.GetArtifactoryDetails().CreateHttpClientDetails()
	resp, err := client.SendPostLeaveBodyOpen(aqlUrl, []byte(aqlQuery), &httpClientsDetails)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errorutils.CheckError(errors.New("Artifactory response: " + resp.Status + "\n"))
	}
	log.Debug("Artifactory response: ", resp.Status)
	return resp.Body, err
}

func ExecAqlSaveToFile(aqlQuery string, flags CommonConf) (*content.ContentReader, error) {
	body, err := ExecAql(aqlQuery, flags)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := body.Close()
		if err != nil {
			log.Warn("Could not close connection:" + err.Error() + ".")
		}
	}()
	log.Debug("Streaming data to file...")
	filePath, err := streamToFile(body)
	if err != nil {
		return nil, err
	}
	log.Debug("Finish streaming data successfully.")
	return content.NewContentReader(filePath, content.DefaultKey), err
}

// Save the reader output into a temp file.
// return the file path.
func streamToFile(reader io.Reader) (string, error) {
	var fd *os.File
	bufio := bufio.NewReaderSize(reader, 65536)
	fd, err := fileutils.CreateTempFile()
	if err != nil {
		return "", err
	}
	defer fd.Close()
	_, err = io.Copy(fd, bufio)
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
	Actual_Md5  string     `json:"actual_md5,omitempty"`
	Actual_Sha1 string     `json:"actual_sha1,omitempty"`
	Size        int64      `json:"size,omitempty"`
	Created     string     `json:"created,omitempty"`
	Modified    string     `json:"modified,omitempty"`
	Properties  []Property `json:"properties,omitempty"`
	Type        string     `json:"type,omitempty"`
}

func (item ResultItem) GetSortKey() string {
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
	url = addSeparator(url, "/", item.Path)
	url = addSeparator(url, "/", item.Name)
	if item.Type == "folder" && !strings.HasSuffix(url, "/") {
		url = url + "/"
	}
	return url
}

// Returns "item.Repo/item.Path/" lowercased.
func (item ResultItem) GetItemRelativeLocation() string {
	return strings.ToLower(addSeparator(item.Repo, "/", item.Path) + "/")
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

type AqlSearchResultItemFilter func(SearchBasedContentItem, *content.ContentReader) (*content.ContentReader, error)

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
			return nil, errorutils.CheckError(errors.New("Reader record is not search-based."))
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
			return nil, errorutils.CheckError(errors.New("Reader record is not search-based."))
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
