# Deploy on Kubernetes

## Build and push ECR image

```bash
task build-push-all
```

## Deploy the application

```bash
kubectl apply -k infra/overlays/dev
```

## Deploy Istio gateway / virtualservice

```bash
kubectl apply -f istio/gateway.yaml
```