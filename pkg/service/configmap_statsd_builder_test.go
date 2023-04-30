package service_test

import (
	"testing"

	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/service"
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
