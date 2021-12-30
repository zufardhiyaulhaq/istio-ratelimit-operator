package ratelimit_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/local/ratelimit"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type V3GatewayBuilderTestCase struct {
	name          string
	config        v1alpha1.LocalRateLimitConfig
	ratelimit     v1alpha1.LocalRateLimit
	expectedError bool
}

var V3GatewayBuilderTestGrid = []V3GatewayBuilderTestCase{
	{
		name: "given correct ratelimit",
		config: v1alpha1.LocalRateLimitConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "gateway-config",
				Namespace: "istio-system",
			},
			Spec: v1alpha1.LocalRateLimitConfigSpec{
				Type: "gateway",
				Selector: v1alpha1.LocalRateLimitConfigSelector{
					IstioVersion: []string{"1.9"},
				},
			},
		},
		ratelimit: v1alpha1.LocalRateLimit{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "hello-zufardhiyaulhaq-dev",
				Namespace: "istio-system",
			},
			Spec: v1alpha1.LocalRateLimitSpec{
				Config: "public-gateway-config",
				Selector: v1alpha1.LocalRateLimitSelector{
					VHost: "hello.zufardhiyaulhaq.dev:443",
				},
				Limit: &v1alpha1.LocalRateLimit_Limit{
					Unit:            "hour",
					RequestsPerUnit: 1,
				},
			},
		},
		expectedError: false,
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
			}
		})
	}
}
