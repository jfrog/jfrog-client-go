package usage

import (
	"encoding/json"
	"errors"
	"net/http"

	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/xray"
)

const (
	minXrayReportUsageVersion = "3.83.0"
	xrayReportUsageApiPath    = "api/v1/usage/events/send"
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
		return errorutils.CheckErrorf("nothing to send.")
	}
	config := serviceManager.Config()
	if config == nil {
		return errorutils.CheckErrorf("expected full config, but no configuration exists.")
	}
	xrDetails := config.GetServiceDetails()
	if xrDetails == nil {
		return errorutils.CheckErrorf("Xray details not configured.")
	}
	xrayVersion, err := xrDetails.GetVersion()
	if err != nil {
		return errors.New("Couldn't get Xray version. Error: " + err.Error())
	}
	if clientutils.ValidateMinimumVersion(clientutils.Xray, xrayVersion, minXrayReportUsageVersion) != nil {
		//nolint:nilerr
		return nil
	}
	url, err := clientutils.BuildUrl(xrDetails.GetUrl(), xrayReportUsageApiPath, make(map[string]string))
	if err != nil {
		return err
	}
	clientDetails := xrDetails.CreateHttpClientDetails()

	bodyContent, err := json.Marshal(events)
	if errorutils.CheckError(err) != nil {
		return err
	}
	clientDetails.SetContentTypeApplicationJson()
	resp, body, err := serviceManager.Client().SendPost(url, bodyContent, &clientDetails)
	if err != nil {
		return errors.New("Couldn't send usage info. Error: " + err.Error())
	}
	return errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK, http.StatusAccepted)
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
