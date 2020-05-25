package services

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/jfrog/gofrog/parallel"
	rthttpclient "github.com/jfrog/jfrog-client-go/artifactory/httpclient"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type DeleteService struct {
	client     *rthttpclient.ArtifactoryHttpClient
	ArtDetails auth.ServiceDetails
	DryRun     bool
	Threads    int
}

func NewDeleteService(client *rthttpclient.ArtifactoryHttpClient) *DeleteService {
	return &DeleteService{client: client}
}

func (ds *DeleteService) GetArtifactoryDetails() auth.ServiceDetails {
	return ds.ArtDetails
}

func (ds *DeleteService) SetArtifactoryDetails(rt auth.ServiceDetails) {
	ds.ArtDetails = rt
}

func (ds *DeleteService) IsDryRun() bool {
	return ds.DryRun
}

func (ds *DeleteService) GetThreads() int {
	return ds.Threads
}

func (ds *DeleteService) SetThreads(threads int) {
	ds.Threads = threads
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

type fileDeleteHandlerFunc func(utils.ResultItem) parallel.TaskFunc

func (ds *DeleteService) createFileHandlerFunc(result *utils.Result) fileDeleteHandlerFunc {
	return func(resultItem utils.ResultItem) parallel.TaskFunc {
		return func(threadId int) error {
			result.TotalCount[threadId]++
			logMsgPrefix := clientutils.GetLogMsgPrefix(threadId, ds.DryRun)
			deletePath, e := utils.BuildArtifactoryUrl(ds.GetArtifactoryDetails().GetUrl(), resultItem.GetItemRelativePath(), make(map[string]string))
			if e != nil {
				return e
			}
			log.Info(logMsgPrefix+"Deleting", resultItem.GetItemRelativePath())
			if ds.DryRun {
				return nil
			}
			httpClientsDetails := ds.GetArtifactoryDetails().CreateHttpClientDetails()
			resp, body, err := ds.client.SendDelete(deletePath, nil, &httpClientsDetails)
			if err != nil {
				log.Error(err)
				return err
			}
			if resp.StatusCode != http.StatusNoContent {
				err = errors.New("Artifactory response: " + resp.Status + "\n" + clientutils.IndentJson(body))
				log.Error(errorutils.CheckError(err))
				return err
			}

			result.SuccessCount[threadId]++
			return nil
		}
	}
}

func (ds *DeleteService) DeleteFiles(deleteItems []utils.ResultItem) (int, error) {
	producerConsumer := parallel.NewBounedRunner(ds.GetThreads(), false)
	errorsQueue := clientutils.NewErrorsQueue(1)
	result := *utils.NewResult(ds.Threads)
	go func() {
		defer producerConsumer.Done()
		for _, deleteItem := range deleteItems {
			fileDeleteHandlerFunc := ds.createFileHandlerFunc(&result)
			producerConsumer.AddTaskWithError(fileDeleteHandlerFunc(deleteItem), errorsQueue.AddError)
		}
	}()
	return ds.performTasks(producerConsumer, errorsQueue, result)
}

func (ds *DeleteService) performTasks(consumer parallel.Runner, errorsQueue *clientutils.ErrorsQueue, result utils.Result) (totalDeleted int, err error) {
	consumer.Run()
	err = errorsQueue.GetError()

	totalDeleted = utils.SumIntArray(result.SuccessCount)
	log.Debug("Deleted", strconv.Itoa(totalDeleted), "artifacts.")
	return
}

type DeleteConfiguration struct {
	ArtDetails auth.ServiceDetails
	DryRun     bool
}

func (conf *DeleteConfiguration) GetArtifactoryDetails() auth.ServiceDetails {
	return conf.ArtDetails
}

func (conf *DeleteConfiguration) SetArtifactoryDetails(art auth.ServiceDetails) {
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
