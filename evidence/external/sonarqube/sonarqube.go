package sonarqube

import (
	"bufio"
	"errors"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type SonarQube struct {
	componentName string
	ServiceConfig
}

type ServiceConfig struct {
	url                  string
	taskAPIPath          string
	projectStatusAPIPath string
	auth                 Authentication
}

type Authentication struct {
	AccessToken string
}

func NewSonarQubeEvidence() *SonarQube {
	sonarQubeURL := os.Getenv("JF_SONARQUBE_URL")
	componentName := os.Getenv("JF_SONARQUBE_COMPONENT_NAME")

	return &SonarQube{
		componentName: componentName,
		ServiceConfig: ServiceConfig{
			url: sonarQubeURL,
			//apiPath: "/api/measures/component", // Example API path
			//apiPath: "/api/project_analyses/search", // Example API path
			taskAPIPath:          "/api/ce/task", // Example API path
			projectStatusAPIPath: "/api/qualitygates/project_status",
			auth: Authentication{
				AccessToken: "",
			},
		},
	}
}

func (sqe *SonarQube) createQueryParam(params map[string]string, key, value string) map[string]string {
	if params != nil {
		params[key] = value
		return params
	}
	return map[string]string{
		//"metricKeys": "coverage,bugs,vulnerabilities",
		//"component":  sqe.componentName,
		//"project":     sqe.componentName,
		//"buildString": "mvn-sonar-2",
		key: value,
	}
}

func getCeTaskUrlFromFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		panic("Failed to open file: " + err.Error())
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "ceTaskUrl=") {
			taskIDs := strings.Split(line, "?id=")
			if len(taskIDs) < 2 {
				log.Error("Invalid ceTaskUrl format in file")
				return "", errors.New("invalid ceTaskUrl format in file")
			}
			return strings.Split(line, "?id=")[1], nil
		}
	}

	if err := scanner.Err(); err != nil {
		panic("Error reading file: " + err.Error())
	}

	log.Error("ceTaskUrl not found in file")
	return "", errors.New("ceTaskUrl not found in file")
}

func createHttpClient() *http.Client {
	// Create a new HTTP client with custom settings
	client := &http.Client{
		Timeout: 30 * time.Second, // Set a timeout for requests
		Transport: &http.Transport{
			MaxIdleConns:      10,               // Maximum idle connections
			IdleConnTimeout:   30 * time.Second, // Idle connection timeout
			DisableKeepAlives: false,            // Enable keep-alive
		},
	}
	return client
}

func (sqe *SonarQube) GetSonarQubeProjectStatus(analysisID string) ([]byte, error) {
	log.Debug("Getting sonarqube project status for analysis: " + analysisID)
	queryParams := sqe.createQueryParam(nil, "analysisId", analysisID)
	url := sqe.ServiceConfig.url + sqe.ServiceConfig.projectStatusAPIPath
	log.Debug("SonarQube URL: " + url)
	// Create a new HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	// Add query parameters to the request
	q := req.URL.Query()
	for key, value := range queryParams {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()
	resp, bytes, err := sqe.sendRequestUsingSonarQubeToken(req)
	if err != nil {
		return bytes, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (sqe *SonarQube) CollectSonarQubePredicate() ([]byte, error) {
	taskID, err := getCeTaskUrlFromFile("target/sonar/report-task.txt")
	if err != nil {
		log.Error("Failed to get ceTaskUrl: " + err.Error())
		return nil, err
	}
	queryParams := sqe.createQueryParam(nil, "id", taskID)
	url := sqe.ServiceConfig.url + sqe.ServiceConfig.taskAPIPath

	// Create a new HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	// Add query parameters to the request
	q := req.URL.Query()
	for key, value := range queryParams {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()

	resp, bytes, err := sqe.sendRequestUsingSonarQubeToken(req)
	if err != nil {
		return bytes, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (sqe *SonarQube) sendRequestUsingSonarQubeToken(req *http.Request) (*http.Response, []byte, error) {
	// Add Authorization header
	sonarQubeToken := os.Getenv("JF_SONARQUBE_TOKEN")
	req.Header.Set("Authorization", "Bearer "+sonarQubeToken)
	httpClient := createHttpClient()

	// Send the request using the standard HTTP client
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	return resp, nil, nil
}
