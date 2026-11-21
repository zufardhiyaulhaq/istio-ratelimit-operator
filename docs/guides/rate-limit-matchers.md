# Rate Limit Matchers Guide

This guide explains the three types of rate limit matchers available for configuring rate limiting on your routes. Each matcher type serves different use cases for controlling traffic to your services.

## Overview

| Matcher Type | Use Case | Rate Limit Scope |
|-------------|----------|------------------|
| Request Header | Limit by header combinations | Per unique combination of header values |
| Unique Request Header | Limit specific header values | Per specific value you define |
| Client IP Address | Limit by client IP | Per unique client IP address |

---

## 1. Request Header Matcher

**What it does:** Applies rate limits based on the combination of request header values. Each unique combination of header values gets its own rate limit counter.

**Best for:** APIs where you want to limit requests per tenant, per API version, or per feature flag.

Every `request_headers` matcher needs `header_name` (the header to read) and `descriptor_key` (the name the ratelimit service uses for it). If the header is missing from a request, Envoy doesn't apply the rule to that request. With `skip_if_absent: true`, Envoy leaves out only that header and still applies the rest of the rule.

### Example: Limit by Tenant

You have a multi-tenant API and want to limit each tenant to 100 requests per minute.

```yaml
apiVersion: ratelimit.zufardhiyaulhaq.com/v1alpha1
kind: GlobalRateLimit
metadata:
  name: tenant-rate-limit
  namespace: istio-system
spec:
  config: "production-ratelimit"
  selector:
    vhost: "api.example.com:443"
  matcher:
    - request_headers:
        header_name: "x-tenant-id"
        descriptor_key: "tenant"
  limit:
    unit: minute
    requests_per_unit: 100
```

**How it works:**
- Tenant A (`x-tenant-id: tenant-a`) gets 100 req/min
- Tenant B (`x-tenant-id: tenant-b`) gets 100 req/min (separate counter)
- Each tenant has independent rate limits

### Example: Limit by Multiple Headers

You want to limit requests per tenant AND per API version.

```yaml
apiVersion: ratelimit.zufardhiyaulhaq.com/v1alpha1
kind: GlobalRateLimit
metadata:
  name: tenant-version-rate-limit
  namespace: istio-system
spec:
  config: "production-ratelimit"
  selector:
    vhost: "api.example.com:443"
  matcher:
    - request_headers:
        header_name: "x-tenant-id"
        descriptor_key: "tenant"
    - request_headers:
        header_name: "x-api-version"
        descriptor_key: "api_version"
  limit:
    unit: minute
    requests_per_unit: 100
```

**How it works:**
- `tenant-a` + `v1` = 100 req/min
- `tenant-a` + `v2` = 100 req/min (separate counter)
- `tenant-b` + `v1` = 100 req/min (separate counter)
- Each combination gets its own limit

### Example: Limit Only Matching Values

You only want to limit downloads under `/files/`, with each file getting its own limit. Set `value` to the values to limit. `*` matches any characters, anywhere in the value.

```yaml
apiVersion: ratelimit.zufardhiyaulhaq.com/v1alpha1
kind: GlobalRateLimit
metadata:
  name: file-download-rate-limit
  namespace: istio-system
spec:
  config: "production-ratelimit"
  selector:
    vhost: "api.example.com:443"
  matcher:
    - request_headers:
        header_name: ":path"
        descriptor_key: "path"
        value: "/files/*"
  limit:
    unit: minute
    requests_per_unit: 10
```

**How it works:**
- `/files/a.pdf` gets 10 req/min
- `/files/b.csv` gets 10 req/min (separate counter)
- `/api/users` isn't limited by this rule

### Example: Share One Limit Across Matching Values

You want all of `/files/*` together to get 10 requests per minute, whichever file is requested. Add `share_threshold: true`. It needs a `value` that contains `*`, and the API server rejects the rule otherwise.

```yaml
  matcher:
    - request_headers:
        header_name: ":path"
        descriptor_key: "path"
        value: "/files/*"
        share_threshold: true
  limit:
    unit: minute
    requests_per_unit: 10
```

