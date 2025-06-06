package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	xscutils "github.com/jfrog/jfrog-client-go/xsc/services/utils"
)

const (
	ConfigProfileMinXscVersion                = "1.11.0"
	ConfigProfileByUrlMinXrayVersion          = "3.110.0"
	ConfigProfileNewSchemaMinXrayVersion      = "3.117.0"
	xscConfigProfileByNameApi                 = "profile"
	xscConfigProfileByUrlApi                  = "profile_repos"
	xscDeprecatedConfigProfileByNameApiSuffix = "api/v1/" + xscConfigProfileByNameApi
	getProfileByUrlBody                       = "{\"repo_url\":\"%s\"}"
)

type ConfigurationProfileService struct {
	client          *jfroghttpclient.JfrogHttpClient
	XscDetails      auth.ServiceDetails
	XrayDetails     auth.ServiceDetails
	ScopeProjectKey string
}

func NewConfigurationProfileService(client *jfroghttpclient.JfrogHttpClient) *ConfigurationProfileService {
	return &ConfigurationProfileService{client: client}
}

type ConfigProfile struct {
	ProfileName    string        `json:"profile_name"`
	GeneralConfig  GeneralConfig `json:"general_config,omitempty"`
	FrogbotConfig  FrogbotConfig `json:"frogbot_config,omitempty"`
	Modules        []Module      `json:"modules"`
	IsDefault      bool          `json:"is_default,omitempty"`
	IsBasicProfile bool
}

type GeneralConfig struct {
	ReleasesRepo           string   `json:"releases_repo,omitempty"`
	AnalyzerManagerVersion string   `json:"analyzer_manager_version,omitempty"`
	ReportAnalytics        bool     `json:"report_analytics,omitempty"`
	ExcludePatterns        []string `json:"exclude_patterns,omitempty"`
	ResultsOutputDir       string   `json:"results_output_dir,omitempty"`
	AllowPartialResults    bool     `json:"allow_partial_results,omitempty"`
}

type FrogbotConfig struct {
	EmailAuthor                         string `json:"email_author,omitempty"`
	AggregateFixes                      bool   `json:"aggregate_fixes,omitempty"`
	AvoidPreviousPrCommentsDeletion     bool   `json:"avoid_previous_pr_comments_deletion,omitempty"`
	AvoidExtraMessages                  bool   `json:"avoid_extra_messages,omitempty"`
	AddSuccessComment                   bool   `json:"add_success_comment,omitempty"`
	BranchNameTemplate                  string `json:"branch_name_template,omitempty"`
	PrTitleTemplate                     string `json:"pr_title_template,omitempty"`
	PrCommentTitle                      string `json:"pr_comment_title,omitempty"`
	CommitMessageTemplate               string `json:"commit_message_template,omitempty"`
	ShowSecretsAsPrComment              bool   `json:"show_secrets_as_pr_comment,omitempty"`
	SkipAutoFix                         bool   `json:"skip_auto_fix,omitempty"`
	IncludeAllRepositoryVulnerabilities bool   `json:"include_all_repository_vulnerabilities,omitempty"`
}

type Module struct {
	ModuleId        int32      `json:"module_id,omitempty"`
	ModuleName      string     `json:"module_name"`
	PathFromRoot    string     `json:"path_from_root"`
	ExcludePatterns []string   `json:"exclude_patterns,omitempty"`
	ScanConfig      ScanConfig `json:"scan_config"`
	DepsRepo        string     `json:"deps_repo,omitempty"`
}

type ScanConfig struct {
	ScaScannerConfig                ScaScannerConfig     `json:"sca_scanner_config,omitempty"`
	ContextualAnalysisScannerConfig CaScannerConfig      `json:"contextual_analysis_scanner_config,omitempty"`
	SastScannerConfig               SastScannerConfig    `json:"sast_scanner_config,omitempty"`
	SecretsScannerConfig            SecretsScannerConfig `json:"secrets_scanner_config,omitempty"`
	IacScannerConfig                IacScannerConfig     `json:"iac_scanner_config,omitempty"`
}

type ScaScannerConfig struct {
	EnableScaScan           bool                    `json:"enable_sca_scan,omitempty"`
	Technology              string                  `json:"technology,omitempty"`
	PackageManagersSettings PackageManagersSettings `json:"package_managers_settings,omitempty"`
	SkipAutoInstall         bool                    `json:"skip_auto_install,omitempty"`
	ExcludePatterns         []string                `json:"exclude_patterns,omitempty"`
}

type PackageManagersSettings struct {
	GradleSettings GradleSettings `json:"gradle_settings,omitempty"`
	MavenSettings  MavenSettings  `json:"maven_settings,omitempty"`
	NpmSettings    NpmSettings    `json:"npm_settings,omitempty"`
	PythonSettings PythonSettings `json:"python_settings,omitempty"`
}

