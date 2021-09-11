package config

import (
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"
	istioClientNetworking "istio.io/client-go/pkg/apis/networking/v1alpha3"
)

type GlobalRateLimitConfigEnvoyFilterBuilder struct {
	Config v1alpha1.GlobalRateLimitConfig
}

func NewGlobalRateLimitConfigEnvoyFilterBuilder() *GlobalRateLimitConfigEnvoyFilterBuilder {
	return &GlobalRateLimitConfigEnvoyFilterBuilder{}
}

func (g *GlobalRateLimitConfigEnvoyFilterBuilder) SetSpec(config v1alpha1.GlobalRateLimitConfig) *GlobalRateLimitConfigEnvoyFilterBuilder {
	g.Config = config
	return g
}

func (g *GlobalRateLimitConfigEnvoyFilterBuilder) Build() ([]*istioClientNetworking.EnvoyFilter, error) {
	var envoyFilters []*istioClientNetworking.EnvoyFilter
	for _, version := range g.Config.Spec.Selector.IstioVersion {
		factory, err := GetGlobalRateLimitConfigEnvoyFilterFactory(version, g.Config)
		if err != nil {
			return nil, err
		}

		envoyfilter, err := factory.Build()
		if err != nil {
			return nil, err
		}

		envoyFilters = append(envoyFilters, envoyfilter)
	}

	return envoyFilters, nil
}
