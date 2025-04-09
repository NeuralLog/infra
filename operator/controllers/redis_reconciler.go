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

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	neurallogv1 "github.com/neurallog/operator/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

// reconcileRedisConfigMap creates or updates the Redis ConfigMap
func (r *TenantReconciler) reconcileRedisConfigMap(ctx context.Context, tenant *neurallogv1.Tenant) (*corev1.ConfigMap, error) {
	logger := log.FromContext(ctx)
	namespaceName := tenant.Status.Namespace

	// Default Redis configuration
	redisConf := `# Redis configuration for NeuralLog
port 6379
bind 0.0.0.0
protected-mode yes
daemonize no

# Memory management
maxmemory 256mb
maxmemory-policy allkeys-lru

# Persistence
appendonly yes
appendfsync everysec

# Logging
loglevel notice
logfile ""`

	// Apply custom configuration if provided
	if tenant.Spec.Redis.Config != nil {
		for key, value := range tenant.Spec.Redis.Config {
			redisConf += fmt.Sprintf("\n%s %s", key, value)
		}
	}

	// Create ConfigMap object
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "redis-config",
			Namespace: namespaceName,
			Labels: map[string]string{
				"app":                    "redis",
				"neurallog.io/tenant":    tenant.Name,
				"neurallog.io/component": "redis",
			},
		},
		Data: map[string]string{
			"redis.conf": redisConf,
		},
	}

	// Set owner reference
	if err := controllerutil.SetControllerReference(tenant, configMap, r.Scheme); err != nil {
		logger.Error(err, "Failed to set owner reference on Redis ConfigMap")
		return nil, err
	}

	// Create or update the ConfigMap
	existingConfigMap := &corev1.ConfigMap{}
	err := r.Get(ctx, client.ObjectKey{Name: configMap.Name, Namespace: configMap.Namespace}, existingConfigMap)
	if err != nil {
		if errors.IsNotFound(err) {
			// Create ConfigMap
			if err := r.Create(ctx, configMap); err != nil {
				logger.Error(err, "Failed to create Redis ConfigMap")
				return nil, err
			}
			logger.Info("Created Redis ConfigMap", "configMap", configMap.Name)
			return configMap, nil
		}
		logger.Error(err, "Failed to get Redis ConfigMap")
		return nil, err
	}

	// Update ConfigMap if it exists
	existingConfigMap.Data = configMap.Data
	if err := r.Update(ctx, existingConfigMap); err != nil {
		logger.Error(err, "Failed to update Redis ConfigMap")
		return nil, err
	}
	logger.Info("Updated Redis ConfigMap", "configMap", existingConfigMap.Name)
	return existingConfigMap, nil
}

// reconcileRedisStatefulSet creates or updates the Redis StatefulSet
func (r *TenantReconciler) reconcileRedisStatefulSet(ctx context.Context, tenant *neurallogv1.Tenant, configMap *corev1.ConfigMap) (*appsv1.StatefulSet, error) {
	logger := log.FromContext(ctx)
	namespaceName := tenant.Status.Namespace

	// Default values
	replicas := int32(1)
	if tenant.Spec.Redis.Replicas != nil {
		replicas = *tenant.Spec.Redis.Replicas
	}

	image := "redis:7-alpine"
	if tenant.Spec.Redis.Image != "" {
		image = tenant.Spec.Redis.Image
	}

	// Resource requirements
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

	// Apply custom resource requirements if provided
	if tenant.Spec.Redis.Resources.CPU.Request != "" {
		resources.Requests[corev1.ResourceCPU] = resource.MustParse(tenant.Spec.Redis.Resources.CPU.Request)
	}
	if tenant.Spec.Redis.Resources.CPU.Limit != "" {
		resources.Limits[corev1.ResourceCPU] = resource.MustParse(tenant.Spec.Redis.Resources.CPU.Limit)
	}
	if tenant.Spec.Redis.Resources.Memory.Request != "" {
		resources.Requests[corev1.ResourceMemory] = resource.MustParse(tenant.Spec.Redis.Resources.Memory.Request)
	}
	if tenant.Spec.Redis.Resources.Memory.Limit != "" {
		resources.Limits[corev1.ResourceMemory] = resource.MustParse(tenant.Spec.Redis.Resources.Memory.Limit)
	}

	// Storage size
	storageSize := "1Gi"
	if tenant.Spec.Redis.Storage != "" {
		storageSize = tenant.Spec.Redis.Storage
	}

	// Create StatefulSet object
	statefulSet := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "redis",
			Namespace: namespaceName,
			Labels: map[string]string{
				"app":                    "redis",
				"neurallog.io/tenant":    tenant.Name,
				"neurallog.io/component": "redis",
			},
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "redis",
				},
			},
			ServiceName: "redis",
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":                    "redis",
						"neurallog.io/tenant":    tenant.Name,
						"neurallog.io/component": "redis",
					},
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
									Name:          "redis",
									ContainerPort: 6379,
								},
							},
							Resources: resources,
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
								corev1.ResourceStorage: resource.MustParse(storageSize),
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

	// Create or update the StatefulSet
	existingStatefulSet := &appsv1.StatefulSet{}
	err := r.Get(ctx, client.ObjectKey{Name: statefulSet.Name, Namespace: statefulSet.Namespace}, existingStatefulSet)
	if err != nil {
		if errors.IsNotFound(err) {
			// Create StatefulSet
			if err := r.Create(ctx, statefulSet); err != nil {
				logger.Error(err, "Failed to create Redis StatefulSet")
				return nil, err
			}
			logger.Info("Created Redis StatefulSet", "statefulSet", statefulSet.Name)
			return statefulSet, nil
		}
		logger.Error(err, "Failed to get Redis StatefulSet")
		return nil, err
	}

	// Update StatefulSet if it exists
	existingStatefulSet.Spec.Replicas = statefulSet.Spec.Replicas
	existingStatefulSet.Spec.Template.Spec.Containers[0].Image = statefulSet.Spec.Template.Spec.Containers[0].Image
	existingStatefulSet.Spec.Template.Spec.Containers[0].Resources = statefulSet.Spec.Template.Spec.Containers[0].Resources
	
	if err := r.Update(ctx, existingStatefulSet); err != nil {
		logger.Error(err, "Failed to update Redis StatefulSet")
		return nil, err
	}
	logger.Info("Updated Redis StatefulSet", "statefulSet", existingStatefulSet.Name)
	return existingStatefulSet, nil
}

