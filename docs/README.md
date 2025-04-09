# NeuralLog Infrastructure Documentation

Welcome to the comprehensive documentation for the NeuralLog infrastructure. This documentation provides detailed information about the infrastructure components, architecture, installation, configuration, and operation of NeuralLog.

## Documentation Structure

- [Architecture Guide](architecture.md): Detailed architecture documentation with diagrams
- [Installation Guide](installation.md): Step-by-step installation instructions
- [Operator Guide](operator.md): Detailed documentation for the Tenant Operator
- [Kubernetes Configuration Guide](kubernetes.md): Documentation for Kubernetes resources
- [Docker Configuration Guide](docker.md): Documentation for Docker configurations
- [Redis Configuration Guide](redis.md): Documentation for Redis setup
- [Auth Service Guide](auth.md): Documentation for the Auth Service
- [Development Guide](development.md): Guide for developers working on the infrastructure
- [Troubleshooting Guide](troubleshooting.md): Common issues and solutions
- [Security Guide](security.md): Security considerations and best practices

## Quick Start

For a quick start, follow these steps:

1. **Prerequisites**: Ensure you have Docker, kubectl, and kind installed
2. **Clone the Repository**: `git clone https://github.com/NeuralLog/infra.git`
3. **Initialize the Environment**: Run `./scripts/Initialize-DevEnvironment.ps1` (Windows) or `./scripts/initialize-dev-environment.sh` (Linux/macOS)
4. **Start the Development Environment**: Run `./scripts/Start-DevEnvironment.ps1` (Windows) or `./scripts/start-dev-env.sh` (Linux/macOS)

For detailed instructions, refer to the [Installation Guide](installation.md).

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
├── operator/             # Kubernetes operator for tenant management
│   ├── api/              # API definitions
│   ├── controllers/      # Controller logic
│   └── config/           # Operator configurations
├── scripts/              # Utility scripts
│   ├── *.ps1             # PowerShell scripts for Windows
│   └── *.sh              # Bash scripts for Linux/macOS
└── docs/                 # Documentation
```

## Contributing

We welcome contributions to the NeuralLog infrastructure. Please refer to the [Development Guide](development.md) for information on how to contribute.

## License

Copyright 2023 NeuralLog Authors.

Licensed under the Apache License, Version 2.0.
