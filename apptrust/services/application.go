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

func (as *ApplicationService) GetApplicationVersionPromotions(applicationKey, applicationVersion string, queryParams map[string]string) (*ApplicationPromotionsResponse, error) {
	restApi := path.Join(applicationDetailsAPI, applicationKey, "versions", applicationVersion, "promotions")
	requestFullUrl, err := clientutils.BuildUrl(as.GetApptrustDetails().GetUrl(), restApi, queryParams)
	if err != nil {
		return nil, errorutils.CheckError(err)
	}

	httpClientDetails := as.GetApptrustDetails().CreateHttpClientDetails()
	httpClientDetails.SetContentTypeApplicationJson()

	log.Debug("Getting Application Version Promotions for:", applicationKey, applicationVersion)
	resp, body, _, err := as.client.SendGet(requestFullUrl, true, &httpClientDetails)
	if err != nil {
		return nil, err
	}

	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return nil, err
	}

	var promotionsResponse ApplicationPromotionsResponse
	if err = json.Unmarshal(body, &promotionsResponse); err != nil {
		return nil, errorutils.CheckError(err)
	}

	return &promotionsResponse, nil
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

type ApplicationPromotionsResponse struct {
	Promotions []ApplicationPromotion `json:"promotions"`
	Total      int                    `json:"total"`
	Limit      int                    `json:"limit"`
	Offset     int                    `json:"offset"`
}

type ApplicationPromotion struct {
	ApplicationKey     string                `json:"application_key"`
	ApplicationVersion string                `json:"application_version"`
	ProjectKey         string                `json:"project_key"`
	Status             PromotionStatus       `json:"status"`
	SourceStage        string                `json:"source_stage"`
	TargetStage        string                `json:"target_stage"`
	CreatedBy          string                `json:"created_by"`
	Created            string                `json:"created"`
	CreatedMillis      int64                 `json:"created_millis"`
	Evaluations        *PromotionEvaluations `json:"evaluations,omitempty"`
}

type PromotionEvaluations struct {
	ExitGate  *GateEvaluation `json:"exit_gate,omitempty"`
	EntryGate *GateEvaluation `json:"entry_gate,omitempty"`
}

type GateEvaluation struct {
	EvalId      string `json:"eval_id"`
	Stage       string `json:"stage"`
	Decision    string `json:"decision"`
	Explanation string `json:"explanation,omitempty"`
}

type PromotionStatus string

const (
	PromotionStatusCompleted PromotionStatus = "COMPLETED"
	PromotionStatusFailed    PromotionStatus = "FAILED"
	PromotionStatusStarted   PromotionStatus = "STARTED"
	PromotionStatusDeleting  PromotionStatus = "DELETING"
)
