package utils

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/jfrog/jfrog-client-go/utils/tests"
	"github.com/stretchr/testify/assert"
)

const (
	localPath  = "localPath"
	repoKey    = "repoKey"
	repoPath   = "repoPath"
	partSize   = SizeGiB
	partSizeMB = 1024
	partNumber = 2
	splitCount = 3
	token      = "token"
	// jfrog-ignore - test url
	partUrl       = "http://dummy-url-part"
	sha1          = "sha1"
	checksumToken = "checksumToken"
)

func TestIsSupported(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check method
		assert.Equal(t, http.MethodGet, r.Method)

		// Check URL
		assert.Equal(t, "/api/v1/uploads/config", r.URL.Path)

		// Send response 200 OK
		w.WriteHeader(http.StatusOK)
		response, err := json.Marshal(getConfigResponse{Supported: true})
		assert.NoError(t, err)
		_, err = w.Write(response)
		assert.NoError(t, err)
	})

	// Create mock multipart upload with server
	multipartUpload, cleanUp := createMockMultipartUpload(t, handler)
	defer cleanUp()

	// Create Artifactory service details
	rtDetails := &dummyArtifactoryServiceDetails{version: minArtifactoryVersion}

	// Execute IsSupported
	supported, err := multipartUpload.IsSupported(rtDetails)
	assert.NoError(t, err)
	assert.True(t, supported)
}

func TestUnsupportedVersion(t *testing.T) {
	// Create Artifactory service details with unsupported Artifactory version
	rtDetails := &dummyArtifactoryServiceDetails{version: "6.0.0"}

	// Create mock multipart upload with server
	client, err := jfroghttpclient.JfrogClientBuilder().Build()
	assert.NoError(t, err)

	// Execute IsSupported
	supported, err := NewMultipartUpload(client, nil, "").IsSupported(rtDetails)
	assert.NoError(t, err)
	assert.False(t, supported)
}

func TestCreateMultipartUpload(t *testing.T) {
	testCases := []struct {
		name                string
		repoKey             string
		repoPath            string
		urlExpectedRepoKey  string
		urlExpectedRepoPath string
	}{
		{"Single word names", repoKey, repoPath, repoKey, repoPath},
		{"Spaced names", "repo with space", "path with space", "repo+with+space", "path+with+space"},
		{"Names contains _", "repo_name", "path_name", "repo_name", "path_name"},
		{"Names contains %", "repo%name", "path%name", "repo%25name", "path%25name"},
		{"Names contains &", "repo&name", "path&name", "repo%26name", "path%26name"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Check method
				assert.Equal(t, http.MethodPost, r.Method)

				// Check URL
				assert.Equal(t, "/api/v1/uploads/create", r.URL.Path)
				assert.Equal(t, fmt.Sprintf("repoKey=%s&repoPath=%s&partSizeMB=%d", testCase.urlExpectedRepoKey, testCase.urlExpectedRepoPath, partSizeMB), r.URL.RawQuery)

				// Send response 200 OK
				w.WriteHeader(http.StatusOK)
				response, err := json.Marshal(createMultipartUploadResponse{Token: token})
				assert.NoError(t, err)
				_, err = w.Write(response)
				assert.NoError(t, err)
			})

			// Create mock multipart upload with server
			multipartUpload, cleanUp := createMockMultipartUpload(t, handler)
			defer cleanUp()

			// Execute CreateMultipartUpload
			actualToken, err := multipartUpload.createMultipartUpload(testCase.repoKey, testCase.repoPath, partSize)
			assert.NoError(t, err)
			assert.Equal(t, token, actualToken)
		})
	}
}

