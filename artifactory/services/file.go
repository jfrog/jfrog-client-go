package services

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	rthttpclient "github.com/jfrog/jfrog-client-go/artifactory/httpclient"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type FileService struct {
	client     *rthttpclient.ArtifactoryHttpClient
	ArtDetails auth.ArtifactoryDetails
}

type FileItem struct {
	Uri        string
	Properties map[string][]string
}

type FileChecksums struct {
	MD5    string
	SHA1   string
	SHA256 string
}

type FileInfo struct {
	Uri               string
	DownloadUri       string
	Repo              string
	Path              string
	RemoteUrl         string
	Created           string
	CreatedBy         string
	LastModified      string
	ModifiedBy        string
	LastUpdated       string
	Size              string
	MimeType          string
	Checksums         FileChecksums
	OriginalChecksums FileChecksums
}

type FileStats struct {
	Uri              string
	LastDownloaded   string
	DownloadCount    string
	LastDownloadedBy string
}

func NewFileService(client *rthttpclient.ArtifactoryHttpClient) *FileService {
	return &FileService{client: client}
}

func (fs *FileService) GetArtifactoryDetails() auth.ArtifactoryDetails {
	return fs.ArtDetails
}

func (fs *FileService) SetArtifactoryDetails(rt auth.ArtifactoryDetails) {
	fs.ArtDetails = rt
}

func (fs *FileService) IsDryRun() bool {
	return false
}

func (fs *FileService) performRequest(relativePath string, query string) (body []byte, err error) {
	storageBaseURL := fs.GetArtifactoryDetails().GetUrl() + "api/storage/"
	requestURL := storageBaseURL + relativePath + query

	log.Info("Getting properties...")
	var resp *http.Response
	resp, body, err = fs.sendGetRequest(relativePath, requestURL)

	if err != nil {
		return body, err
	}
	if err = errorutils.CheckResponseStatus(resp, http.StatusOK); err != nil {
		return body, errorutils.CheckError(err)
	}

	return body, err
}

func (fs *FileService) GetProps(relativePath string, props string) ([]utils.Property, error) {
	propList := strings.Split(props, ",")
	encodedParam := ""
	for _, prop := range propList {
		encodedParam += url.QueryEscape(prop) + ","
	}
	// Remove trailing comma
	if strings.HasSuffix(encodedParam, ",") {
		encodedParam = encodedParam[:len(encodedParam)-1]
	}
	getPropertiesQuery := "?properties=" + encodedParam

	body, err := fs.performRequest(relativePath, getPropertiesQuery)

	itemProperties := []utils.Property{}
	if err != nil {
		return itemProperties, err
	}

	var item FileItem

	if err := json.Unmarshal(body, &item); err != nil {
		return itemProperties, err
	}
	for propKey, propValues := range item.Properties {
		for _, propValue := range propValues {
			propEntry := utils.Property{
				Key:   propKey,
				Value: propValue,
			}
			itemProperties = append(itemProperties, propEntry)
		}
	}
	return itemProperties, nil
}

func (fs *FileService) GetInfo(relativePath string) (FileInfo, error) {
	body, err := fs.performRequest(relativePath, "")
	fileInfo := FileInfo{}
	err = json.Unmarshal(body, &fileInfo)
	return fileInfo, err
}

func (fs *FileService) GetLastModified(relativePath string) (string, error) {
	body, err := fs.performRequest(relativePath, "?lastModified")
	type LastModified struct {
		Uri          string
		LastModified string
	}

	resp := LastModified{}
	err = json.Unmarshal(body, &resp)
	return resp.LastModified, err
}

func (fs *FileService) GetStats(relativePath string) (FileStats, error) {
	body, err := fs.performRequest(relativePath, "?stats")

	stats := FileStats{}
	err = json.Unmarshal(body, &stats)
	return stats, err
}

func (fs *FileService) sendGetRequest(relativePath, getPropertiesURL string) (resp *http.Response, body []byte, err error) {
	log.Info("Getting file info on:", relativePath)
	log.Debug("Sending file request:", getPropertiesURL)
	httpClientsDetails := fs.GetArtifactoryDetails().CreateHttpClientDetails()
	resp, body, _, err = fs.client.SendGet(getPropertiesURL, true, &httpClientsDetails)
	return
}
