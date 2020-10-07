package services

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"sort"
	"time"

	rthttpclient "github.com/jfrog/jfrog-client-go/artifactory/httpclient"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

const (
	WatchBuildAll    WatchBuildType = "all"
	WatchBuildByName WatchBuildType = "byname"

	WatchRepositoriesAll    WatchRepositoriesType = "all"
	WatchRepositoriesByName WatchRepositoriesType = "byname"
)

const WATCH_API_URL = "api/v2/watches"
const XRAY_WATCH_RETRY_CONSECUTIVE_RETRIES = 10           // Retrying to resume the watch 10 times after a stable connection
const XRAY_WATCH_CONNECTION_TIMEOUT = 90 * time.Second    // Expecting \r\n every 30 seconds
const XRAY_WATCH_SLEEP_BETWEEN_RETRIES = 15 * time.Second // 15 seconds sleep between retry
const XRAY_WATCH_STABLE_CONNECTION_WINDOW = 100 * time.Second
const XRAY_WATCH_FATAL_FAIL_STATUS = -1

type XrayWatchService struct {
	client     *rthttpclient.ArtifactoryHttpClient
	ArtDetails auth.ServiceDetails
}

func NewXrayWatchService(client *rthttpclient.ArtifactoryHttpClient) *XrayWatchService {
	return &XrayWatchService{client: client}
}

func (xws *XrayWatchService) GetJfrogHttpClient() *rthttpclient.ArtifactoryHttpClient {
	return xws.client
}

func (xws *XrayWatchService) GetXrayWatchUrl() string {
	url := xws.ArtDetails.GetUrl()
	r := regexp.MustCompile("artifactory/?$")

	return r.ReplaceAllString(url, "xray/") + WATCH_API_URL
}

func (xws *XrayWatchService) Delete(watchName string) error {
	httpClientsDetails := xws.ArtDetails.CreateHttpClientDetails()
	log.Info("Deleting watch...")
	resp, body, err := xws.client.SendDelete(xws.GetXrayWatchUrl()+"/"+watchName, nil, &httpClientsDetails)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errorutils.CheckError(errors.New("Artifactory response: " + resp.Status + "\n" + clientutils.IndentJson(body)))
	}

	log.Debug("Artifactory response:", resp.Status)
	log.Info("Done deleting watches.")
	return nil
}

