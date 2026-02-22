# RateLimitService Feature Gaps Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add structured first-class CRD fields to RateLimitService and GlobalRateLimit for all major envoy/ratelimit features currently only accessible via the `spec.environment` escape hatch.

**Architecture:** Additive CRD changes across 4 phases. Each phase adds new optional types to `api/v1alpha1/`, updates builders in `pkg/service/` to emit the correct env vars and volume mounts, and regenerates deepcopy/CRD manifests. All changes are backward-compatible: existing CRs work identically.

**Tech Stack:** Go 1.23, controller-runtime, kubebuilder, Kubernetes API types (`corev1`, `appsv1`), testify assertions.

---

## Phase 1 — Security & Production Readiness

### Task 1: Add Redis TLS, authSecretRef, and sentinelAuth CRD types

**Files:**
- Modify: `api/v1alpha1/ratelimitservice_types.go:49-59`

**Step 1: Add the new types to ratelimitservice_types.go**

Add these types after `RateLimitServiceSpec_Backend_Redis_Config` (line 59):

```go
type SecretKeyRef struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

type RateLimitServiceSpec_Backend_Redis_TLS struct {
	Enabled                  bool   `json:"enabled,omitempty"`
	SecretRef                string `json:"secretRef,omitempty"`
	SkipHostnameVerification bool   `json:"skipHostnameVerification,omitempty"`
}

type RateLimitServiceSpec_Backend_Redis_SentinelAuth struct {
	SecretRef *SecretKeyRef `json:"secretRef,omitempty"`
}
```

And add new fields to `RateLimitServiceSpec_Backend_Redis` (line 49-54):

```go
type RateLimitServiceSpec_Backend_Redis struct {
	Type         string                                          `json:"type,omitempty"`
	URL          string                                          `json:"url,omitempty"`
	Auth         string                                          `json:"auth,omitempty"`
	AuthSecretRef *SecretKeyRef                                  `json:"authSecretRef,omitempty"`
	Config       *RateLimitServiceSpec_Backend_Redis_Config      `json:"config,omitempty"`
	TLS          *RateLimitServiceSpec_Backend_Redis_TLS         `json:"tls,omitempty"`
	SentinelAuth *RateLimitServiceSpec_Backend_Redis_SentinelAuth `json:"sentinelAuth,omitempty"`
}
```

**Step 2: Verify it compiles**

Run: `cd /home/user/istio-ratelimit-operator && go build ./api/...`
Expected: SUCCESS (no errors)

**Step 3: Commit**

```bash
git add api/v1alpha1/ratelimitservice_types.go
git commit -m "feat(api): add Redis TLS, authSecretRef, sentinelAuth types"
```

---

### Task 2: Add server/gRPC TLS CRD types

**Files:**
- Modify: `api/v1alpha1/ratelimitservice_types.go`

**Step 1: Add server types after the Backend types**

```go
type RateLimitServiceSpec_Server struct {
	GRPC  *RateLimitServiceSpec_Server_GRPC  `json:"grpc,omitempty"`
	Debug *RateLimitServiceSpec_Server_Debug `json:"debug,omitempty"`
}

type RateLimitServiceSpec_Server_GRPC struct {
	Port                  *int32                                       `json:"port,omitempty"`
	TLS                   *RateLimitServiceSpec_Server_GRPC_TLS        `json:"tls,omitempty"`
	ClientTLS             *RateLimitServiceSpec_Server_GRPC_ClientTLS  `json:"clientTls,omitempty"`
	MaxConnectionAge      *string                                      `json:"maxConnectionAge,omitempty"`
	MaxConnectionAgeGrace *string                                      `json:"maxConnectionAgeGrace,omitempty"`
}

type RateLimitServiceSpec_Server_GRPC_TLS struct {
	Enabled   bool   `json:"enabled,omitempty"`
	SecretRef string `json:"secretRef,omitempty"`
}

type RateLimitServiceSpec_Server_GRPC_ClientTLS struct {
	CACertSecretRef string `json:"caCertSecretRef,omitempty"`
	SAN             string `json:"san,omitempty"`
}

type RateLimitServiceSpec_Server_Debug struct {
	Port *int32 `json:"port,omitempty"`
}
```

Add the `Server` field to `RateLimitServiceSpec`:

```go
type RateLimitServiceSpec struct {
	Kubernetes  *RateLimitServiceSpec_Kubernetes `json:"kubernetes,omitempty"`
	Backend     *RateLimitServiceSpec_Backend    `json:"backend,omitempty"`
	Server      *RateLimitServiceSpec_Server     `json:"server,omitempty"`
	Monitoring  *RateLimitServiceSpec_Monitoring `json:"monitoring,omitempty"`
	Environment *map[string]string               `json:"environment,omitempty"`
}
```

**Step 2: Verify it compiles**

Run: `cd /home/user/istio-ratelimit-operator && go build ./api/...`
Expected: SUCCESS

**Step 3: Commit**

```bash
git add api/v1alpha1/ratelimitservice_types.go
git commit -m "feat(api): add server/gRPC TLS types"
```

---

### Task 3: Add Kubernetes scheduling CRD types

**Files:**
- Modify: `api/v1alpha1/ratelimitservice_types.go:32-43`

**Step 1: Add new fields to `RateLimitServiceSpec_Kubernetes`**

```go
type RateLimitServiceSpec_Kubernetes struct {
	ReplicaCount             *int32                                       `json:"replica_count,omitempty"`
	Image                    *string                                      `json:"image,omitempty"`
	Resources                *corev1.ResourceRequirements                 `json:"resources,omitempty"`
	AutoScaling              *RateLimitServiceSpec_Kubernetes_AutoScaling `json:"auto_scaling,omitempty"`
	ExtraLabels              *map[string]string                           `json:"extra_labels,omitempty"`
	Annotations              *map[string]string                           `json:"annotations,omitempty"`
	SecurityContext          *corev1.PodSecurityContext                   `json:"securityContext,omitempty"`
	ContainerSecurityContext *corev1.SecurityContext                      `json:"containerSecurityContext,omitempty"`
	NodeSelector             map[string]string                            `json:"nodeSelector,omitempty"`
	Tolerations              []corev1.Toleration                          `json:"tolerations,omitempty"`
	Affinity                 *corev1.Affinity                             `json:"affinity,omitempty"`
	ImagePullSecrets         []corev1.LocalObjectReference                `json:"imagePullSecrets,omitempty"`
	LivenessProbe            *corev1.Probe                                `json:"livenessProbe,omitempty"`
	PodDisruptionBudget      *RateLimitServiceSpec_Kubernetes_PDB         `json:"podDisruptionBudget,omitempty"`
}

type RateLimitServiceSpec_Kubernetes_PDB struct {
	MinAvailable   *int32 `json:"minAvailable,omitempty"`
	MaxUnavailable *int32 `json:"maxUnavailable,omitempty"`
}
```

**Step 2: Verify it compiles**

Run: `cd /home/user/istio-ratelimit-operator && go build ./api/...`
Expected: SUCCESS

**Step 3: Commit**

```bash
git add api/v1alpha1/ratelimitservice_types.go
git commit -m "feat(api): add Kubernetes scheduling and security types"
```

---

### Task 4: Run make generate for Phase 1 types

**Files:**
- Regenerated: `api/v1alpha1/zz_generated.deepcopy.go`

**Step 1: Run code generation**

Run: `cd /home/user/istio-ratelimit-operator && make generate`
Expected: SUCCESS, `zz_generated.deepcopy.go` updated with DeepCopy methods for all new types.

**Step 2: Verify it compiles**

Run: `cd /home/user/istio-ratelimit-operator && go build ./...`
Expected: SUCCESS

**Step 3: Commit**

```bash
git add api/v1alpha1/zz_generated.deepcopy.go
git commit -m "chore: regenerate deepcopy for Phase 1 types"
```

---

### Task 5: Write failing tests for Redis TLS env vars in EnvBuilder

**Files:**
- Test: `pkg/service/configmap_env_builder_test.go`

**Step 1: Write the failing test**

Add to `pkg/service/configmap_env_builder_test.go`:

