# NeuralLog Tenant Operator Guide

This guide provides detailed information about the NeuralLog Tenant Operator, including its architecture, installation, configuration, and usage.

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Installation](#installation)
- [Tenant Custom Resource](#tenant-custom-resource)
- [Tenant Lifecycle](#tenant-lifecycle)
- [Configuration Options](#configuration-options)
- [Monitoring and Management](#monitoring-and-management)
- [Troubleshooting](#troubleshooting)
- [Advanced Usage](#advanced-usage)
- [Development](#development)

## Overview

The NeuralLog Tenant Operator is a Kubernetes operator that manages multi-tenant deployments of NeuralLog. It provides a simple API for creating and managing tenants, each with their own isolated resources.

### Key Features

- **Tenant Isolation**: Complete isolation between tenants
- **Resource Management**: Configurable resource limits and requests
- **Lifecycle Management**: Automated provisioning, updates, and cleanup
- **Network Policies**: Configurable network policies for security
- **Status Reporting**: Detailed status reporting for tenants

## Architecture

The operator follows the Kubernetes operator pattern:

1. **Custom Resource Definition (CRD)**: Defines the `Tenant` resource
2. **Controller**: Watches for changes to Tenant resources and reconciles the desired state
3. **Reconciliation Logic**: Creates, updates, or deletes Kubernetes resources based on the Tenant specification

### Component Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                     Tenant Operator                         │
│                                                             │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐  │
│  │  Tenant CRD     │  │  Controller     │  │  Reconciler │  │
│  └─────────────────┘  └─────────────────┘  └─────────────┘  │
│                                                             │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐  │
│  │  Redis Manager  │  │  Server Manager │  │  Network    │  │
│  │                 │  │                 │  │  Manager    │  │
│  └─────────────────┘  └─────────────────┘  └─────────────┘  │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### Reconciliation Flow

1. **Watch**: The controller watches for changes to Tenant resources
2. **Reconcile**: When a change is detected, the reconciler is called
3. **Namespace**: The reconciler creates or updates the tenant's namespace
4. **Redis**: The reconciler creates or updates the tenant's Redis resources
5. **Server**: The reconciler creates or updates the tenant's server resources
6. **Network**: The reconciler creates or updates the tenant's network policies
7. **Status**: The reconciler updates the tenant's status

## Installation

### Prerequisites

- Kubernetes cluster (v1.19+)
- kubectl configured to communicate with your cluster
- cert-manager (v1.0.0+)

### Installing the Operator

1. Install the CRDs:

```bash
kubectl apply -f operator/config/crd/bases
```

2. Install the RBAC resources:

```bash
kubectl apply -f operator/config/rbac
```

3. Install the operator:

```bash
kubectl apply -f operator/config/manager
```

### Verifying the Installation

```bash
kubectl get pods -n system
```

You should see the tenant operator pod running.

## Tenant Custom Resource

The Tenant custom resource allows you to define:

- Resource limits and requests
- Server configuration
- Redis configuration
- Network policies

### Example Tenant Resource

```yaml
apiVersion: neurallog.io/v1
kind: Tenant
metadata:
  name: sample-tenant
spec:
  displayName: "Sample Tenant"
  description: "A sample tenant for demonstration purposes"
  
  # Resource limits and requests for the tenant
  resources:
    cpu:
      limit: "4"
      request: "2"
    memory:
      limit: "8Gi"
      request: "4Gi"
    storage:
      limit: "20Gi"
      request: "10Gi"
  
  # Server configuration
  server:
    replicas: 2
    image: "neurallog/server:latest"
    resources:
      cpu:
        limit: "1"
        request: "500m"
      memory:
        limit: "1Gi"
        request: "512Mi"
    env:
      - name: LOG_LEVEL
        value: "debug"
      - name: MAX_CONNECTIONS
        value: "100"
      - name: API_KEY
        valueFrom:
          secretKeyRef:
            name: tenant-secrets
            key: api-key
  
  # Redis configuration
  redis:
    replicas: 1
    image: "redis:7-alpine"
    resources:
      cpu:
        limit: "500m"
        request: "200m"
      memory:
        limit: "1Gi"
        request: "512Mi"
    storage: "5Gi"
    config:
      maxmemory: "512mb"
      maxmemory-policy: "allkeys-lru"
  
  # Network policy configuration
  networkPolicy:
    enabled: true
    allowedNamespaces:
      - "default"
      - "monitoring"
    ingressRules:
      - description: "Allow monitoring tools"
        from:
          app: "prometheus"
        ports:
          - protocol: "TCP"
            port: 3030
    egressRules:
      - description: "Allow external API access"
        to:
          app: "external-api"
        ports:
          - protocol: "TCP"
            port: 443
```

### Tenant Status

The Tenant resource includes a status section that provides information about the tenant's state:

```yaml
status:
  phase: Running
  namespace: tenant-sample-tenant
  serverStatus:
    phase: Running
    readyReplicas: 2
    totalReplicas: 2
    message: "Server is running"
  redisStatus:
    phase: Running
    readyReplicas: 1
    totalReplicas: 1
    message: "Redis is running"
```

## Tenant Lifecycle

### Creating a Tenant

1. Create a Tenant resource:

```bash
kubectl apply -f tenant.yaml
```

2. The operator will:
   - Create a namespace for the tenant
   - Deploy Redis in the namespace
   - Deploy the server in the namespace
   - Configure network policies
   - Update the tenant status

### Updating a Tenant

1. Update the Tenant resource:

```bash
kubectl apply -f tenant.yaml
```

2. The operator will:
   - Update the Redis configuration if needed
   - Update the server configuration if needed
   - Update the network policies if needed
   - Update the tenant status

### Deleting a Tenant

1. Delete the Tenant resource:

```bash
kubectl delete tenant sample-tenant
```

2. The operator will:
   - Delete the tenant's namespace
   - Delete all resources in the namespace
   - Remove the finalizer from the Tenant resource

## Configuration Options

### Server Configuration

| Field | Description | Default |
|-------|-------------|---------|
| `replicas` | Number of server replicas | 1 |
| `image` | Server Docker image | neurallog/server:latest |
| `resources` | Resource limits and requests | See below |
| `env` | Environment variables | See below |

#### Default Resource Limits

```yaml
resources:
  cpu:
    limit: "500m"
    request: "100m"
  memory:
    limit: "512Mi"
    request: "128Mi"
```

#### Default Environment Variables

```yaml
env:
  - name: NODE_ENV
    value: "production"
  - name: PORT
    value: "3030"
  - name: REDIS_URL
    value: "redis://redis:6379"
  - name: LOG_LEVEL
    value: "info"
```

### Redis Configuration

| Field | Description | Default |
|-------|-------------|---------|
| `replicas` | Number of Redis replicas | 1 |
| `image` | Redis Docker image | redis:7-alpine |
| `resources` | Resource limits and requests | See below |
| `storage` | Storage size | 1Gi |
| `config` | Redis configuration | See below |

#### Default Resource Limits

```yaml
resources:
  cpu:
    limit: "300m"
    request: "100m"
  memory:
    limit: "256Mi"
    request: "128Mi"
```

#### Default Redis Configuration

```yaml
config:
  maxmemory: "256mb"
  maxmemory-policy: "allkeys-lru"
```

### Network Policy Configuration

| Field | Description | Default |
|-------|-------------|---------|
| `enabled` | Whether network policies are enabled | true |
| `allowedNamespaces` | Namespaces that can access the tenant | [] |
| `ingressRules` | Custom ingress rules | [] |
| `egressRules` | Custom egress rules | [] |

## Monitoring and Management

### Viewing Tenant Status

```bash
kubectl get tenants
```

Example output:

```
NAME            STATUS    NAMESPACE             AGE
sample-tenant   Running   tenant-sample-tenant  10m
```

### Viewing Tenant Details

```bash
kubectl describe tenant sample-tenant
```

### Viewing Tenant Logs

```bash
# Get the namespace
NAMESPACE=$(kubectl get tenant sample-tenant -o jsonpath='{.status.namespace}')

# View server logs
kubectl logs -n $NAMESPACE -l app=neurallog-server

# View Redis logs
kubectl logs -n $NAMESPACE -l app=redis
```

### Accessing Tenant Resources

```bash
# Get the namespace
NAMESPACE=$(kubectl get tenant sample-tenant -o jsonpath='{.status.namespace}')

# Port forward to the server
kubectl -n $NAMESPACE port-forward svc/neurallog-server 3030:3030
```

## Troubleshooting

### Common Issues

#### Tenant Stuck in Pending State

If a tenant is stuck in the Pending state, check the operator logs:

```bash
kubectl logs -n system -l control-plane=controller-manager
```

#### Redis Not Starting

If Redis is not starting, check the Redis logs:

```bash
# Get the namespace
NAMESPACE=$(kubectl get tenant sample-tenant -o jsonpath='{.status.namespace}')

# View Redis logs
kubectl logs -n $NAMESPACE -l app=redis
```

#### Server Not Starting

If the server is not starting, check the server logs:

```bash
# Get the namespace
NAMESPACE=$(kubectl get tenant sample-tenant -o jsonpath='{.status.namespace}')

# View server logs
kubectl logs -n $NAMESPACE -l app=neurallog-server
```

#### Network Policies Not Working

If network policies are not working, check the network policy configuration:

```bash
# Get the namespace
NAMESPACE=$(kubectl get tenant sample-tenant -o jsonpath='{.status.namespace}')

# View network policies
kubectl get networkpolicies -n $NAMESPACE
```

### Debugging the Operator

To debug the operator, you can increase the log level:

```bash
kubectl edit deployment -n system controller-manager
```

Add the `--zap-log-level=debug` argument to the container args.

## Advanced Usage

### Custom Environment Variables

You can add custom environment variables to the server:

```yaml
spec:
  server:
    env:
      - name: CUSTOM_VAR
        value: "custom-value"
      - name: SECRET_VAR
        valueFrom:
          secretKeyRef:
            name: my-secret
            key: my-key
```

### Custom Redis Configuration

You can customize the Redis configuration:

```yaml
spec:
  redis:
    config:
      maxmemory: "1gb"
      maxmemory-policy: "volatile-lru"
      appendonly: "yes"
      appendfsync: "everysec"
```

### Custom Network Policies

You can add custom network policies:

```yaml
spec:
  networkPolicy:
    ingressRules:
      - description: "Allow monitoring tools"
        from:
          app: "prometheus"
        ports:
          - protocol: "TCP"
            port: 3030
```

## Development

### Building the Operator

1. Build the operator image:

```bash
cd operator
docker build -t neurallog/tenant-operator:latest .
```

2. Push the image to a registry:

```bash
docker push neurallog/tenant-operator:latest
```

### Running Locally

1. Install the CRDs:

```bash
kubectl apply -f operator/config/crd/bases
```

2. Run the operator locally:

```bash
cd operator
go run main.go
```

### Adding New Features

To add new features to the operator:

1. Update the API in `api/v1/tenant_*.go`
2. Update the controller in `controllers/`
3. Generate the CRD:

```bash
cd operator
make manifests
```

4. Build and deploy the operator:

```bash
docker build -t neurallog/tenant-operator:latest .
docker push neurallog/tenant-operator:latest
kubectl apply -f config/manager/manager.yaml
```
