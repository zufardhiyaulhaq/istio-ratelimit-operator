# Global Ratelimit Example

## Prerequisite
#### Ratelimit Service
Global Rate Limit in Envoy uses a gRPC API for requesting quota from a rate limiting service. Istio Ratelimit Operator can help you create rate limiting service with `RateLimitService` object. It's deploying [Envoy ratelimit service](https://github.com/envoyproxy/ratelimit). You only need to provide Redis information:

```
---
apiVersion: ratelimit.zufardhiyaulhaq.com/v1alpha1
kind: RateLimitService
metadata:
  name: public-gateway-ratelimit-service
  namespace: istio-system
spec:
  kubernetes:
    replica_count: 2
    auto_scaling:
      max_replicas: 3
      min_replicas: 2
    resources:
      limits:
        cpu: "256m"
        memory: "256Mi"
      requests:
        cpu: "128m"
        memory: "128Mi"     
  backend:
    redis:
      type: "single"
      url: "172.30.0.13:6379"
---
apiVersion: ratelimit.zufardhiyaulhaq.com/v1alpha1
kind: RateLimitService
metadata:
  name: echo-redis-ratelimit-service
  namespace: default
spec:
  kubernetes:
    replica_count: 2
    auto_scaling:
      max_replicas: 3
      min_replicas: 2
    resources:
      limits:
        cpu: "256m"
        memory: "256Mi"
      requests:
        cpu: "128m"
        memory: "128Mi"     
  backend:
    redis:
      type: "single"
      url: "172.30.0.13:6379"
```

It's support single, sentinel, or clustered Redis. `spec.backend.redis.url` is very depends on `spec.backend.redis.type`. You can check official [Envoy ratelimit service](https://github.com/envoyproxy/ratelimit#redis-type) service for Sentinel and Clustered Redis.

## Gateway
To setup rate limit in Gateway, the first thing you need to do is to make sure the gateway is aware of external rate limit service, you can create `GlobalRateLimitConfig` object to enable that:

```
apiVersion: ratelimit.zufardhiyaulhaq.com/v1alpha1
kind: GlobalRateLimitConfig
metadata:
  name: istio-public-gateway
  namespace: istio-system
spec:
  type: "gateway"
  selector:
    labels:
      "app": "istio-public-gateway"
    istio_version:
      - "1.8"
      - "1.9"
      - "1.10"
  ratelimit:
    spec:
      domain: "public-gateway"
      failure_mode_deny: false
      timeout: "10s"
      service:
        type: "service"
        name: "public-gateway-ratelimit-service"
```

External rate limit service is defined in the `spec.ratelimit.spec.service`. If you using `RateLimitService` object, you can use `type: service` and name to the `RateLimitService` object name.

If you deploy rate limit service by yourself, you can use this configuration instead:
```
  ratelimit:
    spec:
      service:
        type: "fqdn"
        address: "ratelimit.infrastructure.cluster.svc.local"
        port: 8081
```

The next step is to define the rate limit configuration using `GlobalRateLimit` object, for example:

```
apiVersion: ratelimit.zufardhiyaulhaq.com/v1alpha1
kind: GlobalRateLimit
metadata:
  name: helloworld-zufardhiyaulhaq-com-foo-route
  namespace: istio-system
spec:
  config: "istio-public-gateway"
  selector:
    vhost: "helloworld.zufardhiyaulhaq.com:443"
    route: "foo-route"
  matcher:
  - request_headers:
      header_name: ":method"
      descriptor_key: "method"
  - request_headers:
      header_name: ":path"
      descriptor_key: "path"
  - generic_key:
      descriptor_value: "foo-route"
      descriptor_key: "route"
  limit:
    unit: hour
    requests_per_unit: 60
  shadow_mode: false
---
apiVersion: ratelimit.zufardhiyaulhaq.com/v1alpha1
kind: GlobalRateLimit
metadata:
  name: helloworld-zufardhiyaulhaq-com-bar-route
  namespace: istio-system
spec:
  config: "istio-public-gateway"
  selector:
    vhost: "helloworld.zufardhiyaulhaq.com:443"
    route: "bar-route"
  matcher:
  - request_headers:
      header_name: ":method"
      descriptor_key: "method"
  - request_headers:
      header_name: ":path"
      descriptor_key: "path"
  - generic_key:
      descriptor_value: "bar-route"
      descriptor_key: "route"
  limit:
    unit: hour
    requests_per_unit: 120
  shadow_mode: false
```

You must define the `GlobalRateLimitConfig` in the `spec.config`. Also you must define the selector, which is contain two things:
- **vhost**: combination of domain and port from Gateway object
- **route**: route name you define in VirtualService object

Istio Ratelimit Operator will generate two Envoyfilter and descriptors configuration based on [Envoy ratelimit service](https://github.com/envoyproxy/ratelimit). Descriptor example:
```
domain: public-gateway
descriptors:
- key: method
  descriptors:
  - key: path
    descriptors:
    - key: route
      value: foo-route
      rate_limit:
        unit: hour
        requests_per_unit: 60
    - key: route
      value: bar-route
      rate_limit:
        unit: hour
        requests_per_unit: 120
```

## Sidecar
To setup rate limit in Sidecar, it's kinda similar with Gateway. First, create GlobalRateLimitConfig

```
apiVersion: ratelimit.zufardhiyaulhaq.com/v1alpha1
kind: GlobalRateLimitConfig
metadata:
  name: echo-redis
  namespace: default
spec:
  type: "sidecar"
  selector:
    labels:
      "app": "echo-redis"
    istio_version:
      - "1.8"
      - "1.9"
      - "1.10"
  ratelimit:
    spec:
      domain: "echo-redis"
      failure_mode_deny: false
      timeout: "10s"
      service:
        type: "service"
        name: "echo-redis-ratelimit-service"
```

Please make sure `.spec.type` is sidecar and `.spec.selector.sni` is not supported in sidecar. The next step is to define the rate limit configuration using `GlobalRateLimit` object, for example:

```
apiVersion: ratelimit.zufardhiyaulhaq.com/v1alpha1
kind: GlobalRateLimit
metadata:
  name: echo-redis-http-8080
  namespace: default
spec:
  config: "echo-redis"
  selector:
    vhost: "inbound|http|8080"
  matcher:
  - request_headers:
      header_name: ":method"
      descriptor_key: "method"
  - request_headers:
      header_name: ":path"
      descriptor_key: "path"
  - generic_key:
      descriptor_value: "echo-redis"
      descriptor_key: "app"
  - generic_key:
      descriptor_value: "8080"
      descriptor_key: "port"
  limit:
    unit: hour
    requests_per_unit: 60
  shadow_mode: false
```

We only support vhost selector in sidecar, with combination of `inbound|<port-name>|<port-number>`.
