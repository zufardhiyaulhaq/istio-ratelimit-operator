package ratelimit_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/local/ratelimit"

	networking "istio.io/api/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type V3SidecarBuilderTestCase struct {
	name          string
	config        v1alpha1.LocalRateLimitConfig
	ratelimit     v1alpha1.LocalRateLimit
	version       string
	expectedError bool
}

var V3SidecarBuilderTestGrid = []V3SidecarBuilderTestCase{
	{
		name: "given correct ratelimit",
		config: v1alpha1.LocalRateLimitConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "sidecar-config",
				Namespace: "istio-system",
			},
			Spec: v1alpha1.LocalRateLimitConfigSpec{
				Type: "sidecar",
				Selector: v1alpha1.LocalRateLimitConfigSelector{
					IstioVersion: []string{"1.9"},
				},
			},
		},
		ratelimit: v1alpha1.LocalRateLimit{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "hello-zufardhiyaulhaq-dev",
				Namespace: "istio-system",
			},
			Spec: v1alpha1.LocalRateLimitSpec{
				Config: "public-gateway-config",
				Selector: v1alpha1.LocalRateLimitSelector{
					VHost: "hello.zufardhiyaulhaq.dev:443",
				},
				Limit: &v1alpha1.LocalRateLimit_Limit{
					Unit:            "hour",
					RequestsPerUnit: 1,
				},
			},
		},
		version:       "1.9",
		expectedError: false,
	},
	{
		name: "given ratelimit with second unit",
		config: v1alpha1.LocalRateLimitConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "sidecar-config",
				Namespace: "default",
			},
			Spec: v1alpha1.LocalRateLimitConfigSpec{
				Type: "sidecar",
				Selector: v1alpha1.LocalRateLimitConfigSelector{
					Labels:       map[string]string{"app": "my-service"},
					IstioVersion: []string{"1.10"},
				},
			},
		},
		ratelimit: v1alpha1.LocalRateLimit{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "fast-ratelimit",
				Namespace: "default",
			},
			Spec: v1alpha1.LocalRateLimitSpec{
				Config: "sidecar-config",
				Selector: v1alpha1.LocalRateLimitSelector{
					VHost: "inbound|8080||",
				},
				Limit: &v1alpha1.LocalRateLimit_Limit{
					Unit:            "second",
					RequestsPerUnit: 10,
				},
			},
		},
		version:       "1.10",
		expectedError: false,
	},
	{
		name: "given ratelimit with minute unit",
		config: v1alpha1.LocalRateLimitConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "sidecar-config",
				Namespace: "production",
			},
			Spec: v1alpha1.LocalRateLimitConfigSpec{
				Type: "sidecar",
				Selector: v1alpha1.LocalRateLimitConfigSelector{
					Labels:       map[string]string{"app": "api-service"},
					IstioVersion: []string{"1.11"},
				},
			},
		},
		ratelimit: v1alpha1.LocalRateLimit{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "api-ratelimit",
				Namespace: "production",
			},
			Spec: v1alpha1.LocalRateLimitSpec{
				Config: "sidecar-config",
				Selector: v1alpha1.LocalRateLimitSelector{
					VHost: "inbound|9090||",
				},
				Limit: &v1alpha1.LocalRateLimit_Limit{
					Unit:            "minute",
					RequestsPerUnit: 100,
				},
			},
		},
		version:       "1.11",
		expectedError: false,
	},
	{
		name: "given ratelimit with day unit",
		config: v1alpha1.LocalRateLimitConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "sidecar-config",
				Namespace: "staging",
			},
			Spec: v1alpha1.LocalRateLimitConfigSpec{
				Type: "sidecar",
				Selector: v1alpha1.LocalRateLimitConfigSelector{
					Labels:       map[string]string{"app": "worker"},
					IstioVersion: []string{"1.12"},
				},
			},
		},
		ratelimit: v1alpha1.LocalRateLimit{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "daily-ratelimit",
				Namespace: "staging",
			},
			Spec: v1alpha1.LocalRateLimitSpec{
				Config: "sidecar-config",
				Selector: v1alpha1.LocalRateLimitSelector{
					VHost: "inbound|3000||",
				},
				Limit: &v1alpha1.LocalRateLimit_Limit{
					Unit:            "day",
					RequestsPerUnit: 10000,
				},
			},
		},
		version:       "1.12",
		expectedError: false,
	},
}

func TestNewV3SidecarBuilder(t *testing.T) {
	for _, test := range V3SidecarBuilderTestGrid {
		t.Run(test.name, func(t *testing.T) {
			envoyfilter, err := ratelimit.NewV3SidecarBuilder(test.config, test.ratelimit, test.version).
				Build()

			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.ratelimit.Name+"-"+test.version, envoyfilter.Name)
				assert.Equal(t, test.ratelimit.Namespace, envoyfilter.Namespace)
			}
		})
	}
}

