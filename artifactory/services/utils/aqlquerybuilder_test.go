package utils

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var buildAqlSearchQueryDataProvider = []struct {
	pattern     string
	recursive   bool
	expectedAql string
}{
	{"repo-local", true,
		`{"$or":[{"$and":[{"repo":"repo-local","path":{"$match":"*"},"name":{"$match":"*"}}]}]}`},
	{"repo-w*ldcard", true,
		`{"$or":[{"$and":[{"repo":{"$match":"repo-w*"},"path":{"$match":"*"},"name":{"$match":"*ldcard"}}]},{"$and":[{"repo":{"$match":"repo-w*ldcard"},"path":{"$match":"*"},"name":{"$match":"*"}}]}]}`},
	{"repo-local2/a*b*c/dd/", true,
		`{"path":{"$ne":"."},"$or":[{"$and":[{"repo":"repo-local2","path":{"$match":"a*b*c/dd"},"name":{"$match":"*"}}]},{"$and":[{"repo":"repo-local2","path":{"$match":"a*b*c/dd/*"},"name":{"$match":"*"}}]}]}`},
	{"repo-local*/a*b*c/dd/", true,
		`{"path":{"$ne":"."},"$or":[{"$and":[{"repo":{"$match":"repo-local*"},"path":{"$match":"*a*b*c/dd"},"name":{"$match":"*"}}]},{"$and":[{"repo":{"$match":"repo-local*"},"path":{"$match":"*a*b*c/dd/*"},"name":{"$match":"*"}}]}]}`},
	{"repo-local", false,
		`{"$or":[{"$and":[{"repo":"repo-local","path":".","name":{"$match":"*"}}]}]}`},
	{"*repo-local", false,
		`{"$or":[{"$and":[{"path":".","name":{"$match":"*repo-local"}}]},{"$and":[{"repo":{"$match":"*repo-local"},"path":".","name":{"$match":"*"}}]}]}`},
	{"repo-local2/a*b*c/dd/", false,
		`{"path":{"$ne":"."},"$or":[{"$and":[{"repo":"repo-local2","path":{"$match":"a*b*c/dd"},"name":{"$match":"*"}}]}]}`},
	{"*/a*b*c/dd/", false,
		`{"path":{"$ne":"."},"$or":[{"$and":[{"path":{"$match":"*a*b*c/dd"},"name":{"$match":"*"}}]}]}`},
	{"**/a-.*.tar.gz", false,
		`{"$or":[{"$and":[{"path":{"$match":"**"},"name":{"$match":"a-.*.tar.gz"}}]},{"$and":[{"path":".","name":{"$match":"*a-.*.tar.gz"}}]}]}`},
}

func TestBuildAqlSearchQuery(t *testing.T) {
	for _, sample := range buildAqlSearchQueryDataProvider {
		t.Run(sample.pattern+"_recursive_"+strconv.FormatBool(sample.recursive), func(t *testing.T) {
			params := CommonParams{Pattern: sample.pattern, Recursive: sample.recursive, Regexp: false, IncludeDirs: false}
			aqlResult, err := CreateAqlBodyForSpecWithPattern(&params)
			assert.NoError(t, err)
			if aqlResult != sample.expectedAql {
				t.Error("Unexpected download AQL query built. \nExpected: " + sample.expectedAql + " \nGot:      " + aqlResult)
			}
		})
	}
}

func TestCommonParams(t *testing.T) {
	artifactoryParams := CommonParams{}
	assertIsSortLimitSpecBool(t, !includePropertiesInAqlForSpec(&artifactoryParams), false)

	artifactoryParams.SortBy = []string{"Vava", "Bubu"}
	assertIsSortLimitSpecBool(t, !includePropertiesInAqlForSpec(&artifactoryParams), true)

	artifactoryParams.SortBy = nil
	artifactoryParams.Limit = 0
	assertIsSortLimitSpecBool(t, !includePropertiesInAqlForSpec(&artifactoryParams), false)

	artifactoryParams.Limit = -3
	assertIsSortLimitSpecBool(t, !includePropertiesInAqlForSpec(&artifactoryParams), false)

	artifactoryParams.Limit = 3
	assertIsSortLimitSpecBool(t, !includePropertiesInAqlForSpec(&artifactoryParams), true)

	artifactoryParams.SortBy = []string{"Vava", "Bubu"}
	assertIsSortLimitSpecBool(t, !includePropertiesInAqlForSpec(&artifactoryParams), true)
}

