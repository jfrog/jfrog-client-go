package utils

import (
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"regexp"
	"strings"
)

// #nosec G101 -- False positive - no hardcoded credentials.
const CredentialsInUrlRegexp = `(http|https|git)://.+@`

func GetRegExp(regex string) (*regexp.Regexp, error) {
	regExp, err := regexp.Compile(regex)
	if errorutils.CheckError(err) != nil {
		return nil, err
	}
	return regExp, nil
}

// Remove credentials from the URL contained in the input line.
// The credentials are built as 'user:password' or 'token'
// For example:
// line = 'This is a line http://user:password@127.0.0.1:8081/artifactory/path/to/repo'
// credentialsPart = 'http://user:password@'
// Returned value: 'This is a line http://127.0.0.1:8081/artifactory/path/to/repo'
//
// line = 'This is a line http://token@127.0.0.1:8081/artifactory/path/to/repo'
// credentialsPart = 'http://token@'
// Returned value: 'This is a line http://127.0.0.1:8081/artifactory/path/to/repo'
func RemoveCredentials(line, credentialsPart string) string {
	splitResult := strings.Split(credentialsPart, "//")
	return strings.Replace(line, credentialsPart, splitResult[0]+"//", 1)
}
