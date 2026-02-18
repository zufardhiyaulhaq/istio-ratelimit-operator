package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/service"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/settings"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	intstr "k8s.io/apimachinery/pkg/util/intstr"
)

func TestDeploymentBuilder(t *testing.T) {
	setting := settings.Settings{
		RateLimitServiceImage: "foo",
		StatsdExporterImage:   "bar",
	}

	testCases := []struct {
		rateLimitService   v1alpha1.RateLimitService
		expectedDeployment *appsv1.Deployment
	}{
		{
			rateLimitService: v1alpha1.RateLimitService{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "baz",
					Namespace: "fox",
				},
				Spec: v1alpha1.RateLimitServiceSpec{
					Kubernetes: &v1alpha1.RateLimitServiceSpec_Kubernetes{},
				},
			},
			expectedDeployment: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "baz",
					Namespace: "fox",
					Labels: map[string]string{
						"app.kubernetes.io/name":       "baz",
						"app.kubernetes.io/managed-by": "istio-rateltimit-operator",
						"app.kubernetes.io/created-by": "baz",
					},
				},
				Spec: appsv1.DeploymentSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"app.kubernetes.io/name":       "baz",
							"app.kubernetes.io/managed-by": "istio-rateltimit-operator",
							"app.kubernetes.io/created-by": "baz",
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"app.kubernetes.io/name":       "baz",
								"app.kubernetes.io/managed-by": "istio-rateltimit-operator",
								"app.kubernetes.io/created-by": "baz",
							},
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:    "baz",
									Image:   "foo",
									Command: []string{"/bin/ratelimit"},
									Ports: []corev1.ContainerPort{
										{
											Name:          "http",
											ContainerPort: int32(8080),
											Protocol:      corev1.ProtocolTCP,
										},
										{
											Name:          "grpc",
											ContainerPort: int32(8081),
											Protocol:      corev1.ProtocolTCP,
										},
										{
											Name:          "http-admin",
											ContainerPort: int32(6070),
											Protocol:      corev1.ProtocolTCP,
										},
									},
									Env: []corev1.EnvVar{
										{
											Name:  "RUNTIME_ROOT",
											Value: "/data/ratelimit/",
										},
										{
											Name:  "RUNTIME_WATCH_ROOT",
											Value: "false",
										},
										{
											Name:  "RUNTIME_IGNOREDOTFILES",
											Value: "true",
										},
									},
									EnvFrom: []corev1.EnvFromSource{
										{
											ConfigMapRef: &corev1.ConfigMapEnvSource{
												LocalObjectReference: corev1.LocalObjectReference{
													Name: "baz" + "-config-env",
												},
											},
										},
									},
									ReadinessProbe: &corev1.Probe{
										ProbeHandler: corev1.ProbeHandler{
											HTTPGet: &corev1.HTTPGetAction{
												Path: "/healthcheck",
												Port: intstr.FromInt(8080),
											},
										},
										InitialDelaySeconds: int32(5),
										PeriodSeconds:       int32(10),
									},
									VolumeMounts: []corev1.VolumeMount{
										{
											Name:      "baz" + "-config",
											MountPath: "/data/ratelimit/config/",
										},
									},
								},
							},
							Volumes: []corev1.Volume{
								{
									Name: "baz" + "-config",
									VolumeSource: corev1.VolumeSource{
										ConfigMap: &corev1.ConfigMapVolumeSource{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: "baz" + "-config",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			rateLimitService: v1alpha1.RateLimitService{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "baz",
					Namespace: "fox",
				},
				Spec: v1alpha1.RateLimitServiceSpec{
					Kubernetes: &v1alpha1.RateLimitServiceSpec_Kubernetes{
						ExtraLabels: &map[string]string{
							"team": "daf",
						},
					},
				},
			},
			expectedDeployment: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "baz",
					Namespace: "fox",
					Labels: map[string]string{
						"app.kubernetes.io/name":       "baz",
						"app.kubernetes.io/managed-by": "istio-rateltimit-operator",
						"app.kubernetes.io/created-by": "baz",
						"team":                         "daf",
					},
				},
				Spec: appsv1.DeploymentSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"app.kubernetes.io/name":       "baz",
							"app.kubernetes.io/managed-by": "istio-rateltimit-operator",
							"app.kubernetes.io/created-by": "baz",
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"app.kubernetes.io/name":       "baz",
								"app.kubernetes.io/managed-by": "istio-rateltimit-operator",
								"app.kubernetes.io/created-by": "baz",
								"team":                         "daf",
							},
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:    "baz",
									Image:   "foo",
									Command: []string{"/bin/ratelimit"},
									Ports: []corev1.ContainerPort{
										{
											Name:          "http",
											ContainerPort: int32(8080),
											Protocol:      corev1.ProtocolTCP,
										},
										{
											Name:          "grpc",
											ContainerPort: int32(8081),
											Protocol:      corev1.ProtocolTCP,
										},
										{
											Name:          "http-admin",
											ContainerPort: int32(6070),
											Protocol:      corev1.ProtocolTCP,
										},
									},
									Env: []corev1.EnvVar{
										{
											Name:  "RUNTIME_ROOT",
											Value: "/data/ratelimit/",
										},
										{
											Name:  "RUNTIME_WATCH_ROOT",
											Value: "false",
										},
										{
											Name:  "RUNTIME_IGNOREDOTFILES",
											Value: "true",
										},
									},
									EnvFrom: []corev1.EnvFromSource{
										{
											ConfigMapRef: &corev1.ConfigMapEnvSource{
												LocalObjectReference: corev1.LocalObjectReference{
													Name: "baz" + "-config-env",
												},
											},
										},
									},
									ReadinessProbe: &corev1.Probe{
										ProbeHandler: corev1.ProbeHandler{
											HTTPGet: &corev1.HTTPGetAction{
												Path: "/healthcheck",
												Port: intstr.FromInt(8080),
											},
										},
										InitialDelaySeconds: int32(5),
										PeriodSeconds:       int32(10),
									},
									VolumeMounts: []corev1.VolumeMount{
										{
											Name:      "baz" + "-config",
											MountPath: "/data/ratelimit/config/",
										},
									},
								},
							},
							Volumes: []corev1.Volume{
								{
									Name: "baz" + "-config",
									VolumeSource: corev1.VolumeSource{
										ConfigMap: &corev1.ConfigMapVolumeSource{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: "baz" + "-config",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		existingDeployment, err := service.NewDeploymentBuilder(setting).
			SetRateLimitService(tc.rateLimitService).Build()

		assert.NoError(t, err)
		assert.Equal(t, tc.expectedDeployment, existingDeployment)
	}
}

func TestDeploymentBuilder_WithMonitoring(t *testing.T) {
	setting := settings.Settings{
		RateLimitServiceImage: "ratelimit:latest",
		StatsdExporterImage:   "statsd-exporter:latest",
	}

	rateLimitService := v1alpha1.RateLimitService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "monitoring-test",
			Namespace: "default",
		},
		Spec: v1alpha1.RateLimitServiceSpec{
			Kubernetes: &v1alpha1.RateLimitServiceSpec_Kubernetes{},
			Monitoring: &v1alpha1.RateLimitServiceSpec_Monitoring{
				Enabled: true,
			},
		},
	}

	deployment, err := service.NewDeploymentBuilder(setting).
		SetRateLimitService(rateLimitService).
		Build()

	assert.NoError(t, err)
	assert.NotNil(t, deployment)

	// Should have 2 containers: statsd-exporter and ratelimit
	assert.Len(t, deployment.Spec.Template.Spec.Containers, 2)

	// First container should be statsd-exporter
	assert.Equal(t, "monitoring-test-statsd-exporter", deployment.Spec.Template.Spec.Containers[0].Name)
	assert.Equal(t, "statsd-exporter:latest", deployment.Spec.Template.Spec.Containers[0].Image)

	// Second container should be ratelimit
	assert.Equal(t, "monitoring-test", deployment.Spec.Template.Spec.Containers[1].Name)
	assert.Equal(t, "ratelimit:latest", deployment.Spec.Template.Spec.Containers[1].Image)

	// Should have statsd-config volume
	assert.Len(t, deployment.Spec.Template.Spec.Volumes, 2)
}

