package usage

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/jfrog/jfrog-client-go/http/httpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
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
		return errorutils.CheckErrorf("nothing to send.")
	}
	bodyContent, err := json.Marshal(reports)
	if err != nil {
		return errorutils.CheckError(err)
	}
	resp, body, err := sendRequestToEcosystemService(bodyContent)
	if err != nil {
		return errors.New("Couldn't send usage info to Ecosystem. Error: " + err.Error())
	}
	return errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK, http.StatusAccepted)
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
	client, req, err := createEcosystemRequestInfo(content)
	if err != nil {
		return
	}
	log.Debug(fmt.Sprintf("Sending HTTP %s request to: %s", req.Method, req.URL))
	resp, err = client.Do(req)
	err = errorutils.CheckError(err)
	if err != nil {
		return
	}
	if resp == nil {
		err = errorutils.CheckErrorf("Ecosystem-Usage Received empty response from server")
		return
	}
	defer func() {
		if resp.Body != nil {
			err = errors.Join(err, errorutils.CheckError(resp.Body.Close()))
		}
	}()
	respBody, _ = io.ReadAll(resp.Body)
	return
}

func createEcosystemRequestInfo(content []byte) (c *http.Client, req *http.Request, err error) {
	var client *httpclient.HttpClient
	if client, err = httpclient.ClientBuilder().Build(); err != nil {
		return
	}
	if req, err = http.NewRequest(http.MethodPost, ecosystemUsageApiPath, bytes.NewBuffer(content)); err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	return client.GetClient(), req, errorutils.CheckError(err)
}
