package config

import (
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/global/types"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/utils"
	"gopkg.in/yaml.v2"

	networking "istio.io/api/networking/v1alpha3"
	clientnetworking "istio.io/client-go/pkg/apis/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type V3GatewayBuilder struct {
	Config  v1alpha1.GlobalRateLimitConfig
	Version string
}

func NewV3GatewayBuilder(config v1alpha1.GlobalRateLimitConfig, version string) *V3GatewayBuilder {
	return &V3GatewayBuilder{
		Config:  config,
		Version: version,
	}
}

func (g *V3GatewayBuilder) Build() (*clientnetworking.EnvoyFilter, error) {
	httpFilter, err := g.buildHttpFilterPatch()
	if err != nil {
		return nil, err
	}

	cluster, err := g.buildClusterPatch()
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
				httpFilter,
				cluster,
			},
		},
	}

	return envoyfilter, nil
}

func (g *V3GatewayBuilder) buildHttpFilterPatch() (*networking.EnvoyFilter_EnvoyConfigObjectPatch, error) {
	value, err := g.buildHttpFilterPatchValue()
	if err != nil {
		return nil, err
	}

	listener, err := g.buildHttpFilterListener()
	if err != nil {
		return nil, err
	}

	patches := &networking.EnvoyFilter_EnvoyConfigObjectPatch{
		ApplyTo: networking.EnvoyFilter_HTTP_FILTER,
		Match: &networking.EnvoyFilter_EnvoyConfigObjectMatch{
			Context: networking.EnvoyFilter_GATEWAY,
			ObjectTypes: &networking.EnvoyFilter_EnvoyConfigObjectMatch_Listener{
				Listener: listener,
			},
			Proxy: g.buildProxyMatch(),
		},
		Patch: &networking.EnvoyFilter_Patch{
			Operation: networking.EnvoyFilter_Patch_INSERT_BEFORE,
			Value:     utils.ConvertYaml2Struct(value),
		},
	}

	return patches, nil
}

func (g *V3GatewayBuilder) buildHttpFilterListener() (*networking.EnvoyFilter_ListenerMatch, error) {
	listener := &networking.EnvoyFilter_ListenerMatch{
		FilterChain: &networking.EnvoyFilter_ListenerMatch_FilterChainMatch{
			Filter: &networking.EnvoyFilter_ListenerMatch_FilterMatch{
				Name: "envoy.filters.network.http_connection_manager",
				SubFilter: &networking.EnvoyFilter_ListenerMatch_SubFilterMatch{
					Name: "envoy.filters.http.router",
				},
			},
		},
	}

	if g.Config.Spec.Selector.SNI != nil {
		listener.FilterChain.Sni = *g.Config.Spec.Selector.SNI
	}

	return listener, nil
}

func (g *V3GatewayBuilder) buildHttpFilterPatchValue() (string, error) {
	values := types.HttpFilterPatchValues{
		Name: "envoy.filters.http.ratelimit",
		TypedConfig: types.TypedConfig{
			Type:            "type.googleapis.com/envoy.extensions.filters.http.ratelimit.v3.RateLimit",
			Domain:          g.Config.Spec.Ratelimit.Spec.Domain,
			FailureModeDeny: &g.Config.Spec.Ratelimit.Spec.FailureModeDeny,
			Timeout:         g.Config.Spec.Ratelimit.Spec.Timeout,
			RateLimitService: types.RateLimitService{
				TransportAPIVersion: "V3",
				GRPCService: types.GRPCService{
					Timeout: g.Config.Spec.Ratelimit.Spec.Timeout,
					EnvoyGRPC: types.EnvoyGRPC{
						ClusterName: g.Config.Name,
					},
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

func (g *V3GatewayBuilder) buildClusterPatch() (*networking.EnvoyFilter_EnvoyConfigObjectPatch, error) {
	value, err := g.buildClusterPatchValue()
	if err != nil {
		return nil, err
	}

	patches := &networking.EnvoyFilter_EnvoyConfigObjectPatch{
		ApplyTo: networking.EnvoyFilter_CLUSTER,
		Match: &networking.EnvoyFilter_EnvoyConfigObjectMatch{
			Context: networking.EnvoyFilter_GATEWAY,
			ObjectTypes: &networking.EnvoyFilter_EnvoyConfigObjectMatch_Cluster{
				Cluster: &networking.EnvoyFilter_ClusterMatch{
					Service: g.Config.Spec.Ratelimit.Spec.Service.Address,
				},
			},
			Proxy: g.buildProxyMatch(),
		},
		Patch: &networking.EnvoyFilter_Patch{
			Operation: networking.EnvoyFilter_Patch_ADD,
			Value:     utils.ConvertYaml2Struct(value),
		},
	}

	return patches, nil
}

func (g *V3GatewayBuilder) buildClusterPatchValue() (string, error) {
	values := types.ClusterPatchValues{
		Name:                 g.Config.Name,
		Type:                 "STRICT_DNS",
		ConnectTimeout:       g.Config.Spec.Ratelimit.Spec.Timeout,
		HTTP2ProtocolOptions: types.HTTP2ProtocolOptions{},
		LbPolicy:             "ROUND_ROBIN",
		LoadAssignment: types.LoadAssignment{
			ClusterName: g.Config.Name,
			Endpoints: []types.Endpoints{
				{
					LbEndpoints: []types.LbEndpoints{
						{
							Endpoint: types.Endpoint{
								Address: types.Address{
									SocketAddress: types.SocketAddress{
										Address:   g.Config.Spec.Ratelimit.Spec.Service.Address,
										PortValue: g.Config.Spec.Ratelimit.Spec.Service.Port,
									},
								},
							},
						},
					},
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

func (g *V3GatewayBuilder) buildProxyMatch() *networking.EnvoyFilter_ProxyMatch {
	return &networking.EnvoyFilter_ProxyMatch{
		ProxyVersion: utils.WellKnownVersions[g.Version],
	}
}
