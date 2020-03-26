package auth

import (
	"sync"
	"time"

	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
)

var expiryHandleMutex sync.Mutex

// Implement this function and append it to create an interceptor that will run pre request in the http client
type PreRequestInterceptorFunc func(*CommonConfigFields, *httputils.HttpClientDetails) error

type CommonDetails interface {
	GetUrl() string
	GetUser() string
	GetPassword() string
	GetApiKey() string
	GetAccessToken() string
	GetPreRequestInterceptor() []PreRequestInterceptorFunc
	GetClientCertPath() string
	GetClientCertKeyPath() string
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
	AppendPreRequestInterceptor(PreRequestInterceptorFunc)
	SetClientCertPath(certificatePath string)
	SetClientCertKeyPath(certificatePath string)
	SetSshUrl(url string)
	SetSshKeyPath(sshKeyPath string)
	SetSshPassphrase(sshPassphrase string)
	SetSshAuthHeaders(sshAuthHeaders map[string]string)

	IsSshAuthHeaderSet() bool
	IsSshAuthentication() bool
	AuthenticateSsh(sshKey, sshPassphrase string) error
	RunPreRequestInterceptors(httpClientDetails *httputils.HttpClientDetails) error

	CreateHttpClientDetails() httputils.HttpClientDetails
}

type CommonConfigFields struct {
	Url                    string                      `json:"-"`
	User                   string                      `json:"-"`
	Password               string                      `json:"-"`
	ApiKey                 string                      `json:"-"`
	AccessToken            string                      `json:"-"`
	PreRequestInterceptors []PreRequestInterceptorFunc `json:"-"`
	ClientCertPath         string                      `json:"-"`
	ClientCertKeyPath      string                      `json:"-"`
	Version                string                      `json:"-"`
	SshUrl                 string                      `json:"-"`
	SshKeyPath             string                      `json:"-"`
	SshPassphrase          string                      `json:"-"`
	SshAuthHeaders         map[string]string           `json:"-"`
	TokenMutex             sync.Mutex
}

func (ds *CommonConfigFields) GetUrl() string {
	return ds.Url
}

func (ds *CommonConfigFields) GetUser() string {
	return ds.User
}

func (ds *CommonConfigFields) GetPassword() string {
	return ds.Password
}

func (ds *CommonConfigFields) GetApiKey() string {
	return ds.ApiKey
}

func (ds *CommonConfigFields) GetAccessToken() string {
	return ds.AccessToken
}

func (ds *CommonConfigFields) GetPreRequestInterceptor() []PreRequestInterceptorFunc {
	return ds.PreRequestInterceptors
}

func (ds *CommonConfigFields) GetClientCertPath() string {
	return ds.ClientCertPath
}

func (ds *CommonConfigFields) GetClientCertKeyPath() string {
	return ds.ClientCertKeyPath
}

func (ds *CommonConfigFields) GetSshUrl() string {
	return ds.SshUrl
}

func (ds *CommonConfigFields) GetSshKeyPath() string {
	return ds.SshKeyPath
}

func (ds *CommonConfigFields) GetSshPassphrase() string {
	return ds.SshPassphrase
}

func (ds *CommonConfigFields) GetSshAuthHeaders() map[string]string {
	return ds.SshAuthHeaders
}

func (ds *CommonConfigFields) SetUrl(url string) {
	ds.Url = url
}

func (ds *CommonConfigFields) SetUser(user string) {
	ds.User = user
}

func (ds *CommonConfigFields) SetPassword(password string) {
	ds.Password = password
}

func (ds *CommonConfigFields) SetApiKey(apiKey string) {
	ds.ApiKey = apiKey
}

func (ds *CommonConfigFields) SetAccessToken(accessToken string) {
	ds.AccessToken = accessToken
}

func (ds *CommonConfigFields) AppendPreRequestInterceptor(interceptor PreRequestInterceptorFunc) {
	ds.PreRequestInterceptors = append(ds.PreRequestInterceptors, interceptor)
}

