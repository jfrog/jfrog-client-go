package services

import (
	"encoding/json"
	"fmt"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"net/http"
)

const tokensApi = "api/v1/tokens"

type TokenService struct {
	client         *jfroghttpclient.JfrogHttpClient
	ServiceDetails auth.ServiceDetails
}

type TokenParams struct {
	auth.CreateTokenData
}

func NewTokenService(client *jfroghttpclient.JfrogHttpClient) *TokenService {
	return &TokenService{client: client}
}

func (ps *TokenService) CreateAccessToken(params TokenParams) (auth.CreateTokenResponseData, error) {
	// Set request's headers
	httpDetails := ps.ServiceDetails.CreateHttpClientDetails()
	utils.SetContentType("application/json", &httpDetails.Headers)
	utils.AddHeader("Authorization", fmt.Sprintf("Bearer %s", ps.ServiceDetails.GetAccessToken()), &httpDetails.Headers)

	tokenInfo := auth.CreateTokenResponseData{}
	requestContent, err := json.Marshal(params)
	if errorutils.CheckError(err) != nil {
		return tokenInfo, err
	}
	url := fmt.Sprintf("%s%s", ps.ServiceDetails.GetUrl(), tokensApi)
	resp, body, err := ps.client.SendPost(url, requestContent, &httpDetails)
	if err != nil {
		return tokenInfo, err
	}
	if err = errorutils.CheckResponseStatus(resp, http.StatusOK); err != nil {
		return tokenInfo, errorutils.CheckError(errorutils.GenerateResponseError(resp.Status, clientutils.IndentJson(body)))
	}
	err = json.Unmarshal(body, &tokenInfo)
	return tokenInfo, errorutils.CheckError(err)
}
