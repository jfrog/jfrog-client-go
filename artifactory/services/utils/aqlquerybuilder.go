package utils

import (
	"fmt"
	"golang.org/x/exp/slices"
	"strconv"
	"strings"

	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
)

// Returns an AQL body string to search file in Artifactory by pattern, according the specified arguments requirements.
func CreateAqlBodyForSpecWithPattern(params *CommonParams) (string, error) {
	searchPattern := prepareSourceSearchPattern(params.Pattern, params.Target, true)
	repoPathFileTriples, singleRepo, err := createRepoPathFileTriples(searchPattern, params.Recursive)
	if err != nil {
		return "", err
	}
	if params.Transitive && !singleRepo {
		return "", errorutils.CheckErrorf("When searching or downloading with the transitive setting, the pattern must include a single repository only, meaning wildcards are allowed only after the first slash.")
	}
	includeRoot := strings.Count(searchPattern, "/") < 2
	triplesSize := len(repoPathFileTriples)

	propsQueryPart, err := buildPropsQueryPart(params.Props, params.ExcludeProps)
	if err != nil {
		return "", err
	}
	itemTypeQuery := buildItemTypeQueryPart(params)
	nePath := buildNePathPart(triplesSize == 0 || includeRoot)
	excludeQuery, err := buildExcludeQueryPart(params, triplesSize == 0 || params.Recursive, params.Recursive)
	if err != nil {
		return "", err
	}
	releaseBundle, err := buildReleaseBundleQuery(params)
	if err != nil {
		return "", err
	}

	json := fmt.Sprintf(`{%s"$or":[`, propsQueryPart+itemTypeQuery+nePath+excludeQuery+releaseBundle)

	// Get archive search parameters
	archivePathFilePairs := createArchiveSearchParams(params)

	json += handleRepoPathFileTriples(repoPathFileTriples, archivePathFilePairs, triplesSize) + "]}"
	return json, nil
}

