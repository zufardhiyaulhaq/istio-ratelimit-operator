package ratelimit_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/global/ratelimit"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGetConfigFactory(t *testing.T) {
	tests := []struct {
		name          string
		version       string
		config        v1alpha1.GlobalRateLimitConfig
		ratelimit     v1alpha1.GlobalRateLimit
		expectedError bool
		errorContains string
	}{
		{
			name:    "valid gateway config with version 1.9",
			version: "1.9",
			config: v1alpha1.GlobalRateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-config",
					Namespace: "istio-system",
				},
				Spec: v1alpha1.GlobalRateLimitConfigSpec{
					Type: v1alpha1.Gateway,
					Selector: v1alpha1.GlobalRateLimitConfigSelector{
						Labels:       map[string]string{"app": "gateway"},
						IstioVersion: []string{"1.9"},
					},
				},
			},
			ratelimit: v1alpha1.GlobalRateLimit{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ratelimit",
					Namespace: "istio-system",
				},
				Spec: v1alpha1.GlobalRateLimitSpec{
					Config: "test-config",
					Selector: v1alpha1.GlobalRateLimitSelector{
						VHost: "example.com:443",
					},
				},
			},
			expectedError: false,
		},
		{
			name:    "valid sidecar config with version 1.10",
			version: "1.10",
			config: v1alpha1.GlobalRateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-config",
					Namespace: "istio-system",
				},
				Spec: v1alpha1.GlobalRateLimitConfigSpec{
					Type: v1alpha1.Sidecar,
					Selector: v1alpha1.GlobalRateLimitConfigSelector{
						Labels:       map[string]string{"app": "myapp"},
						IstioVersion: []string{"1.10"},
					},
				},
			},
			ratelimit: v1alpha1.GlobalRateLimit{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ratelimit",
					Namespace: "default",
				},
				Spec: v1alpha1.GlobalRateLimitSpec{
					Config: "test-config",
					Selector: v1alpha1.GlobalRateLimitSelector{
						VHost: "inbound|8080||",
					},
				},
			},
			expectedError: false,
		},
		{
			name:    "valid gateway config with version 1.7",
			version: "1.7",
			config: v1alpha1.GlobalRateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-config",
					Namespace: "istio-system",
				},
				Spec: v1alpha1.GlobalRateLimitConfigSpec{
					Type: v1alpha1.Gateway,
					Selector: v1alpha1.GlobalRateLimitConfigSelector{
						Labels:       map[string]string{"app": "gateway"},
						IstioVersion: []string{"1.7"},
					},
				},
			},
			ratelimit: v1alpha1.GlobalRateLimit{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ratelimit",
					Namespace: "istio-system",
				},
				Spec: v1alpha1.GlobalRateLimitSpec{
					Config: "test-config",
					Selector: v1alpha1.GlobalRateLimitSelector{
						VHost: "example.com:443",
					},
				},
			},
			expectedError: false,
		},
		{
			name:    "invalid version format",
			version: "invalid-version",
			config: v1alpha1.GlobalRateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-config",
					Namespace: "istio-system",
				},
				Spec: v1alpha1.GlobalRateLimitConfigSpec{
					Type: v1alpha1.Gateway,
				},
			},
			ratelimit: v1alpha1.GlobalRateLimit{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ratelimit",
					Namespace: "istio-system",
				},
			},
			expectedError: true,
			errorContains: "cannot parse version",
		},
		{
			name:    "unsupported version 1.6",
			version: "1.6",
			config: v1alpha1.GlobalRateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-config",
					Namespace: "istio-system",
				},
				Spec: v1alpha1.GlobalRateLimitConfigSpec{
					Type: v1alpha1.Gateway,
				},
			},
			ratelimit: v1alpha1.GlobalRateLimit{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ratelimit",
					Namespace: "istio-system",
				},
			},
			expectedError: true,
			errorContains: "version not supported",
		},
		{
			name:    "unsupported version 1.5",
			version: "1.5",
			config: v1alpha1.GlobalRateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-config",
					Namespace: "istio-system",
				},
				Spec: v1alpha1.GlobalRateLimitConfigSpec{
					Type: v1alpha1.Gateway,
				},
			},
			ratelimit: v1alpha1.GlobalRateLimit{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ratelimit",
					Namespace: "istio-system",
				},
			},
			expectedError: true,
			errorContains: "version not supported",
		},
		{
			name:    "valid version but unknown type returns error",
			version: "1.9",
			config: v1alpha1.GlobalRateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-config",
					Namespace: "istio-system",
				},
				Spec: v1alpha1.GlobalRateLimitConfigSpec{
					Type: "unknown",
					Selector: v1alpha1.GlobalRateLimitConfigSelector{
						Labels:       map[string]string{"app": "gateway"},
						IstioVersion: []string{"1.9"},
					},
				},
			},
			ratelimit: v1alpha1.GlobalRateLimit{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ratelimit",
					Namespace: "istio-system",
				},
			},
			expectedError: true,
			errorContains: "version not supported",
		},
		{
			name:    "valid sidecar config with version 1.25",
			version: "1.25",
			config: v1alpha1.GlobalRateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-config",
					Namespace: "istio-system",
				},
				Spec: v1alpha1.GlobalRateLimitConfigSpec{
					Type: v1alpha1.Sidecar,
					Selector: v1alpha1.GlobalRateLimitConfigSelector{
						Labels:       map[string]string{"app": "myapp"},
						IstioVersion: []string{"1.25"},
					},
				},
			},
			ratelimit: v1alpha1.GlobalRateLimit{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ratelimit",
					Namespace: "default",
				},
				Spec: v1alpha1.GlobalRateLimitSpec{
					Config: "test-config",
					Selector: v1alpha1.GlobalRateLimitSelector{
						VHost: "inbound|8080||",
					},
				},
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory, err := ratelimit.GetConfigFactory(tt.version, tt.config, tt.ratelimit)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
				assert.Nil(t, factory)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, factory)

				// Verify that the factory can build an EnvoyFilter
				envoyFilter, err := factory.Build()
				assert.NoError(t, err)
				assert.NotNil(t, envoyFilter)
			}
		})
	}
}
