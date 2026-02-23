package service_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/service"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNewEnvBuilder(t *testing.T) {
	builder := service.NewEnvBuilder()
	assert.NotNil(t, builder)
	assert.Equal(t, v1alpha1.RateLimitService{}, builder.RateLimitService)
}

func TestEnvBuilder_SetRateLimitService(t *testing.T) {
	rateLimitService := v1alpha1.RateLimitService{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-env",
			Namespace: "test-namespace",
		},
	}

	builder := service.NewEnvBuilder().SetRateLimitService(rateLimitService)

	assert.Equal(t, rateLimitService, builder.RateLimitService)
}

func TestEnvBuilder_SetRateLimitService_Chaining(t *testing.T) {
	rateLimitService := v1alpha1.RateLimitService{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-env",
			Namespace: "test-namespace",
		},
	}

	builder := service.NewEnvBuilder()
	returnedBuilder := builder.SetRateLimitService(rateLimitService)

	// Verify method chaining returns the same builder
	assert.Same(t, builder, returnedBuilder)
}

func TestEnvBuilder_BuildStatsdEnv(t *testing.T) {
	builder := &service.EnvBuilder{}

	env, err := builder.BuildStatsdEnv()

	assert.NoError(t, err)
	assert.Equal(t, "true", env["USE_STATSD"])
	assert.Equal(t, "localhost", env["STATSD_HOST"])
	assert.Equal(t, "9125", env["STATSD_PORT"])
}

