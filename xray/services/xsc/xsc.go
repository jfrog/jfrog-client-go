package xsc

import (
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/xsc/services"
)

// XscService is the Xray Source Control service in the Xray service, available from v3.107.13.
// This service replaces the Xray Source Control service, which was available as a standalone service.
type XscInnerService struct {
	client      *jfroghttpclient.JfrogHttpClient
	XrayDetails auth.ServiceDetails
}

func NewXscService(client *jfroghttpclient.JfrogHttpClient) *XscInnerService {
	return &XscInnerService{client: client}
}

func (xs *XscInnerService) GetVersion() (string, error) {
	versionService := services.NewVersionService(xs.client)
	versionService.XrayDetails = xs.XrayDetails
	return versionService.GetVersion()
}

func (xs *XscInnerService) AddAnalyticsGeneralEvent(event services.XscAnalyticsGeneralEvent) (string, error) {
	eventService := services.NewAnalyticsEventService(xs.client)
	eventService.XrayDetails = xs.XrayDetails
	return eventService.AddGeneralEvent(event)
}

func (xs *XscInnerService) SendXscLogErrorRequest(errorLog *services.ExternalErrorLog) error {
	logErrorService := services.NewLogErrorEventService(xs.client)
	logErrorService.XrayDetails = xs.XrayDetails
	return logErrorService.SendLogErrorEvent(errorLog)
}

func (xs *XscInnerService) UpdateAnalyticsGeneralEvent(event services.XscAnalyticsGeneralEventFinalize) error {
	eventService := services.NewAnalyticsEventService(xs.client)
	eventService.XrayDetails = xs.XrayDetails
	return eventService.UpdateGeneralEvent(event)
}

func (xs *XscInnerService) GetAnalyticsGeneralEvent(msi string) (*services.XscAnalyticsGeneralEvent, error) {
	eventService := services.NewAnalyticsEventService(xs.client)
	eventService.XrayDetails = xs.XrayDetails
	return eventService.GetGeneralEvent(msi)
}

func (xs *XscInnerService) GetConfigProfileByName(profileName string) (*services.ConfigProfile, error) {
	configProfileService := services.NewConfigurationProfileService(xs.client)
	configProfileService.XrayDetails = xs.XrayDetails
	return configProfileService.GetConfigurationProfileByName(profileName)
}

func (xs *XscInnerService) GetConfigProfileByUrl(repoUrl string) (*services.ConfigProfile, error) {
	configProfileService := services.NewConfigurationProfileService(xs.client)
	configProfileService.XrayDetails = xs.XrayDetails
	return configProfileService.GetConfigurationProfileByUrl(repoUrl)
}
