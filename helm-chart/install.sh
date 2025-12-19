#!/bin/bash

# KubeDB Orchestrator - Quick Install Script

set -e

echo "🚀 Installing KubeDB Orchestrator via Helm..."

# Check if helm is installed
if ! command -v helm &> /dev/null; then
    echo "Helm is not installed. Please install Helm first."
    exit 1
fi

# Check if kubectl is installed
if ! command -v kubectl &> /dev/null; then
    echo "kubectl is not installed. Please install kubectl first."
    exit 1
fi

# Variables
NAMESPACE=${NAMESPACE:-kubedb}
RELEASE_NAME=${RELEASE_NAME:-kubedb-orchestrator}
KUBECONFIG_PATH=${KUBECONFIG_PATH:-""}
VELERO_CREDENTIALS_PATH=${VELERO_CREDENTIALS_PATH:-""}

echo "📦 Configuration:"
echo "  Namespace: $NAMESPACE"
echo "  Release Name: $RELEASE_NAME"

# Create namespace if it doesn't exist
kubectl create namespace $NAMESPACE --dry-run=client -o yaml | kubectl apply -f -

# Create kubeconfig secret if path provided
if [ -n "$KUBECONFIG_PATH" ]; then
    echo "🔐 Creating kubeconfig secret..."
    kubectl create secret generic orchestrator-kubeconfig \
        --from-file=kubeconfig=$KUBECONFIG_PATH \
        --namespace $NAMESPACE \
        --dry-run=client -o yaml | kubectl apply -f -
    echo " Kubeconfig secret created"
else
    echo "No kubeconfig path provided. Make sure orchestrator-kubeconfig secret exists!"
fi

# Create velero credentials secret if path provided
if [ -n "$VELERO_CREDENTIALS_PATH" ]; then
    echo "🔐 Creating Velero credentials secret..."
    kubectl create secret generic velero-credentials \
        --from-file=credentials-velero=$VELERO_CREDENTIALS_PATH \
        --namespace $NAMESPACE \
        --dry-run=client -o yaml | kubectl apply -f -
    echo "Velero credentials secret created"
fi

# Install/Upgrade Helm chart
echo "Installing Helm chart..."
cd "$(dirname "$0")"
helm upgrade --install $RELEASE_NAME ./kubedb-orchestrator \
    --namespace $NAMESPACE \
    --create-namespace \
    --wait \
    --timeout 5m

echo ""
echo "Installation completed!"
echo ""
echo "To check the status:"
echo "  kubectl get pods -n $NAMESPACE -l app.kubernetes.io/name=kubedb-orchestrator"
echo ""
echo "To view logs:"
echo "  kubectl logs -n $NAMESPACE -l app.kubernetes.io/name=kubedb-orchestrator -f"
echo ""
echo "To access the service:"
echo "  kubectl port-forward -n $NAMESPACE svc/$RELEASE_NAME 8080:8080"
echo ""