```go
func TestEnvBuilder_BuildRedisEnv_WithTLS(t *testing.T) {
	builder := &service.EnvBuilder{
		RateLimitService: v1alpha1.RateLimitService{
			Spec: v1alpha1.RateLimitServiceSpec{
				Backend: &v1alpha1.RateLimitServiceSpec_Backend{
					Redis: &v1alpha1.RateLimitServiceSpec_Backend_Redis{
						Type: "single",
						URL:  "redis:6379",
						TLS: &v1alpha1.RateLimitServiceSpec_Backend_Redis_TLS{
							Enabled: true,
						},
					},
				},
			},
		},
	}

	env, err := builder.BuildRedisEnv()
	assert.NoError(t, err)
	assert.Equal(t, "true", env["REDIS_TLS"])
}

func TestEnvBuilder_BuildRedisEnv_WithTLSCert(t *testing.T) {
	builder := &service.EnvBuilder{
		RateLimitService: v1alpha1.RateLimitService{
			Spec: v1alpha1.RateLimitServiceSpec{
				Backend: &v1alpha1.RateLimitServiceSpec_Backend{
					Redis: &v1alpha1.RateLimitServiceSpec_Backend_Redis{
						Type: "single",
						URL:  "redis:6379",
						TLS: &v1alpha1.RateLimitServiceSpec_Backend_Redis_TLS{
							Enabled:   true,
							SecretRef: "redis-tls-secret",
						},
					},
				},
			},
		},
	}

	env, err := builder.BuildRedisEnv()
	assert.NoError(t, err)
	assert.Equal(t, "true", env["REDIS_TLS"])
	assert.Equal(t, "/tls/redis/ca.crt", env["REDIS_TLS_CACERT"])
	assert.Equal(t, "/tls/redis/tls.crt", env["REDIS_TLS_CLIENT_CERT"])
	assert.Equal(t, "/tls/redis/tls.key", env["REDIS_TLS_CLIENT_KEY"])
}

func TestEnvBuilder_BuildRedisEnv_WithTLSSkipVerify(t *testing.T) {
	builder := &service.EnvBuilder{
		RateLimitService: v1alpha1.RateLimitService{
			Spec: v1alpha1.RateLimitServiceSpec{
				Backend: &v1alpha1.RateLimitServiceSpec_Backend{
					Redis: &v1alpha1.RateLimitServiceSpec_Backend_Redis{
						Type: "single",
						URL:  "redis:6379",
						TLS: &v1alpha1.RateLimitServiceSpec_Backend_Redis_TLS{
							Enabled:                  true,
							SkipHostnameVerification: true,
						},
					},
				},
			},
		},
	}

	env, err := builder.BuildRedisEnv()
	assert.NoError(t, err)
	assert.Equal(t, "true", env["REDIS_TLS"])
	assert.Equal(t, "true", env["REDIS_TLS_SKIP_HOSTNAME_VERIFICATION"])
}
```

**Step 2: Run test to verify it fails**

Run: `cd /home/user/istio-ratelimit-operator && go test ./pkg/service/... -run TestEnvBuilder_BuildRedisEnv_WithTLS -v`
Expected: FAIL — `REDIS_TLS` key missing from env map

**Step 3: Implement Redis TLS in BuildRedisEnv**

Add to `pkg/service/configmap_env_builder.go` in `BuildRedisEnv()` after the Config block (line 125):

```go
	if n.RateLimitService.Spec.Backend.Redis.TLS != nil {
		if n.RateLimitService.Spec.Backend.Redis.TLS.Enabled {
			data["REDIS_TLS"] = "true"
		}

		if n.RateLimitService.Spec.Backend.Redis.TLS.SecretRef != "" {
			data["REDIS_TLS_CACERT"] = "/tls/redis/ca.crt"
			data["REDIS_TLS_CLIENT_CERT"] = "/tls/redis/tls.crt"
			data["REDIS_TLS_CLIENT_KEY"] = "/tls/redis/tls.key"
		}

		if n.RateLimitService.Spec.Backend.Redis.TLS.SkipHostnameVerification {
			data["REDIS_TLS_SKIP_HOSTNAME_VERIFICATION"] = "true"
		}
	}
```

**Step 4: Run test to verify it passes**

Run: `cd /home/user/istio-ratelimit-operator && go test ./pkg/service/... -run TestEnvBuilder_BuildRedisEnv_WithTLS -v`
Expected: PASS (all 3 tests)

**Step 5: Commit**

```bash
git add pkg/service/configmap_env_builder.go pkg/service/configmap_env_builder_test.go
git commit -m "feat(service): add Redis TLS env vars to EnvBuilder"
```

---

### Task 6: Write failing tests for gRPC TLS env vars in EnvBuilder

**Files:**
- Test: `pkg/service/configmap_env_builder_test.go`
- Modify: `pkg/service/configmap_env_builder.go`

**Step 1: Write the failing test**

Add to `pkg/service/configmap_env_builder_test.go`:

```go
func TestEnvBuilder_BuildGRPCEnv(t *testing.T) {
	grpcPort := int32(8081)
	builder := &service.EnvBuilder{
		RateLimitService: v1alpha1.RateLimitService{
			Spec: v1alpha1.RateLimitServiceSpec{
				Server: &v1alpha1.RateLimitServiceSpec_Server{
					GRPC: &v1alpha1.RateLimitServiceSpec_Server_GRPC{
						Port: &grpcPort,
						TLS: &v1alpha1.RateLimitServiceSpec_Server_GRPC_TLS{
							Enabled:   true,
							SecretRef: "grpc-tls-secret",
						},
						MaxConnectionAge:      strPtr("30m"),
						MaxConnectionAgeGrace: strPtr("5m"),
					},
				},
			},
		},
	}

	env, err := builder.BuildServerEnv()
	assert.NoError(t, err)
	assert.Equal(t, "8081", env["GRPC_PORT"])
	assert.Equal(t, "/tls/grpc/tls.crt", env["GRPC_SERVER_TLS_CERT"])
	assert.Equal(t, "/tls/grpc/tls.key", env["GRPC_SERVER_TLS_KEY"])
	assert.Equal(t, "30m", env["GRPC_MAX_CONNECTION_AGE"])
	assert.Equal(t, "5m", env["GRPC_MAX_CONNECTION_AGE_GRACE"])
}

func TestEnvBuilder_BuildGRPCEnv_WithClientTLS(t *testing.T) {
	builder := &service.EnvBuilder{
		RateLimitService: v1alpha1.RateLimitService{
			Spec: v1alpha1.RateLimitServiceSpec{
				Server: &v1alpha1.RateLimitServiceSpec_Server{
					GRPC: &v1alpha1.RateLimitServiceSpec_Server_GRPC{
						ClientTLS: &v1alpha1.RateLimitServiceSpec_Server_GRPC_ClientTLS{
							CACertSecretRef: "grpc-ca-secret",
							SAN:             "ratelimit.example.com",
						},
					},
				},
			},
		},
	}

	env, err := builder.BuildServerEnv()
	assert.NoError(t, err)
	assert.Equal(t, "/tls/grpc-client/ca.crt", env["GRPC_SERVER_TLS_CLIENT_CACERT"])
	assert.Equal(t, "ratelimit.example.com", env["GRPC_CLIENT_TLS_SAN"])
}
```

You'll also need this helper if not already present:

```go
func strPtr(s string) *string { return &s }
```

**Step 2: Run test to verify it fails**

Run: `cd /home/user/istio-ratelimit-operator && go test ./pkg/service/... -run TestEnvBuilder_BuildGRPCEnv -v`
Expected: FAIL — `BuildServerEnv` method does not exist

**Step 3: Implement BuildServerEnv**

Add to `pkg/service/configmap_env_builder.go`:

```go
func (n *EnvBuilder) BuildServerEnv() (map[string]string, error) {
	data := make(map[string]string)

	if n.RateLimitService.Spec.Server == nil {
		return data, nil
	}

	if n.RateLimitService.Spec.Server.GRPC != nil {
		grpc := n.RateLimitService.Spec.Server.GRPC

		if grpc.Port != nil {
			data["GRPC_PORT"] = strconv.Itoa(int(*grpc.Port))
		}

		if grpc.TLS != nil && grpc.TLS.Enabled {
			if grpc.TLS.SecretRef != "" {
				data["GRPC_SERVER_TLS_CERT"] = "/tls/grpc/tls.crt"
				data["GRPC_SERVER_TLS_KEY"] = "/tls/grpc/tls.key"
			}
		}

		if grpc.ClientTLS != nil {
			if grpc.ClientTLS.CACertSecretRef != "" {
				data["GRPC_SERVER_TLS_CLIENT_CACERT"] = "/tls/grpc-client/ca.crt"
			}
			if grpc.ClientTLS.SAN != "" {
				data["GRPC_CLIENT_TLS_SAN"] = grpc.ClientTLS.SAN
			}
		}

		if grpc.MaxConnectionAge != nil {
			data["GRPC_MAX_CONNECTION_AGE"] = *grpc.MaxConnectionAge
		}
		if grpc.MaxConnectionAgeGrace != nil {
			data["GRPC_MAX_CONNECTION_AGE_GRACE"] = *grpc.MaxConnectionAgeGrace
		}
	}

	if n.RateLimitService.Spec.Server.Debug != nil && n.RateLimitService.Spec.Server.Debug.Port != nil {
		data["DEBUG_PORT"] = strconv.Itoa(int(*n.RateLimitService.Spec.Server.Debug.Port))
	}

	return data, nil
}
```

And wire it into `BuildEnv()` after the monitoring block:

```go
	if n.RateLimitService.Spec.Server != nil {
		serverEnv, err := n.BuildServerEnv()
		if err != nil {
			return env, err
		}
		for key, value := range serverEnv {
			env[key] = value
		}
	}
```

**Step 4: Run test to verify it passes**

