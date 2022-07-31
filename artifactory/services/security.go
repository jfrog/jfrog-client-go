package services

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
)

const tokenPath = "api/security/token"
const APIKeyPath = "api/security/apiKey"

type SecurityService struct {
	client     *jfroghttpclient.JfrogHttpClient
	ArtDetails auth.ServiceDetails
}

func NewSecurityService(client *jfroghttpclient.JfrogHttpClient) *SecurityService {
	return &SecurityService{client: client}
}

func (ss *SecurityService) getArtifactoryDetails() auth.ServiceDetails {
	return ss.ArtDetails
}

// Create an API key for the current user. Returns an error if API key already exists - use regenerate API key instead.
func (ss *SecurityService) CreateAPIKey() (string, error) {
	httpClientDetails := ss.ArtDetails.CreateHttpClientDetails()
	reqURL, err := utils.BuildArtifactoryUrl(ss.ArtDetails.GetUrl(), APIKeyPath, nil)
	if err != nil {
		return "", err
	}

	resp, body, err := ss.client.SendPost(reqURL, nil, &httpClientDetails)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusCreated {
		return "", errors.New("API key creation failed with status: " + resp.Status + "\n" + clientutils.IndentJson(body))
	}

	return getApiKeyFromBody(body)
}

// RegenerateAPIKey regenerates the API Key in Artifactory
func (ss *SecurityService) RegenerateAPIKey() (string, error) {
	httpClientDetails := ss.ArtDetails.CreateHttpClientDetails()

	reqURL, err := utils.BuildArtifactoryUrl(ss.ArtDetails.GetUrl(), APIKeyPath, nil)
	if err != nil {
		return "", err
	}

	resp, body, err := ss.client.SendPut(reqURL, nil, &httpClientDetails)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("API key regeneration failed with status: " + resp.Status + "\n" + clientutils.IndentJson(body))
	}

	return getApiKeyFromBody(body)
}

// Returns empty string if API Key wasn't generated.
func (ss *SecurityService) GetAPIKey() (string, error) {
	httpClientDetails := ss.ArtDetails.CreateHttpClientDetails()
	reqURL, err := utils.BuildArtifactoryUrl(ss.ArtDetails.GetUrl(), APIKeyPath, nil)
	if err != nil {
		return "", err
	}

	resp, body, _, err := ss.client.SendGet(reqURL, true, &httpClientDetails)
	if err != nil {
		return "", err
	}

	if err = errorutils.CheckResponseStatus(resp, body, http.StatusOK); err != nil {
		return "", err
	}

	return getApiKeyFromBody(body)
}

func getApiKeyFromBody(body []byte) (string, error) {
	var data = make(map[string]interface{})
	if err := json.Unmarshal(body, &data); err != nil {
		return "", errorutils.CheckErrorf("unable to decode json. Error: %w Upstream response: %s", err, string(body))
	}

	if len(data) == 0 {
		return "", nil
	}
	apiKey := data["apiKey"].(string)
	return apiKey, nil
}

func (ss *SecurityService) CreateToken(params CreateTokenParams) (auth.CreateTokenResponseData, error) {
	artifactoryUrl := ss.ArtDetails.GetUrl()
	data := buildCreateTokenUrlValues(params)
	httpClientsDetails := ss.getArtifactoryDetails().CreateHttpClientDetails()
	resp, body, err := ss.client.SendPostForm(artifactoryUrl+tokenPath, data, &httpClientsDetails)
	tokenInfo := auth.CreateTokenResponseData{}
	if err != nil {
		return tokenInfo, err
	}
	if err = errorutils.CheckResponseStatus(resp, body, http.StatusOK); err != nil {
		return tokenInfo, err
	}
	if err = json.Unmarshal(body, &tokenInfo); err != nil {
		return tokenInfo, errorutils.CheckError(err)
	}
	return tokenInfo, err
}

func (ss *SecurityService) GetTokens() (GetTokensResponseData, error) {
	artifactoryUrl := ss.ArtDetails.GetUrl()
	httpClientsDetails := ss.getArtifactoryDetails().CreateHttpClientDetails()
	resp, body, _, err := ss.client.SendGet(artifactoryUrl+tokenPath, true, &httpClientsDetails)
	tokens := GetTokensResponseData{}
	if err != nil {
		return tokens, err
	}
	if err = errorutils.CheckResponseStatus(resp, body, http.StatusOK); err != nil {
		return tokens, err
	}
	if err = json.Unmarshal(body, &tokens); err != nil {
		return tokens, errorutils.CheckError(err)
	}
	return tokens, err
}

