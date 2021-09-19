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

// GlobalRateLimitConfigSpec defines the desired state of GlobalRateLimitConfig
type GlobalRateLimitConfigSpec struct {
	// +kubebuilder:validation:Enum=gateway
	Type string `json:"type"`

	Selector  GlobalRateLimitConfigSelector  `json:"selector"`
	Ratelimit GlobalRateLimitConfigRatelimit `json:"ratelimit"`
}

type GlobalRateLimitConfigSelector struct {
	Labels       map[string]string `json:"labels"`
	IstioVersion []string          `json:"istio_version"`
	SNI          *string           `json:"sni,omitempty"`
}

type GlobalRateLimitConfigRatelimit struct {
	Spec GlobalRateLimitConfigRatelimitSpec `json:"spec"`
}

type GlobalRateLimitConfigRatelimitSpec struct {
	Domain          string                                    `json:"domain"`
	FailureModeDeny bool                                      `json:"failure_mode_deny"`
	Timeout         string                                    `json:"timeout"`
	Service         GlobalRateLimitConfigRatelimitSpecService `json:"service"`
}

type GlobalRateLimitConfigRatelimitSpecService struct {
	Type    string `json:"type"`
	Name    string `json:"name,omitempty"`
	Address string `json:"address,omitempty"`
	Port    int    `json:"port,omitempty"`
}

// GlobalRateLimitConfigStatus defines the observed state of GlobalRateLimitConfig
type GlobalRateLimitConfigStatus struct{}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// GlobalRateLimitConfig is the Schema for the globalratelimitconfigs API
type GlobalRateLimitConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GlobalRateLimitConfigSpec   `json:"spec,omitempty"`
	Status GlobalRateLimitConfigStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// GlobalRateLimitConfigList contains a list of GlobalRateLimitConfig
type GlobalRateLimitConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GlobalRateLimitConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GlobalRateLimitConfig{}, &GlobalRateLimitConfigList{})
}
