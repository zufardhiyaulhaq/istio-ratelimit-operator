package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/service"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/types"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"google.golang.org/protobuf/proto"
)

func TestNewStatsdRegexMatcherFromGlobalRateLimitMatcher(t *testing.T) {
	testCases := []struct {
		matchers []*v1alpha1.GlobalRateLimit_Action
		expected string
	}{
		{
			matchers: []*v1alpha1.GlobalRateLimit_Action{
				{
					RequestHeaders: &v1alpha1.GlobalRateLimit_Action_RequestHeaders{
						HeaderName:    "User-Agent",
						DescriptorKey: "user-agent",
						SkipIfAbsent:  false,
					},
				},
				{
					RemoteAddress: &v1alpha1.GlobalRateLimit_Action_RemoteAddress{},
				},
			},
			expected: "user-agent.remote_address",
		},
		{
			matchers: []*v1alpha1.GlobalRateLimit_Action{
				{
					RequestHeaders: &v1alpha1.GlobalRateLimit_Action_RequestHeaders{
						HeaderName:    "Authorization",
						DescriptorKey: "authorization",
						SkipIfAbsent:  true,
					},
				},
				{
					HeaderValueMatch: &v1alpha1.GlobalRateLimit_Action_HeaderValueMatch{
						DescriptorValue: "content-type",
						ExpectMatch:     proto.Bool(false),
						Headers: []*v1alpha1.GlobalRateLimit_Action_HeaderValueMatch_HeaderMatcher{
							{
								Name:       "Content-Type",
								ExactMatch: "application/json",
							},
						},
					},
				},
			},
			expected: "authorization.header_match_content-type",
		},
		{
			matchers: []*v1alpha1.GlobalRateLimit_Action{},
			expected: "",
		},
	}

	for _, tc := range testCases {
		actual := service.NewStatsdRegexMatcherFromGlobalRateLimitMatcher(tc.matchers)
		if actual != tc.expected {
			t.Errorf("Expected regex '%s' but got '%s' for matchers %+v", tc.expected, actual, tc.matchers)
		}
	}
}

func TestNewStatsdConfigBuilder(t *testing.T) {
	builder := service.NewStatsdConfigBuilder()
	assert.NotNil(t, builder)
	assert.Equal(t, "", builder.Config)
	assert.Equal(t, v1alpha1.RateLimitService{}, builder.RateLimitService)
}

func TestStatsdConfigBuilder_SetRateLimitService(t *testing.T) {
	rateLimitService := v1alpha1.RateLimitService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-statsd",
			Namespace: "test-namespace",
		},
	}

	builder := service.NewStatsdConfigBuilder().SetRateLimitService(rateLimitService)

	assert.Equal(t, rateLimitService, builder.RateLimitService)
}

func TestStatsdConfigBuilder_SetRateLimitService_Chaining(t *testing.T) {
	rateLimitService := v1alpha1.RateLimitService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-statsd",
			Namespace: "test-namespace",
		},
	}

	builder := service.NewStatsdConfigBuilder()
	returnedBuilder := builder.SetRateLimitService(rateLimitService)

	// Verify method chaining returns the same builder
	assert.Same(t, builder, returnedBuilder)
}

func TestStatsdConfigBuilder_SetConfig(t *testing.T) {
	config := "mappings:\n  - name: test"

	builder := service.NewStatsdConfigBuilder().SetConfig(config)

	assert.Equal(t, config, builder.Config)
}

func TestStatsdConfigBuilder_SetConfig_Chaining(t *testing.T) {
	config := "mappings:\n  - name: test"

	builder := service.NewStatsdConfigBuilder()
	returnedBuilder := builder.SetConfig(config)

	// Verify method chaining returns the same builder
	assert.Same(t, builder, returnedBuilder)
}

func TestStatsdConfigBuilder_BuildLabels(t *testing.T) {
	testCases := []struct {
		name             string
		rateLimitService v1alpha1.RateLimitService
		expectedLabels   map[string]string
	}{
		{
			name: "basic labels",
			rateLimitService: v1alpha1.RateLimitService{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-ratelimit",
					Namespace: "default",
				},
			},
			expectedLabels: map[string]string{
				"app.kubernetes.io/name":       "my-ratelimit-statsd-config",
				"app.kubernetes.io/managed-by": "istio-rateltimit-operator",
				"app.kubernetes.io/created-by": "my-ratelimit",
			},
		},
		{
			name: "different name",
			rateLimitService: v1alpha1.RateLimitService{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "prod-statsd",
					Namespace: "production",
				},
			},
			expectedLabels: map[string]string{
				"app.kubernetes.io/name":       "prod-statsd-statsd-config",
				"app.kubernetes.io/managed-by": "istio-rateltimit-operator",
				"app.kubernetes.io/created-by": "prod-statsd",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := service.NewStatsdConfigBuilder().SetRateLimitService(tc.rateLimitService)
			labels := builder.BuildLabels()

			assert.Equal(t, tc.expectedLabels, labels)
		})
	}
}

