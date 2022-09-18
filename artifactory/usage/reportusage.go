package usage

import (
	"encoding/json"
	"errors"
	"fmt"
	versionutil "github.com/jfrog/gofrog/version"
	"github.com/jfrog/jfrog-client-go/artifactory"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
)

const minArtifactoryVersion = "6.9.0"
const ReportUsagePrefix = "Usage Report: "

type ReportUsageAttribute struct {
	AttributeName  string
	AttributeValue string
}

func (rua *ReportUsageAttribute) isEmpty() bool {
	return rua.AttributeName == ""
}

func SendReportUsage(productId, commandName string, serviceManager artifactory.ArtifactoryServicesManager, attributes ...ReportUsageAttribute) error {
	config := serviceManager.GetConfig()
	if config == nil {
		return errorutils.CheckErrorf(ReportUsagePrefix + "Expected full config, but no configuration exists.")
	}
	rtDetails := config.GetServiceDetails()
	if rtDetails == nil {
		return errorutils.CheckErrorf(ReportUsagePrefix + "Artifactory details not configured.")
	}
	url, err := utils.BuildArtifactoryUrl(rtDetails.GetUrl(), "api/system/usage", make(map[string]string))
	if err != nil {
		return errors.New(ReportUsagePrefix + err.Error())
	}
	clientDetails := rtDetails.CreateHttpClientDetails()
	// Check Artifactory version
	artifactoryVersion, err := rtDetails.GetVersion()
	if err != nil {
		return errors.New(ReportUsagePrefix + err.Error())
	}
	if !isVersionCompatible(artifactoryVersion) {
		log.Debug(fmt.Sprintf(ReportUsagePrefix+"Expected Artifactory version %s or above, got %s", minArtifactoryVersion, artifactoryVersion))
		return nil
	}

	bodyContent, err := reportUsageToJson(productId, commandName, attributes...)
	if err != nil {
		return errors.New(ReportUsagePrefix + err.Error())
	}
	utils.AddHeader("Content-Type", "application/json", &clientDetails.Headers)
	resp, body, err := serviceManager.Client().SendPost(url, bodyContent, &clientDetails)
	if err != nil {
		return errors.New(ReportUsagePrefix + err.Error())
	}

	err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK, http.StatusAccepted)
	if err != nil {
		return errorutils.CheckError(err)
	}

	log.Debug(ReportUsagePrefix+"Artifactory response:", resp.Status)
	log.Debug(ReportUsagePrefix + "Usage info sent successfully.")
	return nil
}

// Returns an error if the Artifactory version is not compatible
func isVersionCompatible(artifactoryVersion string) bool {
	// API exists from Artifactory version 6.9.0 and above:
	version := versionutil.NewVersion(artifactoryVersion)
	return version.AtLeast(minArtifactoryVersion)
}

func reportUsageToJson(productId, commandName string, attributes ...ReportUsageAttribute) ([]byte, error) {
	featureInfo := feature{FeatureId: commandName}
	if len(attributes) > 0 {
		featureInfo.Attributes = make(map[string]string, len(attributes))
		for _, attribute := range attributes {
			if !attribute.isEmpty() {
				featureInfo.Attributes[attribute.AttributeName] = attribute.AttributeValue
			}
		}
	}
	params := reportUsageParams{ProductId: productId, Features: []feature{featureInfo}}
	bodyContent, err := json.Marshal(params)
	return bodyContent, errorutils.CheckError(err)
}

type reportUsageParams struct {
	ProductId string    `json:"productId"`
	Features  []feature `json:"features,omitempty"`
}

type feature struct {
	FeatureId  string            `json:"featureId,omitempty"`
	Attributes map[string]string `json:"attributes,omitempty"`
}
