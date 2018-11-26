package _go

import (
	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/errors/httperrors"
	"github.com/jfrog/jfrog-client-go/httpclient"
	"github.com/jfrog/jfrog-client-go/utils/version"
)

func init() {
	register(&publishWithHeader{})
}

// Support for Artifactory older then 6.5.0 API
type publishWithHeader struct {
}

func (pwh *publishWithHeader) isCompatible(artifactoryVersion string) bool {
	propertiesApi := "6.5.0"
	if version.Compare(artifactoryVersion, propertiesApi) < 0 {
		return true
	}
	return false
}

func (pwh *publishWithHeader) PublishPackage(params GoParams, client *httpclient.HttpClient, ArtDetails auth.ArtifactoryDetails) error {
	url, err := utils.BuildArtifactoryUrl(ArtDetails.GetUrl(), "api/go/"+params.GetTargetRepo(), make(map[string]string))
	clientDetails := ArtDetails.CreateHttpClientDetails()
	addHeaders(params, &clientDetails)
	addPropertiesHeaders(params.GetProps(), &clientDetails.Headers)
	resp, body, err := client.UploadFile(params.GetZipPath(), url, clientDetails, 0)
	if err != nil {
		return err
	}
	return httperrors.CheckResponseStatus(resp, body, 201)
}

func addPropertiesHeaders(props string, headers *map[string]string) error {
	properties, err := utils.ParseProperties(props, utils.JoinCommas)
	if err != nil {
		return err
	}
	headersMap := properties.ToHeadersMap()
	for k, v := range headersMap {
		utils.AddHeader("X-ARTIFACTORY-PROPERTY-"+k, v, headers)
	}
	return nil
}