func (ss *SecurityService) GetUserTokens(username string) ([]string, error) {
	var tokens []string
	tokenResponseData, err := ss.GetTokens()
	if err != nil {
		return nil, err
	}
	for _, token := range tokenResponseData.Tokens {
		if strings.HasSuffix(token.Subject, "/"+username) {
			tokens = append(tokens, token.TokenId)
		}
	}
	return tokens, nil
}

func (ss *SecurityService) RefreshToken(params ArtifactoryRefreshTokenParams) (auth.CreateTokenResponseData, error) {
	artifactoryUrl := ss.ArtDetails.GetUrl()
	data := buildRefreshTokenUrlValues(params)
	httpClientsDetails := ss.getArtifactoryDetails().CreateHttpClientDetails()
	resp, body, err := ss.client.SendPostForm(artifactoryUrl+tokenPath, data, &httpClientsDetails)
	tokenInfo := auth.CreateTokenResponseData{}
	if err != nil {
		return tokenInfo, err
	}
	if err = errorutils.CheckResponseStatus(resp, body, http.StatusOK); err != nil {
		return tokenInfo, err
	}
	if err = json.Unmarshal(body, &tokenInfo); err != nil {
		return tokenInfo, errorutils.CheckError(err)
	}
	return tokenInfo, err
}

func (ss *SecurityService) RevokeToken(params RevokeTokenParams) (string, error) {
	artifactoryUrl := ss.ArtDetails.GetUrl()
	requestFullUrl := artifactoryUrl + tokenPath + "/revoke"
	httpClientsDetails := ss.getArtifactoryDetails().CreateHttpClientDetails()
	data := buildRevokeTokenUrlValues(params)
	resp, body, err := ss.client.SendPostForm(requestFullUrl, data, &httpClientsDetails)
	if err != nil {
		return "", err
	}
	if err = errorutils.CheckResponseStatus(resp, body, http.StatusOK); err != nil {
		return "", err
	}
	return string(body), err
}

func buildCreateTokenUrlValues(params CreateTokenParams) url.Values {
	// Gathers required data while avoiding default/ignored values
	data := url.Values{}
	if params.Refreshable {
		data.Set("refreshable", "true")
	}
	if params.Scope != "" {
		data.Set("scope", params.Scope)
	}
	if params.Username != "" {
		data.Set("username", params.Username)
	}
	if params.Audience != "" {
		data.Set("audience", params.Audience)
	}
	if params.ExpiresIn >= 0 {
		data.Set("expires_in", strconv.Itoa(params.ExpiresIn))
	}
	return data
}

func buildRefreshTokenUrlValues(params ArtifactoryRefreshTokenParams) url.Values {
	data := buildCreateTokenUrlValues(params.Token)

	// <grant_type> is used to tell the rest api whether to create or refresh a token.
	// Both operations are performed by the same endpoint.
	data.Set("grant_type", "refresh_token")

	if params.RefreshToken != "" {
		data.Set("refresh_token", params.RefreshToken)
	}
	if params.AccessToken != "" {
		data.Set("access_token", params.AccessToken)
	}
	return data
}

func buildRevokeTokenUrlValues(params RevokeTokenParams) url.Values {
	data := url.Values{}
	if params.Token != "" {
		data.Set("token", params.Token)
	}
	if params.TokenId != "" {
		data.Set("token_id", params.TokenId)
	}
	return data
}

type GetTokensResponseData struct {
	Tokens []Token
}

type Token struct {
	Issuer      string `json:"issuer,omitempty"`
	Subject     string `json:"subject,omitempty"`
	Refreshable bool   `json:"refreshable,omitempty"`
	Expiry      int    `json:"expiry,omitempty"`
	TokenId     string `json:"token_id,omitempty"`
	IssuedAt    int    `json:"issued_at,omitempty"`
}

type CreateTokenParams struct {
	Scope       string
	Username    string
	ExpiresIn   int
	Refreshable bool
	Audience    string
}

type ArtifactoryRefreshTokenParams struct {
	Token        CreateTokenParams
	RefreshToken string
	AccessToken  string
}

type RevokeTokenParams struct {
	Token   string
	TokenId string
}

func NewCreateTokenParams() CreateTokenParams {
	return CreateTokenParams{ExpiresIn: -1}
}

func NewArtifactoryRefreshTokenParams() ArtifactoryRefreshTokenParams {
	return ArtifactoryRefreshTokenParams{Token: NewCreateTokenParams()}
}

func NewRevokeTokenParams() RevokeTokenParams {
	return RevokeTokenParams{}
}
