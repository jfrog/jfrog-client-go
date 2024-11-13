package services

import (
	"encoding/json"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	clientUtils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
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

func (ps *PackageService) GetPackageLeadFile(leadFileRequest LeadFileParams) ([]byte, error) {
	requestUrl, err := clientUtils.BuildUrl(ps.ArtDetails.GetUrl(), apiLeadFile, nil)
	if err != nil {
		return nil, err
	}

	requestContent, err := json.Marshal(leadFileRequest)
	if err != nil {
		return nil, errorutils.CheckError(err)
	}

	httpClientsDetails := ps.ArtDetails.CreateHttpClientDetails()
	httpClientsDetails.SetContentTypeApplicationJson()

	resp, body, err := ps.Client.SendPost(requestUrl, requestContent, &httpClientsDetails)
	if err != nil {
		return nil, err
	}
	return body, errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK)
}

type LeadFileParams struct {
	PackageVersion  string `json:"package_version"`
	PackageName     string `json:"package_name"`
	PackageRepoName string `json:"package_repo_name"`
	PackageType     string `json:"package_type"`
}
