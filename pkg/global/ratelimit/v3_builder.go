package ratelimit

import (
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/global/types"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/utils"
	"gopkg.in/yaml.v2"

	networking "istio.io/api/networking/v1alpha3"
	clientnetworking "istio.io/client-go/pkg/apis/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type V3Builder struct {
	Config    v1alpha1.GlobalRateLimitConfig
	RateLimit v1alpha1.GlobalRateLimit
	Version   string
}

func NewV3Builder(config v1alpha1.GlobalRateLimitConfig, ratelimit v1alpha1.GlobalRateLimit, version string) *V3Builder {
	return &V3Builder{
		Config:    config,
		RateLimit: ratelimit,
		Version:   version,
	}
}

func (g *V3Builder) Build() (*clientnetworking.EnvoyFilter, error) {
	httpRoute, err := g.buildHttpRoutePatch()
	if err != nil {
		return nil, err
	}

	envoyfilter := &clientnetworking.EnvoyFilter{
		TypeMeta: metav1.TypeMeta{
			Kind:       "EnvoyFilter",
			APIVersion: "networking.istio.io/v1alpha3",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      g.buildName(),
			Namespace: g.RateLimit.Namespace,
			Labels: map[string]string{
				"istio/version": "1.9",
			},
		},
		Spec: networking.EnvoyFilter{
			WorkloadSelector: &networking.WorkloadSelector{
				Labels: g.Config.Spec.Selector.Labels,
			},
			ConfigPatches: []*networking.EnvoyFilter_EnvoyConfigObjectPatch{
				httpRoute,
			},
		},
	}

	return envoyfilter, nil
}

func (g *V3Builder) buildName() string {
	return g.RateLimit.Name + "-1.9"
}

func (g *V3Builder) buildHttpRoutePatch() (*networking.EnvoyFilter_EnvoyConfigObjectPatch, error) {
	value, err := g.buildHttpRoutePatchValue()
	if err != nil {
		return nil, err
	}

	routeConfiguration, err := g.buildHttpRouteConfiguration()
	if err != nil {
		return nil, err
	}

	patches := &networking.EnvoyFilter_EnvoyConfigObjectPatch{
		ApplyTo: networking.EnvoyFilter_HTTP_ROUTE,
		Match: &networking.EnvoyFilter_EnvoyConfigObjectMatch{
			Context: g.buildContext(),
			ObjectTypes: &networking.EnvoyFilter_EnvoyConfigObjectMatch_RouteConfiguration{
				RouteConfiguration: routeConfiguration,
			},
			Proxy: g.buildProxyMatch(),
		},
		Patch: &networking.EnvoyFilter_Patch{
			Operation: networking.EnvoyFilter_Patch_MERGE,
			Value:     utils.ConvertYaml2Struct(value),
		},
	}

	return patches, nil
}

func (g *V3Builder) buildHttpRouteConfiguration() (*networking.EnvoyFilter_RouteConfigurationMatch, error) {
	routeConfiguration := &networking.EnvoyFilter_RouteConfigurationMatch{
		Vhost: &networking.EnvoyFilter_RouteConfigurationMatch_VirtualHostMatch{
			Name: g.RateLimit.Spec.Selector.VHost,
			Route: &networking.EnvoyFilter_RouteConfigurationMatch_RouteMatch{
				Action: networking.EnvoyFilter_RouteConfigurationMatch_RouteMatch_ANY,
			},
		},
	}

	if g.RateLimit.Spec.Selector.Route != nil {
		routeConfiguration.Vhost.Route.Name = *g.RateLimit.Spec.Selector.Route
	}

	return routeConfiguration, nil
}

func (g *V3Builder) buildHttpRoutePatchValue() (string, error) {
	values := types.RoutePatchValues{
		Route: types.Route{
			Ratelimits: []types.RateLimits{
				{
					Actions: g.RateLimit.Spec.Matcher,
				},
			},
		},
	}

	bytes, err := yaml.Marshal(&values)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (g *V3Builder) buildContext() networking.EnvoyFilter_PatchContext {
	if g.Config.Spec.Type == "gateway" {
		return networking.EnvoyFilter_GATEWAY
	}

	return networking.EnvoyFilter_GATEWAY
}

func (g *V3Builder) buildProxyMatch() *networking.EnvoyFilter_ProxyMatch {
	return &networking.EnvoyFilter_ProxyMatch{
		ProxyVersion: utils.WellKnownVersions["1.9"],
	}
}
