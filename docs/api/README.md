# NeuralLog Tenant API Documentation

This directory contains the API documentation for the NeuralLog Tenant Operator.

## Overview

The NeuralLog Tenant API is a Kubernetes custom resource API that allows you to create and manage NeuralLog tenants. The API is defined using the Kubernetes Custom Resource Definition (CRD) mechanism.

## API Documentation

The API documentation is available in the following formats:

- **OpenAPI Specification**: [openapi.yaml](openapi.yaml)
- **Interactive Documentation**: [index.html](index.html)

## Using the API Documentation

### OpenAPI Specification

The OpenAPI specification is a machine-readable description of the API. It can be used with various tools to generate client libraries, documentation, and more.

### Interactive Documentation

The interactive documentation provides a user-friendly interface for exploring the API. To use it:

1. Open the `index.html` file in a web browser
2. Browse the API endpoints and models
3. Try out API requests (requires a running Kubernetes cluster with the NeuralLog Tenant Operator installed)

## API Versions

The NeuralLog Tenant API follows Kubernetes API versioning conventions:

- **Alpha Versions** (v1alpha1): May be buggy and are disabled by default
- **Beta Versions** (v1beta1): Well-tested but may have minor changes
- **Stable Versions** (v1): Stable and will not change in incompatible ways

The current version is `v1`.

## API Resources

The NeuralLog Tenant API defines the following resources:

- **Tenant**: Represents a NeuralLog tenant

## API Fields

The Tenant resource includes the following fields:

- **Spec**: Desired state of the tenant
  - Basic information (displayName, description)
  - Server configuration
  - Redis configuration
  - Network policy configuration
  - Monitoring configuration
  - Backup configuration
  - Lifecycle hooks
  - Integrations
- **Status**: Observed state of the tenant
  - Phase
  - Conditions
  - Component status
  - URLs
  - Metrics
  - Backup status

For detailed information about each field, refer to the [API Reference](../api-reference.md).

## Examples

For examples of using the NeuralLog Tenant API, refer to the [examples](../../operator/config/samples) directory.

## API Conventions

The NeuralLog Tenant API follows Kubernetes API conventions:

- **Naming**: CamelCase for field names
- **Versioning**: API versions are in the format `v1`, `v1beta1`, etc.
- **Status**: Status fields are read-only and updated by the controller
- **Spec**: Spec fields are writable and define the desired state
- **Metadata**: Standard Kubernetes metadata (name, namespace, labels, annotations, etc.)

## API Validation

The NeuralLog Tenant API includes validation rules to ensure that resources are valid:

- **Required Fields**: Some fields are required
- **Enumerated Values**: Some fields have a limited set of allowed values
- **Numeric Ranges**: Some numeric fields have minimum and maximum values
- **String Patterns**: Some string fields must match specific patterns
- **Default Values**: Some fields have default values

For detailed information about validation rules, refer to the [API Reference](../api-reference.md#field-validation-rules).

## API Compatibility

The NeuralLog Tenant API maintains compatibility across versions:

- **Never Remove Fields**: Fields are never removed, only deprecated
- **Never Change Field Types**: Field types are never changed
- **Never Change Field Semantics**: Field semantics are never changed
- **Always Add Fields as Optional**: New fields are always optional
- **Always Set Defaults for New Fields**: New fields always have defaults

For detailed information about API compatibility, refer to the [API Reference](../api-reference.md#api-versioning).
