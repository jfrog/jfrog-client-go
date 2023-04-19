package services

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	gitconfig "github.com/go-git/go-git/v5/plumbing/format/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/content"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type GitLfsCleanService struct {
	client     *jfroghttpclient.JfrogHttpClient
	artDetails *auth.ServiceDetails
	DryRun     bool
}

func NewGitLfsCleanService(artDetails auth.ServiceDetails, client *jfroghttpclient.JfrogHttpClient) *GitLfsCleanService {
	return &GitLfsCleanService{artDetails: &artDetails, client: client}
}

func (glc *GitLfsCleanService) GetArtifactoryDetails() auth.ServiceDetails {
	return *glc.artDetails
}

func (glc *GitLfsCleanService) IsDryRun() bool {
	return glc.DryRun
}

func (glc *GitLfsCleanService) GetJfrogHttpClient() *jfroghttpclient.JfrogHttpClient {
	return glc.client
}

func (glc *GitLfsCleanService) GetUnreferencedGitLfsFiles(gitLfsCleanParams GitLfsCleanParams) (*content.ContentReader, error) {
	var err error
	repo := gitLfsCleanParams.GetRepo()
	gitPath := gitLfsCleanParams.GetGitPath()
	if gitPath == "" {
		gitPath, err = os.Getwd()
		if err != nil {
			return nil, errorutils.CheckError(err)
		}
	}
	if len(repo) == 0 {
		repo, err = detectRepo(gitPath, glc.GetArtifactoryDetails().GetUrl())
		if err != nil {
			return nil, err
		}
	}
	log.Info("Searching files from Artifactory repository", repo, "...")
	refsRegex := getRefsRegex(gitLfsCleanParams.GetRef())
	artifactoryLfsFilesReader, err := glc.searchLfsFilesInArtifactory(repo)
	if err != nil {
		return nil, errorutils.CheckError(err)
	}
	defer artifactoryLfsFilesReader.Close()
	log.Info("Collecting files to preserve from Git references matching the pattern", gitLfsCleanParams.GetRef(), "...")
	gitLfsFiles, err := getLfsFilesFromGit(gitPath, refsRegex)
	if err != nil {
		return nil, errorutils.CheckError(err)
	}
	filesToDeleteReader, err := findFilesToDelete(artifactoryLfsFilesReader, gitLfsFiles)
	if err != nil {
		return nil, err
	}
	length, err := filesToDeleteReader.Length()
	if err != nil {
		return nil, err
	}
	log.Info("Found", len(gitLfsFiles), "files to keep, and", length, "to clean")
	return filesToDeleteReader, nil
}

func findFilesToDelete(artifactoryLfsFilesReader *content.ContentReader, gitLfsFiles map[string]struct{}) (*content.ContentReader, error) {
	cw, err := content.NewContentWriter("results", true, false)
	if err != nil {
		return nil, err
	}
	defer cw.Close()
	for resultItem := new(utils.ResultItem); artifactoryLfsFilesReader.NextRecord(resultItem) == nil; resultItem = new(utils.ResultItem) {
		if _, keepFile := gitLfsFiles[resultItem.Name]; !keepFile {
			cw.Write(*resultItem)
		}
	}
	artifactoryLfsFilesReader.Reset()
	return content.NewContentReader(cw.GetFilePath(), cw.GetArrayKey()), nil
}

func lfsConfigUrlExtractor(conf *gitconfig.Config) (*url.URL, error) {
	return url.Parse(conf.Section("lfs").Option("url"))
}

func configLfsUrlExtractor(conf *gitconfig.Config) (*url.URL, error) {
	return url.Parse(conf.Section("remote").Subsection("origin").Option("lfsurl"))
}

func detectRepo(gitPath, rtUrl string) (string, error) {
	repo, err := extractRepo(gitPath, ".lfsconfig", rtUrl, lfsConfigUrlExtractor)
	if err == nil {
		return repo, nil
	}
	errMsg1 := fmt.Sprintf("Cannot detect Git LFS repository from .lfsconfig: %s", err.Error())
	repo, err = extractRepo(gitPath, ".git/config", rtUrl, configLfsUrlExtractor)
	if err == nil {
		return repo, nil
	}
	errMsg2 := fmt.Sprintf("Cannot detect Git LFS repository from .git/config: %s", err.Error())
	suggestedSolution := "You may want to try passing the --repo option manually"
	return "", errorutils.CheckErrorf("%s%s%s", errMsg1, errMsg2, suggestedSolution)
}

func extractRepo(gitPath, configFile, rtUrl string, lfsUrlExtractor lfsUrlExtractorFunc) (string, error) {
	lfsUrl, err := getLfsUrl(gitPath, configFile, lfsUrlExtractor)
	if err != nil {
		return "", err
	}
	artifactoryConfiguredUrl, err := url.Parse(rtUrl)
	if err != nil {
		return "", err
	}
	if artifactoryConfiguredUrl.Scheme != lfsUrl.Scheme || artifactoryConfiguredUrl.Host != lfsUrl.Host {
		return "", fmt.Errorf("configured Git LFS URL %q does not match provided URL %q", lfsUrl.String(), artifactoryConfiguredUrl.String())
	}
	artifactoryConfiguredUrlPath := path.Clean("/"+artifactoryConfiguredUrl.Path+"/api/lfs") + "/"
	lfsUrlPath := path.Clean(lfsUrl.Path)
	if strings.HasPrefix(lfsUrlPath, artifactoryConfiguredUrlPath) {
		return lfsUrlPath[len(artifactoryConfiguredUrlPath):], nil
	}
	return "", fmt.Errorf("configured Git LFS URL %q does not match provided URL %q", lfsUrl.String(), artifactoryConfiguredUrl.String())
}

