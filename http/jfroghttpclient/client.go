package jfroghttpclient

import (
	"io"
	"net/http"
	"net/url"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/httpclient"
	ioutils "github.com/jfrog/jfrog-client-go/utils/io"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
)

type JfrogHttpClient struct {
	httpClient          *httpclient.HttpClient
	JfrogServiceDetails *auth.ServiceDetails
}

func (rtc *JfrogHttpClient) SendGet(url string, followRedirect bool, httpClientsDetails *httputils.HttpClientDetails) (resp *http.Response, respBody []byte, redirectUrl string, err error) {
	err = (*rtc.JfrogServiceDetails).RunPreRequestInterceptors(httpClientsDetails)
	if err != nil {
		return
	}
	return rtc.httpClient.SendGet(url, followRedirect, *httpClientsDetails)
}

func (rtc *JfrogHttpClient) SendPost(url string, content []byte, httpClientsDetails *httputils.HttpClientDetails) (resp *http.Response, body []byte, err error) {
	err = (*rtc.JfrogServiceDetails).RunPreRequestInterceptors(httpClientsDetails)
	if err != nil {
		return
	}
	return rtc.httpClient.SendPost(url, content, *httpClientsDetails)
}

func (rtc *JfrogHttpClient) SendPostLeaveBodyOpen(url string, content []byte, httpClientsDetails *httputils.HttpClientDetails) (*http.Response, error) {
	if err := (*rtc.JfrogServiceDetails).RunPreRequestInterceptors(httpClientsDetails); err != nil {
		return nil, err
	}
	return rtc.httpClient.SendPostLeaveBodyOpen(url, content, *httpClientsDetails)
}

func (rtc *JfrogHttpClient) SendPostForm(url string, data url.Values, httpClientsDetails *httputils.HttpClientDetails) (resp *http.Response, body []byte, err error) {
	httpClientsDetails.Headers["Content-Type"] = "application/x-www-form-urlencoded"
	return rtc.SendPost(url, []byte(data.Encode()), httpClientsDetails)
}

func (rtc *JfrogHttpClient) SendPatch(url string, content []byte, httpClientsDetails *httputils.HttpClientDetails) (resp *http.Response, body []byte, err error) {
	err = (*rtc.JfrogServiceDetails).RunPreRequestInterceptors(httpClientsDetails)
	if err != nil {
		return
	}
	return rtc.httpClient.SendPatch(url, content, *httpClientsDetails)
}

func (rtc *JfrogHttpClient) SendDelete(url string, content []byte, httpClientsDetails *httputils.HttpClientDetails) (resp *http.Response, body []byte, err error) {
	err = (*rtc.JfrogServiceDetails).RunPreRequestInterceptors(httpClientsDetails)
	if err != nil {
		return
	}
	return rtc.httpClient.SendDelete(url, content, *httpClientsDetails)
}

func (rtc *JfrogHttpClient) SendHead(url string, httpClientsDetails *httputils.HttpClientDetails) (resp *http.Response, body []byte, err error) {
	err = (*rtc.JfrogServiceDetails).RunPreRequestInterceptors(httpClientsDetails)
	if err != nil {
		return
	}
	return rtc.httpClient.SendHead(url, *httpClientsDetails)
}

func (rtc *JfrogHttpClient) SendPut(url string, content []byte, httpClientsDetails *httputils.HttpClientDetails) (resp *http.Response, body []byte, err error) {
	err = (*rtc.JfrogServiceDetails).RunPreRequestInterceptors(httpClientsDetails)
	if err != nil {
		return
	}
	return rtc.httpClient.SendPut(url, content, *httpClientsDetails)
}

func (rtc *JfrogHttpClient) Send(method string, url string, content []byte, followRedirect bool, closeBody bool,
	httpClientsDetails *httputils.HttpClientDetails) (resp *http.Response, respBody []byte, redirectUrl string, err error) {
	err = (*rtc.JfrogServiceDetails).RunPreRequestInterceptors(httpClientsDetails)
	if err != nil {
		return
	}
	return rtc.httpClient.Send(method, url, content, followRedirect, closeBody, *httpClientsDetails)
}

func (rtc *JfrogHttpClient) UploadFile(localPath, url, logMsgPrefix string,
	httpClientsDetails *httputils.HttpClientDetails, retries int, progress ioutils.ProgressMgr) (resp *http.Response, body []byte, err error) {
	err = (*rtc.JfrogServiceDetails).RunPreRequestInterceptors(httpClientsDetails)
	if err != nil {
		return
	}
	return rtc.httpClient.UploadFile(localPath, url, logMsgPrefix, *httpClientsDetails, retries, progress)
}

func (rtc *JfrogHttpClient) ReadRemoteFile(downloadPath string, httpClientsDetails *httputils.HttpClientDetails) (ioReaderCloser io.ReadCloser, resp *http.Response, err error) {
	err = (*rtc.JfrogServiceDetails).RunPreRequestInterceptors(httpClientsDetails)
	if err != nil {
		return
	}
	return rtc.httpClient.ReadRemoteFile(downloadPath, *httpClientsDetails)
}

func (rtc *JfrogHttpClient) DownloadFileWithProgress(downloadFileDetails *httpclient.DownloadFileDetails, logMsgPrefix string,
	httpClientsDetails *httputils.HttpClientDetails, retries int, isExplode bool, progress ioutils.ProgressMgr) (resp *http.Response, err error) {
	err = (*rtc.JfrogServiceDetails).RunPreRequestInterceptors(httpClientsDetails)
	if err != nil {
		return
	}
	return rtc.httpClient.DownloadFileWithProgress(downloadFileDetails, logMsgPrefix, *httpClientsDetails, retries, isExplode, progress)
}

func (rtc *JfrogHttpClient) DownloadFile(downloadFileDetails *httpclient.DownloadFileDetails, logMsgPrefix string,
	httpClientsDetails *httputils.HttpClientDetails, retries int, isExplode bool) (resp *http.Response, err error) {
	return rtc.DownloadFileWithProgress(downloadFileDetails, logMsgPrefix, httpClientsDetails, retries, isExplode, nil)
}

func (rtc *JfrogHttpClient) DownloadFileConcurrently(flags httpclient.ConcurrentDownloadFlags,
	logMsgPrefix string, httpClientsDetails *httputils.HttpClientDetails, progress ioutils.ProgressMgr) (resp *http.Response, err error) {
	err = (*rtc.JfrogServiceDetails).RunPreRequestInterceptors(httpClientsDetails)
	if err != nil {
		return
	}
	return rtc.httpClient.DownloadFileConcurrently(flags, logMsgPrefix, *httpClientsDetails, progress)
}

func (rtc *JfrogHttpClient) IsAcceptRanges(downloadUrl string, httpClientsDetails *httputils.HttpClientDetails) (isAcceptRanges bool, resp *http.Response, err error) {
	err = (*rtc.JfrogServiceDetails).RunPreRequestInterceptors(httpClientsDetails)
	if err != nil {
		return
	}
	return rtc.httpClient.IsAcceptRanges(downloadUrl, *httpClientsDetails)
}
