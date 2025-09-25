package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/httpclient"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	clientio "github.com/jfrog/jfrog-client-go/utils/io"
	"github.com/jfrog/jfrog-client-go/utils/io/content"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type DirectDownloadService struct {
	client               *jfroghttpclient.JfrogHttpClient
	Progress             clientio.ProgressMgr
	artDetails           *auth.ServiceDetails
	DryRun               bool
	Threads              int
	saveSummary          bool
	filesTransfersWriter *content.ContentWriter
}

func NewDirectDownloadService(artDetails auth.ServiceDetails, client *jfroghttpclient.JfrogHttpClient) *DirectDownloadService {
	return &DirectDownloadService{artDetails: &artDetails, client: client}
}

func (dds *DirectDownloadService) GetJfrogHttpClient() *jfroghttpclient.JfrogHttpClient {
	return dds.client
}

func (dds *DirectDownloadService) IsDryRun() bool {
	return dds.DryRun
}

func (dds *DirectDownloadService) SetDryRun(isDryRun bool) {
	dds.DryRun = isDryRun
}

func (dds *DirectDownloadService) GetThreads() int {
	return dds.Threads
}

func (dds *DirectDownloadService) SetThreads(threads int) {
	dds.Threads = threads
}

func (dds *DirectDownloadService) SetArtifactoryDetails(artDetails auth.ServiceDetails) {
	dds.artDetails = &artDetails
}

func (dds *DirectDownloadService) SetSaveSummary(saveSummary bool) {
	dds.saveSummary = saveSummary
}

// DownloadFilesWithoutAQL downloads files directly from Artifactory without using AQL queries.
func (dds *DirectDownloadService) DownloadFilesWithoutAQL(downloadParams ...DirectDownloadParams) (totalDownloaded, totalFailed int, err error) {
	summary, err := dds.performDirectDownload(downloadParams...)
	if err != nil {
		return 0, 0, err
	}
	return summary.TotalSucceeded, summary.TotalFailed, nil
}

// DownloadFilesWithoutAQLWithSummary downloads files directly from Artifactory without using AQL queries
// and gives detailed information about each file transfer.
func (dds *DirectDownloadService) DownloadFilesWithoutAQLWithSummary(downloadParams ...DirectDownloadParams) (operationSummary *utils.OperationSummary, err error) {
	return dds.performDirectDownload(downloadParams...)
}

func (dds *DirectDownloadService) performDirectDownload(downloadParams ...DirectDownloadParams) (summary *utils.OperationSummary, err error) {
	summary = &utils.OperationSummary{}

	if dds.saveSummary {
		dds.filesTransfersWriter, err = content.NewContentWriter(content.DefaultKey, true, false)
		if err != nil {
			return nil, err
		}
	}

	for _, params := range downloadParams {
		repo, artifactPath, err := dds.parsePattern(params.GetPattern())
		if err != nil {
			log.Error(err)
			summary.TotalFailed++
			continue
		}

		if dds.isExcluded(artifactPath, params.GetExclusions()) {
			continue
		}

		if dds.containsWildcards(artifactPath) {
			// Handle patterns like "*.zip" or "test?.jar" by listing the directory
			// and downloading matching files one by one
			count, failed, err := dds.handleWildcardDownload(repo, artifactPath, &params)
			summary.TotalSucceeded += count
			summary.TotalFailed += failed
			if err != nil {
				log.Error(err)
			}
		} else {
			// directly download the file in case wildcards are not present
			success, err := dds.downloadSingleFile(repo, artifactPath, &params)
			switch {
			case err != nil:
				log.Error(err)
				summary.TotalFailed++
			case success:
				summary.TotalSucceeded++
			default:
				summary.TotalFailed++
			}
		}
	}

	if dds.saveSummary && dds.filesTransfersWriter != nil {
		if err = dds.filesTransfersWriter.Close(); err != nil {
			return nil, err
		}
		filePath := dds.filesTransfersWriter.GetFilePath()
		log.Debug("Creating content reader from file:", filePath)

		summary.TransferDetailsReader = content.NewContentReader(filePath, content.DefaultKey)
	} else {
		log.Debug("Not creating content reader - saveSummary:", dds.saveSummary, "filesTransfersWriter:", dds.filesTransfersWriter != nil)
	}

	return summary, nil
}