func assertIsSortLimitSpecBool(t *testing.T, actual, expected bool) {
	if actual != expected {
		t.Error("The function includePropertiesInAqlForSpec() expected to return " + strconv.FormatBool(expected) + " but returned " + strconv.FormatBool(actual) + ".")
	}
}

func TestGetQueryReturnFields(t *testing.T) {
	artifactoryParams := CommonParams{}
	minimalFields := []string{"name", "repo", "path", "actual_md5", "actual_sha1", "sha256", "size", "type", "created", "modified"}

	assertEqualFieldsList(t, getQueryReturnFields(&artifactoryParams, ALL), append(minimalFields, "property"))
	assertEqualFieldsList(t, getQueryReturnFields(&artifactoryParams, SYMLINK), append(minimalFields, "property"))
	assertEqualFieldsList(t, getQueryReturnFields(&artifactoryParams, NONE), minimalFields)

	artifactoryParams.SortBy = []string{"Vava"}
	assertEqualFieldsList(t, getQueryReturnFields(&artifactoryParams, NONE), append(minimalFields, "Vava"))
	assertEqualFieldsList(t, getQueryReturnFields(&artifactoryParams, ALL), append(minimalFields, "Vava"))
	assertEqualFieldsList(t, getQueryReturnFields(&artifactoryParams, SYMLINK), append(minimalFields, "Vava"))

	artifactoryParams.SortBy = []string{"Vava", "Bubu"}
	assertEqualFieldsList(t, getQueryReturnFields(&artifactoryParams, ALL), append(minimalFields, "Vava", "Bubu"))
}

func assertEqualFieldsList(t *testing.T, actual, expected []string) {
	if len(actual) != len(expected) {
		t.Error("The function getQueryReturnFields() expected to return the array:\n" + strings.Join(expected, ",") + ".\nbut returned:\n" + strings.Join(actual, ",") + ".")
	}
	for _, v := range actual {
		isFound := false
		for _, t := range expected {
			if v == t {
				isFound = true
				break
			}
		}
		if !isFound {
			t.Error("The function getQueryReturnFields() expected to return the array:\n'" + strings.Join(expected, ",") + "'.\nbut returned:\n'" + strings.Join(actual, ",") + "'.\n" +
				"The field " + v + "is missing!")
		}
	}
}

func TestBuildSortBody(t *testing.T) {
	assertSortBody(t, buildSortQueryPart([]string{"bubu"}, ""), `"$asc":["bubu"]`)
	assertSortBody(t, buildSortQueryPart([]string{"bubu", "kuku"}, ""), `"$asc":["bubu","kuku"]`)
}

func assertSortBody(t *testing.T, actual, expected string) {
	if actual != expected {
		t.Error("The function buildSortQueryPart expected to return the string:\n'" + expected + "'.\nbut returned:\n'" + actual + "'.")
	}
}

func TestCreateAqlQueryForLatestCreated(t *testing.T) {
	actual := CreateAqlQueryForLatestCreated("repo", "name")
	expected := `items.find({` +
		`"type": "` + string(File) + `",` +
		`"repo": "repo",` +
		`"path": {"$match": "name"}` +
		`})` +
		`.sort({"$desc":["created"]})` +
		`.limit(1)`
	if actual != expected {
		t.Error("The function CreateAqlQueryForLatestCreated expected to return the string:\n'" + expected + "'.\nbut returned:\n'" + actual + "'.")
	}
}

func TestCreateAqlQueryForLatestCreatedFolder(t *testing.T) {
	actual := CreateAqlQueryForLatestCreatedFolder("repo", "name")
	expected := `items.find({` +
		`"type": "` + string(Folder) + `",` +
		`"repo": "repo",` +
		`"path": {"$match": "name"}` +
		`})` +
		`.sort({"$desc":["created"]})` +
		`.limit(1)`
	if actual != expected {
		t.Error("The function CreateAqlQueryForLatestCreated expected to return the string:\n'" + expected + "'.\nbut returned:\n'" + actual + "'.")
	}
}

