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
	"net/http"
	"net/http/httptest"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	neurallogv1 "github.com/neurallog/operator/api/v1"
)

var _ = Describe("Tenant Controller", func() {
	// Define utility constants for object names and testing timeouts/durations
	const (
		TenantName      = "test-tenant"
		TenantNamespace = "tenant-test-tenant"
		Timeout         = time.Second * 10
		Interval        = time.Millisecond * 250
	)

	Context("When creating a Tenant", func() {
		It("Should create a namespace, Redis, and Server resources", func() {
			By("Creating a new Tenant")
			ctx := context.Background()
			tenant := &neurallogv1.Tenant{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "neurallog.io/v1",
					Kind:       "Tenant",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: TenantName,
				},
				Spec: neurallogv1.TenantSpec{
					DisplayName: "Test Tenant",
					Description: "A tenant for testing",
					Redis: &neurallogv1.RedisSpec{
						Replicas: 1,
						Image:    "redis:7-alpine",
						Resources: &neurallogv1.ResourceSpec{
							CPU: &neurallogv1.ResourceLimitSpec{
								Request: "100m",
								Limit:   "200m",
							},
							Memory: &neurallogv1.ResourceLimitSpec{
								Request: "128Mi",
								Limit:   "256Mi",
							},
						},
						Storage: "1Gi",
					},
					Server: &neurallogv1.ServerSpec{
						Replicas: 1,
						Image:    "neurallog/server:latest",
						Resources: &neurallogv1.ResourceSpec{
							CPU: &neurallogv1.ResourceLimitSpec{
								Request: "100m",
								Limit:   "300m",
							},
							Memory: &neurallogv1.ResourceLimitSpec{
								Request: "128Mi",
								Limit:   "256Mi",
							},
						},
						Env: []neurallogv1.EnvVar{
							{
								Name:  "TEST_ENV",
								Value: "test-value",
							},
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, tenant)).Should(Succeed())

			// Mock Auth Service
			mockAuthServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case "GET":
					// List tenants
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					fmt.Fprintln(w, `{"status":"success","tenants":[]}`)
				case "POST":
					// Create tenant
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusCreated)
					fmt.Fprintln(w, `{"status":"success","message":"Tenant created successfully"}`)
				case "DELETE":
					// Delete tenant
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					fmt.Fprintln(w, `{"status":"success","message":"Tenant deleted successfully"}`)
				}
			}))
			defer mockAuthServer.Close()

			// Wait for the tenant to be created
			tenantLookupKey := types.NamespacedName{Name: TenantName}
			createdTenant := &neurallogv1.Tenant{}

			// We'll need to retry getting this newly created Tenant, given that creation may not immediately happen.
			Eventually(func() bool {
				err := k8sClient.Get(ctx, tenantLookupKey, createdTenant)
				return err == nil
			}, Timeout, Interval).Should(BeTrue())

			// Let's make sure our Tenant has a status field and is in the Provisioning phase
			Eventually(func() neurallogv1.TenantPhase {
				err := k8sClient.Get(ctx, tenantLookupKey, createdTenant)
				if err != nil {
					return ""
				}
				return createdTenant.Status.Phase
			}, Timeout, Interval).Should(Equal(neurallogv1.TenantProvisioning))

			// Check if the namespace was created
			namespaceLookupKey := types.NamespacedName{Name: TenantNamespace}
			createdNamespace := &corev1.Namespace{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, namespaceLookupKey, createdNamespace)
				return err == nil
			}, Timeout, Interval).Should(BeTrue())

			// Check if the Redis resources were created
			redisConfigMapLookupKey := types.NamespacedName{Name: fmt.Sprintf("%s-redis-config", TenantName), Namespace: TenantNamespace}
			createdRedisConfigMap := &corev1.ConfigMap{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, redisConfigMapLookupKey, createdRedisConfigMap)
				return err == nil
			}, Timeout, Interval).Should(BeTrue())

			redisServiceLookupKey := types.NamespacedName{Name: fmt.Sprintf("%s-redis", TenantName), Namespace: TenantNamespace}
			createdRedisService := &corev1.Service{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, redisServiceLookupKey, createdRedisService)
				return err == nil
			}, Timeout, Interval).Should(BeTrue())

			redisStatefulSetLookupKey := types.NamespacedName{Name: fmt.Sprintf("%s-redis", TenantName), Namespace: TenantNamespace}
			createdRedisStatefulSet := &appsv1.StatefulSet{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, redisStatefulSetLookupKey, createdRedisStatefulSet)
				return err == nil
			}, Timeout, Interval).Should(BeTrue())

			// Check if the Server resources were created
			serverServiceLookupKey := types.NamespacedName{Name: fmt.Sprintf("%s-server", TenantName), Namespace: TenantNamespace}
			createdServerService := &corev1.Service{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, serverServiceLookupKey, createdServerService)
				return err == nil
			}, Timeout, Interval).Should(BeTrue())

			serverDeploymentLookupKey := types.NamespacedName{Name: fmt.Sprintf("%s-server", TenantName), Namespace: TenantNamespace}
			createdServerDeployment := &appsv1.Deployment{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, serverDeploymentLookupKey, createdServerDeployment)
				return err == nil
			}, Timeout, Interval).Should(BeTrue())

			// Verify that the Server deployment has the custom environment variable
			Expect(createdServerDeployment.Spec.Template.Spec.Containers[0].Env).To(ContainElement(corev1.EnvVar{
				Name:  "TEST_ENV",
				Value: "test-value",
			}))

			// Verify that the Redis StatefulSet has the correct resources
			Expect(createdRedisStatefulSet.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().String()).To(Equal("100m"))
			Expect(createdRedisStatefulSet.Spec.Template.Spec.Containers[0].Resources.Limits.Cpu().String()).To(Equal("200m"))
			Expect(createdRedisStatefulSet.Spec.Template.Spec.Containers[0].Resources.Requests.Memory().String()).To(Equal("128Mi"))
			Expect(createdRedisStatefulSet.Spec.Template.Spec.Containers[0].Resources.Limits.Memory().String()).To(Equal("256Mi"))

			// Verify that the Server deployment has the correct resources
			Expect(createdServerDeployment.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().String()).To(Equal("100m"))
			Expect(createdServerDeployment.Spec.Template.Spec.Containers[0].Resources.Limits.Cpu().String()).To(Equal("300m"))
			Expect(createdServerDeployment.Spec.Template.Spec.Containers[0].Resources.Requests.Memory().String()).To(Equal("128Mi"))
			Expect(createdServerDeployment.Spec.Template.Spec.Containers[0].Resources.Limits.Memory().String()).To(Equal("256Mi"))

			// Clean up
			By("Deleting the Tenant")
			Expect(k8sClient.Delete(ctx, tenant)).Should(Succeed())

			// Verify that the Tenant is deleted
			Eventually(func() bool {
				err := k8sClient.Get(ctx, tenantLookupKey, createdTenant)
				return errors.IsNotFound(err)
			}, Timeout, Interval).Should(BeTrue())

			// Verify that the namespace is deleted
			Eventually(func() bool {
				err := k8sClient.Get(ctx, namespaceLookupKey, createdNamespace)
				return errors.IsNotFound(err)
			}, Timeout, Interval).Should(BeTrue())
		})
	})
})
