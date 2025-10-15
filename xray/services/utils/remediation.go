package utils

type CveRemediationResponse map[string][]Option

type PackageVersionKey struct {
	Type      string    `json:"type"`
	Name      string    `json:"name"`
	Namespace string    `json:"namespace"`
	Version   string    `json:"version"`
	Ecosystem Ecosystem `json:"ecosystem"`
}

type Ecosystem string

const (
	// Related ecosystems to packages
	GenericEcosystem Ecosystem = "generic"
	DebianEcosystem  Ecosystem = "debian"
	UbuntuEcosystem  Ecosystem = "ubuntu"
)

type StepType string

const (
	// Step type to indicate no action is required
	None StepType = "None"
	// Step type to indicate how to apply a Fix version remediation
	FixVersion StepType = "FixVersion"
	// Step type to indicate No fix version is available
	NoFixVersion StepType = "NoFixVersion"
	// Step type to indicate the package is not found in catalog
	PackageNotFound StepType = "PackageNotFound"
)

type OptionType string

const (
	// Basic (Fix version on actual component)
	InLock OptionType = "InLock"
	// Direct dependency (Fix version on direct dependency to fix transitive/direct dependency)
	DirectDependency OptionType = "DirectDependency"
)

type Option struct {
	Type        OptionType    `json:"type"`
	Description string        `json:"description"`
	Steps       []OptionStep  `json:"steps"`
	Snippet     []CodeSnippet `json:"snippet"`
}

type OptionStep struct {
	PkgVersion PackageVersionKey
	UpgradeTo  PackageVersionKey
	StepType   StepType
	Badges     []string `json:"badges,omitempty"`
	Party      string   `json:"party,omitempty"`
}

type CodeSnippet struct {
	PackageManager string `json:"package_manager,omitempty"`
	FileName       string `json:"file_name,omitempty"`
	Code           string `json:"code,omitempty"`
}
