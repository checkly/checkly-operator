# docs

The checkly-operator was designed to run inside a kubernetes cluster and listen for events on specific CRDs and ingress resources. With the help of it you can set up:
* [Alert channels](alert-channels.md)
* [Check groups](check-group.md)
* [API Checks](api-checks.md)

## Installation

We currently supply an installation yaml file, this is present in the [releases](https://github.com/checkly/checkly-operator/releases).

The file holds the following resources:
* namespace
* CRDs
* RBAC
* deployment

In order for the operator to work, you need to supply secrets which hold your [checklyhq.com](checklyhq.com) API Key and Account ID. See [docs on how to create an API key and get account ID](https://www.checklyhq.com/docs/integrations/pulumi/#define-your-checkly-account-id-and-api-key). The operator expects the following environment variables:
```
        env:
        - name: CHECKLY_API_KEY
          valueFrom:
            secretKeyRef:
              key: CHECKLY_API_KEY
              name: checkly
        - name: CHECKLY_ACCOUNT_ID
          valueFrom:
            secretKeyRef:
              key: CHECKLY_ACCOUNT_ID
              name: checkly
```

The following steps are an easy example on how to get started with the operator, it is not a production ready method, ex. we're not using any secrets managers, you should not create secrets and commit them to git like in the bellow example, we're only deploying one replica, while the operator does support HA deployments.

If you just want to try out the checkly-operator, you need a local kubernetes installation, the easiest might be [Rancher Desktop](https://rancherdesktop.io/), see [the docs](https://docs.rancherdesktop.io/getting-started/installation/) for the installation, once done, come back to this doc.

### Download install.yaml

First we'll download the provided `install.yaml` files, please change the version number accordingly, we might have newer [releases](https://github.com/checkly/checkly-operator/releases) since we've written these docs.

```bash
export CHECKLY_OPERATOR_RELEASE=v1.4.1
wget "https://github.com/checkly/checkly-operator/releases/download/$CHECKLY_OPERATOR_RELEASE/install-$CHECKLY_OPERATOR_RELEASE.yaml" -O install.yaml
unset CHECKLY_OPERATOR_RELEASE
```

Feel free to edit the `install.yaml` file to your liking, usually you'd want to change:
* checkly-operator deployment replica count
* checkly-operator deployment CPU and Memory resources

You can apply the `install.yaml`, this will create the namespace, we need this to create the secrets in the next step:
```bash
kubectl apply -f install.yaml
```

### Create secret

Grab your [checklyhq.com](checklyhq.com) API key and Account ID, [the official docs](https://www.checklyhq.com/docs/integrations/pulumi/#define-your-checkly-account-id-and-api-key) can help you get this information. Substitute the values into the bellow command:

```bash
export CHECKLY_API_KEY=<api-key-from-checklyhq.com>
export CHECKLY_ACCOUNT_ID=<org-id-from-checklyhq.com>
kubectl create secret generic -n checkly-operator-system checkly \
  --from-literal=CHECKLY_API_KEY=$CHECKLY_API_KEY \
  --from-literal=CHECKLY_ACCOUNT_ID=$CHECKLY_ACCOUNT_ID
unset CHECKLY_API_KEY
unset CHECKLY_ACCOUNT_ID
```

If you check your pod, you should be able to see the pods starting:
```bash
kubectl get pods -n checkly-operator-system
```

## Configuration

We will next create a check group, alert channels and api checks through the custom CRDs. The operator was written with a specific opinion:
* all checks should belong to a specific group
* group configuration should be inherited by each check which is part of it
* alert channels are added to groups, not to individual checks

Based on the above, the order of creation should be:
1. alert channel
2. check group
3. api checks

Reference to resources are done based on the kubernetes internal naming, as in the `metadata.name` field.

Please look at the bellow examples and change the supplied data so it fits your needs the best. Save the example into individual files and apply them when ready:
```bash
kubectl apply -f <name-of-the-file>.yaml
```

### Alert channel

See the [docs](https://www.checklyhq.com/docs/alerting/) on what alert channels are and [alert-channels](alert-channels.md) for the options we support.

The following configuration will send all alerts to an email address:
```yaml
apiVersion: k8s.checklyhq.com/v1alpha1
kind: AlertChannel
metadata:
  name: checkly-operator-test-alertchannel
spec:
  email:
    address: "foo@bar.baz"
```

The `AlertChannel` resource is clustered scoped, it can be used in any check group but it also means that the name of the resource has to be unique in the kubernetes cluster.

Once applied you can check if it worked:
```bash
kubectl get alertchannel checkly-operator-test-alertchannel
```

You can also view the alert channel on the [checklyhq.com dashboard](https://app.checklyhq.com/alert-settings).

### Check group

See the [docs](https://www.checklyhq.com/docs/groups/) on what check groups are and [check-group](check-group.md) for the options we support.

The following configuration adds the above alert channel to all the checks that will be part of the group check:
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
    - checkly-operator-test-alertchannel
```

The `Group` resource is clustered scoped, it can be used in any api check but it also means that the name of the resource has to be unique in the kubernetes cluster.

Once applied you can check if it worked:
```bash
kubectl get group checkly-operator-test-group
```

You can also view the check group on the [checklyhq.com dashboard](https://app.checklyhq.com/).

#### Tags

Any labels added to the `Group` resource will be added as tags to the group, these groups are inherited by the checks.

#### Locations

To see a full list of locations supported by [checklyhq.com](checklyhq.com), see [the docs](https://www.checklyhq.com/docs/monitoring/global-locations/). We're using the location codes, private locations should be technically supported, we just haven't tested them.

### API Checks

See the [docs](https://www.checklyhq.com/docs/api-checks/) on what API checks are and [api-checks](api-checks.md) for the options we support.

We're currently only supporting API checks which perform a GET request.

The following configuration monitors the `https://foo.bar/baz` endpoint, expects a return code of 200, it's added to the above created group and are muted so they do not send an alert:
```yaml
apiVersion: k8s.checklyhq.com/v1alpha1
kind: ApiCheck
metadata:
  name: checkly-operator-test-1
  namespace: default
  labels:
    service: "foo"
spec:
  endpoint: "http://foo.bar/baz"
  success: "200"
  frequency: 10 # Default 5
  muted: true # Default "false"
  group: "checkly-operator-test-group"

```

You can create multiple `ApiCheck` resources which point to the same group:
```yaml
apiVersion: k8s.checklyhq.com/v1alpha1
kind: ApiCheck
metadata:
  name: checkly-operator-test-2
  namespace: default
  labels:
    service: "bar"
spec:
  endpoint: "https://checklyhq.com"
  success: "200"
  frequency: 10 # Default 5
  muted: true # Default "false"
  group: "checkly-operator-test-group"
```

`ApiCheck` resources are namespace scoped, they have to have a unique name in each namespace.

Once applied you can check if it worked:
```bash
kubectl get apichecks -n default
```

You can also view the checks on the [checklyhq.com dashboard](https://app.checklyhq.com/).

### Ingresses

See [ingress](ingress.md) for more details on how we utilize ingress resources.

We can create an ingress object and add annotations to it:
```yaml
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: ingress-sample
  namespace: default
  annotations:
    k8s.checklyhq.com/enabled: "true"
    k8s.checklyhq.com/path: "/baz"
    # k8s.checklyhq.com/endpoint: "foo.baaz" - Default read from spec.rules[0].host
    # k8s.checklyhq.com/success: "200" - Default "200"
    k8s.checklyhq.com/group: "checkly-operator-test-group"
    # k8s.checklyhq.com/muted: "false" # If not set, default "true"
spec:
  rules:
    - host: "foo.bar"
      http:
        paths:
          - path: /
            pathType: ImplementationSpecific
            backend:
              service:
                name: test-service
                port:
                  number: 8080

```

Check if it worked:
```bash
kubectl get ingress -n default
kubectl get apicheck -n default
```

You can also view the checks on the [checklyhq.com dashboard](https://app.checklyhq.com/).
