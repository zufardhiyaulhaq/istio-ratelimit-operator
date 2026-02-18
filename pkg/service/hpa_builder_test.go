package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/service"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNewHorizontalPodAutoscalerBuilder(t *testing.T) {
	builder := service.NewHorizontalPodAutoscalerBuilder()
	assert.NotNil(t, builder)
	assert.Equal(t, v1alpha1.RateLimitService{}, builder.RateLimitService)
}

func TestHorizontalPodAutoscalerBuilder_SetRateLimitService(t *testing.T) {
	rateLimitService := v1alpha1.RateLimitService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-hpa",
			Namespace: "test-namespace",
		},
	}

	builder := service.NewHorizontalPodAutoscalerBuilder().SetRateLimitService(rateLimitService)

	assert.Equal(t, rateLimitService, builder.RateLimitService)
}

func TestHorizontalPodAutoscalerBuilder_SetRateLimitService_Chaining(t *testing.T) {
	rateLimitService := v1alpha1.RateLimitService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-hpa",
			Namespace: "test-namespace",
		},
	}

	builder := service.NewHorizontalPodAutoscalerBuilder()
	returnedBuilder := builder.SetRateLimitService(rateLimitService)

	// Verify method chaining returns the same builder
	assert.Same(t, builder, returnedBuilder)
}

func TestHorizontalPodAutoscalerBuilder_BuildLabels(t *testing.T) {
	testCases := []struct {
		name             string
		rateLimitService v1alpha1.RateLimitService
		expectedLabels   map[string]string
	}{
		{
			name: "basic labels",
			rateLimitService: v1alpha1.RateLimitService{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-ratelimit",
					Namespace: "default",
				},
			},
			expectedLabels: map[string]string{
				"app.kubernetes.io/name":       "my-ratelimit",
				"app.kubernetes.io/managed-by": "istio-rateltimit-operator",
				"app.kubernetes.io/created-by": "my-ratelimit",
			},
		},
		{
			name: "different name",
			rateLimitService: v1alpha1.RateLimitService{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "production-hpa",
					Namespace: "production",
				},
			},
			expectedLabels: map[string]string{
				"app.kubernetes.io/name":       "production-hpa",
				"app.kubernetes.io/managed-by": "istio-rateltimit-operator",
				"app.kubernetes.io/created-by": "production-hpa",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := service.NewHorizontalPodAutoscalerBuilder().SetRateLimitService(tc.rateLimitService)
			labels := builder.BuildLabels()

			assert.Equal(t, tc.expectedLabels, labels)
		})
	}
}

