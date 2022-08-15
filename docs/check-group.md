# check-group

See the [official Checkly docs](https://www.checklyhq.com/docs/groups/) on what check groups are.

We're forcing API checks to be part of groups, this is an opinionated decision, if this does not work for your use case, please raise an Issue so we discuss options or alternate implementations.

## Configuration options

The name of the group is derived from the `metadata.name` of the kubernetes resource. `Group` resources are cluster scoped, meaning they need to be unique in a kubernetes cluster and they don't need a namespace definition.

### Labels

Any `metadata.labels` specified will be transformed into tags, ex. `environment: dev` label will be transformed to `environment:dev` tag, these tags then propagate to Prometheus metrics (if you're using [the checkly prometheus endpoint](https://www.checklyhq.com/docs/integrations/prometheus/)).

### Spec

The `spec` field accepts the following options:

| Option         | Details     | Default |
|--------------|-----------|------------|
| `locations` | Strings; A list of location where the checks should be running, for a list of locations see [doc](https://www.checklyhq.com/docs/monitoring/global-locations/).| `eu-west-1` |
| `alertchannel` | String; A list of alert channels which subscribe to the checks inside the group | none |

### Example

```yaml
apiVersion: k8s.checklyhq.com/v1alpha1
kind: Group
metadata:
  name: checkly-operator-test-group
  labels:
    environment: "local"
spec:
  locations:
    - eu-west-1
    - eu-west-2
  alertchannel:
    - checkly-operator-test-email
    - checkly-operator-test-opsgenie

```

## Referencing

You'll need to reference the name of the check group in the api check configuration. See [api-checks](api-checks.md) for more details.
