package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
)

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

func (us *UserService) GetUser(name string) (user *User, notExists bool, err error) {
	httpDetails := us.ArtDetails.CreateHttpClientDetails()
	url := fmt.Sprintf("%sapi/security/users/%s", us.ArtDetails.GetUrl(), name)
	res, body, _, err := us.client.SendGet(url, true, &httpDetails)
	if err != nil {
		return nil, false, err
	}
	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusNotFound {
			return nil, true, err
		}
		return nil, false, fmt.Errorf("%d %s: %s", res.StatusCode, res.Status, string(body))
	}
	if err := json.Unmarshal(body, user); err != nil {
		return nil, false, errorutils.CheckError(err)
	}
	return user, false, nil
}

func (us *UserService) CreateOrUpdateUser(user User) error {
	httpDetails := us.ArtDetails.CreateHttpClientDetails()
	content, err := json.Marshal(user)
	if err != nil {
		return err
	}

	httpDetails.Headers = map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	url := fmt.Sprintf("%sapi/security/users/%s", us.ArtDetails.GetUrl(), user.Name)
	resp, body, err := us.client.SendPut(url, content, &httpDetails)
	if err != nil {
		return err
	}
	if resp.StatusCode > http.StatusNoContent {
		return fmt.Errorf("%d %s: %s", resp.StatusCode, resp.Status, string(body))
	}
	return nil
}

func (us *UserService) DeleteUser(name string) error {
	httpDetails := us.ArtDetails.CreateHttpClientDetails()
	url := fmt.Sprintf("%sapi/security/users/%s", us.ArtDetails.GetUrl(), name)
	resp, _, err := us.client.SendDelete(url, nil, &httpDetails)
	if resp == nil {
		return fmt.Errorf("no response provided (including status code)")
	}
	if resp.StatusCode > http.StatusNoContent {
		return fmt.Errorf("%d %s", resp.StatusCode, resp.Status)
	}
	return err
}
