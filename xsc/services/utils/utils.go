package utils

import (
	"strings"

	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

const (
	XraySuffix                        = "/xray/"
	XscSuffix                         = "/xsc/"
	XscServiceSuffixInXray            = "api/v1" + XscSuffix
	MinXrayVersionXscTransitionToXray = "3.108.0"
)

// From Xray version 3.108.0, XSC is transitioning to Xray as inner service. This function will return the backward compatibility URL.
func XrayUrlToXscUrl(xrayUrl, xrayVersion string) string {
	if err := utils.ValidateMinimumVersion(utils.Xray, xrayVersion, MinXrayVersionXscTransitionToXray); err != nil {
		log.Debug("Xray version is lower than", MinXrayVersionXscTransitionToXray, "XSC is not an inner service in Xray.")
		return strings.Replace(xrayUrl, XraySuffix, XscSuffix, 1)
	}
	// Newer versions of Xray will have XSC as an inner service.
	return xrayUrl + XscServiceSuffixInXray
}