func TestDeploymentBuilder_WithReplicaCount(t *testing.T) {
	setting := settings.Settings{
		RateLimitServiceImage: "ratelimit:latest",
		StatsdExporterImage:   "statsd-exporter:latest",
	}

	replicaCount := int32(3)

	rateLimitService := v1alpha1.RateLimitService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "replica-test",
			Namespace: "default",
		},
		Spec: v1alpha1.RateLimitServiceSpec{
			Kubernetes: &v1alpha1.RateLimitServiceSpec_Kubernetes{
				ReplicaCount: &replicaCount,
			},
		},
	}

	deployment, err := service.NewDeploymentBuilder(setting).
		SetRateLimitService(rateLimitService).
		Build()

	assert.NoError(t, err)
	assert.NotNil(t, deployment)
	assert.NotNil(t, deployment.Spec.Replicas)
	assert.Equal(t, int32(3), *deployment.Spec.Replicas)
}

func TestDeploymentBuilder_WithResources(t *testing.T) {
	setting := settings.Settings{
		RateLimitServiceImage: "ratelimit:latest",
		StatsdExporterImage:   "statsd-exporter:latest",
	}

	resources := corev1.ResourceRequirements{
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("500m"),
			corev1.ResourceMemory: resource.MustParse("256Mi"),
		},
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("100m"),
			corev1.ResourceMemory: resource.MustParse("128Mi"),
		},
	}

	rateLimitService := v1alpha1.RateLimitService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "resource-test",
			Namespace: "default",
		},
		Spec: v1alpha1.RateLimitServiceSpec{
			Kubernetes: &v1alpha1.RateLimitServiceSpec_Kubernetes{
				Resources: &resources,
			},
		},
	}

	deployment, err := service.NewDeploymentBuilder(setting).
		SetRateLimitService(rateLimitService).
		Build()

	assert.NoError(t, err)
	assert.NotNil(t, deployment)
	assert.Equal(t, resources, deployment.Spec.Template.Spec.Containers[0].Resources)
}

