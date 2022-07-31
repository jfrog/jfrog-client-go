package services

import (
	"io"
	"net/http"

	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
)

type ReadFileService struct {
	client       *jfroghttpclient.JfrogHttpClient
	artDetails   *auth.ServiceDetails
	DryRun       bool
	MinSplitSize int64
	SplitCount   int
}

func NewReadFileService(artDetails auth.ServiceDetails, client *jfroghttpclient.JfrogHttpClient) *ReadFileService {
	return &ReadFileService{artDetails: &artDetails, client: client}
}

func (ds *ReadFileService) GetArtifactoryDetails() auth.ServiceDetails {
	return *ds.artDetails
}

func (ds *ReadFileService) IsDryRun() bool {
	return ds.DryRun
}

func (ds *ReadFileService) GetJfrogHttpClient() *jfroghttpclient.JfrogHttpClient {
	return ds.client
}

func (ds *ReadFileService) SetDryRun(isDryRun bool) {
	ds.DryRun = isDryRun
}

func (ds *ReadFileService) ReadRemoteFile(downloadPath string) (io.ReadCloser, error) {
	readPath, err := utils.BuildArtifactoryUrl(ds.GetArtifactoryDetails().GetUrl(), downloadPath, make(map[string]string))
	if err != nil {
		return nil, err
	}
	httpClientsDetails := ds.GetArtifactoryDetails().CreateHttpClientDetails()
	ioReadCloser, resp, err := ds.client.ReadRemoteFile(readPath, &httpClientsDetails)
	if err != nil {
		return nil, err
	}
	err = errorutils.CheckResponseStatus(resp, []byte{}, http.StatusOK)
	if err != nil {
		return nil, err
	}
	return ioReadCloser, err
}
