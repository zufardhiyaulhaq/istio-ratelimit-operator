package types

type LocalRateLimit_HTTPFilter struct {
	Name        string                     `json:"name,omitempty" yaml:"name,omitempty"`
	TypedConfig LocalRateLimit_TypedConfig `json:"typed_config,omitempty" yaml:"typed_config,omitempty"`
}

type LocalRateLimit_HTTPRoute struct {
	TypedPerFilterConfig LocalRateLimit_TypedPerFilterConfig `json:"typed_per_filter_config,omitempty" yaml:"typed_per_filter_config,omitempty"`
}

type LocalRateLimit_TypedPerFilterConfig struct {
	TypedConfig LocalRateLimit_TypedConfig `json:"envoy.filters.http.local_ratelimit,omitempty" yaml:"envoy.filters.http.local_ratelimit,omitempty"`
}

type LocalRateLimit_TypedConfig struct {
	Type    string         `json:"@type,omitempty" yaml:"@type,omitempty"`
	TypeURL string         `json:"type_url,omitempty" yaml:"type_url,omitempty"`
	Value   LocalRateLimit `json:"value,omitempty" yaml:"value,omitempty"`
}
