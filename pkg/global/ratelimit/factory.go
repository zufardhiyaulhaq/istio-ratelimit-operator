package ratelimit

import (
	"fmt"

	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"
	istioClientNetworking "istio.io/client-go/pkg/apis/networking/v1alpha3"
)

type globalRateLimitEnvoyFilterFactory interface {
	Build() (*istioClientNetworking.EnvoyFilter, error)
}

func GetConfigFactory(version string, config v1alpha1.GlobalRateLimitConfig, ratelimit v1alpha1.GlobalRateLimit) (globalRateLimitEnvoyFilterFactory, error) {
	if version == "1.8" {
		return NewV3Builder(config, ratelimit, version), nil
	}
	if version == "1.9" {
		return NewV3Builder(config, ratelimit, version), nil
	}
	if version == "1.10" {
		return NewV3Builder(config, ratelimit, version), nil
	}

	return nil, fmt.Errorf("version not supported")
}
