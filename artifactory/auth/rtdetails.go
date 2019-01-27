package auth

import (
	"encoding/json"
	"errors"
	"github.com/jfrog/jfrog-client-go/httpclient"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
	"strings"
	"sync"
)

func NewArtifactoryDetails() ArtifactoryDetails {
	return &artifactoryDetails{}
}

var expiryHandleMutex sync.Mutex

type ArtifactoryDetails interface {
	GetUrl() string
	GetUser() string
	GetPassword() string
	GetApiKey() string
	GetAccessToken() string
	GetSshUrl() string
	GetSshKeyPath() string
	GetSshPassphrase() string
	GetSshAuthHeaders() map[string]string
	GetVersion() (string, error)

	SetUrl(url string)
	SetUser(user string)
	SetPassword(password string)
	SetApiKey(apiKey string)
	SetAccessToken(accessToken string)
	SetSshUrl(url string)
	SetSshKeyPath(sshKeyPath string)
	SetSshPassphrase(sshPassphrase string)
	SetSshAuthHeaders(sshAuthHeaders map[string]string)

	IsSshAuthHeaderSet() bool
	IsSshAuthentication() bool
	AuthenticateSsh(sshKey, sshPassphrase string) error
	HandleTokenExpiry(statusCode int, httpClientDetails *httputils.HttpClientDetails) (bool, error)

	CreateHttpClientDetails() httputils.HttpClientDetails
}

type artifactoryDetails struct {
	Url            string            `json:"-"`
	User           string            `json:"-"`
	Password       string            `json:"-"`
	ApiKey         string            `json:"-"`
	AccessToken    string            `json:"-"`
	version        string            `json:"-"`
	SshUrl         string            `json:"-"`
	SshKeyPath     string            `json:"-"`
	SshPassphrase  string            `json:"-"`
	SshAuthHeaders map[string]string `json:"-"`
	TokenMutex     sync.Mutex
}

func (rt *artifactoryDetails) GetUrl() string {
	return rt.Url
}

func (rt *artifactoryDetails) GetUser() string {
	return rt.User
}

func (rt *artifactoryDetails) GetPassword() string {
	return rt.Password
}

func (rt *artifactoryDetails) GetApiKey() string {
	return rt.ApiKey
}

func (rt *artifactoryDetails) GetAccessToken() string {
	return rt.AccessToken
}

func (rt *artifactoryDetails) GetSshUrl() string {
	return rt.SshUrl
}

func (rt *artifactoryDetails) GetSshKeyPath() string {
	return rt.SshKeyPath
}

func (rt *artifactoryDetails) GetSshPassphrase() string {
	return rt.SshPassphrase
}

func (rt *artifactoryDetails) GetSshAuthHeaders() map[string]string {
	return rt.SshAuthHeaders
}

func (rt *artifactoryDetails) SetUrl(url string) {
	rt.Url = url
}

func (rt *artifactoryDetails) SetUser(user string) {
	rt.User = user
}

func (rt *artifactoryDetails) SetPassword(password string) {
	rt.Password = password
}

func (rt *artifactoryDetails) SetApiKey(apiKey string) {
	rt.ApiKey = apiKey
}

func (rt *artifactoryDetails) SetAccessToken(accessToken string) {
	rt.AccessToken = accessToken
}

func (rt *artifactoryDetails) SetSshUrl(sshUrl string) {
	rt.SshUrl = sshUrl
}

func (rt *artifactoryDetails) SetSshKeyPath(sshKeyPath string) {
	rt.SshKeyPath = sshKeyPath
}

func (rt *artifactoryDetails) SetSshPassphrase(sshPassphrase string) {
	rt.SshPassphrase = sshPassphrase
}

func (rt *artifactoryDetails) SetSshAuthHeaders(sshAuthHeaders map[string]string) {
	rt.SshAuthHeaders = sshAuthHeaders
}

func (rt *artifactoryDetails) IsSshAuthHeaderSet() bool {
	return len(rt.SshAuthHeaders) > 0
}

