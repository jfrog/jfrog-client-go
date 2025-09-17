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
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	clientio "github.com/jfrog/jfrog-client-go/utils/io"
	"github.com/jfrog/jfrog-client-go/utils/io/content"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type DirectDownloadService struct {
	client                 *jfroghttpclient.JfrogHttpClient
	Progress               clientio.ProgressMgr
	artDetails             *auth.ServiceDetails
	DryRun                 bool
	Threads                int
	saveSummary            bool
	filesTransfersWriter   *content.ContentWriter
	artifactsDetailsWriter *content.ContentWriter
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

func (dds *DirectDownloadService) DirectDownloadFiles(downloadParams ...DirectDownloadParams) (totalDownloaded, totalFailed int, err error) {
	summary, err := dds.performDirectDownload(downloadParams...)
	if err != nil {
		return 0, 0, err
	}
	return summary.TotalSucceeded, summary.TotalFailed, nil
}

func (dds *DirectDownloadService) DirectDownloadFilesWithSummary(downloadParams ...DirectDownloadParams) (operationSummary *utils.OperationSummary, err error) {
	return dds.performDirectDownload(downloadParams...)
}

func (dds *DirectDownloadService) performDirectDownload(downloadParams ...DirectDownloadParams) (summary *utils.OperationSummary, err error) {
	summary = &utils.OperationSummary{}

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
			count, failed, err := dds.handleWildcardDownload(repo, artifactPath, &params)
			summary.TotalSucceeded += count
			summary.TotalFailed += failed
			if err != nil {
				log.Error(err)
			}
		} else {
			success, err := dds.downloadSingleFile(repo, artifactPath, &params)
			if err != nil {
				log.Error(err)
				summary.TotalFailed++
			} else if success {
				summary.TotalSucceeded++
			} else {
				summary.TotalFailed++
			}
		}
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
	downloadUrl := fmt.Sprintf("%s%s/%s", (*dds.artDetails).GetUrl(), repo, artifactPath)

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
	resp, body, _, err := dds.client.SendGet(downloadUrl, true, &httpClientsDetails)
	if err != nil {
		return false, err
	}
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return false, errorutils.CheckErrorf("Failed to download %s: HTTP %d", downloadUrl, resp.StatusCode)
	}

	out, err := os.Create(localPath)
	if err != nil {
		return false, errorutils.CheckError(err)
	}
	defer out.Close()

	_, err = out.Write(body)
	if err != nil {
		return false, errorutils.CheckError(err)
	}

	if !params.IsSkipChecksum() {
		if err := dds.validateChecksum(downloadUrl, localPath); err != nil {
			log.Warn("Checksum validation failed for", localPath, ":", err)
		}
	}

	log.Info("Downloaded:", downloadUrl, "to", localPath)
	return true, nil
}

func (dds *DirectDownloadService) handleWildcardDownload(repo, pattern string, params *DirectDownloadParams) (int, int, error) {
	artifactoryUrl := strings.TrimSuffix((*dds.artDetails).GetUrl(), "/")

	dir := filepath.Dir(pattern)
	filePattern := filepath.Base(pattern)

	listUrl := fmt.Sprintf("%s/api/storage/%s/%s", artifactoryUrl, repo, dir)

	httpClientsDetails := (*dds.artDetails).CreateHttpClientDetails()
	resp, body, _, err := dds.client.SendGet(listUrl, true, &httpClientsDetails)
	if err != nil {
		return 0, 0, err
	}
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return 0, 0, errorutils.CheckErrorf("Failed to list directory %s: HTTP %d", listUrl, resp.StatusCode)
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
			if !dds.isExcluded(filepath.Join(dir, fileName), params.GetExclusions()) {
				filePath := filepath.Join(dir, fileName)
				success, err := dds.downloadSingleFile(repo, filePath, params)
				if err != nil {
					log.Error("Failed to download", filePath, ":", err)
					failCount++
				} else if success {
					downloadCount++
				} else {
					failCount++
				}
			}
		}
	}

	return downloadCount, failCount, nil
}

func (dds *DirectDownloadService) validateChecksum(downloadUrl, localPath string) error {
	artUrl := (*dds.artDetails).GetUrl()
	repoPath := strings.TrimPrefix(downloadUrl, artUrl)

	storageUrl := fmt.Sprintf("%sapi/storage/%s", artUrl, repoPath)

	httpClientsDetails := (*dds.artDetails).CreateHttpClientDetails()
	resp, body, _, err := dds.client.SendGet(storageUrl, true, &httpClientsDetails)
	if err != nil {
		return err
	}
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return errorutils.CheckErrorf("Failed to get checksum info: HTTP %d", resp.StatusCode)
	}

	var fileInfo utils.FileInfo
	if err := json.Unmarshal(body, &fileInfo); err != nil {
		return errorutils.CheckError(err)
	}

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
