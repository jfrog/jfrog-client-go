package utils

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jfrog/gofrog/crypto"
	"github.com/jfrog/gofrog/safeconvert"

	"github.com/jfrog/gofrog/parallel"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	ioutils "github.com/jfrog/jfrog-client-go/utils/io"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type supportedStatus int
type completionStatus string

const (
	minArtifactoryVersion = "7.82.2"

	// Supported status
	// Multipart upload support is not yet determined
	undetermined supportedStatus = iota
	// Multipart upload is supported
	multipartSupported
	// Multipart upload is not supported
	multipartNotSupported

	// Completion status
	parts             completionStatus = "PARTS"
	queued            completionStatus = "QUEUED"
	processing        completionStatus = "PROCESSING"
	finished          completionStatus = "FINISHED"
	retryableError    completionStatus = "RETRYABLE_ERROR"
	nonRetryableError completionStatus = "NON_RETRYABLE_ERROR"
	aborted           completionStatus = "ABORTED"

	// API constants
	uploadsApi = "/api/v1/uploads/"

	// Sizes and limits constants
	MaxMultipartUploadFileSize       = SizeTiB * 5
	DefaultUploadChunkSize     int64 = SizeMiB * 20

	// Retries and polling constants
	retriesInterval = time.Second * 5
	// A week of retries
	maxPollingRetries      = time.Hour * 168 / retriesInterval
	mergingLoggingInterval = time.Minute
)

var (
	errTooManyAttempts = errors.New("too many upload attempts failed")
	supportedMutex     sync.Mutex
)

type MultipartUpload struct {
	client             *jfroghttpclient.JfrogHttpClient
	httpClientsDetails *httputils.HttpClientDetails
	artifactoryUrl     string
	supportedStatus    supportedStatus
}

func NewMultipartUpload(client *jfroghttpclient.JfrogHttpClient, httpClientsDetails *httputils.HttpClientDetails, artifactoryUrl string) *MultipartUpload {
	return &MultipartUpload{client, httpClientsDetails, strings.TrimSuffix(artifactoryUrl, "/"), undetermined}
}

func (mu *MultipartUpload) IsSupported(serviceDetails auth.ServiceDetails) (supported bool, err error) {
	supportedMutex.Lock()
	defer supportedMutex.Unlock()
	if mu.supportedStatus != undetermined {
		// If the supported status was determined earlier, return true if multipart upload is supported or false if not
		return mu.supportedStatus == multipartSupported, nil
	}

	artifactoryVersion, err := serviceDetails.GetVersion()
	if err != nil {
		return
	}

	if versionErr := utils.ValidateMinimumVersion(utils.Artifactory, artifactoryVersion, minArtifactoryVersion); versionErr != nil {
		log.Debug("Multipart upload is not supported in versions below " + minArtifactoryVersion + ". Proceeding with regular upload...")
		mu.supportedStatus = multipartNotSupported
		return
	}

	url := fmt.Sprintf("%s%sconfig", mu.artifactoryUrl, uploadsApi)
	resp, body, _, err := mu.client.SendGet(url, true, mu.httpClientsDetails)
	if err != nil {
		return
	}
	log.Debug("Artifactory response:", string(body), resp.Status)
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return
	}

	var getConfigResponse getConfigResponse
	err = errorutils.CheckError(json.Unmarshal(body, &getConfigResponse))
	if getConfigResponse.Supported {
		mu.supportedStatus = multipartSupported
	} else {
		mu.supportedStatus = multipartNotSupported
	}
	return getConfigResponse.Supported, err
}

type getConfigResponse struct {
	Supported bool `json:"supported,omitempty"`
}

