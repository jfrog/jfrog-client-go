package _go

import (
	"github.com/jfrog/jfrog-client-go/utils"
	"reflect"
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
			if !strings.EqualFold(test.url, test.expectedUrl) {
				t.Error("Expected:", test.expectedUrl, "Got:", test.url)
			}
		})
	}
}

func TestShouldUseHeaders(t *testing.T) {
	tests := []struct {
		artifactoryVersion string
		expectedResult     string
	}{
		{"6.5.0", "*_go.publishWithMatrixParams"},
		{"6.2.0", "*_go.publishWithHeader"},
		{"5.9.0", "*_go.publishWithHeader"},
		{"6.0.0", "*_go.publishWithHeader"},
		{"6.6.0", "*_go.publishWithMatrixParams"},
		{"6.6.1", "*_go.publishZipAndModApi"},
		{utils.Development, "*_go.publishZipAndModApi"},
		{"6.10.2", "*_go.publishZipAndModApi"},
	}
	for _, test := range tests {
		t.Run(test.artifactoryVersion, func(t *testing.T) {
			result := GetCompatiblePublisher(test.artifactoryVersion)
			if reflect.TypeOf(result).String() != test.expectedResult {
				t.Error("Expected:", test.expectedResult, "Got:", reflect.TypeOf(result).String())
			}
		})
	}
}
