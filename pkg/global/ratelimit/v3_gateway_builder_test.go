package ratelimit_test

import (
	"testing"

	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/assert"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/global/ratelimit"

	proto_types "github.com/gogo/protobuf/types"
	networking "istio.io/api/networking/v1alpha3"
	clientnetworking "istio.io/client-go/pkg/apis/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type V3GatewayBuilderTestCase struct {
	name                string
	config              v1alpha1.GlobalRateLimitConfig
	ratelimit           v1alpha1.GlobalRateLimit
	expectedError       bool
	expectedEnvoyFilter clientnetworking.EnvoyFilter
}

var V3GatewayBuilderTestGrid = []V3GatewayBuilderTestCase{
	{
		name: "given correct ratelimit with request header",
		config: v1alpha1.GlobalRateLimitConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "public-gateway-config",
				Namespace: "istio-system",
			},
			Spec: v1alpha1.GlobalRateLimitConfigSpec{
				Type: "gateway",
				Selector: v1alpha1.GlobalRateLimitConfigSelector{
					IstioVersion: []string{"1.9"},
				},
				Ratelimit: v1alpha1.GlobalRateLimitConfigRatelimit{
					Spec: v1alpha1.GlobalRateLimitConfigRatelimitSpec{
						Domain:          "global",
						FailureModeDeny: false,
						Timeout:         "10s",
						Service: v1alpha1.GlobalRateLimitConfigRatelimitSpecService{
							Address: "grpc-testing.default",
							Port:    3000,
						},
					},
				},
			},
		},
		ratelimit: v1alpha1.GlobalRateLimit{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "hello-zufardhiyaulhaq-dev",
				Namespace: "istio-system",
			},
			Spec: v1alpha1.GlobalRateLimitSpec{
				Config: "public-gateway-config",
				Selector: v1alpha1.GlobalRateLimitSelector{
					VHost: "hello.zufardhiyaulhaq.dev:443",
				},
				Matcher: []*v1alpha1.GlobalRateLimit_Action{
					{
						RequestHeaders: &v1alpha1.GlobalRateLimit_Action_RequestHeaders{
							HeaderName:    ":method",
							DescriptorKey: "hello-zufardhiyaulhaq-dev-header-method",
						},
					},
				},
			},
		},
		expectedError: false,
		expectedEnvoyFilter: clientnetworking.EnvoyFilter{
			Spec: networking.EnvoyFilter{
				ConfigPatches: []*networking.EnvoyFilter_EnvoyConfigObjectPatch{
					{
						Patch: &networking.EnvoyFilter_Patch{
							Value: &proto_types.Struct{
								Fields: map[string]*proto_types.Value{
									"route": {
										Kind: &proto_types.Value_StructValue{
											StructValue: &proto_types.Struct{
												Fields: map[string]*proto_types.Value{
													"rate_limits": {
														Kind: &proto_types.Value_ListValue{
															ListValue: &proto_types.ListValue{
																Values: []*proto_types.Value{
																	{
																		Kind: &proto_types.Value_StructValue{
																			StructValue: &proto_types.Struct{
																				Fields: map[string]*proto_types.Value{
																					"actions": {
																						Kind: &proto_types.Value_ListValue{
																							ListValue: &proto_types.ListValue{
																								Values: []*proto_types.Value{
																									{
																										Kind: &proto_types.Value_StructValue{
																											StructValue: &proto_types.Struct{
																												Fields: map[string]*proto_types.Value{
																													"request_headers": {
																														Kind: &proto_types.Value_StructValue{
																															StructValue: &proto_types.Struct{
																																Fields: map[string]*proto_types.Value{
																																	"descriptor_key": {
																																		Kind: &proto_types.Value_StringValue{
																																			StringValue: "hello-zufardhiyaulhaq-dev-header-method",
																																		},
																																	},
																																	"header_name": {
																																		Kind: &proto_types.Value_StringValue{
																																			StringValue: ":method",
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
	{
		name: "given correct ratelimit with remote address",
		config: v1alpha1.GlobalRateLimitConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "public-gateway-config",
				Namespace: "istio-system",
			},
			Spec: v1alpha1.GlobalRateLimitConfigSpec{
				Type: "gateway",
				Selector: v1alpha1.GlobalRateLimitConfigSelector{
					IstioVersion: []string{"1.9"},
				},
				Ratelimit: v1alpha1.GlobalRateLimitConfigRatelimit{
					Spec: v1alpha1.GlobalRateLimitConfigRatelimitSpec{
						Domain:          "global",
						FailureModeDeny: false,
						Timeout:         "10s",
						Service: v1alpha1.GlobalRateLimitConfigRatelimitSpecService{
							Address: "grpc-testing.default",
							Port:    3000,
						},
					},
				},
			},
		},
		ratelimit: v1alpha1.GlobalRateLimit{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "hello-zufardhiyaulhaq-dev",
				Namespace: "istio-system",
			},
			Spec: v1alpha1.GlobalRateLimitSpec{
				Config: "public-gateway-config",
				Selector: v1alpha1.GlobalRateLimitSelector{
					VHost: "hello.zufardhiyaulhaq.dev:443",
				},
				Matcher: []*v1alpha1.GlobalRateLimit_Action{
					{
						RemoteAddress: &v1alpha1.GlobalRateLimit_Action_RemoteAddress{},
					},
				},
			},
		},
		expectedError: false,
		expectedEnvoyFilter: clientnetworking.EnvoyFilter{
			Spec: networking.EnvoyFilter{
				ConfigPatches: []*networking.EnvoyFilter_EnvoyConfigObjectPatch{
					{
						Patch: &networking.EnvoyFilter_Patch{
							Value: &proto_types.Struct{
								Fields: map[string]*proto_types.Value{
									"route": {
										Kind: &proto_types.Value_StructValue{
											StructValue: &proto_types.Struct{
												Fields: map[string]*proto_types.Value{
													"rate_limits": {
														Kind: &proto_types.Value_ListValue{
															ListValue: &proto_types.ListValue{
																Values: []*proto_types.Value{
																	{
																		Kind: &proto_types.Value_StructValue{
																			StructValue: &proto_types.Struct{
																				Fields: map[string]*proto_types.Value{
																					"actions": {
																						Kind: &proto_types.Value_ListValue{
																							ListValue: &proto_types.ListValue{
																								Values: []*proto_types.Value{
																									{
																										Kind: &proto_types.Value_StructValue{
																											StructValue: &proto_types.Struct{
																												Fields: map[string]*proto_types.Value{
																													"remote_address": {
																														Kind: &proto_types.Value_StructValue{
																															StructValue: &proto_types.Struct{
																																Fields: map[string]*proto_types.Value{},
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
	{
		name: "given correct ratelimit with generic key",
		config: v1alpha1.GlobalRateLimitConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "public-gateway-config",
				Namespace: "istio-system",
			},
			Spec: v1alpha1.GlobalRateLimitConfigSpec{
				Type: "gateway",
				Selector: v1alpha1.GlobalRateLimitConfigSelector{
					IstioVersion: []string{"1.9"},
				},
				Ratelimit: v1alpha1.GlobalRateLimitConfigRatelimit{
					Spec: v1alpha1.GlobalRateLimitConfigRatelimitSpec{
						Domain:          "global",
						FailureModeDeny: false,
						Timeout:         "10s",
						Service: v1alpha1.GlobalRateLimitConfigRatelimitSpecService{
							Address: "grpc-testing.default",
							Port:    3000,
						},
					},
				},
			},
		},
		ratelimit: v1alpha1.GlobalRateLimit{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "hello-zufardhiyaulhaq-dev",
				Namespace: "istio-system",
			},
			Spec: v1alpha1.GlobalRateLimitSpec{
				Config: "public-gateway-config",
				Selector: v1alpha1.GlobalRateLimitSelector{
					VHost: "hello.zufardhiyaulhaq.dev:443",
				},
				Matcher: []*v1alpha1.GlobalRateLimit_Action{
					{
						GenericKey: &v1alpha1.GlobalRateLimit_Action_GenericKey{
							DescriptorKey:   swag.String("foo"),
							DescriptorValue: "bar",
						},
					},
				},
			},
		},
		expectedError: false,
		expectedEnvoyFilter: clientnetworking.EnvoyFilter{
			Spec: networking.EnvoyFilter{
				ConfigPatches: []*networking.EnvoyFilter_EnvoyConfigObjectPatch{
					{
						Patch: &networking.EnvoyFilter_Patch{
							Value: &proto_types.Struct{
								Fields: map[string]*proto_types.Value{
									"route": {
										Kind: &proto_types.Value_StructValue{
											StructValue: &proto_types.Struct{
												Fields: map[string]*proto_types.Value{
													"rate_limits": {
														Kind: &proto_types.Value_ListValue{
															ListValue: &proto_types.ListValue{
																Values: []*proto_types.Value{
																	{
																		Kind: &proto_types.Value_StructValue{
																			StructValue: &proto_types.Struct{
																				Fields: map[string]*proto_types.Value{
																					"actions": {
																						Kind: &proto_types.Value_ListValue{
																							ListValue: &proto_types.ListValue{
																								Values: []*proto_types.Value{
																									{
																										Kind: &proto_types.Value_StructValue{
																											StructValue: &proto_types.Struct{
																												Fields: map[string]*proto_types.Value{
																													"generic_key": {
																														Kind: &proto_types.Value_StructValue{
																															StructValue: &proto_types.Struct{
																																Fields: map[string]*proto_types.Value{
																																	"descriptor_key": {
																																		Kind: &proto_types.Value_StringValue{
																																			StringValue: "foo",
																																		},
																																	},
																																	"descriptor_value": {
																																		Kind: &proto_types.Value_StringValue{
																																			StringValue: "bar",
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
	{
		name: "given correct ratelimit with header value match",
		config: v1alpha1.GlobalRateLimitConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "public-gateway-config",
				Namespace: "istio-system",
			},
			Spec: v1alpha1.GlobalRateLimitConfigSpec{
				Type: "gateway",
				Selector: v1alpha1.GlobalRateLimitConfigSelector{
					IstioVersion: []string{"1.9"},
				},
				Ratelimit: v1alpha1.GlobalRateLimitConfigRatelimit{
					Spec: v1alpha1.GlobalRateLimitConfigRatelimitSpec{
						Domain:          "global",
						FailureModeDeny: false,
						Timeout:         "10s",
						Service: v1alpha1.GlobalRateLimitConfigRatelimitSpecService{
							Address: "grpc-testing.default",
							Port:    3000,
						},
					},
				},
			},
		},
		ratelimit: v1alpha1.GlobalRateLimit{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "hello-zufardhiyaulhaq-dev",
				Namespace: "istio-system",
			},
			Spec: v1alpha1.GlobalRateLimitSpec{
				Config: "public-gateway-config",
				Selector: v1alpha1.GlobalRateLimitSelector{
					VHost: "hello.zufardhiyaulhaq.dev:443",
				},
				Matcher: []*v1alpha1.GlobalRateLimit_Action{
					{
						HeaderValueMatch: &v1alpha1.GlobalRateLimit_Action_HeaderValueMatch{
							DescriptorValue: "foo",
							ExpectMatch:     swag.Bool(true),
							Headers: []*v1alpha1.GlobalRateLimit_Action_HeaderValueMatch_HeaderMatcher{
								{
									Name:       "x-header-foo",
									ExactMatch: "foo",
								},
							},
						},
					},
				},
			},
		},
		expectedError: false,
		expectedEnvoyFilter: clientnetworking.EnvoyFilter{
			Spec: networking.EnvoyFilter{
				ConfigPatches: []*networking.EnvoyFilter_EnvoyConfigObjectPatch{
					{
						Patch: &networking.EnvoyFilter_Patch{
							Value: &proto_types.Struct{
								Fields: map[string]*proto_types.Value{
									"route": {
										Kind: &proto_types.Value_StructValue{
											StructValue: &proto_types.Struct{
												Fields: map[string]*proto_types.Value{
													"rate_limits": {
														Kind: &proto_types.Value_ListValue{
															ListValue: &proto_types.ListValue{
																Values: []*proto_types.Value{
																	{
																		Kind: &proto_types.Value_StructValue{
																			StructValue: &proto_types.Struct{
																				Fields: map[string]*proto_types.Value{
																					"actions": {
																						Kind: &proto_types.Value_ListValue{
																							ListValue: &proto_types.ListValue{
																								Values: []*proto_types.Value{
																									{
																										Kind: &proto_types.Value_StructValue{
																											StructValue: &proto_types.Struct{
																												Fields: map[string]*proto_types.Value{
																													"header_value_match": {
																														Kind: &proto_types.Value_StructValue{
																															StructValue: &proto_types.Struct{
																																Fields: map[string]*proto_types.Value{
																																	"descriptor_value": {
																																		Kind: &proto_types.Value_StringValue{
																																			StringValue: "foo",
																																		},
																																	},
																																	"expect_match": {
																																		Kind: &proto_types.Value_BoolValue{
																																			BoolValue: true,
																																		},
																																	},
																																	"headers": {
																																		Kind: &proto_types.Value_ListValue{
																																			ListValue: &proto_types.ListValue{
																																				Values: []*proto_types.Value{
																																					{
																																						Kind: &proto_types.Value_StructValue{
																																							StructValue: &proto_types.Struct{
																																								Fields: map[string]*proto_types.Value{
																																									"exact_match": {
																																										Kind: &proto_types.Value_StringValue{
																																											StringValue: "foo",
																																										},
																																									},
																																									"name": {
																																										Kind: &proto_types.Value_StringValue{
																																											StringValue: "x-header-foo",
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
		},
	},
}

func TestNewV3GatewayBuilder(t *testing.T) {
	for _, test := range V3GatewayBuilderTestGrid {
		t.Run(test.name, func(t *testing.T) {
			envoyfilter, err := ratelimit.NewV3GatewayBuilder(test.config, test.ratelimit, "1.9").
				Build()

			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.ratelimit.Name+"-"+"1.9", envoyfilter.Name)
				assert.Equal(t, test.ratelimit.Namespace, envoyfilter.Namespace)

				// match value generated
				assert.Equal(t, test.expectedEnvoyFilter.Spec.ConfigPatches[0].Patch.Value, envoyfilter.Spec.ConfigPatches[0].Patch.Value)
			}
		})
	}
}
