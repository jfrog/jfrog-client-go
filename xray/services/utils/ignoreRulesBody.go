package utils

import (
	"fmt"
	"time"
)

func (p *IgnoreRulesGetAllParams) GetParamMap() map[string]string {
	params := make(map[string]string)
	if p == nil {
		return params
	}
	if p.Vulnerability != "" {
		params["vulnerability"] = p.Vulnerability
	}
	if p.License != "" {
		params["license"] = p.License
	}
	if p.Policy != "" {
		params["policy"] = p.Policy
	}
	if p.Watch != "" {
		params["watch"] = p.Watch
	}
	if p.ComponentName != "" {
		params["component_name"] = p.ComponentName
	}
	if p.ComponentVersion != "" {
		params["component_version"] = p.ComponentVersion
	}
	if p.ArtifactName != "" {
		params["artifact_name"] = p.ArtifactName
	}
	if p.ArtifactVersion != "" {
		params["artifact_version"] = p.ArtifactVersion
	}
	if p.BuildName != "" {
		params["build_name"] = p.BuildName
	}
	if p.BuildVersion != "" {
		params["build_version"] = p.BuildVersion
	}
	if p.ReleaseBundleName != "" {
		params["release_bundle_name"] = p.ReleaseBundleName
	}
	if p.ReleaseBundleVersion != "" {
		params["release_bundle_version"] = p.ReleaseBundleVersion
	}
	if p.DockerLayer != "" {
		params["docker_layer"] = p.DockerLayer
	}
	if p.OrderBy != "" {
		params["order_by"] = p.OrderBy
	}
	if p.Direction != "" {
		params["direction"] = p.Direction
	}
	if p.PageNum != 0 {
		params["page_num"] = fmt.Sprintf("%d", p.PageNum)
	}
	if p.NumOfRows != 0 {
		params["num_of_rows"] = fmt.Sprintf("%d", p.NumOfRows)
	}
	if !p.ExpiresBefore.IsZero() {
		params["expires_before"] = p.ExpiresBefore.UTC().Format(time.RFC3339)
	}
	if !p.ExpiresAfter.IsZero() {
		params["expires_after"] = p.ExpiresAfter.UTC().Format(time.RFC3339)
	}
	if p.ProjectKey != "" {
		params["project_key"] = p.ProjectKey
	}

	return params
}

// IgnoreRuleResponse struct representing the entire JSON
type IgnoreRuleResponse struct {
	Data       []IgnoreRuleDetail `json:"data"`
	TotalCount int                `json:"total_count"`
}

// IgnoreRule struct representing an Ignore Rule
type IgnoreRule struct {
	Notes     string        `json:"notes"`
	ExpiresAt *time.Time    `json:"expires_at,omitempty"`
	Filters   IgnoreFilters `json:"ignore_filters"`
}

// IgnoreRuleDetail struct representing an Ignore Rule as returned by the API
type IgnoreRuleDetail struct {
	IgnoreRule
	ID        string    `json:"id"`
	Author    string    `json:"author"`
	Created   time.Time `json:"created"`
	IsExpired bool      `json:"is_expired"`
}

// IgnoreFilters struct representing the "ignore_filters" object
type IgnoreFilters struct {
	ReleaseBundles  []NameVersion        `json:"release_bundles,omitempty"`
	Builds          []NameVersion        `json:"builds,omitempty"`
	Components      []NameVersion        `json:"components,omitempty"`
	Artifacts       []ArtifactDescriptor `json:"artifacts,omitempty"`
	Policies        []string             `json:"policies,omitempty"`
	DockerLayers    []string             `json:"docker_layers,omitempty"`
	Vulnerabilities []string             `json:"vulnerabilities,omitempty"`
	Licenses        []string             `json:"licenses,omitempty"`
	CVEs            []string             `json:"cves,omitempty"`
	Watches         []string             `json:"watches,omitempty"`
	OperationalRisk []string             `json:"operational_risk,omitempty"`
}

// NameVersion struct representing items with a Name / Version combo
type NameVersion struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// ArtifactDescriptor struct representing each item in the "artifacts" array
type ArtifactDescriptor struct {
	NameVersion
	Path string `json:"path"`
}

type IgnoreRulesGetAllParams struct {
	Vulnerability        string    `json:"vulnerability"`
	License              string    `json:"license"`
	Policy               string    `json:"policy"`
	Watch                string    `json:"watch"`
	ComponentName        string    `json:"component_name"`
	ComponentVersion     string    `json:"component_version"`
	ArtifactName         string    `json:"artifact_name"`
	ArtifactVersion      string    `json:"artifact_version"`
	BuildName            string    `json:"build_name"`
	BuildVersion         string    `json:"build_version"`
	ReleaseBundleName    string    `json:"release_bundle_name"`
	ReleaseBundleVersion string    `json:"release_bundle_version"`
	DockerLayer          string    `json:"docker_layer"`
	ExpiresBefore        time.Time `json:"expires_before"`
	ExpiresAfter         time.Time `json:"expires_after"`
	ProjectKey           string    `json:"project_key"`
	OrderBy              string    `json:"order_by"`
	Direction            string    `json:"direction"`
	PageNum              int       `json:"page_num"`
	NumOfRows            int       `json:"num_of_rows"`
}