Run: `cd /home/user/istio-ratelimit-operator && go test ./pkg/service/... -run TestEnvBuilder_BuildGRPCEnv -v`
Expected: PASS

**Step 5: Commit**

```bash
git add pkg/service/configmap_env_builder.go pkg/service/configmap_env_builder_test.go
git commit -m "feat(service): add gRPC server TLS env vars to EnvBuilder"
```

---

### Task 7: Write failing tests for TLS volume mounts in DeploymentBuilder

**Files:**
- Test: `pkg/service/deployment_builder_test.go`
- Modify: `pkg/service/deployment_builder.go`

**Step 1: Write the failing test**

Add to `pkg/service/deployment_builder_test.go`:

```go
func TestDeploymentBuilder_WithRedisTLS(t *testing.T) {
	setting := settings.Settings{
		RateLimitServiceImage: "ratelimit:latest",
		StatsdExporterImage:   "statsd:latest",
	}

	rateLimitService := v1alpha1.RateLimitService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "tls-test",
			Namespace: "default",
		},
		Spec: v1alpha1.RateLimitServiceSpec{
			Backend: &v1alpha1.RateLimitServiceSpec_Backend{
				Redis: &v1alpha1.RateLimitServiceSpec_Backend_Redis{
					Type: "single",
					URL:  "redis:6379",
					TLS: &v1alpha1.RateLimitServiceSpec_Backend_Redis_TLS{
						Enabled:   true,
						SecretRef: "redis-tls-secret",
					},
				},
			},
		},
	}

	deployment, err := service.NewDeploymentBuilder(setting).
		SetRateLimitService(rateLimitService).
		Build()

	assert.NoError(t, err)

	// Should have redis TLS volume
	foundVolume := false
	for _, v := range deployment.Spec.Template.Spec.Volumes {
		if v.Name == "redis-tls" {
			foundVolume = true
			assert.Equal(t, "redis-tls-secret", v.VolumeSource.Secret.SecretName)
		}
	}
	assert.True(t, foundVolume, "redis-tls volume not found")

	// Should have redis TLS volume mount on ratelimit container
	foundMount := false
	for _, vm := range deployment.Spec.Template.Spec.Containers[0].VolumeMounts {
		if vm.Name == "redis-tls" {
			foundMount = true
			assert.Equal(t, "/tls/redis/", vm.MountPath)
			assert.True(t, vm.ReadOnly)
		}
	}
	assert.True(t, foundMount, "redis-tls volume mount not found")
}
```

**Step 2: Run test to verify it fails**

Run: `cd /home/user/istio-ratelimit-operator && go test ./pkg/service/... -run TestDeploymentBuilder_WithRedisTLS -v`
Expected: FAIL — no `redis-tls` volume found

**Step 3: Implement TLS volume mounts in DeploymentBuilder.Build()**

Add to `pkg/service/deployment_builder.go` in `Build()` before the monitoring block (before line 115):

```go
	if n.RateLimitService.Spec.Backend != nil && n.RateLimitService.Spec.Backend.Redis != nil &&
		n.RateLimitService.Spec.Backend.Redis.TLS != nil && n.RateLimitService.Spec.Backend.Redis.TLS.SecretRef != "" {
		deployment.Spec.Template.Spec.Volumes = append(deployment.Spec.Template.Spec.Volumes, corev1.Volume{
			Name: "redis-tls",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: n.RateLimitService.Spec.Backend.Redis.TLS.SecretRef,
				},
			},
		})
		deployment.Spec.Template.Spec.Containers[0].VolumeMounts = append(
			deployment.Spec.Template.Spec.Containers[0].VolumeMounts,
			corev1.VolumeMount{
				Name:      "redis-tls",
				MountPath: "/tls/redis/",
				ReadOnly:  true,
			},
		)
	}

	if n.RateLimitService.Spec.Server != nil && n.RateLimitService.Spec.Server.GRPC != nil {
		grpc := n.RateLimitService.Spec.Server.GRPC
		if grpc.TLS != nil && grpc.TLS.SecretRef != "" {
			deployment.Spec.Template.Spec.Volumes = append(deployment.Spec.Template.Spec.Volumes, corev1.Volume{
				Name: "grpc-tls",
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: grpc.TLS.SecretRef,
					},
				},
			})
			deployment.Spec.Template.Spec.Containers[0].VolumeMounts = append(
				deployment.Spec.Template.Spec.Containers[0].VolumeMounts,
				corev1.VolumeMount{
					Name:      "grpc-tls",
					MountPath: "/tls/grpc/",
					ReadOnly:  true,
				},
			)
		}
		if grpc.ClientTLS != nil && grpc.ClientTLS.CACertSecretRef != "" {
			deployment.Spec.Template.Spec.Volumes = append(deployment.Spec.Template.Spec.Volumes, corev1.Volume{
				Name: "grpc-client-tls",
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: grpc.ClientTLS.CACertSecretRef,
					},
				},
			})
			deployment.Spec.Template.Spec.Containers[0].VolumeMounts = append(
				deployment.Spec.Template.Spec.Containers[0].VolumeMounts,
				corev1.VolumeMount{
					Name:      "grpc-client-tls",
					MountPath: "/tls/grpc-client/",
					ReadOnly:  true,
				},
			)
		}
	}
```

**Step 4: Run test to verify it passes**

Run: `cd /home/user/istio-ratelimit-operator && go test ./pkg/service/... -run TestDeploymentBuilder_WithRedisTLS -v`
Expected: PASS

**Step 5: Commit**

```bash
git add pkg/service/deployment_builder.go pkg/service/deployment_builder_test.go
git commit -m "feat(service): add TLS secret volume mounts to DeploymentBuilder"
```

---

### Task 8: Write failing tests for Kubernetes scheduling in DeploymentBuilder

**Files:**
- Test: `pkg/service/deployment_builder_test.go`
- Modify: `pkg/service/deployment_builder.go`

**Step 1: Write the failing test**

Add to `pkg/service/deployment_builder_test.go`:

```go
func TestDeploymentBuilder_WithScheduling(t *testing.T) {
	setting := settings.Settings{
		RateLimitServiceImage: "ratelimit:latest",
		StatsdExporterImage:   "statsd:latest",
	}

	rateLimitService := v1alpha1.RateLimitService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sched-test",
			Namespace: "default",
		},
		Spec: v1alpha1.RateLimitServiceSpec{
			Kubernetes: &v1alpha1.RateLimitServiceSpec_Kubernetes{
				NodeSelector: map[string]string{"disktype": "ssd"},
				Tolerations: []corev1.Toleration{
					{Key: "dedicated", Operator: corev1.TolerationOpEqual, Value: "ratelimit", Effect: corev1.TaintEffectNoSchedule},
				},
				Affinity: &corev1.Affinity{
					NodeAffinity: &corev1.NodeAffinity{
						RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
							NodeSelectorTerms: []corev1.NodeSelectorTerm{
								{MatchExpressions: []corev1.NodeSelectorRequirement{
									{Key: "zone", Operator: corev1.NodeSelectorOpIn, Values: []string{"us-east-1a"}},
								}},
							},
						},
					},
				},
				SecurityContext: &corev1.PodSecurityContext{
					RunAsNonRoot: boolPtr(true),
				},
				ContainerSecurityContext: &corev1.SecurityContext{
					ReadOnlyRootFilesystem: boolPtr(true),
				},
				ImagePullSecrets: []corev1.LocalObjectReference{
					{Name: "my-registry-secret"},
				},
				Annotations: &map[string]string{
					"prometheus.io/scrape": "true",
				},
			},
		},
	}

	deployment, err := service.NewDeploymentBuilder(setting).
		SetRateLimitService(rateLimitService).
		Build()

	assert.NoError(t, err)
	assert.Equal(t, map[string]string{"disktype": "ssd"}, deployment.Spec.Template.Spec.NodeSelector)
	assert.Len(t, deployment.Spec.Template.Spec.Tolerations, 1)
	assert.NotNil(t, deployment.Spec.Template.Spec.Affinity)
	assert.NotNil(t, deployment.Spec.Template.Spec.SecurityContext)
	assert.True(t, *deployment.Spec.Template.Spec.SecurityContext.RunAsNonRoot)
	assert.True(t, *deployment.Spec.Template.Spec.Containers[0].SecurityContext.ReadOnlyRootFilesystem)
	assert.Len(t, deployment.Spec.Template.Spec.ImagePullSecrets, 1)
	assert.Equal(t, "true", deployment.Spec.Template.ObjectMeta.Annotations["prometheus.io/scrape"])
}

func boolPtr(b bool) *bool { return &b }
```

**Step 2: Run test to verify it fails**

Run: `cd /home/user/istio-ratelimit-operator && go test ./pkg/service/... -run TestDeploymentBuilder_WithScheduling -v`
Expected: FAIL — NodeSelector is nil

**Step 3: Implement Kubernetes scheduling in DeploymentBuilder.Build()**

