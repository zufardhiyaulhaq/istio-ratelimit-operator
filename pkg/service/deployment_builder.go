package service

import (
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/settings"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	intstr "k8s.io/apimachinery/pkg/util/intstr"
)

type DeploymentBuilder struct {
	RateLimitService v1alpha1.RateLimitService
	Settings         settings.Settings
}

func NewDeploymentBuilder(settings settings.Settings) *DeploymentBuilder {
	return &DeploymentBuilder{
		Settings: settings,
	}
}

func (n *DeploymentBuilder) SetRateLimitService(rateLimitService v1alpha1.RateLimitService) *DeploymentBuilder {
	n.RateLimitService = rateLimitService
	return n
}

func (n *DeploymentBuilder) Build() (*appsv1.Deployment, error) {
	image := n.BuildImageInfo()
	env := n.BuildEnv()

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      n.RateLimitService.Name,
			Namespace: n.RateLimitService.Namespace,
			Labels:    n.BuildLabels(),
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: n.BuildLabelsSelector(),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: n.BuildLabels(),
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    n.RateLimitService.Name,
							Image:   image,
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
							Env: env,
							EnvFrom: []corev1.EnvFromSource{
								{
									ConfigMapRef: &corev1.ConfigMapEnvSource{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: n.RateLimitService.Name + "-config-env",
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
									Name:      n.RateLimitService.Name + "-config",
									MountPath: "/data/ratelimit/config/",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: n.RateLimitService.Name + "-config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: n.RateLimitService.Name + "-config",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	if n.RateLimitService.Spec.Monitoring != nil {
		// Only add statsd sidecar when Type is empty (legacy behavior)
		if n.RateLimitService.Spec.Monitoring.Enabled && n.RateLimitService.Spec.Monitoring.Type == "" {
			deployment.Spec.Template.Spec.Containers = append([]corev1.Container{
				{Name: n.RateLimitService.Name + "-statsd-exporter",
					Image:   n.Settings.StatsdExporterImage,
					Command: []string{"/bin/statsd_exporter"},
					Args: []string{
						"--web.enable-lifecycle",
						"--statsd.mapping-config=/etc/prometheus-statsd-exporter/statsd.mappingConf",
					},
					Ports: []corev1.ContainerPort{
						{
							Name:          "http-statsd",
							ContainerPort: int32(9102),
							Protocol:      corev1.ProtocolTCP,
						},
						{
							Name:          "tcp-statsd",
							ContainerPort: int32(9125),
							Protocol:      corev1.ProtocolTCP,
						},
						{
							Name:          "udp-statsd",
							ContainerPort: int32(9125),
							Protocol:      corev1.ProtocolUDP,
						},
					},
					ReadinessProbe: &corev1.Probe{
						ProbeHandler: corev1.ProbeHandler{
							HTTPGet: &corev1.HTTPGetAction{
								Path: "/ready",
								Port: intstr.FromInt(9102),
							},
						},
						InitialDelaySeconds: int32(5),
						PeriodSeconds:       int32(10),
					},
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      n.RateLimitService.Name + "-statsd-config",
							MountPath: "/etc/prometheus-statsd-exporter/",
						},
					}},
			}, deployment.Spec.Template.Spec.Containers...)

			deployment.Spec.Template.Spec.Volumes = append(deployment.Spec.Template.Spec.Volumes, corev1.Volume{
				Name: n.RateLimitService.Name + "-statsd-config",
				VolumeSource: corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: n.RateLimitService.Name + "-statsd-config",
						},
					},
				},
			})
		}
	}

	// TLS volume mounts for Redis
	if n.RateLimitService.Spec.Backend != nil && n.RateLimitService.Spec.Backend.Redis != nil &&
		n.RateLimitService.Spec.Backend.Redis.TLS != nil && n.RateLimitService.Spec.Backend.Redis.TLS.SecretRef != "" {
		deployment.Spec.Template.Spec.Volumes = append(deployment.Spec.Template.Spec.Volumes, corev1.Volume{
			Name: "redis-tls",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: n.RateLimitService.Spec.Backend.Redis.TLS.SecretRef,
				},
			},
		})
		deployment.Spec.Template.Spec.Containers[0].VolumeMounts = append(
			deployment.Spec.Template.Spec.Containers[0].VolumeMounts,
			corev1.VolumeMount{
				Name:      "redis-tls",
				MountPath: "/tls/redis/",
				ReadOnly:  true,
			},
		)
	}

	// TLS volume mounts for Redis PerSecond
	if n.RateLimitService.Spec.Backend != nil && n.RateLimitService.Spec.Backend.Redis != nil &&
		n.RateLimitService.Spec.Backend.Redis.PerSecond != nil && n.RateLimitService.Spec.Backend.Redis.PerSecond.TLS != nil &&
		n.RateLimitService.Spec.Backend.Redis.PerSecond.TLS.SecretRef != "" {
		deployment.Spec.Template.Spec.Volumes = append(deployment.Spec.Template.Spec.Volumes, corev1.Volume{
			Name: "redis-persecond-tls",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: n.RateLimitService.Spec.Backend.Redis.PerSecond.TLS.SecretRef,
				},
			},
		})
		deployment.Spec.Template.Spec.Containers[0].VolumeMounts = append(
			deployment.Spec.Template.Spec.Containers[0].VolumeMounts,
			corev1.VolumeMount{
				Name:      "redis-persecond-tls",
				MountPath: "/tls/redis-persecond/",
				ReadOnly:  true,
			},
		)
	}

	// TLS volume mounts for gRPC
	if n.RateLimitService.Spec.Server != nil && n.RateLimitService.Spec.Server.GRPC != nil {
		grpc := n.RateLimitService.Spec.Server.GRPC
		if grpc.TLS != nil && grpc.TLS.SecretRef != "" {
			deployment.Spec.Template.Spec.Volumes = append(deployment.Spec.Template.Spec.Volumes, corev1.Volume{
				Name: "grpc-tls",
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: grpc.TLS.SecretRef,
					},
				},
			})
			deployment.Spec.Template.Spec.Containers[0].VolumeMounts = append(
				deployment.Spec.Template.Spec.Containers[0].VolumeMounts,
				corev1.VolumeMount{
					Name:      "grpc-tls",
					MountPath: "/tls/grpc/",
					ReadOnly:  true,
				},
			)
		}
		if grpc.ClientTLS != nil && grpc.ClientTLS.CACertSecretRef != "" {
			deployment.Spec.Template.Spec.Volumes = append(deployment.Spec.Template.Spec.Volumes, corev1.Volume{
				Name: "grpc-client-tls",
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: grpc.ClientTLS.CACertSecretRef,
					},
				},
			})
			deployment.Spec.Template.Spec.Containers[0].VolumeMounts = append(
				deployment.Spec.Template.Spec.Containers[0].VolumeMounts,
				corev1.VolumeMount{
					Name:      "grpc-client-tls",
					MountPath: "/tls/grpc-client/",
					ReadOnly:  true,
				},
			)
		}
	}

	if n.RateLimitService.Spec.Kubernetes != nil {
		if n.RateLimitService.Spec.Kubernetes.ReplicaCount != nil {
			deployment.Spec.Replicas = n.RateLimitService.Spec.Kubernetes.ReplicaCount
		}

		if n.RateLimitService.Spec.Kubernetes.Resources != nil {
			deployment.Spec.Template.Spec.Containers[0].Resources = *n.RateLimitService.Spec.Kubernetes.Resources

			if n.RateLimitService.Spec.Monitoring != nil {
				if n.RateLimitService.Spec.Monitoring.Enabled {
					deployment.Spec.Template.Spec.Containers[1].Resources = *n.RateLimitService.Spec.Kubernetes.Resources
				}
			}
		}

		// Kubernetes scheduling fields
		if n.RateLimitService.Spec.Kubernetes.NodeSelector != nil {
			deployment.Spec.Template.Spec.NodeSelector = n.RateLimitService.Spec.Kubernetes.NodeSelector
		}

		if n.RateLimitService.Spec.Kubernetes.Tolerations != nil {
			deployment.Spec.Template.Spec.Tolerations = n.RateLimitService.Spec.Kubernetes.Tolerations
		}

		if n.RateLimitService.Spec.Kubernetes.Affinity != nil {
			deployment.Spec.Template.Spec.Affinity = n.RateLimitService.Spec.Kubernetes.Affinity
		}

		if n.RateLimitService.Spec.Kubernetes.SecurityContext != nil {
			deployment.Spec.Template.Spec.SecurityContext = n.RateLimitService.Spec.Kubernetes.SecurityContext
		}

		if n.RateLimitService.Spec.Kubernetes.ContainerSecurityContext != nil {
			deployment.Spec.Template.Spec.Containers[0].SecurityContext = n.RateLimitService.Spec.Kubernetes.ContainerSecurityContext
		}

		if n.RateLimitService.Spec.Kubernetes.ImagePullSecrets != nil {
			deployment.Spec.Template.Spec.ImagePullSecrets = n.RateLimitService.Spec.Kubernetes.ImagePullSecrets
		}

		if n.RateLimitService.Spec.Kubernetes.Annotations != nil {
			deployment.Spec.Template.ObjectMeta.Annotations = *n.RateLimitService.Spec.Kubernetes.Annotations
		}

		if n.RateLimitService.Spec.Kubernetes.LivenessProbe != nil {
			deployment.Spec.Template.Spec.Containers[0].LivenessProbe = n.RateLimitService.Spec.Kubernetes.LivenessProbe
		}
	}

	return deployment, nil
}

