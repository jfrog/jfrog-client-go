package utils

import (
	"github.com/jfrog/jfrog-client-go/utils/log"
	"testing"
)

func TestBuildParsingNoBuildNumber(t *testing.T) {
	log.SetLogger(log.NewLogger(log.DEBUG, nil))
	buildName, buildNumber := parseBuildNameAndNumber("CLI-Build-Name")
	expectedBuildName, expectedBuildNumber := "CLI-Build-Name", "LATEST"
	if buildName != expectedBuildName {
		t.Error("Unexpected result from 'parseBuildNameAndNumber' method. \nExpected build name: 	" + expectedBuildName + " \nGot:     		 		" + buildName)
	}
	if buildNumber != expectedBuildNumber {
		t.Error("Unexpected result from 'parseBuildNameAndNumber' method. \nExpected build number: 	" + expectedBuildNumber + " \nGot:     			 	" + buildNumber)
	}
}

func TestBuildParsingBuildNumberProvided(t *testing.T) {
	buildName, buildNumber := parseBuildNameAndNumber("CLI-Build-Name/11")
	expectedBuildName, expectedBuildNumber := "CLI-Build-Name", "11"
	if buildName != expectedBuildName {
		t.Error("Unexpected result from 'parseBuildNameAndNumber' method. \nExpected build name: 	" + expectedBuildName + " \nGot:     		 		" + buildName)
	}
	if buildNumber != expectedBuildNumber {
		t.Error("Unexpected result from 'parseBuildNameAndNumber' method. \nExpected build number: 	" + expectedBuildNumber + " \nGot:     			 	" + buildNumber)
	}
}

func TestBuildParsingBuildNumberWithEscapeCharsInTheBuildName(t *testing.T) {
	buildName, buildNumber := parseBuildNameAndNumber("CLI-Build-Name\\/a\\/b\\/c/11")
	expectedBuildName, expectedBuildNumber := "CLI-Build-Name/a/b/c", "11"
	if buildName != expectedBuildName {
		t.Error("Unexpected result from 'parseBuildNameAndNumber' method. \nExpected build name: 	" + expectedBuildName + " \nGot:     		 		" + buildName)
	}
	if buildNumber != expectedBuildNumber {
		t.Error("Unexpected result from 'parseBuildNameAndNumber' method. \nExpected build number: 	" + expectedBuildNumber + " \nGot:     			 	" + buildNumber)
	}
}

func TestBuildParsingBuildNumberWithEscapeCharsInTheBuildNumber(t *testing.T) {
	buildName, buildNumber := parseBuildNameAndNumber("CLI-Build-Name/1\\/2\\/3\\/4")
	expectedBuildName, expectedBuildNumber := "CLI-Build-Name", "1/2/3/4"
	if buildName != expectedBuildName {
		t.Error("Unexpected result from 'parseBuildNameAndNumber' method. \nExpected build name: 	" + expectedBuildName + " \nGot:     		 		" + buildName)
	}
	if buildNumber != expectedBuildNumber {
		t.Error("Unexpected result from 'parseBuildNameAndNumber' method. \nExpected build number: 	" + expectedBuildNumber + " \nGot:     			 	" + buildNumber)
	}
}

func TestBuildParsingBuildNumberWithOnlyEscapeChars(t *testing.T) {
	buildName, buildNumber := parseBuildNameAndNumber("CLI-Build-Name\\/1\\/2\\/3\\/4")
	expectedBuildName, expectedBuildNumber := "CLI-Build-Name/1/2/3/4", "LATEST"
	if buildName != expectedBuildName {
		t.Error("Unexpected result from 'parseBuildNameAndNumber' method. \nExpected build name: 	" + expectedBuildName + " \nGot:     		 		" + buildName)
	}
	if buildNumber != expectedBuildNumber {
		t.Error("Unexpected result from 'parseBuildNameAndNumber' method. \nExpected build number: 	" + expectedBuildNumber + " \nGot:     			 	" + buildNumber)
	}
}
