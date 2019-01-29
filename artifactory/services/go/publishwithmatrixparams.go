package _go

import (
	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	rthttpclient "github.com/jfrog/jfrog-client-go/artifactory/httpclient"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/version"
	"net/http"
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
	if version.Compare(artifactoryVersion, propertiesApi) < 0 {
		return false
	}

	if version.Compare(artifactoryVersion, withoutApi) >= 0 {
		return false
	}
	return true
}

func (pwmp *publishWithMatrixParams) PublishPackage(params GoParams, client *rthttpclient.ArtifactoryHttpClient, ArtDetails auth.ArtifactoryDetails) error {
	url, err := utils.BuildArtifactoryUrl(ArtDetails.GetUrl(), "api/go/"+params.GetTargetRepo(), make(map[string]string))
	clientDetails := ArtDetails.CreateHttpClientDetails()
	addHeaders(params, &clientDetails)

	err = createUrlPath(params.GetModuleId(), params.GetVersion(), params.GetProps(), ".zip", &url)
	if err != nil {
		return err
	}

	resp, _, err := client.UploadFile(params.GetZipPath(), url, params.GetModuleId(), &clientDetails, GoUploadRetries)
	if err != nil {
		return err
	}
	return errorutils.CheckResponseStatus(resp, http.StatusCreated)
}
