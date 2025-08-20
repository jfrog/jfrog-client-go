package utils

import "math"

const (
	SecurityViolation        ViolationType = "Security"
	LicenseViolation         ViolationType = "License"
	OperationalRiskViolation ViolationType = "Operational_Risk"
)

type ViolationType string

type ViolationsRequest struct {
	Filters    *ViolationsFilters `json:"filters,omitempty"`
	Pagination *PaginationOptions `json:"pagination,omitempty"`
}

type PaginationOptions struct {
	// Valid values:  created, summary, severity, type, watcher_name, issue_id (Default: created)
	OrderBy string `json:"order_by,omitempty"`
	// Default if not provided: 25
	Limit int `json:"limit,omitempty"`
	// Default if not provided: 1
	Offset int `json:"offset,omitempty"`
	// Valid values: asc, desc (Default: For ordering by Severity: desc; for all other fields: asc)
	Direction string `json:"direction,omitempty"`
}

type ViolationsFilters struct {
	// Filtering the results for those included in the requested string in the description property
	NameContains string `json:"name_contains,omitempty"`
	// Add additional violation detail properties to the response (role name, policy name, description, remediation, and more)
	IncludeDetails bool `json:"include_details,omitempty"`
	// Filtering by the response for specific violation type. Valid values: Security, License, Operational_Risk
	Type ViolationType `json:"violation_type,omitempty"`
	// Filtering the results for those generated from the selected watch. Default: Any watch.
	WatchName string `json:"watch_name,omitempty"`
	// Filtering the results for those that their severity is equal or higher than min_severity.
	// Valid values: Critical, High, Medium, Low, Information, Unknown (Note:  the values are listed in descending severity order)
	MinSeverity Severity `json:"min_severity,omitempty"`
	// Filter for violations created as of this time. Valid value:  A timestamp in RFC 3339 format
	CreatedFrom string `json:"created_from,omitempty"`
	// Filter for violations created up to this time. Valid value:  A timestamp in RFC 3339 format
	CreatedUntil string `json:"created_until,omitempty"`
	// Filter for violations resulting from the requested Issue ID.
	// Valid values: strings representing the issue ID e.g: XRAY-94620, EXP-1552-00002, GPL-3.0, b1670bb2d3438da6213ed386577fd755bc, b8fdf85cab594e6a3717b4f182b07b
	IssueId string `json:"issue_id,omitempty"`
	// Filter for violations resulting from the requested CVE. Valid values:  a CVE standard identifier format <CVE-YYYY-NNNNNN>
	CveId string `json:"cve_id,omitempty"`
	// Filter for violations found in specific resources.
	Resources ViolationResourceFilters `json:"resources,omitempty"`
}

type ViolationResourceFilters struct {
	Artifacts        []ArtifactResourceFilter        `json:"artifacts,omitempty"`
	Builds           []BuildResourceFilter           `json:"builds,omitempty"`
	ReleaseBundles   []ReleaseBundleResourceFilter   `json:"release_bundles,omitempty"`
	ReleaseBundlesV2 []ReleaseBundleV2ResourceFilter `json:"release_bundles_v2,omitempty"`
	GitRepositories  []GitRepositoriesResourceFilter `json:"git_repositories,omitempty"`
}

type ArtifactResourceFilter struct {
	Repository string `json:"repo"`
	Path       string `json:"path"`
}

type BuildResourceFilter struct {
	Name    string `json:"name"`
	Number  string `json:"number"`
	Project string `json:"project,omitempty"`
}

type ReleaseBundleResourceFilter struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type ReleaseBundleV2ResourceFilter struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Project string `json:"project,omitempty"`
}

type GitRepositoriesResourceFilter struct {
	Name string `json:"name"`
}

func NewViolationsRequest() ViolationsRequest {
	return ViolationsRequest{
		Filters: &ViolationsFilters{},
		Pagination: &PaginationOptions{
			OrderBy:   "created",
			Limit:     math.MaxInt,
			Direction: "asc",
		},
	}
}

