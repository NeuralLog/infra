# NeuralLog Tenant Operator Architecture

This document provides a detailed overview of the NeuralLog Tenant Operator architecture.

## Table of Contents

- [Operator Pattern](#operator-pattern)
- [Operator Architecture](#operator-architecture)
- [Controller-Runtime Framework](#controller-runtime-framework)
- [Reconciliation Loop](#reconciliation-loop)
- [Custom Resource Definitions](#custom-resource-definitions)
- [Component Interactions](#component-interactions)
- [Resource Management](#resource-management)

## Operator Pattern

The Kubernetes Operator pattern is a design pattern that extends the Kubernetes API to create, configure, and manage complex applications. Operators use Custom Resource Definitions (CRDs) to define application-specific resources and controllers to manage these resources.

The NeuralLog Tenant Operator follows this pattern to manage NeuralLog tenants in a Kubernetes cluster. It provides a declarative API for creating and managing tenants, each with their own isolated resources.

## Operator Architecture

The NeuralLog Tenant Operator consists of the following components:

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

1. **Tenant CRD**: Custom Resource Definition for tenants
2. **Controller**: Watches for changes to Tenant resources
3. **Reconciler**: Reconciles the desired state with the actual state
4. **Redis Manager**: Manages Redis resources for tenants
5. **Server Manager**: Manages server resources for tenants
6. **Network Manager**: Manages network policies for tenants

## Controller-Runtime Framework

The NeuralLog Tenant Operator is built using the [controller-runtime](https://github.com/kubernetes-sigs/controller-runtime) framework, which provides high-level abstractions for building Kubernetes controllers.

Key components of the controller-runtime framework used in the operator:

1. **Manager**: The central component that manages controllers, webhooks, and shared caches
2. **Controller**: Watches for changes to resources and triggers reconciliation
3. **Reconciler**: Implements the reconciliation logic
4. **Client**: Provides access to the Kubernetes API
5. **Scheme**: Registers API types with the client
6. **Cache**: Caches resources for efficient access

## Reconciliation Loop

The reconciliation loop is the core of the operator. It is responsible for ensuring that the actual state of the system matches the desired state specified in the Tenant resource.

The reconciliation loop follows these steps:

1. **Fetch**: Fetch the Tenant resource
2. **Validate**: Validate the Tenant resource
3. **Initialize**: Initialize the Tenant status if needed
4. **Reconcile**: Reconcile the Tenant resources
   - Namespace
   - Redis
   - Server
   - Network Policies
5. **Update**: Update the Tenant status
6. **Requeue**: Requeue for periodic reconciliation

```
┌─────────────┐
│   Fetch     │
└─────┬───────┘
      │
      ▼
┌─────────────┐
│  Validate   │
└─────┬───────┘
      │
      ▼
┌─────────────┐
│ Initialize  │
└─────┬───────┘
      │
      ▼
┌─────────────┐
│  Reconcile  │
└─────┬───────┘
      │
      ▼
┌─────────────┐
│   Update    │
└─────┬───────┘
      │
      ▼
┌─────────────┐
│   Requeue   │
└─────────────┘
```

## Custom Resource Definitions

The NeuralLog Tenant Operator defines the following Custom Resource Definitions (CRDs):

1. **Tenant**: Represents a NeuralLog tenant

The Tenant CRD includes:

- **Spec**: Desired state of the tenant
  - Basic information (displayName, description)
  - Server configuration
  - Redis configuration
  - Network policy configuration
  - Monitoring configuration
  - Backup configuration
  - Lifecycle hooks
  - Integrations
- **Status**: Observed state of the tenant
  - Phase
  - Conditions
  - Component status
  - URLs
  - Metrics
  - Backup status

## Component Interactions

The NeuralLog Tenant Operator interacts with various Kubernetes components:

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│  Tenant CRD     │     │  Kubernetes API │     │  Tenant         │
│                 │────▶│                 │◀────│  Controller     │
└─────────────────┘     └─────────────────┘     └─────────────────┘
                                │                        │
                                ▼                        │
┌─────────────────┐     ┌─────────────────┐             │
│  Namespace      │     │  Kubernetes     │             │
│                 │◀────│  Resources      │◀────────────┘
└─────────────────┘     └─────────────────┘
       │                        │
       │                        │
       ▼                        ▼
┌─────────────────┐     ┌─────────────────┐
│  Redis          │     │  Server         │
│  Resources      │     │  Resources      │
└─────────────────┘     └─────────────────┘
```

1. The Tenant Controller watches for changes to Tenant resources
2. When a change is detected, the controller reconciles the Tenant
3. The controller creates or updates Kubernetes resources:
   - Namespace
   - Redis resources (StatefulSet, Service, ConfigMap)
   - Server resources (Deployment, Service)
   - Network policies
4. The controller updates the Tenant status with the current state

## Resource Management

The NeuralLog Tenant Operator manages the following resources for each tenant:

1. **Namespace**: A dedicated namespace for the tenant
2. **Redis Resources**:
   - StatefulSet: Manages Redis pods
   - Service: Exposes Redis
   - ConfigMap: Configures Redis
   - PersistentVolumeClaim: Provides storage for Redis
3. **Server Resources**:
   - Deployment: Manages server pods
   - Service: Exposes the server
   - ConfigMap: Configures the server
4. **Network Policies**:
   - Default deny policy
   - Allow internal traffic policy
   - Allow API access policy
   - Custom policies

The operator uses owner references to ensure that all resources are properly cleaned up when a tenant is deleted.
