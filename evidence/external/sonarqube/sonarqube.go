package sonarqube

import (
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

const sonarAccessTokenKey = "JF_SONARQUBE_ACCESS_TOKEN"

type SonarQube struct {
	Proxy string
	ServiceConfig
}

type ServiceConfig struct {
	url                  string
	taskAPIPath          string
	projectStatusAPIPath string
}

func NewSonarQubeEvidence(sonarQubeURL, proxy string) *SonarQube {
	return &SonarQube{
		Proxy: proxy,
		ServiceConfig: ServiceConfig{
			url:                  sonarQubeURL,
			taskAPIPath:          "/api/ce/task",
			projectStatusAPIPath: "/api/qualitygates/project_status",
		},
	}
}

func (sqe *SonarQube) createQueryParam(params map[string]string, key, value string) map[string]string {
	if params != nil {
		params[key] = value
		return params
	}
	return map[string]string{
		key: value,
	}
}

func createHttpClient(proxy string) *http.Client {
	transport := &http.Transport{
		MaxIdleConns:      10,
		IdleConnTimeout:   30 * time.Second,
		DisableKeepAlives: false,
		Proxy:             http.ProxyFromEnvironment,
	}
	if proxy != "" {
		proxyURL, err := url.Parse(proxy)
		if err != nil {

			// Fallback to environment proxy or no proxy
			log.Error("Failed to parse proxy URL: " + err.Error())
		} else {
			transport.Proxy = http.ProxyURL(proxyURL)
		}
	}
	return &http.Client{
		Timeout:   30 * time.Second,
		Transport: transport,
	}
}

func (sqe *SonarQube) GetSonarAnalysis(analysisID string) ([]byte, error) {
	log.Debug("Fetching sonar analysis for given analysisID: " + analysisID)
	queryParams := sqe.createQueryParam(nil, "analysisId", analysisID)
	sonarServerURL := sqe.ServiceConfig.url + sqe.ServiceConfig.projectStatusAPIPath

	// Create a new HTTP request
	req, err := http.NewRequest("GET", sonarServerURL, nil)
	if err != nil {
		return nil, err
	}

	// Add query parameters to the request
	q := req.URL.Query()
	for key, value := range queryParams {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()
	resp, bytes, err := sqe.sendRequestUsingSonarQubeToken(req, sqe.Proxy)
	if err != nil {
		return bytes, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error("Failed to close response body: " + err.Error())
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (sqe *SonarQube) CollectSonarQubePredicate(taskID string) ([]byte, error) {
	queryParams := sqe.createQueryParam(nil, "id", taskID)
	sonarServerURL := sqe.ServiceConfig.url + sqe.ServiceConfig.taskAPIPath
	req, err := http.NewRequest(http.MethodGet, sonarServerURL, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	for key, value := range queryParams {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()
	resp, bytes, err := sqe.sendRequestUsingSonarQubeToken(req, sqe.Proxy)
	if err != nil {
		return bytes, err
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			log.Error("Failed to close response body: " + cerr.Error())
		}
	}()
	if resp.StatusCode != http.StatusOK {
		return nil, errorutils.CheckErrorf("Failed to get SonarQube task report. Status code: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (sqe *SonarQube) sendRequestUsingSonarQubeToken(req *http.Request, proxy string) (*http.Response, []byte, error) {
	// Add Authorization header
	sonarQubeToken := os.Getenv(sonarAccessTokenKey)
	if sonarQubeToken == "" {
		return nil, nil, errorutils.CheckErrorf("Sonar access token not found in environment variable " + sonarAccessTokenKey)
	}
	req.Header.Set("Authorization", "Bearer "+sonarQubeToken)
	httpClient := createHttpClient(proxy)

	// Send the request using the standard HTTP client
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	return resp, nil, nil
}
