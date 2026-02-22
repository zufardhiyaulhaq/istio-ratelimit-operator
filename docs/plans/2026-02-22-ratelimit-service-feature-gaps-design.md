# RateLimitService Feature Gap Design

**Date:** 2026-02-22
**Status:** Approved

## Problem

The `RateLimitService` CRD and `GlobalRateLimit` CRD in istio-ratelimit-operator only expose a subset of the features provided by [envoy/ratelimit](https://github.com/envoyproxy/ratelimit). Users must fall back to the `spec.environment` escape hatch for many common production needs (TLS, Prometheus, logging), which provides no validation or documentation.

## Approach

Add structured first-class CRD fields for all major feature gaps. All changes are **additive** — existing CRDs continue to work with zero breaking changes. New fields are optional with sensible defaults that preserve current behavior.

### Backward Compatibility Rules

- Every existing field stays where it is with identical behavior.
- `monitoring.enabled: true` (without `monitoring.type`) still deploys the statsd sidecar.
- `backend.redis.auth` (plaintext string) is kept; new `authSecretRef` is preferred and wins if both are set.
- `spec.environment` remains as an escape hatch. First-class fields take precedence if both configure the same setting.

## Feature Gap Summary

| Category | Current State | Gaps |
|----------|--------------|------|
| Backend | Redis basic (type, url, auth, pipeline) | TLS, pool, timeout, per-second, sentinel auth, cache prefix, local cache |
| Server/gRPC | Hardcoded ports, no TLS | TLS/mTLS, connection age, configurable ports |
| Monitoring | Statsd sidecar only | Native Prometheus, OpenTelemetry tracing, DogStatsD |
| Response Headers | None | RateLimit-Limit/Remaining/Reset headers |
| Logging | None | Level, format |
| Global Shadow Mode | None | Service-wide dry-run |
| Descriptor Features | key, value, shadow_mode, rate_limit | unlimited, detailed_metric |
| Kubernetes | Basic (replicas, image, resources, HPA, labels) | Annotations, security context, scheduling, PDB, liveness probe |

## RateLimitService CRD Additions

### Backend

```yaml
spec:
  backend:
    redis:
      # existing fields (unchanged)
      type: "single"
      url: "redis:6379"
      auth: "plaintext"
      config:
        pipeline_window: "1s"
        pipeline_limit: 10

      # new fields
      authSecretRef:                          # preferred over auth
        name: "redis-secret"
        key: "password"
      tls:
        enabled: true
        secretRef: "redis-tls-secret"         # K8s Secret (ca.crt, tls.crt, tls.key)
        skipHostnameVerification: false
      pool:
        size: 10
        onEmptyBehavior: "wait"               # "wait" or "error"
        onEmptyWaitDuration: "1s"
      timeout: "2s"
      sentinelAuth:
        secretRef:
          name: "sentinel-secret"
          key: "password"
      healthCheckActiveConnection: true
      perSecond:
        enabled: true
        url: "redis-persecond:6379"
        authSecretRef:
          name: "redis-ps-secret"
          key: "password"
        tls:
          enabled: true
          secretRef: "redis-ps-tls-secret"
        pool:
          size: 10
    cacheKeyPrefix: "my-service"
    stopCacheKeyIncrementWhenOverlimit: true
```

The operator auto-mounts referenced Secrets as volumes and sets corresponding env vars (`REDIS_TLS=true`, `REDIS_TLS_CACERT=/tls/ca.crt`, etc.).

### Server / gRPC

```yaml
spec:
  server:
    grpc:
      port: 8081
      tls:
        enabled: true
        secretRef: "grpc-tls-secret"          # K8s Secret (tls.crt, tls.key)
      clientTls:
        caCertSecretRef: "grpc-ca-secret"     # K8s Secret (ca.crt)
        san: "ratelimit.example.com"
      maxConnectionAge: "30m"
      maxConnectionAgeGrace: "5m"
    debug:
      port: 6070
```

### Monitoring

```yaml
spec:
  monitoring:
    enabled: true                             # existing (legacy statsd sidecar)
    type: "prometheus"                        # new: "statsd" (default), "prometheus", "dogstatsd"
    prometheus:
      addr: ":9102"
      path: "/metrics"
    tracing:
      enabled: true
      exporterProtocol: "http"                # "http" or "grpc"
      serviceName: "ratelimit"
      serviceNamespace: "default"
      samplingRate: 0.1
    nearLimitRatio: 0.8
    statsFlushInterval: "10s"
```

Behavior:
- `type: "prometheus"`: No statsd sidecar. Ratelimit exposes `/metrics` natively.
- `type` unset + `enabled: true`: Legacy statsd sidecar (current behavior, zero breakage).

### Response Headers

```yaml
spec:
  responseHeaders:
    enabled: true
```

Adds `RateLimit-Limit`, `RateLimit-Remaining`, and `RateLimit-Reset` headers to responses per IETF draft standard.

### Logging

```yaml
spec:
  logging:
    level: "info"                             # debug, info, warning, error
    format: "json"                            # json, text
```

### Global Shadow Mode

```yaml
spec:
  shadowMode: false
```

When `true`, all rate limit checks return "allow" but counters still increment and metrics fire. Useful for first-time deployment observation.

Different from per-descriptor `shadow_mode` on `GlobalRateLimit`:
- Per-descriptor: only that rule is dry-run.
- Global: all rules are dry-run.

### Kubernetes Deployment Enhancements

```yaml
spec:
  kubernetes:
    # existing fields (unchanged)
    replica_count: 3
    image: "envoyproxy/ratelimit:latest"
    resources: {}
    auto_scaling:
      min_replicas: 2
      max_replicas: 10
    extra_labels: {}

    # new fields
    annotations: {}
    securityContext: {}                        # pod-level
    containerSecurityContext: {}               # container-level
    nodeSelector: {}
    tolerations: []
    affinity: {}
    imagePullSecrets: []
    livenessProbe:
      httpGet:
        path: /healthcheck
        port: 8080
    podDisruptionBudget:
      minAvailable: 1
```

## GlobalRateLimit CRD Additions

### Descriptor Features

```yaml
apiVersion: ratelimit.zufardhiyaulhaq.com/v1alpha1
kind: GlobalRateLimit
spec:
  config: "my-ratelimit-config"               # existing
  selector:                                   # existing
    vhost: "app.example.com"
    route: "my-route"
  matcher:                                    # existing
    - request_headers:
        header_name: ":method"
        descriptor_key: "method"
  shadow_mode: false                          # existing
  limit:                                      # existing (extended)
    unit: minute
    requests_per_unit: 100
    unlimited: false                          # new
  identifier: "my-rule"                       # existing
  detailed_metric: false                      # new
```

Changes:
- `limit.unlimited` (bool, default false): When true, the descriptor is whitelisted — no rate limiting applied, no counter incremented. Mutually exclusive with `requests_per_unit`/`unit`.
- `detailed_metric` (bool, default false): Emits per-value metrics instead of per-key. Use only with bounded value sets to avoid cardinality explosion.

Impact on generated config.yaml:

```yaml
domain: my-gateway
descriptors:
  - key: method
    detailed_metric: true
    rate_limit:
      requests_per_unit: 100
      unit: minute
  - key: api_key
    value: "internal-service"
    rate_limit:
      unlimited: true
```

Requires updating `RateLimit_Service_Descriptor` in `pkg/types/ratelimit.go` and the config builder in `pkg/service/configmap_config_builder.go`.

## Implementation Phases

### Phase 1 — Security & Production Readiness
- `backend.redis.tls` with Secret references and auto volume mounting
- `backend.redis.authSecretRef` (Secret-based auth)
- `backend.redis.sentinelAuth` with Secret reference
- `server.grpc.tls` and `server.grpc.clientTls` with Secret references
- `kubernetes.annotations`, `securityContext`, `containerSecurityContext`
- `kubernetes.nodeSelector`, `tolerations`, `affinity`
- `kubernetes.imagePullSecrets`
- `kubernetes.livenessProbe`
- `kubernetes.podDisruptionBudget`

### Phase 2 — Observability
- `monitoring.type` field with `"prometheus"` support (native, no sidecar)
- `monitoring.tracing` (OpenTelemetry config)
- `monitoring.nearLimitRatio`, `monitoring.statsFlushInterval`
- `responseHeaders.enabled`
- `logging.level`, `logging.format`

### Phase 3 — Descriptor Features (GlobalRateLimit CRD)
- `limit.unlimited`
- `detailed_metric`

### Phase 4 — Backend Enhancements
- `backend.redis.pool` (size, onEmptyBehavior, onEmptyWaitDuration)
- `backend.redis.timeout`
- `backend.redis.healthCheckActiveConnection`
- `backend.redis.perSecond` (full per-second Redis instance)
- `backend.cacheKeyPrefix`
- `backend.stopCacheKeyIncrementWhenOverlimit`
- `shadowMode` (global service-level)

## Excluded Features

- **`replaces` / `name` in rate_limit**: Skipped. Rule changes resetting counters is acceptable.
- **Memcache backend**: Experimental upstream, not worth supporting yet.
- **gRPC Unix domain socket**: Niche use case.
