package utils

import "strings"

func GetGitRepoUrlKey(gitRepoUrl string) string {
	if len(gitRepoUrl) == 0 {
		// No git context was provided
		return ""
	}
	if !strings.HasSuffix(gitRepoUrl, ".git") {
		// Append .git to the URL if not included
		gitRepoUrl += ".git"
	}
	// Remove the Http/s protocol from the URL
	if strings.HasPrefix(gitRepoUrl, "http") {
		return strings.TrimPrefix(strings.TrimPrefix(gitRepoUrl, "https://"), "http://")
	}
	// Remove the SSH protocol from the URL
	if strings.Contains(gitRepoUrl, "git@") {
		return strings.Replace(strings.TrimPrefix(strings.TrimPrefix(gitRepoUrl, "ssh://"), "git@"), ":", "/", 1)
	}
	return gitRepoUrl
}
