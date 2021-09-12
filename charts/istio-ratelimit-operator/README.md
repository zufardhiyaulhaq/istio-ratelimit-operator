# istio-ratelimit-operator

Istio ratelimit operator provide an easy way to configure Global or Local Ratelimit in Istio mesh. Istio ratelimit operator also support EnvoyFilter versioning!

![Version: 1.0.0](https://img.shields.io/badge/Version-1.0.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 1.0.0](https://img.shields.io/badge/AppVersion-1.0.0-informational?style=flat-square) [![made with Go](https://img.shields.io/badge/made%20with-Go-brightgreen)](http://golang.org) [![Github master branch build](https://img.shields.io/github/workflow/status/zufardhiyaulhaq/istio-ratelimit-operator/Master)](https://github.com/zufardhiyaulhaq/istio-ratelimit-operator/actions/workflows/master.yml) [![GitHub issues](https://img.shields.io/github/issues/zufardhiyaulhaq/istio-ratelimit-operator)](https://github.com/zufardhiyaulhaq/istio-ratelimit-operator/issues) [![GitHub pull requests](https://img.shields.io/github/issues-pr/zufardhiyaulhaq/istio-ratelimit-operator)](https://github.com/zufardhiyaulhaq/istio-ratelimit-operator/pulls)

## Installing

To install the chart with the release name `my-release`:

```console
helm repo add zufardhiyaulhaq https://charts.zufardhiyaulhaq.com/
helm install my-release zufardhiyaulhaq/istio-ratelimit-operator --values values.yaml
```

## Usage
1. Apply Global ratelimit example
```console
kubectl apply -f examples/global/
```

2. Check Object
```
kubectl get GlobalRateLimitConfig
kubectl get GlobalRateLimit
```

3. Check EnvoyFilter
```
kubectl get envoyfilter
NAME                                            AGE
helloworld-zufardhiyaulhaq-dev-1.8              9m58s
helloworld-zufardhiyaulhaq-dev-1.9              9m54s
public-gateway-1.8                              14m
public-gateway-1.9                              14m
```

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| operator.image | string | `"zufardhiyaulhaq/istio-ratelimit-operator"` |  |
| operator.replica | int | `1` |  |
| operator.tag | string | `"v1.0.0"` |  |
| resources.limits.cpu | string | `"200m"` |  |
| resources.limits.memory | string | `"100Mi"` |  |
| resources.requests.cpu | string | `"100m"` |  |
| resources.requests.memory | string | `"20Mi"` |  |

