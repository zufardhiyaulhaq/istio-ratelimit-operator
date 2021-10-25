package ratelimit_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/global/ratelimit"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type V3BuilderTestCase struct {
	name          string
	config        v1alpha1.GlobalRateLimitConfig
	ratelimit     v1alpha1.GlobalRateLimit
	expectedError bool
}

var mockIstioVersion = "1.9"
var v3BuilderTestGrid = []V3BuilderTestCase{
	{
		name: "given correct ratelimit",
		config: v1alpha1.GlobalRateLimitConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "public-gateway-config",
				Namespace: "istio-system",
			},
			Spec: v1alpha1.GlobalRateLimitConfigSpec{
				Type: "gateway",
				Selector: v1alpha1.GlobalRateLimitConfigSelector{
					IstioVersion: []string{mockIstioVersion},
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
	},
}

func TestNewV3Builder(t *testing.T) {
	for _, test := range v3BuilderTestGrid {
		t.Run(test.name, func(t *testing.T) {
			envoyfilter, err := ratelimit.NewV3Builder(test.config, test.ratelimit, mockIstioVersion).
				Build()

			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.ratelimit.Name+"-"+mockIstioVersion, envoyfilter.Name)
				assert.Equal(t, test.ratelimit.Namespace, envoyfilter.Namespace)
			}
		})
	}
}