func (mu *MultipartUpload) UploadFileConcurrently(localPath, targetPath string, fileSize int64, sha1 string,
	progress ioutils.ProgressMgr, splitCount int, chunkSize int64) (checksumToken string, err error) {
	repoAndPath := strings.SplitN(targetPath, "/", 2)
	repoKey := repoAndPath[0]
	repoPath := repoAndPath[1]
	logMsgPrefix := fmt.Sprintf("[Multipart upload %s] ", repoPath)

	token, err := mu.createMultipartUpload(repoKey, repoPath, calculatePartSize(fileSize, 0, chunkSize))
	if err != nil {
		return
	}

	multipartUploadClient := &httputils.HttpClientDetails{
		AccessToken:           token,
		Transport:             mu.httpClientsDetails.Transport,
		DialTimeout:           mu.httpClientsDetails.DialTimeout,
		OverallRequestTimeout: mu.httpClientsDetails.OverallRequestTimeout,
	}

	var progressReader ioutils.Progress
	if progress != nil {
		progressReader = progress.NewProgressReader(fileSize, "Multipart upload", targetPath)
		progressId := progressReader.GetId()
		defer progress.RemoveProgress(progressId)
	}

	defer func() {
		if err == nil {
			log.Info(logMsgPrefix + "Upload completed successfully!")
		} else {
			err = errors.Join(err, mu.abort(logMsgPrefix, multipartUploadClient))
		}
	}()

	if err = mu.uploadPartsConcurrently(logMsgPrefix, fileSize, chunkSize, splitCount, localPath, progressReader, multipartUploadClient); err != nil {
		return
	}

	if sha1 == "" {
		var checksums map[crypto.Algorithm]string
		if checksums, err = crypto.GetFileChecksums(localPath); errorutils.CheckError(err) != nil {
			return
		}
		sha1 = checksums[crypto.SHA1]
	}

	if progress != nil {
		progressReader = progress.SetMergingState(progressReader.GetId(), false)
	}

	unsignedNumRetries, err := safeconvert.IntToUint(mu.client.GetHttpClient().GetRetries())
	if err != nil {
		return "", fmt.Errorf("failed to convert number of retries to uint64: %w", err)
	}
	log.Info(logMsgPrefix + "Starting parts merge...")
	// The total number of attempts is determined by the number of retries + 1
	return mu.completeAndPollForStatus(logMsgPrefix, unsignedNumRetries+1, sha1, multipartUploadClient, progressReader)
}

func (mu *MultipartUpload) uploadPartsConcurrently(logMsgPrefix string, fileSize, chunkSize int64, splitCount int, localPath string, progressReader ioutils.Progress, multipartUploadClient *httputils.HttpClientDetails) (err error) {
	numberOfParts := calculateNumberOfParts(fileSize, chunkSize)
	unsignedNumOfParts, err := safeconvert.Int64ToUint64(numberOfParts)
	if err != nil {
		return fmt.Errorf("failed to convert number of parts to uint64: %w", err)
	}
	unsignedNumRetries, err := safeconvert.Int64ToUint64(int64(mu.client.GetHttpClient().GetRetries()))
	if err != nil {
		return fmt.Errorf("failed to convert number of retries to uint64: %w", err)
	}
	log.Info(fmt.Sprintf("%sSplitting file to %d parts of %s each, using %d working threads for uploading...", logMsgPrefix, numberOfParts, ConvertIntToStorageSizeString(chunkSize), splitCount))
	producerConsumer := parallel.NewRunner(splitCount, uint(unsignedNumOfParts), false)

	wg := new(sync.WaitGroup)
	wg.Add(int(numberOfParts))
	attemptsAllowed := new(atomic.Uint64)
	attemptsAllowed.Add(unsignedNumOfParts * unsignedNumRetries)
	go func() {
		for i := 0; i < int(numberOfParts); i++ {
			if err = mu.produceUploadTask(producerConsumer, logMsgPrefix, localPath, fileSize, numberOfParts, int64(i), chunkSize, progressReader, multipartUploadClient, attemptsAllowed, wg); err != nil {
				return
			}
		}
	}()
	go func() {
		defer producerConsumer.Done()
		wg.Wait()
	}()
	producerConsumer.Run()
	if attemptsAllowed.Load() == 0 {
		return errorutils.CheckError(errTooManyAttempts)
	}
	return
}

func (mu *MultipartUpload) produceUploadTask(producerConsumer parallel.Runner, logMsgPrefix, localPath string, fileSize, numberOfParts, partId, chunkSize int64, progressReader ioutils.Progress, multipartUploadClient *httputils.HttpClientDetails, attemptsAllowed *atomic.Uint64, wg *sync.WaitGroup) (retErr error) {
	_, retErr = producerConsumer.AddTaskWithError(func(int) error {
		uploadErr := mu.uploadPart(logMsgPrefix, localPath, fileSize, partId, chunkSize, progressReader, multipartUploadClient)
		if uploadErr == nil {
			log.Info(fmt.Sprintf("%sCompleted uploading part %d/%d", logMsgPrefix, partId+1, numberOfParts))
			wg.Done()
		}
		return uploadErr
	}, func(uploadErr error) {
		if attemptsAllowed.Load() == 0 {
			wg.Done()
			return
		}
		log.Warn(fmt.Sprintf("%sPart %d/%d - %s", logMsgPrefix, partId+1, numberOfParts, uploadErr.Error()))
		attemptsAllowed.Add(^uint64(0))

		// Sleep before trying again
		time.Sleep(retriesInterval)
		if err := mu.produceUploadTask(producerConsumer, logMsgPrefix, localPath, fileSize, numberOfParts, partId, chunkSize, progressReader, multipartUploadClient, attemptsAllowed, wg); err != nil {
			retErr = err
		}
	})
	return
}

