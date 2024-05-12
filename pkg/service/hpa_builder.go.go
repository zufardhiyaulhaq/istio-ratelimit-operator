package service

import (
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"

	autoscalingv2 "k8s.io/api/autoscaling/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var DEFAULT_CPU_AVERAGE_UTILIZATION int32 = 80

type HorizontalPodAutoscalerBuilder struct {
	RateLimitService v1alpha1.RateLimitService
}

func NewHorizontalPodAutoscalerBuilder() *HorizontalPodAutoscalerBuilder {
	return &HorizontalPodAutoscalerBuilder{}
}

func (n *HorizontalPodAutoscalerBuilder) SetRateLimitService(rateLimitService v1alpha1.RateLimitService) *HorizontalPodAutoscalerBuilder {
	n.RateLimitService = rateLimitService
	return n
}
func (n *HorizontalPodAutoscalerBuilder) Build() (*autoscalingv2.HorizontalPodAutoscaler, error) {
	HorizontalPodAutoscaler := &autoscalingv2.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      n.RateLimitService.Name,
			Namespace: n.RateLimitService.Namespace,
			Labels:    n.BuildLabels(),
		},
		Spec: autoscalingv2.HorizontalPodAutoscalerSpec{
			MinReplicas: n.RateLimitService.Spec.Kubernetes.AutoScaling.MinReplica,
			MaxReplicas: *n.RateLimitService.Spec.Kubernetes.AutoScaling.MaxReplica,
			ScaleTargetRef: autoscalingv2.CrossVersionObjectReference{
				APIVersion: "apps/v1",
				Kind:       "Deployment",
				Name:       n.RateLimitService.Name,
			},
			Metrics: []autoscalingv2.MetricSpec{
				{
					Type: autoscalingv2.ResourceMetricSourceType,
					Resource: &autoscalingv2.ResourceMetricSource{
						Name: "cpu",
						Target: autoscalingv2.MetricTarget{
							Type:               autoscalingv2.UtilizationMetricType,
							AverageUtilization: &DEFAULT_CPU_AVERAGE_UTILIZATION,
						},
					},
				},
			},
		},
	}

	return HorizontalPodAutoscaler, nil
}

func (n *HorizontalPodAutoscalerBuilder) BuildLabels() map[string]string {
	var labels = map[string]string{
		"app.kubernetes.io/name":       n.RateLimitService.Name,
		"app.kubernetes.io/managed-by": "istio-rateltimit-operator",
		"app.kubernetes.io/created-by": n.RateLimitService.Name,
	}

	return labels
}
