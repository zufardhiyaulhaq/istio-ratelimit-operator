package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetricMapper_String(t *testing.T) {
	tests := []struct {
		name    string
		mapper  *MetricMapper
		want    string
		wantErr bool
	}{
		{
			name: "simple metric mapping",
			mapper: &MetricMapper{
				Mappings: []MetricMapping{
					{
						Match: "ratelimit.service.rate_limit.*.*.total_hits",
						Name:  "ratelimit_service_total_hits",
						Labels: map[string]string{
							"domain":     "$1",
							"descriptor": "$2",
						},
					},
				},
			},
			want: `mappings:
- match: ratelimit.service.rate_limit.*.*.total_hits
  name: ratelimit_service_total_hits
  labels:
    descriptor: $2
    domain: $1
`,
			wantErr: false,
		},
		{
			name: "metric mapping with histogram observer type",
			mapper: &MetricMapper{
				Mappings: []MetricMapping{
					{
						Match:        "ratelimit.service.rate_limit.*.*.response_time",
						Name:         "ratelimit_service_response_time",
						ObserverType: ObserverTypeHistogram,
						Labels: map[string]string{
							"domain": "$1",
						},
					},
				},
			},
			want: `mappings:
- match: ratelimit.service.rate_limit.*.*.response_time
  name: ratelimit_service_response_time
  labels:
    domain: $1
  observer_type: histogram
`,
			wantErr: false,
		},
		{
			name: "metric mapping with summary observer type",
			mapper: &MetricMapper{
				Mappings: []MetricMapping{
					{
						Match:        "ratelimit.service.rate_limit.*.*.latency",
						Name:         "ratelimit_service_latency",
						ObserverType: ObserverTypeSummary,
					},
				},
			},
			want: `mappings:
- match: ratelimit.service.rate_limit.*.*.latency
  name: ratelimit_service_latency
  observer_type: summary
`,
			wantErr: false,
		},
		{
			name: "metric mapping with timer type",
			mapper: &MetricMapper{
				Mappings: []MetricMapping{
					{
						Match:     "ratelimit.service.rate_limit.*.*.timer",
						Name:      "ratelimit_service_timer",
						TimerType: ObserverTypeHistogram,
					},
				},
			},
			want: `mappings:
- match: ratelimit.service.rate_limit.*.*.timer
  name: ratelimit_service_timer
  timer_type: histogram
`,
			wantErr: false,
		},
		{
			name: "metric mapping with counter metric type",
			mapper: &MetricMapper{
				Mappings: []MetricMapping{
					{
						Match:           "ratelimit.service.rate_limit.*.*.counter",
						Name:            "ratelimit_service_counter",
						MatchMetricType: MetricTypeCounter,
					},
				},
			},
			want: `mappings:
- match: ratelimit.service.rate_limit.*.*.counter
  name: ratelimit_service_counter
  match_metric_type: counter
`,
			wantErr: false,
		},
		{
			name: "metric mapping with gauge metric type",
			mapper: &MetricMapper{
				Mappings: []MetricMapping{
					{
						Match:           "ratelimit.service.rate_limit.*.*.gauge",
						Name:            "ratelimit_service_gauge",
						MatchMetricType: MetricTypeGauge,
					},
				},
			},
			want: `mappings:
- match: ratelimit.service.rate_limit.*.*.gauge
  name: ratelimit_service_gauge
  match_metric_type: gauge
`,
			wantErr: false,
		},
		{
			name: "metric mapping with observer metric type",
			mapper: &MetricMapper{
				Mappings: []MetricMapping{
					{
						Match:           "ratelimit.service.rate_limit.*.*.observer",
						Name:            "ratelimit_service_observer",
						MatchMetricType: MetricTypeObserver,
					},
				},
			},
			want: `mappings:
- match: ratelimit.service.rate_limit.*.*.observer
  name: ratelimit_service_observer
  match_metric_type: observer
`,
			wantErr: false,
		},
		{
			name: "metric mapping with timer metric type (deprecated)",
			mapper: &MetricMapper{
				Mappings: []MetricMapping{
					{
						Match:           "ratelimit.service.rate_limit.*.*.timer_metric",
						Name:            "ratelimit_service_timer_metric",
						MatchMetricType: MetricTypeTimer,
					},
				},
			},
			want: `mappings:
- match: ratelimit.service.rate_limit.*.*.timer_metric
  name: ratelimit_service_timer_metric
  match_metric_type: timer
`,
			wantErr: false,
		},
		{
			name: "multiple metric mappings",
			mapper: &MetricMapper{
				Mappings: []MetricMapping{
					{
						Match: "ratelimit.service.rate_limit.*.*.total_hits",
						Name:  "ratelimit_service_total_hits",
						Labels: map[string]string{
							"domain": "$1",
						},
					},
					{
						Match: "ratelimit.service.rate_limit.*.*.over_limit",
						Name:  "ratelimit_service_over_limit",
						Labels: map[string]string{
							"domain": "$1",
						},
					},
					{
						Match: "ratelimit.service.rate_limit.*.*.near_limit",
						Name:  "ratelimit_service_near_limit",
						Labels: map[string]string{
							"domain": "$1",
						},
					},
				},
			},
			want: `mappings:
- match: ratelimit.service.rate_limit.*.*.total_hits
  name: ratelimit_service_total_hits
  labels:
    domain: $1
- match: ratelimit.service.rate_limit.*.*.over_limit
  name: ratelimit_service_over_limit
  labels:
    domain: $1
- match: ratelimit.service.rate_limit.*.*.near_limit
  name: ratelimit_service_near_limit
  labels:
    domain: $1
`,
			wantErr: false,
		},
		{
			name: "metric mapping with all fields",
			mapper: &MetricMapper{
				Mappings: []MetricMapping{
					{
						Match:           "ratelimit.service.rate_limit.*.*.complete",
						Name:            "ratelimit_service_complete",
						ObserverType:    ObserverTypeHistogram,
						TimerType:       ObserverTypeSummary,
						MatchMetricType: MetricTypeCounter,
						Labels: map[string]string{
							"domain":     "$1",
							"descriptor": "$2",
							"service":    "ratelimit",
						},
					},
				},
			},
			want: `mappings:
- match: ratelimit.service.rate_limit.*.*.complete
  name: ratelimit_service_complete
  labels:
    descriptor: $2
    domain: $1
    service: ratelimit
  observer_type: histogram
  timer_type: summary
  match_metric_type: counter
`,
			wantErr: false,
		},
		{
			name:   "empty metric mapper",
			mapper: &MetricMapper{},
			want: `mappings: []
`,
			wantErr: false,
		},
		{
			name: "metric mapping with empty mappings slice",
			mapper: &MetricMapper{
				Mappings: []MetricMapping{},
			},
			want: `mappings: []
`,
			wantErr: false,
		},
		{
			name: "metric mapping without labels",
			mapper: &MetricMapper{
				Mappings: []MetricMapping{
					{
						Match: "ratelimit.service.rate_limit.total",
						Name:  "ratelimit_service_total",
					},
				},
			},
			want: `mappings:
- match: ratelimit.service.rate_limit.total
  name: ratelimit_service_total
`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.mapper.String()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestObserverType_Constants(t *testing.T) {
	tests := []struct {
		name     string
		obsType  ObserverType
		expected string
	}{
		{
			name:     "histogram observer type",
			obsType:  ObserverTypeHistogram,
			expected: "histogram",
		},
		{
			name:     "summary observer type",
			obsType:  ObserverTypeSummary,
			expected: "summary",
		},
		{
			name:     "default observer type",
			obsType:  ObserverTypeDefault,
			expected: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, string(tt.obsType))
		})
	}
}

func TestMetricType_Constants(t *testing.T) {
	tests := []struct {
		name       string
		metricType MetricType
		expected   string
	}{
		{
			name:       "counter metric type",
			metricType: MetricTypeCounter,
			expected:   "counter",
		},
		{
			name:       "gauge metric type",
			metricType: MetricTypeGauge,
			expected:   "gauge",
		},
		{
			name:       "observer metric type",
			metricType: MetricTypeObserver,
			expected:   "observer",
		},
		{
			name:       "timer metric type (deprecated)",
			metricType: MetricTypeTimer,
			expected:   "timer",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, string(tt.metricType))
		})
	}
}

