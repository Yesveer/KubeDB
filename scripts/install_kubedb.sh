#!/bin/bash
set -e
set -x

export KUBECONFIG=kubeconfig.yaml

# Install runtime tools
bash scripts/bootstrap.sh

kubectl version --client
helm version

kubectl apply -f scripts/local-path-storage.yaml
kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.15.3/config/manifests/metallb-native.yaml

helm install kubedb oci://ghcr.io/appscode-charts/kubedb \
  --version v2025.10.17 \
  --namespace kubedb --create-namespace \
  --set global.featureGates.ClickHouse=true \
  --set-file global.license=scripts/licence.txt \
  --set networkPolicy.enabled=false \
  --wait