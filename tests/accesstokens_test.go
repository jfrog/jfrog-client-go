package tests

import (
	"github.com/jfrog/jfrog-client-go/access/services"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/stretchr/testify/assert"
	"testing"
)

const testExpiredInSeconds = 1

func TestAccessTokens(t *testing.T) {
	initAccessTest(t)
	t.Run("createAccessToken", testCreateRefreshableToken)
	t.Run("createAccessTokenWithReference", testAccessTokenWithReference)
	t.Run("refreshToken", testRefreshTokenTest)
}

func testCreateRefreshableToken(t *testing.T) {
	tokenParams := createRefreshableAccessTokenParams(testExpiredInSeconds)
	token, err := testsAccessTokensService.CreateAccessToken(tokenParams)
	assert.NoError(t, err)
	assert.NotEqual(t, "", token.AccessToken, "Access token is empty")
	assert.NotEqual(t, tokenParams.AccessToken, token.AccessToken, "New access token is identical to original one")
	assert.NotEqual(t, "", token.RefreshToken, "Refresh token is empty")
	assert.Equal(t, testExpiredInSeconds, token.ExpiresIn)
	assert.Empty(t, token.ReferenceToken)
}

func testAccessTokenWithReference(t *testing.T) {
	tokenParams := createRefreshableAccessTokenParams(testExpiredInSeconds)
	tokenParams.IncludeReferenceToken = &trueValue
	token, err := testsAccessTokensService.CreateAccessToken(tokenParams)
	assert.NoError(t, err)
	assert.NotEqual(t, "", token.AccessToken, "Access token is empty")
	assert.NotEqual(t, tokenParams.AccessToken, token.AccessToken, "New access token is identical to original one")
	assert.NotEqual(t, "", token.RefreshToken, "Refresh token is empty")
	assert.Equal(t, testExpiredInSeconds, token.ExpiresIn)
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
	assert.Equal(t, token.ExpiresIn, newToken.ExpiresIn, "New access token's expiration is different from original one")
	assert.Empty(t, token.ReferenceToken)
}

func createRefreshableAccessTokenParams(expiredIn int) services.CreateTokenParams {
	tokenParams := services.CreateTokenParams{}
	tokenParams.ExpiresIn = expiredIn
	tokenParams.Refreshable = &trueValue
	tokenParams.Audience = "*@*"
	return tokenParams
}

func createRefreshAccessTokenParams(token auth.CreateTokenResponseData) (refreshParams services.CreateTokenParams) {
	refreshParams = services.CreateTokenParams{}
	refreshParams.ExpiresIn = token.ExpiresIn
	refreshParams.Refreshable = &trueValue
	refreshParams.GrantType = "refresh_token"
	refreshParams.TokenType = "Bearer"
	refreshParams.RefreshToken = token.RefreshToken
	refreshParams.AccessToken = token.AccessToken
	return
}
