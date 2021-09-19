package service

import (
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
			Labels:    n.buildLabels(),
		},
	}

	if n.RateLimitService.Spec.Backend.Redis != nil {
		data, err := n.buildRedisEnv()
		if err != nil {
			return configMap, err
		}

		configMap.Data = data
	}

	return configMap, nil
}

func (n *EnvBuilder) buildLabels() map[string]string {
	var labels = map[string]string{
		"app.kubernetes.io/name":       n.RateLimitService.Name + "-config-env",
		"app.kubernetes.io/managed-by": "istio-rateltimit-operator",
		"app.kubernetes.io/created-by": n.RateLimitService.Name,
	}

	return labels
}

func (n *EnvBuilder) buildRedisEnv() (map[string]string, error) {
	data := make(map[string]string)

	data["REDIS_SOCKET_TYPE"] = "tcp"
	data["USE_STATSD"] = "false"
	data["REDIS_TYPE"] = n.RateLimitService.Spec.Backend.Redis.Type
	data["REDIS_URL"] = n.RateLimitService.Spec.Backend.Redis.URL

	return data, nil
}
