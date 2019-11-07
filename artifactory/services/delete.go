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
	"strings"
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
		tempResultItems, e = removeNotToBeDeletedDirs(*deleteParams.GetFile(), ds, tempResultItems)
		if e != nil {
			err = e
			return
		}
		paths := utils.ReduceDirResult(tempResultItems, utils.FilterTopChainResults)
		resultItems = append(resultItems, paths...)
	case utils.BUILD:
		if resultItemsTemp, e := utils.SearchBySpecWithBuild(deleteParams.GetFile(), ds); e == nil {
			resultItems = append(resultItems, resultItemsTemp...)
		} else {
			err = e
			return
		}
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

// This function receives as an argument the list of files and folders to be deleted from Artifactory.
// In case the search params used to create this list included excludeProps, we might need to remove some directories from this list.
// These directories must be removed, because they include files, which should not be deleted, because of the excludeProps params.
// hese directories must not be deleted from Artifactory.
func removeNotToBeDeletedDirs(specFile utils.ArtifactoryCommonParams, ds *DeleteService, deleteCandidates []utils.ResultItem) ([]utils.ResultItem, error) {
	if specFile.ExcludeProps == "" {
		return deleteCandidates, nil
	}
	specFile.Props = specFile.ExcludeProps
	specFile.ExcludeProps = ""
	remainArtifacts, err := utils.SearchBySpecWithPattern(&specFile, ds, utils.NONE)
	if err != nil {
		return nil, err
	}
	var result []utils.ResultItem
	for _, candidate := range deleteCandidates {
		deleteCandidate := true
		if candidate.Type == "folder" {
			candidatePath := candidate.GetItemRelativePath()
			for _, artifact := range remainArtifacts {
				artifactPath := artifact.GetItemRelativePath()
				if strings.HasPrefix(artifactPath, candidatePath) {
					deleteCandidate = false
					break
				}
			}
		}
		if deleteCandidate {
			result = append(result, candidate)
		}
	}
	return result, nil
}
