package service

import (
	"strconv"

	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type EnvBuilder struct {
	RateLimitService v1alpha1.RateLimitService
}

func NewEnvBuilder() *EnvBuilder {
	return &EnvBuilder{}
}

func (n *EnvBuilder) SetRateLimitService(rateLimitService v1alpha1.RateLimitService) *EnvBuilder {
	n.RateLimitService = rateLimitService
	return n
}

func (n *EnvBuilder) Build() (*corev1.ConfigMap, error) {
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      n.RateLimitService.Name + "-config-env",
			Namespace: n.RateLimitService.Namespace,
			Labels:    n.BuildLabels(),
		},
	}

	data, err := n.BuildEnv()
	if err != nil {
		return configMap, nil
	}

	configMap.Data = data
	return configMap, nil
}

func (n *EnvBuilder) BuildLabels() map[string]string {
	var labels = map[string]string{
		"app.kubernetes.io/name":       n.RateLimitService.Name + "-config-env",
		"app.kubernetes.io/managed-by": "istio-rateltimit-operator",
		"app.kubernetes.io/created-by": n.RateLimitService.Name,
	}

	return labels
}

func (n *EnvBuilder) BuildEnv() (map[string]string, error) {
	env := make(map[string]string)

	defaultEnv, err := n.BuildDefaultEnv()
	if err != nil {
		return env, err
	}

	for key, value := range defaultEnv {
		env[key] = value
	}

	if n.RateLimitService.Spec.Environment != nil {
		for key, value := range *n.RateLimitService.Spec.Environment {
			env[key] = value
		}
	}

	if n.RateLimitService.Spec.Backend.Redis != nil {
		redisEnv, err := n.BuildRedisEnv()
		if err != nil {
			return env, err
		}

		for key, value := range redisEnv {
			env[key] = value
		}
	}

	if n.RateLimitService.Spec.Monitoring != nil {
		if n.RateLimitService.Spec.Monitoring.Type != "" {
			monEnv, err := n.BuildMonitoringEnv()
			if err != nil {
				return env, err
			}
			for key, value := range monEnv {
				env[key] = value
			}
		} else if n.RateLimitService.Spec.Monitoring.Enabled {
			// Legacy statsd sidecar path
			statsdEnv, err := n.BuildStatsdEnv()
			if err != nil {
				return env, err
			}
			for key, value := range statsdEnv {
				env[key] = value
			}
		}
	}

	if n.RateLimitService.Spec.ResponseHeaders != nil {
		rhEnv, err := n.BuildResponseHeadersEnv()
		if err != nil {
			return env, err
		}
		for key, value := range rhEnv {
			env[key] = value
		}
	}

	if n.RateLimitService.Spec.Logging != nil {
		logEnv, err := n.BuildLoggingEnv()
		if err != nil {
			return env, err
		}
		for key, value := range logEnv {
			env[key] = value
		}
	}

	if n.RateLimitService.Spec.ShadowMode {
		smEnv, err := n.BuildShadowModeEnv()
		if err != nil {
			return env, err
		}
		for key, value := range smEnv {
			env[key] = value
		}
	}

	if n.RateLimitService.Spec.Server != nil {
		serverEnv, err := n.BuildServerEnv()
		if err != nil {
			return env, err
		}
		for key, value := range serverEnv {
			env[key] = value
		}
	}

	if n.RateLimitService.Spec.Backend != nil {
		backendEnv, err := n.BuildBackendEnv()
		if err != nil {
			return env, err
		}
		for key, value := range backendEnv {
			env[key] = value
		}
	}

	return env, nil
}

func (n *EnvBuilder) BuildDefaultEnv() (map[string]string, error) {
	data := make(map[string]string)

	data["USE_STATSD"] = "false"

	return data, nil
}

