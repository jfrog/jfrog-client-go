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
)

const (
	cveRemediationAPI      = "api/v1/cveRemediationCDX"
	artifactRemediationAPI = "api/v1/artifactRemediationCDX"
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
func (rs *RemediationService) CveRemediation(bom *cyclonedx.BOM, cves ...string) (*RemediationResponse, error) {
	// Prepare the request body
	request := CveRemediationRequest{
		Bom:  bom,
		Cves: cves,
	}
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal CVE remediation request: %w", err)
	}
	// Get CVE remediation using the Xray service
	return rs.getRemediationResponse(rs.getUrlForRemediationApi(cveRemediationAPI), requestBody)
}

type CveRemediationRequest struct {
	Bom  *cyclonedx.BOM `json:"bom"`
	Cves []string       `json:"cves"`
}

// Get remediation for all the direct dependencies in the BOM
func (rs *RemediationService) ArtifactRemediation(bom *cyclonedx.BOM) (*RemediationResponse, error) {
	// Encode the BOM to JSON format
	encodedBom, err := utils.EncodeBomToJson(bom)
	if err != nil {
		return nil, err
	}
	// Get artifact remediation using the Xray service
	return rs.getRemediationResponse(rs.getUrlForRemediationApi(artifactRemediationAPI), encodedBom)
}

func (rs *RemediationService) getRemediationResponse(url string, requestBody []byte) (*RemediationResponse, error) {
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
	// Parse the response body into a RemediationResponse struct
	var remediationResponse RemediationResponse
	err = json.Unmarshal(body, &remediationResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to parse remediation response from Xray: %w", err)
	}
	return &remediationResponse, nil
}

func (rs *RemediationService) getUrlForRemediationApi(baseEndpoint string) string {
	return utils.AppendScopedProjectKeyParam(rs.XrayDetails.GetUrl()+baseEndpoint, rs.ScopeProjectKey)
}

type RemediationResponse struct {
}
