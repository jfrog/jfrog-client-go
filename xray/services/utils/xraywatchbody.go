package utils

import (
	"errors"
	"sort"

	"github.com/jfrog/jfrog-client-go/utils/errorutils"
)

const (
	// WatchBuildAll is the option where all builds are watched
	WatchBuildAll WatchBuildType = "all"
	// WatchBuildByName is the option where builds are selected by name to be watched
	WatchBuildByName WatchBuildType = "byname"

	// WatchRepositoriesAll is the option where all repositories are watched
	WatchRepositoriesAll WatchRepositoriesType = "all"
	// WatchRepositoriesByName is the option where repositories are selected by name to be watched
	WatchRepositoriesByName WatchRepositoriesType = "byname"
)

// WatchBuildType defines the type of filter for a builds on a watch
type WatchBuildType string

// WatchRepositoriesType defines the type of filter for a repositories on a watch
type WatchRepositoriesType string

// NewXrayWatchParams creates a new struct to configure an Xray watch
func NewXrayWatchParams() XrayWatchParams {
	return XrayWatchParams{}
}

// XrayWatchParams defines all the properties to create an Xray watch
type XrayWatchParams struct {
	Name        string
	Description string
	Active      bool

	Repositories XrayWatchRepositoriesParams

	Builds   XrayWatchBuildsParams
	Policies []XrayAssignedPolicy
}

// XrayWatchRepositoriesParams is a struct that stores the repository configuration for watch
type XrayWatchRepositoriesParams struct {
	Type         WatchRepositoriesType
	All          XrayWatchRepositoryAll
	Repositories map[string]XrayWatchRepository
	XrayWatchPathFilters
}

// XrayWatchRepositoryAll is used to define the parameters when a watch uses all repositories
type XrayWatchRepositoryAll struct {
	Filters xrayWatchFilters
}

// XrayWatchRepository is used to define a specific repository in a watch
type XrayWatchRepository struct {
	Name     string
	BinMgrID string
	Filters  xrayWatchFilters
}

// XrayWatchBuildsParams is a struct that stores the build configuration for watch
type XrayWatchBuildsParams struct {
	Type    WatchBuildType
	All     XrayWatchBuildsAllParams
	ByNames map[string]XrayWatchBuildsByNameParams
}

// XrayWatchBuildsAllParams is used to define the parameters when a watch uses all builds
type XrayWatchBuildsAllParams struct {
	BinMgrID string
	XrayWatchPathFilters
}

// XrayWatchBuildsByNameParams is used to define a specific build in a watch
type XrayWatchBuildsByNameParams struct {
	Name     string
	BinMgrID string
}

// XrayWatchPathFilters is used to define path filters on a repository or a build in a watch
type XrayWatchPathFilters struct {
	ExcludePatterns []string `json:"ExcludePatterns"`
	IncludePatterns []string `json:"IncludePatterns"`
}

// XrayAssignedPolicy struct is used to define a policy associated with a watch
type XrayAssignedPolicy struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// XrayWatchBody is the top level payload to be sent to xray
type XrayWatchBody struct {
	GeneralData      xrayWatchGeneralParams    `json:"general_data"`
	ProjectResources xrayWatchProjectResources `json:"project_resources,omitempty"`
	AssignedPolicies []XrayAssignedPolicy      `json:"assigned_policies,omitempty"`
}

// These structs are internal

type xrayWatchGeneralParams struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"` // Name must be empty on update.
	Description string `json:"description"`
	Active      bool   `json:"active"`
}

type xrayWatchProjectResources struct {
	Resources []xrayWatchProjectResourcesElement `json:"resources"`
}

type xrayWatchProjectResourcesElement struct {
	Name     string            `json:"name,omitempty"`
	BinMgrID string            `json:"bin_mgr_id,oitempty"`
	Type     string            `json:"type"`
	Filters  []xrayWatchFilter `json:"filters,omitempty"`
}

type xrayWatchFilters struct {
	PackageTypes []string
	Names        []string
	Paths        []string
	MimeTypes    []string
	Properties   map[string]string
}

type xrayWatchFilter struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

type xrayWatchFilterPropertyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// CreateBody creates a payload to configure a Watch in Xray
// This can configure repositories and builds
// However, bundles are not supported.
func CreateBody(params XrayWatchParams) (*XrayWatchBody, error) {
	payloadBody := XrayWatchBody{
		GeneralData: xrayWatchGeneralParams{
			Name:        params.Name,
			Description: params.Description,
			Active:      params.Active,
		},
		ProjectResources: xrayWatchProjectResources{
			Resources: []xrayWatchProjectResourcesElement{},
		},
		AssignedPolicies: params.Policies,
	}

	err := configureRepositories(&payloadBody, params)
	if err != nil {
		return nil, err
	}

	err = configureBuilds(&payloadBody, params)
	if err != nil {
		return nil, err
	}

	return &payloadBody, nil
}

