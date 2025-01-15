package jfrogcli

import "github.com/jfrog/jfrog-client-go/jfconnect/services/metrics"

type CommandsCountLabels struct {
	ProductID                            string `json:"product_id"`
	ProductVersion                       string `json:"product_version"`
	FeatureID                            string `json:"feature_id"`
	OIDCUsed                             string `json:"oidc_used"`
	JobID                                string `json:"job_id"`
	RunID                                string `json:"run_id"`
	GitRepo                              string `json:"git_repo"`
	GhTokenForCodeScanningAlertsProvided string `json:"gh_token_for_code_scanning_alerts_provided"`
}

type CommandsCountMetric struct {
	metrics.Metric `json:",inline"`
	Labels         CommandsCountLabels `json:"labels"`
}

func NewCommandsCountMetric() CommandsCountMetric {
	return CommandsCountMetric{
		Metric: metrics.Metric{
			Value: 1,
			Name:  "jfcli_commands_count",
		},
	}
}

func (ccm *CommandsCountMetric) MetricsName() string {
	return ccm.Name
}
