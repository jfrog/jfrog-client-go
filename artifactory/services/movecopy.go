package services

import (
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/jfrog/gofrog/parallel"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/content"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

const (
	MOVE MoveType = "move"
	COPY MoveType = "copy"
)

type MoveCopyService struct {
	moveType   MoveType
	client     *jfroghttpclient.JfrogHttpClient
	DryRun     bool
	artDetails *auth.ServiceDetails
	Threads    int
}

func NewMoveCopyService(artDetails auth.ServiceDetails, client *jfroghttpclient.JfrogHttpClient, moveType MoveType) *MoveCopyService {
	return &MoveCopyService{moveType: moveType, artDetails: &artDetails, client: client}
}

func (mc *MoveCopyService) GetArtifactoryDetails() auth.ServiceDetails {
	return *mc.artDetails
}

func (mc *MoveCopyService) IsDryRun() bool {
	return mc.DryRun
}

func (mc *MoveCopyService) GetJfrogHttpClient() *jfroghttpclient.JfrogHttpClient {
	return mc.client
}

func (mc *MoveCopyService) MoveCopyServiceMoveFilesWrapper(moveSpecs ...MoveCopyParams) (successCount, failedCount int, err error) {
	moveReaders := []*ReaderSpecTuple{}
	defer func() {
		for _, readerSpec := range moveReaders {
			e := readerSpec.Reader.Close()
			if err == nil {
				err = e
			}
		}
	}()
	for i, moveSpec := range moveSpecs {
		// Create reader for each spec.
		var moveReader *content.ContentReader
		moveReader, err = mc.getPathsToMove(moveSpec)
		if err != nil {
			return
		}
		moveReaders = append(moveReaders, &ReaderSpecTuple{moveReader, i})
	}

	var tempAggregatedReader *content.ContentReader
	tempAggregatedReader, err = mergeReaders(moveReaders, content.DefaultKey)
	if err != nil {
		return
	}
	defer func() {
		e := tempAggregatedReader.Close()
		if err == nil {
			err = e
		}
	}()
	aggregatedReader := tempAggregatedReader
	if mc.moveType == MOVE {
		// If move command, reduce top dir chain results.
		aggregatedReader, err = reduceMovePaths(MoveResultItem{}, tempAggregatedReader, false, false)
		if err != nil {
			return
		}
	}
	defer func() {
		e := aggregatedReader.Close()
		if err == nil {
			err = e
		}
	}()
	successCount, failedCount, err = mc.moveFiles(aggregatedReader, moveSpecs)
	if err != nil {
		return
	}

	log.Debug(moveMsgs[mc.moveType].MovedMsg, strconv.Itoa(successCount), "artifacts.")
	if failedCount > 0 {
		err = errorutils.CheckErrorf("Failed " + moveMsgs[mc.moveType].MovingMsg + " " + strconv.Itoa(failedCount) + " artifacts.")
	}

	return
}

func (mc *MoveCopyService) getPathsToMove(moveSpec MoveCopyParams) (resultItems *content.ContentReader, err error) {
	log.Info("Searching artifacts...")
	var tempResultItems *content.ContentReader
	switch moveSpec.GetSpecType() {
	case utils.BUILD:
		resultItems, err = utils.SearchBySpecWithBuild(moveSpec.GetFile(), mc)
	case utils.AQL:
		resultItems, err = utils.SearchBySpecWithAql(moveSpec.GetFile(), mc, utils.NONE)
	case utils.WILDCARD:
		moveSpec.SetIncludeDir(true)
		tempResultItems, err = utils.SearchBySpecWithPattern(moveSpec.GetFile(), mc, utils.NONE)
		if err != nil {
			return
		}
		defer func() {
			e := tempResultItems.Close()
			if err == nil {
				err = e
			}
		}()

		resultItems, err = reduceMovePaths(utils.ResultItem{}, tempResultItems, moveSpec.IsFlat(), clientutils.IsPlaceholdersUsed(moveSpec.Pattern, moveSpec.Target))
		if err != nil {
			return
		}
	}
	if err != nil {
		return
	}

	length, err := resultItems.Length()
	utils.LogSearchResults(length)
	return
}

func reduceMovePaths(readerItem utils.SearchBasedContentItem, cr *content.ContentReader, flat, withPlaceholders bool) (*content.ContentReader, error) {
	// When placeholders are used, the file path shouldn't be taken into account (or in other words, flat = true).
	if flat || withPlaceholders {
		return utils.ReduceBottomChainDirResult(readerItem, cr)
	}
	return utils.ReduceTopChainDirResult(readerItem, cr)
}

func (mc *MoveCopyService) moveFiles(reader *content.ContentReader, params []MoveCopyParams) (successCount, failedCount int, err error) {
	promptMoveCopyMessage(reader, mc.moveType)
	producerConsumer := parallel.NewBounedRunner(mc.GetThreads(), false)
	errorsQueue := clientutils.NewErrorsQueue(1)
	result := *utils.NewResult(mc.Threads)
	go func() {
		defer producerConsumer.Done()
		for resultItem := new(MoveResultItem); reader.NextRecord(resultItem) == nil; resultItem = new(MoveResultItem) {
			fileMoveCopyHandlerFunc := mc.createMoveCopyFileHandlerFunc(&result)
			_, _ = producerConsumer.AddTaskWithError(fileMoveCopyHandlerFunc(resultItem.ResultItem, &params[resultItem.FileSpecId]),
				errorsQueue.AddError)
		}
		if err := reader.GetError(); err != nil {
			errorsQueue.AddError(err)
		}
	}()
	return mc.performTasks(producerConsumer, errorsQueue, result)
}

func (mc *MoveCopyService) performTasks(consumer parallel.Runner, errorsQueue *clientutils.ErrorsQueue, result utils.Result) (totalSuccess, totalFails int, err error) {
	consumer.Run()
	err = errorsQueue.GetError()
	totalSuccess = utils.SumIntArray(result.SuccessCount)
	totalFails = utils.SumIntArray(result.TotalCount) - totalSuccess
	return
}

type fileMoveCopyHandlerFunc func(utils.ResultItem, *MoveCopyParams) parallel.TaskFunc

func (mc *MoveCopyService) createMoveCopyFileHandlerFunc(result *utils.Result) fileMoveCopyHandlerFunc {
	return func(resultItem utils.ResultItem, params *MoveCopyParams) parallel.TaskFunc {
		return func(threadId int) error {
			result.TotalCount[threadId]++
			logMsgPrefix := clientutils.GetLogMsgPrefix(threadId, mc.DryRun)

			// Get destination path.
			destFile, err := getDestinationPath(params.GetFile().Target, params.GetFile().Pattern, resultItem.Path,
				resultItem.GetItemRelativePath(), params.IsFlat())
			if err != nil {
				return err
			}
			if strings.HasSuffix(destFile, "/") {
				if resultItem.Type != "folder" {
					destFile += resultItem.Name
				} else {
					_, err = mc.createPathForMoveAction(destFile, logMsgPrefix)
					if err != nil {
						return err
					}
				}
			}

			// Perform move/copy.
			success, err := mc.moveOrCopyFile(resultItem.GetItemRelativePath(), destFile, logMsgPrefix)
			if err != nil {
				log.Error(err)
				return err
			}
			if success {
				result.SuccessCount[threadId]++
			}
			return nil
		}
	}
}

// Create the destination path of the move/copy.
func getDestinationPath(specTarget, specPattern, sourceItemPath, sourceItemRelativePath string, isFlat bool) (string, error) {
	// Apply placeholders.
	destFile, placeholdersUsed, err := clientutils.BuildTargetPath(specPattern, sourceItemRelativePath, specTarget, true)
	if err != nil {
		return "", err
	}

	// Create raw destination path.
	// When placeholders are used, the file path shouldn't be taken into account (or in other words, flat = true).
	if !isFlat && !placeholdersUsed {
		destPathLocal := specTarget
		if strings.Contains(destPathLocal, "/") {
			file, dir := fileutils.GetFileAndDirFromPath(destPathLocal)
			destPathLocal = clientutils.TrimPath(dir + "/" + sourceItemPath + "/" + file)
		} else {
			destPathLocal = clientutils.TrimPath(destPathLocal + "/" + sourceItemPath + "/")
		}
		return destPathLocal, nil
	}

	return destFile, nil
}

func (mc *MoveCopyService) moveOrCopyFile(sourcePath, destPath, logMsgPrefix string) (bool, error) {
	message := moveMsgs[mc.moveType].MovingMsg + " artifact: " + sourcePath + " to: " + destPath
	moveUrl := mc.GetArtifactoryDetails().GetUrl()
	restApi := path.Join("api", string(mc.moveType), sourcePath)
	params := map[string]string{"to": destPath}
	if mc.IsDryRun() {
		log.Info(logMsgPrefix+"[Dry run]", message)
		params["dry"] = "1"
	} else {
		log.Info(logMsgPrefix + message)
	}
	requestFullUrl, err := utils.BuildArtifactoryUrl(moveUrl, restApi, params)
	if err != nil {
		return false, err
	}
	httpClientsDetails := mc.GetArtifactoryDetails().CreateHttpClientDetails()

	resp, body, err := mc.client.SendPost(requestFullUrl, nil, &httpClientsDetails)
	if err != nil {
		return false, err
	}

	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		log.Error(err)
	}

	log.Debug(logMsgPrefix+"Artifactory response:", resp.Status)
	return resp.StatusCode == http.StatusOK, nil
}

