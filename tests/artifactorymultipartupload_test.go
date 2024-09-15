package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/stretchr/testify/assert"
)

const (
	bigFileSize   = 100 << 20
	propertyKey   = "prop-key"
	propertyValue = "prop-value"
)

func initArtifactoryMultipartUploadTest(t *testing.T) {
	if !*TestMultipartUpload {
		t.Skip("Skipping multipart upload test. To run artifactory test add the '-test.mpu=true' option.")
	}

	supported, err := testsUploadService.MultipartUpload.IsSupported(testsUploadService.ArtDetails)
	assert.NoError(t, err)
	if !supported {
		t.Skip("Skipping multipart upload test. Multipart upload test is not supported in the provided Artifactory server.")
	}
}

func TestArtifactoryMultipartUpload(t *testing.T) {
	initArtifactoryMultipartUploadTest(t)
	t.Run("multipartUpload", multipartUpload)
}

func multipartUpload(t *testing.T) {
	bigFile, cleanup := createBigFile(t)
	defer cleanup()

	// Create upload parameters
	up := services.NewUploadParams()
	props := utils.NewProperties()
	props.AddProperty(propertyKey, propertyValue)
	up.CommonParams = &utils.CommonParams{Pattern: bigFile.Name(), Target: getRtTargetRepo(), TargetProps: props}
	up.Flat = true
	up.MinChecksumDeploy = bigFileSize + 1
	up.MinSplitSize = bigFileSize

	// Upload file and verify success
	summary, err := testsUploadService.UploadFiles(up)
	assert.NoError(t, err)
	assert.Equal(t, 1, summary.TotalSucceeded)
	assert.Zero(t, summary.TotalFailed)

	// Search for the uploaded file in Artifactory
	searchParams := services.NewSearchParams()
	searchParams.Pattern = getRtTargetRepo()
	reader, err := testsSearchService.Search(searchParams)
	defer readerCloseAndAssert(t, reader)
	assert.NoError(t, err)
	length, err := reader.Length()
	assert.NoError(t, err)
	assert.Equal(t, 1, length)

	// Ensure existence of the uploaded file and verify properties
	for item := new(utils.ResultItem); reader.NextRecord(item) == nil; item = new(utils.ResultItem) {
		assert.Equal(t, filepath.Base(bigFile.Name()), item.Name)
		assert.Equal(t, propertyValue, item.GetProperty(propertyKey))
	}
	readerGetErrorAndAssert(t, reader)

	// Cleanup
	artifactoryCleanup(t)
}

func createBigFile(t *testing.T) (bigFile *os.File, cleanUp func()) {
	bigFile, err := fileutils.CreateTempFile()
	assert.NoError(t, err)

	cleanUp = func() {
		assert.NoError(t, bigFile.Close())
		assert.NoError(t, os.Remove(bigFile.Name()))
	}

	data := make([]byte, int(bigFileSize))
	_, err = bigFile.Write(data)
	assert.NoError(t, err)
	return
}
