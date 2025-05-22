package xsc

import (
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/xray/services/utils"
	"github.com/jfrog/jfrog-client-go/xsc/services"
)

// XscService is the Xray Source Control service in the Xray service, available from v3.107.13.
// This service replaces the Xray Source Control service, which was available as a standalone service.
type XscInnerService struct {
	client          *jfroghttpclient.JfrogHttpClient
	XrayDetails     auth.ServiceDetails
	ScopeProjectKey string
}

func NewXscService(client *jfroghttpclient.JfrogHttpClient) *XscInnerService {
	return &XscInnerService{client: client}
}

func (xs *XscInnerService) GetVersion() (string, error) {
	versionService := services.NewVersionService(xs.client)
	versionService.XrayDetails = xs.XrayDetails
	return versionService.GetVersion()
}

func (xs *XscInnerService) AddAnalyticsGeneralEvent(event services.XscAnalyticsGeneralEvent, xrayVersion string) (string, error) {
	eventService := services.NewAnalyticsEventService(xs.client)
	eventService.XrayDetails = xs.XrayDetails
	eventService.ScopeProjectKey = xs.ScopeProjectKey
	return eventService.AddGeneralEvent(event, xrayVersion)
}

func (xs *XscInnerService) SendXscLogErrorRequest(errorLog *services.ExternalErrorLog) error {
	logErrorService := services.NewLogErrorEventService(xs.client)
	logErrorService.XrayDetails = xs.XrayDetails
	logErrorService.ScopeProjectKey = xs.ScopeProjectKey
	return logErrorService.SendLogErrorEvent(errorLog)
}

func (xs *XscInnerService) UpdateAnalyticsGeneralEvent(event services.XscAnalyticsGeneralEventFinalize) error {
	eventService := services.NewAnalyticsEventService(xs.client)
	eventService.XrayDetails = xs.XrayDetails
	eventService.ScopeProjectKey = xs.ScopeProjectKey
	return eventService.UpdateGeneralEvent(event)
}

func (xs *XscInnerService) GetAnalyticsGeneralEvent(msi string) (*services.XscAnalyticsGeneralEvent, error) {
	eventService := services.NewAnalyticsEventService(xs.client)
	eventService.XrayDetails = xs.XrayDetails
	eventService.ScopeProjectKey = xs.ScopeProjectKey
	return eventService.GetGeneralEvent(msi)
}

func (xs *XscInnerService) GetConfigProfileByName(profileName string) (*services.ConfigProfile, error) {
	configProfileService := services.NewConfigurationProfileService(xs.client)
	configProfileService.XrayDetails = xs.XrayDetails
	configProfileService.ScopeProjectKey = xs.ScopeProjectKey
	return configProfileService.GetConfigurationProfileByName(profileName)
}

func (xs *XscInnerService) GetConfigProfileByUrl(repoUrl string) (*services.ConfigProfile, error) {
	configProfileService := services.NewConfigurationProfileService(xs.client)
	configProfileService.XrayDetails = xs.XrayDetails
	configProfileService.ScopeProjectKey = xs.ScopeProjectKey
	return configProfileService.GetConfigurationProfileByUrl(repoUrl)
}

func (xs *XscInnerService) GetResourceWatches(gitRepo, project string) (watches *utils.ResourcesWatchesBody, err error) {
	watchService := services.NewWatchService(xs.client)
	watchService.XrayDetails = xs.XrayDetails
	watchService.ScopeProjectKey = xs.ScopeProjectKey
	return watchService.GetResourceWatches(gitRepo, project)
}
