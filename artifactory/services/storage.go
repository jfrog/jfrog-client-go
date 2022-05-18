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
)

type StorageService struct {
	client     *jfroghttpclient.JfrogHttpClient
	artDetails *auth.ServiceDetails
}

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
	restAPI := path.Join("api", "storage", relativePath)
	folderUrl, err := utils.BuildArtifactoryUrl(s.GetArtifactoryDetails().GetUrl(), restAPI, make(map[string]string))
	if err != nil {
		return nil, err
	}

	httpClientsDetails := s.GetArtifactoryDetails().CreateHttpClientDetails()
	resp, body, _, err := client.SendGet(folderUrl, true, &httpClientsDetails)
	if err != nil {
		return nil, err
	}
	if err = errorutils.CheckResponseStatus(resp, http.StatusOK); err != nil {
		return nil, errorutils.CheckError(err)
	}
	log.Debug("Artifactory response: ", resp.Status)

	result := &utils.FolderInfo{}
	err = json.Unmarshal(body, result)
	return result, errorutils.CheckError(err)
}

func (s *StorageService) FileList(relativePath string) (*utils.FileList, error) {
	client := s.GetJfrogHttpClient()
	restAPI := path.Join("api", "storage", relativePath)
	folderUrl, err := utils.BuildArtifactoryUrl(s.GetArtifactoryDetails().GetUrl(), restAPI, make(map[string]string))
	if err != nil {
		return nil, err
	}
	folderUrl += "?list&listFolders=1"

	httpClientsDetails := s.GetArtifactoryDetails().CreateHttpClientDetails()
	resp, body, _, err := client.SendGet(folderUrl, true, &httpClientsDetails)
	if err != nil {
		return nil, err
	}
	if err = errorutils.CheckResponseStatus(resp, http.StatusOK); err != nil {
		return nil, errorutils.CheckError(err)
	}
	log.Debug("Artifactory response: ", resp.Status)

	result := &utils.FileList{}
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
	if err = errorutils.CheckResponseStatus(resp, http.StatusOK); err != nil {
		return nil, errorutils.CheckError(err)
	}
	log.Debug("Artifactory response: ", resp.Status)

	result := &utils.StorageInfo{}
	err = json.Unmarshal(body, result)
	return result, errorutils.CheckError(err)
}
