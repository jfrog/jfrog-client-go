package services

import (
	"encoding/json"
	"errors"
	"net/http"

	rthttpclient "github.com/jfrog/jfrog-client-go/artifactory/httpclient"
	// "github.com/jfrog/jfrog-client-go/artifactory/services/utils"

	artUtils "github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/jfrog/jfrog-client-go/xray/services/utils"
)

const (
	WatchAPIURL = "api/v2/watches"

	WatchBuildAll    utils.WatchBuildType = "all"
	WatchBuildByName utils.WatchBuildType = "byname"

	WatchRepositoriesAll    utils.WatchRepositoriesType = "all"
	WatchRepositoriesByName utils.WatchRepositoriesType = "byname"
)

type XrayWatchService struct {
	client      *rthttpclient.ArtifactoryHttpClient
	XrayDetails auth.ServiceDetails
}

func NewXrayWatchService(client *rthttpclient.ArtifactoryHttpClient) *XrayWatchService {
	return &XrayWatchService{client: client}
}

func (xws *XrayWatchService) GetJfrogHttpClient() *rthttpclient.ArtifactoryHttpClient {
	return xws.client
}

func (xws *XrayWatchService) GetXrayWatchUrl() string {
	return xws.XrayDetails.GetUrl() + WatchAPIURL
}

func (xws *XrayWatchService) Delete(watchName string) error {
	httpClientsDetails := xws.XrayDetails.CreateHttpClientDetails()
	log.Info("Deleting watch...")
	resp, body, err := xws.client.SendDelete(xws.GetXrayWatchUrl()+"/"+watchName, nil, &httpClientsDetails)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errorutils.CheckError(errors.New("Artifactory response: " + resp.Status + "\n" + clientutils.IndentJson(body)))
	}

	log.Debug("Artifactory response:", resp.Status)
	log.Info("Done deleting watch.")
	return nil
}

func (xws *XrayWatchService) Create(params utils.XrayWatchParams) error {
	payloadBody, err := utils.CreateBody(params)
	if err != nil {
		return errorutils.CheckError(err)
	}

	content, err := json.Marshal(payloadBody)
	if err != nil {
		return errorutils.CheckError(err)
	}

	httpClientsDetails := xws.XrayDetails.CreateHttpClientDetails()
	artUtils.SetContentType("application/json", &httpClientsDetails.Headers)
	var url = xws.GetXrayWatchUrl()
	var resp *http.Response
	var respBody []byte

	log.Info("Creating watch...")
	resp, respBody, err = xws.client.SendPost(url, content, &httpClientsDetails)

	log.Info("Finished request")
	if err != nil {
		log.Info("err: " + err.Error())
		log.Error("error")
		return err
	}
	log.Info("statuscode: " + resp.Status)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return errorutils.CheckError(errors.New("Artifactory response: " + resp.Status + "\n" + clientutils.IndentJson(respBody)))
	}
	log.Debug("Artifactory response:", resp.Status)
	log.Info("Done creating watch.")
	return nil
}

func (xws *XrayWatchService) Update(params utils.XrayWatchParams) error {
	payloadBody, err := utils.CreateBody(params)

	// the update payload must not have a name
	payloadBody.GeneralData.Name = ""

	if err != nil {
		return errorutils.CheckError(err)
	}

	content, err := json.Marshal(payloadBody)
	if err != nil {
		return errorutils.CheckError(err)
	}

	httpClientsDetails := xws.XrayDetails.CreateHttpClientDetails()
	artUtils.SetContentType("application/json", &httpClientsDetails.Headers)
	var url = xws.GetXrayWatchUrl() + "/" + params.Name
	var resp *http.Response
	var respBody []byte

	log.Info("Updating watch...")
	resp, respBody, err = xws.client.SendPut(url, content, &httpClientsDetails)

	log.Info("Finished request")
	if err != nil {
		log.Info("err: " + err.Error())
		log.Error("error")
		return err
	}
	log.Info("statuscode: " + resp.Status)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return errorutils.CheckError(errors.New("Artifactory response: " + resp.Status + "\n" + clientutils.IndentJson(respBody)))
	}
	log.Debug("Artifactory response:", resp.Status)
	log.Info("Done updating watch.")
	return nil
}

func (xws *XrayWatchService) Get(watchName string) (watchResp *utils.XrayWatchParams, err error) {
	httpClientsDetails := xws.XrayDetails.CreateHttpClientDetails()
	log.Info("Getting watch...")
	resp, body, _, err := xws.client.SendGet(xws.GetXrayWatchUrl()+"/"+watchName, true, &httpClientsDetails)
	watch := utils.XrayWatchBody{}

	if err != nil {
		return &utils.XrayWatchParams{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return &utils.XrayWatchParams{}, errorutils.CheckError(errors.New("Artifactory response: " + resp.Status + "\n" + clientutils.IndentJson(body)))
	}

	err = json.Unmarshal(body, &watch)

	if err != nil {
		return &utils.XrayWatchParams{}, errors.New("failed unmarshalling watch " + watchName)
	}

	result := NewXrayWatchParams()
	result.Name = watch.GeneralData.Name
	result.Description = watch.GeneralData.Description
	result.Active = watch.GeneralData.Active
	result.Repositories = utils.XrayWatchRepositoriesParams{
		All:          utils.XrayWatchAll{},
		Repositories: map[string]utils.XrayWatchRepository{},
		XrayWatchPathFilters: utils.XrayWatchPathFilters{
			ExcludePatterns: []string{},
			IncludePatterns: []string{},
		},
	}
	result.Builds = utils.XrayWatchBuildsParams{
		All:     utils.XrayWatchBuildsAllParams{},
		ByNames: map[string]utils.XrayWatchBuildsByNameParams{},
	}
	result.Policies = watch.AssignedPolicies

	utils.UnpackWatchBody(&result, &watch)

	log.Debug("Artifactory response:", resp.Status)
	log.Info("Done getting watch.")

	return &result, nil
}

func NewXrayWatchParams() utils.XrayWatchParams {
	return utils.XrayWatchParams{}
}

func NewXrayPolicy() utils.XrayPolicy {
	return utils.XrayPolicy{}
}

func NewXrayWatchRepository(name string, binMgrID string) utils.XrayWatchRepository {
	return utils.XrayWatchRepository{
		Name:     name,
		BinMgrID: binMgrID,
	}
}