func TestStatsdConfigBuilder_Build(t *testing.T) {
	testCases := []struct {
		name              string
		rateLimitService  v1alpha1.RateLimitService
		config            string
		expectedConfigMap *corev1.ConfigMap
		expectError       bool
	}{
		{
			name: "basic statsd config build",
			rateLimitService: v1alpha1.RateLimitService{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "statsd-cfg",
					Namespace: "default",
				},
			},
			config: "mappings:\n  - name: test_metric",
			expectedConfigMap: &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "statsd-cfg-statsd-config",
					Namespace: "default",
					Labels: map[string]string{
						"app.kubernetes.io/name":       "statsd-cfg-statsd-config",
						"app.kubernetes.io/managed-by": "istio-rateltimit-operator",
						"app.kubernetes.io/created-by": "statsd-cfg",
					},
				},
				Data: map[string]string{
					"statsd.mappingConf": "mappings:\n  - name: test_metric",
				},
			},
			expectError: false,
		},
		{
			name: "statsd config in different namespace",
			rateLimitService: v1alpha1.RateLimitService{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "prod-statsd",
					Namespace: "production",
				},
			},
			config: "mappings:\n  - name: prod_metric",
			expectedConfigMap: &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "prod-statsd-statsd-config",
					Namespace: "production",
					Labels: map[string]string{
						"app.kubernetes.io/name":       "prod-statsd-statsd-config",
						"app.kubernetes.io/managed-by": "istio-rateltimit-operator",
						"app.kubernetes.io/created-by": "prod-statsd",
					},
				},
				Data: map[string]string{
					"statsd.mappingConf": "mappings:\n  - name: prod_metric",
				},
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			configMap, err := service.NewStatsdConfigBuilder().
				SetRateLimitService(tc.rateLimitService).
				SetConfig(tc.config).
				Build()

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedConfigMap, configMap)
			}
		})
	}
}

func TestNewStatsdConfig(t *testing.T) {
	identifier := "test-identifier"

	testCases := []struct {
		name                    string
		rateLimitServiceName    string
		globalRateLimitDomain   string
		globalRateLimitList     []v1alpha1.GlobalRateLimit
		expectedMappingsCount   int
		expectError             bool
	}{
		{
			name:                  "rate limit with identifier",
			rateLimitServiceName:  "my-ratelimit-service",
			globalRateLimitDomain: "my-domain",
			globalRateLimitList: []v1alpha1.GlobalRateLimit{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-rate-limit",
					},
					Spec: v1alpha1.GlobalRateLimitSpec{
						Identifier: &identifier,
						Matcher: []*v1alpha1.GlobalRateLimit_Action{
							{
								RequestHeaders: &v1alpha1.GlobalRateLimit_Action_RequestHeaders{
									HeaderName:    "x-api-key",
									DescriptorKey: "api_key",
								},
							},
						},
					},
				},
			},
			expectedMappingsCount: 12, // 6 metric mappings + 6 default mappings
			expectError:           false,
		},
		{
			name:                  "rate limit without identifier",
			rateLimitServiceName:  "my-ratelimit-service",
			globalRateLimitDomain: "my-domain",
			globalRateLimitList: []v1alpha1.GlobalRateLimit{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-rate-limit",
					},
					Spec: v1alpha1.GlobalRateLimitSpec{
						Identifier: nil,
						Matcher: []*v1alpha1.GlobalRateLimit_Action{
							{
								RequestHeaders: &v1alpha1.GlobalRateLimit_Action_RequestHeaders{
									HeaderName:    "x-api-key",
									DescriptorKey: "api_key",
								},
							},
						},
					},
				},
			},
			expectedMappingsCount: 6, // only default mappings
			expectError:           false,
		},
		{
			name:                  "empty rate limit list",
			rateLimitServiceName:  "my-ratelimit-service",
			globalRateLimitDomain: "my-domain",
			globalRateLimitList:   []v1alpha1.GlobalRateLimit{},
			expectedMappingsCount: 6, // only default mappings
			expectError:           false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			metricMapper, err := service.NewStatsdConfig(tc.rateLimitServiceName, tc.globalRateLimitDomain, tc.globalRateLimitList)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, metricMapper.Mappings, tc.expectedMappingsCount)
			}
		})
	}
}

