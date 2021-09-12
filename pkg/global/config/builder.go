package config

import (
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"
	istioClientNetworking "istio.io/client-go/pkg/apis/networking/v1alpha3"
)

type ConfigBuilder struct {
	Config v1alpha1.GlobalRateLimitConfig
}

func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{}
}

func (g *ConfigBuilder) SetSpec(config v1alpha1.GlobalRateLimitConfig) *ConfigBuilder {
	g.Config = config
	return g
}

func (g *ConfigBuilder) Build() ([]*istioClientNetworking.EnvoyFilter, error) {
	var envoyFilters []*istioClientNetworking.EnvoyFilter
	for _, version := range g.Config.Spec.Selector.IstioVersion {
		factory, err := GetConfigFactory(version, g.Config)
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
