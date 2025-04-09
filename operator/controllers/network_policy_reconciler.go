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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	corev1 "k8s.io/api/core/v1"
	neurallogv1 "github.com/neurallog/operator/api/v1"
	networkingv1 "k8s.io/api/networking/v1"
)

// reconcileNetworkPolicies creates or updates network policies for the tenant
func (r *TenantReconciler) reconcileNetworkPolicies(ctx context.Context, tenant *neurallogv1.Tenant) error {
	logger := log.FromContext(ctx)
	logger.Info("Reconciling Network Policies", "tenant", tenant.Name)

	// Skip if the namespace doesn't exist yet
	if tenant.Status.Namespace == "" {
		logger.Info("Namespace not yet created, skipping Network Policy reconciliation")
		return nil
	}

	// Check if network policies are enabled
	enabled := true
	if tenant.Spec.NetworkPolicy.Enabled != nil {
		enabled = *tenant.Spec.NetworkPolicy.Enabled
	}

	if !enabled {
		logger.Info("Network policies are disabled for this tenant", "tenant", tenant.Name)
		return nil
	}

	// Create default network policies
	if err := r.reconcileDefaultNetworkPolicies(ctx, tenant); err != nil {
		logger.Error(err, "Failed to reconcile default network policies")
		return err
	}

	// Create custom network policies
	if err := r.reconcileCustomNetworkPolicies(ctx, tenant); err != nil {
		logger.Error(err, "Failed to reconcile custom network policies")
		return err
	}

	return nil
}