Add to `pkg/service/deployment_builder.go` in `Build()` inside the `if n.RateLimitService.Spec.Kubernetes != nil` block (after the Resources block, around line 187):

```go
		if n.RateLimitService.Spec.Kubernetes.NodeSelector != nil {
			deployment.Spec.Template.Spec.NodeSelector = n.RateLimitService.Spec.Kubernetes.NodeSelector
		}

		if n.RateLimitService.Spec.Kubernetes.Tolerations != nil {
			deployment.Spec.Template.Spec.Tolerations = n.RateLimitService.Spec.Kubernetes.Tolerations
		}

		if n.RateLimitService.Spec.Kubernetes.Affinity != nil {
			deployment.Spec.Template.Spec.Affinity = n.RateLimitService.Spec.Kubernetes.Affinity
		}

		if n.RateLimitService.Spec.Kubernetes.SecurityContext != nil {
			deployment.Spec.Template.Spec.SecurityContext = n.RateLimitService.Spec.Kubernetes.SecurityContext
		}

		if n.RateLimitService.Spec.Kubernetes.ContainerSecurityContext != nil {
			deployment.Spec.Template.Spec.Containers[0].SecurityContext = n.RateLimitService.Spec.Kubernetes.ContainerSecurityContext
		}

		if n.RateLimitService.Spec.Kubernetes.ImagePullSecrets != nil {
			deployment.Spec.Template.Spec.ImagePullSecrets = n.RateLimitService.Spec.Kubernetes.ImagePullSecrets
		}

		if n.RateLimitService.Spec.Kubernetes.Annotations != nil {
			deployment.Spec.Template.ObjectMeta.Annotations = *n.RateLimitService.Spec.Kubernetes.Annotations
		}

		if n.RateLimitService.Spec.Kubernetes.LivenessProbe != nil {
			deployment.Spec.Template.Spec.Containers[0].LivenessProbe = n.RateLimitService.Spec.Kubernetes.LivenessProbe
		}
```

**Step 4: Run test to verify it passes**

Run: `cd /home/user/istio-ratelimit-operator && go test ./pkg/service/... -run TestDeploymentBuilder_WithScheduling -v`
Expected: PASS

**Step 5: Commit**

```bash
git add pkg/service/deployment_builder.go pkg/service/deployment_builder_test.go
git commit -m "feat(service): add Kubernetes scheduling fields to DeploymentBuilder"
```

---

### Task 9: Run make manifests for Phase 1 and verify all tests pass

**Files:**
- Regenerated: `config/crd/bases/*.yaml`

**Step 1: Run manifests generation**

Run: `cd /home/user/istio-ratelimit-operator && make manifests`
Expected: SUCCESS

**Step 2: Run all tests**

Run: `cd /home/user/istio-ratelimit-operator && go test ./pkg/... -v`
Expected: ALL PASS

**Step 3: Commit**

```bash
git add config/ api/
git commit -m "chore: regenerate CRD manifests for Phase 1"
```

---

## Phase 2 — Observability

### Task 10: Add monitoring.type, prometheus, tracing, responseHeaders, logging, shadowMode CRD types

**Files:**
- Modify: `api/v1alpha1/ratelimitservice_types.go`

**Step 1: Extend RateLimitServiceSpec with new top-level fields**

```go
type RateLimitServiceSpec struct {
	Kubernetes      *RateLimitServiceSpec_Kubernetes      `json:"kubernetes,omitempty"`
	Backend         *RateLimitServiceSpec_Backend         `json:"backend,omitempty"`
	Server          *RateLimitServiceSpec_Server          `json:"server,omitempty"`
	Monitoring      *RateLimitServiceSpec_Monitoring      `json:"monitoring,omitempty"`
	ResponseHeaders *RateLimitServiceSpec_ResponseHeaders `json:"responseHeaders,omitempty"`
	Logging         *RateLimitServiceSpec_Logging         `json:"logging,omitempty"`
	ShadowMode      bool                                  `json:"shadowMode,omitempty"`
	Environment     *map[string]string                    `json:"environment,omitempty"`
}
```

**Step 2: Extend monitoring type with new fields**

```go
type RateLimitServiceSpec_Monitoring struct {
	// +optional
	Enabled bool `json:"enabled,omitempty"`

	// +kubebuilder:validation:Enum=statsd;prometheus;dogstatsd
	// +optional
	Type string `json:"type,omitempty"`

	Prometheus *RateLimitServiceSpec_Monitoring_Prometheus `json:"prometheus,omitempty"`
	Tracing    *RateLimitServiceSpec_Monitoring_Tracing    `json:"tracing,omitempty"`

	NearLimitRatio     *string `json:"nearLimitRatio,omitempty"`
	StatsFlushInterval *string `json:"statsFlushInterval,omitempty"`

	// This API is deprecated
	Statsd *RateLimitServiceSpec_Monitoring_Statsd `json:"statsd,omitempty"`
}

type RateLimitServiceSpec_Monitoring_Prometheus struct {
	Addr string `json:"addr,omitempty"`
	Path string `json:"path,omitempty"`
}

type RateLimitServiceSpec_Monitoring_Tracing struct {
	Enabled          bool    `json:"enabled,omitempty"`
	ExporterProtocol string  `json:"exporterProtocol,omitempty"`
	ServiceName      string  `json:"serviceName,omitempty"`
	ServiceNamespace string  `json:"serviceNamespace,omitempty"`
	SamplingRate     float64 `json:"samplingRate,omitempty"`
}
```

**Step 3: Add new top-level types**

```go
type RateLimitServiceSpec_ResponseHeaders struct {
	Enabled bool `json:"enabled,omitempty"`
}

type RateLimitServiceSpec_Logging struct {
	// +kubebuilder:validation:Enum=debug;info;warning;error
	Level  string `json:"level,omitempty"`
	// +kubebuilder:validation:Enum=json;text
	Format string `json:"format,omitempty"`
}
```

**Step 4: Verify it compiles**

Run: `cd /home/user/istio-ratelimit-operator && go build ./api/...`
Expected: SUCCESS

**Step 5: Run make generate**

Run: `cd /home/user/istio-ratelimit-operator && make generate`
Expected: SUCCESS

**Step 6: Commit**

```bash
git add api/v1alpha1/ratelimitservice_types.go api/v1alpha1/zz_generated.deepcopy.go
git commit -m "feat(api): add monitoring, responseHeaders, logging, shadowMode types"
```

---

### Task 11: Write failing tests for observability env vars in EnvBuilder

**Files:**
- Test: `pkg/service/configmap_env_builder_test.go`
- Modify: `pkg/service/configmap_env_builder.go`

**Step 1: Write the failing tests**

Add to `pkg/service/configmap_env_builder_test.go`:

```go
func TestEnvBuilder_BuildMonitoringEnv_Prometheus(t *testing.T) {
	builder := &service.EnvBuilder{
		RateLimitService: v1alpha1.RateLimitService{
			Spec: v1alpha1.RateLimitServiceSpec{
				Monitoring: &v1alpha1.RateLimitServiceSpec_Monitoring{
					Enabled: true,
					Type:    "prometheus",
					Prometheus: &v1alpha1.RateLimitServiceSpec_Monitoring_Prometheus{
						Addr: ":9102",
						Path: "/metrics",
					},
					NearLimitRatio:     strPtr("0.8"),
					StatsFlushInterval: strPtr("10s"),
				},
			},
		},
	}

	env, err := builder.BuildMonitoringEnv()
	assert.NoError(t, err)
	assert.Equal(t, "false", env["USE_STATSD"])
	assert.Equal(t, "true", env["USE_PROMETHEUS"])
	assert.Equal(t, ":9102", env["PROMETHEUS_ADDR"])
	assert.Equal(t, "/metrics", env["PROMETHEUS_PATH"])
	assert.Equal(t, "0.8", env["NEAR_LIMIT_RATIO"])
	assert.Equal(t, "10s", env["STATS_FLUSH_INTERVAL"])
}

func TestEnvBuilder_BuildMonitoringEnv_Tracing(t *testing.T) {
	builder := &service.EnvBuilder{
		RateLimitService: v1alpha1.RateLimitService{
			Spec: v1alpha1.RateLimitServiceSpec{
				Monitoring: &v1alpha1.RateLimitServiceSpec_Monitoring{
					Tracing: &v1alpha1.RateLimitServiceSpec_Monitoring_Tracing{
						Enabled:          true,
						ExporterProtocol: "http",
						ServiceName:      "ratelimit",
						ServiceNamespace: "default",
						SamplingRate:     0.1,
					},
				},
			},
		},
	}

	env, err := builder.BuildMonitoringEnv()
	assert.NoError(t, err)
	assert.Equal(t, "true", env["TRACING_ENABLED"])
	assert.Equal(t, "http", env["TRACING_EXPORTER_PROTOCOL"])
	assert.Equal(t, "ratelimit", env["TRACING_SERVICE_NAME"])
	assert.Equal(t, "default", env["TRACING_SERVICE_NAMESPACE"])
	assert.Equal(t, "0.1", env["TRACING_SAMPLING_RATE"])
}

func TestEnvBuilder_BuildResponseHeadersEnv(t *testing.T) {
	builder := &service.EnvBuilder{
		RateLimitService: v1alpha1.RateLimitService{
			Spec: v1alpha1.RateLimitServiceSpec{
				ResponseHeaders: &v1alpha1.RateLimitServiceSpec_ResponseHeaders{
					Enabled: true,
				},
			},
		},
	}

	env, err := builder.BuildResponseHeadersEnv()
	assert.NoError(t, err)
	assert.Equal(t, "true", env["RESPONSE_HEADERS_ENABLED"])
}

func TestEnvBuilder_BuildLoggingEnv(t *testing.T) {
	builder := &service.EnvBuilder{
		RateLimitService: v1alpha1.RateLimitService{
			Spec: v1alpha1.RateLimitServiceSpec{
				Logging: &v1alpha1.RateLimitServiceSpec_Logging{
					Level:  "debug",
					Format: "json",
				},
			},
		},
	}

	env, err := builder.BuildLoggingEnv()
	assert.NoError(t, err)
	assert.Equal(t, "debug", env["LOG_LEVEL"])
	assert.Equal(t, "json", env["LOG_FORMAT"])
}

func TestEnvBuilder_BuildShadowModeEnv(t *testing.T) {
	builder := &service.EnvBuilder{
		RateLimitService: v1alpha1.RateLimitService{
			Spec: v1alpha1.RateLimitServiceSpec{
				ShadowMode: true,
			},
		},
	}

	env, err := builder.BuildShadowModeEnv()
	assert.NoError(t, err)
	assert.Equal(t, "true", env["SHADOW_MODE"])
}
```

