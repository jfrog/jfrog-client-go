package utils

import (
	"fmt"
	"net/url"
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
	MinXrayVersionGitIntegrationEvent = "3.135.0"
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

// GetGitRepoUrlKey returns the repository URL in the legacy graph-scan key
// format expected by Xray. Existing HTTP(S) behavior is preserved, while SSH
// and SCP clone URLs are converted to the equivalent host/path key.
func GetGitRepoUrlKey(gitRepoUrl string) string {
	if gitRepoUrl == "" {
		return ""
	}
	if gitRepoKey, ok := sshGitRepoUrlKey(gitRepoUrl); ok {
		gitRepoUrl = gitRepoKey
	} else if strings.HasPrefix(gitRepoUrl, "http") {
		// Preserve the historical behavior for HTTP(S) URLs.
		gitRepoUrl = strings.TrimPrefix(strings.TrimPrefix(gitRepoUrl, "https://"), "http://")
	}
	if !strings.HasSuffix(gitRepoUrl, ".git") {
		gitRepoUrl += ".git"
	}
	return gitRepoUrl
}

func sshGitRepoUrlKey(raw string) (string, bool) {
	if strings.HasPrefix(strings.ToLower(raw), "ssh://") {
		parsed, err := url.Parse(raw)
		if err != nil || parsed.Hostname() == "" {
			return "", false
		}
		return azureDevOpsGitRepoKey(parsed.Hostname(), strings.TrimPrefix(parsed.Path, "/")), true
	}
	if strings.Contains(raw, "://") {
		return "", false
	}
	at := strings.LastIndex(raw, "@")
	if at < 0 {
		return "", false
	}
	value := raw
	value = value[at+1:]
	colon := strings.Index(value, ":")
	if colon <= 0 || colon == len(value)-1 {
		return "", false
	}
	return azureDevOpsGitRepoKey(value[:colon], value[colon+1:]), true
}

func azureDevOpsGitRepoKey(host, repoPath string) string {
	if strings.EqualFold(host, "ssh.dev.azure.com") {
		host = "dev.azure.com"
		parts := strings.Split(repoPath, "/")
		if len(parts) >= 4 && strings.EqualFold(parts[0], "v3") {
			repoPath = parts[1] + "/" + parts[2] + "/_git/" + strings.Join(parts[3:], "/")
		}
	}
	return host + "/" + strings.Trim(repoPath, "/")
}
