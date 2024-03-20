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
	assert.Equal(t, true, *releaseBundleBody.SignImmediately)
	assert.Equal(t, "storing-repo", releaseBundleBody.StoringRepository)
	assert.Equal(t, "Release bundle description", releaseBundleBody.Description)
	assert.Equal(t, "Release notes", releaseBundleBody.ReleaseNotes.Content)
	assert.Equal(t, Asciidoc, releaseBundleBody.ReleaseNotes.Syntax)
	assert.Len(t, releaseBundleBody.BundleSpec.Queries, 0)
}

func TestCreateBundleBodyQuery(t *testing.T) {
	targetProps, err := utils.ParseProperties("a=b;c=d;c=e")
	assert.NoError(t, err)

	releaseBundleParam := ReleaseBundleParams{
		SpecFiles: []*utils.CommonParams{{Pattern: "dist-repo/*", TargetProps: targetProps}},
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
