package utils

import (
	"encoding/json"
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
	aqlBodyWithExclusions, err := createAqlBodyForBuildArtifactsWithExclusions(builds, params)
	assert.NoError(t, err)

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
	aqlBody, err := createAqlBodyForBuildDependenciesWithExclusions(builds, params)
	assert.NoError(t, err)

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
			name:     "build with pattern and props routes to WILDCARD",
			params:   CommonParams{Build: "my-build/1", Pattern: "repo-local/*", Props: "key=value"},
			expected: WILDCARD,
		},
		{
			name:     "build with pattern and excludeProps routes to WILDCARD",
			params:   CommonParams{Build: "my-build/1", Pattern: "repo-local/*", ExcludeProps: "key=value"},
			expected: WILDCARD,
		},
		{
			name:     "build with props but no pattern routes to BUILD",
			params:   CommonParams{Build: "my-build/1", Props: "key=value"},
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

func TestCreateAqlBodyForBuildArtifactsWithPattern(t *testing.T) {
	builds := []Build{{BuildName: "my-build", BuildNumber: "1"}}

	t.Run("no pattern: build-only filter, no $and wrapping", func(t *testing.T) {
		body, err := createAqlBodyForBuildArtifactsWithExclusions(builds, &CommonParams{})
		assert.NoError(t, err)
		assert.Contains(t, body, `"artifact.module.build.name":"my-build"`)
		assert.NotContains(t, body, `"$and":[{"$or"`)
	})

	t.Run("trivial '*' pattern: treated as no pattern", func(t *testing.T) {
		body, err := createAqlBodyForBuildArtifactsWithExclusions(builds, &CommonParams{Pattern: "*"})
		assert.NoError(t, err)
		assert.NotContains(t, body, `"$and":[{"$or"`)
	})

	t.Run("pattern with virtual repo: server-side match via AQL", func(t *testing.T) {
		// A client-side string compare on `Repo` cannot translate virtual → backing local,
		// but AQL can. The pattern must go into the query itself.
		body, err := createAqlBodyForBuildArtifactsWithExclusions(builds, &CommonParams{
			Pattern:   "some-virtual-repo/com/jfrog/*/my-artifact-*.tgz",
			Recursive: true,
		})
		assert.NoError(t, err)
		assert.Contains(t, body, `"artifact.module.build.name":"my-build"`)
		assert.Contains(t, body, `"repo":"some-virtual-repo"`)
		assert.Contains(t, body, `"$and":[{"$or"`, "build and pattern filters should be ANDed together")
	})

	t.Run("invalid pattern surfaces as error", func(t *testing.T) {
		// Pattern starting with '/' is invalid — must start with a repo name or asterisk.
		_, err := createAqlBodyForBuildArtifactsWithExclusions(builds, &CommonParams{
			Pattern:   "/leading-slash-is-invalid",
			Recursive: true,
		})
		assert.Error(t, err, "invalid patterns must not be silently swallowed")
	})

	t.Run("transitive + multi-repo pattern is rejected", func(t *testing.T) {
		// Wildcards before the first slash expand to multiple repos, which transitive mode forbids.
		_, err := createAqlBodyForBuildArtifactsWithExclusions(builds, &CommonParams{
			Pattern:    "repo-*/com/jfrog/*.tgz",
			Recursive:  true,
			Transitive: true,
		})
		assert.Error(t, err, "transitive search with multi-repo wildcards must be rejected")
		assert.Contains(t, err.Error(), "transitive", "error should mention transitive")
	})

	t.Run("transitive + single-repo pattern is accepted", func(t *testing.T) {
		body, err := createAqlBodyForBuildArtifactsWithExclusions(builds, &CommonParams{
			Pattern:    "one-repo/com/jfrog/*.tgz",
			Recursive:  true,
			Transitive: true,
		})
		assert.NoError(t, err)
		assert.Contains(t, body, `"repo":"one-repo"`)
	})

	t.Run("exclusions + pattern: both appear in AQL", func(t *testing.T) {
		body, err := createAqlBodyForBuildArtifactsWithExclusions(builds, &CommonParams{
			Pattern:    "some-repo/com/jfrog/*.jar",
			Recursive:  true,
			Exclusions: []string{"*test*.jar"},
		})
		assert.NoError(t, err)
		assert.Contains(t, body, `"$and":[{"$or"`, "pattern is still present under $and")
		assert.Contains(t, body, `"$nmatch"`, "exclusion should produce $nmatch filter")
		assert.Contains(t, body, `"repo":"some-repo"`)
	})

	t.Run("multi-triple pattern expands into $or", func(t *testing.T) {
		// A recursive pattern with a wildcard mid-path produces multiple triples; each
		// triple carries its own "repo" literal, so the count signals fan-out.
		body, err := createAqlBodyForBuildArtifactsWithExclusions(builds, &CommonParams{
			Pattern:   "some-repo/com/jfrog/*/files/*-linux-amd64.tgz",
			Recursive: true,
		})
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, strings.Count(body, `"repo":"some-repo"`), 2, "recursive pattern should produce multiple triples")
	})

	t.Run("aggregated builds + pattern: all builds ANDed with pattern", func(t *testing.T) {
		aggregated := []Build{
			{BuildName: "my-build", BuildNumber: "1"},
			{BuildName: "my-build", BuildNumber: "2"},
			{BuildName: "upstream-build", BuildNumber: "5"},
		}
		body, err := createAqlBodyForBuildArtifactsWithExclusions(aggregated, &CommonParams{
			Pattern:   "some-repo/com/jfrog/*.jar",
			Recursive: true,
		})
		assert.NoError(t, err)
		assert.Contains(t, body, `"artifact.module.build.name":"my-build","artifact.module.build.number":"1"`)
		assert.Contains(t, body, `"artifact.module.build.name":"my-build","artifact.module.build.number":"2"`)
		assert.Contains(t, body, `"artifact.module.build.name":"upstream-build","artifact.module.build.number":"5"`)
		assert.Contains(t, body, `"$and":[{"$or"`, "aggregated build list must be ANDed with pattern, not replaced")
	})

	t.Run("'*/path' pattern: repo wildcard not emitted as repo filter", func(t *testing.T) {
		// buildInnerQueryPart skips the repo condition when the triple's repo is '*' or '**',
		// so users writing '*/path/...' as a workaround don't get a literal `"repo":"*"` emitted.
		body, err := createAqlBodyForBuildArtifactsWithExclusions(builds, &CommonParams{
			Pattern:   "*/com/jfrog/foo/*.tgz",
			Recursive: true,
		})
		assert.NoError(t, err)
		assert.Contains(t, body, `"$and":[{"$or"`, "pattern should still be present")
		assert.NotContains(t, body, `"repo":"*"`, "repo='*' must not be emitted as an AQL filter")
	})

	t.Run("nil-params wrapper: returns valid JSON with no pattern or exclusions", func(t *testing.T) {
		// The wrapper at createAqlBodyForBuildArtifacts calls the underlying function with
		// nil params and swallows the error; verify its output is a well-formed JSON body
		// containing only the build filter.
		body := createAqlBodyForBuildArtifacts(builds)
		var parsed map[string]any
		assert.NoError(t, json.Unmarshal([]byte(body), &parsed), "body must be valid JSON")
		assert.Contains(t, body, `"artifact.module.build.name":"my-build"`)
		assert.NotContains(t, body, `"$and":[{"$or"`, "no pattern → no $and wrapping")
		assert.NotContains(t, body, `"$nmatch"`, "no exclusions → no $nmatch")
	})

	t.Run("malformed exclusion propagates as error", func(t *testing.T) {
		// Guards the behavior change introduced by this fix: a malformed exclusion pattern
		// used to silently produce a build-only AQL (exclusion dropped); now it aborts.
		_, err := createAqlBodyForBuildArtifactsWithExclusions(builds, &CommonParams{
			Pattern:    "some-repo/*",
			Recursive:  true,
			Exclusions: []string{"/leading-slash-is-invalid"},
		})
		assert.Error(t, err, "malformed exclusion must propagate, not be silently swallowed")
	})
}