func (n *DeploymentBuilder) BuildEnv() []corev1.EnvVar {
	env := []corev1.EnvVar{
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
	}

	// Add Redis auth from secret reference
	if n.RateLimitService.Spec.Backend != nil && n.RateLimitService.Spec.Backend.Redis != nil {
		redis := n.RateLimitService.Spec.Backend.Redis

		// Redis authSecretRef
		if redis.AuthSecretRef != nil {
			env = append(env, corev1.EnvVar{
				Name: "REDIS_AUTH",
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: redis.AuthSecretRef.Name,
						},
						Key: redis.AuthSecretRef.Key,
					},
				},
			})
		}

		// Redis Sentinel auth
		if redis.SentinelAuth != nil && redis.SentinelAuth.SecretRef != nil {
			env = append(env, corev1.EnvVar{
				Name: "REDIS_SENTINEL_PASSWORD",
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: redis.SentinelAuth.SecretRef.Name,
						},
						Key: redis.SentinelAuth.SecretRef.Key,
					},
				},
			})
		}

		// Redis PerSecond authSecretRef
		if redis.PerSecond != nil && redis.PerSecond.AuthSecretRef != nil {
			env = append(env, corev1.EnvVar{
				Name: "REDIS_PERSECOND_AUTH",
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: redis.PerSecond.AuthSecretRef.Name,
						},
						Key: redis.PerSecond.AuthSecretRef.Key,
					},
				},
			})
		}
	}

	return env
}

func (n *DeploymentBuilder) BuildLabelsSelector() map[string]string {
	var labels = map[string]string{
		"app.kubernetes.io/name":       n.RateLimitService.Name,
		"app.kubernetes.io/managed-by": "istio-rateltimit-operator",
		"app.kubernetes.io/created-by": n.RateLimitService.Name,
	}

	return labels
}

func (n *DeploymentBuilder) BuildLabels() map[string]string {
	labelSelector := n.BuildLabelsSelector()

	var labels = map[string]string{}
	for key, value := range labelSelector {
		labels[key] = value
	}

	if n.RateLimitService.Spec.Kubernetes != nil && n.RateLimitService.Spec.Kubernetes.ExtraLabels != nil {
		for key, value := range *n.RateLimitService.Spec.Kubernetes.ExtraLabels {
			labels[key] = value
		}
	}

	return labels
}

func (n *DeploymentBuilder) BuildImageInfo() string {
	if n.RateLimitService.Spec.Kubernetes != nil && n.RateLimitService.Spec.Kubernetes.Image != nil {
		return *n.RateLimitService.Spec.Kubernetes.Image
	}

	return n.Settings.RateLimitServiceImage
}
