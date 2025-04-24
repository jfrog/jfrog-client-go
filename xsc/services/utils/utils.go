package utils

import (
	"fmt"
	"strings"

	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

const (
	XraySuffix                        = "/xray/"
	xscSuffix                         = "/xsc/"
	apiV1Suffix                       = "api/v1"
	XscInXraySuffix                   = apiV1Suffix + xscSuffix
	MinXrayVersionXscTransitionToXray = "3.107.13"
	MinXrayVersionNewGitInfoContext   = "3.117.0"
)

// From Xray version 3.107.13, XSC is transitioning to Xray as inner service. This function will return compatible URL.
func XrayUrlToXscUrl(xrayUrl, xrayVersion string) string {
	if !IsXscXrayInnerService(xrayVersion) {
		log.Debug(fmt.Sprintf("Xray version is lower than %s, XSC is not an inner service in Xray.", MinXrayVersionXscTransitionToXray))
		return strings.Replace(xrayUrl, XraySuffix, xscSuffix, 1) + apiV1Suffix + "/"
	}
	// Newer versions of Xray will have XSC as an inner service.
	return xrayUrl + XscInXraySuffix
}

func IsXscXrayInnerService(xrayVersion string) bool {
	if err := utils.ValidateMinimumVersion(utils.Xray, xrayVersion, MinXrayVersionXscTransitionToXray); err != nil {
		return false
	}
	return true
}

// The platform expects the git repo key to be in the format of the https/http clone Git URL without the protocol.
func GetGitRepoUrlKey(gitRepoHttpUrl string) string {
	if len(gitRepoHttpUrl) == 0 {
		// No git context was provided
		return ""
	}
	if !strings.HasSuffix(gitRepoHttpUrl, ".git") {
		// Append .git to the URL if not included
		gitRepoHttpUrl += ".git"
	}
	// Remove the Http/s protocol from the URL
	if strings.HasPrefix(gitRepoHttpUrl, "http") {
		return strings.TrimPrefix(strings.TrimPrefix(gitRepoHttpUrl, "https://"), "http://")
	}
	return gitRepoHttpUrl
}
