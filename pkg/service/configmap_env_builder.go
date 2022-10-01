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

	data := make(map[string]string)

	defaultEnv, err := n.BuildDefaultEnv()
	if err != nil {
		return configMap, err
	}

	for key, value := range defaultEnv {
		data[key] = value
	}

	if n.RateLimitService.Spec.Kubernetes != nil {
		if n.RateLimitService.Spec.Kubernetes.Environment != nil {
			for key, value := range *n.RateLimitService.Spec.Kubernetes.Environment {
				data[key] = value
			}
		}
	}

	if n.RateLimitService.Spec.Backend.Redis != nil {
		redisEnv, err := n.BuildRedisEnv()
		if err != nil {
			return configMap, err
		}

		for key, value := range redisEnv {
			data[key] = value
		}
	}

	if n.RateLimitService.Spec.Monitoring != nil {
		if n.RateLimitService.Spec.Monitoring.Statsd != nil {
			statsdEnv, err := n.BuildStatsdEnv()
			if err != nil {
				return configMap, err
			}

			for key, value := range statsdEnv {
				data[key] = value
			}
		}
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

	return data, nil
}

func (n *EnvBuilder) BuildStatsdEnv() (map[string]string, error) {
	data := make(map[string]string)

	if n.RateLimitService.Spec.Monitoring.Statsd.Enabled {
		data["USE_STATSD"] = "true"
		data["STATSD_HOST"] = n.RateLimitService.Spec.Monitoring.Statsd.Spec.Host
		data["STATSD_PORT"] = strconv.Itoa(n.RateLimitService.Spec.Monitoring.Statsd.Spec.Port)
	}

	return data, nil
}
