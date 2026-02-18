package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/service"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func TestNewServiceBuilder(t *testing.T) {
	builder := service.NewServiceBuilder()
	assert.NotNil(t, builder)
	assert.Equal(t, v1alpha1.RateLimitService{}, builder.RateLimitService)
}

func TestServiceBuilder_SetRateLimitService(t *testing.T) {
	rateLimitService := v1alpha1.RateLimitService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-service",
			Namespace: "test-namespace",
		},
	}

	builder := service.NewServiceBuilder().SetRateLimitService(rateLimitService)

	assert.Equal(t, rateLimitService, builder.RateLimitService)
}

func TestServiceBuilder_SetRateLimitService_Chaining(t *testing.T) {
	rateLimitService := v1alpha1.RateLimitService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-service",
			Namespace: "test-namespace",
		},
	}

	builder := service.NewServiceBuilder()
	returnedBuilder := builder.SetRateLimitService(rateLimitService)

	// Verify method chaining returns the same builder
	assert.Same(t, builder, returnedBuilder)
}

func TestServiceBuilder_BuildLabels(t *testing.T) {
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
					Name:      "another-service",
					Namespace: "production",
				},
			},
			expectedLabels: map[string]string{
				"app.kubernetes.io/name":       "another-service",
				"app.kubernetes.io/managed-by": "istio-rateltimit-operator",
				"app.kubernetes.io/created-by": "another-service",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := service.NewServiceBuilder().SetRateLimitService(tc.rateLimitService)
			labels := builder.BuildLabels()

			assert.Equal(t, tc.expectedLabels, labels)
		})
	}
}

func TestServiceBuilder_Build(t *testing.T) {
	httpAppProtocol := "http"
	grpcAppProtocol := "grpc"
	tcpAppProtocol := "tcp"
	udpAppProtocol := "udp"

	testCases := []struct {
		name             string
		rateLimitService v1alpha1.RateLimitService
		expectedService  *corev1.Service
		expectError      bool
	}{
		{
			name: "basic service build",
			rateLimitService: v1alpha1.RateLimitService{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "ratelimit-svc",
					Namespace: "default",
				},
			},
			expectedService: &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "ratelimit-svc",
					Namespace: "default",
					Labels: map[string]string{
						"app.kubernetes.io/name":       "ratelimit-svc",
						"app.kubernetes.io/managed-by": "istio-rateltimit-operator",
						"app.kubernetes.io/created-by": "ratelimit-svc",
					},
				},
				Spec: corev1.ServiceSpec{
					Selector: map[string]string{
						"app.kubernetes.io/name":       "ratelimit-svc",
						"app.kubernetes.io/managed-by": "istio-rateltimit-operator",
						"app.kubernetes.io/created-by": "ratelimit-svc",
					},
					Ports: []corev1.ServicePort{
						{
							Name:        "http",
							Port:        int32(8080),
							TargetPort:  intstr.FromInt(8080),
							Protocol:    corev1.ProtocolTCP,
							AppProtocol: &httpAppProtocol,
						},
						{
							Name:        "grpc",
							Port:        int32(8081),
							TargetPort:  intstr.FromInt(8081),
							Protocol:    corev1.ProtocolTCP,
							AppProtocol: &grpcAppProtocol,
						},
						{
							Name:        "http-admin",
							Port:        int32(6070),
							TargetPort:  intstr.FromInt(6070),
							Protocol:    corev1.ProtocolTCP,
							AppProtocol: &httpAppProtocol,
						},
						{
							Name:        "http-statsd-exporter",
							Port:        int32(9102),
							TargetPort:  intstr.FromInt(9102),
							Protocol:    corev1.ProtocolTCP,
							AppProtocol: &httpAppProtocol,
						},
						{
							Name:        "tcp-statsd-exporter",
							Port:        int32(9125),
							TargetPort:  intstr.FromInt(9125),
							Protocol:    corev1.ProtocolTCP,
							AppProtocol: &tcpAppProtocol,
						},
						{
							Name:        "udp-statsd-exporter",
							Port:        int32(9125),
							TargetPort:  intstr.FromInt(9125),
							Protocol:    corev1.ProtocolUDP,
							AppProtocol: &udpAppProtocol,
						},
					},
				},
			},
			expectError: false,
		},
		{
			name: "service in different namespace",
			rateLimitService: v1alpha1.RateLimitService{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "prod-ratelimit",
					Namespace: "production",
				},
			},
			expectedService: &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "prod-ratelimit",
					Namespace: "production",
					Labels: map[string]string{
						"app.kubernetes.io/name":       "prod-ratelimit",
						"app.kubernetes.io/managed-by": "istio-rateltimit-operator",
						"app.kubernetes.io/created-by": "prod-ratelimit",
					},
				},
				Spec: corev1.ServiceSpec{
					Selector: map[string]string{
						"app.kubernetes.io/name":       "prod-ratelimit",
						"app.kubernetes.io/managed-by": "istio-rateltimit-operator",
						"app.kubernetes.io/created-by": "prod-ratelimit",
					},
					Ports: []corev1.ServicePort{
						{
							Name:        "http",
							Port:        int32(8080),
							TargetPort:  intstr.FromInt(8080),
							Protocol:    corev1.ProtocolTCP,
							AppProtocol: &httpAppProtocol,
						},
						{
							Name:        "grpc",
							Port:        int32(8081),
							TargetPort:  intstr.FromInt(8081),
							Protocol:    corev1.ProtocolTCP,
							AppProtocol: &grpcAppProtocol,
						},
						{
							Name:        "http-admin",
							Port:        int32(6070),
							TargetPort:  intstr.FromInt(6070),
							Protocol:    corev1.ProtocolTCP,
							AppProtocol: &httpAppProtocol,
						},
						{
							Name:        "http-statsd-exporter",
							Port:        int32(9102),
							TargetPort:  intstr.FromInt(9102),
							Protocol:    corev1.ProtocolTCP,
							AppProtocol: &httpAppProtocol,
						},
						{
							Name:        "tcp-statsd-exporter",
							Port:        int32(9125),
							TargetPort:  intstr.FromInt(9125),
							Protocol:    corev1.ProtocolTCP,
							AppProtocol: &tcpAppProtocol,
						},
						{
							Name:        "udp-statsd-exporter",
							Port:        int32(9125),
							TargetPort:  intstr.FromInt(9125),
							Protocol:    corev1.ProtocolUDP,
							AppProtocol: &udpAppProtocol,
						},
					},
				},
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc, err := service.NewServiceBuilder().
				SetRateLimitService(tc.rateLimitService).
				Build()

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedService, svc)
			}
		})
	}
}

