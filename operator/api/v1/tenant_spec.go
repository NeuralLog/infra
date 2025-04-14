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

// TenantSpec defines the desired state of a NeuralLog Tenant
type TenantSpec struct {
	// DisplayName is a user-friendly name for the tenant
	// +optional
	DisplayName string `json:"displayName,omitempty"`

	// Description provides additional information about the tenant
	// +optional
	Description string `json:"description,omitempty"`

	// Resources defines the resource limits and requests for the tenant
	// +optional
	Resources ResourceRequirements `json:"resources,omitempty"`

	// Server defines the configuration for the NeuralLog server
	// +optional
	Server ServerSpec `json:"server,omitempty"`

	// Redis defines the configuration for the Redis instance
	// +optional
	Redis RedisSpec `json:"redis,omitempty"`

	// Registry defines the configuration for the Endpoint Registry service
	// +optional
	Registry RegistrySpec `json:"registry,omitempty"`

	// NetworkPolicy defines the network policy configuration for the tenant
	// +optional
	NetworkPolicy NetworkPolicySpec `json:"networkPolicy,omitempty"`
}

// ResourceRequirements defines the resource limits and requests for the tenant
type ResourceRequirements struct {
	// CPU defines the CPU limits and requests
	// +optional
	CPU ResourceLimit `json:"cpu,omitempty"`

	// Memory defines the memory limits and requests
	// +optional
	Memory ResourceLimit `json:"memory,omitempty"`

	// Storage defines the storage limits and requests
	// +optional
	Storage ResourceLimit `json:"storage,omitempty"`
}

// ResourceLimit defines a resource limit and request
type ResourceLimit struct {
	// Limit is the maximum amount of the resource
	// +optional
	Limit string `json:"limit,omitempty"`

	// Request is the minimum amount of the resource
	// +optional
	Request string `json:"request,omitempty"`
}

// ServerSpec defines the configuration for the NeuralLog server
type ServerSpec struct {
	// Replicas is the number of server instances
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`

	// Image is the Docker image for the server
	// +optional
	Image string `json:"image,omitempty"`

	// Resources defines the resource limits and requests for the server
	// +optional
	Resources ResourceRequirements `json:"resources,omitempty"`

	// Env defines additional environment variables for the server
	// +optional
	Env []EnvVar `json:"env,omitempty"`
}

// RedisSpec defines the configuration for the Redis instance
type RedisSpec struct {
	// Replicas is the number of Redis instances (for Redis Sentinel)
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`

	// Image is the Docker image for Redis
	// +optional
	Image string `json:"image,omitempty"`

	// Resources defines the resource limits and requests for Redis
	// +optional
	Resources ResourceRequirements `json:"resources,omitempty"`

	// Storage defines the storage configuration for Redis
	// +optional
	Storage string `json:"storage,omitempty"`

	// Config defines additional Redis configuration
	// +optional
	Config map[string]string `json:"config,omitempty"`
}

// RegistrySpec defines the configuration for the Endpoint Registry service
type RegistrySpec struct {
	// Replicas is the number of Registry instances
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`

	// Image is the Docker image for the Registry
	// +optional
	Image string `json:"image,omitempty"`

	// Resources defines the resource limits and requests for the Registry
	// +optional
	Resources ResourceRequirements `json:"resources,omitempty"`

	// BaseDomain is the base domain for endpoint URLs
	// +optional
	BaseDomain string `json:"baseDomain,omitempty"`
}

// NetworkPolicySpec defines the network policy configuration for the tenant
type NetworkPolicySpec struct {
	// Enabled indicates whether network policies should be created
	// +optional
	Enabled *bool `json:"enabled,omitempty"`

	// AllowedNamespaces is a list of namespaces that can access the tenant
	// +optional
	AllowedNamespaces []string `json:"allowedNamespaces,omitempty"`

	// IngressRules defines additional ingress rules
	// +optional
	IngressRules []NetworkPolicyRule `json:"ingressRules,omitempty"`

	// EgressRules defines additional egress rules
	// +optional
	EgressRules []NetworkPolicyRule `json:"egressRules,omitempty"`
}

// NetworkPolicyRule defines a network policy rule
type NetworkPolicyRule struct {
	// Description provides information about the rule
	// +optional
	Description string `json:"description,omitempty"`

	// From defines the source selector for ingress rules
	// +optional
	From map[string]string `json:"from,omitempty"`

	// To defines the destination selector for egress rules
	// +optional
	To map[string]string `json:"to,omitempty"`

	// Ports defines the ports for the rule
	// +optional
	Ports []NetworkPolicyPort `json:"ports,omitempty"`
}

// NetworkPolicyPort defines a port for a network policy rule
type NetworkPolicyPort struct {
	// Protocol is the protocol for the port
	// +optional
	Protocol string `json:"protocol,omitempty"`

	// Port is the port number
	// +optional
	Port int32 `json:"port,omitempty"`
}

// EnvVar defines an environment variable
type EnvVar struct {
	// Name is the name of the environment variable
	Name string `json:"name"`

	// Value is the value of the environment variable
	// +optional
	Value string `json:"value,omitempty"`

	// ValueFrom defines a source for the environment variable value
	// +optional
	ValueFrom *EnvVarSource `json:"valueFrom,omitempty"`
}

// EnvVarSource defines a source for an environment variable
type EnvVarSource struct {
	// ConfigMapKeyRef references a key in a ConfigMap
	// +optional
	ConfigMapKeyRef *ConfigMapKeySelector `json:"configMapKeyRef,omitempty"`

	// SecretKeyRef references a key in a Secret
	// +optional
	SecretKeyRef *SecretKeySelector `json:"secretKeyRef,omitempty"`
}

// ConfigMapKeySelector references a key in a ConfigMap
type ConfigMapKeySelector struct {
	// Name is the name of the ConfigMap
	Name string `json:"name"`

	// Key is the key in the ConfigMap
	Key string `json:"key"`

	// Optional indicates whether the ConfigMap or key must exist
	// +optional
	Optional *bool `json:"optional,omitempty"`
}

// SecretKeySelector references a key in a Secret
type SecretKeySelector struct {
	// Name is the name of the Secret
	Name string `json:"name"`

	// Key is the key in the Secret
	Key string `json:"key"`

	// Optional indicates whether the Secret or key must exist
	// +optional
	Optional *bool `json:"optional,omitempty"`
}
