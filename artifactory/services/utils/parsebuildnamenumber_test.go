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
	buildName, buildNumber, err := parseNameAndVersion("CLI-Build-Name", true)
	assert.NoError(t, err)
	expectedBuildName, expectedBuildNumber := "CLI-Build-Name", "LATEST"
	if buildName != expectedBuildName {
		t.Error("Unexpected result from 'parseNameAndVersion' method. \nExpected build name: 	" + expectedBuildName + " \nGot:     		 		" + buildName)
	}
	if buildNumber != expectedBuildNumber {
		t.Error("Unexpected result from 'parseNameAndVersion' method. \nExpected build number: 	" + expectedBuildNumber + " \nGot:     			 	" + buildNumber)
	}
}

func TestBuildParsingBuildNumberProvided(t *testing.T) {
	buildName, buildNumber, err := parseNameAndVersion("CLI-Build-Name/11", true)
	assert.NoError(t, err)
	expectedBuildName, expectedBuildNumber := "CLI-Build-Name", "11"
	if buildName != expectedBuildName {
		t.Error("Unexpected result from 'parseNameAndVersion' method. \nExpected build name: 	" + expectedBuildName + " \nGot:     		 		" + buildName)
	}
	if buildNumber != expectedBuildNumber {
		t.Error("Unexpected result from 'parseNameAndVersion' method. \nExpected build number: 	" + expectedBuildNumber + " \nGot:     			 	" + buildNumber)
	}
}

func TestBuildParsingBuildNumberWithEscapeCharsInTheBuildName(t *testing.T) {
	buildName, buildNumber, err := parseNameAndVersion("CLI-Build-Name\\/a\\/b\\/c/11", true)
	assert.NoError(t, err)
	expectedBuildName, expectedBuildNumber := "CLI-Build-Name/a/b/c", "11"
	if buildName != expectedBuildName {
		t.Error("Unexpected result from 'parseNameAndVersion' method. \nExpected build name: 	" + expectedBuildName + " \nGot:     		 		" + buildName)
	}
	if buildNumber != expectedBuildNumber {
		t.Error("Unexpected result from 'parseNameAndVersion' method. \nExpected build number: 	" + expectedBuildNumber + " \nGot:     			 	" + buildNumber)
	}
}

func TestBuildParsingBuildNumberWithEscapeCharsInTheBuildNumber(t *testing.T) {
	buildName, buildNumber, err := parseNameAndVersion("CLI-Build-Name/1\\/2\\/3\\/4", true)
	assert.NoError(t, err)
	expectedBuildName, expectedBuildNumber := "CLI-Build-Name", "1/2/3/4"
	if buildName != expectedBuildName {
		t.Error("Unexpected result from 'parseNameAndVersion' method. \nExpected build name: 	" + expectedBuildName + " \nGot:     		 		" + buildName)
	}
	if buildNumber != expectedBuildNumber {
		t.Error("Unexpected result from 'parseNameAndVersion' method. \nExpected build number: 	" + expectedBuildNumber + " \nGot:     			 	" + buildNumber)
	}
}

func TestBuildParsingBuildNumberWithOnlyEscapeChars(t *testing.T) {
	buildName, buildNumber, err := parseNameAndVersion("CLI-Build-Name\\/1\\/2\\/3\\/4", true)
	assert.NoError(t, err)
	expectedBuildName, expectedBuildNumber := "CLI-Build-Name/1/2/3/4", "LATEST"
	if buildName != expectedBuildName {
		t.Error("Unexpected result from 'parseNameAndVersion' method. \nExpected build name: 	" + expectedBuildName + " \nGot:     		 		" + buildName)
	}
	if buildNumber != expectedBuildNumber {
		t.Error("Unexpected result from 'parseNameAndVersion' method. \nExpected build number: 	" + expectedBuildNumber + " \nGot:     			 	" + buildNumber)
	}
}

func TestBundleParsingNoBundleVersion(t *testing.T) {
	log.SetLogger(log.NewLogger(log.DEBUG, nil))
	_, _, err := parseNameAndVersion("CLI-Bundle-Name", false)
	assert.EqualError(t, err, "No '/' is found in the bundle")

}

func TestBundleParsingBundleVersionProvided(t *testing.T) {
	bundleName, bundleVersion, err := parseNameAndVersion("CLI-Bundle-Name/11", false)
	assert.NoError(t, err)
	expectedBundleName, expectedBundleVersion := "CLI-Bundle-Name", "11"
	if bundleName != expectedBundleName {
		t.Error("Unexpected result from 'parseNameAndVersion' method. \nExpected bundle name: 	" + expectedBundleName + " \nGot:     		 		" + bundleName)
	}
	if bundleVersion != expectedBundleVersion {
		t.Error("Unexpected result from 'parseNameAndVersion' method. \nExpected bundle version: 	" + expectedBundleVersion + " \nGot:     			 	" + bundleVersion)
	}
}

func TestBundleParsingBundleVersionWithEscapeCharsInTheBundleName(t *testing.T) {
	bundleName, bundleVersion, err := parseNameAndVersion("CLI-Bundle-Name\\/a\\/b\\/c/11", false)
	assert.NoError(t, err)
	expectedBundleName, expectedBundleVersion := "CLI-Bundle-Name/a/b/c", "11"
	if bundleName != expectedBundleName {
		t.Error("Unexpected result from 'parseNameAndVersion' method. \nExpected bundle name: 	" + expectedBundleName + " \nGot:     		 		" + bundleName)
	}
	if bundleVersion != expectedBundleVersion {
		t.Error("Unexpected result from 'parseNameAndVersion' method. \nExpected bundle version: 	" + expectedBundleVersion + " \nGot:     			 	" + bundleVersion)
	}
}

func TestBundleParsingBundleVersionWithEscapeCharsInTheBundleVersion(t *testing.T) {
	bundleName, bundleVersion, err := parseNameAndVersion("CLI-Bundle-Name/1\\/2\\/3\\/4", false)
	assert.NoError(t, err)
	expectedBundleName, expectedBundleVersion := "CLI-Bundle-Name", "1/2/3/4"
	if bundleName != expectedBundleName {
		t.Error("Unexpected result from 'parseNameAndVersion' method. \nExpected bundle name: 	" + expectedBundleName + " \nGot:     		 		" + bundleName)
	}
	if bundleVersion != expectedBundleVersion {
		t.Error("Unexpected result from 'parseNameAndVersion' method. \nExpected bundle version: 	" + expectedBundleVersion + " \nGot:     			 	" + bundleVersion)
	}
}

func TestBundleParsingBundleVersionWithOnlyEscapeChars(t *testing.T) {
	_, _, err := parseNameAndVersion("CLI-Bundle-Name\\/1\\/2\\/3\\/4", false)
	assert.EqualError(t, err, "No delimiter char (/) without escaping char was found in the bundle")
}