func (mu *MultipartUpload) uploadPart(logMsgPrefix, localPath string, fileSize, partId, chunkSize int64, progressReader ioutils.Progress, multipartUploadClient *httputils.HttpClientDetails) (err error) {
	file, err := os.Open(localPath)
	if err != nil {
		return errorutils.CheckError(err)
	}
	defer func() {
		if file != nil {
			err = errors.Join(err, errorutils.CheckError(file.Close()))
		}
	}()
	if _, err = file.Seek(partId*chunkSize, io.SeekStart); err != nil {
		return errorutils.CheckError(err)
	}
	partSize := calculatePartSize(fileSize, partId, chunkSize)

	limitReader := io.LimitReader(file, partSize)
	limitReader = bufio.NewReader(limitReader)
	if progressReader != nil {
		limitReader = progressReader.ActionWithProgress(limitReader)
	}

	urlPart, err := mu.generateUrlPart(logMsgPrefix, partId, multipartUploadClient)
	if err != nil {
		return
	}

	resp, body, err := mu.client.GetHttpClient().UploadFileFromReader(limitReader, urlPart, httputils.HttpClientDetails{}, partSize)
	if err != nil {
		return
	}
	log.Debug("Artifactory response:", string(body), resp.Status)
	return errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK)
}

func (mu *MultipartUpload) createMultipartUpload(repoKey, repoPath string, partSize int64) (token string, err error) {
	requestUrl := fmt.Sprintf("%s%screate?repoKey=%s&repoPath=%s&partSizeMB=%d", mu.artifactoryUrl, uploadsApi, url.QueryEscape(repoKey), url.QueryEscape(repoPath), partSize/SizeMiB)
	resp, body, err := mu.client.SendPost(requestUrl, []byte{}, mu.httpClientsDetails)
	if err != nil {
		return
	}
	// We don't log the response body because it includes credentials

	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return
	}

	var createMultipartUploadResponse createMultipartUploadResponse
	err = json.Unmarshal(body, &createMultipartUploadResponse)
	return createMultipartUploadResponse.Token, err
}

type createMultipartUploadResponse struct {
	Token string `json:"token,omitempty"`
}

func (mu *MultipartUpload) generateUrlPart(logMsgPrefix string, partNumber int64, multipartUploadClient *httputils.HttpClientDetails) (partUrl string, err error) {
	url := fmt.Sprintf("%s%surlPart?partNumber=%d", mu.artifactoryUrl, uploadsApi, partNumber+1)
	resp, body, err := mu.client.GetHttpClient().SendPost(url, []byte{}, *multipartUploadClient, logMsgPrefix)
	if err != nil {
		return "", err
	}
	// We don't log the response body because it includes credentials

	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return
	}
	var urlPartResponse urlPartResponse
	err = json.Unmarshal(body, &urlPartResponse)
	return urlPartResponse.Url, errorutils.CheckError(err)
}

type urlPartResponse struct {
	Url string `json:"url,omitempty"`
}

func (mu *MultipartUpload) completeAndPollForStatus(logMsgPrefix string, completionAttemptsLeft uint, sha1 string, multipartUploadClient *httputils.HttpClientDetails, progressReader ioutils.Progress) (checksumToken string, err error) {
	err = mu.completeMultipartUpload(logMsgPrefix, sha1, multipartUploadClient)
	if err != nil {
		return
	}

	checksumToken, err = mu.pollCompletionStatus(logMsgPrefix, completionAttemptsLeft, sha1, multipartUploadClient, progressReader)
	return
}

func (mu *MultipartUpload) pollCompletionStatus(logMsgPrefix string, completionAttemptsLeft uint, sha1 string, multipartUploadClient *httputils.HttpClientDetails, progressReader ioutils.Progress) (checksumToken string, err error) {
	multipartUploadClientWithNodeId := multipartUploadClient.Clone()

	lastMergeLog := time.Now()
	pollingExecutor := &utils.RetryExecutor{
		MaxRetries:               int(maxPollingRetries),
		RetriesIntervalMilliSecs: int(retriesInterval.Milliseconds()),
		LogMsgPrefix:             logMsgPrefix,
		ExecutionHandler: func() (shouldRetry bool, err error) {
			// Get completion status
			status, err := mu.status(logMsgPrefix, multipartUploadClientWithNodeId)
			if err != nil {
				return false, err
			}

			// Parse status
			shouldRetry, shouldRerunComplete, err := parseMultipartUploadStatus(status)
			if err != nil {
				return false, err
			}

			// Rerun complete if needed
			if shouldRerunComplete {
				if completionAttemptsLeft == 0 {
					return false, errorutils.CheckErrorf("multipart upload failed after %d attempts", mu.client.GetHttpClient().GetRetries())
				}
				checksumToken, err = mu.completeAndPollForStatus(logMsgPrefix, completionAttemptsLeft-1, sha1, multipartUploadClient, progressReader)
			}

			// Log status
			if status.Progress != nil {
				if progressReader != nil {
					progressReader.SetProgress(int64(*status.Progress))
				}
				if time.Since(lastMergeLog) > mergingLoggingInterval {
					log.Info(fmt.Sprintf("%sMerging progress: %d%%", logMsgPrefix, *status.Progress))
					lastMergeLog = time.Now()
				}
			}
			checksumToken = status.ChecksumToken
			return
		},
	}
	return checksumToken, pollingExecutor.Execute()
}

