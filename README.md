<p>
  <img height="128" src="https://www.checklyhq.com/images/footer-logo.svg" align="right" />
  <h1>Checkly Operator</h1>
</p>

![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.1-4baaaa.svg)
![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/checkly/pulumi-checkly?label=Version)
[![Build & test](https://github.com/checkly/checkly-operator/actions/workflows/main-merge.yaml/badge.svg)](https://github.com/checkly/checkly-operator/actions/workflows/main-merge.yaml)

A kubernetes operator for [checklyhq.com](https://checklyhq.com). This operator can create Checkly checks and groups based of kubernetes CRDs and Ingress object annotations.

## Running in kubernetes

You can download the combined kubernetes resource manifest file from [the release page](https://github.com/checkly/checkly-operator/releases) for a specific version.

The `CHECKLY_API_KEY` and `CHECKLY_ACCOUNT_ID` should be environment variables attached to the container, make sure you create these as secrets and are available before you apply the resource manifests.

Once you have all the files, you can just apply the generated manifest files:
```bash
$ kubectl apply -f install-(version).yaml
```

### Ingress configuration

We support reading configuration from ingress resources, take a look at the [samples](config/samples/) directory.

> ⚠️  Groups have to exist before the Ingress object, so please at least apply 1 group CRD.

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

Make sure your current `kubectl` context is set to the appropriate kubernetes cluster where you want to test the operator, then run:

```bash
$ kubectl apply -f config/crd/bases/k8s.checklyhq.com_apichecks.yaml
$ kubectl apply -f config/crd/bases/k8s.checklyhq.com_groups.yaml
$ make run
```

If you update any of the types for the CRD, run
```bash
$ make manifests
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
Any checks created on [checklyhq.com](https://checklyhq.com) will be muted.

#### More Resources
Sources used for kick starting this project:
* https://sdk.operatorframework.io/docs/building-operators/golang/tutorial/
* https://kubernetes.io/blog/2021/06/21/writing-a-controller-for-pod-labels/
* https://github.com/checkly/checkly-go-sdk
* https://docs.okd.io/latest/operators/operator_sdk/golang/osdk-golang-tutorial.html
* https://sdk.operatorframework.io/docs/building-operators/golang/advanced-topics/#external-resources
* https://book.kubebuilder.io/cronjob-tutorial/writing-tests.html

## Questions
For questions and support please open a new  [discussion](https://github.com/checkly/checkly-operator/discussions). The issue list of this repo is exclusively for bug reports and feature/docs requests.

## Issues
Please make sure to respect issue requirements and choose the proper [issue template](https://github.com/checkly/checkly-operator/issues/new/choose) when opening an issue. Issues not conforming to the guidelines may be closed.

## License
[MIT](https://github.com/checkly/checkly-operator/blob/main/LICENSE)

<br>
<p align="center">
    <a href="https://www.imgarena.com" target="_blank"></a>
      <img width="100px" src="https://d2czwvv9f7qj0r.cloudfront.net/app/uploads/2021/11/img-logo-grey.svg.gzip" alt="Imgarena" />
    </a>
    <a href="https://checklyhq.com?utm_source=github&utm_medium=sponsor-logo-github&utm_campaign=checkly-operator" target="_blank">
      <img width="75px" src="https://www.checklyhq.com/images/text_racoon_logo.svg" alt="Checkly" />
    </a>
  <br>
  <br>
  <b><sub>From Imgarena & Checkly with ♥️</sub></b>
<p>