func TestPrepareSourceSearchPattern(t *testing.T) {
	newPattern := prepareSourceSearchPattern("/testdata/b/b1/b.in", "/testdata")
	assert.Equal(t, "/testdata/b/b1/b.in", newPattern)

	newPattern = prepareSourceSearchPattern("/testdata/b/b1(b).in", "/testdata")
	assert.Equal(t, "/testdata/b/b1(b).in", newPattern)

	newPattern = prepareSourceSearchPattern("/testdata/b/b1(b.in", "/testdata")
	assert.Equal(t, "/testdata/b/b1(b.in", newPattern)

	newPattern = prepareSourceSearchPattern("/testdata/b/b1/)b.in", "/testdata")
	assert.Equal(t, "/testdata/b/b1/)b.in", newPattern)

	newPattern = prepareSourceSearchPattern("/testdata/b/b1/(*).in", "/testdata/{1}.zip")
	assert.Equal(t, "/testdata/b/b1/*.in", newPattern)

	newPattern = prepareSourceSearchPattern("/testdata/b/b1/(*)", "/testdata/{1}")
	assert.Equal(t, "/testdata/b/b1/*", newPattern)

	newPattern = prepareSourceSearchPattern("/testdata/b/(b1)/(*).in", "/testdata/{2}.zip")
	assert.Equal(t, "/testdata/b/(b1)/*.in", newPattern)

	newPattern = prepareSourceSearchPattern("/testdata/(b/(b1)/(*).in", "/testdata/{2}.zip")
	assert.Equal(t, "/testdata/(b/(b1)/*.in", newPattern)

	newPattern = prepareSourceSearchPattern("/testdata/)b/(b1)/(*).in", "/testdata/{2}.zip")
	assert.Equal(t, "/testdata/)b/(b1)/*.in", newPattern)

	newPattern = prepareSourceSearchPattern("/testdata/)b(/(b1)/(*).in", "/testdata/{2}.zip")
	assert.Equal(t, "/testdata/)b(/(b1)/*.in", newPattern)

	newPattern = prepareSourceSearchPattern("/testdata/)b(/(b1)/(*).in", "/testdata/{1}/{2}.zip")
	assert.Equal(t, "/testdata/)b(/b1/*.in", newPattern)

	newPattern = prepareSourceSearchPattern("/testdata/)b(/(b1)/(*).in", "/testdata/{1}/{1}/{2}.zip")
	assert.Equal(t, "/testdata/)b(/b1/*.in", newPattern)

	newPattern = prepareSourceSearchPattern("/testdata/)b(/(b1)/(*).(in)", "/testdata/{1}/{1}/{3}/{2}.zip")
	assert.Equal(t, "/testdata/)b(/b1/*.in", newPattern)

	newPattern = prepareSourceSearchPattern("/testdata/b/(/(.in", "/testdata")
	assert.Equal(t, "/testdata/b/(/(.in", newPattern)
}

var aqlQueryForBuildDataProvider = []struct {
	artifactsQuery bool
	builds         []Build
	expected       string
}{
	{true, []Build{{"buildName", "buildNumber"}},
		`{"$and":[{"artifact.module.build.name":"buildName","artifact.module.build.number":"buildNumber"}]}`},
	{true, []Build{{"buildName1", "buildNumber1"}, {"buildName2", "buildNumber2"}},
		`{"$and":[{"artifact.module.build.name":"buildName1","artifact.module.build.number":"buildNumber1"}]},{"$and":[{"artifact.module.build.name":"buildName2","artifact.module.build.number":"buildNumber2"}]}`},
	{false, []Build{{"buildName", "buildNumber"}},
		`{"$and":[{"dependency.module.build.name":"buildName","dependency.module.build.number":"buildNumber"}]}`},
	{false, []Build{{"buildName1", "buildNumber1"}, {"buildName2", "buildNumber2"}},
		`{"$and":[{"dependency.module.build.name":"buildName1","dependency.module.build.number":"buildNumber1"}]},{"$and":[{"dependency.module.build.name":"buildName2","dependency.module.build.number":"buildNumber2"}]}`},
}

func TestCreateAqlQueryForBuild(t *testing.T) {
	for _, sample := range aqlQueryForBuildDataProvider {
		t.Run(fmt.Sprintf("%v, artifacts: %v", sample.builds, sample.artifactsQuery), func(t *testing.T) {
			expected := `items.find({"$or":[` + sample.expected + "]})"
			actual := createAqlQueryForBuild("", sample.artifactsQuery, sample.builds)
			assert.Equal(t, expected, actual)
		})
	}
}

var ketValuePartsProvider = []struct {
	key      string
	values   []string
	expected string
}{
	{"key1", []string{"value1"}, `{"@key1":"value1"}`},
	{"key1", []string{"value1", "value2"}, `{"@key1":"value1"},{"@key1":"value2"}`},
	{"key1", []string{"value1", "value2", "value3"}, `{"@key1":"value1"},{"@key1":"value2"},{"@key1":"value3"}`},
}