**Step 2: Run test to verify it fails**

Run: `cd /home/user/istio-ratelimit-operator && go test ./pkg/service/... -run "TestEnvBuilder_Build(Monitoring|ResponseHeaders|Logging|ShadowMode)Env" -v`
Expected: FAIL — methods do not exist

**Step 3: Implement the new env builder methods**

Add to `pkg/service/configmap_env_builder.go`:

```go
func (n *EnvBuilder) BuildMonitoringEnv() (map[string]string, error) {
	data := make(map[string]string)

	if n.RateLimitService.Spec.Monitoring == nil {
		return data, nil
	}

	mon := n.RateLimitService.Spec.Monitoring

	if mon.Type == "prometheus" {
		data["USE_STATSD"] = "false"
		data["USE_PROMETHEUS"] = "true"
		if mon.Prometheus != nil {
			if mon.Prometheus.Addr != "" {
				data["PROMETHEUS_ADDR"] = mon.Prometheus.Addr
			}
			if mon.Prometheus.Path != "" {
				data["PROMETHEUS_PATH"] = mon.Prometheus.Path
			}
		}
	} else {
		data["USE_STATSD"] = "false"
	}

	if mon.NearLimitRatio != nil {
		data["NEAR_LIMIT_RATIO"] = *mon.NearLimitRatio
	}

	if mon.StatsFlushInterval != nil {
		data["STATS_FLUSH_INTERVAL"] = *mon.StatsFlushInterval
	}

	if mon.Tracing != nil && mon.Tracing.Enabled {
		data["TRACING_ENABLED"] = "true"
		if mon.Tracing.ExporterProtocol != "" {
			data["TRACING_EXPORTER_PROTOCOL"] = mon.Tracing.ExporterProtocol
		}
		if mon.Tracing.ServiceName != "" {
			data["TRACING_SERVICE_NAME"] = mon.Tracing.ServiceName
		}
		if mon.Tracing.ServiceNamespace != "" {
			data["TRACING_SERVICE_NAMESPACE"] = mon.Tracing.ServiceNamespace
		}
		if mon.Tracing.SamplingRate > 0 {
			data["TRACING_SAMPLING_RATE"] = strconv.FormatFloat(mon.Tracing.SamplingRate, 'f', -1, 64)
		}
	}

	return data, nil
}

func (n *EnvBuilder) BuildResponseHeadersEnv() (map[string]string, error) {
	data := make(map[string]string)

	if n.RateLimitService.Spec.ResponseHeaders != nil && n.RateLimitService.Spec.ResponseHeaders.Enabled {
		data["RESPONSE_HEADERS_ENABLED"] = "true"
	}

	return data, nil
}

func (n *EnvBuilder) BuildLoggingEnv() (map[string]string, error) {
	data := make(map[string]string)

	if n.RateLimitService.Spec.Logging == nil {
		return data, nil
	}

	if n.RateLimitService.Spec.Logging.Level != "" {
		data["LOG_LEVEL"] = n.RateLimitService.Spec.Logging.Level
	}

	if n.RateLimitService.Spec.Logging.Format != "" {
		data["LOG_FORMAT"] = n.RateLimitService.Spec.Logging.Format
	}

	return data, nil
}

func (n *EnvBuilder) BuildShadowModeEnv() (map[string]string, error) {
	data := make(map[string]string)

	if n.RateLimitService.Spec.ShadowMode {
		data["SHADOW_MODE"] = "true"
	}

	return data, nil
}
```

Then wire all into `BuildEnv()`. Replace the `Monitoring` block in `BuildEnv()` (lines 82-93) with:

```go
	if n.RateLimitService.Spec.Monitoring != nil {
		if n.RateLimitService.Spec.Monitoring.Type != "" {
			monEnv, err := n.BuildMonitoringEnv()
			if err != nil {
				return env, err
			}
			for key, value := range monEnv {
				env[key] = value
			}
		} else if n.RateLimitService.Spec.Monitoring.Enabled {
			// Legacy statsd sidecar path
			statsdEnv, err := n.BuildStatsdEnv()
			if err != nil {
				return env, err
			}
			for key, value := range statsdEnv {
				env[key] = value
			}
		}
	}

	if n.RateLimitService.Spec.ResponseHeaders != nil {
		rhEnv, err := n.BuildResponseHeadersEnv()
		if err != nil {
			return env, err
		}
		for key, value := range rhEnv {
			env[key] = value
		}
	}

	if n.RateLimitService.Spec.Logging != nil {
		logEnv, err := n.BuildLoggingEnv()
		if err != nil {
			return env, err
		}
		for key, value := range logEnv {
			env[key] = value
		}
	}

	if n.RateLimitService.Spec.ShadowMode {
		smEnv, err := n.BuildShadowModeEnv()
		if err != nil {
			return env, err
		}
		for key, value := range smEnv {
			env[key] = value
		}
	}
```

**Step 4: Run test to verify it passes**

Run: `cd /home/user/istio-ratelimit-operator && go test ./pkg/service/... -run "TestEnvBuilder_Build(Monitoring|ResponseHeaders|Logging|ShadowMode)Env" -v`
Expected: PASS

**Step 5: Run existing tests to verify backward compatibility**

Run: `cd /home/user/istio-ratelimit-operator && go test ./pkg/service/... -v`
Expected: ALL PASS (existing statsd tests still work via legacy path)

**Step 6: Commit**

```bash
git add pkg/service/configmap_env_builder.go pkg/service/configmap_env_builder_test.go
git commit -m "feat(service): add observability env vars (prometheus, tracing, logging, headers, shadow mode)"
```

---

### Task 12: Write failing test for Prometheus mode in DeploymentBuilder (skip statsd sidecar)

**Files:**
- Test: `pkg/service/deployment_builder_test.go`
- Modify: `pkg/service/deployment_builder.go`

**Step 1: Write the failing test**

Add to `pkg/service/deployment_builder_test.go`:

```go
func TestDeploymentBuilder_WithPrometheusMonitoring(t *testing.T) {
	setting := settings.Settings{
		RateLimitServiceImage: "ratelimit:latest",
		StatsdExporterImage:   "statsd:latest",
	}

	rateLimitService := v1alpha1.RateLimitService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "prom-test",
			Namespace: "default",
		},
		Spec: v1alpha1.RateLimitServiceSpec{
			Monitoring: &v1alpha1.RateLimitServiceSpec_Monitoring{
				Enabled: true,
				Type:    "prometheus",
			},
		},
	}

	deployment, err := service.NewDeploymentBuilder(setting).
		SetRateLimitService(rateLimitService).
		Build()

	assert.NoError(t, err)
	// Prometheus mode: only 1 container (no statsd sidecar)
	assert.Len(t, deployment.Spec.Template.Spec.Containers, 1)
	assert.Equal(t, "prom-test", deployment.Spec.Template.Spec.Containers[0].Name)
}
```

**Step 2: Run test to verify it fails**

Run: `cd /home/user/istio-ratelimit-operator && go test ./pkg/service/... -run TestDeploymentBuilder_WithPrometheusMonitoring -v`
Expected: FAIL — has 2 containers (statsd sidecar is added because `Enabled: true`)

**Step 3: Guard the statsd sidecar block with a type check**

