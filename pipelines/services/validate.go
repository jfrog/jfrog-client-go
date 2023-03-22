package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gookit/color"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	validatePipResourcePath = "api/v1/validateYaml/json"
)

type ValidateService struct {
	client *jfroghttpclient.JfrogHttpClient
	auth.ServiceDetails
}

func NewValidateService(client *jfroghttpclient.JfrogHttpClient) *ValidateService {
	return &ValidateService{client: client}
}

func (rs *ValidateService) getHttpDetails() httputils.HttpClientDetails {
	httpDetails := rs.ServiceDetails.CreateHttpClientDetails()
	return httpDetails
}

func (vs *ValidateService) ValidatePipeline(data []byte) (string, error) {
	var opMsg string
	var err error
	httpDetails := vs.getHttpDetails()

	// query params
	m := make(map[string]string, 0)

	// URL Construction
	uri := vs.constructValidateAPIURL(m, validatePipResourcePath)

	// headers
	headers := make(map[string]string, 0)
	headers["Content-Type"] = "application/json"
	httpDetails.Headers = headers

	// send post request
	resp, body, err := vs.client.SendPost(uri, data, &httpDetails)
	if err != nil {
		return "", err
	}
	if err := errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return "", err
	}

	// process response
	if nil != resp && resp.StatusCode == http.StatusOK {
		log.Info("Processing resources yaml response ")
		time.Sleep(2 * time.Second)
		log.Info(string(body))
		opMsg, err = processValidatePipResourceResponse(body, "")
		if err != nil {
			log.Error("Failed to process resources validation response")
		}
	}
	return opMsg, nil
}

// processValidatePipResourceResponse process validate pipeline resource response
func processValidatePipResourceResponse(resp []byte, userName string) (string, error) {
	fmt.Println("unfurling response")
	rsc := make(map[string]ValidationResponse)
	err := json.Unmarshal(resp, &rsc)
	if err != nil {
		return "", err
	}
	if v, ok := rsc["pipelines.yml"]; ok {
		if v.IsValid != nil && *v.IsValid {
			userName = "@bhanur"
			log.Info("validation of pipeline resources completed successfully ")
			log.Info("workspace updated with latest resource files for user: ", userName)
			msg := color.Green.Sprintf("validation completed ")
			time.Sleep(2 * time.Second)
			log.Info(msg)
			return msg, nil
		} else {
			log.Error("pipeline resources validation FAILED!! check below errors and try again")
			msg := v.Errors[0].Text
			lnNum := v.Errors[0].LineNumber
			msgs := strings.Split(msg, ":")
			opMsg := color.Red.Sprintf("%s", msgs[0]+":"+strconv.Itoa(lnNum)+":\n"+" "+msgs[1])
			//log.PrintMessage("Please refer pipelines documentation " + coreutils.PrintLink("https://www.jfrog.com/confluence/display/JFROG/Managing+Pipeline+Sources#ManagingPipelineSources-ValidatingYAML") + "")
			log.Error(opMsg)
			return opMsg, nil
		}
	}
	return "", errors.New("pipelines.yml not found")
}

/*
constructPipelinesURL creates URL with all required details to make api call
like headers, queryParams, apiPath
*/
func (rs *ValidateService) constructValidateAPIURL(qParams map[string]string, apiPath string) string {
	uri, err := url.Parse(rs.ServiceDetails.GetUrl() + apiPath)
	if err != nil {
		log.Error("Failed to parse pipelines fetch run status url")
	}

	queryString := uri.Query()
	for k, v := range qParams {
		queryString.Set(k, v)
	}
	uri.RawQuery = queryString.Encode()

	return uri.String()
}
