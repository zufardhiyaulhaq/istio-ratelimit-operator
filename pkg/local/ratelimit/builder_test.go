package ratelimit_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/local/ratelimit"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNewConfigBuilder(t *testing.T) {
	builder := ratelimit.NewConfigBuilder()
	assert.NotNil(t, builder)
	assert.Equal(t, v1alpha1.LocalRateLimit{}, builder.RateLimit)
	assert.Equal(t, v1alpha1.LocalRateLimitConfig{}, builder.Config)
	assert.Nil(t, builder.Versions)
	assert.Nil(t, builder.Labels)
}

func TestConfigBuilder_SetRateLimit(t *testing.T) {
	tests := []struct {
		name      string
		ratelimit v1alpha1.LocalRateLimit
	}{
		{
			name: "set basic ratelimit",
			ratelimit: v1alpha1.LocalRateLimit{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ratelimit",
					Namespace: "default",
				},
				Spec: v1alpha1.LocalRateLimitSpec{
					Config: "test-config",
					Selector: v1alpha1.LocalRateLimitSelector{
						VHost: "test.example.com:443",
					},
					Limit: &v1alpha1.LocalRateLimit_Limit{
						Unit:            "second",
						RequestsPerUnit: 10,
					},
				},
			},
		},
		{
			name: "set ratelimit with route selector",
			ratelimit: v1alpha1.LocalRateLimit{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ratelimit-route",
					Namespace: "istio-system",
				},
				Spec: v1alpha1.LocalRateLimitSpec{
					Config: "gateway-config",
					Selector: v1alpha1.LocalRateLimitSelector{
						VHost: "api.example.com:443",
						Route: stringPtr("my-route"),
					},
					Limit: &v1alpha1.LocalRateLimit_Limit{
						Unit:            "minute",
						RequestsPerUnit: 100,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := ratelimit.NewConfigBuilder()
			result := builder.SetRateLimit(tt.ratelimit)

			assert.Same(t, builder, result, "SetRateLimit should return the same builder for chaining")
			assert.Equal(t, tt.ratelimit, builder.RateLimit)
		})
	}
}

func TestConfigBuilder_SetConfig(t *testing.T) {
	tests := []struct {
		name   string
		config v1alpha1.LocalRateLimitConfig
	}{
		{
			name: "set gateway config",
			config: v1alpha1.LocalRateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "gateway-config",
					Namespace: "istio-system",
				},
				Spec: v1alpha1.LocalRateLimitConfigSpec{
					Type: v1alpha1.Gateway,
					Selector: v1alpha1.LocalRateLimitConfigSelector{
						Labels:       map[string]string{"app": "gateway"},
						IstioVersion: []string{"1.9"},
					},
				},
			},
		},
		{
			name: "set sidecar config",
			config: v1alpha1.LocalRateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sidecar-config",
					Namespace: "default",
				},
				Spec: v1alpha1.LocalRateLimitConfigSpec{
					Type: v1alpha1.Sidecar,
					Selector: v1alpha1.LocalRateLimitConfigSelector{
						Labels:       map[string]string{"app": "my-service"},
						IstioVersion: []string{"1.10", "1.11"},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := ratelimit.NewConfigBuilder()
			result := builder.SetConfig(tt.config)

			assert.Same(t, builder, result, "SetConfig should return the same builder for chaining")
			assert.Equal(t, tt.config, builder.Config)
		})
	}
}

func TestConfigBuilder_SetVersions(t *testing.T) {
	tests := []struct {
		name     string
		versions []string
	}{
		{
			name:     "set single version",
			versions: []string{"1.9"},
		},
		{
			name:     "set multiple versions",
			versions: []string{"1.9", "1.10", "1.11"},
		},
		{
			name:     "set empty versions",
			versions: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := ratelimit.NewConfigBuilder()
			result := builder.SetVersions(tt.versions)

			assert.Same(t, builder, result, "SetVersions should return the same builder for chaining")
			assert.Equal(t, tt.versions, builder.Versions)
		})
	}
}

func TestConfigBuilder_SetLabels(t *testing.T) {
	tests := []struct {
		name   string
		labels map[string]string
	}{
		{
			name:   "set single label",
			labels: map[string]string{"app": "test"},
		},
		{
			name:   "set multiple labels",
			labels: map[string]string{"app": "test", "version": "v1", "environment": "prod"},
		},
		{
			name:   "set empty labels",
			labels: map[string]string{},
		},
		{
			name:   "set nil labels",
			labels: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := ratelimit.NewConfigBuilder()
			result := builder.SetLabels(tt.labels)

			assert.Same(t, builder, result, "SetLabels should return the same builder for chaining")
			assert.Equal(t, tt.labels, builder.Labels)
		})
	}
}

