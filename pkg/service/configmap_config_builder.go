package service

import (
	"fmt"

	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/types"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ConfigBuilder struct {
	Config           string
	RateLimitService v1alpha1.RateLimitService
}

func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{}
}

func (n *ConfigBuilder) SetRateLimitService(rateLimitService v1alpha1.RateLimitService) *ConfigBuilder {
	n.RateLimitService = rateLimitService
	return n
}

func (n *ConfigBuilder) SetConfig(config string) *ConfigBuilder {
	n.Config = config
	return n
}

func (n *ConfigBuilder) Build() (*corev1.ConfigMap, error) {
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      n.RateLimitService.Name + "-config",
			Namespace: n.RateLimitService.Namespace,
			Labels:    n.BuildLabels(),
		},
		Data: map[string]string{
			"config.yaml": n.Config,
		},
	}

	return configMap, nil
}

func (n *ConfigBuilder) BuildLabels() map[string]string {
	var labels = map[string]string{
		"app.kubernetes.io/name":       n.RateLimitService.Name + "-config",
		"app.kubernetes.io/managed-by": "istio-rateltimit-operator",
		"app.kubernetes.io/created-by": n.RateLimitService.Name,
	}

	return labels
}

func NewRateLimitConfig(domain string, descriptors []types.RateLimit_Service_Descriptor) (types.RateLimit_Service_Config, error) {
	config := types.RateLimit_Service_Config{
		Domain:      domain,
		Descriptors: descriptors,
	}

	return config, nil
}

func NewRateLimitDescriptor(globalRateLimitList []v1alpha1.GlobalRateLimit) ([]types.RateLimit_Service_Descriptor, error) {
	var descriptor []types.RateLimit_Service_Descriptor

	for _, globalRateLimit := range globalRateLimitList {
		globalRateLimitDescriptor, err := NewRateLimitDescriptorFromGlobalRateLimit(globalRateLimit)
		if err != nil {
			return descriptor, err
		}

		descriptor = append(descriptor, globalRateLimitDescriptor...)
	}

	descriptor = SyncDescriptors(descriptor)
	return descriptor, nil
}

func SyncDescriptors(descriptorsData []types.RateLimit_Service_Descriptor) []types.RateLimit_Service_Descriptor {
	var descriptors []types.RateLimit_Service_Descriptor
	descriptors = append(descriptors, descriptorsData[0])

	for _, descriptorData := range descriptorsData[1:] {
		shouldAppend := true

		for descriptorIndex, descriptor := range descriptors {
			if descriptor.Key == descriptorData.Key && descriptor.Value == descriptorData.Value {
				descriptors[descriptorIndex].Descriptors = SyncDescriptors(append(descriptors[descriptorIndex].Descriptors, descriptorData.Descriptors...))
				shouldAppend = false
				continue
			}
		}

		if shouldAppend {
			descriptors = append(descriptors, descriptorData)
		}
	}

	return descriptors
}

func NewRateLimitDescriptorFromGlobalRateLimit(globalRateLimit v1alpha1.GlobalRateLimit) ([]types.RateLimit_Service_Descriptor, error) {
	var descriptor []types.RateLimit_Service_Descriptor
	var sanitizeMatchers []*v1alpha1.GlobalRateLimit_Action

	for _, matcher := range globalRateLimit.Spec.Matcher {
		if matcher.RequestHeaders != nil || matcher.GenericKey != nil || matcher.HeaderValueMatch != nil {
			sanitizeMatchers = append(sanitizeMatchers, matcher)
			continue
		}
	}

	if len(sanitizeMatchers) == 0 {
		return descriptor, nil
	}

	descriptor, err := NewRateLimitDescriptorFromMatcher(sanitizeMatchers, globalRateLimit.Spec.Limit)
	if err != nil {
		return descriptor, err
	}

	return descriptor, nil
}

func NewRateLimitDescriptorFromMatcher(matchers []*v1alpha1.GlobalRateLimit_Action, limit *v1alpha1.GlobalRateLimit_Limit) ([]types.RateLimit_Service_Descriptor, error) {
	descriptor := []types.RateLimit_Service_Descriptor{
		{},
	}

	matcher := matchers[0]

	if matcher.RequestHeaders != nil {
		descriptor[0].Key = matcher.RequestHeaders.DescriptorKey

		if len(matchers) == 1 {
			descriptor[0].RateLimit.RequestsPerUnit = limit.RequestsPerUnit
			descriptor[0].RateLimit.Unit = limit.Unit

			return descriptor, nil
		}

		nestedDescriptor, err := NewRateLimitDescriptorFromMatcher(matchers[1:], limit)
		if err != nil {
			return descriptor, fmt.Errorf("error")
		}

		descriptor[0].Descriptors = nestedDescriptor
		return descriptor, nil
	}

	if matcher.GenericKey != nil {
		if matcher.GenericKey.DescriptorKey != nil {
			descriptor[0].Key = *matcher.GenericKey.DescriptorKey
		} else {
			descriptor[0].Key = "generic_key"
		}

		descriptor[0].Value = matcher.GenericKey.DescriptorValue

		if len(matchers) == 1 {
			descriptor[0].RateLimit.RequestsPerUnit = limit.RequestsPerUnit
			descriptor[0].RateLimit.Unit = limit.Unit

			return descriptor, nil
		}

		nestedDescriptor, err := NewRateLimitDescriptorFromMatcher(matchers[1:], limit)
		if err != nil {
			return descriptor, fmt.Errorf("error")
		}

		descriptor[0].Descriptors = nestedDescriptor
		return descriptor, nil
	}

	if matcher.HeaderValueMatch != nil {
		descriptor[0].Key = "header_match"
		descriptor[0].Value = matcher.HeaderValueMatch.DescriptorValue

		if len(matchers) == 1 {
			descriptor[0].RateLimit.RequestsPerUnit = limit.RequestsPerUnit
			descriptor[0].RateLimit.Unit = limit.Unit

			return descriptor, nil
		}

		nestedDescriptor, err := NewRateLimitDescriptorFromMatcher(matchers[1:], limit)
		if err != nil {
			return descriptor, fmt.Errorf("error")
		}

		descriptor[0].Descriptors = nestedDescriptor
		return descriptor, nil
	}

	return descriptor, fmt.Errorf("error")
}
