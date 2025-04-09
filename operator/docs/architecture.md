# NeuralLog Tenant Operator Architecture

This document describes the architecture of the NeuralLog Tenant Operator, including its components, design principles, and interactions with other systems.

## Overview

The NeuralLog Tenant Operator follows the Kubernetes Operator pattern, which extends the Kubernetes API to create, configure, and manage instances of complex applications on behalf of users. The operator is responsible for managing the lifecycle of tenant resources in the NeuralLog platform.

## Components

The Tenant Operator consists of the following components:

### 1. Custom Resource Definition (CRD)

The `Tenant` CRD defines the schema for tenant resources. It includes specifications for:

- Tenant metadata (name, display name, description)
- Server configuration (replicas, image, resources, environment variables)
- Redis configuration (replicas, image, resources, storage, custom configuration)
- Network policy configuration (enabled/disabled, allowed namespaces, custom rules)

### 2. Controller

The controller is the core component of the operator. It watches for changes to Tenant resources and reconciles the desired state with the actual state of the cluster. The controller is implemented using the controller-runtime library and follows the reconciliation pattern.

### 3. Reconcilers

The operator uses several reconcilers to manage different aspects of tenant resources:

#### Namespace Reconciler

- Creates and manages a dedicated namespace for each tenant
- Sets up resource quotas and limits for the namespace
- Adds labels and annotations for tenant identification

#### Redis Reconciler

- Creates and manages Redis resources for each tenant
- Configures Redis StatefulSet, Service, and ConfigMap
- Handles Redis configuration updates
- Monitors Redis health and status

#### Server Reconciler

- Creates and manages Server resources for each tenant
- Configures Server Deployment and Service
- Handles Server configuration updates
- Monitors Server health and status

#### Network Policy Reconciler

- Creates and manages Network Policies for each tenant
- Configures default deny-all policy for tenant isolation
- Allows internal communication within the tenant namespace
- Supports custom ingress and egress rules

#### Auth Service Reconciler

- Integrates with the NeuralLog Auth service
- Creates tenant entries in the Auth service
- Manages tenant authentication and authorization
- Handles tenant deletion in the Auth service

## Workflow

The typical workflow for the Tenant Operator is as follows:

1. A user creates a Tenant resource using kubectl or the NeuralLog Admin UI
2. The controller detects the new Tenant resource and starts the reconciliation process
3. The controller creates a dedicated namespace for the tenant
4. The controller creates Redis resources in the tenant namespace
5. The controller creates Server resources in the tenant namespace
6. The controller creates Network Policies for tenant isolation
7. The controller integrates with the Auth service to set up tenant authentication
8. The controller updates the Tenant status with the current state of the resources
9. The controller continues to monitor the Tenant resource for changes

## Design Principles

The Tenant Operator follows these design principles:

### 1. Tenant Isolation

Each tenant has its own dedicated namespace and resources, ensuring complete isolation between tenants. Network policies are used to enforce isolation at the network level.

### 2. Declarative Configuration

The operator uses a declarative approach to configuration. Users specify the desired state of the tenant resources, and the operator ensures that the actual state matches the desired state.

### 3. Idempotency

The reconciliation process is idempotent, meaning that it can be run multiple times without causing unintended side effects. This ensures that the operator can recover from failures and continue to maintain the desired state.

### 4. Graceful Deletion

When a tenant is deleted, the operator ensures that all associated resources are properly cleaned up, including resources in the tenant namespace and entries in the Auth service.

### 5. Status Reporting

The operator provides detailed status information about tenant resources, including the current state of Redis and Server components, making it easy to monitor and troubleshoot tenant deployments.

## Integration with Other Systems

The Tenant Operator integrates with the following systems:

### 1. Auth Service

The operator integrates with the NeuralLog Auth service to manage tenant authentication and authorization. When a tenant is created, the operator creates a corresponding entry in the Auth service. When a tenant is deleted, the operator removes the entry from the Auth service.

### 2. Monitoring System

The operator exposes metrics for monitoring tenant resources, including resource usage, health status, and reconciliation metrics. These metrics can be collected by Prometheus and visualized in Grafana dashboards.

### 3. Admin UI

The operator provides a REST API that can be used by the NeuralLog Admin UI to manage tenant resources. The Admin UI allows users to create, update, and delete tenants, as well as view tenant status and metrics.

## Conclusion

The NeuralLog Tenant Operator provides a robust and flexible solution for managing multi-tenant deployments of NeuralLog. By following the Kubernetes Operator pattern and adhering to best practices for tenant isolation and resource management, the operator ensures that tenant resources are properly provisioned, configured, and maintained throughout their lifecycle.
