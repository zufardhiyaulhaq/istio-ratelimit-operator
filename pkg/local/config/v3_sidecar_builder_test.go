package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/local/config"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/utils"

	networking "istio.io/api/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type V3SidecarBuilderTestCase struct {
	name          string
	config        v1alpha1.LocalRateLimitConfig
	version       string
	expectedError bool
}

var V3SidecarBuilderTestGrid = []V3SidecarBuilderTestCase{
	{
		name: "given correct ratelimit",
		config: v1alpha1.LocalRateLimitConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "foo",
			},
			Spec: v1alpha1.LocalRateLimitConfigSpec{
				Type: "sidecar",
				Selector: v1alpha1.LocalRateLimitConfigSelector{
					IstioVersion: []string{"1.8"},
					Labels: map[string]string{
						"app": "foo",
					},
				},
			},
		},
		version:       "1.8",
		expectedError: false,
	},
	{
		name: "sidecar with multiple labels",
		config: v1alpha1.LocalRateLimitConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "multi-label-sidecar",
				Namespace: "default",
			},
			Spec: v1alpha1.LocalRateLimitConfigSpec{
				Type: v1alpha1.Sidecar,
				Selector: v1alpha1.LocalRateLimitConfigSelector{
					IstioVersion: []string{"1.9"},
					Labels: map[string]string{
						"app":     "my-service",
						"version": "v1",
						"env":     "production",
					},
				},
			},
		},
		version:       "1.9",
		expectedError: false,
	},
	{
		name: "sidecar with different namespace",
		config: v1alpha1.LocalRateLimitConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "production-sidecar",
				Namespace: "production",
			},
			Spec: v1alpha1.LocalRateLimitConfigSpec{
				Type: v1alpha1.Sidecar,
				Selector: v1alpha1.LocalRateLimitConfigSelector{
					IstioVersion: []string{"1.10"},
					Labels: map[string]string{
						"app": "api-service",
					},
				},
			},
		},
		version:       "1.10",
		expectedError: false,
	},
	{
		name: "sidecar with latest version 1.20",
		config: v1alpha1.LocalRateLimitConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "latest-sidecar",
				Namespace: "default",
			},
			Spec: v1alpha1.LocalRateLimitConfigSpec{
				Type: v1alpha1.Sidecar,
				Selector: v1alpha1.LocalRateLimitConfigSelector{
					IstioVersion: []string{"1.20"},
					Labels: map[string]string{
						"app": "service",
					},
				},
			},
		},
		version:       "1.20",
		expectedError: false,
	},
}

func TestNewV3SidecarBuilder(t *testing.T) {
	for _, test := range V3SidecarBuilderTestGrid {
		t.Run(test.name, func(t *testing.T) {
			envoyfilter, err := config.NewV3SidecarBuilder(test.config, test.version).
				Build()

			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.config.Name+"-"+test.version, envoyfilter.Name)
				assert.Equal(t, test.config.Namespace, envoyfilter.Namespace)
			}
		})
	}
}

func TestV3SidecarBuilder_EnvoyFilterStructure(t *testing.T) {
	cfg := v1alpha1.LocalRateLimitConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "structure-test",
			Namespace: "test-namespace",
		},
		Spec: v1alpha1.LocalRateLimitConfigSpec{
			Type: v1alpha1.Sidecar,
			Selector: v1alpha1.LocalRateLimitConfigSelector{
				IstioVersion: []string{"1.8"},
				Labels: map[string]string{
					"app": "test-app",
				},
			},
		},
	}

	envoyfilter, err := config.NewV3SidecarBuilder(cfg, "1.8").Build()
	assert.NoError(t, err)

	// Verify TypeMeta
	assert.Equal(t, "EnvoyFilter", envoyfilter.Kind)
	assert.Equal(t, "networking.istio.io/v1alpha3", envoyfilter.APIVersion)

	// Verify ObjectMeta
	assert.Equal(t, "structure-test-1.8", envoyfilter.Name)
	assert.Equal(t, "test-namespace", envoyfilter.Namespace)
	assert.Equal(t, "1.8", envoyfilter.Labels["istio/version"])

	// Verify WorkloadSelector
	assert.NotNil(t, envoyfilter.Spec.WorkloadSelector)
	assert.Equal(t, "test-app", envoyfilter.Spec.WorkloadSelector.Labels["app"])

	// Verify ConfigPatches
	assert.Len(t, envoyfilter.Spec.ConfigPatches, 1)
	patch := envoyfilter.Spec.ConfigPatches[0]
	assert.Equal(t, networking.EnvoyFilter_HTTP_FILTER, patch.ApplyTo)
	assert.Equal(t, networking.EnvoyFilter_Patch_INSERT_BEFORE, patch.Patch.Operation)

	// Verify Match context - should be SIDECAR_INBOUND for sidecar
	assert.Equal(t, networking.EnvoyFilter_SIDECAR_INBOUND, patch.Match.Context)
}

