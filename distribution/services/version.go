package services

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
)

type VersionService struct {
	client      *jfroghttpclient.JfrogHttpClient
	DistDetails auth.ServiceDetails
}

func NewVersionService(client *jfroghttpclient.JfrogHttpClient) *VersionService {
	return &VersionService{client: client}
}

func (vs *VersionService) GetDistDetails() auth.ServiceDetails {
	return vs.DistDetails
}

func (vs *VersionService) GetDistributionVersion() (string, error) {
	httpDetails := vs.DistDetails.CreateHttpClientDetails()
	resp, body, _, err := vs.client.SendGet(vs.DistDetails.GetUrl()+"api/v1/system/info", true, &httpDetails)
	if err != nil {
		return "", errors.New("failed while attempting to get JFrog Distribution version: " + err.Error())
	}

	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return "", errors.New("got unexpected server response while attempting to get JFrog Distribution version:\n" + err.Error())
	}
	var version distributionVersion
	if err = json.Unmarshal(body, &version); err != nil {
		return "", errorutils.CheckErrorf("couldn't parse JFrog Distribution server version version response: " + err.Error())
	}
	return strings.TrimSpace(version.Version), nil
}

type distributionVersion struct {
	Version string `json:"version,omitempty"`
}
