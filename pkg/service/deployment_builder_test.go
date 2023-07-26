package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/service"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/settings"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
				ObjectMeta: v1.ObjectMeta{
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
				ObjectMeta: v1.ObjectMeta{
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