func TestDeploymentBuilder_WithResourcesAndMonitoring(t *testing.T) {
	setting := settings.Settings{
		RateLimitServiceImage: "ratelimit:latest",
		StatsdExporterImage:   "statsd-exporter:latest",
	}

	resources := corev1.ResourceRequirements{
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("500m"),
			corev1.ResourceMemory: resource.MustParse("256Mi"),
		},
	}

	rateLimitService := v1alpha1.RateLimitService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "resource-monitoring-test",
			Namespace: "default",
		},
		Spec: v1alpha1.RateLimitServiceSpec{
			Kubernetes: &v1alpha1.RateLimitServiceSpec_Kubernetes{
				Resources: &resources,
			},
			Monitoring: &v1alpha1.RateLimitServiceSpec_Monitoring{
				Enabled: true,
			},
		},
	}

	deployment, err := service.NewDeploymentBuilder(setting).
		SetRateLimitService(rateLimitService).
		Build()

	assert.NoError(t, err)
	assert.NotNil(t, deployment)

	// Both containers should have resources
	assert.Equal(t, resources, deployment.Spec.Template.Spec.Containers[0].Resources)
	assert.Equal(t, resources, deployment.Spec.Template.Spec.Containers[1].Resources)
}

func TestDeploymentBuilder_WithCustomImage(t *testing.T) {
	setting := settings.Settings{
		RateLimitServiceImage: "default-ratelimit:latest",
		StatsdExporterImage:   "statsd-exporter:latest",
	}

	customImage := "custom-ratelimit:v1.0.0"

	rateLimitService := v1alpha1.RateLimitService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "custom-image-test",
			Namespace: "default",
		},
		Spec: v1alpha1.RateLimitServiceSpec{
			Kubernetes: &v1alpha1.RateLimitServiceSpec_Kubernetes{
				Image: &customImage,
			},
		},
	}

	deployment, err := service.NewDeploymentBuilder(setting).
		SetRateLimitService(rateLimitService).
		Build()

	assert.NoError(t, err)
	assert.NotNil(t, deployment)
	assert.Equal(t, customImage, deployment.Spec.Template.Spec.Containers[0].Image)
}

func TestDeploymentBuilder_BuildImageInfo(t *testing.T) {
	testCases := []struct {
		name          string
		settings      settings.Settings
		customImage   *string
		expectedImage string
	}{
		{
			name: "use default image from settings",
			settings: settings.Settings{
				RateLimitServiceImage: "default-image:v1",
			},
			customImage:   nil,
			expectedImage: "default-image:v1",
		},
		{
			name: "use custom image when specified",
			settings: settings.Settings{
				RateLimitServiceImage: "default-image:v1",
			},
			customImage:   stringPtr("custom-image:v2"),
			expectedImage: "custom-image:v2",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rateLimitService := v1alpha1.RateLimitService{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "image-test",
					Namespace: "default",
				},
				Spec: v1alpha1.RateLimitServiceSpec{
					Kubernetes: &v1alpha1.RateLimitServiceSpec_Kubernetes{
						Image: tc.customImage,
					},
				},
			}

			builder := service.NewDeploymentBuilder(tc.settings).
				SetRateLimitService(rateLimitService)

			image := builder.BuildImageInfo()
			assert.Equal(t, tc.expectedImage, image)
		})
	}
}

func stringPtr(s string) *string {
	return &s
}

func TestDeploymentBuilder_NilKubernetes(t *testing.T) {
	setting := settings.Settings{
		RateLimitServiceImage: "test:latest",
		StatsdExporterImage:   "statsd:latest",
	}

	rateLimitService := v1alpha1.RateLimitService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-service",
			Namespace: "default",
		},
		Spec: v1alpha1.RateLimitServiceSpec{
			Kubernetes: nil, // nil Kubernetes spec
		},
	}

	// Should not panic when Kubernetes is nil
	assert.NotPanics(t, func() {
		_, err := service.NewDeploymentBuilder(setting).
			SetRateLimitService(rateLimitService).
			Build()
		assert.NoError(t, err)
	})
}