**How it works:**
- `/files/a.pdf`, `/files/b.csv` and `/files/c.txt` share one counter of 10 req/min
- After 5 requests for `/files/a.pdf` and 5 for `/files/b.csv`, a request for `/files/c.txt` is limited

On clusters older than Kubernetes 1.25 the API server doesn't run this check. There, the operator leaves `share_threshold` out of the ratelimit config instead, so the other rules keep working.

### Example: Per-Value Metrics

By default, the ratelimit service's metric keys are built from the descriptor keys only, so all tenants are counted together in one metric. Add `value_to_metric: true` to put the header value in the metric name, so each tenant gets its own metrics.

```yaml
  matcher:
    - request_headers:
        header_name: "x-tenant-id"
        descriptor_key: "tenant"
        value_to_metric: true
  limit:
    unit: minute
    requests_per_unit: 100
```

**How it works:**
- Without it: one metric, `production.tenant`
- With it: `production.tenant_tenant-a`, `production.tenant_tenant-b`, and so on (here `production` is the domain of the `GlobalRateLimitConfig`)
- Every new value creates new metrics, so avoid it for headers with many distinct values, such as user IDs
- With `share_threshold`, the metric uses the wildcard prefix of `value` instead of each full value, since those values share one counter
- `remote_address` takes `value_to_metric: true` as well (`remote_address: {value_to_metric: true}`), but every client IP then gets its own metrics
- If the rule also sets `detailed_metric: true`, that wins on the last matcher and `value_to_metric` is ignored there

`value`, `share_threshold` and `value_to_metric` are only used by the ratelimit service. The operator doesn't send them to Envoy.

---

## 2. Unique Request Header Matcher

**What it does:** Applies rate limits only when a header matches specific criteria you define. Unlike Request Header matcher, this targets specific values rather than any value.

**Best for:** Applying different limits to specific customers, blocking/limiting known bad actors, or creating premium tiers.

### Matching Options

| Match Type | Description | Example |
|------------|-------------|---------|
| `exact_match` | Exact value match | `"premium"` |
| `regex_match` | Regex pattern (multiple values) | `"cust-123\|cust-456\|cust-789"` |
| `prefix_match` | Starts with | `"enterprise-"` |
| `suffix_match` | Ends with | `"-premium"` |
| `contains_match` | Contains substring | `"vip"` |
| `present_match` | Header exists (any value) | `true` |

### Example: Premium vs Standard Tier

You want to give premium customers 1000 req/min and standard customers 100 req/min.

```yaml
# Premium tier - 1000 req/min
apiVersion: ratelimit.zufardhiyaulhaq.com/v1alpha1
kind: GlobalRateLimit
metadata:
  name: premium-tier-limit
  namespace: istio-system
spec:
  config: "production-ratelimit"
  selector:
    vhost: "api.example.com:443"
  matcher:
    - header_value_match:
        descriptor_value: "premium-tier"
        headers:
          - name: "x-subscription-tier"
            exact_match: "premium"
  limit:
    unit: minute
    requests_per_unit: 1000
---
# Standard tier - 100 req/min
apiVersion: ratelimit.zufardhiyaulhaq.com/v1alpha1
kind: GlobalRateLimit
metadata:
  name: standard-tier-limit
  namespace: istio-system
spec:
  config: "production-ratelimit"
  selector:
    vhost: "api.example.com:443"
  matcher:
    - header_value_match:
        descriptor_value: "standard-tier"
        headers:
          - name: "x-subscription-tier"
            exact_match: "standard"
  limit:
    unit: minute
    requests_per_unit: 100
```

**How it works:**
- Requests with `x-subscription-tier: premium` = 1000 req/min
- Requests with `x-subscription-tier: standard` = 100 req/min
- Requests without this header = no rate limit (unless you add a catch-all)

### Example: Rate Limit Specific Customer

A specific customer is causing issues and you need to throttle them.

