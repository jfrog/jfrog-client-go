package _go

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestCreateUrlPath(t *testing.T) {
	tests := []struct {
		name              string
		pathInArtifactory string
		props             string
		goApiUrl          string
		expectedUrl       string
	}{
		{"withBuildProperties", "go-repo/github.com/jfrog/test/@v/v1.1.1.zip", "build.name=a;build.number=1", "http://test.url/api/go/", "http://test.url/api/go/go-repo/github.com/jfrog/test/@v/v1.1.1.zip;build.name=a;build.number=1"},
		{"withoutBuildProperties", "go-repo/github.com/jfrog/test/@v/v1.1.1.zip", "", "http://test.url/api/go/", "http://test.url/api/go/go-repo/github.com/jfrog/test/@v/v1.1.1.zip"},
		{"withoutBuildPropertiesModExtension", "go-repo/github.com/jfrog/test/@v/v1.1.1.mod", "", "http://test.url/api/go/", "http://test.url/api/go/go-repo/github.com/jfrog/test/@v/v1.1.1.mod"},
	}
	for _, test := range tests {
		test := test // Create a local copy of the test variable,fixing Implicit memory aliasing in for loop.
		t.Run(test.name, func(t *testing.T) {
			assert.NoError(t, CreateUrlPath(test.pathInArtifactory, test.props, &test.goApiUrl))
			// The props might have a different order each time, so we split the URLs and check if the lists are equal (ignoring the order)
			assert.ElementsMatch(t, strings.Split(test.goApiUrl, ";"), strings.Split(test.expectedUrl, ";"))
		})
	}
}
