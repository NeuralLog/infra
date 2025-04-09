# NeuralLog Tenant API Reference

This document provides a comprehensive reference for the NeuralLog Tenant API.

## Table of Contents

- [Tenant Resource](#tenant-resource)
- [Spec Fields](#spec-fields)
  - [Basic Fields](#basic-fields)
  - [Upgrade Strategy](#upgrade-strategy)
  - [Resources](#resources)
  - [Server Configuration](#server-configuration)
  - [Redis Configuration](#redis-configuration)
  - [Network Policy Configuration](#network-policy-configuration)
  - [Monitoring Configuration](#monitoring-configuration)
  - [Backup Configuration](#backup-configuration)
  - [Lifecycle Hooks](#lifecycle-hooks)
  - [Integrations](#integrations)
- [Status Fields](#status-fields)
  - [Basic Status Fields](#basic-status-fields)
  - [Component Status](#component-status)
  - [URL Status](#url-status)
  - [Metrics Status](#metrics-status)
  - [Backup Status](#backup-status)
- [Examples](#examples)
  - [Basic Tenant](#basic-tenant)
  - [Production Tenant](#production-tenant)
  - [High Availability Tenant](#high-availability-tenant)
  - [Secure Tenant](#secure-tenant)
- [Field Validation Rules](#field-validation-rules)
- [API Versioning](#api-versioning)

## Tenant Resource

The Tenant resource is a Kubernetes custom resource that represents a NeuralLog tenant.

```yaml
apiVersion: neurallog.io/v1
kind: Tenant
metadata:
  name: tenant-name
spec:
  # Tenant specification
status:
  # Tenant status
```

## Spec Fields

### Basic Fields

| Field | Type | Description | Default |
|-------|------|-------------|---------|
| `displayName` | string | User-friendly name for the tenant | - |
| `description` | string | Additional information about the tenant | - |
| `version` | string | Version of NeuralLog to deploy | latest |

### Upgrade Strategy

| Field | Type | Description | Default |
|-------|------|-------------|---------|
| `upgradeStrategy.type` | string | Type of upgrade strategy (RollingUpdate or Recreate) | RollingUpdate |
| `upgradeStrategy.maxUnavailable` | string | Maximum number of pods that can be unavailable during the update | 25% |
| `upgradeStrategy.maxSurge` | string | Maximum number of pods that can be scheduled above the desired number of pods | 25% |

### Resources

| Field | Type | Description | Default |
|-------|------|-------------|---------|
| `resources.cpu.limit` | string | Maximum CPU allocation for the tenant | - |
| `resources.cpu.request` | string | Minimum CPU allocation for the tenant | - |
| `resources.memory.limit` | string | Maximum memory allocation for the tenant | - |
| `resources.memory.request` | string | Minimum memory allocation for the tenant | - |
| `resources.storage.limit` | string | Maximum storage allocation for the tenant | - |
| `resources.storage.request` | string | Minimum storage allocation for the tenant | - |

### Server Configuration

#### Basic Server Configuration

| Field | Type | Description | Default |
|-------|------|-------------|---------|
| `server.replicas` | integer | Number of server replicas | 1 |
| `server.image` | string | Docker image for the server | neurallog/server:latest |
| `server.resources.cpu.limit` | string | Maximum CPU allocation for the server | 500m |
| `server.resources.cpu.request` | string | Minimum CPU allocation for the server | 100m |
| `server.resources.memory.limit` | string | Maximum memory allocation for the server | 512Mi |
| `server.resources.memory.request` | string | Minimum memory allocation for the server | 128Mi |
| `server.env` | array | Environment variables for the server | - |
| `server.logLevel` | string | Log level for the server (debug, info, warn, error) | info |

#### Server Deployment Configuration

| Field | Type | Description | Default |
|-------|------|-------------|---------|
| `server.deployment.strategy` | string | Deployment strategy (RollingUpdate or Recreate) | RollingUpdate |
| `server.deployment.rollingUpdate.maxUnavailable` | string | Maximum number of pods that can be unavailable during the update | 25% |
| `server.deployment.rollingUpdate.maxSurge` | string | Maximum number of pods that can be scheduled above the desired number of pods | 25% |

#### Server Autoscaling Configuration

| Field | Type | Description | Default |
|-------|------|-------------|---------|
| `server.autoscaling.enabled` | boolean | Whether autoscaling is enabled | false |
| `server.autoscaling.minReplicas` | integer | Minimum number of replicas | 1 |
| `server.autoscaling.maxReplicas` | integer | Maximum number of replicas | 10 |
| `server.autoscaling.targetCPUUtilizationPercentage` | integer | Target CPU utilization percentage | 80 |
| `server.autoscaling.targetMemoryUtilizationPercentage` | integer | Target memory utilization percentage | - |

#### Server Affinity Configuration

| Field | Type | Description | Default |
|-------|------|-------------|---------|
| `server.affinity.nodeAffinity` | object | Node affinity configuration | - |
| `server.affinity.podAffinity` | object | Pod affinity configuration | - |
| `server.affinity.podAntiAffinity` | object | Pod anti-affinity configuration | - |

#### Server Security Context Configuration

| Field | Type | Description | Default |
|-------|------|-------------|---------|
| `server.securityContext.runAsUser` | integer | User ID to run as | - |
| `server.securityContext.runAsGroup` | integer | Group ID to run as | - |
| `server.securityContext.runAsNonRoot` | boolean | Whether to run as a non-root user | true |
| `server.securityContext.readOnlyRootFilesystem` | boolean | Whether to use a read-only root filesystem | true |
| `server.securityContext.allowPrivilegeEscalation` | boolean | Whether to allow privilege escalation | false |
| `server.securityContext.capabilities.drop` | array | Capabilities to drop | ["ALL"] |

#### Server Probe Configuration

| Field | Type | Description | Default |
|-------|------|-------------|---------|
| `server.probes.liveness.path` | string | HTTP path for liveness probe | /health |
| `server.probes.liveness.port` | integer | Port for liveness probe | 3030 |
| `server.probes.liveness.initialDelaySeconds` | integer | Initial delay for liveness probe | 15 |
| `server.probes.liveness.periodSeconds` | integer | Period for liveness probe | 10 |
| `server.probes.liveness.timeoutSeconds` | integer | Timeout for liveness probe | 5 |
| `server.probes.liveness.successThreshold` | integer | Success threshold for liveness probe | 1 |
| `server.probes.liveness.failureThreshold` | integer | Failure threshold for liveness probe | 3 |
| `server.probes.readiness` | object | Readiness probe configuration | - |
| `server.probes.startup` | object | Startup probe configuration | - |

#### Server API Configuration

| Field | Type | Description | Default |
|-------|------|-------------|---------|
| `server.api.authentication.enabled` | boolean | Whether authentication is enabled | true |
| `server.api.authentication.type` | string | Authentication type (apiKey, jwt, oauth2) | apiKey |
| `server.api.authentication.apiKey` | object | API key configuration | - |
| `server.api.authentication.jwt` | object | JWT configuration | - |
| `server.api.authentication.oauth2` | object | OAuth2 configuration | - |
| `server.api.cors.enabled` | boolean | Whether CORS is enabled | true |
| `server.api.cors.allowOrigins` | array | Allowed origins | ["*"] |
| `server.api.cors.allowMethods` | array | Allowed methods | ["GET", "POST", "PUT", "DELETE", "OPTIONS"] |
| `server.api.cors.allowHeaders` | array | Allowed headers | ["Content-Type", "Authorization"] |
| `server.api.cors.exposeHeaders` | array | Exposed headers | - |
| `server.api.cors.maxAge` | integer | Max age in seconds | 86400 |
| `server.api.rateLimit.enabled` | boolean | Whether rate limiting is enabled | true |
| `server.api.rateLimit.requestsPerSecond` | integer | Number of requests per second | 100 |
| `server.api.rateLimit.burstSize` | integer | Burst size | 200 |

### Redis Configuration

#### Basic Redis Configuration

| Field | Type | Description | Default |
|-------|------|-------------|---------|
| `redis.replicas` | integer | Number of Redis replicas | 1 |
| `redis.image` | string | Docker image for Redis | redis:7-alpine |
| `redis.resources.cpu.limit` | string | Maximum CPU allocation for Redis | 300m |
| `redis.resources.cpu.request` | string | Minimum CPU allocation for Redis | 100m |
| `redis.resources.memory.limit` | string | Maximum memory allocation for Redis | 256Mi |
| `redis.resources.memory.request` | string | Minimum memory allocation for Redis | 128Mi |
| `redis.storage` | string | Storage size for Redis | 1Gi |
| `redis.config` | object | Additional Redis configuration | - |
| `redis.mode` | string | Redis mode (standalone, sentinel, cluster) | standalone |

#### Redis Persistence Configuration

| Field | Type | Description | Default |
|-------|------|-------------|---------|
| `redis.persistence.type` | string | Persistence type (aof, rdb, both) | aof |
| `redis.persistence.fsync` | string | Fsync policy (everysec, always, no) | everysec |
| `redis.persistence.savePoints` | array | RDB save points | ["900 1", "300 10", "60 10000"] |

#### Redis Security Configuration

| Field | Type | Description | Default |
|-------|------|-------------|---------|
| `redis.security.enabled` | boolean | Whether security is enabled | true |
| `redis.security.authSecret` | object | Reference to a secret containing the Redis password | - |
| `redis.security.tls.enabled` | boolean | Whether TLS is enabled | false |
| `redis.security.tls.certSecret` | object | Reference to a secret containing the TLS certificate | - |

#### Redis Advanced Configuration

| Field | Type | Description | Default |
|-------|------|-------------|---------|
| `redis.advanced.maxmemoryPolicy` | string | Max memory policy | allkeys-lru |
| `redis.advanced.clientOutputBufferLimit` | string | Client output buffer limit | normal 0 0 0 |
| `redis.advanced.databases` | integer | Number of databases | 16 |
| `redis.advanced.tcpKeepalive` | integer | TCP keepalive interval in seconds | 300 |

### Network Policy Configuration

| Field | Type | Description | Default |
|-------|------|-------------|---------|
| `networkPolicy.enabled` | boolean | Whether network policies are enabled | true |
| `networkPolicy.allowedNamespaces` | array | Namespaces that can access the tenant | - |
| `networkPolicy.ingressRules` | array | Custom ingress rules | - |
| `networkPolicy.egressRules` | array | Custom egress rules | - |

### Monitoring Configuration

| Field | Type | Description | Default |
|-------|------|-------------|---------|
| `monitoring.enabled` | boolean | Whether monitoring is enabled | true |
| `monitoring.prometheus.scrape` | boolean | Whether Prometheus should scrape metrics | true |
| `monitoring.prometheus.port` | integer | Port to scrape metrics from | 3030 |
| `monitoring.prometheus.path` | string | Path to scrape metrics from | /metrics |
| `monitoring.prometheus.interval` | string | Interval at which to scrape metrics | 15s |
| `monitoring.alerts.enabled` | boolean | Whether alerts are enabled | true |
| `monitoring.alerts.receivers` | array | Alert receivers | - |

### Backup Configuration

| Field | Type | Description | Default |
|-------|------|-------------|---------|
| `backup.enabled` | boolean | Whether backups are enabled | true |
| `backup.schedule` | string | Cron schedule for backups | 0 2 * * * |
| `backup.retention.count` | integer | Number of backups to retain | 7 |
| `backup.retention.days` | integer | Number of days to retain backups | 30 |
| `backup.storage.type` | string | Type of storage (s3, gcs, azure, local) | local |
| `backup.storage.bucket` | string | Storage bucket | - |
| `backup.storage.prefix` | string | Storage prefix | - |
| `backup.storage.secretRef` | object | Reference to a secret containing storage credentials | - |

### Lifecycle Hooks

| Field | Type | Description | Default |
|-------|------|-------------|---------|
| `lifecycle.preCreate` | object | Hook executed before the tenant is created | - |
| `lifecycle.postCreate` | object | Hook executed after the tenant is created | - |
| `lifecycle.preDelete` | object | Hook executed before the tenant is deleted | - |
| `lifecycle.postDelete` | object | Hook executed after the tenant is deleted | - |

### Integrations

| Field | Type | Description | Default |
|-------|------|-------------|---------|
| `integrations.ingress.enabled` | boolean | Whether ingress is enabled | false |
| `integrations.ingress.annotations` | object | Annotations for the ingress | - |
| `integrations.ingress.hosts` | array | Ingress hosts | - |
| `integrations.ingress.tls` | array | Ingress TLS configuration | - |
| `integrations.serviceMesh.enabled` | boolean | Whether service mesh is enabled | false |
| `integrations.serviceMesh.type` | string | Type of service mesh (istio, linkerd, consul) | - |

## Status Fields

### Basic Status Fields

| Field | Type | Description |
|-------|------|-------------|
| `status.conditions` | array | Latest available observations of the tenant's state |
| `status.phase` | string | Current phase of the tenant (Pending, Provisioning, Running, Failed, Terminating) |
| `status.namespace` | string | Namespace created for the tenant |
| `status.observedGeneration` | integer | Most recent generation observed by the controller |
| `status.lastReconcileTime` | string | Last time the tenant was reconciled |

### Component Status

| Field | Type | Description |
|-------|------|-------------|
| `status.serverStatus` | object | Status of the server deployment |
| `status.redisStatus` | object | Status of the Redis deployment |
| `status.components` | object | Status of all components |

Component status fields:

| Field | Type | Description |
|-------|------|-------------|
| `phase` | string | Current phase of the component (Pending, Running, Failed) |
| `message` | string | Additional information about the component status |
| `readyReplicas` | integer | Number of ready replicas |
| `totalReplicas` | integer | Total number of replicas |
| `lastTransitionTime` | string | Last time the component phase changed |
| `url` | string | URL for accessing the component |
| `version` | string | Version of the component |
| `metrics` | object | Component metrics |

### URL Status

| Field | Type | Description |
|-------|------|-------------|
| `status.urls.server` | string | URL for accessing the server |
| `status.urls.api` | string | URL for accessing the API |
| `status.urls.dashboard` | string | URL for accessing the dashboard |

### Metrics Status

| Field | Type | Description |
|-------|------|-------------|
| `status.metrics.cpu` | string | CPU usage in millicores |
| `status.metrics.memory` | string | Memory usage |
| `status.metrics.storage` | string | Storage usage |
| `status.metrics.connections` | integer | Number of connections |
| `status.metrics.requestsPerSecond` | number | Number of requests per second |
| `status.metrics.averageResponseTime` | number | Average response time in milliseconds |

### Backup Status

| Field | Type | Description |
|-------|------|-------------|
| `status.backupStatus.lastBackupTime` | string | Last time a backup was taken |
| `status.backupStatus.lastBackupSize` | string | Size of the last backup |
| `status.backupStatus.lastBackupStatus` | string | Status of the last backup |
| `status.backupStatus.backupCount` | integer | Number of backups |
| `status.backupStatus.nextBackupTime` | string | Next scheduled backup time |

## Examples

### Basic Tenant

```yaml
apiVersion: neurallog.io/v1
kind: Tenant
metadata:
  name: basic-tenant
spec:
  displayName: "Basic Tenant"
  description: "A basic tenant configuration"
```

### Production Tenant

```yaml
apiVersion: neurallog.io/v1
kind: Tenant
metadata:
  name: production-tenant
spec:
  displayName: "Production Tenant"
  description: "A production tenant configuration"
  version: "1.2.3"
  
  resources:
    cpu:
      limit: "4"
      request: "2"
    memory:
      limit: "8Gi"
      request: "4Gi"
  
  server:
    replicas: 3
    image: "neurallog/server:stable"
    resources:
      cpu:
        limit: "2"
        request: "1"
      memory:
        limit: "4Gi"
        request: "2Gi"
  
  redis:
    replicas: 3
    image: "redis:7-alpine"
    storage: "20Gi"
    resources:
      cpu:
        limit: "2"
        request: "1"
      memory:
        limit: "4Gi"
        request: "2Gi"
  
  monitoring:
    enabled: true
    prometheus:
      scrape: true
  
  backup:
    enabled: true
    schedule: "0 2 * * *"
    retention:
      count: 30
      days: 90
```

### High Availability Tenant

```yaml
apiVersion: neurallog.io/v1
kind: Tenant
metadata:
  name: ha-tenant
spec:
  displayName: "High Availability Tenant"
  description: "A high availability tenant configuration"
  
  server:
    replicas: 5
    autoscaling:
      enabled: true
      minReplicas: 3
      maxReplicas: 10
      targetCPUUtilizationPercentage: 70
    
    affinity:
      podAntiAffinity:
        requiredDuringSchedulingIgnoredDuringExecution:
        - labelSelector:
            matchLabels:
              app: neurallog-server
          topologyKey: "kubernetes.io/hostname"
  
  redis:
    replicas: 3
    mode: "sentinel"
    persistence:
      type: "both"
      fsync: "everysec"
```

### Secure Tenant

```yaml
apiVersion: neurallog.io/v1
kind: Tenant
metadata:
  name: secure-tenant
spec:
  displayName: "Secure Tenant"
  description: "A secure tenant configuration"
  
  server:
    securityContext:
      runAsNonRoot: true
      readOnlyRootFilesystem: true
      allowPrivilegeEscalation: false
      capabilities:
        drop:
        - ALL
    
    api:
      authentication:
        enabled: true
        type: "jwt"
        jwt:
          issuer: "neurallog"
          audience: "api"
          secretRef:
            name: jwt-secret
            key: jwt-key
      
      rateLimit:
        enabled: true
        requestsPerSecond: 50
  
  redis:
    security:
      enabled: true
      authSecret:
        name: redis-secret
        key: redis-password
      tls:
        enabled: true
        certSecret:
          name: redis-tls
          key: tls.crt
  
  networkPolicy:
    enabled: true
    allowedNamespaces:
      - "monitoring"
    ingressRules:
      - description: "Allow API access"
        from:
          app: "api-gateway"
        ports:
          - protocol: "TCP"
            port: 3030
  
  integrations:
    ingress:
      enabled: true
      annotations:
        kubernetes.io/ingress.class: "nginx"
        cert-manager.io/cluster-issuer: "letsencrypt-prod"
      hosts:
        - host: "secure-tenant.example.com"
          paths:
            - path: "/"
              pathType: "Prefix"
      tls:
        - secretName: "secure-tenant-tls"
          hosts:
            - "secure-tenant.example.com"
```

## Field Validation Rules

### Basic Fields

| Field | Validation |
|-------|------------|
| `displayName` | Max length: 63, Pattern: `^[a-zA-Z0-9]([a-zA-Z0-9\-\_\.]*[a-zA-Z0-9])?$` |
| `version` | Pattern: `^v?([0-9]+)(\.[0-9]+)?(\.[0-9]+)?(-([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?(\+([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?$` |

### Server Configuration

| Field | Validation |
|-------|------------|
| `server.deployment.strategy` | Enum: RollingUpdate, Recreate |
| `server.logLevel` | Enum: debug, info, warn, error |
| `server.api.authentication.type` | Enum: apiKey, jwt, oauth2 |
| `server.api.cors.maxAge` | Minimum: 0 |
| `server.api.rateLimit.requestsPerSecond` | Minimum: 1 |
| `server.api.rateLimit.burstSize` | Minimum: 1 |

### Redis Configuration

| Field | Validation |
|-------|------------|
| `redis.mode` | Enum: standalone, sentinel, cluster |
| `redis.persistence.type` | Enum: aof, rdb, both |
| `redis.persistence.fsync` | Enum: everysec, always, no |
| `redis.advanced.maxmemoryPolicy` | Enum: allkeys-lru, volatile-lru, allkeys-random, volatile-random, volatile-ttl, noeviction |
| `redis.advanced.databases` | Minimum: 1 |
| `redis.advanced.tcpKeepalive` | Minimum: 0 |

### Backup Configuration

| Field | Validation |
|-------|------------|
| `backup.storage.type` | Enum: s3, gcs, azure, local |
| `backup.retention.count` | Minimum: 1 |
| `backup.retention.days` | Minimum: 1 |

### Integrations

| Field | Validation |
|-------|------------|
| `integrations.ingress.hosts[].paths[].pathType` | Enum: Exact, Prefix, ImplementationSpecific |
| `integrations.serviceMesh.type` | Enum: istio, linkerd, consul |

## API Versioning

The NeuralLog Tenant API follows Kubernetes API versioning conventions:

- **Alpha Versions** (v1alpha1): May be buggy and are disabled by default
- **Beta Versions** (v1beta1): Well-tested but may have minor changes
- **Stable Versions** (v1): Stable and will not change in incompatible ways

### Version Support Matrix

| Version | Status | Supported Until |
|---------|--------|-----------------|
| v1alpha1 | Deprecated | 2023-01-01 |
| v1beta1 | Deprecated | 2023-06-01 |
| v1 | Current | - |
