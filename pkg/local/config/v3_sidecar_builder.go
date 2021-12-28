package config

import (
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/local/types"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/utils"
	"gopkg.in/yaml.v2"

	networking "istio.io/api/networking/v1alpha3"
	clientnetworking "istio.io/client-go/pkg/apis/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type V3SidecarBuilder struct {
	Config  v1alpha1.LocalRateLimitConfig
	Version string
}

func NewV3SidecarBuilder(config v1alpha1.LocalRateLimitConfig, version string) *V3SidecarBuilder {
	return &V3SidecarBuilder{
		Config:  config,
		Version: version,
	}
}

func (g *V3SidecarBuilder) Build() (*clientnetworking.EnvoyFilter, error) {
	configPatches, err := g.buildConfigPatches()
	if err != nil {
		return nil, err
	}

	envoyfilter := &clientnetworking.EnvoyFilter{
		TypeMeta: metav1.TypeMeta{
			Kind:       "EnvoyFilter",
			APIVersion: "networking.istio.io/v1alpha3",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      g.Config.Name + "-" + g.Version,
			Namespace: g.Config.Namespace,
			Labels: map[string]string{
				"istio/version": g.Version,
			},
		},
		Spec: networking.EnvoyFilter{
			WorkloadSelector: &networking.WorkloadSelector{
				Labels: g.Config.Spec.Selector.Labels,
			},
			ConfigPatches: []*networking.EnvoyFilter_EnvoyConfigObjectPatch{
				configPatches,
			},
		},
	}

	return envoyfilter, nil
}

func (g *V3SidecarBuilder) buildConfigPatches() (*networking.EnvoyFilter_EnvoyConfigObjectPatch, error) {
	match, err := g.buildMatch()
	if err != nil {
		return nil, err
	}

	patch, err := g.buildPatch()
	if err != nil {
		return nil, err
	}

	configPatches := &networking.EnvoyFilter_EnvoyConfigObjectPatch{
		ApplyTo: networking.EnvoyFilter_HTTP_FILTER,
		Match:   match,
		Patch:   patch,
	}

	return configPatches, nil
}

func (g *V3SidecarBuilder) buildMatch() (*networking.EnvoyFilter_EnvoyConfigObjectMatch, error) {
	match := &networking.EnvoyFilter_EnvoyConfigObjectMatch{
		Context: networking.EnvoyFilter_SIDECAR_INBOUND,
		ObjectTypes: &networking.EnvoyFilter_EnvoyConfigObjectMatch_Listener{
			Listener: &networking.EnvoyFilter_ListenerMatch{
				FilterChain: &networking.EnvoyFilter_ListenerMatch_FilterChainMatch{
					Filter: &networking.EnvoyFilter_ListenerMatch_FilterMatch{
						Name: "envoy.filters.network.http_connection_manager",
						SubFilter: &networking.EnvoyFilter_ListenerMatch_SubFilterMatch{
							Name: "envoy.filters.http.router",
						},
					},
				},
			},
		},
		Proxy: g.buildProxyMatch(),
	}

	return match, nil
}

func (g *V3SidecarBuilder) buildPatch() (*networking.EnvoyFilter_Patch, error) {
	value, err := g.buildPatchValue()
	if err != nil {
		return nil, err
	}

	patch := &networking.EnvoyFilter_Patch{
		Operation: networking.EnvoyFilter_Patch_INSERT_BEFORE,
		Value:     utils.ConvertYaml2Struct(value),
	}

	return patch, nil
}

func (g *V3SidecarBuilder) buildPatchValue() (string, error) {
	values := types.LocalRateLimitConfig_Value{
		Name: "envoy.filters.http.ratelimit",
		TypedConfig: types.LocalRateLimitConfig_TypedConfig{
			Type: "type.googleapis.com/envoy.extensions.filters.http.local_ratelimit.v3.LocalRateLimit",
		},
	}

	bytes, err := yaml.Marshal(&values)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (g *V3SidecarBuilder) buildProxyMatch() *networking.EnvoyFilter_ProxyMatch {
	return &networking.EnvoyFilter_ProxyMatch{
		ProxyVersion: utils.WellKnownVersions[g.Version],
	}
}
