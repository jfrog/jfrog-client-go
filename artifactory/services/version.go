package services

import (
	"encoding/json"
	"errors"
	rthttpclient "github.com/jfrog/jfrog-client-go/artifactory/httpclient"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"net/http"
	"strings"
)

type VersionService struct {
	client     *rthttpclient.ArtifactoryHttpClient
	ArtDetails auth.CommonDetails
}

func NewVersionService(client *rthttpclient.ArtifactoryHttpClient) *VersionService {
	return &VersionService{client: client}
}

func (vs *VersionService) GetArtifactoryDetails() auth.CommonDetails {
	return vs.ArtDetails
}

func (vs *VersionService) SetArtifactoryDetails(rt auth.CommonDetails) {
	vs.ArtDetails = rt
}

func (vs *VersionService) GetJfrogHttpClient() (*rthttpclient.ArtifactoryHttpClient, error) {
	return vs.client, nil
}

func (vs *VersionService) IsDryRun() bool {
	return false
}

func (vs *VersionService) GetArtifactoryVersion() (string, error) {
	httpDetails := vs.ArtDetails.CreateHttpClientDetails()
	resp, body, _, err := vs.client.SendGet(vs.ArtDetails.GetUrl()+"api/system/version", true, &httpDetails)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errorutils.CheckError(errors.New("Artifactory response: " + resp.Status + "\n" + utils.IndentJson(body)))
	}
	var version artifactoryVersion
	err = json.Unmarshal(body, &version)
	if err != nil {
		return "", errorutils.CheckError(err)
	}
	return strings.TrimSpace(version.Version), nil
}

type artifactoryVersion struct {
	Version string `json:"version,omitempty"`
}
