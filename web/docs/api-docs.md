# NeuralLog Admin API Documentation

This document describes the API endpoints provided by the NeuralLog Admin interface for managing tenants.

## API Overview

The NeuralLog Admin API is a RESTful API that allows you to manage tenants in your Kubernetes cluster. The API is built on top of the Kubernetes API and uses the Tenant custom resource definition (CRD).

### Base URL

The API is available at the following base URL:

```
/api
```

### Authentication

The API uses the same authentication mechanism as the Kubernetes API. When running in a Kubernetes cluster, it uses the service account token. When running locally, it uses the kubeconfig file.

### Response Format

All API responses are in JSON format. Successful responses have a 2xx status code, while error responses have a 4xx or 5xx status code.

Error responses have the following format:

```json
{
  "error": "Error message"
}
```

## Tenant Endpoints

### List Tenants

Retrieves a list of all tenants.

**Endpoint:** `GET /api/tenants`

**Response:**

```json
{
  "items": [
    {
      "metadata": {
        "name": "tenant-1",
        "creationTimestamp": "2023-01-01T00:00:00Z"
      },
      "spec": {
        "displayName": "Tenant 1",
        "description": "Description of Tenant 1",
        "server": {
          "replicas": 1,
          "image": "neurallog/server:latest"
        },
        "redis": {
          "replicas": 1,
          "image": "redis:7-alpine",
          "storage": "1Gi"
        },
        "networkPolicy": {
          "enabled": true
        }
      },
      "status": {
        "phase": "Running",
        "namespace": "tenant-tenant-1",
        "serverStatus": {
          "phase": "Running",
          "readyReplicas": 1,
          "totalReplicas": 1
        },
        "redisStatus": {
          "phase": "Running",
          "readyReplicas": 1,
          "totalReplicas": 1
        }
      }
    }
  ]
}
```

### Get Tenant

Retrieves a specific tenant by name.

**Endpoint:** `GET /api/tenants/{name}`

**Parameters:**
- `name` (path): The name of the tenant.

**Response:**

```json
{
  "metadata": {
    "name": "tenant-1",
    "creationTimestamp": "2023-01-01T00:00:00Z"
  },
  "spec": {
    "displayName": "Tenant 1",
    "description": "Description of Tenant 1",
    "server": {
      "replicas": 1,
      "image": "neurallog/server:latest"
    },
    "redis": {
      "replicas": 1,
      "image": "redis:7-alpine",
      "storage": "1Gi"
    },
    "networkPolicy": {
      "enabled": true
    }
  },
  "status": {
    "phase": "Running",
    "namespace": "tenant-tenant-1",
    "serverStatus": {
      "phase": "Running",
      "readyReplicas": 1,
      "totalReplicas": 1
    },
    "redisStatus": {
      "phase": "Running",
      "readyReplicas": 1,
      "totalReplicas": 1
    }
  }
}
```

### Create Tenant

Creates a new tenant.

**Endpoint:** `POST /api/tenants`

**Request Body:**

```json
{
  "metadata": {
    "name": "tenant-1"
  },
  "spec": {
    "displayName": "Tenant 1",
    "description": "Description of Tenant 1",
    "server": {
      "replicas": 1,
      "image": "neurallog/server:latest",
      "resources": {
        "cpu": {
          "request": "100m",
          "limit": "500m"
        },
        "memory": {
          "request": "128Mi",
          "limit": "512Mi"
        }
      },
      "env": [
        {
          "name": "ENV_VAR_1",
          "value": "value1"
        }
      ]
    },
    "redis": {
      "replicas": 1,
      "image": "redis:7-alpine",
      "storage": "1Gi",
      "resources": {
        "cpu": {
          "request": "100m",
          "limit": "500m"
        },
        "memory": {
          "request": "128Mi",
          "limit": "512Mi"
        }
      }
    },
    "networkPolicy": {
      "enabled": true
    }
  }
}
```

**Response:**

