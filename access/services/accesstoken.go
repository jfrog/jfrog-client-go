package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
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
	Audience              string `json:"audience,omitempty"`
	ProviderType          string `json:"provider_type,omitempty"`
	IdentityMappingName   string `json:"identity_mapping_name,omitempty"`
	IncludeReferenceToken *bool  `json:"include_reference_token,omitempty"`
	Repo                  string `json:"repo,omitempty"`
	Revision              string `json:"revision,omitempty"`
	Branch                string `json:"branch,omitempty"`
	ApplicationKey        string `json:"application_key,omitempty"`
}

type GetTokensParams struct {
	Description     string `url:"description,omitempty"`
	Username        string `url:"username,omitempty"`
	Refreshable     *bool  `url:"refreshable,omitempty"`
	TokenId         string `url:"token_id,omitempty"`
	OrderBy         string `url:"order_by,omitempty"`
	DescendingOrder *bool  `url:"descending_order,omitempty"`
}

type TokenInfos struct {
	Tokens []TokenInfo `json:"tokens"`
}

type TokenInfo struct {
	TokenId     string `json:"token_id"`
	Subject     string `json:"subject"`
	Expiry      int64  `json:"expiry,omitempty"`
	IssuedAt    int64  `json:"issued_at"`
	Issuer      string `json:"issuer"`
	Description string `json:"description,omitempty"`
	Refreshable bool   `json:"refreshable,omitempty"`
	Scope       string `json:"scope,omitempty"`
	LastUsed    int64  `json:"last_used,omitempty"`
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

func (ps *TokenService) GetTokens(params GetTokensParams) ([]TokenInfo, error) {
	httpDetails := ps.ServiceDetails.CreateHttpClientDetails()
	requestUrl := fmt.Sprintf("%s%s", ps.ServiceDetails.GetUrl(), tokensApi)

	// Build query parameters manually
	queryParams := url.Values{}
	if params.Description != "" {
		queryParams.Add("description", params.Description)
	}
	if params.Username != "" {
		queryParams.Add("username", params.Username)
	}
	if params.Refreshable != nil {
		queryParams.Add("refreshable", strconv.FormatBool(*params.Refreshable))
	}
	if params.TokenId != "" {
		queryParams.Add("token_id", params.TokenId)
	}
	if params.OrderBy != "" {
		queryParams.Add("order_by", params.OrderBy)
	}
	if params.DescendingOrder != nil {
		queryParams.Add("descending_order", strconv.FormatBool(*params.DescendingOrder))
	}

	if queryString := queryParams.Encode(); queryString != "" {
		requestUrl += "?" + queryString
	}

	resp, body, _, err := ps.client.SendGet(requestUrl, true, &httpDetails)
	if err != nil {
		return nil, err
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return nil, err
	}

	var tokenInfos TokenInfos
	err = json.Unmarshal(body, &tokenInfos)
	return tokenInfos.Tokens, errorutils.CheckError(err)
}

func (ps *TokenService) GetTokenByID(tokenId string) (*TokenInfo, error) {
	if tokenId == "" {
		return nil, errorutils.CheckErrorf("token ID cannot be empty")
	}

	httpDetails := ps.ServiceDetails.CreateHttpClientDetails()
	url := fmt.Sprintf("%s%s/%s", ps.ServiceDetails.GetUrl(), tokensApi, tokenId)

	resp, body, _, err := ps.client.SendGet(url, true, &httpDetails)
	if err != nil {
		return nil, err
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return nil, err
	}

	var tokenInfo TokenInfo
	err = json.Unmarshal(body, &tokenInfo)
	return &tokenInfo, errorutils.CheckError(err)
}

func (ps *TokenService) RevokeTokenByID(tokenId string) error {
	if tokenId == "" {
		return errorutils.CheckErrorf("token ID cannot be empty")
	}

	httpDetails := ps.ServiceDetails.CreateHttpClientDetails()
	url := fmt.Sprintf("%s%s/%s", ps.ServiceDetails.GetUrl(), tokensApi, tokenId)

	resp, body, err := ps.client.SendDelete(url, nil, &httpDetails)
	if err != nil {
		return err
	}
	return errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK, http.StatusNoContent)
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
