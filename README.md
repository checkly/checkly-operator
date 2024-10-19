# checkly-operator

[![Build and push](https://github.com/checkly/checkly-operator/actions/workflows/main-merge.yaml/badge.svg)](https://github.com/checkly/checkly-operator/actions/workflows/main-merge.yaml)

A kubernetes operator for [checklyhq.com](https://checklyhq.com).

The operator can create checklyhq.com checks, groups and alert channels based of kubernetes CRDs and Ingress object annotations.

## Documentation
Please see our [docs](docs/README.md) for more details on how to install and use the operator.

## Development
### direnv

We're using [direnv](https://direnv.net/) to manage environment variables for this project (or export them manually and you can skip this step). Make sure you generate a checkly API key and you get the account ID as well.

```
touch .envrc
echo "export CHECKLY_API_KEY=foorbarbaz" > .envrc
echo "export CHECKLY_ACCOUNT_ID=randomnumbers" >> .envrc
direnv allow .
```

### Makefile

Make sure your current kubectl context is set to the appropriate kubernetes cluster where you want to test the operator, then run

```bash
kubectl apply -f config/crd/bases/k8s.checklyhq.com_apichecks.yaml
kubectl apply -f config/crd/bases/k8s.checklyhq.com_groups.yaml
kubectl apply -f config/crd/bases/k8s.checklyhq.com_alertchannels.yaml
make run
```

If you update any of the types for the CRD, run
```bash
make manifests
```
and re-apply the CRD.

### Testing the controller

#### Unit and integration tests
* Make sure your kubectl context is set to your local k8s cluster
* Run `USE_EXISTING_CLUSTER=true make test`
* To see coverage run `go tool cover -html=cover.out`

#### Running locally
See [docs](docs/README.md) for details.

## Source material

Sources used for kick starting this project:
* https://sdk.operatorframework.io/docs/building-operators/golang/tutorial/
* https://kubernetes.io/blog/2021/06/21/writing-a-controller-for-pod-labels/
* https://github.com/checkly/checkly-go-sdk
* https://docs.okd.io/latest/operators/operator_sdk/golang/osdk-golang-tutorial.html
* https://sdk.operatorframework.io/docs/building-operators/golang/advanced-topics/#external-resources
* https://book.kubebuilder.io/cronjob-tutorial/writing-tests.html


### Versions

We're using the following versions of packages:
* operator-sdk 1.33.0
* golang 1.22

