# Consularis

*Consul Key/Value controller for Kubernetes*

Consularis implements a Consul Object custom resource in Kubernetes containing one or more Key/Value pairs as well as a controller for that resource.
This enables you to use Kubernetes API objects to set the KV state in Consul.

Consularis will only work against Consul 0.7.x + as it relies on Key/Value Transactions which are only available in version 0.7 and up.

By default Consularis will assume it's running inside Kubernetes and tries to grab configuration from the cluster.
Run `./consularis -h` for details on how to run it outside Kubernetes.

## Build & Development
To build Consularis you first need to make sure you have all the dependencies.
```
brew install dep
brew upgrade dep
dep ensure
```
To build it, run `go build`.

If you make any changes to the specs under `pkg/apis/consularis.io/*` you need to re-generate the code.
```
hack/update-codegen.sh
```

## Kubernetes Custom Resource Definitions
Consularis manages Consul Key/Values using [Custom Resource Definitions](https://kubernetes.io/docs/concepts/api-extension/custom-resources/) (CRD).
At startup time Consularis will create the CRD if it does not already exist. 

The CRD Consularis manages is called ConsulObjects and exists under the `consulobjects.consularis.io` CRD.
You can check existing CRDs using `kubectl get crd` and you can list the `consulobjects` by running `kubectl get consulobjects` or `kubectl get co`.

The ConsulObjects CRD supports all the regular kubectl API interactions.


| kubectl cmd | Result |
| -------- | -------- |
| `kubectl get consulobjects` &#124; `co`  | Lists all ConsulObjects |
| `kubectl apply -f /path/to/my-kv.yaml` | Will create or update the Key/Values `my-kv` manages |
| `kubectl describe co my-kv` | Describes the object and lists the Key/Values it manages |
| `kubectl delete -f my-kv.yaml` &#124; `kubectl delete consulobject my-kv` | Deletes an object using the file that was used to create it or directly by referencing the API object |

> Note that `kubectl describe` will try to do some fancy "pretty yaml" formatting that can break when you have "\n" characters in the `value`. To get around that you can use `kubectl get co my-kv -o json` to output `my-kv` in JSON format instead which preserves the original `value`.

Consularis runs a Kubernetes controller that watches the ConsulObjects API for any create, update and delete events and updates Consul KV accordingly.
It will also re-sync existing ConsulObjects in regular intervals which means if you make a manual change in the Consul UI it will almost certainly be overritten.

### Defining Key/Value resources
To add new KV object you declare a `consulobject` resource:
```
apiVersion: consularis.io/v1alpha1
kind: ConsulObject
metadata:
  name: web
spec:
  kv:
    - key: web/cors_origin
      value: https://*.example.org
    - key: web/disable_error_pages
      value: 'true'
    - key: web/enable_test_environment
      value: 'true'
    - key: web/enforce_ssl
      flag: 0
      value: 'false'
```

The Keys and Values are arbitrary as long as you follow YAML spec and it also supports Flags.

## TODO
- [ ] Give a nicer error message if you run Consularis outside of Kubernetes without flags.
- [ ] Create basic example Dockerfile