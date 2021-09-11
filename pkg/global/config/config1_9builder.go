package config

import (
	"context"
	"fmt"
	"os"

	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/utils"
	"gopkg.in/yaml.v2"

	types "github.com/gogo/protobuf/types"
	istioAPINetworking "istio.io/api/networking/v1alpha3"
	istioClientNetworking "istio.io/client-go/pkg/apis/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"

	versionedclient "istio.io/client-go/pkg/clientset/versioned"
)

const http2_protocol_options string = "{}"

type GlobalRateLimitConfig1_9EnvoyFilterBuilder struct {
	Config v1alpha1.GlobalRateLimitConfig
}

func NewGlobalRateLimitConfig1_9EnvoyFilterBuilder(config v1alpha1.GlobalRateLimitConfig) *GlobalRateLimitConfig1_9EnvoyFilterBuilder {
	return &GlobalRateLimitConfig1_9EnvoyFilterBuilder{
		Config: config,
	}
}

func (g *GlobalRateLimitConfig1_9EnvoyFilterBuilder) Build() (*istioClientNetworking.EnvoyFilter, error) {
	var patchContext istioAPINetworking.EnvoyFilter_PatchContext
	if g.Config.Spec.Type == "gateway" {
		patchContext = istioAPINetworking.EnvoyFilter_GATEWAY
	}

	envoyfilter := &istioClientNetworking.EnvoyFilter{
		ObjectMeta: metav1.ObjectMeta{
			Name:      g.Config.Name,
			Namespace: g.Config.Namespace,
			Labels: map[string]string{
				"generator": "istio-rateltimit-operator",
			},
		},
		Spec: istioAPINetworking.EnvoyFilter{
			WorkloadSelector: &istioAPINetworking.WorkloadSelector{
				Labels: g.Config.Spec.Selector.Labels,
			},
			ConfigPatches: []*istioAPINetworking.EnvoyFilter_EnvoyConfigObjectPatch{
				{
					ApplyTo: istioAPINetworking.EnvoyFilter_HTTP_FILTER,
					Match: &istioAPINetworking.EnvoyFilter_EnvoyConfigObjectMatch{
						Context: patchContext,
						ObjectTypes: &istioAPINetworking.EnvoyFilter_EnvoyConfigObjectMatch_Listener{
							Listener: &istioAPINetworking.EnvoyFilter_ListenerMatch{
								FilterChain: &istioAPINetworking.EnvoyFilter_ListenerMatch_FilterChainMatch{
									Filter: &istioAPINetworking.EnvoyFilter_ListenerMatch_FilterMatch{
										Name: "envoy.filters.network.http_connection_manager",
										SubFilter: &istioAPINetworking.EnvoyFilter_ListenerMatch_SubFilterMatch{
											Name: "envoy.filters.http.router",
										},
									},
								},
							},
						},
						Proxy: &istioAPINetworking.EnvoyFilter_ProxyMatch{
							ProxyVersion: utils.WellKnownVersions["1.9"],
						},
					},
					Patch: &istioAPINetworking.EnvoyFilter_Patch{
						Operation: istioAPINetworking.EnvoyFilter_Patch_INSERT_BEFORE,
						Value: &types.Struct{
							Fields: map[string]*types.Value{
								"name": {
									Kind: &types.Value_StringValue{
										StringValue: "envoy.filters.http.ratelimit",
									},
								},
								"typed_config": {
									Kind: &types.Value_StructValue{
										StructValue: &types.Struct{
											Fields: map[string]*types.Value{
												"@type": {
													Kind: &types.Value_StringValue{
														StringValue: "type.googleapis.com/envoy.extensions.filters.http.ratelimit.v3.RateLimit",
													},
												},
												"domain": {
													Kind: &types.Value_StringValue{
														StringValue: g.Config.Spec.Ratelimit.Spec.Domain,
													},
												},
												"failure_mode_deny": {
													Kind: &types.Value_BoolValue{
														BoolValue: g.Config.Spec.Ratelimit.Spec.FailureModeDeny,
													},
												},
												"rate_limit_service": {
													Kind: &types.Value_StructValue{
														StructValue: &types.Struct{
															Fields: map[string]*types.Value{
																"transport_api_version": {
																	Kind: &types.Value_StringValue{
																		StringValue: "V3",
																	},
																},
																"grpc_service": {
																	Kind: &types.Value_StructValue{
																		StructValue: &types.Struct{
																			Fields: map[string]*types.Value{
																				"timeout": {
																					Kind: &types.Value_StringValue{
																						StringValue: g.Config.Spec.Ratelimit.Spec.Timeout,
																					},
																				},
																				"envoy_grpc": {
																					Kind: &types.Value_StructValue{
																						StructValue: &types.Struct{
																							Fields: map[string]*types.Value{
																								"cluster_name": {
																									Kind: &types.Value_StringValue{
																										StringValue: g.Config.Name,
																									},
																								},
																							},
																						},
																					},
																				},
																			},
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
				{
					ApplyTo: istioAPINetworking.EnvoyFilter_CLUSTER,
					Match: &istioAPINetworking.EnvoyFilter_EnvoyConfigObjectMatch{
						Context: patchContext,
						ObjectTypes: &istioAPINetworking.EnvoyFilter_EnvoyConfigObjectMatch_Cluster{
							Cluster: &istioAPINetworking.EnvoyFilter_ClusterMatch{
								Service: g.Config.Spec.Ratelimit.Spec.Service.Address,
							},
						},
						Proxy: &istioAPINetworking.EnvoyFilter_ProxyMatch{
							ProxyVersion: utils.WellKnownVersions["1.9"],
						},
					},
					Patch: &istioAPINetworking.EnvoyFilter_Patch{
						Operation: istioAPINetworking.EnvoyFilter_Patch_ADD,
						Value: &types.Struct{
							Fields: map[string]*types.Value{
								"name": {
									Kind: &types.Value_StringValue{
										StringValue: g.Config.Name,
									},
								},
								"type": {
									Kind: &types.Value_StringValue{
										StringValue: "STRICT_DNS",
									},
								},
								"connect_timeout": {
									Kind: &types.Value_StringValue{
										StringValue: g.Config.Spec.Ratelimit.Spec.Timeout,
									},
								},
								"http2_protocol_options": {
									Kind: &types.Value_StructValue{
										StructValue: &types.Struct{
											Fields: map[string]*types.Value{},
										},
									},
								},
								"lb_policy": {
									Kind: &types.Value_StringValue{
										StringValue: "ROUND_ROBIN",
									},
								},
								"load_assignment": {
									Kind: &types.Value_StructValue{
										StructValue: &types.Struct{
											Fields: map[string]*types.Value{
												"cluster_name": {
													Kind: &types.Value_StringValue{
														StringValue: g.Config.Name,
													},
												},
												"endpoints": {
													Kind: &types.Value_ListValue{
														ListValue: &types.ListValue{
															Values: []*types.Value{
																{
																	Kind: &types.Value_StructValue{
																		StructValue: &types.Struct{
																			Fields: map[string]*types.Value{
																				"lb_endpoints": {
																					Kind: &types.Value_ListValue{
																						ListValue: &types.ListValue{
																							Values: []*types.Value{
																								{
																									Kind: &types.Value_StructValue{
																										StructValue: &types.Struct{
																											Fields: map[string]*types.Value{
																												"endpoint": {
																													Kind: &types.Value_StructValue{
																														StructValue: &types.Struct{
																															Fields: map[string]*types.Value{
																																"address": {
																																	Kind: &types.Value_StructValue{
																																		StructValue: &types.Struct{
																																			Fields: map[string]*types.Value{
																																				"socket_address": {
																																					Kind: &types.Value_StructValue{
																																						StructValue: &types.Struct{
																																							Fields: map[string]*types.Value{
																																								"address": {
																																									Kind: &types.Value_StringValue{
																																										StringValue: g.Config.Spec.Ratelimit.Spec.Service.Address,
																																									},
																																								},
																																								"port_value": {
																																									Kind: &types.Value_NumberValue{
																																										NumberValue: float64(g.Config.Spec.Ratelimit.Spec.Service.Port),
																																									},
																																								},
																																							},
																																						},
																																					},
																																				},
																																			},
																																		},
																																	},
																																},
																															},
																														},
																													},
																												},
																											},
																										},
																									},
																								},
																							},
																						},
																					},
																				},
																			},
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	kubeconfig := os.Getenv("KUBECONFIG")
	restConfig, _ := clientcmd.BuildConfigFromFlags("", kubeconfig)
	ic, _ := versionedclient.NewForConfig(restConfig)
	ctx := context.Background()
	_, err := ic.NetworkingV1alpha3().EnvoyFilters(g.Config.Namespace).Create(ctx, envoyfilter, v1.CreateOptions{})
	fmt.Println(err)

	bytes, _ := yaml.Marshal(envoyfilter)
	fmt.Println(string(bytes))

	return envoyfilter, nil
}

var RateLimit1_9Patch = `{"name": "envoy.filters.http.ratelimit","typed_config": {"@type": "type.googleapis.com/envoy.extensions.filters.http.ratelimit.v3.RateLimit","domain": "%s","failure_mode_deny": %t,"rate_limit_service": {"grpc_service": {"envoy_grpc": {"cluster_name": "%s"},"timeout": "%s"},"transport_api_version": "V3"}}}`
