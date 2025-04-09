# Security

This document provides guidance on security considerations for the NeuralLog Tenant Operator.

## Table of Contents

- [RBAC](#rbac)
  - [Service Accounts](#service-accounts)
  - [Roles and ClusterRoles](#roles-and-clusterroles)
  - [RoleBindings and ClusterRoleBindings](#rolebindings-and-clusterrolebindings)
  - [Least Privilege Principle](#least-privilege-principle)
- [Pod Security](#pod-security)
  - [Security Context](#security-context)
  - [Pod Security Standards](#pod-security-standards)
  - [Container Images](#container-images)
- [Network Security](#network-security)
  - [Network Policies](#network-policies)
  - [TLS](#tls)
  - [API Server Communication](#api-server-communication)
- [Secret Management](#secret-management)
  - [Kubernetes Secrets](#kubernetes-secrets)
  - [External Secret Management](#external-secret-management)
  - [Secret Rotation](#secret-rotation)
- [Certificate Management](#certificate-management)
  - [Webhook Certificates](#webhook-certificates)
  - [Client Certificates](#client-certificates)
  - [Certificate Rotation](#certificate-rotation)
- [Audit Logging](#audit-logging)
  - [Kubernetes Audit Logging](#kubernetes-audit-logging)
  - [Operator Audit Logging](#operator-audit-logging)
- [Vulnerability Management](#vulnerability-management)
  - [Dependency Scanning](#dependency-scanning)
  - [Container Scanning](#container-scanning)
  - [Code Scanning](#code-scanning)
- [Security Best Practices](#security-best-practices)

## RBAC

Role-Based Access Control (RBAC) is a method of regulating access to resources based on the roles of individual users. The NeuralLog Tenant Operator uses RBAC to control access to Kubernetes resources.

### Service Accounts

Service accounts are used to authenticate the operator with the Kubernetes API server:

```yaml
# config/rbac/service_account.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: controller-manager
  namespace: system
```

### Roles and ClusterRoles

Roles and ClusterRoles define permissions for resources:

```yaml
# config/rbac/role.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: manager-role
rules:
# Tenant resources
- apiGroups:
  - neurallog.io
  resources:
  - tenants
  verbs:
  - create
  - delete
  - get
  - list
  - patch
  - update
  - watch
- apiGroups:
  - neurallog.io
  resources:
  - tenants/status
  verbs:
  - get
  - patch
  - update
# Namespace resources
- apiGroups:
  - ""
  resources:
  - namespaces
  verbs:
  - create
  - delete
  - get
  - list
  - patch
  - update
  - watch
# Deployment resources
- apiGroups:
  - apps
  resources:
  - deployments
  verbs:
  - create
  - delete
  - get
  - list
  - patch
  - update
  - watch
# StatefulSet resources
- apiGroups:
  - apps
  resources:
  - statefulsets
  verbs:
  - create
  - delete
  - get
  - list
  - patch
  - update
  - watch
# Service resources
- apiGroups:
  - ""
  resources:
  - services
  verbs:
  - create
  - delete
  - get
  - list
  - patch
  - update
  - watch
# ConfigMap resources
- apiGroups:
  - ""
  resources:
  - configmaps
  verbs:
  - create
  - delete
  - get
  - list
  - patch
  - update
  - watch
# Secret resources
- apiGroups:
  - ""
  resources:
  - secrets
  verbs:
  - create
  - delete
  - get
  - list
  - patch
  - update
  - watch
# NetworkPolicy resources
- apiGroups:
  - networking.k8s.io
  resources:
  - networkpolicies
  verbs:
  - create
  - delete
  - get
  - list
  - patch
  - update
  - watch
# Events
- apiGroups:
  - ""
  resources:
  - events
  verbs:
  - create
  - patch
```

### RoleBindings and ClusterRoleBindings

RoleBindings and ClusterRoleBindings bind roles to service accounts:

```yaml
# config/rbac/role_binding.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: manager-rolebinding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: manager-role
subjects:
- kind: ServiceAccount
  name: controller-manager
  namespace: system
```

### Least Privilege Principle

The least privilege principle states that a subject should be given only the privileges needed to complete its task. The NeuralLog Tenant Operator follows this principle by:

1. **Limiting Permissions**: Only requesting permissions that are needed
2. **Using Namespaced Resources**: Using namespaced resources when possible
3. **Avoiding Cluster-Wide Permissions**: Avoiding cluster-wide permissions when possible
4. **Using Service Accounts**: Using service accounts with limited permissions

## Pod Security

Pod security ensures that pods run with appropriate security settings.

### Security Context

Security context defines privilege and access control settings for pods and containers:

```yaml
# config/manager/manager.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: controller-manager
  namespace: system
spec:
  template:
    spec:
      securityContext:
        runAsNonRoot: true
      containers:
      - name: manager
        securityContext:
          allowPrivilegeEscalation: false
          capabilities:
            drop:
            - ALL
          readOnlyRootFilesystem: true
          runAsNonRoot: true
          seccompProfile:
            type: RuntimeDefault
```

### Pod Security Standards

Pod Security Standards define different levels of security for pods:

1. **Privileged**: Unrestricted policy, providing the widest possible level of permissions
2. **Baseline**: Minimally restrictive policy which prevents known privilege escalations
3. **Restricted**: Heavily restricted policy, following current Pod hardening best practices

The NeuralLog Tenant Operator follows the Restricted policy:

```yaml
# config/manager/manager.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: controller-manager
  namespace: system
  labels:
    pod-security.kubernetes.io/enforce: restricted
    pod-security.kubernetes.io/audit: restricted
    pod-security.kubernetes.io/warn: restricted
```

### Container Images

Container images should be secure and minimal:

1. **Use Minimal Base Images**: Use minimal base images like Alpine or distroless
2. **Use Specific Tags**: Use specific tags instead of latest
3. **Scan Images for Vulnerabilities**: Scan images for vulnerabilities
4. **Use Multi-Stage Builds**: Use multi-stage builds to minimize image size

```dockerfile
# Build stage
FROM golang:1.20-alpine AS build

WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY main.go main.go
COPY api/ api/
COPY controllers/ controllers/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o manager main.go

# Final stage
FROM gcr.io/distroless/static:nonroot

WORKDIR /
COPY --from=build /workspace/manager .
USER 65532:65532

ENTRYPOINT ["/manager"]
```

## Network Security

Network security ensures that network communication is secure.

### Network Policies

Network policies define how pods communicate with each other and other network endpoints:

```yaml
# config/network/network_policy.yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: default-deny
  namespace: system
spec:
  podSelector: {}
  policyTypes:
  - Ingress
  - Egress
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-webhook
  namespace: system
spec:
  podSelector:
    matchLabels:
      control-plane: controller-manager
  policyTypes:
  - Ingress
  ingress:
  - from:
    - namespaceSelector:
        matchLabels:
          kubernetes.io/metadata.name: kube-system
    ports:
    - protocol: TCP
      port: 9443
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-metrics
  namespace: system
spec:
  podSelector:
    matchLabels:
      control-plane: controller-manager
  policyTypes:
  - Ingress
  ingress:
  - from:
    - namespaceSelector:
        matchLabels:
          kubernetes.io/metadata.name: monitoring
    ports:
    - protocol: TCP
      port: 8080
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-api-server
  namespace: system
spec:
  podSelector:
    matchLabels:
      control-plane: controller-manager
  policyTypes:
  - Egress
  egress:
  - to:
    - ipBlock:
        cidr: <api-server-ip>/32
    ports:
    - protocol: TCP
      port: 443
```

### TLS

Transport Layer Security (TLS) encrypts network communication:

1. **Webhook Server**: The webhook server uses TLS to secure communication with the API server
2. **Metrics Server**: The metrics server can use TLS to secure communication with Prometheus
3. **API Server Communication**: Communication with the API server uses TLS

```go
// main.go
func main() {
    // ... existing code
    
    mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
        Scheme:             scheme,
        MetricsBindAddress: metricsAddr,
        Port:               9443, // HTTPS port for webhook server
        CertDir:            "/tmp/k8s-webhook-server/serving-certs", // Certificate directory
        LeaderElection:     enableLeaderElection,
        LeaderElectionID:   "neurallog-operator.neurallog.io",
    })
    
    // ... existing code
}
```

### API Server Communication

Communication with the API server should be secure:

1. **Use TLS**: Use TLS for API server communication
2. **Use Service Account Tokens**: Use service account tokens for authentication
3. **Use RBAC**: Use RBAC for authorization
4. **Limit API Server Access**: Limit API server access to necessary operations

## Secret Management

Secret management ensures that sensitive information is secure.

### Kubernetes Secrets

Kubernetes Secrets store sensitive information:

```yaml
# config/manager/webhook_secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: webhook-server-cert
  namespace: system
type: kubernetes.io/tls
data:
  tls.crt: <base64-encoded-cert>
  tls.key: <base64-encoded-key>
```

### External Secret Management

External secret management systems provide additional security:

1. **HashiCorp Vault**: A tool for secrets management, encryption as a service, and privileged access management
2. **AWS Secrets Manager**: A service to store, manage, and retrieve secrets
3. **Google Secret Manager**: A service to store API keys, passwords, certificates, and other sensitive data
4. **Azure Key Vault**: A service to safeguard cryptographic keys and other secrets

```go
// Example using HashiCorp Vault
import (
    "github.com/hashicorp/vault/api"
)

func getSecret(path string, key string) (string, error) {
    // Create a Vault client
    client, err := api.NewClient(&api.Config{
        Address: "https://vault.example.com",
    })
    if err != nil {
        return "", err
    }

    // Read the secret
    secret, err := client.Logical().Read(path)
    if err != nil {
        return "", err
    }

    // Get the value
    value, ok := secret.Data[key].(string)
    if !ok {
        return "", fmt.Errorf("secret not found")
    }

    return value, nil
}
```

### Secret Rotation

Secret rotation ensures that secrets are regularly updated:

1. **Automatic Rotation**: Automatically rotate secrets on a schedule
2. **Manual Rotation**: Manually rotate secrets when needed
3. **Rotation Policies**: Define policies for secret rotation
4. **Rotation Monitoring**: Monitor secret rotation

```go
// Example secret rotation
func rotateSecret(ctx context.Context, client client.Client, secretName string, namespace string) error {
    // Get the secret
    secret := &corev1.Secret{}
    if err := client.Get(ctx, types.NamespacedName{Name: secretName, Namespace: namespace}, secret); err != nil {
        return err
    }

    // Generate a new secret
    newSecret, err := generateSecret()
    if err != nil {
        return err
    }

    // Update the secret
    secret.Data = newSecret
    if err := client.Update(ctx, secret); err != nil {
        return err
    }

    return nil
}
```

## Certificate Management

Certificate management ensures that certificates are secure and up-to-date.

### Webhook Certificates

Webhook certificates secure communication between the webhook server and the API server:

```yaml
# config/certmanager/certificate.yaml
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: serving-cert
  namespace: system
spec:
  dnsNames:
  - webhook-service.system.svc
  - webhook-service.system.svc.cluster.local
  issuerRef:
    kind: Issuer
    name: selfsigned-issuer
  secretName: webhook-server-cert
---
apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
  name: selfsigned-issuer
  namespace: system
spec:
  selfSigned: {}
```

### Client Certificates

Client certificates authenticate the operator with the API server:

```go
// Example using client certificates
import (
    "crypto/tls"
    "crypto/x509"
    "io/ioutil"
    "net/http"
)

func createClient() (*http.Client, error) {
    // Load client certificate
    cert, err := tls.LoadX509KeyPair("client.crt", "client.key")
    if err != nil {
        return nil, err
    }

    // Load CA certificate
    caCert, err := ioutil.ReadFile("ca.crt")
    if err != nil {
        return nil, err
    }

    // Create CA certificate pool
    caCertPool := x509.NewCertPool()
    caCertPool.AppendCertsFromPEM(caCert)

    // Create TLS configuration
    tlsConfig := &tls.Config{
        Certificates: []tls.Certificate{cert},
        RootCAs:      caCertPool,
    }

    // Create HTTP client
    client := &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: tlsConfig,
        },
    }

    return client, nil
}
```

### Certificate Rotation

Certificate rotation ensures that certificates are regularly updated:

1. **Automatic Rotation**: Automatically rotate certificates on a schedule
2. **Manual Rotation**: Manually rotate certificates when needed
3. **Rotation Policies**: Define policies for certificate rotation
4. **Rotation Monitoring**: Monitor certificate rotation

```go
// Example certificate rotation using cert-manager
func rotateCertificate(ctx context.Context, client client.Client, certificateName string, namespace string) error {
    // Get the certificate
    certificate := &certmanagerv1.Certificate{}
    if err := client.Get(ctx, types.NamespacedName{Name: certificateName, Namespace: namespace}, certificate); err != nil {
        return err
    }

    // Trigger certificate rotation
    certificate.Spec.RenewBefore = &metav1.Duration{Duration: 24 * time.Hour}
    if err := client.Update(ctx, certificate); err != nil {
        return err
    }

    return nil
}
```

## Audit Logging

Audit logging records actions for security analysis.

### Kubernetes Audit Logging

Kubernetes audit logging records API server requests:

```yaml
# audit-policy.yaml
apiVersion: audit.k8s.io/v1
kind: Policy
rules:
- level: Metadata
  resources:
  - group: neurallog.io
    resources: ["tenants"]
- level: RequestResponse
  resources:
  - group: neurallog.io
    resources: ["tenants/status"]
- level: Metadata
  resources:
  - group: ""
    resources: ["namespaces"]
- level: Metadata
  resources:
  - group: apps
    resources: ["deployments", "statefulsets"]
- level: Metadata
  resources:
  - group: ""
    resources: ["services", "configmaps", "secrets"]
- level: Metadata
  resources:
  - group: networking.k8s.io
    resources: ["networkpolicies"]
```

### Operator Audit Logging

Operator audit logging records operator actions:

```go
// Example audit logging
func auditLog(ctx context.Context, action string, resource string, name string, namespace string, result string) {
    logger := log.FromContext(ctx)
    logger.Info("Audit",
        "action", action,
        "resource", resource,
        "name", name,
        "namespace", namespace,
        "result", result,
        "user", "system:serviceaccount:system:controller-manager",
        "timestamp", time.Now().Format(time.RFC3339),
    )
}
```

## Vulnerability Management

Vulnerability management identifies and addresses security vulnerabilities.

### Dependency Scanning

Dependency scanning identifies vulnerabilities in dependencies:

```bash
# Example using Go vulnerability scanning
go list -m all | xargs go list -json | nancy sleuth
```

### Container Scanning

Container scanning identifies vulnerabilities in container images:

```bash
# Example using Trivy
trivy image neurallog/tenant-operator:latest
```

### Code Scanning

Code scanning identifies vulnerabilities in code:

```bash
# Example using gosec
gosec ./...
```

## Security Best Practices

Follow these security best practices:

1. **Follow the Principle of Least Privilege**: Grant only the permissions needed
2. **Use RBAC**: Use RBAC for authorization
3. **Use Service Accounts**: Use service accounts with limited permissions
4. **Use Security Context**: Use security context to restrict pod and container privileges
5. **Use Network Policies**: Use network policies to restrict network communication
6. **Use TLS**: Use TLS for network communication
7. **Use Secrets**: Use secrets for sensitive information
8. **Rotate Secrets and Certificates**: Regularly rotate secrets and certificates
9. **Use Audit Logging**: Use audit logging to record actions
10. **Scan for Vulnerabilities**: Regularly scan for vulnerabilities
11. **Keep Dependencies Updated**: Keep dependencies updated to address vulnerabilities
12. **Use Minimal Container Images**: Use minimal container images to reduce attack surface
13. **Use Pod Security Standards**: Use pod security standards to enforce security policies
14. **Monitor Security Events**: Monitor security events to detect and respond to security incidents
15. **Follow Security Best Practices**: Follow security best practices for Kubernetes and Go