func (rt *artifactoryDetails) IsSshAuthentication() bool {
	return fileutils.IsSshUrl(rt.Url) || rt.SshUrl != ""
}

func (rt *artifactoryDetails) AuthenticateSsh(sshKeyPath, sshPassphrase string) error {
	// If SshUrl is unset, set it and use it to authenticate.
	// The SshUrl variable could be used again later if there's a need to reauthenticate (Url is being overwritten with baseUrl).
	if rt.SshUrl == "" {
		rt.SshUrl = rt.Url
	}

	sshHeaders, baseUrl, err := sshAuthentication(rt.SshUrl, sshKeyPath, sshPassphrase)
	if err != nil {
		return err
	}

	// Set base url as the connection url
	rt.Url = baseUrl
	rt.SetSshAuthHeaders(sshHeaders)
	return nil
}

// Checks if a token has expired.
// If so, acquires a new token from server (if one wasn't acquired yet) and returns true.
// Otherwise, or in case of an error, returns false.
func (rt *artifactoryDetails) HandleTokenExpiry(statusCode int, httpClientDetails *httputils.HttpClientDetails) (bool, error) {
	// If an unauthorized ssh connection -> ssh token has expired.
	if statusCode == http.StatusUnauthorized && rt.IsSshAuthentication() {
		return rt.handleSshTokenExpiry(httpClientDetails)
	}
	return false, nil
}

// Handles the process of acquiring a new ssh token from server (if one wasn't acquired yet) and returns true.
// Returns false if an error has occurred.
func (rt *artifactoryDetails) handleSshTokenExpiry(httpClientDetails *httputils.HttpClientDetails) (bool, error) {
	// Lock expiryHandleMutex to make sure only one authentication is made
	expiryHandleMutex.Lock()
	// Reauthenticate if a new token wasn't acquired (by another thread) while waiting at mutex.
	// Otherwise, token has already changed -> get new token and return true without authenticating.
	if rt.GetSshAuthHeaders()["Authorization"] == httpClientDetails.Headers["Authorization"] {
		// Obtain a new token and return true (false for error).
		err := rt.AuthenticateSsh(rt.GetSshKeyPath(), rt.GetSshPassphrase())
		if err != nil {
			expiryHandleMutex.Unlock()
			return false, err
		}
	}
	expiryHandleMutex.Unlock()

	// Copy new token from the mutual headers map in artifactoryDetails to the private headers map in httpClientDetails
	utils.MergeMaps(rt.GetSshAuthHeaders(), httpClientDetails.Headers)

	return true, nil
}

func (rt *artifactoryDetails) CreateHttpClientDetails() httputils.HttpClientDetails {
	return httputils.HttpClientDetails{
		User:        rt.User,
		Password:    rt.Password,
		ApiKey:      rt.ApiKey,
		AccessToken: rt.AccessToken,
		Headers:     utils.CopyMap(rt.GetSshAuthHeaders())}
}

func (rt *artifactoryDetails) GetVersion() (string, error) {
	var err error
	if rt.version == "" {
		rt.version, err = rt.getArtifactoryVersion()
		if err != nil {
			return "", err
		}
		log.Debug("The Artifactory version is:", rt.version)
	}
	return rt.version, nil
}

func (rt *artifactoryDetails) getArtifactoryVersion() (string, error) {
	client, err := httpclient.ClientBuilder().Build()
	if err != nil {
		return "", err
	}
	resp, body, _, err := client.SendGet(rt.GetUrl()+"api/system/version", true, rt.CreateHttpClientDetails())
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errorutils.CheckError(errors.New("Artifactory response: " + resp.Status + "\n" + utils.IndentJson(body)))
	}
	var version artifactoryVersion
	err = json.Unmarshal(body, &version)
	if err != nil {
		return "", errorutils.CheckError(err)
	}
	return strings.TrimSpace(version.Version), nil
}

type artifactoryVersion struct {
	Version string `json:"version,omitempty"`
}