func createArchiveSearchParams(params *CommonParams) []RepoPathFile {
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

func createAqlBodyForBuildArtifacts(builds []Build) string {
	buildArtifactsItem := `{"$and":[{"artifact.module.build.name":"%s","artifact.module.build.number":"%s"}]}`
	var items []string
	for _, build := range builds {
		items = append(items, fmt.Sprintf(buildArtifactsItem, build.BuildName, build.BuildNumber))
	}
	return `{"$or":[` + strings.Join(items, ",") + "]}"
}

func createAqlBodyForBuildDependencies(builds []Build) string {
	buildDependenciesItem := `{"$and":[{"dependency.module.build.name":"%s","dependency.module.build.number":"%s"}]}`
	var items []string
	for _, build := range builds {
		items = append(items, fmt.Sprintf(buildDependenciesItem, build.BuildName, build.BuildNumber))
	}
	return `{"$or":[` + strings.Join(items, ",") + "]}"
}

func createAqlQueryForBuild(includeQueryPart string, artifactsQuery bool, builds []Build) string {
	var queryBody string
	if artifactsQuery {
		queryBody = createAqlBodyForBuildArtifacts(builds)
	} else {
		queryBody = createAqlBodyForBuildDependencies(builds)
	}
	itemsPart := `items.find(%s)%s`
	return fmt.Sprintf(itemsPart, queryBody, includeQueryPart)
}

// noinspection GoUnusedExportedFunction
func CreateAqlQueryForYarn(npmName, npmVersion string) string {
	itemsPart :=
		`items.find({` +
			`"@npm.name":"%s",` +
			`"$or": [` +
			// sometimes the npm.version in the repository is written with "v" prefix, so we search both syntaxes
			`{"@npm.version":"%[2]s"},` +
			`{"@npm.version":"v%[2]s"}` +
			`]` +
			`})%s`
	return fmt.Sprintf(itemsPart, npmName, npmVersion, buildIncludeQueryPart([]string{"name", "repo", "path", "actual_sha1", "actual_md5", "sha256"}))
}

func CreateAqlQueryForPypi(repo, file string) string {
	itemsPart :=
		`items.find({` +
			`"repo": "%s",` +
			`"$or": [{` +
			`"$and":[{` +
			`"path": {"$match": "*"},` +
			`"name": {"$match": "%s"}` +
			`}]` +
			`}]` +
			`})%s`
	return fmt.Sprintf(itemsPart, repo, file, buildIncludeQueryPart([]string{"name", "repo", "path", "actual_md5", "actual_sha1", "sha256"}))
}

func CreateAqlQueryForLatestCreated(repo, path string) string {
	itemsPart :=
		`items.find({` +
			`"repo": "%s",` +
			`"path": {"$match": "%s"}` +
			`})` +
			`.sort({%s})` +
			`.limit(1)`
	return fmt.Sprintf(itemsPart, repo, path, buildSortQueryPart([]string{"created"}, "desc"))
}

func prepareSearchPattern(pattern string, repositoryExists bool) string {
	addWildcardIfNeeded(&pattern, repositoryExists)
	// Remove parenthesis
	pattern = strings.Replace(pattern, "(", "", -1)
	pattern = strings.Replace(pattern, ")", "", -1)
	return pattern
}

func buildPropsQueryPart(props, excludeProps string) (string, error) {
	propsQuery := ""
	properties, err := ParseProperties(props)
	if err != nil {
		return "", err
	}
	for key, values := range properties.ToMap() {
		propsQuery += buildKeyAllValQueryPart(key, values) + `,`
	}

	excludePropsQuery := ""
	excludeProperties, err := ParseProperties(excludeProps)
	if err != nil {
		return "", err
	}
	excludePropsLen := excludeProperties.KeysLen()
	if excludePropsLen > 0 {
		excludePropsQuery = `"$or":[`
		for key, values := range excludeProperties.ToMap() {
			for _, value := range values {
				excludePropsQuery += `{` + buildExcludedKeyValQueryPart(key, value) + `},`
			}
		}
		excludePropsQuery = strings.TrimSuffix(excludePropsQuery, ",") + `],`
	}
	return propsQuery + excludePropsQuery, nil
}

func buildKeyValQueryPart(key string, propValues []string) string {
	var items []string
	for _, value := range propValues {
		items = append(items, fmt.Sprintf(`{"@%s":%s}`, key, getAqlValue(value)))
	}
	return `"$or":[` + strings.Join(items, ",") + "]"
}

func buildKeyAllValQueryPart(key string, propValues []string) string {
	var items []string
	for _, value := range propValues {
		items = append(items, fmt.Sprintf(`{"@%s":%s}`, key, getAqlValue(value)))
	}
	return `"$and":[` + strings.Join(items, ",") + "]"
}

func buildExcludedKeyValQueryPart(key string, value string) string {
	return fmt.Sprintf(`"@%s":{"$ne":%s}`, key, getAqlValue(value))
}

func buildItemTypeQueryPart(params *CommonParams) string {
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

func buildExcludeQueryPart(params *CommonParams, useLocalPath, recursive bool) (string, error) {
	excludeQuery := ""
	var excludeTriples []RepoPathFile
	for _, exclusion := range params.GetExclusions() {
		repoPathFileTriples, _, err := createRepoPathFileTriples(prepareSearchPattern(exclusion, true), recursive)
		if err != nil {
			return "", err
		}
		excludeTriples = append(excludeTriples, repoPathFileTriples...)
	}

	for _, excludeTriple := range excludeTriples {
		excludePath := excludeTriple.path
		if !useLocalPath && excludePath == "." {
			excludePath = "*"
		}
		excludeRepoStr := ""

		// repo="*" may cause an error to be returned from Artifactory in transitive search.
		if excludeTriple.repo != "" && excludeTriple.repo != "*" {
			excludeRepoStr = fmt.Sprintf(`"repo":{"$nmatch":"%s"},`, excludeTriple.repo)
		}
		excludeQuery += fmt.Sprintf(`"$or":[{%s"path":{"$nmatch":"%s"},"name":{"$nmatch":"%s"}}],`, excludeRepoStr, excludePath, excludeTriple.file)
	}
	return excludeQuery, nil
}

func buildReleaseBundleQuery(params *CommonParams) (string, error) {
	bundleName, bundleVersion, err := ParseNameAndVersion(params.Bundle, false)
	if bundleName == "" || err != nil {
		return "", err
	}
	itemsPart := `"$and":` +
		`[{` +
		`"release_artifact.release.name":%s,` +
		`"release_artifact.release.version":%s` +
		`}],`
	return fmt.Sprintf(itemsPart, getAqlValue(bundleName), getAqlValue(bundleVersion)), nil
}

// Creates a list of basic required return fields. The list will include the sortBy field if needed.
// If requiredArtifactProps is NONE or 'includePropertiesInAqlForSpec' return false,
// "property" field won't be included due to a limitation in the AQL implementation in Artifactory.
func getQueryReturnFields(specFile *CommonParams, requiredArtifactProps RequiredArtifactProps) []string {
	returnFields := []string{"name", "repo", "path", "actual_md5", "actual_sha1", "sha256", "size", "type", "modified", "created"}
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
func includePropertiesInAqlForSpec(specFile *CommonParams) bool {
	return !(len(specFile.SortBy) > 0 || specFile.Limit > 0)
}

func appendMissingFields(fields []string, defaultFields []string) []string {
	for _, field := range fields {
		if !slices.Contains(defaultFields, field) {
			defaultFields = append(defaultFields, field)
		}
	}
	return defaultFields
}

func prepareFieldsForQuery(fields []string) []string {
	// Since a slice is basically a pointer, we don't want to modify the underlying fields array because it might be used again (like in delete service)
	// We will create new slice with the quoted values and will return it.
	var queryFields []string
	for _, val := range fields {
		queryFields = append(queryFields, `"`+val+`"`)
	}
	return queryFields
}

// Creates an aql query from a spec file.
func BuildQueryFromSpecFile(specFile *CommonParams, requiredArtifactProps RequiredArtifactProps) string {
	aqlBody := specFile.Aql.ItemsFind
	query := fmt.Sprintf(`items.find(%s)%s`, aqlBody, buildIncludeQueryPart(getQueryReturnFields(specFile, requiredArtifactProps)))
	query = appendSortQueryPart(specFile, query)
	query = appendOffsetQueryPart(specFile, query)
	query = appendTransitiveQueryPart(specFile, query)
	query = appendLimitQueryPart(specFile, query)
	return query
}

func appendOffsetQueryPart(specFile *CommonParams, query string) string {
	if specFile.Offset > 0 {
		query = fmt.Sprintf(`%s.offset(%s)`, query, strconv.Itoa(specFile.Offset))
	}
	return query
}

func appendLimitQueryPart(specFile *CommonParams, query string) string {
	if specFile.Limit > 0 {
		query = fmt.Sprintf(`%s.limit(%s)`, query, strconv.Itoa(specFile.Limit))
	}
	return query
}

func appendSortQueryPart(specFile *CommonParams, query string) string {
	if len(specFile.SortBy) > 0 {
		query = fmt.Sprintf(`%s.sort({%s})`, query, buildSortQueryPart(specFile.SortBy, specFile.SortOrder))
	}
	return query
}

func appendTransitiveQueryPart(specFile *CommonParams, query string) string {
	if specFile.Transitive {
		query = fmt.Sprintf(`%s.transitive()`, query)
	}
	return query
}

func buildSortQueryPart(sortFields []string, sortOrder string) string {
	if sortOrder == "" {
		sortOrder = "asc"
	}
	return fmt.Sprintf(`"$%s":[%s]`, sortOrder, strings.Join(prepareFieldsForQuery(sortFields), `,`))
}

func createPropsQuery(aqlBody, propKey string, propValues []string) string {
	propKeyValQueryPart := buildKeyValQueryPart(propKey, propValues)
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

func prepareSourceSearchPattern(pattern, target string, repositoryExists bool) string {
	addWildcardIfNeeded(&pattern, repositoryExists)
	pattern = utils.RemovePlaceholderParentheses(pattern, target)
	return pattern
}

func addWildcardIfNeeded(pattern *string, repositoryExists bool) {
	if strings.HasSuffix(*pattern, "/") || (*pattern == "" && repositoryExists) {
		*pattern += "*"
	}
}
