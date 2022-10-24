package types

import "gopkg.in/yaml.v2"

type MetricMapper struct {
	Mappings []MetricMapping `yaml:"mappings"`
}

type MetricMapping struct {
	Match           string            `yaml:"match"`
	Name            string            `yaml:"name"`
	Labels          map[string]string `yaml:"labels,omitempty"`
	ObserverType    ObserverType      `yaml:"observer_type,omitempty"`
	TimerType       ObserverType      `yaml:"timer_type,omitempty"`
	MatchMetricType MetricType        `yaml:"match_metric_type,omitempty"`
}

type ObserverType string

const (
	ObserverTypeHistogram ObserverType = "histogram"
	ObserverTypeSummary   ObserverType = "summary"
	ObserverTypeDefault   ObserverType = ""
)

type MetricType string

const (
	MetricTypeCounter  MetricType = "counter"
	MetricTypeGauge    MetricType = "gauge"
	MetricTypeObserver MetricType = "observer"
	MetricTypeTimer    MetricType = "timer" // DEPRECATED
)

func (n *MetricMapper) String() (string, error) {
	bytes, err := yaml.Marshal(n)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
