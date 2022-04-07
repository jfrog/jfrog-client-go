package tests

import (
	"github.com/jfrog/jfrog-client-go/access/services"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAccessTokens(t *testing.T) {
	initAccessTest(t)
	t.Run("createAccessToken", testCreateToken)
}

func testCreateToken(t *testing.T) {
	tokenParams := services.TokenParams{}
	tokenParams.ExpiresIn = 60
	tokenParams.Refreshable = &trueValue
	tokenParams.RefreshToken = "e0134455-7d7f-4da5-ba0d-9f04c358ba2d"
	tokenParams.GrantType = "refresh_token"
	tokenParams.TokenType = "Bearer"
	//tokenParams.AccessToken = "eyJ2ZXIiOiIyIiwidHlwIjoiSldUIiwiYWxnIjoiUlMyNTYiLCJraWQiOiJ5TE9LOU9wSkRBbkh1U09YY0ZmQkNGUUVRTmI1clRfSEpKS3dBZWk3dWY4In0.eyJleHQiOiJ7XCJyZXZvY2FibGVcIjpcInRydWVcIn0iLCJzdWIiOiJqZmFjQDAxZnp2MjNrYTF2Z3MxMGVqM2oydDMxNzc0XC91c2Vyc1wvYWRtaW4iLCJzY3AiOiJhcHBsaWVkLXBlcm1pc3Npb25zXC91c2VyIiwiYXVkIjoiKkAqIiwiaXNzIjoiamZmZUAwMDAiLCJleHAiOjE2NDkyNTQyMTYsImlhdCI6MTY0OTI1NDE1NiwianRpIjoiZDYwZTgxZmUtODYyYi00MWQwLThmMmEtMDZmMTc3MjJjNjdkIn0.D7L8DsLCWfUImNSi9g-S_Va5mBQhxu5jzQajejSRzqmPkrODRo1QlTRXMNQrlEk9vDDDWFlFOqesRZYTCx0bTYs041JDLdWJu29iu33CbAmTbiYb14RVLfplJj1jDhary_VswVzfWwPx0O83RK7oT8ni3JXlfRPpgFMQ54wzBp9k5LgNSJfk2hFexNIKDHxVjpsZf6f8v2mfw5DBm58heWMxEu5uooxS9vWDR9R1Kmw90pKRmvDgGUo5O3rpZXZh4usEK0cwfF9Nny3TZ0Ij-U57BMl_uZQ_llOHMoLZ6k6oXBHZdsZP4A3gmeFXivZILSySvIS0hWz7_J7uFtw0lQ"
	token, err := testsAccessTokensService.CreateAccessToken(tokenParams)
	assert.NoError(t, err)
	assert.NotEqual(t, "", token.AccessToken, "Access token is empty")
	// TODO: check why 31536000
	assert.Equal(t, 31536000, token.ExpiresIn)
}
