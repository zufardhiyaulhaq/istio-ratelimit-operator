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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// LocalRateLimitConfigSpec defines the desired state of LocalRateLimitConfig
type LocalRateLimitConfigSpec struct {
	// +kubebuilder:validation:Enum=gateway;sidecar
	Type ConfigContext `json:"type"`

	Selector LocalRateLimitConfigSelector `json:"selector"`
}

type LocalRateLimitConfigSelector struct {
	Labels       map[string]string `json:"labels"`
	IstioVersion []string          `json:"istio_version"`
	SNI          *string           `json:"sni,omitempty"`
}

// LocalRateLimitConfigStatus defines the observed state of LocalRateLimitConfig
type LocalRateLimitConfigStatus struct{}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// LocalRateLimitConfig is the Schema for the localratelimitconfigs API
type LocalRateLimitConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LocalRateLimitConfigSpec   `json:"spec,omitempty"`
	Status LocalRateLimitConfigStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// LocalRateLimitConfigList contains a list of LocalRateLimitConfig
type LocalRateLimitConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LocalRateLimitConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LocalRateLimitConfig{}, &LocalRateLimitConfigList{})
}
