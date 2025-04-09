# NeuralLog Admin User Manual

This user manual provides detailed instructions for using the NeuralLog Admin interface to manage tenants in your Kubernetes cluster.

## Table of Contents

1. [Dashboard Overview](#dashboard-overview)
2. [Managing Tenants](#managing-tenants)
   - [Creating a Tenant](#creating-a-tenant)
   - [Viewing Tenant Details](#viewing-tenant-details)
   - [Editing a Tenant](#editing-a-tenant)
   - [Deleting a Tenant](#deleting-a-tenant)
3. [Tenant Configuration](#tenant-configuration)
   - [Server Configuration](#server-configuration)
   - [Redis Configuration](#redis-configuration)
   - [Network Policies](#network-policies)
4. [Monitoring Tenant Status](#monitoring-tenant-status)
5. [Troubleshooting](#troubleshooting)

## Dashboard Overview

The NeuralLog Admin dashboard provides a centralized view of all tenants in your Kubernetes cluster.

![Dashboard Overview](./images/dashboard.png)

The dashboard includes:

1. **Navigation Bar**: Access to the main sections of the application.
2. **Tenant List**: A table showing all tenants with their key information.
3. **Status Indicators**: Visual indicators of tenant health and status.
4. **Action Buttons**: Buttons for creating, editing, and deleting tenants.

## Managing Tenants

### Creating a Tenant

To create a new tenant:

1. Click the **New Tenant** button in the top-right corner of the dashboard.

   ![New Tenant Button](./images/new-tenant-button.png)

2. Fill in the tenant details in the form:

   ![Create Tenant Form](./images/create-tenant-form.png)

   - **Tenant Name**: A unique identifier for the tenant (lowercase alphanumeric characters, hyphens allowed).
   - **Display Name**: A human-readable name for the tenant.
   - **Description**: A brief description of the tenant's purpose.

3. Configure the server settings:
   - **Replicas**: The number of server instances to run.
   - **Image**: The Docker image to use for the server.
   - **Resources**: CPU and memory requests/limits.

4. Configure the Redis settings:
   - **Replicas**: The number of Redis instances to run.
   - **Image**: The Docker image to use for Redis.
   - **Storage**: The amount of persistent storage to allocate.
   - **Resources**: CPU and memory requests/limits.

5. Configure network policies:
   - **Enable Network Policies**: Toggle to enable or disable network isolation.

6. Click the **Create Tenant** button to create the tenant.

### Viewing Tenant Details

To view details for a specific tenant:

1. On the dashboard, click on the tenant's name in the tenant list.

   ![Tenant List](./images/tenant-list.png)

2. The tenant details page shows comprehensive information about the tenant:

   ![Tenant Details](./images/tenant-details.png)

   - **Overview**: Basic tenant information and status.
   - **Server Status**: Status of the server deployment.
   - **Redis Status**: Status of the Redis statefulset.
   - **Configuration**: Current configuration settings.

### Editing a Tenant

To edit an existing tenant:

1. On the dashboard, click the **Edit** button for the tenant you want to modify.

   ![Edit Button](./images/edit-button.png)

2. Update the tenant's configuration in the form:

   ![Edit Tenant Form](./images/edit-tenant-form.png)

3. Click the **Update Tenant** button to save your changes.

### Deleting a Tenant

To delete a tenant:

1. On the dashboard, click the **Delete** button for the tenant you want to remove.

   ![Delete Button](./images/delete-button.png)

2. Confirm the deletion in the confirmation dialog:

   ![Delete Confirmation](./images/delete-confirmation.png)

3. Click the **Delete** button to permanently remove the tenant.

## Tenant Configuration

### Server Configuration

The server configuration controls the NeuralLog server deployment for the tenant:

![Server Configuration](./images/server-configuration.png)

- **Replicas**: The number of server instances to run. Increasing this number improves availability and scalability.
- **Image**: The Docker image to use for the server. Format: `repository/image:tag`.
- **Resources**:
  - **CPU Request/Limit**: The amount of CPU resources to request/limit for each server instance.
  - **Memory Request/Limit**: The amount of memory to request/limit for each server instance.
- **Environment Variables**: Custom environment variables for the server.

### Redis Configuration

The Redis configuration controls the Redis statefulset for the tenant:

![Redis Configuration](./images/redis-configuration.png)

- **Replicas**: The number of Redis instances to run. For high availability, use at least 2 replicas.
- **Image**: The Docker image to use for Redis. Format: `repository/image:tag`.
- **Storage**: The amount of persistent storage to allocate for Redis data.
- **Resources**:
  - **CPU Request/Limit**: The amount of CPU resources to request/limit for each Redis instance.
  - **Memory Request/Limit**: The amount of memory to request/limit for each Redis instance.

### Network Policies

Network policies control the network traffic to and from the tenant's resources:

![Network Policies](./images/network-policies.png)

- **Enable Network Policies**: When enabled, network traffic is restricted according to the defined policies.
- **Default Policy**: By default, all ingress traffic is denied except from within the tenant's namespace.

## Monitoring Tenant Status

The tenant list and details pages provide status information for each tenant:

![Tenant Status](./images/tenant-status.png)

- **Phase**: The overall status of the tenant (Provisioning, Running, Degraded, Failed).
- **Server Status**: The status of the server deployment (Running, Degraded, Failed).
- **Redis Status**: The status of the Redis statefulset (Running, Degraded, Failed).
- **Ready Replicas**: The number of ready replicas for the server and Redis.
- **Total Replicas**: The total number of replicas for the server and Redis.
- **Message**: Additional status information or error messages.

## Troubleshooting

### Common Issues

#### Tenant Creation Fails

If tenant creation fails, check the following:

1. Ensure the tenant name is unique and follows the naming convention.
2. Verify that the Kubernetes API server is accessible.
3. Check that the Tenant Operator is running correctly.

#### Tenant Status Shows "Degraded"

If a tenant's status shows "Degraded":

1. Check the server and Redis status for more specific information.
2. Verify that the requested resources are available in the cluster.
3. Check the Kubernetes events for the tenant's namespace.

#### Cannot Access Tenant Services

If you cannot access a tenant's services:

1. Verify that the tenant is in the "Running" state.
2. Check the network policies if enabled.
3. Ensure that the service endpoints are correctly configured.

### Getting Help

If you encounter issues that you cannot resolve:

1. Check the [Troubleshooting Guide](./troubleshooting.md) for more detailed information.
2. Review the Kubernetes logs for the Tenant Operator and tenant resources.
3. Contact the NeuralLog support team for assistance.
