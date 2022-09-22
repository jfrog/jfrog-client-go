package tests

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/io/content"
	"github.com/jfrog/jfrog-client-go/utils/log"

	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils/tests"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	testutils "github.com/jfrog/jfrog-client-go/utils/tests"
	"github.com/stretchr/testify/assert"
)

func TestArtifactoryUpload(t *testing.T) {
	initArtifactoryTest(t)
	t.Run("flat", flatUpload)
	t.Run("recursive", recursiveUpload)
	t.Run("placeholder", placeholderUpload)
	t.Run("includeDirs", includeDirsUpload)
	t.Run("explode", explodeUpload)
	t.Run("props", propsUpload)
	t.Run("summary", summaryUpload)
}

func flatUpload(t *testing.T) {
	workingDir, _ := createWorkingDir(t)
	defer testutils.RemoveAllAndAssert(t, workingDir)

	pattern := filepath.Join(workingDir, "out", "*")
	up := services.NewUploadParams()
	up.CommonParams = &utils.CommonParams{Pattern: pattern, Recursive: true, Target: getRtTargetRepo()}
	up.Flat = true
	summary, err := testsUploadService.UploadFiles(up)
	if err != nil {
		t.Error(err)
	}
	if summary.TotalSucceeded != 1 {
		t.Error("Expected to upload 1 file.")
	}
	if summary.TotalFailed != 0 {
		t.Error("Failed to upload", summary.TotalFailed, "files.")
	}
	searchParams := services.NewSearchParams()
	searchParams.CommonParams = &utils.CommonParams{}
	searchParams.Pattern = getRtTargetRepo()
	reader, err := testsSearchService.Search(searchParams)
	defer readerCloseAndAssert(t, reader)
	if err != nil {
		t.Error(err)
	}
	for item := new(utils.ResultItem); reader.NextRecord(item) == nil; item = new(utils.ResultItem) {
		if item.Path != "." {
			t.Error("Expected path to be root due to using the flat flag.", "Got:", item.Path)
		}
	}
	readerGetErrorAndAssert(t, reader)
	length, err := reader.Length()
	assert.NoError(t, err)
	if length > 1 {
		t.Error("Expected single file.")
	}
	artifactoryCleanup(t)
}

func recursiveUpload(t *testing.T) {
	workingDir, _ := createWorkingDir(t)
	defer testutils.RemoveAllAndAssert(t, workingDir)

	pattern := filepath.Join(workingDir, "*")
	up := services.NewUploadParams()
	up.CommonParams = &utils.CommonParams{Pattern: pattern, Recursive: true, Target: getRtTargetRepo()}
	up.Flat = true
	summary, err := testsUploadService.UploadFiles(up)
	if err != nil {
		t.Error(err)
	}
	if summary.TotalSucceeded != 1 {
		t.Error("Expected to upload 1 file.")
	}
	if summary.TotalFailed != 0 {
		t.Error("Failed to upload", summary.TotalFailed, "files.")
	}
	searchParams := services.NewSearchParams()
	searchParams.CommonParams = &utils.CommonParams{}
	searchParams.Pattern = getRtTargetRepo()
	reader, err := testsSearchService.Search(searchParams)
	defer readerCloseAndAssert(t, reader)
	if err != nil {
		t.Error(err)
	}
	for item := new(utils.ResultItem); reader.NextRecord(item) == nil; item = new(utils.ResultItem) {
		if item.Path != "." {
			t.Error("Expected path to be root(flat by default).", "Got:", item.Path)
		}
		if item.Name != "a.in" {
			t.Error("Missing File a.in")
		}
	}
	readerGetErrorAndAssert(t, reader)
	length, err := reader.Length()
	assert.NoError(t, err)
	if length > 1 {
		t.Error("Expected single file.")
	}
	artifactoryCleanup(t)
}

func placeholderUpload(t *testing.T) {
	workingDir, _ := createWorkingDir(t)
	defer testutils.RemoveAllAndAssert(t, workingDir)

	pattern := filepath.Join(workingDir, "(*).in")
	up := services.NewUploadParams()
	up.CommonParams = &utils.CommonParams{Pattern: pattern, Recursive: true, Target: getRtTargetRepo() + "{1}"}
	up.Flat = true
	summary, err := testsUploadService.UploadFiles(up)
	if err != nil {
		t.Error(err)
	}
	if summary.TotalSucceeded != 1 {
		t.Error("Expected to upload 1 file.")
	}
	if summary.TotalFailed != 0 {
		t.Error("Failed to upload", summary.TotalFailed, "files.")
	}
	searchParams := services.NewSearchParams()
	searchParams.CommonParams = &utils.CommonParams{}
	searchParams.Pattern = getRtTargetRepo()
	reader, err := testsSearchService.Search(searchParams)
	defer readerCloseAndAssert(t, reader)
	if err != nil {
		t.Error(err)
	}
	for item := new(utils.ResultItem); reader.NextRecord(item) == nil; item = new(utils.ResultItem) {
		if item.Path != "out" {
			t.Error("Expected path to be out.", "Got:", item.Path)
		}
		if item.Name != "a" {
			t.Error("Missing File a")
		}
	}
	readerGetErrorAndAssert(t, reader)
	length, err := reader.Length()
	assert.NoError(t, err)
	if length > 1 {
		t.Error("Expected single file.")
	}
	artifactoryCleanup(t)
}

