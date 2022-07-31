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

type UserParams struct {
	UserDetails     User
	ReplaceIfExists bool
	ClearGroups     bool
}

func NewUserParams() UserParams {
	return UserParams{}
}

// application/vnd.org.jfrog.artifactory.security.User+json
type User struct {
	Name                     string    `json:"name,omitempty" csv:"username,omitempty"`
	Email                    string    `json:"email,omitempty" csv:"email,omitempty"`
	Password                 string    `json:"password,omitempty" csv:"password,omitempty"`
	Admin                    *bool     `json:"admin,omitempty" csv:"admin,omitempty"`
	ProfileUpdatable         *bool     `json:"profileUpdatable,omitempty" csv:"profileUpdatable,omitempty"`
	DisableUIAccess          *bool     `json:"disableUIAccess,omitempty" csv:"disableUIAccess,omitempty"`
	InternalPasswordDisabled *bool     `json:"internalPasswordDisabled,omitempty" csv:"internalPasswordDisabled,omitempty"`
	LastLoggedIn             string    `json:"lastLoggedIn,omitempty" csv:"lastLoggedIn,omitempty"`
	Realm                    string    `json:"realm,omitempty" csv:"realm,omitempty"`
	Groups                   *[]string `json:"groups,omitempty" csv:"groups,omitempty"`
	ShouldInvite             *bool     `json:"shouldInvite,omitempty" csv:"shouldInvite,omitempty"`
	Source                   string    `json:"source,omitempty" csv:"source,omitempty"`
	WatchManager             *bool     `json:"watchManager,omitempty" csv:"watchManager,omitempty"`
	ReportsManager           *bool     `json:"reportsManager,omitempty" csv:"reportsManager,omitempty"`
	PolicyManager            *bool     `json:"policyManager,omitempty" csv:"policyManager,omitempty"`
	ProjectAdmin             *bool     `json:"projectAdmin,omitempty" csv:"projectAdmin,omitempty"`
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

func (us *UserService) GetUser(params UserParams) (u *User, err error) {
	httpDetails := us.ArtDetails.CreateHttpClientDetails()
	url := fmt.Sprintf("%sapi/security/users/%s", us.ArtDetails.GetUrl(), params.UserDetails.Name)
	resp, body, _, err := us.client.SendGet(url, true, &httpDetails)
	if err != nil {
		return nil, err
	}
	// The case the requested user is not found
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return nil, err
	}
	var user User
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, errorutils.CheckError(err)
	}
	return &user, nil
}

func (us *UserService) GetAllUsers() ([]*User, error) {
	httpDetails := us.ArtDetails.CreateHttpClientDetails()
	url := fmt.Sprintf("%sapi/security/users", us.ArtDetails.GetUrl())
	resp, body, _, err := us.client.SendGet(url, true, &httpDetails)
	if err != nil {
		return nil, err
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return nil, err
	}
	var users []*User
	if err := json.Unmarshal(body, &users); err != nil {
		return nil, errorutils.CheckError(err)
	}
	return users, nil
}

func (us *UserService) CreateUser(params UserParams) error {
	// Checks if the user already exist and act according to ReplaceIfExists parameter.
	if !params.ReplaceIfExists {
		user, err := us.GetUser(params)
		if err != nil {
			return err
		}
		if user != nil {
			return errorutils.CheckErrorf("user '%s' already exists", user.Name)
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
	return errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK, http.StatusCreated)
}

func (us *UserService) UpdateUser(params UserParams) error {
	if params.ClearGroups {
		params.UserDetails.Groups = &[]string{}
	}
	url, content, httpDetails, err := us.createOrUpdateUserRequest(params.UserDetails)
	if err != nil {
		return err
	}
	resp, body, err := us.client.SendPost(url, content, &httpDetails)
	if err != nil {
		return err
	}
	return errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK, http.StatusCreated)
}

func (us *UserService) createOrUpdateUserRequest(user User) (url string, requestContent []byte, httpDetails httputils.HttpClientDetails, err error) {
	httpDetails = us.ArtDetails.CreateHttpClientDetails()
	requestContent, err = json.Marshal(user)
	if errorutils.CheckError(err) != nil {
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
	resp, body, err := us.client.SendDelete(url, nil, &httpDetails)
	if err != nil {
		return err
	}
	if resp == nil {
		return errorutils.CheckErrorf("no response provided (including status code)")
	}
	return errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK)
}
