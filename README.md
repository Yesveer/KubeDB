# KubeDB Orchestrator

A Database as a Service (DBaaS) orchestrator for managing KubeDB instances, backups, and MetalLB configurations.

## 🏗️ Architecture

This is a Go-based HTTP service that provides:
- KubeDB installation and management
- Database backup/restore via Velero
- MetalLB configuration management
- License management for KubeDB

## 🚀 Quick Start

### 1. Build Docker Image

```bash
docker build -t amiteshhsingh/kubedb-orchestrator:v1 .
docker push amiteshhsingh/kubedb-orchestrator:v1
```

### 2. Deploy with Helm

See the [`../helm-chart/`](../helm-chart/) directory for deployment instructions.

```bash
cd ../helm-chart
KUBECONFIG_PATH=/path/to/kubeconfig ./install.sh
```

## 📦 What Gets Built

The Docker image includes:
- Compiled Go binary (`/server`)
- Installation scripts (`/scripts`)
- Alpine base with `kubectl` installed
- Required certificates and bash

**NOT included in image:**
- Kubernetes manifests (*.yaml)
- Helm chart
- Documentation
- Test files

## 🔧 Development

### Prerequisites
- Go 1.24+
- Docker

### Local Build
```bash
go mod download
go build -o server ./cmd/server
```

### Run Locally
```bash
export COMPASS_BASE_URL="http://your-compass-url"
export KUBECONFIG_PATH="/path/to/kubeconfig"
export MONGO_URI="mongodb://..."
./server
```

## 📂 Project Structure

```
KubeDB/
├── cmd/server/              # Application entry point
├── internal/
│   ├── backup/             # Velero backup logic
│   ├── config/             # Configuration management
│   ├── db/                 # Database connections
│   ├── handlers/           # HTTP handlers
│   ├── installer/          # KubeDB installation
│   ├── kubeconfig/         # Kubeconfig handling
│   ├── licence/            # License management
│   ├── metallb/            # MetalLB configuration
│   ├── models/             # Data models
│   ├── repository/         # Data access layer
│   └── routes/             # Route definitions
├── scripts/                # Installation scripts
├── Dockerfile              # Docker build
├── .dockerignore           # Docker build exclusions
├── go.mod                  # Go dependencies
└── go.sum                  # Dependency checksums
```

## 🌐 API Endpoints

The service exposes various endpoints for:
- Database management
- Backup/restore operations
- KubeDB operations
- MetalLB configuration
- License management

See the handlers in `internal/handlers/` for details.

## 🐳 Dockerfile

Multi-stage build:
1. **Builder stage**: Compiles Go application
2. **Runtime stage**: Alpine Linux with kubectl and the binary

## 📝 Environment Variables

Required:
- `COMPASS_BASE_URL` - Compass API base URL
- `KUBECONFIG_PATH` - Path to kubeconfig file
- `MONGO_URI` - MongoDB connection string

Optional:
- `INSTALL_SCRIPT_PATH` - Custom install script path

## 🚢 Deployment

**Use the Helm chart for deployment!**

The Helm chart provides:
- Proper RBAC configuration
- Secret management
- Resource limits
- ConfigMaps for scripts
- Service and Ingress options

See [`../helm-chart/README.md`](../helm-chart/README.md) for complete deployment guide.

## 🔐 Security Notes

- Service requires cluster-admin level permissions
- Kubeconfig must be provided via Kubernetes Secret
- Credentials should never be in the Docker image

## 📚 Related Documentation

- [Deployment Guide](../DEPLOYMENT.md)
- [Helm Chart Documentation](../helm-chart/README.md)
- [Complete Summary](../SUMMARY.md)
- [Quick Reference](../QUICKREF.sh)

## 🤝 Contributing

1. Make changes to the code
2. Test locally
3. Build and push new Docker image
4. Update Helm chart values if needed
5. Deploy and test

## 📄 License

[Your License]
