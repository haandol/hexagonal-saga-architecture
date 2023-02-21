# Deploy on Kubernetes

## Prerequisites

- [aws-cli](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-install.html)
- [eksctl](https://eksctl.io/introduction/#installation)
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)

## Build and push ECR image

```bash
task build-all
task push-all
```

## Create ServiceAccount for the migration

```bash
export CLUSTER_NAME=$(aws eks list-clusters --query 'clusters[0]' --output text)
export ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)

eksctl utils associate-iam-oidc-provider --region=ap-northeast-2 --cluster=${CLUSTER_NAME} --approve

eksctl create iamserviceaccount --name migration-service-account --namespace default --cluster ${CLUSTER_NAME} --role-name "MigrationSaRole" \
    --attach-policy-arn arn:aws:iam::aws:policy/SecretsManagerReadWrite --approve
```

## Deploy the application

```bash
kubectl apply -k infra/overlays/dev
```

## Deploy Istio gateway / virtualservice

```bash
kubectl apply -f infra/istio/gateway.yaml
```