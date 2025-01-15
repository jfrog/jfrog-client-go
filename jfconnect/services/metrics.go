package services

type VisibilityMetric interface {
	MetricValue() int
	MetricName() string
}

type Metric struct {
	Value int    `json:"value"`
	Name  string `json:"metrics_name"`
}

func (m *Metric) MetricValue() int {
	return m.Value
}

func (m *Metric) MetricName() string {
	return m.Name
}
