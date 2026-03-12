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
    - request_headers:
        header_name: "x-api-version"
  limit:
    unit: minute
    requests_per_unit: 100
```

**How it works:**
- `tenant-a` + `v1` = 100 req/min
- `tenant-a` + `v2` = 100 req/min (separate counter)
- `tenant-b` + `v1` = 100 req/min (separate counter)
- Each combination gets its own limit

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
    - remote_address: true
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
    - remote_address: true
    - request_headers:
        header_name: "x-tenant-id"
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

### Combining Matchers

You can combine multiple matchers in a single rule. The rate limit applies to the unique combination of all matcher values:

```yaml
matcher:
  - remote_address: true           # + client IP
  - request_headers:
      header_name: "x-tenant-id"   # + tenant ID
  - request_headers:
      header_name: "x-api-key"     # + API key
```

This creates a rate limit counter for each unique combination of: `IP + tenant + API key`.

---

## Common Patterns

### Pattern 1: Layered Rate Limits

Apply multiple rate limits with different scopes:

```yaml
# Global limit per IP (protection against abuse)
- matcher:
    - remote_address: true
  limit:
    unit: minute
    requests_per_unit: 1000

# Per-tenant limit (fair usage)
- matcher:
    - request_headers:
        header_name: "x-tenant-id"
  limit:
    unit: minute
    requests_per_unit: 500

# Per-tenant per-endpoint limit (resource protection)
- matcher:
    - request_headers:
        header_name: "x-tenant-id"
    - request_headers:
        header_name: ":path"
        header_value: "/api/v1/expensive-operation"
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
        header_value: "enterprise"
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
