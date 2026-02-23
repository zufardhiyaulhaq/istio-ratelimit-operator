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

func TestDeploymentBuilder_WithPrometheusMonitoring(t *testing.T) {
	setting := settings.Settings{
		RateLimitServiceImage: "ratelimit:latest",
		StatsdExporterImage:   "statsd:latest",
	}

	rateLimitService := v1alpha1.RateLimitService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "prom-test",
			Namespace: "default",
		},
		Spec: v1alpha1.RateLimitServiceSpec{
			Monitoring: &v1alpha1.RateLimitServiceSpec_Monitoring{
				Enabled: true,
				Type:    "prometheus",
			},
		},
	}

	deployment, err := service.NewDeploymentBuilder(setting).
		SetRateLimitService(rateLimitService).
		Build()

	assert.NoError(t, err)
	// Prometheus mode: only 1 container (no statsd sidecar)
	assert.Len(t, deployment.Spec.Template.Spec.Containers, 1)
	assert.Equal(t, "prom-test", deployment.Spec.Template.Spec.Containers[0].Name)
}

func TestDeploymentBuilder_WithRedisTLS(t *testing.T) {
	setting := settings.Settings{
		RateLimitServiceImage: "ratelimit:latest",
		StatsdExporterImage:   "statsd:latest",
	}

	rateLimitService := v1alpha1.RateLimitService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "tls-test",
			Namespace: "default",
		},
		Spec: v1alpha1.RateLimitServiceSpec{
			Backend: &v1alpha1.RateLimitServiceSpec_Backend{
				Redis: &v1alpha1.RateLimitServiceSpec_Backend_Redis{
					Type: "single",
					URL:  "redis:6379",
					TLS: &v1alpha1.RateLimitServiceSpec_Backend_Redis_TLS{
						Enabled:   true,
						SecretRef: "redis-tls-secret",
					},
				},
			},
		},
	}

	deployment, err := service.NewDeploymentBuilder(setting).
		SetRateLimitService(rateLimitService).
		Build()

	assert.NoError(t, err)

	// Should have redis TLS volume
	foundVolume := false
	for _, v := range deployment.Spec.Template.Spec.Volumes {
		if v.Name == "redis-tls" {
			foundVolume = true
			assert.Equal(t, "redis-tls-secret", v.VolumeSource.Secret.SecretName)
		}
	}
	assert.True(t, foundVolume, "redis-tls volume not found")

	// Should have redis TLS volume mount on ratelimit container
	foundMount := false
	for _, vm := range deployment.Spec.Template.Spec.Containers[0].VolumeMounts {
		if vm.Name == "redis-tls" {
			foundMount = true
			assert.Equal(t, "/tls/redis/", vm.MountPath)
			assert.True(t, vm.ReadOnly)
		}
	}
	assert.True(t, foundMount, "redis-tls volume mount not found")
}

func TestDeploymentBuilder_WithGRPCTLS(t *testing.T) {
	setting := settings.Settings{
		RateLimitServiceImage: "ratelimit:latest",
		StatsdExporterImage:   "statsd:latest",
	}

	rateLimitService := v1alpha1.RateLimitService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "grpc-tls-test",
			Namespace: "default",
		},
		Spec: v1alpha1.RateLimitServiceSpec{
			Server: &v1alpha1.RateLimitServiceSpec_Server{
				GRPC: &v1alpha1.RateLimitServiceSpec_Server_GRPC{
					TLS: &v1alpha1.RateLimitServiceSpec_Server_GRPC_TLS{
						Enabled:   true,
						SecretRef: "grpc-tls-secret",
					},
					ClientTLS: &v1alpha1.RateLimitServiceSpec_Server_GRPC_ClientTLS{
						CACertSecretRef: "grpc-ca-secret",
					},
				},
			},
		},
	}

	deployment, err := service.NewDeploymentBuilder(setting).
		SetRateLimitService(rateLimitService).
		Build()

	assert.NoError(t, err)

	// Should have grpc-tls volume
	foundGRPCVolume := false
	foundClientVolume := false
	for _, v := range deployment.Spec.Template.Spec.Volumes {
		if v.Name == "grpc-tls" {
			foundGRPCVolume = true
			assert.Equal(t, "grpc-tls-secret", v.VolumeSource.Secret.SecretName)
		}
		if v.Name == "grpc-client-tls" {
			foundClientVolume = true
			assert.Equal(t, "grpc-ca-secret", v.VolumeSource.Secret.SecretName)
		}
	}
	assert.True(t, foundGRPCVolume, "grpc-tls volume not found")
	assert.True(t, foundClientVolume, "grpc-client-tls volume not found")
}