func TestUploadPartsConcurrentlyTooManyAttempts(t *testing.T) {
	// Create temp file
	tempFile, cleanUp := createTempFile(t)
	defer cleanUp()

	// Write something to the file
	buf := make([]byte, DefaultUploadChunkSize*3)
	_, err := rand.Read(buf)
	assert.NoError(t, err)
	_, err = tempFile.Write(buf)
	assert.NoError(t, err)

	var multipartUpload *MultipartUpload
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		// Generate part URL for upload
		case http.MethodPost:
			// Check URL
			assert.Equal(t, "/api/v1/uploads/urlPart", r.URL.Path)

			// Send response 200 OK
			w.WriteHeader(http.StatusOK)
			response, unmarshalErr := json.Marshal(urlPartResponse{Url: multipartUpload.artifactoryUrl})
			assert.NoError(t, unmarshalErr)
			_, err = w.Write(response)
			assert.NoError(t, err)
		// Fail the upload to trigger retry
		case http.MethodPut:
			assert.Equal(t, "/", r.URL.Path)

			// Send response 502 OK
			w.WriteHeader(http.StatusBadGateway)
		default:
			assert.Fail(t, "unexpected method "+r.Method)
		}
	})

	// Create mock multipart upload with server
	multipartUpload, cleanUp = createMockMultipartUpload(t, handler)
	defer cleanUp()

	// Execute uploadPartsConcurrently
	fileSize := int64(len(buf))
	err = multipartUpload.uploadPartsConcurrently("", fileSize, DefaultUploadChunkSize, splitCount, tempFile.Name(), nil, &httputils.HttpClientDetails{})
	assert.ErrorIs(t, err, errTooManyAttempts)
}

func TestGenerateUrlPart(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check method
		assert.Equal(t, http.MethodPost, r.Method)

		// Check URL
		assert.Equal(t, "/api/v1/uploads/urlPart", r.URL.Path)
		assert.Equal(t, fmt.Sprintf("partNumber=%d", partNumber+1), r.URL.RawQuery)

		// Send response 200 OK
		w.WriteHeader(http.StatusOK)
		response, err := json.Marshal(urlPartResponse{Url: partUrl})
		assert.NoError(t, err)
		_, err = w.Write(response)
		assert.NoError(t, err)
	})

	// Create mock multipart upload with server
	multipartUpload, cleanUp := createMockMultipartUpload(t, handler)
	defer cleanUp()

	// Execute GenerateUrlPart
	actualPartUrl, err := multipartUpload.generateUrlPart("", partNumber, &httputils.HttpClientDetails{})
	assert.NoError(t, err)
	assert.Equal(t, partUrl, actualPartUrl)
}

func TestCompleteMultipartUpload(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check method
		assert.Equal(t, http.MethodPost, r.Method)

		// Check URL
		assert.Equal(t, "/api/v1/uploads/complete", r.URL.Path)
		assert.Equal(t, fmt.Sprintf("sha1=%s", sha1), r.URL.RawQuery)

		// Send response 202 Accepted
		w.WriteHeader(http.StatusAccepted)
	})

	// Create mock multipart upload with server
	multipartUpload, cleanUp := createMockMultipartUpload(t, handler)
	defer cleanUp()

	// Execute completeMultipartUpload
	err := multipartUpload.completeMultipartUpload("", sha1, &httputils.HttpClientDetails{})
	assert.NoError(t, err)
}

func TestStatus(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check method
		assert.Equal(t, http.MethodPost, r.Method)

		// Check URL
		assert.Equal(t, "/api/v1/uploads/status", r.URL.Path)

		// Send response 200 OK
		w.WriteHeader(http.StatusOK)
		response, err := json.Marshal(statusResponse{Status: finished, Progress: utils.Pointer(100), ChecksumToken: checksumToken})
		assert.NoError(t, err)
		_, err = w.Write(response)
		assert.NoError(t, err)
	})

	// Create mock multipart upload with server
	multipartUpload, cleanUp := createMockMultipartUpload(t, handler)
	defer cleanUp()

	// Execute status
	status, err := multipartUpload.status("", &httputils.HttpClientDetails{})
	assert.NoError(t, err)
	assert.Equal(t, statusResponse{Status: finished, Progress: utils.Pointer(100), ChecksumToken: checksumToken}, status)
}

func TestStatusServiceUnavailable(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check method
		assert.Equal(t, http.MethodPost, r.Method)

		// Check URL
		assert.Equal(t, "/api/v1/uploads/status", r.URL.Path)

		// Send response 503 Service unavailable
		w.WriteHeader(http.StatusServiceUnavailable)
		_, err := w.Write([]byte("Service unavailable"))
		assert.NoError(t, err)
	})

	// Create mock multipart upload with server
	multipartUpload, cleanUp := createMockMultipartUpload(t, handler)
	defer cleanUp()

	// Execute status
	status, err := multipartUpload.status("", &httputils.HttpClientDetails{})
	assert.NoError(t, err)
	assert.Equal(t, statusResponse{Status: retryableError, Error: "Artifactory is unavailable."}, status)
}

