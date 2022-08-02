package services

import (
	"encoding/json"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
	"path"
	"strconv"
)

type StorageService struct {
	client     *jfroghttpclient.JfrogHttpClient
	artDetails *auth.ServiceDetails
}

const StorageRestApi = "api/storage/"

func NewStorageService(artDetails auth.ServiceDetails, client *jfroghttpclient.JfrogHttpClient) *StorageService {
	return &StorageService{artDetails: &artDetails, client: client}
}

func (s *StorageService) GetArtifactoryDetails() auth.ServiceDetails {
	return *s.artDetails
}

func (s *StorageService) GetJfrogHttpClient() *jfroghttpclient.JfrogHttpClient {
	return s.client
}

func (s *StorageService) FolderInfo(relativePath string) (*utils.FolderInfo, error) {
	client := s.GetJfrogHttpClient()
	restAPI := path.Join(StorageRestApi, relativePath)
	folderUrl, err := utils.BuildArtifactoryUrl(s.GetArtifactoryDetails().GetUrl(), restAPI, make(map[string]string))
	if err != nil {
		return nil, err
	}

	httpClientsDetails := s.GetArtifactoryDetails().CreateHttpClientDetails()
	resp, body, _, err := client.SendGet(folderUrl, true, &httpClientsDetails)
	if err != nil {
		return nil, err
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return nil, err
	}
	log.Debug("Artifactory response: ", resp.Status)

	result := &utils.FolderInfo{}
	err = json.Unmarshal(body, result)
	return result, errorutils.CheckError(err)
}

func (s *StorageService) FileList(relativePath string, optionalParams utils.FileListParams) (*utils.FileListResponse, error) {
	client := s.GetJfrogHttpClient()
	restAPI := path.Join(StorageRestApi, relativePath)

	// Convert params to map:
	params := make(map[string]string)
	params["list"] = "true"
	addParamIfTrue(params, "deep", optionalParams.Deep)
	addParamIfTrue(params, "listFolders", optionalParams.ListFolders)
	addParamIfTrue(params, "mdTimestamps", optionalParams.MetadataTimestamps)
	addParamIfTrue(params, "includeRootPath", optionalParams.IncludeRootPath)
	if optionalParams.Depth > 0 {
		params["depth"] = strconv.Itoa(optionalParams.Depth)
	}

	folderUrl, err := utils.BuildArtifactoryUrl(s.GetArtifactoryDetails().GetUrl(), restAPI, params)
	if err != nil {
		return nil, err
	}

	httpClientsDetails := s.GetArtifactoryDetails().CreateHttpClientDetails()
	resp, body, _, err := client.SendGet(folderUrl, true, &httpClientsDetails)
	if err != nil {
		return nil, err
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return nil, err
	}
	log.Debug("Artifactory response: ", resp.Status)

	result := &utils.FileListResponse{}
	err = json.Unmarshal(body, result)
	return result, errorutils.CheckError(err)
}

func (s *StorageService) StorageInfo() (*utils.StorageInfo, error) {
	client := s.GetJfrogHttpClient()
	url := s.GetArtifactoryDetails().GetUrl() + "api/storageinfo"

	httpClientsDetails := s.GetArtifactoryDetails().CreateHttpClientDetails()
	resp, body, _, err := client.SendGet(url, true, &httpClientsDetails)
	if err != nil {
		return nil, err
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return nil, err
	}
	log.Debug("Artifactory response: ", resp.Status)

	result := &utils.StorageInfo{}
	err = json.Unmarshal(body, result)
	return result, errorutils.CheckError(err)
}

func (s *StorageService) StorageInfoRefresh() error {
	client := s.GetJfrogHttpClient()
	url := s.GetArtifactoryDetails().GetUrl() + "api/storageinfo/calculate"

	httpClientsDetails := s.GetArtifactoryDetails().CreateHttpClientDetails()
	resp, body, err := client.SendPost(url, nil, &httpClientsDetails)
	if err != nil {
		return err
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusAccepted); err != nil {
		return err
	}
	log.Debug("Artifactory response: ", resp.Status)
	return nil
}

func addParamIfTrue(params map[string]string, paramName string, value bool) {
	if value {
		params[paramName] = "1"
	}
}
