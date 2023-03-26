# Local Ratelimit Example

## Sidecar
To setup rate limit in Sidecar, the first thing you need to do is to make sure that the extension is enabled, you can create `LocalRateLimitConfig` object to enable that:

```
apiVersion: ratelimit.zufardhiyaulhaq.com/v1alpha1
kind: LocalRateLimitConfig
metadata:
  name: podinfo
  namespace: development
spec:
  type: "sidecar"
  selector:
    labels:
      app: podinfo
    istio_version:
      - "1.14"
      - "1.15"
      - "1.16"
      - "1.17"
```

You must add your pod label in the `.spec.selector`. The next step is to define the rate limit configuration using `LocalRateLimit` object, for example:

```
apiVersion: ratelimit.zufardhiyaulhaq.com/v1alpha1
kind: LocalRateLimit
metadata:
  name: podinfo-http-9898
  namespace: development
spec:
  config: "podinfo"
  selector:
    vhost: "inbound|http|9898"
  limit:
    unit: hour
    requests_per_unit: 1
```

Please use this `inbound|<port-name>|<port-number>` combination of `vhost` when applying in the sidecar.

## Gateway
It's similar like setuping rate limit for sidecar, the first thing is you need to create `LocalRateLimitConfig` object to enable the ratelimit:

```
apiVersion: ratelimit.zufardhiyaulhaq.com/v1alpha1
kind: LocalRateLimitConfig
metadata:
  name: gateway
  namespace: istio-system
spec:
  type: "gateway"
  selector:
    labels:
      app: istio-ingressgateway
      istio: ingressgateway
    istio_version:
      - "1.14"
      - "1.15"
      - "1.16"
      - "1.17"
```

You can also add SNI matching in this `LocalRateLimitConfig` by configuring `.spec.selector.sni`. The next step is to define the rate limit configuration using `LocalRateLimit` object, for example:

```
apiVersion: ratelimit.zufardhiyaulhaq.com/v1alpha1
kind: LocalRateLimit
metadata:
  name: podinfo-default-route
  namespace: istio-system
spec:
  config: "gateway"
  selector:
    vhost: "podinfo.e2e.dev:80"
    route: "default-route"
  limit:
    unit: hour
    requests_per_unit: 1
```

You must define the `LocalRateLimitConfig` in the `spec.config`. Also you must define the selector, which is contain two things:
- **vhost**: combination of domain and port from Gateway object
- **route**: route name you define in VirtualService object