func (ds *CommonConfigFields) SetClientCertPath(certificatePath string) {
	ds.ClientCertPath = certificatePath
}

func (ds *CommonConfigFields) SetClientCertKeyPath(certificatePath string) {
	ds.ClientCertKeyPath = certificatePath
}

func (ds *CommonConfigFields) SetSshUrl(sshUrl string) {
	ds.SshUrl = sshUrl
}

func (ds *CommonConfigFields) SetSshKeyPath(sshKeyPath string) {
	ds.SshKeyPath = sshKeyPath
}

func (ds *CommonConfigFields) SetSshPassphrase(sshPassphrase string) {
	ds.SshPassphrase = sshPassphrase
}

func (ds *CommonConfigFields) SetSshAuthHeaders(sshAuthHeaders map[string]string) {
	ds.SshAuthHeaders = sshAuthHeaders
}

func (ds *CommonConfigFields) IsSshAuthHeaderSet() bool {
	return len(ds.SshAuthHeaders) > 0
}

func (ds *CommonConfigFields) IsSshAuthentication() bool {
	return fileutils.IsSshUrl(ds.Url) || ds.SshUrl != ""
}

func (ds *CommonConfigFields) AuthenticateSsh(sshKeyPath, sshPassphrase string) error {
	// If SshUrl is unset, set it and use it to authenticate.
	// The SshUrl variable could be used again later if there's a need to reauthenticate (Url is being overwritten with baseUrl).
	if ds.SshUrl == "" {
		ds.SshUrl = ds.Url
	}

	sshHeaders, baseUrl, err := SshAuthentication(ds.SshUrl, sshKeyPath, sshPassphrase)
	if err != nil {
		return err
	}

	// Set base url as the connection url
	ds.Url = baseUrl
	ds.SetSshAuthHeaders(sshHeaders)
	return nil
}

// Runs an interceptor before sending a request via the http client
func (ds *CommonConfigFields) RunPreRequestInterceptors(httpClientDetails *httputils.HttpClientDetails) error {
	for _, exec := range ds.PreRequestInterceptors {
		err := exec(ds, httpClientDetails)
		if err != nil {
			return err
		}
	}
	return nil
}

// Handles the process of acquiring a new ssh token
func SshTokenRefreshPreRequestInterceptor(fields *CommonConfigFields, httpClientDetails *httputils.HttpClientDetails) error {
	if !fields.IsSshAuthentication() {
		return nil
	}
	curToken := httpClientDetails.Headers["Authorization"]
	timeLeft, err := GetTokenMinutesLeft(curToken)
	if err != nil || timeLeft > RefreshBeforeExpiryMinutes {
		return err
	}

	// Lock expiryHandleMutex to make sure only one authentication is made
	expiryHandleMutex.Lock()
	defer expiryHandleMutex.Unlock()
	// Reauthenticate only if a new token wasn't acquired (by another thread) while waiting at mutex.
	if fields.GetSshAuthHeaders()["Authorization"] == curToken {
		// If token isn't already expired, Wait to make sure requests using the current token are sent before it is refreshed and becomes invalid
		if timeLeft != 0 {
			time.Sleep(WaitBeforeRefreshSeconds * time.Second)
		}

		// Obtain a new token and return true (false for error).
		err := fields.AuthenticateSsh(fields.GetSshKeyPath(), fields.GetSshPassphrase())
		if err != nil {
			return err
		}
	}

	// Copy new token from the mutual headers map in CommonDetails to the private headers map in httpClientDetails
	utils.MergeMaps(fields.GetSshAuthHeaders(), httpClientDetails.Headers)
	return nil
}

func (ds *CommonConfigFields) CreateHttpClientDetails() httputils.HttpClientDetails {
	return httputils.HttpClientDetails{
		User:        ds.User,
		Password:    ds.Password,
		ApiKey:      ds.ApiKey,
		AccessToken: ds.AccessToken,
		Headers:     utils.CopyMap(ds.GetSshAuthHeaders())}
}