func (vr ViolationsRequest) FilterByNameContains(name string) ViolationsRequest {
	vr.Filters.NameContains = name
	return vr
}

func (vr ViolationsRequest) FilterByType(violationType ViolationType) ViolationsRequest {
	vr.Filters.Type = violationType
	return vr
}

func (vr ViolationsRequest) FilterByWatchName(watchName string) ViolationsRequest {
	vr.Filters.WatchName = watchName
	return vr
}

func (vr ViolationsRequest) FilterByMinSeverity(severity Severity) ViolationsRequest {
	vr.Filters.MinSeverity = severity
	return vr
}

func (vr ViolationsRequest) FilterByCreatedFrom(createdFrom string) ViolationsRequest {
	vr.Filters.CreatedFrom = createdFrom
	return vr
}

func (vr ViolationsRequest) FilterByCreatedUntil(createdUntil string) ViolationsRequest {
	vr.Filters.CreatedUntil = createdUntil
	return vr
}

func (vr ViolationsRequest) FilterByIssueId(issueId string) ViolationsRequest {
	vr.Filters.IssueId = issueId
	return vr
}

func (vr ViolationsRequest) FilterByCveId(cveId string) ViolationsRequest {
	vr.Filters.CveId = cveId
	return vr
}

func (vr ViolationsRequest) FilterByArtifacts(artifacts ...ArtifactResourceFilter) ViolationsRequest {
	if vr.Filters.Resources.Artifacts == nil {
		vr.Filters.Resources.Artifacts = []ArtifactResourceFilter{}
	}
	vr.Filters.Resources.Artifacts = append(vr.Filters.Resources.Artifacts, artifacts...)
	return vr
}

func (vr ViolationsRequest) FilterByBuilds(builds ...BuildResourceFilter) ViolationsRequest {
	if vr.Filters.Resources.Builds == nil {
		vr.Filters.Resources.Builds = []BuildResourceFilter{}
	}
	vr.Filters.Resources.Builds = append(vr.Filters.Resources.Builds, builds...)
	return vr
}

func (vr ViolationsRequest) FilterByReleaseBundles(releaseBundles ...ReleaseBundleResourceFilter) ViolationsRequest {
	if vr.Filters.Resources.ReleaseBundles == nil {
		vr.Filters.Resources.ReleaseBundles = []ReleaseBundleResourceFilter{}
	}
	vr.Filters.Resources.ReleaseBundles = append(vr.Filters.Resources.ReleaseBundles, releaseBundles...)
	return vr
}

func (vr ViolationsRequest) FilterByReleaseBundleV2(releaseBundles ...ReleaseBundleV2ResourceFilter) ViolationsRequest {
	if vr.Filters.Resources.ReleaseBundlesV2 == nil {
		vr.Filters.Resources.ReleaseBundlesV2 = []ReleaseBundleV2ResourceFilter{}
	}
	vr.Filters.Resources.ReleaseBundlesV2 = append(vr.Filters.Resources.ReleaseBundlesV2, releaseBundles...)
	return vr
}

func (vr ViolationsRequest) FilterByGitRepositories(gitRepositories ...GitRepositoriesResourceFilter) ViolationsRequest {
	if vr.Filters.Resources.GitRepositories == nil {
		vr.Filters.Resources.GitRepositories = []GitRepositoriesResourceFilter{}
	}
	vr.Filters.Resources.GitRepositories = append(vr.Filters.Resources.GitRepositories, gitRepositories...)
	return vr
}

func (vr ViolationsRequest) IncludeDetails(include bool) ViolationsRequest {
	vr.Filters.IncludeDetails = include
	return vr
}

func (vr ViolationsRequest) SetPaginationOptions(orderBy string, limit, offset int, direction string) ViolationsRequest {
	vr.Pagination.OrderBy = orderBy
	vr.Pagination.Limit = limit
	vr.Pagination.Offset = offset
	vr.Pagination.Direction = direction
	return vr
}
