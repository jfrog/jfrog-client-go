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
		`{"$or":[{"$and":[{"repo":{"$match":"*"},"path":".","name":{"$match":"*repo-local"}}]},{"$and":[{"repo":{"$match":"*repo-local"},"path":".","name":{"$match":"*"}}]}]}`},
	{"repo-local2/a*b*c/dd/", false,
		`{"path":{"$ne":"."},"$or":[{"$and":[{"repo":"repo-local2","path":{"$match":"a*b*c/dd"},"name":{"$match":"*"}}]}]}`},
	{"*/a*b*c/dd/", false,
		`{"path":{"$ne":"."},"$or":[{"$and":[{"repo":{"$match":"*"},"path":{"$match":"*a*b*c/dd"},"name":{"$match":"*"}}]}]}`},
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
