# Bug Fixes Design Document

## Overview

This document outlines bugs discovered in the `pkg/` directory and their proposed fixes.

## Bugs and Fixes

### Bug 1: Panic Risk in YAML Type Assertion

**Location:** `pkg/utils/struct.go:45`

**Problem:** Unsafe type assertion `k.(string)` panics if YAML key is not a string.

**Fix:** Add type assertion with ok check, skip non-string keys.

```go
keyStr, ok := k.(string)
if !ok {
    continue
}
m[keyStr] = convertValue(v)
```

---

### Bug 2: Index Out of Bounds in SyncDescriptors

**Location:** `pkg/service/configmap_config_builder.go:84-86`

**Problem:** Accessing `descriptorsData[0]` and `descriptorsData[1:]` without length check.

**Fix:** Return early if slice is empty.

```go
func SyncDescriptors(descriptorsData []types.RateLimit_Service_Descriptor) []types.RateLimit_Service_Descriptor {
    if len(descriptorsData) == 0 {
        return nil
    }
    // ... existing code
}
```

---

### Bug 3: Nil Pointer Dereference in HPA Builder

**Location:** `pkg/service/hpa_builder.go.go:32-33`

**Problem:** No nil checks before accessing nested fields `Kubernetes.AutoScaling.MaxReplica`.

**Fix:** The caller should ensure these fields are set, but defensive checks should be added or documented as precondition.

---

### Bug 4: Nil Pointer Dereference in Deployment Labels

**Location:** `pkg/service/deployment_builder.go:229`

**Problem:** Accessing `Spec.Kubernetes.ExtraLabels` without checking if `Kubernetes` is nil.

**Fix:** Add nil check for `Kubernetes` before accessing `ExtraLabels`.

```go
if n.RateLimitService.Spec.Kubernetes != nil && n.RateLimitService.Spec.Kubernetes.ExtraLabels != nil {
```

---

### Bug 5: Wrong Container Index for Resources

**Location:** `pkg/service/deployment_builder.go:179-183`

**Problem:** When monitoring is enabled, containers are ordered `[statsd-exporter, ratelimit]`. Resources are applied to index 0 (statsd-exporter) first, then index 1. But the intent is to apply resources to the ratelimit container.

**Fix:** Find the ratelimit container by name instead of using hardcoded indices, or reorder the logic to apply resources before prepending the statsd-exporter container.

---

### Bug 6: Typo "rateltimit" in Labels

**Location:** Multiple files:
- `pkg/service/deployment_builder.go:214`
- `pkg/service/service_builder.go:87`
- `pkg/service/hpa_builder.go.go:60`
- `pkg/service/configmap_config_builder.go:50`
- `pkg/service/configmap_env_builder.go` (if exists)
- `pkg/service/configmap_statsd_builder.go` (if exists)

**Problem:** Label value is `"istio-rateltimit-operator"` instead of `"istio-ratelimit-operator"`.

**Fix:** Replace all occurrences with correct spelling.

**Note:** This is a breaking change for existing deployments. Resources with old labels won't match selectors with new labels. Consider:
1. Accept this as a fix and update documentation
2. Keep the typo for backwards compatibility (not recommended)

---

### Bug 7: Double .go Extension in Filename

**Location:** `pkg/service/hpa_builder.go.go`

**Problem:** File has `.go.go` extension.

**Fix:** Rename to `pkg/service/hpa_builder.go`.

---

## Testing Strategy

Each fix should have:
1. A failing test demonstrating the bug
2. The fix applied
3. Test passes

## Risk Assessment

| Bug | Severity | Risk of Fix |
|-----|----------|-------------|
| Bug 1 (type assertion) | High - panic | Low |
| Bug 2 (index bounds) | High - panic | Low |
| Bug 3 (nil deref HPA) | Medium - panic in specific path | Low |
| Bug 4 (nil deref labels) | High - panic | Low |
| Bug 5 (wrong container) | Medium - wrong behavior | Medium |
| Bug 6 (typo) | Low - cosmetic | Medium - breaking change |
| Bug 7 (filename) | Low - cosmetic | Low |

## Recommended Fix Order

1. Bug 1, 2, 4 - High severity panics, low risk fixes
2. Bug 5 - Medium severity, requires careful testing
3. Bug 3 - Medium severity, needs caller analysis
4. Bug 7 - Rename file
5. Bug 6 - Address typo last (breaking change discussion)
