# NeuralLog Admin Troubleshooting Guide

This guide provides solutions for common issues you might encounter when using the NeuralLog Admin interface and managing tenants.

## Table of Contents

1. [Installation Issues](#installation-issues)
2. [Connection Issues](#connection-issues)
3. [Tenant Management Issues](#tenant-management-issues)
4. [Kubernetes Integration Issues](#kubernetes-integration-issues)
5. [Performance Issues](#performance-issues)
6. [Collecting Diagnostic Information](#collecting-diagnostic-information)

## Installation Issues

### Web Interface Won't Start

**Symptoms:**
- Error when running `npm start`
- Blank page when accessing the web interface

**Possible Causes and Solutions:**

1. **Node.js Version Mismatch**
   - **Cause:** The application requires Node.js v16 or higher.
   - **Solution:** Upgrade Node.js to the latest LTS version.
   ```bash
   # Check your Node.js version
   node --version
   
   # Install nvm (Node Version Manager) if needed
   # Then install and use the latest LTS version
   nvm install --lts
   nvm use --lts
   ```

2. **Missing Dependencies**
   - **Cause:** Not all dependencies were installed correctly.
   - **Solution:** Reinstall dependencies.
   ```bash
   npm ci
   ```

3. **Port Conflict**
   - **Cause:** Another application is using port 3000.
   - **Solution:** Change the port in the `.env.local` file.
   ```
   PORT=3001
   ```

### Build Errors

**Symptoms:**
- Errors during `npm run build`
- TypeScript compilation errors

**Possible Causes and Solutions:**

1. **TypeScript Errors**
   - **Cause:** Type errors in the codebase.
   - **Solution:** Fix the type errors or temporarily disable strict type checking.
   ```bash
   # Check for type errors
   npx tsc --noEmit
   ```

2. **Outdated Dependencies**
   - **Cause:** Dependencies are outdated or incompatible.
   - **Solution:** Update dependencies.
   ```bash
   npm update
   ```

## Connection Issues

### Cannot Connect to Kubernetes API

**Symptoms:**
- "Failed to fetch tenants" error
- Empty tenant list
- API errors in the browser console

**Possible Causes and Solutions:**

1. **Kubernetes API Server Unreachable**
   - **Cause:** The Kubernetes API server is not accessible.
   - **Solution:** Verify the API server is running and accessible.
   ```bash
   kubectl cluster-info
   ```

2. **RBAC Permissions**
   - **Cause:** Insufficient permissions to access Kubernetes resources.
   - **Solution:** Ensure the service account has the necessary permissions.
   ```bash
   # Check RBAC permissions
   kubectl auth can-i list tenants.neurallog.io --as=system:serviceaccount:default:neurallog-operator
   ```

3. **Incorrect API URL**
   - **Cause:** The Kubernetes API URL is incorrect in the configuration.
   - **Solution:** Update the API URL in the `.env.local` file.
   ```
   KUBERNETES_API_URL=https://your-kubernetes-api-server:6443
   ```

4. **CORS Issues**
   - **Cause:** CORS policies preventing API access.
   - **Solution:** Configure CORS headers on the API server or use a proxy.

## Tenant Management Issues

### Tenant Creation Fails

**Symptoms:**
- Error message when creating a tenant
- Tenant appears in the list but with "Failed" status

**Possible Causes and Solutions:**

1. **Invalid Tenant Name**
   - **Cause:** The tenant name contains invalid characters or is already in use.
   - **Solution:** Use a unique name with only lowercase alphanumeric characters and hyphens.

2. **Resource Constraints**
   - **Cause:** The cluster doesn't have enough resources to create the tenant.
   - **Solution:** Reduce the resource requests or add more resources to the cluster.
   ```bash
   # Check available resources
   kubectl describe nodes | grep -A 5 "Allocated resources"
   ```

3. **Operator Issues**
   - **Cause:** The Tenant Operator is not running or has errors.
   - **Solution:** Check the operator logs.
   ```bash
   kubectl logs -l app=neurallog-operator -n <operator-namespace>
   ```

### Tenant Status Shows "Degraded"

**Symptoms:**
- Tenant status shows "Degraded"
- Some components are not running correctly

**Possible Causes and Solutions:**

1. **Pod Failures**
   - **Cause:** Pods are failing to start or are crashing.
   - **Solution:** Check the pod logs.
   ```bash
   kubectl logs -l app=neurallog-server -n tenant-<tenant-name>
   kubectl logs -l app=redis -n tenant-<tenant-name>
   ```

2. **Resource Limits**
   - **Cause:** Pods are hitting resource limits.
   - **Solution:** Increase resource limits or optimize the application.
   ```bash
   # Check resource usage
   kubectl top pods -n tenant-<tenant-name>
   ```

3. **Storage Issues**
   - **Cause:** Persistent volume claims are not being fulfilled.
   - **Solution:** Check the PVC status and storage class.
   ```bash
   kubectl get pvc -n tenant-<tenant-name>
   ```

## Kubernetes Integration Issues

### CRD Not Found

**Symptoms:**
- "No matches for kind 'Tenant' in version 'neurallog.io/v1'" error
- Tenant creation fails with API errors

**Possible Causes and Solutions:**

1. **CRD Not Installed**
   - **Cause:** The Tenant CRD is not installed in the cluster.
   - **Solution:** Install the CRD.
   ```bash
   kubectl apply -f ../operator/config/crd/bases/neurallog.io_tenants.yaml
   ```

2. **API Group Mismatch**
   - **Cause:** The API group in the code doesn't match the CRD.
   - **Solution:** Ensure the API group in the code matches the CRD.
   ```bash
   # Check the API group in the CRD
   kubectl get crd tenants.neurallog.io -o jsonpath='{.spec.group}'
   ```

### Operator Not Reconciling

**Symptoms:**
- Tenants are created but not reconciled
- No resources are created for the tenant

**Possible Causes and Solutions:**

1. **Operator Not Running**
   - **Cause:** The Tenant Operator is not running.
   - **Solution:** Check the operator deployment.
   ```bash
   kubectl get pods -l app=neurallog-operator -n <operator-namespace>
   ```

2. **Operator Errors**
   - **Cause:** The operator is encountering errors during reconciliation.
   - **Solution:** Check the operator logs.
   ```bash
   kubectl logs -l app=neurallog-operator -n <operator-namespace>
   ```

3. **RBAC Issues**
   - **Cause:** The operator doesn't have permission to create resources.
   - **Solution:** Ensure the operator's service account has the necessary permissions.
   ```bash
   kubectl get clusterrole neurallog-operator -o yaml
   kubectl get clusterrolebinding neurallog-operator -o yaml
   ```

## Performance Issues

### Slow UI Response

**Symptoms:**
- UI is slow to load or respond
- Operations take a long time to complete

**Possible Causes and Solutions:**

1. **Large Number of Tenants**
   - **Cause:** The system is managing a large number of tenants.
   - **Solution:** Implement pagination or filtering in the UI.

2. **Network Latency**
   - **Cause:** High network latency between the UI and the Kubernetes API.
   - **Solution:** Deploy the UI closer to the Kubernetes API or optimize API calls.

3. **Resource Constraints**
   - **Cause:** The UI or API server is resource-constrained.
   - **Solution:** Increase resources for the UI or API server.

### High Resource Usage

**Symptoms:**
- High CPU or memory usage
- Pods being OOMKilled

**Possible Causes and Solutions:**

1. **Inefficient Queries**
   - **Cause:** Inefficient queries to the Kubernetes API.
   - **Solution:** Optimize queries and implement caching.

2. **Memory Leaks**
   - **Cause:** Memory leaks in the application.
   - **Solution:** Identify and fix memory leaks.
   ```bash
   # Check memory usage
   kubectl top pods -l app=neurallog-admin
   ```

3. **Insufficient Resources**
   - **Cause:** Not enough resources allocated to the pods.
   - **Solution:** Increase resource requests and limits.

## Collecting Diagnostic Information

When troubleshooting issues, it's helpful to collect diagnostic information to share with support:

### System Information

```bash
# Kubernetes version
kubectl version

# Node information
kubectl describe nodes

# Cluster events
kubectl get events --sort-by='.lastTimestamp'
```

### Tenant Information

```bash
# List all tenants
kubectl get tenants.neurallog.io

# Describe a specific tenant
kubectl describe tenant <tenant-name>

# Get tenant status
kubectl get tenant <tenant-name> -o jsonpath='{.status}'
```

### Operator Information

```bash
# Operator pods
kubectl get pods -l app=neurallog-operator -n <operator-namespace>

# Operator logs
kubectl logs -l app=neurallog-operator -n <operator-namespace>

# Operator events
kubectl get events -n <operator-namespace> --field-selector involvedObject.name=neurallog-operator
```

### Tenant Resources

```bash
# Namespace resources
kubectl get all -n tenant-<tenant-name>

# Pod logs
kubectl logs -l app=neurallog-server -n tenant-<tenant-name>
kubectl logs -l app=redis -n tenant-<tenant-name>

# Persistent volumes
kubectl get pv,pvc -n tenant-<tenant-name>
```

### Web Interface Logs

```bash
# Check Next.js logs
npm run start > nextjs.log 2>&1
```

## Getting Additional Help

If you're still experiencing issues after trying the solutions in this guide, please:

1. Check the [GitHub Issues](https://github.com/NeuralLog/infra/issues) for similar problems and solutions.
2. Create a new issue with detailed information about your problem, including:
   - Steps to reproduce
   - Error messages
   - Diagnostic information
   - Environment details (Kubernetes version, browser, etc.)
3. Contact the NeuralLog support team for assistance.
