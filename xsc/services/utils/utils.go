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
	} else {
		gitRepoUrl = stripHTTPScheme(gitRepoUrl)
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
		return gitRepoKeyFromHostPath(parsed.Hostname(), strings.TrimPrefix(parsed.Path, "/")), true
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
	return gitRepoKeyFromHostPath(value[:colon], value[colon+1:]), true
}

func stripHTTPScheme(raw string) string {
	lower := strings.ToLower(raw)
	switch {
	case strings.HasPrefix(lower, "https://"):
		return raw[len("https://"):]
	case strings.HasPrefix(lower, "http://"):
		return raw[len("http://"):]
	default:
		return raw
	}
}

func gitRepoKeyFromHostPath(host, repoPath string) string {
	repoPath = strings.Trim(repoPath, "/")
	if rewritten, ok := azureDevOpsHTTPSGitRepoPath(host, repoPath); ok {
		return rewritten
	}
	return host + "/" + repoPath
}

func azureDevOpsHTTPSGitRepoPath(host, repoPath string) (string, bool) {
	if !isAzureDevOpsSSHHost(host) {
		return "", false
	}
	parts := strings.Split(repoPath, "/")
	if len(parts) < 4 || !strings.EqualFold(parts[0], "v3") {
		return "", false
	}
	return "dev.azure.com/" + parts[1] + "/" + parts[2] + "/_git/" + strings.Join(parts[3:], "/"), true
}

func isAzureDevOpsSSHHost(host string) bool {
	return strings.EqualFold(host, "ssh.dev.azure.com") ||
		strings.EqualFold(host, "vs-ssh.visualstudio.com") ||
		strings.EqualFold(host, "dev.azure.com")
}