func (mu *MultipartUpload) completeMultipartUpload(logMsgPrefix, sha1 string, multipartUploadClient *httputils.HttpClientDetails) error {
	url := fmt.Sprintf("%s%scomplete?sha1=%s", mu.artifactoryUrl, uploadsApi, sha1)
	resp, body, err := mu.client.GetHttpClient().SendPost(url, []byte{}, *multipartUploadClient, logMsgPrefix)
	if err != nil {
		return err
	}
	log.Debug("Artifactory response:", string(body), resp.Status)
	return errorutils.CheckResponseStatusWithBody(resp, body, http.StatusAccepted)
}

func (mu *MultipartUpload) status(logMsgPrefix string, multipartUploadClientWithNodeId *httputils.HttpClientDetails) (status statusResponse, err error) {
	url := fmt.Sprintf("%s%sstatus", mu.artifactoryUrl, uploadsApi)
	resp, body, err := mu.client.GetHttpClient().SendPost(url, []byte{}, *multipartUploadClientWithNodeId, logMsgPrefix)
	// If the Artifactory node returns a "Service unavailable" error (status 503), attempt to retry the upload completion process on a different node.
	if resp != nil && resp.StatusCode == http.StatusServiceUnavailable {
		unavailableNodeErr := fmt.Sprintf("%sArtifactory is unavailable.", logMsgPrefix)
		return statusResponse{Status: retryableError, Error: unavailableNodeErr}, nil
	}
	if err != nil {
		return
	}
	log.Debug("Artifactory response:", string(body), resp.Status)
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return
	}
	err = errorutils.CheckError(json.Unmarshal(body, &status))
	return
}

type statusResponse struct {
	Status        completionStatus `json:"status,omitempty"`
	Error         string           `json:"error,omitempty"`
	Progress      *int             `json:"progress,omitempty"`
	ChecksumToken string           `json:"checksumToken,omitempty"`
}

func (mu *MultipartUpload) abort(logMsgPrefix string, multipartUploadClient *httputils.HttpClientDetails) (err error) {
	log.Info("Aborting multipart upload...")
	url := fmt.Sprintf("%s%sabort", mu.artifactoryUrl, uploadsApi)
	resp, body, err := mu.client.GetHttpClient().SendPost(url, []byte{}, *multipartUploadClient, logMsgPrefix)
	if err != nil {
		return
	}
	log.Debug("Artifactory response:", string(body), resp.Status)
	return errorutils.CheckResponseStatusWithBody(resp, body, http.StatusNoContent)
}

// Calculates the part size based on the file size, the part number and the requested chunk size.
// fileSize - the file size
// partNumber - the current part number
// requestedChunkSize - chunk size requested by the user, or default.
func calculatePartSize(fileSize, partNumber, requestedChunkSize int64) int64 {
	partOffset := partNumber * requestedChunkSize
	if partOffset+requestedChunkSize > fileSize {
		return fileSize - partOffset
	}
	return requestedChunkSize
}

// Calculates the number of parts based on the file size and the requested chunks size.
// fileSize - the file size
func calculateNumberOfParts(fileSize, chunkSize int64) int64 {
	return (fileSize + chunkSize - 1) / chunkSize
}

func parseMultipartUploadStatus(status statusResponse) (shouldKeepPolling, shouldRerunComplete bool, err error) {
	switch status.Status {
	case queued, processing:
		// File merging had not yet been completed - keep polling
		return true, false, nil
	case retryableError:
		// Retryable error was received - stop polling and rerun the /complete API again
		log.Warn(fmt.Printf("received error upon multipart upload completion process: '%s', retrying...", status.Error))
		return false, true, nil
	case finished, aborted:
		// Upload finished or aborted
		return false, false, nil
	case nonRetryableError:
		// Fatal error occurred - stop the entire process
		return false, false, errorutils.CheckErrorf("received non retryable error upon multipart upload completion process: '%s'", status.Error)
	default:
		// Unexpected status - stop the entire process
		return false, false, errorutils.CheckErrorf("received unexpected status upon multipart upload completion process: '%s', error: '%s'", status.Status, status.Error)
	}
}
