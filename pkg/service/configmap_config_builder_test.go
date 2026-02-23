package service_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/service"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/types"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNewRateLimitDescriptorFromMatcher(t *testing.T) {
	fakeDescriptorKey01 := "foo"
	fakeDescriptorKey02 := "bar"
	fakeExpectedMatch := true

	type args struct {
		matchers   []*v1alpha1.GlobalRateLimit_Action
		limit      *v1alpha1.GlobalRateLimit_Limit
		shadowMode bool
	}
	tests := []struct {
		name    string
		args    args
		want    []types.RateLimit_Service_Descriptor
		wantErr bool
	}{
		{
			name: "simple Request Header",
			args: args{
				matchers: []*v1alpha1.GlobalRateLimit_Action{
					{
						RequestHeaders: &v1alpha1.GlobalRateLimit_Action_RequestHeaders{
							HeaderName:    "foo",
							DescriptorKey: "foo",
						},
					},
				},
				limit: &v1alpha1.GlobalRateLimit_Limit{
					Unit:            "hour",
					RequestsPerUnit: 1,
				},
				shadowMode: false,
			},
			want: []types.RateLimit_Service_Descriptor{
				{
					Key: "foo",
					RateLimit: v1alpha1.GlobalRateLimit_Limit{
						Unit:            "hour",
						RequestsPerUnit: 1,
					},
					ShadowMode: false,
				},
			},
			wantErr: false,
		},
		{
			name: "nested Request Header",
			args: args{
				matchers: []*v1alpha1.GlobalRateLimit_Action{
					{
						RequestHeaders: &v1alpha1.GlobalRateLimit_Action_RequestHeaders{
							HeaderName:    "foo",
							DescriptorKey: "foo",
						},
					},
					{
						RequestHeaders: &v1alpha1.GlobalRateLimit_Action_RequestHeaders{
							HeaderName:    "bar",
							DescriptorKey: "bar",
						},
					},
				},
				limit: &v1alpha1.GlobalRateLimit_Limit{
					Unit:            "hour",
					RequestsPerUnit: 1,
				},
				shadowMode: false,
			},
			want: []types.RateLimit_Service_Descriptor{
				{
					Key: "foo",
					Descriptors: []types.RateLimit_Service_Descriptor{
						{
							Key: "bar",
							RateLimit: v1alpha1.GlobalRateLimit_Limit{
								Unit:            "hour",
								RequestsPerUnit: 1,
							},
							ShadowMode: false,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "simple remote address",
			args: args{
				matchers: []*v1alpha1.GlobalRateLimit_Action{
					{
						RemoteAddress: &v1alpha1.GlobalRateLimit_Action_RemoteAddress{},
					},
				},
				limit: &v1alpha1.GlobalRateLimit_Limit{
					Unit:            "hour",
					RequestsPerUnit: 1,
				},
				shadowMode: false,
			},
			want: []types.RateLimit_Service_Descriptor{
				{
					Key: "remote_address",
					RateLimit: v1alpha1.GlobalRateLimit_Limit{
						Unit:            "hour",
						RequestsPerUnit: 1,
					},
					ShadowMode: false,
				},
			},
			wantErr: false,
		},
		{
			name: "simple Generic key",
			args: args{
				matchers: []*v1alpha1.GlobalRateLimit_Action{
					{
						GenericKey: &v1alpha1.GlobalRateLimit_Action_GenericKey{
							DescriptorValue: "foo",
							DescriptorKey:   &fakeDescriptorKey01,
						},
					},
				},
				limit: &v1alpha1.GlobalRateLimit_Limit{
					Unit:            "hour",
					RequestsPerUnit: 1,
				},
				shadowMode: false,
			},
			want: []types.RateLimit_Service_Descriptor{
				{
					Key:   "foo",
					Value: "foo",
					RateLimit: v1alpha1.GlobalRateLimit_Limit{
						Unit:            "hour",
						RequestsPerUnit: 1,
					},
					ShadowMode: false,
				},
			},
			wantErr: false,
		},
		{
			name: "nested generic key",
			args: args{
				matchers: []*v1alpha1.GlobalRateLimit_Action{
					{
						GenericKey: &v1alpha1.GlobalRateLimit_Action_GenericKey{
							DescriptorValue: "foo",
							DescriptorKey:   &fakeDescriptorKey01,
						},
					},
					{
						GenericKey: &v1alpha1.GlobalRateLimit_Action_GenericKey{
							DescriptorValue: "bar",
							DescriptorKey:   &fakeDescriptorKey02,
						},
					},
				},
				limit: &v1alpha1.GlobalRateLimit_Limit{
					Unit:            "hour",
					RequestsPerUnit: 1,
				},
				shadowMode: false,
			},
			want: []types.RateLimit_Service_Descriptor{
				{
					Key:   "foo",
					Value: "foo",
					Descriptors: []types.RateLimit_Service_Descriptor{
						{
							Key:   "bar",
							Value: "bar",
							RateLimit: v1alpha1.GlobalRateLimit_Limit{
								Unit:            "hour",
								RequestsPerUnit: 1,
							},
							ShadowMode: false,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "simple header value match",
			args: args{
				matchers: []*v1alpha1.GlobalRateLimit_Action{
					{
						HeaderValueMatch: &v1alpha1.GlobalRateLimit_Action_HeaderValueMatch{
							DescriptorValue: "foo",
							ExpectMatch:     &fakeExpectedMatch,
							Headers: []*v1alpha1.GlobalRateLimit_Action_HeaderValueMatch_HeaderMatcher{
								{
									Name:       "bar",
									ExactMatch: "baz",
								},
							},
						},
					},
				},
				limit: &v1alpha1.GlobalRateLimit_Limit{
					Unit:            "hour",
					RequestsPerUnit: 1,
				},
				shadowMode: false,
			},
			want: []types.RateLimit_Service_Descriptor{
				{
					Key:   "header_match",
					Value: "foo",
					RateLimit: v1alpha1.GlobalRateLimit_Limit{
						Unit:            "hour",
						RequestsPerUnit: 1,
					},
					ShadowMode: false,
				},
			},
			wantErr: false,
		},
		{
			name: "simple header value match",
			args: args{
				matchers: []*v1alpha1.GlobalRateLimit_Action{
					{
						HeaderValueMatch: &v1alpha1.GlobalRateLimit_Action_HeaderValueMatch{
							DescriptorValue: "foo",
							ExpectMatch:     &fakeExpectedMatch,
							Headers: []*v1alpha1.GlobalRateLimit_Action_HeaderValueMatch_HeaderMatcher{
								{
									Name:       "bar",
									ExactMatch: "baz",
								},
							},
						},
					},
					{
						HeaderValueMatch: &v1alpha1.GlobalRateLimit_Action_HeaderValueMatch{
							DescriptorValue: "qux",
							ExpectMatch:     &fakeExpectedMatch,
							Headers: []*v1alpha1.GlobalRateLimit_Action_HeaderValueMatch_HeaderMatcher{
								{
									Name:       "quux",
									ExactMatch: "quuz",
								},
							},
						},
					},
				},
				limit: &v1alpha1.GlobalRateLimit_Limit{
					Unit:            "hour",
					RequestsPerUnit: 1,
				},
				shadowMode: false,
			},
			want: []types.RateLimit_Service_Descriptor{
				{
					Key:   "header_match",
					Value: "foo",
					Descriptors: []types.RateLimit_Service_Descriptor{
						{
							Key:   "header_match",
							Value: "qux",
							RateLimit: v1alpha1.GlobalRateLimit_Limit{
								Unit:            "hour",
								RequestsPerUnit: 1,
							},
							ShadowMode: false,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "simple Request Header with shadow mode true",
			args: args{
				matchers: []*v1alpha1.GlobalRateLimit_Action{
					{
						RequestHeaders: &v1alpha1.GlobalRateLimit_Action_RequestHeaders{
							HeaderName:    "foo",
							DescriptorKey: "foo",
						},
					},
				},
				limit: &v1alpha1.GlobalRateLimit_Limit{
					Unit:            "hour",
					RequestsPerUnit: 1,
				},
				shadowMode: true,
			},
			want: []types.RateLimit_Service_Descriptor{
				{
					Key: "foo",
					RateLimit: v1alpha1.GlobalRateLimit_Limit{
						Unit:            "hour",
						RequestsPerUnit: 1,
					},
					ShadowMode: true,
				},
			},
			wantErr: false,
		},
		{
			name: "nested Request Header with shadow mode true",
			args: args{
				matchers: []*v1alpha1.GlobalRateLimit_Action{
					{
						RequestHeaders: &v1alpha1.GlobalRateLimit_Action_RequestHeaders{
							HeaderName:    "foo",
							DescriptorKey: "foo",
						},
					},
					{
						RequestHeaders: &v1alpha1.GlobalRateLimit_Action_RequestHeaders{
							HeaderName:    "bar",
							DescriptorKey: "bar",
						},
					},
				},
				limit: &v1alpha1.GlobalRateLimit_Limit{
					Unit:            "hour",
					RequestsPerUnit: 1,
				},
				shadowMode: true,
			},
			want: []types.RateLimit_Service_Descriptor{
				{
					Key: "foo",
					Descriptors: []types.RateLimit_Service_Descriptor{
						{
							Key: "bar",
							RateLimit: v1alpha1.GlobalRateLimit_Limit{
								Unit:            "hour",
								RequestsPerUnit: 1,
							},
							ShadowMode: true,
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.NewRateLimitDescriptorFromMatcher(tt.args.matchers, tt.args.limit, tt.args.shadowMode)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRateLimitDescriptorFromMatcher() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRateLimitDescriptorFromMatcher() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSyncDescriptors(t *testing.T) {
	type args struct {
		descriptorsData []types.RateLimit_Service_Descriptor
	}
	tests := []struct {
		name string
		args args
		want []types.RateLimit_Service_Descriptor
	}{
		{
			name: "give two descriptor data that doesn't colapse each other",
			args: args{
				descriptorsData: []types.RateLimit_Service_Descriptor{
					{
						Key: "foo",
						RateLimit: v1alpha1.GlobalRateLimit_Limit{
							Unit:            "hour",
							RequestsPerUnit: 1,
						},
					},
					{
						Key: "bar",
						RateLimit: v1alpha1.GlobalRateLimit_Limit{
							Unit:            "hour",
							RequestsPerUnit: 1,
						},
					},
				},
			},
			want: []types.RateLimit_Service_Descriptor{
				{
					Key: "foo",
					RateLimit: v1alpha1.GlobalRateLimit_Limit{
						Unit:            "hour",
						RequestsPerUnit: 1,
					},
				},
				{
					Key: "bar",
					RateLimit: v1alpha1.GlobalRateLimit_Limit{
						Unit:            "hour",
						RequestsPerUnit: 1,
					},
				},
			},
		},
		{
			name: "give two descriptor data that collapse each other",
			args: args{
				descriptorsData: []types.RateLimit_Service_Descriptor{
					{
						Key: "foo",
						Descriptors: []types.RateLimit_Service_Descriptor{
							{
								Key: "bar",
								RateLimit: v1alpha1.GlobalRateLimit_Limit{
									Unit:            "hour",
									RequestsPerUnit: 1,
								},
							},
						},
					},
					{
						Key: "foo",
						Descriptors: []types.RateLimit_Service_Descriptor{
							{
								Key: "baz",
								RateLimit: v1alpha1.GlobalRateLimit_Limit{
									Unit:            "hour",
									RequestsPerUnit: 2,
								},
							},
						},
					},
				},
			},
			want: []types.RateLimit_Service_Descriptor{
				{
					Key: "foo",
					Descriptors: []types.RateLimit_Service_Descriptor{
						{
							Key: "bar",
							RateLimit: v1alpha1.GlobalRateLimit_Limit{
								Unit:            "hour",
								RequestsPerUnit: 1,
							},
						},
						{
							Key: "baz",
							RateLimit: v1alpha1.GlobalRateLimit_Limit{
								Unit:            "hour",
								RequestsPerUnit: 2,
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := service.SyncDescriptors(tt.args.descriptorsData); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SyncDescriptors() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewConfigBuilder(t *testing.T) {
	builder := service.NewConfigBuilder()
	assert.NotNil(t, builder)
	assert.Equal(t, "", builder.Config)
	assert.Equal(t, v1alpha1.RateLimitService{}, builder.RateLimitService)
}

func TestConfigBuilder_SetRateLimitService(t *testing.T) {
	rateLimitService := v1alpha1.RateLimitService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-config",
			Namespace: "test-namespace",
		},
	}

	builder := service.NewConfigBuilder().SetRateLimitService(rateLimitService)

	assert.Equal(t, rateLimitService, builder.RateLimitService)
}

func TestConfigBuilder_SetRateLimitService_Chaining(t *testing.T) {
	rateLimitService := v1alpha1.RateLimitService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-config",
			Namespace: "test-namespace",
		},
	}

	builder := service.NewConfigBuilder()
	returnedBuilder := builder.SetRateLimitService(rateLimitService)

	// Verify method chaining returns the same builder
	assert.Same(t, builder, returnedBuilder)
}

func TestConfigBuilder_SetConfig(t *testing.T) {
	config := "domain: test\ndescriptors:\n  - key: foo"

	builder := service.NewConfigBuilder().SetConfig(config)

	assert.Equal(t, config, builder.Config)
}

func TestConfigBuilder_SetConfig_Chaining(t *testing.T) {
	config := "domain: test"

	builder := service.NewConfigBuilder()
	returnedBuilder := builder.SetConfig(config)

	// Verify method chaining returns the same builder
	assert.Same(t, builder, returnedBuilder)
}

func TestConfigBuilder_BuildLabels(t *testing.T) {
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
				"app.kubernetes.io/name":       "my-ratelimit-config",
				"app.kubernetes.io/managed-by": "istio-rateltimit-operator",
				"app.kubernetes.io/created-by": "my-ratelimit",
			},
		},
		{
			name: "different name",
			rateLimitService: v1alpha1.RateLimitService{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "production-rl",
					Namespace: "production",
				},
			},
			expectedLabels: map[string]string{
				"app.kubernetes.io/name":       "production-rl-config",
				"app.kubernetes.io/managed-by": "istio-rateltimit-operator",
				"app.kubernetes.io/created-by": "production-rl",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := service.NewConfigBuilder().SetRateLimitService(tc.rateLimitService)
			labels := builder.BuildLabels()

			assert.Equal(t, tc.expectedLabels, labels)
		})
	}
}

func TestConfigBuilder_Build(t *testing.T) {
	testCases := []struct {
		name              string
		rateLimitService  v1alpha1.RateLimitService
		config            string
		expectedConfigMap *corev1.ConfigMap
		expectError       bool
	}{
		{
			name: "basic config build",
			rateLimitService: v1alpha1.RateLimitService{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "ratelimit-cfg",
					Namespace: "default",
				},
			},
			config: "domain: test\ndescriptors:\n  - key: foo",
			expectedConfigMap: &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "ratelimit-cfg-config",
					Namespace: "default",
					Labels: map[string]string{
						"app.kubernetes.io/name":       "ratelimit-cfg-config",
						"app.kubernetes.io/managed-by": "istio-rateltimit-operator",
						"app.kubernetes.io/created-by": "ratelimit-cfg",
					},
				},
				Data: map[string]string{
					"config.yaml": "domain: test\ndescriptors:\n  - key: foo",
				},
			},
			expectError: false,
		},
		{
			name: "config in different namespace",
			rateLimitService: v1alpha1.RateLimitService{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "prod-config",
					Namespace: "production",
				},
			},
			config: "domain: production",
			expectedConfigMap: &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "prod-config-config",
					Namespace: "production",
					Labels: map[string]string{
						"app.kubernetes.io/name":       "prod-config-config",
						"app.kubernetes.io/managed-by": "istio-rateltimit-operator",
						"app.kubernetes.io/created-by": "prod-config",
					},
				},
				Data: map[string]string{
					"config.yaml": "domain: production",
				},
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			configMap, err := service.NewConfigBuilder().
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

func TestNewRateLimitConfig(t *testing.T) {
	testCases := []struct {
		name           string
		domain         string
		descriptors    []types.RateLimit_Service_Descriptor
		expectedConfig types.RateLimit_Service_Config
		expectError    bool
	}{
		{
			name:   "basic config",
			domain: "my-domain",
			descriptors: []types.RateLimit_Service_Descriptor{
				{
					Key: "foo",
					RateLimit: v1alpha1.GlobalRateLimit_Limit{
						Unit:            "minute",
						RequestsPerUnit: 10,
					},
				},
			},
			expectedConfig: types.RateLimit_Service_Config{
				Domain: "my-domain",
				Descriptors: []types.RateLimit_Service_Descriptor{
					{
						Key: "foo",
						RateLimit: v1alpha1.GlobalRateLimit_Limit{
							Unit:            "minute",
							RequestsPerUnit: 10,
						},
					},
				},
			},
			expectError: false,
		},
		{
			name:        "empty descriptors",
			domain:      "empty-domain",
			descriptors: []types.RateLimit_Service_Descriptor{},
			expectedConfig: types.RateLimit_Service_Config{
				Domain:      "empty-domain",
				Descriptors: []types.RateLimit_Service_Descriptor{},
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config, err := service.NewRateLimitConfig(tc.domain, tc.descriptors)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedConfig, config)
			}
		})
	}
}

func TestNewRateLimitDescriptor(t *testing.T) {
	testCases := []struct {
		name                string
		globalRateLimitList []v1alpha1.GlobalRateLimit
		expectedDescriptors []types.RateLimit_Service_Descriptor
		expectError         bool
	}{
		{
			name: "single rate limit with request headers",
			globalRateLimitList: []v1alpha1.GlobalRateLimit{
				{
					Spec: v1alpha1.GlobalRateLimitSpec{
						Matcher: []*v1alpha1.GlobalRateLimit_Action{
							{
								RequestHeaders: &v1alpha1.GlobalRateLimit_Action_RequestHeaders{
									HeaderName:    "x-api-key",
									DescriptorKey: "api_key",
								},
							},
						},
						Limit: &v1alpha1.GlobalRateLimit_Limit{
							Unit:            "minute",
							RequestsPerUnit: 100,
						},
					},
				},
			},
			expectedDescriptors: []types.RateLimit_Service_Descriptor{
				{
					Key: "api_key",
					RateLimit: v1alpha1.GlobalRateLimit_Limit{
						Unit:            "minute",
						RequestsPerUnit: 100,
					},
				},
			},
			expectError: false,
		},
		{
			name: "multiple rate limits with request headers",
			globalRateLimitList: []v1alpha1.GlobalRateLimit{
				{
					Spec: v1alpha1.GlobalRateLimitSpec{
						Matcher: []*v1alpha1.GlobalRateLimit_Action{
							{
								RequestHeaders: &v1alpha1.GlobalRateLimit_Action_RequestHeaders{
									HeaderName:    "x-api-key",
									DescriptorKey: "api_key",
								},
							},
						},
						Limit: &v1alpha1.GlobalRateLimit_Limit{
							Unit:            "minute",
							RequestsPerUnit: 100,
						},
					},
				},
				{
					Spec: v1alpha1.GlobalRateLimitSpec{
						Matcher: []*v1alpha1.GlobalRateLimit_Action{
							{
								RequestHeaders: &v1alpha1.GlobalRateLimit_Action_RequestHeaders{
									HeaderName:    "x-user-id",
									DescriptorKey: "user_id",
								},
							},
						},
						Limit: &v1alpha1.GlobalRateLimit_Limit{
							Unit:            "hour",
							RequestsPerUnit: 1000,
						},
					},
				},
			},
			expectedDescriptors: []types.RateLimit_Service_Descriptor{
				{
					Key: "api_key",
					RateLimit: v1alpha1.GlobalRateLimit_Limit{
						Unit:            "minute",
						RequestsPerUnit: 100,
					},
				},
				{
					Key: "user_id",
					RateLimit: v1alpha1.GlobalRateLimit_Limit{
						Unit:            "hour",
						RequestsPerUnit: 1000,
					},
				},
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			descriptors, err := service.NewRateLimitDescriptor(tc.globalRateLimitList)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedDescriptors, descriptors)
			}
		})
	}
}

func TestNewRateLimitDescriptorFromGlobalRateLimit(t *testing.T) {
	testCases := []struct {
		name                string
		globalRateLimit     v1alpha1.GlobalRateLimit
		expectedDescriptors []types.RateLimit_Service_Descriptor
		expectError         bool
	}{
		{
			name: "rate limit with request headers",
			globalRateLimit: v1alpha1.GlobalRateLimit{
				Spec: v1alpha1.GlobalRateLimitSpec{
					Matcher: []*v1alpha1.GlobalRateLimit_Action{
						{
							RequestHeaders: &v1alpha1.GlobalRateLimit_Action_RequestHeaders{
								HeaderName:    "x-user-id",
								DescriptorKey: "user_id",
							},
						},
					},
					Limit: &v1alpha1.GlobalRateLimit_Limit{
						Unit:            "hour",
						RequestsPerUnit: 1000,
					},
				},
			},
			expectedDescriptors: []types.RateLimit_Service_Descriptor{
				{
					Key: "user_id",
					RateLimit: v1alpha1.GlobalRateLimit_Limit{
						Unit:            "hour",
						RequestsPerUnit: 1000,
					},
				},
			},
			expectError: false,
		},
		{
			name: "rate limit with empty matchers",
			globalRateLimit: v1alpha1.GlobalRateLimit{
				Spec: v1alpha1.GlobalRateLimitSpec{
					Matcher: []*v1alpha1.GlobalRateLimit_Action{},
					Limit: &v1alpha1.GlobalRateLimit_Limit{
						Unit:            "hour",
						RequestsPerUnit: 1000,
					},
				},
			},
			expectedDescriptors: nil,
			expectError:         false,
		},
		{
			name: "rate limit with invalid matcher (all nil)",
			globalRateLimit: v1alpha1.GlobalRateLimit{
				Spec: v1alpha1.GlobalRateLimitSpec{
					Matcher: []*v1alpha1.GlobalRateLimit_Action{
						{},
					},
					Limit: &v1alpha1.GlobalRateLimit_Limit{
						Unit:            "hour",
						RequestsPerUnit: 1000,
					},
				},
			},
			expectedDescriptors: nil,
			expectError:         false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			descriptors, err := service.NewRateLimitDescriptorFromGlobalRateLimit(tc.globalRateLimit)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedDescriptors, descriptors)
			}
		})
	}
}

func TestNewRateLimitDescriptorFromMatcher_GenericKeyWithoutDescriptorKey(t *testing.T) {
	matchers := []*v1alpha1.GlobalRateLimit_Action{
		{
			GenericKey: &v1alpha1.GlobalRateLimit_Action_GenericKey{
				DescriptorValue: "my-value",
				DescriptorKey:   nil,
			},
		},
	}
	limit := &v1alpha1.GlobalRateLimit_Limit{
		Unit:            "minute",
		RequestsPerUnit: 50,
	}

	descriptors, err := service.NewRateLimitDescriptorFromMatcher(matchers, limit, false)

	assert.NoError(t, err)
	assert.Len(t, descriptors, 1)
	assert.Equal(t, "generic_key", descriptors[0].Key)
	assert.Equal(t, "my-value", descriptors[0].Value)
}

func TestNewRateLimitDescriptorFromMatcher_NestedRemoteAddress(t *testing.T) {
	matchers := []*v1alpha1.GlobalRateLimit_Action{
		{
			RemoteAddress: &v1alpha1.GlobalRateLimit_Action_RemoteAddress{},
		},
		{
			RequestHeaders: &v1alpha1.GlobalRateLimit_Action_RequestHeaders{
				HeaderName:    "x-api-key",
				DescriptorKey: "api_key",
			},
		},
	}
	limit := &v1alpha1.GlobalRateLimit_Limit{
		Unit:            "second",
		RequestsPerUnit: 10,
	}

	descriptors, err := service.NewRateLimitDescriptorFromMatcher(matchers, limit, false)

	assert.NoError(t, err)
	assert.Len(t, descriptors, 1)
	assert.Equal(t, "remote_address", descriptors[0].Key)
	assert.Len(t, descriptors[0].Descriptors, 1)
	assert.Equal(t, "api_key", descriptors[0].Descriptors[0].Key)
}

func TestSyncDescriptors_EmptySlice(t *testing.T) {
	// Empty slice should not panic, should return nil
	assert.NotPanics(t, func() {
		result := service.SyncDescriptors([]types.RateLimit_Service_Descriptor{})
		assert.Nil(t, result)
	})
}

func TestNewRateLimitDescriptor_Unlimited(t *testing.T) {
	globalRateLimitList := []v1alpha1.GlobalRateLimit{
		{
			Spec: v1alpha1.GlobalRateLimitSpec{
				Matcher: []*v1alpha1.GlobalRateLimit_Action{
					{
						GenericKey: &v1alpha1.GlobalRateLimit_Action_GenericKey{
							DescriptorValue: "internal-service",
						},
					},
				},
				Limit: &v1alpha1.GlobalRateLimit_Limit{
					Unlimited: true,
				},
				DetailedMetric: true,
			},
		},
	}

	descriptors, err := service.NewRateLimitDescriptor(globalRateLimitList)
	assert.NoError(t, err)
	assert.Len(t, descriptors, 1)
	assert.True(t, descriptors[0].RateLimit.Unlimited)
	assert.True(t, descriptors[0].DetailedMetric)
}
