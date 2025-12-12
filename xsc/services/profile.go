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
	ProfileName   string        `json:"profile_name"`
	GeneralConfig GeneralConfig `json:"general_config,omitempty"`
	FrogbotConfig FrogbotConfig `json:"frogbot_config,omitempty"`
	Modules       []Module      `json:"modules"`
}

type GeneralConfig struct {
	ScannersDownloadPath    string   `json:"scanners_download_path,omitempty"`
	GeneralExcludePatterns  []string `json:"general_exclude_patterns,omitempty"`
	FailUponAnyScannerError bool     `json:"fail_upon_any_scanner_error,omitempty"`
}

type FrogbotConfig struct {
	AggregateFixes                      bool   `json:"aggregate_fixes,omitempty"`
	HideSuccessBannerForNoIssues        bool   `json:"hide_success_banner_for_no_issues,omitempty"`
	BranchNameTemplate                  string `json:"branch_name_template,omitempty"`
	PrTitleTemplate                     string `json:"pr_title_template,omitempty"`
	CommitMessageTemplate               string `json:"commit_message_template,omitempty"`
	ShowSecretsAsPrComment              bool   `json:"show_secrets_as_pr_comment,omitempty"`
	CreateAutoFixPr                     bool   `json:"create_auto_fix_pr,omitempty"`
	IncludeVulnerabilitiesAndViolations bool   `json:"include_vulnerabilities_and_violations,omitempty"`
	MinSeverityToDisplay                string `json:"min_severity_to_display,omitempty"`
	DisplayFixableOnly                  bool   `json:"display_fixable_only,omitempty"`
}

type Module struct {
	ModuleId     int32      `json:"module_id,omitempty"`
	ModuleName   string     `json:"module_name"`
	PathFromRoot string     `json:"path_from_root"`
	ScanConfig   ScanConfig `json:"scan_config"`
}

type ScanConfig struct {
	ScaScannerConfig                ScaScannerConfig     `json:"sca_scanner_config,omitempty"`
	ContextualAnalysisScannerConfig CaScannerConfig      `json:"contextual_analysis_scanner_config,omitempty"`
	SastScannerConfig               SastScannerConfig    `json:"sast_scanner_config,omitempty"`
	SecretsScannerConfig            SecretsScannerConfig `json:"secrets_scanner_config,omitempty"`
	IacScannerConfig                IacScannerConfig     `json:"iac_scanner_config,omitempty"`
}

type ScaScannerConfig struct {
	EnableScaScan   bool     `json:"enable_sca_scan,omitempty"`
	ExcludePatterns []string `json:"exclude_patterns,omitempty"`
}

type CaScannerConfig struct {
	EnableCaScan    bool     `json:"enable_ca_scan,omitempty"`
	ExcludePatterns []string `json:"exclude_patterns,omitempty"`
}

type SastScannerConfig struct {
	EnableSastScan  bool     `json:"enable_sast_scan,omitempty"`
	ExcludePatterns []string `json:"exclude_patterns,omitempty"`
	ExcludeRules    []string `json:"exclude_rules,omitempty"`
}

type SecretsScannerConfig struct {
	EnableSecretsScan   bool     `json:"enable_secrets_scan,omitempty"`
	ValidateSecrets     bool     `json:"validate_secrets,omitempty"`
	ExcludePatterns     []string `json:"exclude_patterns,omitempty"`
	EnableCustomSecrets bool     `json:"enable_custom_secrets,omitempty"`
}

type IacScannerConfig struct {
	EnableIacScan   bool     `json:"enable_iac_scan,omitempty"`
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
