# NeuralLog Admin Architecture

## Overview

The NeuralLog Admin interface is a web application that allows administrators to manage NeuralLog tenants in a Kubernetes environment. This document provides a high-level overview of the system architecture.

## System Components

![NeuralLog Admin Architecture](./images/architecture-diagram.png)

### Web Frontend

The web frontend is a Next.js application that provides a user interface for managing NeuralLog tenants. It communicates with the Kubernetes API server to perform operations on tenant resources.

**Key Technologies:**
- Next.js (React framework)
- Chakra UI (Component library)
- SWR (Data fetching)

### Kubernetes Backend

The Kubernetes backend consists of the following components:

1. **Kubernetes API Server**: The central management point for the Kubernetes cluster.
2. **Tenant Operator**: A custom Kubernetes operator that manages tenant resources.
3. **Tenant Custom Resource Definition (CRD)**: A custom resource that defines a NeuralLog tenant.

### Tenant Resources

Each tenant consists of the following resources:

1. **Namespace**: A dedicated Kubernetes namespace for the tenant.
2. **Server Deployment**: A deployment of the NeuralLog server.
3. **Redis StatefulSet**: A StatefulSet for the tenant's Redis instance.
4. **Network Policies**: Optional network policies to restrict communication.

## Data Flow

1. **User Interaction**: The user interacts with the web frontend to create, view, update, or delete tenants.
2. **API Requests**: The frontend makes API requests to the Next.js API routes.
3. **Kubernetes API**: The Next.js API routes communicate with the Kubernetes API server.
4. **Tenant Operator**: The Tenant Operator watches for changes to tenant resources and reconciles the desired state.
5. **Resource Management**: The Tenant Operator creates, updates, or deletes the necessary Kubernetes resources for each tenant.

## Security Model

The NeuralLog Admin interface uses the following security measures:

1. **Kubernetes RBAC**: Role-Based Access Control for API access.
2. **Tenant Isolation**: Each tenant has its own namespace and resources.
3. **Network Policies**: Optional network policies to restrict communication between tenants.

## Scalability

The system is designed to scale in the following ways:

1. **Horizontal Scaling**: The web frontend can be scaled horizontally.
2. **Tenant Scaling**: Each tenant's server and Redis instances can be scaled independently.
3. **Multi-Cluster Support**: The system can manage tenants across multiple Kubernetes clusters.
