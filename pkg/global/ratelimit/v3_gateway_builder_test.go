package ratelimit_test

import (
	"testing"

	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/assert"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/global/ratelimit"
	"google.golang.org/protobuf/proto"

	proto_types "github.com/golang/protobuf/ptypes/struct"
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
	for i := range V3GatewayBuilderTestGrid {
		test := &V3GatewayBuilderTestGrid[i]
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
				assert.True(t, proto.Equal(test.expectedEnvoyFilter.Spec.ConfigPatches[0].Patch.Value, envoyfilter.Spec.ConfigPatches[0].Patch.Value))
			}
		})
	}
}

func TestV3GatewayBuilder_WithRouteSelector(t *testing.T) {
	routeName := "test-route"
	config := v1alpha1.GlobalRateLimitConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "public-gateway-config",
			Namespace: "istio-system",
		},
		Spec: v1alpha1.GlobalRateLimitConfigSpec{
			Type: "gateway",
			Selector: v1alpha1.GlobalRateLimitConfigSelector{
				Labels:       map[string]string{"app": "gateway"},
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
			Config: "public-gateway-config",
			Selector: v1alpha1.GlobalRateLimitSelector{
				VHost: "example.com:443",
				Route: &routeName,
			},
			Matcher: []*v1alpha1.GlobalRateLimit_Action{
				{
					RemoteAddress: &v1alpha1.GlobalRateLimit_Action_RemoteAddress{},
				},
			},
		},
	}

	envoyfilter, err := ratelimit.NewV3GatewayBuilder(config, rateLimitObj, "1.9").Build()

	assert.NoError(t, err)
	assert.NotNil(t, envoyfilter)
	assert.Equal(t, "test-ratelimit-1.9", envoyfilter.Name)

	// Verify the route configuration match includes the route name
	routeConfig := envoyfilter.Spec.ConfigPatches[0].Match.GetRouteConfiguration()
	assert.NotNil(t, routeConfig)
	assert.NotNil(t, routeConfig.Vhost)
	assert.NotNil(t, routeConfig.Vhost.Route)
	assert.Equal(t, routeName, routeConfig.Vhost.Route.Name)
}

func TestV3GatewayBuilder_WithSourceCluster(t *testing.T) {
	config := v1alpha1.GlobalRateLimitConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "public-gateway-config",
			Namespace: "istio-system",
		},
		Spec: v1alpha1.GlobalRateLimitConfigSpec{
			Type: "gateway",
			Selector: v1alpha1.GlobalRateLimitConfigSelector{
				Labels:       map[string]string{"app": "gateway"},
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
			Config: "public-gateway-config",
			Selector: v1alpha1.GlobalRateLimitSelector{
				VHost: "example.com:443",
			},
			Matcher: []*v1alpha1.GlobalRateLimit_Action{
				{
					SourceCluster: &v1alpha1.GlobalRateLimit_Action_SourceCluster{},
				},
			},
		},
	}

	expectedValue := &proto_types.Struct{
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
																							"source_cluster": {
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
	}

	envoyfilter, err := ratelimit.NewV3GatewayBuilder(config, rateLimitObj, "1.9").Build()

	assert.NoError(t, err)
	assert.NotNil(t, envoyfilter)
	assert.True(t, proto.Equal(expectedValue, envoyfilter.Spec.ConfigPatches[0].Patch.Value))
}

func TestV3GatewayBuilder_WithDestinationCluster(t *testing.T) {
	config := v1alpha1.GlobalRateLimitConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "public-gateway-config",
			Namespace: "istio-system",
		},
		Spec: v1alpha1.GlobalRateLimitConfigSpec{
			Type: "gateway",
			Selector: v1alpha1.GlobalRateLimitConfigSelector{
				Labels:       map[string]string{"app": "gateway"},
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
			Config: "public-gateway-config",
			Selector: v1alpha1.GlobalRateLimitSelector{
				VHost: "example.com:443",
			},
			Matcher: []*v1alpha1.GlobalRateLimit_Action{
				{
					DestinationCluster: &v1alpha1.GlobalRateLimit_Action_DestinationCluster{},
				},
			},
		},
	}

	expectedValue := &proto_types.Struct{
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
																							"destination_cluster": {
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
	}

	envoyfilter, err := ratelimit.NewV3GatewayBuilder(config, rateLimitObj, "1.9").Build()

	assert.NoError(t, err)
	assert.NotNil(t, envoyfilter)
	assert.True(t, proto.Equal(expectedValue, envoyfilter.Spec.ConfigPatches[0].Patch.Value))
}

