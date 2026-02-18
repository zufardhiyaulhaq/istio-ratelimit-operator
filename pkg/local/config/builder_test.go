package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/local/config"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNewConfigBuilder(t *testing.T) {
	builder := config.NewConfigBuilder()
	assert.NotNil(t, builder)
}

func TestConfigBuilder_SetConfig(t *testing.T) {
	testCases := []struct {
		name   string
		config v1alpha1.LocalRateLimitConfig
	}{
		{
			name: "set gateway config",
			config: v1alpha1.LocalRateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-gateway",
					Namespace: "test-namespace",
				},
				Spec: v1alpha1.LocalRateLimitConfigSpec{
					Type: v1alpha1.Gateway,
					Selector: v1alpha1.LocalRateLimitConfigSelector{
						IstioVersion: []string{"1.8"},
						Labels: map[string]string{
							"app": "test-app",
						},
					},
				},
			},
		},
		{
			name: "set sidecar config",
			config: v1alpha1.LocalRateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-sidecar",
					Namespace: "test-namespace",
				},
				Spec: v1alpha1.LocalRateLimitConfigSpec{
					Type: v1alpha1.Sidecar,
					Selector: v1alpha1.LocalRateLimitConfigSelector{
						IstioVersion: []string{"1.9"},
						Labels: map[string]string{
							"app": "test-sidecar-app",
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := config.NewConfigBuilder()
			result := builder.SetConfig(tc.config)

			// SetConfig should return the builder for chaining
			assert.Same(t, builder, result)
			// Config should be set
			assert.Equal(t, tc.config.Name, builder.Config.Name)
			assert.Equal(t, tc.config.Namespace, builder.Config.Namespace)
			assert.Equal(t, tc.config.Spec.Type, builder.Config.Spec.Type)
		})
	}
}

func TestConfigBuilder_Build(t *testing.T) {
	testCases := []struct {
		name                   string
		config                 v1alpha1.LocalRateLimitConfig
		expectedError          bool
		expectedEnvoyFilterCnt int
	}{
		{
			name: "build gateway config with single version",
			config: v1alpha1.LocalRateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-gateway",
					Namespace: "test-namespace",
				},
				Spec: v1alpha1.LocalRateLimitConfigSpec{
					Type: v1alpha1.Gateway,
					Selector: v1alpha1.LocalRateLimitConfigSelector{
						IstioVersion: []string{"1.8"},
						Labels: map[string]string{
							"app": "test-app",
						},
					},
				},
			},
			expectedError:          false,
			expectedEnvoyFilterCnt: 1,
		},
		{
			name: "build sidecar config with single version",
			config: v1alpha1.LocalRateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-sidecar",
					Namespace: "test-namespace",
				},
				Spec: v1alpha1.LocalRateLimitConfigSpec{
					Type: v1alpha1.Sidecar,
					Selector: v1alpha1.LocalRateLimitConfigSelector{
						IstioVersion: []string{"1.9"},
						Labels: map[string]string{
							"app": "test-app",
						},
					},
				},
			},
			expectedError:          false,
			expectedEnvoyFilterCnt: 1,
		},
		{
			name: "build gateway config with multiple versions",
			config: v1alpha1.LocalRateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-gateway-multi",
					Namespace: "test-namespace",
				},
				Spec: v1alpha1.LocalRateLimitConfigSpec{
					Type: v1alpha1.Gateway,
					Selector: v1alpha1.LocalRateLimitConfigSelector{
						IstioVersion: []string{"1.8", "1.9", "1.10"},
						Labels: map[string]string{
							"app": "test-app",
						},
					},
				},
			},
			expectedError:          false,
			expectedEnvoyFilterCnt: 3,
		},
		{
			name: "build sidecar config with multiple versions",
			config: v1alpha1.LocalRateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-sidecar-multi",
					Namespace: "test-namespace",
				},
				Spec: v1alpha1.LocalRateLimitConfigSpec{
					Type: v1alpha1.Sidecar,
					Selector: v1alpha1.LocalRateLimitConfigSelector{
						IstioVersion: []string{"1.10", "1.11"},
						Labels: map[string]string{
							"app": "test-app",
						},
					},
				},
			},
			expectedError:          false,
			expectedEnvoyFilterCnt: 2,
		},
		{
			name: "build config with invalid version",
			config: v1alpha1.LocalRateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-invalid",
					Namespace: "test-namespace",
				},
				Spec: v1alpha1.LocalRateLimitConfigSpec{
					Type: v1alpha1.Gateway,
					Selector: v1alpha1.LocalRateLimitConfigSelector{
						IstioVersion: []string{"invalid-version"},
						Labels: map[string]string{
							"app": "test-app",
						},
					},
				},
			},
			expectedError:          true,
			expectedEnvoyFilterCnt: 0,
		},
		{
			name: "build config with unsupported old version",
			config: v1alpha1.LocalRateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-old-version",
					Namespace: "test-namespace",
				},
				Spec: v1alpha1.LocalRateLimitConfigSpec{
					Type: v1alpha1.Gateway,
					Selector: v1alpha1.LocalRateLimitConfigSelector{
						IstioVersion: []string{"1.5"},
						Labels: map[string]string{
							"app": "test-app",
						},
					},
				},
			},
			expectedError:          true,
			expectedEnvoyFilterCnt: 0,
		},
		{
			name: "build config with empty versions",
			config: v1alpha1.LocalRateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-empty-versions",
					Namespace: "test-namespace",
				},
				Spec: v1alpha1.LocalRateLimitConfigSpec{
					Type: v1alpha1.Gateway,
					Selector: v1alpha1.LocalRateLimitConfigSelector{
						IstioVersion: []string{},
						Labels: map[string]string{
							"app": "test-app",
						},
					},
				},
			},
			expectedError:          false,
			expectedEnvoyFilterCnt: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := config.NewConfigBuilder().SetConfig(tc.config)
			envoyFilters, err := builder.Build()

			if tc.expectedError {
				assert.Error(t, err)
				assert.Nil(t, envoyFilters)
			} else {
				assert.NoError(t, err)
				assert.Len(t, envoyFilters, tc.expectedEnvoyFilterCnt)

				// Verify each envoy filter has correct naming
				for i, ef := range envoyFilters {
					expectedName := tc.config.Name + "-" + tc.config.Spec.Selector.IstioVersion[i]
					assert.Equal(t, expectedName, ef.Name)
					assert.Equal(t, tc.config.Namespace, ef.Namespace)
					assert.Equal(t, "EnvoyFilter", ef.Kind)
					assert.Equal(t, "networking.istio.io/v1alpha3", ef.APIVersion)
				}
			}
		})
	}
}

func TestConfigBuilder_BuildChaining(t *testing.T) {
	cfg := v1alpha1.LocalRateLimitConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-chaining",
			Namespace: "test-namespace",
		},
		Spec: v1alpha1.LocalRateLimitConfigSpec{
			Type: v1alpha1.Gateway,
			Selector: v1alpha1.LocalRateLimitConfigSelector{
				IstioVersion: []string{"1.8"},
				Labels: map[string]string{
					"app": "test-app",
				},
			},
		},
	}

	// Test fluent API chaining
	envoyFilters, err := config.NewConfigBuilder().
		SetConfig(cfg).
		Build()

	assert.NoError(t, err)
	assert.Len(t, envoyFilters, 1)
	assert.Equal(t, "test-chaining-1.8", envoyFilters[0].Name)
}
