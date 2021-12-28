package config

import (
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"
	clientnetworking "istio.io/client-go/pkg/apis/networking/v1alpha3"
)

type ConfigBuilder struct {
	Config v1alpha1.LocalRateLimitConfig
}

func (g *ConfigBuilder) SetConfig(config v1alpha1.LocalRateLimitConfig) *ConfigBuilder {
	g.Config = config
	return g
}

func (g *ConfigBuilder) Build() ([]*clientnetworking.EnvoyFilter, error) {
	var envoyFilters []*clientnetworking.EnvoyFilter
	for _, version := range g.Config.Spec.Selector.IstioVersion {
		factory, err := NewConfigFactory(version, g.Config)
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

func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{}
}
