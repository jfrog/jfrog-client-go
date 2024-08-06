package utils

import (
	"fmt"
	"time"
)

type IgnoreRuleParams struct {
	Notes         string        `json:"notes"`
	ExpiresAt     time.Time     `json:"expires_at,omitempty"`
	IgnoreFilters IgnoreFilters `json:"ignore_filters"`
}

type IgnoreRuleBody struct {
	Id        string    `json:"id,omitempty"`
	Author    string    `json:"author,omitempty"`
	Created   time.Time `json:"created,omitempty"`
	IsExpired bool      `json:"is_expired,omitempty"`
	IgnoreRuleParams
}

type IgnoreFilters struct {
	Vulnerabilities  []string                      `json:"vulnerabilities,omitempty"`
	Licenses         []string                      `json:"licenses,omitempty"`
	CVEs             []string                      `json:"cves,omitempty"`
	Policies         []string                      `json:"policies,omitempty"`
	Watches          []string                      `json:"watches,omitempty"`
	DockerLayers     []string                      `json:"docker-layers,omitempty"`
	OperationalRisks []string                      `json:"operational_risk,omitempty"`
	Exposures        []ExposuresFilterName         `json:"exposures,omitempty"`
	ReleaseBundles   []IgnoreFilterNameVersion     `json:"release-bundles,omitempty"`
	Builds           []IgnoreFilterNameVersion     `json:"builds,omitempty"`
	Components       []IgnoreFilterNameVersion     `json:"components,omitempty"`
	Artifacts        []IgnoreFilterNameVersionPath `json:"artifacts,omitempty"`
}

type IgnoreFilterNameVersion struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
}

type IgnoreFilterNameVersionPath struct {
	IgnoreFilterNameVersion
	Path string `json:"path,omitempty"`
}

type ExposuresFilterName struct {
	Catagories []ExposuresCatagories `json:"catagories,omitempty"`
	Scanners   []string              `json:"scanners,omitempty"`
	FilePath   []string              `json:"file_path,omitempty"`
}

type ExposuresCatagories struct {
	Secrets      bool `json:"secrets,omitempty"`
	Services     bool `json:"services,omitempty"`
	Applications bool `json:"applications,omitempty"`
	Iac          bool `json:"iac,omitempty"`
}

func NewIgnoreRuleParams() IgnoreRuleParams {
	return IgnoreRuleParams{}
}

func CreateIgnoreRuleBody(ignoreRuleParams IgnoreRuleParams) IgnoreRuleBody {
	return IgnoreRuleBody{
		IgnoreRuleParams: ignoreRuleParams,
	}
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

// IgnoreRuleResponse struct representing the entire JSON
type IgnoreRuleResponse struct {
	Data       []IgnoreRuleBody `json:"data"`
	TotalCount int              `json:"total_count"`
}

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
