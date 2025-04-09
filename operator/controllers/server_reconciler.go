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
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	neurallogv1 "github.com/neurallog/operator/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

// reconcileServerDeployment creates or updates the Server Deployment
func (r *TenantReconciler) reconcileServerDeployment(ctx context.Context, tenant *neurallogv1.Tenant) (*appsv1.Deployment, error) {
	logger := log.FromContext(ctx)
	namespaceName := tenant.Status.Namespace

	// Default values
	replicas := int32(1)
	if tenant.Spec.Server.Replicas != nil {
		replicas = *tenant.Spec.Server.Replicas
	}

	image := "neurallog/server:latest"
	if tenant.Spec.Server.Image != "" {
		image = tenant.Spec.Server.Image
	}

	// Resource requirements
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

	// Apply custom resource requirements if provided
	if tenant.Spec.Server.Resources.CPU.Request != "" {
		resources.Requests[corev1.ResourceCPU] = resource.MustParse(tenant.Spec.Server.Resources.CPU.Request)
	}
	if tenant.Spec.Server.Resources.CPU.Limit != "" {
		resources.Limits[corev1.ResourceCPU] = resource.MustParse(tenant.Spec.Server.Resources.CPU.Limit)
	}
	if tenant.Spec.Server.Resources.Memory.Request != "" {
		resources.Requests[corev1.ResourceMemory] = resource.MustParse(tenant.Spec.Server.Resources.Memory.Request)
	}
	if tenant.Spec.Server.Resources.Memory.Limit != "" {
		resources.Limits[corev1.ResourceMemory] = resource.MustParse(tenant.Spec.Server.Resources.Memory.Limit)
	}

	// Environment variables
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
			Value: "redis://redis:6379",
		},
		{
			Name:  "LOG_LEVEL",
			Value: "info",
		},
		{
			Name:  "TENANT_ID",
			Value: tenant.Name,
		},
	}

	// Add custom environment variables if provided
	if tenant.Spec.Server.Env != nil {
		for _, envVar := range tenant.Spec.Server.Env {
			kubeEnvVar := corev1.EnvVar{
				Name:  envVar.Name,
				Value: envVar.Value,
			}
			if envVar.ValueFrom != nil {
				kubeEnvVar.ValueFrom = &corev1.EnvVarSource{}
				if envVar.ValueFrom.ConfigMapKeyRef != nil {
					kubeEnvVar.ValueFrom.ConfigMapKeyRef = &corev1.ConfigMapKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: envVar.ValueFrom.ConfigMapKeyRef.Name,
						},
						Key:      envVar.ValueFrom.ConfigMapKeyRef.Key,
						Optional: envVar.ValueFrom.ConfigMapKeyRef.Optional,
					}
				}
				if envVar.ValueFrom.SecretKeyRef != nil {
					kubeEnvVar.ValueFrom.SecretKeyRef = &corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: envVar.ValueFrom.SecretKeyRef.Name,
						},
						Key:      envVar.ValueFrom.SecretKeyRef.Key,
						Optional: envVar.ValueFrom.SecretKeyRef.Optional,
					}
				}
			}
			env = append(env, kubeEnvVar)
		}
	}

	// Create Deployment object
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "neurallog-server",
			Namespace: namespaceName,
			Labels: map[string]string{
				"app":                    "neurallog-server",
				"neurallog.io/tenant":    tenant.Name,
				"neurallog.io/component": "server",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "neurallog-server",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":                    "neurallog-server",
						"neurallog.io/tenant":    tenant.Name,
						"neurallog.io/component": "server",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "server",
							Image: image,
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: 3030,
								},
							},
							Resources: resources,
							Env:       env,
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "tmp-volume",
									MountPath: "/tmp",
								},
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
								FailureThreshold:    3,
								SuccessThreshold:    1,
							},
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
								FailureThreshold:    3,
								SuccessThreshold:    1,
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

	// Create or update the Deployment
	existingDeployment := &appsv1.Deployment{}
	err := r.Get(ctx, client.ObjectKey{Name: deployment.Name, Namespace: deployment.Namespace}, existingDeployment)
	if err != nil {
		if errors.IsNotFound(err) {
			// Create Deployment
			if err := r.Create(ctx, deployment); err != nil {
				logger.Error(err, "Failed to create Server Deployment")
				return nil, err
			}
			logger.Info("Created Server Deployment", "deployment", deployment.Name)
			return deployment, nil
		}
		logger.Error(err, "Failed to get Server Deployment")
		return nil, err
	}

	// Update Deployment if it exists
	existingDeployment.Spec.Replicas = deployment.Spec.Replicas
	existingDeployment.Spec.Template.Spec.Containers[0].Image = deployment.Spec.Template.Spec.Containers[0].Image
	existingDeployment.Spec.Template.Spec.Containers[0].Resources = deployment.Spec.Template.Spec.Containers[0].Resources
	existingDeployment.Spec.Template.Spec.Containers[0].Env = deployment.Spec.Template.Spec.Containers[0].Env
	
	if err := r.Update(ctx, existingDeployment); err != nil {
		logger.Error(err, "Failed to update Server Deployment")
		return nil, err
	}
	logger.Info("Updated Server Deployment", "deployment", existingDeployment.Name)
	return existingDeployment, nil
}

