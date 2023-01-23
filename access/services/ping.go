package services

import (
	"net/http"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type PingService struct {
	client         *jfroghttpclient.JfrogHttpClient
	ServiceDetails auth.ServiceDetails
}

func NewPingService(client *jfroghttpclient.JfrogHttpClient) *PingService {
	return &PingService{client: client}
}

func (ps *PingService) Ping() ([]byte, error) {
	httpDetails := ps.ServiceDetails.CreateHttpClientDetails()
	resp, body, _, err := ps.client.SendGet(ps.ServiceDetails.GetUrl()+"api/v1/system/ping", true, &httpDetails)
	if err != nil {
		return nil, err
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return body, err
	}
	log.Debug("JFrog Access response:", resp.Status)
	return body, nil
}
