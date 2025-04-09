# NeuralLog Troubleshooting Guide

This guide provides solutions for common issues that may arise when working with the NeuralLog infrastructure.

## Table of Contents

- [Development Environment Issues](#development-environment-issues)
- [Docker Issues](#docker-issues)
- [Kubernetes Issues](#kubernetes-issues)
- [Redis Issues](#redis-issues)
- [Server Issues](#server-issues)
- [Operator Issues](#operator-issues)
- [Network Issues](#network-issues)
- [Storage Issues](#storage-issues)
- [Performance Issues](#performance-issues)
- [Security Issues](#security-issues)
- [Logging and Monitoring](#logging-and-monitoring)
- [Common Error Messages](#common-error-messages)

## Development Environment Issues

### Docker Not Running

**Symptoms:**
- Error message: "Cannot connect to the Docker daemon"
- Docker commands fail

**Solutions:**
1. Start Docker Desktop (Windows/macOS)
2. Start Docker service (Linux):
   ```bash
   sudo systemctl start docker
   ```
3. Check Docker status:
   ```bash
   docker info
   ```

### Docker Compose Not Working

**Symptoms:**
- Error message: "docker-compose command not found"
- Docker Compose commands fail

**Solutions:**
1. Install Docker Compose:
   ```bash
   pip install docker-compose
   ```
2. Check Docker Compose version:
   ```bash
   docker-compose --version
   ```

### Development Environment Not Starting

**Symptoms:**
- Error when running `Start-DevEnvironment.ps1` or `start-dev-env.sh`
- Docker Compose errors

**Solutions:**
1. Check Docker is running
2. Check Docker Compose is installed
3. Check for port conflicts:
   ```bash
   netstat -tuln | grep 3030
   netstat -tuln | grep 6379
   netstat -tuln | grep 8081
   ```
4. Stop any conflicting services
5. Try stopping and removing existing containers:
   ```bash
   docker-compose -f docker/dev/docker-compose.dev.yml down
   ```
6. Check logs:
   ```bash
   docker-compose -f docker/dev/docker-compose.dev.yml logs
   ```

## Docker Issues

### Container Not Starting

**Symptoms:**
- Container status is "Created" or "Exited"
- Error in container logs

**Solutions:**
1. Check container logs:
   ```bash
   docker logs <container_id>
   ```
2. Check container status:
   ```bash
   docker inspect <container_id>
   ```
3. Check for resource constraints:
   ```bash
   docker info
   ```
4. Check for port conflicts:
   ```bash
   netstat -tuln
   ```
5. Try running with interactive mode:
   ```bash
   docker run -it --rm <image> sh
   ```

### Image Build Failing

**Symptoms:**
- Error when building Docker image
- Build process fails

**Solutions:**
1. Check Dockerfile syntax
2. Check for missing files
3. Check for network issues
4. Try building with verbose output:
   ```bash
   docker build --progress=plain -t <image> .
   ```
5. Check Docker daemon logs:
   ```bash
   docker system info
   ```

### Volume Mount Issues

**Symptoms:**
- Files not appearing in container
- Changes not persisting

**Solutions:**
1. Check volume mount syntax
2. Check file permissions
3. Check volume exists:
   ```bash
   docker volume ls
   docker volume inspect <volume_name>
   ```
4. Try using absolute paths
5. Check SELinux/AppArmor settings (Linux)

## Kubernetes Issues

### Pod Not Starting

**Symptoms:**
- Pod status is "Pending" or "CrashLoopBackOff"
- Error in pod logs

**Solutions:**
1. Check pod status:
   ```bash
   kubectl get pods -n <namespace>
   ```
2. Check pod details:
   ```bash
   kubectl describe pod <pod_name> -n <namespace>
   ```
3. Check pod logs:
   ```bash
   kubectl logs <pod_name> -n <namespace>
   ```
4. Check for resource constraints:
   ```bash
   kubectl describe node
   ```
5. Check for image pull issues:
   ```bash
   kubectl describe pod <pod_name> -n <namespace> | grep "Image:"
   ```

### Service Not Accessible

**Symptoms:**
- Cannot connect to service
- Service endpoint not responding

**Solutions:**
1. Check service status:
   ```bash
   kubectl get svc -n <namespace>
   ```
2. Check service details:
   ```bash
   kubectl describe svc <service_name> -n <namespace>
   ```
3. Check endpoints:
   ```bash
   kubectl get endpoints <service_name> -n <namespace>
   ```
4. Check pod selector:
   ```bash
   kubectl get pods -l <selector> -n <namespace>
   ```
5. Try port forwarding:
   ```bash
   kubectl port-forward svc/<service_name> <local_port>:<service_port> -n <namespace>
   ```

### PersistentVolumeClaim Not Binding

**Symptoms:**
- PVC status is "Pending"
- Pod cannot start due to unbound PVC

**Solutions:**
1. Check PVC status:
   ```bash
   kubectl get pvc -n <namespace>
   ```
2. Check PVC details:
   ```bash
   kubectl describe pvc <pvc_name> -n <namespace>
   ```
3. Check storage class:
   ```bash
   kubectl get storageclass
   ```
4. Check for available PVs:
   ```bash
   kubectl get pv
   ```
5. Check storage provisioner:
   ```bash
   kubectl get storageclass <storage_class> -o yaml
   ```

### Kind Cluster Issues

**Symptoms:**
- Kind cluster not starting
- Kind cluster not accessible

**Solutions:**
1. Check Docker is running
2. Check kind version:
   ```bash
   kind version
   ```
3. Check for existing clusters:
   ```bash
   kind get clusters
   ```
4. Delete and recreate cluster:
   ```bash
   kind delete cluster --name <cluster_name>
   kind create cluster --name <cluster_name>
   ```
5. Check kind logs:
   ```bash
   kind export logs --name <cluster_name>
   ```

## Redis Issues

### Redis Not Starting

**Symptoms:**
- Redis container not starting
- Redis pod not starting

**Solutions:**
1. Check Redis logs:
   ```bash
   # Docker
   docker logs <redis_container_id>
   
   # Kubernetes
   kubectl logs <redis_pod_name> -n <namespace>
   ```
2. Check Redis configuration:
   ```bash
   # Docker
   docker exec <redis_container_id> cat /etc/redis/redis.conf
   
   # Kubernetes
   kubectl exec <redis_pod_name> -n <namespace> -- cat /etc/redis/redis.conf
   ```
3. Try starting Redis with default configuration:
   ```bash
   # Docker
   docker run -it --rm redis:7-alpine redis-server
   
   # Kubernetes
   kubectl run redis-test --image=redis:7-alpine -n <namespace>
   ```
4. Check for permission issues:
   ```bash
   # Docker
   docker exec <redis_container_id> ls -la /data
   
   # Kubernetes
   kubectl exec <redis_pod_name> -n <namespace> -- ls -la /data
   ```

### Redis Connection Issues

**Symptoms:**
- Server cannot connect to Redis
- Redis connection errors

**Solutions:**
1. Check Redis is running:
   ```bash
   # Docker
   docker ps | grep redis
   
   # Kubernetes
   kubectl get pods -l app=redis -n <namespace>
   ```
2. Check Redis service:
   ```bash
   # Kubernetes
   kubectl get svc redis -n <namespace>
   ```
3. Test Redis connection:
   ```bash
   # Docker
   docker exec <server_container_id> redis-cli -h redis ping
   
   # Kubernetes
   kubectl exec <server_pod_name> -n <namespace> -- redis-cli -h redis ping
   ```
4. Check Redis URL:
   ```bash
   # Docker
   docker exec <server_container_id> env | grep REDIS_URL
   
   # Kubernetes
   kubectl exec <server_pod_name> -n <namespace> -- env | grep REDIS_URL
   ```
5. Check network connectivity:
   ```bash
   # Docker
   docker exec <server_container_id> ping redis
   
   # Kubernetes
   kubectl exec <server_pod_name> -n <namespace> -- ping redis
   ```

### Redis Data Loss

**Symptoms:**
- Data not persisting after restart
- Missing data

**Solutions:**
1. Check Redis persistence configuration:
   ```bash
   # Docker
   docker exec <redis_container_id> redis-cli config get appendonly
   
   # Kubernetes
   kubectl exec <redis_pod_name> -n <namespace> -- redis-cli config get appendonly
   ```
2. Check Redis data directory:
   ```bash
   # Docker
   docker exec <redis_container_id> ls -la /data
   
   # Kubernetes
   kubectl exec <redis_pod_name> -n <namespace> -- ls -la /data
   ```
3. Check Redis volume mount:
   ```bash
   # Docker
   docker inspect <redis_container_id> | grep -A 10 Mounts
   
   # Kubernetes
   kubectl describe pod <redis_pod_name> -n <namespace> | grep -A 10 Volumes
   ```
4. Enable Redis persistence:
   ```bash
   # Docker
   docker exec <redis_container_id> redis-cli config set appendonly yes
   
   # Kubernetes
   kubectl exec <redis_pod_name> -n <namespace> -- redis-cli config set appendonly yes
   ```

## Server Issues

### Server Not Starting

**Symptoms:**
- Server container not starting
- Server pod not starting

**Solutions:**
1. Check server logs:
   ```bash
   # Docker
   docker logs <server_container_id>
   
   # Kubernetes
   kubectl logs <server_pod_name> -n <namespace>
   ```
2. Check environment variables:
   ```bash
   # Docker
   docker exec <server_container_id> env
   
   # Kubernetes
   kubectl exec <server_pod_name> -n <namespace> -- env
   ```
3. Check for missing dependencies:
   ```bash
   # Docker
   docker exec <server_container_id> ls -la /app/node_modules
   
   # Kubernetes
   kubectl exec <server_pod_name> -n <namespace> -- ls -la /app/node_modules
   ```
4. Check for file permission issues:
   ```bash
   # Docker
   docker exec <server_container_id> ls -la /app
   
   # Kubernetes
   kubectl exec <server_pod_name> -n <namespace> -- ls -la /app
   ```
5. Try running with debug logging:
   ```bash
   # Docker
   docker run -e LOG_LEVEL=debug -p 3030:3030 neurallog/server:latest
   
   # Kubernetes
   kubectl set env deployment/neurallog-server LOG_LEVEL=debug -n <namespace>
   ```

### Server API Errors

**Symptoms:**
- API endpoints returning errors
- HTTP 500 errors

**Solutions:**
1. Check server logs:
   ```bash
   # Docker
   docker logs <server_container_id>
   
   # Kubernetes
   kubectl logs <server_pod_name> -n <namespace>
   ```
2. Check Redis connection:
   ```bash
   # Docker
   docker exec <server_container_id> redis-cli -h redis ping
   
   # Kubernetes
   kubectl exec <server_pod_name> -n <namespace> -- redis-cli -h redis ping
   ```
3. Check API endpoint:
   ```bash
   # Docker
   curl http://localhost:3030/health
   
   # Kubernetes
   kubectl port-forward svc/neurallog-server 3030:3030 -n <namespace>
   curl http://localhost:3030/health
   ```
4. Check server configuration:
   ```bash
   # Docker
   docker exec <server_container_id> env
   
   # Kubernetes
   kubectl exec <server_pod_name> -n <namespace> -- env
   ```

### Server Performance Issues

**Symptoms:**
- Slow API responses
- High CPU/memory usage

**Solutions:**
1. Check resource usage:
   ```bash
   # Docker
   docker stats <server_container_id>
   
   # Kubernetes
   kubectl top pod <server_pod_name> -n <namespace>
   ```
2. Check Redis performance:
   ```bash
   # Docker
   docker exec <redis_container_id> redis-cli info
   
   # Kubernetes
   kubectl exec <redis_pod_name> -n <namespace> -- redis-cli info
   ```
3. Check for memory leaks:
   ```bash
   # Docker
   docker exec <server_container_id> node --inspect
   
   # Kubernetes
   kubectl port-forward <server_pod_name> 9229:9229 -n <namespace>
   ```
4. Increase resource limits:
   ```bash
   # Docker
   docker run --memory=1g --cpus=1 -p 3030:3030 neurallog/server:latest
   
   # Kubernetes
   kubectl set resources deployment/neurallog-server -c server --limits=cpu=1,memory=1Gi --requests=cpu=500m,memory=512Mi -n <namespace>
   ```

## Operator Issues

### Operator Not Starting

**Symptoms:**
- Operator pod not starting
- Operator pod in CrashLoopBackOff

**Solutions:**
1. Check operator logs:
   ```bash
   kubectl logs -l control-plane=controller-manager -n system
   ```
2. Check operator deployment:
   ```bash
   kubectl describe deployment controller-manager -n system
   ```
3. Check RBAC permissions:
   ```bash
   kubectl get clusterrole manager-role -o yaml
   kubectl get clusterrolebinding manager-rolebinding -o yaml
   ```
4. Check CRDs:
   ```bash
   kubectl get crd tenants.neurallog.io
   ```
5. Try running the operator locally:
   ```bash
   cd operator
   make run
   ```

### Tenant Not Being Reconciled

**Symptoms:**
- Tenant status not updating
- Tenant resources not being created

**Solutions:**
1. Check tenant status:
   ```bash
   kubectl get tenant <tenant_name>
   kubectl describe tenant <tenant_name>
   ```
2. Check operator logs:
   ```bash
   kubectl logs -l control-plane=controller-manager -n system
   ```
3. Check for RBAC issues:
   ```bash
   kubectl auth can-i create namespace --as=system:serviceaccount:system:controller-manager
   kubectl auth can-i create deployment --as=system:serviceaccount:system:controller-manager -n <namespace>
   ```
4. Check for validation issues:
   ```bash
   kubectl get tenant <tenant_name> -o yaml
   ```
5. Try deleting and recreating the tenant:
   ```bash
   kubectl delete tenant <tenant_name>
   kubectl apply -f <tenant_file>
   ```

### CRD Issues

**Symptoms:**
- CRD not being recognized
- Validation errors when creating tenants

**Solutions:**
1. Check CRD installation:
   ```bash
   kubectl get crd tenants.neurallog.io
   ```
2. Check CRD definition:
   ```bash
   kubectl get crd tenants.neurallog.io -o yaml
   ```
3. Check for validation errors:
   ```bash
   kubectl apply -f <tenant_file> --validate=true
   ```
4. Reinstall CRD:
   ```bash
   kubectl delete crd tenants.neurallog.io
   kubectl apply -f operator/config/crd/bases/neurallog.io_tenants.yaml
   ```
5. Check operator logs for CRD issues:
   ```bash
   kubectl logs -l control-plane=controller-manager -n system
   ```

## Network Issues

### Service Discovery Issues

**Symptoms:**
- Services cannot find each other
- DNS resolution failures

**Solutions:**
1. Check DNS resolution:
   ```bash
   # Docker
   docker exec <container_id> nslookup redis
   
   # Kubernetes
   kubectl exec <pod_name> -n <namespace> -- nslookup redis
   ```
2. Check service endpoints:
   ```bash
   kubectl get endpoints -n <namespace>
   ```
3. Check network policies:
   ```bash
   kubectl get networkpolicies -n <namespace>
   ```
4. Check CoreDNS:
   ```bash
   kubectl get pods -n kube-system -l k8s-app=kube-dns
   kubectl logs -l k8s-app=kube-dns -n kube-system
   ```
5. Try using IP address instead of hostname:
   ```bash
   kubectl get svc -n <namespace> -o wide
   ```

### Network Policy Issues

**Symptoms:**
- Connections being blocked
- Services cannot communicate

**Solutions:**
1. Check network policies:
   ```bash
   kubectl get networkpolicies -n <namespace>
   kubectl describe networkpolicy <policy_name> -n <namespace>
   ```
2. Temporarily disable network policies:
   ```bash
   kubectl delete networkpolicy --all -n <namespace>
   ```
3. Test connectivity:
   ```bash
   kubectl exec <pod_name> -n <namespace> -- curl http://neurallog-server:3030/health
   ```
4. Check pod labels:
   ```bash
   kubectl get pods --show-labels -n <namespace>
   ```
5. Check network policy selectors:
   ```bash
   kubectl get networkpolicy <policy_name> -n <namespace> -o yaml
   ```

## Storage Issues

### PersistentVolume Issues

**Symptoms:**
- PVCs stuck in Pending state
- Pods cannot start due to volume issues

**Solutions:**
1. Check PVC status:
   ```bash
   kubectl get pvc -n <namespace>
   ```
2. Check PV status:
   ```bash
   kubectl get pv
   ```
3. Check storage class:
   ```bash
   kubectl get storageclass
   ```
4. Check for storage provisioner issues:
   ```bash
   kubectl get pods -n kube-system -l app=<provisioner>
   ```
5. Try using a different storage class:
   ```yaml
   apiVersion: v1
   kind: PersistentVolumeClaim
   metadata:
     name: redis-data
   spec:
     accessModes:
       - ReadWriteOnce
     storageClassName: standard
     resources:
       requests:
         storage: 1Gi
   ```

### Volume Mount Issues

**Symptoms:**
- Files not appearing in container
- Permission denied errors

**Solutions:**
1. Check volume mounts:
   ```bash
   kubectl describe pod <pod_name> -n <namespace> | grep -A 10 Volumes
   ```
2. Check file permissions:
   ```bash
   kubectl exec <pod_name> -n <namespace> -- ls -la <mount_path>
   ```
3. Check for SELinux/AppArmor issues (Linux):
   ```bash
   kubectl exec <pod_name> -n <namespace> -- dmesg | grep denied
   ```
4. Try using an init container to set permissions:
   ```yaml
   initContainers:
   - name: set-permissions
     image: busybox
     command: ["sh", "-c", "chmod -R 777 /data"]
     volumeMounts:
     - name: data
       mountPath: /data
   ```
5. Check for volume type issues:
   ```bash
   kubectl describe pv <pv_name>
   ```

## Performance Issues

### High CPU Usage

**Symptoms:**
- High CPU usage
- Slow response times

**Solutions:**
1. Check CPU usage:
   ```bash
   # Docker
   docker stats <container_id>
   
   # Kubernetes
   kubectl top pod <pod_name> -n <namespace>
   ```
2. Check for CPU-intensive processes:
   ```bash
   # Docker
   docker exec <container_id> top -b -n 1
   
   # Kubernetes
   kubectl exec <pod_name> -n <namespace> -- top -b -n 1
   ```
3. Increase CPU limits:
   ```bash
   # Docker
   docker run --cpus=2 -p 3030:3030 neurallog/server:latest
   
   # Kubernetes
   kubectl set resources deployment/neurallog-server -c server --limits=cpu=2 --requests=cpu=1 -n <namespace>
   ```
4. Profile the application:
   ```bash
   # Docker
   docker exec <container_id> node --prof
   
   # Kubernetes
   kubectl exec <pod_name> -n <namespace> -- node --prof
   ```
5. Check for infinite loops or blocking operations

### High Memory Usage

**Symptoms:**
- High memory usage
- Out of memory errors

**Solutions:**
1. Check memory usage:
   ```bash
   # Docker
   docker stats <container_id>
   
   # Kubernetes
   kubectl top pod <pod_name> -n <namespace>
   ```
2. Check for memory leaks:
   ```bash
   # Docker
   docker exec <container_id> node --inspect
   
   # Kubernetes
   kubectl port-forward <pod_name> 9229:9229 -n <namespace>
   ```
3. Increase memory limits:
   ```bash
   # Docker
   docker run --memory=2g -p 3030:3030 neurallog/server:latest
   
   # Kubernetes
   kubectl set resources deployment/neurallog-server -c server --limits=memory=2Gi --requests=memory=1Gi -n <namespace>
   ```
4. Check for memory-intensive operations:
   ```bash
   # Docker
   docker exec <container_id> node --inspect
   
   # Kubernetes
   kubectl exec <pod_name> -n <namespace> -- node --inspect
   ```
5. Enable garbage collection logging:
   ```bash
   # Docker
   docker run -e NODE_OPTIONS="--trace-gc" -p 3030:3030 neurallog/server:latest
   
   # Kubernetes
   kubectl set env deployment/neurallog-server NODE_OPTIONS="--trace-gc" -n <namespace>
   ```

## Security Issues

### Authentication Issues

**Symptoms:**
- Authentication failures
- Unauthorized access

**Solutions:**
1. Check authentication configuration:
   ```bash
   # Docker
   docker exec <container_id> env | grep AUTH
   
   # Kubernetes
   kubectl exec <pod_name> -n <namespace> -- env | grep AUTH
   ```
2. Check for missing secrets:
   ```bash
   kubectl get secrets -n <namespace>
   ```
3. Check for incorrect credentials:
   ```bash
   kubectl describe secret <secret_name> -n <namespace>
   ```
4. Check server logs for authentication errors:
   ```bash
   # Docker
   docker logs <container_id> | grep auth
   
   # Kubernetes
   kubectl logs <pod_name> -n <namespace> | grep auth
   ```
5. Try with explicit credentials:
   ```bash
   curl -u username:password http://localhost:3030/api
   ```

### Authorization Issues

**Symptoms:**
- Permission denied errors
- Forbidden access

**Solutions:**
1. Check RBAC configuration:
   ```bash
   kubectl get role -n <namespace>
   kubectl get rolebinding -n <namespace>
   ```
2. Check service account permissions:
   ```bash
   kubectl auth can-i <verb> <resource> --as=system:serviceaccount:<namespace>:<serviceaccount>
   ```
3. Check for missing role bindings:
   ```bash
   kubectl get rolebinding -n <namespace>
   ```
4. Check for incorrect role definitions:
   ```bash
   kubectl get role <role_name> -n <namespace> -o yaml
   ```
5. Add necessary permissions:
   ```yaml
   apiVersion: rbac.authorization.k8s.io/v1
   kind: Role
   metadata:
     name: <role_name>
     namespace: <namespace>
   rules:
   - apiGroups: [""]
     resources: ["pods"]
     verbs: ["get", "list", "watch"]
   ```

## Logging and Monitoring

### Missing Logs

**Symptoms:**
- Logs not appearing
- Incomplete logs

**Solutions:**
1. Check logging configuration:
   ```bash
   # Docker
   docker exec <container_id> env | grep LOG
   
   # Kubernetes
   kubectl exec <pod_name> -n <namespace> -- env | grep LOG
   ```
2. Increase log level:
   ```bash
   # Docker
   docker run -e LOG_LEVEL=debug -p 3030:3030 neurallog/server:latest
   
   # Kubernetes
   kubectl set env deployment/neurallog-server LOG_LEVEL=debug -n <namespace>
   ```
3. Check for log rotation issues:
   ```bash
   # Docker
   docker exec <container_id> ls -la /var/log
   
   # Kubernetes
   kubectl exec <pod_name> -n <namespace> -- ls -la /var/log
   ```
4. Check for disk space issues:
   ```bash
   # Docker
   docker exec <container_id> df -h
   
   # Kubernetes
   kubectl exec <pod_name> -n <namespace> -- df -h
   ```
5. Use a logging sidecar:
   ```yaml
   containers:
   - name: logging-sidecar
     image: busybox
     command: ["sh", "-c", "tail -f /var/log/app.log"]
     volumeMounts:
     - name: logs
       mountPath: /var/log
   ```

### Monitoring Issues

**Symptoms:**
- Metrics not appearing
- Incomplete metrics

**Solutions:**
1. Check monitoring configuration:
   ```bash
   # Docker
   docker exec <container_id> env | grep METRICS
   
   # Kubernetes
   kubectl exec <pod_name> -n <namespace> -- env | grep METRICS
   ```
2. Check metrics endpoint:
   ```bash
   # Docker
   curl http://localhost:3030/metrics
   
   # Kubernetes
   kubectl port-forward svc/neurallog-server 3030:3030 -n <namespace>
   curl http://localhost:3030/metrics
   ```
3. Check Prometheus configuration:
   ```bash
   kubectl get configmap prometheus-config -n monitoring -o yaml
   ```
4. Check for scrape errors:
   ```bash
   kubectl logs -l app=prometheus -n monitoring
   ```
5. Add Prometheus annotations:
   ```yaml
   metadata:
     annotations:
       prometheus.io/scrape: "true"
       prometheus.io/port: "3030"
       prometheus.io/path: "/metrics"
   ```

## Common Error Messages

### "Error: ECONNREFUSED"

**Cause:** Cannot connect to a service

**Solutions:**
1. Check if the service is running
2. Check if the service is accessible
3. Check network connectivity
4. Check for firewall/network policy issues
5. Check service name and port

### "Error: ENOENT"

**Cause:** File or directory not found

**Solutions:**
1. Check if the file exists
2. Check file permissions
3. Check volume mounts
4. Check working directory
5. Check for typos in file paths

### "Error: ETIMEDOUT"

**Cause:** Connection timed out

**Solutions:**
1. Check if the service is running
2. Check network connectivity
3. Check for firewall/network policy issues
4. Check for service overload
5. Increase timeout settings

### "CrashLoopBackOff"

**Cause:** Pod repeatedly crashing

**Solutions:**
1. Check pod logs
2. Check for missing dependencies
3. Check for configuration errors
4. Check for resource constraints
5. Check for permission issues

### "ImagePullBackOff"

**Cause:** Cannot pull container image

**Solutions:**
1. Check image name and tag
2. Check image registry accessibility
3. Check for authentication issues
4. Check for network connectivity
5. Try pulling the image manually
