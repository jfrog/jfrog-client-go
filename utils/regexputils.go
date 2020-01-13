package utils

import (
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"regexp"
	"strings"
)

const CredentialsInUrlRegexp = `((http|https):\/\/[%|\w]+:[%|\w]+@)`

func GetRegExp(regex string) (*regexp.Regexp, error) {
	regExp, err := regexp.Compile(regex)
	if errorutils.WrapError(err) != nil {
		return nil, err
	}
	return regExp, nil
}

// Mask the credentials information from the completeUrl, contained in credentialsPart.
// The credentials are built as user:password
// For example:
// completeUrl = http://user:password@127.0.0.1:8081/artifactory/path/to/repo
// credentialsPart = http://user:password@
// Returned value: http://***:***@127.0.0.1:8081/artifactory/path/to/repo
func MaskCredentials(completeUrl, credentialsPart string) string {
	splitResult := strings.Split(credentialsPart, "//")
	return strings.Replace(completeUrl, credentialsPart, splitResult[0]+"//***.***@", 1)
}