func TestBuildKeyValQueryPart(t *testing.T) {
	for _, sample := range ketValuePartsProvider {
		t.Run(sample.expected, func(t *testing.T) {
			expected := `"$or":[` + sample.expected + "]"
			actual := buildKeyValQueryPart(sample.key, sample.values)
			assert.Equal(t, expected, actual)
		})
	}
}

var encodeForBuildInfoRepositoryProvider = []struct {
	value            string
	expectedEncoding string
}{
	// Shouldn't encode
	{"", ""},
	{"a", "a"},
	{"a b", "a b"},
	{"a.b", "a.b"},
	{"a&b", "a&b"},

	// Should encode
	{"a/b", "a%2Fb"},
	{"a\\b", "a%5Cb"},
	{"a:b", "a%3Ab"},
	{"a|b", "a%7Cb"},
	{"a*b", "a%2Ab"},
	{"a?b", "a%3Fb"},
	{"a  /  b", "a %20%2F%20 b"},

	// Should convert whitespace to space
	{"a\tb", "a b"},
	{"a\nb", "a b"},
}

func TestEncodeForBuildInfoRepository(t *testing.T) {
	for _, testCase := range encodeForBuildInfoRepositoryProvider {
		t.Run(testCase.value, func(t *testing.T) {
			assert.Equal(t, testCase.expectedEncoding, encodeForBuildInfoRepository(testCase.value))
		})
	}
}

func TestBuildAqlSearchQueryWithExclusions(t *testing.T) {
	params := CommonParams{
		Pattern:    "repo-local/*",
		Recursive:  true,
		Exclusions: []string{"*.json"},
	}
	aqlResult, err := CreateAqlBodyForSpecWithPattern(&params)
	assert.NoError(t, err)
	assert.Contains(t, aqlResult, `"$nmatch"`)
	assert.Contains(t, aqlResult, `*.json`)
}

func TestCreateAqlBodyForBuildArtifactsWithExclusions(t *testing.T) {
	builds := []Build{{"myBuild", "123"}}

	aqlBodyNoExclusions := createAqlBodyForBuildArtifacts(builds)
	assert.NotContains(t, aqlBodyNoExclusions, `"$nmatch"`)
	assert.Contains(t, aqlBodyNoExclusions, `artifact.module.build.name`)
	assert.Contains(t, aqlBodyNoExclusions, `myBuild`)

	params := &CommonParams{
		Exclusions: []string{"*.json"},
		Recursive:  true,
	}
	aqlBodyWithExclusions := createAqlBodyForBuildArtifactsWithExclusions(builds, params)

	assert.Contains(t, aqlBodyWithExclusions, `"$nmatch"`, "FIX VERIFIED: createAqlBodyForBuildArtifactsWithExclusions includes exclusions")
	assert.Contains(t, aqlBodyWithExclusions, `*.json`)
	assert.Contains(t, aqlBodyWithExclusions, `artifact.module.build.name`)
	assert.Contains(t, aqlBodyWithExclusions, `artifact.module.build.number`)
	assert.Contains(t, aqlBodyWithExclusions, `myBuild`)
	assert.Contains(t, aqlBodyWithExclusions, `123`)
}

func TestCreateAqlBodyForBuildDependenciesWithExclusions(t *testing.T) {
	builds := []Build{{"myBuild", "123"}}

	// Test with exclusions
	params := &CommonParams{
		Exclusions: []string{"*.xml", "test-*"},
		Recursive:  true,
	}
	aqlBody := createAqlBodyForBuildDependenciesWithExclusions(builds, params)

	assert.Contains(t, aqlBody, `"$nmatch"`)
	assert.Contains(t, aqlBody, `*.xml`)
	assert.Contains(t, aqlBody, `test-*`)
	assert.Contains(t, aqlBody, `dependency.module.build.name`)
	assert.Contains(t, aqlBody, `dependency.module.build.number`)
}

