# checkly-operator

[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=checkly_checkly-operator&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=checkly_checkly-operator) [![Build and push](https://github.com/checkly/checkly-operator/actions/workflows/main-merge.yaml/badge.svg)](https://github.com/checkly/checkly-operator/actions/workflows/main-merge.yaml) [![Code Coverage](https://sonarcloud.io/api/project_badges/measure?project=checkly_checkly-operator&metric=coverage)](https://sonarcloud.io/summary/new_code?id=checkly_checkly-operator)

A kubernetes operator for [checklyhq.com](https://checklyhq.com).

The operator can create checklyhq.com checks and groups based of kubernetes CRDs and Ingress object annotations.

## Development

Sources used for kick starting this project:
* https://sdk.operatorframework.io/docs/building-operators/golang/tutorial/
* https://kubernetes.io/blog/2021/06/21/writing-a-controller-for-pod-labels/
* https://github.com/checkly/checkly-go-sdk
* https://docs.okd.io/latest/operators/operator_sdk/golang/osdk-golang-tutorial.html
* https://sdk.operatorframework.io/docs/building-operators/golang/advanced-topics/#external-resources
* https://book.kubebuilder.io/cronjob-tutorial/writing-tests.html

## Running in kubernetes

You can download the combined kubernetes resource manifest file from [the release page](https://github.com/checkly/checkly-operator/releases) for a specific version.

The `CHECKLY_API_KEY` and `CHECKLY_ACCOUNT_ID` should be environment variables attached to the container, make sure you create these as secrets and are available before you apply the resource manifests.

Once you have all the files, you can just apply the generated manifest files:
```
kubectl apply -f install-(version).yaml
```

### Ingress configuration

We support reading configuration from ingress resources, take a look at the [samples](config/samples/) directory.

> **Warning**
> Groups have to exist before the Ingress object, so please at least apply 1 group CRD.

## Running locally

## Versions

We're using the following versions of packages:
* operator-sdk 1.22.0
* golang 1.18

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
Set your default kubernetes context to the cluster you desire to work on.

Modify the [config/samples/checkly_v1alpha1_apicheck.yaml](config/samples/checkly_v1alpha1_apicheck.yaml) or [config/samples/checkly_v1alpha1_group.yaml](config/samples/checkly_v1alpha1_group.yaml) and apply it.

#### Current settings

Any checks created on checklyhq.com will be muted.
