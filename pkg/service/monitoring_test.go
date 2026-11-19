package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/service"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/settings"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func monitoringRateLimitService(monitoring *v1alpha1.RateLimitServiceSpec_Monitoring) v1alpha1.RateLimitService {
	return v1alpha1.RateLimitService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mon-test",
			Namespace: "default",
		},
		Spec: v1alpha1.RateLimitServiceSpec{
			Monitoring: monitoring,
		},
	}
}

func monitoringSettings() settings.Settings {
	return settings.Settings{
		RateLimitServiceImage: "ratelimit:latest",
		StatsdExporterImage:   "statsd:latest",
	}
}

func findEnv(env []corev1.EnvVar, name string) *corev1.EnvVar {
	for i := range env {
		if env[i].Name == name {
			return &env[i]
		}
	}
	return nil
}

func TestUsesStatsdExporter(t *testing.T) {
	tests := []struct {
		name       string
		monitoring *v1alpha1.RateLimitServiceSpec_Monitoring
		want       bool
	}{
		{"nil", nil, false},
		{"legacy disabled", &v1alpha1.RateLimitServiceSpec_Monitoring{}, false},
		{"legacy enabled", &v1alpha1.RateLimitServiceSpec_Monitoring{Enabled: true}, true},
		{"statsd", &v1alpha1.RateLimitServiceSpec_Monitoring{Type: "statsd"}, true},
		{"prometheus", &v1alpha1.RateLimitServiceSpec_Monitoring{Enabled: true, Type: "prometheus"}, false},
		{"dogstatsd", &v1alpha1.RateLimitServiceSpec_Monitoring{Enabled: true, Type: "dogstatsd"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, service.UsesStatsdExporter(tt.monitoring))
		})
	}
}

func TestPrometheusPort(t *testing.T) {
	tests := []struct {
		name       string
		monitoring *v1alpha1.RateLimitServiceSpec_Monitoring
		want       int32
	}{
		{"nil", nil, 0},
		{"statsd", &v1alpha1.RateLimitServiceSpec_Monitoring{Type: "statsd"}, 0},
		{"default addr", &v1alpha1.RateLimitServiceSpec_Monitoring{Type: "prometheus"}, 9090},
		{"custom port", &v1alpha1.RateLimitServiceSpec_Monitoring{
			Type:       "prometheus",
			Prometheus: &v1alpha1.RateLimitServiceSpec_Monitoring_Prometheus{Addr: ":9191"},
		}, 9191},
		{"host and port", &v1alpha1.RateLimitServiceSpec_Monitoring{
			Type:       "prometheus",
			Prometheus: &v1alpha1.RateLimitServiceSpec_Monitoring_Prometheus{Addr: "0.0.0.0:9292"},
		}, 9292},
		{"invalid addr falls back", &v1alpha1.RateLimitServiceSpec_Monitoring{
			Type:       "prometheus",
			Prometheus: &v1alpha1.RateLimitServiceSpec_Monitoring_Prometheus{Addr: "not-an-addr"},
		}, 9090},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, service.PrometheusPort(tt.monitoring))
		})
	}
}

func TestEnvBuilder_BuildEnv_MonitoringTypeStatsd(t *testing.T) {
	builder := &service.EnvBuilder{
		RateLimitService: monitoringRateLimitService(&v1alpha1.RateLimitServiceSpec_Monitoring{Type: "statsd"}),
	}

	env, err := builder.BuildEnv()
	assert.NoError(t, err)
	assert.Equal(t, "true", env["USE_STATSD"])
	assert.Equal(t, "localhost", env["STATSD_HOST"])
	assert.Equal(t, "9125", env["STATSD_PORT"])
	assert.NotContains(t, env, "USE_DOG_STATSD")
	assert.NotContains(t, env, "USE_PROMETHEUS")
}

func TestEnvBuilder_BuildEnv_MonitoringTypeDogStatsd(t *testing.T) {
	builder := &service.EnvBuilder{
		RateLimitService: monitoringRateLimitService(&v1alpha1.RateLimitServiceSpec_Monitoring{Type: "dogstatsd"}),
	}

	env, err := builder.BuildEnv()
	assert.NoError(t, err)
	// envoy/ratelimit exits when more than one sink is enabled
	assert.Equal(t, "false", env["USE_STATSD"])
	assert.Equal(t, "true", env["USE_DOG_STATSD"])
	assert.NotContains(t, env, "USE_PROMETHEUS")
}

