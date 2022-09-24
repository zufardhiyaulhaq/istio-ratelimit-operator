package service_test

import (
	"reflect"
	"testing"

	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/service"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/types"
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
