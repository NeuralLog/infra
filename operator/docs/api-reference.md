# NeuralLog Tenant Operator API Reference

This document provides a detailed reference for the NeuralLog Tenant Operator API.

## Tenant

The `Tenant` custom resource is the primary API object for the NeuralLog Tenant Operator. It represents a tenant in the NeuralLog platform and defines the configuration for tenant resources.

### API Group and Version

```
apiVersion: neurallog.io/v1
kind: Tenant
```

### Spec

The `spec` field defines the desired state of the tenant.

#### Top-Level Fields

| Field | Type | Description | Required |
|-------|------|-------------|----------|
| `displayName` | string | A user-friendly name for the tenant | No |
| `description` | string | A description of the tenant | No |
| `resources` | [ResourceRequirements](#resourcerequirements) | Resource limits and requests for the tenant | No |
| `server` | [ServerSpec](#serverspec) | Configuration for the NeuralLog server | No |
| `redis` | [RedisSpec](#redisspec) | Configuration for the Redis instance | No |
| `networkPolicy` | [NetworkPolicySpec](#networkpolicyspec) | Configuration for network policies | No |

#### ResourceRequirements

The `resources` field defines resource limits and requests for the tenant.

| Field | Type | Description | Required |
|-------|------|-------------|----------|
| `cpu` | [ResourceLimit](#resourcelimit) | CPU limits and requests | No |
| `memory` | [ResourceLimit](#resourcelimit) | Memory limits and requests | No |
| `storage` | [ResourceLimit](#resourcelimit) | Storage limits and requests | No |

#### ResourceLimit

The `resourceLimit` field defines a resource limit and request.

| Field | Type | Description | Required |
|-------|------|-------------|----------|
| `limit` | string | The maximum amount of the resource | No |
| `request` | string | The minimum amount of the resource | No |

#### ServerSpec

The `server` field defines the configuration for the NeuralLog server.

| Field | Type | Description | Required |
|-------|------|-------------|----------|
| `replicas` | int32 | The number of server instances | No |
| `image` | string | The Docker image for the server | No |
| `resources` | [ResourceRequirements](#resourcerequirements) | Resource limits and requests for the server | No |
| `env` | [][EnvVar](#envvar) | Environment variables for the server | No |

#### RedisSpec

The `redis` field defines the configuration for the Redis instance.

| Field | Type | Description | Required |
|-------|------|-------------|----------|
| `replicas` | int32 | The number of Redis instances (for Redis Sentinel) | No |
| `image` | string | The Docker image for Redis | No |
| `resources` | [ResourceRequirements](#resourcerequirements) | Resource limits and requests for Redis | No |
| `storage` | string | The storage configuration for Redis | No |
| `config` | map[string]string | Additional Redis configuration | No |

#### NetworkPolicySpec

The `networkPolicy` field defines the network policy configuration for the tenant.

| Field | Type | Description | Required |
|-------|------|-------------|----------|
| `enabled` | bool | Whether network policies should be created | No |
| `allowedNamespaces` | []string | A list of namespaces that can access the tenant | No |
| `ingressRules` | [][NetworkPolicyRule](#networkpolicyrule) | Additional ingress rules | No |
| `egressRules` | [][NetworkPolicyRule](#networkpolicyrule) | Additional egress rules | No |

#### NetworkPolicyRule

The `networkPolicyRule` field defines a network policy rule.

| Field | Type | Description | Required |
|-------|------|-------------|----------|
| `description` | string | A description of the rule | No |
| `from` | map[string]string | The source selector for ingress rules | No |
| `to` | map[string]string | The destination selector for egress rules | No |
| `ports` | [][NetworkPolicyPort](#networkpolicyport) | The ports for the rule | No |

#### NetworkPolicyPort

The `networkPolicyPort` field defines a port for a network policy rule.

| Field | Type | Description | Required |
|-------|------|-------------|----------|
| `protocol` | string | The protocol for the port | No |
| `port` | int32 | The port number | No |

#### EnvVar

The `envVar` field defines an environment variable.

| Field | Type | Description | Required |
|-------|------|-------------|----------|
| `name` | string | The name of the environment variable | Yes |
| `value` | string | The value of the environment variable | No |
| `valueFrom` | [EnvVarSource](#envvarsource) | A source for the environment variable value | No |

#### EnvVarSource

The `envVarSource` field defines a source for an environment variable.

| Field | Type | Description | Required |
|-------|------|-------------|----------|
| `configMapKeyRef` | [ConfigMapKeySelector](#configmapkeyselector) | A reference to a key in a ConfigMap | No |
| `secretKeyRef` | [SecretKeySelector](#secretkeyselector) | A reference to a key in a Secret | No |

#### ConfigMapKeySelector

The `configMapKeySelector` field references a key in a ConfigMap.

| Field | Type | Description | Required |
|-------|------|-------------|----------|
| `name` | string | The name of the ConfigMap | Yes |
| `key` | string | The key in the ConfigMap | Yes |
| `optional` | bool | Whether the ConfigMap or key must exist | No |

#### SecretKeySelector

The `secretKeySelector` field references a key in a Secret.

| Field | Type | Description | Required |
|-------|------|-------------|----------|
| `name` | string | The name of the Secret | Yes |
| `key` | string | The key in the Secret | Yes |
| `optional` | bool | Whether the Secret or key must exist | No |

### Status

The `status` field represents the observed state of the tenant.

| Field | Type | Description |
|-------|------|-------------|
| `conditions` | []metav1.Condition | The latest available observations of the tenant's state |
| `phase` | [TenantPhase](#tenantphase) | The current phase of the tenant |
| `namespace` | string | The namespace created for the tenant |
| `serverStatus` | [ComponentStatus](#componentstatus) | The status of the server deployment |
| `redisStatus` | [ComponentStatus](#componentstatus) | The status of the Redis deployment |

#### TenantPhase

The `tenantPhase` field represents the phase of a tenant.

| Value | Description |
|-------|-------------|
| `Pending` | The tenant is being created |
| `Provisioning` | The tenant resources are being provisioned |
| `Running` | The tenant is running |
| `Failed` | The tenant creation failed |
| `Terminating` | The tenant is being deleted |

#### ComponentStatus

The `componentStatus` field represents the status of a component.

| Field | Type | Description |
|-------|------|-------------|
| `phase` | [ComponentPhase](#componentphase) | The current phase of the component |
| `message` | string | Additional information about the component status |
| `readyReplicas` | int32 | The number of ready replicas |
| `totalReplicas` | int32 | The total number of replicas |

#### ComponentPhase

The `componentPhase` field represents the phase of a component.

| Value | Description |
|-------|-------------|
| `Pending` | The component is being created |
| `Provisioning` | The component is being provisioned |
| `Running` | The component is running |
| `Degraded` | The component is running but not all replicas are ready |
| `Failed` | The component creation failed |

## Examples

### Basic Tenant

```yaml
apiVersion: neurallog.io/v1
kind: Tenant
metadata:
  name: example-tenant
spec:
  displayName: Example Tenant
  description: An example tenant for demonstration purposes
```

### Tenant with Server Configuration

```yaml
apiVersion: neurallog.io/v1
kind: Tenant
metadata:
  name: example-tenant
spec:
  displayName: Example Tenant
  description: An example tenant for demonstration purposes
  server:
    replicas: 2
    image: neurallog/server:latest
    resources:
      cpu:
        request: 100m
        limit: 500m
      memory:
        request: 128Mi
        limit: 512Mi
    env:
      - name: LOG_LEVEL
        value: debug
```

### Tenant with Redis Configuration

```yaml
apiVersion: neurallog.io/v1
kind: Tenant
metadata:
  name: example-tenant
spec:
  displayName: Example Tenant
  description: An example tenant for demonstration purposes
  redis:
    replicas: 1
    image: redis:7-alpine
    resources:
      cpu:
        request: 100m
        limit: 300m
      memory:
        request: 128Mi
        limit: 256Mi
    storage: 1Gi
    config:
      maxmemory-policy: allkeys-lru
```

### Tenant with Network Policy Configuration

```yaml
apiVersion: neurallog.io/v1
kind: Tenant
metadata:
  name: example-tenant
spec:
  displayName: Example Tenant
  description: An example tenant for demonstration purposes
  networkPolicy:
    enabled: true
    allowedNamespaces:
      - default
      - monitoring
    ingressRules:
      - description: Allow ingress from app namespace
        from:
          app: web-frontend
        ports:
          - protocol: TCP
            port: 80
    egressRules:
      - description: Allow egress to database
        to:
          app: database
        ports:
          - protocol: TCP
            port: 5432
```

### Complete Tenant Example

```yaml
apiVersion: neurallog.io/v1
kind: Tenant
metadata:
  name: example-tenant
spec:
  displayName: Example Tenant
  description: An example tenant for demonstration purposes
  server:
    replicas: 2
    image: neurallog/server:latest
    resources:
      cpu:
        request: 100m
        limit: 500m
      memory:
        request: 128Mi
        limit: 512Mi
    env:
      - name: LOG_LEVEL
        value: debug
      - name: API_KEY
        valueFrom:
          secretKeyRef:
            name: api-keys
            key: example-tenant
  redis:
    replicas: 1
    image: redis:7-alpine
    resources:
      cpu:
        request: 100m
        limit: 300m
      memory:
        request: 128Mi
        limit: 256Mi
    storage: 1Gi
    config:
      maxmemory-policy: allkeys-lru
  networkPolicy:
    enabled: true
    allowedNamespaces:
      - default
      - monitoring
    ingressRules:
      - description: Allow ingress from app namespace
        from:
          app: web-frontend
        ports:
          - protocol: TCP
            port: 80
    egressRules:
      - description: Allow egress to database
        to:
          app: database
        ports:
          - protocol: TCP
            port: 5432
```
