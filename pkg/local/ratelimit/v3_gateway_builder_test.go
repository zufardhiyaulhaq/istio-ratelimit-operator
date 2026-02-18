package ratelimit_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/local/ratelimit"

	networking "istio.io/api/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type V3GatewayBuilderTestCase struct {
	name          string
	config        v1alpha1.LocalRateLimitConfig
	ratelimit     v1alpha1.LocalRateLimit
	version       string
	expectedError bool
}

var V3GatewayBuilderTestGrid = []V3GatewayBuilderTestCase{
	{
		name: "given correct ratelimit",
		config: v1alpha1.LocalRateLimitConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "gateway-config",
				Namespace: "istio-system",
			},
			Spec: v1alpha1.LocalRateLimitConfigSpec{
				Type: "gateway",
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
		name: "given ratelimit with route selector",
		config: v1alpha1.LocalRateLimitConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "gateway-config",
				Namespace: "istio-system",
			},
			Spec: v1alpha1.LocalRateLimitConfigSpec{
				Type: "gateway",
				Selector: v1alpha1.LocalRateLimitConfigSelector{
					Labels:       map[string]string{"app": "ingressgateway"},
					IstioVersion: []string{"1.10"},
				},
			},
		},
		ratelimit: v1alpha1.LocalRateLimit{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "api-ratelimit",
				Namespace: "istio-system",
			},
			Spec: v1alpha1.LocalRateLimitSpec{
				Config: "gateway-config",
				Selector: v1alpha1.LocalRateLimitSelector{
					VHost: "api.example.com:443",
					Route: stringPtr("api-route"),
				},
				Limit: &v1alpha1.LocalRateLimit_Limit{
					Unit:            "minute",
					RequestsPerUnit: 100,
				},
			},
		},
		version:       "1.10",
		expectedError: false,
	},
	{
		name: "given ratelimit with second unit",
		config: v1alpha1.LocalRateLimitConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "gateway-config",
				Namespace: "istio-system",
			},
			Spec: v1alpha1.LocalRateLimitConfigSpec{
				Type: "gateway",
				Selector: v1alpha1.LocalRateLimitConfigSelector{
					Labels:       map[string]string{"app": "gateway"},
					IstioVersion: []string{"1.11"},
				},
			},
		},
		ratelimit: v1alpha1.LocalRateLimit{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "fast-ratelimit",
				Namespace: "default",
			},
			Spec: v1alpha1.LocalRateLimitSpec{
				Config: "gateway-config",
				Selector: v1alpha1.LocalRateLimitSelector{
					VHost: "fast.example.com:443",
				},
				Limit: &v1alpha1.LocalRateLimit_Limit{
					Unit:            "second",
					RequestsPerUnit: 10,
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
				Name:      "gateway-config",
				Namespace: "istio-system",
			},
			Spec: v1alpha1.LocalRateLimitConfigSpec{
				Type: "gateway",
				Selector: v1alpha1.LocalRateLimitConfigSelector{
					Labels:       map[string]string{"app": "gateway"},
					IstioVersion: []string{"1.12"},
				},
			},
		},
		ratelimit: v1alpha1.LocalRateLimit{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "daily-ratelimit",
				Namespace: "production",
			},
			Spec: v1alpha1.LocalRateLimitSpec{
				Config: "gateway-config",
				Selector: v1alpha1.LocalRateLimitSelector{
					VHost: "daily.example.com:443",
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

func TestNewV3GatewayBuilder(t *testing.T) {
	for _, test := range V3GatewayBuilderTestGrid {
		t.Run(test.name, func(t *testing.T) {
			envoyfilter, err := ratelimit.NewV3GatewayBuilder(test.config, test.ratelimit, test.version).
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

func TestNewV3GatewayBuilder_EnvoyFilterStructure(t *testing.T) {
	config := v1alpha1.LocalRateLimitConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "gateway-config",
			Namespace: "istio-system",
		},
		Spec: v1alpha1.LocalRateLimitConfigSpec{
			Type: "gateway",
			Selector: v1alpha1.LocalRateLimitConfigSelector{
				Labels:       map[string]string{"app": "ingressgateway"},
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

	envoyfilter, err := ratelimit.NewV3GatewayBuilder(config, rl, "1.9").Build()
	assert.NoError(t, err)

	// Verify TypeMeta
	assert.Equal(t, "EnvoyFilter", envoyfilter.Kind)
	assert.Equal(t, "networking.istio.io/v1alpha3", envoyfilter.APIVersion)

	// Verify ObjectMeta
	assert.Equal(t, "test-ratelimit-1.9", envoyfilter.Name)
	assert.Equal(t, "istio-system", envoyfilter.Namespace)
	assert.Equal(t, "1.9", envoyfilter.Labels["istio/version"])

	// Verify WorkloadSelector
	assert.NotNil(t, envoyfilter.Spec.WorkloadSelector)
	assert.Equal(t, "ingressgateway", envoyfilter.Spec.WorkloadSelector.Labels["app"])

	// Verify ConfigPatches
	assert.Len(t, envoyfilter.Spec.ConfigPatches, 1)
	configPatch := envoyfilter.Spec.ConfigPatches[0]
	assert.Equal(t, networking.EnvoyFilter_HTTP_ROUTE, configPatch.ApplyTo)

	// Verify Match
	assert.NotNil(t, configPatch.Match)
	assert.Equal(t, networking.EnvoyFilter_GATEWAY, configPatch.Match.Context)

	// Verify Patch
	assert.NotNil(t, configPatch.Patch)
	assert.Equal(t, networking.EnvoyFilter_Patch_MERGE, configPatch.Patch.Operation)
}

func TestNewV3GatewayBuilder_WithRouteSelector(t *testing.T) {
	routeName := "my-route"
	config := v1alpha1.LocalRateLimitConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "gateway-config",
			Namespace: "istio-system",
		},
		Spec: v1alpha1.LocalRateLimitConfigSpec{
			Type: "gateway",
			Selector: v1alpha1.LocalRateLimitConfigSelector{
				Labels:       map[string]string{"app": "ingressgateway"},
				IstioVersion: []string{"1.9"},
			},
		},
	}

	rl := v1alpha1.LocalRateLimit{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "route-ratelimit",
			Namespace: "istio-system",
		},
		Spec: v1alpha1.LocalRateLimitSpec{
			Config: "gateway-config",
			Selector: v1alpha1.LocalRateLimitSelector{
				VHost: "test.example.com:443",
				Route: &routeName,
			},
			Limit: &v1alpha1.LocalRateLimit_Limit{
				Unit:            "second",
				RequestsPerUnit: 10,
			},
		},
	}

	envoyfilter, err := ratelimit.NewV3GatewayBuilder(config, rl, "1.9").Build()
	assert.NoError(t, err)

	// Verify the route configuration match includes route name
	configPatch := envoyfilter.Spec.ConfigPatches[0]
	assert.NotNil(t, configPatch.Match)
	routeConfig := configPatch.Match.GetRouteConfiguration()
	assert.NotNil(t, routeConfig)
	assert.NotNil(t, routeConfig.Vhost)
	assert.Equal(t, "test.example.com:443", routeConfig.Vhost.Name)
	assert.NotNil(t, routeConfig.Vhost.Route)
	assert.Equal(t, "my-route", routeConfig.Vhost.Route.Name)
}

func TestNewV3GatewayBuilder_DifferentVersions(t *testing.T) {
	versions := []string{"1.9", "1.10", "1.11", "1.12", "1.20", "1.25", "1.26"}

	for _, version := range versions {
		t.Run("version_"+version, func(t *testing.T) {
			config := v1alpha1.LocalRateLimitConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "gateway-config",
					Namespace: "istio-system",
				},
				Spec: v1alpha1.LocalRateLimitConfigSpec{
					Type: "gateway",
					Selector: v1alpha1.LocalRateLimitConfigSelector{
						Labels:       map[string]string{"app": "gateway"},
						IstioVersion: []string{version},
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

			envoyfilter, err := ratelimit.NewV3GatewayBuilder(config, rl, version).Build()
			assert.NoError(t, err)
			assert.Equal(t, "test-ratelimit-"+version, envoyfilter.Name)
			assert.Equal(t, version, envoyfilter.Labels["istio/version"])
		})
	}
}

func TestNewV3GatewayBuilder_ProxyMatch(t *testing.T) {
	config := v1alpha1.LocalRateLimitConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "gateway-config",
			Namespace: "istio-system",
		},
		Spec: v1alpha1.LocalRateLimitConfigSpec{
			Type: "gateway",
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

	envoyfilter, err := ratelimit.NewV3GatewayBuilder(config, rl, "1.9").Build()
	assert.NoError(t, err)

	configPatch := envoyfilter.Spec.ConfigPatches[0]
	assert.NotNil(t, configPatch.Match.Proxy)
	assert.NotEmpty(t, configPatch.Match.Proxy.ProxyVersion)
}
