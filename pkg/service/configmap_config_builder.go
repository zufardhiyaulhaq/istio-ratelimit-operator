package service

import (
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ConfigBuilder struct {
	Name      string
	Namespace string
	Spec      v1alpha1.RateLimitServiceSpec
}

func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{}
}

func (n *ConfigBuilder) SetName(name string) *ConfigBuilder {
	n.Name = name
	return n
}

func (n *ConfigBuilder) SetNamespace(namespace string) *ConfigBuilder {
	n.Namespace = namespace
	return n
}

func (n *ConfigBuilder) SetSpec(spec v1alpha1.RateLimitServiceSpec) *ConfigBuilder {
	n.Spec = spec
	return n
}

func (n *ConfigBuilder) Build() (*corev1.ConfigMap, error) {
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      n.Name + "-config",
			Namespace: n.Namespace,
			Labels: map[string]string{
				"app.kubernetes.io/name":       n.Name + "-config",
				"app.kubernetes.io/created-by": "istio-rateltimit-operator",
				"app.kubernetes.io/managed-by": "istio-rateltimit-operator",
			},
		},
		Data: map[string]string{
			"config.yaml": "lol",
		},
	}

	return configMap, nil
}