func TestEnvBuilder_BuildEnv_MonitoringTypePrometheus(t *testing.T) {
	builder := &service.EnvBuilder{
		RateLimitService: monitoringRateLimitService(&v1alpha1.RateLimitServiceSpec_Monitoring{Type: "prometheus"}),
	}

	env, err := builder.BuildEnv()
	assert.NoError(t, err)
	assert.Equal(t, "false", env["USE_STATSD"])
	assert.Equal(t, "true", env["USE_PROMETHEUS"])
	assert.NotContains(t, env, "USE_DOG_STATSD")
}

func TestDeploymentBuilder_MonitoringTypeStatsd_AddsSidecar(t *testing.T) {
	deployment, err := service.NewDeploymentBuilder(monitoringSettings()).
		SetRateLimitService(monitoringRateLimitService(&v1alpha1.RateLimitServiceSpec_Monitoring{Type: "statsd"})).
		Build()

	assert.NoError(t, err)
	assert.Len(t, deployment.Spec.Template.Spec.Containers, 2)
	assert.Equal(t, "mon-test-statsd-exporter", deployment.Spec.Template.Spec.Containers[1].Name)
}

func TestDeploymentBuilder_MonitoringTypePrometheus_ExposesMetricsPort(t *testing.T) {
	deployment, err := service.NewDeploymentBuilder(monitoringSettings()).
		SetRateLimitService(monitoringRateLimitService(&v1alpha1.RateLimitServiceSpec_Monitoring{
			Type:       "prometheus",
			Prometheus: &v1alpha1.RateLimitServiceSpec_Monitoring_Prometheus{Addr: ":9191"},
		})).
		Build()

	assert.NoError(t, err)
	assert.Len(t, deployment.Spec.Template.Spec.Containers, 1)

	var metricsPort *corev1.ContainerPort
	for i, port := range deployment.Spec.Template.Spec.Containers[0].Ports {
		if port.Name == "http-metrics" {
			metricsPort = &deployment.Spec.Template.Spec.Containers[0].Ports[i]
		}
	}
	assert.NotNil(t, metricsPort, "http-metrics port should exist")
	assert.Equal(t, int32(9191), metricsPort.ContainerPort)
}

func TestDeploymentBuilder_PrometheusWithResources_DoesNotPanic(t *testing.T) {
	rateLimitService := monitoringRateLimitService(&v1alpha1.RateLimitServiceSpec_Monitoring{
		Enabled: true,
		Type:    "prometheus",
	})
	rateLimitService.Spec.Kubernetes = &v1alpha1.RateLimitServiceSpec_Kubernetes{
		Resources: &corev1.ResourceRequirements{
			Limits: corev1.ResourceList{corev1.ResourceCPU: resource.MustParse("100m")},
		},
	}

	assert.NotPanics(t, func() {
		deployment, err := service.NewDeploymentBuilder(monitoringSettings()).
			SetRateLimitService(rateLimitService).
			Build()
		assert.NoError(t, err)
		assert.Len(t, deployment.Spec.Template.Spec.Containers, 1)
		assert.Equal(t, resource.MustParse("100m"), deployment.Spec.Template.Spec.Containers[0].Resources.Limits[corev1.ResourceCPU])
	})
}

func TestDeploymentBuilder_MonitoringTypeDogStatsd_UsesNodeIP(t *testing.T) {
	deployment, err := service.NewDeploymentBuilder(monitoringSettings()).
		SetRateLimitService(monitoringRateLimitService(&v1alpha1.RateLimitServiceSpec_Monitoring{Type: "dogstatsd"})).
		Build()

	assert.NoError(t, err)
	assert.Len(t, deployment.Spec.Template.Spec.Containers, 1)

	statsdHost := findEnv(deployment.Spec.Template.Spec.Containers[0].Env, "STATSD_HOST")
	assert.NotNil(t, statsdHost, "STATSD_HOST should be set for dogstatsd")
	assert.NotNil(t, statsdHost.ValueFrom)
	assert.Equal(t, "status.hostIP", statsdHost.ValueFrom.FieldRef.FieldPath)
}

