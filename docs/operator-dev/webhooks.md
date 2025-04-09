# Webhooks

This document provides guidance on developing and using webhooks in the NeuralLog Tenant Operator.

## Table of Contents

- [Introduction to Webhooks](#introduction-to-webhooks)
- [Types of Webhooks](#types-of-webhooks)
  - [Validation Webhooks](#validation-webhooks)
  - [Mutation Webhooks](#mutation-webhooks)
  - [Conversion Webhooks](#conversion-webhooks)
- [Webhook Development](#webhook-development)
  - [Setting Up Webhooks](#setting-up-webhooks)
  - [Implementing Webhooks](#implementing-webhooks)
  - [Registering Webhooks](#registering-webhooks)
- [Webhook Testing](#webhook-testing)
  - [Unit Testing Webhooks](#unit-testing-webhooks)
  - [Integration Testing Webhooks](#integration-testing-webhooks)
- [Certificate Management](#certificate-management)
  - [Generating Certificates](#generating-certificates)
  - [Configuring Certificates](#configuring-certificates)
- [Webhook Deployment](#webhook-deployment)
  - [Deploying Webhooks](#deploying-webhooks)
  - [Troubleshooting Webhooks](#troubleshooting-webhooks)
- [Best Practices](#best-practices)

## Introduction to Webhooks

Webhooks are HTTP callbacks that are triggered by events in the Kubernetes API server. They allow the operator to intercept and modify requests to the API server before they are processed.

Webhooks are used for:

1. **Validation**: Validate resources before they are created or updated
2. **Mutation**: Modify resources before they are created or updated
3. **Conversion**: Convert resources between different API versions

## Types of Webhooks

### Validation Webhooks

Validation webhooks validate resources before they are created, updated, or deleted. They can reject invalid resources.

```go
// +kubebuilder:webhook:path=/validate-neurallog-io-v1-tenant,mutating=false,failurePolicy=fail,sideEffects=None,groups=neurallog.io,resources=tenants,verbs=create;update;delete,versions=v1,name=vtenant.kb.io,admissionReviewVersions=v1
```

### Mutation Webhooks

Mutation webhooks modify resources before they are created or updated. They can set default values, add labels, or make other changes.

```go
// +kubebuilder:webhook:path=/mutate-neurallog-io-v1-tenant,mutating=true,failurePolicy=fail,sideEffects=None,groups=neurallog.io,resources=tenants,verbs=create;update,versions=v1,name=mtenant.kb.io,admissionReviewVersions=v1
```

### Conversion Webhooks

Conversion webhooks convert resources between different API versions. They are used when you have multiple API versions.

```go
// +kubebuilder:webhook:path=/convert,mutating=false,failurePolicy=fail,sideEffects=None,groups=neurallog.io,resources=tenants,verbs=create;update,versions=v2,name=convert.tenant.neurallog.io,admissionReviewVersions=v1
```

## Webhook Development

### Setting Up Webhooks

To set up webhooks, you need to:

1. Add webhook markers to your API types
2. Implement webhook functions
3. Register webhooks with the manager

### Implementing Webhooks

#### Validation Webhook

```go
// TenantValidator validates Tenants
type TenantValidator struct {
}

// ValidateCreate implements webhook.Validator
func (v *TenantValidator) ValidateCreate(ctx context.Context, obj runtime.Object) error {
    tenant := obj.(*neurallogv1.Tenant)
    
    // Validate tenant
    if tenant.Spec.Server.Replicas != nil && *tenant.Spec.Server.Replicas < 1 {
        return fmt.Errorf("server replicas must be at least 1")
    }
    
    return nil
}

// ValidateUpdate implements webhook.Validator
func (v *TenantValidator) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) error {
    oldTenant := oldObj.(*neurallogv1.Tenant)
    newTenant := newObj.(*neurallogv1.Tenant)
    
    // Validate tenant update
    if newTenant.Spec.Server.Replicas != nil && *newTenant.Spec.Server.Replicas < 1 {
        return fmt.Errorf("server replicas must be at least 1")
    }
    
    // Prevent changes to immutable fields
    if oldTenant.Spec.Redis.Storage != newTenant.Spec.Redis.Storage {
        return fmt.Errorf("redis storage is immutable")
    }
    
    return nil
}

// ValidateDelete implements webhook.Validator
func (v *TenantValidator) ValidateDelete(ctx context.Context, obj runtime.Object) error {
    tenant := obj.(*neurallogv1.Tenant)
    
    // Validate tenant deletion
    if tenant.Status.Phase == neurallogv1.TenantProvisioning {
        return fmt.Errorf("cannot delete tenant while it is being provisioned")
    }
    
    return nil
}
```

#### Mutation Webhook

```go
// TenantDefaulter defaults Tenants
type TenantDefaulter struct {
}

// Default implements webhook.Defaulter
func (d *TenantDefaulter) Default(ctx context.Context, obj runtime.Object) error {
    tenant := obj.(*neurallogv1.Tenant)
    
    // Set default values
    if tenant.Spec.Server.Replicas == nil {
        tenant.Spec.Server.Replicas = ptr.To(int32(1))
    }
    
    if tenant.Spec.Server.Image == "" {
        tenant.Spec.Server.Image = "neurallog/server:latest"
    }
    
    if tenant.Spec.Redis.Replicas == nil {
        tenant.Spec.Redis.Replicas = ptr.To(int32(1))
    }
    
    if tenant.Spec.Redis.Image == "" {
        tenant.Spec.Redis.Image = "redis:7-alpine"
    }
    
    if tenant.Spec.Redis.Storage == "" {
        tenant.Spec.Redis.Storage = "1Gi"
    }
    
    return nil
}
```

#### Conversion Webhook

```go
// ConvertTo converts this Tenant to the Hub version (v1).
func (src *v2.Tenant) ConvertTo(dstRaw runtime.Object) error {
    dst := dstRaw.(*v1.Tenant)
    
    // Convert ObjectMeta
    dst.ObjectMeta = src.ObjectMeta
    
    // Convert Spec
    dst.Spec.DisplayName = src.Spec.DisplayName
    dst.Spec.Description = src.Spec.Description
    
    // Convert Server
    dst.Spec.Server.Replicas = src.Spec.Server.Replicas
    dst.Spec.Server.Image = src.Spec.Server.Image
    
    // Convert Redis
    dst.Spec.Redis.Replicas = src.Spec.Redis.Replicas
    dst.Spec.Redis.Image = src.Spec.Redis.Image
    dst.Spec.Redis.Storage = src.Spec.Redis.Storage
    
    // Convert Status
    dst.Status.Phase = v1.TenantPhase(src.Status.Phase)
    dst.Status.Namespace = src.Status.Namespace
    
    return nil
}

// ConvertFrom converts from the Hub version (v1) to this version.
func (dst *v2.Tenant) ConvertFrom(srcRaw runtime.Object) error {
    src := srcRaw.(*v1.Tenant)
    
    // Convert ObjectMeta
    dst.ObjectMeta = src.ObjectMeta
    
    // Convert Spec
    dst.Spec.DisplayName = src.Spec.DisplayName
    dst.Spec.Description = src.Spec.Description
    
    // Convert Server
    dst.Spec.Server.Replicas = src.Spec.Server.Replicas
    dst.Spec.Server.Image = src.Spec.Server.Image
    
    // Convert Redis
    dst.Spec.Redis.Replicas = src.Spec.Redis.Replicas
    dst.Spec.Redis.Image = src.Spec.Redis.Image
    dst.Spec.Redis.Storage = src.Spec.Redis.Storage
    
    // Convert Status
    dst.Status.Phase = v2.TenantPhase(src.Status.Phase)
    dst.Status.Namespace = src.Status.Namespace
    
    return nil
}
```

### Registering Webhooks

Register webhooks with the manager in the main.go file:

```go
func main() {
    // ... existing code
    
    if err = (&neurallogv1.Tenant{}).SetupWebhookWithManager(mgr); err != nil {
        setupLog.Error(err, "unable to create webhook", "webhook", "Tenant")
        os.Exit(1)
    }
    
    // ... existing code
}
```

## Webhook Testing

### Unit Testing Webhooks

Unit tests for webhooks focus on the webhook functions:

```go
func TestTenantValidator_ValidateCreate(t *testing.T) {
    validator := &webhooks.TenantValidator{}
    
    tests := []struct {
        name    string
        tenant  *neurallogv1.Tenant
        wantErr bool
    }{
        {
            name: "valid tenant",
            tenant: &neurallogv1.Tenant{
                ObjectMeta: metav1.ObjectMeta{
                    Name: "valid-tenant",
                },
                Spec: neurallogv1.TenantSpec{
                    DisplayName: "Valid Tenant",
                    Server: neurallogv1.ServerSpec{
                        Replicas: ptr.To(int32(1)),
                    },
                },
            },
            wantErr: false,
        },
        {
            name: "invalid replicas",
            tenant: &neurallogv1.Tenant{
                ObjectMeta: metav1.ObjectMeta{
                    Name: "invalid-tenant",
                },
                Spec: neurallogv1.TenantSpec{
                    DisplayName: "Invalid Tenant",
                    Server: neurallogv1.ServerSpec{
                        Replicas: ptr.To(int32(0)), // Invalid: Minimum is 1
                    },
                },
            },
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validator.ValidateCreate(context.Background(), tt.tenant)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### Integration Testing Webhooks

Integration tests for webhooks use the envtest package:

```go
var _ = Describe("Tenant webhooks", func() {
    Context("When creating a Tenant", func() {
        It("Should validate the tenant", func() {
            ctx := context.Background()
            
            // Create an invalid tenant
            tenant := &neurallogv1.Tenant{
                ObjectMeta: metav1.ObjectMeta{
                    Name: "invalid-tenant",
                },
                Spec: neurallogv1.TenantSpec{
                    DisplayName: "Invalid Tenant",
                    Server: neurallogv1.ServerSpec{
                        Replicas: ptr.To(int32(0)), // Invalid: Minimum is 1
                    },
                },
            }
            
            // Try to create the tenant
            err := k8sClient.Create(ctx, tenant)
            
            // Expect an error
            Expect(err).To(HaveOccurred())
            Expect(err.Error()).To(ContainSubstring("server replicas must be at least 1"))
        })
        
        It("Should default the tenant", func() {
            ctx := context.Background()
            
            // Create a tenant without defaults
            tenant := &neurallogv1.Tenant{
                ObjectMeta: metav1.ObjectMeta{
                    Name: "default-tenant",
                },
                Spec: neurallogv1.TenantSpec{
                    DisplayName: "Default Tenant",
                },
            }
            
            // Create the tenant
            Expect(k8sClient.Create(ctx, tenant)).To(Succeed())
            
            // Get the tenant
            createdTenant := &neurallogv1.Tenant{}
            Expect(k8sClient.Get(ctx, types.NamespacedName{Name: "default-tenant"}, createdTenant)).To(Succeed())
            
            // Check defaults
            Expect(createdTenant.Spec.Server.Replicas).NotTo(BeNil())
            Expect(*createdTenant.Spec.Server.Replicas).To(Equal(int32(1)))
            Expect(createdTenant.Spec.Server.Image).To(Equal("neurallog/server:latest"))
            Expect(createdTenant.Spec.Redis.Replicas).NotTo(BeNil())
            Expect(*createdTenant.Spec.Redis.Replicas).To(Equal(int32(1)))
            Expect(createdTenant.Spec.Redis.Image).To(Equal("redis:7-alpine"))
            Expect(createdTenant.Spec.Redis.Storage).To(Equal("1Gi"))
            
            // Cleanup
            Expect(k8sClient.Delete(ctx, tenant)).To(Succeed())
        })
    })
})
```

## Certificate Management

### Generating Certificates

Webhooks require TLS certificates for secure communication. The cert-manager project can generate and manage certificates:

```yaml
# config/certmanager/certificate.yaml
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: serving-cert
  namespace: system
spec:
  dnsNames:
  - webhook-service.system.svc
  - webhook-service.system.svc.cluster.local
  issuerRef:
    kind: Issuer
    name: selfsigned-issuer
  secretName: webhook-server-cert
```

### Configuring Certificates

Configure the webhook server to use the certificates:

```go
func main() {
    // ... existing code
    
    mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
        Scheme:             scheme,
        MetricsBindAddress: metricsAddr,
        Port:               9443,
        LeaderElection:     enableLeaderElection,
        LeaderElectionID:   "neurallog-operator.neurallog.io",
        CertDir:            "/tmp/k8s-webhook-server/serving-certs", // Certificate directory
    })
    
    // ... existing code
}
```

## Webhook Deployment

### Deploying Webhooks

Deploy webhooks using the following manifests:

```yaml
# config/webhook/manifests.yaml
apiVersion: admissionregistration.k8s.io/v1
kind: ValidatingWebhookConfiguration
metadata:
  name: validating-webhook-configuration
webhooks:
- admissionReviewVersions:
  - v1
  clientConfig:
    service:
      name: webhook-service
      namespace: system
      path: /validate-neurallog-io-v1-tenant
  failurePolicy: Fail
  name: vtenant.kb.io
  rules:
  - apiGroups:
    - neurallog.io
    apiVersions:
    - v1
    operations:
    - CREATE
    - UPDATE
    - DELETE
    resources:
    - tenants
  sideEffects: None
---
apiVersion: admissionregistration.k8s.io/v1
kind: MutatingWebhookConfiguration
metadata:
  name: mutating-webhook-configuration
webhooks:
- admissionReviewVersions:
  - v1
  clientConfig:
    service:
      name: webhook-service
      namespace: system
      path: /mutate-neurallog-io-v1-tenant
  failurePolicy: Fail
  name: mtenant.kb.io
  rules:
  - apiGroups:
    - neurallog.io
    apiVersions:
    - v1
    operations:
    - CREATE
    - UPDATE
    resources:
    - tenants
  sideEffects: None
```

### Troubleshooting Webhooks

Common webhook issues and solutions:

1. **Certificate Issues**: Ensure certificates are valid and properly configured
2. **Webhook Service Issues**: Ensure the webhook service is running and accessible
3. **Webhook Configuration Issues**: Ensure webhook configurations are correct
4. **Webhook Logic Issues**: Ensure webhook logic is correct

## Best Practices

Follow these best practices for webhook development:

1. **Keep Webhooks Simple**: Keep webhook logic simple and focused
2. **Handle Errors Gracefully**: Handle errors gracefully and provide clear error messages
3. **Test Webhooks Thoroughly**: Test webhooks thoroughly to ensure they work correctly
4. **Use Appropriate Failure Policies**: Use appropriate failure policies for webhooks
5. **Document Webhook Behavior**: Document webhook behavior for users
6. **Monitor Webhook Performance**: Monitor webhook performance to ensure they don't impact API server performance
7. **Use Conversion Webhooks for API Versioning**: Use conversion webhooks for API versioning
8. **Use Validation Webhooks for Complex Validation**: Use validation webhooks for complex validation
9. **Use Mutation Webhooks for Defaulting**: Use mutation webhooks for defaulting
10. **Use Cert-Manager for Certificate Management**: Use cert-manager for certificate management