func (n *EnvBuilder) BuildRedisEnv() (map[string]string, error) {
	data := make(map[string]string)

	data["REDIS_SOCKET_TYPE"] = "tcp"
	data["REDIS_TYPE"] = n.RateLimitService.Spec.Backend.Redis.Type
	data["REDIS_URL"] = n.RateLimitService.Spec.Backend.Redis.URL

	if n.RateLimitService.Spec.Backend.Redis.Auth != "" {
		data["REDIS_AUTH"] = n.RateLimitService.Spec.Backend.Redis.Auth
	}

	if n.RateLimitService.Spec.Backend.Redis.Config != nil {
		if n.RateLimitService.Spec.Backend.Redis.Config.PipelineLimit != nil {
			data["REDIS_PIPELINE_LIMIT"] = strconv.Itoa(*n.RateLimitService.Spec.Backend.Redis.Config.PipelineLimit)
		}

		if n.RateLimitService.Spec.Backend.Redis.Config.PipelineWindow != nil {
			data["REDIS_PIPELINE_WINDOW"] = *n.RateLimitService.Spec.Backend.Redis.Config.PipelineWindow
		}
	}

	if n.RateLimitService.Spec.Backend.Redis.TLS != nil {
		if n.RateLimitService.Spec.Backend.Redis.TLS.Enabled {
			data["REDIS_TLS"] = "true"
		}

		if n.RateLimitService.Spec.Backend.Redis.TLS.SecretRef != "" {
			data["REDIS_TLS_CACERT"] = "/tls/redis/ca.crt"
			data["REDIS_TLS_CLIENT_CERT"] = "/tls/redis/tls.crt"
			data["REDIS_TLS_CLIENT_KEY"] = "/tls/redis/tls.key"
		}

		if n.RateLimitService.Spec.Backend.Redis.TLS.SkipHostnameVerification {
			data["REDIS_TLS_SKIP_HOSTNAME_VERIFICATION"] = "true"
		}
	}

	if n.RateLimitService.Spec.Backend.Redis.Pool != nil {
		if n.RateLimitService.Spec.Backend.Redis.Pool.Size > 0 {
			data["REDIS_POOL_SIZE"] = strconv.Itoa(n.RateLimitService.Spec.Backend.Redis.Pool.Size)
		}
		if n.RateLimitService.Spec.Backend.Redis.Pool.OnEmptyBehavior != "" {
			data["REDIS_POOL_ON_EMPTY_BEHAVIOR"] = n.RateLimitService.Spec.Backend.Redis.Pool.OnEmptyBehavior
		}
		if n.RateLimitService.Spec.Backend.Redis.Pool.OnEmptyWaitDuration != "" {
			data["REDIS_POOL_ON_EMPTY_WAIT_DURATION"] = n.RateLimitService.Spec.Backend.Redis.Pool.OnEmptyWaitDuration
		}
	}

	if n.RateLimitService.Spec.Backend.Redis.Timeout != nil {
		data["REDIS_TIMEOUT"] = *n.RateLimitService.Spec.Backend.Redis.Timeout
	}

	if n.RateLimitService.Spec.Backend.Redis.HealthCheckActiveConnection {
		data["REDIS_HEALTH_CHECK_ACTIVE_CONNECTION"] = "true"
	}

	if n.RateLimitService.Spec.Backend.Redis.PerSecond != nil {
		ps := n.RateLimitService.Spec.Backend.Redis.PerSecond
		if ps.Enabled {
			data["REDIS_PERSECOND"] = "true"
			data["REDIS_PERSECOND_SOCKET_TYPE"] = "tcp"
			if ps.URL != "" {
				data["REDIS_PERSECOND_URL"] = ps.URL
			}
			if ps.TLS != nil && ps.TLS.Enabled {
				data["REDIS_PERSECOND_TLS"] = "true"
				if ps.TLS.SecretRef != "" {
					data["REDIS_PERSECOND_TLS_CACERT"] = "/tls/redis-persecond/ca.crt"
					data["REDIS_PERSECOND_TLS_CLIENT_CERT"] = "/tls/redis-persecond/tls.crt"
					data["REDIS_PERSECOND_TLS_CLIENT_KEY"] = "/tls/redis-persecond/tls.key"
				}
			}
			if ps.Pool != nil && ps.Pool.Size > 0 {
				data["REDIS_PERSECOND_POOL_SIZE"] = strconv.Itoa(ps.Pool.Size)
			}
		}
	}

	return data, nil
}

func (n *EnvBuilder) BuildBackendEnv() (map[string]string, error) {
	data := make(map[string]string)

	if n.RateLimitService.Spec.Backend == nil {
		return data, nil
	}

	if n.RateLimitService.Spec.Backend.CacheKeyPrefix != "" {
		data["CACHE_KEY_PREFIX"] = n.RateLimitService.Spec.Backend.CacheKeyPrefix
	}

	if n.RateLimitService.Spec.Backend.StopCacheKeyIncrementWhenOverlimit {
		data["STOP_CACHE_KEY_INCREMENT_WHEN_OVERLIMIT"] = "true"
	}

	return data, nil
}

func (n *EnvBuilder) BuildStatsdEnv() (map[string]string, error) {
	data := make(map[string]string)
	data["USE_STATSD"] = "true"
	data["STATSD_HOST"] = "localhost"
	data["STATSD_PORT"] = "9125"

	return data, nil
}