type GradleSettings struct {
	ExcludeTestDeps bool `json:"exclude_test_deps,omitempty"`
	UseWrapper      bool `json:"use_wrapper,omitempty"`
}

type MavenSettings struct {
	UseWrapper bool `json:"use_wrapper,omitempty"`
}

type NpmSettings struct {
	DepType          string `json:"dep_type,omitempty"`
	PnpmMaxTreeDepth int32  `json:"pnpm_max_tree_depth,omitempty"`
}

type PythonSettings struct {
	RequirementsFile string `json:"requirements_file,omitempty"`
}

type CaScannerConfig struct {
	EnableCaScan    bool     `json:"enable_ca_scan,omitempty"`
	ExcludePatterns []string `json:"exclude_patterns,omitempty"`
}

type SastScannerConfig struct {
	EnableSastScan  bool     `json:"enable_sast_scan,omitempty"`
	Language        string   `json:"language,omitempty"`
	IncludePatterns []string `json:"Include_patterns,omitempty"`
	ExcludePatterns []string `json:"exclude_patterns,omitempty"`
	ExcludeRules    []string `json:"exclude_rules,omitempty"`
}

type SecretsScannerConfig struct {
	EnableSecretsScan bool     `json:"enable_secrets_scan,omitempty"`
	IncludePatterns   []string `json:"Include_patterns,omitempty"`
	ExcludePatterns   []string `json:"exclude_patterns,omitempty"`
}

type IacScannerConfig struct {
	EnableIacScan   bool     `json:"enable_iac_scan,omitempty"`
	IncludePatterns []string `json:"Include_patterns,omitempty"`
	ExcludePatterns []string `json:"exclude_patterns,omitempty"`
}

func (cp *ConfigurationProfileService) sendConfigProfileByNameRequest(profileName string) (url string, resp *http.Response, body []byte, err error) {
	if cp.XrayDetails != nil {
		httpDetails := cp.XrayDetails.CreateHttpClientDetails()
		url = fmt.Sprintf("%s%s%s/%s", utils.AddTrailingSlashIfNeeded(cp.XrayDetails.GetUrl()), xscutils.XscInXraySuffix, xscConfigProfileByNameApi, profileName)
		resp, body, _, err = cp.client.SendGet(utils.AppendScopedProjectKeyParam(url, cp.ScopeProjectKey), true, &httpDetails)
		return
	}
	// Backward compatibility
	httpDetails := cp.XscDetails.CreateHttpClientDetails()
	url = fmt.Sprintf("%s%s/%s", utils.AddTrailingSlashIfNeeded(cp.XscDetails.GetUrl()), xscDeprecatedConfigProfileByNameApiSuffix, profileName)
	resp, body, _, err = cp.client.SendGet(url, true, &httpDetails)
	return
}

func (cp *ConfigurationProfileService) GetConfigurationProfileByName(profileName string) (*ConfigProfile, error) {
	url, res, body, err := cp.sendConfigProfileByNameRequest(profileName)
	if err != nil {
		return nil, fmt.Errorf("failed to send GET query to '%s': %q", url, err)
	}
	if err = errorutils.CheckResponseStatusWithBody(res, body, http.StatusOK); err != nil {
		return nil, err
	}

	var profile ConfigProfile
	err = errorutils.CheckError(json.Unmarshal(body, &profile))
	return &profile, err
}

func (cp *ConfigurationProfileService) sendConfigProfileByUrlRequest(repoUrl string) (url string, resp *http.Response, body []byte, err error) {
	if cp.XrayDetails == nil {
		err = errors.New("received empty Xray details")
		return
	}
	httpDetails := cp.XrayDetails.CreateHttpClientDetails()
	url = fmt.Sprintf("%s%s%s", utils.AddTrailingSlashIfNeeded(cp.XrayDetails.GetUrl()), xscutils.XscInXraySuffix, xscConfigProfileByUrlApi)
	requestContent := []byte(fmt.Sprintf(getProfileByUrlBody, repoUrl))
	resp, body, err = cp.client.SendPost(utils.AppendScopedProjectKeyParam(url, cp.ScopeProjectKey), requestContent, &httpDetails)
	return
}

func (cp *ConfigurationProfileService) GetConfigurationProfileByUrl(url string) (*ConfigProfile, error) {
	url, res, body, err := cp.sendConfigProfileByUrlRequest(url)
	if err != nil {
		return nil, fmt.Errorf("failed to send POST query to '%s': %q", url, err)
	}
	if err = errorutils.CheckResponseStatusWithBody(res, body, http.StatusOK); err != nil {
		return nil, err
	}

	var profile ConfigProfile
	err = errorutils.CheckError(json.Unmarshal(body, &profile))
	return &profile, err
}
