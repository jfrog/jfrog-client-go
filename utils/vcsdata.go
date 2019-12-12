package utils

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
)

type (
	VcsDetails struct {
		vcsRootDirectory sync.Map //key: path to dot git, value: vcsData
		vcsDirectory     sync.Map //key: path to folder, value: index of vcsRootDirectory
		vcsDirectorySize *int32   // Size of vcs folders entries
	}
	vcsData struct {
		url      string
		revision string
	}
)

const MAX_ENTRIES = 10000

func NewVcsDetals() *VcsDetails {
	return &VcsDetails{vcsRootDirectory: sync.Map{}, vcsDirectory: sync.Map{}, vcsDirectorySize: new(int32)}
}

func (this *VcsDetails) increment(num int32) {
	atomic.AddInt32(this.vcsDirectorySize, num)
}

func (this *VcsDetails) get() int32 {
	return atomic.LoadInt32(this.vcsDirectorySize)
}

/*
	Start search for '.git' dir for the current path, incase there is one, extract the details and add a new entry to the cache(hash-map,key:path,value:pointer to the git details).
	otherwise, search in the parent folder and try:
	1. search for .git, and save the details for the current dir and all subpath
	2. .git not found, go to parent dir and repeat
	3. not found on the root directory, add all subpath to cache with nil as a value
*/
func (this *VcsDetails) GetVcsData(path string) (revision, refUrl string, err error) {
	keys := strings.Split(path, string(os.PathSeparator))
	var subPath string
	var subPaths []string
	var vcsDataResult *vcsData
	for i := len(keys); i > 0; i-- {
		subPath = strings.Join(keys[:i], string(os.PathSeparator))
		// Try to get from cache
		if searchResult, found := this.searchCache(subPath); found {
			if data, ok := searchResult.(*vcsData); ok {
				if data != nil {
					revision, refUrl, vcsDataResult = data.revision, data.url, data
				}
			}
			break
		}
		// Begin dir search
		revision, refUrl, err = tryGetGitDetails(subPath, this)
		if revision != "" || refUrl != "" {
			vcsDataResult = &vcsData{revision: revision, url: refUrl}
			this.vcsRootDirectory.Store(subPath, vcsDataResult)
			break
		}
		if err != nil {
			return
		}
		subPaths = append(subPaths, subPath)
	}
	if size := len(subPaths); size > 0 {
		this.healthCheack()
		for _, v := range subPaths {
			this.vcsDirectory.Store(v, vcsDataResult)
		}
		this.increment(int32(size))
	}
	return
}

func (this *VcsDetails) healthCheack() {
	if this.get() > MAX_ENTRIES {
		this.vcsDirectory = sync.Map{}
		this.vcsDirectorySize = new(int32)
	}
}

func tryGetGitDetails(path string, this *VcsDetails) (string, string, error) {
	dotGitPath := filepath.Join(path, ".git")
	exists, err := fileutils.IsDirExists(dotGitPath, false)
	if exists {
		return extractGitInfo(dotGitPath)
	}
	return "", "", err
}

func extractGitInfo(path string) (revision string, refUrl string, err error) {
	dotGitPath := filepath.Join(path, "HEAD")
	file, er := os.Open(dotGitPath)
	if er != nil {
		err = er
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
		return
	}
	if revision == "" {
		dotGitRevision, er := os.Open(filepath.Join(path, refUrl))
		if er != nil {
			err = er
			return
		}
		defer dotGitRevision.Close()
		scanner = bufio.NewScanner(dotGitRevision)
		for scanner.Scan() {
			text := scanner.Text()
			revision = strings.TrimSpace(text)
			break
		}
		if err = scanner.Err(); err != nil {
			return
		}
	}
	return
}

func (this *VcsDetails) searchCache(path string) (gitData interface{}, found bool) {
	if gitData, found = this.vcsDirectory.Load(path); found {
		return
	}
	gitData, found = this.vcsRootDirectory.Load(path)
	return
}
