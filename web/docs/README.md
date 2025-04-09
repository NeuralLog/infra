# NeuralLog Admin Documentation

Welcome to the NeuralLog Admin documentation. This documentation will help you understand, install, and use the NeuralLog Admin interface to manage tenants in your Kubernetes cluster.

## Table of Contents

1. [Architecture Overview](./architecture.md)
2. [Getting Started](./getting-started.md)
3. [User Manual](./user-manual.md)
4. [API Documentation](./api-docs.md)
5. [Troubleshooting Guide](./troubleshooting.md)
6. [Glossary](./glossary.md)

## Introduction

NeuralLog Admin is a web-based interface for managing NeuralLog tenants in a Kubernetes cluster. It allows administrators to create, view, update, and delete tenants, as well as monitor their status.

Each tenant consists of:
- A dedicated Kubernetes namespace
- A NeuralLog server deployment
- A Redis statefulset for data storage
- Optional network policies for isolation

The NeuralLog Admin interface communicates with the Kubernetes API to manage these resources through a custom Tenant resource definition.

## Quick Links

- [Architecture Overview](./architecture.md): Understand the system architecture and components.
- [Getting Started](./getting-started.md): Install and set up the NeuralLog Admin interface.
- [User Manual](./user-manual.md): Learn how to use the NeuralLog Admin interface.
- [API Documentation](./api-docs.md): Explore the API endpoints for programmatic access.
- [Troubleshooting Guide](./troubleshooting.md): Solve common issues.
- [Glossary](./glossary.md): Definitions of terms used in the documentation.

## Screenshots

![Dashboard](./images/dashboard.png)
*The NeuralLog Admin dashboard showing a list of tenants*

![Tenant Creation](./images/create-tenant-form.png)
*Creating a new tenant*

![Tenant Details](./images/tenant-details.png)
*Viewing tenant details*

## Contributing

We welcome contributions to the NeuralLog Admin interface and documentation. Please see the [Contributing Guide](../CONTRIBUTING.md) for more information.

## License

NeuralLog Admin is licensed under the [MIT License](../LICENSE).