func (xws *XrayWatchService) Create(params XrayWatchParams) error {
	payloadBody, err := CreateBody(params)
	if err != nil {
		return errorutils.CheckError(err)
	}

	content, err := json.Marshal(payloadBody)
	if err != nil {
		return errorutils.CheckError(err)
	}

	httpClientsDetails := xws.ArtDetails.CreateHttpClientDetails()
	utils.SetContentType("application/json", &httpClientsDetails.Headers)
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

func CreateBody(params XrayWatchParams) (*XrayWatchBody, error) {
	payloadBody := XrayWatchBody{
		GeneralData: XrayWatchGeneralParams{
			Name:        params.Name,
			Description: params.Description,
			Active:      params.Active,
		},
		ProjectResources: XrayWatchProjectResources{
			Resources: []XrayWatchProjectResourcesElement{},
		},
		AssignedPolicies: params.Policies,
	}

	err := ConfigureRepositories(&payloadBody, params)
	if err != nil {
		return nil, err
	}

	err = ConfigureBuilds(&payloadBody, params)
	if err != nil {
		return nil, err
	}

	err = ConfigureBundles(&payloadBody, params)
	if err != nil {
		return nil, err
	}

	return &payloadBody, nil
}

func ConfigureRepositories(payloadBody *XrayWatchBody, params XrayWatchParams) error {
	if params.Repositories.Type == WatchRepositoriesAll {
		allFilters := XrayWatchProjectResourcesElement{
			Type:          "all-repos",
			StringFilters: []XrayWatchFilterString{},
		}

		allFilters.StringFilters = append(allFilters.StringFilters, CreateFilters(params.Repositories.All.Filters, params.Repositories)...)

		payloadBody.ProjectResources.Resources = append(payloadBody.ProjectResources.Resources, allFilters)
	} else if params.Repositories.Type == WatchRepositoriesByName {
		for _, repository := range params.Repositories.Repositories {
			repo := XrayWatchProjectResourcesElement{
				Type:          "repository",
				Name:          repository.Name,
				Bin_Mgr_ID:    repository.Bin_Mgr_ID,
				StringFilters: repository.StringFilters,
			}
			if repo.StringFilters == nil {
				repo.StringFilters = []XrayWatchFilterString{}
			}
			repo.StringFilters = append(repo.StringFilters, CreateFilters(repository.Filters, params.Repositories)...)

			payloadBody.ProjectResources.Resources = append(payloadBody.ProjectResources.Resources, repo)
		}
	}

	return nil
}

func CreateFilters(filters XrayWatchFilters, repo XrayWatchRepositoriesParams) []XrayWatchFilterString {
	result := []XrayWatchFilterString{}

	for _, packageType := range filters.PackageTypes {
		filter := XrayWatchFilterString{
			Type:  "package-type",
			Value: packageType,
		}
		result = append(result, filter)
	}

	for _, name := range filters.Names {
		filter := XrayWatchFilterString{
			Type:  "regex",
			Value: name,
		}
		result = append(result, filter)
	}

	for _, path := range filters.Paths {
		filter := XrayWatchFilterString{
			Type:  "path-regex",
			Value: path,
		}
		result = append(result, filter)
	}

	for _, mimeType := range filters.MimeTypes {
		filter := XrayWatchFilterString{
			Type:  "mime-type",
			Value: mimeType,
		}
		result = append(result, filter)
	}

	for key, value := range filters.Properties {
		filter := XrayWatchFilterString{
			Type: "property",
			Value: XrayWatchFilterPropertyValue{
				Key:   key,
				Value: value,
			},
		}
		result = append(result, filter)
	}

	if repo.ExcludePatterns != nil || repo.IncludePatterns != nil {
		filter := XrayWatchFilterString{
			Type: "path-ant-patterns",
			Value: XrayWatchPathFilters{
				ExcludePatterns: repo.ExcludePatterns,
				IncludePatterns: repo.IncludePatterns,
			},
		}
		result = append(result, filter)
	}

	return result
}

func ConfigureBuilds(payloadBody *XrayWatchBody, params XrayWatchParams) error {
	if params.Builds.Type == WatchBuildAll {
		allBuilds := XrayWatchProjectResourcesElement{
			Name:          "All Builds",
			Type:          "all-builds",
			Bin_Mgr_ID:    params.Builds.All.Bin_Mgr_ID,
			StringFilters: []XrayWatchFilterString{},
		}

		if params.Builds.All.ExcludePatterns != nil || params.Builds.All.IncludePatterns != nil {
			filters := []XrayWatchFilterString{{
				Type: "ant-patterns",
				Value: XrayWatchPathFilters{
					ExcludePatterns: params.Builds.All.ExcludePatterns,
					IncludePatterns: params.Builds.All.IncludePatterns,
				}},
			}
			allBuilds.StringFilters = filters
		}

		payloadBody.ProjectResources.Resources = append(payloadBody.ProjectResources.Resources, allBuilds)
	} else if params.Builds.Type == WatchBuildByName {
		for _, byName := range params.Builds.ByNames {
			build := XrayWatchProjectResourcesElement{
				Type:       "build",
				Name:       byName.Name,
				Bin_Mgr_ID: byName.Bin_Mgr_ID,
			}

			payloadBody.ProjectResources.Resources = append(payloadBody.ProjectResources.Resources, build)
		}
	}

	return nil
}

func ConfigureBundles(payloadBody *XrayWatchBody, params XrayWatchParams) error {
	// to be implemented
	return nil
}

func (xws *XrayWatchService) Update(params XrayWatchParams) error {
	payloadBody, err := CreateBody(params)

	// the update payload must not have a name
	payloadBody.GeneralData.Name = ""

	if err != nil {
		return errorutils.CheckError(err)
	}

	content, err := json.Marshal(payloadBody)
	if err != nil {
		return errorutils.CheckError(err)
	}

	httpClientsDetails := xws.ArtDetails.CreateHttpClientDetails()
	utils.SetContentType("application/json", &httpClientsDetails.Headers)
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

func (xws *XrayWatchService) Get(watchName string) (watchesResp *XrayWatchParams, err error) {
	httpClientsDetails := xws.ArtDetails.CreateHttpClientDetails()
	log.Info("Getting watch...")
	resp, body, _, err := xws.client.SendGet(xws.GetXrayWatchUrl()+"/"+watchName, true, &httpClientsDetails)
	watches := XrayWatchBody{}

	if err != nil {
		return &XrayWatchParams{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return &XrayWatchParams{}, errorutils.CheckError(errors.New("Artifactory response: " + resp.Status + "\n" + clientutils.IndentJson(body)))
	}

	err = json.Unmarshal(body, &watches)

	if err != nil {
		return &XrayWatchParams{}, errors.New("failed unmarshalling watches " + watchName)
	}

	result := XrayWatchParams{
		Name:        watches.GeneralData.Name,
		Description: watches.GeneralData.Description,
		Active:      watches.GeneralData.Active,
		Repositories: XrayWatchRepositoriesParams{
			Type:         "",             // WatchRepositoriesType
			All:          XrayWatchAll{}, // XrayWatchAll
			Repositories: map[string]XrayWatchRepository{},
			XrayWatchPathFilters: XrayWatchPathFilters{
				ExcludePatterns: []string{},
				IncludePatterns: []string{},
			},
		},
		Builds: XrayWatchBuildsParams{
			Type:    "", //       WatchBuildType
			All:     XrayWatchBuildsAllParams{},
			ByNames: map[string]XrayWatchBuildsByNameParams{},
		},
		Policies: watches.AssignedPolicies,
	}

	unpackWatchBody(&result, &watches)

	log.Debug("Artifactory response:", resp.Status)
	log.Info("Done getting watches.")

	return &result, nil
}

func unpackWatchBody(watches *XrayWatchParams, body *XrayWatchBody) {
	for _, resource := range body.ProjectResources.Resources {
		if resource.Type == "all-repos" {
			watches.Repositories.Type = WatchRepositoriesAll
			unpackFilters(resource.StringFilters, &watches.Repositories.All.Filters, &watches.Repositories)
		}
		if resource.Type == "repository" {
			watches.Repositories.Type = WatchRepositoriesByName
			repository := XrayWatchRepository{
				Name:       resource.Name,
				Bin_Mgr_ID: resource.Bin_Mgr_ID,
			}
			unpackFilters(resource.StringFilters, &repository.Filters, &watches.Repositories)
			watches.Repositories.Repositories[repository.Name] = repository
		}
		if resource.Type == "all-builds" {
			watches.Builds.Type = WatchBuildAll
			watches.Builds.All.Bin_Mgr_ID = resource.Bin_Mgr_ID

			for _, filter := range resource.StringFilters {
				if filter.Type == "ant-patterns" {
					pathFilters := filter.Value.(map[string]interface{})

					if pathFilters["ExcludePatterns"] != nil {
						for _, path := range pathFilters["ExcludePatterns"].([]interface{}) {
							watches.Builds.All.ExcludePatterns = append(watches.Builds.All.ExcludePatterns, path.(string))
						}
					}
					if pathFilters["IncludePatterns"] != nil {
						for _, path := range pathFilters["IncludePatterns"].([]interface{}) {
							watches.Builds.All.IncludePatterns = append(watches.Builds.All.IncludePatterns, path.(string))
						}
					}
				}
			}

		}
		if resource.Type == "build" {
			watches.Builds.Type = WatchBuildByName
			watches.Builds.ByNames[resource.Name] = XrayWatchBuildsByNameParams{
				Name:       resource.Name,
				Bin_Mgr_ID: resource.Bin_Mgr_ID,
			}
		}
	}

	// Sort all the properties so they are returned in a consistent format

	sort.Strings(watches.Repositories.ExcludePatterns)
	sort.Strings(watches.Repositories.IncludePatterns)
}

func unpackFilters(filters []XrayWatchFilterString, output *XrayWatchFilters, repos *XrayWatchRepositoriesParams) {

	for _, filter := range filters {
		if filter.Type == "package-type" {
			output.PackageTypes = append(output.PackageTypes, filter.Value.(string))
		}
		if filter.Type == "regex" {
			output.Names = append(output.Names, filter.Value.(string))
		}
		if filter.Type == "path-regex" {
			output.Paths = append(output.Paths, filter.Value.(string))
		}
		if filter.Type == "mime-type" {
			output.MimeTypes = append(output.MimeTypes, filter.Value.(string))
		}
		if filter.Type == "property" {
			output.Properties = map[string]string{}
			filterParams := filter.Value.(map[string]interface{})
			key := filterParams["key"].(string)
			value := filterParams["value"].(string)
			output.Properties[key] = value
		}

		if filter.Type == "path-ant-patterns" {
			// The path filters are defined once for repositories, either all, or by name
			// So, we only add the paths once

			pathFilters := filter.Value.(map[string]interface{})

			if len(repos.ExcludePatterns) == 0 && pathFilters["ExcludePatterns"] != nil {
				for _, path := range pathFilters["ExcludePatterns"].([]interface{}) {
					repos.ExcludePatterns = append(repos.ExcludePatterns, path.(string))
				}
			}
			if len(repos.IncludePatterns) == 0 && pathFilters["IncludePatterns"] != nil {
				for _, path := range pathFilters["IncludePatterns"].([]interface{}) {
					repos.IncludePatterns = append(repos.IncludePatterns, path.(string))
				}
			}
		}
	}

	// Sorting so that outputs are consistent
	// Not sure if this is the best solution.
	sort.Strings(output.PackageTypes)
	sort.Strings(output.Names)
	sort.Strings(output.Paths)
	sort.Strings(output.MimeTypes)
}

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func NewXrayWatchParams() XrayWatchParams {
	return XrayWatchParams{}
}

type XrayWatchParams struct {
	Name        string
	Description string
	Active      bool

	Repositories XrayWatchRepositoriesParams

	Builds   XrayWatchBuildsParams
	Policies []XrayWatchPolicy
}

type XrayWatchBody struct {
	GeneralData      XrayWatchGeneralParams    `json:"general_data"`
	ProjectResources XrayWatchProjectResources `json:"project_resources,omitempty"`
	AssignedPolicies []XrayWatchPolicy         `json:"assigned_policies,omitempty"`
}

type WatchBuildType string
type WatchRepositoriesType string

type XrayWatchGeneralParams struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"` // Must be empty on update.
	Description string `json:"description"`
	Active      bool   `json:"active"`
}

type XrayWatchRepositoriesParams struct {
	Type         WatchRepositoriesType
	All          XrayWatchAll
	Repositories map[string]XrayWatchRepository
	XrayWatchPathFilters
}

type XrayWatchAll struct {
	Filters XrayWatchFilters
}

type XrayWatchFilters struct {
	PackageTypes []string
	Names        []string
	Paths        []string
	MimeTypes    []string
	Properties   map[string]string
}

type XrayWatchBuildsParams struct {
	Type    WatchBuildType
	All     XrayWatchBuildsAllParams
	ByNames map[string]XrayWatchBuildsByNameParams
}

type XrayWatchBuildsAllParams struct {
	Bin_Mgr_ID string `json:"bin_mgr_id"`
	XrayWatchPathFilters
}

type XrayWatchBuildsByNameParams struct {
	Name       string
	Bin_Mgr_ID string
}

type XrayWatchFilterString struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

type XrayWatchFilterProperty struct {
	Type  string                       `json:"type"`
	Value XrayWatchFilterPropertyValue `json:"value"`
}

type XrayWatchFilterPropertyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type XrayWatchProjectResources struct {
	Resources []XrayWatchProjectResourcesElement `json:"resources"`
}

type XrayWatchProjectResourcesElement struct {
	Name          string                  `json:"name,omitempty"`
	Bin_Mgr_ID    string                  `json:"bin_mgr_id,omitempty"`
	Type          string                  `json:"type"`
	StringFilters []XrayWatchFilterString `json:"filters,omitempty"`
}

type XrayWatchRepository struct {
	Name          string                  `json:"name"`
	StringFilters []XrayWatchFilterString `json:"filters"`
	Bin_Mgr_ID    string                  `json:"bin_mgr_id`
	Filters       XrayWatchFilters
}

type XrayWatchPathFilters struct {
	ExcludePatterns []string `json:"ExcludePatterns"`
	IncludePatterns []string `json:"IncludePatterns"`
}

func NewXrayWatchRepository(name string, bin_mgr_id string) XrayWatchRepository {
	return XrayWatchRepository{
		Name:       name,
		Bin_Mgr_ID: bin_mgr_id,
	}
}

type XrayWatchPolicy struct {
	Name string `json:"name"`
	Type string `json:"type"`
}