func TestV3SidecarBuilder_ProxyMatch(t *testing.T) {
	versions := []string{"1.8", "1.9", "1.10", "1.20"}

	for _, version := range versions {
		t.Run("version_"+version, func(t *testing.T) {
			cfg := v1alpha1.LocalRateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "proxy-match-test",
					Namespace: "test-namespace",
				},
				Spec: v1alpha1.LocalRateLimitConfigSpec{
					Type: v1alpha1.Sidecar,
					Selector: v1alpha1.LocalRateLimitConfigSelector{
						IstioVersion: []string{version},
						Labels: map[string]string{
							"app": "test-app",
						},
					},
				},
			}

			envoyfilter, err := config.NewV3SidecarBuilder(cfg, version).Build()
			assert.NoError(t, err)

			patch := envoyfilter.Spec.ConfigPatches[0]
			assert.NotNil(t, patch.Match.Proxy)
			assert.Equal(t, utils.WellKnownVersions[version], patch.Match.Proxy.ProxyVersion)
		})
	}
}

func TestV3SidecarBuilder_FilterChainMatch(t *testing.T) {
	cfg := v1alpha1.LocalRateLimitConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "filter-chain-test",
			Namespace: "test-namespace",
		},
		Spec: v1alpha1.LocalRateLimitConfigSpec{
			Type: v1alpha1.Sidecar,
			Selector: v1alpha1.LocalRateLimitConfigSelector{
				IstioVersion: []string{"1.8"},
				Labels: map[string]string{
					"app": "test-app",
				},
			},
		},
	}

	envoyfilter, err := config.NewV3SidecarBuilder(cfg, "1.8").Build()
	assert.NoError(t, err)

	patch := envoyfilter.Spec.ConfigPatches[0]
	listener := patch.Match.GetListener()
	assert.NotNil(t, listener)
	assert.NotNil(t, listener.FilterChain)
	assert.NotNil(t, listener.FilterChain.Filter)
	assert.Equal(t, "envoy.filters.network.http_connection_manager", listener.FilterChain.Filter.Name)
	assert.NotNil(t, listener.FilterChain.Filter.SubFilter)
	assert.Equal(t, "envoy.filters.http.router", listener.FilterChain.Filter.SubFilter.Name)
}

func TestV3SidecarBuilder_PatchValue(t *testing.T) {
	cfg := v1alpha1.LocalRateLimitConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "patch-value-test",
			Namespace: "test-namespace",
		},
		Spec: v1alpha1.LocalRateLimitConfigSpec{
			Type: v1alpha1.Sidecar,
			Selector: v1alpha1.LocalRateLimitConfigSelector{
				IstioVersion: []string{"1.8"},
				Labels: map[string]string{
					"app": "test-app",
				},
			},
		},
	}

	envoyfilter, err := config.NewV3SidecarBuilder(cfg, "1.8").Build()
	assert.NoError(t, err)

	patch := envoyfilter.Spec.ConfigPatches[0]
	assert.NotNil(t, patch.Patch.Value)
}

func TestV3SidecarBuilder_DiffersFromGateway(t *testing.T) {
	cfg := v1alpha1.LocalRateLimitConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "compare-test",
			Namespace: "test-namespace",
		},
		Spec: v1alpha1.LocalRateLimitConfigSpec{
			Type: v1alpha1.Sidecar,
			Selector: v1alpha1.LocalRateLimitConfigSelector{
				IstioVersion: []string{"1.8"},
				Labels: map[string]string{
					"app": "test-app",
				},
			},
		},
	}

	sidecarFilter, err := config.NewV3SidecarBuilder(cfg, "1.8").Build()
	assert.NoError(t, err)

	// Change type to gateway for comparison
	cfg.Spec.Type = v1alpha1.Gateway
	gatewayFilter, err := config.NewV3GatewayBuilder(cfg, "1.8").Build()
	assert.NoError(t, err)

	// Both should have the same name and namespace
	assert.Equal(t, sidecarFilter.Name, gatewayFilter.Name)
	assert.Equal(t, sidecarFilter.Namespace, gatewayFilter.Namespace)

	// But different contexts
	sidecarPatch := sidecarFilter.Spec.ConfigPatches[0]
	gatewayPatch := gatewayFilter.Spec.ConfigPatches[0]
	assert.Equal(t, networking.EnvoyFilter_SIDECAR_INBOUND, sidecarPatch.Match.Context)
	assert.Equal(t, networking.EnvoyFilter_GATEWAY, gatewayPatch.Match.Context)
}

func TestV3SidecarBuilder_WorkloadSelector(t *testing.T) {
	testCases := []struct {
		name           string
		labels         map[string]string
		expectedLabels map[string]string
	}{
		{
			name: "single label",
			labels: map[string]string{
				"app": "my-service",
			},
			expectedLabels: map[string]string{
				"app": "my-service",
			},
		},
		{
			name: "multiple labels",
			labels: map[string]string{
				"app":     "my-service",
				"version": "v2",
				"tier":    "backend",
			},
			expectedLabels: map[string]string{
				"app":     "my-service",
				"version": "v2",
				"tier":    "backend",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := v1alpha1.LocalRateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "workload-selector-test",
					Namespace: "test-namespace",
				},
				Spec: v1alpha1.LocalRateLimitConfigSpec{
					Type: v1alpha1.Sidecar,
					Selector: v1alpha1.LocalRateLimitConfigSelector{
						IstioVersion: []string{"1.8"},
						Labels:       tc.labels,
					},
				},
			}

			envoyfilter, err := config.NewV3SidecarBuilder(cfg, "1.8").Build()
			assert.NoError(t, err)

			assert.Equal(t, tc.expectedLabels, envoyfilter.Spec.WorkloadSelector.Labels)
		})
	}
}
