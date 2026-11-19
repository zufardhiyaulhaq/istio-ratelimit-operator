package service

import (
	"net"
	"strconv"

	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"
)

const defaultPrometheusPort = int32(9090)

// UsesStatsdExporter reports whether the ratelimit pod runs the statsd-exporter
// sidecar. That is the legacy path (monitoring.enabled with no type) and type statsd.
func UsesStatsdExporter(monitoring *v1alpha1.RateLimitServiceSpec_Monitoring) bool {
	if monitoring == nil {
		return false
	}

	switch monitoring.Type {
	case "":
		return monitoring.Enabled
	case "statsd":
		return true
	default:
		return false
	}
}

// PrometheusPort returns the port envoy/ratelimit serves Prometheus metrics on,
// or 0 when monitoring type is not prometheus. It follows PROMETHEUS_ADDR and
// falls back to the upstream default of 9090.
func PrometheusPort(monitoring *v1alpha1.RateLimitServiceSpec_Monitoring) int32 {
	if monitoring == nil || monitoring.Type != "prometheus" {
		return 0
	}

	if monitoring.Prometheus == nil || monitoring.Prometheus.Addr == "" {
		return defaultPrometheusPort
	}

	_, portString, err := net.SplitHostPort(monitoring.Prometheus.Addr)
	if err != nil {
		return defaultPrometheusPort
	}

	port, err := strconv.ParseInt(portString, 10, 32)
	if err != nil || port <= 0 || port > 65535 {
		return defaultPrometheusPort
	}

	return int32(port)
}
