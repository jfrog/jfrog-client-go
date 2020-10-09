package services

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	rthttpclient "github.com/jfrog/jfrog-client-go/artifactory/httpclient"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
)

type VersionService struct {
	client      *rthttpclient.ArtifactoryHttpClient
	XrayDetails auth.ServiceDetails
}

func NewVersionService(client *rthttpclient.ArtifactoryHttpClient) *VersionService {
	return &VersionService{client: client}
}

func (vs *VersionService) GetXrayDetails() auth.ServiceDetails {
	return vs.XrayDetails
}

func (vs *VersionService) GetXrayVersion() (string, error) {
	httpDetails := vs.XrayDetails.CreateHttpClientDetails()
	resp, body, _, err := vs.client.SendGet(vs.XrayDetails.GetUrl()+"api/v1/system/version", true, &httpDetails)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errorutils.CheckError(errors.New("Xray response: " + resp.Status + "\n" + utils.IndentJson(body)))
	}
	var version xrayVersion
	err = json.Unmarshal(body, &version)
	if err != nil {
		return "", errorutils.CheckError(err)
	}
	return strings.TrimSpace(version.Version), nil
}

func (vs *VersionService) GetXrayRevision() (string, error) {
	httpDetails := vs.XrayDetails.CreateHttpClientDetails()
	resp, body, _, err := vs.client.SendGet(vs.XrayDetails.GetUrl()+"api/v1/system/version", true, &httpDetails)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errorutils.CheckError(errors.New("Xray response: " + resp.Status + "\n" + utils.IndentJson(body)))
	}
	var version xrayVersion
	err = json.Unmarshal(body, &version)
	if err != nil {
		return "", errorutils.CheckError(err)
	}
	return strings.TrimSpace(version.Revision), nil
}

type xrayVersion struct {
	Version  string `json:"xray_version,omitempty"`
	Revision string `json:"xray_revision,omitempty"`
}
