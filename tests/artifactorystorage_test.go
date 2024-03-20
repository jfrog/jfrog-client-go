package tests

import (
	"errors"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/jfrog/jfrog-client-go/artifactory/services"
	servicesutils "github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/stretchr/testify/assert"
)

func TestArtifactoryStorage(t *testing.T) {
	initArtifactoryTest(t)
	uploadDummyFile(t)
	t.Run("folder info test", folderInfoTest)
	t.Run("file info test", fileInfoTest)
	t.Run("file list test", fileListTest)
	t.Run("storage info test", storageInfoTest)

	artifactoryCleanup(t)
}

func folderInfoTest(t *testing.T) {
	info, err := testsStorageService.FolderInfo(getRtTargetRepo() + "test/")
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, utils.AddTrailingSlashIfNeeded(*RtUrl)+path.Join(services.StorageRestApi, getRtTargetRepo()+"test/"), info.Uri)
	assert.Equal(t, strings.TrimSuffix(getRtTargetRepo(), "/"), info.Repo)
	assert.Equal(t, "/test", info.Path)
	assert.NotEmpty(t, info.Created)
	assert.NotEmpty(t, info.CreatedBy)
	assert.NotEmpty(t, info.LastModified)
	assert.NotEmpty(t, info.LastUpdated)
	assert.Len(t, info.Children, 1)
	assert.Equal(t, "/a.in", info.Children[0].Uri)
	assert.False(t, info.Children[0].Folder)
}

func fileInfoTest(t *testing.T) {
	info, err := testsStorageService.FileInfo(getRtTargetRepo() + "test/a.in")
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, utils.AddTrailingSlashIfNeeded(*RtUrl)+path.Join(services.StorageRestApi, getRtTargetRepo()+"test/a.in"), info.Uri)
	assert.Equal(t, strings.TrimSuffix(getRtTargetRepo(), "/"), info.Repo)
	assert.Equal(t, "/test/a.in", info.Path)
	assert.NotEmpty(t, info.Created)
	assert.NotEmpty(t, info.CreatedBy)
	assert.NotEmpty(t, info.LastModified)
	assert.NotEmpty(t, info.LastUpdated)
}

func fileListTest(t *testing.T) {
	params := servicesutils.NewFileListParams()
	params.Deep = true
	params.Depth = 2
	params.ListFolders = true
	params.MetadataTimestamps = true
	params.IncludeRootPath = true

	fileList, err := testsStorageService.FileList(getRtTargetRepo(), params)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, utils.AddTrailingSlashIfNeeded(*RtUrl)+path.Join(services.StorageRestApi, getRtTargetRepo()), fileList.Uri)
	assert.NotEmpty(t, fileList.Created)
	assert.Len(t, fileList.Files, 5)
	for _, file := range fileList.Files {
		if strings.HasSuffix(file.Uri, "a.in") {
			assert.Equal(t, "/test/a.in", file.Uri)
			assert.NotEmpty(t, file.Size)
			assert.NotEmpty(t, file.LastModified)
			assert.False(t, file.Folder)
			assert.NotEmpty(t, file.Sha1)
			assert.NotEmpty(t, file.Sha2)
			assert.NotEmpty(t, file.MetadataTimestamps.Properties)
			return
		}
	}
	assert.Fail(t, "could not find the expected dummy file")
}

func storageInfoTest(t *testing.T) {
	err := testsStorageService.StorageInfoRefresh()
	if !assert.NoError(t, err) {
		return
	}
	info, err := waitForRepoInStorageInfo(testsStorageService, getRtTargetRepoKey(), 5)
	if !assert.NoError(t, err) {
		return
	}

	assert.NotEmpty(t, info.BinariesCount)
	assert.NotEmpty(t, info.BinariesSize)
	assert.NotEmpty(t, info.ArtifactsSize)
	assert.NotEmpty(t, info.Optimization)
	assert.NotEmpty(t, info.ItemsCount)
	assert.NotEmpty(t, info.ArtifactsCount)
	assert.NotEmpty(t, info.StorageType)
	assert.NotEmpty(t, info.StorageDirectory)

	for _, repoSummary := range info.RepositoriesSummaryList {
		if repoSummary.RepoKey == getRtTargetRepoKey() {
			assert.NotEmpty(t, repoSummary.RepoType)
			assert.NotEmpty(t, repoSummary.FoldersCount)
			assert.NotEmpty(t, repoSummary.FilesCount)
			assert.NotEmpty(t, repoSummary.UsedSpace)
			assert.NotEmpty(t, repoSummary.UsedSpaceInBytes)
			assert.NotEmpty(t, repoSummary.ItemsCount)
			assert.NotEmpty(t, repoSummary.PackageType)
			assert.NotEmpty(t, repoSummary.ProjectKey)
			assert.NotEmpty(t, repoSummary.Percentage)
			return
		}
	}
	assert.Fail(t, "could not find the summary of repo '"+getRtTargetRepoKey()+"' in the storage info")
}

func waitForRepoInStorageInfo(testsStorageService *services.StorageService, repositoryKey string, timeoutInSeconds int) (*servicesutils.StorageInfo, error) {
	for i := 0; i < timeoutInSeconds; i++ {
		storageInfo, err := testsStorageService.StorageInfo()
		if err != nil {
			return nil, err
		}

		if _, err := storageInfo.FindRepositoryWithKey(repositoryKey); err == nil {
			return storageInfo, nil
		}

		time.Sleep(1 * time.Second)
	}

	return nil, errors.New("Failed to find a repositoryKey with the key " + repositoryKey + " within " + string(rune(timeoutInSeconds)) + " seconds")
}