func TestDeploymentBuilder_WithPerSecondTLS(t *testing.T) {
	setting := settings.Settings{
		RateLimitServiceImage: "ratelimit:latest",
		StatsdExporterImage:   "statsd:latest",
	}

	rateLimitService := v1alpha1.RateLimitService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "ps-tls-test",
			Namespace: "default",
		},
		Spec: v1alpha1.RateLimitServiceSpec{
			Backend: &v1alpha1.RateLimitServiceSpec_Backend{
				Redis: &v1alpha1.RateLimitServiceSpec_Backend_Redis{
					Type: "single",
					URL:  "redis:6379",
					PerSecond: &v1alpha1.RateLimitServiceSpec_Backend_Redis_PerSecond{
						Enabled: true,
						URL:     "redis-ps:6379",
						TLS: &v1alpha1.RateLimitServiceSpec_Backend_Redis_TLS{
							Enabled:   true,
							SecretRef: "redis-ps-tls-secret",
						},
					},
				},
			},
		},
	}

	deployment, err := service.NewDeploymentBuilder(setting).
		SetRateLimitService(rateLimitService).
		Build()

	assert.NoError(t, err)

	foundVolume := false
	for _, v := range deployment.Spec.Template.Spec.Volumes {
		if v.Name == "redis-persecond-tls" {
			foundVolume = true
			assert.Equal(t, "redis-ps-tls-secret", v.VolumeSource.Secret.SecretName)
		}
	}
	assert.True(t, foundVolume, "redis-persecond-tls volume not found")

	foundMount := false
	for _, vm := range deployment.Spec.Template.Spec.Containers[0].VolumeMounts {
		if vm.Name == "redis-persecond-tls" {
			foundMount = true
			assert.Equal(t, "/tls/redis-persecond/", vm.MountPath)
			assert.True(t, vm.ReadOnly)
		}
	}
	assert.True(t, foundMount, "redis-persecond-tls volume mount not found")
}

func boolPtr(b bool) *bool { return &b }

func TestDeploymentBuilder_WithScheduling(t *testing.T) {
	setting := settings.Settings{
		RateLimitServiceImage: "ratelimit:latest",
		StatsdExporterImage:   "statsd:latest",
	}

	rateLimitService := v1alpha1.RateLimitService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sched-test",
			Namespace: "default",
		},
		Spec: v1alpha1.RateLimitServiceSpec{
			Kubernetes: &v1alpha1.RateLimitServiceSpec_Kubernetes{
				NodeSelector: map[string]string{"disktype": "ssd"},
				Tolerations: []corev1.Toleration{
					{Key: "dedicated", Operator: corev1.TolerationOpEqual, Value: "ratelimit", Effect: corev1.TaintEffectNoSchedule},
				},
				Affinity: &corev1.Affinity{
					NodeAffinity: &corev1.NodeAffinity{
						RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
							NodeSelectorTerms: []corev1.NodeSelectorTerm{
								{MatchExpressions: []corev1.NodeSelectorRequirement{
									{Key: "zone", Operator: corev1.NodeSelectorOpIn, Values: []string{"us-east-1a"}},
								}},
							},
						},
					},
				},
				SecurityContext: &corev1.PodSecurityContext{
					RunAsNonRoot: boolPtr(true),
				},
				ContainerSecurityContext: &corev1.SecurityContext{
					ReadOnlyRootFilesystem: boolPtr(true),
				},
				ImagePullSecrets: []corev1.LocalObjectReference{
					{Name: "my-registry-secret"},
				},
				Annotations: &map[string]string{
					"prometheus.io/scrape": "true",
				},
			},
		},
	}

	deployment, err := service.NewDeploymentBuilder(setting).
		SetRateLimitService(rateLimitService).
		Build()

	assert.NoError(t, err)
	assert.Equal(t, map[string]string{"disktype": "ssd"}, deployment.Spec.Template.Spec.NodeSelector)
	assert.Len(t, deployment.Spec.Template.Spec.Tolerations, 1)
	assert.NotNil(t, deployment.Spec.Template.Spec.Affinity)
	assert.NotNil(t, deployment.Spec.Template.Spec.SecurityContext)
	assert.True(t, *deployment.Spec.Template.Spec.SecurityContext.RunAsNonRoot)
	assert.True(t, *deployment.Spec.Template.Spec.Containers[0].SecurityContext.ReadOnlyRootFilesystem)
	assert.Len(t, deployment.Spec.Template.Spec.ImagePullSecrets, 1)
	assert.Equal(t, "true", deployment.Spec.Template.ObjectMeta.Annotations["prometheus.io/scrape"])
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