func TestDeploymentBuilder_MonitoringTypeDogStatsd_EnvironmentOverridesHost(t *testing.T) {
	rateLimitService := monitoringRateLimitService(&v1alpha1.RateLimitServiceSpec_Monitoring{Type: "dogstatsd"})
	rateLimitService.Spec.Environment = &map[string]string{"STATSD_HOST": "datadog-agent.datadog.svc"}

	deployment, err := service.NewDeploymentBuilder(monitoringSettings()).
		SetRateLimitService(rateLimitService).
		Build()

	assert.NoError(t, err)
	// Container env would win over the env ConfigMap, so it must not be set here
	assert.Nil(t, findEnv(deployment.Spec.Template.Spec.Containers[0].Env, "STATSD_HOST"))
}

func TestServiceBuilder_MonitoringTypePrometheus_ExposesMetricsPort(t *testing.T) {
	svc, err := service.NewServiceBuilder().
		SetRateLimitService(monitoringRateLimitService(&v1alpha1.RateLimitServiceSpec_Monitoring{Type: "prometheus"})).
		Build()

	assert.NoError(t, err)

	var metricsPort *corev1.ServicePort
	for i, port := range svc.Spec.Ports {
		if port.Name == "http-metrics" {
			metricsPort = &svc.Spec.Ports[i]
		}
	}
	assert.NotNil(t, metricsPort, "http-metrics service port should exist")
	assert.Equal(t, int32(9090), metricsPort.Port)
	assert.Equal(t, 9090, metricsPort.TargetPort.IntValue())
}

func TestServiceBuilder_NoPrometheus_NoMetricsPort(t *testing.T) {
	svc, err := service.NewServiceBuilder().
		SetRateLimitService(monitoringRateLimitService(&v1alpha1.RateLimitServiceSpec_Monitoring{Enabled: true})).
		Build()

	assert.NoError(t, err)
	for _, port := range svc.Spec.Ports {
		assert.NotEqual(t, "http-metrics", port.Name)
	}
}

func TestDeploymentBuilder_StatsdSidecar_RateLimitContainerKeepsSettings(t *testing.T) {
	readOnly := true
	rateLimitService := monitoringRateLimitService(&v1alpha1.RateLimitServiceSpec_Monitoring{Enabled: true})
	rateLimitService.Spec.Backend = &v1alpha1.RateLimitServiceSpec_Backend{
		Redis: &v1alpha1.RateLimitServiceSpec_Backend_Redis{
			Type: "single",
			URL:  "redis:6379",
			TLS: &v1alpha1.RateLimitServiceSpec_Backend_Redis_TLS{
				Enabled:   true,
				SecretRef: "redis-tls",
			},
		},
	}
	rateLimitService.Spec.Kubernetes = &v1alpha1.RateLimitServiceSpec_Kubernetes{
		ContainerSecurityContext: &corev1.SecurityContext{ReadOnlyRootFilesystem: &readOnly},
		LivenessProbe:            &corev1.Probe{PeriodSeconds: 7},
	}

	deployment, err := service.NewDeploymentBuilder(monitoringSettings()).
		SetRateLimitService(rateLimitService).
		Build()
	assert.NoError(t, err)

	containers := deployment.Spec.Template.Spec.Containers
	assert.Len(t, containers, 2)
	assert.Equal(t, "mon-test", containers[0].Name)
	assert.Equal(t, "mon-test-statsd-exporter", containers[1].Name)

	foundRedisTLS := false
	for _, vm := range containers[0].VolumeMounts {
		if vm.Name == "redis-tls" {
			foundRedisTLS = true
		}
	}
	assert.True(t, foundRedisTLS, "redis-tls must be mounted on the ratelimit container")
	assert.Equal(t, &readOnly, containers[0].SecurityContext.ReadOnlyRootFilesystem)
	assert.Equal(t, int32(7), containers[0].LivenessProbe.PeriodSeconds)

	for _, vm := range containers[1].VolumeMounts {
		assert.NotEqual(t, "redis-tls", vm.Name)
	}
	assert.Nil(t, containers[1].SecurityContext)
	assert.Nil(t, containers[1].LivenessProbe)
}
