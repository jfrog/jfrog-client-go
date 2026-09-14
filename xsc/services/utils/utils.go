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
// and SCP clone URLs are converted to the equivalent host/path key, including
// Azure DevOps v3 → _git rewriting. Malformed SSH/SCP input returns "".
func GetGitRepoUrlKey(gitRepoUrl string) string {
	if gitRepoUrl == "" {
		return ""
	}
	if gitRepoKey, ok := sshGitRepoUrlKey(gitRepoUrl); ok {
		return ensureGitSuffix(gitRepoKey)
	}
	if isSSHCloneURL(gitRepoUrl) {
		return ""
	}
	return ensureGitSuffix(stripHTTPScheme(gitRepoUrl))
}

// GitCloneHostPath returns host and repository path for HTTP(S), ssh://, and
// SCP clone URLs without Xray git-repo-key rewrites (Azure v3 → _git).
func GitCloneHostPath(raw string) (host, repoPath string, ok bool) {
	return parseGitCloneHostPath(raw)
}

func sshGitRepoUrlKey(raw string) (string, bool) {
	if !isSSHCloneURL(raw) {
		return "", false
	}
	host, repoPath, ok := parseGitCloneHostPath(raw)
	if !ok {
		return "", false
	}
	return gitRepoKeyFromHostPath(host, repoPath), true
}

func parseGitCloneHostPath(raw string) (host, repoPath string, ok bool) {
	lower := strings.ToLower(raw)
	switch {
	case strings.HasPrefix(lower, "http://"), strings.HasPrefix(lower, "https://"):
		parsed, err := url.Parse(raw)
		if err != nil || parsed.Host == "" {
			return "", "", false
		}
		return parsed.Host, strings.Trim(parsed.EscapedPath(), "/"), true
	case strings.HasPrefix(lower, "ssh://"):
		parsed, err := url.Parse(raw)
		if err != nil || parsed.Hostname() == "" {
			return "", "", false
		}
		repoPath = strings.Trim(parsed.Path, "/")
		if repoPath == "" {
			return "", "", false
		}
		return parsed.Hostname(), repoPath, true
	}
	if strings.Contains(raw, "://") {
		return "", "", false
	}
	at := strings.LastIndex(raw, "@")
	if at < 0 {
		return "", "", false
	}
	value := raw[at+1:]
	colon := strings.Index(value, ":")
	if colon <= 0 || colon == len(value)-1 {
		return "", "", false
	}
	repoPath = strings.Trim(value[colon+1:], "/")
	if repoPath == "" {
		return "", "", false
	}
	return value[:colon], repoPath, true
}

func isSSHCloneURL(raw string) bool {
	lower := strings.ToLower(raw)
	if strings.HasPrefix(lower, "ssh://") {
		return true
	}
	if strings.Contains(raw, "://") {
		return false
	}
	at := strings.LastIndex(raw, "@")
	if at < 0 {
		return false
	}
	colon := strings.Index(raw[at+1:], ":")
	return colon > 0
}

func ensureGitSuffix(gitRepoUrl string) string {
	if !strings.HasSuffix(gitRepoUrl, ".git") {
		return gitRepoUrl + ".git"
	}
	return gitRepoUrl
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