// reconcileDefaultNetworkPolicies creates or updates default network policies
func (r *TenantReconciler) reconcileDefaultNetworkPolicies(ctx context.Context, tenant *neurallogv1.Tenant) error {
	logger := log.FromContext(ctx)
	namespaceName := tenant.Status.Namespace

	// Create default deny all ingress policy
	denyAllPolicy := &networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "default-deny-all",
			Namespace: namespaceName,
			Labels: map[string]string{
				"neurallog.io/tenant":    tenant.Name,
				"neurallog.io/component": "network-policy",
				"neurallog.io/policy":    "default",
			},
		},
		Spec: networkingv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{},
			PolicyTypes: []networkingv1.PolicyType{
				networkingv1.PolicyTypeIngress,
			},
		},
	}

	// Set owner reference
	if err := controllerutil.SetControllerReference(tenant, denyAllPolicy, r.Scheme); err != nil {
		logger.Error(err, "Failed to set owner reference on default deny network policy")
		return err
	}

	// Create or update the network policy
	existingPolicy := &networkingv1.NetworkPolicy{}
	err := r.Get(ctx, client.ObjectKey{Name: denyAllPolicy.Name, Namespace: denyAllPolicy.Namespace}, existingPolicy)
	if err != nil {
		if errors.IsNotFound(err) {
			// Create network policy
			if err := r.Create(ctx, denyAllPolicy); err != nil {
				logger.Error(err, "Failed to create default deny network policy")
				return err
			}
			logger.Info("Created default deny network policy", "policy", denyAllPolicy.Name)
		} else {
			logger.Error(err, "Failed to get default deny network policy")
			return err
		}
	} else {
		// Update network policy if it exists
		existingPolicy.Spec = denyAllPolicy.Spec
		if err := r.Update(ctx, existingPolicy); err != nil {
			logger.Error(err, "Failed to update default deny network policy")
			return err
		}
		logger.Info("Updated default deny network policy", "policy", existingPolicy.Name)
	}

	// Create allow internal traffic policy
	allowInternalPolicy := &networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "allow-internal-traffic",
			Namespace: namespaceName,
			Labels: map[string]string{
				"neurallog.io/tenant":    tenant.Name,
				"neurallog.io/component": "network-policy",
				"neurallog.io/policy":    "default",
			},
		},
		Spec: networkingv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{},
			Ingress: []networkingv1.NetworkPolicyIngressRule{
				{
					From: []networkingv1.NetworkPolicyPeer{
						{
							PodSelector: &metav1.LabelSelector{},
						},
					},
				},
			},
			PolicyTypes: []networkingv1.PolicyType{
				networkingv1.PolicyTypeIngress,
			},
		},
	}

	// Set owner reference
	if err := controllerutil.SetControllerReference(tenant, allowInternalPolicy, r.Scheme); err != nil {
		logger.Error(err, "Failed to set owner reference on allow internal traffic network policy")
		return err
	}

	// Create or update the network policy
	existingInternalPolicy := &networkingv1.NetworkPolicy{}
	err = r.Get(ctx, client.ObjectKey{Name: allowInternalPolicy.Name, Namespace: allowInternalPolicy.Namespace}, existingInternalPolicy)
	if err != nil {
		if errors.IsNotFound(err) {
			// Create network policy
			if err := r.Create(ctx, allowInternalPolicy); err != nil {
				logger.Error(err, "Failed to create allow internal traffic network policy")
				return err
			}
			logger.Info("Created allow internal traffic network policy", "policy", allowInternalPolicy.Name)
		} else {
			logger.Error(err, "Failed to get allow internal traffic network policy")
			return err
		}
	} else {
		// Update network policy if it exists
		existingInternalPolicy.Spec = allowInternalPolicy.Spec
		if err := r.Update(ctx, existingInternalPolicy); err != nil {
			logger.Error(err, "Failed to update allow internal traffic network policy")
			return err
		}
		logger.Info("Updated allow internal traffic network policy", "policy", existingInternalPolicy.Name)
	}

	// Create allow API access policy
	allowApiPolicy := &networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "allow-api-access",
			Namespace: namespaceName,
			Labels: map[string]string{
				"neurallog.io/tenant":    tenant.Name,
				"neurallog.io/component": "network-policy",
				"neurallog.io/policy":    "default",
			},
		},
		Spec: networkingv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "neurallog-server",
				},
			},
			Ingress: []networkingv1.NetworkPolicyIngressRule{
				{
					Ports: []networkingv1.NetworkPolicyPort{
						{
							Protocol: func() *corev1.Protocol {
								p := corev1.ProtocolTCP
								return &p
							}(),
							Port: &intstr.IntOrString{
								Type:   intstr.String,
								StrVal: "http",
							},
						},
					},
				},
			},
			PolicyTypes: []networkingv1.PolicyType{
				networkingv1.PolicyTypeIngress,
			},
		},
	}

	// Add allowed namespaces if specified
	if len(tenant.Spec.NetworkPolicy.AllowedNamespaces) > 0 {
		for _, ns := range tenant.Spec.NetworkPolicy.AllowedNamespaces {
			allowApiPolicy.Spec.Ingress[0].From = append(allowApiPolicy.Spec.Ingress[0].From, networkingv1.NetworkPolicyPeer{
				NamespaceSelector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"kubernetes.io/metadata.name": ns,
					},
				},
			})
		}
	}

	// Set owner reference
	if err := controllerutil.SetControllerReference(tenant, allowApiPolicy, r.Scheme); err != nil {
		logger.Error(err, "Failed to set owner reference on allow API access network policy")
		return err
	}

	// Create or update the network policy
	existingApiPolicy := &networkingv1.NetworkPolicy{}
	err = r.Get(ctx, client.ObjectKey{Name: allowApiPolicy.Name, Namespace: allowApiPolicy.Namespace}, existingApiPolicy)
	if err != nil {
		if errors.IsNotFound(err) {
			// Create network policy
			if err := r.Create(ctx, allowApiPolicy); err != nil {
				logger.Error(err, "Failed to create allow API access network policy")
				return err
			}
			logger.Info("Created allow API access network policy", "policy", allowApiPolicy.Name)
		} else {
			logger.Error(err, "Failed to get allow API access network policy")
			return err
		}
	} else {
		// Update network policy if it exists
		existingApiPolicy.Spec = allowApiPolicy.Spec
		if err := r.Update(ctx, existingApiPolicy); err != nil {
			logger.Error(err, "Failed to update allow API access network policy")
			return err
		}
		logger.Info("Updated allow API access network policy", "policy", existingApiPolicy.Name)
	}

	return nil
}

