package utils

import (
	"bufio"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	submoduleDotGitPrefix = "gitdir: "
)

type manager struct {
	path                string
	err                 error
	revision            string
	url                 string
	branch              string
	message             string
	submoduleDotGitPath string
}

func NewGitManager(path string) *manager {
	dotGitPath := filepath.Join(path, ".git")
	return &manager{path: dotGitPath}
}

func (m *manager) ReadConfig() error {
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
func (m *manager) handleSubmoduleIfNeeded() {
	exists, err := fileutils.IsFileExists(m.path, false)
	if err != nil {
		m.err = err
		return
	}
	if !exists {
		// .git is a directory, continue extracting vcs details.
		return
	}

	// Saving .git file path
	m.submoduleDotGitPath = m.path

	content, err := ioutil.ReadFile(m.path)
	if err != nil {
		m.err = errorutils.CheckError(err)
		return
	}

	line := string(content)
	// Expecting git submodule to have exactly one line, with a prefix and the path to the actual submodule's git.
	if !strings.HasPrefix(line, submoduleDotGitPrefix) {
		m.err = errorutils.CheckErrorf("failed to parse .git path for submodule")
		return
	}

	// Extract path by removing prefix.
	actualRelativePath := strings.TrimSpace(line[strings.Index(line, ":")+1:])
	actualAbsPath := filepath.Join(filepath.Dir(m.path), actualRelativePath)
	exists, err = fileutils.IsDirExists(actualAbsPath, false)
	if err != nil {
		m.err = err
		return
	}
	if !exists {
		m.err = errorutils.CheckErrorf("path found in .git file '" + m.path + "' does not exist: '" + actualAbsPath + "'")
		return
	}

	// Actual .git directory found.
	m.path = actualAbsPath
}

func (m *manager) GetUrl() string {
	return m.url
}

func (m *manager) GetRevision() string {
	return m.revision
}

func (m *manager) GetBranch() string {
	return m.branch
}

func (m *manager) GetMessage() string {
	return m.message
}

func (m *manager) readUrl() {
	if m.err != nil {
		return
	}
	dotGitPath := filepath.Join(m.path, "config")
	file, err := os.Open(dotGitPath)
	if errorutils.CheckError(err) != nil {
		m.err = err
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var IsNextLineUrl bool
	var originUrl string
	for scanner.Scan() {
		if IsNextLineUrl {
			text := scanner.Text()
			strings.HasPrefix(text, "url")
			originUrl = strings.TrimSpace(strings.SplitAfter(text, "=")[1])
			break
		}
		if scanner.Text() == "[remote \"origin\"]" {
			IsNextLineUrl = true
		}
	}
	if err := scanner.Err(); err != nil {
		errorutils.CheckError(err)
		m.err = err
		return
	}
	if !strings.HasSuffix(originUrl, ".git") {
		originUrl += ".git"
	}
	m.url = originUrl

	// Mask url if required
	regExp, err := GetRegExp(CredentialsInUrlRegexp)
	if err != nil {
		m.err = err
		return
	}
	matchedResult := regExp.FindString(originUrl)
	if matchedResult == "" {
		return
	}
	m.url = RemoveCredentials(originUrl, matchedResult)
}

func (m *manager) getRevisionAndBranchPath() (revision, refUrl string, err error) {
	dotGitPath := filepath.Join(m.path, "HEAD")
	file, e := os.Open(dotGitPath)
	if errorutils.CheckError(e) != nil {
		err = e
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasPrefix(text, "ref") {
			refUrl = strings.TrimSpace(strings.SplitAfter(text, ":")[1])
			break
		}
		revision = text
	}
	if err = scanner.Err(); err != nil {
		errorutils.CheckError(err)
	}
	return
}

func (m *manager) readRevisionAndBranch() {
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
		splitRefArr := strings.Split(ref, "/")
		m.branch = splitRefArr[len(splitRefArr)-1]
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

func (m *manager) readRevisionFromRef(refPath string) {
	revision := ""
	file, err := os.Open(refPath)
	if errorutils.CheckError(err) != nil {
		m.err = err
		return
	}
	defer file.Close()

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
	return
}

func (m *manager) readRevisionFromPackedRef(ref string) {
	packedRefPath := filepath.Join(m.path, "packed-refs")
	exists, err := fileutils.IsFileExists(packedRefPath, false)
	if err != nil {
		m.err = err
		return
	}
	if exists {
		file, err := os.Open(packedRefPath)
		if errorutils.CheckError(err) != nil {
			m.err = err
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			// Expecting to find the revision (the full extended SHA-1, or a unique leading substring) followed by the ref.
			if strings.HasSuffix(line, ref) {
				split := strings.Split(line, " ")
				if len(split) == 2 {
					m.revision = split[0]
				} else {
					m.err = errorutils.CheckErrorf("failed fetching revision for ref :" + ref + " - Unexpected line structure in packed-refs file")
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
	return
}

func (m *manager) readMessage() {
	if m.err != nil {
		return
	}
	var err error
	m.message, err = m.doReadMessage()
	if err != nil {
		log.Debug("Latest commit message was not extracted due to", err.Error())
	}
}

func (m *manager) doReadMessage() (string, error) {
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

func (m *manager) getPathHandleSubmodule() (path string) {
	if m.submoduleDotGitPath == "" {
		path = m.path
	} else {
		path = m.submoduleDotGitPath
	}
	path = strings.TrimSuffix(path, filepath.Join("", ".git"))
	return
}
