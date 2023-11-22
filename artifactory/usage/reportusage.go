package usage

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jfrog/jfrog-client-go/artifactory"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
)

const minArtifactoryVersion = "6.9.0"

type ReportUsageAttribute struct {
	AttributeName  string
	AttributeValue string
}

func (rua *ReportUsageAttribute) isEmpty() bool {
	return rua.AttributeName == ""
}

func validateAndGetUsageServerInfo(serviceManager artifactory.ArtifactoryServicesManager) (url string, clientDetails httputils.HttpClientDetails, err error) {
	config := serviceManager.GetConfig()
	if config == nil {
		err = errorutils.CheckErrorf("Expected full config, but no configuration exists.")
		return
	}
	rtDetails := config.GetServiceDetails()
	if rtDetails == nil {
		err = errorutils.CheckErrorf("Artifactory details not configured.")
		return
	}
	// Check Artifactory version
	artifactoryVersion, err := rtDetails.GetVersion()
	if err != nil {
		err = errors.New("Couldn't get Artifactory version. Error: " + err.Error())
		return
	}
	if e := clientutils.ValidateMinimumVersion(clientutils.Artifactory, artifactoryVersion, minArtifactoryVersion); e != nil {
		return
	}
	url, err = clientutils.BuildUrl(rtDetails.GetUrl(), "api/system/usage", make(map[string]string))
	if err != nil {
		return
	}
	clientDetails = rtDetails.CreateHttpClientDetails()
	return
}

func sendReport(url string, serviceManager artifactory.ArtifactoryServicesManager, clientDetails httputils.HttpClientDetails, bodyContent []byte) error {
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

func ReportUsageToArtifactory(productId string, serviceManager artifactory.ArtifactoryServicesManager, features ...Feature) error {
	url, clientDetails, err := validateAndGetUsageServerInfo(serviceManager)
	if err != nil || url == "" {
		return err
	}
	bodyContent, err := usageFeaturesToJson(productId, features...)
	if err != nil {
		return err
	}
	return sendReport(url, serviceManager, clientDetails, bodyContent)
}

func SendReportUsage(productId, commandName string, serviceManager artifactory.ArtifactoryServicesManager, attributes ...ReportUsageAttribute) error {
	url, clientDetails, err := validateAndGetUsageServerInfo(serviceManager)
	if err != nil || url == "" {
		return err
	}
	bodyContent, err := reportUsageToJson(productId, commandName, attributes...)
	if err != nil {
		return err
	}
	return sendReport(url, serviceManager, clientDetails, bodyContent)
}

func usageFeaturesToJson(productId string, features ...Feature) ([]byte, error) {
	params := reportUsageParams{ProductId: productId, Features: features}
	bodyContent, err := json.Marshal(params)
	return bodyContent, errorutils.CheckError(err)
}

func reportUsageToJson(productId, commandName string, attributes ...ReportUsageAttribute) ([]byte, error) {
	featureInfo := Feature{FeatureId: commandName}
	if len(attributes) > 0 {
		featureInfo.Attributes = make(map[string]string, len(attributes))
		for _, attribute := range attributes {
			if !attribute.isEmpty() {
				featureInfo.Attributes[attribute.AttributeName] = attribute.AttributeValue
			}
		}
	}
	return usageFeaturesToJson(productId, featureInfo)
}

type reportUsageParams struct {
	ProductId string    `json:"productId"`
	Features  []Feature `json:"features,omitempty"`
}

type Feature struct {
	FeatureId  string            `json:"featureId,omitempty"`
	Attributes map[string]string `json:"attributes,omitempty"`
	ClientId   string            `json:"uniqueClientId,omitempty"`
}
