package utils

import (
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
		`{"path":{"$ne":"."},"$or":[{"$and":[{"repo":{"$match":"repo-local*"},"path":{"$match":"a*b*c/dd"},"name":{"$match":"*"}}]},{"$and":[{"repo":{"$match":"repo-local*"},"path":{"$match":"a*b*c/dd/*"},"name":{"$match":"*"}}]}]}`},
	{"repo-local", false,
		`{"$or":[{"$and":[{"repo":"repo-local","path":".","name":{"$match":"*"}}]}]}`},
	{"*repo-local", false,
		`{"$or":[{"$and":[{"repo":{"$match":"*"},"path":".","name":{"$match":"*repo-local"}}]},{"$and":[{"repo":{"$match":"*repo-local"},"path":".","name":{"$match":"*"}}]}]}`},
	{"repo-local2/a*b*c/dd/", false,
		`{"path":{"$ne":"."},"$or":[{"$and":[{"repo":"repo-local2","path":{"$match":"a*b*c/dd"},"name":{"$match":"*"}}]}]}`},
	{"*/a*b*c/dd/", false,
		`{"path":{"$ne":"."},"$or":[{"$and":[{"repo":{"$match":"*"},"path":{"$match":"a*b*c/dd"},"name":{"$match":"*"}}]}]}`},
}

func TestBuildAqlSearchQuery(t *testing.T) {
	for _, sample := range buildAqlSearchQueryDataProvider {
		t.Run(sample.pattern+"_recursive_"+strconv.FormatBool(sample.recursive), func(t *testing.T) {
			params := ArtifactoryCommonParams{Pattern: sample.pattern, Recursive: sample.recursive, Regexp: false, IncludeDirs: false}
			aqlResult, err := CreateAqlBodyForSpecWithPattern(&params)
			assert.NoError(t, err)
			if aqlResult != sample.expectedAql {
				t.Error("Unexpected download AQL query built. \nExpected: " + sample.expectedAql + " \nGot:      " + aqlResult)
			}
		})
	}
}

func TestArtifactoryCommonParams(t *testing.T) {
	artifactoryParams := ArtifactoryCommonParams{}
	assertIsSortLimitSpecBool(!includePropertiesInAqlForSpec(&artifactoryParams), false, t)

	artifactoryParams.SortBy = []string{"Vava", "Bubu"}
	assertIsSortLimitSpecBool(!includePropertiesInAqlForSpec(&artifactoryParams), true, t)

	artifactoryParams.SortBy = nil
	artifactoryParams.Limit = 0
	assertIsSortLimitSpecBool(!includePropertiesInAqlForSpec(&artifactoryParams), false, t)

	artifactoryParams.Limit = -3
	assertIsSortLimitSpecBool(!includePropertiesInAqlForSpec(&artifactoryParams), false, t)

	artifactoryParams.Limit = 3
	assertIsSortLimitSpecBool(!includePropertiesInAqlForSpec(&artifactoryParams), true, t)

	artifactoryParams.SortBy = []string{"Vava", "Bubu"}
	assertIsSortLimitSpecBool(!includePropertiesInAqlForSpec(&artifactoryParams), true, t)
}

func assertIsSortLimitSpecBool(actual, expected bool, t *testing.T) {
	if actual != expected {
		t.Error("The function includePropertiesInAqlForSpec() expected to return " + strconv.FormatBool(expected) + " but returned " + strconv.FormatBool(actual) + ".")
	}
}

func TestGetQueryReturnFields(t *testing.T) {
	artifactoryParams := ArtifactoryCommonParams{}
	minimalFields := []string{"name", "repo", "path", "actual_md5", "actual_sha1", "size", "type", "created", "modified"}

	assertEqualFieldsList(getQueryReturnFields(&artifactoryParams, ALL), append(minimalFields, "property"), t)
	assertEqualFieldsList(getQueryReturnFields(&artifactoryParams, SYMLINK), append(minimalFields, "property"), t)
	assertEqualFieldsList(getQueryReturnFields(&artifactoryParams, NONE), append(minimalFields), t)

	artifactoryParams.SortBy = []string{"Vava"}
	assertEqualFieldsList(getQueryReturnFields(&artifactoryParams, NONE), append(minimalFields, "Vava"), t)
	assertEqualFieldsList(getQueryReturnFields(&artifactoryParams, ALL), append(minimalFields, "Vava"), t)
	assertEqualFieldsList(getQueryReturnFields(&artifactoryParams, SYMLINK), append(minimalFields, "Vava"), t)

	artifactoryParams.SortBy = []string{"Vava", "Bubu"}
	assertEqualFieldsList(getQueryReturnFields(&artifactoryParams, ALL), append(minimalFields, "Vava", "Bubu"), t)
}