func TestDeploymentBuilder_WithAuthSecretRef(t *testing.T) {
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
			Backend: &v1alpha1.RateLimitServiceSpec_Backend{
				Redis: &v1alpha1.RateLimitServiceSpec_Backend_Redis{
					Type: "single",
					URL:  "redis:6379",
					AuthSecretRef: &v1alpha1.SecretKeyRef{
						Name: "redis-auth-secret",
						Key:  "password",
					},
				},
			},
		},
	}

	deployment, err := service.NewDeploymentBuilder(setting).
		SetRateLimitService(rateLimitService).
		Build()

	assert.NoError(t, err)

	// Find the REDIS_AUTH env var
	var redisAuthEnv *corev1.EnvVar
	for i, env := range deployment.Spec.Template.Spec.Containers[0].Env {
		if env.Name == "REDIS_AUTH" {
			redisAuthEnv = &deployment.Spec.Template.Spec.Containers[0].Env[i]
			break
		}
	}

	assert.NotNil(t, redisAuthEnv, "REDIS_AUTH env var should exist")
	assert.NotNil(t, redisAuthEnv.ValueFrom)
	assert.NotNil(t, redisAuthEnv.ValueFrom.SecretKeyRef)
	assert.Equal(t, "redis-auth-secret", redisAuthEnv.ValueFrom.SecretKeyRef.Name)
	assert.Equal(t, "password", redisAuthEnv.ValueFrom.SecretKeyRef.Key)
}

func TestDeploymentBuilder_WithSentinelAuth(t *testing.T) {
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
			Backend: &v1alpha1.RateLimitServiceSpec_Backend{
				Redis: &v1alpha1.RateLimitServiceSpec_Backend_Redis{
					Type: "sentinel",
					URL:  "sentinel:26379",
					SentinelAuth: &v1alpha1.RateLimitServiceSpec_Backend_Redis_SentinelAuth{
						SecretRef: &v1alpha1.SecretKeyRef{
							Name: "sentinel-auth-secret",
							Key:  "sentinel-password",
						},
					},
				},
			},
		},
	}

	deployment, err := service.NewDeploymentBuilder(setting).
		SetRateLimitService(rateLimitService).
		Build()

	assert.NoError(t, err)

	// Find the REDIS_SENTINEL_PASSWORD env var
	var sentinelAuthEnv *corev1.EnvVar
	for i, env := range deployment.Spec.Template.Spec.Containers[0].Env {
		if env.Name == "REDIS_SENTINEL_PASSWORD" {
			sentinelAuthEnv = &deployment.Spec.Template.Spec.Containers[0].Env[i]
			break
		}
	}

	assert.NotNil(t, sentinelAuthEnv, "REDIS_SENTINEL_PASSWORD env var should exist")
	assert.NotNil(t, sentinelAuthEnv.ValueFrom)
	assert.NotNil(t, sentinelAuthEnv.ValueFrom.SecretKeyRef)
	assert.Equal(t, "sentinel-auth-secret", sentinelAuthEnv.ValueFrom.SecretKeyRef.Name)
	assert.Equal(t, "sentinel-password", sentinelAuthEnv.ValueFrom.SecretKeyRef.Key)
}

func TestDeploymentBuilder_WithPerSecondAuthSecretRef(t *testing.T) {
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
			Backend: &v1alpha1.RateLimitServiceSpec_Backend{
				Redis: &v1alpha1.RateLimitServiceSpec_Backend_Redis{
					Type: "single",
					URL:  "redis:6379",
					PerSecond: &v1alpha1.RateLimitServiceSpec_Backend_Redis_PerSecond{
						Enabled: true,
						URL:     "redis-persecond:6379",
						AuthSecretRef: &v1alpha1.SecretKeyRef{
							Name: "redis-persecond-auth-secret",
							Key:  "password",
						},
					},
				},
			},
		},
	}

	deployment, err := service.NewDeploymentBuilder(setting).
		SetRateLimitService(rateLimitService).
		Build()

	assert.NoError(t, err)

	// Find the REDIS_PERSECOND_AUTH env var
	var perSecondAuthEnv *corev1.EnvVar
	for i, env := range deployment.Spec.Template.Spec.Containers[0].Env {
		if env.Name == "REDIS_PERSECOND_AUTH" {
			perSecondAuthEnv = &deployment.Spec.Template.Spec.Containers[0].Env[i]
			break
		}
	}

	assert.NotNil(t, perSecondAuthEnv, "REDIS_PERSECOND_AUTH env var should exist")
	assert.NotNil(t, perSecondAuthEnv.ValueFrom)
	assert.NotNil(t, perSecondAuthEnv.ValueFrom.SecretKeyRef)
	assert.Equal(t, "redis-persecond-auth-secret", perSecondAuthEnv.ValueFrom.SecretKeyRef.Name)
	assert.Equal(t, "password", perSecondAuthEnv.ValueFrom.SecretKeyRef.Key)
}
