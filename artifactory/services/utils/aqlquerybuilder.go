package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// Returns an AQL body string to search file in Artifactory by pattern, according the the specified arguments requirements.
func createAqlBodyForSpecWithPattern(params *ArtifactoryCommonParams) (string, error) {
	searchPattern := prepareSearchPattern(params.Pattern, true)
	pathFilePairs := createRepoPathFileTriples(searchPattern, params.Recursive)
	includeRoot := strings.Count(searchPattern, "/") < 2
	pathPairsSize := len(pathFilePairs)

	propsQueryPart, err := buildPropsQueryPart(params.Props)
	if err != nil {
		return "", err
	}
	itemTypeQuery := buildItemTypeQueryPart(params)
	nePath := buildNePathPart(pathPairsSize == 0 || includeRoot)
	excludeQuery := buildExcludeQueryPart(params.ExcludePatterns, pathPairsSize == 0 || params.Recursive, params.Recursive)

	json := fmt.Sprintf(`{%s"$or":[`, propsQueryPart+itemTypeQuery+nePath+excludeQuery)

	// Get archive search parameters
	archivePathFilePairs := createArchiveSearchParams(params)

	json += handleRepoPathFileTriples(pathFilePairs, archivePathFilePairs, pathPairsSize) + "]}"
	return json, nil
}

func createArchiveSearchParams(params *ArtifactoryCommonParams) []RepoPathFile {
	var archivePathFilePairs []RepoPathFile

	if params.ArchiveEntries != "" {
		archiveSearchPattern := prepareSearchPattern(params.ArchiveEntries, false)
		archivePathFilePairs = createPathFilePairs("", archiveSearchPattern, true)
	}

	return archivePathFilePairs
}

// Handle building aql query when having PathFilePairs
func handleRepoPathFileTriples(pathFilePairs []RepoPathFile, archivePathFilePairs []RepoPathFile, pathPairSize int) string {
	var query string
	archivePathPairSize := len(archivePathFilePairs)

	for i := 0; i < pathPairSize; i++ {
		if archivePathPairSize > 0 {
			query += handleArchiveSearch(pathFilePairs[i], archivePathFilePairs)
		} else {
			query += buildInnerQueryPart(pathFilePairs[i])
		}

		if i+1 < pathPairSize {
			query += ","
		}
	}

	return query
}

// Handle building aql query including archive search
func handleArchiveSearch(triple RepoPathFile, archivePathFilePairs []RepoPathFile) string {
	var query string
	archivePathPairSize := len(archivePathFilePairs)
	for i := 0; i < archivePathPairSize; i++ {
		query += buildInnerArchiveQueryPart(triple, archivePathFilePairs[i].path, archivePathFilePairs[i].file)

		if i+1 < archivePathPairSize {
			query += ","
		}
	}
	return query
}

func createAqlBodyForBuild(buildName, buildNumber string) string {
	itemsPart :=
		`{` +
			`"artifact.module.build.name":"%s",` +
			`"artifact.module.build.number":"%s"` +
			`}`
	return fmt.Sprintf(itemsPart, buildName, buildNumber)
}

func createAqlQueryForBuild(buildName, buildNumber, includeQueryPart string) string {
	queryBody := createAqlBodyForBuild(buildName, buildNumber)
	itemsPart := `items.find(%s)%s`
	return fmt.Sprintf(itemsPart, queryBody, includeQueryPart)
}

//noinspection GoUnusedExportedFunction
func CreateAqlQueryForNpm(npmName, npmVersion string) string {
	itemsPart :=
		`items.find({` +
			`"@npm.name":"%s",` +
			`"@npm.version":"%s"` +
			`})%s`
	return fmt.Sprintf(itemsPart, npmName, npmVersion, buildIncludeQueryPart([]string{"name", "repo", "path", "actual_sha1", "actual_md5"}))
}