```yaml
apiVersion: ratelimit.zufardhiyaulhaq.com/v1alpha1
kind: GlobalRateLimit
metadata:
  name: throttle-problem-customer
  namespace: istio-system
spec:
  config: "production-ratelimit"
  selector:
    vhost: "api.example.com:443"
  matcher:
    - header_value_match:
        descriptor_value: "throttled-customer"
        headers:
          - name: "x-customer-id"
            exact_match: "cust-12345"
  limit:
    unit: minute
    requests_per_unit: 10  # Heavily throttled
```

### Example: Rate Limit Multiple Customers (Single Rule)

You want to apply the same rate limit to multiple specific customers in one GlobalRateLimit.

**Option 1: Using regex_match**

```yaml
apiVersion: ratelimit.zufardhiyaulhaq.com/v1alpha1
kind: GlobalRateLimit
metadata:
  name: throttle-multiple-customers
  namespace: istio-system
spec:
  config: "production-ratelimit"
  selector:
    vhost: "api.example.com:443"
  matcher:
    - header_value_match:
        descriptor_value: "throttled-customers"
        headers:
          - name: "x-customer-id"
            regex_match: "^(cust-123|cust-456|cust-789)$"
  limit:
    unit: minute
    requests_per_unit: 10
```

**Option 2: Using prefix_match**

Rate limit all customers with IDs starting with "trial-":

```yaml
apiVersion: ratelimit.zufardhiyaulhaq.com/v1alpha1
kind: GlobalRateLimit
metadata:
  name: throttle-trial-customers
  namespace: istio-system
spec:
  config: "production-ratelimit"
  selector:
    vhost: "api.example.com:443"
  matcher:
    - header_value_match:
        descriptor_value: "trial-customers"
        headers:
          - name: "x-customer-id"
            prefix_match: "trial-"
  limit:
    unit: minute
    requests_per_unit: 50
```

**Option 3: Using contains_match**

Rate limit all customers whose ID contains "demo":

```yaml
apiVersion: ratelimit.zufardhiyaulhaq.com/v1alpha1
kind: GlobalRateLimit
metadata:
  name: throttle-demo-customers
  namespace: istio-system
spec:
  config: "production-ratelimit"
  selector:
    vhost: "api.example.com:443"
  matcher:
    - header_value_match:
        descriptor_value: "demo-customers"
        headers:
          - name: "x-customer-id"
            contains_match: "demo"
  limit:
    unit: minute
    requests_per_unit: 30
```

**How it works:**
- All matching customers share the SAME rate limit counter
- `cust-123`, `cust-456`, `cust-789` together = 10 req/min total (not each!)
- If you want separate limits per customer, use **Request Header matcher** instead

---

## 3. Client IP Address Matcher

**What it does:** Applies rate limits based on the client's IP address. Each unique IP gets its own rate limit counter.

**Best for:** Protecting against DDoS, preventing brute force attacks, or limiting anonymous users.

### Example: Limit Requests Per IP

You want to limit each IP address to 60 requests per minute.

```yaml
apiVersion: ratelimit.zufardhiyaulhaq.com/v1alpha1
kind: GlobalRateLimit
metadata:
  name: per-ip-rate-limit
  namespace: istio-system
spec:
  config: "production-ratelimit"
  selector:
    vhost: "api.example.com:443"
  matcher:
    - remote_address: {}
  limit:
    unit: minute
    requests_per_unit: 60
```

**How it works:**
- IP `192.168.1.100` gets 60 req/min
- IP `192.168.1.101` gets 60 req/min (separate counter)
- Each client IP has independent rate limits

### Example: Combine IP with Header

You want to limit per IP per tenant (useful when tenants share IPs like corporate networks).

```yaml
apiVersion: ratelimit.zufardhiyaulhaq.com/v1alpha1
kind: GlobalRateLimit
metadata:
  name: ip-tenant-rate-limit
  namespace: istio-system
spec:
  config: "production-ratelimit"
  selector:
    vhost: "api.example.com:443"
  matcher:
    - remote_address: {}
    - request_headers:
        header_name: "x-tenant-id"
        descriptor_key: "tenant"
  limit:
    unit: minute
    requests_per_unit: 100
```