type lfsUrlExtractorFunc func(conf *gitconfig.Config) (*url.URL, error)

func getLfsUrl(gitPath, configFile string, lfsUrlExtractor lfsUrlExtractorFunc) (lfsUrl *url.URL, err error) {
	lfsConf, err := os.Open(path.Join(gitPath, configFile))
	if err != nil {
		return nil, errorutils.CheckError(err)
	}
	defer func() {
		e := lfsConf.Close()
		if err == nil {
			err = errorutils.CheckError(e)
		}
	}()
	conf := gitconfig.New()
	err = gitconfig.NewDecoder(lfsConf).Decode(conf)
	if err != nil {
		return nil, errorutils.CheckError(err)
	}
	lfsUrl, err = lfsUrlExtractor(conf)
	return lfsUrl, errorutils.CheckError(err)
}

func getRefsRegex(refs string) string {
	replacer := strings.NewReplacer(",", "|", "\\*", ".*")
	return replacer.Replace(regexp.QuoteMeta(refs))
}

func (glc *GitLfsCleanService) searchLfsFilesInArtifactory(repo string) (*content.ContentReader, error) {
	spec := &utils.CommonParams{Pattern: repo, Target: "", Props: "", ExcludeProps: "", Build: "", Recursive: true, Regexp: false, IncludeDirs: false}
	return utils.SearchBySpecWithPattern(spec, glc, utils.NONE)
}

func getLfsFilesFromGit(path, refMatch string) (map[string]struct{}, error) {
	// a hash set of sha2 sums, to make lookup faster later
	results := make(map[string]struct{}, 0)
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, errorutils.CheckError(err)
	}
	log.Debug("Opened Git repo at", path, "for reading")
	refs, err := repo.References()
	if err != nil {
		return nil, errorutils.CheckError(err)
	}
	// look for every Git LFS pointer file that exists in any ref (branch,
	// remote branch, tag, etc.) who's name matches the regex refMatch
	err = refs.ForEach(func(ref *plumbing.Reference) error {
		// go-git recognizes three types of refs: regular hash refs,
		// symbolic refs (e.g. HEAD), and invalid refs. We only care
		// about the first type here.
		if ref.Type() != plumbing.HashReference {
			return nil
		}
		log.Debug("Checking ref", ref.Name().String())
		match, err := regexp.MatchString(refMatch, ref.Name().String())
		if err != nil || !match {
			return errorutils.CheckError(err)
		}
		var commit *object.Commit
		commit, err = repo.CommitObject(ref.Hash())
		if err != nil {
			commit, err = checkAnnotatedTag(ref, repo)
			if err != nil {
				return err
			}
		}
		files, err := commit.Files()
		if err != nil {
			return errorutils.CheckError(err)
		}
		err = files.ForEach(func(file *object.File) error {
			return collectLfsFileFromGit(results, file)
		})
		return errorutils.CheckError(err)
	})
	return results, errorutils.CheckError(err)
}

// checkAnnotatedTag checks the case of an annotated tag in which the commit is within the tag object
func checkAnnotatedTag(ref *plumbing.Reference, repo *git.Repository) (*object.Commit, error) {
	// Get the annotated tag object.
	annotatedTag, err := repo.TagObject(ref.Hash())
	if err != nil {
		return nil, errorutils.CheckError(err)
	}
	// Find the commit object associated with the annotated tag.
	commit, err := repo.CommitObject(annotatedTag.Target)
	return commit, errorutils.CheckError(err)
}

func collectLfsFileFromGit(results map[string]struct{}, file *object.File) error {
	// A Git LFS pointer is a small file containing a sha2. Any file bigger
	// than a kilobyte is extremely unlikely to be such a pointer.
	if file.Size > 1024 {
		return nil
	}
	lines, err := file.Lines()
	if err != nil {
		return errorutils.CheckError(err)
	}
	// the line containing the sha2 we're looking for will match this regex
	regex := "^oid sha256:[[:alnum:]]{64}$"
	for _, line := range lines {
		if !strings.HasPrefix(line, "oid ") {
			continue
		}
		match, err := regexp.MatchString(regex, line)
		if err != nil || !match {
			return errorutils.CheckError(err)
		}
		result := line[strings.Index(line, ":")+1:]
		log.Debug("Found file", result)
		results[result] = struct{}{}
		break
	}
	return nil
}

type GitLfsCleanParams struct {
	Refs    string
	Repo    string
	GitPath string
}

func (glc *GitLfsCleanParams) GetRef() string {
	return glc.Refs
}

func (glc *GitLfsCleanParams) GetRepo() string {
	return glc.Repo
}

func (glc *GitLfsCleanParams) GetGitPath() string {
	return glc.GitPath
}

func NewGitLfsCleanParams() GitLfsCleanParams {
	return GitLfsCleanParams{}
}
