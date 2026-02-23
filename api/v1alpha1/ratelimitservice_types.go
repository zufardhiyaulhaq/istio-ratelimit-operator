/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RateLimitServiceSpec defines the desired state of RateLimitService
type RateLimitServiceSpec struct {
	Kubernetes      *RateLimitServiceSpec_Kubernetes      `json:"kubernetes,omitempty"`
	Backend         *RateLimitServiceSpec_Backend         `json:"backend,omitempty"`
	Server          *RateLimitServiceSpec_Server          `json:"server,omitempty"`
	Monitoring      *RateLimitServiceSpec_Monitoring      `json:"monitoring,omitempty"`
	ResponseHeaders *RateLimitServiceSpec_ResponseHeaders `json:"responseHeaders,omitempty"`
	Logging         *RateLimitServiceSpec_Logging         `json:"logging,omitempty"`
	ShadowMode      bool                                  `json:"shadowMode,omitempty"`
	Environment     *map[string]string                    `json:"environment,omitempty"`
}

type RateLimitServiceSpec_Kubernetes struct {
	ReplicaCount             *int32                                       `json:"replica_count,omitempty"`
	Image                    *string                                      `json:"image,omitempty"`
	Resources                *corev1.ResourceRequirements                 `json:"resources,omitempty"`
	AutoScaling              *RateLimitServiceSpec_Kubernetes_AutoScaling `json:"auto_scaling,omitempty"`
	ExtraLabels              *map[string]string                           `json:"extra_labels,omitempty"`
	Annotations              *map[string]string                           `json:"annotations,omitempty"`
	SecurityContext          *corev1.PodSecurityContext                   `json:"securityContext,omitempty"`
	ContainerSecurityContext *corev1.SecurityContext                      `json:"containerSecurityContext,omitempty"`
	NodeSelector             map[string]string                            `json:"nodeSelector,omitempty"`
	Tolerations              []corev1.Toleration                          `json:"tolerations,omitempty"`
	Affinity                 *corev1.Affinity                             `json:"affinity,omitempty"`
	ImagePullSecrets         []corev1.LocalObjectReference                `json:"imagePullSecrets,omitempty"`
	LivenessProbe            *corev1.Probe                                `json:"livenessProbe,omitempty"`
	PodDisruptionBudget      *RateLimitServiceSpec_Kubernetes_PDB         `json:"podDisruptionBudget,omitempty"`
}

type RateLimitServiceSpec_Kubernetes_AutoScaling struct {
	MaxReplica *int32 `json:"max_replicas,omitempty"`
	MinReplica *int32 `json:"min_replicas,omitempty"`
}

type RateLimitServiceSpec_Kubernetes_PDB struct {
	MinAvailable   *int32 `json:"minAvailable,omitempty"`
	MaxUnavailable *int32 `json:"maxUnavailable,omitempty"`
}

type RateLimitServiceSpec_Backend struct {
	Redis                              *RateLimitServiceSpec_Backend_Redis `json:"redis,omitempty"`
	CacheKeyPrefix                     string                              `json:"cacheKeyPrefix,omitempty"`
	StopCacheKeyIncrementWhenOverlimit bool                                `json:"stopCacheKeyIncrementWhenOverlimit,omitempty"`
}

type RateLimitServiceSpec_Backend_Redis struct {
	Type                        string                                           `json:"type,omitempty"`
	URL                         string                                           `json:"url,omitempty"`
	Auth                        string                                           `json:"auth,omitempty"`
	AuthSecretRef               *SecretKeyRef                                    `json:"authSecretRef,omitempty"`
	Config                      *RateLimitServiceSpec_Backend_Redis_Config       `json:"config,omitempty"`
	TLS                         *RateLimitServiceSpec_Backend_Redis_TLS          `json:"tls,omitempty"`
	SentinelAuth                *RateLimitServiceSpec_Backend_Redis_SentinelAuth `json:"sentinelAuth,omitempty"`
	Pool                        *RateLimitServiceSpec_Backend_Redis_Pool         `json:"pool,omitempty"`
	Timeout                     *string                                          `json:"timeout,omitempty"`
	HealthCheckActiveConnection bool                                             `json:"healthCheckActiveConnection,omitempty"`
	PerSecond                   *RateLimitServiceSpec_Backend_Redis_PerSecond    `json:"perSecond,omitempty"`
}

type RateLimitServiceSpec_Backend_Redis_Config struct {
	PipelineWindow *string `json:"pipeline_window,omitempty"`
	PipelineLimit  *int    `json:"pipeline_limit,omitempty"`
}

// SecretKeyRef references a key in a Secret
type SecretKeyRef struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

type RateLimitServiceSpec_Backend_Redis_TLS struct {
	Enabled                  bool   `json:"enabled,omitempty"`
	SecretRef                string `json:"secretRef,omitempty"`
	SkipHostnameVerification bool   `json:"skipHostnameVerification,omitempty"`
}

