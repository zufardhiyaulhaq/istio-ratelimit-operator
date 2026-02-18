package ratelimit_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/local/ratelimit"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGetConfigFactory(t *testing.T) {
	tests := []struct {
		name          string
		version       string
		config        v1alpha1.LocalRateLimitConfig
		ratelimit     v1alpha1.LocalRateLimit
		expectedType  string
		expectedError bool
		errorContains string
	}{
		{
			name:    "valid gateway config with version 1.9",
			version: "1.9",
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
			expectedType:  "*ratelimit.V3GatewayBuilder",
			expectedError: false,
		},
		{
			name:    "valid sidecar config with version 1.9",
			version: "1.9",
			config: v1alpha1.LocalRateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sidecar-config",
					Namespace: "default",
				},
				Spec: v1alpha1.LocalRateLimitConfigSpec{
					Type: v1alpha1.Sidecar,
					Selector: v1alpha1.LocalRateLimitConfigSelector{
						Labels:       map[string]string{"app": "my-service"},
						IstioVersion: []string{"1.9"},
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
			expectedType:  "*ratelimit.V3SidecarBuilder",
			expectedError: false,
		},
		{
			name:    "valid gateway config with version 1.10",
			version: "1.10",
			config: v1alpha1.LocalRateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "gateway-config",
					Namespace: "istio-system",
				},
				Spec: v1alpha1.LocalRateLimitConfigSpec{
					Type: v1alpha1.Gateway,
					Selector: v1alpha1.LocalRateLimitConfigSelector{
						Labels:       map[string]string{"app": "gateway"},
						IstioVersion: []string{"1.10"},
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
			expectedType:  "*ratelimit.V3GatewayBuilder",
			expectedError: false,
		},
		{
			name:    "valid sidecar config with version 1.7",
			version: "1.7",
			config: v1alpha1.LocalRateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sidecar-config",
					Namespace: "default",
				},
				Spec: v1alpha1.LocalRateLimitConfigSpec{
					Type: v1alpha1.Sidecar,
					Selector: v1alpha1.LocalRateLimitConfigSelector{
						Labels:       map[string]string{"app": "my-service"},
						IstioVersion: []string{"1.7"},
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
						Unit:            "minute",
						RequestsPerUnit: 100,
					},
				},
			},
			expectedType:  "*ratelimit.V3SidecarBuilder",
			expectedError: false,
		},
		{
			name:    "invalid version format",
			version: "invalid-version",
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
			expectedType:  "",
			expectedError: true,
			errorContains: "cannot parse version",
		},
		{
			name:    "unsupported version 1.5",
			version: "1.5",
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
			expectedType:  "",
			expectedError: true,
			errorContains: "version not supported",
		},
		{
			name:    "unsupported version 1.6",
			version: "1.6",
			config: v1alpha1.LocalRateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sidecar-config",
					Namespace: "default",
				},
				Spec: v1alpha1.LocalRateLimitConfigSpec{
					Type: v1alpha1.Sidecar,
					Selector: v1alpha1.LocalRateLimitConfigSelector{
						Labels:       map[string]string{"app": "my-service"},
						IstioVersion: []string{"1.6"},
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
			expectedType:  "",
			expectedError: true,
			errorContains: "version not supported",
		},
		{
			name:    "unknown config type",
			version: "1.9",
			config: v1alpha1.LocalRateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "unknown-config",
					Namespace: "istio-system",
				},
				Spec: v1alpha1.LocalRateLimitConfigSpec{
					Type: "unknown",
					Selector: v1alpha1.LocalRateLimitConfigSelector{
						Labels:       map[string]string{"app": "test"},
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
					Config: "unknown-config",
					Selector: v1alpha1.LocalRateLimitSelector{
						VHost: "test.example.com:443",
					},
					Limit: &v1alpha1.LocalRateLimit_Limit{
						Unit:            "second",
						RequestsPerUnit: 10,
					},
				},
			},
			expectedType:  "",
			expectedError: true,
			errorContains: "version not supported",
		},
		{
			name:    "valid gateway config with latest version 1.26",
			version: "1.26",
			config: v1alpha1.LocalRateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "gateway-config",
					Namespace: "istio-system",
				},
				Spec: v1alpha1.LocalRateLimitConfigSpec{
					Type: v1alpha1.Gateway,
					Selector: v1alpha1.LocalRateLimitConfigSelector{
						Labels:       map[string]string{"app": "gateway"},
						IstioVersion: []string{"1.26"},
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
						Unit:            "day",
						RequestsPerUnit: 10000,
					},
				},
			},
			expectedType:  "*ratelimit.V3GatewayBuilder",
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory, err := ratelimit.GetConfigFactory(tt.version, tt.config, tt.ratelimit)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, factory)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, factory)
			}
		})
	}
}

func TestGetConfigFactory_BuildsValidEnvoyFilter(t *testing.T) {
	tests := []struct {
		name      string
		version   string
		config    v1alpha1.LocalRateLimitConfig
		ratelimit v1alpha1.LocalRateLimit
	}{
		{
			name:    "gateway factory builds valid envoy filter",
			version: "1.9",
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
		},
		{
			name:    "sidecar factory builds valid envoy filter",
			version: "1.10",
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
						Unit:            "minute",
						RequestsPerUnit: 100,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory, err := ratelimit.GetConfigFactory(tt.version, tt.config, tt.ratelimit)
			assert.NoError(t, err)
			assert.NotNil(t, factory)

			envoyFilter, err := factory.Build()
			assert.NoError(t, err)
			assert.NotNil(t, envoyFilter)
			assert.Equal(t, tt.ratelimit.Name+"-"+tt.version, envoyFilter.Name)
			assert.Equal(t, tt.ratelimit.Namespace, envoyFilter.Namespace)
			assert.Equal(t, "EnvoyFilter", envoyFilter.Kind)
			assert.Equal(t, "networking.istio.io/v1alpha3", envoyFilter.APIVersion)
		})
	}
}

func TestGetConfigFactory_VersionBoundary(t *testing.T) {
	config := v1alpha1.LocalRateLimitConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "gateway-config",
			Namespace: "istio-system",
		},
		Spec: v1alpha1.LocalRateLimitConfigSpec{
			Type: v1alpha1.Gateway,
			Selector: v1alpha1.LocalRateLimitConfigSelector{
				Labels:       map[string]string{"app": "gateway"},
				IstioVersion: []string{"1.7"},
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

	// Test version 1.7.0 (exactly at boundary >= 1.7.x)
	factory, err := ratelimit.GetConfigFactory("1.7.0", config, rl)
	assert.NoError(t, err)
	assert.NotNil(t, factory)

	// Test version 1.7.1 (above boundary)
	factory, err = ratelimit.GetConfigFactory("1.7.1", config, rl)
	assert.NoError(t, err)
	assert.NotNil(t, factory)

	// Test version 1.6.99 (below boundary)
	factory, err = ratelimit.GetConfigFactory("1.6.99", config, rl)
	assert.Error(t, err)
	assert.Nil(t, factory)
	assert.Contains(t, err.Error(), "version not supported")
}
