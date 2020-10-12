package utils

import (
	"errors"
	"sort"
)

const (
	WatchBuildAll    WatchBuildType = "all"
	WatchBuildByName WatchBuildType = "byname"

	WatchRepositoriesAll    WatchRepositoriesType = "all"
	WatchRepositoriesByName WatchRepositoriesType = "byname"
)

func NewXrayWatchParams() XrayWatchParams {
	return XrayWatchParams{}
}

type XrayWatchParams struct {
	Name        string
	Description string
	Active      bool

	Repositories XrayWatchRepositoriesParams

	Builds   XrayWatchBuildsParams
	Policies []XrayPolicy
}

type XrayWatchBody struct {
	GeneralData      XrayWatchGeneralParams    `json:"general_data"`
	ProjectResources XrayWatchProjectResources `json:"project_resources,omitempty"`
	AssignedPolicies []XrayPolicy              `json:"assigned_policies,omitempty"`
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
	BinMgrID string `json:"bin_mgr_id"`
	XrayWatchPathFilters
}

type XrayWatchBuildsByNameParams struct {
	Name     string
	BinMgrID string
}

type XrayWatchFilter struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

type XrayWatchFilterPropertyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type XrayWatchProjectResources struct {
	Resources []XrayWatchProjectResourcesElement `json:"resources"`
}

type XrayWatchProjectResourcesElement struct {
	Name          string            `json:"name,omitempty"`
	BinMgrID      string            `json:"bin_mgr_id,oitempty"`
	Type          string            `json:"type"`
	StringFilters []XrayWatchFilter `json:"filters,omitempty"`
}

type XrayWatchRepository struct {
	Name          string            `json:"name"`
	StringFilters []XrayWatchFilter `json:"filters"`
	BinMgrID      string            `json:"bin_mgr_id"`
	Filters       XrayWatchFilters
}

type XrayWatchPathFilters struct {
	ExcludePatterns []string `json:"ExcludePatterns"`
	IncludePatterns []string `json:"IncludePatterns"`
}

type XrayPolicy struct {
	Name string `json:"name"`
	Type string `json:"type"`
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
	switch params.Repositories.Type {

	case WatchRepositoriesAll:
		allFilters := XrayWatchProjectResourcesElement{
			Type:          "all-repos",
			StringFilters: []XrayWatchFilter{},
		}

		allFilters.StringFilters = append(allFilters.StringFilters, CreateFilters(params.Repositories.All.Filters, params.Repositories)...)

		payloadBody.ProjectResources.Resources = append(payloadBody.ProjectResources.Resources, allFilters)

	case WatchRepositoriesByName:
		for _, repository := range params.Repositories.Repositories {
			repo := XrayWatchProjectResourcesElement{
				Type:          "repository",
				Name:          repository.Name,
				BinMgrID:      repository.BinMgrID,
				StringFilters: repository.StringFilters,
			}
			if repo.StringFilters == nil {
				repo.StringFilters = []XrayWatchFilter{}
			}
			repo.StringFilters = append(repo.StringFilters, CreateFilters(repository.Filters, params.Repositories)...)

			payloadBody.ProjectResources.Resources = append(payloadBody.ProjectResources.Resources, repo)
		}
	case "":
		// empty is fine
	default:
		return errors.New("Invalid Repository Type. Must be " + string(WatchRepositoriesAll) + " or " + string(WatchRepositoriesByName))
	}

	return nil
}

