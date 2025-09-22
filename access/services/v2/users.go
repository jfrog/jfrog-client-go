package v2

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
)

const usersApi = "api/v2/users"

type UserParams struct {
	CommonUserParams
}

type UserResponse struct {
	CommonUserParams
	Status string `json:"status,omitempty"`
}

func NewUserParams() UserParams {
	return UserParams{}
}

type CommonUserParams struct {
	Username                 string    `json:"username,omitempty"`
	Email                    string    `json:"email,omitempty"`
	Password                 string    `json:"password,omitempty"`
	Admin                    *bool     `json:"admin,omitempty"`
	ProfileUpdatable         *bool     `json:"profile_updatable,omitempty"`
	DisableUIAccess          *bool     `json:"disable_ui_access,omitempty"`
	InternalPasswordDisabled *bool     `json:"internal_password_disabled,omitempty"`
	Realm                    string    `json:"realm,omitempty"`
	Groups                   *[]string `json:"groups,omitempty"`
}

type UserService struct {
	client         *jfroghttpclient.JfrogHttpClient
	ServiceDetails auth.ServiceDetails
}

func NewUserService(client *jfroghttpclient.JfrogHttpClient) *UserService {
	return &UserService{client: client}
}

func (us *UserService) getBaseUrl() string {
	return fmt.Sprintf("%s%s", us.ServiceDetails.GetUrl(), usersApi)
}

func (us *UserService) GetAll() ([]UserResponse, error) {
	httpDetails := us.ServiceDetails.CreateHttpClientDetails()
	url := us.getBaseUrl()
	resp, body, _, err := us.client.SendGet(url, true, &httpDetails)
	if err != nil {
		return nil, err
	}

	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return nil, err
	}
	var users []UserResponse
	err = json.Unmarshal(body, &users)
	return users, errorutils.CheckError(err)
}

func (us *UserService) Get(username string) (u *UserResponse, err error) {
	httpDetails := us.ServiceDetails.CreateHttpClientDetails()
	url := fmt.Sprintf("%s/%s", us.getBaseUrl(), username)
	resp, body, _, err := us.client.SendGet(url, true, &httpDetails)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return nil, err
	}
	var user UserResponse
	err = json.Unmarshal(body, &user)
	return &user, errorutils.CheckError(err)
}

func (us *UserService) Create(params UserParams) error {
	user, err := us.Get(params.Username)
	if err != nil {
		return err
	}
	if user != nil {
		return errorutils.CheckErrorf("user '%s' already exists", user.Username)
	}
	content, httpDetails, err := us.createOrUpdateRequest(params)
	if err != nil {
		return err
	}
	resp, body, err := us.client.SendPost(us.getBaseUrl(), content, &httpDetails)
	if err != nil {
		return err
	}
	return errorutils.CheckResponseStatusWithBody(resp, body, http.StatusCreated)
}

func (us *UserService) Update(params UserParams) error {
	content, httpDetails, err := us.createOrUpdateRequest(params)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/%s", us.getBaseUrl(), params.Username)
	resp, body, err := us.client.SendPatch(url, content, &httpDetails)
	if err != nil {
		return err
	}
	return errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK)
}

func (us *UserService) createOrUpdateRequest(user UserParams) (requestContent []byte, httpDetails httputils.HttpClientDetails, err error) {
	httpDetails = us.ServiceDetails.CreateHttpClientDetails()
	requestContent, err = json.Marshal(user)
	if errorutils.CheckError(err) != nil {
		return
	}
	httpDetails.Headers = map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}
	return
}

func (us *UserService) Delete(username string) error {
	httpDetails := us.ServiceDetails.CreateHttpClientDetails()
	url := fmt.Sprintf("%s/%s", us.getBaseUrl(), username)
	resp, body, err := us.client.SendDelete(url, nil, &httpDetails)
	if err != nil {
		return err
	}
	if resp == nil {
		return errorutils.CheckErrorf("no response provided (including status code)")
	}
	return errorutils.CheckResponseStatusWithBody(resp, body, http.StatusNoContent)
}
