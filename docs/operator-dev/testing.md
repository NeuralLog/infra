# Testing the NeuralLog Tenant Operator

This document provides guidance on testing the NeuralLog Tenant Operator, including unit tests, integration tests, and end-to-end tests.

## Table of Contents

- [Testing Strategy](#testing-strategy)
- [Unit Testing](#unit-testing)
  - [Testing API Types](#testing-api-types)
  - [Testing Controllers](#testing-controllers)
  - [Testing Webhooks](#testing-webhooks)
- [Integration Testing](#integration-testing)
  - [Testing with envtest](#testing-with-envtest)
  - [Testing with a Fake Client](#testing-with-a-fake-client)
- [End-to-End Testing](#end-to-end-testing)
  - [Testing with kind](#testing-with-kind)
  - [Testing with Ginkgo and Gomega](#testing-with-ginkgo-and-gomega)
- [Test Fixtures](#test-fixtures)
- [Test Mocks](#test-mocks)
- [Test Coverage](#test-coverage)
- [Continuous Integration](#continuous-integration)
- [Best Practices](#best-practices)

## Testing Strategy

The NeuralLog Tenant Operator uses a multi-layered testing strategy:

1. **Unit Tests**: Test individual components in isolation
2. **Integration Tests**: Test components together with a simulated Kubernetes API
3. **End-to-End Tests**: Test the entire operator in a real Kubernetes cluster

This strategy ensures that the operator is thoroughly tested at all levels.

## Unit Testing

Unit tests verify that individual components work correctly in isolation. They are fast, reliable, and provide immediate feedback.

### Testing API Types

Unit tests for API types focus on validation, defaulting, and conversion.

```go
package v1_test

import (
    "testing"

    "github.com/stretchr/testify/assert"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/utils/ptr"

    neurallogv1 "github.com/neurallog/infra/operator/api/v1"
)

func TestTenantValidation(t *testing.T) {
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
            err := tt.tenant.ValidateCreate()
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}

func TestTenantDefaulting(t *testing.T) {
    tenant := &neurallogv1.Tenant{
        ObjectMeta: metav1.ObjectMeta{
            Name: "default-tenant",
        },
        Spec: neurallogv1.TenantSpec{
            DisplayName: "Default Tenant",
        },
    }

    tenant.Default()

    // Check that defaults were set
    assert.Equal(t, ptr.To(int32(1)), tenant.Spec.Server.Replicas)
    assert.Equal(t, "neurallog/server:latest", tenant.Spec.Server.Image)
}
```

### Testing Controllers

Unit tests for controllers focus on the reconciliation logic.

```go
package controllers_test

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/runtime"
    "k8s.io/apimachinery/pkg/types"
    "k8s.io/client-go/kubernetes/scheme"
    ctrl "sigs.k8s.io/controller-runtime"
    "sigs.k8s.io/controller-runtime/pkg/client/fake"

    neurallogv1 "github.com/neurallog/infra/operator/api/v1"
    "github.com/neurallog/infra/operator/controllers"
)

func TestTenantReconciler_Reconcile(t *testing.T) {
    // Register the scheme
    s := runtime.NewScheme()
    _ = scheme.AddToScheme(s)
    _ = neurallogv1.AddToScheme(s)

    // Create a tenant
    tenant := &neurallogv1.Tenant{
        ObjectMeta: metav1.ObjectMeta{
            Name: "test-tenant",
        },
        Spec: neurallogv1.TenantSpec{
            DisplayName: "Test Tenant",
        },
    }

    // Create a fake client
    client := fake.NewClientBuilder().
        WithScheme(s).
        WithObjects(tenant).
        Build()

    // Create a reconciler
    reconciler := &controllers.TenantReconciler{
        Client: client,
        Scheme: s,
    }

    // Reconcile
    req := ctrl.Request{
        NamespacedName: types.NamespacedName{
            Name: "test-tenant",
        },
    }
    result, err := reconciler.Reconcile(context.Background(), req)

    // Check the result
    assert.NoError(t, err)
    assert.Equal(t, ctrl.Result{}, result)

    // Check that the tenant was updated
    updatedTenant := &neurallogv1.Tenant{}
    err = client.Get(context.Background(), types.NamespacedName{Name: "test-tenant"}, updatedTenant)
    assert.NoError(t, err)
    assert.Equal(t, neurallogv1.TenantPending, updatedTenant.Status.Phase)

    // Check that a namespace was created
    namespace := &corev1.Namespace{}
    err = client.Get(context.Background(), types.NamespacedName{Name: "tenant-test-tenant"}, namespace)
    assert.NoError(t, err)
    assert.Equal(t, "test-tenant", namespace.Labels["neurallog.io/tenant"])
}
```

### Testing Webhooks

Unit tests for webhooks focus on validation, mutation, and conversion.

```go
package webhooks_test

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/utils/ptr"

    neurallogv1 "github.com/neurallog/infra/operator/api/v1"
    "github.com/neurallog/infra/operator/webhooks"
)

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

## Integration Testing

Integration tests verify that components work correctly together. They use a simulated Kubernetes API to test the operator without a real cluster.

### Testing with envtest

The `envtest` package provides a Kubernetes API server for testing.

```go
package controllers_test

import (
    "context"
    "path/filepath"
    "testing"
    "time"

    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/types"
    "k8s.io/client-go/kubernetes/scheme"
    "k8s.io/client-go/rest"
    ctrl "sigs.k8s.io/controller-runtime"
    "sigs.k8s.io/controller-runtime/pkg/client"
    "sigs.k8s.io/controller-runtime/pkg/envtest"
    logf "sigs.k8s.io/controller-runtime/pkg/log"
    "sigs.k8s.io/controller-runtime/pkg/log/zap"

    neurallogv1 "github.com/neurallog/infra/operator/api/v1"
    "github.com/neurallog/infra/operator/controllers"
)

var cfg *rest.Config
var k8sClient client.Client
var testEnv *envtest.Environment

func TestControllers(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Controller Suite")
}

var _ = BeforeSuite(func() {
    logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

    By("bootstrapping test environment")
    testEnv = &envtest.Environment{
        CRDDirectoryPaths:     []string{filepath.Join("..", "config", "crd", "bases")},
        ErrorIfCRDPathMissing: true,
    }

    var err error
    cfg, err = testEnv.Start()
    Expect(err).NotTo(HaveOccurred())
    Expect(cfg).NotTo(BeNil())

    err = neurallogv1.AddToScheme(scheme.Scheme)
    Expect(err).NotTo(HaveOccurred())

    k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
    Expect(err).NotTo(HaveOccurred())
    Expect(k8sClient).NotTo(BeNil())

    k8sManager, err := ctrl.NewManager(cfg, ctrl.Options{
        Scheme: scheme.Scheme,
    })
    Expect(err).ToNot(HaveOccurred())

    err = (&controllers.TenantReconciler{
        Client: k8sManager.GetClient(),
        Scheme: k8sManager.GetScheme(),
    }).SetupWithManager(k8sManager)
    Expect(err).ToNot(HaveOccurred())

    go func() {
        err = k8sManager.Start(ctrl.SetupSignalHandler())
        Expect(err).ToNot(HaveOccurred())
    }()
})

var _ = AfterSuite(func() {
    By("tearing down the test environment")
    err := testEnv.Stop()
    Expect(err).NotTo(HaveOccurred())
})

var _ = Describe("Tenant controller", func() {
    Context("When creating a Tenant", func() {
        It("Should create a namespace", func() {
            ctx := context.Background()
            tenant := &neurallogv1.Tenant{
                ObjectMeta: metav1.ObjectMeta{
                    Name: "test-tenant",
                },
                Spec: neurallogv1.TenantSpec{
                    DisplayName: "Test Tenant",
                },
            }

            Expect(k8sClient.Create(ctx, tenant)).Should(Succeed())

            // Wait for the namespace to be created
            namespaceName := types.NamespacedName{Name: "tenant-test-tenant"}
            namespace := &corev1.Namespace{}
            Eventually(func() bool {
                err := k8sClient.Get(ctx, namespaceName, namespace)
                return err == nil
            }, 10*time.Second, 1*time.Second).Should(BeTrue())

            // Check that the namespace has the correct labels
            Expect(namespace.Labels["neurallog.io/tenant"]).Should(Equal("test-tenant"))

            // Cleanup
            Expect(k8sClient.Delete(ctx, tenant)).Should(Succeed())
        })
    })
})
```

### Testing with a Fake Client

The `fake` client package provides a simulated client for testing.

```go
package controllers_test

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/runtime"
    "k8s.io/apimachinery/pkg/types"
    "k8s.io/client-go/kubernetes/scheme"
    ctrl "sigs.k8s.io/controller-runtime"
    "sigs.k8s.io/controller-runtime/pkg/client/fake"

    neurallogv1 "github.com/neurallog/infra/operator/api/v1"
    "github.com/neurallog/infra/operator/controllers"
)

func TestTenantReconciler_ReconcileWithFakeClient(t *testing.T) {
    // Register the scheme
    s := runtime.NewScheme()
    _ = scheme.AddToScheme(s)
    _ = neurallogv1.AddToScheme(s)

    // Create a tenant
    tenant := &neurallogv1.Tenant{
        ObjectMeta: metav1.ObjectMeta{
            Name: "test-tenant",
        },
        Spec: neurallogv1.TenantSpec{
            DisplayName: "Test Tenant",
        },
    }

    // Create a fake client
    client := fake.NewClientBuilder().
        WithScheme(s).
        WithObjects(tenant).
        Build()

    // Create a reconciler
    reconciler := &controllers.TenantReconciler{
        Client: client,
        Scheme: s,
    }

    // Reconcile
    req := ctrl.Request{
        NamespacedName: types.NamespacedName{
            Name: "test-tenant",
        },
    }
    result, err := reconciler.Reconcile(context.Background(), req)

    // Check the result
    assert.NoError(t, err)
    assert.Equal(t, ctrl.Result{}, result)

    // Check that the tenant was updated
    updatedTenant := &neurallogv1.Tenant{}
    err = client.Get(context.Background(), types.NamespacedName{Name: "test-tenant"}, updatedTenant)
    assert.NoError(t, err)
    assert.Equal(t, neurallogv1.TenantPending, updatedTenant.Status.Phase)

    // Check that a namespace was created
    namespace := &corev1.Namespace{}
    err = client.Get(context.Background(), types.NamespacedName{Name: "tenant-test-tenant"}, namespace)
    assert.NoError(t, err)
    assert.Equal(t, "test-tenant", namespace.Labels["neurallog.io/tenant"])
}
```

## End-to-End Testing

End-to-end tests verify that the operator works correctly in a real Kubernetes cluster. They are slower but provide the most realistic testing.

### Testing with kind

The `kind` tool creates a Kubernetes cluster in Docker for testing.

```go
package e2e_test

import (
    "context"
    "os"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/types"
    "k8s.io/client-go/kubernetes/scheme"
    "k8s.io/client-go/tools/clientcmd"
    "sigs.k8s.io/controller-runtime/pkg/client"

    neurallogv1 "github.com/neurallog/infra/operator/api/v1"
)

func TestE2E(t *testing.T) {
    // Skip if not running in CI
    if os.Getenv("CI") != "true" {
        t.Skip("Skipping E2E tests in local environment")
    }

    // Get the kubeconfig
    kubeconfig := os.Getenv("KUBECONFIG")
    if kubeconfig == "" {
        t.Skip("KUBECONFIG not set")
    }

    // Create a client
    config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
    assert.NoError(t, err)

    // Register the scheme
    s := scheme.Scheme
    err = neurallogv1.AddToScheme(s)
    assert.NoError(t, err)

    // Create a client
    c, err := client.New(config, client.Options{Scheme: s})
    assert.NoError(t, err)

    // Create a tenant
    tenant := &neurallogv1.Tenant{
        ObjectMeta: metav1.ObjectMeta{
            Name: "e2e-tenant",
        },
        Spec: neurallogv1.TenantSpec{
            DisplayName: "E2E Tenant",
        },
    }

    // Create the tenant
    ctx := context.Background()
    err = c.Create(ctx, tenant)
    assert.NoError(t, err)

    // Wait for the namespace to be created
    namespaceName := types.NamespacedName{Name: "tenant-e2e-tenant"}
    namespace := &corev1.Namespace{}
    for i := 0; i < 10; i++ {
        err = c.Get(ctx, namespaceName, namespace)
        if err == nil {
            break
        }
        time.Sleep(1 * time.Second)
    }
    assert.NoError(t, err)
    assert.Equal(t, "e2e-tenant", namespace.Labels["neurallog.io/tenant"])

    // Cleanup
    err = c.Delete(ctx, tenant)
    assert.NoError(t, err)
}
```

### Testing with Ginkgo and Gomega

The `ginkgo` and `gomega` packages provide a BDD-style testing framework.

```go
package e2e_test

import (
    "context"
    "os"
    "testing"
    "time"

    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/types"
    "k8s.io/client-go/kubernetes/scheme"
    "k8s.io/client-go/tools/clientcmd"
    "sigs.k8s.io/controller-runtime/pkg/client"

    neurallogv1 "github.com/neurallog/infra/operator/api/v1"
)

func TestE2E(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "E2E Suite")
}

var _ = Describe("Tenant operator", func() {
    var c client.Client
    var ctx context.Context

    BeforeEach(func() {
        // Skip if not running in CI
        if os.Getenv("CI") != "true" {
            Skip("Skipping E2E tests in local environment")
        }

        // Get the kubeconfig
        kubeconfig := os.Getenv("KUBECONFIG")
        if kubeconfig == "" {
            Skip("KUBECONFIG not set")
        }

        // Create a client
        config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
        Expect(err).NotTo(HaveOccurred())

        // Register the scheme
        s := scheme.Scheme
        err = neurallogv1.AddToScheme(s)
        Expect(err).NotTo(HaveOccurred())

        // Create a client
        c, err = client.New(config, client.Options{Scheme: s})
        Expect(err).NotTo(HaveOccurred())

        ctx = context.Background()
    })

    Context("When creating a Tenant", func() {
        var tenant *neurallogv1.Tenant

        BeforeEach(func() {
            tenant = &neurallogv1.Tenant{
                ObjectMeta: metav1.ObjectMeta{
                    Name: "e2e-tenant",
                },
                Spec: neurallogv1.TenantSpec{
                    DisplayName: "E2E Tenant",
                },
            }

            err := c.Create(ctx, tenant)
            Expect(err).NotTo(HaveOccurred())
        })

        AfterEach(func() {
            err := c.Delete(ctx, tenant)
            Expect(err).NotTo(HaveOccurred())
        })

        It("Should create a namespace", func() {
            // Wait for the namespace to be created
            namespaceName := types.NamespacedName{Name: "tenant-e2e-tenant"}
            namespace := &corev1.Namespace{}
            Eventually(func() error {
                return c.Get(ctx, namespaceName, namespace)
            }, 10*time.Second, 1*time.Second).Should(Succeed())

            // Check that the namespace has the correct labels
            Expect(namespace.Labels["neurallog.io/tenant"]).To(Equal("e2e-tenant"))
        })
    })
})
```

## Test Fixtures

Test fixtures provide consistent test data. They can be defined in separate files or in the test files themselves.

```go
package fixtures

import (
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/utils/ptr"

    neurallogv1 "github.com/neurallog/infra/operator/api/v1"
)

// NewTenant creates a new Tenant with the given name
func NewTenant(name string) *neurallogv1.Tenant {
    return &neurallogv1.Tenant{
        ObjectMeta: metav1.ObjectMeta{
            Name: name,
        },
        Spec: neurallogv1.TenantSpec{
            DisplayName: name,
            Server: neurallogv1.ServerSpec{
                Replicas: ptr.To(int32(1)),
                Image:    "neurallog/server:latest",
            },
            Redis: neurallogv1.RedisSpec{
                Replicas: ptr.To(int32(1)),
                Image:    "redis:7-alpine",
                Storage:  "1Gi",
            },
        },
    }
}

// NewTenantWithStatus creates a new Tenant with the given name and status
func NewTenantWithStatus(name string, phase neurallogv1.TenantPhase) *neurallogv1.Tenant {
    tenant := NewTenant(name)
    tenant.Status.Phase = phase
    tenant.Status.Namespace = "tenant-" + name
    return tenant
}
```

## Test Mocks

Test mocks provide simulated implementations of interfaces for testing. They can be created manually or using a mocking library like `gomock`.

```go
package mocks

import (
    "context"

    "k8s.io/apimachinery/pkg/runtime"
    "sigs.k8s.io/controller-runtime/pkg/client"
)

// MockClient is a mock implementation of client.Client
type MockClient struct {
    GetFunc    func(ctx context.Context, key client.ObjectKey, obj client.Object) error
    ListFunc   func(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error
    CreateFunc func(ctx context.Context, obj client.Object, opts ...client.CreateOption) error
    DeleteFunc func(ctx context.Context, obj client.Object, opts ...client.DeleteOption) error
    UpdateFunc func(ctx context.Context, obj client.Object, opts ...client.UpdateOption) error
    PatchFunc  func(ctx context.Context, obj client.Object, patch client.Patch, opts ...client.PatchOption) error
    StatusFunc *MockStatusWriter
}

// Get implements client.Client
func (m *MockClient) Get(ctx context.Context, key client.ObjectKey, obj client.Object) error {
    if m.GetFunc != nil {
        return m.GetFunc(ctx, key, obj)
    }
    return nil
}

// List implements client.Client
func (m *MockClient) List(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error {
    if m.ListFunc != nil {
        return m.ListFunc(ctx, list, opts...)
    }
    return nil
}

// Create implements client.Client
func (m *MockClient) Create(ctx context.Context, obj client.Object, opts ...client.CreateOption) error {
    if m.CreateFunc != nil {
        return m.CreateFunc(ctx, obj, opts...)
    }
    return nil
}

// Delete implements client.Client
func (m *MockClient) Delete(ctx context.Context, obj client.Object, opts ...client.DeleteOption) error {
    if m.DeleteFunc != nil {
        return m.DeleteFunc(ctx, obj, opts...)
    }
    return nil
}

// Update implements client.Client
func (m *MockClient) Update(ctx context.Context, obj client.Object, opts ...client.UpdateOption) error {
    if m.UpdateFunc != nil {
        return m.UpdateFunc(ctx, obj, opts...)
    }
    return nil
}

// Patch implements client.Client
func (m *MockClient) Patch(ctx context.Context, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
    if m.PatchFunc != nil {
        return m.PatchFunc(ctx, obj, patch, opts...)
    }
    return nil
}

// Status implements client.Client
func (m *MockClient) Status() client.StatusWriter {
    if m.StatusFunc != nil {
        return m.StatusFunc
    }
    return &MockStatusWriter{}
}

// Scheme implements client.Client
func (m *MockClient) Scheme() *runtime.Scheme {
    return runtime.NewScheme()
}

// RESTMapper implements client.Client
func (m *MockClient) RESTMapper() meta.RESTMapper {
    return nil
}

// MockStatusWriter is a mock implementation of client.StatusWriter
type MockStatusWriter struct {
    UpdateFunc func(ctx context.Context, obj client.Object, opts ...client.UpdateOption) error
    PatchFunc  func(ctx context.Context, obj client.Object, patch client.Patch, opts ...client.PatchOption) error
}

// Update implements client.StatusWriter
func (m *MockStatusWriter) Update(ctx context.Context, obj client.Object, opts ...client.UpdateOption) error {
    if m.UpdateFunc != nil {
        return m.UpdateFunc(ctx, obj, opts...)
    }
    return nil
}

// Patch implements client.StatusWriter
func (m *MockStatusWriter) Patch(ctx context.Context, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
    if m.PatchFunc != nil {
        return m.PatchFunc(ctx, obj, patch, opts...)
    }
    return nil
}
```

## Test Coverage

Test coverage measures how much of the code is covered by tests. It helps identify areas that need more testing.

```bash
# Run tests with coverage
go test -coverprofile=coverage.out ./...

# View coverage report
go tool cover -html=coverage.out
```

### Coverage Goals

- **Unit Tests**: Aim for 80-90% coverage
- **Integration Tests**: Aim for 70-80% coverage
- **End-to-End Tests**: Aim for key scenarios

### Coverage Report

The coverage report shows which lines of code are covered by tests:

- **Green**: Covered by tests
- **Red**: Not covered by tests
- **Gray**: Not executable (comments, imports, etc.)

## Continuous Integration

Continuous Integration (CI) runs tests automatically when code is pushed to the repository. It helps catch issues early.

```yaml
# .github/workflows/test.yml
name: Test

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    - name: Set up Go
      uses: actions/setup-go@v2
      with:
        go-version: 1.20.x
    - name: Install dependencies
      run: go mod download
    - name: Run tests
      run: go test -v -coverprofile=coverage.out ./...
    - name: Upload coverage
      uses: codecov/codecov-action@v1
      with:
        file: ./coverage.out
        flags: unittests
        fail_ci_if_error: true
```

## Best Practices

Follow these best practices for testing:

1. **Write Tests First**: Write tests before implementing features
2. **Test Edge Cases**: Test edge cases and error conditions
3. **Keep Tests Fast**: Keep tests fast for quick feedback
4. **Use Test Fixtures**: Use test fixtures for consistent testing
5. **Mock External Dependencies**: Mock external dependencies for isolation
6. **Test Public API**: Focus on testing the public API
7. **Test One Thing at a Time**: Each test should test one thing
8. **Use Descriptive Test Names**: Use descriptive test names
9. **Clean Up After Tests**: Clean up resources after tests
10. **Run Tests in CI**: Run tests in CI for every change
