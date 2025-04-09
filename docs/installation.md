# NeuralLog Installation Guide

This guide provides detailed instructions for installing and configuring the NeuralLog infrastructure.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Installation Options](#installation-options)
- [Development Environment Setup](#development-environment-setup)
- [Test Kubernetes Environment Setup](#test-kubernetes-environment-setup)
- [Production Environment Setup](#production-environment-setup)
- [Tenant Operator Installation](#tenant-operator-installation)
- [Tenant Creation](#tenant-creation)
- [Verification](#verification)
- [Upgrading](#upgrading)
- [Uninstallation](#uninstallation)

## Prerequisites

### Required Software

- **Docker**: [Docker Desktop](https://www.docker.com/products/docker-desktop/) (Windows/macOS) or Docker Engine (Linux)
- **kubectl**: [Kubernetes CLI](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- **kind**: [Kubernetes in Docker](https://kind.sigs.k8s.io/docs/user/quick-start/#installation) (for local testing)
- **Git**: [Git](https://git-scm.com/downloads)
- **Node.js**: [Node.js](https://nodejs.org/) (for development)

### Hardware Requirements

- **Development**: 4 CPU cores, 8GB RAM, 20GB disk space
- **Test**: 4 CPU cores, 8GB RAM, 20GB disk space
- **Production**: 8 CPU cores, 16GB RAM, 50GB disk space (minimum)

### Operating System Support

- **Windows**: Windows 10/11 with PowerShell 5.1+
- **macOS**: macOS 10.15+ with Bash
- **Linux**: Ubuntu 20.04+, CentOS 8+, or other modern distributions with Bash

## Installation Options

NeuralLog can be installed in several ways:

1. **Development Environment**: Local development using Docker Compose
2. **Test Kubernetes Environment**: Local Kubernetes cluster using kind
3. **Production Environment**: Production Kubernetes cluster

## Development Environment Setup

### Step 1: Clone the Repository

```bash
git clone https://github.com/NeuralLog/infra.git
cd infra
```

### Step 2: Initialize the Development Environment

#### Windows (PowerShell)

```powershell
.\scripts\Initialize-DevEnvironment.ps1
```

#### Linux/macOS (Bash)

```bash
chmod +x scripts/*.sh
./scripts/initialize-dev-environment.sh
```

This script will:
- Check for required prerequisites
- Verify Docker is running
- Check for the NeuralLog/server repository
- Pull required Docker images

### Step 3: Start the Development Environment

#### Windows (PowerShell)

```powershell
.\scripts\Start-DevEnvironment.ps1
```

#### Linux/macOS (Bash)

```bash
./scripts/start-dev-env.sh
```

This will start:
- NeuralLog server on http://localhost:3030
- Redis on port 6379
- Redis Commander on http://localhost:8081

### Step 4: Verify the Installation

1. Open a web browser and navigate to http://localhost:3030/health
2. You should see a health check response indicating the server is running
3. Navigate to http://localhost:8081 to access Redis Commander

### Step 5: Stop the Development Environment

When you're done, you can stop the development environment:

#### Windows (PowerShell)

```powershell
.\scripts\Stop-DevEnvironment.ps1
```

#### Linux/macOS (Bash)

```bash
./scripts/stop-dev-env.sh
```

## Test Kubernetes Environment Setup

### Step 1: Clone the Repository

```bash
git clone https://github.com/NeuralLog/infra.git
cd infra
```

### Step 2: Set Up the Test Kubernetes Cluster

#### Windows (PowerShell)

```powershell
.\scripts\Setup-TestCluster.ps1
```

#### Linux/macOS (Bash)

```bash
chmod +x scripts/*.sh
./scripts/setup-test-cluster.sh
```

This script will:
- Create a kind cluster
- Create the neurallog namespace
- Apply the Kubernetes configurations
- Wait for deployments to be ready

### Step 3: Verify the Installation

```bash
kubectl get pods -n neurallog
```

You should see the NeuralLog server and Redis pods running.

### Step 4: Access the Server

```bash
kubectl port-forward svc/neurallog-server 3030:3030 -n neurallog
```

Then open a web browser and navigate to http://localhost:3030/health

### Step 5: Clean Up the Test Environment

When you're done, you can clean up the test environment:

#### Windows (PowerShell)

```powershell
.\scripts\Cleanup-TestCluster.ps1
```

#### Linux/macOS (Bash)

```bash
./scripts/cleanup-test-cluster.sh
```

## Production Environment Setup

### Step 1: Clone the Repository

```bash
git clone https://github.com/NeuralLog/infra.git
cd infra
```

### Step 2: Configure the Production Environment

1. Create a production overlay:

```bash
mkdir -p kubernetes/overlays/production
```

2. Create a kustomization.yaml file:

```bash
cat > kubernetes/overlays/production/kustomization.yaml << EOF
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
- ../../base

namespace: neurallog-prod

patches:
- path: patches/resource-limits.yaml
  target:
    kind: Deployment
    name: neurallog-server

images:
- name: neurallog/server
  newTag: stable
EOF
```

3. Create a resource limits patch:

```bash
mkdir -p kubernetes/overlays/production/patches
cat > kubernetes/overlays/production/patches/resource-limits.yaml << EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: neurallog-server
spec:
  template:
    spec:
      containers:
      - name: server
        resources:
          requests:
            cpu: 1
            memory: 2Gi
          limits:
            cpu: 2
            memory: 4Gi
EOF
```

### Step 3: Build and Push the Server Image

```bash
./scripts/build-server-image.sh stable true
```

### Step 4: Apply the Kubernetes Configurations

```bash
kubectl apply -k kubernetes/overlays/production
```

### Step 5: Verify the Installation

```bash
kubectl get pods -n neurallog-prod
```

You should see the NeuralLog server and Redis pods running.

## Tenant Operator Installation

### Step 1: Install the CRDs

```bash
kubectl apply -f operator/config/crd/bases
```

### Step 2: Install the RBAC Resources

```bash
kubectl apply -f operator/config/rbac
```

### Step 3: Install the Operator

```bash
kubectl apply -f operator/config/manager
```

### Step 4: Verify the Installation

```bash
kubectl get pods -n system
```

You should see the tenant operator pod running.

## Tenant Creation

### Step 1: Create a Tenant Resource

Create a file named `tenant.yaml`:

```yaml
apiVersion: neurallog.io/v1
kind: Tenant
metadata:
  name: tenant1
spec:
  displayName: "Tenant 1"
  description: "First tenant"
  
  server:
    replicas: 1
    
  redis:
    storage: "5Gi"
```

### Step 2: Apply the Tenant Resource

```bash
kubectl apply -f tenant.yaml
```

### Step 3: Verify the Tenant Creation

```bash
kubectl get tenants
```

You should see the tenant with its status.

```bash
kubectl get pods -n tenant-tenant1
```

You should see the tenant's server and Redis pods running.

## Verification

### Verify the Server

```bash
# Get the namespace
NAMESPACE=$(kubectl get tenant tenant1 -o jsonpath='{.status.namespace}')

# Port forward to the server
kubectl -n $NAMESPACE port-forward svc/neurallog-server 3030:3030
```

Then open a web browser and navigate to http://localhost:3030/health

### Verify Redis

```bash
# Get the namespace
NAMESPACE=$(kubectl get tenant tenant1 -o jsonpath='{.status.namespace}')

# Port forward to Redis Commander
kubectl -n $NAMESPACE port-forward svc/redis 6379:6379
```

Then you can connect to Redis using a Redis client.

## Upgrading

### Upgrading the Server

1. Build and push a new server image:

```bash
./scripts/build-server-image.sh new-version true
```

2. Update the image in the deployment:

```bash
kubectl set image deployment/neurallog-server server=neurallog/server:new-version -n neurallog
```

### Upgrading the Operator

1. Apply the new operator manifests:

```bash
kubectl apply -f operator/config/manager
```

## Uninstallation

### Uninstall Tenants

```bash
kubectl delete tenants --all
```

### Uninstall the Operator

```bash
kubectl delete -f operator/config/manager
kubectl delete -f operator/config/rbac
kubectl delete -f operator/config/crd/bases
```

### Uninstall the NeuralLog Infrastructure

```bash
kubectl delete -k kubernetes/overlays/production
```
