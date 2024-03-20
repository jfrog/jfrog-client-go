package utils

import (
	"bufio"
	"bytes"
	"errors"
	ioutils "github.com/jfrog/gofrog/io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type GitManager struct {
	path                string
	err                 error
	revision            string
	url                 string
	branch              string
	message             string
	submoduleDotGitPath string
}

func NewGitManager(path string) *GitManager {
	dotGitPath := filepath.Join(path, ".git")
	return &GitManager{path: dotGitPath}
}

func (m *GitManager) ExecGit(args ...string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("git", args...)
	cmd.Stdin = nil
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return strings.TrimSpace(stdout.String()), strings.TrimSpace(stderr.String()), errorutils.CheckError(err)
}

func (m *GitManager) ReadConfig() error {
	if m.path == "" {
		return errorutils.CheckErrorf(".git path must be defined")
	}
	if !fileutils.IsPathExists(m.path, false) {
		return errorutils.CheckErrorf(".git path must exist in order to collect vcs details")
	}

	m.handleSubmoduleIfNeeded()
	m.readRevisionAndBranch()
	m.readUrl()
	if m.revision != "" {
		m.readMessage()
	}
	return m.err
}

// If .git is a file and not a directory, assume it is a git submodule and extract the actual .git directory of the submodule.
// The actual .git directory is under the parent project's .git/modules directory.
func (m *GitManager) handleSubmoduleIfNeeded() {
	exists, err := fileutils.IsFileExists(m.path, false)
	if err != nil {
		m.err = err
		return
	}
	if !exists {
		// .git is a directory, continue extracting vcs details.
		return
	}
	// ask git for where the .git directory is directly for submodules and worktrees
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("git", "rev-parse", "--git-common-dir")
	cmd.Dir = filepath.Dir(m.path)
	cmd.Stdin = nil
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if m.err = errors.Join(m.err, err); m.err != nil {
		return
	}
	resolvedGitPath := strings.TrimSpace(stdout.String())
	exists, err = fileutils.IsDirExists(resolvedGitPath, false)
	if m.err = errors.Join(m.err, err); m.err != nil {
		return
	}
	if !exists {
		m.err = errorutils.CheckErrorf("path found in .git file '" + m.path + "' does not exist: '" + resolvedGitPath + "'")
		return
	}
	m.path = resolvedGitPath
}

func (m *GitManager) GetUrl() string {
	return m.url
}

func (m *GitManager) GetRevision() string {
	return m.revision
}

func (m *GitManager) GetBranch() string {
	return m.branch
}

func (m *GitManager) GetMessage() string {
	return m.message
}

func (m *GitManager) readUrl() {
	if m.err != nil {
		return
	}
	dotGitPath := filepath.Join(m.path, "config")
	file, err := os.Open(dotGitPath)
	if err != nil {
		m.err = err
		return
	}
	defer func() {
		m.err = errors.Join(m.err, errorutils.CheckError(file.Close()))
	}()

	scanner := bufio.NewScanner(file)
	var IsNextLineUrl bool
	var originUrl string
	for scanner.Scan() {
		if IsNextLineUrl {
			text := strings.TrimSpace(scanner.Text())
			if strings.HasPrefix(text, "url") {
				originUrl = strings.TrimSpace(strings.SplitAfter(text, "=")[1])
				break
			}
		}
		if scanner.Text() == "[remote \"origin\"]" {
			IsNextLineUrl = true
		}
	}
	if err := scanner.Err(); err != nil {
		m.err = errorutils.CheckError(err)
		return
	}
	if !strings.HasSuffix(originUrl, ".git") {
		originUrl += ".git"
	}
	m.url = originUrl

	// Mask url if required
	matchedResult := regexp.MustCompile(CredentialsInUrlRegexp).FindString(originUrl)
	if matchedResult == "" {
		return
	}
	m.url = RemoveCredentials(originUrl, matchedResult)
}