In `pkg/service/deployment_builder.go`, change the monitoring block (around line 115) from:

```go
	if n.RateLimitService.Spec.Monitoring != nil {
		if n.RateLimitService.Spec.Monitoring.Enabled {
```

to:

```go
	if n.RateLimitService.Spec.Monitoring != nil {
		if n.RateLimitService.Spec.Monitoring.Enabled && n.RateLimitService.Spec.Monitoring.Type == "" {
```

This means: only add the statsd sidecar when `type` is unset (legacy behavior). When `type` is `"prometheus"` or `"dogstatsd"`, the ratelimit binary exposes metrics natively — no sidecar needed.

**Step 4: Run test to verify it passes**

Run: `cd /home/user/istio-ratelimit-operator && go test ./pkg/service/... -run TestDeploymentBuilder_WithPrometheusMonitoring -v`
Expected: PASS

**Step 5: Verify legacy statsd test still passes**

Run: `cd /home/user/istio-ratelimit-operator && go test ./pkg/service/... -run TestDeploymentBuilder_WithMonitoring -v`
Expected: PASS (legacy path uses `Enabled: true` with no Type, so sidecar is still added)

**Step 6: Commit**

```bash
git add pkg/service/deployment_builder.go pkg/service/deployment_builder_test.go
git commit -m "feat(service): skip statsd sidecar when monitoring.type is set"
```

---

### Task 13: Run make manifests for Phase 2 and verify all tests pass

**Step 1: Run manifests generation**

Run: `cd /home/user/istio-ratelimit-operator && make manifests`
Expected: SUCCESS

**Step 2: Run all tests**

Run: `cd /home/user/istio-ratelimit-operator && go test ./pkg/... -v`
Expected: ALL PASS

**Step 3: Commit**

```bash
git add config/ api/
git commit -m "chore: regenerate CRD manifests for Phase 2"
```

---

## Phase 3 — Descriptor Features

### Task 14: Add limit.unlimited and detailed_metric to GlobalRateLimit CRD

**Files:**
- Modify: `api/v1alpha1/globalratelimit_types.go:38-41`

**Step 1: Add `Unlimited` to `GlobalRateLimit_Limit`**

```go
type GlobalRateLimit_Limit struct {
	Unit            string `json:"unit,omitempty" yaml:"unit,omitempty"`
	RequestsPerUnit int    `json:"requests_per_unit,omitempty" yaml:"requests_per_unit,omitempty"`
	Unlimited       bool   `json:"unlimited,omitempty" yaml:"unlimited,omitempty"`
}
```

**Step 2: Add `DetailedMetric` to `GlobalRateLimitSpec`**

```go
type GlobalRateLimitSpec struct {
	Config         string                    `json:"config"`
	Selector       GlobalRateLimitSelector   `json:"selector"`
	Matcher        []*GlobalRateLimit_Action `json:"matcher"`
	ShadowMode     bool                      `json:"shadow_mode,omitempty"`
	Limit          *GlobalRateLimit_Limit    `json:"limit,omitempty"`
	Identifier     *string                   `json:"identifier,omitempty"`
	DetailedMetric bool                      `json:"detailed_metric,omitempty"`
}
```

**Step 3: Verify it compiles**

Run: `cd /home/user/istio-ratelimit-operator && go build ./api/...`
Expected: SUCCESS

**Step 4: Run make generate**

Run: `cd /home/user/istio-ratelimit-operator && make generate`
Expected: SUCCESS

**Step 5: Commit**

```bash
git add api/v1alpha1/globalratelimit_types.go api/v1alpha1/zz_generated.deepcopy.go
git commit -m "feat(api): add limit.unlimited and detailed_metric to GlobalRateLimit"
```

---

### Task 15: Update RateLimit_Service_Descriptor type

**Files:**
- Modify: `pkg/types/ratelimit.go:13-19`

**Step 1: Add `DetailedMetric` field to `RateLimit_Service_Descriptor`**

```go
type RateLimit_Service_Descriptor struct {
	Key            string                         `json:"key,omitempty" yaml:"key,omitempty"`
	Value          string                         `json:"value,omitempty" yaml:"value,omitempty"`
	ShadowMode     bool                           `json:"shadow_mode,omitempty" yaml:"shadow_mode,omitempty"`
	DetailedMetric bool                           `json:"detailed_metric,omitempty" yaml:"detailed_metric,omitempty"`
	RateLimit      v1alpha1.GlobalRateLimit_Limit `json:"rate_limit,omitempty" yaml:"rate_limit,omitempty"`
	Descriptors    []RateLimit_Service_Descriptor `json:"descriptors,omitempty" yaml:"descriptors,omitempty"`
}
```

**Step 2: Verify it compiles**

Run: `cd /home/user/istio-ratelimit-operator && go build ./pkg/...`
Expected: SUCCESS

**Step 3: Commit**

```bash
git add pkg/types/ratelimit.go
git commit -m "feat(types): add DetailedMetric to RateLimit_Service_Descriptor"
```

---

### Task 16: Write failing tests for unlimited and detailed_metric in configmap_config_builder

**Files:**
- Test: `pkg/service/configmap_config_builder_test.go`
- Modify: `pkg/service/configmap_config_builder.go`

**Step 1: Write the failing test**

Add to `pkg/service/configmap_config_builder_test.go`:

```go
func TestNewRateLimitDescriptor_Unlimited(t *testing.T) {
	globalRateLimitList := []v1alpha1.GlobalRateLimit{
		{
			Spec: v1alpha1.GlobalRateLimitSpec{
				Matcher: []*v1alpha1.GlobalRateLimit_Action{
					{
						GenericKey: &v1alpha1.GlobalRateLimit_Action_GenericKey{
							DescriptorValue: "internal-service",
						},
					},
				},
				Limit: &v1alpha1.GlobalRateLimit_Limit{
					Unlimited: true,
				},
				DetailedMetric: true,
			},
		},
	}

	descriptors, err := service.NewRateLimitDescriptor(globalRateLimitList)
	assert.NoError(t, err)
	assert.Len(t, descriptors, 1)
	assert.True(t, descriptors[0].RateLimit.Unlimited)
	assert.True(t, descriptors[0].DetailedMetric)
}
```

**Step 2: Run test to verify it fails**

Run: `cd /home/user/istio-ratelimit-operator && go test ./pkg/service/... -run TestNewRateLimitDescriptor_Unlimited -v`
Expected: FAIL — `DetailedMetric` is always false (not passed through)

**Step 3: Pass detailed_metric through the descriptor builder**

In `pkg/service/configmap_config_builder.go`, update `NewRateLimitDescriptorFromGlobalRateLimit` (line 108) to pass through `DetailedMetric`:

```go
func NewRateLimitDescriptorFromGlobalRateLimit(globalRateLimit v1alpha1.GlobalRateLimit) ([]types.RateLimit_Service_Descriptor, error) {
	var descriptor []types.RateLimit_Service_Descriptor
	var sanitizeMatchers []*v1alpha1.GlobalRateLimit_Action

	for _, matcher := range globalRateLimit.Spec.Matcher {
		if matcher.RequestHeaders != nil || matcher.GenericKey != nil || matcher.HeaderValueMatch != nil || matcher.RemoteAddress != nil {
			sanitizeMatchers = append(sanitizeMatchers, matcher)
			continue
		}
	}

	if len(sanitizeMatchers) == 0 {
		return descriptor, nil
	}

	descriptor, err := NewRateLimitDescriptorFromMatcher(sanitizeMatchers, globalRateLimit.Spec.Limit, globalRateLimit.Spec.ShadowMode)
	if err != nil {
		return descriptor, err
	}

	// Apply detailed_metric to leaf descriptors
	if globalRateLimit.Spec.DetailedMetric {
		applyDetailedMetric(descriptor)
	}

	return descriptor, nil
}

func applyDetailedMetric(descriptors []types.RateLimit_Service_Descriptor) {
	for i := range descriptors {
		if len(descriptors[i].Descriptors) == 0 {
			// Leaf descriptor — apply detailed_metric here
			descriptors[i].DetailedMetric = true
		} else {
			applyDetailedMetric(descriptors[i].Descriptors)
		}
	}
}
```

The `Unlimited` field flows automatically because `GlobalRateLimit_Limit` is copied directly into `descriptor[0].RateLimit` in `NewRateLimitDescriptorFromMatcher`. Since we added `Unlimited bool` to the struct, it's already passed through.

**Step 4: Run test to verify it passes**

Run: `cd /home/user/istio-ratelimit-operator && go test ./pkg/service/... -run TestNewRateLimitDescriptor_Unlimited -v`
Expected: PASS

**Step 5: Run all existing configmap_config tests**

Run: `cd /home/user/istio-ratelimit-operator && go test ./pkg/service/... -run TestNewRateLimitDescriptor -v`
Expected: ALL PASS

**Step 6: Commit**

```bash
git add pkg/service/configmap_config_builder.go pkg/service/configmap_config_builder_test.go
git commit -m "feat(service): pass unlimited and detailed_metric through descriptor builder"
```

