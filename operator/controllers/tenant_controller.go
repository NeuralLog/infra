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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	neurallogv1 "github.com/neurallog/operator/api/v1"
	appsv1 "k8s.io/api/apps/v1"
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

	// Reconcile Registry resources
	if err := r.reconcileRegistry(ctx, tenant); err != nil {
		logger.Error(err, "Failed to reconcile Registry")
		return ctrl.Result{}, err
	}

	// Reconcile Network Policies
	if err := r.reconcileNetworkPolicies(ctx, tenant); err != nil {
		logger.Error(err, "Failed to reconcile Network Policies")
		return ctrl.Result{}, err
	}

	// Reconcile Auth Service integration
	if err := r.reconcileAuthService(ctx, tenant); err != nil {
		logger.Error(err, "Failed to reconcile Auth Service integration")
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

	// Delete the tenant from the Auth service
	if err := r.deleteTenantFromAuthService(ctx, tenant.Name); err != nil {
		logger.Error(err, "Failed to delete tenant from Auth service")
		// Don't return an error here, as we still want to remove the finalizer
		// The tenant will be garbage collected by the Auth service
	} else {
		logger.Info("Deleted tenant from Auth service", "tenant", tenant.Name)
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
						"neurallog.io/tenant":     tenant.Name,
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

	if tenant.Status.Namespace == "" {
		logger.Info("Namespace not yet created, skipping Redis reconciliation")
		return nil
	}

	// Create or update Redis ConfigMap
	configMap, err := r.reconcileRedisConfigMap(ctx, tenant)
	if err != nil {
		logger.Error(err, "Failed to reconcile Redis ConfigMap")
		return err
	}

	// Create or update Redis Service
	service, err := r.reconcileRedisService(ctx, tenant)
	if err != nil {
		logger.Error(err, "Failed to reconcile Redis Service")
		return err
	}

	// Create or update Redis StatefulSet
	statefulSet, err := r.reconcileRedisStatefulSet(ctx, tenant, configMap)
	if err != nil {
		logger.Error(err, "Failed to reconcile Redis StatefulSet")
		return err
	}

	// Update Redis status in tenant.Status.RedisStatus
	if tenant.Status.RedisStatus == nil {
		tenant.Status.RedisStatus = &neurallogv1.ComponentStatus{}
	}

	// Update status based on StatefulSet
	tenant.Status.RedisStatus.TotalReplicas = statefulSet.Status.Replicas
	tenant.Status.RedisStatus.ReadyReplicas = statefulSet.Status.ReadyReplicas

	// Determine phase based on replicas
	if statefulSet.Status.ReadyReplicas == 0 {
		tenant.Status.RedisStatus.Phase = neurallogv1.ComponentProvisioning
		tenant.Status.RedisStatus.Message = "Redis is being provisioned"
	} else if statefulSet.Status.ReadyReplicas < statefulSet.Status.Replicas {
		tenant.Status.RedisStatus.Phase = neurallogv1.ComponentDegraded
		tenant.Status.RedisStatus.Message = fmt.Sprintf("Redis is degraded: %d/%d replicas ready", statefulSet.Status.ReadyReplicas, statefulSet.Status.Replicas)
	} else {
		tenant.Status.RedisStatus.Phase = neurallogv1.ComponentRunning
		tenant.Status.RedisStatus.Message = "Redis is running"
	}

	// Update tenant status
	if err := r.Status().Update(ctx, tenant); err != nil {
		logger.Error(err, "Failed to update tenant status with Redis status")
		return err
	}

	return nil
}

// reconcileServer creates or updates Server resources for the tenant
func (r *TenantReconciler) reconcileServer(ctx context.Context, tenant *neurallogv1.Tenant) error {
	logger := log.FromContext(ctx)
	logger.Info("Reconciling Server resources", "tenant", tenant.Name)

	if tenant.Status.Namespace == "" {
		logger.Info("Namespace not yet created, skipping Server reconciliation")
		return nil
	}

	// Create or update Server Service
	service, err := r.reconcileServerService(ctx, tenant)
	if err != nil {
		logger.Error(err, "Failed to reconcile Server Service")
		return err
	}

	// Create or update Server Deployment
	deployment, err := r.reconcileServerDeployment(ctx, tenant)
	if err != nil {
		logger.Error(err, "Failed to reconcile Server Deployment")
		return err
	}

	// Update Server status in tenant.Status.ServerStatus
	if tenant.Status.ServerStatus == nil {
		tenant.Status.ServerStatus = &neurallogv1.ComponentStatus{}
	}

	// Update status based on Deployment
	tenant.Status.ServerStatus.TotalReplicas = deployment.Status.Replicas
	tenant.Status.ServerStatus.ReadyReplicas = deployment.Status.ReadyReplicas

	// Determine phase based on replicas
	if deployment.Status.ReadyReplicas == 0 {
		tenant.Status.ServerStatus.Phase = neurallogv1.ComponentProvisioning
		tenant.Status.ServerStatus.Message = "Server is being provisioned"
	} else if deployment.Status.ReadyReplicas < deployment.Status.Replicas {
		tenant.Status.ServerStatus.Phase = neurallogv1.ComponentDegraded
		tenant.Status.ServerStatus.Message = fmt.Sprintf("Server is degraded: %d/%d replicas ready", deployment.Status.ReadyReplicas, deployment.Status.Replicas)
	} else {
		tenant.Status.ServerStatus.Phase = neurallogv1.ComponentRunning
		tenant.Status.ServerStatus.Message = "Server is running"
	}

	// Update tenant status
	if err := r.Status().Update(ctx, tenant); err != nil {
		logger.Error(err, "Failed to update tenant status with Server status")
		return err
	}

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

// reconcileAuthService integrates with the Auth service to manage tenant authentication
func (r *TenantReconciler) reconcileAuthService(ctx context.Context, tenant *neurallogv1.Tenant) error {
	logger := log.FromContext(ctx)
	logger.Info("Reconciling Auth Service integration", "tenant", tenant.Name)

	// Skip if the tenant is being deleted
	if !tenant.ObjectMeta.DeletionTimestamp.IsZero() {
		return nil
	}

	// Check if the tenant exists in the Auth service
	exists, err := r.tenantExistsInAuthService(ctx, tenant.Name)
	if err != nil {
		logger.Error(err, "Failed to check if tenant exists in Auth service")
		return err
	}

	// If the tenant doesn't exist in the Auth service, create it
	if !exists {
		if err := r.createTenantInAuthService(ctx, tenant.Name); err != nil {
			logger.Error(err, "Failed to create tenant in Auth service")
			return err
		}
		logger.Info("Created tenant in Auth service", "tenant", tenant.Name)
	}

	return nil
}

// tenantExistsInAuthService checks if a tenant exists in the Auth service
func (r *TenantReconciler) tenantExistsInAuthService(ctx context.Context, tenantId string) (bool, error) {
	logger := log.FromContext(ctx)

	// Make a request to the Auth service to list tenants
	resp, err := http.Get("http://auth:3000/api/tenants")
	if err != nil {
		logger.Error(err, "Failed to connect to Auth service")
		return false, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err, "Failed to read response from Auth service")
		return false, err
	}

	// Parse the response
	var response struct {
		Status  string   `json:"status"`
		Tenants []string `json:"tenants"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		logger.Error(err, "Failed to parse response from Auth service")
		return false, err
	}

	// Check if the tenant exists
	for _, t := range response.Tenants {
		if t == tenantId {
			return true, nil
		}
	}

	return false, nil
}

// createTenantInAuthService creates a tenant in the Auth service
func (r *TenantReconciler) createTenantInAuthService(ctx context.Context, tenantId string) error {
	logger := log.FromContext(ctx)

	// Create the request body
	reqBody, err := json.Marshal(map[string]string{
		"tenantId":    tenantId,
		"adminUserId": "system", // Default admin user
	})
	if err != nil {
		logger.Error(err, "Failed to marshal request body")
		return err
	}

	// Make a request to the Auth service to create a tenant
	resp, err := http.Post("http://auth:3000/api/tenants", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		logger.Error(err, "Failed to connect to Auth service")
		return err
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusCreated {
		// Read the response body for error details
		body, _ := io.ReadAll(resp.Body)
		logger.Error(nil, "Failed to create tenant in Auth service", "statusCode", resp.StatusCode, "response", string(body))
		return fmt.Errorf("failed to create tenant in Auth service: %d", resp.StatusCode)
	}

	return nil
}

// deleteTenantFromAuthService deletes a tenant from the Auth service
func (r *TenantReconciler) deleteTenantFromAuthService(ctx context.Context, tenantId string) error {
	logger := log.FromContext(ctx)

	// Create a DELETE request
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("http://auth:3000/api/tenants/%s", tenantId), nil)
	if err != nil {
		logger.Error(err, "Failed to create DELETE request")
		return err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(err, "Failed to connect to Auth service")
		return err
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		// Read the response body for error details
		body, _ := io.ReadAll(resp.Body)
		logger.Error(nil, "Failed to delete tenant from Auth service", "statusCode", resp.StatusCode, "response", string(body))
		return fmt.Errorf("failed to delete tenant from Auth service: %d", resp.StatusCode)
	}

	return nil
}

// reconcileRedisConfigMap creates or updates the Redis ConfigMap for the tenant
func (r *TenantReconciler) reconcileRedisConfigMap(ctx context.Context, tenant *neurallogv1.Tenant) (*corev1.ConfigMap, error) {
	logger := log.FromContext(ctx)

	// Define ConfigMap name
	configMapName := fmt.Sprintf("%s-redis-config", tenant.Name)

	// Create Redis configuration
	redisConfig := `# Redis configuration for NeuralLog tenant
port 6379
bind 0.0.0.0
protected-mode yes
daemonize no

# Memory management
`

	// Add memory limit if specified
	maxMemory := "256mb"
	if tenant.Spec.Redis != nil && tenant.Spec.Redis.Resources != nil && tenant.Spec.Redis.Resources.Memory != nil && tenant.Spec.Redis.Resources.Memory.Limit != "" {
		maxMemory = tenant.Spec.Redis.Resources.Memory.Limit
	}
	redisConfig += fmt.Sprintf("maxmemory %s\n", maxMemory)
	redisConfig += "maxmemory-policy allkeys-lru\n\n"

	// Add persistence configuration
	redisConfig += `# Persistence
appendonly yes
appendfsync everysec

# Logging
loglevel notice
logfile ""
`

	// Add custom configuration if specified
	if tenant.Spec.Redis != nil && tenant.Spec.Redis.Config != nil {
		for key, value := range tenant.Spec.Redis.Config {
			redisConfig += fmt.Sprintf("%s %s\n", key, value)
		}
	}

	// Define ConfigMap
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName,
			Namespace: tenant.Status.Namespace,
			Labels: map[string]string{
				"app":                     "redis",
				"neurallog.io/tenant":     tenant.Name,
				"neurallog.io/component":  "redis",
				"neurallog.io/managed-by": "tenant-operator",
			},
		},
		Data: map[string]string{
			"redis.conf": redisConfig,
		},
	}

	// Set owner reference
	if err := controllerutil.SetControllerReference(tenant, configMap, r.Scheme); err != nil {
		logger.Error(err, "Failed to set owner reference on Redis ConfigMap")
		return nil, err
	}

	// Create or update ConfigMap
	existingConfigMap := &corev1.ConfigMap{}
	err := r.Get(ctx, client.ObjectKey{Name: configMapName, Namespace: tenant.Status.Namespace}, existingConfigMap)
	if err != nil {
		if errors.IsNotFound(err) {
			// Create ConfigMap
			if err := r.Create(ctx, configMap); err != nil {
				logger.Error(err, "Failed to create Redis ConfigMap")
				return nil, err
			}
			logger.Info("Created Redis ConfigMap", "configMap", configMapName)
			return configMap, nil
		}
		logger.Error(err, "Failed to get Redis ConfigMap")
		return nil, err
	}

	// Update ConfigMap if needed
	if existingConfigMap.Data["redis.conf"] != configMap.Data["redis.conf"] {
		existingConfigMap.Data = configMap.Data
		if err := r.Update(ctx, existingConfigMap); err != nil {
			logger.Error(err, "Failed to update Redis ConfigMap")
			return nil, err
		}
		logger.Info("Updated Redis ConfigMap", "configMap", configMapName)
	}

	return existingConfigMap, nil
}

// reconcileRedisService creates or updates the Redis Service for the tenant
func (r *TenantReconciler) reconcileRedisService(ctx context.Context, tenant *neurallogv1.Tenant) (*corev1.Service, error) {
	logger := log.FromContext(ctx)

	// Define Service name
	serviceName := fmt.Sprintf("%s-redis", tenant.Name)

	// Define Service
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: tenant.Status.Namespace,
			Labels: map[string]string{
				"app":                     "redis",
				"neurallog.io/tenant":     tenant.Name,
				"neurallog.io/component":  "redis",
				"neurallog.io/managed-by": "tenant-operator",
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app":                 "redis",
				"neurallog.io/tenant": tenant.Name,
			},
			Ports: []corev1.ServicePort{
				{
					Port:       6379,
					TargetPort: intstr.FromString("redis"),
					Name:       "redis",
				},
			},
			ClusterIP: "None", // Headless service for StatefulSet
		},
	}

	// Set owner reference
	if err := controllerutil.SetControllerReference(tenant, service, r.Scheme); err != nil {
		logger.Error(err, "Failed to set owner reference on Redis Service")
		return nil, err
	}

	// Create or update Service
	existingService := &corev1.Service{}
	err := r.Get(ctx, client.ObjectKey{Name: serviceName, Namespace: tenant.Status.Namespace}, existingService)
	if err != nil {
		if errors.IsNotFound(err) {
			// Create Service
			if err := r.Create(ctx, service); err != nil {
				logger.Error(err, "Failed to create Redis Service")
				return nil, err
			}
			logger.Info("Created Redis Service", "service", serviceName)
			return service, nil
		}
		logger.Error(err, "Failed to get Redis Service")
		return nil, err
	}

	// No need to update the service as it's a headless service with minimal configuration
	return existingService, nil
}

// reconcileRedisStatefulSet creates or updates the Redis StatefulSet for the tenant
func (r *TenantReconciler) reconcileRedisStatefulSet(ctx context.Context, tenant *neurallogv1.Tenant, configMap *corev1.ConfigMap) (*appsv1.StatefulSet, error) {
	logger := log.FromContext(ctx)

	// Define StatefulSet name
	statefulSetName := fmt.Sprintf("%s-redis", tenant.Name)

	// Define labels
	labels := map[string]string{
		"app":                     "redis",
		"neurallog.io/tenant":     tenant.Name,
		"neurallog.io/component":  "redis",
		"neurallog.io/managed-by": "tenant-operator",
	}

	// Set replicas
	replicas := int32(1)
	if tenant.Spec.Redis != nil && tenant.Spec.Redis.Replicas > 0 {
		replicas = tenant.Spec.Redis.Replicas
	}

	// Set image
	image := "redis:7-alpine"
	if tenant.Spec.Redis != nil && tenant.Spec.Redis.Image != "" {
		image = tenant.Spec.Redis.Image
	}

	// Set resource requirements
	resources := corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("100m"),
			corev1.ResourceMemory: resource.MustParse("128Mi"),
		},
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("300m"),
			corev1.ResourceMemory: resource.MustParse("256Mi"),
		},
	}

	if tenant.Spec.Redis != nil && tenant.Spec.Redis.Resources != nil {
		if tenant.Spec.Redis.Resources.CPU != nil {
			if tenant.Spec.Redis.Resources.CPU.Request != "" {
				resources.Requests[corev1.ResourceCPU] = resource.MustParse(tenant.Spec.Redis.Resources.CPU.Request)
			}
			if tenant.Spec.Redis.Resources.CPU.Limit != "" {
				resources.Limits[corev1.ResourceCPU] = resource.MustParse(tenant.Spec.Redis.Resources.CPU.Limit)
			}
		}
		if tenant.Spec.Redis.Resources.Memory != nil {
			if tenant.Spec.Redis.Resources.Memory.Request != "" {
				resources.Requests[corev1.ResourceMemory] = resource.MustParse(tenant.Spec.Redis.Resources.Memory.Request)
			}
			if tenant.Spec.Redis.Resources.Memory.Limit != "" {
				resources.Limits[corev1.ResourceMemory] = resource.MustParse(tenant.Spec.Redis.Resources.Memory.Limit)
			}
		}
	}

	// Set storage
	storage := "1Gi"
	if tenant.Spec.Redis != nil && tenant.Spec.Redis.Storage != "" {
		storage = tenant.Spec.Redis.Storage
	}

	// Define StatefulSet
	statefulSet := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      statefulSetName,
			Namespace: tenant.Status.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: fmt.Sprintf("%s-redis", tenant.Name),
			Replicas:    &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":                 "redis",
					"neurallog.io/tenant": tenant.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "redis",
							Image: image,
							Command: []string{
								"redis-server",
								"/etc/redis/redis.conf",
							},
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 6379,
									Name:          "redis",
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "redis-data",
									MountPath: "/data",
								},
								{
									Name:      "redis-config",
									MountPath: "/etc/redis",
								},
							},
							Resources: resources,
							LivenessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									TCPSocket: &corev1.TCPSocketAction{
										Port: intstr.FromString("redis"),
									},
								},
								InitialDelaySeconds: 15,
								PeriodSeconds:       20,
							},
							ReadinessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									Exec: &corev1.ExecAction{
										Command: []string{
											"redis-cli",
											"ping",
										},
									},
								},
								InitialDelaySeconds: 5,
								PeriodSeconds:       10,
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "redis-config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: configMap.Name,
									},
									Items: []corev1.KeyToPath{
										{
											Key:  "redis.conf",
											Path: "redis.conf",
										},
									},
								},
							},
						},
					},
				},
			},
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "redis-data",
					},
					Spec: corev1.PersistentVolumeClaimSpec{
						AccessModes: []corev1.PersistentVolumeAccessMode{
							corev1.ReadWriteOnce,
						},
						Resources: corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceStorage: resource.MustParse(storage),
							},
						},
					},
				},
			},
		},
	}

	// Set owner reference
	if err := controllerutil.SetControllerReference(tenant, statefulSet, r.Scheme); err != nil {
		logger.Error(err, "Failed to set owner reference on Redis StatefulSet")
		return nil, err
	}

	// Create or update StatefulSet
	existingStatefulSet := &appsv1.StatefulSet{}
	err := r.Get(ctx, client.ObjectKey{Name: statefulSetName, Namespace: tenant.Status.Namespace}, existingStatefulSet)
	if err != nil {
		if errors.IsNotFound(err) {
			// Create StatefulSet
			if err := r.Create(ctx, statefulSet); err != nil {
				logger.Error(err, "Failed to create Redis StatefulSet")
				return nil, err
			}
			logger.Info("Created Redis StatefulSet", "statefulSet", statefulSetName)
			return statefulSet, nil
		}
		logger.Error(err, "Failed to get Redis StatefulSet")
		return nil, err
	}

	// Update StatefulSet if needed
	updated := false

	// Check if replicas need to be updated
	if existingStatefulSet.Spec.Replicas == nil || *existingStatefulSet.Spec.Replicas != replicas {
		existingStatefulSet.Spec.Replicas = &replicas
		updated = true
	}

	// Check if image needs to be updated
	if existingStatefulSet.Spec.Template.Spec.Containers[0].Image != image {
		existingStatefulSet.Spec.Template.Spec.Containers[0].Image = image
		updated = true
	}

	// Check if resources need to be updated
	if !resourcesEqual(existingStatefulSet.Spec.Template.Spec.Containers[0].Resources, resources) {
		existingStatefulSet.Spec.Template.Spec.Containers[0].Resources = resources
		updated = true
	}

	// Check if ConfigMap reference needs to be updated
	for i, volume := range existingStatefulSet.Spec.Template.Spec.Volumes {
		if volume.Name == "redis-config" && volume.ConfigMap != nil && volume.ConfigMap.Name != configMap.Name {
			existingStatefulSet.Spec.Template.Spec.Volumes[i].ConfigMap.Name = configMap.Name
			updated = true
		}
	}

	// Update StatefulSet if needed
	if updated {
		if err := r.Update(ctx, existingStatefulSet); err != nil {
			logger.Error(err, "Failed to update Redis StatefulSet")
			return nil, err
		}
		logger.Info("Updated Redis StatefulSet", "statefulSet", statefulSetName)
	}

	return existingStatefulSet, nil
}

// resourcesEqual compares two ResourceRequirements for equality
func resourcesEqual(a, b corev1.ResourceRequirements) bool {
	// Compare CPU requests
	if !resourceValueEqual(a.Requests.Cpu(), b.Requests.Cpu()) {
		return false
	}

	// Compare Memory requests
	if !resourceValueEqual(a.Requests.Memory(), b.Requests.Memory()) {
		return false
	}

	// Compare CPU limits
	if !resourceValueEqual(a.Limits.Cpu(), b.Limits.Cpu()) {
		return false
	}

	// Compare Memory limits
	if !resourceValueEqual(a.Limits.Memory(), b.Limits.Memory()) {
		return false
	}

	return true
}

// resourceValueEqual compares two resource.Quantity pointers for equality
func resourceValueEqual(a, b *resource.Quantity) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.Equal(*b)
}

// reconcileServerService creates or updates the Server Service for the tenant
func (r *TenantReconciler) reconcileServerService(ctx context.Context, tenant *neurallogv1.Tenant) (*corev1.Service, error) {
	logger := log.FromContext(ctx)

	// Define Service name
	serviceName := fmt.Sprintf("%s-server", tenant.Name)

	// Define Service
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: tenant.Status.Namespace,
			Labels: map[string]string{
				"app":                     "neurallog-server",
				"neurallog.io/tenant":     tenant.Name,
				"neurallog.io/component":  "server",
				"neurallog.io/managed-by": "tenant-operator",
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app":                 "neurallog-server",
				"neurallog.io/tenant": tenant.Name,
			},
			Ports: []corev1.ServicePort{
				{
					Port:       3030,
					TargetPort: intstr.FromString("http"),
					Name:       "http",
				},
			},
			Type: corev1.ServiceTypeClusterIP,
		},
	}

	// Set owner reference
	if err := controllerutil.SetControllerReference(tenant, service, r.Scheme); err != nil {
		logger.Error(err, "Failed to set owner reference on Server Service")
		return nil, err
	}

	// Create or update Service
	existingService := &corev1.Service{}
	err := r.Get(ctx, client.ObjectKey{Name: serviceName, Namespace: tenant.Status.Namespace}, existingService)
	if err != nil {
		if errors.IsNotFound(err) {
			// Create Service
			if err := r.Create(ctx, service); err != nil {
				logger.Error(err, "Failed to create Server Service")
				return nil, err
			}
			logger.Info("Created Server Service", "service", serviceName)
			return service, nil
		}
		logger.Error(err, "Failed to get Server Service")
		return nil, err
	}

	// Update Service if needed
	updated := false

	// Check if selector needs to be updated
	if !reflect.DeepEqual(existingService.Spec.Selector, service.Spec.Selector) {
		existingService.Spec.Selector = service.Spec.Selector
		updated = true
	}

	// Check if ports need to be updated
	if !reflect.DeepEqual(existingService.Spec.Ports, service.Spec.Ports) {
		existingService.Spec.Ports = service.Spec.Ports
		updated = true
	}

	// Update Service if needed
	if updated {
		if err := r.Update(ctx, existingService); err != nil {
			logger.Error(err, "Failed to update Server Service")
			return nil, err
		}
		logger.Info("Updated Server Service", "service", serviceName)
	}

	return existingService, nil
}

// reconcileServerDeployment creates or updates the Server Deployment for the tenant
func (r *TenantReconciler) reconcileServerDeployment(ctx context.Context, tenant *neurallogv1.Tenant) (*appsv1.Deployment, error) {
	logger := log.FromContext(ctx)

	// Define Deployment name
	deploymentName := fmt.Sprintf("%s-server", tenant.Name)

	// Define labels
	labels := map[string]string{
		"app":                     "neurallog-server",
		"neurallog.io/tenant":     tenant.Name,
		"neurallog.io/component":  "server",
		"neurallog.io/managed-by": "tenant-operator",
	}

	// Set replicas
	replicas := int32(1)
	if tenant.Spec.Server != nil && tenant.Spec.Server.Replicas > 0 {
		replicas = tenant.Spec.Server.Replicas
	}

	// Set image
	image := "neurallog/server:latest"
	if tenant.Spec.Server != nil && tenant.Spec.Server.Image != "" {
		image = tenant.Spec.Server.Image
	}

	// Set resource requirements
	resources := corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("100m"),
			corev1.ResourceMemory: resource.MustParse("128Mi"),
		},
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("500m"),
			corev1.ResourceMemory: resource.MustParse("512Mi"),
		},
	}

	if tenant.Spec.Server != nil && tenant.Spec.Server.Resources != nil {
		if tenant.Spec.Server.Resources.CPU != nil {
			if tenant.Spec.Server.Resources.CPU.Request != "" {
				resources.Requests[corev1.ResourceCPU] = resource.MustParse(tenant.Spec.Server.Resources.CPU.Request)
			}
			if tenant.Spec.Server.Resources.CPU.Limit != "" {
				resources.Limits[corev1.ResourceCPU] = resource.MustParse(tenant.Spec.Server.Resources.CPU.Limit)
			}
		}
		if tenant.Spec.Server.Resources.Memory != nil {
			if tenant.Spec.Server.Resources.Memory.Request != "" {
				resources.Requests[corev1.ResourceMemory] = resource.MustParse(tenant.Spec.Server.Resources.Memory.Request)
			}
			if tenant.Spec.Server.Resources.Memory.Limit != "" {
				resources.Limits[corev1.ResourceMemory] = resource.MustParse(tenant.Spec.Server.Resources.Memory.Limit)
			}
		}
	}

	// Define environment variables
	env := []corev1.EnvVar{
		{
			Name:  "NODE_ENV",
			Value: "production",
		},
		{
			Name:  "PORT",
			Value: "3030",
		},
		{
			Name:  "REDIS_URL",
			Value: fmt.Sprintf("redis://%s-redis:6379", tenant.Name),
		},
		{
			Name:  "LOG_LEVEL",
			Value: "info",
		},
		{
			Name:  "TENANT_ID",
			Value: tenant.Name,
		},
		{
			Name:  "AUTH_URL",
			Value: "http://auth:3000",
		},
	}

	// Add custom environment variables if specified
	if tenant.Spec.Server != nil && tenant.Spec.Server.Env != nil {
		for _, customEnv := range tenant.Spec.Server.Env {
			// Create a new EnvVar
			newEnv := corev1.EnvVar{
				Name: customEnv.Name,
			}

			// Set value or valueFrom
			if customEnv.Value != "" {
				newEnv.Value = customEnv.Value
			} else if customEnv.ValueFrom != nil {
				newEnv.ValueFrom = &corev1.EnvVarSource{}

				// Set ConfigMapKeyRef if specified
				if customEnv.ValueFrom.ConfigMapKeyRef != nil {
					newEnv.ValueFrom.ConfigMapKeyRef = &corev1.ConfigMapKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: customEnv.ValueFrom.ConfigMapKeyRef.Name,
						},
						Key:      customEnv.ValueFrom.ConfigMapKeyRef.Key,
						Optional: customEnv.ValueFrom.ConfigMapKeyRef.Optional,
					}
				}

				// Set SecretKeyRef if specified
				if customEnv.ValueFrom.SecretKeyRef != nil {
					newEnv.ValueFrom.SecretKeyRef = &corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: customEnv.ValueFrom.SecretKeyRef.Name,
						},
						Key:      customEnv.ValueFrom.SecretKeyRef.Key,
						Optional: customEnv.ValueFrom.SecretKeyRef.Optional,
					}
				}
			}

			// Add the environment variable
			env = append(env, newEnv)
		}
	}

	// Define Deployment
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentName,
			Namespace: tenant.Status.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":                 "neurallog-server",
					"neurallog.io/tenant": tenant.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            "server",
							Image:           image,
							ImagePullPolicy: corev1.PullIfNotPresent,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 3030,
									Name:          "http",
								},
							},
							Env:       env,
							Resources: resources,
							ReadinessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/health",
										Port: intstr.FromString("http"),
									},
								},
								InitialDelaySeconds: 5,
								PeriodSeconds:       10,
								TimeoutSeconds:      5,
								SuccessThreshold:    1,
								FailureThreshold:    3,
							},
							LivenessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/health",
										Port: intstr.FromString("http"),
									},
								},
								InitialDelaySeconds: 15,
								PeriodSeconds:       20,
								TimeoutSeconds:      5,
								SuccessThreshold:    1,
								FailureThreshold:    3,
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "tmp-volume",
									MountPath: "/tmp",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "tmp-volume",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
					},
				},
			},
		},
	}

	// Set owner reference
	if err := controllerutil.SetControllerReference(tenant, deployment, r.Scheme); err != nil {
		logger.Error(err, "Failed to set owner reference on Server Deployment")
		return nil, err
	}

	// Create or update Deployment
	existingDeployment := &appsv1.Deployment{}
	err := r.Get(ctx, client.ObjectKey{Name: deploymentName, Namespace: tenant.Status.Namespace}, existingDeployment)
	if err != nil {
		if errors.IsNotFound(err) {
			// Create Deployment
			if err := r.Create(ctx, deployment); err != nil {
				logger.Error(err, "Failed to create Server Deployment")
				return nil, err
			}
			logger.Info("Created Server Deployment", "deployment", deploymentName)
			return deployment, nil
		}
		logger.Error(err, "Failed to get Server Deployment")
		return nil, err
	}

	// Update Deployment if needed
	updated := false

	// Check if replicas need to be updated
	if existingDeployment.Spec.Replicas == nil || *existingDeployment.Spec.Replicas != replicas {
		existingDeployment.Spec.Replicas = &replicas
		updated = true
	}

	// Check if image needs to be updated
	if existingDeployment.Spec.Template.Spec.Containers[0].Image != image {
		existingDeployment.Spec.Template.Spec.Containers[0].Image = image
		updated = true
	}

	// Check if resources need to be updated
	if !resourcesEqual(existingDeployment.Spec.Template.Spec.Containers[0].Resources, resources) {
		existingDeployment.Spec.Template.Spec.Containers[0].Resources = resources
		updated = true
	}

	// Check if environment variables need to be updated
	if !envVarsEqual(existingDeployment.Spec.Template.Spec.Containers[0].Env, env) {
		existingDeployment.Spec.Template.Spec.Containers[0].Env = env
		updated = true
	}

	// Update Deployment if needed
	if updated {
		if err := r.Update(ctx, existingDeployment); err != nil {
			logger.Error(err, "Failed to update Server Deployment")
			return nil, err
		}
		logger.Info("Updated Server Deployment", "deployment", deploymentName)
	}

	return existingDeployment, nil
}

// envVarsEqual compares two slices of EnvVar for equality
func envVarsEqual(a, b []corev1.EnvVar) bool {
	if len(a) != len(b) {
		return false
	}

	// Create maps for easier comparison
	aMap := make(map[string]corev1.EnvVar)
	for _, env := range a {
		aMap[env.Name] = env
	}

	bMap := make(map[string]corev1.EnvVar)
	for _, env := range b {
		bMap[env.Name] = env
	}

	// Check if all keys in a are in b with the same values
	for name, envA := range aMap {
		envB, ok := bMap[name]
		if !ok {
			return false
		}

		// Compare Value
		if envA.Value != envB.Value {
			return false
		}

		// Compare ValueFrom
		if !envVarSourceEqual(envA.ValueFrom, envB.ValueFrom) {
			return false
		}
	}

	// Check if all keys in b are in a
	for name := range bMap {
		_, ok := aMap[name]
		if !ok {
			return false
		}
	}

	return true
}

// envVarSourceEqual compares two EnvVarSource pointers for equality
func envVarSourceEqual(a, b *corev1.EnvVarSource) bool {
	// If both are nil, they are equal
	if a == nil && b == nil {
		return true
	}

	// If only one is nil, they are not equal
	if a == nil || b == nil {
		return false
	}

	// Compare ConfigMapKeyRef
	if a.ConfigMapKeyRef != nil && b.ConfigMapKeyRef != nil {
		if a.ConfigMapKeyRef.Name != b.ConfigMapKeyRef.Name ||
			a.ConfigMapKeyRef.Key != b.ConfigMapKeyRef.Key ||
			(a.ConfigMapKeyRef.Optional != nil && b.ConfigMapKeyRef.Optional != nil && *a.ConfigMapKeyRef.Optional != *b.ConfigMapKeyRef.Optional) {
			return false
		}
	} else if a.ConfigMapKeyRef != nil || b.ConfigMapKeyRef != nil {
		return false
	}

	// Compare SecretKeyRef
	if a.SecretKeyRef != nil && b.SecretKeyRef != nil {
		if a.SecretKeyRef.Name != b.SecretKeyRef.Name ||
			a.SecretKeyRef.Key != b.SecretKeyRef.Key ||
			(a.SecretKeyRef.Optional != nil && b.SecretKeyRef.Optional != nil && *a.SecretKeyRef.Optional != *b.SecretKeyRef.Optional) {
			return false
		}
	} else if a.SecretKeyRef != nil || b.SecretKeyRef != nil {
		return false
	}

	return true
}

// SetupWithManager sets up the controller with the Manager.
func (r *TenantReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&neurallogv1.Tenant{}).
		Complete(r)
}
