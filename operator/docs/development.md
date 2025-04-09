# NeuralLog Tenant Operator Development Guide

This guide provides instructions for setting up a development environment for the NeuralLog Tenant Operator and contributing to the project.

## Prerequisites

Before you begin, make sure you have the following tools installed:

- [Go](https://golang.org/doc/install) (v1.21 or later)
- [Docker](https://docs.docker.com/get-docker/)
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [kind](https://kind.sigs.k8s.io/docs/user/quick-start/) or [minikube](https://minikube.sigs.k8s.io/docs/start/) for local Kubernetes development
- [kubebuilder](https://book.kubebuilder.io/quick-start.html) (v3.0.0 or later)
- [kustomize](https://kubectl.docs.kubernetes.io/installation/kustomize/) (v4.0.0 or later)

## Setting Up the Development Environment

### 1. Clone the Repository

```bash
git clone https://github.com/neurallog/infra.git
cd infra/operator
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Set Up a Local Kubernetes Cluster

Using kind:

```bash
kind create cluster --name neurallog-dev
```

Or using minikube:

```bash
minikube start --driver=docker
```

### 4. Install the CRDs

```bash
make install
```

## Development Workflow

### 1. Making Changes

The operator code is organized as follows:

- `api/v1/`: Contains the API definitions for the Tenant CRD
- `controllers/`: Contains the controller logic for reconciling Tenant resources
- `config/`: Contains Kubernetes manifests for deploying the operator
- `docs/`: Contains documentation

When making changes to the API, you need to regenerate the code:

```bash
make generate
```

### 2. Building the Operator

To build the operator locally:

```bash
make build
```

### 3. Running the Operator Locally

You can run the operator locally for development purposes:

```bash
make run
```

This will run the operator outside of the Kubernetes cluster, but it will still connect to the cluster using your kubeconfig.

### 4. Building and Pushing the Docker Image

To build and push the Docker image:

```bash
make docker-build docker-push IMG=<your-registry>/neurallog-operator:tag
```

### 5. Deploying the Operator

To deploy the operator to your Kubernetes cluster:

```bash
make deploy IMG=<your-registry>/neurallog-operator:tag
```

### 6. Testing

#### Running Unit Tests

```bash
make test
```

#### Running Integration Tests

```bash
make test-integration
```

#### Running End-to-End Tests

```bash
make test-e2e
```

## Debugging

### Viewing Operator Logs

If the operator is deployed in the cluster:

```bash
kubectl logs -n neurallog-system deployment/neurallog-operator-controller-manager -c manager
```

If running locally, the logs will be printed to the console.

### Debugging with Delve

You can use [Delve](https://github.com/go-delve/delve) to debug the operator:

```bash
dlv debug ./main.go
```

## Code Style and Conventions

### Go Code Style

We follow the standard Go code style guidelines:

- Use `gofmt` or `goimports` to format your code
- Follow the [Effective Go](https://golang.org/doc/effective_go) guidelines
- Use meaningful variable and function names
- Add comments for exported functions and types

### Commit Messages

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification for commit messages:

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

Types include:
- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation changes
- `style`: Changes that do not affect the meaning of the code
- `refactor`: Code changes that neither fix a bug nor add a feature
- `test`: Adding or modifying tests
- `chore`: Changes to the build process or auxiliary tools

### Pull Requests

When submitting a pull request:

1. Make sure all tests pass
2. Update documentation if necessary
3. Add a clear description of the changes
4. Reference any related issues

## Project Structure

```
operator/
├── api/                  # API definitions
│   └── v1/               # v1 API
├── bin/                  # Compiled binaries
├── config/               # Kubernetes manifests
│   ├── crd/              # Custom Resource Definitions
│   ├── default/          # Default configuration
│   ├── manager/          # Manager configuration
│   ├── rbac/             # RBAC configuration
│   └── samples/          # Sample resources
├── controllers/          # Controller logic
├── docs/                 # Documentation
├── hack/                 # Scripts and tools
├── main.go               # Main entry point
└── Makefile              # Build and development commands
```

## Common Development Tasks

### Adding a New API Field

1. Add the field to the appropriate struct in `api/v1/tenant_types.go`
2. Run `make generate` to update the generated code
3. Update the controller logic to handle the new field
4. Update the documentation in `docs/api-reference.md`
5. Add tests for the new field

### Adding a New Reconciliation Feature

1. Create a new file in the `controllers` directory for the feature
2. Implement the reconciliation logic
3. Update the main reconciliation loop in `controllers/tenant_controller.go`
4. Add tests for the new feature
5. Update the documentation

### Updating the CRD

After making changes to the API:

1. Run `make manifests` to update the CRD manifests
2. Run `make install` to apply the updated CRD to the cluster

## Troubleshooting

### Common Issues

#### CRD Not Updating

If changes to the CRD are not being applied:

```bash
kubectl delete crd tenants.neurallog.io
make install
```

#### Controller Not Reconciling

Check the controller logs for errors:

```bash
kubectl logs -n neurallog-system deployment/neurallog-operator-controller-manager -c manager
```

#### Permission Issues

If the controller is encountering permission issues, check the RBAC configuration:

```bash
kubectl describe clusterrole neurallog-operator-manager-role
kubectl describe clusterrolebinding neurallog-operator-manager-rolebinding
```

## Additional Resources

- [Kubernetes Operator Pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)
- [Kubebuilder Book](https://book.kubebuilder.io/)
- [Controller Runtime](https://github.com/kubernetes-sigs/controller-runtime)
- [Kubernetes API Conventions](https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md)
