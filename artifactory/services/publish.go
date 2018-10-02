package services

import (
	"encoding/json"
	"errors"
	"github.com/jfrog/jfrog-cli-go/jfrog-client/artifactory/auth"
	"github.com/jfrog/jfrog-cli-go/jfrog-client/artifactory/services/utils"
	"github.com/jfrog/jfrog-cli-go/jfrog-client/httpclient"
	clientutils "github.com/jfrog/jfrog-cli-go/jfrog-client/utils"
	"github.com/jfrog/jfrog-cli-go/jfrog-client/utils/errorutils"
	"github.com/jfrog/jfrog-cli-go/jfrog-client/utils/log"
	"github.com/jfrog/jfrog-cli-go/jfrog-client/artifactory/buildinfo"
)

type buildInfoPublishService struct {
	client     *httpclient.HttpClient
	ArtDetails auth.ArtifactoryDetails
	DryRun     bool
}

func NewBuildInfoPublishService(client *httpclient.HttpClient) *buildInfoPublishService {
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
	resp, body, err := bip.client.SendPut(bip.ArtDetails.GetUrl()+"api/build/", content, httpClientsDetails)
	if err != nil {
		return err
	}
	if resp.StatusCode != 204 {
		return errorutils.CheckError(errors.New("Artifactory response: " + resp.Status + "\n" + clientutils.IndentJson(body)))
	}

	log.Debug("Artifactory response:", resp.Status)
	log.Info("Build info successfully deployed. Browse it in Artifactory under " + bip.getArtifactoryDetails().GetUrl() + "webapp/builds/" + build.Name + "/" + build.Number)
	return nil
}
