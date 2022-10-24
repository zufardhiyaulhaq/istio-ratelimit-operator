package service

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/types"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type StatsdConfigBuilder struct {
	Config           string
	RateLimitService v1alpha1.RateLimitService
}

func NewStatsdConfigBuilder() *StatsdConfigBuilder {
	return &StatsdConfigBuilder{}
}

func (n *StatsdConfigBuilder) SetRateLimitService(rateLimitService v1alpha1.RateLimitService) *StatsdConfigBuilder {
	n.RateLimitService = rateLimitService
	return n
}

func (n *StatsdConfigBuilder) SetConfig(config string) *StatsdConfigBuilder {
	n.Config = config
	return n
}

func (n *StatsdConfigBuilder) Build() (*corev1.ConfigMap, error) {
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      n.RateLimitService.Name + "-statsd-config",
			Namespace: n.RateLimitService.Namespace,
			Labels:    n.BuildLabels(),
		},
		Data: map[string]string{
			"statsd.mappingConf": n.Config,
		},
	}

	return configMap, nil
}

func (n *StatsdConfigBuilder) BuildLabels() map[string]string {
	var labels = map[string]string{
		"app.kubernetes.io/name":       n.RateLimitService.Name + "-statsd-config",
		"app.kubernetes.io/managed-by": "istio-rateltimit-operator",
		"app.kubernetes.io/created-by": n.RateLimitService.Name,
	}

	return labels
}

func NewStatsdConfig(rateLimitServiceName string, globalRateLimitDomain string, globalRateLimitList []v1alpha1.GlobalRateLimit) (types.MetricMapper, error) {
	metricMapper := types.MetricMapper{}

	for _, globalRateLimit := range globalRateLimitList {
		if globalRateLimit.Spec.Identifier != nil {
			metricMappings, err := NewMetricMappingFromGlobalRateLimit(rateLimitServiceName, globalRateLimitDomain, globalRateLimit)
			if err != nil {
				return metricMapper, err
			}

			metricMapper.Mappings = append(metricMapper.Mappings, metricMappings...)
		}
	}

	metricMapper.Mappings = append(metricMapper.Mappings, NewDefaultMetricMapping()...)

	return metricMapper, nil
}

