//go:build itest

package tests

import (
	"fmt"
	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	accessauth "github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const tokenRevokeSuccessResponse = "Token revoked"
const tokenNotFoundResponse = "Token not found"

func TestToken(t *testing.T) {
	initArtifactoryTest(t)
	t.Run("CreateToken", createTokenTest)
	t.Run("RevokeToken", revokeTokenTest)
	t.Run("RevokeToken: token not found", revokeTokenNotFoundTest)
	t.Run("RefreshToken", refreshTokenTest)
	t.Run("GetTokens", getTokensTest)
	t.Run("GetUserTokens", getUserTokensTest)
}

func TestAPIKey(t *testing.T) {
	initArtifactoryTest(t)
	t.Run("Create API Key", createAPIKeyTest)
	t.Run("Regenerate API Key", regenerateAPIKeyTest)
	t.Run("Get API Key", getAPIKeyTest)
	t.Run("Get Empty API Key", getEmptyAPIKeyTest)
}

func createAPIKeyTest(t *testing.T) {
	expectedApiKey := "new-api-key"
	tls := createArtifactoryTLSServer(t, http.MethodPost, expectedApiKey, http.StatusCreated)
	defer tls.Close()

	apiKeyService, err := createDummySecurityService(tls.URL, true)
	if !assert.NoError(t, err) {
		return
	}
	key, err := apiKeyService.CreateAPIKey()
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, expectedApiKey, key)
}

func regenerateAPIKeyTest(t *testing.T) {
	expectedApiKey := "new-api-key"
	tls := createArtifactoryTLSServer(t, http.MethodPut, expectedApiKey, http.StatusOK)
	defer tls.Close()

	apiKeyService, err := createDummySecurityService(tls.URL, true)
	if !assert.NoError(t, err) {
		return
	}
	key, err := apiKeyService.RegenerateAPIKey()
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, expectedApiKey, key)
}

func getAPIKeyTest(t *testing.T) {
	expectedApiKey := "new-api-key"
	getAPIKeyTestCore(t, expectedApiKey)
}

// The GetAPIKey service returns empty string if an API Key wasn't generated.
func getEmptyAPIKeyTest(t *testing.T) {
	expectedApiKey := ""
	getAPIKeyTestCore(t, expectedApiKey)
}

func getAPIKeyTestCore(t *testing.T, expectedApiKey string) {
	tls := createArtifactoryTLSServer(t, http.MethodGet, expectedApiKey, http.StatusOK)
	defer tls.Close()

	apiKeyService, err := createDummySecurityService(tls.URL, true)
	if !assert.NoError(t, err) {
		return
	}
	key, err := apiKeyService.GetAPIKey()
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, expectedApiKey, key)
}

func createArtifactoryTLSServer(t *testing.T, expectedRequest, expectedApiKey string, expectedStatusCode int) *httptest.Server {
	returnValue := fmt.Sprintf(`{"apiKey": "%s"}`, expectedApiKey)
	if expectedApiKey == "" {
		returnValue = `{}`
	}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, expectedRequest, r.Method)
		assert.Equal(t, "/"+services.APIKeyPath, r.URL.Path)
		w.WriteHeader(expectedStatusCode)
		_, err := w.Write([]byte(returnValue))
		assert.NoError(t, err)
	})
	return httptest.NewTLSServer(handler)
}

func createDummySecurityService(tlsUrl string, setApiKey bool) (*services.SecurityService, error) {
	rtDetails := auth.NewArtifactoryDetails()
	rtDetails.SetUrl(tlsUrl + "/")
	rtDetails.SetUser("fake-user")

	if setApiKey {
		rtDetails.SetApiKey("fake-api-key")
	} else {
		rtDetails.SetPassword("fake-pass")
	}

	client, err := jfroghttpclient.JfrogClientBuilder().
		SetInsecureTls(true).
		SetClientCertPath(rtDetails.GetClientCertPath()).
		SetClientCertKeyPath(rtDetails.GetClientCertKeyPath()).
		AppendPreRequestInterceptor(rtDetails.RunPreRequestFunctions).
		Build()
	if err != nil {
		return nil, err
	}

	apiKeyService := services.NewSecurityService(client)
	apiKeyService.ArtDetails = rtDetails
	return apiKeyService, nil
}

