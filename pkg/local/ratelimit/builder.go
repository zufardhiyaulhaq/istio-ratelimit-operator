package ratelimit

import (
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"
	clientnetworking "istio.io/client-go/pkg/apis/networking/v1alpha3"
)

type ConfigBuilder struct {
	RateLimit v1alpha1.LocalRateLimit
	Config    v1alpha1.LocalRateLimitConfig
	Versions  []string
	Labels    map[string]string
}

func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{}
}

func (g *ConfigBuilder) SetRateLimit(ratelimit v1alpha1.LocalRateLimit) *ConfigBuilder {
	g.RateLimit = ratelimit
	return g
}

func (g *ConfigBuilder) SetConfig(config v1alpha1.LocalRateLimitConfig) *ConfigBuilder {
	g.Config = config
	return g
}

func (g *ConfigBuilder) SetVersions(versions []string) *ConfigBuilder {
	g.Versions = versions
	return g
}

func (g *ConfigBuilder) SetLabels(labels map[string]string) *ConfigBuilder {
	g.Labels = labels
	return g
}

func (g *ConfigBuilder) Build() ([]*clientnetworking.EnvoyFilter, error) {
	var envoyFilters []*clientnetworking.EnvoyFilter
	for _, version := range g.Config.Spec.Selector.IstioVersion {
		factory, err := GetConfigFactory(version, g.Config, g.RateLimit)
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
