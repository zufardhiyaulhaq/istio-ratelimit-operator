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

type Config1_9Builder struct {
	Config v1alpha1.GlobalRateLimitConfig
}

func NewConfig1_9Builder(config v1alpha1.GlobalRateLimitConfig) *Config1_9Builder {
	return &Config1_9Builder{
		Config: config,
	}
}

func (g *Config1_9Builder) Build() (*clientnetworking.EnvoyFilter, error) {
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
			Name:      g.buildName(),
			Namespace: g.Config.Namespace,
			Labels: map[string]string{
				"app.kubernetes.io/created-by": "istio-rateltimit-operator",
				"app.kubernetes.io/managed-by": "istio-rateltimit-operator",
				"istio/version":                "1.9",
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

func (g *Config1_9Builder) buildName() string {
	return g.Config.Name + "-1.9"
}

func (g *Config1_9Builder) buildHttpFilterPatch() (*networking.EnvoyFilter_EnvoyConfigObjectPatch, error) {
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
			Context: g.buildContext(),
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

func (g *Config1_9Builder) buildHttpFilterListener() (*networking.EnvoyFilter_ListenerMatch, error) {
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

func (g *Config1_9Builder) buildHttpFilterPatchValue() (string, error) {
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

func (g *Config1_9Builder) buildClusterPatch() (*networking.EnvoyFilter_EnvoyConfigObjectPatch, error) {
	value, err := g.buildClusterPatchValue()
	if err != nil {
		return nil, err
	}

	patches := &networking.EnvoyFilter_EnvoyConfigObjectPatch{
		ApplyTo: networking.EnvoyFilter_CLUSTER,
		Match: &networking.EnvoyFilter_EnvoyConfigObjectMatch{
			Context: g.buildContext(),
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

func (g *Config1_9Builder) buildClusterPatchValue() (string, error) {
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

func (g *Config1_9Builder) buildContext() networking.EnvoyFilter_PatchContext {
	if g.Config.Spec.Type == "gateway" {
		return networking.EnvoyFilter_GATEWAY
	}

	return networking.EnvoyFilter_GATEWAY
}

func (g *Config1_9Builder) buildProxyMatch() *networking.EnvoyFilter_ProxyMatch {
	return &networking.EnvoyFilter_ProxyMatch{
		ProxyVersion: utils.WellKnownVersions["1.9"],
	}
}
