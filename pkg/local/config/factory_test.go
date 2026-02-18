package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/local/config"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNewConfigFactory(t *testing.T) {
	testCases := []struct {
		name          string
		version       string
		config        v1alpha1.LocalRateLimitConfig
		expectedError bool
		errorContains string
	}{
		{
			name:    "valid gateway config with version 1.7",
			version: "1.7",
			config: v1alpha1.LocalRateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-gateway",
					Namespace: "test-namespace",
				},
				Spec: v1alpha1.LocalRateLimitConfigSpec{
					Type: v1alpha1.Gateway,
					Selector: v1alpha1.LocalRateLimitConfigSelector{
						Labels: map[string]string{
							"app": "test-app",
						},
					},
				},
			},
			expectedError: false,
		},
		{
			name:    "valid gateway config with version 1.8",
			version: "1.8",
			config: v1alpha1.LocalRateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-gateway",
					Namespace: "test-namespace",
				},
				Spec: v1alpha1.LocalRateLimitConfigSpec{
					Type: v1alpha1.Gateway,
					Selector: v1alpha1.LocalRateLimitConfigSelector{
						Labels: map[string]string{
							"app": "test-app",
						},
					},
				},
			},
			expectedError: false,
		},
		{
			name:    "valid sidecar config with version 1.9",
			version: "1.9",
			config: v1alpha1.LocalRateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-sidecar",
					Namespace: "test-namespace",
				},
				Spec: v1alpha1.LocalRateLimitConfigSpec{
					Type: v1alpha1.Sidecar,
					Selector: v1alpha1.LocalRateLimitConfigSelector{
						Labels: map[string]string{
							"app": "test-app",
						},
					},
				},
			},
			expectedError: false,
		},
		{
			name:    "valid sidecar config with version 1.10",
			version: "1.10",
			config: v1alpha1.LocalRateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-sidecar",
					Namespace: "test-namespace",
				},
				Spec: v1alpha1.LocalRateLimitConfigSpec{
					Type: v1alpha1.Sidecar,
					Selector: v1alpha1.LocalRateLimitConfigSelector{
						Labels: map[string]string{
							"app": "test-app",
						},
					},
				},
			},
			expectedError: false,
		},
		{
			name:    "valid gateway config with recent version 1.20",
			version: "1.20",
			config: v1alpha1.LocalRateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-gateway",
					Namespace: "test-namespace",
				},
				Spec: v1alpha1.LocalRateLimitConfigSpec{
					Type: v1alpha1.Gateway,
					Selector: v1alpha1.LocalRateLimitConfigSelector{
						Labels: map[string]string{
							"app": "test-app",
						},
					},
				},
			},
			expectedError: false,
		},
		{
			name:    "invalid version format",
			version: "invalid-version",
			config: v1alpha1.LocalRateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-invalid",
					Namespace: "test-namespace",
				},
				Spec: v1alpha1.LocalRateLimitConfigSpec{
					Type: v1alpha1.Gateway,
					Selector: v1alpha1.LocalRateLimitConfigSelector{
						Labels: map[string]string{
							"app": "test-app",
						},
					},
				},
			},
			expectedError: true,
			errorContains: "cannot parse version",
		},
		{
			name:    "unsupported old version 1.5",
			version: "1.5",
			config: v1alpha1.LocalRateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-old",
					Namespace: "test-namespace",
				},
				Spec: v1alpha1.LocalRateLimitConfigSpec{
					Type: v1alpha1.Gateway,
					Selector: v1alpha1.LocalRateLimitConfigSelector{
						Labels: map[string]string{
							"app": "test-app",
						},
					},
				},
			},
			expectedError: true,
			errorContains: "version not supported",
		},
		{
			name:    "unsupported old version 1.6",
			version: "1.6",
			config: v1alpha1.LocalRateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-old",
					Namespace: "test-namespace",
				},
				Spec: v1alpha1.LocalRateLimitConfigSpec{
					Type: v1alpha1.Gateway,
					Selector: v1alpha1.LocalRateLimitConfigSelector{
						Labels: map[string]string{
							"app": "test-app",
						},
					},
				},
			},
			expectedError: true,
			errorContains: "version not supported",
		},
		{
			name:    "unsupported type returns error",
			version: "1.8",
			config: v1alpha1.LocalRateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-unknown-type",
					Namespace: "test-namespace",
				},
				Spec: v1alpha1.LocalRateLimitConfigSpec{
					Type: "unknown",
					Selector: v1alpha1.LocalRateLimitConfigSelector{
						Labels: map[string]string{
							"app": "test-app",
						},
					},
				},
			},
			expectedError: true,
			errorContains: "version not supported",
		},
		{
			name:    "empty version string",
			version: "",
			config: v1alpha1.LocalRateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-empty",
					Namespace: "test-namespace",
				},
				Spec: v1alpha1.LocalRateLimitConfigSpec{
					Type: v1alpha1.Gateway,
					Selector: v1alpha1.LocalRateLimitConfigSelector{
						Labels: map[string]string{
							"app": "test-app",
						},
					},
				},
			},
			expectedError: true,
			errorContains: "cannot parse version",
		},
		{
			name:    "version with patch number",
			version: "1.8.5",
			config: v1alpha1.LocalRateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-patch",
					Namespace: "test-namespace",
				},
				Spec: v1alpha1.LocalRateLimitConfigSpec{
					Type: v1alpha1.Gateway,
					Selector: v1alpha1.LocalRateLimitConfigSelector{
						Labels: map[string]string{
							"app": "test-app",
						},
					},
				},
			},
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			factory, err := config.NewConfigFactory(tc.version, tc.config)

			if tc.expectedError {
				assert.Error(t, err)
				assert.Nil(t, factory)
				if tc.errorContains != "" {
					assert.Contains(t, err.Error(), tc.errorContains)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, factory)
			}
		})
	}
}