func (n *EnvBuilder) BuildMonitoringEnv() (map[string]string, error) {
	data := make(map[string]string)

	if n.RateLimitService.Spec.Monitoring == nil {
		return data, nil
	}

	mon := n.RateLimitService.Spec.Monitoring

	if mon.Type == "prometheus" {
		data["USE_STATSD"] = "false"
		data["USE_PROMETHEUS"] = "true"
		if mon.Prometheus != nil {
			if mon.Prometheus.Addr != "" {
				data["PROMETHEUS_ADDR"] = mon.Prometheus.Addr
			}
			if mon.Prometheus.Path != "" {
				data["PROMETHEUS_PATH"] = mon.Prometheus.Path
			}
		}
	} else {
		data["USE_STATSD"] = "false"
	}

	if mon.NearLimitRatio != nil {
		data["NEAR_LIMIT_RATIO"] = *mon.NearLimitRatio
	}

	if mon.StatsFlushInterval != nil {
		data["STATS_FLUSH_INTERVAL"] = *mon.StatsFlushInterval
	}

	if mon.Tracing != nil && mon.Tracing.Enabled {
		data["TRACING_ENABLED"] = "true"
		if mon.Tracing.ExporterProtocol != "" {
			data["TRACING_EXPORTER_PROTOCOL"] = mon.Tracing.ExporterProtocol
		}
		if mon.Tracing.ServiceName != "" {
			data["TRACING_SERVICE_NAME"] = mon.Tracing.ServiceName
		}
		if mon.Tracing.ServiceNamespace != "" {
			data["TRACING_SERVICE_NAMESPACE"] = mon.Tracing.ServiceNamespace
		}
		if mon.Tracing.SamplingRate != "" {
			data["TRACING_SAMPLING_RATE"] = mon.Tracing.SamplingRate
		}
	}

	return data, nil
}

func (n *EnvBuilder) BuildResponseHeadersEnv() (map[string]string, error) {
	data := make(map[string]string)

	if n.RateLimitService.Spec.ResponseHeaders != nil && n.RateLimitService.Spec.ResponseHeaders.Enabled {
		data["RESPONSE_HEADERS_ENABLED"] = "true"
	}

	return data, nil
}

func (n *EnvBuilder) BuildLoggingEnv() (map[string]string, error) {
	data := make(map[string]string)

	if n.RateLimitService.Spec.Logging == nil {
		return data, nil
	}

	if n.RateLimitService.Spec.Logging.Level != "" {
		data["LOG_LEVEL"] = n.RateLimitService.Spec.Logging.Level
	}

	if n.RateLimitService.Spec.Logging.Format != "" {
		data["LOG_FORMAT"] = n.RateLimitService.Spec.Logging.Format
	}

	return data, nil
}

func (n *EnvBuilder) BuildShadowModeEnv() (map[string]string, error) {
	data := make(map[string]string)

	if n.RateLimitService.Spec.ShadowMode {
		data["SHADOW_MODE"] = "true"
	}

	return data, nil
}

func (n *EnvBuilder) BuildServerEnv() (map[string]string, error) {
	data := make(map[string]string)

	if n.RateLimitService.Spec.Server == nil {
		return data, nil
	}

	if n.RateLimitService.Spec.Server.GRPC != nil {
		grpc := n.RateLimitService.Spec.Server.GRPC

		if grpc.Port != nil {
			data["GRPC_PORT"] = strconv.Itoa(int(*grpc.Port))
		}

		if grpc.TLS != nil && grpc.TLS.Enabled {
			if grpc.TLS.SecretRef != "" {
				data["GRPC_SERVER_TLS_CERT"] = "/tls/grpc/tls.crt"
				data["GRPC_SERVER_TLS_KEY"] = "/tls/grpc/tls.key"
			}
		}

		if grpc.ClientTLS != nil {
			if grpc.ClientTLS.CACertSecretRef != "" {
				data["GRPC_SERVER_TLS_CLIENT_CACERT"] = "/tls/grpc-client/ca.crt"
			}
			if grpc.ClientTLS.SAN != "" {
				data["GRPC_CLIENT_TLS_SAN"] = grpc.ClientTLS.SAN
			}
		}

		if grpc.MaxConnectionAge != nil {
			data["GRPC_MAX_CONNECTION_AGE"] = *grpc.MaxConnectionAge
		}
		if grpc.MaxConnectionAgeGrace != nil {
			data["GRPC_MAX_CONNECTION_AGE_GRACE"] = *grpc.MaxConnectionAgeGrace
		}
	}

	if n.RateLimitService.Spec.Server.Debug != nil && n.RateLimitService.Spec.Server.Debug.Port != nil {
		data["DEBUG_PORT"] = strconv.Itoa(int(*n.RateLimitService.Spec.Server.Debug.Port))
	}

	return data, nil
}