---

### Task 17: Run make manifests for Phase 3 and verify all tests pass

**Step 1: Run manifests generation**

Run: `cd /home/user/istio-ratelimit-operator && make manifests`
Expected: SUCCESS

**Step 2: Run all tests**

Run: `cd /home/user/istio-ratelimit-operator && go test ./pkg/... -v`
Expected: ALL PASS

**Step 3: Commit**

```bash
git add config/ api/
git commit -m "chore: regenerate CRD manifests for Phase 3"
```

---

## Phase 4 — Backend Enhancements

### Task 18: Add Redis pool, timeout, perSecond, cacheKeyPrefix CRD types

**Files:**
- Modify: `api/v1alpha1/ratelimitservice_types.go`

**Step 1: Add new types**

```go
type RateLimitServiceSpec_Backend_Redis_Pool struct {
	Size                int    `json:"size,omitempty"`
	OnEmptyBehavior     string `json:"onEmptyBehavior,omitempty"`
	OnEmptyWaitDuration string `json:"onEmptyWaitDuration,omitempty"`
}

type RateLimitServiceSpec_Backend_Redis_PerSecond struct {
	Enabled       bool                                    `json:"enabled,omitempty"`
	URL           string                                  `json:"url,omitempty"`
	AuthSecretRef *SecretKeyRef                           `json:"authSecretRef,omitempty"`
	TLS           *RateLimitServiceSpec_Backend_Redis_TLS `json:"tls,omitempty"`
	Pool          *RateLimitServiceSpec_Backend_Redis_Pool `json:"pool,omitempty"`
}
```

**Step 2: Add new fields to `RateLimitServiceSpec_Backend_Redis`**

```go
type RateLimitServiceSpec_Backend_Redis struct {
	Type                           string                                           `json:"type,omitempty"`
	URL                            string                                           `json:"url,omitempty"`
	Auth                           string                                           `json:"auth,omitempty"`
	AuthSecretRef                  *SecretKeyRef                                    `json:"authSecretRef,omitempty"`
	Config                         *RateLimitServiceSpec_Backend_Redis_Config       `json:"config,omitempty"`
	TLS                            *RateLimitServiceSpec_Backend_Redis_TLS          `json:"tls,omitempty"`
	SentinelAuth                   *RateLimitServiceSpec_Backend_Redis_SentinelAuth `json:"sentinelAuth,omitempty"`
	Pool                           *RateLimitServiceSpec_Backend_Redis_Pool         `json:"pool,omitempty"`
	Timeout                        *string                                          `json:"timeout,omitempty"`
	HealthCheckActiveConnection    bool                                             `json:"healthCheckActiveConnection,omitempty"`
	PerSecond                      *RateLimitServiceSpec_Backend_Redis_PerSecond    `json:"perSecond,omitempty"`
}
```

**Step 3: Add new fields to `RateLimitServiceSpec_Backend`**

```go
type RateLimitServiceSpec_Backend struct {
	Redis                                *RateLimitServiceSpec_Backend_Redis `json:"redis,omitempty"`
	CacheKeyPrefix                       string                              `json:"cacheKeyPrefix,omitempty"`
	StopCacheKeyIncrementWhenOverlimit   bool                                `json:"stopCacheKeyIncrementWhenOverlimit,omitempty"`
}
```

**Step 4: Verify it compiles**

Run: `cd /home/user/istio-ratelimit-operator && go build ./api/...`
Expected: SUCCESS

**Step 5: Run make generate**

Run: `cd /home/user/istio-ratelimit-operator && make generate`
Expected: SUCCESS

**Step 6: Commit**

```bash
git add api/v1alpha1/ratelimitservice_types.go api/v1alpha1/zz_generated.deepcopy.go
git commit -m "feat(api): add Redis pool, timeout, perSecond, cacheKeyPrefix types"
```

---

### Task 19: Write failing tests for Phase 4 Redis env vars in EnvBuilder

**Files:**
- Test: `pkg/service/configmap_env_builder_test.go`
- Modify: `pkg/service/configmap_env_builder.go`

**Step 1: Write the failing tests**

Add to `pkg/service/configmap_env_builder_test.go`:

```go
func TestEnvBuilder_BuildRedisEnv_WithPool(t *testing.T) {
	builder := &service.EnvBuilder{
		RateLimitService: v1alpha1.RateLimitService{
			Spec: v1alpha1.RateLimitServiceSpec{
				Backend: &v1alpha1.RateLimitServiceSpec_Backend{
					Redis: &v1alpha1.RateLimitServiceSpec_Backend_Redis{
						Type: "single",
						URL:  "redis:6379",
						Pool: &v1alpha1.RateLimitServiceSpec_Backend_Redis_Pool{
							Size: 10,
						},
						Timeout: strPtr("2s"),
						HealthCheckActiveConnection: true,
					},
				},
			},
		},
	}

	env, err := builder.BuildRedisEnv()
	assert.NoError(t, err)
	assert.Equal(t, "10", env["REDIS_POOL_SIZE"])
	assert.Equal(t, "2s", env["REDIS_TIMEOUT"])
	assert.Equal(t, "true", env["REDIS_HEALTH_CHECK_ACTIVE_CONNECTION"])
}

func TestEnvBuilder_BuildRedisEnv_WithPerSecond(t *testing.T) {
	builder := &service.EnvBuilder{
		RateLimitService: v1alpha1.RateLimitService{
			Spec: v1alpha1.RateLimitServiceSpec{
				Backend: &v1alpha1.RateLimitServiceSpec_Backend{
					Redis: &v1alpha1.RateLimitServiceSpec_Backend_Redis{
						Type: "single",
						URL:  "redis:6379",
						PerSecond: &v1alpha1.RateLimitServiceSpec_Backend_Redis_PerSecond{
							Enabled: true,
							URL:     "redis-ps:6379",
							TLS: &v1alpha1.RateLimitServiceSpec_Backend_Redis_TLS{
								Enabled:   true,
								SecretRef: "redis-ps-tls",
							},
							Pool: &v1alpha1.RateLimitServiceSpec_Backend_Redis_Pool{
								Size: 5,
							},
						},
					},
				},
			},
		},
	}

	env, err := builder.BuildRedisEnv()
	assert.NoError(t, err)
	assert.Equal(t, "true", env["REDIS_PERSECOND"])
	assert.Equal(t, "redis-ps:6379", env["REDIS_PERSECOND_URL"])
	assert.Equal(t, "tcp", env["REDIS_PERSECOND_SOCKET_TYPE"])
	assert.Equal(t, "true", env["REDIS_PERSECOND_TLS"])
	assert.Equal(t, "/tls/redis-persecond/ca.crt", env["REDIS_PERSECOND_TLS_CACERT"])
	assert.Equal(t, "/tls/redis-persecond/tls.crt", env["REDIS_PERSECOND_TLS_CLIENT_CERT"])
	assert.Equal(t, "/tls/redis-persecond/tls.key", env["REDIS_PERSECOND_TLS_CLIENT_KEY"])
	assert.Equal(t, "5", env["REDIS_PERSECOND_POOL_SIZE"])
}

func TestEnvBuilder_BuildBackendEnv(t *testing.T) {
	builder := &service.EnvBuilder{
		RateLimitService: v1alpha1.RateLimitService{
			Spec: v1alpha1.RateLimitServiceSpec{
				Backend: &v1alpha1.RateLimitServiceSpec_Backend{
					Redis: &v1alpha1.RateLimitServiceSpec_Backend_Redis{
						Type: "single",
						URL:  "redis:6379",
					},
					CacheKeyPrefix:                     "my-svc",
					StopCacheKeyIncrementWhenOverlimit: true,
				},
			},
		},
	}

	env, err := builder.BuildBackendEnv()
	assert.NoError(t, err)
	assert.Equal(t, "my-svc", env["CACHE_KEY_PREFIX"])
	assert.Equal(t, "true", env["STOP_CACHE_KEY_INCREMENT_WHEN_OVERLIMIT"])
}
```

**Step 2: Run test to verify it fails**

Run: `cd /home/user/istio-ratelimit-operator && go test ./pkg/service/... -run "TestEnvBuilder_Build(RedisEnv_WithPool|RedisEnv_WithPerSecond|BackendEnv)" -v`
Expected: FAIL

**Step 3: Implement the new env vars**

Add to `BuildRedisEnv()` in `pkg/service/configmap_env_builder.go` (after the TLS block):

