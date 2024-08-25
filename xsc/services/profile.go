package services

import (
	"encoding/json"
	"fmt"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"net/http"
)

const (
	ConfigProfileMinXscVersion = "1.11.0"
	xscConfigProfileApi        = "api/v1/profile"
)

type ConfigurationProfileService struct {
	client     *jfroghttpclient.JfrogHttpClient
	XscDetails auth.ServiceDetails
}

func NewConfigurationProfileService(client *jfroghttpclient.JfrogHttpClient) *ConfigurationProfileService {
	return &ConfigurationProfileService{client: client}
}

type ConfigProfile struct {
	ProfileName   string        `json:"profile_name"`
	FrogbotConfig FrogbotConfig `json:"frogbot_config,omitempty"`
	Modules       []Module      `json:"modules"`
	IsDefault     bool          `json:"is_default,omitempty"`
}

type FrogbotConfig struct {
	EmailAuthor                     string `json:"email_author,omitempty"`
	AggregateFixes                  bool   `json:"aggregate_fixes,omitempty"`
	AvoidPreviousPrCommentsDeletion bool   `json:"avoid_previous_pr_comments_deletion,omitempty"`
	BranchNameTemplate              string `json:"branch_name_template,omitempty"`
	PrTitleTemplate                 string `json:"pr_title_template,omitempty"`
	PrCommentTitle                  string `json:"pr_comment_title,omitempty"`
	CommitMessageTemplate           string `json:"commit_message_template,omitempty"`
	ShowSecretsAsPrComment          bool   `json:"show_secrets_as_pr_comment,omitempty"`
}

type Module struct {
	ModuleId                 int32      `json:"module_id,omitempty"`
	ModuleName               string     `json:"module_name"`
	PathFromRoot             string     `json:"path_from_root"`
	ReleasesRepo             string     `json:"releases_repo,omitempty"`
	AnalyzerManagerVersion   string     `json:"analyzer_manager_version,omitempty"`
	AdditionalPathsForModule []string   `json:"additional_paths_for_module,omitempty"`
	ExcludePaths             []string   `json:"exclude_paths,omitempty"`
	ScanConfig               ScanConfig `json:"scan_config"`
	ProtectedBranches        []string   `json:"protected_branches,omitempty"`
	IncludeExcludeMode       int32      `json:"include_exclude_mode,omitempty"`
	IncludeExcludePattern    string     `json:"include_exclude_pattern,omitempty"`
	ReportAnalytics          bool       `json:"report_analytics,omitempty"`
}

type ScanConfig struct {
	ScanTimeout                  int32                     `json:"scan_timeout,omitempty"`
	ExcludePattern               string                    `json:"exclude_pattern,omitempty"`
	EnableScaScan                bool                      `json:"enable_sca_scan,omitempty"`
	EnableContextualAnalysisScan bool                      `json:"enable_contextual_analysis_scan,omitempty"`
	SastScannerConfig            SastScannerConfig         `json:"sast_scanner_config,omitempty"`
	SecretsScannerConfig         SecretsScannerConfig      `json:"secrets_scanner_config,omitempty"`
	IacScannerConfig             IacScannerConfig          `json:"iac_scanner_config,omitempty"`
	ApplicationsScannerConfig    ApplicationsScannerConfig `json:"applications_scanner_config,omitempty"`
	ServicesScannerConfig        ServicesScannerConfig     `json:"services_scanner_config,omitempty"`
}

type SastScannerConfig struct {
	EnableSastScan  bool     `json:"enable_sast_scan,omitempty"`
	Language        string   `json:"language,omitempty"`
	ExcludePatterns []string `json:"exclude_patterns,omitempty"`
	ExcludeRules    []string `json:"exclude_rules,omitempty"`
}

type SecretsScannerConfig struct {
	EnableSecretsScan bool     `json:"enable_secrets_scan,omitempty"`
	ExcludePatterns   []string `json:"exclude_patterns,omitempty"`
}

type IacScannerConfig struct {
	EnableIacScan   bool     `json:"enable_iac_scan,omitempty"`
	ExcludePatterns []string `json:"exclude_patterns,omitempty"`
}

type ApplicationsScannerConfig struct {
	EnableApplicationsScan bool     `json:"enable_applications_scan,omitempty"`
	ExcludePatterns        []string `json:"exclude_patterns,omitempty"`
}

type ServicesScannerConfig struct {
	EnableServicesScan bool     `json:"enable_services_scan,omitempty"`
	ExcludePatterns    []string `json:"exclude_patterns,omitempty"`
}

func (cp *ConfigurationProfileService) GetConfigurationProfile(profileName string) (*ConfigProfile, error) {
	httpDetails := cp.XscDetails.CreateHttpClientDetails()
	url := fmt.Sprintf("%s%s/%s", utils.AddTrailingSlashIfNeeded(cp.XscDetails.GetUrl()), xscConfigProfileApi, profileName)
	res, body, _, err := cp.client.SendGet(url, true, &httpDetails)
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
