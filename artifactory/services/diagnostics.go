package services

import (
	//"github.com/jfrog/gofrog/parallel"

	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	rthttpclient "github.com/jfrog/jfrog-client-go/artifactory/httpclient"

	//"github.com/jfrog/jfrog-client-go/artifactory/services/fspatterns"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	//clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	ioutils "github.com/jfrog/jfrog-client-go/utils/io"

	//"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	//"github.com/jfrog/jfrog-client-go/utils/io/fileutils/checksum"
	//"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"errors"
	"net/http"
	"time"

	"github.com/jfrog/jfrog-client-go/utils/log"
	//"os"
	//"regexp"
	//"sort"
	//"strconv"
	//"strings"
)

type DiagnosticsService struct {
	client     *rthttpclient.ArtifactoryHttpClient
	Progress   ioutils.Progress
	ArtDetails auth.ArtifactoryDetails
	DryRun     bool
	Threads    int
}

func NewDiagnosticsService(client *rthttpclient.ArtifactoryHttpClient) *DiagnosticsService {
	return &DiagnosticsService{client: client}
}

func (ds *DiagnosticsService) GetThreads() int {
	return ds.Threads
}

func (ds *DiagnosticsService) SetThreads(threads int) {
	ds.Threads = threads
}

func (ds *DiagnosticsService) GetJfrogHttpClient() *rthttpclient.ArtifactoryHttpClient {
	return ds.client
}

func (ds *DiagnosticsService) SetArtDetails(artDetails auth.ArtifactoryDetails) {
	ds.ArtDetails = artDetails
}

func (ds *DiagnosticsService) SetDryRun(isDryRun bool) {
	ds.DryRun = isDryRun
}

func (ds *DiagnosticsService) GetSystemInfo() ([]byte, error) {
	url, err := utils.BuildArtifactoryUrl(ds.ArtDetails.GetUrl(), "api/system", nil)
	if err != nil {
		return nil, err
	}
	httpClientDetails := ds.ArtDetails.CreateHttpClientDetails()
	resp, respBody, _, err := ds.client.SendGet(url, true, &httpClientDetails)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return respBody, errorutils.CheckError(errors.New("Artifactory response: " + resp.Status))
	}
	log.Debug("Artifactory response: ", resp.Status)
	return respBody, nil
}

func (ds *DiagnosticsService) GetSysConf() ([]byte, error) {
	url, err := utils.BuildArtifactoryUrl(ds.ArtDetails.GetUrl(), "api/system/configuration", nil)
	if err != nil {
		return nil, err
	}
	httpClientDetails := ds.ArtDetails.CreateHttpClientDetails()
	resp, respBody, _, err := ds.client.SendGet(url, true, &httpClientDetails)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return respBody, errorutils.CheckError(errors.New("Artifactory response: " + resp.Status))
	}
	log.Debug("Artifactory response: ", resp.Status)
	return respBody, nil
}

func (ds *DiagnosticsService) GetWebServerInfo() ([]byte, error) {
	url, err := utils.BuildArtifactoryUrl(ds.ArtDetails.GetUrl(), "api/system/configuration/webServer", nil)
	if err != nil {
		return nil, err
	}
	httpClientDetails := ds.ArtDetails.CreateHttpClientDetails()
	resp, respBody, _, err := ds.client.SendGet(url, true, &httpClientDetails)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return respBody, errorutils.CheckError(errors.New("Artifactory response: " + resp.Status))
	}
	log.Debug("Artifactory response: ", resp.Status)
	return respBody, nil
}

func (ds *DiagnosticsService) StartDiagnosis() (*DiagnosisResult, error) {
	version, _ := ds.ArtDetails.GetVersion()
	result := &DiagnosisResult{
		Version: version,
		Url:     ds.ArtDetails.GetUrl(),
		Tasks:   make(map[string]*DiagnosisTask),
	}
	var tasks = make(map[string](func() ([]byte, error)))
	tasks["SystemInfo"] = ds.GetSystemInfo
	tasks["GetSysConf"] = ds.GetSysConf
	tasks["GetWebServerInfo"] = ds.GetWebServerInfo
	for taskName, task := range tasks {
		starttime := time.Now()
		taskResult := &DiagnosisTask{
			Name: taskName,
		}
		out, err := task()
		if err != nil {
			taskResult.Err = err
		} else {
			taskResult.Output = out
		}
		taskResult.ElapsedTime = time.Since(starttime)
		result.Tasks[taskName] = taskResult
	}
	return result, nil
}

type DiagnosisTask struct {
	Name        string
	Output      []byte
	ElapsedTime time.Duration
	Err         error
}

type DiagnosisResult struct {
	Url     string
	Version string
	Tasks   map[string]*DiagnosisTask
}
