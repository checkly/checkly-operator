# alert-channels

See the [official checkly docs](https://www.checklyhq.com/docs/alerting/) on what Alert channels are.

## Configuration options

The name of the Alert channel derives from the `metadata.name` of the created kubernetes resource.

We're supporting the email, OpsGenie and webhook configurations. You can not specify all in a config as each alert channel can only have one channel, if you want to alert to multiple channels, create a resource for each and later reference them in the check group configuration.

### Spec

| Option         | Details     | Default |
|--------------|-----------|------------|
| `sendRecovery` | Bool; Should recovery alerts be sent | none |
| `sendFailure` | Bool; Should failure alerts be sent | none |
| `sslExpiry` | Bool; Should ssl expiry check alerts be sent | none |
| `sendDegraded` | Bool; Should degraded alerts be sent | none |
| `sslExpiryThreshold` | int; At what moment in time to start alerting on SSL certificates. | none |
| `email.address` | string; Which email address should the alert be sent to | none |
| `opsgenie.apikey.name` | string; Name of the secret or configmap which holds the Opsgenie API key | none |
| `opsgenie.apikey.namespace` | string; Namespace of the secret or configmap | none |
| `opsgenie.apikey.fieldPath` | string; Key inside the secret or configmap | none |
| `webhook.name` | string; Name for the webhook | none |
| `webhook.url` | string; URL for the webhook | none |
| `webhook.webhookType` | string; TODO: can't determine what this is | none |
| `webhook.method` | string; HTTP type for the webhook (POST/GET/PUT/HEAD/DELETE/PATCH) | none |
| `webhook.template` | string; Template for webhook message | none |
| `webhook.name` | string; Name for the webhook | none |
| `webhook.webhookSecret.name` | string; Name of the secret or configmap which holds the Opsgenie API key | none |
| `webhook.webhookSecret.namespace` | string; Namespace of the secret or configmap | none |
| `webhook.webhookSecret.fieldPath` | string; Key inside the secret or configmap | none |
| `webhook.(*).headers.key` | string; Name for the header key | none |
| `webhook.(*).headers.value` | string; Value of the header | none |
| `webhook.(*).headers.locked` | bool; Is the header value visible in the checklyhq console | none |
| `webhook.(*).queryParameters.key` | string; Name for the query parameter key | none |
| `webhook.(*).queryParameters.value` | string; Value of the query parameter | none |
| `webhook.(*).queryParameters.locked` | bool; Is the query parameter value visible in the checklyhq console | none |

### Email

You can send alerts to an email address of your liking, all you need to do is set the `spec.email.address` field.
Example:
```yaml
apiVersion: k8s.checklyhq.com/v1alpha1
kind: AlertChannel
metadata:
  name: checkly-operator-test-email
spec:
  sendRecovery: false
  sendFailure: true
  sslExpiry: true
  sslExpiryThreshold: 30
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
  sendRecovery: false
  sendFailure: true
  sslExpiry: true
  sslExpiryThreshold: 30
  opsgenie:
    apikey:
      name: test-secret # Name of the secret or configmap which holds the API key
      namespace: default # Namespace of the secret or configmap
      fieldPath: "API_KEY" # Key inside the secret or configmap
     priority: "P3" # P1, P2, P3, P4, P5 are the options
     region: "EU" # Your OpsGenie region
```

### Webhook
The webhook integration supports all the fields which are supported by the checkly-go-sdk, see [details](https://pkg.go.dev/github.com/checkly/checkly-go-sdk#AlertChannelWebhook). For other fields and their options, please see [the official docs](https://www.checklyhq.com/docs/alerting-and-retries/webhooks/).

The `WebhookSecret` is an optional field and it requires a kubernetes secret (just like the OpsGenie integration).

Minimum required fields:
* `name` - string
* `url` - string
* `method` - string

```yaml
apiVersion: k8s.checklyhq.com/v1alpha1
kind: AlertChannel
metadata:
  name: checkly-operator-test-webhook
spec:
  sendRecovery: false
  sendFailure: true
  sslExpiry: true
  sslExpiryThreshold: 30
  webhook:
    name: foo # Name of the webhook
    url: https://foo.bar # URL of the webhook
    webhookType : baz # Type of webhook
    method: POST # Method of webhook
    template: testing # Checkly webhook template
    webhookSecret:
      name: test-secret # Name of the secret or configmap which holds the webhook secret
      namespace: default # Namespace of the secret or configmap
      fieldPath: "SECRET_KEY" # Key inside the secret or configmap
    headers:
      - key: "foo"
        value: "bar"
        locked: true # Not visible in the UI
    queryParameters:
      - key: "bar"
        value: "baz"
        locked: false # Visible in the UI
```

## Referencing

You'll need to reference the name of the alert channel in the group check configuration. See [check-group](check-group.md) for more details.
