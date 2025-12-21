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

helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update

helm install kube-prom-stack prometheus-community/kube-prometheus-stack \
  -n monitoring \
  --create-namespace \
  --set grafana.enabled=true \
  --set prometheus.prometheusSpec.serviceMonitorSelectorNilUsesHelmValues=false \
  --set prometheus.prometheusSpec.podMonitorSelectorNilUsesHelmValues=false \
  --wait