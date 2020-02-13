package httpclient

import (
	"errors"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/httpclient"
	ioutils "github.com/jfrog/jfrog-client-go/utils/io"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"io"
	"net/http"
	"net/url"
)

type ArtifactoryHttpClient struct {
	httpClient *httpclient.HttpClient
	ArtDetails *auth.CommonDetails
}

func (rtc *ArtifactoryHttpClient) SendGet(url string, followRedirect bool, httpClientsDetails *httputils.HttpClientDetails) (resp *http.Response, respBody []byte, redirectUrl string, err error) {
	isNewToken := false
	for i := 0; i < 2; i++ {
		resp, respBody, redirectUrl, err = rtc.httpClient.SendGet(url, followRedirect, *httpClientsDetails)
		if err != nil {
			return
		}
		// Check if token expired, if so obtain a new one
		isNewToken, err = (*rtc.ArtDetails).HandleTokenExpiry(resp.StatusCode, httpClientsDetails)
		// Return if no new token was acquired or an error occurred
		if !isNewToken || err != nil {
			return
		}
	}
	err = errors.New("failed to obtain a new authentication token after one has expired; " + resp.Status)
	return
}

func (rtc *ArtifactoryHttpClient) SendPost(url string, content []byte, httpClientsDetails *httputils.HttpClientDetails) (resp *http.Response, body []byte, err error) {
	isNewToken := false
	for i := 0; i < 2; i++ {
		resp, body, err = rtc.httpClient.SendPost(url, content, *httpClientsDetails)
		if err != nil {
			return
		}
		// Check if token expired, if so obtain a new one
		isNewToken, err = (*rtc.ArtDetails).HandleTokenExpiry(resp.StatusCode, httpClientsDetails)
		// Return if no new token was acquired or an error occurred
		if !isNewToken || err != nil {
			return
		}
	}
	err = errors.New("failed to obtain a new authentication token after one has expired; " + resp.Status)
	return
}

func (rtc *ArtifactoryHttpClient) SendPostForm(url string, data url.Values, httpClientsDetails *httputils.HttpClientDetails) (resp *http.Response, body []byte, err error) {
	httpClientsDetails.Headers["Content-Type"] = "application/x-www-form-urlencoded"
	return rtc.SendPost(url, []byte(data.Encode()), httpClientsDetails)
}

func (rtc *ArtifactoryHttpClient) SendPatch(url string, content []byte, httpClientsDetails *httputils.HttpClientDetails) (resp *http.Response, body []byte, err error) {
	isNewToken := false
	for i := 0; i < 2; i++ {
		resp, body, err = rtc.httpClient.SendPatch(url, content, *httpClientsDetails)
		if err != nil {
			return
		}
		// Check if token expired, if so obtain a new one
		isNewToken, err = (*rtc.ArtDetails).HandleTokenExpiry(resp.StatusCode, httpClientsDetails)
		// Return if no new token was acquired or an error occurred
		if !isNewToken || err != nil {
			return
		}
	}
	err = errors.New("failed to obtain a new authentication token after one has expired; " + resp.Status)
	return
}

func (rtc *ArtifactoryHttpClient) SendDelete(url string, content []byte, httpClientsDetails *httputils.HttpClientDetails) (resp *http.Response, body []byte, err error) {
	isNewToken := false
	for i := 0; i < 2; i++ {
		resp, body, err = rtc.httpClient.SendDelete(url, content, *httpClientsDetails)
		if err != nil {
			return
		}
		// Check if token expired, if so obtain a new one
		isNewToken, err = (*rtc.ArtDetails).HandleTokenExpiry(resp.StatusCode, httpClientsDetails)
		// Return if no new token was acquired or an error occurred
		if !isNewToken || err != nil {
			return
		}
	}
	err = errors.New("failed to obtain a new authentication token after one has expired; " + resp.Status)
	return
}

func (rtc *ArtifactoryHttpClient) SendHead(url string, httpClientsDetails *httputils.HttpClientDetails) (resp *http.Response, body []byte, err error) {
	isNewToken := false
	for i := 0; i < 2; i++ {
		resp, body, err = rtc.httpClient.SendHead(url, *httpClientsDetails)
		if err != nil {
			return
		}
		// Check if token expired, if so obtain a new one
		isNewToken, err = (*rtc.ArtDetails).HandleTokenExpiry(resp.StatusCode, httpClientsDetails)
		// Return if no new token was acquired or an error occurred
		if !isNewToken || err != nil {
			return
		}
	}
	err = errors.New("failed to obtain a new authentication token after one has expired; " + resp.Status)
	return
}

func (rtc *ArtifactoryHttpClient) SendPut(url string, content []byte, httpClientsDetails *httputils.HttpClientDetails) (resp *http.Response, body []byte, err error) {
	isNewToken := false
	for i := 0; i < 2; i++ {
		resp, body, err = rtc.httpClient.SendPut(url, content, *httpClientsDetails)
		if err != nil {
			return
		}
		// Check if token expired, if so obtain a new one
		isNewToken, err = (*rtc.ArtDetails).HandleTokenExpiry(resp.StatusCode, httpClientsDetails)
		// Return if no new token was acquired or an error occurred
		if !isNewToken || err != nil {
			return
		}
	}
	err = errors.New("failed to obtain a new authentication token after one has expired; " + resp.Status)
	return
}

