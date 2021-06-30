package tests

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

// Teardown should revoke these tokens
var tokensToRevoke []string

const tokenRevokeSuccessResponse = "Token revoked"
const tokenNotFoundResponse = "Token not found"

func TestToken(t *testing.T) {
	initArtifactoryTest(t)
	t.Run("CreateToken", createTokenTest)
	t.Run("RevokeToken", revokeTokenTest)
	t.Run("RevokeToken: token not found", revokeTokenNotFoundTest)
	t.Run("RefreshToken", refreshTokenTest)
	t.Run("GetTokens", getTokensTest)
	teardown()
}

func createArtifactoryTLSServer(t *testing.T, expectedRequest string, expectedStatusCode int) *httptest.Server {
	returnValue := `{}`
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, expectedRequest, r.Method)
		w.WriteHeader(expectedStatusCode)
		_, err := w.Write([]byte(returnValue))
		assert.NoError(t, err)
	})
	return httptest.NewTLSServer(handler)
}

func createDummySecurityService(tlsUrl string) (*services.SecurityService, error) {
	rtDetails := auth.NewArtifactoryDetails()
	rtDetails.SetUrl(tlsUrl + "/")
	rtDetails.SetUser("fake-user")
	rtDetails.SetPassword("fake-pass")

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

func getUserTokenstest(t *testing.T) {
	token, err := createToken()
	if err != nil {
		t.Error(err)
	}
	tokens, err := testsSecurityService.GetUserTokens("anonymous")
	if len(tokens) != 1 {
		t.Error("Failed to get tokens of anonymous user")
	}
	if tokens[0] != token.AccessToken {
		t.Error("Retried user token doesn't match expected token value")
	}

	tokensToRevoke = append(tokensToRevoke, token.RefreshToken)

	params := services.NewCreateTokenParams()
	params.Username = "test-user"
	params.Scope = "api:* member-of-groups:readers"
	params.Refreshable = true  // We need to use the refresh token to revoke these tokens on teardown
	params.Audience = "jfrt@*" // Allow token to be accepted by all instances of Artifactory.

	token1, err := testsSecurityService.CreateToken(params)
	if err != nil {
		t.Error(err)
	}

	token2, err := testsSecurityService.CreateToken(params)
	if err != nil {
		t.Error(err)
	}
	tokens, err = testsSecurityService.GetUserTokens("test-user")
	if len(tokens) != 2 {
		t.Error("Failed to get tokens of test-user")
	}
	if tokens[0] != token1.AccessToken || tokens[1] != token2.AccessToken {
		t.Error("Retried user token doesn't match expected token value")
	}
	tokensToRevoke = append(tokensToRevoke, token1.RefreshToken)
	tokensToRevoke = append(tokensToRevoke, token2.RefreshToken)
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
	params.Refreshable = true  // We need to use the refresh token to revoke these tokens on teardown
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
