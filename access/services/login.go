package services

import (
	"encoding/json"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"net/http"
	"path"
	"time"
)

const (
	baseClientLoginApi       = "api/v2/authentication/jfrog_client_login/"
	requestApi               = "request"
	tokenApi                 = "token"
	defaultMaxWait           = 5 * time.Minute
	defaultWaitSleepInterval = 3 * time.Second
)

var MaxWait = defaultMaxWait

type LoginService struct {
	client         *jfroghttpclient.JfrogHttpClient
	ServiceDetails auth.ServiceDetails
}

type LoginAuthRequestBody struct {
	Session string `json:"session,omitempty"`
}

func NewLoginService(client *jfroghttpclient.JfrogHttpClient) *LoginService {
	return &LoginService{client: client}
}

func (ls *LoginService) SendLoginAuthenticationRequest(uuid string) error {
	restAPI := path.Join(baseClientLoginApi, requestApi)
	fullUrl, err := utils.BuildArtifactoryUrl(ls.ServiceDetails.GetUrl(), restAPI, make(map[string]string))
	if err != nil {
		return err
	}
	data := LoginAuthRequestBody{
		Session: uuid,
	}
	requestContent, err := json.Marshal(data)
	if err != nil {
		return errorutils.CheckError(err)
	}
	httpClientsDetails := ls.ServiceDetails.CreateHttpClientDetails()
	utils.SetContentType("application/json", &httpClientsDetails.Headers)
	resp, body, err := ls.client.SendPost(fullUrl, requestContent, &httpClientsDetails)
	if err != nil {
		return err
	}
	return errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK)
}

func (ls *LoginService) GetLoginAuthenticationToken(uuid string) (token auth.CommonTokenParams, err error) {
	pollingAction := func() (shouldStop bool, responseBody []byte, err error) {
		var resp *http.Response
		resp, responseBody, err = ls.getLoginAuthenticationToken(uuid)
		if err != nil {
			return true, nil, err
		}
		switch resp.StatusCode {
		case http.StatusOK:
			return true, responseBody, nil
		case http.StatusBadRequest:
			return false, responseBody, nil
		default:
			return true, nil, errorutils.CheckErrorf("received unexpected response when getting the login authentication token")
		}
	}
	pollingExecutor := &httputils.PollingExecutor{
		Timeout:         MaxWait,
		PollingInterval: defaultWaitSleepInterval,
		PollingAction:   pollingAction,
		MsgPrefix:       "Attempting to get authentication token from the JFrog platform...",
	}
	finalRespBody, err := pollingExecutor.Execute()
	if err != nil {
		return auth.CommonTokenParams{}, err
	}
	err = errorutils.CheckError(json.Unmarshal(finalRespBody, &token))
	return
}

func (ls *LoginService) getLoginAuthenticationToken(uuid string) (resp *http.Response, body []byte, err error) {
	restAPI := path.Join(baseClientLoginApi, tokenApi, uuid)
	fullUrl, err := utils.BuildArtifactoryUrl(ls.ServiceDetails.GetUrl(), restAPI, make(map[string]string))
	if err != nil {
		return
	}

	httpClientsDetails := ls.ServiceDetails.CreateHttpClientDetails()
	resp, body, _, err = ls.client.SendGet(fullUrl, true, &httpClientsDetails)
	return
}
