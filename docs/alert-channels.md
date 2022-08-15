# alert-channels

See the [official checkly docs](https://www.checklyhq.com/docs/alerting/) on what Alert channels are.

## Configuration options

The name of the Alert channel derives from the `metadata.name` of the created kubernetes resource.

We're supporting the email and OpsGenie configurations. You can not specify both in a config as each alert channel can only have one channel, if you want to alert to multiple channels, create a resource for each and later reference them in the check group configuration.

### Email

You can send alerts to an email address of your liking, all you need to do is set the `spec.email.address` field.
Example:
```yaml
apiVersion: k8s.checklyhq.com/v1alpha1
kind: AlertChannel
metadata:
  name: checkly-operator-test-email
spec:
  email:
    address: "foo@bar.baz"
```

### OpsGenie

The OpsGenie integration requires an API key to work. See [docs](https://www.checklyhq.com/docs/integrations/opsgenie/) on how to get the OpsGenie API key and determine your region.

You have the option of saving this information either as a kubernetes secret or configmap resource, we recommend a secret.

Once the above information is available, here's an example on how to setup the integration via our CRD:
```yaml
apiVersion: k8s.checklyhq.com/v1alpha1
kind: AlertChannel
metadata:
  name: checkly-operator-test-opsgenie
spec:
  opsgenie:
    apikey:
      name: test-secret # Name of the secret or configmap which holds the API key
      namespace: default # Namespace of the secret or configmap
      fieldPath: "API_KEY" # Key inside the secret or configmap
     priority: "P3" # P1, P2, P3, P4, P5 are the options
     region: "EU" # Your OpsGenie region
```

## Referencing

You'll need to reference the name of the alert channel in the group check configuration. See [check-group](check-group.md) for more details.
