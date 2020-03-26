package usage

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jfrog/jfrog-client-go/artifactory"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	versionutil "github.com/jfrog/jfrog-client-go/utils/version"
	"github.com/pkg/errors"
)

const minArtifactoryVersion = "6.9.0"
const usagePrefix = "Usage Report: "

func SendReportUsage(productId, commandName string, serviceManager *artifactory.ArtifactoryServicesManager) error {
	config := serviceManager.GetConfig()
	if config == nil {
		return errorutils.CheckError(errors.New(usagePrefix + "Expected full config, but no configuration exists."))
	}
	rtDetails := config.GetCommonDetails()
	if rtDetails == nil {
		return errorutils.CheckError(errors.New(usagePrefix + "Artifactory details not configured."))
	}
	url, err := utils.BuildArtifactoryUrl(rtDetails.GetUrl(), "api/system/usage", make(map[string]string))
	if err != nil {
		return err
	}
	clientDetails := rtDetails.CreateHttpClientDetails()
	// Check Artifactory version
	artifactoryVersion, err := rtDetails.GetVersion()
	if err != nil {
		return errors.New(usagePrefix + err.Error())
	}
	if !isVersionCompatible(artifactoryVersion) {
		log.Debug(fmt.Sprintf(usagePrefix+"Expected Artifactory version %s or above, got %s", minArtifactoryVersion, artifactoryVersion))
		return nil
	}

	bodyContent, err := reportUsageToJson(productId, commandName)
	if err != nil {
		return errors.New(usagePrefix + err.Error())
	}
	utils.AddHeader("Content-Type", "application/json", &clientDetails.Headers)
	resp, body, err := serviceManager.Client().SendPost(url, bodyContent, &clientDetails)
	if err != nil {
		return errors.New(usagePrefix + err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(usagePrefix + "Artifactory response: " + resp.Status + "\n" + clientutils.IndentJson(body))
	}

	log.Debug(usagePrefix+"Artifactory response:", resp.Status)
	log.Debug(usagePrefix + "Usage info sent successfully.")
	return nil
}

// Returns an error if the Artifactory version is not compatible
func isVersionCompatible(artifactoryVersion string) bool {
	// API exists from Artifactory version 6.9.0 and above:
	version := versionutil.NewVersion(artifactoryVersion)
	return version.AtLeast(minArtifactoryVersion)
}

func reportUsageToJson(productId, commandName string) ([]byte, error) {
	featureInfo := feature{FeatureId: commandName}
	params := reportUsageParams{ProductId: productId, Features: []feature{featureInfo}}
	bodyContent, err := json.Marshal(params)
	return bodyContent, errorutils.CheckError(err)
}

type reportUsageParams struct {
	ProductId string    `json:"productId"`
	Features  []feature `json:"features,omitempty"`
}

type feature struct {
	FeatureId string `json:"featureId,omitempty"`
}