func TestGetSpecType_BuildWithPattern(t *testing.T) {
	tests := []struct {
		name     string
		params   CommonParams
		expected SpecType
	}{
		{
			name:     "build only routes to BUILD",
			params:   CommonParams{Build: "my-build/1"},
			expected: BUILD,
		},
		{
			name:     "build with empty pattern routes to BUILD",
			params:   CommonParams{Build: "my-build/1", Pattern: ""},
			expected: BUILD,
		},
		{
			name:     "build with wildcard-all pattern routes to BUILD",
			params:   CommonParams{Build: "my-build/1", Pattern: "*"},
			expected: BUILD,
		},
		{
			name:     "build with specific pattern routes to BUILD",
			params:   CommonParams{Build: "my-build/1", Pattern: "docker-local-ash/*"},
			expected: BUILD,
		},
		{
			name:     "build with repo path pattern routes to BUILD",
			params:   CommonParams{Build: "my-build/1", Pattern: "repo-local/path/to/*.jar"},
			expected: BUILD,
		},
		{
			name:     "pattern only routes to WILDCARD",
			params:   CommonParams{Pattern: "repo-local/*"},
			expected: WILDCARD,
		},
		{
			name:     "aql routes to AQL",
			params:   CommonParams{Aql: Aql{ItemsFind: `{"repo":"test"}`}},
			expected: AQL,
		},
		{
			name:     "aql with build routes to AQL (AQL takes precedence)",
			params:   CommonParams{Aql: Aql{ItemsFind: `{"repo":"test"}`}, Build: "my-build/1"},
			expected: AQL,
		},
		{
			name:     "build with IncludeDeps and pattern routes to WILDCARD",
			params:   CommonParams{Build: "my-build/1", Pattern: "repo-local/*", IncludeDeps: true},
			expected: WILDCARD,
		},
		{
			name:     "build with IncludeDeps and no pattern routes to BUILD",
			params:   CommonParams{Build: "my-build/1", IncludeDeps: true},
			expected: BUILD,
		},
		{
			name:     "build with IncludeDeps and wildcard pattern routes to BUILD",
			params:   CommonParams{Build: "my-build/1", Pattern: "*", IncludeDeps: true},
			expected: BUILD,
		},
		{
			name:     "build with ExcludeArtifacts but no IncludeDeps routes to BUILD",
			params:   CommonParams{Build: "my-build/1", ExcludeArtifacts: true},
			expected: BUILD,
		},
		{
			name:     "empty params routes to WILDCARD",
			params:   CommonParams{},
			expected: WILDCARD,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.params.GetSpecType()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAqlMatch(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		pattern string
		match   bool
	}{
		{"exact match", "repo-local", "repo-local", true},
		{"exact mismatch", "repo-local", "other-repo", false},
		{"wildcard all", "anything", "*", true},
		{"wildcard suffix", "file.jar", "*.jar", true},
		{"wildcard suffix mismatch", "file.txt", "*.jar", false},
		{"wildcard prefix", "test-file", "test-*", true},
		{"wildcard prefix mismatch", "prod-file", "test-*", false},
		{"wildcard middle", "test-123-file", "test-*-file", true},
		{"wildcard middle mismatch", "test-123-other", "test-*-file", false},
		{"multiple wildcards", "a/b/c/d", "a/*/c/*", true},
		{"dot path", ".", "*", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := aqlMatch(tt.value, tt.pattern)
			assert.Equal(t, tt.match, result)
		})
	}
}

func TestMatchResultItemToTriples(t *testing.T) {
	tests := []struct {
		name    string
		item    ResultItem
		triples []RepoPathFile
		match   bool
	}{
		{
			name: "item matches single triple",
			item: ResultItem{Repo: "docker-local-ash", Path: "path/to", Name: "file.jar"},
			triples: []RepoPathFile{
				{repo: "docker-local-ash", path: "*", file: "*"},
			},
			match: true,
		},
		{
			name: "item does not match repo",
			item: ResultItem{Repo: "other-repo", Path: "path/to", Name: "file.jar"},
			triples: []RepoPathFile{
				{repo: "docker-local-ash", path: "*", file: "*"},
			},
			match: false,
		},
		{
			name: "item matches one of multiple triples",
			item: ResultItem{Repo: "repo-b", Path: "some/path", Name: "test.txt"},
			triples: []RepoPathFile{
				{repo: "repo-a", path: "*", file: "*"},
				{repo: "repo-b", path: "*", file: "*.txt"},
			},
			match: true,
		},
		{
			name:    "empty triples matches nothing",
			item:    ResultItem{Repo: "repo-a", Path: ".", Name: "file.jar"},
			triples: []RepoPathFile{},
			match:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matchResultItemToTriples(&tt.item, tt.triples)
			assert.Equal(t, tt.match, result)
		})
	}
}
