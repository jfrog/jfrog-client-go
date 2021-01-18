package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
)

type GroupParams struct {
	GroupDetails      Group
	ReplaceExistGroup bool
	IncludeUsers      bool
}

func NewGroupParams() GroupParams {
	return GroupParams{}
}

type Group struct {
	Name            string   `json:"name,omitempty"`
	Description     string   `json:"description,omitempty"`
	AutoJoin        bool     `json:"autoJoin,omitempty"`
	AdminPrivileges bool     `json:"adminPrivileges,omitempty"`
	Realm           string   `json:"realm,omitempty"`
	RealmAttributes string   `json:"realmAttributes,omitempty"`
	UsersNames      []string `json:"userNames,omitempty"`
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

func (gs *GroupService) GetGroup(params GroupParams) (g *Group, notExists bool, err error) {
	httpDetails := gs.ArtDetails.CreateHttpClientDetails()
	url := fmt.Sprintf("%sapi/security/groups/%s?includeUsers=%t", gs.ArtDetails.GetUrl(), params.GroupDetails.Name, params.IncludeUsers)
	res, body, _, err := gs.client.SendGet(url, true, &httpDetails)
	if err != nil {
		return nil, false, err
	}
	if res.StatusCode != http.StatusOK {
		// If the requseted group isn't exists.
		if res.StatusCode == http.StatusNotFound {
			return nil, true, err
		}
		// Other errors from the server
		return nil, false, fmt.Errorf("%d %s: %s", res.StatusCode, res.Status, string(body))
	}
	var group Group
	if err := json.Unmarshal(body, &group); err != nil {
		return nil, false, errorutils.CheckError(err)
	}
	return &group, false, nil
}

func (gs *GroupService) CreateGroup(params GroupParams) error {
	// Checks if the group allready exists in the system and act according to ReplaceExistGroup parameter.
	if !params.ReplaceExistGroup {
		_, notExists, err := gs.GetGroup(params)
		if err != nil {
			return err
		}
		if !notExists {
			return fmt.Errorf("Group %s is allready exists in the system", params.GroupDetails.Name)
		}
	}
	url, content, httpDetails, err := gs.createOrUpdateGroupRequest(params.GroupDetails)
	if err != nil {
		return err
	}
	resp, body, err := gs.client.SendPut(url, content, &httpDetails)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%d %s: %s", resp.StatusCode, resp.Status, string(body))
	}
	return nil
}

func (gs *GroupService) UpdateGroup(params GroupParams) error {
	url, content, httpDetails, err := gs.createOrUpdateGroupRequest(params.GroupDetails)
	if err != nil {
		return err
	}
	resp, body, err := gs.client.SendPost(url, content, &httpDetails)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%d %s: %s", resp.StatusCode, resp.Status, string(body))
	}
	return nil
}

func (gs *GroupService) createOrUpdateGroupRequest(group Group) (url string, requestContent []byte, httpDetails httputils.HttpClientDetails, err error) {
	httpDetails = gs.ArtDetails.CreateHttpClientDetails()
	requestContent, err = json.Marshal(group)
	if err != nil {
		return
	}

	httpDetails.Headers = map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}
	url = fmt.Sprintf("%sapi/security/groups/%s", gs.ArtDetails.GetUrl(), group.Name)
	return
}

func (gs *GroupService) DeleteGroup(name string) error {
	httpDetails := gs.ArtDetails.CreateHttpClientDetails()
	url := fmt.Sprintf("%sapi/security/groups/%s", gs.ArtDetails.GetUrl(), name)
	resp, _, err := gs.client.SendDelete(url, nil, &httpDetails)
	if resp == nil {
		return fmt.Errorf("no response provided (including status code)")
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%d %s", resp.StatusCode, resp.Status)
	}
	return err
}
