package types

type LocalRateLimitConfig_Value struct {
	Name        string                           `json:"name,omitempty" yaml:"name,omitempty"`
	TypedConfig LocalRateLimitConfig_TypedConfig `json:"typed_config,omitempty" yaml:"typed_config,omitempty"`
}

type LocalRateLimitConfig_TypedConfig struct {
	Type string `json:"@type,omitempty" yaml:"@type,omitempty"`
}