func (dds *DirectDownloadService) parsePattern(pattern string) (repo, artifactPath string, err error) {
	parts := strings.SplitN(pattern, "/", 2)
	if len(parts) < 2 {
		return "", "", errorutils.CheckErrorf("Invalid pattern format: %s. Should be 'repo/path/to/artifact'", pattern)
	}
	return parts[0], parts[1], nil
}

func (dds *DirectDownloadService) containsWildcards(path string) bool {
	return strings.ContainsAny(path, "*?")
}

func (dds *DirectDownloadService) isExcluded(path string, exclusions []string) bool {
	for _, exclusion := range exclusions {
		if matched, _ := filepath.Match(exclusion, path); matched {
			log.Debug("Artifact excluded by pattern:", path, "matches", exclusion)
			return true
		}
	}
	return false
}

func (dds *DirectDownloadService) downloadSingleFile(repo, artifactPath string, params *DirectDownloadParams) (bool, error) {
	downloadPath := fmt.Sprintf("%s/%s", repo, artifactPath)
	downloadUrl, err := clientutils.BuildUrl((*dds.artDetails).GetUrl(), downloadPath, make(map[string]string))
	if err != nil {
		return false, err
	}

	targetPath := params.GetTarget()
	if targetPath == "" {
		targetPath = "./"
	}

	var localPath string
	if params.IsFlat() {
		localPath = filepath.Join(targetPath, filepath.Base(artifactPath))
	} else {
		localPath = filepath.Join(targetPath, artifactPath)
	}

	localDir := filepath.Dir(localPath)
	if err := os.MkdirAll(localDir, 0755); err != nil {
		return false, errorutils.CheckError(err)
	}

	if dds.DryRun {
		log.Info("[Dry run] Would download:", downloadUrl, "to", localPath)
		return true, nil
	}

	httpClientsDetails := (*dds.artDetails).CreateHttpClientDetails()

	downloadFileDetails := &httpclient.DownloadFileDetails{
		DownloadPath:  downloadUrl,
		LocalPath:     filepath.Dir(localPath),
		LocalFileName: filepath.Base(localPath),
		SkipChecksum:  params.IsSkipChecksum(),
	}

	resp, err := dds.client.DownloadFile(downloadFileDetails, "", &httpClientsDetails, false, false)
	if err != nil {
		return false, err
	}

	if resp.StatusCode != http.StatusOK {
		return false, errorutils.CheckErrorf("Failed to download %s: HTTP %d", downloadUrl, resp.StatusCode)
	}

	var sha256 string
	if !params.IsSkipChecksum() {
		// Fetch file info from artifactory to verify the download was successful
		fileInfo, err := dds.getFileInfo(downloadUrl)
		if err == nil && fileInfo != nil {
			sha256 = fileInfo.Checksums.Sha256
			if err := dds.validateChecksumWithInfo(fileInfo, localPath); err != nil {
				log.Warn("Checksum validation failed for", localPath, ":", err)
			}
		} else if err != nil {
			log.Warn("Failed to get file info for checksum validation:", err)
		}
	}

	log.Info("Downloaded:", downloadUrl, "to", localPath)

	if dds.saveSummary && dds.filesTransfersWriter != nil {
		rtUrl := strings.TrimSuffix((*dds.artDetails).GetUrl(), "/")

		sourcePath := downloadUrl
		if strings.HasPrefix(sourcePath, rtUrl) {
			sourcePath = strings.TrimPrefix(sourcePath, rtUrl)
			if !strings.HasPrefix(sourcePath, "/") {
				sourcePath = "/" + sourcePath
			}
		}

		fileTransferDetails := clientutils.FileTransferDetails{
			SourcePath: sourcePath,
			TargetPath: localPath,
			RtUrl:      rtUrl,
			Sha256:     sha256,
		}
		log.Debug("Writing file transfer details - Source:", sourcePath, "Target:", localPath, "RtUrl:", rtUrl, "SHA256:", sha256)
		dds.filesTransfersWriter.Write(fileTransferDetails)
	}

	return true, nil
}

