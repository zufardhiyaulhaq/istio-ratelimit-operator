package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/global/config"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type GlobalRateLimitConfig1_9EnvoyFilterBuilderTestCase struct {
	name          string
	config        v1alpha1.GlobalRateLimitConfig
	expectedError bool
}

var testGrid = []GlobalRateLimitConfig1_9EnvoyFilterBuilderTestCase{
	{
		name: "given correct ratelimit",
		config: v1alpha1.GlobalRateLimitConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "istio-system",
				Labels: map[string]string{
					"generator": "istio-rateltimit-operator",
				},
			},
			Spec: v1alpha1.GlobalRateLimitConfigSpec{
				Type: "gateway",
				Selector: v1alpha1.GlobalRateLimitConfigSelector{
					IstioVersion: []string{"1.9"},
					Labels: map[string]string{
						"app": "istio-public-gateway",
					},
				},
				Ratelimit: v1alpha1.GlobalRateLimitConfigRatelimit{
					Spec: v1alpha1.GlobalRateLimitConfigRatelimitSpec{
						Domain:          "foo",
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
		expectedError: true,
	},
}

func TestConfigurationBuilder(t *testing.T) {
	for _, test := range testGrid {
		t.Run(test.name, func(t *testing.T) {
			_, err := config.NewGlobalRateLimitConfig1_9EnvoyFilterBuilder(test.config).
				Build()

			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
