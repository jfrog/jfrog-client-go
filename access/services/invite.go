package services

import (
	"encoding/json"
	"fmt"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"net/http"
)

const (
	inviteApi           = "api/v1/users/invite"
	InviteCliSourceName = "cli"
)

type InviteService struct {
	client         *jfroghttpclient.JfrogHttpClient
	ServiceDetails auth.ServiceDetails
}

type InvitedUser struct {
	InvitedEmail string `json:"invited_email,omitempty" csv:"invited_email,omitempty"`
	Source       string `json:"source,omitempty" csv:"source,omitempty"`
}

func NewInviteService(client *jfroghttpclient.JfrogHttpClient) *InviteService {
	return &InviteService{client: client}
}

func (us *InviteService) InviteUser(email, source string) error {
	httpDetails := us.ServiceDetails.CreateHttpClientDetails()
	url := fmt.Sprintf("%s%s", us.ServiceDetails.GetUrl(), inviteApi)
	data := InvitedUser{
		InvitedEmail: email,
		Source:       source,
	}
	requestContent, err := json.Marshal(data)
	if err != nil {
		return errorutils.CheckError(err)
	}
	utils.SetContentType("application/json", &httpDetails.Headers)
	resp, body, err := us.client.SendPost(url, requestContent, &httpDetails)
	if err != nil {
		return err
	}
	if resp == nil {
		return errorutils.CheckErrorf("no response was returned for the request sent")
	}
	return errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK)
}
