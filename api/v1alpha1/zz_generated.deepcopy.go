//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/api/core/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GlobalRateLimit) DeepCopyInto(out *GlobalRateLimit) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GlobalRateLimit.
func (in *GlobalRateLimit) DeepCopy() *GlobalRateLimit {
	if in == nil {
		return nil
	}
	out := new(GlobalRateLimit)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *GlobalRateLimit) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GlobalRateLimitConfig) DeepCopyInto(out *GlobalRateLimitConfig) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GlobalRateLimitConfig.
func (in *GlobalRateLimitConfig) DeepCopy() *GlobalRateLimitConfig {
	if in == nil {
		return nil
	}
	out := new(GlobalRateLimitConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *GlobalRateLimitConfig) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GlobalRateLimitConfigList) DeepCopyInto(out *GlobalRateLimitConfigList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]GlobalRateLimitConfig, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GlobalRateLimitConfigList.
func (in *GlobalRateLimitConfigList) DeepCopy() *GlobalRateLimitConfigList {
	if in == nil {
		return nil
	}
	out := new(GlobalRateLimitConfigList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *GlobalRateLimitConfigList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GlobalRateLimitConfigRatelimit) DeepCopyInto(out *GlobalRateLimitConfigRatelimit) {
	*out = *in
	out.Spec = in.Spec
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GlobalRateLimitConfigRatelimit.
func (in *GlobalRateLimitConfigRatelimit) DeepCopy() *GlobalRateLimitConfigRatelimit {
	if in == nil {
		return nil
	}
	out := new(GlobalRateLimitConfigRatelimit)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GlobalRateLimitConfigRatelimitSpec) DeepCopyInto(out *GlobalRateLimitConfigRatelimitSpec) {
	*out = *in
	out.Service = in.Service
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GlobalRateLimitConfigRatelimitSpec.
func (in *GlobalRateLimitConfigRatelimitSpec) DeepCopy() *GlobalRateLimitConfigRatelimitSpec {
	if in == nil {
		return nil
	}
	out := new(GlobalRateLimitConfigRatelimitSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GlobalRateLimitConfigRatelimitSpecService) DeepCopyInto(out *GlobalRateLimitConfigRatelimitSpecService) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GlobalRateLimitConfigRatelimitSpecService.
func (in *GlobalRateLimitConfigRatelimitSpecService) DeepCopy() *GlobalRateLimitConfigRatelimitSpecService {
	if in == nil {
		return nil
	}
	out := new(GlobalRateLimitConfigRatelimitSpecService)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GlobalRateLimitConfigSelector) DeepCopyInto(out *GlobalRateLimitConfigSelector) {
	*out = *in
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.IstioVersion != nil {
		in, out := &in.IstioVersion, &out.IstioVersion
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.SNI != nil {
		in, out := &in.SNI, &out.SNI
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GlobalRateLimitConfigSelector.
func (in *GlobalRateLimitConfigSelector) DeepCopy() *GlobalRateLimitConfigSelector {
	if in == nil {
		return nil
	}
	out := new(GlobalRateLimitConfigSelector)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GlobalRateLimitConfigSpec) DeepCopyInto(out *GlobalRateLimitConfigSpec) {
	*out = *in
	in.Selector.DeepCopyInto(&out.Selector)
	out.Ratelimit = in.Ratelimit
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GlobalRateLimitConfigSpec.
func (in *GlobalRateLimitConfigSpec) DeepCopy() *GlobalRateLimitConfigSpec {
	if in == nil {
		return nil
	}
	out := new(GlobalRateLimitConfigSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GlobalRateLimitConfigStatus) DeepCopyInto(out *GlobalRateLimitConfigStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GlobalRateLimitConfigStatus.
func (in *GlobalRateLimitConfigStatus) DeepCopy() *GlobalRateLimitConfigStatus {
	if in == nil {
		return nil
	}
	out := new(GlobalRateLimitConfigStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GlobalRateLimitList) DeepCopyInto(out *GlobalRateLimitList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]GlobalRateLimit, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GlobalRateLimitList.
func (in *GlobalRateLimitList) DeepCopy() *GlobalRateLimitList {
	if in == nil {
		return nil
	}
	out := new(GlobalRateLimitList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *GlobalRateLimitList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GlobalRateLimitSelector) DeepCopyInto(out *GlobalRateLimitSelector) {
	*out = *in
	if in.Route != nil {
		in, out := &in.Route, &out.Route
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GlobalRateLimitSelector.
func (in *GlobalRateLimitSelector) DeepCopy() *GlobalRateLimitSelector {
	if in == nil {
		return nil
	}
	out := new(GlobalRateLimitSelector)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GlobalRateLimitSpec) DeepCopyInto(out *GlobalRateLimitSpec) {
	*out = *in
	in.Selector.DeepCopyInto(&out.Selector)
	if in.Matcher != nil {
		in, out := &in.Matcher, &out.Matcher
		*out = make([]*GlobalRateLimit_Action, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(GlobalRateLimit_Action)
				(*in).DeepCopyInto(*out)
			}
		}
	}
	if in.Limit != nil {
		in, out := &in.Limit, &out.Limit
		*out = new(GlobalRateLimit_Limit)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GlobalRateLimitSpec.
func (in *GlobalRateLimitSpec) DeepCopy() *GlobalRateLimitSpec {
	if in == nil {
		return nil
	}
	out := new(GlobalRateLimitSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GlobalRateLimitStatus) DeepCopyInto(out *GlobalRateLimitStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GlobalRateLimitStatus.
func (in *GlobalRateLimitStatus) DeepCopy() *GlobalRateLimitStatus {
	if in == nil {
		return nil
	}
	out := new(GlobalRateLimitStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GlobalRateLimit_Action) DeepCopyInto(out *GlobalRateLimit_Action) {
	*out = *in
	if in.SourceCluster != nil {
		in, out := &in.SourceCluster, &out.SourceCluster
		*out = new(GlobalRateLimit_Action_SourceCluster)
		**out = **in
	}
	if in.DestinationCluster != nil {
		in, out := &in.DestinationCluster, &out.DestinationCluster
		*out = new(GlobalRateLimit_Action_DestinationCluster)
		**out = **in
	}
	if in.RequestHeaders != nil {
		in, out := &in.RequestHeaders, &out.RequestHeaders
		*out = new(GlobalRateLimit_Action_RequestHeaders)
		**out = **in
	}
	if in.RemoteAddress != nil {
		in, out := &in.RemoteAddress, &out.RemoteAddress
		*out = new(GlobalRateLimit_Action_RemoteAddress)
		**out = **in
	}
	if in.GenericKey != nil {
		in, out := &in.GenericKey, &out.GenericKey
		*out = new(GlobalRateLimit_Action_GenericKey)
		(*in).DeepCopyInto(*out)
	}
	if in.HeaderValueMatch != nil {
		in, out := &in.HeaderValueMatch, &out.HeaderValueMatch
		*out = new(GlobalRateLimit_Action_HeaderValueMatch)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GlobalRateLimit_Action.
func (in *GlobalRateLimit_Action) DeepCopy() *GlobalRateLimit_Action {
	if in == nil {
		return nil
	}
	out := new(GlobalRateLimit_Action)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GlobalRateLimit_Action_DestinationCluster) DeepCopyInto(out *GlobalRateLimit_Action_DestinationCluster) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GlobalRateLimit_Action_DestinationCluster.
func (in *GlobalRateLimit_Action_DestinationCluster) DeepCopy() *GlobalRateLimit_Action_DestinationCluster {
	if in == nil {
		return nil
	}
	out := new(GlobalRateLimit_Action_DestinationCluster)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GlobalRateLimit_Action_GenericKey) DeepCopyInto(out *GlobalRateLimit_Action_GenericKey) {
	*out = *in
	if in.DescriptorKey != nil {
		in, out := &in.DescriptorKey, &out.DescriptorKey
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GlobalRateLimit_Action_GenericKey.
func (in *GlobalRateLimit_Action_GenericKey) DeepCopy() *GlobalRateLimit_Action_GenericKey {
	if in == nil {
		return nil
	}
	out := new(GlobalRateLimit_Action_GenericKey)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GlobalRateLimit_Action_HeaderValueMatch) DeepCopyInto(out *GlobalRateLimit_Action_HeaderValueMatch) {
	*out = *in
	if in.ExpectMatch != nil {
		in, out := &in.ExpectMatch, &out.ExpectMatch
		*out = new(bool)
		**out = **in
	}
	if in.Headers != nil {
		in, out := &in.Headers, &out.Headers
		*out = make([]*GlobalRateLimit_Action_HeaderValueMatch_HeaderMatcher, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(GlobalRateLimit_Action_HeaderValueMatch_HeaderMatcher)
				(*in).DeepCopyInto(*out)
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GlobalRateLimit_Action_HeaderValueMatch.
func (in *GlobalRateLimit_Action_HeaderValueMatch) DeepCopy() *GlobalRateLimit_Action_HeaderValueMatch {
	if in == nil {
		return nil
	}
	out := new(GlobalRateLimit_Action_HeaderValueMatch)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GlobalRateLimit_Action_HeaderValueMatch_HeaderMatcher) DeepCopyInto(out *GlobalRateLimit_Action_HeaderValueMatch_HeaderMatcher) {
	*out = *in
	if in.SafeRegexMatch != nil {
		in, out := &in.SafeRegexMatch, &out.SafeRegexMatch
		*out = new(GlobalRateLimit_Action_HeaderValueMatch_HeaderMatcher_RegexMatcher)
		**out = **in
	}
	if in.RangeMatch != nil {
		in, out := &in.RangeMatch, &out.RangeMatch
		*out = new(GlobalRateLimit_Action_HeaderValueMatch_HeaderMatcher_Int64Range)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GlobalRateLimit_Action_HeaderValueMatch_HeaderMatcher.
func (in *GlobalRateLimit_Action_HeaderValueMatch_HeaderMatcher) DeepCopy() *GlobalRateLimit_Action_HeaderValueMatch_HeaderMatcher {
	if in == nil {
		return nil
	}
	out := new(GlobalRateLimit_Action_HeaderValueMatch_HeaderMatcher)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GlobalRateLimit_Action_HeaderValueMatch_HeaderMatcher_Int64Range) DeepCopyInto(out *GlobalRateLimit_Action_HeaderValueMatch_HeaderMatcher_Int64Range) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GlobalRateLimit_Action_HeaderValueMatch_HeaderMatcher_Int64Range.
func (in *GlobalRateLimit_Action_HeaderValueMatch_HeaderMatcher_Int64Range) DeepCopy() *GlobalRateLimit_Action_HeaderValueMatch_HeaderMatcher_Int64Range {
	if in == nil {
		return nil
	}
	out := new(GlobalRateLimit_Action_HeaderValueMatch_HeaderMatcher_Int64Range)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GlobalRateLimit_Action_HeaderValueMatch_HeaderMatcher_RegexMatcher) DeepCopyInto(out *GlobalRateLimit_Action_HeaderValueMatch_HeaderMatcher_RegexMatcher) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GlobalRateLimit_Action_HeaderValueMatch_HeaderMatcher_RegexMatcher.
func (in *GlobalRateLimit_Action_HeaderValueMatch_HeaderMatcher_RegexMatcher) DeepCopy() *GlobalRateLimit_Action_HeaderValueMatch_HeaderMatcher_RegexMatcher {
	if in == nil {
		return nil
	}
	out := new(GlobalRateLimit_Action_HeaderValueMatch_HeaderMatcher_RegexMatcher)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GlobalRateLimit_Action_RemoteAddress) DeepCopyInto(out *GlobalRateLimit_Action_RemoteAddress) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GlobalRateLimit_Action_RemoteAddress.
func (in *GlobalRateLimit_Action_RemoteAddress) DeepCopy() *GlobalRateLimit_Action_RemoteAddress {
	if in == nil {
		return nil
	}
	out := new(GlobalRateLimit_Action_RemoteAddress)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GlobalRateLimit_Action_RequestHeaders) DeepCopyInto(out *GlobalRateLimit_Action_RequestHeaders) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GlobalRateLimit_Action_RequestHeaders.
func (in *GlobalRateLimit_Action_RequestHeaders) DeepCopy() *GlobalRateLimit_Action_RequestHeaders {
	if in == nil {
		return nil
	}
	out := new(GlobalRateLimit_Action_RequestHeaders)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GlobalRateLimit_Action_SourceCluster) DeepCopyInto(out *GlobalRateLimit_Action_SourceCluster) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GlobalRateLimit_Action_SourceCluster.
func (in *GlobalRateLimit_Action_SourceCluster) DeepCopy() *GlobalRateLimit_Action_SourceCluster {
	if in == nil {
		return nil
	}
	out := new(GlobalRateLimit_Action_SourceCluster)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GlobalRateLimit_Limit) DeepCopyInto(out *GlobalRateLimit_Limit) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GlobalRateLimit_Limit.
func (in *GlobalRateLimit_Limit) DeepCopy() *GlobalRateLimit_Limit {
	if in == nil {
		return nil
	}
	out := new(GlobalRateLimit_Limit)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LocalRateLimit) DeepCopyInto(out *LocalRateLimit) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LocalRateLimit.
func (in *LocalRateLimit) DeepCopy() *LocalRateLimit {
	if in == nil {
		return nil
	}
	out := new(LocalRateLimit)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *LocalRateLimit) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LocalRateLimitConfig) DeepCopyInto(out *LocalRateLimitConfig) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LocalRateLimitConfig.
func (in *LocalRateLimitConfig) DeepCopy() *LocalRateLimitConfig {
	if in == nil {
		return nil
	}
	out := new(LocalRateLimitConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *LocalRateLimitConfig) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LocalRateLimitConfigList) DeepCopyInto(out *LocalRateLimitConfigList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]LocalRateLimitConfig, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LocalRateLimitConfigList.
func (in *LocalRateLimitConfigList) DeepCopy() *LocalRateLimitConfigList {
	if in == nil {
		return nil
	}
	out := new(LocalRateLimitConfigList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *LocalRateLimitConfigList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LocalRateLimitConfigSelector) DeepCopyInto(out *LocalRateLimitConfigSelector) {
	*out = *in
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.IstioVersion != nil {
		in, out := &in.IstioVersion, &out.IstioVersion
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.SNI != nil {
		in, out := &in.SNI, &out.SNI
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LocalRateLimitConfigSelector.
func (in *LocalRateLimitConfigSelector) DeepCopy() *LocalRateLimitConfigSelector {
	if in == nil {
		return nil
	}
	out := new(LocalRateLimitConfigSelector)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LocalRateLimitConfigSpec) DeepCopyInto(out *LocalRateLimitConfigSpec) {
	*out = *in
	in.Selector.DeepCopyInto(&out.Selector)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LocalRateLimitConfigSpec.
func (in *LocalRateLimitConfigSpec) DeepCopy() *LocalRateLimitConfigSpec {
	if in == nil {
		return nil
	}
	out := new(LocalRateLimitConfigSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LocalRateLimitConfigStatus) DeepCopyInto(out *LocalRateLimitConfigStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LocalRateLimitConfigStatus.
func (in *LocalRateLimitConfigStatus) DeepCopy() *LocalRateLimitConfigStatus {
	if in == nil {
		return nil
	}
	out := new(LocalRateLimitConfigStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LocalRateLimitList) DeepCopyInto(out *LocalRateLimitList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]LocalRateLimit, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LocalRateLimitList.
func (in *LocalRateLimitList) DeepCopy() *LocalRateLimitList {
	if in == nil {
		return nil
	}
	out := new(LocalRateLimitList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *LocalRateLimitList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LocalRateLimitSelector) DeepCopyInto(out *LocalRateLimitSelector) {
	*out = *in
	if in.Route != nil {
		in, out := &in.Route, &out.Route
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LocalRateLimitSelector.
func (in *LocalRateLimitSelector) DeepCopy() *LocalRateLimitSelector {
	if in == nil {
		return nil
	}
	out := new(LocalRateLimitSelector)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LocalRateLimitSpec) DeepCopyInto(out *LocalRateLimitSpec) {
	*out = *in
	in.Selector.DeepCopyInto(&out.Selector)
	if in.Limit != nil {
		in, out := &in.Limit, &out.Limit
		*out = new(LocalRateLimit_Limit)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LocalRateLimitSpec.
func (in *LocalRateLimitSpec) DeepCopy() *LocalRateLimitSpec {
	if in == nil {
		return nil
	}
	out := new(LocalRateLimitSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LocalRateLimitStatus) DeepCopyInto(out *LocalRateLimitStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LocalRateLimitStatus.
func (in *LocalRateLimitStatus) DeepCopy() *LocalRateLimitStatus {
	if in == nil {
		return nil
	}
	out := new(LocalRateLimitStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LocalRateLimit_Limit) DeepCopyInto(out *LocalRateLimit_Limit) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LocalRateLimit_Limit.
func (in *LocalRateLimit_Limit) DeepCopy() *LocalRateLimit_Limit {
	if in == nil {
		return nil
	}
	out := new(LocalRateLimit_Limit)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RateLimitService) DeepCopyInto(out *RateLimitService) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RateLimitService.
func (in *RateLimitService) DeepCopy() *RateLimitService {
	if in == nil {
		return nil
	}
	out := new(RateLimitService)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *RateLimitService) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RateLimitServiceList) DeepCopyInto(out *RateLimitServiceList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]RateLimitService, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RateLimitServiceList.
func (in *RateLimitServiceList) DeepCopy() *RateLimitServiceList {
	if in == nil {
		return nil
	}
	out := new(RateLimitServiceList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *RateLimitServiceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RateLimitServiceSpec) DeepCopyInto(out *RateLimitServiceSpec) {
	*out = *in
	if in.Kubernetes != nil {
		in, out := &in.Kubernetes, &out.Kubernetes
		*out = new(RateLimitServiceSpec_Kubernetes)
		(*in).DeepCopyInto(*out)
	}
	if in.Backend != nil {
		in, out := &in.Backend, &out.Backend
		*out = new(RateLimitServiceSpec_Backend)
		(*in).DeepCopyInto(*out)
	}
	if in.Monitoring != nil {
		in, out := &in.Monitoring, &out.Monitoring
		*out = new(RateLimitServiceSpec_Monitoring)
		(*in).DeepCopyInto(*out)
	}
	if in.Environment != nil {
		in, out := &in.Environment, &out.Environment
		*out = new(map[string]string)
		if **in != nil {
			in, out := *in, *out
			*out = make(map[string]string, len(*in))
			for key, val := range *in {
				(*out)[key] = val
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RateLimitServiceSpec.
func (in *RateLimitServiceSpec) DeepCopy() *RateLimitServiceSpec {
	if in == nil {
		return nil
	}
	out := new(RateLimitServiceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RateLimitServiceSpec_Backend) DeepCopyInto(out *RateLimitServiceSpec_Backend) {
	*out = *in
	if in.Redis != nil {
		in, out := &in.Redis, &out.Redis
		*out = new(RateLimitServiceSpec_Backend_Redis)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RateLimitServiceSpec_Backend.
func (in *RateLimitServiceSpec_Backend) DeepCopy() *RateLimitServiceSpec_Backend {
	if in == nil {
		return nil
	}
	out := new(RateLimitServiceSpec_Backend)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RateLimitServiceSpec_Backend_Redis) DeepCopyInto(out *RateLimitServiceSpec_Backend_Redis) {
	*out = *in
	if in.Config != nil {
		in, out := &in.Config, &out.Config
		*out = new(RateLimitServiceSpec_Backend_Redis_Config)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RateLimitServiceSpec_Backend_Redis.
func (in *RateLimitServiceSpec_Backend_Redis) DeepCopy() *RateLimitServiceSpec_Backend_Redis {
	if in == nil {
		return nil
	}
	out := new(RateLimitServiceSpec_Backend_Redis)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RateLimitServiceSpec_Backend_Redis_Config) DeepCopyInto(out *RateLimitServiceSpec_Backend_Redis_Config) {
	*out = *in
	if in.PipelineWindow != nil {
		in, out := &in.PipelineWindow, &out.PipelineWindow
		*out = new(string)
		**out = **in
	}
	if in.PipelineLimit != nil {
		in, out := &in.PipelineLimit, &out.PipelineLimit
		*out = new(int)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RateLimitServiceSpec_Backend_Redis_Config.
func (in *RateLimitServiceSpec_Backend_Redis_Config) DeepCopy() *RateLimitServiceSpec_Backend_Redis_Config {
	if in == nil {
		return nil
	}
	out := new(RateLimitServiceSpec_Backend_Redis_Config)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RateLimitServiceSpec_Kubernetes) DeepCopyInto(out *RateLimitServiceSpec_Kubernetes) {
	*out = *in
	if in.ReplicaCount != nil {
		in, out := &in.ReplicaCount, &out.ReplicaCount
		*out = new(int32)
		**out = **in
	}
	if in.Image != nil {
		in, out := &in.Image, &out.Image
		*out = new(string)
		**out = **in
	}
	if in.Resources != nil {
		in, out := &in.Resources, &out.Resources
		*out = new(v1.ResourceRequirements)
		(*in).DeepCopyInto(*out)
	}
	if in.AutoScaling != nil {
		in, out := &in.AutoScaling, &out.AutoScaling
		*out = new(RateLimitServiceSpec_Kubernetes_AutoScaling)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RateLimitServiceSpec_Kubernetes.
func (in *RateLimitServiceSpec_Kubernetes) DeepCopy() *RateLimitServiceSpec_Kubernetes {
	if in == nil {
		return nil
	}
	out := new(RateLimitServiceSpec_Kubernetes)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RateLimitServiceSpec_Kubernetes_AutoScaling) DeepCopyInto(out *RateLimitServiceSpec_Kubernetes_AutoScaling) {
	*out = *in
	if in.MaxReplica != nil {
		in, out := &in.MaxReplica, &out.MaxReplica
		*out = new(int32)
		**out = **in
	}
	if in.MinReplica != nil {
		in, out := &in.MinReplica, &out.MinReplica
		*out = new(int32)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RateLimitServiceSpec_Kubernetes_AutoScaling.
func (in *RateLimitServiceSpec_Kubernetes_AutoScaling) DeepCopy() *RateLimitServiceSpec_Kubernetes_AutoScaling {
	if in == nil {
		return nil
	}
	out := new(RateLimitServiceSpec_Kubernetes_AutoScaling)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RateLimitServiceSpec_Monitoring) DeepCopyInto(out *RateLimitServiceSpec_Monitoring) {
	*out = *in
	if in.Statsd != nil {
		in, out := &in.Statsd, &out.Statsd
		*out = new(RateLimitServiceSpec_Monitoring_Statsd)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RateLimitServiceSpec_Monitoring.
func (in *RateLimitServiceSpec_Monitoring) DeepCopy() *RateLimitServiceSpec_Monitoring {
	if in == nil {
		return nil
	}
	out := new(RateLimitServiceSpec_Monitoring)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RateLimitServiceSpec_Monitoring_Statsd) DeepCopyInto(out *RateLimitServiceSpec_Monitoring_Statsd) {
	*out = *in
	out.Spec = in.Spec
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RateLimitServiceSpec_Monitoring_Statsd.
func (in *RateLimitServiceSpec_Monitoring_Statsd) DeepCopy() *RateLimitServiceSpec_Monitoring_Statsd {
	if in == nil {
		return nil
	}
	out := new(RateLimitServiceSpec_Monitoring_Statsd)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RateLimitServiceSpec_Monitoring_Statsd_Spec) DeepCopyInto(out *RateLimitServiceSpec_Monitoring_Statsd_Spec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RateLimitServiceSpec_Monitoring_Statsd_Spec.
func (in *RateLimitServiceSpec_Monitoring_Statsd_Spec) DeepCopy() *RateLimitServiceSpec_Monitoring_Statsd_Spec {
	if in == nil {
		return nil
	}
	out := new(RateLimitServiceSpec_Monitoring_Statsd_Spec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RateLimitServiceStatus) DeepCopyInto(out *RateLimitServiceStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RateLimitServiceStatus.
func (in *RateLimitServiceStatus) DeepCopy() *RateLimitServiceStatus {
	if in == nil {
		return nil
	}
	out := new(RateLimitServiceStatus)
	in.DeepCopyInto(out)
	return out
}