func createTokenTest(t *testing.T) {
	username := getUniqueUsername()
	token, err := createToken(username)
	if err != nil {
		t.Error(err)
	}
	if token.AccessToken == "" {
		t.Error("Failed to create access token")
	}
	revokeTokenCleanup(t, token.RefreshToken)
}

func revokeTokenTest(t *testing.T) {
	username := getUniqueUsername()
	token, err := createToken(username)
	if err != nil {
		t.Error(err)
	}
	responseText, err := revokeToken(token.RefreshToken)
	if err != nil {
		t.Error(err)
	}
	if responseText != tokenRevokeSuccessResponse {
		t.Error("Token was not revoked: ", responseText)
	}
}

func revokeTokenNotFoundTest(t *testing.T) {
	responseText, err := revokeToken("faketoken")
	if err != nil {
		t.Error(err)
	}
	if responseText != tokenNotFoundResponse {
		t.Error("Expected Response: ", tokenNotFoundResponse, ". Got", responseText)
	}
}

func refreshTokenTest(t *testing.T) {
	username := getUniqueUsername()
	token, err := createToken(username)
	if err != nil {
		t.Error(err)
	}
	params := services.NewArtifactoryRefreshTokenParams()
	params.RefreshToken = token.RefreshToken
	params.AccessToken = token.AccessToken
	newToken, err := testsSecurityService.RefreshToken(params)
	if err != nil {
		t.Error(err)
	}
	revokeTokenCleanup(t, newToken.RefreshToken)
}

func getTokensTest(t *testing.T) {
	username := getUniqueUsername()
	token, err := createToken(username)
	if err != nil {
		t.Error(err)
	}
	defer revokeTokenCleanup(t, token.RefreshToken)
	tokens, err := testsSecurityService.GetTokens()
	if err != nil {
		t.Error(err)
	}
	if len(tokens.Tokens) < 1 {
		if err != nil {
			t.Error("Failed to get tokens")
		}
	}
}

func getUserTokensTest(t *testing.T) {
	username := getUniqueUsername()
	token, err := createToken(username)
	if err != nil {
		t.Error(err)
		return
	}
	tokens, err := testsSecurityService.GetUserTokens(username)
	if err != nil {
		t.Error(err)
		return
	}
	defer revokeTokenCleanup(t, token.RefreshToken)
	assert.Len(t, tokens, 1)

	username2 := username + "-second"
	token1, err := createToken(username2)
	if err != nil {
		t.Error(err)
		return
	}
	defer revokeTokenCleanup(t, token1.RefreshToken)

	token2, err := createToken(username2)
	if err != nil {
		t.Error(err)
		return
	}
	defer revokeTokenCleanup(t, token2.RefreshToken)

	tokens, err = testsSecurityService.GetUserTokens(username2)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Len(t, tokens, 2)
}

// Util function to revoke a token
func revokeToken(token string) (string, error) {
	params := services.NewRevokeTokenParams()
	params.Token = token
	return testsSecurityService.RevokeToken(params)
}

// Util function to create a token
func createToken(username string) (accessauth.CreateTokenResponseData, error) {
	params := services.NewCreateTokenParams()
	params.Username = username
	params.Scope = "api:* member-of-groups:readers"
	params.Refreshable = true  // We need to use the refresh token to revoke these tokens on teardown
	params.Audience = "jfrt@*" // Allow token to be accepted by all instances of Artifactory.
	return testsSecurityService.CreateToken(params)
}

func revokeTokenCleanup(t *testing.T, refreshToken string) {
	log.Debug("Revoking Token with refresh token:", refreshToken)
	responseText, err := revokeToken(refreshToken)
	if err != nil {
		t.Error(err)
	}
	if responseText != tokenRevokeSuccessResponse {
		t.Error("Token was not revoked: ", responseText)
	}
}

func getUniqueUsername() string {
	return getUniqueField("user")
}