```go
	if n.RateLimitService.Spec.Backend.Redis.Pool != nil {
		if n.RateLimitService.Spec.Backend.Redis.Pool.Size > 0 {
			data["REDIS_POOL_SIZE"] = strconv.Itoa(n.RateLimitService.Spec.Backend.Redis.Pool.Size)
		}
	}

	if n.RateLimitService.Spec.Backend.Redis.Timeout != nil {
		data["REDIS_TIMEOUT"] = *n.RateLimitService.Spec.Backend.Redis.Timeout
	}

	if n.RateLimitService.Spec.Backend.Redis.HealthCheckActiveConnection {
		data["REDIS_HEALTH_CHECK_ACTIVE_CONNECTION"] = "true"
	}

	if n.RateLimitService.Spec.Backend.Redis.PerSecond != nil {
		ps := n.RateLimitService.Spec.Backend.Redis.PerSecond
		if ps.Enabled {
			data["REDIS_PERSECOND"] = "true"
			data["REDIS_PERSECOND_SOCKET_TYPE"] = "tcp"
			if ps.URL != "" {
				data["REDIS_PERSECOND_URL"] = ps.URL
			}
			if ps.TLS != nil && ps.TLS.Enabled {
				data["REDIS_PERSECOND_TLS"] = "true"
				if ps.TLS.SecretRef != "" {
					data["REDIS_PERSECOND_TLS_CACERT"] = "/tls/redis-persecond/ca.crt"
					data["REDIS_PERSECOND_TLS_CLIENT_CERT"] = "/tls/redis-persecond/tls.crt"
					data["REDIS_PERSECOND_TLS_CLIENT_KEY"] = "/tls/redis-persecond/tls.key"
				}
			}
			if ps.Pool != nil && ps.Pool.Size > 0 {
				data["REDIS_PERSECOND_POOL_SIZE"] = strconv.Itoa(ps.Pool.Size)
			}
		}
	}
```

Add new method `BuildBackendEnv()`:

```go
func (n *EnvBuilder) BuildBackendEnv() (map[string]string, error) {
	data := make(map[string]string)

	if n.RateLimitService.Spec.Backend == nil {
		return data, nil
	}

	if n.RateLimitService.Spec.Backend.CacheKeyPrefix != "" {
		data["CACHE_KEY_PREFIX"] = n.RateLimitService.Spec.Backend.CacheKeyPrefix
	}

	if n.RateLimitService.Spec.Backend.StopCacheKeyIncrementWhenOverlimit {
		data["STOP_CACHE_KEY_INCREMENT_WHEN_OVERLIMIT"] = "true"
	}

	return data, nil
}
```

Wire `BuildBackendEnv()` into `BuildEnv()`:

```go
	if n.RateLimitService.Spec.Backend != nil {
		backendEnv, err := n.BuildBackendEnv()
		if err != nil {
			return env, err
		}
		for key, value := range backendEnv {
			env[key] = value
		}
	}
```

**Step 4: Run test to verify it passes**

Run: `cd /home/user/istio-ratelimit-operator && go test ./pkg/service/... -run "TestEnvBuilder_Build(RedisEnv_WithPool|RedisEnv_WithPerSecond|BackendEnv)" -v`
Expected: PASS

**Step 5: Commit**

```bash
git add pkg/service/configmap_env_builder.go pkg/service/configmap_env_builder_test.go
git commit -m "feat(service): add Redis pool, timeout, perSecond, cacheKeyPrefix env vars"
```

---

### Task 20: Write failing test for perSecond TLS volume mounts in DeploymentBuilder

**Files:**
- Test: `pkg/service/deployment_builder_test.go`
- Modify: `pkg/service/deployment_builder.go`

**Step 1: Write the failing test**

Add to `pkg/service/deployment_builder_test.go`:

```go
func TestDeploymentBuilder_WithPerSecondTLS(t *testing.T) {
	setting := settings.Settings{
		RateLimitServiceImage: "ratelimit:latest",
		StatsdExporterImage:   "statsd:latest",
	}

	rateLimitService := v1alpha1.RateLimitService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "ps-tls-test",
			Namespace: "default",
		},
		Spec: v1alpha1.RateLimitServiceSpec{
			Backend: &v1alpha1.RateLimitServiceSpec_Backend{
				Redis: &v1alpha1.RateLimitServiceSpec_Backend_Redis{
					Type: "single",
					URL:  "redis:6379",
					PerSecond: &v1alpha1.RateLimitServiceSpec_Backend_Redis_PerSecond{
						Enabled: true,
						URL:     "redis-ps:6379",
						TLS: &v1alpha1.RateLimitServiceSpec_Backend_Redis_TLS{
							Enabled:   true,
							SecretRef: "redis-ps-tls-secret",
						},
					},
				},
			},
		},
	}

	deployment, err := service.NewDeploymentBuilder(setting).
		SetRateLimitService(rateLimitService).
		Build()

	assert.NoError(t, err)

	foundVolume := false
	for _, v := range deployment.Spec.Template.Spec.Volumes {
		if v.Name == "redis-persecond-tls" {
			foundVolume = true
			assert.Equal(t, "redis-ps-tls-secret", v.VolumeSource.Secret.SecretName)
		}
	}
	assert.True(t, foundVolume, "redis-persecond-tls volume not found")

	foundMount := false
	for _, vm := range deployment.Spec.Template.Spec.Containers[0].VolumeMounts {
		if vm.Name == "redis-persecond-tls" {
			foundMount = true
			assert.Equal(t, "/tls/redis-persecond/", vm.MountPath)
			assert.True(t, vm.ReadOnly)
		}
	}
	assert.True(t, foundMount, "redis-persecond-tls volume mount not found")
}
```

**Step 2: Run test to verify it fails**

Run: `cd /home/user/istio-ratelimit-operator && go test ./pkg/service/... -run TestDeploymentBuilder_WithPerSecondTLS -v`
Expected: FAIL

**Step 3: Implement perSecond TLS volume mounts**

Add to `pkg/service/deployment_builder.go` in `Build()` after the Redis TLS volume mount block:

```go
	if n.RateLimitService.Spec.Backend != nil && n.RateLimitService.Spec.Backend.Redis != nil &&
		n.RateLimitService.Spec.Backend.Redis.PerSecond != nil && n.RateLimitService.Spec.Backend.Redis.PerSecond.TLS != nil &&
		n.RateLimitService.Spec.Backend.Redis.PerSecond.TLS.SecretRef != "" {
		deployment.Spec.Template.Spec.Volumes = append(deployment.Spec.Template.Spec.Volumes, corev1.Volume{
			Name: "redis-persecond-tls",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: n.RateLimitService.Spec.Backend.Redis.PerSecond.TLS.SecretRef,
				},
			},
		})
		deployment.Spec.Template.Spec.Containers[0].VolumeMounts = append(
			deployment.Spec.Template.Spec.Containers[0].VolumeMounts,
			corev1.VolumeMount{
				Name:      "redis-persecond-tls",
				MountPath: "/tls/redis-persecond/",
				ReadOnly:  true,
			},
		)
	}
```

**Step 4: Run test to verify it passes**

Run: `cd /home/user/istio-ratelimit-operator && go test ./pkg/service/... -run TestDeploymentBuilder_WithPerSecondTLS -v`
Expected: PASS

**Step 5: Commit**

```bash
git add pkg/service/deployment_builder.go pkg/service/deployment_builder_test.go
git commit -m "feat(service): add perSecond Redis TLS volume mounts to DeploymentBuilder"
```

---

### Task 21: Run make manifests for Phase 4, full test run, and final verification

**Step 1: Run manifests generation**

Run: `cd /home/user/istio-ratelimit-operator && make manifests`
Expected: SUCCESS

**Step 2: Run ALL tests across the entire project**

Run: `cd /home/user/istio-ratelimit-operator && go test ./... -v`
Expected: ALL PASS

**Step 3: Verify build**

Run: `cd /home/user/istio-ratelimit-operator && make build`
Expected: SUCCESS — `bin/manager` is built

**Step 4: Commit CRD manifests**

```bash
git add config/ api/
git commit -m "chore: regenerate CRD manifests for Phase 4"
```

---

## Summary of All Files Modified

| File | Phases | Purpose |
|------|--------|---------|
| `api/v1alpha1/ratelimitservice_types.go` | 1,2,4 | New CRD types for all RateLimitService features |
| `api/v1alpha1/globalratelimit_types.go` | 3 | `limit.unlimited`, `detailed_metric` |
| `api/v1alpha1/zz_generated.deepcopy.go` | 1,2,3,4 | Auto-regenerated by `make generate` |
| `pkg/types/ratelimit.go` | 3 | `DetailedMetric` on descriptor type |
| `pkg/service/configmap_env_builder.go` | 1,2,4 | New env var methods: server, monitoring, logging, backend |
| `pkg/service/configmap_env_builder_test.go` | 1,2,4 | TDD tests for all new env vars |
| `pkg/service/configmap_config_builder.go` | 3 | Pass `unlimited`/`detailed_metric` through descriptors |
| `pkg/service/configmap_config_builder_test.go` | 3 | TDD test for unlimited/detailed_metric |
| `pkg/service/deployment_builder.go` | 1,2,4 | TLS volumes, scheduling fields, prometheus guard |
| `pkg/service/deployment_builder_test.go` | 1,2,4 | TDD tests for all deployment changes |
| `config/crd/bases/*.yaml` | 1,2,3,4 | Auto-regenerated CRD manifests |
