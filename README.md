# checkly-operator

## Development

Sources used for kick starting this project:
* https://sdk.operatorframework.io/docs/building-operators/golang/tutorial/
* https://kubernetes.io/blog/2021/06/21/writing-a-controller-for-pod-labels/
* https://github.com/checkly/checkly-go-sdk

### Running locally

Make sure your current kubectl context is set to the appropriate kubernetes cluster where you want to test the operator, then run

```bash
kubectl apply -f config/crd/bases/checkly.imgarena.com_apichecks.yaml
make run
```

If you update any of the types for the CRD, run
```bash
make generate
make manifests
```
and re-apply the CRD.
