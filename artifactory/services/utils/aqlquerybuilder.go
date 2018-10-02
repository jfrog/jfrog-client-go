package utils

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

// Returns an AQL body string to search file in Artifactory according the the specified arguments requirements.
func createAqlBodyForSpec(params *ArtifactoryCommonParams) (string, error) {
	var itemType string
	if params.IncludeDirs {
		itemType = "any"
	}
	searchPattern := prepareSearchPattern(params.Pattern, true)
	repoIndex := strings.Index(searchPattern, "/")

	repo := searchPattern[:repoIndex]
	searchPattern = searchPattern[repoIndex+1:]

	pathFilePairs := createPathFilePairs(searchPattern, params.Recursive)
	includeRoot := strings.LastIndex(searchPattern, "/") < 0
	pathPairsSize := len(pathFilePairs)
	propsQueryPart, err := buildPropsQueryPart(params.Props)
	if err != nil {
		return "", err
	}
	itemTypeQuery := buildItemTypeQueryPart(itemType)
	nePath := buildNePathPart(pathPairsSize == 0 || includeRoot)
	excludeQuery := buildExcludeQueryPart(params.ExcludePatterns, pathPairsSize == 0 || params.Recursive, params.Recursive)

	json := fmt.Sprintf(`{"repo": "%s",%s"$or": [`, repo, propsQueryPart+itemTypeQuery+nePath+excludeQuery)

	// Get archive search parameters
	archivePathFilePairs := createArchiveSearchParams(params)

	if pathPairsSize == 0 {
		json += handleEmptyPathFilePairs(searchPattern, archivePathFilePairs)
	} else {
		json += handlePathFilePairs(pathFilePairs, archivePathFilePairs, pathPairsSize)
	}

	json += "]}"
	return json, nil
}

func createArchiveSearchParams(params *ArtifactoryCommonParams) []PathFilePair {
	var archivePathFilePairs []PathFilePair

	if params.ArchiveEntries != "" {
		archiveSearchPattern := prepareSearchPattern(params.ArchiveEntries, false)
		archivePathFilePairs = createPathFilePairs(archiveSearchPattern, true)
	}

	return archivePathFilePairs
}

// Handle building aql query when having PathFilePairs
func handlePathFilePairs(pathFilePairs []PathFilePair, archivePathFilePairs []PathFilePair, pathPairSize int) string {
	var query string
	archivePathPairSize := len(archivePathFilePairs)

	for i := 0; i < pathPairSize; i++ {
		if archivePathPairSize > 0 {
			query += handleArchiveSearch(pathFilePairs[i].path, pathFilePairs[i].file, archivePathFilePairs)
		} else {
			query += buildInnerQueryPart(pathFilePairs[i].path, pathFilePairs[i].file)
		}

		if i+1 < pathPairSize {
			query += ","
		}
	}

	return query
}

// Handle building aql query when not having PathFilePairs
func handleEmptyPathFilePairs(searchPattern string, archivePathFilePairs []PathFilePair) string {
	var query string
	if len(archivePathFilePairs) > 0 {
		// Archive search
		query = handleArchiveSearch(".", searchPattern, archivePathFilePairs)
	} else {
		// No archive search
		query = buildInnerQueryPart(".", searchPattern)
	}

	return query
}

// Handle building aql query including archive search
func handleArchiveSearch(path, name string, archivePathFilePairs []PathFilePair) string {
	var query string
	archivePathPairSize := len(archivePathFilePairs)
	for i := 0; i < archivePathPairSize; i++ {
		query += buildInnerArchiveQueryPart(path, name, archivePathFilePairs[i].path, archivePathFilePairs[i].file)

		if i+1 < archivePathPairSize {
			query += ","
		}
	}
	return query
}

func createAqlQueryForBuild(buildName, buildNumber string) string {
	itemsPart :=
		`items.find({` +
			`"artifact.module.build.name": "%s",` +
			`"artifact.module.build.number": "%s"` +
			`})%s`
	return fmt.Sprintf(itemsPart, buildName, buildNumber, buildIncludeQueryPart([]string{"name", "repo", "path", "actual_sha1"}))
}

func CreateAqlQueryForNpm(npmName, npmVersion string) string {
	itemsPart :=
		`items.find({` +
			`"@npm.name": "%s",` +
			`"@npm.version": "%s"` +
			`})%s`
	return fmt.Sprintf(itemsPart, npmName, npmVersion, buildIncludeQueryPart([]string{"name", "repo", "path", "actual_sha1", "actual_md5"}))
}

func prepareSearchPattern(pattern string, repositoryExists bool) string {
	if repositoryExists && !strings.Contains(pattern, "/") {
		pattern += "/"
	}
	if strings.HasSuffix(pattern, "/") {
		pattern += "*"
	}

	// Remove parenthesis
	pattern = strings.Replace(pattern, "(", "", -1)
	pattern = strings.Replace(pattern, ")", "", -1)
	return pattern
}

