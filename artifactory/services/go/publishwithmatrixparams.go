package _go

import (
	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/errors/httperrors"
	"github.com/jfrog/jfrog-client-go/httpclient"
	"github.com/jfrog/jfrog-client-go/utils/version"
)

func init() {
	register(&publishWithMatrixParams{})
}

// Support for Artifactory version at least 6.5.0 and below 6.6.0
type publishWithMatrixParams struct {
}

func (pwmp *publishWithMatrixParams) isCompatible(artifactoryVersion string) (bool, error) {
	atLeastWithProps, err := version.NewVersion(artifactoryVersion).IsAtLeast(propertiesApi)
	if err != nil {
		return false, err
	}

	lessThanWithoutApi, err := version.NewVersion(artifactoryVersion).IsLessThan(withoutApi)
	if err != nil {
		return false, err
	}

	return atLeastWithProps && lessThanWithoutApi, nil
}

func (pwmp *publishWithMatrixParams) PublishPackage(params GoParams, client *httpclient.HttpClient, ArtDetails auth.ArtifactoryDetails) error {
	url, err := utils.BuildArtifactoryUrl(ArtDetails.GetUrl(), "api/go/"+params.GetTargetRepo(), make(map[string]string))
	clientDetails := ArtDetails.CreateHttpClientDetails()
	addHeaders(params, &clientDetails)

	err = createUrlPath(params.GetModuleId(), params.GetVersion(), params.GetProps(), ".zip", &url)
	if err != nil {
		return err
	}

	resp, body, err := client.UploadFile(params.GetZipPath(), url, clientDetails, 0)
	if err != nil {
		return err
	}
	return httperrors.CheckResponseStatus(resp, body, 201)
}
