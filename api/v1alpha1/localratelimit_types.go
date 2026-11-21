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
	"strings"

	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/local/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// LocalRateLimitSpec defines the desired state of LocalRateLimit
type LocalRateLimitSpec struct {
	Config   string                 `json:"config"`
	Selector LocalRateLimitSelector `json:"selector"`
	Limit    *LocalRateLimit_Limit  `json:"limit,omitempty"`
}

type LocalRateLimitSelector struct {
	VHost string  `json:"vhost"`
	Route *string `json:"route,omitempty"`
}

type LocalRateLimit_Limit struct {
	// Unit is the token bucket fill interval: second, minute, hour, day, week
	// (7 days), month (30 days) or year (365 days). Case-insensitive.
	// +kubebuilder:validation:Pattern=`^(?i)(second|minute|hour|day|week|month|year)$`
	Unit            string `json:"unit,omitempty" yaml:"unit,omitempty"`
	RequestsPerUnit int    `json:"requests_per_unit,omitempty" yaml:"requests_per_unit,omitempty"`
}

// ToTokenBucket builds an Envoy token bucket that allows RequestsPerUnit
// requests per unit: it holds RequestsPerUnit tokens and refills all of them
// at the end of each unit.
func (l LocalRateLimit_Limit) ToTokenBucket() *types.LocalRateLimit_TokenBucket {
	var interval string

	switch strings.ToLower(l.Unit) {
	case "second":
		interval = "1s"
	case "minute":
		interval = "60s"
	case "hour":
		interval = "3600s"
	case "day":
		interval = "86400s"
	case "week":
		interval = "604800s"
	case "month":
		interval = "2592000s"
	case "year":
		interval = "31536000s"
	default:
		interval = "1s"
	}

	return &types.LocalRateLimit_TokenBucket{
		MaxTokens:     l.RequestsPerUnit,
		TokensPerFill: l.RequestsPerUnit,
		FillInterval:  interval,
	}
}

// LocalRateLimitStatus defines the observed state of LocalRateLimit
type LocalRateLimitStatus struct{}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// LocalRateLimit is the Schema for the localratelimits API
type LocalRateLimit struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LocalRateLimitSpec   `json:"spec,omitempty"`
	Status LocalRateLimitStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// LocalRateLimitList contains a list of LocalRateLimit
type LocalRateLimitList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LocalRateLimit `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LocalRateLimit{}, &LocalRateLimitList{})
}