func includeDirsUpload(t *testing.T) {
	workingDir, _ := createWorkingDir(t)
	defer testutils.RemoveAllAndAssert(t, workingDir)

	pattern := filepath.Join(workingDir, "*")
	up := services.NewUploadParams()
	up.CommonParams = &utils.CommonParams{Pattern: pattern, IncludeDirs: true, Recursive: false, Target: getRtTargetRepo()}
	up.Flat = true
	summary, err := testsUploadService.UploadFiles(up)
	if err != nil {
		t.Error(err)
	}
	if summary.TotalSucceeded != 0 {
		t.Error("Expected to upload 1 file.")
	}
	if summary.TotalFailed != 0 {
		t.Error("Failed to upload", summary.TotalFailed, "files.")
	}
	searchParams := services.NewSearchParams()
	searchParams.CommonParams = &utils.CommonParams{}
	searchParams.Pattern = getRtTargetRepo()
	searchParams.IncludeDirs = true
	reader, err := testsSearchService.Search(searchParams)
	defer readerCloseAndAssert(t, reader)
	if err != nil {
		t.Error(err)
	}
	for item := new(utils.ResultItem); reader.NextRecord(item) == nil; item = new(utils.ResultItem) {
		if item.Name == "." {
			continue
		}
		if item.Path != "." {
			t.Error("Expected path to be root(flat by default).", "Got:", item.Path)
		}
		if item.Name != "out" {
			t.Error("Missing directory out.")
		}
	}
	readerGetErrorAndAssert(t, reader)
	length, err := reader.Length()
	assert.NoError(t, err)
	if length < 2 {
		t.Error("Expected to get at least two items, default and the out folder.")
	}
	artifactoryCleanup(t)
}

func explodeUpload(t *testing.T) {
	workingDir, filePath := createWorkingDir(t)
	defer testutils.RemoveAllAndAssert(t, workingDir)

	err := fileutils.ZipFolderFiles(filePath, filepath.Join(workingDir, "zipFile.zip"))
	if err != nil {
		t.Fatal(err)
	}
	err = os.Remove(filePath)
	if err != nil {
		t.Fatal(err)
	}
	pattern := filepath.Join(workingDir, "*.zip")
	up := services.NewUploadParams()
	up.CommonParams = &utils.CommonParams{Pattern: pattern, IncludeDirs: true, Recursive: false, Target: getRtTargetRepo()}
	up.Flat = true
	up.ExplodeArchive = true
	summary, err := testsUploadService.UploadFiles(up)
	if err != nil {
		t.Error(err)
	}
	if summary.TotalSucceeded != 1 {
		t.Error("Expected to upload 1 file.")
	}
	if summary.TotalFailed != 0 {
		t.Error("Failed to upload", summary.TotalFailed, "files.")
	}
	searchParams := services.NewSearchParams()
	searchParams.CommonParams = &utils.CommonParams{}
	searchParams.Pattern = getRtTargetRepo()
	searchParams.IncludeDirs = true
	reader, err := testsSearchService.Search(searchParams)
	defer readerCloseAndAssert(t, reader)
	if err != nil {
		t.Error(err)
	}
	for item := new(utils.ResultItem); reader.NextRecord(item) == nil; item = new(utils.ResultItem) {
		if item.Name == "." {
			continue
		}
		if item.Name != "a.in" {
			t.Error("Missing file a.in")
		}
	}
	readerGetErrorAndAssert(t, reader)
	length, err := reader.Length()
	assert.NoError(t, err)
	if length < 2 {
		t.Error("Expected to get at least two items, default and the out folder.")
	}
	artifactoryCleanup(t)
}

func propsUpload(t *testing.T) {
	workingDir, _ := createWorkingDir(t)
	defer testutils.RemoveAllAndAssert(t, workingDir)

	// Upload a.in with property key1=val1
	pattern := filepath.Join(workingDir, "out", "*")
	targetProps, err := utils.ParseProperties("key1=val1")
	assert.NoError(t, err)
	up := services.NewUploadParams()
	up.CommonParams = &utils.CommonParams{Pattern: pattern, Target: getRtTargetRepo(), TargetProps: targetProps}
	up.Flat = true
	summary, err := testsUploadService.UploadFiles(up)
	assert.NoError(t, err)
	assert.Equal(t, 1, summary.TotalSucceeded)
	assert.Equal(t, 0, summary.TotalFailed)

	// Search a.in with property key1=val1
	searchParams := services.NewSearchParams()
	searchParams.CommonParams = &utils.CommonParams{}
	searchParams.Pattern = getRtTargetRepo()
	searchParams.Props = "key1=val1"
	reader, err := testsSearchService.Search(searchParams)
	defer readerCloseAndAssert(t, reader)
	if err != nil {
		t.Error(err)
	}
	length, err := reader.Length()
	assert.NoError(t, err)
	assert.Equal(t, 1, length)

	// Assert property key and value exist in the search results
	item := new(utils.ResultItem)
	err = reader.NextRecord(item)
	assert.NoError(t, err)
	assert.Len(t, item.Properties, 1)
	assert.Equal(t, "key1", item.Properties[0].Key)
	assert.Equal(t, "val1", item.Properties[0].Value)

	artifactoryCleanup(t)
}

