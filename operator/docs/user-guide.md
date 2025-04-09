# NeuralLog Tenant Operator User Guide

This guide provides instructions for using the NeuralLog Tenant Operator to manage tenants in the NeuralLog platform.

## Introduction

The NeuralLog Tenant Operator allows you to create and manage isolated tenant environments within a Kubernetes cluster. Each tenant gets its own dedicated namespace, server deployment, Redis instance, and network policies.

## Prerequisites

Before you begin, make sure you have:

- A Kubernetes cluster with the NeuralLog Tenant Operator installed
- `kubectl` configured to communicate with your cluster
- Appropriate permissions to create and manage Tenant resources

## Creating a Tenant

### Basic Tenant

To create a basic tenant, create a YAML file with the following content:

```yaml
apiVersion: neurallog.io/v1
kind: Tenant
metadata:
  name: my-tenant
spec:
  displayName: My Tenant
  description: My first NeuralLog tenant
```

Apply the YAML file to your cluster:

```bash
kubectl apply -f my-tenant.yaml
```

This will create a tenant with default configuration for the server and Redis components.

### Customized Tenant

For more control over the tenant configuration, you can specify additional parameters:

```yaml
apiVersion: neurallog.io/v1
kind: Tenant
metadata:
  name: my-tenant
spec:
  displayName: My Tenant
  description: My customized NeuralLog tenant
  
  # Server configuration
  server:
    replicas: 2
    image: neurallog/server:latest
    resources:
      cpu:
        request: 200m
        limit: 500m
      memory:
        request: 256Mi
        limit: 512Mi
    env:
      - name: LOG_LEVEL
        value: debug
  
  # Redis configuration
  redis:
    replicas: 1
    image: redis:7-alpine
    resources:
      cpu:
        request: 100m
        limit: 300m
      memory:
        request: 128Mi
        limit: 256Mi
    storage: 5Gi
    config:
      maxmemory-policy: allkeys-lru
  
  # Network policy configuration
  networkPolicy:
    enabled: true
    allowedNamespaces:
      - default
      - monitoring
```

Apply the YAML file to your cluster:

```bash
kubectl apply -f my-customized-tenant.yaml
```

## Viewing Tenants

### Listing Tenants

To list all tenants in the cluster:

```bash
kubectl get tenants
```

Example output:

```
NAME        DISPLAY NAME   NAMESPACE        PHASE     AGE
my-tenant   My Tenant      tenant-my-tenant Running   5m
```

### Viewing Tenant Details

To view detailed information about a tenant:

```bash
kubectl describe tenant my-tenant
```

Example output:

```
Name:         my-tenant
Namespace:    
Labels:       <none>
Annotations:  <none>
API Version:  neurallog.io/v1
Kind:         Tenant
Metadata:
  Creation Timestamp:  2023-07-15T12:34:56Z
  Generation:          1
  Resource Version:    12345678
  UID:                 abcdef12-3456-7890-abcd-ef1234567890
Spec:
  Description:  My first NeuralLog tenant
  Display Name: My Tenant
  Network Policy:
    Enabled:  true
  Redis:
    Image:     redis:7-alpine
    Replicas:  1
    Storage:   1Gi
  Server:
    Image:     neurallog/server:latest
    Replicas:  1
Status:
  Namespace:  tenant-my-tenant
  Phase:      Running
  Redis Status:
    Message:        Redis is running
    Phase:          Running
    Ready Replicas: 1
    Total Replicas: 1
  Server Status:
    Message:        Server is running
    Phase:          Running
    Ready Replicas: 1
    Total Replicas: 1
Events:
  Type    Reason    Age   From              Message
  ----    ------    ----  ----              -------
  Normal  Created   5m    tenant-controller  Created namespace tenant-my-tenant
  Normal  Created   5m    tenant-controller  Created Redis ConfigMap
  Normal  Created   5m    tenant-controller  Created Redis Service
  Normal  Created   5m    tenant-controller  Created Redis StatefulSet
  Normal  Created   4m    tenant-controller  Created Server Service
  Normal  Created   4m    tenant-controller  Created Server Deployment
  Normal  Running   3m    tenant-controller  Tenant is running
```

### Viewing Tenant Resources

To view the resources created for a tenant, first get the tenant's namespace:

```bash
NAMESPACE=$(kubectl get tenant my-tenant -o jsonpath='{.status.namespace}')
```

Then, you can list the resources in that namespace:

```bash
kubectl -n $NAMESPACE get all
```

Example output:

```
NAME                                READY   STATUS    RESTARTS   AGE
pod/my-tenant-redis-0               1/1     Running   0          5m
pod/my-tenant-server-6b7f8c9d4-xyz  1/1     Running   0          4m

NAME                   TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)    AGE
service/my-tenant-redis   ClusterIP   None             <none>        6379/TCP   5m
service/my-tenant-server  ClusterIP   10.96.123.456    <none>        3030/TCP   4m

NAME                           READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/my-tenant-server   1/1     1            1           4m

NAME                                      DESIRED   CURRENT   READY   AGE
replicaset.apps/my-tenant-server-6b7f8c9d4   1         1         1       4m

NAME                              READY   AGE
statefulset.apps/my-tenant-redis   1/1     5m
```