func TestV3GatewayBuilder_WithMultipleActions(t *testing.T) {
	config := v1alpha1.GlobalRateLimitConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "public-gateway-config",
			Namespace: "istio-system",
		},
		Spec: v1alpha1.GlobalRateLimitConfigSpec{
			Type: "gateway",
			Selector: v1alpha1.GlobalRateLimitConfigSelector{
				Labels:       map[string]string{"app": "gateway"},
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
			Config: "public-gateway-config",
			Selector: v1alpha1.GlobalRateLimitSelector{
				VHost: "example.com:443",
			},
			Matcher: []*v1alpha1.GlobalRateLimit_Action{
				{
					RemoteAddress: &v1alpha1.GlobalRateLimit_Action_RemoteAddress{},
				},
				{
					RequestHeaders: &v1alpha1.GlobalRateLimit_Action_RequestHeaders{
						HeaderName:    ":path",
						DescriptorKey: "path",
					},
				},
			},
		},
	}

	envoyfilter, err := ratelimit.NewV3GatewayBuilder(config, rateLimitObj, "1.9").Build()

	assert.NoError(t, err)
	assert.NotNil(t, envoyfilter)

	// Verify there are two actions in the rate_limits
	patchValue := envoyfilter.Spec.ConfigPatches[0].Patch.Value
	routeField := patchValue.GetFields()["route"]
	rateLimitsField := routeField.GetStructValue().GetFields()["rate_limits"]
	actionsField := rateLimitsField.GetListValue().GetValues()[0].GetStructValue().GetFields()["actions"]
	actions := actionsField.GetListValue().GetValues()

	assert.Len(t, actions, 2)
}

func TestV3GatewayBuilder_VerifyEnvoyFilterStructure(t *testing.T) {
	config := v1alpha1.GlobalRateLimitConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "public-gateway-config",
			Namespace: "istio-system",
		},
		Spec: v1alpha1.GlobalRateLimitConfigSpec{
			Type: "gateway",
			Selector: v1alpha1.GlobalRateLimitConfigSelector{
				Labels:       map[string]string{"app": "gateway", "istio": "ingressgateway"},
				IstioVersion: []string{"1.10"},
			},
		},
	}

	rateLimitObj := v1alpha1.GlobalRateLimit{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-ratelimit",
			Namespace: "istio-system",
		},
		Spec: v1alpha1.GlobalRateLimitSpec{
			Config: "public-gateway-config",
			Selector: v1alpha1.GlobalRateLimitSelector{
				VHost: "example.com:443",
			},
			Matcher: []*v1alpha1.GlobalRateLimit_Action{
				{
					RemoteAddress: &v1alpha1.GlobalRateLimit_Action_RemoteAddress{},
				},
			},
		},
	}

	envoyfilter, err := ratelimit.NewV3GatewayBuilder(config, rateLimitObj, "1.10").Build()

	assert.NoError(t, err)
	assert.NotNil(t, envoyfilter)

	// Verify TypeMeta
	assert.Equal(t, "EnvoyFilter", envoyfilter.Kind)
	assert.Equal(t, "networking.istio.io/v1alpha3", envoyfilter.APIVersion)

	// Verify ObjectMeta
	assert.Equal(t, "test-ratelimit-1.10", envoyfilter.Name)
	assert.Equal(t, "istio-system", envoyfilter.Namespace)
	assert.Equal(t, map[string]string{"istio/version": "1.10"}, envoyfilter.Labels)

	// Verify WorkloadSelector
	assert.Equal(t, config.Spec.Selector.Labels, envoyfilter.Spec.WorkloadSelector.Labels)

	// Verify ConfigPatches
	assert.Len(t, envoyfilter.Spec.ConfigPatches, 1)
	patch := envoyfilter.Spec.ConfigPatches[0]

	// Verify ApplyTo
	assert.Equal(t, networking.EnvoyFilter_HTTP_ROUTE, patch.ApplyTo)

	// Verify Match context
	assert.Equal(t, networking.EnvoyFilter_GATEWAY, patch.Match.Context)

	// Verify Proxy match
	assert.NotNil(t, patch.Match.Proxy)
	assert.Equal(t, `^1\.10.*`, patch.Match.Proxy.ProxyVersion)

	// Verify RouteConfiguration match
	routeConfig := patch.Match.GetRouteConfiguration()
	assert.NotNil(t, routeConfig)
	assert.Equal(t, "example.com:443", routeConfig.Vhost.Name)
	assert.Equal(t, networking.EnvoyFilter_RouteConfigurationMatch_RouteMatch_ANY, routeConfig.Vhost.Route.Action)

	// Verify Patch operation
	assert.Equal(t, networking.EnvoyFilter_Patch_MERGE, patch.Patch.Operation)
}

