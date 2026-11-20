package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/service"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func deploymentSpecWith(container corev1.Container) appsv1.DeploymentSpec {
	return appsv1.DeploymentSpec{
		Template: corev1.PodTemplateSpec{
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{container},
			},
		},
	}
}

func TestCheckDeploymentSpecDiff(t *testing.T) {
	base := corev1.Container{
		Name:  "ratelimit",
		Image: "envoyproxy/ratelimit:3dfe2f9a",
		Ports: []corev1.ContainerPort{{Name: "http", ContainerPort: 8080, Protocol: corev1.ProtocolTCP}},
		Env:   []corev1.EnvVar{{Name: "RUNTIME_ROOT", Value: "/data/ratelimit/"}},
	}

	t.Run("identical", func(t *testing.T) {
		assert.False(t, checkDeploymentSpecDiff(deploymentSpecWith(base), deploymentSpecWith(*base.DeepCopy())))
	})

	t.Run("image changed", func(t *testing.T) {
		changed := *base.DeepCopy()
		changed.Image = "envoyproxy/ratelimit:other"
		assert.True(t, checkDeploymentSpecDiff(deploymentSpecWith(base), deploymentSpecWith(changed)))
	})

	t.Run("env added", func(t *testing.T) {
		changed := *base.DeepCopy()
		changed.Env = append(changed.Env, corev1.EnvVar{
			Name: "STATSD_HOST",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{APIVersion: "v1", FieldPath: "status.hostIP"},
			},
		})
		assert.True(t, checkDeploymentSpecDiff(deploymentSpecWith(base), deploymentSpecWith(changed)))
	})

	t.Run("port added", func(t *testing.T) {
		changed := *base.DeepCopy()
		changed.Ports = append(changed.Ports, corev1.ContainerPort{Name: "http-metrics", ContainerPort: 9090, Protocol: corev1.ProtocolTCP})
		assert.True(t, checkDeploymentSpecDiff(deploymentSpecWith(base), deploymentSpecWith(changed)))
	})

	t.Run("container order changed", func(t *testing.T) {
		sidecar := corev1.Container{Name: "ratelimit-statsd-exporter", Image: "statsd"}
		existing := deploymentSpecWith(sidecar)
		existing.Template.Spec.Containers = append(existing.Template.Spec.Containers, base)
		desired := deploymentSpecWith(base)
		desired.Template.Spec.Containers = append(desired.Template.Spec.Containers, sidecar)
		assert.True(t, checkDeploymentSpecDiff(existing, desired))
	})
}

func TestCheckDeploymentTemplateHashDiff(t *testing.T) {
	withHash := func(hash string) *appsv1.Deployment {
		deployment := &appsv1.Deployment{}
		if hash != "" {
			deployment.ObjectMeta = metav1.ObjectMeta{Annotations: map[string]string{service.PodTemplateHashAnnotation: hash}}
		}
		return deployment
	}

	assert.False(t, checkDeploymentTemplateHashDiff(withHash("abc"), withHash("abc")))
	assert.True(t, checkDeploymentTemplateHashDiff(withHash("abc"), withHash("def")))
	// Deployments created before the hash existed get updated once
	assert.True(t, checkDeploymentTemplateHashDiff(withHash(""), withHash("abc")))
}