func buildPropsQueryPart(props string) (string, error) {
	if props == "" {
		return "", nil
	}
	properties, err := ParseProperties(props, JoinCommas)
	if err != nil {
		return "", err
	}
	query := ""
	for _, v := range properties.Properties {
		query += buildKeyValQueryPart(v.Key, v.Value) + `,`
	}
	return query, nil
}

func buildKeyValQueryPart(key string, value string) string {
	return fmt.Sprintf(`"@%s": {"$match" : "%s"}`, key, value)
}

func buildItemTypeQueryPart(itemType string) string {
	if itemType != "" {
		return fmt.Sprintf(`"type": {"$eq": "%s"},`, itemType)
	}
	return ""
}

func buildNePathPart(includeRoot bool) string {
	if !includeRoot {
		return `"path": {"$ne": "."},`
	}
	return ""
}

func buildInnerQueryPart(path, name string) string {
	innerQueryPattern := `{"$and":` +
		`[{` +
		`"path": {"$match": "%s"},` +
		`"name": {"$match": "%s"}` +
		`}]}`
	return fmt.Sprintf(innerQueryPattern, path, name)
}

func buildInnerArchiveQueryPart(path, name, archivePath, archiveName string) string {
	innerQueryPattern := `{"$and":` +
		`[{` +
		`"path": {"$match": "%s"},` +
		`"name": {"$match": "%s"},` +
		`"archive.entry.path": {"$match": "%s"},` +
		`"archive.entry.name": {"$match": "%s"}` +
		`}]}`
	return fmt.Sprintf(innerQueryPattern, path, name, archivePath, archiveName)
}

func buildExcludeQueryPart(excludePatterns []string, useLocalPath, recursive bool) string {
	if excludePatterns == nil {
		return ""
	}
	excludeQuery := ""
	var excludePairs []PathFilePair
	for _, excludePattern := range excludePatterns {
		excludePairs = append(excludePairs, createPathFilePairs(prepareSearchPattern(excludePattern, false), recursive)...)
	}

	for _, excludePair := range excludePairs {
		excludePath := excludePair.path
		if !useLocalPath && excludePath == "." {
			excludePath = "*"
		}
		excludeQuery += fmt.Sprintf(`"$or": [{"path": {"$nmatch": "%s"}, "name": {"$nmatch": "%s"}}],`, excludePath, excludePair.file)
	}
	return excludeQuery
}

// We need to translate the provided download pattern to an AQL query.
// In Artifactory, for each artifact the name and path of the artifact are saved separately including folders.
// We therefore need to build an AQL query that covers all possible folders the provided
// pattern can include.
// For example, the pattern a/*b*c*/ can include the two following folders:
// a/b/c, a/bc/, a/x/y/z/b/c/
// To achieve that, this function parses the pattern by splitting it by its * characters.
// The end result is a list of PathFilePair structs.
// Each struct represent a possible path and folder name pair to be included in AQL query with an "or" relationship.
func createPathFolderPairs(searchPattern string) []PathFilePair {
	// Remove parenthesis
	searchPattern = searchPattern[:len(searchPattern)-1]
	searchPattern = strings.Replace(searchPattern, "(", "", -1)
	searchPattern = strings.Replace(searchPattern, ")", "", -1)

	index := strings.Index(searchPattern, "/")
	searchPattern = searchPattern[index+1:]

	index = strings.LastIndex(searchPattern, "/")
	lastSlashPath := searchPattern
	path := "."
	if index != -1 {
		lastSlashPath = searchPattern[index+1:]
		path = searchPattern[:index]
	}

	pairs := []PathFilePair{{path: path, file: lastSlashPath}}
	for i := 0; i < len(lastSlashPath); i++ {
		if string(lastSlashPath[i]) == "*" {
			pairs = append(pairs, PathFilePair{path: filepath.Join(path, lastSlashPath[:i+1]), file: lastSlashPath[i:]})
		}
	}
	return pairs
}

// We need to translate the provided download pattern to an AQL query.
// In Artifactory, for each artifact the name and path of the artifact are saved separately.
// We therefore need to build an AQL query that covers all possible paths and names the provided
// pattern can include.
// For example, the pattern a/* can include the two following file:
// a/file1.tgz and also a/b/file2.tgz
// To achieve that, this function parses the pattern by splitting it by its * characters.
// The end result is a list of PathFilePair structs.
// Each struct represent a possible path and file name pair to be included in AQL query with an "or" relationship.
func createPathFilePairs(pattern string, recursive bool) []PathFilePair {
	var defaultPath string
	if recursive {
		defaultPath = "*"
	} else {
		defaultPath = "."
	}

	pairs := []PathFilePair{}
	if pattern == "*" {
		pairs = append(pairs, PathFilePair{defaultPath, "*"})
		return pairs
	}

	slashIndex := strings.LastIndex(pattern, "/")
	var path string
	var name string
	if slashIndex < 0 {
		pairs = append(pairs, PathFilePair{".", pattern})
		path = ""
		name = pattern
	} else {
		path = pattern[:slashIndex]
		name = pattern[slashIndex+1:]
		pairs = append(pairs, PathFilePair{path, name})
	}
	if !recursive {
		return pairs
	}
	if name == "*" {
		path += "/*"
		pairs = append(pairs, PathFilePair{path, "*"})
		return pairs
	}
	pattern = name

	sections := strings.Split(pattern, "*")
	size := len(sections)
	for i := 0; i < size; i++ {
		options := []string{}
		if i+1 < size {
			options = append(options, sections[i]+"*/")
		}
		for _, option := range options {
			str := ""
			for j := 0; j < size; j++ {
				if j > 0 {
					str += "*"
				}
				if j == i {
					str += option
				} else {
					str += sections[j]
				}
			}
			split := strings.Split(str, "/")
			filePath := split[0]
			fileName := split[1]
			if fileName == "" {
				fileName = "*"
			}
			if path != "" {
				if !strings.HasSuffix(path, "/") {
					path += "/"
				}
			}
			pairs = append(pairs, PathFilePair{path + filePath, fileName})
		}
	}
	return pairs
}