func (rtc *ArtifactoryHttpClient) Send(method string, url string, content []byte, followRedirect bool, closeBody bool,
	httpClientsDetails *httputils.HttpClientDetails) (resp *http.Response, respBody []byte, redirectUrl string, err error) {
	isNewToken := false
	for i := 0; i < 2; i++ {
		resp, respBody, redirectUrl, err = rtc.httpClient.Send(method, url, content, followRedirect, closeBody, *httpClientsDetails)
		if err != nil {
			return
		}
		// Check if token expired, if so obtain a new one
		isNewToken, err = (*rtc.ArtDetails).HandleTokenExpiry(resp.StatusCode, httpClientsDetails)
		// Return if no new token was acquired or an error occurred
		if !isNewToken || err != nil {
			return
		}
	}
	err = errors.New("failed to obtain a new authentication token after one has expired; " + resp.Status)
	return
}

func (rtc *ArtifactoryHttpClient) UploadFile(localPath, url, logMsgPrefix string,
	httpClientsDetails *httputils.HttpClientDetails, retries int, progress ioutils.Progress) (resp *http.Response, body []byte, err error) {
	isNewToken := false
	for i := 0; i < 2; i++ {
		resp, body, err = rtc.httpClient.UploadFile(localPath, url, logMsgPrefix, *httpClientsDetails, retries, progress)
		if err != nil {
			return
		}
		// Check if token expired, if so obtain a new one
		isNewToken, err = (*rtc.ArtDetails).HandleTokenExpiry(resp.StatusCode, httpClientsDetails)
		// Return if no new token was acquired or an error occurred
		if !isNewToken || err != nil {
			return
		}
	}
	err = errors.New("failed to obtain a new authentication token after one has expired; " + resp.Status)
	return
}

func (rtc *ArtifactoryHttpClient) ReadRemoteFile(downloadPath string, httpClientsDetails *httputils.HttpClientDetails) (ioReaderCloser io.ReadCloser, resp *http.Response, err error) {
	isNewToken := false
	for i := 0; i < 2; i++ {
		ioReaderCloser, resp, err = rtc.httpClient.ReadRemoteFile(downloadPath, *httpClientsDetails)
		if err != nil {
			return
		}
		// Check if token expired, if so obtain a new one
		isNewToken, err = (*rtc.ArtDetails).HandleTokenExpiry(resp.StatusCode, httpClientsDetails)
		// Return if no new token was acquired or an error occurred
		if !isNewToken || err != nil {
			return
		}
	}
	err = errors.New("failed to obtain a new authentication token after one has expired; " + resp.Status)
	return
}

func (rtc *ArtifactoryHttpClient) DownloadFileWithProgress(downloadFileDetails *httpclient.DownloadFileDetails, logMsgPrefix string,
	httpClientsDetails *httputils.HttpClientDetails, retries int, isExplode bool, progress ioutils.Progress) (resp *http.Response, err error) {
	isNewToken := false
	for i := 0; i < 2; i++ {
		resp, err = rtc.httpClient.DownloadFileWithProgress(downloadFileDetails, logMsgPrefix, *httpClientsDetails,
			retries, isExplode, progress)
		if err != nil {
			return
		}
		// Check if token expired, if so obtain a new one
		isNewToken, err = (*rtc.ArtDetails).HandleTokenExpiry(resp.StatusCode, httpClientsDetails)
		// Return if no new token was acquired or an error occurred
		if !isNewToken || err != nil {
			return
		}
	}
	err = errors.New("failed to obtain a new authentication token after one has expired; " + resp.Status)
	return
}

func (rtc *ArtifactoryHttpClient) DownloadFile(downloadFileDetails *httpclient.DownloadFileDetails, logMsgPrefix string,
	httpClientsDetails *httputils.HttpClientDetails, retries int, isExplode bool) (resp *http.Response, err error) {
	return rtc.DownloadFileWithProgress(downloadFileDetails, logMsgPrefix, httpClientsDetails, retries, isExplode, nil)
}

func (rtc *ArtifactoryHttpClient) DownloadFileConcurrently(flags httpclient.ConcurrentDownloadFlags,
	logMsgPrefix string, httpClientsDetails *httputils.HttpClientDetails, progress ioutils.Progress) (resp *http.Response, err error) {
	isNewToken := false
	for i := 0; i < 2; i++ {
		resp, err = rtc.httpClient.DownloadFileConcurrently(flags, logMsgPrefix, *httpClientsDetails, progress)
		if err != nil {
			return
		}
		// Check if token expired, if so obtain a new one
		isNewToken, err = (*rtc.ArtDetails).HandleTokenExpiry(resp.StatusCode, httpClientsDetails)
		// Return if no new token was acquired or an error occurred
		if !isNewToken || err != nil {
			return
		}
	}
	err = errors.New("failed to obtain a new authentication token after one has expired; " + resp.Status)
	return
}

func (rtc *ArtifactoryHttpClient) IsAcceptRanges(downloadUrl string, httpClientsDetails *httputils.HttpClientDetails) (isAcceptRanges bool, resp *http.Response, err error) {
	isNewToken := false
	for i := 0; i < 2; i++ {
		isAcceptRanges, resp, err = rtc.httpClient.IsAcceptRanges(downloadUrl, *httpClientsDetails)
		if err != nil {
			return
		}
		// Check if token expired, if so obtain a new one
		isNewToken, err = (*rtc.ArtDetails).HandleTokenExpiry(resp.StatusCode, httpClientsDetails)
		// Return if no new token was acquired or an error occurred
		if !isNewToken || err != nil {
			return
		}
	}
	err = errors.New("failed to obtain a new authentication token after one has expired; " + resp.Status)
	return
}