func TestNewV3SidecarBuilder_EnvoyFilterStructure(t *testing.T) {
	config := v1alpha1.LocalRateLimitConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sidecar-config",
			Namespace: "default",
		},
		Spec: v1alpha1.LocalRateLimitConfigSpec{
			Type: "sidecar",
			Selector: v1alpha1.LocalRateLimitConfigSelector{
				Labels:       map[string]string{"app": "my-service"},
				IstioVersion: []string{"1.9"},
			},
		},
	}

	rl := v1alpha1.LocalRateLimit{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-ratelimit",
			Namespace: "default",
		},
		Spec: v1alpha1.LocalRateLimitSpec{
			Config: "sidecar-config",
			Selector: v1alpha1.LocalRateLimitSelector{
				VHost: "inbound|8080||",
			},
			Limit: &v1alpha1.LocalRateLimit_Limit{
				Unit:            "second",
				RequestsPerUnit: 10,
			},
		},
	}

	envoyfilter, err := ratelimit.NewV3SidecarBuilder(config, rl, "1.9").Build()
	assert.NoError(t, err)

	// Verify TypeMeta
	assert.Equal(t, "EnvoyFilter", envoyfilter.Kind)
	assert.Equal(t, "networking.istio.io/v1alpha3", envoyfilter.APIVersion)

	// Verify ObjectMeta
	assert.Equal(t, "test-ratelimit-1.9", envoyfilter.Name)
	assert.Equal(t, "default", envoyfilter.Namespace)
	assert.Equal(t, "1.9", envoyfilter.Labels["istio/version"])

	// Verify WorkloadSelector
	assert.NotNil(t, envoyfilter.Spec.WorkloadSelector)
	assert.Equal(t, "my-service", envoyfilter.Spec.WorkloadSelector.Labels["app"])

	// Verify ConfigPatches
	assert.Len(t, envoyfilter.Spec.ConfigPatches, 1)
	configPatch := envoyfilter.Spec.ConfigPatches[0]
	assert.Equal(t, networking.EnvoyFilter_HTTP_ROUTE, configPatch.ApplyTo)

	// Verify Match - Sidecar should use SIDECAR_INBOUND context
	assert.NotNil(t, configPatch.Match)
	assert.Equal(t, networking.EnvoyFilter_SIDECAR_INBOUND, configPatch.Match.Context)

	// Verify Patch
	assert.NotNil(t, configPatch.Patch)
	assert.Equal(t, networking.EnvoyFilter_Patch_MERGE, configPatch.Patch.Operation)
}

func TestNewV3SidecarBuilder_DifferentVersions(t *testing.T) {
	versions := []string{"1.9", "1.10", "1.11", "1.12", "1.20", "1.25", "1.26"}

	for _, version := range versions {
		t.Run("version_"+version, func(t *testing.T) {
			config := v1alpha1.LocalRateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sidecar-config",
					Namespace: "default",
				},
				Spec: v1alpha1.LocalRateLimitConfigSpec{
					Type: "sidecar",
					Selector: v1alpha1.LocalRateLimitConfigSelector{
						Labels:       map[string]string{"app": "my-service"},
						IstioVersion: []string{version},
					},
				},
			}

			rl := v1alpha1.LocalRateLimit{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ratelimit",
					Namespace: "default",
				},
				Spec: v1alpha1.LocalRateLimitSpec{
					Config: "sidecar-config",
					Selector: v1alpha1.LocalRateLimitSelector{
						VHost: "inbound|8080||",
					},
					Limit: &v1alpha1.LocalRateLimit_Limit{
						Unit:            "second",
						RequestsPerUnit: 10,
					},
				},
			}

			envoyfilter, err := ratelimit.NewV3SidecarBuilder(config, rl, version).Build()
			assert.NoError(t, err)
			assert.Equal(t, "test-ratelimit-"+version, envoyfilter.Name)
			assert.Equal(t, version, envoyfilter.Labels["istio/version"])
		})
	}
}

func TestNewV3SidecarBuilder_VHostFormats(t *testing.T) {
	vhosts := []string{
		"inbound|8080||",
		"inbound|9090||",
		"inbound|3000||",
		"inbound|80||",
		"inbound|443||",
	}

	for _, vhost := range vhosts {
		t.Run("vhost_"+vhost, func(t *testing.T) {
			config := v1alpha1.LocalRateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sidecar-config",
					Namespace: "default",
				},
				Spec: v1alpha1.LocalRateLimitConfigSpec{
					Type: "sidecar",
					Selector: v1alpha1.LocalRateLimitConfigSelector{
						Labels:       map[string]string{"app": "my-service"},
						IstioVersion: []string{"1.9"},
					},
				},
			}

			rl := v1alpha1.LocalRateLimit{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ratelimit",
					Namespace: "default",
				},
				Spec: v1alpha1.LocalRateLimitSpec{
					Config: "sidecar-config",
					Selector: v1alpha1.LocalRateLimitSelector{
						VHost: vhost,
					},
					Limit: &v1alpha1.LocalRateLimit_Limit{
						Unit:            "second",
						RequestsPerUnit: 10,
					},
				},
			}

			envoyfilter, err := ratelimit.NewV3SidecarBuilder(config, rl, "1.9").Build()
			assert.NoError(t, err)

			configPatch := envoyfilter.Spec.ConfigPatches[0]
			routeConfig := configPatch.Match.GetRouteConfiguration()
			assert.NotNil(t, routeConfig)
			assert.NotNil(t, routeConfig.Vhost)
			assert.Equal(t, vhost, routeConfig.Vhost.Name)
		})
	}
}