func TestV3GatewayBuilder_DifferentVersions(t *testing.T) {
	tests := []struct {
		version              string
		expectedName         string
		expectedProxyVersion string
	}{
		{
			version:              "1.9",
			expectedName:         "test-ratelimit-1.9",
			expectedProxyVersion: `^1\.9.*`,
		},
		{
			version:              "1.10",
			expectedName:         "test-ratelimit-1.10",
			expectedProxyVersion: `^1\.10.*`,
		},
		{
			version:              "1.25",
			expectedName:         "test-ratelimit-1.25",
			expectedProxyVersion: `^1\.25.*`,
		},
	}

	config := v1alpha1.GlobalRateLimitConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "public-gateway-config",
			Namespace: "istio-system",
		},
		Spec: v1alpha1.GlobalRateLimitConfigSpec{
			Type: "gateway",
			Selector: v1alpha1.GlobalRateLimitConfigSelector{
				Labels: map[string]string{"app": "gateway"},
			},
		},
	}

	rateLimitObj := v1alpha1.GlobalRateLimit{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-ratelimit",
			Namespace: "istio-system",
		},
		Spec: v1alpha1.GlobalRateLimitSpec{
			Config: "public-gateway-config",
			Selector: v1alpha1.GlobalRateLimitSelector{
				VHost: "example.com:443",
			},
			Matcher: []*v1alpha1.GlobalRateLimit_Action{
				{
					RemoteAddress: &v1alpha1.GlobalRateLimit_Action_RemoteAddress{},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			envoyfilter, err := ratelimit.NewV3GatewayBuilder(config, rateLimitObj, tt.version).Build()

			assert.NoError(t, err)
			assert.NotNil(t, envoyfilter)
			assert.Equal(t, tt.expectedName, envoyfilter.Name)
			assert.Equal(t, map[string]string{"istio/version": tt.version}, envoyfilter.Labels)
			assert.Equal(t, tt.expectedProxyVersion, envoyfilter.Spec.ConfigPatches[0].Match.Proxy.ProxyVersion)
		})
	}
}

func TestV3GatewayBuilder_WithRequestHeadersSkipIfAbsent(t *testing.T) {
	config := v1alpha1.GlobalRateLimitConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "public-gateway-config",
			Namespace: "istio-system",
		},
		Spec: v1alpha1.GlobalRateLimitConfigSpec{
			Type: "gateway",
			Selector: v1alpha1.GlobalRateLimitConfigSelector{
				Labels:       map[string]string{"app": "gateway"},
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
			Config: "public-gateway-config",
			Selector: v1alpha1.GlobalRateLimitSelector{
				VHost: "example.com:443",
			},
			Matcher: []*v1alpha1.GlobalRateLimit_Action{
				{
					RequestHeaders: &v1alpha1.GlobalRateLimit_Action_RequestHeaders{
						HeaderName:    "x-custom-header",
						DescriptorKey: "custom-header",
						SkipIfAbsent:  true,
					},
				},
			},
		},
	}

	envoyfilter, err := ratelimit.NewV3GatewayBuilder(config, rateLimitObj, "1.9").Build()

	assert.NoError(t, err)
	assert.NotNil(t, envoyfilter)

	// Verify the request headers action contains skip_if_absent
	patchValue := envoyfilter.Spec.ConfigPatches[0].Patch.Value
	routeField := patchValue.GetFields()["route"]
	rateLimitsField := routeField.GetStructValue().GetFields()["rate_limits"]
	actionsField := rateLimitsField.GetListValue().GetValues()[0].GetStructValue().GetFields()["actions"]
	action := actionsField.GetListValue().GetValues()[0]
	requestHeaders := action.GetStructValue().GetFields()["request_headers"]

	assert.NotNil(t, requestHeaders)
	skipIfAbsent := requestHeaders.GetStructValue().GetFields()["skip_if_absent"]
	assert.NotNil(t, skipIfAbsent)
	assert.True(t, skipIfAbsent.GetBoolValue())
}
