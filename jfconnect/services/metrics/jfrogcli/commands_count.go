package jfrogcli

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
	Value  int                 `json:"value"`
	Name   string              `json:"metrics_name"`
	Labels CommandsCountLabels `json:"labels"`
}

func NewCommandsCountMetric() CommandsCountMetric {
	return CommandsCountMetric{
		Value: 1,
		Name:  "jfcli_commands_count",
	}
}

func (ccm *CommandsCountMetric) MetricsName() string {
	return ccm.Name
}
