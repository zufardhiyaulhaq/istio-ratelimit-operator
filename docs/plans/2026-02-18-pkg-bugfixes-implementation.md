# Pkg Bugfixes Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Fix 7 bugs in the pkg/ directory that cause panics, wrong behavior, or cosmetic issues

**Architecture:** Apply minimal fixes following TDD - write failing test, implement fix, verify. Focus on high-severity panic bugs first, then medium severity, finally cosmetic issues.

**Tech Stack:** Go 1.23, testify/assert

---

### Task 1: Fix YAML Type Assertion Panic

**Files:**
- Modify: `pkg/utils/struct.go:44-46`
- Test: `pkg/utils/struct_test.go`

**Step 1: Write failing test for non-string YAML keys**

Add to `pkg/utils/struct_test.go`:

```go
func TestConvertYaml2Struct_NonStringKeys(t *testing.T) {
	// YAML with integer key (unusual but valid YAML)
	yamlWithIntKey := `
123: "value"
normal_key: "normal_value"
`
	// Should not panic, should return nil or skip non-string keys
	result := ConvertYaml2Struct(yamlWithIntKey)
	// The function should handle this gracefully
	assert.NotPanics(t, func() {
		ConvertYaml2Struct(yamlWithIntKey)
	})
}
```

**Step 2: Run test to verify it fails**

Run: `go test -v ./pkg/utils/... -run TestConvertYaml2Struct_NonStringKeys`
Expected: PANIC with "interface conversion: interface {} is int, not string"

**Step 3: Implement fix**

In `pkg/utils/struct.go`, change lines 44-46 from:

```go
for k, v := range val {
    m[k.(string)] = convertValue(v)
}
```

To:

```go
for k, v := range val {
    keyStr, ok := k.(string)
    if !ok {
        continue
    }
    m[keyStr] = convertValue(v)
}
```

**Step 4: Run test to verify it passes**

Run: `go test -v ./pkg/utils/... -run TestConvertYaml2Struct_NonStringKeys`
Expected: PASS

**Step 5: Run all utils tests**

Run: `go test -v ./pkg/utils/...`
Expected: All tests pass

---

### Task 2: Fix Index Out of Bounds in SyncDescriptors

**Files:**
- Modify: `pkg/service/configmap_config_builder.go:82-86`
- Test: `pkg/service/configmap_config_builder_test.go`

**Step 1: Write failing test for empty descriptors**

Add to `pkg/service/configmap_config_builder_test.go`:

```go
func TestSyncDescriptors_EmptySlice(t *testing.T) {
	// Empty slice should not panic
	assert.NotPanics(t, func() {
		result := service.SyncDescriptors([]types.RateLimit_Service_Descriptor{})
		assert.Nil(t, result)
	})
}
```

**Step 2: Run test to verify it fails**

Run: `go test -v ./pkg/service/... -run TestSyncDescriptors_EmptySlice`
Expected: PANIC with "index out of range"

**Step 3: Implement fix**

In `pkg/service/configmap_config_builder.go`, change lines 82-84 from:

```go
func SyncDescriptors(descriptorsData []types.RateLimit_Service_Descriptor) []types.RateLimit_Service_Descriptor {
	var descriptors []types.RateLimit_Service_Descriptor
	descriptors = append(descriptors, descriptorsData[0])
```

To:

```go
func SyncDescriptors(descriptorsData []types.RateLimit_Service_Descriptor) []types.RateLimit_Service_Descriptor {
	if len(descriptorsData) == 0 {
		return nil
	}
	var descriptors []types.RateLimit_Service_Descriptor
	descriptors = append(descriptors, descriptorsData[0])
```

**Step 4: Run test to verify it passes**

Run: `go test -v ./pkg/service/... -run TestSyncDescriptors_EmptySlice`
Expected: PASS

**Step 5: Run all service tests**

Run: `go test -v ./pkg/service/...`
Expected: All tests pass

---

### Task 3: Fix Nil Pointer Dereference in BuildLabels

**Files:**
- Modify: `pkg/service/deployment_builder.go:229`
- Test: `pkg/service/deployment_builder_test.go`

**Step 1: Write failing test for nil Kubernetes spec**

Add to `pkg/service/deployment_builder_test.go`:

```go
func TestDeploymentBuilder_NilKubernetes(t *testing.T) {
	setting := settings.Settings{
		RateLimitServiceImage: "test:latest",
		StatsdExporterImage:   "statsd:latest",
	}

	rateLimitService := v1alpha1.RateLimitService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-service",
			Namespace: "default",
		},
		Spec: v1alpha1.RateLimitServiceSpec{
			Kubernetes: nil, // nil Kubernetes spec
		},
	}

	// Should not panic when Kubernetes is nil
	assert.NotPanics(t, func() {
		_, err := service.NewDeploymentBuilder(setting).
			SetRateLimitService(rateLimitService).
			Build()
		assert.NoError(t, err)
	})
}
```

**Step 2: Run test to verify it fails**

Run: `go test -v ./pkg/service/... -run TestDeploymentBuilder_NilKubernetes`
Expected: PANIC with "nil pointer dereference"

**Step 3: Implement fix**

In `pkg/service/deployment_builder.go`, change line 229 from:

```go
	if n.RateLimitService.Spec.Kubernetes.ExtraLabels != nil {
```

To:

```go
	if n.RateLimitService.Spec.Kubernetes != nil && n.RateLimitService.Spec.Kubernetes.ExtraLabels != nil {
```

**Step 4: Run test to verify it passes**

Run: `go test -v ./pkg/service/... -run TestDeploymentBuilder_NilKubernetes`
Expected: PASS

**Step 5: Run all service tests**