// Create destPath in Artifactory
func (mc *MoveCopyService) createPathForMoveAction(destPath, logMsgPrefix string) (bool, error) {
	if mc.IsDryRun() {
		log.Info(logMsgPrefix+"[Dry run]", "Create path:", destPath)
		return true, nil
	}

	return mc.createPathInArtifactory(destPath, logMsgPrefix)
}

func (mc *MoveCopyService) createPathInArtifactory(destPath, logMsgPrefix string) (bool, error) {
	rtUrl := mc.GetArtifactoryDetails().GetUrl()
	requestFullUrl, err := utils.BuildArtifactoryUrl(rtUrl, destPath, map[string]string{})
	if err != nil {
		return false, err
	}
	httpClientsDetails := mc.GetArtifactoryDetails().CreateHttpClientDetails()
	resp, body, err := mc.client.SendPut(requestFullUrl, nil, &httpClientsDetails)
	if err != nil {
		return false, err
	}

	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusCreated); err != nil {
		log.Error(err)
	}

	log.Debug(logMsgPrefix+"Artifactory response:", resp.Status)
	return resp.StatusCode == http.StatusOK, nil
}

// Receives multiple 'ReaderSpecTuple' items and merge them into a single 'ContentReader' of 'MoveResultItem'.
// Each item in the reader, keeps the index of its corresponding MoveSpec.
func mergeReaders(arr []*ReaderSpecTuple, arrayKey string) (contentReader *content.ContentReader, err error) {
	cw, err := content.NewContentWriter(arrayKey, true, false)
	if err != nil {
		return nil, err
	}
	defer func() {
		e := cw.Close()
		if err == nil {
			err = e
		}
	}()
	for _, tuple := range arr {
		cr := tuple.Reader
		for item := new(utils.ResultItem); cr.NextRecord(item) == nil; item = new(utils.ResultItem) {
			writeItem := &MoveResultItem{*item, tuple.MoveSpec}
			cw.Write(*writeItem)
		}
		if err := cr.GetError(); err != nil {
			return nil, err
		}
	}
	contentReader = content.NewContentReader(cw.GetFilePath(), arrayKey)
	return contentReader, nil
}