func NewMetricMappingFromGlobalRateLimit(rateLimitServiceName string, globalRateLimitDomain string, globalRateLimit v1alpha1.GlobalRateLimit) ([]types.MetricMapping, error) {
	metricMappings := []types.MetricMapping{}

	regexMatcher := NewStatsdRegexMatcherFromGlobalRateLimitMatcher(globalRateLimit.Spec.Matcher)

	nearLimitMetricMapping := types.MetricMapping{
		Name:      "ratelimit_service_rate_limit_near_limit",
		Match:     "ratelimit.service.rate_limit." + globalRateLimitDomain + "." + regexMatcher + ".near_limit",
		TimerType: types.ObserverTypeHistogram,
		Labels: prometheus.Labels{
			"identifier":              *globalRateLimit.Spec.Identifier,
			"rate_limit_service_name": rateLimitServiceName,
			"global_rate_limit_name":  globalRateLimit.Name,
		},
	}

	overLimitMetricMapping := types.MetricMapping{
		Name:      "ratelimit_service_rate_limit_over_limit",
		Match:     "ratelimit.service.rate_limit." + globalRateLimitDomain + "." + regexMatcher + ".over_limit",
		TimerType: types.ObserverTypeHistogram,
		Labels: prometheus.Labels{
			"identifier":              *globalRateLimit.Spec.Identifier,
			"rate_limit_service_name": rateLimitServiceName,
			"global_rate_limit_name":  globalRateLimit.Name,
		},
	}

	overLimitWithLocalCacheMetricMapping := types.MetricMapping{
		Name:      "ratelimit_service_rate_limit_over_limit_with_local_cache",
		Match:     "ratelimit.service.rate_limit." + globalRateLimitDomain + "." + regexMatcher + ".over_limit_with_local_cache",
		TimerType: types.ObserverTypeHistogram,
		Labels: prometheus.Labels{
			"identifier":              *globalRateLimit.Spec.Identifier,
			"rate_limit_service_name": rateLimitServiceName,
			"global_rate_limit_name":  globalRateLimit.Name,
		},
	}

	totalHitsMetricMapping := types.MetricMapping{
		Name:      "ratelimit_service_rate_limit_total_hits",
		Match:     "ratelimit.service.rate_limit." + globalRateLimitDomain + "." + regexMatcher + ".total_hits",
		TimerType: types.ObserverTypeHistogram,
		Labels: prometheus.Labels{
			"identifier":              *globalRateLimit.Spec.Identifier,
			"rate_limit_service_name": rateLimitServiceName,
			"global_rate_limit_name":  globalRateLimit.Name,
		},
	}

	withinLimitMetricMapping := types.MetricMapping{
		Name:      "ratelimit_service_rate_limit_within_limit",
		Match:     "ratelimit.service.rate_limit." + globalRateLimitDomain + "." + regexMatcher + ".within_limit",
		TimerType: types.ObserverTypeHistogram,
		Labels: prometheus.Labels{
			"identifier":              *globalRateLimit.Spec.Identifier,
			"rate_limit_service_name": rateLimitServiceName,
			"global_rate_limit_name":  globalRateLimit.Name,
		},
	}

	shadowModeMetricMapping := types.MetricMapping{
		Name:      "ratelimit_service_rate_limit_shadow_mode",
		Match:     "ratelimit.service.rate_limit." + globalRateLimitDomain + "." + regexMatcher + ".shadow_mode",
		TimerType: types.ObserverTypeHistogram,
		Labels: prometheus.Labels{
			"identifier":              *globalRateLimit.Spec.Identifier,
			"rate_limit_service_name": rateLimitServiceName,
			"global_rate_limit_name":  globalRateLimit.Name,
		},
	}

	metricMappings = append(metricMappings, nearLimitMetricMapping, overLimitMetricMapping, overLimitWithLocalCacheMetricMapping, totalHitsMetricMapping, withinLimitMetricMapping, shadowModeMetricMapping)

	return metricMappings, nil
}

func NewStatsdRegexMatcherFromGlobalRateLimitMatcher(matchers []*v1alpha1.GlobalRateLimit_Action) string {
	var regex string

	matchersLength := len(matchers)
	for index, matcher := range matchers {
		if matcher.RequestHeaders != nil {
			regex = regex + matcher.RequestHeaders.DescriptorKey
		}
		if matcher.GenericKey != nil {
			if matcher.GenericKey.DescriptorKey != nil {
				regex = regex + *matcher.GenericKey.DescriptorKey + "_" + matcher.GenericKey.DescriptorValue
			} else {
				regex = regex + "generic_key" + "_" + matcher.GenericKey.DescriptorValue
			}

		}
		if matcher.HeaderValueMatch != nil {
			regex = regex + "header_match" + "_" + matcher.HeaderValueMatch.DescriptorValue
		}

		if index+1 != matchersLength {
			regex = regex + "."
		}
	}

	return regex
}

func NewDefaultMetricMapping() []types.MetricMapping {
	return []types.MetricMapping{
		{
			Name:            "ratelimit_service_should_rate_limit_error",
			Match:           "ratelimit.service.call.should_rate_limit.*",
			MatchMetricType: types.MetricTypeCounter,
			Labels: prometheus.Labels{
				"err_type": "$1",
			},
		},
		{
			Name:            "ratelimit_service_total_requests",
			Match:           "ratelimit_server.*.total_requests",
			MatchMetricType: types.MetricTypeCounter,
			Labels: prometheus.Labels{
				"grpc_method": "$1",
			},
		},
		{
			Name:      "ratelimit_service_response_time_seconds",
			Match:     "ratelimit_server.*.response_time",
			TimerType: types.ObserverTypeHistogram,
			Labels: prometheus.Labels{
				"grpc_method": "$1",
			},
		},
		{
			Name:  "ratelimit_service_config_load_success",
			Match: "ratelimit.service.config_load_success",
		},
		{
			Name:  "ratelimit_service_config_load_error",
			Match: "ratelimit.service.config_load_error",
		},
	}
}