**How it works:**
- `192.168.1.100` + `tenant-a` = 100 req/min
- `192.168.1.100` + `tenant-b` = 100 req/min (separate counter)
- `192.168.1.101` + `tenant-a` = 100 req/min (separate counter)

---

## Quick Reference

### When to Use Each Matcher

| Scenario | Recommended Matcher | Match Type |
|----------|-------------------|------------|
| Limit per tenant/customer (each gets own limit) | Request Header | `request_headers` |
| Limit per API key (each gets own limit) | Request Header | `request_headers` |
| Different limits for premium/standard users | Header Value Match | `exact_match` |
| Throttle a specific bad actor | Header Value Match | `exact_match` |
| Throttle multiple specific customers (shared limit) | Header Value Match | `regex_match` |
| Limit all trial/demo customers | Header Value Match | `prefix_match` or `contains_match` |
| Protect against DDoS | Client IP Address | `remote_address` |
| Prevent brute force on login | Client IP Address | `remote_address` |
| Limit anonymous users | Client IP Address | `remote_address` |

### Important: Separate vs Shared Limits

| If you want... | Use |
|----------------|-----|
| Each customer has **own** rate limit counter | `request_headers` (any value) |
| Multiple customers **share** one rate limit counter | `header_value_match` with `regex_match` |
| All values matching a wildcard **share** one counter | `request_headers` with `value: "prefix*"` and `share_threshold: true` |

### Combining Matchers

You can combine multiple matchers in a single rule. The rate limit applies to the unique combination of all matcher values:

```yaml
matcher:
  - remote_address: {}             # + client IP
  - request_headers:
      header_name: "x-tenant-id"   # + tenant ID
      descriptor_key: "tenant"
  - request_headers:
      header_name: "x-api-key"     # + API key
      descriptor_key: "api_key"
```

This creates a rate limit counter for each unique combination of: `IP + tenant + API key`.

---

## Common Patterns

### Pattern 1: Layered Rate Limits

Apply multiple rate limits with different scopes:

```yaml
# Global limit per IP (protection against abuse)
- matcher:
    - remote_address: {}
  limit:
    unit: minute
    requests_per_unit: 1000

# Per-tenant limit (fair usage)
- matcher:
    - request_headers:
        header_name: "x-tenant-id"
        descriptor_key: "tenant"
  limit:
    unit: minute
    requests_per_unit: 500

# Per-tenant per-endpoint limit (resource protection)
- matcher:
    - request_headers:
        header_name: "x-tenant-id"
        descriptor_key: "tenant"
    - request_headers:
        header_name: ":path"
        descriptor_key: "path"
        # :path includes the query string, so match it with "*"
        # and share one counter across all query strings
        value: "/api/v1/expensive-operation*"
        share_threshold: true
  limit:
    unit: minute
    requests_per_unit: 10
```

### Pattern 2: Whitelist with Unlimited

Allow specific customers unlimited access using the `unlimited` field:

```yaml
apiVersion: ratelimit.zufardhiyaulhaq.com/v1alpha1
kind: GlobalRateLimit
metadata:
  name: enterprise-unlimited
spec:
  config: "production-ratelimit"
  selector:
    vhost: "api.example.com:443"
  matcher:
    - request_headers:
        header_name: "x-subscription-tier"
        descriptor_key: "subscription_tier"
        value: "enterprise"
  limit:
    unlimited: true  # No rate limit for enterprise customers
```

---

## Troubleshooting

### Rate limit not applying?

1. **Check header names are exact match** - Header names are case-sensitive
2. **Verify the selector matches your route** - Check `vhost` matches your Istio gateway
3. **Confirm GlobalRateLimitConfig exists** - The `config` field must reference a valid config

### All users sharing the same limit?

You're probably missing a unique identifier in your matcher. Add a header or IP matcher to create per-user counters.

### Need to see what's being rate limited?

Enable `detailed_metric: true` to get per-value metrics:

```yaml
spec:
  detailed_metric: true  # Adds metrics per unique header value
```

To choose which matchers add their value to the metrics, set `value_to_metric: true` on those `request_headers` or `remote_address` matchers instead. See [Per-Value Metrics](#example-per-value-metrics).
