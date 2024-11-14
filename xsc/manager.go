package xsc

import (
	"strings"

	"github.com/jfrog/jfrog-client-go/config"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/jfrog/jfrog-client-go/xsc/services"
)


// XscServicesManager defines the http client and general configuration
type XscServicesManager struct {
	client *jfroghttpclient.JfrogHttpClient
	config config.Config
}

// New creates a service manager to interact with Xsc
func New(config config.Config) (*XscServicesManager, error) {
	details := config.GetServiceDetails()
	var err error
	manager := &XscServicesManager{config: config}
	manager.client, err = jfroghttpclient.JfrogClientBuilder().
		SetCertificatesPath(config.GetCertificatesPath()).
		SetInsecureTls(config.IsInsecureTls()).
		SetContext(config.GetContext()).
		SetDialTimeout(config.GetDialTimeout()).
		SetOverallRequestTimeout(config.GetOverallRequestTimeout()).
		SetClientCertPath(details.GetClientCertPath()).
		SetClientCertKeyPath(details.GetClientCertKeyPath()).
		AppendPreRequestInterceptor(details.RunPreRequestFunctions).
		SetRetries(config.GetHttpRetries()).
		SetRetryWaitMilliSecs(config.GetHttpRetryWaitMilliSecs()).
		Build()
	return manager, err
}

// Client will return the http client
func (sm *XscServicesManager) Client() *jfroghttpclient.JfrogHttpClient {
	return sm.client
}

func (sm *XscServicesManager) Config() config.Config {
	return sm.config
}

// GetVersion will return the Xsc version
func (sm *XscServicesManager) GetVersion() (string, error) {
	versionService := services.NewVersionService(sm.client)
	versionService.XscDetails = sm.config.GetServiceDetails()
	return versionService.GetVersion()
}

// AddAnalyticsGeneralEvent will send an analytics metrics general event to Xsc and return MSI (multi scan id) generated by Xsc.
func (sm *XscServicesManager) AddAnalyticsGeneralEvent(event services.XscAnalyticsGeneralEvent) (string, error) {
	eventService := services.NewAnalyticsEventService(sm.client)
	eventService.XscDetails = sm.config.GetServiceDetails()
	return eventService.AddGeneralEvent(event)
}

func (sm *XscServicesManager) SendXscLogErrorRequest(errorLog *services.ExternalErrorLog) error {
	logErrorService := services.NewLogErrorEventService(sm.client)
	logErrorService.XscDetails = sm.config.GetServiceDetails()
	return logErrorService.SendLogErrorEvent(errorLog)
}

// UpdateAnalyticsGeneralEvent upon completion of the scan and we have all the results to report on,
// we send a finalized analytics metrics event with information matching an existing event's msi.
func (sm *XscServicesManager) UpdateAnalyticsGeneralEvent(event services.XscAnalyticsGeneralEventFinalize) error {
	eventService := services.NewAnalyticsEventService(sm.client)
	eventService.XscDetails = sm.config.GetServiceDetails()
	return eventService.UpdateGeneralEvent(event)
}

// GetAnalyticsGeneralEvent returns general event that match the msi provided.
func (sm *XscServicesManager) GetAnalyticsGeneralEvent(msi string) (*services.XscAnalyticsGeneralEvent, error) {
	eventService := services.NewAnalyticsEventService(sm.client)
	eventService.XscDetails = sm.config.GetServiceDetails()
	return eventService.GetGeneralEvent(msi)
}

func (sm *XscServicesManager) GetConfigProfile(profileName string) (*services.ConfigProfile, error) {
	configProfileService := services.NewConfigurationProfileService(sm.client)
	configProfileService.XscDetails = sm.config.GetServiceDetails()
	return configProfileService.GetConfigurationProfile(profileName)
}