func TestMetricMapping_Fields(t *testing.T) {
	tests := []struct {
		name            string
		mapping         MetricMapping
		wantMatch       string
		wantName        string
		wantLabelsCount int
		wantObsType     ObserverType
		wantTimerType   ObserverType
		wantMetricType  MetricType
	}{
		{
			name: "mapping with all fields populated",
			mapping: MetricMapping{
				Match:           "test.match.pattern",
				Name:            "test_metric_name",
				Labels:          map[string]string{"label1": "value1", "label2": "value2"},
				ObserverType:    ObserverTypeHistogram,
				TimerType:       ObserverTypeSummary,
				MatchMetricType: MetricTypeCounter,
			},
			wantMatch:       "test.match.pattern",
			wantName:        "test_metric_name",
			wantLabelsCount: 2,
			wantObsType:     ObserverTypeHistogram,
			wantTimerType:   ObserverTypeSummary,
			wantMetricType:  MetricTypeCounter,
		},
		{
			name: "mapping with minimal fields",
			mapping: MetricMapping{
				Match: "minimal.match",
				Name:  "minimal_name",
			},
			wantMatch:       "minimal.match",
			wantName:        "minimal_name",
			wantLabelsCount: 0,
			wantObsType:     ObserverTypeDefault,
			wantTimerType:   ObserverTypeDefault,
			wantMetricType:  "",
		},
		{
			name:            "empty mapping",
			mapping:         MetricMapping{},
			wantMatch:       "",
			wantName:        "",
			wantLabelsCount: 0,
			wantObsType:     ObserverTypeDefault,
			wantTimerType:   ObserverTypeDefault,
			wantMetricType:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantMatch, tt.mapping.Match)
			assert.Equal(t, tt.wantName, tt.mapping.Name)
			assert.Len(t, tt.mapping.Labels, tt.wantLabelsCount)
			assert.Equal(t, tt.wantObsType, tt.mapping.ObserverType)
			assert.Equal(t, tt.wantTimerType, tt.mapping.TimerType)
			assert.Equal(t, tt.wantMetricType, tt.mapping.MatchMetricType)
		})
	}
}

func TestMetricMapper_Mappings_Access(t *testing.T) {
	mapper := &MetricMapper{
		Mappings: []MetricMapping{
			{
				Match: "first.match",
				Name:  "first_metric",
			},
			{
				Match: "second.match",
				Name:  "second_metric",
			},
		},
	}

	assert.Len(t, mapper.Mappings, 2)
	assert.Equal(t, "first.match", mapper.Mappings[0].Match)
	assert.Equal(t, "second.match", mapper.Mappings[1].Match)
}