func TestNewMetricMappingFromGlobalRateLimit(t *testing.T) {
	identifier := "test-identifier"

	testCases := []struct {
		name                  string
		rateLimitServiceName  string
		globalRateLimitDomain string
		globalRateLimit       v1alpha1.GlobalRateLimit
		expectedMappingsCount int
		expectError           bool
	}{
		{
			name:                  "basic metric mapping",
			rateLimitServiceName:  "my-ratelimit-service",
			globalRateLimitDomain: "my-domain",
			globalRateLimit: v1alpha1.GlobalRateLimit{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-rate-limit",
				},
				Spec: v1alpha1.GlobalRateLimitSpec{
					Identifier: &identifier,
					Matcher: []*v1alpha1.GlobalRateLimit_Action{
						{
							RequestHeaders: &v1alpha1.GlobalRateLimit_Action_RequestHeaders{
								HeaderName:    "x-api-key",
								DescriptorKey: "api_key",
							},
						},
					},
				},
			},
			expectedMappingsCount: 6, // near_limit, over_limit, over_limit_with_local_cache, total_hits, within_limit, shadow_mode
			expectError:           false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mappings, err := service.NewMetricMappingFromGlobalRateLimit(tc.rateLimitServiceName, tc.globalRateLimitDomain, tc.globalRateLimit)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, mappings, tc.expectedMappingsCount)

				// Verify each mapping has the correct labels
				for _, mapping := range mappings {
					assert.Equal(t, *tc.globalRateLimit.Spec.Identifier, mapping.Labels["identifier"])
					assert.Equal(t, tc.rateLimitServiceName, mapping.Labels["rate_limit_service_name"])
					assert.Equal(t, tc.globalRateLimit.Name, mapping.Labels["global_rate_limit_name"])
				}
			}
		})
	}
}

func TestNewDefaultMetricMapping(t *testing.T) {
	mappings := service.NewDefaultMetricMapping()

	assert.Len(t, mappings, 6)

	// Verify mapping names
	expectedNames := []string{
		"ratelimit_service_should_rate_limit_error",
		"ratelimit_service_total_requests",
		"ratelimit_service_response_time_seconds",
		"ratelimit_service_config_load_success",
		"ratelimit_service_config_load_error",
		"ratelimit_service_global_shadow_mode",
	}

	for i, mapping := range mappings {
		assert.Equal(t, expectedNames[i], mapping.Name)
	}
}

func TestNewStatsdRegexMatcherFromGlobalRateLimitMatcher_GenericKeyWithDescriptorKey(t *testing.T) {
	descriptorKey := "my-key"
	matchers := []*v1alpha1.GlobalRateLimit_Action{
		{
			GenericKey: &v1alpha1.GlobalRateLimit_Action_GenericKey{
				DescriptorKey:   &descriptorKey,
				DescriptorValue: "my-value",
			},
		},
	}

	result := service.NewStatsdRegexMatcherFromGlobalRateLimitMatcher(matchers)

	assert.Equal(t, "my-key_my-value", result)
}

func TestNewStatsdRegexMatcherFromGlobalRateLimitMatcher_GenericKeyWithoutDescriptorKey(t *testing.T) {
	matchers := []*v1alpha1.GlobalRateLimit_Action{
		{
			GenericKey: &v1alpha1.GlobalRateLimit_Action_GenericKey{
				DescriptorKey:   nil,
				DescriptorValue: "my-value",
			},
		},
	}

	result := service.NewStatsdRegexMatcherFromGlobalRateLimitMatcher(matchers)

	assert.Equal(t, "generic_key_my-value", result)
}

func TestMetricMappingTypes(t *testing.T) {
	// Test that the types package has the expected constants
	assert.Equal(t, types.ObserverType("histogram"), types.ObserverTypeHistogram)
	assert.Equal(t, types.MetricType("counter"), types.MetricTypeCounter)
}

func TestMetricMappingStructure(t *testing.T) {
	identifier := "test-id"

	globalRateLimit := v1alpha1.GlobalRateLimit{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-rate-limit",
		},
		Spec: v1alpha1.GlobalRateLimitSpec{
			Identifier: &identifier,
			Matcher: []*v1alpha1.GlobalRateLimit_Action{
				{
					RequestHeaders: &v1alpha1.GlobalRateLimit_Action_RequestHeaders{
						HeaderName:    "x-test",
						DescriptorKey: "test_key",
					},
				},
			},
		},
	}

	mappings, err := service.NewMetricMappingFromGlobalRateLimit("test-service", "test-domain", globalRateLimit)

	assert.NoError(t, err)
	assert.Len(t, mappings, 6)

	// Check specific mapping properties
	nearLimitMapping := mappings[0]
	assert.Equal(t, "ratelimit_service_rate_limit_near_limit", nearLimitMapping.Name)
	assert.Contains(t, nearLimitMapping.Match, "test-domain")
	assert.Contains(t, nearLimitMapping.Match, "test_key")
	assert.Contains(t, nearLimitMapping.Match, "near_limit")
	assert.Equal(t, types.ObserverTypeHistogram, nearLimitMapping.TimerType)
	assert.Equal(t, "test-id", nearLimitMapping.Labels["identifier"])
	assert.Equal(t, "test-service", nearLimitMapping.Labels["rate_limit_service_name"])
	assert.Equal(t, "test-rate-limit", nearLimitMapping.Labels["global_rate_limit_name"])
}
