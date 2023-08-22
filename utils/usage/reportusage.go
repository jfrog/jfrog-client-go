package usage

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/http/httpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
)

const (
	ecosystemUsageApiPath = "https://usage-ecosystem.jfrog.io/api/usage/report"
)

type ReportEcosystemUsageData struct {
	ProductId string   `json:"productId"`
	AccountId string   `json:"accountId"`
	ClientId  string   `json:"clientId,omitempty"`
	Features  []string `json:"features"`
}

func SendEcosystemUsageReports(reports ...ReportEcosystemUsageData) error {
	if len(reports) == 0 {
		return errorutils.CheckErrorf("Nothing to send.")
	}
	bodyContent, err := json.Marshal(reports)
	if err != nil {
		return errorutils.CheckError(err)
	}
	if err != nil {
		return err
	}
	resp, body, err := sendRequestToEcosystemService(bodyContent)
	if err != nil {
		return errors.New("Couldn't send usage info to Ecosystem. Error: " + err.Error())
	}

	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK, http.StatusAccepted); err != nil {
		return err
	}
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
		err = errorutils.CheckErrorf("expected at least one feature to report usage on")
	}
	return
}

func sendRequestToEcosystemService(content []byte) (resp *http.Response, respBody []byte, err error) {
	var client *httpclient.HttpClient
	if client, err = httpclient.ClientBuilder().Build(); err != nil {
		return
	}
	details := httputils.HttpClientDetails{}
	utils.AddHeader("Content-Type", "application/json", &details.Headers)
	return client.SendPost(ecosystemUsageApiPath, content, details, "Ecosystem-Usage")
}
