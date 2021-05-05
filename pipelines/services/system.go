package services

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils"
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

	if resp.StatusCode != http.StatusOK {
		err := errors.New("Pipelines response: " + resp.Status + "\n" + utils.IndentJson(body))
		if resp.StatusCode == http.StatusNotFound {
			return nil, errorutils.CheckError(&PipelinesNotAvailableError{InnerError: err})
		}
		return nil, errorutils.CheckError(err)
	}
	var sysInfo PipelinesSystemInfo
	err = json.Unmarshal(body, &sysInfo)
	return &sysInfo, errorutils.CheckError(err)
}

type PipelinesSystemInfo struct {
	ServiceId string `json:"serviceId,omitempty"`
	Version   string `json:"version,omitempty"`
}

type PipelinesNotAvailableError struct {
	InnerError error
}

func (*PipelinesNotAvailableError) Error() string {
	return "Pipelines: Pipelines is not aviable at the moment."
}