func promptMoveCopyMessage(reader *content.ContentReader, moveType MoveType) {
	length, err := reader.Length()
	if err != nil {
		return
	}
	var msgSuffix = "artifacts."
	if length == 1 {
		msgSuffix = "artifact."
	}
	log.Info("Preparing to", moveType, strconv.Itoa(length), msgSuffix)
}

var moveMsgs = map[MoveType]MoveOptions{
	MOVE: {MovingMsg: "Moving", MovedMsg: "Moved"},
	COPY: {MovingMsg: "Copying", MovedMsg: "Copied"},
}

type MoveOptions struct {
	MovingMsg string
	MovedMsg  string
}

type MoveType string

type MoveCopyParams struct {
	*utils.CommonParams
	Flat bool
}

// Tuple of a 'ResultItem' and its corresponding file-spec's index.
// We have to keep the file-spec index for each item as the file-spec's data is required for the actual move/copy, and
// this operation uses 'content.ContentReader' to hold all items.
// This is the item used in the 'ContentReader' and 'ContentWriter' of the move/copy.
type MoveResultItem struct {
	utils.ResultItem `json:"resultItem,omitempty"`
	FileSpecId       int `json:"fileSpecId,omitempty"`
}

// Tuple of a ContentReader and its corresponding file-spec index.
type ReaderSpecTuple struct {
	Reader   *content.ContentReader
	MoveSpec int
}

func (mc *MoveCopyParams) GetFile() *utils.CommonParams {
	return mc.CommonParams
}

func (mc *MoveCopyParams) SetIncludeDir(isIncludeDir bool) {
	mc.GetFile().IncludeDirs = isIncludeDir
}

func (mc *MoveCopyParams) IsFlat() bool {
	return mc.Flat
}

func (mc *MoveCopyService) GetThreads() int {
	return mc.Threads
}

func (mc *MoveCopyService) SetThreads(threads int) {
	mc.Threads = threads
}

func NewMoveCopyParams() MoveCopyParams {
	return MoveCopyParams{CommonParams: &utils.CommonParams{}}
}