## Updating a Tenant

To update a tenant, modify the YAML file and apply it again:

```bash
kubectl apply -f my-tenant.yaml
```

The operator will reconcile the changes and update the tenant resources accordingly.

### Scaling the Server

To scale the server deployment, update the `replicas` field in the `server` section:

```yaml
spec:
  server:
    replicas: 3
```

Apply the updated YAML file:

```bash
kubectl apply -f my-tenant.yaml
```

### Changing Resource Limits

To change resource limits, update the `resources` section:

```yaml
spec:
  server:
    resources:
      cpu:
        request: 300m
        limit: 600m
      memory:
        request: 384Mi
        limit: 768Mi
```

Apply the updated YAML file:

```bash
kubectl apply -f my-tenant.yaml
```

### Adding Environment Variables

To add environment variables to the server, update the `env` section:

```yaml
spec:
  server:
    env:
      - name: LOG_LEVEL
        value: debug
      - name: FEATURE_FLAG_ENABLE_X
        value: "true"
```

Apply the updated YAML file:

```bash
kubectl apply -f my-tenant.yaml
```

## Accessing Tenant Services

### Port Forwarding

To access a tenant's server from your local machine, you can use port forwarding:

```bash
NAMESPACE=$(kubectl get tenant my-tenant -o jsonpath='{.status.namespace}')
kubectl -n $NAMESPACE port-forward svc/my-tenant-server 3030:3030
```

Then, you can access the server at `http://localhost:3030`.

### Creating an Ingress

To expose a tenant's server to the internet, you can create an Ingress resource:

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: my-tenant-ingress
  namespace: tenant-my-tenant
spec:
  rules:
  - host: my-tenant.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: my-tenant-server
            port:
              number: 3030
```

Apply the Ingress resource:

```bash
kubectl apply -f my-tenant-ingress.yaml
```

## Monitoring a Tenant

### Viewing Logs

To view the logs for a tenant's server:

```bash
NAMESPACE=$(kubectl get tenant my-tenant -o jsonpath='{.status.namespace}')
kubectl -n $NAMESPACE logs deployment/my-tenant-server
```

To view the logs for a tenant's Redis instance:

```bash
NAMESPACE=$(kubectl get tenant my-tenant -o jsonpath='{.status.namespace}')
kubectl -n $NAMESPACE logs statefulset/my-tenant-redis
```

### Monitoring Resource Usage

To monitor resource usage for a tenant, you can use the Kubernetes Dashboard or a monitoring solution like Prometheus and Grafana.

## Deleting a Tenant

To delete a tenant and all its resources:

```bash
kubectl delete tenant my-tenant
```

This will delete the tenant's namespace and all resources within it, as well as remove the tenant from the Auth service.

## Troubleshooting

### Tenant Creation Issues

If a tenant is not being created properly, check the tenant status:

```bash
kubectl describe tenant my-tenant
```

Look for error messages in the `Status` section and the `Events` section.

### Pod Issues

If a tenant's pods are not starting properly, check the pod status:

```bash
NAMESPACE=$(kubectl get tenant my-tenant -o jsonpath='{.status.namespace}')
kubectl -n $NAMESPACE describe pod
```

Look for error messages in the `Events` section.

### Operator Issues

If the operator is not functioning properly, check the operator logs:

```bash
kubectl -n neurallog-system logs deployment/neurallog-operator-controller-manager -c manager
```

## Best Practices

### Resource Management

- Set appropriate resource requests and limits for server and Redis components
- Monitor resource usage and adjust limits as needed
- Use horizontal pod autoscaling for server components when possible

### Security

- Enable network policies to isolate tenant resources
- Use secrets for sensitive configuration values
- Regularly update server and Redis images to include security patches

### Backup and Recovery

- Regularly backup tenant data
- Test recovery procedures
- Consider using Redis persistence for critical data

## Advanced Topics

### Custom Redis Configuration

You can customize the Redis configuration by adding key-value pairs to the `config` section:

```yaml
spec:
  redis:
    config:
      maxmemory-policy: allkeys-lru
      appendonly: "yes"
      appendfsync: everysec
```

### Custom Network Policies

You can define custom network policies for a tenant:

```yaml
spec:
  networkPolicy:
    enabled: true
    allowedNamespaces:
      - default
      - monitoring
    ingressRules:
      - description: Allow ingress from app namespace
        from:
          app: web-frontend
        ports:
          - protocol: TCP
            port: 80
    egressRules:
      - description: Allow egress to database
        to:
          app: database
        ports:
          - protocol: TCP
            port: 5432
```

### Using Environment Variables from ConfigMaps and Secrets

You can use environment variables from ConfigMaps and Secrets:

```yaml
spec:
  server:
    env:
      - name: API_KEY
        valueFrom:
          secretKeyRef:
            name: api-keys
            key: my-tenant
      - name: CONFIG_VALUE
        valueFrom:
          configMapKeyRef:
            name: tenant-config
            key: config-value
```

## Conclusion

The NeuralLog Tenant Operator provides a powerful and flexible way to manage tenant resources in the NeuralLog platform. By following this guide, you should be able to create, manage, and troubleshoot tenants effectively.
