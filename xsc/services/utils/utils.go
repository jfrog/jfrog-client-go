package utils

import (
	"fmt"
	"strings"

	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

const (
	XraySuffix                        = "/xray/"
	XscSuffix                         = "/xsc/"
	MinXrayVersionXscTransitionToXray = "3.108.0"
)

// From Xray version 3.108.0, XSC is transitioning to Xray as inner service. This function will return the backward compatibility URL.
func XrayUrlToXscUrl(xrayUrl, xrayVersion string) string {
	if err := utils.ValidateMinimumVersion(utils.Xray, xrayVersion, MinXrayVersionXscTransitionToXray); err != nil {
		log.Debug(fmt.Sprintf("Xray version is lower than %s, XSC is not an inner service in Xray.", MinXrayVersionXscTransitionToXray))
		return strings.Replace(xrayUrl, XraySuffix, XscSuffix, 1)
	}
	// Newer versions of Xray will have XSC as an inner service.
	return xrayUrl + "api/v1" + XscSuffix
}