func TestServiceBuilder_Build_ServicePorts(t *testing.T) {
	rateLimitService := v1alpha1.RateLimitService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-svc",
			Namespace: "default",
		},
	}

	svc, err := service.NewServiceBuilder().
		SetRateLimitService(rateLimitService).
		Build()

	assert.NoError(t, err)
	assert.NotNil(t, svc)

	// Verify the number of ports
	assert.Len(t, svc.Spec.Ports, 6)

	// Verify each port configuration
	portChecks := []struct {
		name        string
		port        int32
		targetPort  int
		protocol    corev1.Protocol
		appProtocol string
	}{
		{"http", 8080, 8080, corev1.ProtocolTCP, "http"},
		{"grpc", 8081, 8081, corev1.ProtocolTCP, "grpc"},
		{"http-admin", 6070, 6070, corev1.ProtocolTCP, "http"},
		{"http-statsd-exporter", 9102, 9102, corev1.ProtocolTCP, "http"},
		{"tcp-statsd-exporter", 9125, 9125, corev1.ProtocolTCP, "tcp"},
		{"udp-statsd-exporter", 9125, 9125, corev1.ProtocolUDP, "udp"},
	}

	for i, check := range portChecks {
		t.Run("port_"+check.name, func(t *testing.T) {
			port := svc.Spec.Ports[i]
			assert.Equal(t, check.name, port.Name)
			assert.Equal(t, check.port, port.Port)
			assert.Equal(t, intstr.FromInt(check.targetPort), port.TargetPort)
			assert.Equal(t, check.protocol, port.Protocol)
			assert.NotNil(t, port.AppProtocol)
			assert.Equal(t, check.appProtocol, *port.AppProtocol)
		})
	}
}

func TestServiceBuilder_Build_Metadata(t *testing.T) {
	rateLimitService := v1alpha1.RateLimitService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "metadata-test",
			Namespace: "test-ns",
		},
	}

	svc, err := service.NewServiceBuilder().
		SetRateLimitService(rateLimitService).
		Build()

	assert.NoError(t, err)
	assert.NotNil(t, svc)

	// Verify metadata
	assert.Equal(t, "metadata-test", svc.ObjectMeta.Name)
	assert.Equal(t, "test-ns", svc.ObjectMeta.Namespace)

	// Verify labels are set correctly
	assert.Equal(t, "metadata-test", svc.ObjectMeta.Labels["app.kubernetes.io/name"])
	assert.Equal(t, "istio-rateltimit-operator", svc.ObjectMeta.Labels["app.kubernetes.io/managed-by"])
	assert.Equal(t, "metadata-test", svc.ObjectMeta.Labels["app.kubernetes.io/created-by"])

	// Verify selector matches labels
	assert.Equal(t, svc.ObjectMeta.Labels, svc.Spec.Selector)
}
