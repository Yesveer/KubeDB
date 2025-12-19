# KubeDB Orchestrator Helm Chart

Complete Helm chart for deploying KubeDB Orchestrator - a Database as a Service (DBaaS) management system.

## 📁 Structure

```
helm-chart/
├── install.sh                          # Quick installation script
├── kubedb-orchestrator/
│   ├── Chart.yaml                      # Chart metadata
│   ├── values.yaml                     # Default values
│   ├── values-production.yaml          # Production overrides example
│   ├── .helmignore                     # Helm ignore patterns
│   ├── README.md                       # Detailed chart documentation
│   └── templates/
│       ├── _helpers.tpl                # Template helpers/functions
│       ├── NOTES.txt                   # Post-install instructions
│       ├── configmap.yaml              # Application configuration
│       ├── configmap-scripts.yaml      # Installation scripts
│       ├── deployment.yaml             # Main deployment
│       ├── service.yaml                # Kubernetes service
│       ├── serviceaccount.yaml         # Service account
│       ├── clusterrole.yaml            # RBAC permissions
│       └── clusterrolebinding.yaml     # RBAC binding
```

## 🚀 Quick Start

### Prerequisites

- Kubernetes 1.19+
- Helm 3.0+
- kubectl configured to access your cluster

### Method 1: Using the Install Script (Recommended)

```bash
cd helm-chart

# With kubeconfig
KUBECONFIG_PATH=/path/to/kubeconfig ./install.sh

# With kubeconfig and Velero credentials
KUBECONFIG_PATH=/path/to/kubeconfig \
VELERO_CREDENTIALS_PATH=/path/to/credentials \
NAMESPACE=my-namespace \
./install.sh
```

### Method 2: Manual Installation

1. **Create the kubeconfig secret:**
```bash
kubectl create secret generic orchestrator-kubeconfig \
  --from-file=kubeconfig=/path/to/your/kubeconfig \
  --namespace default
```

2. **(Optional) Create Velero credentials:**
```bash
kubectl create secret generic velero-credentials \
  --from-file=credentials-velero=/path/to/credentials \
  --namespace default
```

3. **Install the chart:**
```bash
helm install kubedb-orchestrator ./kubedb-orchestrator \
  --namespace default
```

### Method 3: With Custom Values

```bash
helm install kubedb-orchestrator ./kubedb-orchestrator \
  --namespace default \
  --values ./kubedb-orchestrator/values-production.yaml
```

## ⚙️ Configuration

### Key Configuration Options

| Parameter | Description | Default |
|-----------|-------------|---------|
| `replicaCount` | Number of replicas | `1` |
| `image.repository` | Docker image repository | `amiteshhsingh/kubedb-orchestrator` |
| `image.tag` | Docker image tag | `v1` |
| `image.pullPolicy` | Image pull policy | `Always` |
| `service.type` | Kubernetes service type | `ClusterIP` |
| `service.port` | Service port | `8080` |
| `resources.limits.cpu` | CPU limit | `500m` |
| `resources.limits.memory` | Memory limit | `512Mi` |
| `resources.requests.cpu` | CPU request | `100m` |
| `resources.requests.memory` | Memory request | `128Mi` |
| `config.compassBaseUrl` | Compass API base URL | `http://172.16.109.237:32500` |
| `config.kubeconfigPath` | Path to kubeconfig in pod | `/data/kubeconfig` |
| `config.mongoUri` | MongoDB connection URI | See values.yaml |
| `secrets.kubeconfigSecretName` | Kubeconfig secret name | `orchestrator-kubeconfig` |
| `secrets.veleroCredentialsSecretName` | Velero credentials secret | `velero-credentials` |
| `rbac.create` | Create RBAC resources | `true` |
| `serviceAccount.create` | Create service account | `true` |

### Example Custom Configuration

Create `my-values.yaml`:

```yaml
# Scale to 2 replicas
replicaCount: 2

# Use a different image version
image:
  tag: "v2"
  pullPolicy: IfNotPresent

# Increase resources
resources:
  limits:
    cpu: 1000m
    memory: 1Gi
  requests:
    cpu: 200m
    memory: 256Mi

# Use NodePort for external access
service:
  type: NodePort

# Custom configuration
config:
  compassBaseUrl: "http://my-compass.example.com"
  mongoUri: "mongodb://user:pass@mongo.example.com:27017/mydb"

# Enable health checks (requires /health endpoint in app)
livenessProbe:
  enabled: true

readinessProbe:
  enabled: true
```

Install with custom values:

```bash
helm install kubedb-orchestrator ./kubedb-orchestrator -f my-values.yaml
```

## 📦 Managing the Deployment

### Check Installation Status

```bash
helm status kubedb-orchestrator -n default
```

### View Pod Status

```bash
kubectl get pods -n default -l app.kubernetes.io/name=kubedb-orchestrator
```

### View Logs

```bash
kubectl logs -n default -l app.kubernetes.io/name=kubedb-orchestrator -f
```

### Upgrade

```bash
helm upgrade kubedb-orchestrator ./kubedb-orchestrator \
  --namespace default \
  --values my-values.yaml
```

### Rollback

```bash
# List releases
helm history kubedb-orchestrator -n default

# Rollback to previous version
helm rollback kubedb-orchestrator -n default

# Rollback to specific revision
helm rollback kubedb-orchestrator 1 -n default
```

### Uninstall

```bash
helm uninstall kubedb-orchestrator -n default
```

## 🌐 Accessing the Application

### Port Forward (for ClusterIP)

```bash
kubectl port-forward svc/kubedb-orchestrator 8080:8080 -n default
```

Then access at: `http://localhost:8080`

### NodePort

If using NodePort service:

```bash
export NODE_PORT=$(kubectl get svc kubedb-orchestrator -n default -o jsonpath='{.spec.ports[0].nodePort}')
export NODE_IP=$(kubectl get nodes -o jsonpath='{.items[0].status.addresses[0].address}')
echo "Access at: http://$NODE_IP:$NODE_PORT"
```

### LoadBalancer

If using LoadBalancer service:

```bash
export SERVICE_IP=$(kubectl get svc kubedb-orchestrator -n default -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
echo "Access at: http://$SERVICE_IP:8080"
```

## 🔧 Troubleshooting

### Pod Not Starting

```bash
# Check pod status
kubectl describe pod -n default -l app.kubernetes.io/name=kubedb-orchestrator

# Check logs
kubectl logs -n default -l app.kubernetes.io/name=kubedb-orchestrator --previous
```

### Secret Not Found

Verify secrets exist:

```bash
kubectl get secrets -n default
```

Create missing secrets:

```bash
kubectl create secret generic orchestrator-kubeconfig \
  --from-file=kubeconfig=/path/to/kubeconfig \
  --namespace default
```

### RBAC Issues

Check if RBAC resources were created:

```bash
kubectl get clusterrole kubedb-orchestrator
kubectl get clusterrolebinding kubedb-orchestrator
kubectl get serviceaccount kubedb-orchestrator -n default
```

### View All Resources

```bash
kubectl get all -n default -l app.kubernetes.io/name=kubedb-orchestrator
```

## 🔐 Security Notes

1. **Secrets Management**: The chart expects secrets to be created manually before installation. For production, consider using:
   - External Secrets Operator
   - Sealed Secrets
   - Vault

2. **RBAC**: The chart creates a ClusterRole with broad permissions. Review and adjust `rbac.rules` in `values.yaml` for your security requirements.

3. **Network Policies**: Consider adding NetworkPolicy resources to restrict pod communication.

## 🧪 Testing the Chart

### Lint the Chart

```bash
helm lint ./kubedb-orchestrator
```

### Dry Run

```bash
helm install kubedb-orchestrator ./kubedb-orchestrator \
  --namespace default \
  --dry-run --debug
```

### Template Rendering

```bash
helm template kubedb-orchestrator ./kubedb-orchestrator \
  --namespace default \
  --values my-values.yaml > rendered.yaml
```

## 📝 Development

### Modifying the Chart

1. Edit values in `values.yaml` or templates in `templates/`
2. Update `Chart.yaml` version
3. Test with `helm lint` and `helm template`
4. Upgrade the release:

```bash
helm upgrade kubedb-orchestrator ./kubedb-orchestrator \
  --namespace default \
  --values my-values.yaml
```

### Packaging the Chart

```bash
helm package ./kubedb-orchestrator
```

This creates `kubedb-orchestrator-0.1.0.tgz`

### Publishing to a Helm Repository

```bash
# Index the chart
helm repo index .

# Upload to your chart repository
# (method depends on your repository type)
```

## 📚 Additional Resources

- [Helm Documentation](https://helm.sh/docs/)
- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [KubeDB Documentation](https://kubedb.com/docs/)

## 🤝 Contributing

When contributing to the chart:

1. Test all changes with `helm lint`
2. Update documentation
3. Increment version in `Chart.yaml`
4. Test installation in a clean namespace

## 📄 License

[Your License Here]

