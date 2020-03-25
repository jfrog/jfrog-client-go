package auth

import (
	"sync"
	"time"

	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
)

var expiryHandleMutex sync.Mutex

// Implement this function and pass it to the CommonDetails struct to handle access token refresh if needed
type TokenRefreshHandlerFunc func(currentAccessToken string) (newAccessToken string, err error)

type CommonDetails interface {
	GetUrl() string
	GetUser() string
	GetPassword() string
	GetApiKey() string
	GetAccessToken() string
	GetTokenRefreshHandler() TokenRefreshHandlerFunc
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
	SetTokenRefreshHandler(TokenRefreshHandlerFunc)
	SetClientCertPath(certificatePath string)
	SetClientCertKeyPath(certificatePath string)
	SetSshUrl(url string)
	SetSshKeyPath(sshKeyPath string)
	SetSshPassphrase(sshPassphrase string)
	SetSshAuthHeaders(sshAuthHeaders map[string]string)

	IsSshAuthHeaderSet() bool
	IsSshAuthentication() bool
	AuthenticateSsh(sshKey, sshPassphrase string) error
	HandleTokenExpiry(httpClientDetails *httputils.HttpClientDetails) error

	CreateHttpClientDetails() httputils.HttpClientDetails
}

type CommonConfigFields struct {
	Url                 string                  `json:"-"`
	User                string                  `json:"-"`
	Password            string                  `json:"-"`
	ApiKey              string                  `json:"-"`
	AccessToken         string                  `json:"-"`
	TokenRefreshHandler TokenRefreshHandlerFunc `json:"-"`
	ClientCertPath      string                  `json:"-"`
	ClientCertKeyPath   string                  `json:"-"`
	Version             string                  `json:"-"`
	SshUrl              string                  `json:"-"`
	SshKeyPath          string                  `json:"-"`
	SshPassphrase       string                  `json:"-"`
	SshAuthHeaders      map[string]string       `json:"-"`
	TokenMutex          sync.Mutex
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

func (ds *CommonConfigFields) GetTokenRefreshHandler() TokenRefreshHandlerFunc {
	return ds.TokenRefreshHandler
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

func (ds *CommonConfigFields) SetTokenRefreshHandler(tokenRefreshHandler TokenRefreshHandlerFunc) {
	ds.TokenRefreshHandler = tokenRefreshHandler
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

// Checks if a token expired or about to expire
// If so, acquire a new token from server (if one wasn't acquired yet)
func (ds *CommonConfigFields) HandleTokenExpiry(httpClientDetails *httputils.HttpClientDetails) error {
	if ds.IsSshAuthentication() {
		return ds.handleSshTokenExpiry(httpClientDetails)
	}
	if ds.GetAccessToken() != "" && ds.GetTokenRefreshHandler() != nil {
		return ds.handleAccessTokenExpiry(httpClientDetails)
	}
	return nil
}

// Handles the process of acquiring a new ssh token
func (ds *CommonConfigFields) handleSshTokenExpiry(httpClientDetails *httputils.HttpClientDetails) error {
	curToken := httpClientDetails.Headers["Authorization"]
	timeLeft, err := GetTokenMinutesLeft(curToken)
	if err != nil || timeLeft > refreshBeforeExpiryMinutes {
		return err
	}

	// Lock expiryHandleMutex to make sure only one authentication is made
	expiryHandleMutex.Lock()
	// Reauthenticate only if a new token wasn't acquired (by another thread) while waiting at mutex.
	if ds.GetSshAuthHeaders()["Authorization"] == curToken {
		// If token isn't already expired, Wait to make sure requests using the current token are sent before it is refreshed and becomes invalid
		if timeLeft != 0 {
			time.Sleep(WaitBeforeRefreshSeconds * time.Second)
		}

		// Obtain a new token and return true (false for error).
		err := ds.AuthenticateSsh(ds.GetSshKeyPath(), ds.GetSshPassphrase())
		if err != nil {
			expiryHandleMutex.Unlock()
			return err
		}
	}
	expiryHandleMutex.Unlock()

	// Copy new token from the mutual headers map in CommonDetails to the private headers map in httpClientDetails
	utils.MergeMaps(ds.GetSshAuthHeaders(), httpClientDetails.Headers)
	return nil
}

func (ds *CommonConfigFields) handleAccessTokenExpiry(httpClientDetails *httputils.HttpClientDetails) error {
	timeLeft, err := GetTokenMinutesLeft(httpClientDetails.AccessToken)
	if err != nil || timeLeft > refreshBeforeExpiryMinutes {
		return err
	}

	// Lock expiryHandleMutex to make sure only one thread is trying to refresh
	expiryHandleMutex.Lock()
	// Refresh only if a new token wasn't acquired (by another thread) while waiting at mutex.
	if ds.AccessToken == httpClientDetails.AccessToken {
		// Call a predefined handler to manage the refresh process
		newAccessToken, err := ds.TokenRefreshHandler(httpClientDetails.AccessToken)
		if err != nil {
			return err
		}
		if newAccessToken != "" && newAccessToken != httpClientDetails.AccessToken {
			ds.AccessToken = newAccessToken
		}
	}
	expiryHandleMutex.Unlock()

	// Copy new token from the mutual struct CommonDetails to the private struct in httpClientDetails
	httpClientDetails.AccessToken = ds.AccessToken
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
