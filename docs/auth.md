# Auth Service

The Auth Service provides centralized authentication and authorization capabilities for the NeuralLog platform. It uses OpenFGA (Fine-Grained Authorization) to manage permissions and supports multi-tenancy with proper isolation.

## Overview

The Auth Service is a critical component of the NeuralLog infrastructure, providing:

- **Authentication**: Verifying user identities through Auth0 integration
- **Authorization**: Managing permissions and role-based access control (RBAC)
- **Multi-tenancy**: Supporting isolated tenant environments with proper security boundaries
- **Fine-grained access control**: Controlling access at resource, organization, and tenant levels
- **API Key Management**: Secure API key generation and verification using Zero-Knowledge Proofs (ZKP)

## Architecture

The Auth Service follows a centralized architecture with global and tenant-specific components:

### Global Components (Shared Across All Tenants)

1. **Auth Service**: A Node.js application that provides the API for authentication and authorization
2. **OpenFGA**: A single global instance for fine-grained authorization across all tenants
3. **Auth0**: A single global identity provider for user authentication
4. **PostgreSQL**: A single global database for OpenFGA and Auth Service data

### Tenant-Specific Components (Isolated Per Tenant)

1. **Web Server**: Dedicated web application instance per tenant
2. **Logs Server**: Dedicated logs server instance per tenant
3. **Redis**: One Redis instance per tenant, shared between auth and logs services

### Deployment Options

The Auth Service supports two deployment modes:

1. **Local Mode**: For development and self-hosted deployments
   - Uses a single OpenFGA instance with proper tenant namespacing
   - Suitable for development and testing

2. **Kubernetes Mode**: For production deployments
   - Uses a single global OpenFGA instance with proper tenant namespacing
   - Each tenant gets dedicated web, logs, and Redis instances in isolated namespaces
   - Provides better scalability, reliability, and security isolation

## Tenant Isolation

NeuralLog implements a hybrid isolation model combining infrastructure and logical isolation:

### Infrastructure Isolation

- **Dedicated Instances**: Each tenant gets dedicated web, logs, and Redis instances
- **Namespace Isolation**: In Kubernetes, each tenant's components run in isolated namespaces
- **Network Isolation**: Network policies restrict communication between tenant namespaces
- **Resource Isolation**: Resource quotas and limits prevent resource contention

### Logical Isolation

- **Authorization Model**: OpenFGA enforces tenant boundaries through its authorization model
- **Parent Relationships**: Each resource has a parent relationship to a tenant
- **Tenant Context**: All permission checks include tenant context
- **Membership Requirements**: Users must be members of a tenant to access its resources
- **Data Namespacing**: Even in shared components, data is properly namespaced by tenant ID

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

- `POST /api/auth/login`: Authenticate a user with username and password
- `POST /api/auth/m2m`: Authenticate a machine-to-machine client
- `POST /api/auth/validate`: Validate a token and get user information
- `POST /api/auth/exchange-token`: Exchange an Auth0 token for a server access token

### Authorization

- `POST /api/auth/check`: Check if a user has permission to access a resource
- `POST /api/auth/grant`: Grant a permission to a user
- `POST /api/auth/revoke`: Revoke a permission from a user

### Role Management

- `POST /api/roles`: Create a new role
- `GET /api/roles`: List all roles
- `GET /api/roles/:roleId`: Get a specific role
- `PUT /api/roles/:roleId`: Update a role
- `DELETE /api/roles/:roleId`: Delete a role
- `POST /api/roles/:roleId/assign`: Assign a role to a user
- `POST /api/roles/:roleId/revoke`: Revoke a role from a user

### User Management

- `GET /api/users`: List all users
- `GET /api/users/:userId`: Get a specific user
- `DELETE /api/users/:userId`: Delete a user
- `PUT /api/users/:userId/roles`: Update a user's roles

### API Key Management

- `POST /api/apikeys`: Create a new API key
- `GET /api/apikeys`: List all API keys for the current user
- `DELETE /api/apikeys/:keyId`: Revoke an API key
- `POST /api/apikeys/verify`: Verify an API key

### Tenant Management

- `POST /api/tenants`: Create a new tenant
- `GET /api/tenants`: List all tenants
- `GET /api/tenants/:tenantId`: Get a specific tenant
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

Other services can integrate with the Auth Service using these patterns:

### Authentication Integration

1. **Token Exchange**: Exchange Auth0 tokens for server access tokens
2. **API Key Verification**: Verify API keys for machine-to-machine authentication
3. **Token Validation**: Validate tokens to get user information

### Authorization Integration

1. **Permission Checking**: Check if a user has permission to access a resource
2. **Role-Based Access Control**: Use roles to control access to resources
3. **Resource-Level Permissions**: Control access at the resource level

### Implementation Guidelines

1. **Include Tenant Context**: Always include the tenant ID in request headers
2. **Use Bearer Tokens**: Include tokens in the Authorization header
3. **Implement Middleware**: Use middleware to check permissions for all API endpoints
4. **Cache Results**: Cache permission checks for performance
5. **Handle Errors**: Properly handle authorization errors

### Example: Checking Permissions

```javascript
// Check if a user has permission to access a resource
async function checkPermission(userId, permission, resourceType, resourceId) {
  const response = await fetch('http://auth:3040/api/auth/check', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${serverToken}`,
      'X-Tenant-ID': tenantId
    },
    body: JSON.stringify({
      user: `user:${userId}`,
      permission,
      resourceType,
      resourceId
    })
  });

  const { allowed } = await response.json();
  return allowed;
}
```

### Example: Role-Based Access Control

```javascript
// Assign a role to a user
async function assignRole(userId, roleId, organizationId) {
  const response = await fetch(`http://auth:3040/api/roles/${roleId}/assign`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${serverToken}`,
      'X-Tenant-ID': tenantId
    },
    body: JSON.stringify({
      userId,
      organizationId
    })
  });

  return response.ok;
}
```

### Example: API Key Management

```javascript
// Create a new API key
async function createApiKey(name, scopes) {
  const response = await fetch('http://auth:3040/api/apikeys', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${serverToken}`,
      'X-Tenant-ID': tenantId
    },
    body: JSON.stringify({
      name,
      scopes
    })
  });

  return await response.json();
}
```
