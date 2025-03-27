package tests

import (
	"encoding/json"
	accessAuth "github.com/jfrog/jfrog-client-go/access/auth"
	"github.com/jfrog/jfrog-client-go/access/services"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

const testExpiredInSeconds = 1

func TestAccessTokens(t *testing.T) {
	initAccessTest(t)
	t.Run("createAccessToken", testCreateRefreshableToken)
	t.Run("createAccessTokenWithReference", testAccessTokenWithReference)
	t.Run("refreshToken", testRefreshTokenTest)
	t.Run("exchangeOIDCToken", testExchangeOidcToken)
}

// Mocks exchange response from the server as it requires setting up a full OIDC flow
// which currently is not supported by CLI commands.
func testExchangeOidcToken(t *testing.T) {
	initAccessTest(t)

	// Create mock server
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/access/api/v1/oidc/token", r.URL.Path)

		// Verify request body
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		defer func() {
			assert.NoError(t, r.Body.Close())
		}()

		var req services.CreateOidcTokenParams
		err = json.Unmarshal(body, &req)
		assert.NoError(t, err)
		assert.Equal(t, "mockOidcTokenID", req.OidcTokenID)

		// Simulate response
		resp := auth.OidcTokenResponseData{
			CommonTokenParams: auth.CommonTokenParams{
				AccessToken: "mockAccessToken",
			},
			IssuedTokenType: "mockIssuedTokenType",
			Username:        "mockUsername",
		}
		responseBody, err := json.Marshal(resp)
		assert.NoError(t, err)

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(responseBody)
		assert.NoError(t, err)
	})
	ts := httptest.NewServer(handler)
	defer ts.Close()

	// Setup JFrog client
	client, err := jfroghttpclient.JfrogClientBuilder().
		SetInsecureTls(true).
		Build()
	assert.NoError(t, err, "Failed to create JFrog client")

	// Setup TokenService
	service := services.NewTokenService(client)
	serverDetails := accessAuth.NewAccessDetails()
	serverDetails.SetUrl(ts.URL + "/access")
	service.ServiceDetails = serverDetails

	// Define OIDC token parameters
	params := services.CreateOidcTokenParams{
		GrantType:        "authorization_code",
		SubjectTokenType: "Generic",
		OidcTokenID:      "mockOidcTokenID",
		ProviderName:     "mockProviderName",
		ProjectKey:       "mockProjectKey",
		JobId:            "mockJobId",
		RunId:            "mockRunId",
		Repo:             "mockRepo",
		ApplicationKey:   "mockApplicationKey",
		Audience:         "mockAudience",
	}

	// Execute ExchangeOidcToken
	response, err := service.ExchangeOidcToken(params)

	// Verify response
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func testCreateRefreshableToken(t *testing.T) {
	tokenParams := createRefreshableAccessTokenParams(testExpiredInSeconds)
	token, err := testsAccessTokensService.CreateAccessToken(tokenParams)
	assert.NoError(t, err)
	assert.NotEqual(t, "", token.AccessToken, "Access token is empty")
	assert.NotEqual(t, tokenParams.AccessToken, token.AccessToken, "New access token is identical to original one")
	assert.NotEqual(t, "", token.RefreshToken, "Refresh token is empty")
	assert.EqualValues(t, testExpiredInSeconds, *token.ExpiresIn)
	assert.Empty(t, token.ReferenceToken)
}

func testAccessTokenWithReference(t *testing.T) {
	tokenParams := createRefreshableAccessTokenParams(testExpiredInSeconds)
	tokenParams.IncludeReferenceToken = utils.Pointer(true)
	token, err := testsAccessTokensService.CreateAccessToken(tokenParams)
	assert.NoError(t, err)
	assert.NotEqual(t, "", token.AccessToken, "Access token is empty")
	assert.NotEqual(t, tokenParams.AccessToken, token.AccessToken, "New access token is identical to original one")
	assert.NotEqual(t, "", token.RefreshToken, "Refresh token is empty")
	assert.EqualValues(t, testExpiredInSeconds, *token.ExpiresIn)
	assert.NotEmpty(t, token.ReferenceToken)
}

func testRefreshTokenTest(t *testing.T) {
	// Create token
	tokenParams := createRefreshableAccessTokenParams(testExpiredInSeconds)
	token, err := testsAccessTokensService.CreateAccessToken(tokenParams)
	assert.NoError(t, err)
	// Refresh token
	refreshTokenParams := createRefreshAccessTokenParams(token)
	newToken, err := testsAccessTokensService.RefreshAccessToken(refreshTokenParams)
	assert.NoError(t, err)
	// Validate
	assert.NotEqual(t, token.AccessToken, newToken.AccessToken, "New access token is identical to original one")
	assert.NotEqual(t, token.RefreshToken, newToken.RefreshToken, "New refresh token is identical to original one")
	assert.EqualValues(t, token.ExpiresIn, newToken.ExpiresIn, "New access token's expiration is different from original one")
	assert.Empty(t, token.ReferenceToken)
}

func createRefreshableAccessTokenParams(expiredIn uint) services.CreateTokenParams {
	tokenParams := services.CreateTokenParams{}
	tokenParams.ExpiresIn = &expiredIn
	tokenParams.Refreshable = utils.Pointer(true)
	tokenParams.Audience = "*@*"
	return tokenParams
}

func createRefreshAccessTokenParams(token auth.CreateTokenResponseData) (refreshParams services.CreateTokenParams) {
	refreshParams = services.CreateTokenParams{}
	refreshParams.ExpiresIn = token.ExpiresIn
	refreshParams.Refreshable = utils.Pointer(true)
	refreshParams.GrantType = "refresh_token"
	refreshParams.TokenType = "Bearer"
	refreshParams.RefreshToken = token.RefreshToken
	refreshParams.AccessToken = token.AccessToken
	return
}