type PathFilePair struct {
	path string
	file string
}

// Creates a list of basic required return fields. The list will include the sortBy field if needed.
// If requiredArtifactProps is NONE or sortBy is configured, "property" field won't be included due to a limitation in the AQL implementation in Artifactory.
func getQueryReturnFields(specFile *ArtifactoryCommonParams, requiredArtifactProps RequiredArtifactProps) []string {
	returnFields := []string{"name", "repo", "path", "actual_md5", "actual_sha1", "size", "type"}
	if specIncludesSortOrLimit(specFile) {
		// Sort dose not work when property is in the include section. in this case we will append properties in later stage.
		return appendMissingFields(specFile.SortBy, returnFields)
	}
	if requiredArtifactProps != NONE {
		// If any prop is needed we just adding all the properties to the result, in order to prevent the second props query.
		return append(returnFields, "property")
	}
	return returnFields
}

func specIncludesSortOrLimit(specFile *ArtifactoryCommonParams) bool {
	return len(specFile.SortBy) > 0 || specFile.Limit > 0
}

func appendMissingFields(fields []string, defaultFields []string) []string {
	for _, field := range fields {
		if !stringIsInSlice(field, defaultFields) {
			defaultFields = append(defaultFields, field)
		}
	}
	return defaultFields
}

func stringIsInSlice(string string, strings []string) bool {
	for _, v := range strings {
		if v == string {
			return true
		}
	}
	return false
}

func prepareFieldsForQuery(fields []string) []string {
	for i, val := range fields {
		fields[i] = `"` + val + `"`
	}
	return fields
}

// Creates an aql query from a spec file.
// If the spec includes sortBy, the produced AQL won't includes property in the include section,
// due to an Artifactory limitation to limitation related to using sort with props in an AQL statement - this mean the result wont contain properties.
// Same will happen if requiredArtifactProps is 'NONE'.
func buildQueryFromSpecFile(specFile *ArtifactoryCommonParams, requiredArtifactProps RequiredArtifactProps) string {
	aqlBody := specFile.Aql.ItemsFind
	query := fmt.Sprintf(`items.find(%s)%s`, aqlBody, buildIncludeQueryPart(getQueryReturnFields(specFile, requiredArtifactProps)))
	query = appendSortQueryPart(specFile, query)
	query = appendOffsetQueryPart(specFile, query)
	return appendLimitQueryPart(specFile, query)
}

func appendOffsetQueryPart(specFile *ArtifactoryCommonParams, query string) string {
	if specFile.Offset > 0 {
		query = fmt.Sprintf(`%s.offset(%s)`, query, strconv.Itoa(specFile.Offset))
	}
	return query
}

func appendLimitQueryPart(specFile *ArtifactoryCommonParams, query string) string {
	if specFile.Limit > 0 {
		query = fmt.Sprintf(`%s.limit(%s)`, query, strconv.Itoa(specFile.Limit))
	}
	return query
}

func appendSortQueryPart(specFile *ArtifactoryCommonParams, query string) string {
	if len(specFile.SortBy) > 0 {
		query = fmt.Sprintf(`%s.sort({%s})`, query, buildSortQueryPart(specFile.SortBy, specFile.SortOrder))
	}
	return query
}

func buildSortQueryPart(sortFields []string, sortOrder string) string {
	if sortOrder == "" {
		sortOrder = "asc"
	}
	return fmt.Sprintf(`"$%s":[%s]`, sortOrder, strings.Join(prepareFieldsForQuery(sortFields), `,`))
}

func createPropsQuery(aqlBody, propKey, propVal string) string {
	propKeyValQueryPart := buildKeyValQueryPart(propKey, propVal)
	propsQuery :=
		`items.find({` +
			`"$and" :[%s,{%s}]` +
			`})%s`
	return fmt.Sprintf(propsQuery, aqlBody, propKeyValQueryPart, buildIncludeQueryPart([]string{"name", "repo", "path", "actual_sha1", "property"}))
}

func buildIncludeQueryPart(fieldsToInclude []string) string {
	fieldsToInclude = prepareFieldsForQuery(fieldsToInclude)
	return fmt.Sprintf(`.include(%s)`, strings.Join(fieldsToInclude, `,`))
}
