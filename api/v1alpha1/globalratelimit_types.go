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

// GlobalRateLimitSpec defines the desired state of GlobalRateLimit
type GlobalRateLimitSpec struct {
	Config     string                    `json:"config"`
	Selector   GlobalRateLimitSelector   `json:"selector"`
	Matcher    []*GlobalRateLimit_Action `json:"matcher"`
	ShadowMode bool                      `json:"shadow_mode,omitempty"`
	Limit      *GlobalRateLimit_Limit    `json:"limit,omitempty"`
}

type GlobalRateLimitSelector struct {
	VHost string  `json:"vhost"`
	Route *string `json:"route,omitempty"`
}

type GlobalRateLimit_Limit struct {
	Unit            string `json:"unit,omitempty" yaml:"unit,omitempty"`
	RequestsPerUnit int    `json:"requests_per_unit,omitempty" yaml:"requests_per_unit,omitempty"`
}

type GlobalRateLimit_Action struct {
	SourceCluster      *GlobalRateLimit_Action_SourceCluster      `json:"source_cluster,omitempty" yaml:"source_cluster,omitempty"`
	DestinationCluster *GlobalRateLimit_Action_DestinationCluster `json:"destination_cluster,omitempty" yaml:"destination_cluster,omitempty"`
	RequestHeaders     *GlobalRateLimit_Action_RequestHeaders     `json:"request_headers,omitempty" yaml:"request_headers,omitempty"`
	RemoteAddress      *GlobalRateLimit_Action_RemoteAddress      `json:"remote_address,omitempty" yaml:"remote_address,omitempty"`
	GenericKey         *GlobalRateLimit_Action_GenericKey         `json:"generic_key,omitempty" yaml:"generic_key,omitempty"`
	HeaderValueMatch   *GlobalRateLimit_Action_HeaderValueMatch   `json:"header_value_match,omitempty" yaml:"header_value_match,omitempty"`
}

type GlobalRateLimit_Action_SourceCluster struct{}
type GlobalRateLimit_Action_DestinationCluster struct{}
type GlobalRateLimit_Action_RemoteAddress struct{}

type GlobalRateLimit_Action_RequestHeaders struct {
	// The header name to be queried from the request headers. The header’s
	// value is used to populate the value of the descriptor entry for the
	// descriptor_key.
	HeaderName string `json:"header_name,omitempty" yaml:"header_name,omitempty"`
	// The key to use in the descriptor entry.
	DescriptorKey string `json:"descriptor_key,omitempty" yaml:"descriptor_key,omitempty"`
	// If set to true, Envoy skips the descriptor while calling rate limiting service
	// when header is not present in the request. By default it skips calling the
	// rate limiting service if this header is not present in the request.
	SkipIfAbsent bool `json:"skip_if_absent,omitempty" yaml:"skip_if_absent,omitempty"`
}

type GlobalRateLimit_Action_GenericKey struct {
	// The value to use in the descriptor entry.
	DescriptorValue string `json:"descriptor_value,omitempty" yaml:"descriptor_value,omitempty"`
	// An optional key to use in the descriptor entry. If not set it defaults
	// to 'generic_key' as the descriptor key.
	DescriptorKey *string `json:"descriptor_key,omitempty" yaml:"descriptor_key,omitempty"`
}

type GlobalRateLimit_Action_HeaderValueMatch struct {
	// The value to use in the descriptor entry.
	DescriptorValue string `json:"descriptor_value,omitempty" yaml:"descriptor_value,omitempty"`
	// If set to true, the action will append a descriptor entry when the
	// request matches the headers. If set to false, the action will append a
	// descriptor entry when the request does not match the headers. The
	// default value is true.
	ExpectMatch *bool `json:"expect_match,omitempty" yaml:"expect_match,omitempty"`
	// Specifies a set of headers that the rate limit action should match
	// on. The action will check the request’s headers against all the
	// specified headers in the config. A match will happen if all the
	// headers in the config are present in the request with the same values
	// (or based on presence if the value field is not in the config).
	Headers []*GlobalRateLimit_Action_HeaderValueMatch_HeaderMatcher `json:"headers,omitempty" yaml:"headers,omitempty"`
}

type GlobalRateLimit_Action_HeaderValueMatch_HeaderMatcher struct {
	// Specifies the name of the header in the request.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// Specifies how the header match will be performed to route the request.
	ExactMatch     string                                                              `json:"exact_match,omitempty" yaml:"exact_match,omitempty"`
	RegexMatch     string                                                              `json:"regex_match,omitempty" yaml:"regex_match,omitempty"`
	SafeRegexMatch *GlobalRateLimit_Action_HeaderValueMatch_HeaderMatcher_RegexMatcher `json:"safe_regex_match,omitempty" yaml:"safe_regex_match,omitempty"`
	RangeMatch     *GlobalRateLimit_Action_HeaderValueMatch_HeaderMatcher_Int64Range   `json:"range_match,omitempty" yaml:"range_match,omitempty"`
	PresentMatch   bool                                                                `json:"present_match,omitempty" yaml:"present_match,omitempty"`
	PrefixMatch    string                                                              `json:"prefix_match,omitempty" yaml:"prefix_match,omitempty"`
	SuffixMatch    string                                                              `json:"suffix_match,omitempty" yaml:"suffix_match,omitempty"`
	ContainsMatch  string                                                              `json:"contains_match,omitempty" yaml:"contains_match,omitempty"`
	// If specified, the match result will be inverted before checking. Defaults to false.
	//
	// Examples:
	//
	// * The regex ``\d{3}`` does not match the value *1234*, so it will match when inverted.
	// * The range [-10,0) will match the value -1, so it will not match when inverted.
	InvertMatch bool `json:"invert_match,omitempty" yaml:"invert_match,omitempty"`
}

type GlobalRateLimit_Action_HeaderValueMatch_HeaderMatcher_RegexMatcher struct {
	Regex string `json:"regex,omitempty" yaml:"regex,omitempty"`
}

type GlobalRateLimit_Action_HeaderValueMatch_HeaderMatcher_Int64Range struct {
	// start of the range (inclusive)
	Start int64 `json:"start,omitempty" yaml:"start,omitempty"`
	// end of the range (exclusive)
	End int64 `json:"end,omitempty" yaml:"end,omitempty"`
}

// GlobalRateLimitStatus defines the observed state of GlobalRateLimit
type GlobalRateLimitStatus struct{}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// GlobalRateLimit is the Schema for the globalratelimits API
type GlobalRateLimit struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GlobalRateLimitSpec   `json:"spec,omitempty"`
	Status GlobalRateLimitStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// GlobalRateLimitList contains a list of GlobalRateLimit
type GlobalRateLimitList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GlobalRateLimit `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GlobalRateLimit{}, &GlobalRateLimitList{})
}
