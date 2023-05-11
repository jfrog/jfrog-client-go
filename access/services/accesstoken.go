package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
)

// #nosec G101 -- False positive - no hardcoded credentials.
const tokensApi = "api/v1/tokens"

type TokenService struct {
	client         *jfroghttpclient.JfrogHttpClient
	ServiceDetails auth.ServiceDetails
}

type CreateTokenParams struct {
	auth.CommonTokenParams
	Description           string `json:"description,omitempty"`
	IncludeReferenceToken *bool  `json:"include_reference_token,omitempty"`
	Username              string `json:"username,omitempty"`
}

func NewCreateTokenParams(params CreateTokenParams) CreateTokenParams {
	return CreateTokenParams{
		CommonTokenParams:     params.CommonTokenParams,
		Description:           params.Description,
		IncludeReferenceToken: params.IncludeReferenceToken,
		Username:              params.Username,
	}
}

func NewTokenService(client *jfroghttpclient.JfrogHttpClient) *TokenService {
	return &TokenService{client: client}
}

// Create an access token for the JFrog Platform
func (ps *TokenService) CreateAccessToken(params CreateTokenParams) (auth.CreateTokenResponseData, error) {
	return ps.createAccessToken(params)
}

// Refresh an existing access token without having to provide the old token.
// The Refresh Token is the same API endpoint as Create Token, with a specific grant type: refresh_token
func (ps *TokenService) RefreshAccessToken(token CreateTokenParams) (auth.CreateTokenResponseData, error) {
	// Validate provided parameters
	if token.RefreshToken == "" {
		return auth.CreateTokenResponseData{}, errorutils.CheckErrorf("error: trying to refresh token, but 'refresh_token' field wasn't provided. ")
	}
	// Set refresh required parameters
	var trueValue = true
	params := NewCreateTokenParams(token)
	params.GrantType = "refresh_token"
	params.Refreshable = &trueValue

	return ps.createAccessToken(params)
}

// createAccessToken is used to create & refresh access tokens.
func (ps *TokenService) createAccessToken(params CreateTokenParams) (auth.CreateTokenResponseData, error) {
	// Create output response variable
	tokenInfo := auth.CreateTokenResponseData{}

	// Set the request headers
	httpDetails := ps.ServiceDetails.CreateHttpClientDetails()
	utils.SetContentType("application/json", &httpDetails.Headers)
	err := ps.addAccessTokenAuthorizationHeader(params, &httpDetails)
	if err != nil {
		return tokenInfo, err
	}

	// Marshall the request body
	requestContent, err := json.Marshal(params)
	if errorutils.CheckError(err) != nil {
		return tokenInfo, err
	}
	url := fmt.Sprintf("%s%s", ps.ServiceDetails.GetUrl(), tokensApi)
	resp, body, err := ps.client.SendPost(url, requestContent, &httpDetails)
	if err != nil {
		return tokenInfo, err
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return tokenInfo, err
	}

	// Unmarshall the response body and return
	err = json.Unmarshal(body, &tokenInfo)
	return tokenInfo, errorutils.CheckError(err)
}

// Use AccessToken from ServiceDetails (which is the default behaviour)
// If that is not present then we can use the token we are refreshing as the token
func (ps *TokenService) addAccessTokenAuthorizationHeader(params CreateTokenParams, httpDetails *httputils.HttpClientDetails) error {
	access := ps.ServiceDetails.GetAccessToken()
	if access == "" {
		access = params.AccessToken
	}
	if access == "" {
		return errorutils.CheckErrorf("failed: adding accessToken authorization, but No accessToken was provided. ")
	}
	utils.AddHeader("Authorization", fmt.Sprintf("Bearer %s", access), &httpDetails.Headers)
	return nil
}
