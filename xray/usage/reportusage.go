package usage

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/jfrog/jfrog-client-go/xray"
)

const (
	minXrayVersion   = "3.81.4"
	xrayUsageApiPath = "api/v1/usage/events/send"
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
		return errorutils.CheckErrorf("Nothing to send.")
	}
	config := serviceManager.Config()
	if config == nil {
		return errorutils.CheckErrorf("Expected full config, but no configuration exists.")
	}
	xrDetails := config.GetServiceDetails()
	if xrDetails == nil {
		return errorutils.CheckErrorf("Xray details not configured.")
	}
	xrayVersion, err := xrDetails.GetVersion()
	if err != nil {
		return errors.New("Couldn't get Xray version. Error: " + err.Error())
	}
	if e := clientutils.ValidateMinimumVersion(clientutils.Xray, xrayVersion, minXrayVersion); e != nil {
		log.Debug("Usage Report:", e.Error())
		return nil
	}
	url, err := clientutils.BuildUrl(xrDetails.GetUrl(), xrayUsageApiPath, make(map[string]string))
	if err != nil {
		return errors.New(err.Error())
	}
	clientDetails := xrDetails.CreateHttpClientDetails()

	bodyContent, err := json.Marshal(events)
	if errorutils.CheckError(err) != nil {
		return errors.New(err.Error())
	}
	utils.AddHeader("Content-Type", "application/json", &clientDetails.Headers)
	resp, body, err := serviceManager.Client().SendPost(url, bodyContent, &clientDetails)
	if err != nil {
		return errors.New("Couldn't send usage info. Error: " + err.Error())
	}
	err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK, http.StatusAccepted)
	if err != nil {
		return err
	}
	return nil
}

func CreateUsageEvent(productId, featureId string, additionalAttributes ...ReportUsageAttribute) (event ReportXrayEventData) {
	event = ReportXrayEventData{ProductId: productId, EventId: GetExpectedXrayEventName(productId, featureId), Origin: "API_CLI"}
	if len(additionalAttributes) == 0 {
		return
	}
	event.Attributes = make(map[string]string, len(additionalAttributes))
	for _, attribute := range additionalAttributes {
		if !attribute.isEmpty() {
			event.Attributes[attribute.AttributeName] = attribute.AttributeValue
		}
	}
	return
}

func GetExpectedXrayEventName(productId, commandName string) string {
	return "server_" + productId + "_" + commandName
}
