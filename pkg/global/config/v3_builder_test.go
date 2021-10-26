package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/global/config"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type V3BuilderTestCase struct {
	name          string
	config        v1alpha1.GlobalRateLimitConfig
	expectedError bool
}

var mockIstioVersion = "1.8"
var v3BuilderTestGrid = []V3BuilderTestCase{
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
					IstioVersion: []string{mockIstioVersion},
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

func TestNewV3Builder(t *testing.T) {
	for _, test := range v3BuilderTestGrid {
		t.Run(test.name, func(t *testing.T) {
			envoyfilter, err := config.NewV3Builder(test.config, mockIstioVersion).
				Build()

			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.config.Name+"-"+mockIstioVersion, envoyfilter.Name)
				assert.Equal(t, test.config.Namespace, envoyfilter.Namespace)
			}
		})
	}
}
