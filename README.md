# NeuralLog Infrastructure

This repository contains infrastructure configurations for the NeuralLog system, including Kubernetes manifests, Docker configurations, and Redis setup.

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
├── scripts/              # Utility scripts
│   ├── *.ps1             # PowerShell scripts for Windows
│   └── *.sh              # Bash scripts for Linux/macOS
└── README.md             # Repository documentation
```

## Getting Started

### Prerequisites

- [Docker Desktop](https://www.docker.com/products/docker-desktop/)
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation) (for local Kubernetes testing)
- [Git](https://git-scm.com/downloads)
- [Node.js](https://nodejs.org/) (for development)

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

To start the development environment:

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

To stop the development environment:

```powershell
# Windows (PowerShell)
.\scripts\Stop-DevEnvironment.ps1

# Linux/macOS (Bash)
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

## Docker Configurations

Docker configurations for the NeuralLog server are located in the `docker` directory:

- **Production**: `docker/server/Dockerfile`
- **Development**: `docker/dev/Dockerfile.dev`

The development environment is configured using Docker Compose in `docker/dev/docker-compose.dev.yml`.
