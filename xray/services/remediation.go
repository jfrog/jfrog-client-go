package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	xrayUtils "github.com/jfrog/jfrog-client-go/xray/services/utils"
)

const (
	cveRemediationAPI = "api/v1/cveRemediationCDX"
)

type RemediationService struct {
	client          *jfroghttpclient.JfrogHttpClient
	XrayDetails     auth.ServiceDetails
	ScopeProjectKey string
}

func NewRemediationService(client *jfroghttpclient.JfrogHttpClient) *RemediationService {
	return &RemediationService{client: client}
}

// Get remediation for the specified CVEs in the context of the provided BOM
func (rs *RemediationService) RemediationByCve(bom *cyclonedx.BOM) (xrayUtils.CveRemediationResponse, error) {
	// Encode the BOM to JSON format
	encodedBom, err := utils.EncodeBomToJson(bom)
	if err != nil {
		return nil, fmt.Errorf("failed to encode CycloneDX BOM: %w", err)
	}
	// Get CVE remediation using the Xray service
	body, err := rs.getRemediationResponse(rs.getUrlForRemediationApi(cveRemediationAPI), encodedBom)
	if err != nil {
		return nil, err
	}
	// Decode the response back to a CveRemediationResponse object
	var response xrayUtils.CveRemediationResponse
	if err = errorutils.CheckError(json.Unmarshal(body, &response)); err != nil {
		return nil, fmt.Errorf("failed to decode CVE remediation response: %w", err)
	}
	return response, nil
}

func (rs *RemediationService) getRemediationResponse(url string, requestBody []byte) ([]byte, error) {
	httpDetails := rs.XrayDetails.CreateHttpClientDetails()
	resp, body, err := rs.client.SendPost(url, requestBody, &httpDetails)
	if err != nil {
		return nil, fmt.Errorf("failed while attempting to get remediation from Xray: %w", err)
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		if resp.StatusCode == http.StatusUnauthorized {
			err = fmt.Errorf("%s\nHint: It appears that the credentials provided do not have sufficient permissions for JFrog Xray. This could be due to either incorrect credentials or limited permissions restricted to Artifactory only", err.Error())
		}
		return nil, fmt.Errorf("got unexpected Xray server response while attempting to get remediation:\n%w", err)
	}
	return body, nil
}

func (rs *RemediationService) getUrlForRemediationApi(baseEndpoint string) string {
	return utils.AppendScopedProjectKeyParam(rs.XrayDetails.GetUrl()+baseEndpoint, rs.ScopeProjectKey)
}
