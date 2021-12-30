package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/local/config"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type V3GatewayBuilderTestCase struct {
	name          string
	config        v1alpha1.LocalRateLimitConfig
	expectedError bool
}

var V3GatewayBuilderTestGrid = []V3GatewayBuilderTestCase{
	{
		name: "given correct ratelimit",
		config: v1alpha1.LocalRateLimitConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "foo",
			},
			Spec: v1alpha1.LocalRateLimitConfigSpec{
				Type: "gateway",
				Selector: v1alpha1.LocalRateLimitConfigSelector{
					IstioVersion: []string{"1.8"},
					Labels: map[string]string{
						"app": "foo",
					},
				},
			},
		},
		expectedError: false,
	},
}

func TestNewV3GatewayBuilder(t *testing.T) {
	for _, test := range V3GatewayBuilderTestGrid {
		t.Run(test.name, func(t *testing.T) {
			envoyfilter, err := config.NewV3GatewayBuilder(test.config, "1.8").
				Build()

			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.config.Name+"-"+"1.8", envoyfilter.Name)
				assert.Equal(t, test.config.Namespace, envoyfilter.Namespace)
			}
		})
	}
}