func TestEnvBuilder_BuildRedisEnv(t *testing.T) {
	mockPipelineLimit := 1
	mockPipelineWindow := "1s"

	type fields struct {
		RateLimitService v1alpha1.RateLimitService
	}
	tests := []struct {
		name    string
		fields  fields
		want    map[string]string
		wantErr bool
	}{
		{
			name: "simple redis",
			fields: fields{
				RateLimitService: v1alpha1.RateLimitService{
					Spec: v1alpha1.RateLimitServiceSpec{
						Backend: &v1alpha1.RateLimitServiceSpec_Backend{
							Redis: &v1alpha1.RateLimitServiceSpec_Backend_Redis{
								Type: "single",
								URL:  "127.0.0.1:6379",
							},
						},
					},
				},
			},
			want: map[string]string{
				"REDIS_SOCKET_TYPE": "tcp",
				"REDIS_TYPE":        "single",
				"REDIS_URL":         "127.0.0.1:6379",
			},
			wantErr: false,
		},
		{
			name: "auth redis",
			fields: fields{
				RateLimitService: v1alpha1.RateLimitService{
					Spec: v1alpha1.RateLimitServiceSpec{
						Backend: &v1alpha1.RateLimitServiceSpec_Backend{
							Redis: &v1alpha1.RateLimitServiceSpec_Backend_Redis{
								Type: "single",
								URL:  "127.0.0.1:6379",
								Auth: "password",
							},
						},
					},
				},
			},
			want: map[string]string{
				"REDIS_SOCKET_TYPE": "tcp",
				"REDIS_TYPE":        "single",
				"REDIS_URL":         "127.0.0.1:6379",
				"REDIS_AUTH":        "password",
			},
			wantErr: false,
		},
		{
			name: "redis pipeline limit",
			fields: fields{
				RateLimitService: v1alpha1.RateLimitService{
					Spec: v1alpha1.RateLimitServiceSpec{
						Backend: &v1alpha1.RateLimitServiceSpec_Backend{
							Redis: &v1alpha1.RateLimitServiceSpec_Backend_Redis{
								Type: "single",
								URL:  "127.0.0.1:6379",
								Auth: "password",
								Config: &v1alpha1.RateLimitServiceSpec_Backend_Redis_Config{
									PipelineLimit: &mockPipelineLimit,
								},
							},
						},
					},
				},
			},
			want: map[string]string{
				"REDIS_SOCKET_TYPE":    "tcp",
				"REDIS_TYPE":           "single",
				"REDIS_URL":            "127.0.0.1:6379",
				"REDIS_AUTH":           "password",
				"REDIS_PIPELINE_LIMIT": "1",
			},
			wantErr: false,
		},
		{
			name: "redis pipeline window",
			fields: fields{
				RateLimitService: v1alpha1.RateLimitService{
					Spec: v1alpha1.RateLimitServiceSpec{
						Backend: &v1alpha1.RateLimitServiceSpec_Backend{
							Redis: &v1alpha1.RateLimitServiceSpec_Backend_Redis{
								Type: "single",
								URL:  "127.0.0.1:6379",
								Auth: "password",
								Config: &v1alpha1.RateLimitServiceSpec_Backend_Redis_Config{
									PipelineWindow: &mockPipelineWindow,
								},
							},
						},
					},
				},
			},
			want: map[string]string{
				"REDIS_SOCKET_TYPE":     "tcp",
				"REDIS_TYPE":            "single",
				"REDIS_URL":             "127.0.0.1:6379",
				"REDIS_AUTH":            "password",
				"REDIS_PIPELINE_WINDOW": "1s",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &service.EnvBuilder{
				RateLimitService: tt.fields.RateLimitService,
			}
			got, err := n.BuildRedisEnv()
			if (err != nil) != tt.wantErr {
				t.Errorf("EnvBuilder.buildRedisEnv() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EnvBuilder.buildRedisEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnvBuilder_BuildDefaultEnv(t *testing.T) {
	type fields struct {
		RateLimitService v1alpha1.RateLimitService
	}
	tests := []struct {
		name    string
		fields  fields
		want    map[string]string
		wantErr bool
	}{
		{
			name:   "default env",
			fields: fields{},
			want: map[string]string{
				"USE_STATSD": "false",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &service.EnvBuilder{
				RateLimitService: tt.fields.RateLimitService,
			}
			got, err := n.BuildDefaultEnv()
			if (err != nil) != tt.wantErr {
				t.Errorf("EnvBuilder.buildDefaultEnv() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EnvBuilder.buildDefaultEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnvBuilder_BuildLabels(t *testing.T) {
	type fields struct {
		RateLimitService v1alpha1.RateLimitService
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]string
	}{
		{
			name: "have correct name",
			fields: fields{
				RateLimitService: v1alpha1.RateLimitService{
					ObjectMeta: v1.ObjectMeta{
						Name: "foo",
					},
				},
			},
			want: map[string]string{
				"app.kubernetes.io/name":       "foo-config-env",
				"app.kubernetes.io/managed-by": "istio-rateltimit-operator",
				"app.kubernetes.io/created-by": "foo",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &service.EnvBuilder{
				RateLimitService: tt.fields.RateLimitService,
			}
			if got := n.BuildLabels(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EnvBuilder.BuildLabels() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnvBuilder_BuildRedisEnv_WithTLS(t *testing.T) {
	builder := &service.EnvBuilder{
		RateLimitService: v1alpha1.RateLimitService{
			Spec: v1alpha1.RateLimitServiceSpec{
				Backend: &v1alpha1.RateLimitServiceSpec_Backend{
					Redis: &v1alpha1.RateLimitServiceSpec_Backend_Redis{
						Type: "single",
						URL:  "redis:6379",
						TLS: &v1alpha1.RateLimitServiceSpec_Backend_Redis_TLS{
							Enabled: true,
						},
					},
				},
			},
		},
	}

	env, err := builder.BuildRedisEnv()
	assert.NoError(t, err)
	assert.Equal(t, "true", env["REDIS_TLS"])
}

func TestEnvBuilder_BuildRedisEnv_WithTLSCert(t *testing.T) {
	builder := &service.EnvBuilder{
		RateLimitService: v1alpha1.RateLimitService{
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
		},
	}

	env, err := builder.BuildRedisEnv()
	assert.NoError(t, err)
	assert.Equal(t, "true", env["REDIS_TLS"])
	assert.Equal(t, "/tls/redis/ca.crt", env["REDIS_TLS_CACERT"])
	assert.Equal(t, "/tls/redis/tls.crt", env["REDIS_TLS_CLIENT_CERT"])
	assert.Equal(t, "/tls/redis/tls.key", env["REDIS_TLS_CLIENT_KEY"])
}

func TestEnvBuilder_BuildRedisEnv_WithTLSSkipVerify(t *testing.T) {
	builder := &service.EnvBuilder{
		RateLimitService: v1alpha1.RateLimitService{
			Spec: v1alpha1.RateLimitServiceSpec{
				Backend: &v1alpha1.RateLimitServiceSpec_Backend{
					Redis: &v1alpha1.RateLimitServiceSpec_Backend_Redis{
						Type: "single",
						URL:  "redis:6379",
						TLS: &v1alpha1.RateLimitServiceSpec_Backend_Redis_TLS{
							Enabled:                  true,
							SkipHostnameVerification: true,
						},
					},
				},
			},
		},
	}

	env, err := builder.BuildRedisEnv()
	assert.NoError(t, err)
	assert.Equal(t, "true", env["REDIS_TLS"])
	assert.Equal(t, "true", env["REDIS_TLS_SKIP_HOSTNAME_VERIFICATION"])
}

func strPtr(s string) *string { return &s }

func TestEnvBuilder_BuildServerEnv(t *testing.T) {
	grpcPort := int32(8081)
	builder := &service.EnvBuilder{
		RateLimitService: v1alpha1.RateLimitService{
			Spec: v1alpha1.RateLimitServiceSpec{
				Server: &v1alpha1.RateLimitServiceSpec_Server{
					GRPC: &v1alpha1.RateLimitServiceSpec_Server_GRPC{
						Port: &grpcPort,
						TLS: &v1alpha1.RateLimitServiceSpec_Server_GRPC_TLS{
							Enabled:   true,
							SecretRef: "grpc-tls-secret",
						},
						MaxConnectionAge:      strPtr("30m"),
						MaxConnectionAgeGrace: strPtr("5m"),
					},
				},
			},
		},
	}

	env, err := builder.BuildServerEnv()
	assert.NoError(t, err)
	assert.Equal(t, "8081", env["GRPC_PORT"])
	assert.Equal(t, "/tls/grpc/tls.crt", env["GRPC_SERVER_TLS_CERT"])
	assert.Equal(t, "/tls/grpc/tls.key", env["GRPC_SERVER_TLS_KEY"])
	assert.Equal(t, "30m", env["GRPC_MAX_CONNECTION_AGE"])
	assert.Equal(t, "5m", env["GRPC_MAX_CONNECTION_AGE_GRACE"])
}

func TestEnvBuilder_BuildServerEnv_WithClientTLS(t *testing.T) {
	builder := &service.EnvBuilder{
		RateLimitService: v1alpha1.RateLimitService{
			Spec: v1alpha1.RateLimitServiceSpec{
				Server: &v1alpha1.RateLimitServiceSpec_Server{
					GRPC: &v1alpha1.RateLimitServiceSpec_Server_GRPC{
						ClientTLS: &v1alpha1.RateLimitServiceSpec_Server_GRPC_ClientTLS{
							CACertSecretRef: "grpc-ca-secret",
							SAN:             "ratelimit.example.com",
						},
					},
				},
			},
		},
	}

	env, err := builder.BuildServerEnv()
	assert.NoError(t, err)
	assert.Equal(t, "/tls/grpc-client/ca.crt", env["GRPC_SERVER_TLS_CLIENT_CACERT"])
	assert.Equal(t, "ratelimit.example.com", env["GRPC_CLIENT_TLS_SAN"])
}

func TestEnvBuilder_BuildMonitoringEnv_Prometheus(t *testing.T) {
	builder := &service.EnvBuilder{
		RateLimitService: v1alpha1.RateLimitService{
			Spec: v1alpha1.RateLimitServiceSpec{
				Monitoring: &v1alpha1.RateLimitServiceSpec_Monitoring{
					Enabled: true,
					Type:    "prometheus",
					Prometheus: &v1alpha1.RateLimitServiceSpec_Monitoring_Prometheus{
						Addr: ":9102",
						Path: "/metrics",
					},
					NearLimitRatio:     strPtr("0.8"),
					StatsFlushInterval: strPtr("10s"),
				},
			},
		},
	}

	env, err := builder.BuildMonitoringEnv()
	assert.NoError(t, err)
	assert.Equal(t, "false", env["USE_STATSD"])
	assert.Equal(t, "true", env["USE_PROMETHEUS"])
	assert.Equal(t, ":9102", env["PROMETHEUS_ADDR"])
	assert.Equal(t, "/metrics", env["PROMETHEUS_PATH"])
	assert.Equal(t, "0.8", env["NEAR_LIMIT_RATIO"])
	assert.Equal(t, "10s", env["STATS_FLUSH_INTERVAL"])
}

func TestEnvBuilder_BuildMonitoringEnv_Tracing(t *testing.T) {
	builder := &service.EnvBuilder{
		RateLimitService: v1alpha1.RateLimitService{
			Spec: v1alpha1.RateLimitServiceSpec{
				Monitoring: &v1alpha1.RateLimitServiceSpec_Monitoring{
					Tracing: &v1alpha1.RateLimitServiceSpec_Monitoring_Tracing{
						Enabled:          true,
						ExporterProtocol: "http",
						ServiceName:      "ratelimit",
						ServiceNamespace: "default",
						SamplingRate:     "0.1",
					},
				},
			},
		},
	}

	env, err := builder.BuildMonitoringEnv()
	assert.NoError(t, err)
	assert.Equal(t, "true", env["TRACING_ENABLED"])
	assert.Equal(t, "http", env["TRACING_EXPORTER_PROTOCOL"])
	assert.Equal(t, "ratelimit", env["TRACING_SERVICE_NAME"])
	assert.Equal(t, "default", env["TRACING_SERVICE_NAMESPACE"])
	assert.Equal(t, "0.1", env["TRACING_SAMPLING_RATE"])
}

func TestEnvBuilder_BuildResponseHeadersEnv(t *testing.T) {
	builder := &service.EnvBuilder{
		RateLimitService: v1alpha1.RateLimitService{
			Spec: v1alpha1.RateLimitServiceSpec{
				ResponseHeaders: &v1alpha1.RateLimitServiceSpec_ResponseHeaders{
					Enabled: true,
				},
			},
		},
	}

	env, err := builder.BuildResponseHeadersEnv()
	assert.NoError(t, err)
	assert.Equal(t, "true", env["RESPONSE_HEADERS_ENABLED"])
}

func TestEnvBuilder_BuildLoggingEnv(t *testing.T) {
	builder := &service.EnvBuilder{
		RateLimitService: v1alpha1.RateLimitService{
			Spec: v1alpha1.RateLimitServiceSpec{
				Logging: &v1alpha1.RateLimitServiceSpec_Logging{
					Level:  "debug",
					Format: "json",
				},
			},
		},
	}

	env, err := builder.BuildLoggingEnv()
	assert.NoError(t, err)
	assert.Equal(t, "debug", env["LOG_LEVEL"])
	assert.Equal(t, "json", env["LOG_FORMAT"])
}

func TestEnvBuilder_BuildShadowModeEnv(t *testing.T) {
	builder := &service.EnvBuilder{
		RateLimitService: v1alpha1.RateLimitService{
			Spec: v1alpha1.RateLimitServiceSpec{
				ShadowMode: true,
			},
		},
	}

	env, err := builder.BuildShadowModeEnv()
	assert.NoError(t, err)
	assert.Equal(t, "true", env["SHADOW_MODE"])
}

func TestEnvBuilder_BuildRedisEnv_WithPool(t *testing.T) {
	timeout := "2s"
	builder := &service.EnvBuilder{
		RateLimitService: v1alpha1.RateLimitService{
			Spec: v1alpha1.RateLimitServiceSpec{
				Backend: &v1alpha1.RateLimitServiceSpec_Backend{
					Redis: &v1alpha1.RateLimitServiceSpec_Backend_Redis{
						Type: "single",
						URL:  "redis:6379",
						Pool: &v1alpha1.RateLimitServiceSpec_Backend_Redis_Pool{
							Size: 10,
						},
						Timeout:                     &timeout,
						HealthCheckActiveConnection: true,
					},
				},
			},
		},
	}

	env, err := builder.BuildRedisEnv()
	assert.NoError(t, err)
	assert.Equal(t, "10", env["REDIS_POOL_SIZE"])
	assert.Equal(t, "2s", env["REDIS_TIMEOUT"])
	assert.Equal(t, "true", env["REDIS_HEALTH_CHECK_ACTIVE_CONNECTION"])
}

func TestEnvBuilder_BuildRedisEnv_WithPerSecond(t *testing.T) {
	builder := &service.EnvBuilder{
		RateLimitService: v1alpha1.RateLimitService{
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
								SecretRef: "redis-ps-tls",
							},
							Pool: &v1alpha1.RateLimitServiceSpec_Backend_Redis_Pool{
								Size: 5,
							},
						},
					},
				},
			},
		},
	}

	env, err := builder.BuildRedisEnv()
	assert.NoError(t, err)
	assert.Equal(t, "true", env["REDIS_PERSECOND"])
	assert.Equal(t, "redis-ps:6379", env["REDIS_PERSECOND_URL"])
	assert.Equal(t, "tcp", env["REDIS_PERSECOND_SOCKET_TYPE"])
	assert.Equal(t, "true", env["REDIS_PERSECOND_TLS"])
	assert.Equal(t, "/tls/redis-persecond/ca.crt", env["REDIS_PERSECOND_TLS_CACERT"])
	assert.Equal(t, "/tls/redis-persecond/tls.crt", env["REDIS_PERSECOND_TLS_CLIENT_CERT"])
	assert.Equal(t, "/tls/redis-persecond/tls.key", env["REDIS_PERSECOND_TLS_CLIENT_KEY"])
	assert.Equal(t, "5", env["REDIS_PERSECOND_POOL_SIZE"])
}

func TestEnvBuilder_BuildBackendEnv(t *testing.T) {
	builder := &service.EnvBuilder{
		RateLimitService: v1alpha1.RateLimitService{
			Spec: v1alpha1.RateLimitServiceSpec{
				Backend: &v1alpha1.RateLimitServiceSpec_Backend{
					Redis: &v1alpha1.RateLimitServiceSpec_Backend_Redis{
						Type: "single",
						URL:  "redis:6379",
					},
					CacheKeyPrefix:                     "my-svc",
					StopCacheKeyIncrementWhenOverlimit: true,
				},
			},
		},
	}

	env, err := builder.BuildBackendEnv()
	assert.NoError(t, err)
	assert.Equal(t, "my-svc", env["CACHE_KEY_PREFIX"])
	assert.Equal(t, "true", env["STOP_CACHE_KEY_INCREMENT_WHEN_OVERLIMIT"])
}

func TestEnvBuilder_Build(t *testing.T) {
	type fields struct {
		RateLimitService v1alpha1.RateLimitService
	}
	type expectations struct {
		wantError bool
		config    corev1.ConfigMap
	}
	tests := []struct {
		name         string
		fields       fields
		expectations expectations
	}{
		{
			name: "generate the simplest configmap",
			fields: fields{
				RateLimitService: v1alpha1.RateLimitService{
					ObjectMeta: v1.ObjectMeta{
						Name:      "foo",
						Namespace: "bar",
					},
					Spec: v1alpha1.RateLimitServiceSpec{
						Backend: &v1alpha1.RateLimitServiceSpec_Backend{
							Redis: &v1alpha1.RateLimitServiceSpec_Backend_Redis{
								Type: "single",
								URL:  "127.0.0.1:6379",
							},
						},
					},
				},
			},
			expectations: expectations{
				wantError: false,
				config: corev1.ConfigMap{
					ObjectMeta: v1.ObjectMeta{
						Name:      "foo-config-env",
						Namespace: "bar",
						Labels: map[string]string{
							"app.kubernetes.io/created-by": "foo",
							"app.kubernetes.io/managed-by": "istio-rateltimit-operator",
							"app.kubernetes.io/name":       "foo-config-env",
						},
					},
					Data: map[string]string{
						"REDIS_SOCKET_TYPE": "tcp",
						"REDIS_TYPE":        "single",
						"REDIS_URL":         "127.0.0.1:6379",
						"USE_STATSD":        "false",
					},
				},
			},
		},
		{
			name: "custom environment variable",
			fields: fields{
				RateLimitService: v1alpha1.RateLimitService{
					ObjectMeta: v1.ObjectMeta{
						Name:      "foo",
						Namespace: "bar",
					},
					Spec: v1alpha1.RateLimitServiceSpec{
						Backend: &v1alpha1.RateLimitServiceSpec_Backend{
							Redis: &v1alpha1.RateLimitServiceSpec_Backend_Redis{
								Type: "single",
								URL:  "127.0.0.1:6379",
							},
						},
						Environment: &map[string]string{
							"CACHE_KEY_PREFIX": "foo",
						},
					},
				},
			},
			expectations: expectations{
				wantError: false,
				config: corev1.ConfigMap{
					ObjectMeta: v1.ObjectMeta{
						Name:      "foo-config-env",
						Namespace: "bar",
						Labels: map[string]string{
							"app.kubernetes.io/created-by": "foo",
							"app.kubernetes.io/managed-by": "istio-rateltimit-operator",
							"app.kubernetes.io/name":       "foo-config-env",
						},
					},
					Data: map[string]string{
						"REDIS_SOCKET_TYPE": "tcp",
						"REDIS_TYPE":        "single",
						"REDIS_URL":         "127.0.0.1:6379",
						"USE_STATSD":        "false",
						"CACHE_KEY_PREFIX":  "foo",
					},
				},
			},
		},
		{
			name: "custom environment variable should replace default environment variable",
			fields: fields{
				RateLimitService: v1alpha1.RateLimitService{
					ObjectMeta: v1.ObjectMeta{
						Name:      "foo",
						Namespace: "bar",
					},
					Spec: v1alpha1.RateLimitServiceSpec{
						Backend: &v1alpha1.RateLimitServiceSpec_Backend{
							Redis: &v1alpha1.RateLimitServiceSpec_Backend_Redis{
								Type: "single",
								URL:  "127.0.0.1:6379",
							},
						},
						Environment: &map[string]string{
							"USE_STATSD": "true",
						},
					},
				},
			},
			expectations: expectations{
				wantError: false,
				config: corev1.ConfigMap{
					ObjectMeta: v1.ObjectMeta{
						Name:      "foo-config-env",
						Namespace: "bar",
						Labels: map[string]string{
							"app.kubernetes.io/created-by": "foo",
							"app.kubernetes.io/managed-by": "istio-rateltimit-operator",
							"app.kubernetes.io/name":       "foo-config-env",
						},
					},
					Data: map[string]string{
						"REDIS_SOCKET_TYPE": "tcp",
						"REDIS_TYPE":        "single",
						"REDIS_URL":         "127.0.0.1:6379",
						"USE_STATSD":        "true",
					},
				},
			},
		},
		{
			name: "defined API should replace custom environment variable",
			fields: fields{
				RateLimitService: v1alpha1.RateLimitService{
					ObjectMeta: v1.ObjectMeta{
						Name:      "foo",
						Namespace: "bar",
					},
					Spec: v1alpha1.RateLimitServiceSpec{
						Backend: &v1alpha1.RateLimitServiceSpec_Backend{
							Redis: &v1alpha1.RateLimitServiceSpec_Backend_Redis{
								Type: "single",
								URL:  "127.0.0.1:6379",
							},
						},
						Environment: &map[string]string{
							"REDIS_TYPE": "sentinel",
						},
					},
				},
			},
			expectations: expectations{
				wantError: false,
				config: corev1.ConfigMap{
					ObjectMeta: v1.ObjectMeta{
						Name:      "foo-config-env",
						Namespace: "bar",
						Labels: map[string]string{
							"app.kubernetes.io/created-by": "foo",
							"app.kubernetes.io/managed-by": "istio-rateltimit-operator",
							"app.kubernetes.io/name":       "foo-config-env",
						},
					},
					Data: map[string]string{
						"REDIS_SOCKET_TYPE": "tcp",
						"REDIS_TYPE":        "single",
						"REDIS_URL":         "127.0.0.1:6379",
						"USE_STATSD":        "false",
					},
				},
			},
		},
		{
			name: "full environment",
			fields: fields{
				RateLimitService: v1alpha1.RateLimitService{
					ObjectMeta: v1.ObjectMeta{
						Name:      "foo",
						Namespace: "bar",
					},
					Spec: v1alpha1.RateLimitServiceSpec{
						Backend: &v1alpha1.RateLimitServiceSpec_Backend{
							Redis: &v1alpha1.RateLimitServiceSpec_Backend_Redis{
								Type: "single",
								URL:  "127.0.0.1:6379",
							},
						},
						Monitoring: &v1alpha1.RateLimitServiceSpec_Monitoring{
							Enabled: true,
						},
					},
				},
			},
			expectations: expectations{
				wantError: false,
				config: corev1.ConfigMap{
					ObjectMeta: v1.ObjectMeta{
						Name:      "foo-config-env",
						Namespace: "bar",
						Labels: map[string]string{
							"app.kubernetes.io/created-by": "foo",
							"app.kubernetes.io/managed-by": "istio-rateltimit-operator",
							"app.kubernetes.io/name":       "foo-config-env",
						},
					},
					Data: map[string]string{
						"REDIS_SOCKET_TYPE": "tcp",
						"REDIS_TYPE":        "single",
						"REDIS_URL":         "127.0.0.1:6379",
						"USE_STATSD":        "true",
						"STATSD_HOST":       "localhost",
						"STATSD_PORT":       "9125",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &service.EnvBuilder{
				RateLimitService: tt.fields.RateLimitService,
			}

			got, err := n.Build()
			if tt.expectations.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if !reflect.DeepEqual(*got, tt.expectations.config) {
					t.Errorf("EnvBuilder.BuildLabels() = %v, want %v", *got, tt.expectations.config)
				}
			}
		})
	}
}

func TestEnvBuilder_BuildRedisEnv_WithPoolBehavior(t *testing.T) {
	builder := &service.EnvBuilder{
		RateLimitService: v1alpha1.RateLimitService{
			Spec: v1alpha1.RateLimitServiceSpec{
				Backend: &v1alpha1.RateLimitServiceSpec_Backend{
					Redis: &v1alpha1.RateLimitServiceSpec_Backend_Redis{
						Type: "single",
						URL:  "redis:6379",
						Pool: &v1alpha1.RateLimitServiceSpec_Backend_Redis_Pool{
							Size:                10,
							OnEmptyBehavior:     "wait",
							OnEmptyWaitDuration: "1s",
						},
					},
				},
			},
		},
	}

	env, err := builder.BuildRedisEnv()
	assert.NoError(t, err)
	assert.Equal(t, "10", env["REDIS_POOL_SIZE"])
	assert.Equal(t, "wait", env["REDIS_POOL_ON_EMPTY_BEHAVIOR"])
	assert.Equal(t, "1s", env["REDIS_POOL_ON_EMPTY_WAIT_DURATION"])
}
