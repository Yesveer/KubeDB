#!/bin/bash
set -euo pipefail

S3_URL="$1"
ACCESS_KEY="$2"
SECRET_KEY="$3"

echo "🚀 Installing Velero..."
echo "👉 S3 URL: $S3_URL"

export KUBECONFIG="$(pwd)/kubeconfig.yaml"

# Check kubectl access
kubectl get ns >/dev/null

# Install velero binary (NO sudo)
if ! command -v velero &>/dev/null; then
  echo "📦 Installing velero binary"
  wget -q https://github.com/vmware-tanzu/velero/releases/download/v1.14.0/velero-v1.14.0-linux-amd64.tar.gz
  tar -xzf velero-v1.14.0-linux-amd64.tar.gz
  chmod +x velero-v1.14.0-linux-amd64/velero
  mv velero-v1.14.0-linux-amd64/velero /usr/local/bin/velero
fi

# Namespace
kubectl create namespace velero --dry-run=client -o yaml | kubectl apply -f -

# Credentials
cat <<EOF > credentials-velero
[default]
aws_access_key_id=${ACCESS_KEY}
aws_secret_access_key=${SECRET_KEY}
EOF

kubectl delete secret cloud-credentials -n velero --ignore-not-found
kubectl create secret generic cloud-credentials \
  -n velero \
  --from-file=cloud=credentials-velero

# Install velero
velero install \
  --provider aws \
  --plugins velero/velero-plugin-for-aws:v1.9.0 \
  --bucket velero \
  --secret-file ./credentials-velero \
  --backup-location-config region=minio,s3ForcePathStyle=true,s3Url=${S3_URL} \
  --use-node-agent \
  --wait

# Verify
velero backup-location get | grep -q Available

echo "✅ Velero installed successfully"
