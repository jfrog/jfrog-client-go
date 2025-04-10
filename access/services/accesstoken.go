package services

import (
	"encoding/json"
	"fmt"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"net/http"
)

// #nosec G101 -- False positive - no hardcoded credentials.
const (
	// jfrog-ignore - not a real token
	tokensApi = "api/v1/tokens"
	// jfrog-ignore - not a real token
	oidcTokensApi = "api/v1/oidc/token"
)

type TokenService struct {
	client         *jfroghttpclient.JfrogHttpClient
	ServiceDetails auth.ServiceDetails
}

type CreateTokenParams struct {
	auth.CommonTokenParams
	IncludeReferenceToken *bool  `json:"include_reference_token,omitempty"`
	Username              string `json:"username,omitempty"`
	ProjectKey            string `json:"project_key,omitempty"`
	Description           string `json:"description,omitempty"`
}

type CreateOidcTokenParams struct {
	GrantType             string `json:"grant_type,omitempty"`
	SubjectTokenType      string `json:"subject_token_type,omitempty"`
	OidcTokenID           string `json:"subject_token,omitempty"`
	ProviderName          string `json:"provider_name,omitempty"`
	ProjectKey            string `json:"project_key,omitempty"`
	JobId                 string `json:"job_id,omitempty"`
	RunId                 string `json:"run_id,omitempty"`
	Repo                  string `json:"repo,omitempty"`
	ApplicationKey        string `json:"application_key,omitempty"`
	Audience              string `json:"audience,omitempty"`
	IdentityMappingName   string `json:"identity_mapping_name,omitempty"`
	IncludeReferenceToken *bool  `json:"include_reference_token,omitempty"`
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

func (ps *TokenService) RefreshAccessToken(params CreateTokenParams) (auth.CreateTokenResponseData, error) {
	refreshParams, err := prepareForRefresh(params)
	if err != nil {
		return auth.CreateTokenResponseData{}, err
	}
	return ps.createAccessToken(*refreshParams)
}

// createAccessToken is used to create & refresh access tokens.
func (ps *TokenService) createAccessToken(params CreateTokenParams) (tokenInfo auth.CreateTokenResponseData, err error) {
	httpDetails := ps.ServiceDetails.CreateHttpClientDetails()
	httpDetails.SetContentTypeApplicationJson()
	if err = ps.handleUnauthenticated(params, &httpDetails); err != nil {
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

func (ps *TokenService) handleUnauthenticated(params CreateTokenParams, httpDetails *httputils.HttpClientDetails) error {
	// Creating access tokens using username and password is available since Artifactory 7.63.2,
	// by enabling "Enable Token Generation via API" in the UI.
	if httpDetails.AccessToken != "" || (httpDetails.User != "" && httpDetails.Password != "") {
		return nil
	}
	// Use token from params if provided.
	if params.AccessToken != "" {
		httpDetails.AccessToken = params.AccessToken
		return nil
	}
	return errorutils.CheckErrorf("cannot create access token without credentials")
}

func (ps *TokenService) ExchangeOidcToken(params CreateOidcTokenParams) (auth.OidcTokenResponseData, error) {
	var tokenInfo auth.OidcTokenResponseData
	httpDetails := ps.ServiceDetails.CreateHttpClientDetails()
	httpDetails.SetContentTypeApplicationJson()
	requestContent, err := json.Marshal(params)
	if errorutils.CheckError(err) != nil {
		return tokenInfo, err
	}
	url := fmt.Sprintf("%s%s", ps.ServiceDetails.GetUrl(), oidcTokensApi)
	resp, body, err := ps.client.SendPost(url, requestContent, &httpDetails)
	if err != nil {
		return tokenInfo, err
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return tokenInfo, fmt.Errorf("failed to exchange OIDC token: %w", err)
	}
	err = json.Unmarshal(body, &tokenInfo)
	return tokenInfo, errorutils.CheckError(err)
}

func prepareForRefresh(p CreateTokenParams) (*CreateTokenParams, error) {
	// Validate provided parameters
	if p.RefreshToken == "" {
		return nil, errorutils.CheckErrorf("trying to refresh token, but 'refresh_token' field wasn't provided")
	}

	params := NewCreateTokenParams(p)
	// Set refresh required parameters
	params.GrantType = "refresh_token"
	params.Refreshable = clientutils.Pointer(true)
	return &params, nil
}
