package _go

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestCreateUrlPath(t *testing.T) {

	tests := []struct {
		name        string
		extension   string
		moduleId    string
		version     string
		props       string
		url         string
		expectedUrl string
	}{
		{"withBuildProperties", ".zip", "github.com/jfrog/test", "v1.1.1", "build.name=a;build.number=1", "http://test.url/", "http://test.url//github.com/jfrog/test/@v/v1.1.1.zip;build.name=a;build.number=1"},
		{"withoutBuildProperties", ".zip", "github.com/jfrog/test", "v1.1.1", "", "http://test.url/", "http://test.url//github.com/jfrog/test/@v/v1.1.1.zip"},
		{"withoutBuildPropertiesModExtension", ".mod", "github.com/jfrog/test", "v1.1.1", "", "http://test.url/", "http://test.url//github.com/jfrog/test/@v/v1.1.1.mod"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			CreateUrlPath(test.moduleId, test.version, test.props, test.extension, &test.url)
			// The props might have a different order each time, so we split the URLs and check if the lists are equal (ignoring the order)
			assert.ElementsMatch(t, strings.Split(test.url, ";"), strings.Split(test.expectedUrl, ";"))
		})
	}
}
