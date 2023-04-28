# bootstrap-operator

Spin up k8s cluster using cluster API.

## Development

### Code and manifest generation

Code and manifest generation is required after changing types. To do it run following commands from repository root:

```sh
make generate manifests
```

### Running operator locally

In order to run operator locally first prepare `bootstrap.yaml` manifest and then run following commands from repository root:

```sh
# Start kind cluster:
kind create cluster
kubectl cluster-info --context kind-kind

# Apply custom resources:
kubectl apply -f example/crds
kubectl apply -f config/crd/bases

# Apply bootstrap config:
kubectl apply -f ../bootstrap.yaml

# Install cert-manager:
helm install cert-manager jetstack/cert-manager --namespace bootstrap --create-namespace --version v1.11.0 --set installCRDs=true

# Run operator:
go run cmd/main.go --namespace=bootstrap

# Delete kind cluster afterwards:
kind delete cluster
```