func TestNewV3SidecarBuilder_ProxyMatch(t *testing.T) {
	config := v1alpha1.LocalRateLimitConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sidecar-config",
			Namespace: "default",
		},
		Spec: v1alpha1.LocalRateLimitConfigSpec{
			Type: "sidecar",
			Selector: v1alpha1.LocalRateLimitConfigSelector{
				Labels:       map[string]string{"app": "my-service"},
				IstioVersion: []string{"1.9"},
			},
		},
	}

	rl := v1alpha1.LocalRateLimit{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-ratelimit",
			Namespace: "default",
		},
		Spec: v1alpha1.LocalRateLimitSpec{
			Config: "sidecar-config",
			Selector: v1alpha1.LocalRateLimitSelector{
				VHost: "inbound|8080||",
			},
			Limit: &v1alpha1.LocalRateLimit_Limit{
				Unit:            "second",
				RequestsPerUnit: 10,
			},
		},
	}

	envoyfilter, err := ratelimit.NewV3SidecarBuilder(config, rl, "1.9").Build()
	assert.NoError(t, err)

	configPatch := envoyfilter.Spec.ConfigPatches[0]
	assert.NotNil(t, configPatch.Match.Proxy)
	assert.NotEmpty(t, configPatch.Match.Proxy.ProxyVersion)
}

func TestNewV3SidecarBuilder_DifferentNamespaces(t *testing.T) {
	namespaces := []string{"default", "production", "staging", "istio-system", "custom-namespace"}

	for _, ns := range namespaces {
		t.Run("namespace_"+ns, func(t *testing.T) {
			config := v1alpha1.LocalRateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sidecar-config",
					Namespace: ns,
				},
				Spec: v1alpha1.LocalRateLimitConfigSpec{
					Type: "sidecar",
					Selector: v1alpha1.LocalRateLimitConfigSelector{
						Labels:       map[string]string{"app": "my-service"},
						IstioVersion: []string{"1.9"},
					},
				},
			}

			rl := v1alpha1.LocalRateLimit{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ratelimit",
					Namespace: ns,
				},
				Spec: v1alpha1.LocalRateLimitSpec{
					Config: "sidecar-config",
					Selector: v1alpha1.LocalRateLimitSelector{
						VHost: "inbound|8080||",
					},
					Limit: &v1alpha1.LocalRateLimit_Limit{
						Unit:            "second",
						RequestsPerUnit: 10,
					},
				},
			}

			envoyfilter, err := ratelimit.NewV3SidecarBuilder(config, rl, "1.9").Build()
			assert.NoError(t, err)
			assert.Equal(t, ns, envoyfilter.Namespace)
		})
	}
}

func TestNewV3SidecarBuilder_MultipleLabels(t *testing.T) {
	config := v1alpha1.LocalRateLimitConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sidecar-config",
			Namespace: "default",
		},
		Spec: v1alpha1.LocalRateLimitConfigSpec{
			Type: "sidecar",
			Selector: v1alpha1.LocalRateLimitConfigSelector{
				Labels: map[string]string{
					"app":     "my-service",
					"version": "v1",
					"env":     "production",
				},
				IstioVersion: []string{"1.9"},
			},
		},
	}

	rl := v1alpha1.LocalRateLimit{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-ratelimit",
			Namespace: "default",
		},
		Spec: v1alpha1.LocalRateLimitSpec{
			Config: "sidecar-config",
			Selector: v1alpha1.LocalRateLimitSelector{
				VHost: "inbound|8080||",
			},
			Limit: &v1alpha1.LocalRateLimit_Limit{
				Unit:            "second",
				RequestsPerUnit: 10,
			},
		},
	}

	envoyfilter, err := ratelimit.NewV3SidecarBuilder(config, rl, "1.9").Build()
	assert.NoError(t, err)

	assert.Equal(t, "my-service", envoyfilter.Spec.WorkloadSelector.Labels["app"])
	assert.Equal(t, "v1", envoyfilter.Spec.WorkloadSelector.Labels["version"])
	assert.Equal(t, "production", envoyfilter.Spec.WorkloadSelector.Labels["env"])
}
