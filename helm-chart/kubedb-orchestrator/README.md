# KubeDB Orchestrator Helm Chart

A Helm chart for deploying KubeDB Orchestrator - a DBaaS management system.

## Prerequisites

- Kubernetes 1.19+
- Helm 3.0+

## Installation

### 1. Create the kubeconfig secret

```bash
kubectl create secret generic orchestrator-kubeconfig \
  --from-file=kubeconfig=path/to/your/kubeconfig \
  --namespace kubedb
```

### 2. (Optional) Create Velero credentials secret

```bash
kubectl create secret generic velero-credentials \
  --from-file=credentials-velero=path/to/credentials \
  --namespace kubedb
```

### 3. Install the chart

```bash
helm install kubedb-orchestrator ./kubedb-orchestrator \
  --namespace kubedb \
  --create-namespace
```

### 4. Install with custom values

```bash
helm install kubedb-orchestrator ./kubedb-orchestrator \
  --namespace kubedb \
  --values custom-values.yaml
```

## Configuration

The following table lists the configurable parameters of the KubeDB Orchestrator chart and their kubedb values.

| Parameter | Description | Default |
|-----------|-------------|---------|
| `replicaCount` | Number of replicas | `1` |
| `image.repository` | Image repository | `amiteshhsingh/kubedb-orchestrator` |
| `image.pullPolicy` | Image pull policy | `Always` |
| `image.tag` | Image tag | `v1` |
| `service.type` | Service type | `ClusterIP` |
| `service.port` | Service port | `8080` |
| `resources.limits.cpu` | CPU limit | `500m` |
| `resources.limits.memory` | Memory limit | `512Mi` |
| `resources.requests.cpu` | CPU request | `100m` |
| `resources.requests.memory` | Memory request | `128Mi` |
| `config.compassBaseUrl` | Compass base URL | `http://172.16.109.237:32500` |
| `config.kubeconfigPath` | Kubeconfig path | `/data/kubeconfig` |
| `config.mongoUri` | MongoDB URI | (see values.yaml) |
| `secrets.kubeconfigSecretName` | Kubeconfig secret name | `orchestrator-kubeconfig` |
| `secrets.veleroCredentialsSecretName` | Velero credentials secret name | `velero-credentials` |
| `rbac.create` | Create RBAC resources | `true` |
| `serviceAccount.create` | Create service account | `true` |

## Customization Example

Create a `custom-values.yaml` file:

```yaml
image:
  tag: "v2"
  pullPolicy: IfNotPresent

replicaCount: 2

resources:
  limits:
    cpu: 1000m
    memory: 1Gi
  requests:
    cpu: 200m
    memory: 256Mi

config:
  compassBaseUrl: "http://your-compass-url:port"
  mongoUri: "mongodb://user:pass@host:port/db"

service:
  type: NodePort

livenessProbe:
  enabled: true

readinessProbe:
  enabled: true
```

Then install:

```bash
helm install kubedb-orchestrator ./kubedb-orchestrator -f custom-values.yaml
```

## Upgrade

```bash
helm upgrade kubedb-orchestrator ./kubedb-orchestrator \
  --namespace kubedb \
  --values custom-values.yaml
```

## Uninstall

```bash
helm uninstall kubedb-orchestrator --namespace kubedb
```

## Accessing the Application

For ClusterIP service:

```bash
kubectl port-forward svc/kubedb-orchestrator 8080:8080 -n kubedb
```

Then access at `http://localhost:8080`

## Troubleshooting

### Check pod status
```bash
kubectl get pods -n kubedb -l app.kubernetes.io/name=kubedb-orchestrator
```

### View logs
```bash
kubectl logs -n kubedb -l app.kubernetes.io/name=kubedb-orchestrator
```

### Check service
```bash
kubectl get svc -n kubedb kubedb-orchestrator
```

### Verify secrets
```bash
kubectl get secrets -n kubedb
```