func summaryUpload(t *testing.T) {
	pattern := filepath.Join("testdata", "a", "*")
	up := services.NewUploadParams()
	up.CommonParams = &utils.CommonParams{Pattern: pattern, Recursive: true, Target: getRtTargetRepo()}
	up.Flat = true
	testsUploadService.SetSaveSummary(true)
	defer testsUploadService.SetSaveSummary(false)
	summary, err := testsUploadService.UploadFiles(up)
	if err != nil {
		t.Error(err)
	}
	defer func() {
		assert.NoError(t, summary.Close())
	}()
	if summary.TotalSucceeded != 1 {
		t.Error("Expected to upload 1 file.")
	}
	if summary.TotalFailed != 0 {
		t.Error("Failed to upload", summary.TotalFailed, "files.")
	}
	var transfers []clientutils.FileTransferDetails
	for item := new(clientutils.FileTransferDetails); summary.TransferDetailsReader.NextRecord(item) == nil; item = new(clientutils.FileTransferDetails) {
		transfers = append(transfers, *item)
	}
	expectedSha256 := "4eb341b5d2762a853d79cc25e622aa8b978eb6e12c3259e2d99dc9dc60d82c5d"
	assert.Len(t, transfers, 1)
	assert.Equal(t, filepath.Join("testdata", "a", "a.in"), transfers[0].SourcePath)
	assert.Equal(t, testsUploadService.ArtDetails.GetUrl()+getRtTargetRepo()+"a.in", transfers[0].RtUrl+transfers[0].TargetPath)
	assert.Equal(t, expectedSha256, transfers[0].Sha256)
	var artifacts []utils.ArtifactDetails
	for item := new(utils.ArtifactDetails); summary.ArtifactsDetailsReader.NextRecord(item) == nil; item = new(utils.ArtifactDetails) {
		artifacts = append(artifacts, *item)
	}
	assert.Len(t, artifacts, 1)
	assert.Equal(t, getRtTargetRepo()+"a.in", artifacts[0].ArtifactoryPath)
	artifactoryCleanup(t)
}

func createWorkingDir(t *testing.T) (string, string) {
	workingDir, relativePath, err := tests.CreateFileWithContent("a.in", "/out/")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	return workingDir, relativePath
}

func readerCloseAndAssert(t *testing.T, reader *content.ContentReader) {
	assert.NoError(t, reader.Close(), "Couldn't close reader")
}

func readerGetErrorAndAssert(t *testing.T, reader *content.ContentReader) {
	assert.NoError(t, reader.GetError(), "Couldn't get reader error")
}

func TestUploadFilesWithFailure(t *testing.T) {
	// Create Artifactory mock server
	port := startArtifactoryMockServer(createUploadFilesWithFailureHandlers())
	client, err := jfroghttpclient.JfrogClientBuilder().
		Build()
	if err != nil {
		t.Error(err)
	}

	// Create Artifactory mock details
	rtDetails := auth.NewArtifactoryDetails()
	rtDetails.SetUrl("http://localhost:" + strconv.Itoa(port))
	rtDetails.SetUser("user")
	rtDetails.SetPassword("password")

	// Create upload service
	params := services.NewUploadParams()
	dir, err := os.Getwd()
	assert.NoError(t, err)
	params.Pattern = filepath.Join(dir, "testdata", "upload", "folder*")
	params.Target = "/generic"
	params.Flat = true
	params.Recursive = true
	service := services.NewUploadService(client)
	service.Threads = 1
	service.SetServiceDetails(rtDetails)

	// Upload files
	summary, err := service.UploadFiles(params)

	// Check for expected results
	assert.Error(t, err)
	assert.Equal(t, summary.TotalSucceeded, 1)
	assert.Equal(t, summary.TotalFailed, 1)
}

// Creates handlers for TestUploadFilesWithFailure mock server.
// The first upload request returns 200, and the rest return 404.
func createUploadFilesWithFailureHandlers() *testutils.HttpServerHandlers {
	handlers := testutils.HttpServerHandlers{}
	counter := 0
	handlers["/generic"] = func(w http.ResponseWriter, r *http.Request) {
		if counter == 0 {
			fmt.Fprintln(w, "{\"checksums\":{\"sha256\":\"123\"}}")
			w.WriteHeader(http.StatusOK)
			counter++
		} else {
			http.Error(w, "404 page not found", http.StatusNotFound)
		}
	}
	return &handlers
}

func startArtifactoryMockServer(handlers *testutils.HttpServerHandlers) int {
	port, err := testutils.StartHttpServer(*handlers)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	return port
}
