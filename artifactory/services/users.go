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

type UsersParams struct {
	UserDetails       User
	ReplaceExistUsers bool
}

func NewUsersParams() UsersParams {
	return UsersParams{}
}

// application/vnd.org.jfrog.artifactory.security.User+json
type User struct {
	Name                     string   `json:"name,omitempty"`
	Email                    string   `json:"email,omitempty"`
	Password                 string   `json:"password,omitempty"`
	Admin                    bool     `json:"admin,omitempty"`
	ProfileUpdatable         bool     `json:"profileUpdatable,omitempty"`
	DisableUIAccess          bool     `json:"disableUIAccess,omitempty"`
	InternalPasswordDisabled bool     `json:"internalPasswordDisabled,omitempty"`
	LastLoggedIn             string   `json:"lastLoggedIn,omitempty"`
	Realm                    string   `json:"realm,omitempty"`
	Groups                   []string `json:"groups,omitempty"`
}

type UserService struct {
	client     *jfroghttpclient.JfrogHttpClient
	ArtDetails auth.ServiceDetails
}

func NewUserService(client *jfroghttpclient.JfrogHttpClient) *UserService {
	return &UserService{client: client}
}

func (us *UserService) SetArtifactoryDetails(rt auth.ServiceDetails) {
	us.ArtDetails = rt
}

func (us *UserService) GetUser(params UsersParams) (u *User, notExists bool, err error) {
	httpDetails := us.ArtDetails.CreateHttpClientDetails()
	url := fmt.Sprintf("%sapi/security/users/%s", us.ArtDetails.GetUrl(), params.UserDetails.Name)
	res, body, _, err := us.client.SendGet(url, true, &httpDetails)
	if err != nil {
		return nil, false, err
	}
	if res.StatusCode != http.StatusOK {
		// The case the requseted user is not found
		if res.StatusCode == http.StatusNotFound {
			return nil, true, err
		}
		return nil, false, fmt.Errorf("%d %s: %s", res.StatusCode, res.Status, string(body))
	}
	var user User
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, false, errorutils.CheckError(err)
	}
	return &user, false, nil
}

func (us *UserService) CreateUser(params UsersParams) error {
	user := params.UserDetails
	// Checks if the user allready exists in the system and act according to replaceExistUsers parameter.
	if !params.ReplaceExistUsers {
		_, notExists, err := us.GetUser(params)
		if err != nil {
			return err
		}
		if !notExists {
			return fmt.Errorf("User %s is allready exists in the system", user.Name)
		}
	}
	url, content, httpDetails, err := us.createOrUpdateUserRequest(params.UserDetails)
	if err != nil {
		return err
	}
	resp, body, err := us.client.SendPut(url, content, &httpDetails)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("%d %s: %s", resp.StatusCode, resp.Status, string(body))
	}
	return nil
}

func (us *UserService) UpdateUser(params UsersParams) error {
	url, content, httpDetails, err := us.createOrUpdateUserRequest(params.UserDetails)
	if err != nil {
		return err
	}
	resp, body, err := us.client.SendPost(url, content, &httpDetails)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("%d %s: %s", resp.StatusCode, resp.Status, string(body))
	}
	return nil
}

func (us *UserService) createOrUpdateUserRequest(user User) (url string, requestContent []byte, httpDetails httputils.HttpClientDetails, err error) {
	httpDetails = us.ArtDetails.CreateHttpClientDetails()
	requestContent, err = json.Marshal(user)
	if err != nil {
		return
	}

	httpDetails.Headers = map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	url = fmt.Sprintf("%sapi/security/users/%s", us.ArtDetails.GetUrl(), user.Name)
	return
}

func (us *UserService) DeleteUser(name string) error {
	httpDetails := us.ArtDetails.CreateHttpClientDetails()
	url := fmt.Sprintf("%sapi/security/users/%s", us.ArtDetails.GetUrl(), name)
	resp, _, err := us.client.SendDelete(url, nil, &httpDetails)
	if resp == nil {
		return fmt.Errorf("no response provided (including status code)")
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%d %s", resp.StatusCode, resp.Status)
	}
	return err
}