// The checksum token returned after a retried /complete must not be overwritten by the earlier RETRYABLE_ERROR status.
func TestCompleteAndPollForStatus(t *testing.T) {
	retryable := statusResponse{Status: retryableError, Error: "Injected S3 SlowDown", Progress: utils.Pointer(0)}
	done := statusResponse{Status: finished, Progress: utils.Pointer(100), ChecksumToken: checksumToken}
	nonRetryable := statusResponse{Status: nonRetryableError, Error: "fatal"}

	testCases := []struct {
		name                   string
		attemptsLeft           uint
		statuses               []statusResponse
		expectedToken          string
		expectedErr            string
		expectedCompleteCalls  int
		expectedStatusRequests int
	}{
		{name: "finished immediately", attemptsLeft: 1, statuses: []statusResponse{done}, expectedToken: checksumToken, expectedCompleteCalls: 1, expectedStatusRequests: 1},
		{name: "retry then finished", attemptsLeft: 1, statuses: []statusResponse{retryable, done}, expectedToken: checksumToken, expectedCompleteCalls: 2, expectedStatusRequests: 2},
		{name: "two retries then finished", attemptsLeft: 2, statuses: []statusResponse{retryable, retryable, done}, expectedToken: checksumToken, expectedCompleteCalls: 3, expectedStatusRequests: 3},
		{name: "retries exhausted", attemptsLeft: 0, statuses: []statusResponse{retryable}, expectedErr: "multipart upload failed after", expectedCompleteCalls: 1, expectedStatusRequests: 1},
		{name: "retries exhausted after a retry", attemptsLeft: 1, statuses: []statusResponse{retryable, retryable}, expectedErr: "multipart upload failed after", expectedCompleteCalls: 2, expectedStatusRequests: 2},
		{name: "retry then non retryable error", attemptsLeft: 1, statuses: []statusResponse{retryable, nonRetryable}, expectedErr: "received non retryable error", expectedCompleteCalls: 2, expectedStatusRequests: 2},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			completeCalls := 0
			statusRequests := 0
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				switch r.URL.Path {
				case "/api/v1/uploads/complete":
					completeCalls++
					w.WriteHeader(http.StatusAccepted)
				case "/api/v1/uploads/status":
					if !assert.Less(t, statusRequests, len(testCase.statuses), "unexpected extra status request") {
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					response, err := json.Marshal(testCase.statuses[statusRequests])
					assert.NoError(t, err)
					statusRequests++
					w.WriteHeader(http.StatusOK)
					_, err = w.Write(response)
					assert.NoError(t, err)
				default:
					assert.Fail(t, "unexpected request", r.URL.Path)
				}
			})

			multipartUpload, cleanUp := createMockMultipartUpload(t, handler)
			defer cleanUp()

			actualToken, err := multipartUpload.completeAndPollForStatus("", testCase.attemptsLeft, sha1, &httputils.HttpClientDetails{}, nil)
			if testCase.expectedErr != "" {
				assert.ErrorContains(t, err, testCase.expectedErr)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, testCase.expectedToken, actualToken)
			assert.Equal(t, testCase.expectedCompleteCalls, completeCalls)
			assert.Equal(t, testCase.expectedStatusRequests, statusRequests)
		})
	}
}

func TestAbort(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check method
		assert.Equal(t, http.MethodPost, r.Method)

		// Check URL
		assert.Equal(t, "/api/v1/uploads/abort", r.URL.Path)

		// Send response 204 No Content
		w.WriteHeader(http.StatusNoContent)
	})

	// Create mock multipart upload with server
	multipartUpload, cleanUp := createMockMultipartUpload(t, handler)
	defer cleanUp()

	// Execute status
	clientDetails := &httputils.HttpClientDetails{}
	err := multipartUpload.abort("", clientDetails)
	assert.NoError(t, err)
}

