# istio-ratelimit-operator

Istio ratelimit operator provide an easy way to configure Global or Local Ratelimit in Istio mesh. Istio ratelimit operator also support EnvoyFilter versioning!

![Version: 2.15.0](https://img.shields.io/badge/Version-2.15.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 2.15.0](https://img.shields.io/badge/AppVersion-2.15.0-informational?style=flat-square) [![made with Go](https://img.shields.io/badge/made%20with-Go-brightgreen)](http://golang.org) [![GitHub issues](https://img.shields.io/github/issues/zufardhiyaulhaq/istio-ratelimit-operator)](https://github.com/zufardhiyaulhaq/istio-ratelimit-operator/issues) [![GitHub pull requests](https://img.shields.io/github/issues-pr/zufardhiyaulhaq/istio-ratelimit-operator)](https://github.com/zufardhiyaulhaq/istio-ratelimit-operator/pulls)[![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/istio-ratelimit-operator)](https://artifacthub.io/packages/search?repo=istio-ratelimit-operator)

## Installation

To install the chart with the release name `my-istio-ratelimit-operator`:

```console
helm repo add istio-ratelimit-operator https://zufardhiyaulhaq.com/istio-ratelimit-operator/charts/releases/
helm install my-istio-ratelimit-operator istio-ratelimit-operator/istio-ratelimit-operator --version 2.15.0 --values values.yaml
```

## Usage
1. Apply Global ratelimit example
```console
kubectl apply -f examples/global/gateway/ratelimitservice/
```

2. Check Object
```
kubectl get GlobalRateLimitConfig
NAME                   AGE
istio-public-gateway   3m1s

kubectl get GlobalRateLimit
NAME                                          AGE
helloworld-zufardhiyaulhaq-com-bar-route      2m57s
helloworld-zufardhiyaulhaq-com-baz-route      2m56s
helloworld-zufardhiyaulhaq-com-corge-route    2m53s
helloworld-zufardhiyaulhaq-com-foo-route      2m57s
helloworld-zufardhiyaulhaq-com-garply-route   2m51s
helloworld-zufardhiyaulhaq-com-grault-route   2m52s
helloworld-zufardhiyaulhaq-com-quux-route     2m54s
helloworld-zufardhiyaulhaq-com-qux-route      2m55s

kubectl get RateLimitService
NAME                               AGE
public-gateway-ratelimit-service   2m33s
```

3. Check EnvoyFilter
```
kubectl get envoyfilter
NAME                                                                                                                                         AGE
helloworld-zufardhiyaulhaq-com-bar-route-1.8                                                                                                 3m7s
helloworld-zufardhiyaulhaq-com-bar-route-1.9                                                                                                 3m7s
helloworld-zufardhiyaulhaq-com-baz-route-1.8                                                                                                 3m7s
helloworld-zufardhiyaulhaq-com-baz-route-1.9                                                                                                 3m6s
helloworld-zufardhiyaulhaq-com-corge-route-1.8                                                                                               3m8s
helloworld-zufardhiyaulhaq-com-corge-route-1.9                                                                                               3m4s
helloworld-zufardhiyaulhaq-com-foo-route-1.8                                                                                                 3m8s
helloworld-zufardhiyaulhaq-com-foo-route-1.9                                                                                                 3m7s
helloworld-zufardhiyaulhaq-com-garply-route-1.8                                                                                              3m10s
helloworld-zufardhiyaulhaq-com-garply-route-1.9                                                                                              3m8s
helloworld-zufardhiyaulhaq-com-grault-route-1.8                                                                                              3m10s
helloworld-zufardhiyaulhaq-com-grault-route-1.9                                                                                              3m7s
helloworld-zufardhiyaulhaq-com-quux-route-1.8                                                                                                3m9s
helloworld-zufardhiyaulhaq-com-quux-route-1.9                                                                                                3m5s
helloworld-zufardhiyaulhaq-com-qux-route-1.8                                                                                                 3m11s
helloworld-zufardhiyaulhaq-com-qux-route-1.9                                                                                                 3m6s

istio-public-gateway-1.8                                                                                                                     3m11s
istio-public-gateway-1.9                                                                                                                     3m11s
```

4. Check Ratelimit
```
kubectl get service
NAME                               TYPE           CLUSTER-IP      EXTERNAL-IP      PORT(S)                                                           AGE
public-gateway-ratelimit-service   ClusterIP      10.32.214.174   <none>           8080/TCP,8081/TCP,6070/TCP                                        4m17s

kubectl get deployment
NAME                               READY   UP-TO-DATE   AVAILABLE   AGE
public-gateway-ratelimit-service   2/2     2            2           4m53s

kubectl get configmap
NAME                                          DATA   AGE
public-gateway-ratelimit-service-config       1      5m14s
public-gateway-ratelimit-service-config-env   4      5m14s

kubectl port-forward svc/public-gateway-ratelimit-service 6070:6070
curl http://127.0.0.1:6070/rlconfig

public-gateway.path.corge-route_corge-route: unit=HOUR requests_per_unit=120
public-gateway.path.quux-route_quux-route: unit=HOUR requests_per_unit=60
public-gateway.method.garply-route_garply-route: unit=HOUR requests_per_unit=120
public-gateway.method.path.bar-route_bar-route: unit=HOUR requests_per_unit=120
public-gateway.method.path.foo-route_foo-route: unit=HOUR requests_per_unit=60
public-gateway.method.machineid.qux-route_qux-route: unit=HOUR requests_per_unit=90
public-gateway.method.machineid.baz-route_baz-route: unit=HOUR requests_per_unit=90
public-gateway.method.grault-route_grault-route: unit=HOUR requests_per_unit=60
```

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| extraLabels | object | `{}` |  |
| operator.image | string | `"zufardhiyaulhaq/istio-ratelimit-operator"` |  |
| operator.replica | int | `1` |  |
| operator.tag | string | `"v2.15.0"` |  |
| resources.limits.cpu | string | `"512m"` |  |
| resources.limits.memory | string | `"512Mi"` |  |
| resources.requests.cpu | string | `"256m"` |  |
| resources.requests.memory | string | `"256Mi"` |  |
| serviceAccount.imagePullSecrets | list | `[]` |  |
| settings.ratelimitservice.image | string | `"envoyproxy/ratelimit:5e1be594"` |  |
| settings.statsdExporter.image | string | `"prom/statsd-exporter:v0.26.1"` |  |

## Supported Releases

| Operator Version | Istio Version |
|-----|------|
| 2.13.0 | <= 1.21.x |
| 2.12.0 | <= 1.20.x |
| 2.11.2 | <= 1.17.x |
| 2.9.0 | <= 1.16.x |
| 2.8.0 | <= 1.15.x |
| 2.5.1 | <= 1.13.x |
