package tests

import (
	"github.com/jfrog/jfrog-client-go/access/services"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAccessTokens(t *testing.T) {
	initAccessTest(t)
	t.Run("invite", testCreateToken)
}

func testCreateToken(t *testing.T) {
	tokenParams := services.TokenParams{}
	tokenParams.ExpiresIn = 0
	token, err := testsAccessTokensService.CreateAccessToken(tokenParams)
	assert.NoError(t, err)
	assert.NotEqual(t, "", token.AccessToken, "Access token is empty")
	//TODO: check why
	assert.Equal(t, 31536000, token.ExpiresIn)
}
