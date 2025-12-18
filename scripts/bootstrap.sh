#!/bin/bash
set -e
set -x

# Set kubeconfig
export KUBECONFIG=kubeconfig.yaml

echo "Installing kubectl"
if ! command -v kubectl &>/dev/null; then
  curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl
  chmod +x kubectl
  mv kubectl /usr/local/bin/
fi

echo "Installing helm"
if ! command -v helm &>/dev/null; then
  curl -L https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 | bash
fi