func TestHorizontalPodAutoscalerBuilder_Build(t *testing.T) {
	minReplica := int32(2)
	maxReplica := int32(10)
	defaultCPUUtilization := int32(80)

	testCases := []struct {
		name             string
		rateLimitService v1alpha1.RateLimitService
		expectedHPA      *autoscalingv2.HorizontalPodAutoscaler
		expectError      bool
	}{
		{
			name: "basic HPA build",
			rateLimitService: v1alpha1.RateLimitService{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "ratelimit-hpa",
					Namespace: "default",
				},
				Spec: v1alpha1.RateLimitServiceSpec{
					Kubernetes: &v1alpha1.RateLimitServiceSpec_Kubernetes{
						AutoScaling: &v1alpha1.RateLimitServiceSpec_Kubernetes_AutoScaling{
							MinReplica: &minReplica,
							MaxReplica: &maxReplica,
						},
					},
				},
			},
			expectedHPA: &autoscalingv2.HorizontalPodAutoscaler{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "ratelimit-hpa",
					Namespace: "default",
					Labels: map[string]string{
						"app.kubernetes.io/name":       "ratelimit-hpa",
						"app.kubernetes.io/managed-by": "istio-rateltimit-operator",
						"app.kubernetes.io/created-by": "ratelimit-hpa",
					},
				},
				Spec: autoscalingv2.HorizontalPodAutoscalerSpec{
					MinReplicas: &minReplica,
					MaxReplicas: maxReplica,
					ScaleTargetRef: autoscalingv2.CrossVersionObjectReference{
						APIVersion: "apps/v1",
						Kind:       "Deployment",
						Name:       "ratelimit-hpa",
					},
					Metrics: []autoscalingv2.MetricSpec{
						{
							Type: autoscalingv2.ResourceMetricSourceType,
							Resource: &autoscalingv2.ResourceMetricSource{
								Name: "cpu",
								Target: autoscalingv2.MetricTarget{
									Type:               autoscalingv2.UtilizationMetricType,
									AverageUtilization: &defaultCPUUtilization,
								},
							},
						},
					},
				},
			},
			expectError: false,
		},
		{
			name: "HPA in production namespace",
			rateLimitService: v1alpha1.RateLimitService{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "prod-ratelimit",
					Namespace: "production",
				},
				Spec: v1alpha1.RateLimitServiceSpec{
					Kubernetes: &v1alpha1.RateLimitServiceSpec_Kubernetes{
						AutoScaling: &v1alpha1.RateLimitServiceSpec_Kubernetes_AutoScaling{
							MinReplica: &minReplica,
							MaxReplica: &maxReplica,
						},
					},
				},
			},
			expectedHPA: &autoscalingv2.HorizontalPodAutoscaler{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "prod-ratelimit",
					Namespace: "production",
					Labels: map[string]string{
						"app.kubernetes.io/name":       "prod-ratelimit",
						"app.kubernetes.io/managed-by": "istio-rateltimit-operator",
						"app.kubernetes.io/created-by": "prod-ratelimit",
					},
				},
				Spec: autoscalingv2.HorizontalPodAutoscalerSpec{
					MinReplicas: &minReplica,
					MaxReplicas: maxReplica,
					ScaleTargetRef: autoscalingv2.CrossVersionObjectReference{
						APIVersion: "apps/v1",
						Kind:       "Deployment",
						Name:       "prod-ratelimit",
					},
					Metrics: []autoscalingv2.MetricSpec{
						{
							Type: autoscalingv2.ResourceMetricSourceType,
							Resource: &autoscalingv2.ResourceMetricSource{
								Name: "cpu",
								Target: autoscalingv2.MetricTarget{
									Type:               autoscalingv2.UtilizationMetricType,
									AverageUtilization: &defaultCPUUtilization,
								},
							},
						},
					},
				},
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hpa, err := service.NewHorizontalPodAutoscalerBuilder().
				SetRateLimitService(tc.rateLimitService).
				Build()

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedHPA, hpa)
			}
		})
	}
}

func TestHorizontalPodAutoscalerBuilder_Build_ScaleTargetRef(t *testing.T) {
	minReplica := int32(1)
	maxReplica := int32(5)

	rateLimitService := v1alpha1.RateLimitService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "scale-target-test",
			Namespace: "default",
		},
		Spec: v1alpha1.RateLimitServiceSpec{
			Kubernetes: &v1alpha1.RateLimitServiceSpec_Kubernetes{
				AutoScaling: &v1alpha1.RateLimitServiceSpec_Kubernetes_AutoScaling{
					MinReplica: &minReplica,
					MaxReplica: &maxReplica,
				},
			},
		},
	}

	hpa, err := service.NewHorizontalPodAutoscalerBuilder().
		SetRateLimitService(rateLimitService).
		Build()

	assert.NoError(t, err)
	assert.NotNil(t, hpa)

	// Verify ScaleTargetRef
	assert.Equal(t, "apps/v1", hpa.Spec.ScaleTargetRef.APIVersion)
	assert.Equal(t, "Deployment", hpa.Spec.ScaleTargetRef.Kind)
	assert.Equal(t, "scale-target-test", hpa.Spec.ScaleTargetRef.Name)
}

