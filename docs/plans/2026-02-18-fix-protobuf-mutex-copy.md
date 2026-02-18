# Fix Protobuf Mutex Copy Errors Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Fix Go vet errors caused by copying sync.Mutex embedded in protobuf types when comparing/assigning EnvoyFilter.Spec

**Architecture:** Replace `equality.Semantic.DeepEqual` with `proto.Equal` for comparison, and use `proto.Reset` + `proto.Merge` for assignment. Fix test range loops to use index-based iteration to avoid copying structs with embedded mutexes.

**Tech Stack:** Go, google.golang.org/protobuf/proto

---

### Task 1: Fix globalratelimit_controller.go

**Files:**
- Modify: `internal/controller/globalratelimit_controller.go:19-38` (imports)
- Modify: `internal/controller/globalratelimit_controller.go:127-128` (comparison and assignment)

**Step 1: Update imports**

Remove unused `equality` import, add `proto` import:

```go
import (
	"context"
	"fmt"
	"time"

	"google.golang.org/protobuf/proto"

	clientnetworking "istio.io/client-go/pkg/apis/networking/v1alpha3"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	funk "github.com/thoas/go-funk"
	ratelimitv1alpha1 "github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"

	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/global/ratelimit"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/utils"
)
```

**Step 2: Fix comparison and assignment (lines 127-128)**

Change from:
```go
if !equality.Semantic.DeepEqual(createdEnvoyFilter.Spec, envoyFilter.Spec) {
    createdEnvoyFilter.Spec = envoyFilter.Spec
```

To:
```go
if !proto.Equal(&createdEnvoyFilter.Spec, &envoyFilter.Spec) {
    createdEnvoyFilter.Spec.Reset()
    proto.Merge(&createdEnvoyFilter.Spec, &envoyFilter.Spec)
```

**Step 3: Run vet to verify fix**

Run: `go vet ./internal/controller/globalratelimit_controller.go`
Expected: No errors

---

### Task 2: Fix globalratelimitconfig_controller.go

**Files:**
- Modify: `internal/controller/globalratelimitconfig_controller.go:19-38` (imports)
- Modify: `internal/controller/globalratelimitconfig_controller.go:139-140` (comparison and assignment)

**Step 1: Update imports**

Remove unused `equality` import, add `proto` import:

```go
import (
	"context"
	"fmt"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/global/config"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/utils"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	funk "github.com/thoas/go-funk"
	ratelimitv1alpha1 "github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"
	clientnetworking "istio.io/client-go/pkg/apis/networking/v1alpha3"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)
```

**Step 2: Fix comparison and assignment (lines 139-140)**

Change from:
```go
if !equality.Semantic.DeepEqual(createdEnvoyFilter.Spec, envoyFilter.Spec) {
    createdEnvoyFilter.Spec = envoyFilter.Spec
```

To:
```go
if !proto.Equal(&createdEnvoyFilter.Spec, &envoyFilter.Spec) {
    createdEnvoyFilter.Spec.Reset()
    proto.Merge(&createdEnvoyFilter.Spec, &envoyFilter.Spec)
```

**Step 3: Run vet to verify fix**

Run: `go vet ./internal/controller/globalratelimitconfig_controller.go`
Expected: No errors

---

### Task 3: Fix localratelimit_controller.go

**Files:**
- Modify: `internal/controller/localratelimit_controller.go:19-37` (imports)
- Modify: `internal/controller/localratelimit_controller.go:126-127` (comparison and assignment)

**Step 1: Update imports**

Remove unused `equality` import, add `proto` import:

```go
import (
	"context"
	"fmt"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/local/ratelimit"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/utils"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	funk "github.com/thoas/go-funk"
	ratelimitv1alpha1 "github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"
	clientnetworking "istio.io/client-go/pkg/apis/networking/v1alpha3"
	ctrl "sigs.k8s.io/controller-runtime"
)
```

**Step 2: Fix comparison and assignment (lines 126-127)**

Change from:
```go
if !equality.Semantic.DeepEqual(createdEnvoyFilter.Spec, envoyFilter.Spec) {
    createdEnvoyFilter.Spec = envoyFilter.Spec
```