func assertEqualFieldsList(actual, expected []string, t *testing.T) {
	if len(actual) != len(expected) {
		t.Error("The function getQueryReturnFields() expected to return the array:\n" + strings.Join(expected[:], ",") + ".\nbut returned:\n" + strings.Join(actual[:], ",") + ".")
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
			t.Error("The function getQueryReturnFields() expected to return the array:\n'" + strings.Join(expected[:], ",") + "'.\nbut returned:\n'" + strings.Join(actual[:], ",") + "'.\n" +
				"The field " + v + "is missing!")
		}
	}
}

func TestBuildSortBody(t *testing.T) {
	assertSortBody(buildSortQueryPart([]string{"bubu"}, ""), `"$asc":["bubu"]`, t)
	assertSortBody(buildSortQueryPart([]string{"bubu", "kuku"}, ""), `"$asc":["bubu","kuku"]`, t)
}

func assertSortBody(actual, expected string, t *testing.T) {
	if actual != expected {
		t.Error("The function buildSortQueryPart expected to return the string:\n'" + expected + "'.\nbut returned:\n'" + actual + "'.")
	}
}

func TestPrepareSourceSearchPattern(t *testing.T) {
	newPattern := prepareSourceSearchPattern("/testsdata/b/b1/b.in", "/testsdata", true)
	assert.Equal(t, "/testsdata/b/b1/b.in", newPattern)

	newPattern = prepareSourceSearchPattern("/testsdata/b/b1(b).in", "/testsdata", true)
	assert.Equal(t, "/testsdata/b/b1(b).in", newPattern)

	newPattern = prepareSourceSearchPattern("/testsdata/b/b1(b.in", "/testsdata", true)
	assert.Equal(t, "/testsdata/b/b1(b.in", newPattern)

	newPattern = prepareSourceSearchPattern("/testsdata/b/b1/)b.in", "/testsdata", true)
	assert.Equal(t, "/testsdata/b/b1/)b.in", newPattern)

	newPattern = prepareSourceSearchPattern("/testsdata/b/b1/(*).in", "/testsdata/{1}.zip", true)
	assert.Equal(t, "/testsdata/b/b1/*.in", newPattern)

	newPattern = prepareSourceSearchPattern("/testsdata/b/b1/(*)", "/testsdata/{1}", true)
	assert.Equal(t, "/testsdata/b/b1/*", newPattern)

	newPattern = prepareSourceSearchPattern("/testsdata/b/(b1)/(*).in", "/testsdata/{2}.zip", true)
	assert.Equal(t, "/testsdata/b/(b1)/*.in", newPattern)

	newPattern = prepareSourceSearchPattern("/testsdata/(b/(b1)/(*).in", "/testsdata/{2}.zip", true)
	assert.Equal(t, "/testsdata/(b/(b1)/*.in", newPattern)

	newPattern = prepareSourceSearchPattern("/testsdata/)b/(b1)/(*).in", "/testsdata/{2}.zip", true)
	assert.Equal(t, "/testsdata/)b/(b1)/*.in", newPattern)

	newPattern = prepareSourceSearchPattern("/testsdata/)b(/(b1)/(*).in", "/testsdata/{2}.zip", true)
	assert.Equal(t, "/testsdata/)b(/(b1)/*.in", newPattern)

	newPattern = prepareSourceSearchPattern("/testsdata/)b(/(b1)/(*).in", "/testsdata/{1}/{2}.zip", true)
	assert.Equal(t, "/testsdata/)b(/b1/*.in", newPattern)

	newPattern = prepareSourceSearchPattern("/testsdata/)b(/(b1)/(*).in", "/testsdata/{1}/{1}/{2}.zip", true)
	assert.Equal(t, "/testsdata/)b(/b1/*.in", newPattern)

	newPattern = prepareSourceSearchPattern("/testsdata/)b(/(b1)/(*).(in)", "/testsdata/{1}/{1}/{3}/{2}.zip", true)
	assert.Equal(t, "/testsdata/)b(/b1/*.in", newPattern)

	newPattern = prepareSourceSearchPattern("/testsdata/b/(/(.in", "/testsdata", true)
	assert.Equal(t, "/testsdata/b/(/(.in", newPattern)
}