```json
{
  "metadata": {
    "name": "tenant-1",
    "creationTimestamp": "2023-01-01T00:00:00Z"
  },
  "spec": {
    "displayName": "Tenant 1",
    "description": "Description of Tenant 1",
    "server": {
      "replicas": 1,
      "image": "neurallog/server:latest",
      "resources": {
        "cpu": {
          "request": "100m",
          "limit": "500m"
        },
        "memory": {
          "request": "128Mi",
          "limit": "512Mi"
        }
      },
      "env": [
        {
          "name": "ENV_VAR_1",
          "value": "value1"
        }
      ]
    },
    "redis": {
      "replicas": 1,
      "image": "redis:7-alpine",
      "storage": "1Gi",
      "resources": {
        "cpu": {
          "request": "100m",
          "limit": "500m"
        },
        "memory": {
          "request": "128Mi",
          "limit": "512Mi"
        }
      }
    },
    "networkPolicy": {
      "enabled": true
    }
  }
}
```

### Update Tenant

Updates an existing tenant.

**Endpoint:** `PUT /api/tenants/{name}`

**Parameters:**
- `name` (path): The name of the tenant.

**Request Body:**

```json
{
  "metadata": {
    "name": "tenant-1"
  },
  "spec": {
    "displayName": "Updated Tenant 1",
    "description": "Updated description of Tenant 1",
    "server": {
      "replicas": 2,
      "image": "neurallog/server:latest",
      "resources": {
        "cpu": {
          "request": "200m",
          "limit": "1000m"
        },
        "memory": {
          "request": "256Mi",
          "limit": "1Gi"
        }
      },
      "env": [
        {
          "name": "ENV_VAR_1",
          "value": "updated-value1"
        }
      ]
    },
    "redis": {
      "replicas": 1,
      "image": "redis:7-alpine",
      "storage": "2Gi",
      "resources": {
        "cpu": {
          "request": "200m",
          "limit": "1000m"
        },
        "memory": {
          "request": "256Mi",
          "limit": "1Gi"
        }
      }
    },
    "networkPolicy": {
      "enabled": false
    }
  }
}
```

**Response:**

```json
{
  "metadata": {
    "name": "tenant-1",
    "creationTimestamp": "2023-01-01T00:00:00Z"
  },
  "spec": {
    "displayName": "Updated Tenant 1",
    "description": "Updated description of Tenant 1",
    "server": {
      "replicas": 2,
      "image": "neurallog/server:latest",
      "resources": {
        "cpu": {
          "request": "200m",
          "limit": "1000m"
        },
        "memory": {
          "request": "256Mi",
          "limit": "1Gi"
        }
      },
      "env": [
        {
          "name": "ENV_VAR_1",
          "value": "updated-value1"
        }
      ]
    },
    "redis": {
      "replicas": 1,
      "image": "redis:7-alpine",
      "storage": "2Gi",
      "resources": {
        "cpu": {
          "request": "200m",
          "limit": "1000m"
        },
        "memory": {
          "request": "256Mi",
          "limit": "1Gi"
        }
      }
    },
    "networkPolicy": {
      "enabled": false
    }
  }
}
```

### Delete Tenant

Deletes a tenant.

**Endpoint:** `DELETE /api/tenants/{name}`

**Parameters:**
- `name` (path): The name of the tenant.

**Response:**

```json
{
  "message": "Tenant deleted successfully"
}
```

## Tenant Resource Schema

### Tenant

The Tenant resource has the following schema:

```typescript
interface Tenant {
  metadata: {
    name: string;
    creationTimestamp?: string;
  };
  spec: {
    displayName: string;
    description: string;
    server?: {
      replicas: number;
      image: string;
      resources?: {
        cpu?: {
          request?: string;
          limit?: string;
        };
        memory?: {
          request?: string;
          limit?: string;
        };
      };
      env?: Array<{
        name: string;
        value: string;
      }>;
    };
    redis?: {
      replicas: number;
      image: string;
      storage: string;
      resources?: {
        cpu?: {
          request?: string;
          limit?: string;
        };
        memory?: {
          request?: string;
          limit?: string;
        };
      };
    };
    networkPolicy?: {
      enabled: boolean;
    };
  };
  status?: {
    phase: string;
    namespace: string;
    serverStatus?: {
      phase: string;
      readyReplicas: number;
      totalReplicas: number;
      message?: string;
    };
    redisStatus?: {
      phase: string;
      readyReplicas: number;
      totalReplicas: number;
      message?: string;
    };
  };
}
```