var calculatePartSizeProvider = []struct {
	fileSize         int64
	partNumber       int64
	expectedPartSize int64
}{
	{DefaultUploadChunkSize - 1, 0, DefaultUploadChunkSize - 1},
	{DefaultUploadChunkSize, 0, DefaultUploadChunkSize},
	{DefaultUploadChunkSize + 1, 0, DefaultUploadChunkSize},

	{DefaultUploadChunkSize*2 - 1, 1, DefaultUploadChunkSize - 1},
	{DefaultUploadChunkSize * 2, 1, DefaultUploadChunkSize},
	{DefaultUploadChunkSize*2 + 1, 1, DefaultUploadChunkSize},
}

func TestCalculatePartSize(t *testing.T) {
	for _, testCase := range calculatePartSizeProvider {
		t.Run(fmt.Sprintf("fileSize: %d partNumber: %d", testCase.fileSize, testCase.partNumber), func(t *testing.T) {
			assert.Equal(t, testCase.expectedPartSize, calculatePartSize(testCase.fileSize, testCase.partNumber, DefaultUploadChunkSize))
		})
	}
}

var calculateNumberOfPartsProvider = []struct {
	fileSize              int64
	expectedNumberOfParts int64
}{
	{0, 0},
	{1, 1},
	{DefaultUploadChunkSize - 1, 1},
	{DefaultUploadChunkSize, 1},
	{DefaultUploadChunkSize + 1, 2},

	{DefaultUploadChunkSize*2 - 1, 2},
	{DefaultUploadChunkSize * 2, 2},
	{DefaultUploadChunkSize*2 + 1, 3},
}

func TestCalculateNumberOfParts(t *testing.T) {
	for _, testCase := range calculateNumberOfPartsProvider {
		t.Run(fmt.Sprintf("fileSize: %d", testCase.fileSize), func(t *testing.T) {
			assert.Equal(t, testCase.expectedNumberOfParts, calculateNumberOfParts(testCase.fileSize, DefaultUploadChunkSize))
		})
	}
}

var parseMultipartUploadStatusProvider = []struct {
	status              completionStatus
	shouldKeepPolling   bool
	shouldRerunComplete bool
	expectedError       string
}{
	{queued, true, false, ""},
	{processing, true, false, ""},
	{parts, false, false, "received unexpected status upon multipart upload completion process: 'PARTS', error: 'Some error'"},
	{finished, false, false, ""},
	{aborted, false, false, ""},
	{retryableError, false, true, ""},
	{nonRetryableError, false, false, "received non retryable error upon multipart upload completion process: 'Some error'"},
}

func TestParseMultipartUploadStatus(t *testing.T) {
	previousLog := tests.RedirectLogOutputToNil()
	defer func() {
		log.SetLogger(previousLog)
	}()

	for _, testCase := range parseMultipartUploadStatusProvider {
		t.Run(string(testCase.status), func(t *testing.T) {

			shouldKeepPolling, shouldRerunComplete, err := parseMultipartUploadStatus(statusResponse{Status: testCase.status, Error: "Some error"})
			if testCase.expectedError != "" {
				assert.EqualError(t, err, testCase.expectedError)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, testCase.shouldKeepPolling, shouldKeepPolling)
			assert.Equal(t, testCase.shouldRerunComplete, shouldRerunComplete)
		})
	}
}

func createTempFile(t *testing.T) (tempFile *os.File, cleanUp func()) {
	// Create a temporary file
	tempFile, err := fileutils.CreateTempFile()
	assert.NoError(t, err)
	cleanUp = func() {
		assert.NoError(t, tempFile.Close())
		assert.NoError(t, fileutils.RemovePath(localPath))
	}
	return
}

func createMockMultipartUpload(t *testing.T, handler http.Handler) (multipartUpload *MultipartUpload, cleanUp func()) {
	ts := httptest.NewServer(handler)
	cleanUp = ts.Close

	client, err := jfroghttpclient.JfrogClientBuilder().Build()
	assert.NoError(t, err)
	multipartUpload = NewMultipartUpload(client, &httputils.HttpClientDetails{}, ts.URL)
	return
}

type dummyArtifactoryServiceDetails struct {
	auth.CommonConfigFields
	version string
}

func (dasd *dummyArtifactoryServiceDetails) GetVersion() (string, error) {
	return dasd.version, nil
}