func TestHorizontalPodAutoscalerBuilder_Build_Metrics(t *testing.T) {
	minReplica := int32(1)
	maxReplica := int32(5)
	expectedCPUUtilization := int32(80)

	rateLimitService := v1alpha1.RateLimitService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "metrics-test",
			Namespace: "default",
		},
		Spec: v1alpha1.RateLimitServiceSpec{
			Kubernetes: &v1alpha1.RateLimitServiceSpec_Kubernetes{
				AutoScaling: &v1alpha1.RateLimitServiceSpec_Kubernetes_AutoScaling{
					MinReplica: &minReplica,
					MaxReplica: &maxReplica,
				},
			},
		},
	}

	hpa, err := service.NewHorizontalPodAutoscalerBuilder().
		SetRateLimitService(rateLimitService).
		Build()

	assert.NoError(t, err)
	assert.NotNil(t, hpa)

	// Verify metrics configuration
	assert.Len(t, hpa.Spec.Metrics, 1)

	metric := hpa.Spec.Metrics[0]
	assert.Equal(t, autoscalingv2.ResourceMetricSourceType, metric.Type)
	assert.NotNil(t, metric.Resource)
	assert.Equal(t, "cpu", string(metric.Resource.Name))
	assert.Equal(t, autoscalingv2.UtilizationMetricType, metric.Resource.Target.Type)
	assert.NotNil(t, metric.Resource.Target.AverageUtilization)
	assert.Equal(t, expectedCPUUtilization, *metric.Resource.Target.AverageUtilization)
}

func TestHorizontalPodAutoscalerBuilder_Build_Replicas(t *testing.T) {
	testCases := []struct {
		name        string
		minReplica  int32
		maxReplica  int32
		expectedMin int32
		expectedMax int32
	}{
		{
			name:        "small scale",
			minReplica:  1,
			maxReplica:  3,
			expectedMin: 1,
			expectedMax: 3,
		},
		{
			name:        "medium scale",
			minReplica:  2,
			maxReplica:  10,
			expectedMin: 2,
			expectedMax: 10,
		},
		{
			name:        "large scale",
			minReplica:  5,
			maxReplica:  100,
			expectedMin: 5,
			expectedMax: 100,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			minReplica := tc.minReplica
			maxReplica := tc.maxReplica

			rateLimitService := v1alpha1.RateLimitService{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "replica-test",
					Namespace: "default",
				},
				Spec: v1alpha1.RateLimitServiceSpec{
					Kubernetes: &v1alpha1.RateLimitServiceSpec_Kubernetes{
						AutoScaling: &v1alpha1.RateLimitServiceSpec_Kubernetes_AutoScaling{
							MinReplica: &minReplica,
							MaxReplica: &maxReplica,
						},
					},
				},
			}

			hpa, err := service.NewHorizontalPodAutoscalerBuilder().
				SetRateLimitService(rateLimitService).
				Build()

			assert.NoError(t, err)
			assert.NotNil(t, hpa)
			assert.NotNil(t, hpa.Spec.MinReplicas)
			assert.Equal(t, tc.expectedMin, *hpa.Spec.MinReplicas)
			assert.Equal(t, tc.expectedMax, hpa.Spec.MaxReplicas)
		})
	}
}

func TestHorizontalPodAutoscalerBuilder_Build_Metadata(t *testing.T) {
	minReplica := int32(1)
	maxReplica := int32(5)

	rateLimitService := v1alpha1.RateLimitService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "metadata-test",
			Namespace: "test-ns",
		},
		Spec: v1alpha1.RateLimitServiceSpec{
			Kubernetes: &v1alpha1.RateLimitServiceSpec_Kubernetes{
				AutoScaling: &v1alpha1.RateLimitServiceSpec_Kubernetes_AutoScaling{
					MinReplica: &minReplica,
					MaxReplica: &maxReplica,
				},
			},
		},
	}

	hpa, err := service.NewHorizontalPodAutoscalerBuilder().
		SetRateLimitService(rateLimitService).
		Build()

	assert.NoError(t, err)
	assert.NotNil(t, hpa)

	// Verify metadata
	assert.Equal(t, "metadata-test", hpa.ObjectMeta.Name)
	assert.Equal(t, "test-ns", hpa.ObjectMeta.Namespace)

	// Verify labels are set correctly
	assert.Equal(t, "metadata-test", hpa.ObjectMeta.Labels["app.kubernetes.io/name"])
	assert.Equal(t, "istio-rateltimit-operator", hpa.ObjectMeta.Labels["app.kubernetes.io/managed-by"])
	assert.Equal(t, "metadata-test", hpa.ObjectMeta.Labels["app.kubernetes.io/created-by"])
}

func TestDefaultCPUAverageUtilization(t *testing.T) {
	// Verify the default CPU average utilization value
	assert.Equal(t, int32(80), service.DEFAULT_CPU_AVERAGE_UTILIZATION)
}
