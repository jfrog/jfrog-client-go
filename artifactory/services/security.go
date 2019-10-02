package services

import (
	"encoding/json"
	"errors"
	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	rthttpclient "github.com/jfrog/jfrog-client-go/artifactory/httpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"net/http"
	"net/url"
	"strconv"
)

const tokenPath = "api/security/token"

type SecurityService struct {
	client     *rthttpclient.ArtifactoryHttpClient
	ArtDetails auth.ArtifactoryDetails
}

func NewSecurityService(client *rthttpclient.ArtifactoryHttpClient) *SecurityService {
	return &SecurityService{client: client}
}

func (ss *SecurityService) getArtifactoryDetails() auth.ArtifactoryDetails {
	return ss.ArtDetails
}

func (ss *SecurityService) CreateToken(params CreateTokenParams) (CreateTokenResponseData, error) {
	artifactoryUrl := ss.ArtDetails.GetUrl()
	data := buildCreateTokenUrlValues(params)
	httpClientsDetails := ss.getArtifactoryDetails().CreateHttpClientDetails()
	resp, body, err := ss.client.SendPostForm(artifactoryUrl+tokenPath, data, &httpClientsDetails)
	tokenInfo := CreateTokenResponseData{}
	if err != nil {
		return tokenInfo, err
	}
	if resp.StatusCode != http.StatusOK {
		return tokenInfo, errorutils.CheckError(
			errors.New("Artifactory response: " + resp.Status + "\n" + clientutils.IndentJson(body)))
	}
	err = json.Unmarshal(body, &tokenInfo)
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
	if resp.StatusCode != http.StatusOK {
		return tokens, errorutils.CheckError(
			errors.New("Artifactory response: " + resp.Status + "\n" + clientutils.IndentJson(body)))
	}
	err = json.Unmarshal(body, &tokens)
	return tokens, err
}

func (ss *SecurityService) RefreshToken(params RefreshTokenParams) (CreateTokenResponseData, error) {
	artifactoryUrl := ss.ArtDetails.GetUrl()
	data := buildRefreshTokenUrlValues(params)
	httpClientsDetails := ss.getArtifactoryDetails().CreateHttpClientDetails()
	resp, body, err := ss.client.SendPostForm(artifactoryUrl+tokenPath, data, &httpClientsDetails)
	tokenInfo := CreateTokenResponseData{}
	if err != nil {
		return tokenInfo, err
	}
	if resp.StatusCode != http.StatusOK {
		return tokenInfo, errorutils.CheckError(
			errors.New("Artifactory response: " + resp.Status + "\n" + clientutils.IndentJson(body)))
	}
	err = json.Unmarshal(body, &tokenInfo)
	return tokenInfo, err
}

func (ss *SecurityService) RevokeToken(params RevokeTokenParams) error {
	artifactoryUrl := ss.ArtDetails.GetUrl()
	requestFullUrl := artifactoryUrl + tokenPath + "/revoke"
	httpClientsDetails := ss.getArtifactoryDetails().CreateHttpClientDetails()
	data := buildRevokeTokenUrlValues(params)
	resp, body, err := ss.client.SendPostForm(requestFullUrl, data, &httpClientsDetails)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errorutils.CheckError(
			errors.New("Artifactory response: " + resp.Status + "\n" + clientutils.IndentJson(body)))
	}
	return err
}

func buildCreateTokenUrlValues(params CreateTokenParams) url.Values {
	data := url.Values{}
	data.Set("refreshable", strconv.FormatBool(params.Refreshable))
	if params.Scope != "" {
		data.Set("scope", params.Scope)
	}
	if params.Username != "" {
		data.Set("username", params.Username)
	}
	if params.Audience != "" {
		data.Set("audience", params.Audience)
	}
	if params.ExpiresIn != 0 {
		data.Set("expires_in", strconv.Itoa(params.ExpiresIn))
	}
	return data
}

func buildRefreshTokenUrlValues(params RefreshTokenParams) url.Values {
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

type CreateTokenResponseData struct {
	Scope        string `json:"scope,omitempty"`
	AccessToken  string `json:"access_token,omitempty"`
	ExpiresIn    int    `json:"expires_in,omitempty"`
	TokenType    string `json:"token_type,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type GetTokensResponseData struct {
	Tokens []struct {
		Issuer      string `json:"issuer,omitempty"`
		Subject     string `json:"subject,omitempty"`
		Refreshable bool   `json:"refreshable,omitempty"`
		Expiry      int    `json:"expiry,omitempty"`
		TokenId     string `json:"token_id,omitempty"`
		IssuedAt    int    `json:"issued_at,omitempty"`
	}
}

type CreateTokenParams struct {
	Scope       string
	Username    string
	ExpiresIn   int
	Refreshable bool
	Audience    string
}

type RefreshTokenParams struct {
	Token        CreateTokenParams
	RefreshToken string
	AccessToken  string
}

type RevokeTokenParams struct {
	Token   string
	TokenId string
}

func NewCreateTokenParams() CreateTokenParams {
	return CreateTokenParams{}
}

func NewRefreshTokenParams() RefreshTokenParams {
	return RefreshTokenParams{}
}

func NewRevokeTokenParams() RevokeTokenParams {
	return RevokeTokenParams{}
}