func TestConfigBuilder_Build(t *testing.T) {
	tests := []struct {
		name          string
		config        v1alpha1.LocalRateLimitConfig
		ratelimit     v1alpha1.LocalRateLimit
		expectedCount int
		expectedError bool
		errorContains string
	}{
		{
			name: "build gateway envoy filters for single version",
			config: v1alpha1.LocalRateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "gateway-config",
					Namespace: "istio-system",
				},
				Spec: v1alpha1.LocalRateLimitConfigSpec{
					Type: v1alpha1.Gateway,
					Selector: v1alpha1.LocalRateLimitConfigSelector{
						Labels:       map[string]string{"app": "gateway"},
						IstioVersion: []string{"1.9"},
					},
				},
			},
			ratelimit: v1alpha1.LocalRateLimit{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ratelimit",
					Namespace: "istio-system",
				},
				Spec: v1alpha1.LocalRateLimitSpec{
					Config: "gateway-config",
					Selector: v1alpha1.LocalRateLimitSelector{
						VHost: "test.example.com:443",
					},
					Limit: &v1alpha1.LocalRateLimit_Limit{
						Unit:            "second",
						RequestsPerUnit: 10,
					},
				},
			},
			expectedCount: 1,
			expectedError: false,
		},
		{
			name: "build gateway envoy filters for multiple versions",
			config: v1alpha1.LocalRateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "gateway-config",
					Namespace: "istio-system",
				},
				Spec: v1alpha1.LocalRateLimitConfigSpec{
					Type: v1alpha1.Gateway,
					Selector: v1alpha1.LocalRateLimitConfigSelector{
						Labels:       map[string]string{"app": "gateway"},
						IstioVersion: []string{"1.9", "1.10", "1.11"},
					},
				},
			},
			ratelimit: v1alpha1.LocalRateLimit{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ratelimit",
					Namespace: "istio-system",
				},
				Spec: v1alpha1.LocalRateLimitSpec{
					Config: "gateway-config",
					Selector: v1alpha1.LocalRateLimitSelector{
						VHost: "test.example.com:443",
					},
					Limit: &v1alpha1.LocalRateLimit_Limit{
						Unit:            "minute",
						RequestsPerUnit: 100,
					},
				},
			},
			expectedCount: 3,
			expectedError: false,
		},
		{
			name: "build sidecar envoy filters for single version",
			config: v1alpha1.LocalRateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sidecar-config",
					Namespace: "default",
				},
				Spec: v1alpha1.LocalRateLimitConfigSpec{
					Type: v1alpha1.Sidecar,
					Selector: v1alpha1.LocalRateLimitConfigSelector{
						Labels:       map[string]string{"app": "my-service"},
						IstioVersion: []string{"1.10"},
					},
				},
			},
			ratelimit: v1alpha1.LocalRateLimit{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sidecar-ratelimit",
					Namespace: "default",
				},
				Spec: v1alpha1.LocalRateLimitSpec{
					Config: "sidecar-config",
					Selector: v1alpha1.LocalRateLimitSelector{
						VHost: "inbound|8080||",
					},
					Limit: &v1alpha1.LocalRateLimit_Limit{
						Unit:            "hour",
						RequestsPerUnit: 1000,
					},
				},
			},
			expectedCount: 1,
			expectedError: false,
		},
		{
			name: "error with invalid version",
			config: v1alpha1.LocalRateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "gateway-config",
					Namespace: "istio-system",
				},
				Spec: v1alpha1.LocalRateLimitConfigSpec{
					Type: v1alpha1.Gateway,
					Selector: v1alpha1.LocalRateLimitConfigSelector{
						Labels:       map[string]string{"app": "gateway"},
						IstioVersion: []string{"invalid-version"},
					},
				},
			},
			ratelimit: v1alpha1.LocalRateLimit{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ratelimit",
					Namespace: "istio-system",
				},
				Spec: v1alpha1.LocalRateLimitSpec{
					Config: "gateway-config",
					Selector: v1alpha1.LocalRateLimitSelector{
						VHost: "test.example.com:443",
					},
					Limit: &v1alpha1.LocalRateLimit_Limit{
						Unit:            "second",
						RequestsPerUnit: 10,
					},
				},
			},
			expectedCount: 0,
			expectedError: true,
			errorContains: "cannot parse version",
		},
		{
			name: "error with unsupported version",
			config: v1alpha1.LocalRateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "gateway-config",
					Namespace: "istio-system",
				},
				Spec: v1alpha1.LocalRateLimitConfigSpec{
					Type: v1alpha1.Gateway,
					Selector: v1alpha1.LocalRateLimitConfigSelector{
						Labels:       map[string]string{"app": "gateway"},
						IstioVersion: []string{"1.5"},
					},
				},
			},
			ratelimit: v1alpha1.LocalRateLimit{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ratelimit",
					Namespace: "istio-system",
				},
				Spec: v1alpha1.LocalRateLimitSpec{
					Config: "gateway-config",
					Selector: v1alpha1.LocalRateLimitSelector{
						VHost: "test.example.com:443",
					},
					Limit: &v1alpha1.LocalRateLimit_Limit{
						Unit:            "second",
						RequestsPerUnit: 10,
					},
				},
			},
			expectedCount: 0,
			expectedError: true,
			errorContains: "version not supported",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := ratelimit.NewConfigBuilder().
				SetConfig(tt.config).
				SetRateLimit(tt.ratelimit)

			envoyFilters, err := builder.Build()

			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
				assert.Len(t, envoyFilters, tt.expectedCount)

				for i, filter := range envoyFilters {
					expectedVersion := tt.config.Spec.Selector.IstioVersion[i]
					assert.Equal(t, tt.ratelimit.Name+"-"+expectedVersion, filter.Name)
					assert.Equal(t, tt.ratelimit.Namespace, filter.Namespace)
					assert.Equal(t, expectedVersion, filter.Labels["istio/version"])
				}
			}
		})
	}
}

