# Development Environment Setup

This document provides instructions for setting up a development environment for the NeuralLog Tenant Operator.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Setting Up the Development Environment](#setting-up-the-development-environment)
- [Building the Operator](#building-the-operator)
- [Running the Operator Locally](#running-the-operator-locally)
- [Debugging the Operator](#debugging-the-operator)
- [Development Workflow](#development-workflow)
- [IDE Setup](#ide-setup)
- [Troubleshooting](#troubleshooting)

## Prerequisites

Before you begin, ensure you have the following installed:

- [Go](https://golang.org/dl/) (version 1.20 or later)
- [Docker](https://docs.docker.com/get-docker/)
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation) (for local Kubernetes cluster)
- [kubebuilder](https://book.kubebuilder.io/quick-start.html#installation) (for operator development)
- [kustomize](https://kubectl.docs.kubernetes.io/installation/kustomize/) (for Kubernetes manifests)
- [controller-gen](https://book.kubebuilder.io/reference/controller-gen.html) (for generating CRDs and code)

## Setting Up the Development Environment

### Clone the Repository

```bash
git clone https://github.com/neurallog/infra.git
cd infra/operator
```

### Install Dependencies

```bash
go mod download
```

### Generate CRDs and Code

```bash
make generate
make manifests
```

### Create a Kind Cluster

```bash
kind create cluster --name neurallog-dev
```

### Install CRDs

```bash
make install
```

## Building the Operator

### Build the Operator Binary

```bash
make build
```

### Build the Operator Image

```bash
make docker-build IMG=neurallog/tenant-operator:dev
```

### Push the Operator Image

```bash
make docker-push IMG=neurallog/tenant-operator:dev
```

## Running the Operator Locally

### Run the Operator Outside the Cluster

```bash
make run
```

This will run the operator locally, connecting to the Kubernetes cluster specified in your kubeconfig.

### Run the Operator in the Cluster

```bash
make deploy IMG=neurallog/tenant-operator:dev
```

This will deploy the operator to the Kubernetes cluster.

### Verify the Operator

```bash
kubectl get pods -n neurallog-system
```

You should see the operator pod running.

## Debugging the Operator

### Debugging with Delve

You can use [Delve](https://github.com/go-delve/delve) to debug the operator:

```bash
dlv debug main.go -- --zap-log-level=debug
```

### Debugging with VS Code

1. Install the Go extension for VS Code
2. Create a launch configuration in `.vscode/launch.json`:

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Operator",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/main.go",
            "args": ["--zap-log-level=debug"]
        }
    ]
}
```

3. Press F5 to start debugging

### Debugging with GoLand

1. Create a new Run/Debug Configuration
2. Set the configuration type to "Go Build"
3. Set the package to "main"
4. Add the following program arguments: "--zap-log-level=debug"
5. Click "Debug" to start debugging

### Viewing Operator Logs

```bash
# For locally running operator
# The logs will be printed to the console

# For operator running in the cluster
kubectl logs -l control-plane=controller-manager -n neurallog-system
```

## Development Workflow

### Making Changes

1. Make changes to the code
2. Generate CRDs and code:

```bash
make generate
make manifests
```

3. Build and run the operator:

```bash
make build
make run
```

4. Test your changes

### Testing Changes

1. Create a test tenant:

```bash
kubectl apply -f config/samples/neurallog_v1_tenant.yaml
```

2. Verify the tenant is created:

```bash
kubectl get tenants
kubectl describe tenant sample-tenant
```

3. Verify the tenant resources are created:

```bash
kubectl get all -n tenant-sample-tenant
```

### Cleaning Up

1. Delete the test tenant:

```bash
kubectl delete -f config/samples/neurallog_v1_tenant.yaml
```

2. Uninstall the CRDs:

```bash
make uninstall
```

3. Delete the kind cluster:

```bash
kind delete cluster --name neurallog-dev
```

## IDE Setup

### VS Code

1. Install the following extensions:
   - Go
   - Kubernetes
   - YAML
   - Docker

2. Configure the Go extension:
   - Go: Lint Tool = "golangci-lint"
   - Go: Format Tool = "gofmt"
   - Go: Test Flags = "-v"

3. Configure the YAML extension:
   - YAML: Schemas = Add Kubernetes schema

### GoLand

1. Install the following plugins:
   - Kubernetes
   - YAML/Ansible Support
   - Docker

2. Configure Go tools:
   - Settings > Tools > File Watchers > Add "go fmt"
   - Settings > Tools > File Watchers > Add "golangci-lint"

3. Configure Kubernetes:
   - Settings > Languages & Frameworks > Kubernetes > Enable Kubernetes Support

## Troubleshooting

### Common Issues

#### CRD Not Found

If you get an error like "no matches for kind "Tenant" in version "neurallog.io/v1"", ensure you have installed the CRDs:

```bash
make install
```

#### Controller Not Starting

If the controller doesn't start, check the logs:

```bash
kubectl logs -l control-plane=controller-manager -n neurallog-system
```

#### Permission Issues

If you get permission errors, ensure the controller has the necessary RBAC permissions:

```bash
kubectl apply -f config/rbac/role.yaml
kubectl apply -f config/rbac/role_binding.yaml
```

#### Image Pull Issues

If the operator pod is stuck in "ImagePullBackOff", ensure the image is available:

```bash
docker push neurallog/tenant-operator:dev
```

#### Go Module Issues

If you get Go module errors, try:

```bash
go mod tidy
go mod verify
```
