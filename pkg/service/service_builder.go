package service

import (
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type ServiceBuilder struct {
	RateLimitService v1alpha1.RateLimitService
}

func NewServiceBuilder() *ServiceBuilder {
	return &ServiceBuilder{}
}

func (n *ServiceBuilder) SetRateLimitService(rateLimitService v1alpha1.RateLimitService) *ServiceBuilder {
	n.RateLimitService = rateLimitService
	return n
}

func (n *ServiceBuilder) Build() (*corev1.Service, error) {
	Service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      n.RateLimitService.Name,
			Namespace: n.RateLimitService.Namespace,
			Labels:    n.buildLabels(),
		},
		Spec: corev1.ServiceSpec{
			Selector: n.buildLabels(),
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					Port:       int32(8080),
					TargetPort: intstr.FromInt(8080),
				},
				{
					Name:       "grpc",
					Port:       int32(8081),
					TargetPort: intstr.FromInt(8081),
				},
				{
					Name:       "http-admin",
					Port:       int32(6070),
					TargetPort: intstr.FromInt(6070),
				},
			},
		},
	}

	return Service, nil
}

func (n *ServiceBuilder) buildLabels() map[string]string {
	var labels = map[string]string{
		"app.kubernetes.io/name":       n.RateLimitService.Name,
		"app.kubernetes.io/managed-by": "istio-rateltimit-operator",
		"app.kubernetes.io/created-by": n.RateLimitService.Name,
	}

	return labels
}