// reconcileRedisService creates or updates the Redis Service
func (r *TenantReconciler) reconcileRedisService(ctx context.Context, tenant *neurallogv1.Tenant) (*corev1.Service, error) {
	logger := log.FromContext(ctx)
	namespaceName := tenant.Status.Namespace

	// Create Service object
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "redis",
			Namespace: namespaceName,
			Labels: map[string]string{
				"app":                    "redis",
				"neurallog.io/tenant":    tenant.Name,
				"neurallog.io/component": "redis",
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": "redis",
			},
			Ports: []corev1.ServicePort{
				{
					Name:       "redis",
					Port:       6379,
					TargetPort: intstr.FromString("redis"),
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

	// Create or update the Service
	existingService := &corev1.Service{}
	err := r.Get(ctx, client.ObjectKey{Name: service.Name, Namespace: service.Namespace}, existingService)
	if err != nil {
		if errors.IsNotFound(err) {
			// Create Service
			if err := r.Create(ctx, service); err != nil {
				logger.Error(err, "Failed to create Redis Service")
				return nil, err
			}
			logger.Info("Created Redis Service", "service", service.Name)
			return service, nil
		}
		logger.Error(err, "Failed to get Redis Service")
		return nil, err
	}

	// Update Service if it exists
	existingService.Spec.Ports = service.Spec.Ports
	if err := r.Update(ctx, existingService); err != nil {
		logger.Error(err, "Failed to update Redis Service")
		return nil, err
	}
	logger.Info("Updated Redis Service", "service", existingService.Name)
	return existingService, nil
}

// updateRedisStatus updates the Redis status in the tenant
func (r *TenantReconciler) updateRedisStatus(ctx context.Context, tenant *neurallogv1.Tenant, statefulSet *appsv1.StatefulSet) error {
	logger := log.FromContext(ctx)

	// Update Redis status
	tenant.Status.RedisStatus = neurallogv1.ComponentStatus{
		TotalReplicas: *statefulSet.Spec.Replicas,
		ReadyReplicas: statefulSet.Status.ReadyReplicas,
	}

	// Set phase based on readiness
	if statefulSet.Status.ReadyReplicas == 0 {
		tenant.Status.RedisStatus.Phase = neurallogv1.ComponentPending
		tenant.Status.RedisStatus.Message = "Redis is being provisioned"
	} else if statefulSet.Status.ReadyReplicas < *statefulSet.Spec.Replicas {
		tenant.Status.RedisStatus.Phase = neurallogv1.ComponentPending
		tenant.Status.RedisStatus.Message = fmt.Sprintf("Redis is scaling up (%d/%d replicas ready)", statefulSet.Status.ReadyReplicas, *statefulSet.Spec.Replicas)
	} else {
		tenant.Status.RedisStatus.Phase = neurallogv1.ComponentRunning
		tenant.Status.RedisStatus.Message = "Redis is running"
	}

	// Update tenant status
	if err := r.Status().Update(ctx, tenant); err != nil {
		logger.Error(err, "Failed to update Redis status")
		return err
	}

	return nil
}