// reconcileServerService creates or updates the Server Service
func (r *TenantReconciler) reconcileServerService(ctx context.Context, tenant *neurallogv1.Tenant) (*corev1.Service, error) {
	logger := log.FromContext(ctx)
	namespaceName := tenant.Status.Namespace

	// Create Service object
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "neurallog-server",
			Namespace: namespaceName,
			Labels: map[string]string{
				"app":                    "neurallog-server",
				"neurallog.io/tenant":    tenant.Name,
				"neurallog.io/component": "server",
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": "neurallog-server",
			},
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					Port:       3030,
					TargetPort: intstr.FromString("http"),
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

	// Create or update the Service
	existingService := &corev1.Service{}
	err := r.Get(ctx, client.ObjectKey{Name: service.Name, Namespace: service.Namespace}, existingService)
	if err != nil {
		if errors.IsNotFound(err) {
			// Create Service
			if err := r.Create(ctx, service); err != nil {
				logger.Error(err, "Failed to create Server Service")
				return nil, err
			}
			logger.Info("Created Server Service", "service", service.Name)
			return service, nil
		}
		logger.Error(err, "Failed to get Server Service")
		return nil, err
	}

	// Update Service if it exists
	existingService.Spec.Ports = service.Spec.Ports
	if err := r.Update(ctx, existingService); err != nil {
		logger.Error(err, "Failed to update Server Service")
		return nil, err
	}
	logger.Info("Updated Server Service", "service", existingService.Name)
	return existingService, nil
}

// updateServerStatus updates the Server status in the tenant
func (r *TenantReconciler) updateServerStatus(ctx context.Context, tenant *neurallogv1.Tenant, deployment *appsv1.Deployment) error {
	logger := log.FromContext(ctx)

	// Update Server status
	tenant.Status.ServerStatus = neurallogv1.ComponentStatus{
		TotalReplicas: *deployment.Spec.Replicas,
		ReadyReplicas: deployment.Status.ReadyReplicas,
	}

	// Set phase based on readiness
	if deployment.Status.ReadyReplicas == 0 {
		tenant.Status.ServerStatus.Phase = neurallogv1.ComponentPending
		tenant.Status.ServerStatus.Message = "Server is being provisioned"
	} else if deployment.Status.ReadyReplicas < *deployment.Spec.Replicas {
		tenant.Status.ServerStatus.Phase = neurallogv1.ComponentPending
		tenant.Status.ServerStatus.Message = fmt.Sprintf("Server is scaling up (%d/%d replicas ready)", deployment.Status.ReadyReplicas, *deployment.Spec.Replicas)
	} else {
		tenant.Status.ServerStatus.Phase = neurallogv1.ComponentRunning
		tenant.Status.ServerStatus.Message = "Server is running"
	}

	// Update tenant status
	if err := r.Status().Update(ctx, tenant); err != nil {
		logger.Error(err, "Failed to update Server status")
		return err
	}

	return nil
}
