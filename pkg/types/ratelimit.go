package types

import (
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"
	"gopkg.in/yaml.v2"
)

type RateLimit_Service_Config struct {
	Domain      string                         `json:"domain,omitempty" yaml:"domain,omitempty"`
	Descriptors []RateLimit_Service_Descriptor `json:"descriptors,omitempty" yaml:"descriptors,omitempty"`
}

type RateLimit_Service_Descriptor struct {
	Key         string                         `json:"key,omitempty" yaml:"key,omitempty"`
	Value       string                         `json:"value,omitempty" yaml:"value,omitempty"`
	ShadowMode  bool                           `json:"shadow_mode,omitempty" yaml:"shadow_mode,omitempty"`
	RateLimit   v1alpha1.GlobalRateLimit_Limit `json:"rate_limit,omitempty" yaml:"rate_limit,omitempty"`
	Descriptors []RateLimit_Service_Descriptor `json:"descriptors,omitempty" yaml:"descriptors,omitempty"`
}

func (n *RateLimit_Service_Config) String() (string, error) {
	bytes, err := yaml.Marshal(n)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
