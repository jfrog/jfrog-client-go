package _go

import (
	"net/http"

	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/version"
)

func init() {
	register(&publishWithMatrixParams{})
}

// Support for Artifactory version between 6.5.0 and 6.6.1 API
type publishWithMatrixParams struct {
}

func (pwmp *publishWithMatrixParams) isCompatible(artifactoryVersion string) bool {
	propertiesApi := "6.5.0"
	withoutApi := "6.6.1"
	version := version.NewVersion(artifactoryVersion)
	if version.Compare(propertiesApi) > 0 {
		return false
	}
	if version.Compare(withoutApi) <= 0 {
		return false
	}
	return true
}

func (pwmp *publishWithMatrixParams) PublishPackage(params GoParams, client *jfroghttpclient.JfrogHttpClient, ArtDetails auth.ServiceDetails) (*utils.OperationSummary, error) {
	url, err := utils.BuildArtifactoryUrl(ArtDetails.GetUrl(), "api/go/"+params.GetTargetRepo(), make(map[string]string))
	clientDetails := ArtDetails.CreateHttpClientDetails()
	addHeaders(params, &clientDetails)

	err = CreateUrlPath(params.GetModuleId(), params.GetVersion(), params.GetProps(), ".zip", &url)
	if err != nil {
		return nil, err
	}

	resp, _, err := client.UploadFile(params.GetZipPath(), url, "", &clientDetails, GoUploadRetries, nil)
	if err != nil {
		return nil, err
	}
	return nil, errorutils.CheckResponseStatus(resp, http.StatusCreated)
}