func prepareSearchPattern(pattern string, repositoryExists bool) string {
	if strings.HasSuffix(pattern, "/") || (pattern == "" && repositoryExists) {
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
	return fmt.Sprintf(`"@%s":%s`, key, getAqlValue(value))
}

func buildItemTypeQueryPart(params *ArtifactoryCommonParams) string {
	if params.IncludeDirs {
		return `"type":"any",`
	}
	return ""
}

func buildNePathPart(includeRoot bool) string {
	if !includeRoot {
		return `"path":{"$ne":"."},`
	}
	return ""
}

func buildInnerQueryPart(triple RepoPathFile) string {
	innerQueryPattern := `{"$and":` +
		`[{` +
		`"repo":%s,` +
		`"path":%s,` +
		`"name":%s` +
		`}]}`
	return fmt.Sprintf(innerQueryPattern, getAqlValue(triple.repo), getAqlValue(triple.path), getAqlValue(triple.file))
}

func buildInnerArchiveQueryPart(triple RepoPathFile, archivePath, archiveName string) string {
	innerQueryPattern := `{"$and":` +
		`[{` +
		`"repo":%s,` +
		`"path":%s,` +
		`"name":%s,` +
		`"archive.entry.path":%s,` +
		`"archive.entry.name":%s` +
		`}]}`
	return fmt.Sprintf(innerQueryPattern, getAqlValue(triple.repo), getAqlValue(triple.path), getAqlValue(triple.file), getAqlValue(archivePath), getAqlValue(archiveName))
}

func buildExcludeQueryPart(excludePatterns []string, useLocalPath, recursive bool) string {
	if excludePatterns == nil {
		return ""
	}
	excludeQuery := ""
	var excludePairs []RepoPathFile
	for _, excludePattern := range excludePatterns {
		excludePairs = append(excludePairs, createPathFilePairs("", prepareSearchPattern(excludePattern, false), recursive)...)
	}

	for _, excludePair := range excludePairs {
		excludePath := excludePair.path
		if !useLocalPath && excludePath == "." {
			excludePath = "*"
		}
		excludeQuery += fmt.Sprintf(`"$or":[{"path":{"$nmatch": "%s"},"name":{"$nmatch":"%s"}}],`, excludePath, excludePair.file)
	}
	return excludeQuery
}

// Creates a list of basic required return fields. The list will include the sortBy field if needed.
// If requiredArtifactProps is NONE or 'includePropertiesInAqlForSpec' return false,
// "property" field won't be included due to a limitation in the AQL implementation in Artifactory.
func getQueryReturnFields(specFile *ArtifactoryCommonParams, requiredArtifactProps RequiredArtifactProps) []string {
	returnFields := []string{"name", "repo", "path", "actual_md5", "actual_sha1", "size", "type"}
	if !includePropertiesInAqlForSpec(specFile) {
		// Sort dose not work when property is in the include section. in this case we will append properties in later stage.
		return appendMissingFields(specFile.SortBy, returnFields)
	}
	if requiredArtifactProps != NONE {
		// If any prop is needed we just add all the properties to the result.
		return append(returnFields, "property")
	}
	return returnFields
}

// If specFile includes sortBy or limit, the produced AQL won't include property in the include section.
// This due to an Artifactory limitation related to using these flags with props in an AQL statement.
// Meaning - the result won't contain properties.
func includePropertiesInAqlForSpec(specFile *ArtifactoryCommonParams) bool {
	return !(len(specFile.SortBy) > 0 || specFile.Limit > 0)
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
			`"$and":[%s,{%s}]` +
			`})%s`
	return fmt.Sprintf(propsQuery, aqlBody, propKeyValQueryPart, buildIncludeQueryPart([]string{"name", "repo", "path", "actual_sha1", "property"}))
}

func buildIncludeQueryPart(fieldsToInclude []string) string {
	fieldsToInclude = prepareFieldsForQuery(fieldsToInclude)
	return fmt.Sprintf(`.include(%s)`, strings.Join(fieldsToInclude, `,`))
}

// Optimization - If value is a wildcard pattern, return `{"$match":"value"}`. Otherwise, return `"value"`.
func getAqlValue(val string) string {
	var aqlValuePattern string
	if strings.Contains(val, "*") {
		aqlValuePattern = `{"$match":"%s"}`
	} else {
		aqlValuePattern = `"%s"`
	}
	return fmt.Sprintf(aqlValuePattern, val)
}
