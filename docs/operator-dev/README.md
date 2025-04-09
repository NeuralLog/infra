# NeuralLog Tenant Operator Developer Guide

Welcome to the NeuralLog Tenant Operator Developer Guide. This guide provides comprehensive information for developers working on extending and maintaining the NeuralLog Tenant Operator.

## Table of Contents

1. [Architecture Overview](architecture.md)
   - Operator Architecture
   - Controller-Runtime Framework
   - Reconciliation Loop
   - Custom Resource Definitions

2. [Development Environment Setup](development-environment.md)
   - Prerequisites
   - Setting Up the Development Environment
   - Building the Operator
   - Running the Operator Locally
   - Debugging the Operator

3. [API Development](api-development.md)
   - API Design Principles
   - Adding New API Fields
   - Versioning the API
   - Validation Rules
   - Conversion Webhooks
   - OpenAPI Documentation

4. [Controller Development](controller-development.md)
   - Controller Structure
   - Reconciliation Loop
   - Error Handling
   - Status Updates
   - Finalizers
   - Owner References
   - Event Recording

5. [Resource Management](resource-management.md)
   - Creating Resources
   - Updating Resources
   - Deleting Resources
   - Resource Ownership
   - Resource Status
   - Resource Cleanup

6. [Testing](testing.md)
   - Unit Testing
   - Integration Testing
   - End-to-End Testing
   - Test Fixtures
   - Test Mocks
   - Test Coverage

7. [Webhooks](webhooks.md)
   - Validation Webhooks
   - Mutation Webhooks
   - Conversion Webhooks
   - Webhook Development
   - Webhook Testing
   - Certificate Management

8. [Metrics and Monitoring](metrics-monitoring.md)
   - Prometheus Metrics
   - Custom Metrics
   - Alerting
   - Logging
   - Tracing
   - Dashboards

9. [Security](security.md)
   - RBAC
   - Pod Security
   - Network Security
   - Secret Management
   - Certificate Management
   - Security Best Practices

10. [Performance](performance.md)
    - Resource Usage
    - Scaling
    - Caching
    - Rate Limiting
    - Backoff and Retry
    - Performance Testing

11. [Deployment](deployment.md)
    - Operator Lifecycle Manager (OLM)
    - Helm Charts
    - Kustomize
    - Operator SDK
    - Continuous Deployment
    - Versioning and Upgrades

12. [Troubleshooting](troubleshooting.md)
    - Common Issues
    - Debugging Techniques
    - Logging
    - Profiling
    - Tracing
    - Support

## Getting Started

To get started with operator development, follow these steps:

1. Read the [Architecture Overview](architecture.md) to understand the operator architecture
2. Set up your [Development Environment](development-environment.md)
3. Learn about [API Development](api-development.md) and [Controller Development](controller-development.md)
4. Explore the [Testing](testing.md) guide to ensure your changes are well-tested

## Contributing

We welcome contributions to the NeuralLog Tenant Operator. Please follow these guidelines:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for your changes
5. Ensure all tests pass
6. Submit a pull request

## License

Copyright 2023 NeuralLog Authors.

Licensed under the Apache License, Version 2.0.