func configureRepositories(payloadBody *XrayWatchBody, params XrayWatchParams) error {
	switch params.Repositories.Type {

	case WatchRepositoriesAll:
		allFilters := xrayWatchProjectResourcesElement{
			Type:    "all-repos",
			Filters: []xrayWatchFilter{},
		}

		allFilters.Filters = append(allFilters.Filters, createFilters(params.Repositories.All.Filters, params.Repositories)...)

		payloadBody.ProjectResources.Resources = append(payloadBody.ProjectResources.Resources, allFilters)

	case WatchRepositoriesByName:
		for _, repository := range params.Repositories.Repositories {
			repo := xrayWatchProjectResourcesElement{
				Type:     "repository",
				Name:     repository.Name,
				BinMgrID: repository.BinMgrID,
				Filters:  []xrayWatchFilter{},
			}

			repo.Filters = append(repo.Filters, createFilters(repository.Filters, params.Repositories)...)

			payloadBody.ProjectResources.Resources = append(payloadBody.ProjectResources.Resources, repo)
		}
	case "":
		// Empty is fine
	default:
		return errorutils.CheckError(errors.New("Invalid Repository Type. Must be " + string(WatchRepositoriesAll) + " or " + string(WatchRepositoriesByName)))
	}

	return nil
}

func createFilters(filters xrayWatchFilters, repo XrayWatchRepositoriesParams) []xrayWatchFilter {
	result := []xrayWatchFilter{}

	for _, packageType := range filters.PackageTypes {
		filter := xrayWatchFilter{
			Type:  "package-type",
			Value: packageType,
		}
		result = append(result, filter)
	}

	for _, name := range filters.Names {
		filter := xrayWatchFilter{
			Type:  "regex",
			Value: name,
		}
		result = append(result, filter)
	}

	for _, path := range filters.Paths {
		filter := xrayWatchFilter{
			Type:  "path-regex",
			Value: path,
		}
		result = append(result, filter)
	}

	for _, mimeType := range filters.MimeTypes {
		filter := xrayWatchFilter{
			Type:  "mime-type",
			Value: mimeType,
		}
		result = append(result, filter)
	}

	for key, value := range filters.Properties {
		filter := xrayWatchFilter{
			Type: "property",
			Value: xrayWatchFilterPropertyValue{
				Key:   key,
				Value: value,
			},
		}
		result = append(result, filter)
	}

	if repo.ExcludePatterns != nil || repo.IncludePatterns != nil {
		filter := xrayWatchFilter{
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

func configureBuilds(payloadBody *XrayWatchBody, params XrayWatchParams) error {
	switch params.Builds.Type {
	case WatchBuildAll:
		allBuilds := xrayWatchProjectResourcesElement{
			Name:     "All Builds",
			Type:     "all-builds",
			BinMgrID: params.Builds.All.BinMgrID,
			Filters:  []xrayWatchFilter{},
		}

		if params.Builds.All.ExcludePatterns != nil || params.Builds.All.IncludePatterns != nil {
			filters := []xrayWatchFilter{{
				Type: "ant-patterns",
				Value: XrayWatchPathFilters{
					ExcludePatterns: params.Builds.All.ExcludePatterns,
					IncludePatterns: params.Builds.All.IncludePatterns,
				}},
			}
			allBuilds.Filters = filters
		}

		payloadBody.ProjectResources.Resources = append(payloadBody.ProjectResources.Resources, allBuilds)

	case WatchBuildByName:
		for _, byName := range params.Builds.ByNames {
			build := xrayWatchProjectResourcesElement{
				Type:     "build",
				Name:     byName.Name,
				BinMgrID: byName.BinMgrID,
			}

			payloadBody.ProjectResources.Resources = append(payloadBody.ProjectResources.Resources, build)
		}
	case "":
		// Empty is fine
	default:
		return errorutils.CheckError(errors.New("Invalid Build Type. Must be " + string(WatchBuildAll) + " or " + string(WatchBuildByName)))
	}

	return nil
}

// UnpackWatchBody unpacks a payload response from Xray.
// It transforms the data into the params object so that a consumer can interact with a watch in a consistent way.
func UnpackWatchBody(watch *XrayWatchParams, body *XrayWatchBody) {
	for _, resource := range body.ProjectResources.Resources {
		switch resource.Type {

		case "all-repos":
			watch.Repositories.Type = WatchRepositoriesAll
			unpackFilters(resource.Filters, &watch.Repositories.All.Filters, &watch.Repositories)

		case "repository":
			watch.Repositories.Type = WatchRepositoriesByName
			repository := XrayWatchRepository{
				Name:     resource.Name,
				BinMgrID: resource.BinMgrID,
			}
			unpackFilters(resource.Filters, &repository.Filters, &watch.Repositories)
			watch.Repositories.Repositories[repository.Name] = repository

		case "all-builds":
			watch.Builds.Type = WatchBuildAll
			watch.Builds.All.BinMgrID = resource.BinMgrID

			for _, filter := range resource.Filters {
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

func unpackFilters(filters []xrayWatchFilter, output *xrayWatchFilters, repos *XrayWatchRepositoriesParams) {

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

// NewXrayWatchRepository creates a new repository struct to configure an Xray Watch
func NewXrayWatchRepository(name string, binMgrID string) XrayWatchRepository {
	return XrayWatchRepository{
		Name:     name,
		BinMgrID: binMgrID,
	}
}
