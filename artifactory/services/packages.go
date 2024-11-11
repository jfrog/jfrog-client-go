package services

import (
	"encoding/json"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
)

const apiLeadFile = "api/packagesSearch/leadFile"

type PackageService struct {
	Client     *jfroghttpclient.JfrogHttpClient
	ArtDetails auth.ServiceDetails
}

func NewPackageService(client *jfroghttpclient.JfrogHttpClient) *PackageService {
	return &PackageService{Client: client}
}

func (ps *PackageService) GetJfrogHttpClient() *jfroghttpclient.JfrogHttpClient {
	return ps.Client
}

func (ps *PackageService) GetPackageLeadFile(leadFileRequest LeadFileRequest) ([]byte, error) {
	var url = ps.ArtDetails.GetUrl() + apiLeadFile
	log.Info("Sending API request to get LeadFile for package: ", leadFileRequest.PackageName+" version: ", leadFileRequest.PackageVersion)

	requestContent, err := json.Marshal(leadFileRequest)
	if err != nil {
		return nil, errorutils.CheckError(err)
	}

	httpClientsDetails := ps.ArtDetails.CreateHttpClientDetails()
	httpClientsDetails.SetContentTypeApplicationJson()

	resp, body, err := ps.Client.SendPost(url, requestContent, &httpClientsDetails)
	if err != nil {
		return nil, err
	}
	return body, errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK)
}

type LeadFileRequest struct {
	PackageVersion  string `json:"package_version"`
	PackageName     string `json:"package_name"`
	PackageRepoName string `json:"package_repo_name"`
	PackageType     string `json:"package_type"`
}