func TestConfigFactory_Build(t *testing.T) {
	testCases := []struct {
		name               string
		version            string
		config             v1alpha1.LocalRateLimitConfig
		expectedError      bool
		expectedName       string
		expectedNamespace  string
		expectedLabelCount int
	}{
		{
			name:    "build gateway envoy filter",
			version: "1.8",
			config: v1alpha1.LocalRateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "gateway-test",
					Namespace: "istio-system",
				},
				Spec: v1alpha1.LocalRateLimitConfigSpec{
					Type: v1alpha1.Gateway,
					Selector: v1alpha1.LocalRateLimitConfigSelector{
						Labels: map[string]string{
							"app":     "istio-ingressgateway",
							"version": "v1",
						},
					},
				},
			},
			expectedError:      false,
			expectedName:       "gateway-test-1.8",
			expectedNamespace:  "istio-system",
			expectedLabelCount: 2,
		},
		{
			name:    "build sidecar envoy filter",
			version: "1.9",
			config: v1alpha1.LocalRateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sidecar-test",
					Namespace: "default",
				},
				Spec: v1alpha1.LocalRateLimitConfigSpec{
					Type: v1alpha1.Sidecar,
					Selector: v1alpha1.LocalRateLimitConfigSelector{
						Labels: map[string]string{
							"app": "my-service",
						},
					},
				},
			},
			expectedError:      false,
			expectedName:       "sidecar-test-1.9",
			expectedNamespace:  "default",
			expectedLabelCount: 1,
		},
		{
			name:    "build gateway with latest version",
			version: "1.20",
			config: v1alpha1.LocalRateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "latest-gateway",
					Namespace: "production",
				},
				Spec: v1alpha1.LocalRateLimitConfigSpec{
					Type: v1alpha1.Gateway,
					Selector: v1alpha1.LocalRateLimitConfigSelector{
						Labels: map[string]string{
							"app": "gateway",
						},
					},
				},
			},
			expectedError:      false,
			expectedName:       "latest-gateway-1.20",
			expectedNamespace:  "production",
			expectedLabelCount: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			factory, err := config.NewConfigFactory(tc.version, tc.config)
			assert.NoError(t, err)
			assert.NotNil(t, factory)

			envoyFilter, err := factory.Build()

			if tc.expectedError {
				assert.Error(t, err)
				assert.Nil(t, envoyFilter)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, envoyFilter)
				assert.Equal(t, tc.expectedName, envoyFilter.Name)
				assert.Equal(t, tc.expectedNamespace, envoyFilter.Namespace)
				assert.Equal(t, "EnvoyFilter", envoyFilter.Kind)
				assert.Equal(t, "networking.istio.io/v1alpha3", envoyFilter.APIVersion)
				assert.Len(t, envoyFilter.Spec.WorkloadSelector.Labels, tc.expectedLabelCount)
				assert.Equal(t, tc.version, envoyFilter.Labels["istio/version"])
			}
		})
	}
}

func TestConfigFactory_VersionBoundary(t *testing.T) {
	baseConfig := v1alpha1.LocalRateLimitConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "version-boundary-test",
			Namespace: "test-namespace",
		},
		Spec: v1alpha1.LocalRateLimitConfigSpec{
			Type: v1alpha1.Gateway,
			Selector: v1alpha1.LocalRateLimitConfigSelector{
				Labels: map[string]string{
					"app": "test-app",
				},
			},
		},
	}

	// Test version 1.7 (boundary - should work, >= 1.7.x)
	factory, err := config.NewConfigFactory("1.7", baseConfig)
	assert.NoError(t, err)
	assert.NotNil(t, factory)

	// Test version 1.6 (just below boundary - should fail)
	factory, err = config.NewConfigFactory("1.6", baseConfig)
	assert.Error(t, err)
	assert.Nil(t, factory)
	assert.Contains(t, err.Error(), "version not supported")
}