### Field Descriptions

#### Metadata

- `name`: The name of the tenant. Must be unique and consist of lowercase alphanumeric characters and hyphens.
- `creationTimestamp`: The time when the tenant was created. This is set by the API server.

#### Spec

- `displayName`: A human-readable name for the tenant.
- `description`: A description of the tenant.
- `server`: Configuration for the NeuralLog server.
  - `replicas`: The number of server instances to run.
  - `image`: The Docker image to use for the server.
  - `resources`: Resource requests and limits for the server.
    - `cpu`: CPU resources.
      - `request`: The amount of CPU requested.
      - `limit`: The maximum amount of CPU allowed.
    - `memory`: Memory resources.
      - `request`: The amount of memory requested.
      - `limit`: The maximum amount of memory allowed.
  - `env`: Environment variables for the server.
    - `name`: The name of the environment variable.
    - `value`: The value of the environment variable.
- `redis`: Configuration for Redis.
  - `replicas`: The number of Redis instances to run.
  - `image`: The Docker image to use for Redis.
  - `storage`: The amount of persistent storage to allocate for Redis data.
  - `resources`: Resource requests and limits for Redis.
    - `cpu`: CPU resources.
      - `request`: The amount of CPU requested.
      - `limit`: The maximum amount of CPU allowed.
    - `memory`: Memory resources.
      - `request`: The amount of memory requested.
      - `limit`: The maximum amount of memory allowed.
- `networkPolicy`: Configuration for network policies.
  - `enabled`: Whether to enable network policies.

#### Status

- `phase`: The overall status of the tenant (Provisioning, Running, Degraded, Failed).
- `namespace`: The Kubernetes namespace for the tenant.
- `serverStatus`: Status of the server deployment.
  - `phase`: The status of the server deployment (Provisioning, Running, Degraded, Failed).
  - `readyReplicas`: The number of ready server replicas.
  - `totalReplicas`: The total number of server replicas.
  - `message`: Additional status information or error messages.
- `redisStatus`: Status of the Redis statefulset.
  - `phase`: The status of the Redis statefulset (Provisioning, Running, Degraded, Failed).
  - `readyReplicas`: The number of ready Redis replicas.
  - `totalReplicas`: The total number of Redis replicas.
  - `message`: Additional status information or error messages.

## Error Codes

The API may return the following error codes:

- `400 Bad Request`: The request was invalid or malformed.
- `401 Unauthorized`: Authentication is required.
- `403 Forbidden`: The authenticated user does not have permission to perform the requested operation.
- `404 Not Found`: The requested resource was not found.
- `409 Conflict`: The request conflicts with the current state of the server.
- `500 Internal Server Error`: An error occurred on the server.

## API Usage Examples

### Creating a Tenant

```bash
curl -X POST http://localhost:3000/api/tenants \
  -H "Content-Type: application/json" \
  -d '{
    "metadata": {
      "name": "example-tenant"
    },
    "spec": {
      "displayName": "Example Tenant",
      "description": "An example tenant",
      "server": {
        "replicas": 1,
        "image": "neurallog/server:latest"
      },
      "redis": {
        "replicas": 1,
        "image": "redis:7-alpine",
        "storage": "1Gi"
      },
      "networkPolicy": {
        "enabled": true
      }
    }
  }'
```

### Updating a Tenant

```bash
curl -X PUT http://localhost:3000/api/tenants/example-tenant \
  -H "Content-Type: application/json" \
  -d '{
    "metadata": {
      "name": "example-tenant"
    },
    "spec": {
      "displayName": "Updated Example Tenant",
      "description": "An updated example tenant",
      "server": {
        "replicas": 2,
        "image": "neurallog/server:latest"
      },
      "redis": {
        "replicas": 1,
        "image": "redis:7-alpine",
        "storage": "2Gi"
      },
      "networkPolicy": {
        "enabled": true
      }
    }
  }'
```

### Deleting a Tenant

```bash
curl -X DELETE http://localhost:3000/api/tenants/example-tenant
```
