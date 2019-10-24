package tests

import (
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"testing"
)

// Teardown should revoke these tokens
var tokensToRevoke []string

const tokenRevokeSuccessResponse = "Token revoked"
const tokenNotFoundResponse = "Token not found"

func TestToken(t *testing.T) {
	t.Run("CreateToken", createTokenTest)
	t.Run("RevokeToken", revokeTokenTest)
	t.Run("RevokeToken: token not found", revokeTokenNotFoundTest)
	t.Run("RefreshToken", refreshTokenTest)
	t.Run("GetTokens", getTokensTest)
	teardown()
}

func createTokenTest(t *testing.T) {
	token, err := createToken()
	if err != nil {
		t.Error(err)
	}
	if token.AccessToken == "" {
		t.Error("Failed to create access token")
	}
	tokensToRevoke = append(tokensToRevoke, token.RefreshToken)
}

func revokeTokenTest(t *testing.T) {
	token, err := createToken()
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
	token, err := createToken()
	if err != nil {
		t.Error(err)
	}
	params := services.NewRefreshTokenParams()
	params.RefreshToken = token.RefreshToken
	params.AccessToken = token.AccessToken
	newToken, err := testsSecurityService.RefreshToken(params)
	if err != nil {
		t.Error(err)
	}
	tokensToRevoke = append(tokensToRevoke, newToken.RefreshToken)
}

func getTokensTest(t *testing.T) {
	token, err := createToken()
	if err != nil {
		t.Error(err)
	}
	tokens, err := testsSecurityService.GetTokens()
	if err != nil {
		t.Error(err)
	}
	if len(tokens.Tokens) < 1 {
		if err != nil {
			t.Error("Failed to get tokens")
		}
	}
	tokensToRevoke = append(tokensToRevoke, token.RefreshToken)
}

// Util function to revoke a token
func revokeToken(token string) (string, error) {
	params := services.NewRevokeTokenParams()
	params.Token = token
	return testsSecurityService.RevokeToken(params)
}

// Util function to create a token
func createToken() (services.CreateTokenResponseData, error) {
	params := services.NewCreateTokenParams()
	params.Username = "anonymous"
	params.Scope = "api:* member-of-groups:readers"
	params.Refreshable = true // We need to use the refresh token to revoke these tokens on teardown
	params.Audience = "jfrt@*" // Allow token to be accepted by all instances of Artifactory.
	return testsSecurityService.CreateToken(params)
}

func revokeAllTokens() {
	for _, element := range tokensToRevoke {
		log.Debug("Revoking Token: ", element)
		responseText, err := revokeToken(element)
		if err != nil {
			log.Error(err)
		}
		if responseText != tokenRevokeSuccessResponse {
			log.Error("Token was not revoked: ", responseText)
		}
	}
}

func teardown() {
	log.Info("REVOKING ALL ", len(tokensToRevoke), " tokens")
	revokeAllTokens()
}
