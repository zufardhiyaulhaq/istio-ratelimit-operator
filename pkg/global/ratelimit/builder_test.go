package ratelimit_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/global/ratelimit"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNewConfigBuilder(t *testing.T) {
	builder := ratelimit.NewConfigBuilder()
	assert.NotNil(t, builder)
}

func TestConfigBuilder_SetRateLimit(t *testing.T) {
	builder := ratelimit.NewConfigBuilder()

	rateLimitObj := v1alpha1.GlobalRateLimit{
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
	}

	result := builder.SetRateLimit(rateLimitObj)

	// Verify fluent interface returns the builder
	assert.Equal(t, builder, result)
	// Verify the ratelimit is set
	assert.Equal(t, rateLimitObj.Name, builder.RateLimit.Name)
	assert.Equal(t, rateLimitObj.Namespace, builder.RateLimit.Namespace)
	assert.Equal(t, rateLimitObj.Spec.Config, builder.RateLimit.Spec.Config)
}

func TestConfigBuilder_SetConfig(t *testing.T) {
	builder := ratelimit.NewConfigBuilder()

	config := v1alpha1.GlobalRateLimitConfig{
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
	}

	result := builder.SetConfig(config)

	// Verify fluent interface returns the builder
	assert.Equal(t, builder, result)
	// Verify the config is set
	assert.Equal(t, config.Name, builder.Config.Name)
	assert.Equal(t, config.Spec.Type, builder.Config.Spec.Type)
}

func TestConfigBuilder_SetVersions(t *testing.T) {
	builder := ratelimit.NewConfigBuilder()

	versions := []string{"1.9", "1.10", "1.11"}

	result := builder.SetVersions(versions)

	// Verify fluent interface returns the builder
	assert.Equal(t, builder, result)
	// Verify the versions are set
	assert.Equal(t, versions, builder.Versions)
}

func TestConfigBuilder_SetLabels(t *testing.T) {
	builder := ratelimit.NewConfigBuilder()

	labels := map[string]string{
		"app":     "test",
		"version": "v1",
	}

	result := builder.SetLabels(labels)

	// Verify fluent interface returns the builder
	assert.Equal(t, builder, result)
	// Verify the labels are set
	assert.Equal(t, labels, builder.Labels)
}

func TestConfigBuilder_FluentInterface(t *testing.T) {
	rateLimitObj := v1alpha1.GlobalRateLimit{
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
	}

	config := v1alpha1.GlobalRateLimitConfig{
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
	}

	versions := []string{"1.9", "1.10"}
	labels := map[string]string{"app": "test"}

	// Test chaining all setters
	builder := ratelimit.NewConfigBuilder().
		SetRateLimit(rateLimitObj).
		SetConfig(config).
		SetVersions(versions).
		SetLabels(labels)

	assert.Equal(t, rateLimitObj.Name, builder.RateLimit.Name)
	assert.Equal(t, config.Name, builder.Config.Name)
	assert.Equal(t, versions, builder.Versions)
	assert.Equal(t, labels, builder.Labels)
}

