# API Development

This document provides guidance on developing and extending the NeuralLog Tenant Operator API.

## Table of Contents

- [API Design Principles](#api-design-principles)
- [API Structure](#api-structure)
- [Adding New API Fields](#adding-new-api-fields)
- [Versioning the API](#versioning-the-api)
- [Validation Rules](#validation-rules)
- [Conversion Webhooks](#conversion-webhooks)
- [OpenAPI Documentation](#openapi-documentation)
- [API Conventions](#api-conventions)
- [Testing API Changes](#testing-api-changes)
- [API Compatibility](#api-compatibility)

## API Design Principles

The NeuralLog Tenant Operator API follows these design principles:

1. **Kubernetes-Native**: Follow Kubernetes API conventions
2. **Declarative**: Focus on desired state, not imperative commands
3. **Extensible**: Easy to extend with new fields and capabilities
4. **Backward Compatible**: Maintain compatibility across versions
5. **Well-Documented**: Clear documentation for all fields
6. **Validated**: Strong validation for all fields
7. **Defaulted**: Sensible defaults for optional fields

## API Structure

The NeuralLog Tenant Operator API is structured as follows:

```
Tenant
├── Spec
│   ├── DisplayName
│   ├── Description
│   ├── Version
│   ├── UpgradeStrategy
│   ├── Resources
│   ├── Server
│   ├── Redis
│   ├── NetworkPolicy
│   ├── Monitoring
│   ├── Backup
│   ├── Lifecycle
│   └── Integrations
└── Status
    ├── Conditions
    ├── Phase
    ├── Namespace
    ├── ServerStatus
    ├── RedisStatus
    ├── ObservedGeneration
    ├── LastReconcileTime
    ├── Components
    ├── URLs
    ├── Metrics
    └── BackupStatus
```

## Adding New API Fields

To add a new field to the API:

1. **Identify the Need**: Determine why the field is needed
2. **Choose the Right Location**: Decide where the field should be added
3. **Define the Field**: Define the field type and documentation
4. **Add Validation**: Add validation rules for the field
5. **Set Defaults**: Set sensible defaults for the field
6. **Update Documentation**: Update the API documentation
7. **Update Tests**: Update tests to cover the new field

### Example: Adding a New Field

```go
// Add a new field to the ServerSpec
type ServerSpec struct {
    // Existing fields...

    // LogRotation defines the log rotation configuration
    // +optional
    LogRotation *LogRotationSpec `json:"logRotation,omitempty"`
}

// LogRotationSpec defines the log rotation configuration
type LogRotationSpec struct {
    // Enabled indicates whether log rotation is enabled
    // +optional
    // +kubebuilder:default=true
    Enabled *bool `json:"enabled,omitempty"`

    // MaxSize is the maximum size of the log file before it gets rotated
    // +optional
    // +kubebuilder:default="100Mi"
    MaxSize string `json:"maxSize,omitempty"`

    // MaxAge is the maximum number of days to retain old log files
    // +optional
    // +kubebuilder:default=7
    MaxAge *int32 `json:"maxAge,omitempty"`

    // MaxBackups is the maximum number of old log files to retain
    // +optional
    // +kubebuilder:default=10
    MaxBackups *int32 `json:"maxBackups,omitempty"`

    // Compress indicates whether the rotated log files should be compressed
    // +optional
    // +kubebuilder:default=true
    Compress *bool `json:"compress,omitempty"`
}
```

## Versioning the API

The NeuralLog Tenant Operator API follows Kubernetes API versioning conventions:

1. **Alpha Versions** (v1alpha1): May be buggy and are disabled by default
2. **Beta Versions** (v1beta1): Well-tested but may have minor changes
3. **Stable Versions** (v1): Stable and will not change in incompatible ways

### Creating a New API Version

To create a new API version:

1. Create a new directory for the version:

```bash
mkdir -p api/v2
```

2. Copy the existing API files:

```bash
cp api/v1/*.go api/v2/
```

3. Update the package name:

```go
package v2
```

4. Update the API group and version in the `groupversion_info.go` file:

```go
// Package v2 contains API Schema definitions for the neurallog v2 API group
// +kubebuilder:object:generate=true
// +groupName=neurallog.io
package v2

import (
    "k8s.io/apimachinery/pkg/runtime/schema"
    "sigs.k8s.io/controller-runtime/pkg/scheme"
)

var (
    // GroupVersion is group version used to register these objects
    GroupVersion = schema.GroupVersion{Group: "neurallog.io", Version: "v2"}

    // SchemeBuilder is used to add go types to the GroupVersionKind scheme
    SchemeBuilder = &scheme.Builder{GroupVersion: GroupVersion}

    // AddToScheme adds the types in this group-version to the given scheme.
    AddToScheme = SchemeBuilder.AddToScheme
)
```

5. Update the API types with new fields and changes

6. Implement conversion webhooks (see [Conversion Webhooks](#conversion-webhooks))

7. Update the main.go file to register the new API version:

```go
func init() {
    utilruntime.Must(clientgoscheme.AddToScheme(scheme))

    utilruntime.Must(neurallogv1.AddToScheme(scheme))
    utilruntime.Must(neurallogv2.AddToScheme(scheme)) // Add the new version
    //+kubebuilder:scaffold:scheme
}
```

## Validation Rules

Validation rules ensure that API objects are valid. Use kubebuilder markers to add validation rules:

```go
// +kubebuilder:validation:Minimum=1
// +kubebuilder:validation:Maximum=100
// +kubebuilder:validation:ExclusiveMaximum=false
// +kubebuilder:validation:Format=date-time
// +kubebuilder:validation:Pattern=`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`
// +kubebuilder:validation:MaxLength=63
// +kubebuilder:validation:Enum=RollingUpdate;Recreate
```

### Common Validation Rules

- **Required Fields**: Mark fields as required in the struct tag
- **Numeric Ranges**: Use Minimum and Maximum markers
- **String Patterns**: Use Pattern marker for regex validation
- **Enumerations**: Use Enum marker for allowed values
- **String Length**: Use MaxLength and MinLength markers
- **Default Values**: Use default marker for default values

### Example: Adding Validation Rules

```go
type ServerSpec struct {
    // Replicas is the number of server replicas
    // +optional
    // +kubebuilder:validation:Minimum=1
    // +kubebuilder:validation:Maximum=100
    // +kubebuilder:default=1
    Replicas *int32 `json:"replicas,omitempty"`

    // Image is the Docker image for the server
    // +optional
    // +kubebuilder:default="neurallog/server:latest"
    Image string `json:"image,omitempty"`

    // LogLevel defines the log level for the server
    // +optional
    // +kubebuilder:validation:Enum=debug;info;warn;error
    // +kubebuilder:default=info
    LogLevel string `json:"logLevel,omitempty"`
}
```

## Conversion Webhooks

Conversion webhooks convert objects between different API versions. They are required when you have multiple API versions.

### Implementing Conversion Webhooks

1. Add the conversion webhook marker to the API types:

```go
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:webhook:path=/convert,mutating=false,failurePolicy=fail,sideEffects=None,groups=neurallog.io,resources=tenants,verbs=create;update,versions=v2,name=convert.tenant.neurallog.io,admissionReviewVersions=v1
type Tenant struct {
    metav1.TypeMeta   `json:",inline"`
    metav1.ObjectMeta `json:"metadata,omitempty"`

    Spec   TenantSpec   `json:"spec,omitempty"`
    Status TenantStatus `json:"status,omitempty"`
}
```

2. Implement the conversion functions:

```go
// ConvertTo converts this Tenant to the Hub version (v1).
func (src *Tenant) ConvertTo(dstRaw runtime.Object) error {
    dst := dstRaw.(*v1.Tenant)
    
    // Convert ObjectMeta
    dst.ObjectMeta = src.ObjectMeta
    
    // Convert Spec
    dst.Spec.DisplayName = src.Spec.DisplayName
    dst.Spec.Description = src.Spec.Description
    // ... convert other fields
    
    // Convert Status
    dst.Status.Phase = v1.TenantPhase(src.Status.Phase)
    dst.Status.Namespace = src.Status.Namespace
    // ... convert other fields
    
    return nil
}

// ConvertFrom converts from the Hub version (v1) to this version.
func (dst *Tenant) ConvertFrom(srcRaw runtime.Object) error {
    src := srcRaw.(*v1.Tenant)
    
    // Convert ObjectMeta
    dst.ObjectMeta = src.ObjectMeta
    
    // Convert Spec
    dst.Spec.DisplayName = src.Spec.DisplayName
    dst.Spec.Description = src.Spec.Description
    // ... convert other fields
    
    // Convert Status
    dst.Status.Phase = TenantPhase(src.Status.Phase)
    dst.Status.Namespace = src.Status.Namespace
    // ... convert other fields
    
    return nil
}
```

3. Register the conversion webhook in main.go:

```go
func main() {
    // ... existing code
    
    if err = (&neurallogv2.Tenant{}).SetupWebhookWithManager(mgr); err != nil {
        setupLog.Error(err, "unable to create webhook", "webhook", "Tenant")
        os.Exit(1)
    }
    
    // ... existing code
}
```

## OpenAPI Documentation

OpenAPI documentation is generated from the API types using controller-gen. Use kubebuilder markers to add documentation:

```go
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Namespace",type="string",JSONPath=".status.namespace"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
```

### Generating OpenAPI Documentation

```bash
make manifests
```

This will generate the CRD manifests with OpenAPI documentation.

### Example: Adding Documentation

```go
// Tenant is the Schema for the tenants API
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Namespace",type="string",JSONPath=".status.namespace"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type Tenant struct {
    metav1.TypeMeta   `json:",inline"`
    metav1.ObjectMeta `json:"metadata,omitempty"`

    // Spec defines the desired state of a NeuralLog Tenant
    Spec TenantSpec `json:"spec,omitempty"`

    // Status defines the observed state of a NeuralLog Tenant
    // +optional
    Status TenantStatus `json:"status,omitempty"`
}
```

## API Conventions

Follow these conventions when developing the API:

1. **Naming Conventions**:
   - Use CamelCase for field names
   - Use singular nouns for resource names
   - Use descriptive names that reflect the purpose

2. **Documentation Conventions**:
   - Document all fields with clear descriptions
   - Include examples where appropriate
   - Document default values and validation rules

3. **Field Conventions**:
   - Use pointers for optional fields with defaults
   - Use non-pointers for required fields
   - Use appropriate types (string, int32, bool, etc.)

4. **Status Conventions**:
   - Use conditions for complex status
   - Use phase for simple status
   - Include timestamps for status changes

## Testing API Changes

Test API changes thoroughly:

1. **Unit Tests**: Test conversion, validation, and defaulting
2. **Integration Tests**: Test API with the controller
3. **End-to-End Tests**: Test API in a real cluster

### Example: Testing API Changes

```go
func TestTenantValidation(t *testing.T) {
    tenant := &v1.Tenant{
        Spec: v1.TenantSpec{
            Server: v1.ServerSpec{
                Replicas: ptr.To(int32(0)), // Invalid: Minimum is 1
            },
        },
    }
    
    err := tenant.ValidateCreate()
    if err == nil {
        t.Error("Expected validation error, got nil")
    }
}
```

## API Compatibility

Maintain API compatibility across versions:

1. **Never Remove Fields**: Only add new fields
2. **Never Change Field Types**: Only add new fields with new types
3. **Never Change Field Semantics**: Only add new fields with new semantics
4. **Always Add Fields as Optional**: New fields should be optional
5. **Always Set Defaults for New Fields**: New fields should have defaults

### Breaking Changes

Avoid these breaking changes:

1. Removing fields
2. Changing field types
3. Changing field semantics
4. Making optional fields required
5. Changing default values

If a breaking change is necessary, create a new API version and implement conversion webhooks.
