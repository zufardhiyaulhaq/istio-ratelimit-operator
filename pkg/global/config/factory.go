package config

import (
	"fmt"

	"github.com/Masterminds/semver"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"
	clientnetworking "istio.io/client-go/pkg/apis/networking/v1alpha3"
)

type ConfigFactory interface {
	Build() (*clientnetworking.EnvoyFilter, error)
}

func NewConfigFactory(v string, config v1alpha1.GlobalRateLimitConfig) (ConfigFactory, error) {
	version, err := semver.NewVersion(v)
	if err != nil {
		return nil, fmt.Errorf("cannot parse version")
	}

	versionConstrain, err := semver.NewConstraint(">= 1.7.x")
	if err != nil {
		return nil, fmt.Errorf("cannot parse version constrain")
	}

	valid, _ := versionConstrain.Validate(version)
	if valid {
		if config.Spec.Type == v1alpha1.Gateway {
			return NewV3GatewayBuilder(config, v), nil
		}
	}

	return nil, fmt.Errorf("version not supported")
}
