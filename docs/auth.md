# Auth Service

The Auth Service provides authentication and authorization capabilities for the NeuralLog platform. It uses OpenFGA (Fine-Grained Authorization) to manage permissions and supports multi-tenancy.

## Overview

The Auth Service is a critical component of the NeuralLog infrastructure, providing:

- **Authentication**: Verifying user identities
- **Authorization**: Managing permissions and access control
- **Multi-tenancy**: Supporting isolated tenant environments
- **Fine-grained access control**: Controlling access at a granular level

## Architecture

The Auth Service consists of two main components:

1. **Auth Service**: A Node.js application that provides the API for authentication and authorization
2. **OpenFGA**: A fine-grained authorization system that stores and evaluates authorization models and relationships

### Deployment Options

The Auth Service supports two deployment modes:

1. **Local Mode**: For development and self-hosted deployments
   - Uses a single OpenFGA instance
   - Suitable for development and testing

2. **Kubernetes Mode**: For production deployments
   - Can use either:
     - A single global OpenFGA instance (recommended)
     - Tenant-specific OpenFGA instances (advanced use case)
   - Provides better scalability and reliability

## Tenant Isolation

Tenant isolation is implemented through the authorization model:

- Each resource has a parent relationship to a tenant
- All permission checks include tenant context
- Users must be members of a tenant to access its resources
- The authorization model enforces isolation at the data level

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `NODE_ENV` | Environment (development/production) | `development` |
| `PORT` | Port to listen on | `3000` |
| `LOG_LEVEL` | Logging level | `info` |
| `OPENFGA_ADAPTER_TYPE` | Adapter type (local/kubernetes) | `local` in dev, `kubernetes` in prod |
| `OPENFGA_HOST` | OpenFGA host (local mode) | `localhost` |
| `OPENFGA_PORT` | OpenFGA port (local mode) | `8080` |
| `OPENFGA_GLOBAL_API_URL` | Global OpenFGA URL (kubernetes mode) | `http://openfga.openfga-system.svc.cluster.local:8080` |
| `OPENFGA_USE_TENANT_SPECIFIC_INSTANCES` | Whether to use tenant-specific OpenFGA instances | `false` |
| `CACHE_TTL` | Cache TTL in seconds | `300` |
| `CACHE_CHECK_PERIOD` | Cache check period in seconds | `60` |

## API Endpoints

### Authentication

- `POST /api/auth/check`: Check if a user has permission to access a resource
- `POST /api/auth/grant`: Grant a permission to a user
- `POST /api/auth/revoke`: Revoke a permission from a user

### Tenant Management

- `POST /api/tenants`: Create a new tenant
- `GET /api/tenants`: List all tenants
- `DELETE /api/tenants/:tenantId`: Delete a tenant

## Local Development

To run the Auth Service locally:

```bash
# From the infra directory
docker-compose -f docker/auth/docker-compose.yml up -d
```

This will start:
- Auth Service on http://localhost:3040
- OpenFGA on http://localhost:8080
- PostgreSQL on port 5432

## Kubernetes Deployment

To deploy the Auth Service to Kubernetes:

```bash
# From the infra directory
kubectl apply -k kubernetes/base/auth
kubectl apply -k kubernetes/base/openfga
```

## Integration with Other Services

Other services can integrate with the Auth Service by:

1. Making API calls to check permissions
2. Including the tenant ID in the request headers
3. Using the appropriate relations for their resources

Example:

```javascript
// Check if a user has permission to access a resource
const response = await fetch('http://auth:3000/api/auth/check', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'X-Tenant-ID': 'tenant123'
  },
  body: JSON.stringify({
    user: 'user:alice',
    relation: 'reader',
    object: 'log:mylog'
  })
});

const { allowed } = await response.json();
if (allowed) {
  // User has permission
} else {
  // User does not have permission
}
```
