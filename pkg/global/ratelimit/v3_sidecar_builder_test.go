package ratelimit_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/global/ratelimit"
	"google.golang.org/protobuf/proto"

	proto_types "github.com/golang/protobuf/ptypes/struct"
	networking "istio.io/api/networking/v1alpha3"
	clientnetworking "istio.io/client-go/pkg/apis/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type V3SidecarBuilderTestCase struct {
	name                string
	config              v1alpha1.GlobalRateLimitConfig
	ratelimit           v1alpha1.GlobalRateLimit
	expectedError       bool
	expectedEnvoyFilter clientnetworking.EnvoyFilter
}

var V3SidecarBuilderTestGrid = []V3SidecarBuilderTestCase{
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
}

func TestNewV3SidecarBuilder(t *testing.T) {
	for _, test := range V3SidecarBuilderTestGrid {
		t.Run(test.name, func(t *testing.T) {
			envoyfilter, err := ratelimit.NewV3SidecarBuilder(test.config, test.ratelimit, "1.9").
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

func TestV3SidecarBuilder_VerifyEnvoyFilterStructure(t *testing.T) {
	config := v1alpha1.GlobalRateLimitConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sidecar-config",
			Namespace: "default",
		},
		Spec: v1alpha1.GlobalRateLimitConfigSpec{
			Type: "sidecar",
			Selector: v1alpha1.GlobalRateLimitConfigSelector{
				Labels:       map[string]string{"app": "myapp", "version": "v1"},
				IstioVersion: []string{"1.10"},
			},
		},
	}

	rateLimitObj := v1alpha1.GlobalRateLimit{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-ratelimit",
			Namespace: "default",
		},
		Spec: v1alpha1.GlobalRateLimitSpec{
			Config: "sidecar-config",
			Selector: v1alpha1.GlobalRateLimitSelector{
				VHost: "inbound|8080||",
			},
			Matcher: []*v1alpha1.GlobalRateLimit_Action{
				{
					RemoteAddress: &v1alpha1.GlobalRateLimit_Action_RemoteAddress{},
				},
			},
		},
	}

	envoyfilter, err := ratelimit.NewV3SidecarBuilder(config, rateLimitObj, "1.10").Build()

	assert.NoError(t, err)
	assert.NotNil(t, envoyfilter)

	// Verify TypeMeta
	assert.Equal(t, "EnvoyFilter", envoyfilter.Kind)
	assert.Equal(t, "networking.istio.io/v1alpha3", envoyfilter.APIVersion)

	// Verify ObjectMeta
	assert.Equal(t, "test-ratelimit-1.10", envoyfilter.Name)
	assert.Equal(t, "default", envoyfilter.Namespace)
	assert.Equal(t, map[string]string{"istio/version": "1.10"}, envoyfilter.Labels)

	// Verify WorkloadSelector
	assert.Equal(t, config.Spec.Selector.Labels, envoyfilter.Spec.WorkloadSelector.Labels)

	// Verify ConfigPatches
	assert.Len(t, envoyfilter.Spec.ConfigPatches, 1)
	patch := envoyfilter.Spec.ConfigPatches[0]

	// Verify ApplyTo
	assert.Equal(t, networking.EnvoyFilter_HTTP_ROUTE, patch.ApplyTo)

	// Verify Match context - SIDECAR_INBOUND for sidecar
	assert.Equal(t, networking.EnvoyFilter_SIDECAR_INBOUND, patch.Match.Context)

	// Verify Proxy match
	assert.NotNil(t, patch.Match.Proxy)
	assert.Equal(t, `^1\.10.*`, patch.Match.Proxy.ProxyVersion)

	// Verify RouteConfiguration match
	routeConfig := patch.Match.GetRouteConfiguration()
	assert.NotNil(t, routeConfig)
	assert.Equal(t, "inbound|8080||", routeConfig.Vhost.Name)
	assert.Equal(t, networking.EnvoyFilter_RouteConfigurationMatch_RouteMatch_ANY, routeConfig.Vhost.Route.Action)

	// Verify Patch operation
	assert.Equal(t, networking.EnvoyFilter_Patch_MERGE, patch.Patch.Operation)
}

func TestV3SidecarBuilder_WithRemoteAddress(t *testing.T) {
	config := v1alpha1.GlobalRateLimitConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sidecar-config",
			Namespace: "default",
		},
		Spec: v1alpha1.GlobalRateLimitConfigSpec{
			Type: "sidecar",
			Selector: v1alpha1.GlobalRateLimitConfigSelector{
				Labels:       map[string]string{"app": "myapp"},
				IstioVersion: []string{"1.9"},
			},
		},
	}

	rateLimitObj := v1alpha1.GlobalRateLimit{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-ratelimit",
			Namespace: "default",
		},
		Spec: v1alpha1.GlobalRateLimitSpec{
			Config: "sidecar-config",
			Selector: v1alpha1.GlobalRateLimitSelector{
				VHost: "inbound|8080||",
			},
			Matcher: []*v1alpha1.GlobalRateLimit_Action{
				{
					RemoteAddress: &v1alpha1.GlobalRateLimit_Action_RemoteAddress{},
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
	}

	envoyfilter, err := ratelimit.NewV3SidecarBuilder(config, rateLimitObj, "1.9").Build()

	assert.NoError(t, err)
	assert.NotNil(t, envoyfilter)
	assert.True(t, proto.Equal(expectedValue, envoyfilter.Spec.ConfigPatches[0].Patch.Value))
}

func TestV3SidecarBuilder_WithGenericKey(t *testing.T) {
	config := v1alpha1.GlobalRateLimitConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sidecar-config",
			Namespace: "default",
		},
		Spec: v1alpha1.GlobalRateLimitConfigSpec{
			Type: "sidecar",
			Selector: v1alpha1.GlobalRateLimitConfigSelector{
				Labels:       map[string]string{"app": "myapp"},
				IstioVersion: []string{"1.9"},
			},
		},
	}

	descriptorKey := "service-key"
	rateLimitObj := v1alpha1.GlobalRateLimit{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-ratelimit",
			Namespace: "default",
		},
		Spec: v1alpha1.GlobalRateLimitSpec{
			Config: "sidecar-config",
			Selector: v1alpha1.GlobalRateLimitSelector{
				VHost: "inbound|8080||",
			},
			Matcher: []*v1alpha1.GlobalRateLimit_Action{
				{
					GenericKey: &v1alpha1.GlobalRateLimit_Action_GenericKey{
						DescriptorKey:   &descriptorKey,
						DescriptorValue: "my-service",
					},
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
																							"generic_key": {
																								Kind: &proto_types.Value_StructValue{
																									StructValue: &proto_types.Struct{
																										Fields: map[string]*proto_types.Value{
																											"descriptor_key": {
																												Kind: &proto_types.Value_StringValue{
																													StringValue: "service-key",
																												},
																											},
																											"descriptor_value": {
																												Kind: &proto_types.Value_StringValue{
																													StringValue: "my-service",
																												},
																											},
																										},
																									},
																								},
																							},
																						},
																					},
																				},
																			},
																		},
																	},
																},
															},
														},
													},
												},
											},
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

	envoyfilter, err := ratelimit.NewV3SidecarBuilder(config, rateLimitObj, "1.9").Build()

	assert.NoError(t, err)
	assert.NotNil(t, envoyfilter)
	assert.True(t, proto.Equal(expectedValue, envoyfilter.Spec.ConfigPatches[0].Patch.Value))
}

