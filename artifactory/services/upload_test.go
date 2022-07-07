package services

import (
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDebianProperties(t *testing.T) {
	var debianPaths = []struct {
		in       string
		expected string
	}{
		{"dist/comp/arch", ";deb.distribution=dist;deb.component=comp;deb.architecture=arch"},
		{"dist1,dist2/comp/arch", ";deb.distribution=dist1,dist2;deb.component=comp;deb.architecture=arch"},
		{"dist/comp1,comp2/arch", ";deb.distribution=dist;deb.component=comp1,comp2;deb.architecture=arch"},
		{"dist/comp/arch1,arch2", ";deb.distribution=dist;deb.component=comp;deb.architecture=arch1,arch2"},
		{"dist1,dist2/comp1,comp2/arch1,arch2", ";deb.distribution=dist1,dist2;deb.component=comp1,comp2;deb.architecture=arch1,arch2"},
	}

	for _, v := range debianPaths {
		result := getDebianProps(v.in)
		if result != v.expected {
			t.Errorf("getDebianProps(\"%s\") => '%s', want '%s'", v.in, result, v.expected)
		}
	}
}

func TestBuildUploadUrls(t *testing.T) {
	var testsParams = []struct {
		targetPath                  string
		targetProps                 string
		buildProps                  string
		expectedTargetPathWithProps string
	}{
		{"repo1/file1", "k1=v1", "k2=v2", "http://localhost:8881/artifactory/repo1/file1;k1=v1;k2=v2"},
		{"repo1/file@1", "k1=v1", "k2=v2", "http://localhost:8881/artifactory/repo1/file@1;k1=v1;k2=v2"},
		{"repo1/file;1", "k1=v1", "k2=v2", "http://localhost:8881/artifactory/repo1/file%3B1;k1=v1;k2=v2"},
		{"repo1/file,1", "k1=v1", "k2=v2", "http://localhost:8881/artifactory/repo1/file,1;k1=v1;k2=v2"},
		{"repo1/file^1", "k1=v1", "k2=v2", "http://localhost:8881/artifactory/repo1/file%5E1;k1=v1;k2=v2"},
		{"repo1/file:1", "k1=v1", "k2=v2", "http://localhost:8881/artifactory/repo1/file:1;k1=v1;k2=v2"},
		{"repo1/file1", "", "k2=v2", "http://localhost:8881/artifactory/repo1/file1;k2=v2"},
		{"repo1/file1", "k1=v1", "", "http://localhost:8881/artifactory/repo1/file1;k1=v1"},
		{"repo1/file1", "", "", "http://localhost:8881/artifactory/repo1/file1"},
	}

	for _, v := range testsParams {
		targetProps, e := utils.ParseProperties(v.targetProps)
		assert.NoError(t, e)
		actualTargetPathWithProps, e := buildUploadUrls("http://localhost:8881/artifactory/", v.targetPath, v.buildProps, "", targetProps)
		assert.NoError(t, e)
		assert.Equal(t, v.expectedTargetPathWithProps, actualTargetPathWithProps)
	}
}
