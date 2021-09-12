package types

type HttpFilterPatchValues struct {
	Config      Config      `json:"config,omitempty" yaml:"config,omitempty"`
	TypedConfig TypedConfig `json:"typed_config,omitempty" yaml:"typed_config,omitempty"`
	Name        string      `json:"name,omitempty" yaml:"name,omitempty"`
}

type Config struct {
	Domain           string           `json:"domain,omitempty" yaml:"domain,omitempty"`
	FailureModeDeny  *bool            `json:"failure_mode_deny,omitempty" yaml:"failure_mode_deny,omitempty"`
	Timeout          string           `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	RateLimitService RateLimitService `json:"rate_limit_service,omitempty" yaml:"rate_limit_service,omitempty"`
}

type TypedConfig struct {
	Type             string           `json:"@type,omitempty" yaml:"@type,omitempty"`
	Domain           string           `json:"domain,omitempty" yaml:"domain,omitempty"`
	FailureModeDeny  *bool            `json:"failure_mode_deny,omitempty" yaml:"failure_mode_deny,omitempty"`
	Timeout          string           `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	RateLimitService RateLimitService `json:"rate_limit_service,omitempty" yaml:"rate_limit_service,omitempty"`
}

type RateLimitService struct {
	GRPCService         GRPCService `json:"grpc_service,omitempty" yaml:"grpc_service,omitempty"`
	TransportAPIVersion string      `json:"transport_api_version,omitempty" yaml:"transport_api_version,omitempty"`
}

type GRPCService struct {
	EnvoyGRPC EnvoyGRPC `json:"envoy_grpc,omitempty" yaml:"envoy_grpc,omitempty"`
	Timeout   string    `json:"timeout" yaml:"timeout,omitempty"`
}

type EnvoyGRPC struct {
	ClusterName string `json:"cluster_name,omitempty" yaml:"cluster_name,omitempty"`
}
