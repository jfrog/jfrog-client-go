package utils

import (
	"testing"

	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetLogger(log.NewLogger(log.INFO, nil))
}

func TestBuildParsingNoBuildNumber(t *testing.T) {
	buildName, buildNumber, err := ParseNameAndVersion("CLI-Build-Name", true)
	assert.NoError(t, err)
	expectedBuildName, expectedBuildNumber := "CLI-Build-Name", LatestBuildNumberKey
	if buildName != expectedBuildName {
		t.Error("Unexpected result from 'ParseNameAndVersion' method. \nExpected build name: 	" + expectedBuildName + " \nGot:     		 		" + buildName)
	}
	if buildNumber != expectedBuildNumber {
		t.Error("Unexpected result from 'ParseNameAndVersion' method. \nExpected build number: 	" + expectedBuildNumber + " \nGot:     			 	" + buildNumber)
	}
}

func TestBuildParsingBuildNumberProvided(t *testing.T) {
	buildName, buildNumber, err := ParseNameAndVersion("CLI-Build-Name/11", true)
	assert.NoError(t, err)
	expectedBuildName, expectedBuildNumber := "CLI-Build-Name", "11"
	if buildName != expectedBuildName {
		t.Error("Unexpected result from 'ParseNameAndVersion' method. \nExpected build name: 	" + expectedBuildName + " \nGot:     		 		" + buildName)
	}
	if buildNumber != expectedBuildNumber {
		t.Error("Unexpected result from 'ParseNameAndVersion' method. \nExpected build number: 	" + expectedBuildNumber + " \nGot:     			 	" + buildNumber)
	}
}

func TestBuildParsingBuildNumberWithEscapeCharsInTheBuildName(t *testing.T) {
	buildName, buildNumber, err := ParseNameAndVersion("CLI-Build-Name\\/a\\/b\\/c/11", true)
	assert.NoError(t, err)
	expectedBuildName, expectedBuildNumber := "CLI-Build-Name/a/b/c", "11"
	if buildName != expectedBuildName {
		t.Error("Unexpected result from 'ParseNameAndVersion' method. \nExpected build name: 	" + expectedBuildName + " \nGot:     		 		" + buildName)
	}
	if buildNumber != expectedBuildNumber {
		t.Error("Unexpected result from 'ParseNameAndVersion' method. \nExpected build number: 	" + expectedBuildNumber + " \nGot:     			 	" + buildNumber)
	}
}

func TestBuildParsingBuildNumberWithEscapeCharsInTheBuildNumber(t *testing.T) {
	buildName, buildNumber, err := ParseNameAndVersion("CLI-Build-Name/1\\/2\\/3\\/4", true)
	assert.NoError(t, err)
	expectedBuildName, expectedBuildNumber := "CLI-Build-Name", "1/2/3/4"
	if buildName != expectedBuildName {
		t.Error("Unexpected result from 'ParseNameAndVersion' method. \nExpected build name: 	" + expectedBuildName + " \nGot:     		 		" + buildName)
	}
	if buildNumber != expectedBuildNumber {
		t.Error("Unexpected result from 'ParseNameAndVersion' method. \nExpected build number: 	" + expectedBuildNumber + " \nGot:     			 	" + buildNumber)
	}
}

func TestBuildParsingBuildNumberWithOnlyEscapeChars(t *testing.T) {
	buildName, buildNumber, err := ParseNameAndVersion("CLI-Build-Name\\/1\\/2\\/3\\/4", true)
	assert.NoError(t, err)
	expectedBuildName, expectedBuildNumber := "CLI-Build-Name/1/2/3/4", LatestBuildNumberKey
	if buildName != expectedBuildName {
		t.Error("Unexpected result from 'ParseNameAndVersion' method. \nExpected build name: 	" + expectedBuildName + " \nGot:     		 		" + buildName)
	}
	if buildNumber != expectedBuildNumber {
		t.Error("Unexpected result from 'ParseNameAndVersion' method. \nExpected build number: 	" + expectedBuildNumber + " \nGot:     			 	" + buildNumber)
	}
}

func TestBundleParsingNoBundleVersion(t *testing.T) {
	log.SetLogger(log.NewLogger(log.DEBUG, nil))
	_, _, err := ParseNameAndVersion("CLI-Bundle-Name", false)
	assert.EqualError(t, err, "no '/' is found in 'CLI-Bundle-Name'")
}

func TestBundleParsingBundleVersionProvided(t *testing.T) {
	bundleName, bundleVersion, err := ParseNameAndVersion("CLI-Bundle-Name/11", false)
	assert.NoError(t, err)
	expectedBundleName, expectedBundleVersion := "CLI-Bundle-Name", "11"
	if bundleName != expectedBundleName {
		t.Error("Unexpected result from 'ParseNameAndVersion' method. \nExpected bundle name: 	" + expectedBundleName + " \nGot:     		 		" + bundleName)
	}
	if bundleVersion != expectedBundleVersion {
		t.Error("Unexpected result from 'ParseNameAndVersion' method. \nExpected bundle version: 	" + expectedBundleVersion + " \nGot:     			 	" + bundleVersion)
	}
}

func TestBundleParsingBundleVersionWithEscapeCharsInTheBundleName(t *testing.T) {
	bundleName, bundleVersion, err := ParseNameAndVersion("CLI-Bundle-Name\\/a\\/b\\/c/11", false)
	assert.NoError(t, err)
	expectedBundleName, expectedBundleVersion := "CLI-Bundle-Name/a/b/c", "11"
	if bundleName != expectedBundleName {
		t.Error("Unexpected result from 'ParseNameAndVersion' method. \nExpected bundle name: 	" + expectedBundleName + " \nGot:     		 		" + bundleName)
	}
	if bundleVersion != expectedBundleVersion {
		t.Error("Unexpected result from 'ParseNameAndVersion' method. \nExpected bundle version: 	" + expectedBundleVersion + " \nGot:     			 	" + bundleVersion)
	}
}

func TestBundleParsingBundleVersionWithEscapeCharsInTheBundleVersion(t *testing.T) {
	bundleName, bundleVersion, err := ParseNameAndVersion("CLI-Bundle-Name/1\\/2\\/3\\/4", false)
	assert.NoError(t, err)
	expectedBundleName, expectedBundleVersion := "CLI-Bundle-Name", "1/2/3/4"
	if bundleName != expectedBundleName {
		t.Error("Unexpected result from 'ParseNameAndVersion' method. \nExpected bundle name: 	" + expectedBundleName + " \nGot:     		 		" + bundleName)
	}
	if bundleVersion != expectedBundleVersion {
		t.Error("Unexpected result from 'ParseNameAndVersion' method. \nExpected bundle version: 	" + expectedBundleVersion + " \nGot:     			 	" + bundleVersion)
	}
}

func TestBundleParsingBundleVersionWithOnlyEscapeChars(t *testing.T) {
	_, _, err := ParseNameAndVersion("CLI-Bundle-Name\\/1\\/2\\/3\\/4", false)
	assert.EqualError(t, err, "no delimiter char (/) without escaping char was found in 'CLI-Bundle-Name\\/1\\/2\\/3\\/4'")
}