To:
```go
if !proto.Equal(&createdEnvoyFilter.Spec, &envoyFilter.Spec) {
    createdEnvoyFilter.Spec.Reset()
    proto.Merge(&createdEnvoyFilter.Spec, &envoyFilter.Spec)
```

**Step 3: Run vet to verify fix**

Run: `go vet ./internal/controller/localratelimit_controller.go`
Expected: No errors

---

### Task 4: Fix localratelimitconfig_controller.go

**Files:**
- Modify: `internal/controller/localratelimitconfig_controller.go:19-37` (imports)
- Modify: `internal/controller/localratelimitconfig_controller.go:113-114` (comparison and assignment)

**Step 1: Update imports**

Remove unused `equality` import, add `proto` import:

```go
import (
	"context"
	"fmt"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/local/config"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/utils"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	funk "github.com/thoas/go-funk"
	ratelimitv1alpha1 "github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"
	clientnetworking "istio.io/client-go/pkg/apis/networking/v1alpha3"
	ctrl "sigs.k8s.io/controller-runtime"
)
```

**Step 2: Fix comparison and assignment (lines 113-114)**

Change from:
```go
if !equality.Semantic.DeepEqual(createdEnvoyFilter.Spec, envoyFilter.Spec) {
    createdEnvoyFilter.Spec = envoyFilter.Spec
```

To:
```go
if !proto.Equal(&createdEnvoyFilter.Spec, &envoyFilter.Spec) {
    createdEnvoyFilter.Spec.Reset()
    proto.Merge(&createdEnvoyFilter.Spec, &envoyFilter.Spec)
```

**Step 3: Run vet to verify fix**

Run: `go vet ./internal/controller/localratelimitconfig_controller.go`
Expected: No errors

---

### Task 5: Fix v3_gateway_builder_test.go range loop

**Files:**
- Modify: `pkg/global/ratelimit/v3_gateway_builder_test.go:524` (range loop)

**Step 1: Fix range loop**

Change from:
```go
func TestNewV3GatewayBuilder(t *testing.T) {
	for _, test := range V3GatewayBuilderTestGrid {
		t.Run(test.name, func(t *testing.T) {
```

To:
```go
func TestNewV3GatewayBuilder(t *testing.T) {
	for i := range V3GatewayBuilderTestGrid {
		test := &V3GatewayBuilderTestGrid[i]
		t.Run(test.name, func(t *testing.T) {
```

**Step 2: Run vet to verify fix**

Run: `go vet ./pkg/global/ratelimit/...`
Expected: No copylocks error for v3_gateway_builder_test.go

---

### Task 6: Fix v3_sidecar_builder_test.go range loop

**Files:**
- Modify: `pkg/global/ratelimit/v3_sidecar_builder_test.go:148` (range loop)

**Step 1: Fix range loop**

Change from:
```go
func TestNewV3SidecarBuilder(t *testing.T) {
	for _, test := range V3SidecarBuilderTestGrid {
		t.Run(test.name, func(t *testing.T) {
```

To:
```go
func TestNewV3SidecarBuilder(t *testing.T) {
	for i := range V3SidecarBuilderTestGrid {
		test := &V3SidecarBuilderTestGrid[i]
		t.Run(test.name, func(t *testing.T) {
```

**Step 2: Run vet to verify fix**

Run: `go vet ./pkg/global/ratelimit/...`
Expected: No copylocks errors

---

### Task 7: Run full verification

**Step 1: Run go vet on entire project**

Run: `make vet`
Expected: No errors

**Step 2: Run all tests**

Run: `make test`
Expected: All tests pass

**Step 3: Run build**

Run: `make build`
Expected: Build succeeds

---

### Task 8: Commit changes

**Step 1: Stage and commit**

```bash
git add internal/controller/*.go pkg/global/ratelimit/*_test.go
git commit -m "fix: resolve protobuf mutex copy errors in controllers and tests

- Replace equality.Semantic.DeepEqual with proto.Equal for EnvoyFilter.Spec comparison
- Use proto.Reset + proto.Merge instead of direct assignment to avoid copying mutex
- Fix test range loops to use index-based iteration

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>"
```
