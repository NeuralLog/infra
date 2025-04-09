# Controller Development

This document provides guidance on developing and extending the NeuralLog Tenant Operator controllers.

## Table of Contents

- [Controller Architecture](#controller-architecture)
- [Reconciliation Loop](#reconciliation-loop)
- [Error Handling](#error-handling)
- [Status Updates](#status-updates)
- [Finalizers](#finalizers)
- [Owner References](#owner-references)
- [Event Recording](#event-recording)
- [Caching and Indexing](#caching-and-indexing)
- [Rate Limiting and Backoff](#rate-limiting-and-backoff)
- [Testing Controllers](#testing-controllers)
- [Debugging Controllers](#debugging-controllers)
- [Best Practices](#best-practices)

## Controller Architecture

The NeuralLog Tenant Operator uses the controller-runtime framework to implement controllers. The main components are:

1. **Manager**: Manages controllers, webhooks, and shared caches
2. **Controller**: Watches for changes to resources and triggers reconciliation
3. **Reconciler**: Implements the reconciliation logic
4. **Client**: Provides access to the Kubernetes API

```go
func main() {
    // Create a manager
    mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
        Scheme:             scheme,
        MetricsBindAddress: metricsAddr,
        Port:               9443,
        LeaderElection:     enableLeaderElection,
        LeaderElectionID:   "neurallog-operator.neurallog.io",
    })
    if err != nil {
        setupLog.Error(err, "unable to start manager")
        os.Exit(1)
    }

    // Create a controller
    if err = (&controllers.TenantReconciler{
        Client: mgr.GetClient(),
        Scheme: mgr.GetScheme(),
    }).SetupWithManager(mgr); err != nil {
        setupLog.Error(err, "unable to create controller", "controller", "Tenant")
        os.Exit(1)
    }

    // Start the manager
    if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
        setupLog.Error(err, "problem running manager")
        os.Exit(1)
    }
}
```

## Reconciliation Loop

The reconciliation loop is the core of the controller. It is responsible for ensuring that the actual state of the system matches the desired state specified in the Tenant resource.

```go
// Reconcile is part of the main kubernetes reconciliation loop
func (r *TenantReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    logger := log.FromContext(ctx)

    // Fetch the Tenant instance
    tenant := &neurallogv1.Tenant{}
    if err := r.Get(ctx, req.NamespacedName, tenant); err != nil {
        if errors.IsNotFound(err) {
            // Object not found, return
            return ctrl.Result{}, nil
        }
        // Error reading the object
        return ctrl.Result{}, err
    }

    // Initialize status if needed
    if tenant.Status.Phase == "" {
        tenant.Status.Phase = neurallogv1.TenantPending
        if err := r.Status().Update(ctx, tenant); err != nil {
            logger.Error(err, "Failed to update Tenant status")
            return ctrl.Result{}, err
        }
    }

    // Add finalizer if needed
    if !controllerutil.ContainsFinalizer(tenant, tenantFinalizer) {
        controllerutil.AddFinalizer(tenant, tenantFinalizer)
        if err := r.Update(ctx, tenant); err != nil {
            logger.Error(err, "Failed to add finalizer")
            return ctrl.Result{}, err
        }
    }

    // Check if the Tenant is being deleted
    if !tenant.ObjectMeta.DeletionTimestamp.IsZero() {
        return r.reconcileDelete(ctx, tenant)
    }

    // Reconcile the Tenant
    return r.reconcileNormal(ctx, tenant)
}
```

### Reconcile Normal

The `reconcileNormal` function handles the normal reconciliation flow:

```go
func (r *TenantReconciler) reconcileNormal(ctx context.Context, tenant *neurallogv1.Tenant) (ctrl.Result, error) {
    logger := log.FromContext(ctx)

    // Reconcile namespace
    namespace, err := r.reconcileNamespace(ctx, tenant)
    if err != nil {
        logger.Error(err, "Failed to reconcile namespace")
        return ctrl.Result{}, err
    }

    // Update tenant status with namespace
    tenant.Status.Namespace = namespace.Name

    // Reconcile Redis
    redis, err := r.reconcileRedis(ctx, tenant)
    if err != nil {
        logger.Error(err, "Failed to reconcile Redis")
        return ctrl.Result{}, err
    }

    // Update tenant status with Redis status
    tenant.Status.RedisStatus = r.getRedisStatus(redis)

    // Reconcile server
    server, err := r.reconcileServer(ctx, tenant)
    if err != nil {
        logger.Error(err, "Failed to reconcile server")
        return ctrl.Result{}, err
    }

    // Update tenant status with server status
    tenant.Status.ServerStatus = r.getServerStatus(server)

    // Reconcile network policies
    if err := r.reconcileNetworkPolicies(ctx, tenant); err != nil {
        logger.Error(err, "Failed to reconcile network policies")
        return ctrl.Result{}, err
    }

    // Update tenant phase
    tenant.Status.Phase = r.getTenantPhase(tenant)

    // Update tenant status
    if err := r.Status().Update(ctx, tenant); err != nil {
        logger.Error(err, "Failed to update Tenant status")
        return ctrl.Result{}, err
    }

    // Requeue for periodic reconciliation
    return ctrl.Result{RequeueAfter: 5 * time.Minute}, nil
}
```

### Reconcile Delete

The `reconcileDelete` function handles the deletion flow:

```go
func (r *TenantReconciler) reconcileDelete(ctx context.Context, tenant *neurallogv1.Tenant) (ctrl.Result, error) {
    logger := log.FromContext(ctx)

    // Delete namespace
    namespace := &corev1.Namespace{
        ObjectMeta: metav1.ObjectMeta{
            Name: tenant.Status.Namespace,
        },
    }
    if err := r.Delete(ctx, namespace); err != nil {
        if !errors.IsNotFound(err) {
            logger.Error(err, "Failed to delete namespace")
            return ctrl.Result{}, err
        }
    }

    // Remove finalizer
    controllerutil.RemoveFinalizer(tenant, tenantFinalizer)
    if err := r.Update(ctx, tenant); err != nil {
        logger.Error(err, "Failed to remove finalizer")
        return ctrl.Result{}, err
    }

    return ctrl.Result{}, nil
}
```

## Error Handling

Proper error handling is crucial for controllers. Follow these guidelines:

1. **Return Errors**: Return errors to trigger requeuing
2. **Log Errors**: Log errors with context
3. **Categorize Errors**: Distinguish between transient and permanent errors
4. **Handle Not Found Errors**: Handle not found errors gracefully
5. **Requeue with Backoff**: Requeue with backoff for transient errors

```go
// Example error handling
if err := r.Get(ctx, req.NamespacedName, tenant); err != nil {
    if errors.IsNotFound(err) {
        // Object not found, return without error
        return ctrl.Result{}, nil
    }
    // Error reading the object, requeue
    logger.Error(err, "Failed to get Tenant")
    return ctrl.Result{}, err
}

// Example transient error handling
if err := r.Create(ctx, namespace); err != nil {
    if errors.IsAlreadyExists(err) {
        // Namespace already exists, continue
        logger.Info("Namespace already exists", "namespace", namespace.Name)
    } else {
        // Error creating namespace, requeue
        logger.Error(err, "Failed to create namespace")
        return ctrl.Result{}, err
    }
}

// Example permanent error handling
if err := validateTenant(tenant); err != nil {
    // Tenant is invalid, update status and don't requeue
    tenant.Status.Phase = neurallogv1.TenantFailed
    tenant.Status.Message = err.Error()
    if updateErr := r.Status().Update(ctx, tenant); updateErr != nil {
        logger.Error(updateErr, "Failed to update Tenant status")
        return ctrl.Result{}, updateErr
    }
    // Don't requeue for permanent errors
    return ctrl.Result{}, nil
}
```

## Status Updates

Status updates inform users about the state of their resources. Follow these guidelines:

1. **Update Status Frequently**: Update status after each reconciliation step
2. **Use Conditions**: Use conditions for complex status
3. **Include Timestamps**: Include timestamps for status changes
4. **Provide Detailed Messages**: Provide detailed messages for status changes
5. **Handle Status Update Errors**: Handle status update errors gracefully

```go
// Example status update
tenant.Status.Phase = neurallogv1.TenantRunning
tenant.Status.Message = "Tenant is running"
tenant.Status.LastTransitionTime = metav1.Now()

// Example condition update
condition := metav1.Condition{
    Type:               "Ready",
    Status:             metav1.ConditionTrue,
    LastTransitionTime: metav1.Now(),
    Reason:             "TenantReady",
    Message:            "Tenant is ready",
}
meta.SetStatusCondition(&tenant.Status.Conditions, condition)

// Example status update error handling
if err := r.Status().Update(ctx, tenant); err != nil {
    logger.Error(err, "Failed to update Tenant status")
    return ctrl.Result{}, err
}
```

## Finalizers

Finalizers ensure that resources are cleaned up before deletion. Follow these guidelines:

1. **Add Finalizers Early**: Add finalizers during the first reconciliation
2. **Remove Finalizers Last**: Remove finalizers after all cleanup is done
3. **Handle Finalizer Errors**: Handle finalizer errors gracefully
4. **Use Descriptive Finalizer Names**: Use descriptive finalizer names

```go
// Example finalizer management
const tenantFinalizer = "tenant.neurallog.io/finalizer"

// Add finalizer
if !controllerutil.ContainsFinalizer(tenant, tenantFinalizer) {
    controllerutil.AddFinalizer(tenant, tenantFinalizer)
    if err := r.Update(ctx, tenant); err != nil {
        logger.Error(err, "Failed to add finalizer")
        return ctrl.Result{}, err
    }
}

// Check if the Tenant is being deleted
if !tenant.ObjectMeta.DeletionTimestamp.IsZero() {
    // Tenant is being deleted, clean up resources
    
    // Remove finalizer
    controllerutil.RemoveFinalizer(tenant, tenantFinalizer)
    if err := r.Update(ctx, tenant); err != nil {
        logger.Error(err, "Failed to remove finalizer")
        return ctrl.Result{}, err
    }
}
```

## Owner References

Owner references establish parent-child relationships between resources. Follow these guidelines:

1. **Set Owner References**: Set owner references for all created resources
2. **Use Controller References**: Use controller references for resources that should be deleted with the owner
3. **Handle Owner Reference Errors**: Handle owner reference errors gracefully

```go
// Example owner reference
if err := controllerutil.SetControllerReference(tenant, namespace, r.Scheme); err != nil {
    logger.Error(err, "Failed to set owner reference")
    return nil, err
}
```

## Event Recording

Events provide visibility into controller actions. Follow these guidelines:

1. **Record Significant Events**: Record events for significant actions
2. **Use Appropriate Event Types**: Use Normal for expected events, Warning for issues
3. **Provide Detailed Messages**: Provide detailed messages for events
4. **Include Resource References**: Include resource references in events

```go
// Example event recording
r.Recorder.Event(tenant, corev1.EventTypeNormal, "Created", fmt.Sprintf("Created namespace %s", namespace.Name))
r.Recorder.Event(tenant, corev1.EventTypeWarning, "Failed", fmt.Sprintf("Failed to create namespace: %s", err.Error()))
```

## Caching and Indexing

Caching and indexing improve controller performance. Follow these guidelines:

1. **Use Client Cache**: Use the client cache for reads
2. **Use Direct Client for Writes**: Use the direct client for writes
3. **Add Indexes for Lookups**: Add indexes for frequent lookups
4. **Handle Cache Errors**: Handle cache errors gracefully

```go
// Example cache setup
func (r *TenantReconciler) SetupWithManager(mgr ctrl.Manager) error {
    // Add index for tenant namespace
    if err := mgr.GetFieldIndexer().IndexField(context.Background(), &corev1.Namespace{}, "tenant", func(obj client.Object) []string {
        namespace := obj.(*corev1.Namespace)
        tenant := namespace.Labels["neurallog.io/tenant"]
        if tenant == "" {
            return nil
        }
        return []string{tenant}
    }); err != nil {
        return err
    }

    return ctrl.NewControllerManagedBy(mgr).
        For(&neurallogv1.Tenant{}).
        Owns(&corev1.Namespace{}).
        Complete(r)
}

// Example index usage
var namespaces corev1.NamespaceList
if err := r.List(ctx, &namespaces, client.MatchingFields{"tenant": tenant.Name}); err != nil {
    logger.Error(err, "Failed to list namespaces")
    return ctrl.Result{}, err
}
```

## Rate Limiting and Backoff

Rate limiting and backoff prevent overloading the API server. Follow these guidelines:

1. **Use Rate Limiters**: Use rate limiters for reconciliation
2. **Use Exponential Backoff**: Use exponential backoff for retries
3. **Set Appropriate Requeue Times**: Set appropriate requeue times for periodic reconciliation
4. **Handle Rate Limiting Errors**: Handle rate limiting errors gracefully

```go
// Example rate limiter setup
func (r *TenantReconciler) SetupWithManager(mgr ctrl.Manager) error {
    return ctrl.NewControllerManagedBy(mgr).
        For(&neurallogv1.Tenant{}).
        WithOptions(controller.Options{
            RateLimiter: workqueue.NewMaxOfRateLimiter(
                workqueue.NewItemExponentialFailureRateLimiter(5*time.Millisecond, 1000*time.Second),
                &workqueue.BucketRateLimiter{Limiter: rate.NewLimiter(rate.Limit(10), 100)},
            ),
        }).
        Complete(r)
}

// Example requeue with backoff
return ctrl.Result{RequeueAfter: time.Duration(2^attempt) * time.Second}, nil
```

## Testing Controllers

Testing controllers is crucial for reliability. Follow these guidelines:

1. **Write Unit Tests**: Test individual functions
2. **Write Integration Tests**: Test controller with fake client
3. **Write End-to-End Tests**: Test controller in a real cluster
4. **Use Test Fixtures**: Use test fixtures for consistent testing
5. **Mock External Dependencies**: Mock external dependencies for testing

```go
// Example controller test
func TestTenantReconciler_Reconcile(t *testing.T) {
    // Create a fake client
    tenant := &neurallogv1.Tenant{
        ObjectMeta: metav1.ObjectMeta{
            Name: "test-tenant",
        },
        Spec: neurallogv1.TenantSpec{
            DisplayName: "Test Tenant",
        },
    }
    
    // Create a reconciler with the fake client
    r := &TenantReconciler{
        Client: fake.NewClientBuilder().WithObjects(tenant).Build(),
        Scheme: scheme.Scheme,
    }
    
    // Reconcile
    _, err := r.Reconcile(context.Background(), ctrl.Request{
        NamespacedName: types.NamespacedName{
            Name: "test-tenant",
        },
    })
    
    // Check for errors
    if err != nil {
        t.Errorf("Reconcile() error = %v", err)
    }
    
    // Check that the tenant was updated
    updatedTenant := &neurallogv1.Tenant{}
    if err := r.Get(context.Background(), types.NamespacedName{Name: "test-tenant"}, updatedTenant); err != nil {
        t.Errorf("Failed to get updated tenant: %v", err)
    }
    
    // Check that the tenant status was updated
    if updatedTenant.Status.Phase != neurallogv1.TenantPending {
        t.Errorf("Expected tenant phase to be Pending, got %s", updatedTenant.Status.Phase)
    }
}
```

## Debugging Controllers

Debugging controllers can be challenging. Follow these guidelines:

1. **Use Verbose Logging**: Use verbose logging for debugging
2. **Use Debugger**: Use a debugger for step-by-step debugging
3. **Inspect Resources**: Inspect resources for state
4. **Check Events**: Check events for controller actions
5. **Use Metrics**: Use metrics for performance debugging

```go
// Example verbose logging
logger.V(1).Info("Reconciling Tenant", "tenant", tenant.Name)
logger.V(2).Info("Creating namespace", "namespace", namespace.Name)
logger.V(3).Info("Namespace spec", "spec", namespace.Spec)

// Example resource inspection
kubectl get tenant test-tenant -o yaml
kubectl get namespace tenant-test-tenant -o yaml
kubectl get events --field-selector involvedObject.name=test-tenant
```

## Best Practices

Follow these best practices for controller development:

1. **Keep Reconciliation Idempotent**: Reconciliation should be idempotent
2. **Minimize API Calls**: Minimize API calls for performance
3. **Handle Edge Cases**: Handle edge cases gracefully
4. **Use Status for Visibility**: Use status for visibility
5. **Record Events for Significant Actions**: Record events for significant actions
6. **Use Owner References**: Use owner references for resource management
7. **Use Finalizers for Cleanup**: Use finalizers for cleanup
8. **Test Thoroughly**: Test thoroughly for reliability
9. **Document Controller Behavior**: Document controller behavior
10. **Follow Kubernetes Patterns**: Follow Kubernetes patterns for consistency

```go
// Example idempotent reconciliation
if err := r.Get(ctx, types.NamespacedName{Name: namespaceName}, namespace); err != nil {
    if errors.IsNotFound(err) {
        // Namespace doesn't exist, create it
        namespace = &corev1.Namespace{
            ObjectMeta: metav1.ObjectMeta{
                Name: namespaceName,
                Labels: map[string]string{
                    "neurallog.io/tenant": tenant.Name,
                },
            },
        }
        if err := controllerutil.SetControllerReference(tenant, namespace, r.Scheme); err != nil {
            logger.Error(err, "Failed to set owner reference")
            return nil, err
        }
        if err := r.Create(ctx, namespace); err != nil {
            logger.Error(err, "Failed to create namespace")
            return nil, err
        }
        logger.Info("Created namespace", "namespace", namespace.Name)
    } else {
        // Error getting namespace
        logger.Error(err, "Failed to get namespace")
        return nil, err
    }
} else {
    // Namespace exists, update it if needed
    if namespace.Labels["neurallog.io/tenant"] != tenant.Name {
        namespace.Labels["neurallog.io/tenant"] = tenant.Name
        if err := r.Update(ctx, namespace); err != nil {
            logger.Error(err, "Failed to update namespace")
            return nil, err
        }
        logger.Info("Updated namespace", "namespace", namespace.Name)
    }
}
```
