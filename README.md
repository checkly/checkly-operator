# checkly-operator

## Development

Sources used for kick starting this project:
* https://sdk.operatorframework.io/docs/building-operators/golang/tutorial/
* https://kubernetes.io/blog/2021/06/21/writing-a-controller-for-pod-labels/
* https://github.com/checkly/checkly-go-sdk
* https://docs.okd.io/latest/operators/operator_sdk/golang/osdk-golang-tutorial.html
* https://sdk.operatorframework.io/docs/building-operators/golang/advanced-topics/#external-resources

## Running locally

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
kubectl apply -f config/crd/bases/checkly.imgarena.com_apichecks.yaml
make run
```

If you update any of the types for the CRD, run
```bash
make manifests
```
and re-apply the CRD.

### Testing the controller

Set your default kubernetes context to the cluster you desire to work on.

Modify the [config/samples/checkly_v1alpha1_apicheck.yaml] and apply it.

#### Current settings

Any checks created on checkly.com will be muted and disabled.