func CreateFilters(filters XrayWatchFilters, repo XrayWatchRepositoriesParams) []XrayWatchFilter {
	result := []XrayWatchFilter{}

	for _, packageType := range filters.PackageTypes {
		filter := XrayWatchFilter{
			Type:  "package-type",
			Value: packageType,
		}
		result = append(result, filter)
	}

	for _, name := range filters.Names {
		filter := XrayWatchFilter{
			Type:  "regex",
			Value: name,
		}
		result = append(result, filter)
	}

	for _, path := range filters.Paths {
		filter := XrayWatchFilter{
			Type:  "path-regex",
			Value: path,
		}
		result = append(result, filter)
	}

	for _, mimeType := range filters.MimeTypes {
		filter := XrayWatchFilter{
			Type:  "mime-type",
			Value: mimeType,
		}
		result = append(result, filter)
	}

	for key, value := range filters.Properties {
		filter := XrayWatchFilter{
			Type: "property",
			Value: XrayWatchFilterPropertyValue{
				Key:   key,
				Value: value,
			},
		}
		result = append(result, filter)
	}

	if repo.ExcludePatterns != nil || repo.IncludePatterns != nil {
		filter := XrayWatchFilter{
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
	switch params.Builds.Type {
	case WatchBuildAll:
		allBuilds := XrayWatchProjectResourcesElement{
			Name:          "All Builds",
			Type:          "all-builds",
			BinMgrID:      params.Builds.All.BinMgrID,
			StringFilters: []XrayWatchFilter{},
		}

		if params.Builds.All.ExcludePatterns != nil || params.Builds.All.IncludePatterns != nil {
			filters := []XrayWatchFilter{{
				Type: "ant-patterns",
				Value: XrayWatchPathFilters{
					ExcludePatterns: params.Builds.All.ExcludePatterns,
					IncludePatterns: params.Builds.All.IncludePatterns,
				}},
			}
			allBuilds.StringFilters = filters
		}

		payloadBody.ProjectResources.Resources = append(payloadBody.ProjectResources.Resources, allBuilds)

	case WatchBuildByName:
		for _, byName := range params.Builds.ByNames {
			build := XrayWatchProjectResourcesElement{
				Type:     "build",
				Name:     byName.Name,
				BinMgrID: byName.BinMgrID,
			}

			payloadBody.ProjectResources.Resources = append(payloadBody.ProjectResources.Resources, build)
		}
	case "":
		// empty is fine
	default:
		return errors.New("Invalid Build Type. Must be " + string(WatchBuildAll) + " or " + string(WatchBuildByName))
	}

	return nil
}

func ConfigureBundles(payloadBody *XrayWatchBody, params XrayWatchParams) error {
	// placeholder method to support bundles in a future release
	return nil
}

func UnpackWatchBody(watch *XrayWatchParams, body *XrayWatchBody) {
	for _, resource := range body.ProjectResources.Resources {
		switch resource.Type {

		case "all-repos":
			watch.Repositories.Type = WatchRepositoriesAll
			unpackFilters(resource.StringFilters, &watch.Repositories.All.Filters, &watch.Repositories)

		case "repository":
			watch.Repositories.Type = WatchRepositoriesByName
			repository := XrayWatchRepository{
				Name:     resource.Name,
				BinMgrID: resource.BinMgrID,
			}
			unpackFilters(resource.StringFilters, &repository.Filters, &watch.Repositories)
			watch.Repositories.Repositories[repository.Name] = repository

		case "all-builds":
			watch.Builds.Type = WatchBuildAll
			watch.Builds.All.BinMgrID = resource.BinMgrID

			for _, filter := range resource.StringFilters {
				if filter.Type == "ant-patterns" {
					pathFilters := filter.Value.(map[string]interface{})

					if value, ok := pathFilters["ExcludePatterns"]; ok {
						for _, path := range value.([]interface{}) {
							watch.Builds.All.ExcludePatterns = append(watch.Builds.All.ExcludePatterns, path.(string))
						}
					}
					if value, ok := pathFilters["IncludePatterns"]; ok {
						for _, path := range value.([]interface{}) {
							watch.Builds.All.IncludePatterns = append(watch.Builds.All.IncludePatterns, path.(string))
						}
					}
				}
			}

		case "build":
			watch.Builds.Type = WatchBuildByName
			watch.Builds.ByNames[resource.Name] = XrayWatchBuildsByNameParams{
				Name:     resource.Name,
				BinMgrID: resource.BinMgrID,
			}
		}
	}

	// Sort all the properties so they are returned in a consistent format
	sort.Strings(watch.Repositories.ExcludePatterns)
	sort.Strings(watch.Repositories.IncludePatterns)
}

func unpackFilters(filters []XrayWatchFilter, output *XrayWatchFilters, repos *XrayWatchRepositoriesParams) {

	for _, filter := range filters {
		switch filter.Type {

		case "package-type":
			output.PackageTypes = append(output.PackageTypes, filter.Value.(string))

		case "regex":
			output.Names = append(output.Names, filter.Value.(string))

		case "path-regex":
			output.Paths = append(output.Paths, filter.Value.(string))

		case "mime-type":
			output.MimeTypes = append(output.MimeTypes, filter.Value.(string))

		case "property":
			output.Properties = map[string]string{}
			filterParams := filter.Value.(map[string]interface{})
			key := filterParams["key"].(string)
			value := filterParams["value"].(string)
			output.Properties[key] = value

		case "path-ant-patterns":
			// The path filters are defined once for repositories, either all, or by name
			// However, in each repository, the data is stored in the filter.
			// So, if we have 5 repositories, the exclude and include patterns will exist in 5 filters
			// When unpacking, we only want to store them once, rather than 5 times.

			pathFilters := filter.Value.(map[string]interface{})

			if len(repos.ExcludePatterns) == 0 {
				if val, ok := pathFilters["ExcludePatterns"]; ok {
					for _, path := range val.([]interface{}) {
						repos.ExcludePatterns = append(repos.ExcludePatterns, path.(string))
					}
				}
			}
			if len(repos.IncludePatterns) == 0 {
				if val, ok := pathFilters["IncludePatterns"]; ok {
					for _, path := range val.([]interface{}) {
						repos.IncludePatterns = append(repos.IncludePatterns, path.(string))
					}
				}
			}
		}
	}

	// Sorting so that outputs are consistent
	sort.Strings(output.PackageTypes)
	sort.Strings(output.Names)
	sort.Strings(output.Paths)
	sort.Strings(output.MimeTypes)
}
