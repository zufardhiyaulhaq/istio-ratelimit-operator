package types

import (
	"testing"

	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"
)

func TestRateLimit_Service_Config_String(t *testing.T) {
	type fields struct {
		Domain      string
		Descriptors []RateLimit_Service_Descriptor
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "simple ratelimit config",
			fields: fields{
				Domain: "foo",
				Descriptors: []RateLimit_Service_Descriptor{
					{
						Key:        "bar",
						Value:      "baz",
						ShadowMode: false,
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
  shadow_mode: false
`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &RateLimit_Service_Config{
				Domain:      tt.fields.Domain,
				Descriptors: tt.fields.Descriptors,
			}
			got, err := n.String()
			if (err != nil) != tt.wantErr {
				t.Errorf("RateLimit_Service_Config.String() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RateLimit_Service_Config.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
