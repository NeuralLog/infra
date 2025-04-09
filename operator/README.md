# NeuralLog Tenant Operator

[![Build Status](https://github.com/neurallog/infra/actions/workflows/operator-build.yml/badge.svg)](https://github.com/neurallog/infra/actions/workflows/operator-build.yml)
[![Docker Image](https://github.com/neurallog/infra/actions/workflows/operator-docker.yml/badge.svg)](https://github.com/neurallog/infra/actions/workflows/operator-docker.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/neurallog/infra/operator)](https://goreportcard.com/report/github.com/neurallog/infra/operator)
[![License](https://img.shields.io/github/license/neurallog/infra)](https://github.com/neurallog/infra/blob/main/LICENSE)

The NeuralLog Tenant Operator is a Kubernetes operator that manages multi-tenant deployments of NeuralLog. It provides a simple API for creating and managing tenants, each with their own isolated resources.

## Overview

The operator introduces a new custom resource definition (CRD) called `Tenant`. Each Tenant resource represents a complete NeuralLog deployment with its own:

- Dedicated namespace
- Server deployment
- Redis instance
- Network policies
- Auth service integration

## Architecture

The operator follows the Kubernetes operator pattern:

1. **Custom Resource Definition (CRD)**: Defines the `Tenant` resource
2. **Controller**: Watches for changes to Tenant resources and reconciles the desired state
3. **Reconciliation Logic**: Creates, updates, or deletes Kubernetes resources based on the Tenant specification

## Tenant Custom Resource

The Tenant custom resource allows you to define:

- Resource limits and requests
- Server configuration
- Redis configuration
- Network policies

Example:

```yaml
apiVersion: neurallog.io/v1
kind: Tenant
metadata:
  name: sample-tenant
spec:
  displayName: "Sample Tenant"
  description: "A sample tenant for demonstration purposes"

  # Server configuration
  server:
    replicas: 2
    image: "neurallog/server:latest"

  # Redis configuration
  redis:
    storage: "5Gi"
```

## Installation

### Prerequisites

- Kubernetes cluster (v1.19+)
- kubectl configured to communicate with your cluster
- cert-manager (v1.0.0+)

### Installing the Operator

1. Install the CRDs:

```bash
kubectl apply -f config/crd/bases
```

2. Install the operator:

```bash
kubectl apply -f config/manager
```

## Usage

### Creating a Tenant

1. Create a Tenant resource:

```bash
kubectl apply -f config/samples/neurallog_v1_tenant.yaml
```

2. Check the status of the tenant:

```bash
kubectl get tenants
```

3. Access the tenant's server:

```bash
# Get the namespace
NAMESPACE=$(kubectl get tenant sample-tenant -o jsonpath='{.status.namespace}')

# Port forward to the server
kubectl -n $NAMESPACE port-forward svc/neurallog-server 3030:3030
```

### Updating a Tenant

Edit the Tenant resource and apply the changes:

```bash
kubectl apply -f config/samples/neurallog_v1_tenant.yaml
```

The operator will reconcile the changes and update the resources accordingly.

### Deleting a Tenant

Delete the Tenant resource:

```bash
kubectl delete tenant sample-tenant
```

The operator will clean up all resources associated with the tenant.

## Development

### Prerequisites

- Go (1.20+)
- Kubebuilder (3.0.0+)
- Docker

### Building the Operator

1. Build the operator image:

```bash
docker build -t neurallog/tenant-operator:latest .
```

2. Push the image to a registry:

```bash
docker push neurallog/tenant-operator:latest
```

### Running Locally

1. Install the CRDs:

```bash
make install
```

2. Run the operator locally:

```bash
make run
```

## License

Copyright 2023 NeuralLog Authors.

Licensed under the Apache License, Version 2.0.