func TestConfigBuilder_BuildChaining(t *testing.T) {
	config := v1alpha1.LocalRateLimitConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "gateway-config",
			Namespace: "istio-system",
		},
		Spec: v1alpha1.LocalRateLimitConfigSpec{
			Type: v1alpha1.Gateway,
			Selector: v1alpha1.LocalRateLimitConfigSelector{
				Labels:       map[string]string{"app": "gateway"},
				IstioVersion: []string{"1.9"},
			},
		},
	}

	rl := v1alpha1.LocalRateLimit{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-ratelimit",
			Namespace: "istio-system",
		},
		Spec: v1alpha1.LocalRateLimitSpec{
			Config: "gateway-config",
			Selector: v1alpha1.LocalRateLimitSelector{
				VHost: "test.example.com:443",
			},
			Limit: &v1alpha1.LocalRateLimit_Limit{
				Unit:            "second",
				RequestsPerUnit: 10,
			},
		},
	}

	versions := []string{"1.9"}
	labels := map[string]string{"custom": "label"}

	// Test method chaining
	envoyFilters, err := ratelimit.NewConfigBuilder().
		SetConfig(config).
		SetRateLimit(rl).
		SetVersions(versions).
		SetLabels(labels).
		Build()

	assert.NoError(t, err)
	assert.Len(t, envoyFilters, 1)
}

func TestConfigBuilder_BuildEmptyVersions(t *testing.T) {
	config := v1alpha1.LocalRateLimitConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "gateway-config",
			Namespace: "istio-system",
		},
		Spec: v1alpha1.LocalRateLimitConfigSpec{
			Type: v1alpha1.Gateway,
			Selector: v1alpha1.LocalRateLimitConfigSelector{
				Labels:       map[string]string{"app": "gateway"},
				IstioVersion: []string{},
			},
		},
	}

	rl := v1alpha1.LocalRateLimit{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-ratelimit",
			Namespace: "istio-system",
		},
		Spec: v1alpha1.LocalRateLimitSpec{
			Config: "gateway-config",
			Selector: v1alpha1.LocalRateLimitSelector{
				VHost: "test.example.com:443",
			},
			Limit: &v1alpha1.LocalRateLimit_Limit{
				Unit:            "second",
				RequestsPerUnit: 10,
			},
		},
	}

	builder := ratelimit.NewConfigBuilder().
		SetConfig(config).
		SetRateLimit(rl)

	envoyFilters, err := builder.Build()

	assert.NoError(t, err)
	assert.Empty(t, envoyFilters)
}

// Helper function for pointer to string
func stringPtr(s string) *string {
	return &s
}
