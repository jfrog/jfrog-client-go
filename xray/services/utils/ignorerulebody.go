package utils

import "time"

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
