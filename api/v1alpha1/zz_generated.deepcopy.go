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
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GlobalRateLimit) DeepCopyInto(out *GlobalRateLimit) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
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
func (in *GlobalRateLimitSpec) DeepCopyInto(out *GlobalRateLimitSpec) {
	*out = *in
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
