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
	Kubernetes *RateLimitServiceSpec_Kubernetes `json:"kubernetes,omitempty"`
	Backend    *RateLimitServiceSpec_Backend    `json:"backend,omitempty"`
}

type RateLimitServiceSpec_Kubernetes struct {
	ReplicaCount *int32                                       `json:"replica_count,omitempty"`
	Resources    *corev1.ResourceRequirements                 `json:"resources,omitempty"`
	AutoScaling  *RateLimitServiceSpec_Kubernetes_AutoScaling `json:"auto_scaling,omitempty"`
}

type RateLimitServiceSpec_Kubernetes_AutoScaling struct {
	MaxReplica *int32 `json:"max_replicas,omitempty"`
	MinReplica *int32 `json:"min_replicas,omitempty"`
}

type RateLimitServiceSpec_Backend struct {
	Redis *RateLimitService_RedisBackend `json:"redis,omitempty"`
}

type RateLimitService_RedisBackend struct {
	Type string `json:"type,omitempty"`
	URL  string `json:"url,omitempty"`
}

type RateLimitServiceStatus struct {
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
