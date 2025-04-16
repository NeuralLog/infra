# NeuralLog Infrastructure

This repository contains infrastructure configurations for the NeuralLog system, including Kubernetes manifests, Docker configurations, Redis setup, and authentication services.

## Repository Structure

```
infra/
├── kubernetes/           # Kubernetes configurations
│   ├── base/             # Base Kubernetes resources
│   │   ├── server/       # Server resources
│   │   ├── redis/        # Redis resources
│   │   ├── auth/         # Auth service resources
│   │   └── openfga/      # OpenFGA resources
│   └── overlays/         # Kustomize overlays
│       └── test/         # Test environment
├── docker/               # Docker configurations
│   ├── server/           # Server Docker files
│   ├── auth/             # Auth service Docker files
│   └── dev/              # Development Docker files
├── redis/                # Redis configuration
│   ├── conf/             # Redis config files
│   └── scripts/          # Redis scripts
├── scripts/              # Utility scripts
│   ├── *.ps1             # PowerShell scripts for Windows
│   └── *.sh              # Bash scripts for Linux/macOS
├── docker-compose.*.yml  # Docker Compose files for different components
└── README.md             # Repository documentation
```

## Getting Started

### Prerequisites

- [Docker Desktop](https://www.docker.com/products/docker-desktop/)
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation) (for local Kubernetes testing)
- [Git](https://git-scm.com/downloads)
- [Node.js](https://nodejs.org/) v22 or later (for development)
- [PowerShell](https://docs.microsoft.com/en-us/powershell/) (for Windows users)

### Initializing the Development Environment

To initialize the development environment and check prerequisites:

```powershell
# Windows (PowerShell)
.\scripts\Initialize-DevEnvironment.ps1

# Linux/macOS (Bash)
./scripts/initialize-dev-environment.sh
```

This script will:
- Check for required prerequisites
- Verify Docker is running
- Check for the NeuralLog/server repository
- Pull required Docker images

### Development Environment

There are two ways to start the development environment:

#### 1. Using the legacy development environment

```powershell
# Windows (PowerShell)
.\scripts\Start-DevEnvironment.ps1

# Linux/macOS (Bash)
./scripts/start-dev-env.sh
```

This will start:
- NeuralLog server on http://localhost:3030
- Redis on port 6379
- Redis Commander on http://localhost:8081

#### 2. Using the complete development environment

```powershell
# Windows (PowerShell)
.\scripts\Start-All.ps1
```

This will start:
- Verdaccio private npm registry on http://localhost:4873
- Redis on port 6379
- PostgreSQL on port 5432
- OpenFGA on http://localhost:8080
- Auth service on http://localhost:3040
- Logs server on http://localhost:3030

### Stopping the Development Environment

```powershell
# For the legacy environment
.\scripts\Stop-DevEnvironment.ps1

# For the complete environment
.\scripts\Stop-All.ps1

# Linux/macOS (Bash, legacy environment only)
./scripts/stop-dev-env.sh
```

### Building Docker Images

To build and optionally push the server Docker image:

```powershell
# Windows (PowerShell)
.\scripts\Build-ServerImage.ps1 -Tag "latest" -Push

# Linux/macOS (Bash)
./scripts/build-server-image.sh latest true
```

If no tag is provided, "latest" will be used.

## Test Kubernetes Environment

This repository includes configurations for a test Kubernetes environment using kind (Kubernetes in Docker).

### Setting Up the Test Environment

1. Create a kind cluster and deploy NeuralLog:
   ```powershell
   # Windows (PowerShell)
   .\scripts\Setup-TestCluster.ps1

   # Linux/macOS (Bash)
   ./scripts/setup-test-cluster.sh
   ```

2. Verify the deployment:
   ```bash
   kubectl get pods -n neurallog
   ```

3. Access the server:
   ```bash
   kubectl port-forward svc/neurallog-server 3030:3030 -n neurallog
   ```

4. Clean up the test environment:
   ```powershell
   # Windows (PowerShell)
   .\scripts\Cleanup-TestCluster.ps1

   # Linux/macOS (Bash)
   ./scripts/cleanup-test-cluster.sh
   ```

## Kubernetes Configuration

The Kubernetes configuration uses Kustomize to manage base resources and environment-specific overlays:

- **Base**: Contains the core resources (deployments, services, etc.)
- **Overlays**: Contains environment-specific configurations (test, staging, production)

To apply the Kubernetes configuration:

```bash
kubectl apply -k kubernetes/overlays/test
```

## Redis Configuration

The Redis configuration is located in the `redis/conf` directory. The default configuration is suitable for development and testing purposes.

Redis scripts for backup and restore are available in the `redis/scripts` directory:

```bash
# Backup Redis data
./redis/scripts/backup-redis.sh

# Restore Redis data
./redis/scripts/restore-redis.sh <backup-file>
```

## Auth Service

The Auth Service provides authentication and authorization capabilities for the NeuralLog platform. It uses OpenFGA (Fine-Grained Authorization) to manage permissions and supports multi-tenancy.

For more information, see the [Auth Service documentation](docs/auth.md).

### Running the Auth Service Locally

```powershell
# Start the Auth Service and OpenFGA
docker-compose -f docker-compose.auth.yml up -d
```

### Deploying to Kubernetes

```bash
# Deploy the Auth Service and OpenFGA
kubectl apply -k kubernetes/base/auth
kubectl apply -k kubernetes/base/openfga
```

## Docker Configurations

Docker configurations for the NeuralLog server are located in the `docker` directory:

- **Production**: `docker/server/Dockerfile`
- **Development**: `docker/dev/Dockerfile.dev`

The development environment is configured using Docker Compose files in the root directory:

- **docker-compose.combined.yml**: Combined configuration for all components
- **docker-compose.web.yml**: Configuration for the web application and Verdaccio
- **docker-compose.server.yml**: Configuration for the logs server and Redis
- **docker-compose.auth.yml**: Configuration for the auth service, PostgreSQL, and OpenFGA

Additionally, there's a development-specific Docker Compose file in `docker/dev/docker-compose.dev.yml`.

## Package Management

NeuralLog uses private packages for sharing code between components. These packages are published to a private Verdaccio registry.

### Shared Package

The `@neurallog/shared` package contains common types and utilities used across all NeuralLog components.

#### Publishing the Shared Package

```powershell
# Windows (PowerShell)
.\scripts\Publish-Shared.ps1
```

This script will:
- Start Verdaccio if it's not already running
- Configure npm to use the private registry for the @neurallog scope
- Build the shared package
- Publish the package to the private registry

#### Updating the Shared Package in All Repositories

```powershell
# Windows (PowerShell)
.\scripts\Update-Shared.ps1 [version]
```

This script will install the specified version (or latest by default) of the shared package in all repositories.

### TypeScript SDK

The `@neurallog/sdk` package provides a client library for interacting with the NeuralLog server.

#### Publishing the SDK

```powershell
# Windows (PowerShell)
.\scripts\Publish-SDK.ps1
```

This script will:
- Start Verdaccio if it's not already running
- Configure npm to use the private registry for the @neurallog scope
- Build the SDK
- Publish the SDK to the private registry

## Repository Management

Since NeuralLog is not a monorepo, we provide scripts to help manage the multiple repositories:

### Checking Repository Status

```powershell
# Windows (PowerShell)
.\scripts\Repo-Status.ps1
```

This script will show the git status of all repositories.

### Pulling All Repositories

```powershell
# Windows (PowerShell)
.\scripts\Pull-All.ps1
```

This script will pull the latest changes for all repositories.

### Pushing All Repositories

```powershell
# Windows (PowerShell)
.\scripts\Push-All.ps1 "Commit message"
```

This script will commit and push changes for all repositories with the specified commit message.

## Documentation

Detailed documentation is available in the [docs](./docs) directory:

- [API Reference](./docs/api.md)
- [Configuration](./docs/configuration.md)
- [Architecture](./docs/architecture.md)
- [Examples](./docs/examples)

For integration guides and tutorials, visit the [NeuralLog Documentation Site](https://neurallog.github.io/docs/).

## Contributing

Contributions are welcome! Please read our [Contributing Guide](./CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## Related NeuralLog Components

- [NeuralLog Auth](https://github.com/NeuralLog/auth) - Authentication and authorization
- [NeuralLog Server](https://github.com/NeuralLog/server) - Core server functionality
- [NeuralLog Web](https://github.com/NeuralLog/web) - Web interface components
- [NeuralLog TypeScript Client SDK](https://github.com/NeuralLog/typescript-client-sdk) - TypeScript client SDK
- [NeuralLog Java Client SDK](https://github.com/NeuralLog/Java-client-sdk) - Java client SDK
- [NeuralLog Python SDK](https://github.com/NeuralLog/python) - Python SDK
- [NeuralLog C# SDK](https://github.com/NeuralLog/csharp) - C# SDK
- [NeuralLog Go SDK](https://github.com/NeuralLog/go) - Go SDK

## License

MIT
