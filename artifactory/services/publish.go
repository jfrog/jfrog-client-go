package services

import (
	"encoding/json"
	"errors"
	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	"github.com/jfrog/jfrog-client-go/artifactory/buildinfo"
	rthttpclient "github.com/jfrog/jfrog-client-go/artifactory/httpclient"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
)

type buildInfoPublishService struct {
	client     *rthttpclient.ArtifactoryHttpClient
	ArtDetails auth.ArtifactoryDetails
	DryRun     bool
}

func NewBuildInfoPublishService(client *rthttpclient.ArtifactoryHttpClient) *buildInfoPublishService {
	return &buildInfoPublishService{client: client}
}

func (bip *buildInfoPublishService) getArtifactoryDetails() auth.ArtifactoryDetails {
	return bip.ArtDetails
}

func (bip *buildInfoPublishService) isDryRun() bool {
	return bip.DryRun
}

func (bip *buildInfoPublishService) PublishBuildInfo(build *buildinfo.BuildInfo) error {
	content, err := json.Marshal(build)
	if errorutils.CheckError(err) != nil {
		return err
	}
	if bip.isDryRun() {
		log.Output(clientutils.IndentJson(content))
		return nil
	}
	httpClientsDetails := bip.getArtifactoryDetails().CreateHttpClientDetails()
	utils.SetContentType("application/vnd.org.jfrog.artifactory+json", &httpClientsDetails.Headers)
	log.Info("Deploying build info...")
	resp, body, err := bip.client.SendPut(bip.ArtDetails.GetUrl()+"api/build/", content, &httpClientsDetails)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusNoContent {
		return errorutils.CheckError(errors.New("Artifactory response: " + resp.Status + "\n" + clientutils.IndentJson(body)))
	}

	log.Debug("Artifactory response:", resp.Status)
	log.Info("Build info successfully deployed. Browse it in Artifactory under " + bip.getArtifactoryDetails().GetUrl() + "webapp/builds/" + build.Name + "/" + build.Number)
	return nil
}
