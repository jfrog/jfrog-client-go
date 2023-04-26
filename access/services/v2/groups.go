package v2

import (
	"encoding/json"
	"fmt"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"net/http"
)

const groupsApi = "api/v2/groups"

type GroupParams struct {
	GroupDetails
}

type GroupDetails struct {
	Name            string   `json:"name,omitempty"`
	Description     string   `json:"description,omitempty"`
	AutoJoin        *bool    `json:"auto_join,omitempty"`
	AdminPrivileges *bool    `json:"admin_privileges,omitempty"`
	Realm           string   `json:"realm,omitempty"`
	RealmAttributes string   `json:"realm_attributes,omitempty"`
	ExternalId      string   `json:"external_id,omitempty"`
	Members         []string `json:"members,omitempty"`
}

type GroupListItem struct {
	GroupName string `json:"group_name"`
	Uri       string `json:"uri"`
}
type GroupList struct {
	Cursor string          `json:"cursor,omitempty"`
	Groups []GroupListItem `json:"groups"`
}

func NewGroupParams() GroupParams {
	return GroupParams{}
}

type GroupService struct {
	client         *jfroghttpclient.JfrogHttpClient
	ServiceDetails auth.ServiceDetails
}

func NewGroupService(client *jfroghttpclient.JfrogHttpClient) *GroupService {
	return &GroupService{client: client}
}

func (gs *GroupService) getBaseUrl() string {
	return fmt.Sprintf("%s%s", gs.ServiceDetails.GetUrl(), groupsApi)
}
func (gs *GroupService) GetAll() ([]GroupListItem, error) {
	httpDetails := gs.ServiceDetails.CreateHttpClientDetails()
	url := gs.getBaseUrl()
	resp, body, _, err := gs.client.SendGet(url, true, &httpDetails)
	if err != nil {
		return nil, err
	}

	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return nil, err
	}
	var groupList GroupList
	err = json.Unmarshal(body, &groupList)
	return groupList.Groups, errorutils.CheckError(err)
}

func (gs *GroupService) Get(name string) (u *GroupDetails, err error) {
	httpDetails := gs.ServiceDetails.CreateHttpClientDetails()
	url := fmt.Sprintf("%s/%s", gs.getBaseUrl(), name)
	resp, body, _, err := gs.client.SendGet(url, true, &httpDetails)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return nil, err
	}
	var group GroupDetails
	err = json.Unmarshal(body, &group)
	return &group, errorutils.CheckError(err)
}

func (gs *GroupService) Create(params GroupParams) error {
	group, err := gs.Get(params.Name)
	if err != nil {
		return err
	}
	if group != nil {
		return errorutils.CheckErrorf("group '%s' already exists", group.Name)
	}
	content, httpDetails, err := gs.createOrUpdateRequest(params)
	if err != nil {
		return err
	}
	resp, body, err := gs.client.SendPost(gs.getBaseUrl(), content, &httpDetails)
	if err != nil {
		return err
	}
	return errorutils.CheckResponseStatusWithBody(resp, body, http.StatusCreated)
}

func (gs *GroupService) Update(params GroupParams) error {
	content, httpDetails, err := gs.createOrUpdateRequest(params)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/%s", gs.getBaseUrl(), params.Name)
	resp, body, err := gs.client.SendPatch(url, content, &httpDetails)
	if err != nil {
		return err
	}
	return errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK)
}

func (gs *GroupService) createOrUpdateRequest(group GroupParams) (requestContent []byte, httpDetails httputils.HttpClientDetails, err error) {
	httpDetails = gs.ServiceDetails.CreateHttpClientDetails()
	requestContent, err = json.Marshal(group)
	if errorutils.CheckError(err) != nil {
		return
	}
	httpDetails.Headers = map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}
	return
}

func (gs *GroupService) Delete(name string) error {
	httpDetails := gs.ServiceDetails.CreateHttpClientDetails()
	url := fmt.Sprintf("%s/%s", gs.getBaseUrl(), name)
	resp, body, err := gs.client.SendDelete(url, nil, &httpDetails)
	if err != nil {
		return err
	}
	if resp == nil {
		return errorutils.CheckErrorf("no response provided (including status code)")
	}
	return errorutils.CheckResponseStatusWithBody(resp, body, http.StatusNoContent)
}
