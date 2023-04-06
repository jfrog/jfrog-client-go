package services

import (
	"encoding/json"
	"fmt"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"net/http"
)

// #nosec G101 -- False positive - no hardcoded credentials.
const tokensApi = "api/v1/tokens"

type TokenService struct {
	client         *jfroghttpclient.JfrogHttpClient
	ServiceDetails auth.ServiceDetails
}

type CreateTokenParams struct {
	auth.CommonTokenParams
	IncludeReferenceToken *bool `json:"include_reference_token,omitempty"`
}

func NewCreateTokenParams(params CreateTokenParams) CreateTokenParams {
	return CreateTokenParams{CommonTokenParams: params.CommonTokenParams, IncludeReferenceToken: params.IncludeReferenceToken}
}

func NewTokenService(client *jfroghttpclient.JfrogHttpClient) *TokenService {
	return &TokenService{client: client}
}

func (ps *TokenService) CreateAccessToken(params CreateTokenParams) (auth.CreateTokenResponseData, error) {
	return ps.createAccessToken(params)
}

func (ps *TokenService) RefreshAccessToken(token CreateTokenParams) (auth.CreateTokenResponseData, error) {
	param, err := createRefreshTokenRequestParams(token)
	if err != nil {
		return auth.CreateTokenResponseData{}, err
	}
	return ps.createAccessToken(*param)
}

// createAccessToken is used to create & refresh access tokens.
func (ps *TokenService) createAccessToken(params CreateTokenParams) (auth.CreateTokenResponseData, error) {
	// Set the request headers
	tokenInfo := auth.CreateTokenResponseData{}
	httpDetails := ps.ServiceDetails.CreateHttpClientDetails()
	utils.SetContentType("application/json", &httpDetails.Headers)
	err := ps.addAccessTokenAuthorizationHeader(params, &httpDetails)
	if err != nil {
		return tokenInfo, err
	}
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
	err = json.Unmarshal(body, &tokenInfo)
	return tokenInfo, errorutils.CheckError(err)
}

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

func createRefreshTokenRequestParams(p CreateTokenParams) (*CreateTokenParams, error) {
	var trueValue = true
	// Validate provided parameters
	if p.RefreshToken == "" {
		return nil, errorutils.CheckErrorf("error: trying to refresh token, but 'refresh_token' field wasn't provided. ")
	}
	params := NewCreateTokenParams(p)
	// Set refresh required parameters
	params.GrantType = "refresh_token"
	params.Refreshable = &trueValue
	return &params, nil
}
