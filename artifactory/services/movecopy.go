package services

import (
	"errors"
	"net/http"
	"path"
	"strconv"
	"strings"

	rthttpclient "github.com/jfrog/jfrog-client-go/artifactory/httpclient"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

const (
	MOVE MoveType = "move"
	COPY MoveType = "copy"
)

type MoveCopyService struct {
	moveType   MoveType
	client     *rthttpclient.ArtifactoryHttpClient
	DryRun     bool
	ArtDetails auth.ServiceDetails
}

func NewMoveCopyService(client *rthttpclient.ArtifactoryHttpClient, moveType MoveType) *MoveCopyService {
	return &MoveCopyService{moveType: moveType, client: client}
}

func (mc *MoveCopyService) GetArtifactoryDetails() auth.ServiceDetails {
	return mc.ArtDetails
}

func (mc *MoveCopyService) SetArtifactoryDetails(rt auth.ServiceDetails) {
	mc.ArtDetails = rt
}

func (mc *MoveCopyService) IsDryRun() bool {
	return mc.DryRun
}

func (mc *MoveCopyService) GetJfrogHttpClient() (*rthttpclient.ArtifactoryHttpClient, error) {
	return mc.client, nil
}

func (mc *MoveCopyService) MoveCopyServiceMoveFilesWrapper(moveSpec MoveCopyParams) (successCount, failedCount int, err error) {
	var resultItems []utils.ResultItem
	log.Info("Searching items...")

	switch moveSpec.GetSpecType() {
	case utils.BUILD:
		resultItems, err = utils.SearchBySpecWithBuild(moveSpec.GetFile(), mc)
	case utils.AQL:
		resultItems, err = utils.SearchBySpecWithAql(moveSpec.GetFile(), mc, utils.NONE)
	case utils.WILDCARD:
		moveSpec.SetIncludeDir(true)
		resultItems, err = utils.SearchBySpecWithPattern(moveSpec.GetFile(), mc, utils.NONE)
	}
	if err != nil {
		return 0, 0, err
	}

	successCount, failedCount, err = mc.moveFiles(resultItems, moveSpec)
	if err != nil {
		return 0, 0, err
	}

	log.Debug(moveMsgs[mc.moveType].MovedMsg, strconv.Itoa(successCount), "artifacts.")
	if failedCount > 0 {
		err = errorutils.CheckError(errors.New("Failed " + moveMsgs[mc.moveType].MovingMsg + " " + strconv.Itoa(failedCount) + " artifacts."))
	}
	return
}

func reduceMovePaths(resultItems []utils.ResultItem, params MoveCopyParams) []utils.ResultItem {
	if params.IsFlat() {
		return utils.ReduceDirResult(resultItems, utils.FilterBottomChainResults)
	}
	return utils.ReduceDirResult(resultItems, utils.FilterTopChainResults)
}

func (mc *MoveCopyService) moveFiles(resultItems []utils.ResultItem, params MoveCopyParams) (successCount, failedCount int, err error) {
	successCount = 0
	failedCount = 0
	resultItems = reduceMovePaths(resultItems, params)
	utils.LogSearchResults(len(resultItems))
	for _, v := range resultItems {
		destPathLocal := params.GetFile().Target
		if !params.IsFlat() {
			if strings.Contains(destPathLocal, "/") {
				file, dir := fileutils.GetFileAndDirFromPath(destPathLocal)
				destPathLocal = clientutils.TrimPath(dir + "/" + v.Path + "/" + file)
			} else {
				destPathLocal = clientutils.TrimPath(destPathLocal + "/" + v.Path + "/")
			}
		}
		destFile, err := clientutils.BuildTargetPath(params.GetFile().Pattern, v.GetItemRelativePath(), destPathLocal, true)
		if err != nil {
			log.Error(err)
			continue
		}
		if strings.HasSuffix(destFile, "/") {
			if v.Type != "folder" {
				destFile += v.Name
			} else {
				mc.createPathForMoveAction(destFile)
			}
		}
		success, err := mc.moveFile(v.GetItemRelativePath(), destFile)
		if err != nil {
			log.Error(err)
			continue
		}

		successCount += clientutils.Bool2Int(success)
		failedCount += clientutils.Bool2Int(!success)
	}
	return
}

func (mc *MoveCopyService) moveFile(sourcePath, destPath string) (bool, error) {
	message := moveMsgs[mc.moveType].MovingMsg + " artifact: " + sourcePath + " to: " + destPath
	moveUrl := mc.GetArtifactoryDetails().GetUrl()
	restApi := path.Join("api", string(mc.moveType), sourcePath)
	params := map[string]string{"to": destPath}
	if mc.IsDryRun() {
		log.Info("[Dry run]", message)
		params["dry"] = "1"
	} else {
		log.Info(message)
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

	if resp.StatusCode != http.StatusOK {
		log.Error("Artifactory response: " + resp.Status + "\n" + clientutils.IndentJson(body))
	}

	log.Debug("Artifactory response:", resp.Status)
	return resp.StatusCode == http.StatusOK, nil
}

// Create destPath in Artifactory
func (mc *MoveCopyService) createPathForMoveAction(destPath string) (bool, error) {
	if mc.IsDryRun() == true {
		log.Info("[Dry run]", "Create path:", destPath)
		return true, nil
	}

	return mc.createPathInArtifactory(destPath, mc)
}

func (mc *MoveCopyService) createPathInArtifactory(destPath string, conf utils.CommonConf) (bool, error) {
	rtUrl := conf.GetArtifactoryDetails().GetUrl()
	requestFullUrl, err := utils.BuildArtifactoryUrl(rtUrl, destPath, map[string]string{})
	if err != nil {
		return false, err
	}
	httpClientsDetails := conf.GetArtifactoryDetails().CreateHttpClientDetails()
	resp, body, err := mc.client.SendPut(requestFullUrl, nil, &httpClientsDetails)
	if err != nil {
		return false, err
	}

	if resp.StatusCode != http.StatusCreated {
		log.Error("Artifactory response: " + resp.Status + "\n" + clientutils.IndentJson(body))
	}

	log.Debug("Artifactory response:", resp.Status)
	return resp.StatusCode == http.StatusOK, nil
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
	*utils.ArtifactoryCommonParams
	Flat bool
}

func (mc *MoveCopyParams) GetFile() *utils.ArtifactoryCommonParams {
	return mc.ArtifactoryCommonParams
}

func (mc *MoveCopyParams) SetIncludeDir(isIncludeDir bool) {
	mc.GetFile().IncludeDirs = isIncludeDir
}

func (mc *MoveCopyParams) IsFlat() bool {
	return mc.Flat
}

func NewMoveCopyParams() MoveCopyParams {
	return MoveCopyParams{ArtifactoryCommonParams: &utils.ArtifactoryCommonParams{}}
}
