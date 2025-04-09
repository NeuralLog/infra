/*
Copyright 2023 NeuralLog Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	neurallogv1 "github.com/neurallog/operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TenantReconciler reconciles a Tenant object
type TenantReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=neurallog.io,resources=tenants,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=neurallog.io,resources=tenants/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=neurallog.io,resources=tenants/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=namespaces,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=serviceaccounts,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=persistentvolumeclaims,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=networking.k8s.io,resources=networkpolicies,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=roles,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=rolebindings,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *TenantReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Reconciling Tenant", "tenant", req.NamespacedName)

	// Fetch the Tenant instance
	tenant := &neurallogv1.Tenant{}
	err := r.Get(ctx, req.NamespacedName, tenant)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Return and don't requeue
			logger.Info("Tenant resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		logger.Error(err, "Failed to get Tenant")
		return ctrl.Result{}, err
	}

	// Initialize status if it's a new tenant
	if tenant.Status.Phase == "" {
		tenant.Status.Phase = neurallogv1.TenantPending
		if err := r.Status().Update(ctx, tenant); err != nil {
			logger.Error(err, "Failed to update Tenant status")
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	// Add finalizer if it doesn't exist
	if !controllerutil.ContainsFinalizer(tenant, "neurallog.io/finalizer") {
		controllerutil.AddFinalizer(tenant, "neurallog.io/finalizer")
		if err := r.Update(ctx, tenant); err != nil {
			logger.Error(err, "Failed to add finalizer")
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	// Check if the tenant is being deleted
	if !tenant.ObjectMeta.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, tenant)
	}

	// Update status to Provisioning if it's still Pending
	if tenant.Status.Phase == neurallogv1.TenantPending {
		tenant.Status.Phase = neurallogv1.TenantProvisioning
		if err := r.Status().Update(ctx, tenant); err != nil {
			logger.Error(err, "Failed to update Tenant status")
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	// Create or update the namespace
	namespace, err := r.reconcileNamespace(ctx, tenant)
	if err != nil {
		logger.Error(err, "Failed to reconcile namespace")
		return ctrl.Result{}, err
	}

	// Update the namespace in the status if it's not set
	if tenant.Status.Namespace == "" {
		tenant.Status.Namespace = namespace.Name
		if err := r.Status().Update(ctx, tenant); err != nil {
			logger.Error(err, "Failed to update Tenant status with namespace")
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	// Reconcile Redis resources
	if err := r.reconcileRedis(ctx, tenant); err != nil {
		logger.Error(err, "Failed to reconcile Redis")
		return ctrl.Result{}, err
	}

	// Reconcile Server resources
	if err := r.reconcileServer(ctx, tenant); err != nil {
		logger.Error(err, "Failed to reconcile Server")
		return ctrl.Result{}, err
	}

	// Reconcile Network Policies
	if err := r.reconcileNetworkPolicies(ctx, tenant); err != nil {
		logger.Error(err, "Failed to reconcile Network Policies")
		return ctrl.Result{}, err
	}

	// Update status to Running if everything is provisioned
	if tenant.Status.Phase == neurallogv1.TenantProvisioning {
		tenant.Status.Phase = neurallogv1.TenantRunning
		if err := r.Status().Update(ctx, tenant); err != nil {
			logger.Error(err, "Failed to update Tenant status")
			return ctrl.Result{}, err
		}
	}

	// Requeue to check status periodically
	return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
}

// reconcileDelete handles the deletion of a Tenant
func (r *TenantReconciler) reconcileDelete(ctx context.Context, tenant *neurallogv1.Tenant) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Reconciling Tenant deletion", "tenant", tenant.Name)

	// Update status to Terminating if it's not already
	if tenant.Status.Phase != neurallogv1.TenantTerminating {
		tenant.Status.Phase = neurallogv1.TenantTerminating
		if err := r.Status().Update(ctx, tenant); err != nil {
			logger.Error(err, "Failed to update Tenant status")
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	// Delete the namespace if it exists
	if tenant.Status.Namespace != "" {
		namespace := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: tenant.Status.Namespace,
			},
		}
		err := r.Delete(ctx, namespace)
		if err != nil && !errors.IsNotFound(err) {
			logger.Error(err, "Failed to delete namespace")
			return ctrl.Result{}, err
		}
		logger.Info("Deleted namespace", "namespace", tenant.Status.Namespace)
	}

	// Remove finalizer
	controllerutil.RemoveFinalizer(tenant, "neurallog.io/finalizer")
	if err := r.Update(ctx, tenant); err != nil {
		logger.Error(err, "Failed to remove finalizer")
		return ctrl.Result{}, err
	}

	logger.Info("Successfully deleted Tenant", "tenant", tenant.Name)
	return ctrl.Result{}, nil
}

// reconcileNamespace creates or updates the namespace for the tenant
func (r *TenantReconciler) reconcileNamespace(ctx context.Context, tenant *neurallogv1.Tenant) (*corev1.Namespace, error) {
	logger := log.FromContext(ctx)

	// Generate namespace name if not set in status
	namespaceName := tenant.Status.Namespace
	if namespaceName == "" {
		namespaceName = fmt.Sprintf("tenant-%s", tenant.Name)
	}

	// Check if namespace exists
	namespace := &corev1.Namespace{}
	err := r.Get(ctx, client.ObjectKey{Name: namespaceName}, namespace)
	if err != nil {
		if errors.IsNotFound(err) {
			// Create namespace
			namespace = &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespaceName,
					Labels: map[string]string{
						"neurallog.io/tenant": tenant.Name,
						"neurallog.io/managed-by": "tenant-operator",
					},
				},
			}
			if err := controllerutil.SetOwnerReference(tenant, namespace, r.Scheme); err != nil {
				logger.Error(err, "Failed to set owner reference on namespace")
				return nil, err
			}
			if err := r.Create(ctx, namespace); err != nil {
				logger.Error(err, "Failed to create namespace")
				return nil, err
			}
			logger.Info("Created namespace", "namespace", namespaceName)
		} else {
			logger.Error(err, "Failed to get namespace")
			return nil, err
		}
	} else {
		// Namespace exists, ensure it has the correct labels
		updated := false
		if namespace.Labels == nil {
			namespace.Labels = make(map[string]string)
			updated = true
		}
		if namespace.Labels["neurallog.io/tenant"] != tenant.Name {
			namespace.Labels["neurallog.io/tenant"] = tenant.Name
			updated = true
		}
		if namespace.Labels["neurallog.io/managed-by"] != "tenant-operator" {
			namespace.Labels["neurallog.io/managed-by"] = "tenant-operator"
			updated = true
		}
		if updated {
			if err := r.Update(ctx, namespace); err != nil {
				logger.Error(err, "Failed to update namespace")
				return nil, err
			}
			logger.Info("Updated namespace", "namespace", namespaceName)
		}
	}

	return namespace, nil
}

// reconcileRedis creates or updates Redis resources for the tenant
func (r *TenantReconciler) reconcileRedis(ctx context.Context, tenant *neurallogv1.Tenant) error {
	logger := log.FromContext(ctx)
	logger.Info("Reconciling Redis resources", "tenant", tenant.Name)

	// TODO: Implement Redis reconciliation
	// This would include:
	// - Creating/updating Redis ConfigMap
	// - Creating/updating Redis StatefulSet
	// - Creating/updating Redis Service
	// - Updating Redis status in tenant.Status.RedisStatus

	return nil
}

// reconcileServer creates or updates Server resources for the tenant
func (r *TenantReconciler) reconcileServer(ctx context.Context, tenant *neurallogv1.Tenant) error {
	logger := log.FromContext(ctx)
	logger.Info("Reconciling Server resources", "tenant", tenant.Name)

	// TODO: Implement Server reconciliation
	// This would include:
	// - Creating/updating Server ConfigMap
	// - Creating/updating Server Deployment
	// - Creating/updating Server Service
	// - Updating Server status in tenant.Status.ServerStatus

	return nil
}

// reconcileNetworkPolicies creates or updates Network Policies for the tenant
func (r *TenantReconciler) reconcileNetworkPolicies(ctx context.Context, tenant *neurallogv1.Tenant) error {
	logger := log.FromContext(ctx)
	logger.Info("Reconciling Network Policies", "tenant", tenant.Name)

	// TODO: Implement Network Policy reconciliation
	// This would include:
	// - Creating/updating Network Policies based on tenant.Spec.NetworkPolicy

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TenantReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&neurallogv1.Tenant{}).
		Complete(r)
}