func TestConfigBuilder_Build(t *testing.T) {
	tests := []struct {
		name           string
		config         v1alpha1.GlobalRateLimitConfig
		ratelimit      v1alpha1.GlobalRateLimit
		expectedError  bool
		expectedCount  int
		expectedNames  []string
	}{
		{
			name: "build single gateway envoy filter",
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
					Matcher: []*v1alpha1.GlobalRateLimit_Action{
						{
							RemoteAddress: &v1alpha1.GlobalRateLimit_Action_RemoteAddress{},
						},
					},
				},
			},
			expectedError: false,
			expectedCount: 1,
			expectedNames: []string{"test-ratelimit-1.9"},
		},
		{
			name: "build multiple gateway envoy filters for multiple versions",
			config: v1alpha1.GlobalRateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-config",
					Namespace: "istio-system",
				},
				Spec: v1alpha1.GlobalRateLimitConfigSpec{
					Type: v1alpha1.Gateway,
					Selector: v1alpha1.GlobalRateLimitConfigSelector{
						Labels:       map[string]string{"app": "gateway"},
						IstioVersion: []string{"1.9", "1.10", "1.11"},
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
					Matcher: []*v1alpha1.GlobalRateLimit_Action{
						{
							RemoteAddress: &v1alpha1.GlobalRateLimit_Action_RemoteAddress{},
						},
					},
				},
			},
			expectedError: false,
			expectedCount: 3,
			expectedNames: []string{"test-ratelimit-1.9", "test-ratelimit-1.10", "test-ratelimit-1.11"},
		},
		{
			name: "build sidecar envoy filter",
			config: v1alpha1.GlobalRateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-config",
					Namespace: "default",
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
					Matcher: []*v1alpha1.GlobalRateLimit_Action{
						{
							RemoteAddress: &v1alpha1.GlobalRateLimit_Action_RemoteAddress{},
						},
					},
				},
			},
			expectedError: false,
			expectedCount: 1,
			expectedNames: []string{"test-ratelimit-1.10"},
		},
		{
			name: "error with invalid version",
			config: v1alpha1.GlobalRateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-config",
					Namespace: "istio-system",
				},
				Spec: v1alpha1.GlobalRateLimitConfigSpec{
					Type: v1alpha1.Gateway,
					Selector: v1alpha1.GlobalRateLimitConfigSelector{
						Labels:       map[string]string{"app": "gateway"},
						IstioVersion: []string{"invalid"},
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
			expectedError: true,
		},
		{
			name: "error with unsupported version 1.6",
			config: v1alpha1.GlobalRateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-config",
					Namespace: "istio-system",
				},
				Spec: v1alpha1.GlobalRateLimitConfigSpec{
					Type: v1alpha1.Gateway,
					Selector: v1alpha1.GlobalRateLimitConfigSelector{
						Labels:       map[string]string{"app": "gateway"},
						IstioVersion: []string{"1.6"},
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
			expectedError: true,
		},
		{
			name: "empty versions results in empty slice",
			config: v1alpha1.GlobalRateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-config",
					Namespace: "istio-system",
				},
				Spec: v1alpha1.GlobalRateLimitConfigSpec{
					Type: v1alpha1.Gateway,
					Selector: v1alpha1.GlobalRateLimitConfigSelector{
						Labels:       map[string]string{"app": "gateway"},
						IstioVersion: []string{},
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
			expectedCount: 0,
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
				assert.Nil(t, envoyFilters)
			} else {
				assert.NoError(t, err)
				assert.Len(t, envoyFilters, tt.expectedCount)

				// Verify each envoy filter has the expected name
				for i, expectedName := range tt.expectedNames {
					assert.Equal(t, expectedName, envoyFilters[i].Name)
					assert.Equal(t, tt.ratelimit.Namespace, envoyFilters[i].Namespace)
				}
			}
		})
	}
}

func TestConfigBuilder_Build_VerifyEnvoyFilterContent(t *testing.T) {
	config := v1alpha1.GlobalRateLimitConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-config",
			Namespace: "istio-system",
		},
		Spec: v1alpha1.GlobalRateLimitConfigSpec{
			Type: v1alpha1.Gateway,
			Selector: v1alpha1.GlobalRateLimitConfigSelector{
				Labels:       map[string]string{"app": "gateway", "istio": "ingressgateway"},
				IstioVersion: []string{"1.9"},
			},
		},
	}

	rateLimitObj := v1alpha1.GlobalRateLimit{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-ratelimit",
			Namespace: "istio-system",
		},
		Spec: v1alpha1.GlobalRateLimitSpec{
			Config: "test-config",
			Selector: v1alpha1.GlobalRateLimitSelector{
				VHost: "example.com:443",
			},
			Matcher: []*v1alpha1.GlobalRateLimit_Action{
				{
					RequestHeaders: &v1alpha1.GlobalRateLimit_Action_RequestHeaders{
						HeaderName:    ":method",
						DescriptorKey: "method",
					},
				},
			},
		},
	}

	builder := ratelimit.NewConfigBuilder().
		SetConfig(config).
		SetRateLimit(rateLimitObj)

	envoyFilters, err := builder.Build()

	assert.NoError(t, err)
	assert.Len(t, envoyFilters, 1)

	envoyFilter := envoyFilters[0]

	// Verify metadata
	assert.Equal(t, "test-ratelimit-1.9", envoyFilter.Name)
	assert.Equal(t, "istio-system", envoyFilter.Namespace)
	assert.Equal(t, "EnvoyFilter", envoyFilter.Kind)
	assert.Equal(t, "networking.istio.io/v1alpha3", envoyFilter.APIVersion)

	// Verify workload selector labels
	assert.Equal(t, config.Spec.Selector.Labels, envoyFilter.Spec.WorkloadSelector.Labels)

	// Verify config patches exist
	assert.Len(t, envoyFilter.Spec.ConfigPatches, 1)
}
