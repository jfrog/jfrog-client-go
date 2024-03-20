package distribution

import (
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreatePathMappingsFromPatternAndTarget(t *testing.T) {
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
			specFile := &utils.CommonParams{Pattern: test.specPattern, Target: test.specTarget}
			pathMappings := CreatePathMappingsFromPatternAndTarget(specFile.Pattern, specFile.Target)
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

func TestCreatePathMappings(t *testing.T) {
	tests := []struct {
		specInput             string
		specOutput            string
		expectedMappingInput  string
		expectedMappingOutput string
	}{
		{"", "", "", ""},
		{"repo/path/file.in", "", "", ""},
		{"", "repo/path/file.out", "", ""},
		{"a/b/c", "a/b/x", "a/b/c", "a/b/x"},
	}

	for _, test := range tests {
		t.Run(test.specInput, func(t *testing.T) {
			specFile := &utils.CommonParams{
				PathMapping: utils.PathMapping{
					Input:  test.specInput,
					Output: test.specOutput,
				},
			}
			pathMappings := CreatePathMappings(specFile.PathMapping.Input, specFile.PathMapping.Output)
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
