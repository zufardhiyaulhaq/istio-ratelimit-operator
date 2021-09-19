package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/global/config"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Config1_9BuilderTestCase struct {
	name          string
	config        v1alpha1.GlobalRateLimitConfig
	expectedError bool
}

var config1_9BuilderTestGrid = []Config1_9BuilderTestCase{
	{
		name: "given correct ratelimit",
		config: v1alpha1.GlobalRateLimitConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "istio-system",
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
		expectedError: false,
	},
}

func TestNewConfig1_9Builder(t *testing.T) {
	for _, test := range config1_9BuilderTestGrid {
		t.Run(test.name, func(t *testing.T) {
			envoyfilter, err := config.NewConfig1_9Builder(test.config).
				Build()

			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.config.Name+"-1.9", envoyfilter.Name)
				assert.Equal(t, test.config.Namespace, envoyfilter.Namespace)
			}
		})
	}
}
