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

func TestEnvBuilder_BuildStatsdEnv(t *testing.T) {
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
			name: "statsd not enabled",
			fields: fields{
				RateLimitService: v1alpha1.RateLimitService{
					Spec: v1alpha1.RateLimitServiceSpec{
						Monitoring: &v1alpha1.RateLimitServiceSpec_Monitoring{
							Statsd: &v1alpha1.RateLimitServiceSpec_Monitoring_Statsd{
								Enabled: false,
							},
						},
					},
				},
			},
			want:    make(map[string]string),
			wantErr: false,
		},
		{
			name: "statsd enabled",
			fields: fields{
				RateLimitService: v1alpha1.RateLimitService{
					Spec: v1alpha1.RateLimitServiceSpec{
						Monitoring: &v1alpha1.RateLimitServiceSpec_Monitoring{
							Statsd: &v1alpha1.RateLimitServiceSpec_Monitoring_Statsd{
								Enabled: true,
								Spec: v1alpha1.RateLimitServiceSpec_Monitoring_Statsd_Spec{
									Host: "foo",
									Port: 8125,
								},
							},
						},
					},
				},
			},
			want: map[string]string{
				"USE_STATSD":  "true",
				"STATSD_HOST": "foo",
				"STATSD_PORT": "8125",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &service.EnvBuilder{
				RateLimitService: tt.fields.RateLimitService,
			}
			got, err := n.BuildStatsdEnv()
			if (err != nil) != tt.wantErr {
				t.Errorf("EnvBuilder.buildStatsdEnv() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EnvBuilder.buildStatsdEnv() = %v, want %v", got, tt.want)
			}
		})
	}
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
							Statsd: &v1alpha1.RateLimitServiceSpec_Monitoring_Statsd{
								Enabled: true,
								Spec: v1alpha1.RateLimitServiceSpec_Monitoring_Statsd_Spec{
									Host: "statsd",
									Port: 8125,
								},
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
						"USE_STATSD":        "true",
						"STATSD_HOST":       "statsd",
						"STATSD_PORT":       "8125",
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
