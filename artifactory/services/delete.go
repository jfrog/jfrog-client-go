package services

import (
	"net/http"
	"strconv"

	"github.com/jfrog/gofrog/parallel"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/content"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type DeleteService struct {
	client     *jfroghttpclient.JfrogHttpClient
	artDetails *auth.ServiceDetails
	DryRun     bool
	Threads    int
}

func NewDeleteService(artDetails auth.ServiceDetails, client *jfroghttpclient.JfrogHttpClient) *DeleteService {
	return &DeleteService{artDetails: &artDetails, client: client}
}

func (ds *DeleteService) GetArtifactoryDetails() auth.ServiceDetails {
	return *ds.artDetails
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

func (ds *DeleteService) GetJfrogHttpClient() *jfroghttpclient.JfrogHttpClient {
	return ds.client
}

func (ds *DeleteService) GetPathsToDelete(deleteParams DeleteParams) (resultItems *content.ContentReader, err error) {
	log.Info("Searching artifacts...")
	var tempResultItems, toBeDeletedDirs *content.ContentReader
	switch deleteParams.GetSpecType() {
	case utils.AQL:
		resultItems, err = utils.SearchBySpecWithAql(deleteParams.GetFile(), ds, utils.NONE)
		if err != nil {
			return
		}
	case utils.WILDCARD:
		deleteParams.SetIncludeDirs(true)
		tempResultItems, err = utils.SearchBySpecWithPattern(deleteParams.GetFile(), ds, utils.NONE)
		if err != nil {
			return
		}
		defer func() {
			e := tempResultItems.Close()
			if err == nil {
				err = e
			}
		}()
		toBeDeletedDirs, err = removeNotToBeDeletedDirs(deleteParams.GetFile(), ds, tempResultItems)
		if err != nil {
			return
		}
		// The 'removeNotToBeDeletedDirs' should filter out any folders that should not be deleted, if no action is needed, nil will be returned.
		// As a result, we should keep the flow with tempResultItems reader instead.
		if toBeDeletedDirs == nil {
			toBeDeletedDirs = tempResultItems
		}
		defer func() {
			e := toBeDeletedDirs.Close()
			if err == nil {
				err = e
			}
		}()
		resultItems, err = utils.ReduceTopChainDirResult(utils.ResultItem{}, toBeDeletedDirs)
		if err != nil {
			return
		}
	case utils.BUILD:
		resultItems, err = utils.SearchBySpecWithBuild(deleteParams.GetFile(), ds)
	}
	length, err := resultItems.Length()
	if err != nil {
		return
	}
	utils.LogSearchResults(length)
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
			if err = errorutils.CheckResponseStatus(resp, body, http.StatusNoContent); err != nil {
				log.Error(err)
				return err
			}
			result.SuccessCount[threadId]++
			return nil
		}
	}
}

func (ds *DeleteService) DeleteFiles(deleteItems *content.ContentReader) (int, error) {
	producerConsumer := parallel.NewBounedRunner(ds.GetThreads(), false)
	errorsQueue := clientutils.NewErrorsQueue(1)
	result := *utils.NewResult(ds.Threads)
	go func() {
		defer producerConsumer.Done()
		for deleteItem := new(utils.ResultItem); deleteItems.NextRecord(deleteItem) == nil; deleteItem = new(utils.ResultItem) {
			fileDeleteHandlerFunc := ds.createFileHandlerFunc(&result)
			_, _ = producerConsumer.AddTaskWithError(fileDeleteHandlerFunc(*deleteItem), errorsQueue.AddError)
		}
		if err := deleteItems.GetError(); err != nil {
			errorsQueue.AddError(err)
		}
		deleteItems.Reset()
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
	*utils.CommonParams
}

func (ds *DeleteParams) GetFile() *utils.CommonParams {
	return ds.CommonParams
}

func (ds *DeleteParams) SetIncludeDirs(includeDirs bool) {
	ds.IncludeDirs = includeDirs
}

func NewDeleteParams() DeleteParams {
	return DeleteParams{CommonParams: &utils.CommonParams{}}
}

// This function receives as an argument a reader within the list of files and dirs to be deleted from Artifactory.
// In case the search params used to create this list included excludeProps, we might need to remove some directories from this list.
// These directories must be removed, because they include files, which should not be deleted, because of the excludeProps params.
// These directories must not be deleted from Artifactory.
// In case of no excludeProps filed in the file spec, nil will be return so all deleteCandidates will get deleted.
func removeNotToBeDeletedDirs(specFile *utils.CommonParams, ds *DeleteService, deleteCandidates *content.ContentReader) (contentReader *content.ContentReader, err error) {
	length, err := deleteCandidates.Length()
	if err != nil || specFile.ExcludeProps == "" || length == 0 {
		return nil, err
	}
	// Send AQL to get all artifacts that include the exclude props.
	resultWriter, err := content.NewContentWriter(content.DefaultKey, true, false)
	if err != nil {
		return nil, err
	}
	bufferFiles, err := utils.FilterCandidateToBeDeleted(deleteCandidates, resultWriter, "folder")
	if len(bufferFiles) > 0 {
		defer func() {
			for _, file := range bufferFiles {
				e := file.Close()
				if err == nil {
					err = errorutils.CheckError(e)
				}
			}
		}()
		if err != nil {
			return nil, err
		}
		var artifactNotToBeDeleteReader *content.ContentReader
		artifactNotToBeDeleteReader, err := getSortedArtifactsToNotDelete(specFile, ds)
		if err != nil {
			return nil, err
		}
		defer func() {
			e := artifactNotToBeDeleteReader.Close()
			if err == nil {
				err = e
			}
		}()
		if err = utils.WriteCandidateDirsToBeDeleted(bufferFiles, artifactNotToBeDeleteReader, resultWriter); err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}
	if err = resultWriter.Close(); err != nil {
		return nil, err
	}
	return content.NewContentReader(resultWriter.GetFilePath(), content.DefaultKey), err
}

func getSortedArtifactsToNotDelete(specFile *utils.CommonParams, ds *DeleteService) (*content.ContentReader, error) {
	specFile.Props = specFile.ExcludeProps
	specFile.ExcludeProps = ""
	tempResults, err := utils.SearchBySpecWithPattern(specFile, ds, utils.NONE)
	if err != nil {
		return nil, err
	}
	resultWriter, err := content.NewContentWriter(content.DefaultKey, true, false)
	if err != nil {
		return nil, err
	}
	// Note that we have to sort the result by ourselves, and not relay on Artifactory's OrderBy, because of 2 main reasons:
	// 1. Go sorts strings differently from Artifactory's database, when the strings include special chars, such as dashes.
	// 2. Artifactory sorts by database columns, so directories will be sorted differently than files,
	//    because the path and name cols have different values.
	sortedResults, err := utils.FilterCandidateToBeDeleted(tempResults, resultWriter, "file")
	if err != nil {
		return nil, err
	}
	return content.MergeSortedReaders(utils.ResultItem{}, sortedResults, true)
}
