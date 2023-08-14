package usage

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	versionutil "github.com/jfrog/gofrog/version"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/http/httpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/jfrog/jfrog-client-go/xray"
)

const (
	minXrayVersion = "3.80.0"

	xrayUsageApiPath     = "api/v1/usage/events/send"
	ecosystemUsageApiUrl = "http://usage-ecosystem.jfrog.info/api/usage/report"
	ReportUsagePrefix    = "Usage Report: "
)

type ReportUsageAttribute struct {
	AttributeName  string
	AttributeValue string
}

func (rua *ReportUsageAttribute) isEmpty() bool {
	return rua.AttributeName == ""
}

type ReportXrayEventData struct {
	Attributes map[string]string `json:"data,omitempty"`
	ProductId  string            `json:"product_name"`
	EventId    string            `json:"event_name"`
	Origin     string            `json:"origin,omitempty"`
}

func SendXrayUsageEvents(serviceManager xray.XrayServicesManager, events ...ReportXrayEventData) error {
	if len(events) == 0 {
		return errorutils.CheckErrorf(ReportUsagePrefix + "Nothing to send.")
	}
	config := serviceManager.Config()
	if config == nil {
		return errorutils.CheckErrorf(ReportUsagePrefix + "Expected full config, but no configuration exists.")
	}
	xrDetails := config.GetServiceDetails()
	if xrDetails == nil {
		return errorutils.CheckErrorf(ReportUsagePrefix + "Xray details not configured.")
	}
	xrayVersion, err := xrDetails.GetVersion()
	if err != nil {
		return errors.New(ReportUsagePrefix + "Couldn't get Xray version. Error: " + err.Error())
	}
	if !isVersionCompatible(xrayVersion) {
		log.Debug(fmt.Sprintf(ReportUsagePrefix+"Expected Xray version %s or above, got %s", minXrayVersion, xrayVersion))
		return nil
	}

	url, err := utils.BuildArtifactoryUrl(xrDetails.GetUrl(), xrayUsageApiPath, make(map[string]string))
	if err != nil {
		return errors.New(ReportUsagePrefix + err.Error())
	}
	clientDetails := xrDetails.CreateHttpClientDetails()

	bodyContent, err := xrayUsageEventsToJson(events...)
	if err != nil {
		return errors.New(ReportUsagePrefix + err.Error())
	}
	utils.AddHeader("Content-Type", "application/json", &clientDetails.Headers)
	resp, body, err := serviceManager.Client().SendPost(url, bodyContent, &clientDetails)
	if err != nil {
		return errors.New(ReportUsagePrefix + "Couldn't send usage info. Error: " + err.Error())
	}

	err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK, http.StatusAccepted)
	if err != nil {
		return err
	}

	log.Debug(ReportUsagePrefix+"Usage info sent successfully.", "Xray response:", resp.Status)
	return nil
}

// Returns an error if the Xray version is not compatible to run usage api
func isVersionCompatible(xrayVersion string) bool {
	// API exists from Xray version 3.80.0 and above:
	version := versionutil.NewVersion(xrayVersion)
	return version.AtLeast(minXrayVersion)
}

func CreateUsageEvent(productId, featureId string, additionalAttributes ...ReportUsageAttribute) ReportXrayEventData {
	reportInfo := ReportXrayEventData{ProductId: productId, EventId: getExpectedEventName(productId, featureId), Origin: "API"}

	if len(additionalAttributes) > 0 {
		reportInfo.Attributes = make(map[string]string, len(additionalAttributes))
		for _, attribute := range additionalAttributes {
			if !attribute.isEmpty() {
				reportInfo.Attributes[attribute.AttributeName] = attribute.AttributeValue
			}
		}
	}
	return reportInfo
}

func xrayUsageEventsToJson(events ...ReportXrayEventData) ([]byte, error) {
	bodyContent, err := json.Marshal(events)
	return bodyContent, errorutils.CheckError(err)
}

func getExpectedEventName(productId, commandName string) string {
	return "server_" + productId + "_" + commandName
}

type ReportEcosystemUsageData struct {
	ProductId string   `json:"productId"`
	AccountId string   `json:"accountId"`
	ClientId  string   `json:"clientId,omitempty"`
	Features  []string `json:"features"`
}

func SendEcosystemUsageReports(reports ...ReportEcosystemUsageData) error {
	if len(reports) == 0 {
		return errorutils.CheckErrorf(ReportUsagePrefix + "Nothing to send.")
	}
	bodyContent, err := ecosystemUsageReportsToJson(reports...)
	if err != nil {
		return errors.New(ReportUsagePrefix + err.Error())
	}

	resp, body, err := sendRequestToEcosystemService(bodyContent)
	if err != nil {
		return errors.New(ReportUsagePrefix + "Couldn't send usage info to Ecosystem. Error: " + err.Error())
	}

	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK, http.StatusAccepted); err != nil {
		return err
	}

	log.Debug(ReportUsagePrefix+"Usage info sent successfully.", "Ecosystem response:", resp.Status)
	return nil
}

func CreateUsageData(productId, accountId, clientId string, features ...string) (reportInfo ReportEcosystemUsageData, err error) {
	reportInfo = ReportEcosystemUsageData{ProductId: productId, AccountId: accountId, ClientId: clientId, Features: []string{}}
	for _, feature := range features {
		if feature != "" {
			reportInfo.Features = append(reportInfo.Features, feature)
		}
	}
	if len(reportInfo.Features) == 0 {
		err = fmt.Errorf("expected at least one feature to report usage on")
	}
	return
}

func ecosystemUsageReportsToJson(events ...ReportEcosystemUsageData) ([]byte, error) {
	bodyContent, err := json.Marshal(events)
	return bodyContent, errorutils.CheckError(err)
}

func sendRequestToEcosystemService(content []byte) (resp *http.Response, respBody []byte, err error) {
	var client *httpclient.HttpClient
	if client, err = httpclient.ClientBuilder().Build(); err != nil {
		return
	}

	details := httputils.HttpClientDetails{}
	utils.AddHeader("Content-Type", "application/json", &details.Headers)
	return client.SendPost(ecosystemUsageApiUrl, content, details, "")
}