# NeuralLog Development Guide

This guide provides detailed information for developers working on the NeuralLog infrastructure.

## Table of Contents

- [Development Environment Setup](#development-environment-setup)
- [Repository Structure](#repository-structure)
- [Development Workflow](#development-workflow)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Debugging](#debugging)
- [Documentation](#documentation)
- [Pull Requests](#pull-requests)
- [Continuous Integration](#continuous-integration)
- [Release Process](#release-process)
- [Advanced Development Topics](#advanced-development-topics)

## Development Environment Setup

### Prerequisites

- **Docker**: [Docker Desktop](https://www.docker.com/products/docker-desktop/) (Windows/macOS) or Docker Engine (Linux)
- **kubectl**: [Kubernetes CLI](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- **kind**: [Kubernetes in Docker](https://kind.sigs.k8s.io/docs/user/quick-start/#installation) (for local testing)
- **Git**: [Git](https://git-scm.com/downloads)
- **Node.js**: [Node.js](https://nodejs.org/) (for development)
- **Go**: [Go](https://golang.org/dl/) (for operator development)

### Setting Up the Development Environment

1. Clone the repository:

```bash
git clone https://github.com/NeuralLog/infra.git
cd infra
```

2. Initialize the development environment:

#### Windows (PowerShell)

```powershell
.\scripts\Initialize-DevEnvironment.ps1
```

#### Linux/macOS (Bash)

```bash
chmod +x scripts/*.sh
./scripts/initialize-dev-environment.sh
```

3. Start the development environment:

#### Windows (PowerShell)

```powershell
.\scripts\Start-DevEnvironment.ps1
```

#### Linux/macOS (Bash)

```bash
./scripts/start-dev-env.sh
```

### IDE Setup

#### Visual Studio Code

1. Install Visual Studio Code
2. Install the following extensions:
   - Docker
   - Kubernetes
   - YAML
   - Go
   - ESLint
   - Prettier

#### JetBrains IDEs

1. Install IntelliJ IDEA or GoLand
2. Install the following plugins:
   - Docker
   - Kubernetes
   - YAML
   - Go (for GoLand)
   - Node.js (for IntelliJ IDEA)

## Repository Structure

```
infra/
├── kubernetes/           # Kubernetes configurations
│   ├── base/             # Base Kubernetes resources
│   │   ├── server/       # Server resources
│   │   └── redis/        # Redis resources
│   └── overlays/         # Kustomize overlays
│       └── test/         # Test environment
├── docker/               # Docker configurations
│   ├── server/           # Server Docker files
│   └── dev/              # Development Docker files
├── redis/                # Redis configuration
│   ├── conf/             # Redis config files
│   └── scripts/          # Redis scripts
├── operator/             # Kubernetes operator for tenant management
│   ├── api/              # API definitions
│   ├── controllers/      # Controller logic
│   └── config/           # Operator configurations
├── scripts/              # Utility scripts
│   ├── *.ps1             # PowerShell scripts for Windows
│   └── *.sh              # Bash scripts for Linux/macOS
└── docs/                 # Documentation
```

## Development Workflow

### Feature Development

1. Create a new branch for the feature:

```bash
git checkout -b feature/my-feature
```

2. Make changes to the code
3. Test the changes
4. Commit the changes:

```bash
git add .
git commit -m "Add my feature"
```

5. Push the changes:

```bash
git push origin feature/my-feature
```

6. Create a pull request

### Bug Fixes

1. Create a new branch for the bug fix:

```bash
git checkout -b fix/my-bug-fix
```

2. Make changes to the code
3. Test the changes
4. Commit the changes:

```bash
git add .
git commit -m "Fix my bug"
```

5. Push the changes:

```bash
git push origin fix/my-bug-fix
```

6. Create a pull request

### Code Reviews

1. All code changes must be reviewed before merging
2. Reviewers should check for:
   - Code quality
   - Test coverage
   - Documentation
   - Security issues
   - Performance issues

## Coding Standards

### General Guidelines

- Use consistent formatting
- Write clear and concise code
- Add comments for complex logic
- Follow the principle of least surprise
- Keep functions and methods small and focused
- Use meaningful variable and function names

### YAML Guidelines

- Use 2 spaces for indentation
- Use lowercase for keys
- Use kebab-case for resource names
- Add comments for complex configurations

### Go Guidelines

- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` to format code
- Use `golint` to check for style issues
- Use `go vet` to check for potential errors
- Write tests for all code

### Shell Script Guidelines

- Use `shellcheck` to check for issues
- Add error handling
- Add comments for complex logic
- Make scripts executable
- Add a shebang line

## Testing

### Testing the Infrastructure

#### Testing Kubernetes Configurations

```bash
# Validate Kubernetes configurations
kubectl kustomize kubernetes/base | kubectl apply --dry-run=client -f -

# Test with kind
./scripts/Setup-TestCluster.ps1
```

#### Testing Docker Configurations

```bash
# Build the server image
./scripts/Build-ServerImage.ps1

# Test the development environment
./scripts/Start-DevEnvironment.ps1
```

### Testing the Operator

#### Unit Tests

```bash
cd operator
go test ./...
```

#### Integration Tests

```bash
cd operator
make test
```

#### End-to-End Tests

```bash
# Install the operator
kubectl apply -f operator/config/crd/bases
kubectl apply -f operator/config/rbac
kubectl apply -f operator/config/manager

# Create a test tenant
kubectl apply -f operator/config/samples/neurallog_v1_tenant.yaml

# Verify the tenant
kubectl get tenants
kubectl get pods -n tenant-sample-tenant
```

## Debugging

### Debugging Kubernetes Resources

```bash
# Get pod logs
kubectl logs -l app=neurallog-server -n neurallog

# Describe resources
kubectl describe deployment neurallog-server -n neurallog
kubectl describe statefulset redis -n neurallog

# Port forward to services
kubectl port-forward svc/neurallog-server 3030:3030 -n neurallog
```

### Debugging Docker Containers

```bash
# Get container logs
docker logs <container_id>

# Exec into containers
docker exec -it <container_id> sh

# Inspect containers
docker inspect <container_id>
```

### Debugging the Operator

```bash
# Run the operator locally
cd operator
make run

# Get operator logs
kubectl logs -l control-plane=controller-manager -n system
```

## Documentation

### Documentation Guidelines

- Write clear and concise documentation
- Use Markdown for documentation
- Include examples and code snippets
- Keep documentation up-to-date
- Document all configuration options

### Generating Documentation

```bash
# Generate operator API documentation
cd operator
make docs
```

## Pull Requests

### Pull Request Guidelines

- Create a descriptive title
- Include a detailed description
- Reference related issues
- Include screenshots or GIFs for UI changes
- Include test results
- Ensure all tests pass
- Ensure documentation is updated

### Pull Request Template

```markdown
## Description

[Description of the changes]

## Related Issues

[Related issues]

## Type of Change

- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Checklist

- [ ] I have tested these changes
- [ ] I have updated the documentation
- [ ] I have added tests for my changes
- [ ] All tests pass
```

## Continuous Integration

### CI/CD Pipeline

The CI/CD pipeline includes the following steps:

1. **Lint**: Check code style and formatting
2. **Build**: Build Docker images and Go binaries
3. **Test**: Run unit and integration tests
4. **Deploy**: Deploy to test environment
5. **E2E Tests**: Run end-to-end tests
6. **Release**: Release to production

### GitHub Actions

```yaml
name: CI/CD

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    - name: Lint
      run: |
        # Lint code

  build:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    - name: Build
      run: |
        # Build code

  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    - name: Test
      run: |
        # Test code
```

## Release Process

### Versioning

NeuralLog follows [Semantic Versioning](https://semver.org/):

- **Major**: Incompatible API changes
- **Minor**: Backwards-compatible new features
- **Patch**: Backwards-compatible bug fixes

### Release Steps

1. Update version numbers
2. Update CHANGELOG.md
3. Create a release branch
4. Build and test the release
5. Create a GitHub release
6. Push Docker images to Docker Hub

### Release Checklist

- [ ] Update version numbers
- [ ] Update CHANGELOG.md
- [ ] Create release branch
- [ ] Build and test release
- [ ] Create GitHub release
- [ ] Push Docker images

## Advanced Development Topics

### Custom Resource Definitions (CRDs)

CRDs are defined in the `operator/api/v1` directory:

```go
// Tenant is the Schema for the tenants API
type Tenant struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TenantSpec   `json:"spec,omitempty"`
	Status TenantStatus `json:"status,omitempty"`
}
```

### Controllers

Controllers are defined in the `operator/controllers` directory:

```go
// TenantReconciler reconciles a Tenant object
type TenantReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// Reconcile is part of the main kubernetes reconciliation loop
func (r *TenantReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// Reconciliation logic
}
```

### Kubernetes Custom Resources

Custom resources are defined in YAML:

```yaml
apiVersion: neurallog.io/v1
kind: Tenant
metadata:
  name: sample-tenant
spec:
  displayName: "Sample Tenant"
  description: "A sample tenant for demonstration purposes"
```

### Multi-Tenancy

Multi-tenancy is implemented using the Tenant operator:

1. Each tenant has a dedicated namespace
2. Each tenant has dedicated resources
3. Network policies isolate tenants
4. Resource quotas prevent resource starvation
