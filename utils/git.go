package utils

import (
	"bufio"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"os"
	"path/filepath"
	"strings"
)

type manager struct {
	path     string
	err      error
	revision string
	url      string
}

func NewGitManager(path string) *manager {
	dotGitPath := filepath.Join(path, ".git")
	return &manager{path: dotGitPath}
}

func (m *manager) ReadConfig() error {
	if m.path == "" {
		return errorutils.NewError(".git path must be defined.")
	}
	m.readRevision()
	m.readUrl()
	return m.err
}

func (m *manager) GetUrl() string {
	return m.url
}

func (m *manager) GetRevision() string {
	return m.revision
}

func (m *manager) readUrl() {
	if m.err != nil {
		return
	}
	dotGitPath := filepath.Join(m.path, "config")
	file, err := os.Open(dotGitPath)
	if err != nil {
		m.err = errorutils.WrapError(err)
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
		m.err = errorutils.WrapError(err)
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
	m.url = MaskCredentials(originUrl, matchedResult)
}

func (m *manager) getRevisionOrBranchPath() (revision, refUrl string, err error) {
	dotGitPath := filepath.Join(m.path, "HEAD")
	file, e := os.Open(dotGitPath)
	if e != nil {
		err = errorutils.WrapError(e)
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
		err = errorutils.WrapError(err)
	}
	return
}

func (m *manager) readRevision() {
	if m.err != nil {
		return
	}
	// This function will either return the revision or the branch ref:
	revision, ref, err := m.getRevisionOrBranchPath()
	if err != nil {
		m.err = err
		return
	}
	// If the revision was returned, then we're done:
	if revision != "" {
		m.revision = revision
		return
	}

	// Since the revision was not returned, then we'll fetch it, by using the ref:
	dotGitPath := filepath.Join(m.path, ref)
	file, err := os.Open(dotGitPath)
	if err != nil {
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
		m.err = errorutils.WrapError(err)
		return
	}

	m.revision = revision
	return
}
