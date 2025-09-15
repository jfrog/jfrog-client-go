package services

import (
	"encoding/json"
	"net/http"
	"path"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

const (
	applicationDetailsAPI = "apptrust/api/v1/applications"
)

type ApplicationService struct {
	client          *jfroghttpclient.JfrogHttpClient
	apptrustDetails auth.ServiceDetails
}

func NewApplicationService(apptrustDetails auth.ServiceDetails, client *jfroghttpclient.JfrogHttpClient) *ApplicationService {
	return &ApplicationService{apptrustDetails: apptrustDetails, client: client}
}

func (as *ApplicationService) GetApptrustDetails() auth.ServiceDetails {
	return as.apptrustDetails
}

func (as *ApplicationService) GetApplicationDetails(applicationKey string) (*Application, error) {
	restApi := path.Join(applicationDetailsAPI, applicationKey)
	requestFullUrl, err := clientutils.BuildUrl(as.GetApptrustDetails().GetUrl(), restApi, nil)
	if err != nil {
		return nil, errorutils.CheckError(err)
	}

	httpClientDetails := as.GetApptrustDetails().CreateHttpClientDetails()
	httpClientDetails.SetContentTypeApplicationJson()

	log.Debug("Getting Application Details for:", applicationKey)
	resp, body, _, err := as.client.SendGet(requestFullUrl, true, &httpClientDetails)
	if err != nil {
		return nil, err
	}

	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return nil, err
	}

	var applicationResponse ApplicationResponse
	if err = json.Unmarshal(body, &applicationResponse); err != nil {
		return nil, errorutils.CheckError(err)
	}

	return &applicationResponse.Application, nil
}

type ApplicationResponse struct {
	Application Application `json:"application"`
}

type Application struct {
	ApplicationName string `json:"application_name"`
	ApplicationKey  string `json:"application_key"`
	ProjectName     string `json:"project_name"`
	ProjectKey      string `json:"project_key"`
}
