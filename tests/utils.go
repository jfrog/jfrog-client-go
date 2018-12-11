package tests

import (
	"errors"
	"flag"
	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/httpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var RtUrl *string
var RtUser *string
var RtPassword *string
var RtApiKey *string
var RtSshKeyPath *string
var RtSshPassphrase *string
var RtAccessToken *string
var LogLevel *string
var testsUploadService *services.UploadService
var testsSearchService *services.SearchService
var testsDeleteService *services.DeleteService
var testsDownloadService *services.DownloadService

const (
	RtTargetRepo              = "jfrog-cli-tests-repo1/"
	SpecsTestRepositoryConfig = "specs_test_repository_config.json"
	RepoDetailsUrl            = "api/repositories/"
)

func init() {
	RtUrl = flag.String("rt.url", "http://localhost:8081/artifactory/", "Artifactory url")
	RtUser = flag.String("rt.user", "admin", "Artifactory username")
	RtPassword = flag.String("rt.password", "password", "Artifactory password")
	RtApiKey = flag.String("rt.apikey", "", "Artifactory user API key")
	RtSshKeyPath = flag.String("rt.sshKeyPath", "", "Ssh key file path")
	RtSshPassphrase = flag.String("rt.sshPassphrase", "", "Ssh key passphrase")
	RtAccessToken = flag.String("rt.accessToken", "", "Artifactory access token")
	LogLevel = flag.String("log-level", "INFO", "Sets the log level")
}

func createArtifactorySearchManager() {
	testsSearchService = services.NewSearchService(httpclient.NewDefaultHttpClient())
	testsSearchService.ArtDetails = getArtDetails()
}

func createArtifactoryDeleteManager() {
	testsDeleteService = services.NewDeleteService(httpclient.NewDefaultHttpClient())
	testsDeleteService.ArtDetails = getArtDetails()
}

func createArtifactoryUploadManager() {
	testsUploadService = services.NewUploadService(httpclient.NewDefaultHttpClient())
	testsUploadService.ArtDetails = getArtDetails()
	testsUploadService.Threads = 3
}

func createArtifactoryDownloadManager() {
	testsDownloadService = services.NewDownloadService(httpclient.NewDefaultHttpClient())
	testsDownloadService.ArtDetails = getArtDetails()
	testsDownloadService.SetThreads(3)
}

func getArtDetails() auth.ArtifactoryDetails {
	rtDetails := auth.NewArtifactoryDetails()
	rtDetails.SetUrl(clientutils.AddTrailingSlashIfNeeded(*RtUrl))
	if !fileutils.IsSshUrl(rtDetails.GetUrl()) {
		if *RtApiKey != "" {
			rtDetails.SetApiKey(*RtApiKey)
		} else if *RtAccessToken != "" {
			rtDetails.SetAccessToken(*RtAccessToken)
		} else {
			rtDetails.SetUser(*RtUser)
			rtDetails.SetPassword(*RtPassword)
		}
		return rtDetails
	}

	err := rtDetails.AuthenticateSsh(*RtSshKeyPath, *RtSshPassphrase)
	if err != nil {
		log.Error("Failed while attempting to authenticate with Artifactory: " + err.Error())
		os.Exit(1)
	}
	return rtDetails
}

func artifactoryCleanup(t *testing.T) {
	params := &utils.ArtifactoryCommonParams{Pattern: RtTargetRepo}
	toDelete, err := testsDeleteService.GetPathsToDelete(services.DeleteParams{ArtifactoryCommonParams: params})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	deleteItems := make([]services.DeleteItem, len(toDelete))
	for i, item := range toDelete {
		deleteItems[i] = item
	}
	deletedCount, err := testsDeleteService.DeleteFiles(deleteItems, testsDeleteService)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if len(toDelete) != deletedCount {
		t.Errorf("Failed to delete files from Artifactory expected %d items to be deleted got %d.", len(toDelete), deletedCount)
	}
}

func createReposIfNeeded() error {
	var err error
	var repoConfig string
	repo := RtTargetRepo
	if strings.HasSuffix(repo, "/") {
		repo = repo[0:strings.LastIndex(repo, "/")]
	}
	if !isRepoExist(repo) {
		repoConfig = filepath.Join(getTestDataPath(), "reposconfig", SpecsTestRepositoryConfig)
		err = execCreateRepoRest(repoConfig, repo)
		if err != nil {
			log.Error(err)
			return err
		}
	}
	return nil
}

func isRepoExist(repoName string) bool {
	artDetails := getArtDetails()
	artHttpDetails := artDetails.CreateHttpClientDetails()
	client := httpclient.NewDefaultHttpClient()
	resp, _, _, err := client.SendGet(artDetails.GetUrl()+RepoDetailsUrl+repoName, true, artHttpDetails)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	if resp.StatusCode != http.StatusBadRequest {
		return true
	}
	return false
}

func execCreateRepoRest(repoConfig, repoName string) error {
	content, err := ioutil.ReadFile(repoConfig)
	if err != nil {
		return err
	}
	artDetails := getArtDetails()
	artHttpDetails := artDetails.CreateHttpClientDetails()

	artHttpDetails.Headers = map[string]string{"Content-Type": "application/json"}
	client := httpclient.NewDefaultHttpClient()
	resp, _, err := client.SendPut(artDetails.GetUrl()+"api/repositories/"+repoName, content, artHttpDetails)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return errors.New("Fail to create repository. Reason local repository with key: " + repoName + " already exist\n")
	}
	log.Info("Repository", repoName, "created.")
	return nil
}

func getTestDataPath() string {
	dir, _ := os.Getwd()
	return filepath.Join(dir, "testsdata")
}

func FixWinPath(filePath string) string {
	fixedPath := strings.Replace(filePath, "\\", "\\\\", -1)
	return fixedPath
}
