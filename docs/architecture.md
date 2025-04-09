# NeuralLog Architecture Guide

This document provides a detailed overview of the NeuralLog infrastructure architecture, including components, interactions, and design decisions.

## Table of Contents

- [System Overview](#system-overview)
- [Multi-Tenant Architecture](#multi-tenant-architecture)
- [Component Architecture](#component-architecture)
  - [Server Component](#server-component)
  - [Redis Component](#redis-component)
  - [Tenant Operator](#tenant-operator)
- [Network Architecture](#network-architecture)
- [Storage Architecture](#storage-architecture)
- [Security Architecture](#security-architecture)
- [Scalability Architecture](#scalability-architecture)
- [High Availability Architecture](#high-availability-architecture)

## System Overview

NeuralLog is a system for intelligent log processing with automated actions. The infrastructure is designed to support a multi-tenant deployment model where each tenant has completely isolated resources.

### Key Design Principles

1. **Complete Tenant Isolation**: Each tenant has dedicated resources with no shared components
2. **Kubernetes-Native**: All components are designed to run on Kubernetes
3. **Scalability**: Components can scale independently to handle varying loads
4. **Observability**: Built-in monitoring and logging
5. **Security**: Network policies and resource isolation for security
6. **Automation**: Automated provisioning and management through the Tenant Operator

### High-Level Architecture Diagram

```
                                  ┌─────────────────────┐
                                  │   Tenant Operator   │
                                  └─────────────────────┘
                                            │
                                            ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                             Kubernetes Cluster                           │
│                                                                         │
│  ┌─────────────────────────┐      ┌─────────────────────────┐          │
│  │      Tenant 1           │      │      Tenant 2           │          │
│  │  ┌─────────┐ ┌─────────┐│      │  ┌─────────┐ ┌─────────┐│          │
│  │  │ Server  │ │  Redis  ││      │  │ Server  │ │  Redis  ││          │
│  │  └─────────┘ └─────────┘│      │  └─────────┘ └─────────┘│          │
│  └─────────────────────────┘      └─────────────────────────┘          │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

## Multi-Tenant Architecture

NeuralLog uses a multi-tenant architecture with complete isolation between tenants.

### Tenant Isolation

Each tenant has:

1. **Dedicated Namespace**: A separate Kubernetes namespace for all tenant resources
2. **Dedicated Server**: A dedicated NeuralLog server deployment
3. **Dedicated Redis**: A dedicated Redis instance for data storage
4. **Network Isolation**: Network policies that restrict communication between tenants
5. **Resource Quotas**: Dedicated resource quotas to prevent resource starvation

### Tenant Provisioning

Tenants are provisioned through the Tenant Operator, which:

1. Creates a dedicated namespace for the tenant
2. Deploys Redis StatefulSet in the namespace
3. Deploys NeuralLog server in the namespace
4. Configures network policies for isolation
5. Sets up resource quotas based on tenant specification

## Component Architecture

### Server Component

The NeuralLog server is the core component responsible for log processing and actions.

#### Server Architecture Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                     NeuralLog Server                        │
│                                                             │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐  │
│  │  Log Ingestion  │  │  Rule Engine    │  │  Actions    │  │
│  └─────────────────┘  └─────────────────┘  └─────────────┘  │
│                                                             │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐  │
│  │  API Endpoints  │  │  MCP Server     │  │  Storage    │  │
│  └─────────────────┘  └─────────────────┘  └─────────────┘  │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

#### Server Components

1. **Log Ingestion**: Collects and normalizes logs from various sources
2. **Rule Engine**: Processes logs based on defined rules
3. **Actions**: Executes actions based on rule matches
4. **API Endpoints**: Provides RESTful APIs for management
5. **MCP Server**: Hosts the Model Control Protocol server for AI model interaction
6. **Storage**: Interfaces with Redis for data persistence

#### Server Deployment

The server is deployed as a Kubernetes Deployment with:

- Configurable replicas for horizontal scaling
- Resource limits and requests
- Health checks for liveness and readiness
- Environment variables for configuration
- Volume mounts for temporary storage

### Redis Component

Redis serves as the primary data store for NeuralLog.

#### Redis Architecture

Redis is deployed as a StatefulSet with:

- Persistent storage for data durability
- Custom configuration for performance optimization
- Health checks for liveness and readiness
- Resource limits and requests

#### Redis Data Model

Redis stores:

1. **Recent Logs**: High-performance storage for recent log entries
2. **Pattern Indexes**: Efficient indexing for fast pattern matching and queries
3. **Rule Storage**: Storage for rule definitions and execution history
4. **Real-time Features**: Supporting pub/sub for real-time log streaming and notifications

### Tenant Operator

The Tenant Operator is a Kubernetes operator that manages the lifecycle of NeuralLog tenants.

#### Operator Architecture

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

#### Operator Components

1. **Tenant CRD**: Custom Resource Definition for tenants
2. **Controller**: Watches for changes to Tenant resources
3. **Reconciler**: Reconciles the desired state with the actual state
4. **Redis Manager**: Manages Redis resources for tenants
5. **Server Manager**: Manages server resources for tenants
6. **Network Manager**: Manages network policies for tenants

## Network Architecture

### Network Flow

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  External   │     │  Ingress    │     │  Server     │
│  Client     │────▶│  Controller │────▶│  Service    │
└─────────────┘     └─────────────┘     └─────────────┘
                                               │
                                               ▼
                                        ┌─────────────┐
                                        │  Server     │
                                        │  Pod        │
                                        └─────────────┘
                                               │
                                               ▼
                                        ┌─────────────┐
                                        │  Redis      │
                                        │  Service    │
                                        └─────────────┘
                                               │
                                               ▼
                                        ┌─────────────┐
                                        │  Redis      │
                                        │  Pod        │
                                        └─────────────┘
```

### Network Policies

Network policies are used to restrict communication between components:

1. **Default Deny**: By default, all ingress traffic is denied
2. **Allow Internal**: Traffic within a tenant's namespace is allowed
3. **Allow API**: Traffic to the server's API endpoint is allowed from specified sources
4. **Custom Rules**: Additional rules can be defined in the tenant specification

## Storage Architecture

### Persistent Storage

Redis uses persistent storage for data durability:

1. **StatefulSet**: Redis is deployed as a StatefulSet with persistent volume claims
2. **Storage Class**: The storage class can be configured based on the environment
3. **Storage Size**: The storage size can be configured in the tenant specification

### Backup and Restore

Redis backup and restore scripts are provided for data protection:

1. **Backup Script**: Creates Redis backups
2. **Restore Script**: Restores from backups

## Security Architecture

### Authentication and Authorization

1. **API Authentication**: The server API uses authentication for access control
2. **Kubernetes RBAC**: The operator uses Kubernetes RBAC for authorization

### Network Security

1. **Network Policies**: Network policies restrict communication between components
2. **TLS**: TLS can be configured for secure communication

### Container Security

1. **Non-Root User**: Containers run as non-root users
2. **Read-Only Filesystem**: Containers use read-only filesystems where possible
3. **Resource Limits**: Resource limits prevent resource exhaustion

## Scalability Architecture

### Horizontal Scaling

1. **Server Replicas**: The server can be scaled horizontally by increasing replicas
2. **Redis Replicas**: Redis can be scaled for read performance (future)

### Vertical Scaling

1. **Resource Limits**: Resource limits can be adjusted for vertical scaling
2. **Node Selection**: Node selection can be used for placement on appropriate nodes

## High Availability Architecture

### Component Redundancy

1. **Server Redundancy**: Multiple server replicas for high availability
2. **Redis Redundancy**: Redis can be configured for high availability (future)

### Failure Recovery

1. **Health Checks**: Health checks detect component failures
2. **Automatic Restart**: Failed components are automatically restarted
3. **Data Persistence**: Redis persistence ensures data durability
