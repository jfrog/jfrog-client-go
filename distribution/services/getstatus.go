package services

import (
	"encoding/json"
	"errors"
	"github.com/jfrog/jfrog-client-go/utils/distribution"
	"net/http"
	"strings"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type DistributionStatusService struct {
	client      *jfroghttpclient.JfrogHttpClient
	DistDetails auth.ServiceDetails
}

func NewDistributionStatusService(client *jfroghttpclient.JfrogHttpClient) *DistributionStatusService {
	return &DistributionStatusService{client: client}
}

func (ds *DistributionStatusService) GetDistDetails() auth.ServiceDetails {
	return ds.DistDetails
}

func (ds *DistributionStatusService) GetStatus(distributionStatusParams DistributionStatusParams) (*[]distribution.DistributionStatusResponse, error) {
	if err := ds.checkParameters(distributionStatusParams); err != nil {
		return nil, err
	}
	return ds.execGetStatus(distributionStatusParams.Name, distributionStatusParams.Version, distributionStatusParams.TrackerId)
}

func (ds *DistributionStatusService) checkParameters(distributionStatusParams DistributionStatusParams) error {
	var err error
	if distributionStatusParams.Name == "" && (distributionStatusParams.Version != "" || distributionStatusParams.TrackerId != "") {
		err = errors.New("missing distribution name parameter")
	}
	if distributionStatusParams.Version == "" && distributionStatusParams.TrackerId != "" {
		err = errors.New("missing distribution version parameter")
	}
	return errorutils.CheckError(err)
}

func (ds *DistributionStatusService) execGetStatus(name, version, trackerId string) (*[]distribution.DistributionStatusResponse, error) {
	httpClientsDetails := ds.DistDetails.CreateHttpClientDetails()
	url := ds.BuildUrlForGetStatus(ds.DistDetails.GetUrl(), name, version, trackerId)

	resp, body, _, err := ds.client.SendGet(url, true, &httpClientsDetails)
	if err != nil {
		return nil, err
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return nil, err
	}
	log.Debug("Distribution response:", resp.Status)
	log.Debug(utils.IndentJson(body))
	var distributionStatusResponse []distribution.DistributionStatusResponse
	stringBody := string(body)
	if !strings.HasPrefix(stringBody, "[") {
		stringBody = "[" + stringBody + "]"
	}
	err = json.Unmarshal([]byte(stringBody), &distributionStatusResponse)
	return &distributionStatusResponse, errorutils.CheckError(err)
}

func (ds *DistributionStatusService) BuildUrlForGetStatus(url, name, version, trackerId string) string {
	url += "api/v1/release_bundle"
	if name == "" {
		return url + "/distribution"
	}
	url += "/" + name
	if version == "" {
		return url + "/distribution"
	}
	url += "/" + version + "/distribution"
	if trackerId != "" {
		return url + "/" + trackerId
	}
	return url
}

type DistributionStatusParams struct {
	Name      string
	Version   string
	TrackerId string
}

func NewDistributionStatusParams() DistributionStatusParams {
	return DistributionStatusParams{}
}
