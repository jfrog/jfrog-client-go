package services

import (
	"errors"
	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	rthttpclient "github.com/jfrog/jfrog-client-go/artifactory/httpclient"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
)

type DeleteService struct {
	client     *rthttpclient.ArtifactoryHttpClient
	ArtDetails auth.ArtifactoryDetails
	DryRun     bool
}

func NewDeleteService(client *rthttpclient.ArtifactoryHttpClient) *DeleteService {
	return &DeleteService{client: client}
}

func (ds *DeleteService) GetArtifactoryDetails() auth.ArtifactoryDetails {
	return ds.ArtDetails
}

func (ds *DeleteService) SetArtifactoryDetails(rt auth.ArtifactoryDetails) {
	ds.ArtDetails = rt
}

func (ds *DeleteService) IsDryRun() bool {
	return ds.DryRun
}

func (ds *DeleteService) GetJfrogHttpClient() (*rthttpclient.ArtifactoryHttpClient, error) {
	return ds.client, nil
}

func (ds *DeleteService) GetPathsToDelete(deleteParams DeleteParams) (resultItems []utils.ResultItem, err error) {
	log.Info("Searching artifacts...")
	switch deleteParams.GetSpecType() {
	case utils.AQL:
		if resultItemsTemp, e := utils.SearchBySpecWithAql(deleteParams.GetFile(), ds, utils.NONE); e == nil {
			resultItems = append(resultItems, resultItemsTemp...)
		} else {
			err = e
			return
		}
	case utils.WILDCARD:
		deleteParams.SetIncludeDirs(true)
		tempResultItems, e := utils.SearchBySpecWithPattern(deleteParams.GetFile(), ds, utils.NONE)
		if e != nil {
			err = e
			return
		}
		paths := utils.ReduceDirResult(tempResultItems, utils.FilterTopChainResults)
		resultItems = append(resultItems, paths...)
	}
	utils.LogSearchResults(len(resultItems))
	return
}

func (ds *DeleteService) DeleteFiles(deleteItems []utils.ResultItem) (int, error) {
	deletedCount := 0
	for _, v := range deleteItems {
		fileUrl, err := utils.BuildArtifactoryUrl(ds.GetArtifactoryDetails().GetUrl(), v.GetItemRelativePath(), make(map[string]string))
		if err != nil {
			return deletedCount, err
		}
		if ds.IsDryRun() {
			log.Info("[Dry run] Deleting:", v.GetItemRelativePath())
			continue
		}

		log.Info("Deleting:", v.GetItemRelativePath())
		httpClientsDetails := ds.GetArtifactoryDetails().CreateHttpClientDetails()
		resp, body, err := ds.client.SendDelete(fileUrl, nil, &httpClientsDetails)
		if err != nil {
			log.Error(err)
			continue
		}
		if resp.StatusCode != http.StatusNoContent {
			log.Error(errorutils.CheckError(errors.New("Artifactory response: " + resp.Status + "\n" + clientutils.IndentJson(body))))
			continue
		}
		deletedCount++
		log.Debug("Artifactory response:", resp.Status)
	}
	return deletedCount, nil
}

type DeleteConfiguration struct {
	ArtDetails auth.ArtifactoryDetails
	DryRun     bool
}

func (conf *DeleteConfiguration) GetArtifactoryDetails() auth.ArtifactoryDetails {
	return conf.ArtDetails
}

func (conf *DeleteConfiguration) SetArtifactoryDetails(art auth.ArtifactoryDetails) {
	conf.ArtDetails = art
}

func (conf *DeleteConfiguration) IsDryRun() bool {
	return conf.DryRun
}

type DeleteParams struct {
	*utils.ArtifactoryCommonParams
}

func (ds *DeleteParams) GetFile() *utils.ArtifactoryCommonParams {
	return ds.ArtifactoryCommonParams
}

func (ds *DeleteParams) SetIncludeDirs(includeDirs bool) {
	ds.IncludeDirs = includeDirs
}

func NewDeleteParams() DeleteParams {
	return DeleteParams{ArtifactoryCommonParams: &utils.ArtifactoryCommonParams{}}
}
