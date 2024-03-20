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
	for i, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.NoError(t, CreateUrlPath(tests[i].pathInArtifactory, tests[i].props, &tests[i].goApiUrl))
			// The props might have a different order each time, so we split the URLs and check if the lists are equal (ignoring the order)
			assert.ElementsMatch(t, strings.Split(tests[i].goApiUrl, ";"), strings.Split(tests[i].expectedUrl, ";"))
		})
	}
}
