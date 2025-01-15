package metrics

type VisibilityMetric interface {
	MetricsName() string
}

type Metric struct {
	Value int    `json:"value"`
	Name  string `json:"metrics_name"`
}
