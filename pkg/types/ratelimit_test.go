package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"
)

func TestRateLimit_Service_Config_String(t *testing.T) {
	tests := []struct {
		name    string
		config  *RateLimit_Service_Config
		want    string
		wantErr bool
	}{
		{
			name: "simple ratelimit config",
			config: &RateLimit_Service_Config{
				Domain: "foo",
				Descriptors: []RateLimit_Service_Descriptor{
					{
						Key:   "bar",
						Value: "baz",
						RateLimit: v1alpha1.GlobalRateLimit_Limit{
							Unit:            "hour",
							RequestsPerUnit: 1,
						},
					},
				},
			},
			want: `domain: foo
descriptors:
- key: bar
  value: baz
  rate_limit:
    unit: hour
    requests_per_unit: 1
`,
			wantErr: false,
		},
		{
			name: "config with shadow mode enabled",
			config: &RateLimit_Service_Config{
				Domain: "shadow-domain",
				Descriptors: []RateLimit_Service_Descriptor{
					{
						Key:        "header_match",
						Value:      "test-value",
						ShadowMode: true,
						RateLimit: v1alpha1.GlobalRateLimit_Limit{
							Unit:            "minute",
							RequestsPerUnit: 100,
						},
					},
				},
			},
			want: `domain: shadow-domain
descriptors:
- key: header_match
  value: test-value
  shadow_mode: true
  rate_limit:
    unit: minute
    requests_per_unit: 100
`,
			wantErr: false,
		},
		{
			name: "config with nested descriptors",
			config: &RateLimit_Service_Config{
				Domain: "nested-domain",
				Descriptors: []RateLimit_Service_Descriptor{
					{
						Key:   "parent_key",
						Value: "parent_value",
						Descriptors: []RateLimit_Service_Descriptor{
							{
								Key:   "child_key",
								Value: "child_value",
								RateLimit: v1alpha1.GlobalRateLimit_Limit{
									Unit:            "second",
									RequestsPerUnit: 10,
								},
							},
						},
					},
				},
			},
			want: `domain: nested-domain
descriptors:
- key: parent_key
  value: parent_value
  descriptors:
  - key: child_key
    value: child_value
    rate_limit:
      unit: second
      requests_per_unit: 10
`,
			wantErr: false,
		},
		{
			name: "config with multiple descriptors",
			config: &RateLimit_Service_Config{
				Domain: "multi-domain",
				Descriptors: []RateLimit_Service_Descriptor{
					{
						Key:   "first_key",
						Value: "first_value",
						RateLimit: v1alpha1.GlobalRateLimit_Limit{
							Unit:            "hour",
							RequestsPerUnit: 1000,
						},
					},
					{
						Key:   "second_key",
						Value: "second_value",
						RateLimit: v1alpha1.GlobalRateLimit_Limit{
							Unit:            "day",
							RequestsPerUnit: 10000,
						},
					},
				},
			},
			want: `domain: multi-domain
descriptors:
- key: first_key
  value: first_value
  rate_limit:
    unit: hour
    requests_per_unit: 1000
- key: second_key
  value: second_value
  rate_limit:
    unit: day
    requests_per_unit: 10000
`,
			wantErr: false,
		},
		{
			name: "config with key only (no value)",
			config: &RateLimit_Service_Config{
				Domain: "key-only-domain",
				Descriptors: []RateLimit_Service_Descriptor{
					{
						Key: "generic_key",
						RateLimit: v1alpha1.GlobalRateLimit_Limit{
							Unit:            "minute",
							RequestsPerUnit: 50,
						},
					},
				},
			},
			want: `domain: key-only-domain
descriptors:
- key: generic_key
  rate_limit:
    unit: minute
    requests_per_unit: 50
`,
			wantErr: false,
		},
		{
			name:   "empty config",
			config: &RateLimit_Service_Config{},
			want: `{}
`,
			wantErr: false,
		},
		{
			name: "config with domain only",
			config: &RateLimit_Service_Config{
				Domain: "domain-only",
			},
			want: `domain: domain-only
`,
			wantErr: false,
		},
		{
			name: "deeply nested descriptors",
			config: &RateLimit_Service_Config{
				Domain: "deep-domain",
				Descriptors: []RateLimit_Service_Descriptor{
					{
						Key:   "level1",
						Value: "value1",
						Descriptors: []RateLimit_Service_Descriptor{
							{
								Key:   "level2",
								Value: "value2",
								Descriptors: []RateLimit_Service_Descriptor{
									{
										Key:   "level3",
										Value: "value3",
										RateLimit: v1alpha1.GlobalRateLimit_Limit{
											Unit:            "second",
											RequestsPerUnit: 5,
										},
									},
								},
							},
						},
					},
				},
			},
			want: `domain: deep-domain
descriptors:
- key: level1
  value: value1
  descriptors:
  - key: level2
    value: value2
    descriptors:
    - key: level3
      value: value3
      rate_limit:
        unit: second
        requests_per_unit: 5
`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.config.String()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRateLimit_Service_Descriptor_Fields(t *testing.T) {
	tests := []struct {
		name       string
		descriptor RateLimit_Service_Descriptor
		wantKey    string
		wantValue  string
		wantShadow bool
	}{
		{
			name: "descriptor with all fields",
			descriptor: RateLimit_Service_Descriptor{
				Key:        "test_key",
				Value:      "test_value",
				ShadowMode: true,
				RateLimit: v1alpha1.GlobalRateLimit_Limit{
					Unit:            "minute",
					RequestsPerUnit: 100,
				},
			},
			wantKey:    "test_key",
			wantValue:  "test_value",
			wantShadow: true,
		},
		{
			name: "descriptor with shadow mode disabled",
			descriptor: RateLimit_Service_Descriptor{
				Key:        "no_shadow_key",
				Value:      "no_shadow_value",
				ShadowMode: false,
			},
			wantKey:    "no_shadow_key",
			wantValue:  "no_shadow_value",
			wantShadow: false,
		},
		{
			name:       "empty descriptor",
			descriptor: RateLimit_Service_Descriptor{},
			wantKey:    "",
			wantValue:  "",
			wantShadow: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantKey, tt.descriptor.Key)
			assert.Equal(t, tt.wantValue, tt.descriptor.Value)
			assert.Equal(t, tt.wantShadow, tt.descriptor.ShadowMode)
		})
	}
}