Run: `go test -v ./pkg/service/...`
Expected: All tests pass

---

### Task 4: Fix Wrong Container Index for Resources

**Files:**
- Modify: `pkg/service/deployment_builder.go:173-186`
- Test: `pkg/service/deployment_builder_test.go`

**Step 1: Verify existing test behavior**

The existing test `TestDeploymentBuilder_WithResourcesAndMonitoring` checks that both containers have resources. Looking at the code:
- When monitoring is enabled, statsd-exporter is prepended to containers (index 0)
- Then resources are applied to containers[0] (statsd-exporter) and containers[1] (ratelimit)

This is actually the INTENDED behavior based on the test expectations. The design document may have misunderstood the intent.

**Step 2: Verify with existing tests**

Run: `go test -v ./pkg/service/... -run TestDeploymentBuilder_WithResourcesAndMonitoring`
Expected: PASS (confirms current behavior is intentional)

**Note:** After review, Bug 5 appears to be by design - both containers get the same resources. No fix needed. Skip to Task 5.

---

### Task 5: Rename File with Double .go Extension

**Files:**
- Rename: `pkg/service/hpa_builder.go.go` -> `pkg/service/hpa_builder.go`
- Update: `pkg/service/hpa_builder_test.go` (if imports change)

**Step 1: Rename the file**

Run: `git mv pkg/service/hpa_builder.go.go pkg/service/hpa_builder.go`

**Step 2: Verify build still works**

Run: `go build ./pkg/service/...`
Expected: Build succeeds

**Step 3: Run tests**

Run: `go test -v ./pkg/service/...`
Expected: All tests pass

---

### Task 6: Fix Typo "rateltimit" in Labels (Breaking Change)

**Note:** This is a breaking change. Existing resources with old labels won't match selectors with new labels. Proceed with caution.

**Files:**
- Modify: `pkg/service/deployment_builder.go:214`
- Modify: `pkg/service/service_builder.go:87`
- Modify: `pkg/service/hpa_builder.go:60`
- Modify: `pkg/service/configmap_config_builder.go:50`
- Modify: `pkg/service/configmap_env_builder.go:46`
- Modify: `pkg/service/configmap_statsd_builder.go:49`
- Update: All corresponding test files

**Step 1: Find all occurrences**

Run: `grep -r "rateltimit" pkg/`
Expected: List of all files with typo

**Step 2: Update deployment_builder.go**

Change line 214 from:
```go
"app.kubernetes.io/managed-by": "istio-rateltimit-operator",
```
To:
```go
"app.kubernetes.io/managed-by": "istio-ratelimit-operator",
```

**Step 3: Update service_builder.go**

Change line 87 from:
```go
"app.kubernetes.io/managed-by": "istio-rateltimit-operator",
```
To:
```go
"app.kubernetes.io/managed-by": "istio-ratelimit-operator",
```

**Step 4: Update hpa_builder.go**

Change line 60 from:
```go
"app.kubernetes.io/managed-by": "istio-rateltimit-operator",
```
To:
```go
"app.kubernetes.io/managed-by": "istio-ratelimit-operator",
```

**Step 5: Update configmap_config_builder.go**

Change line 50 from:
```go
"app.kubernetes.io/managed-by": "istio-rateltimit-operator",
```
To:
```go
"app.kubernetes.io/managed-by": "istio-ratelimit-operator",
```

**Step 6: Update configmap_env_builder.go**

Change line 46 from:
```go
"app.kubernetes.io/managed-by": "istio-rateltimit-operator",
```
To:
```go
"app.kubernetes.io/managed-by": "istio-ratelimit-operator",
```

**Step 7: Update configmap_statsd_builder.go**

Change line 49 from:
```go
"app.kubernetes.io/managed-by": "istio-rateltimit-operator",
```
To:
```go
"app.kubernetes.io/managed-by": "istio-ratelimit-operator",
```

**Step 8: Update test expectations**

Update all test files that assert the old typo value to use the corrected spelling.

**Step 9: Run all tests**

Run: `go test ./pkg/service/...`
Expected: All tests pass

---

### Task 7: Full Verification

**Step 1: Run go vet**

Run: `make vet`
Expected: No errors

**Step 2: Run all tests**

Run: `make test`
Expected: All tests pass

**Step 3: Run build**

Run: `make build`
Expected: Build succeeds

---

### Task 8: Commit Changes

**Step 1: Stage bug fix files**

```bash
git add pkg/utils/struct.go pkg/utils/struct_test.go
git add pkg/service/configmap_config_builder.go pkg/service/configmap_config_builder_test.go
git add pkg/service/deployment_builder.go pkg/service/deployment_builder_test.go
git add pkg/service/hpa_builder.go pkg/service/hpa_builder_test.go
git add pkg/service/service_builder.go pkg/service/service_builder_test.go
git add pkg/service/configmap_env_builder.go pkg/service/configmap_env_builder_test.go
git add pkg/service/configmap_statsd_builder.go pkg/service/configmap_statsd_builder_test.go
```

**Step 2: Commit**

```bash
git commit -m "fix: resolve panic bugs and typos in pkg/ directory

- Fix YAML type assertion panic in struct.go (skip non-string keys)
- Fix index out of bounds in SyncDescriptors (handle empty slice)
- Fix nil pointer dereference in BuildLabels (check Kubernetes != nil)
- Rename hpa_builder.go.go to hpa_builder.go
- Fix typo 'rateltimit' -> 'ratelimit' in managed-by labels

BREAKING CHANGE: Label 'app.kubernetes.io/managed-by' changed from
'istio-rateltimit-operator' to 'istio-ratelimit-operator'. Existing
resources with old labels need manual update or recreation.

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>"
```
