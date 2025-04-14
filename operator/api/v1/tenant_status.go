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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TenantStatus defines the observed state of a NeuralLog Tenant
type TenantStatus struct {
	// Conditions represent the latest available observations of the tenant's state
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// Phase represents the current phase of the tenant
	// +optional
	Phase TenantPhase `json:"phase,omitempty"`

	// Namespace is the namespace created for the tenant
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// ServerStatus represents the status of the server deployment
	// +optional
	ServerStatus ComponentStatus `json:"serverStatus,omitempty"`

	// RedisStatus represents the status of the Redis deployment
	// +optional
	RedisStatus ComponentStatus `json:"redisStatus,omitempty"`

	// RegistryStatus represents the status of the registry deployment
	// +optional
	RegistryStatus ComponentStatus `json:"registryStatus,omitempty"`
}

// TenantPhase represents the phase of a tenant
type TenantPhase string

const (
	// TenantPending means the tenant is being created
	TenantPending TenantPhase = "Pending"

	// TenantProvisioning means the tenant resources are being provisioned
	TenantProvisioning TenantPhase = "Provisioning"

	// TenantRunning means the tenant is running
	TenantRunning TenantPhase = "Running"

	// TenantFailed means the tenant creation failed
	TenantFailed TenantPhase = "Failed"

	// TenantTerminating means the tenant is being deleted
	TenantTerminating TenantPhase = "Terminating"
)

// ComponentStatus represents the status of a component
type ComponentStatus struct {
	// Phase represents the current phase of the component
	// +optional
	Phase ComponentPhase `json:"phase,omitempty"`

	// Message provides additional information about the component status
	// +optional
	Message string `json:"message,omitempty"`

	// ReadyReplicas is the number of ready replicas
	// +optional
	ReadyReplicas int32 `json:"readyReplicas,omitempty"`

	// TotalReplicas is the total number of replicas
	// +optional
	TotalReplicas int32 `json:"totalReplicas,omitempty"`
}

// ComponentPhase represents the phase of a component
type ComponentPhase string

const (
	// ComponentPending means the component is being created
	ComponentPending ComponentPhase = "Pending"

	// ComponentProvisioning means the component is being provisioned
	ComponentProvisioning ComponentPhase = "Provisioning"

	// ComponentRunning means the component is running
	ComponentRunning ComponentPhase = "Running"

	// ComponentDegraded means the component is running but not all replicas are ready
	ComponentDegraded ComponentPhase = "Degraded"

	// ComponentFailed means the component creation failed
	ComponentFailed ComponentPhase = "Failed"
)