func TestV3SidecarBuilder_DifferentVersions(t *testing.T) {
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
			Name:      "sidecar-config",
			Namespace: "default",
		},
		Spec: v1alpha1.GlobalRateLimitConfigSpec{
			Type: "sidecar",
			Selector: v1alpha1.GlobalRateLimitConfigSelector{
				Labels: map[string]string{"app": "myapp"},
			},
		},
	}

	rateLimitObj := v1alpha1.GlobalRateLimit{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-ratelimit",
			Namespace: "default",
		},
		Spec: v1alpha1.GlobalRateLimitSpec{
			Config: "sidecar-config",
			Selector: v1alpha1.GlobalRateLimitSelector{
				VHost: "inbound|8080||",
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
			envoyfilter, err := ratelimit.NewV3SidecarBuilder(config, rateLimitObj, tt.version).Build()

			assert.NoError(t, err)
			assert.NotNil(t, envoyfilter)
			assert.Equal(t, tt.expectedName, envoyfilter.Name)
			assert.Equal(t, map[string]string{"istio/version": tt.version}, envoyfilter.Labels)
			assert.Equal(t, tt.expectedProxyVersion, envoyfilter.Spec.ConfigPatches[0].Match.Proxy.ProxyVersion)
		})
	}
}

func TestV3SidecarBuilder_WithMultipleActions(t *testing.T) {
	config := v1alpha1.GlobalRateLimitConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sidecar-config",
			Namespace: "default",
		},
		Spec: v1alpha1.GlobalRateLimitConfigSpec{
			Type: "sidecar",
			Selector: v1alpha1.GlobalRateLimitConfigSelector{
				Labels:       map[string]string{"app": "myapp"},
				IstioVersion: []string{"1.9"},
			},
		},
	}

	rateLimitObj := v1alpha1.GlobalRateLimit{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-ratelimit",
			Namespace: "default",
		},
		Spec: v1alpha1.GlobalRateLimitSpec{
			Config: "sidecar-config",
			Selector: v1alpha1.GlobalRateLimitSelector{
				VHost: "inbound|8080||",
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
				{
					SourceCluster: &v1alpha1.GlobalRateLimit_Action_SourceCluster{},
				},
			},
		},
	}

	envoyfilter, err := ratelimit.NewV3SidecarBuilder(config, rateLimitObj, "1.9").Build()

	assert.NoError(t, err)
	assert.NotNil(t, envoyfilter)

	// Verify there are three actions in the rate_limits
	patchValue := envoyfilter.Spec.ConfigPatches[0].Patch.Value
	routeField := patchValue.GetFields()["route"]
	rateLimitsField := routeField.GetStructValue().GetFields()["rate_limits"]
	actionsField := rateLimitsField.GetListValue().GetValues()[0].GetStructValue().GetFields()["actions"]
	actions := actionsField.GetListValue().GetValues()

	assert.Len(t, actions, 3)
}

