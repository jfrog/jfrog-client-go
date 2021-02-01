package utils

import (
	"testing"

	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/stretchr/testify/assert"
)

func TestCreateBundleBody(t *testing.T) {
	releaseBundleParam := ReleaseBundleParams{
		SignImmediately:    true,
		StoringRepository:  "storing-repo",
		Description:        "Release bundle description",
		ReleaseNotes:       "Release notes",
		ReleaseNotesSyntax: Asciidoc,
	}

	releaseBundleBody, err := CreateBundleBody(releaseBundleParam, true)
	assert.NoError(t, err)
	assert.NotNil(t, releaseBundleBody)
	assert.Equal(t, true, releaseBundleBody.DryRun)
	assert.Equal(t, true, releaseBundleBody.SignImmediately)
	assert.Equal(t, "storing-repo", releaseBundleBody.StoringRepository)
	assert.Equal(t, "Release bundle description", releaseBundleBody.Description)
	assert.Equal(t, "Release notes", releaseBundleBody.ReleaseNotes.Content)
	assert.Equal(t, ReleaseNotesSyntax(Asciidoc), releaseBundleBody.ReleaseNotes.Syntax)
	assert.Len(t, releaseBundleBody.BundleSpec.Queries, 0)
}

func TestCreateBundleBodyQuery(t *testing.T) {
	releaseBundleParam := ReleaseBundleParams{
		SpecFiles: []*utils.ArtifactoryCommonParams{{Pattern: "dist-repo/*", TargetProps: "a=b;c=d;c=e"}},
	}

	releaseBundleBody, err := CreateBundleBody(releaseBundleParam, true)
	assert.NoError(t, err)
	assert.NotNil(t, releaseBundleBody)
	assert.Len(t, releaseBundleBody.BundleSpec.Queries, 1)
	query := releaseBundleBody.BundleSpec.Queries[0]
	assert.Contains(t, query.Aql, "dist-repo")
	props := query.AddedProps
	assert.Len(t, props, 2)
	for _, prop := range props {
		switch prop.Key {
		case "a":
			assert.Equal(t, []string{"b"}, prop.Values)
		case "c":
			assert.ElementsMatch(t, []string{"d", "e"}, prop.Values)
		default:
			assert.Fail(t, "Unexpected key "+prop.Key)
		}
	}
}

func TestCreatePathMappings(t *testing.T) {
	tests := []struct {
		specPattern           string
		specTarget            string
		expectedMappingInput  string
		expectedMappingOutput string
	}{
		{"", "", "", ""},
		{"repo/path/file.in", "", "", ""},
		{"a/b/c", "a/b/x", "^a/b/c$", "a/b/x"},
		{"a/(b)/c", "a/d/c", "^a/(b)/c$", "a/d/c"},
		{"a/(*)/c", "a/d/c", "^a/(.*)/c$", "a/d/c"},
		{"a/(b)/c", "a/(d)/c", "^a/(b)/c$", "a/(d)/c"},
		{"a/(b)/c", "a/b/c/{1}", "^a/(b)/c$", "a/b/c/$1"},
		{"a/(b)/(c)", "a/b/c/{1}/{2}", "^a/(b)/(c)$", "a/b/c/$1/$2"},
		{"a/(b)/(c)", "a/b/c/{2}/{1}", "^a/(b)/(c)$", "a/b/c/$2/$1"},
	}

	for _, test := range tests {
		t.Run(test.specPattern, func(t *testing.T) {
			specFile := &utils.ArtifactoryCommonParams{Pattern: test.specPattern, Target: test.specTarget}
			pathMappings := createPathMappings(specFile)
			if test.expectedMappingInput == "" {
				assert.Empty(t, pathMappings)
				return
			}
			assert.Len(t, pathMappings, 1)
			actualPathMapping := pathMappings[0]
			assert.Equal(t, test.expectedMappingInput, actualPathMapping.Input)
			assert.Equal(t, test.expectedMappingOutput, actualPathMapping.Output)
		})
	}
}
