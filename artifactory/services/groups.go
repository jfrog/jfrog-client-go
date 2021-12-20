package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
)

type GroupParams struct {
	GroupDetails    Group
	ReplaceIfExists bool
	IncludeUsers    bool
}

func NewGroupParams() GroupParams {
	return GroupParams{}
}

type Group struct {
	Name            string   `json:"name,omitempty"`
	Description     string   `json:"description,omitempty"`
	AutoJoin        *bool    `json:"autoJoin,omitempty"`
	AdminPrivileges *bool    `json:"adminPrivileges,omitempty"`
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

func (gs *GroupService) GetGroup(params GroupParams) (g *Group, err error) {
	httpDetails := gs.ArtDetails.CreateHttpClientDetails()
	url := fmt.Sprintf("%sapi/security/groups/%s?includeUsers=%t", gs.ArtDetails.GetUrl(), params.GroupDetails.Name, params.IncludeUsers)
	resp, body, _, err := gs.client.SendGet(url, true, &httpDetails)
	if err != nil {
		return nil, err
	}
	// If the requseted group doesn't exists.
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if err = errorutils.CheckResponseStatus(resp, http.StatusOK); err != nil {
		return nil, errorutils.CheckError(errorutils.GenerateResponseError(resp.Status, clientutils.IndentJson(body)))
	}
	var group Group
	if err := json.Unmarshal(body, &group); err != nil {
		return nil, errorutils.CheckError(err)
	}
	return &group, nil
}

func (gs *GroupService) CreateGroup(params GroupParams) error {
	// Checks if the group already exists and act according to ReplaceIfExists parameter.
	if !params.ReplaceIfExists {
		group, err := gs.GetGroup(params)
		if err != nil {
			return err
		}
		if group != nil {
			err = errors.New("Group '" + params.GroupDetails.Name + "' already exists.")
			return errorutils.CheckError(&GroupAlreadyExistsError{InnerError: err})
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
	if err = errorutils.CheckResponseStatus(resp, http.StatusOK, http.StatusCreated); err != nil {
		return errorutils.CheckError(errorutils.GenerateResponseError(resp.Status, clientutils.IndentJson(body)))
	}
	return nil
}

type GroupAlreadyExistsError struct {
	InnerError error
}

func (*GroupAlreadyExistsError) Error() string {
	return "Artifactory: Group already exists."
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
	if err = errorutils.CheckResponseStatus(resp, http.StatusOK); err != nil {
		return errorutils.CheckError(errorutils.GenerateResponseError(resp.Status, clientutils.IndentJson(body)))
	}
	return nil
}

func (gs *GroupService) createOrUpdateGroupRequest(group Group) (url string, requestContent []byte, httpDetails httputils.HttpClientDetails, err error) {
	httpDetails = gs.ArtDetails.CreateHttpClientDetails()
	requestContent, err = json.Marshal(group)
	if errorutils.CheckError(err) != nil {
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
	resp, body, err := gs.client.SendDelete(url, nil, &httpDetails)
	if resp == nil {
		return errorutils.CheckErrorf("no response provided (including status code)")
	}
	if err = errorutils.CheckResponseStatus(resp, http.StatusOK); err != nil {
		return errorutils.CheckError(errorutils.GenerateResponseError(resp.Status, clientutils.IndentJson(body)))
	}
	return err
}