func TestV3SidecarBuilder_WithSourceCluster(t *testing.T) {
	config := v1alpha1.GlobalRateLimitConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sidecar-config",
			Namespace: "default",
		},
		Spec: v1alpha1.GlobalRateLimitConfigSpec{
			Type: "sidecar",
			Selector: v1alpha1.GlobalRateLimitConfigSelector{
				Labels:       map[string]string{"app": "myapp"},
				IstioVersion: []string{"1.9"},
			},
		},
	}

	rateLimitObj := v1alpha1.GlobalRateLimit{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-ratelimit",
			Namespace: "default",
		},
		Spec: v1alpha1.GlobalRateLimitSpec{
			Config: "sidecar-config",
			Selector: v1alpha1.GlobalRateLimitSelector{
				VHost: "inbound|8080||",
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

	envoyfilter, err := ratelimit.NewV3SidecarBuilder(config, rateLimitObj, "1.9").Build()

	assert.NoError(t, err)
	assert.NotNil(t, envoyfilter)
	assert.True(t, proto.Equal(expectedValue, envoyfilter.Spec.ConfigPatches[0].Patch.Value))
}

func TestV3SidecarBuilder_WithDestinationCluster(t *testing.T) {
	config := v1alpha1.GlobalRateLimitConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sidecar-config",
			Namespace: "default",
		},
		Spec: v1alpha1.GlobalRateLimitConfigSpec{
			Type: "sidecar",
			Selector: v1alpha1.GlobalRateLimitConfigSelector{
				Labels:       map[string]string{"app": "myapp"},
				IstioVersion: []string{"1.9"},
			},
		},
	}

	rateLimitObj := v1alpha1.GlobalRateLimit{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-ratelimit",
			Namespace: "default",
		},
		Spec: v1alpha1.GlobalRateLimitSpec{
			Config: "sidecar-config",
			Selector: v1alpha1.GlobalRateLimitSelector{
				VHost: "inbound|8080||",
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

	envoyfilter, err := ratelimit.NewV3SidecarBuilder(config, rateLimitObj, "1.9").Build()

	assert.NoError(t, err)
	assert.NotNil(t, envoyfilter)
	assert.True(t, proto.Equal(expectedValue, envoyfilter.Spec.ConfigPatches[0].Patch.Value))
}