// reconcileCustomNetworkPolicies creates or updates custom network policies
func (r *TenantReconciler) reconcileCustomNetworkPolicies(ctx context.Context, tenant *neurallogv1.Tenant) error {
	logger := log.FromContext(ctx)
	namespaceName := tenant.Status.Namespace

	// Process custom ingress rules
	for i, rule := range tenant.Spec.NetworkPolicy.IngressRules {
		policyName := fmt.Sprintf("custom-ingress-%d", i)

		// Create network policy for the ingress rule
		ingressPolicy := &networkingv1.NetworkPolicy{
			ObjectMeta: metav1.ObjectMeta{
				Name:      policyName,
				Namespace: namespaceName,
				Labels: map[string]string{
					"neurallog.io/tenant":    tenant.Name,
					"neurallog.io/component": "network-policy",
					"neurallog.io/policy":    "custom",
				},
			},
			Spec: networkingv1.NetworkPolicySpec{
				PodSelector: metav1.LabelSelector{},
				Ingress:     []networkingv1.NetworkPolicyIngressRule{{}},
				PolicyTypes: []networkingv1.PolicyType{
					networkingv1.PolicyTypeIngress,
				},
			},
		}

		// Add from selectors if specified
		if len(rule.From) > 0 {
			ingressPolicy.Spec.Ingress[0].From = []networkingv1.NetworkPolicyPeer{
				{
					PodSelector: &metav1.LabelSelector{
						MatchLabels: rule.From,
					},
				},
			}
		}

		// Add ports if specified
		if len(rule.Ports) > 0 {
			for _, port := range rule.Ports {
				networkPort := networkingv1.NetworkPolicyPort{}

				if port.Protocol != "" {
					protocol := corev1.Protocol(port.Protocol)
					networkPort.Protocol = &protocol
				}

				if port.Port != 0 {
					portVal := intstr.FromInt(int(port.Port))
					networkPort.Port = &portVal
				}

				ingressPolicy.Spec.Ingress[0].Ports = append(ingressPolicy.Spec.Ingress[0].Ports, networkPort)
			}
		}

		// Set owner reference
		if err := controllerutil.SetControllerReference(tenant, ingressPolicy, r.Scheme); err != nil {
			logger.Error(err, "Failed to set owner reference on custom ingress network policy")
			return err
		}

		// Create or update the network policy
		existingPolicy := &networkingv1.NetworkPolicy{}
		err := r.Get(ctx, client.ObjectKey{Name: ingressPolicy.Name, Namespace: ingressPolicy.Namespace}, existingPolicy)
		if err != nil {
			if errors.IsNotFound(err) {
				// Create network policy
				if err := r.Create(ctx, ingressPolicy); err != nil {
					logger.Error(err, "Failed to create custom ingress network policy")
					return err
				}
				logger.Info("Created custom ingress network policy", "policy", ingressPolicy.Name)
			} else {
				logger.Error(err, "Failed to get custom ingress network policy")
				return err
			}
		} else {
			// Update network policy if it exists
			existingPolicy.Spec = ingressPolicy.Spec
			if err := r.Update(ctx, existingPolicy); err != nil {
				logger.Error(err, "Failed to update custom ingress network policy")
				return err
			}
			logger.Info("Updated custom ingress network policy", "policy", existingPolicy.Name)
		}
	}

	// Process custom egress rules
	for i, rule := range tenant.Spec.NetworkPolicy.EgressRules {
		policyName := fmt.Sprintf("custom-egress-%d", i)

		// Create network policy for the egress rule
		egressPolicy := &networkingv1.NetworkPolicy{
			ObjectMeta: metav1.ObjectMeta{
				Name:      policyName,
				Namespace: namespaceName,
				Labels: map[string]string{
					"neurallog.io/tenant":    tenant.Name,
					"neurallog.io/component": "network-policy",
					"neurallog.io/policy":    "custom",
				},
			},
			Spec: networkingv1.NetworkPolicySpec{
				PodSelector: metav1.LabelSelector{},
				Egress:      []networkingv1.NetworkPolicyEgressRule{{}},
				PolicyTypes: []networkingv1.PolicyType{
					networkingv1.PolicyTypeEgress,
				},
			},
		}

		// Add to selectors if specified
		if len(rule.To) > 0 {
			egressPolicy.Spec.Egress[0].To = []networkingv1.NetworkPolicyPeer{
				{
					PodSelector: &metav1.LabelSelector{
						MatchLabels: rule.To,
					},
				},
			}
		}

		// Add ports if specified
		if len(rule.Ports) > 0 {
			for _, port := range rule.Ports {
				networkPort := networkingv1.NetworkPolicyPort{}

				if port.Protocol != "" {
					protocol := corev1.Protocol(port.Protocol)
					networkPort.Protocol = &protocol
				}

				if port.Port != 0 {
					portVal := intstr.FromInt(int(port.Port))
					networkPort.Port = &portVal
				}

				egressPolicy.Spec.Egress[0].Ports = append(egressPolicy.Spec.Egress[0].Ports, networkPort)
			}
		}

		// Set owner reference
		if err := controllerutil.SetControllerReference(tenant, egressPolicy, r.Scheme); err != nil {
			logger.Error(err, "Failed to set owner reference on custom egress network policy")
			return err
		}

		// Create or update the network policy
		existingPolicy := &networkingv1.NetworkPolicy{}
		err := r.Get(ctx, client.ObjectKey{Name: egressPolicy.Name, Namespace: egressPolicy.Namespace}, existingPolicy)
		if err != nil {
			if errors.IsNotFound(err) {
				// Create network policy
				if err := r.Create(ctx, egressPolicy); err != nil {
					logger.Error(err, "Failed to create custom egress network policy")
					return err
				}
				logger.Info("Created custom egress network policy", "policy", egressPolicy.Name)
			} else {
				logger.Error(err, "Failed to get custom egress network policy")
				return err
			}
		} else {
			// Update network policy if it exists
			existingPolicy.Spec = egressPolicy.Spec
			if err := r.Update(ctx, existingPolicy); err != nil {
				logger.Error(err, "Failed to update custom egress network policy")
				return err
			}
			logger.Info("Updated custom egress network policy", "policy", existingPolicy.Name)
		}
	}

	return nil
}
