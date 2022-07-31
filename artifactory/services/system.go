package services

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

const (
	apiSystem         = "api/system/"
	runningNodeStatus = "RUNNING"
)

type SystemService struct {
	client     *jfroghttpclient.JfrogHttpClient
	artDetails *auth.ServiceDetails
}

func NewSystemService(artDetails auth.ServiceDetails, client *jfroghttpclient.JfrogHttpClient) *SystemService {
	return &SystemService{artDetails: &artDetails, client: client}
}

func (ss *SystemService) GetArtifactoryDetails() auth.ServiceDetails {
	return *ss.artDetails
}

func (ss *SystemService) GetJfrogHttpClient() *jfroghttpclient.JfrogHttpClient {
	return ss.client
}

func (ss *SystemService) IsDryRun() bool {
	return false
}

func (ss *SystemService) GetVersion() (string, error) {
	body, err := ss.sendGet("version")
	if err != nil {
		return "", err
	}
	var version artifactoryVersion
	err = json.Unmarshal(body, &version)
	if err != nil {
		return "", errorutils.CheckError(err)
	}
	return strings.TrimSpace(version.Version), nil
}

func (ss *SystemService) GetServiceId() (string, error) {
	body, err := ss.sendGet("service_id")
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (ss *SystemService) GetRunningNodes() ([]string, error) {
	body, err := ss.sendGet("status")
	if err != nil {
		return []string{}, err
	}
	var status artifactoryStatus
	err = json.Unmarshal(body, &status)
	if err != nil {
		return []string{}, errorutils.CheckError(err)
	}
	var runningNodes []string
	for _, node := range status.Nodes {
		if node.State == runningNodeStatus {
			runningNodes = append(runningNodes, strings.TrimSpace(node.Id))
		}
	}
	return runningNodes, nil
}

func (ss *SystemService) GetConfigDescriptor() (string, error) {
	log.Info("Fetching config descriptor from Artifactory...")
	body, err := ss.sendGet("configuration")
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (ss *SystemService) ActivateKeyEncryption() error {
	log.Info("Activating key encryption in Artifactory...")
	if err := ss.sendEmptyPost("encrypt"); err != nil {
		return err
	}
	log.Info("Artifactory key encryption activated")
	return nil
}

func (ss *SystemService) DeactivateKeyEncryption() error {
	log.Info("Deactivating key encryption in Artifactory...")
	if err := ss.sendEmptyPost("decrypt"); err != nil {
		return err
	}
	log.Info("Artifactory key encryption deactivated")
	return nil
}

func (ss *SystemService) sendGet(endpoint string) ([]byte, error) {
	httpDetails := (*ss.artDetails).CreateHttpClientDetails()
	resp, body, _, err := ss.client.SendGet((*ss.artDetails).GetUrl()+apiSystem+endpoint, true, &httpDetails)
	if err != nil {
		return nil, err
	}
	if err = errorutils.CheckResponseStatus(resp, body, http.StatusOK, http.StatusCreated); err != nil {
		return nil, err
	}
	log.Debug("Artifactory response:", resp.Status)
	return body, nil
}

func (ss *SystemService) sendEmptyPost(endpoint string) error {
	httpDetails := (*ss.artDetails).CreateHttpClientDetails()
	resp, body, err := ss.client.SendPost((*ss.artDetails).GetUrl()+apiSystem+endpoint, nil, &httpDetails)
	if err != nil {
		return err
	}
	if err = errorutils.CheckResponseStatus(resp, body, http.StatusOK, http.StatusCreated); err != nil {
		return err
	}
	log.Debug("Artifactory response:", string(body), resp.Status)
	return nil
}

type artifactoryVersion struct {
	Version string `json:"version,omitempty"`
}

type artifactoryStatus struct {
	Nodes []struct {
		Id    string `json:"id,omitempty"`
		State string `json:"state,omitempty"`
	} `json:"nodes,omitempty"`
}
