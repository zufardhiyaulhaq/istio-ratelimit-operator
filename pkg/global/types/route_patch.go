package types

import (
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"
)

type RoutePatchValues struct {
	Route Route `json:"route,omitempty" yaml:"route,omitempty"`
}

type Route struct {
	Ratelimits []RateLimits `json:"rate_limits,omitempty" yaml:"rate_limits,omitempty"`
}

type RateLimits struct {
	Actions []*v1alpha1.GlobalRateLimit_Action `json:"actions,omitempty" yaml:"actions,omitempty"`
}
