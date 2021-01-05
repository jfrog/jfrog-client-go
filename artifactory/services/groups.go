package services

import (
	"encoding/json"
	"fmt"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"net/http"
)

type Group struct {
	Name            string `json:"name,omitempty"`
	Description     string `json:"description,omitempty"`
	AutoJoin        bool   `json:"autoJoin,omitempty"`
	AdminPrivileges bool   `json:"adminPrivileges,omitempty"`
	Realm           string `json:"realm,omitempty"`
	RealmAttributes string `json:"realmAttributes,omitempty"`
}

type GroupService struct {
	client     *jfroghttpclient.JfrogHttpClient
	ArtDetails auth.ServiceDetails
}

func NewGroupService(client *jfroghttpclient.JfrogHttpClient) *GroupService {
	return &GroupService{client: client}
}
func (gs *GroupService) SetArtifactoryDetails(rt auth.ServiceDetails) {
	gs.ArtDetails = rt
}

func (gs *GroupService) GetGroup(name string) (*Group, error) {
	httpDetails := gs.ArtDetails.CreateHttpClientDetails()
	url := fmt.Sprintf("%sapi/security/groups/%s", gs.ArtDetails.GetUrl(), name)
	res, body, _, err := gs.client.SendGet(url, true, &httpDetails)
	if err != nil {
		return nil, err
	}
	if res.StatusCode > 204 {
		return nil, fmt.Errorf("%d %s: %s", res.StatusCode, res.Status, string(body))
	}
	var group Group
	if err := json.Unmarshal(body, &group); err != nil {
		return nil, errorutils.CheckError(err)
	}
	return &group, nil
}

func (gs *GroupService) CreateOrUpdateGroup(group Group) error {
	httpDetails := gs.ArtDetails.CreateHttpClientDetails()
	content, err := json.Marshal(group)
	if err != nil {
		return err
	}

	httpDetails.Headers = map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}
	//
	url := fmt.Sprintf("%sapi/security/groups/%s", gs.ArtDetails.GetUrl(), group.Name)
	resp, body, err := gs.client.SendPut(url, content, &httpDetails)
	if err != nil {
		return err
	}
	if resp.StatusCode > http.StatusNoContent {
		return fmt.Errorf("%d %s: %s", resp.StatusCode, resp.Status, string(body))
	}
	return nil
}

func (gs *GroupService) DeleteGroup(name string) error {
	httpDetails := gs.ArtDetails.CreateHttpClientDetails()
	url := fmt.Sprintf("%sapi/security/groups/%s", gs.ArtDetails.GetUrl(), name)
	resp, _, err := gs.client.SendDelete(url, nil, &httpDetails)
	if resp == nil {
		return fmt.Errorf("no response provided (including status code)")
	}
	if resp.StatusCode > 204 {
		return fmt.Errorf("%d %s", resp.StatusCode, resp.Status)
	}
	return err
}

func (gs *GroupService) GroupExits(name string) (bool, error) {
	// usage of HEAD chokes on an internal proxy issue. Apparently multiple services
	// are running in 1 container
	group, err := gs.GetGroup(name)
	return err != nil && group != nil, err
}
