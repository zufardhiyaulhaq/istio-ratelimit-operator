# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

istio-ratelimit-operator is a Kubernetes operator that simplifies rate limiting configuration for Istio service mesh. It manages:
- **Global rate limiting**: Requires external envoy/ratelimit service (Redis-backed)
- **Local rate limiting**: Per-pod rate limiting without external dependencies

The operator creates and manages Istio EnvoyFilters that configure Envoy proxies (in gateways or sidecars) to connect to rate limit services.

## Common Commands

```bash
# Build
make build                    # Build manager binary to bin/manager
make manifests               # Generate CRD manifests and RBAC
make generate                # Generate DeepCopy methods

# Test
make test                    # Run unit tests with coverage
go test ./pkg/... -v         # Run specific package tests
go test ./pkg/global/config/... -run TestBuildGateway  # Run single test

# Lint
make lint                    # Run golangci-lint (installs v1.50.1 if missing)

# Run locally
make run                     # Run controller against current kubeconfig cluster
make install                 # Install CRDs into cluster

# Docker
make docker-build IMG=<tag>  # Build container image
```

## Architecture

### Custom Resource Definitions (CRDs)

Located in `api/v1alpha1/`, the operator manages 5 CRDs in the `ratelimit.zufardhiyaulhaq.com` group:

| CRD | Purpose |
|-----|---------|
| `GlobalRateLimitConfig` | Configures global ratelimit connection settings (domain, service address, failure mode). Creates EnvoyFilter for HTTP filter config. |
| `GlobalRateLimit` | Defines rate limit rules (matcher actions, limits). References a GlobalRateLimitConfig. Creates EnvoyFilter for route-level actions. |
| `RateLimitService` | Deploys envoy/ratelimit service (Deployment, Service, ConfigMaps). Aggregates GlobalRateLimit rules into ratelimit config. |
| `LocalRateLimitConfig` | Configures local (per-pod) ratelimit settings. Creates EnvoyFilter. |
| `LocalRateLimit` | Defines local rate limit rules. References a LocalRateLimitConfig. |

### Controller Flow

```
GlobalRateLimitConfig  ─────────────────────────────────→  EnvoyFilter (HTTP filter patch)
         │
         │ (referenced by)
         ▼
GlobalRateLimit  ───────────────────────────────────────→  EnvoyFilter (route patch)
         │
         │ (aggregated by)
         ▼
RateLimitService  ──→  Deployment + Service + ConfigMaps (envoy/ratelimit service)
```

### Key Packages

- `internal/controller/`: Kubernetes controllers using controller-runtime
- `pkg/global/config/`: Builds EnvoyFilters for GlobalRateLimitConfig (gateway/sidecar variants)
- `pkg/global/ratelimit/`: Builds EnvoyFilters for GlobalRateLimit route actions
- `pkg/local/config/`: Builds EnvoyFilters for LocalRateLimitConfig
- `pkg/local/ratelimit/`: Builds EnvoyFilters for LocalRateLimit
- `pkg/service/`: Builders for RateLimitService resources (Deployment, Service, ConfigMaps)
- `pkg/utils/version.go`: Istio version mapping for EnvoyFilter compatibility

### EnvoyFilter Versioning

The operator creates version-specific EnvoyFilters for each Istio version specified in `spec.selector.istio_version`. EnvoyFilter names follow the pattern `{name}-{version}` (e.g., `my-config-1.25`). Supported versions are in `pkg/utils/version.go`.

### Context Types

Both Global and Local configs support two context types:
- `gateway`: Targets Istio ingress gateway
- `sidecar`: Targets application sidecars

## Testing

Unit tests use the standard Go testing library with `testify` assertions. Test files follow the `*_test.go` convention and are colocated with source files.

E2E tests use Python scripts in `e2e/scripts/` and run against k3d clusters with real Istio installations. See Makefile targets like `e2e.global.gateway`.

## Helm Chart

The operator is deployed via Helm chart in `charts/istio-ratelimit-operator/`. CRDs are in `crds/crds.yaml`. After CRD changes, regenerate with `make manifests` and copy to the chart.
