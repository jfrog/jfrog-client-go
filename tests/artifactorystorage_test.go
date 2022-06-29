package tests

import (
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	servicesutils "github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/stretchr/testify/assert"
	"path"
	"strings"
	"testing"
)

func TestArtifactoryStorage(t *testing.T) {
	initArtifactoryTest(t)
	uploadDummyFile(t)
	t.Run("file info test", fileInfoTest)
	t.Run("file list test", fileListTest)
	t.Run("storage info test", storageInfoTest)

	artifactoryCleanup(t)
}

func fileInfoTest(t *testing.T) {
	info, err := testsStorageService.FolderInfo(getRtTargetRepo() + "test/")
	if err != nil {
		assert.NoError(t, err)
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

func fileListTest(t *testing.T) {
	params := servicesutils.NewFileListParams()
	params.Deep = true
	params.Depth = 2
	params.ListFolders = true
	params.MetadataTimestamps = true
	params.IncludeRootPath = true

	fileList, err := testsStorageService.FileList(getRtTargetRepo(), params)
	if err != nil {
		assert.NoError(t, err)
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
	info, err := testsStorageService.StorageInfo()
	if err != nil {
		assert.NoError(t, err)
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
	assert.NotEmpty(t, info.TotalSpace)
	assert.NotEmpty(t, info.UsedSpace)
	assert.NotEmpty(t, info.FreeSpace)

	// Verifying the repositories summary was filled correctly.
	// Cannot check any repo that is created as part of the test suite because it does not reflect in storage info immediately.
	// Cannot verify the first value because the repos "TOTAL" is returned there.
	// Checking the 2nd value.
	assert.True(t, len(info.RepositoriesSummaryList) > 1)
	repo := info.RepositoriesSummaryList[1]
	assert.NotEmpty(t, repo.RepoKey)
	assert.NotEmpty(t, repo.RepoType)
	assert.NotEmpty(t, repo.FoldersCount)
	assert.NotEmpty(t, repo.FilesCount)
	assert.NotEmpty(t, repo.UsedSpace)
	assert.NotEmpty(t, repo.UsedSpaceInBytes)
	assert.NotEmpty(t, repo.ItemsCount)
}
