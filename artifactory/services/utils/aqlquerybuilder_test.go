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

// TestCreateAqlBodyForBuildArtifactsWithProperties tests the property-based AQL query generation
// The WithExclusions functions now use property-based queries by default (RTDEV-64748)
func TestCreateAqlBodyForBuildArtifactsWithProperties(t *testing.T) {
	tests := []struct {
		name     string
		builds   []Build
		expected string
	}{
		{
			name: "Single build (no exclusions)",
			builds: []Build{
				{BuildName: "my-build", BuildNumber: "123"},
			},
			expected: `{"$or":[{"$and":[{"@build.name":"my-build","@build.number":"123"}]}]}`,
		},
		{
			name: "Multiple builds (no exclusions)",
			builds: []Build{
				{BuildName: "build-1", BuildNumber: "1"},
				{BuildName: "build-2", BuildNumber: "2"},
			},
			expected: `{"$or":[{"$and":[{"@build.name":"build-1","@build.number":"1"}]},{"$and":[{"@build.name":"build-2","@build.number":"2"}]}]}`,
		},
		{
			name: "Build with special characters",
			builds: []Build{
				{BuildName: "docker-build", BuildNumber: "20220607-2"},
			},
			expected: `{"$or":[{"$and":[{"@build.name":"docker-build","@build.number":"20220607-2"}]}]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use WithExclusions with nil params (no exclusions) - now uses property-based queries
			result := createAqlBodyForBuildArtifactsWithExclusions(tt.builds, nil)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestCreateAqlBodyForBuildDependenciesWithExclusions tests the property-based AQL query for dependencies
func TestCreateAqlBodyForBuildDependenciesWithProperties(t *testing.T) {
	tests := []struct {
		name     string
		builds   []Build
		expected string
	}{
		{
			name: "Single build dependency (no exclusions)",
			builds: []Build{
				{BuildName: "my-build", BuildNumber: "456"},
			},
			expected: `{"$or":[{"$and":[{"@build.name":"my-build","@build.number":"456"}]}]}`,
		},
		{
			name: "Multiple builds dependencies",
			builds: []Build{
				{BuildName: "build-A", BuildNumber: "10"},
				{BuildName: "build-B", BuildNumber: "20"},
				{BuildName: "build-C", BuildNumber: "30"},
			},
			expected: `{"$or":[{"$and":[{"@build.name":"build-A","@build.number":"10"}]},{"$and":[{"@build.name":"build-B","@build.number":"20"}]},{"$and":[{"@build.name":"build-C","@build.number":"30"}]}]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use WithExclusions with nil params (no exclusions) - now uses property-based queries
			result := createAqlBodyForBuildDependenciesWithExclusions(tt.builds, nil)
			assert.Equal(t, tt.expected, result)
		})
	}
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

// Updated to use property-based queries (@build.name) instead of JOINs (artifact.module.build.name) - RTDEV-64748
var aqlQueryForBuildDataProvider = []struct {
	artifactsQuery bool
	builds         []Build
	expected       string
}{
	{true, []Build{{"buildName", "buildNumber"}},
		`{"$and":[{"@build.name":"buildName","@build.number":"buildNumber"}]}`},
	{true, []Build{{"buildName1", "buildNumber1"}, {"buildName2", "buildNumber2"}},
		`{"$and":[{"@build.name":"buildName1","@build.number":"buildNumber1"}]},{"$and":[{"@build.name":"buildName2","@build.number":"buildNumber2"}]}`},
	{false, []Build{{"buildName", "buildNumber"}},
		`{"$and":[{"@build.name":"buildName","@build.number":"buildNumber"}]}`},
	{false, []Build{{"buildName1", "buildNumber1"}, {"buildName2", "buildNumber2"}},
		`{"$and":[{"@build.name":"buildName1","@build.number":"buildNumber1"}]},{"$and":[{"@build.name":"buildName2","@build.number":"buildNumber2"}]}`},
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
	// Uses property-based query now (RTDEV-64748 fix)
	assert.Contains(t, aqlBodyNoExclusions, `@build.name`)
	assert.Contains(t, aqlBodyNoExclusions, `myBuild`)

	params := &CommonParams{
		Exclusions: []string{"*.json"},
		Recursive:  true,
	}
	aqlBodyWithExclusions := createAqlBodyForBuildArtifactsWithExclusions(builds, params)

	assert.Contains(t, aqlBodyWithExclusions, `"$nmatch"`, "FIX VERIFIED: createAqlBodyForBuildArtifactsWithExclusions includes exclusions")
	assert.Contains(t, aqlBodyWithExclusions, `*.json`)
	// Uses property-based query now (RTDEV-64748 fix)
	assert.Contains(t, aqlBodyWithExclusions, `@build.name`)
	assert.Contains(t, aqlBodyWithExclusions, `@build.number`)
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
	// Uses property-based query now (RTDEV-64748 fix)
	assert.Contains(t, aqlBody, `@build.name`)
	assert.Contains(t, aqlBody, `@build.number`)
}
