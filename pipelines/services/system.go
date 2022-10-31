package services

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type SystemService struct {
	client *jfroghttpclient.JfrogHttpClient
	auth.ServiceDetails
}

func NewSystemService(client *jfroghttpclient.JfrogHttpClient) *SystemService {
	return &SystemService{client: client}
}

func (ss *SystemService) GetSystemInfo() (*PipelinesSystemInfo, error) {
	log.Debug("Getting Pipelines System Info...")
	httpDetails := ss.ServiceDetails.CreateHttpClientDetails()
	resp, body, _, err := ss.client.SendGet(ss.ServiceDetails.GetUrl()+"api/v1/system/info", true, &httpDetails)
	if err != nil {
		return nil, err
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		if resp.StatusCode == http.StatusNotFound {
			return nil, &PipelinesNotAvailableError{InnerError: err}
		}
		return nil, errors.New("failed while attempting to get Pipelines version:\n" + err.Error())
	}
	var sysInfo PipelinesSystemInfo
	err = json.Unmarshal(body, &sysInfo)
	if err != nil {
		return nil, errorutils.CheckErrorf("couldn't parse Pipelines server response: " + err.Error())
	}
	return &sysInfo, nil
}

type PipelinesSystemInfo struct {
	ServiceId string `json:"serviceId,omitempty"`
	Version   string `json:"version,omitempty"`
}

type PipelinesNotAvailableError struct {
	InnerError error
}

func (*PipelinesNotAvailableError) Error() string {
	return "Pipelines: Pipelines is not available at the moment."
}
