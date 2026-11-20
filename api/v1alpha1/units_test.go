package v1alpha1

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestLocalRateLimitLimitToTokenBucket(t *testing.T) {
	tests := []struct {
		unit string
		want string
	}{
		{"second", "1s"},
		{"minute", "60s"},
		{"hour", "3600s"},
		{"day", "86400s"},
		{"week", "604800s"},
		{"month", "2592000s"},
		{"year", "31536000s"},
		{"WEEK", "604800s"},
		{"Month", "2592000s"},
		{"", "1s"},
	}

	for _, tt := range tests {
		t.Run(tt.unit, func(t *testing.T) {
			bucket := LocalRateLimit_Limit{Unit: tt.unit, RequestsPerUnit: 10}.ToTokenBucket()
			assert.Equal(t, tt.want, bucket.FillInterval)
			assert.Equal(t, 10, bucket.MaxTokens)
			assert.Equal(t, 1, bucket.TokensPerFill)
		})
	}
}

type crdSchema struct {
	Pattern    string               `yaml:"pattern"`
	Properties map[string]crdSchema `yaml:"properties"`
}

type crdFile struct {
	Spec struct {
		Versions []struct {
			Schema struct {
				OpenAPIV3Schema crdSchema `yaml:"openAPIV3Schema"`
			} `yaml:"schema"`
		} `yaml:"versions"`
	} `yaml:"spec"`
}

// crdPattern returns the pattern the generated CRD puts on the property at
// path, so the tests check what the API server enforces.
func crdPattern(t *testing.T, name string, path ...string) *regexp.Regexp {
	t.Helper()

	data, err := os.ReadFile(filepath.Join("..", "..", "config", "crd", "bases", name))
	require.NoError(t, err)

	var crd crdFile
	require.NoError(t, yaml.Unmarshal(data, &crd))
	require.NotEmpty(t, crd.Spec.Versions)

	schema := crd.Spec.Versions[0].Schema.OpenAPIV3Schema
	for _, p := range path {
		next, ok := schema.Properties[p]
		require.True(t, ok, "property %q not found in %s", p, name)
		schema = next
	}
	require.NotEmpty(t, schema.Pattern)

	return regexp.MustCompile(schema.Pattern)
}

func TestLimitUnitPattern(t *testing.T) {
	patterns := map[string]*regexp.Regexp{
		"GlobalRateLimit": crdPattern(t, "ratelimit.zufardhiyaulhaq.com_globalratelimits.yaml", "spec", "limit", "unit"),
		"LocalRateLimit":  crdPattern(t, "ratelimit.zufardhiyaulhaq.com_localratelimits.yaml", "spec", "limit", "unit"),
	}

	valid := []string{"second", "minute", "hour", "day", "week", "month", "year", "HOUR", "Week", "MONTH"}
	invalid := []string{"", "weeks", "fortnight", "1h", " hour", "hour "}

	for name, pattern := range patterns {
		for _, unit := range valid {
			assert.True(t, pattern.MatchString(unit), "%s should accept %q", name, unit)
		}
		for _, unit := range invalid {
			assert.False(t, pattern.MatchString(unit), "%s should reject %q", name, unit)
		}
	}
}

func TestWeekResetDayPattern(t *testing.T) {
	pattern := crdPattern(t, "ratelimit.zufardhiyaulhaq.com_ratelimitservices.yaml", "spec", "window", "weekResetDay")

	for _, day := range []string{"Sunday", "monday", "TUESDAY", "Wednesday", "thursday", "Friday", "saturday"} {
		assert.True(t, pattern.MatchString(day), "should accept %q", day)
	}
	for _, day := range []string{"", "Mon", "funday", "Monday "} {
		assert.False(t, pattern.MatchString(day), "should reject %q", day)
	}
}
