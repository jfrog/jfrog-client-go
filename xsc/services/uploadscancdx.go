package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	xscutils "github.com/jfrog/jfrog-client-go/xsc/services/utils"
)

const uploadScanCdxAPIUrl = "scan-cdx/upload"

type UploadScanCdxService struct {
	client          *jfroghttpclient.JfrogHttpClient
	XrayDetails     auth.ServiceDetails
	ScopeProjectKey string
}

type UploadScanCdxParams struct {
	RepoName string `json:"repo_name"`
	RepoPath string `json:"repo_path"`
	FileName string `json:"file_name"`
	Bom      string `json:"bom"`
}

type UploadScanCdxResponse struct {
	Repository string `json:"repository"`
	Path       string `json:"path"`
}

func NewUploadScanCdxService(client *jfroghttpclient.JfrogHttpClient) *UploadScanCdxService {
	return &UploadScanCdxService{client: client}
}

func (us *UploadScanCdxService) getUploadScanCdxURL() string {
	return utils.AppendScopedProjectKeyParam(utils.AddTrailingSlashIfNeeded(us.XrayDetails.GetUrl())+xscutils.XscInXraySuffix+uploadScanCdxAPIUrl, us.ScopeProjectKey)
}

func (us *UploadScanCdxService) Upload(params UploadScanCdxParams) (response *UploadScanCdxResponse, err error) {
	httpClientsDetails := us.XrayDetails.CreateHttpClientDetails()
	httpClientsDetails.SetContentTypeApplicationJson()

	requestBody, err := json.Marshal(params)
	if errorutils.CheckError(err) != nil {
		return
	}

	resp, body, err := us.client.SendPost(us.getUploadScanCdxURL(), requestBody, &httpClientsDetails)
	if err != nil {
		return
	}

	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusCreated); err != nil {
		err = fmt.Errorf("got unexpected server response while attempting to upload cdx for repo %s:\n%s", params.RepoName, err.Error())
		return
	}

	response = &UploadScanCdxResponse{}
	if err = json.Unmarshal(body, response); err != nil {
		err = errorutils.CheckErrorf("couldn't parse JFrog Xray upload scan cdx response: %s", err.Error())
	}
	return
}
