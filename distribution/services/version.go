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
	client      *rthttpclient.ArtifactoryHttpClient
	DistDetails auth.ServiceDetails
}

func NewVersionService(client *rthttpclient.ArtifactoryHttpClient) *VersionService {
	return &VersionService{client: client}
}

func (vs *VersionService) GetDistDetails() auth.ServiceDetails {
	return vs.DistDetails
}

func (vs *VersionService) GetDistributionVersion() (string, error) {
	httpDetails := vs.DistDetails.CreateHttpClientDetails()
	resp, body, _, err := vs.client.SendGet(vs.DistDetails.GetUrl()+"api/v1/system/info", true, &httpDetails)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errorutils.CheckError(errors.New("Distribution response: " + resp.Status + "\n" + utils.IndentJson(body)))
	}
	var version distributionVersion
	err = json.Unmarshal(body, &version)
	if err != nil {
		return "", errorutils.CheckError(err)
	}
	return strings.TrimSpace(version.Version), nil
}

type distributionVersion struct {
	Version string `json:"version,omitempty"`
}
