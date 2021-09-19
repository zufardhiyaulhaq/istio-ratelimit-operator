package service

import (
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	intstr "k8s.io/apimachinery/pkg/util/intstr"
)

type DeploymentBuilder struct {
	RateLimitService v1alpha1.RateLimitService
}

func NewDeploymentBuilder() *DeploymentBuilder {
	return &DeploymentBuilder{}
}

func (n *DeploymentBuilder) SetRateLimitService(rateLimitService v1alpha1.RateLimitService) *DeploymentBuilder {
	n.RateLimitService = rateLimitService
	return n
}

func (n *DeploymentBuilder) Build() (*appsv1.Deployment, error) {

	image := "zufardhiyaulhaq/ratelimit"
	imageTag := "v1.0.0"
	env := n.BuildEnv()

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      n.RateLimitService.Name,
			Namespace: n.RateLimitService.Namespace,
			Labels:    n.buildLabels(),
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: n.buildLabels(),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: n.buildLabels(),
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    n.RateLimitService.Name,
							Image:   image + ":" + imageTag,
							Command: []string{"/bin/ratelimit"},
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: int32(8080),
								},
								{
									Name:          "grpc",
									ContainerPort: int32(8081),
								},
								{
									Name:          "http-admin",
									ContainerPort: int32(6070),
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
								Handler: corev1.Handler{
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
									MountPath: "/data/ratelimit/config/config.yaml",
									SubPath:   "config.yaml",
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
						{
							Name: n.RateLimitService.Name + "-config-env",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: n.RateLimitService.Name + "-config-env",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	if n.RateLimitService.Spec.Kubernetes != nil {
		if n.RateLimitService.Spec.Kubernetes.ReplicaCount != nil {
			deployment.Spec.Replicas = n.RateLimitService.Spec.Kubernetes.ReplicaCount
		}

		if n.RateLimitService.Spec.Kubernetes.Resources != nil {
			deployment.Spec.Template.Spec.Containers[0].Resources = *n.RateLimitService.Spec.Kubernetes.Resources
		}
	}

	return deployment, nil
}

func (n *DeploymentBuilder) BuildEnv() []corev1.EnvVar {
	env := []corev1.EnvVar{
		{
			Name:  "RUNTIME_ROOT",
			Value: "/data",
		},
		{
			Name:  "RUNTIME_SUBDIRECTORY",
			Value: "ratelimit",
		},
	}

	return env
}

func (n *DeploymentBuilder) buildLabels() map[string]string {
	var labels = map[string]string{
		"app.kubernetes.io/name":       n.RateLimitService.Name,
		"app.kubernetes.io/managed-by": "istio-rateltimit-operator",
		"app.kubernetes.io/created-by": n.RateLimitService.Name,
	}

	return labels
}
