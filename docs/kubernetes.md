# NeuralLog Kubernetes Configuration Guide

This guide provides detailed information about the Kubernetes configurations used in the NeuralLog infrastructure.

## Table of Contents

- [Overview](#overview)
- [Directory Structure](#directory-structure)
- [Base Resources](#base-resources)
  - [Server Resources](#server-resources)
  - [Redis Resources](#redis-resources)
- [Overlays](#overlays)
  - [Test Overlay](#test-overlay)
  - [Production Overlay](#production-overlay)
- [Kustomize Usage](#kustomize-usage)
- [Resource Management](#resource-management)
- [Network Configuration](#network-configuration)
- [Storage Configuration](#storage-configuration)
- [Security Configuration](#security-configuration)
- [Custom Resources](#custom-resources)
- [Advanced Configuration](#advanced-configuration)

## Overview

NeuralLog uses Kubernetes for container orchestration and management. The Kubernetes configurations are organized using Kustomize, which allows for environment-specific customizations.

### Key Features

- **Kustomize**: Environment-specific customizations
- **Base Resources**: Core resources shared across environments
- **Overlays**: Environment-specific configurations
- **Resource Management**: Configurable resource limits and requests
- **Network Configuration**: Network policies for security
- **Storage Configuration**: Persistent storage for data durability

## Directory Structure

```
kubernetes/
├── base/                 # Base Kubernetes resources
│   ├── kustomization.yaml
│   ├── server/           # Server resources
│   │   ├── deployment.yaml
│   │   └── service.yaml
│   └── redis/            # Redis resources
│       ├── statefulset.yaml
│       ├── service.yaml
│       └── configmap.yaml
└── overlays/             # Kustomize overlays
    ├── test/             # Test environment
    │   ├── kustomization.yaml
    │   └── patches/
    │       └── resource-limits.yaml
    └── production/       # Production environment (example)
        ├── kustomization.yaml
        └── patches/
            └── resource-limits.yaml
```

## Base Resources

The base directory contains the core Kubernetes resources that are shared across all environments.

### Server Resources

#### Deployment

The server deployment defines the NeuralLog server container:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: neurallog-server
  labels:
    app: neurallog-server
spec:
  replicas: 1
  selector:
    matchLabels:
      app: neurallog-server
  template:
    metadata:
      labels:
        app: neurallog-server
    spec:
      containers:
      - name: server
        image: neurallog/server:latest
        ports:
        - containerPort: 3030
          name: http
        env:
        - name: NODE_ENV
          value: "production"
        - name: PORT
          value: "3030"
        - name: REDIS_URL
          value: "redis://redis:6379"
        - name: LOG_LEVEL
          value: "info"
        resources:
          requests:
            cpu: 100m
            memory: 128Mi
          limits:
            cpu: 500m
            memory: 512Mi
        livenessProbe:
          httpGet:
            path: /health
            port: http
          initialDelaySeconds: 15
          periodSeconds: 20
        readinessProbe:
          httpGet:
            path: /health
            port: http
          initialDelaySeconds: 5
          periodSeconds: 10
        volumeMounts:
        - name: tmp-volume
          mountPath: /tmp
      volumes:
      - name: tmp-volume
        emptyDir: {}
```

#### Service

The server service exposes the NeuralLog server:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: neurallog-server
  labels:
    app: neurallog-server
spec:
  selector:
    app: neurallog-server
  ports:
  - port: 3030
    targetPort: http
    name: http
  type: ClusterIP
```

### Redis Resources

#### StatefulSet

The Redis StatefulSet defines the Redis container:

```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: redis
  labels:
    app: redis
spec:
  serviceName: redis
  replicas: 1
  selector:
    matchLabels:
      app: redis
  template:
    metadata:
      labels:
        app: redis
    spec:
      containers:
      - name: redis
        image: redis:7-alpine
        command:
        - redis-server
        - /etc/redis/redis.conf
        ports:
        - containerPort: 6379
          name: redis
        volumeMounts:
        - name: redis-data
          mountPath: /data
        - name: redis-config
          mountPath: /etc/redis
        resources:
          requests:
            cpu: 100m
            memory: 128Mi
          limits:
            cpu: 300m
            memory: 256Mi
        livenessProbe:
          tcpSocket:
            port: redis
          initialDelaySeconds: 15
          periodSeconds: 20
        readinessProbe:
          exec:
            command:
            - redis-cli
            - ping
          initialDelaySeconds: 5
          periodSeconds: 10
      volumes:
      - name: redis-config
        configMap:
          name: redis-config
          items:
          - key: redis.conf
            path: redis.conf
  volumeClaimTemplates:
  - metadata:
      name: redis-data
    spec:
      accessModes: [ "ReadWriteOnce" ]
      resources:
        requests:
          storage: 1Gi
```

#### Service

The Redis service exposes the Redis instance:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: redis
  labels:
    app: redis
spec:
  selector:
    app: redis
  ports:
  - port: 6379
    targetPort: redis
    name: redis
  clusterIP: None
```

#### ConfigMap

The Redis ConfigMap contains the Redis configuration:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: redis-config
  labels:
    app: redis
data:
  redis.conf: |
    # Redis configuration for NeuralLog
    port 6379
    bind 0.0.0.0
    protected-mode yes
    daemonize no
    
    # Memory management
    maxmemory 256mb
    maxmemory-policy allkeys-lru
    
    # Persistence
    appendonly yes
    appendfsync everysec
    
    # Logging
    loglevel notice
    logfile ""
```

## Overlays

Overlays provide environment-specific customizations using Kustomize.

### Test Overlay

The test overlay provides configurations for the test environment:

#### Kustomization

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
- ../../base

namespace: neurallog

patches:
- path: patches/resource-limits.yaml
  target:
    kind: Deployment
    name: neurallog-server

images:
- name: neurallog/server
  newTag: latest
```

#### Resource Limits Patch

```yaml
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
            cpu: 50m
            memory: 64Mi
          limits:
            cpu: 200m
            memory: 256Mi
```

### Production Overlay

The production overlay provides configurations for the production environment:

#### Kustomization

```yaml
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
```

#### Resource Limits Patch

```yaml
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
```

## Kustomize Usage

Kustomize is used to manage environment-specific configurations:

### Applying Base Resources

```bash
kubectl apply -k kubernetes/base
```

### Applying Test Overlay

```bash
kubectl apply -k kubernetes/overlays/test
```

### Applying Production Overlay

```bash
kubectl apply -k kubernetes/overlays/production
```

### Building Manifests

```bash
kubectl kustomize kubernetes/overlays/test > test-manifests.yaml
```

## Resource Management

Resource management is configured using resource limits and requests:

### CPU Resources

- **Requests**: Minimum CPU allocation
- **Limits**: Maximum CPU allocation

### Memory Resources

- **Requests**: Minimum memory allocation
- **Limits**: Maximum memory allocation

### Storage Resources

- **Requests**: Minimum storage allocation

## Network Configuration

Network configuration is managed using Kubernetes Services and NetworkPolicies:

### Services

- **ClusterIP**: Internal service for the server
- **Headless**: Headless service for Redis

### NetworkPolicies

Network policies are used to restrict communication between components:

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: default-deny
spec:
  podSelector: {}
  policyTypes:
  - Ingress
```

## Storage Configuration

Storage is configured using PersistentVolumeClaims:

### Redis Storage

```yaml
volumeClaimTemplates:
- metadata:
    name: redis-data
  spec:
    accessModes: [ "ReadWriteOnce" ]
    resources:
      requests:
        storage: 1Gi
```

## Security Configuration

Security is configured using various Kubernetes features:

### Pod Security

- **Non-Root User**: Containers run as non-root users
- **Read-Only Filesystem**: Containers use read-only filesystems where possible

### Network Security

- **Network Policies**: Network policies restrict communication between components

### Secret Management

- **Kubernetes Secrets**: Sensitive information is stored in Kubernetes Secrets

## Custom Resources

Custom resources are defined using CustomResourceDefinitions (CRDs):

### Tenant CRD

```yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: tenants.neurallog.io
spec:
  group: neurallog.io
  names:
    kind: Tenant
    plural: tenants
    singular: tenant
  scope: Cluster
  versions:
  - name: v1
    served: true
    storage: true
    schema:
      openAPIV3Schema:
        type: object
        properties:
          spec:
            type: object
            properties:
              # Tenant specification
          status:
            type: object
            properties:
              # Tenant status
```

## Advanced Configuration

### Horizontal Pod Autoscaling

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: neurallog-server
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: neurallog-server
  minReplicas: 1
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 80
```

### Pod Disruption Budget

```yaml
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: neurallog-server
spec:
  minAvailable: 1
  selector:
    matchLabels:
      app: neurallog-server
```

### Affinity and Anti-Affinity

```yaml
spec:
  template:
    spec:
      affinity:
        podAntiAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
          - weight: 100
            podAffinityTerm:
              labelSelector:
                matchExpressions:
                - key: app
                  operator: In
                  values:
                  - neurallog-server
              topologyKey: kubernetes.io/hostname
```
