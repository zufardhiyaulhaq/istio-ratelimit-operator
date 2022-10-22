# Contributing
By participating to this project, you agree to abide our [code of conduct](https://github.com/zufardhiyaulhaq/istio-ratelimit-operator/blob/master/.github/CODE_OF_CONDUCT.md).

## Development
For small things like fixing typos in documentation, you can [make edits through GitHub](https://help.github.com/articles/editing-files-in-another-user-s-repository/), which will handle forking and making a pull request (PR) for you. For anything bigger or more complex, you'll probably want to set up a development environment on your machine, a quick procedure for which is as folows:

### Setup your machine
Prerequisites:
- make
- [Go 1.17](https://golang.org/doc/install)
- [operator-sdk v1.11.0](https://sdk.operatorframework.io/)

Fork and clone **[istio-ratelimit-operator](https://github.com/zufardhiyaulhaq/istio-ratelimit-operator)** repository.

- deploy CRDs
```
kubectl apply -f config/crd/bases/
```

- Run istio-ratelimit-operator locally
```
make install run
```

### Submit a pull request
As you are ready with your code contribution, push your branch to your `istio-ratelimit-operator` fork and open a pull request against the **master** branch.
