package services

import (
	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/httpclient"
	"io"
)

type ReadFileService struct {
	client       *httpclient.HttpClient
	ArtDetails   auth.ArtifactoryDetails
	DryRun       bool
	MinSplitSize int64
	SplitCount   int
	Retries      int
}

func NewReadFileService(client *httpclient.HttpClient) *ReadFileService {
	return &ReadFileService{client: client}
}

func (ds *ReadFileService) GetArtifactoryDetails() auth.ArtifactoryDetails {
	return ds.ArtDetails
}

func (ds *ReadFileService) SetArtifactoryDetails(rt auth.ArtifactoryDetails) {
	ds.ArtDetails = rt
}

func (ds *ReadFileService) IsDryRun() bool {
	return ds.DryRun
}

func (ds *ReadFileService) GetJfrogHttpClient() (*httpclient.HttpClient, error) {
	return ds.client, nil
}

func (ds *ReadFileService) SetArtDetails(artDetails auth.ArtifactoryDetails) {
	ds.ArtDetails = artDetails
}

func (ds *ReadFileService) SetDryRun(isDryRun bool) {
	ds.DryRun = isDryRun
}

func (ds *ReadFileService) setMinSplitSize(minSplitSize int64) {
	ds.MinSplitSize = minSplitSize
}

func (ds *ReadFileService) ReadRemoteFile(downloadPath string) (io.ReadCloser, error) {
	readPath, err := utils.BuildArtifactoryUrl(ds.ArtDetails.GetUrl(), downloadPath, make(map[string]string))
	if err != nil {
		return nil, err
	}
	httpClientsDetails := ds.ArtDetails.CreateHttpClientDetails()
	return ds.client.ReadRemoteFile(readPath, httpClientsDetails, ds.Retries)
}