type RateLimitServiceSpec_Backend_Redis_SentinelAuth struct {
	SecretRef *SecretKeyRef `json:"secretRef,omitempty"`
}

type RateLimitServiceSpec_Backend_Redis_Pool struct {
	Size                int    `json:"size,omitempty"`
	OnEmptyBehavior     string `json:"onEmptyBehavior,omitempty"`
	OnEmptyWaitDuration string `json:"onEmptyWaitDuration,omitempty"`
}

type RateLimitServiceSpec_Backend_Redis_PerSecond struct {
	Enabled       bool                                     `json:"enabled,omitempty"`
	URL           string                                   `json:"url,omitempty"`
	AuthSecretRef *SecretKeyRef                            `json:"authSecretRef,omitempty"`
	TLS           *RateLimitServiceSpec_Backend_Redis_TLS  `json:"tls,omitempty"`
	Pool          *RateLimitServiceSpec_Backend_Redis_Pool `json:"pool,omitempty"`
}

// Server configuration for gRPC and debug ports
type RateLimitServiceSpec_Server struct {
	GRPC  *RateLimitServiceSpec_Server_GRPC  `json:"grpc,omitempty"`
	Debug *RateLimitServiceSpec_Server_Debug `json:"debug,omitempty"`
}

type RateLimitServiceSpec_Server_GRPC struct {
	Port                  *int32                                      `json:"port,omitempty"`
	TLS                   *RateLimitServiceSpec_Server_GRPC_TLS       `json:"tls,omitempty"`
	ClientTLS             *RateLimitServiceSpec_Server_GRPC_ClientTLS `json:"clientTls,omitempty"`
	MaxConnectionAge      *string                                     `json:"maxConnectionAge,omitempty"`
	MaxConnectionAgeGrace *string                                     `json:"maxConnectionAgeGrace,omitempty"`
}

type RateLimitServiceSpec_Server_GRPC_TLS struct {
	Enabled   bool   `json:"enabled,omitempty"`
	SecretRef string `json:"secretRef,omitempty"`
}

type RateLimitServiceSpec_Server_GRPC_ClientTLS struct {
	CACertSecretRef string `json:"caCertSecretRef,omitempty"`
	SAN             string `json:"san,omitempty"`
}

type RateLimitServiceSpec_Server_Debug struct {
	Port *int32 `json:"port,omitempty"`
}

type RateLimitServiceSpec_Monitoring struct {
	// +optional
	Enabled bool `json:"enabled,omitempty"`

	// +kubebuilder:validation:Enum=statsd;prometheus;dogstatsd
	// +optional
	Type string `json:"type,omitempty"`

	Prometheus *RateLimitServiceSpec_Monitoring_Prometheus `json:"prometheus,omitempty"`
	Tracing    *RateLimitServiceSpec_Monitoring_Tracing    `json:"tracing,omitempty"`

	NearLimitRatio     *string `json:"nearLimitRatio,omitempty"`
	StatsFlushInterval *string `json:"statsFlushInterval,omitempty"`

	// This API is deprecated
	Statsd *RateLimitServiceSpec_Monitoring_Statsd `json:"statsd,omitempty"`
}

type RateLimitServiceSpec_Monitoring_Prometheus struct {
	Addr string `json:"addr,omitempty"`
	Path string `json:"path,omitempty"`
}

type RateLimitServiceSpec_Monitoring_Tracing struct {
	Enabled          bool   `json:"enabled,omitempty"`
	ExporterProtocol string `json:"exporterProtocol,omitempty"`
	ServiceName      string `json:"serviceName,omitempty"`
	ServiceNamespace string `json:"serviceNamespace,omitempty"`
	SamplingRate     string `json:"samplingRate,omitempty"`
}

// This API is deprecated
type RateLimitServiceSpec_Monitoring_Statsd struct {
	Enabled bool                                        `json:"enabled,omitempty"`
	Spec    RateLimitServiceSpec_Monitoring_Statsd_Spec `json:"spec,omitempty"`
}

type RateLimitServiceSpec_Monitoring_Statsd_Spec struct {
	Host string `json:"host,omitempty"`
	Port int    `json:"port,omitempty"`
}

type RateLimitServiceSpec_ResponseHeaders struct {
	Enabled bool `json:"enabled,omitempty"`
}

type RateLimitServiceSpec_Logging struct {
	// +kubebuilder:validation:Enum=debug;info;warning;error
	Level string `json:"level,omitempty"`
	// +kubebuilder:validation:Enum=json;text
	Format string `json:"format,omitempty"`
}

type RateLimitServiceStatus struct {
	TriggerStatsdExporterReload bool `json:"trigger_statsd_exporter_reload,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// RateLimitService is the Schema for the ratelimitservices API
type RateLimitService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RateLimitServiceSpec   `json:"spec,omitempty"`
	Status RateLimitServiceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RateLimitServiceList contains a list of RateLimitService
type RateLimitServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RateLimitService `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RateLimitService{}, &RateLimitServiceList{})
}
