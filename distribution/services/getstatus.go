package services

import (
	"encoding/json"
	"errors"
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

func (ds *DistributionStatusService) GetStatus(distributionStatusParams DistributionStatusParams) (*[]DistributionStatusResponse, error) {
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

func (ds *DistributionStatusService) execGetStatus(name, version, trackerId string) (*[]DistributionStatusResponse, error) {
	httpClientsDetails := ds.DistDetails.CreateHttpClientDetails()
	url := ds.BuildUrlForGetStatus(ds.DistDetails.GetUrl(), name, version, trackerId)

	resp, body, _, err := ds.client.SendGet(url, true, &httpClientsDetails)
	if err != nil {
		return nil, err
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return nil, err
	}
	log.Debug("Distribution response: ", resp.Status)
	log.Debug(utils.IndentJson(body))
	var distributionStatusResponse []DistributionStatusResponse
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

type DistributionType string

const (
	Distribute                 DistributionType = "distribute"
	DeleteReleaseBundleVersion DistributionType = "delete_release_bundle_version"
)

type DistributionStatus string

const (
	NotDistributed DistributionStatus = "Not distributed"
	InProgress     DistributionStatus = "In progress"
	Completed      DistributionStatus = "Completed"
	Failed         DistributionStatus = "Failed"
)

type DistributionStatusResponse struct {
	Id                json.Number              `json:"distribution_id"`
	FriendlyId        json.Number              `json:"distribution_friendly_id,omitempty"`
	Type              DistributionType         `json:"type,omitempty"`
	Name              string                   `json:"release_bundle_name,omitempty"`
	Version           string                   `json:"release_bundle_version,omitempty"`
	Status            DistributionStatus       `json:"status,omitempty"`
	DistributionRules []DistributionRulesBody  `json:"distribution_rules,omitempty"`
	Sites             []DistributionSiteStatus `json:"sites,omitempty"`
}

type DistributionSiteStatus struct {
	Status            string            `json:"status,omitempty"`
	Error             string            `json:"general_error,omitempty"`
	TargetArtifactory TargetArtifactory `json:"target_artifactory,omitempty"`
	TotalFiles        json.Number       `json:"total_files,omitempty"`
	TotalBytes        json.Number       `json:"total_bytes,omitempty"`
	DistributedBytes  json.Number       `json:"distributed_bytes,omitempty"`
	DistributedFiles  json.Number       `json:"distributed_files,omitempty"`
	FileErrors        []string          `json:"file_errors,omitempty"`
	FilesInProgress   []string          `json:"files_in_progress,omitempty"`
}

type TargetArtifactory struct {
	ServiceId string `json:"service_id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
}