// handleWildcardDownload deals with patterns like "*.zip" or "logs/2024*.txt" while downloading the artifacts.
func (dds *DirectDownloadService) handleWildcardDownload(repo, pattern string, params *DirectDownloadParams) (int, int, error) {
	dir := filepath.Dir(pattern)
	filePattern := filepath.Base(pattern)

	storagePath := fmt.Sprintf("api/storage/%s/%s", repo, dir)
	listUrl, err := clientutils.BuildUrl((*dds.artDetails).GetUrl(), storagePath, make(map[string]string))
	if err != nil {
		return 0, 0, err
	}

	httpClientsDetails := (*dds.artDetails).CreateHttpClientDetails()
	resp, body, _, err := dds.client.SendGet(listUrl, true, &httpClientsDetails)
	if err != nil {
		return 0, 0, err
	}
	defer func() {
		if resp != nil && resp.Body != nil {
			if err := resp.Body.Close(); err != nil {
				log.Warn("Failed to close response body:", err)
			}
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return 0, 0, errorutils.CheckErrorf("Failed to list directory %s: Status %d", listUrl, resp.StatusCode)
	}

	var storageInfo struct {
		Children []struct {
			Uri    string `json:"uri"`
			Folder bool   `json:"folder"`
		} `json:"children"`
	}

	if err := json.Unmarshal(body, &storageInfo); err != nil {
		return 0, 0, err
	}

	downloadCount := 0
	failCount := 0

	for _, child := range storageInfo.Children {
		if child.Folder {
			continue
		}

		fileName := strings.TrimPrefix(child.Uri, "/")
		matched, err := filepath.Match(filePattern, fileName)
		if err != nil {
			return downloadCount, failCount, err
		}

		if matched {
			// Check if the matched file is not excluded
			if !dds.isExcluded(filepath.Join(dir, fileName), params.GetExclusions()) {
				filePath := filepath.Join(dir, fileName)
				// Download the matched file using the direct API
				success, err := dds.downloadSingleFile(repo, filePath, params)
				switch {
				case err != nil:
					log.Error("Failed to download", filePath, ":", err)
					failCount++
				case success:
					downloadCount++
				default:
					failCount++
				}
			}
		}
	}

	return downloadCount, failCount, nil
}

// getFileInfo fetches the details about a file, including its checksums
func (dds *DirectDownloadService) getFileInfo(downloadUrl string) (*utils.FileInfo, error) {
	artUrl := (*dds.artDetails).GetUrl()

	repoPath := strings.TrimPrefix(downloadUrl, artUrl)
	repoPath = strings.TrimPrefix(repoPath, "/")

	storagePath := fmt.Sprintf("api/storage/%s", repoPath)
	storageUrl, err := clientutils.BuildUrl(artUrl, storagePath, make(map[string]string))
	if err != nil {
		return nil, err
	}

	httpClientsDetails := (*dds.artDetails).CreateHttpClientDetails()
	resp, body, _, err := dds.client.SendGet(storageUrl, true, &httpClientsDetails)
	if err != nil {
		return nil, err
	}
	defer func() {
		if resp != nil && resp.Body != nil {
			if err := resp.Body.Close(); err != nil {
				log.Warn("Failed to close response body:", err)
			}
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, errorutils.CheckErrorf("Failed to get file info: HTTP %d", resp.StatusCode)
	}

	var fileInfo utils.FileInfo
	if err := json.Unmarshal(body, &fileInfo); err != nil {
		return nil, errorutils.CheckError(err)
	}

	return &fileInfo, nil
}

// validateChecksumWithInfo validates the downloaded file's checksums
func (dds *DirectDownloadService) validateChecksumWithInfo(fileInfo *utils.FileInfo, localPath string) error {
	localFileDetails, err := fileutils.GetFileDetails(localPath, true)
	if err != nil {
		return err
	}

	if fileInfo.Checksums.Md5 != "" && localFileDetails.Checksum.Md5 != fileInfo.Checksums.Md5 {
		return errorutils.CheckErrorf("MD5 checksum mismatch for %s. Expected: %s, Got: %s",
			localPath, fileInfo.Checksums.Md5, localFileDetails.Checksum.Md5)
	}

	if fileInfo.Checksums.Sha1 != "" && localFileDetails.Checksum.Sha1 != fileInfo.Checksums.Sha1 {
		return errorutils.CheckErrorf("SHA1 checksum mismatch for %s. Expected: %s, Got: %s",
			localPath, fileInfo.Checksums.Sha1, localFileDetails.Checksum.Sha1)
	}

	if fileInfo.Checksums.Sha256 != "" && localFileDetails.Checksum.Sha256 != fileInfo.Checksums.Sha256 {
		return errorutils.CheckErrorf("SHA256 checksum mismatch for %s. Expected: %s, Got: %s",
			localPath, fileInfo.Checksums.Sha256, localFileDetails.Checksum.Sha256)
	}

	log.Debug("Checksum validation passed for:", localPath)
	return nil
}