func (m *GitManager) getRevisionAndBranchPath() (revision, refUrl string, err error) {
	dotGitPath := filepath.Join(m.path, "HEAD")
	file, err := os.Open(dotGitPath)
	if errorutils.CheckError(err) != nil {
		return
	}
	defer ioutils.Close(file, &err)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasPrefix(text, "ref") {
			refUrl = strings.TrimSpace(strings.SplitAfter(text, ":")[1])
			break
		}
		revision = text
	}
	err = errorutils.CheckError(scanner.Err())
	return
}

func (m *GitManager) readRevisionAndBranch() {
	if m.err != nil {
		return
	}
	// This function will either return the revision or the branch ref:
	revision, ref, err := m.getRevisionAndBranchPath()
	if err != nil {
		m.err = err
		return
	}
	if ref != "" {
		// Get branch short name (refs/heads/master > master)
		m.branch = plumbing.ReferenceName(ref).Short()
	}
	// If the revision was returned, then we're done:
	if revision != "" {
		m.revision = revision
		return
	}

	// Else, if found ref try getting revision using it.
	refPath := filepath.Join(m.path, ref)
	exists, err := fileutils.IsFileExists(refPath, false)
	if err != nil {
		m.err = err
		return
	}
	if exists {
		m.readRevisionFromRef(refPath)
		return
	}
	// Otherwise, try to find .git/packed-refs and look for the HEAD there
	m.readRevisionFromPackedRef(ref)
}

func (m *GitManager) readRevisionFromRef(refPath string) {
	revision := ""
	file, err := os.Open(refPath)
	if err != nil {
		m.err = err
		return
	}
	defer func() {
		m.err = errors.Join(m.err, errorutils.CheckError(file.Close()))
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		revision = strings.TrimSpace(text)
		break
	}
	if err := scanner.Err(); err != nil {
		m.err = errorutils.CheckError(err)
		return
	}
	m.revision = revision
}

func (m *GitManager) readRevisionFromPackedRef(ref string) {
	packedRefPath := filepath.Join(m.path, "packed-refs")
	exists, err := fileutils.IsFileExists(packedRefPath, false)
	if err != nil {
		m.err = err
		return
	}
	if exists {
		file, err := os.Open(packedRefPath)
		if err != nil {
			m.err = err
			return
		}
		defer func() {
			m.err = errors.Join(m.err, errorutils.CheckError(file.Close()))
		}()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			// Expecting to find the revision (the full extended SHA-1, or a unique leading substring) followed by the ref.
			if strings.HasSuffix(line, ref) {
				split := strings.Split(line, " ")
				if len(split) == 2 {
					m.revision = split[0]
				} else {
					m.err = errors.Join(err, errorutils.CheckErrorf("failed fetching revision for ref :"+ref+" - Unexpected line structure in packed-refs file"))
				}
				return
			}
		}
		if err = scanner.Err(); err != nil {
			m.err = errorutils.CheckError(err)
			return
		}
	}
	log.Debug("No packed-refs file was found. Assuming git repository is empty")
}

func (m *GitManager) readMessage() {
	if m.err != nil {
		return
	}
	var err error
	m.message, err = m.doReadMessage()
	if err != nil {
		log.Debug("Latest commit message was not extracted due to", err.Error())
	}
}

func (m *GitManager) doReadMessage() (string, error) {
	path := m.getPathHandleSubmodule()
	gitRepo, err := git.PlainOpenWithOptions(path, &git.PlainOpenOptions{DetectDotGit: false})
	if errorutils.CheckError(err) != nil {
		return "", err
	}
	hash, err := gitRepo.ResolveRevision(plumbing.Revision(m.revision))
	if errorutils.CheckError(err) != nil {
		return "", err
	}
	message, err := gitRepo.CommitObject(*hash)
	if errorutils.CheckError(err) != nil {
		return "", err
	}
	return strings.TrimSpace(message.Message), nil
}

func (m *GitManager) getPathHandleSubmodule() (path string) {
	if m.submoduleDotGitPath == "" {
		path = m.path
	} else {
		path = m.submoduleDotGitPath
	}
	path = strings.TrimSuffix(path, filepath.Join("", ".git"))
	return
}
