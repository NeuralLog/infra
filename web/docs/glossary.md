# NeuralLog Admin Glossary

This glossary provides definitions for terms used in the NeuralLog Admin interface and documentation.

## A

### API Server
The central management point for the Kubernetes cluster. It exposes the Kubernetes API that the NeuralLog Admin interface uses to manage tenants.

### Authentication
The process of verifying the identity of a user or system. The NeuralLog Admin interface uses Kubernetes authentication mechanisms.

## C

### Container
A lightweight, standalone, executable software package that includes everything needed to run an application: code, runtime, system tools, system libraries, and settings.

### CRD (Custom Resource Definition)
A Kubernetes extension mechanism that allows you to define custom resources. The NeuralLog Admin interface uses a Tenant CRD to define and manage tenants.

### Custom Resource
An extension of the Kubernetes API that represents a customization of a particular Kubernetes installation. The Tenant resource is a custom resource used by the NeuralLog Admin interface.

## D

### Deployment
A Kubernetes resource that manages a replicated application. The NeuralLog server is deployed as a Kubernetes Deployment for each tenant.

### Docker
A platform for developing, shipping, and running applications in containers.

## E

### Environment Variable
A variable that is part of the environment in which a process runs. Environment variables can be used to configure the NeuralLog server and Redis.

## K

### Kubernetes
An open-source container orchestration platform for automating deployment, scaling, and management of containerized applications.

### Kubernetes API
The API that serves as the foundation of the Kubernetes control plane. The NeuralLog Admin interface interacts with the Kubernetes API to manage tenants.

## N

### Namespace
A Kubernetes concept for organizing resources in a cluster. Each tenant in NeuralLog has its own namespace for isolation.

### Network Policy
A Kubernetes resource that specifies how groups of pods are allowed to communicate with each other and other network endpoints. Network policies can be enabled for tenants to restrict communication.

### Next.js
A React framework used to build the NeuralLog Admin web interface.

### NeuralLog
A platform for managing and analyzing logs using neural networks and machine learning.

### NeuralLog Admin
The web interface for managing NeuralLog tenants in a Kubernetes cluster.

### NeuralLog Server
The server component of NeuralLog that processes and analyzes logs.

## O

### Operator
A Kubernetes extension that uses custom resources to manage applications and their components. The Tenant Operator manages the lifecycle of NeuralLog tenants.

## P

### Pod
The smallest deployable unit in Kubernetes. A pod represents a single instance of a running process in a cluster.

### Persistent Volume (PV)
A piece of storage in the cluster that has been provisioned by an administrator or dynamically provisioned using Storage Classes. Used by Redis for persistent storage.

### Persistent Volume Claim (PVC)
A request for storage by a user. Redis uses PVCs to request storage for its data.

## R

### RBAC (Role-Based Access Control)
A method of regulating access to resources based on the roles of individual users. Kubernetes RBAC is used to control access to the Kubernetes API.

### Redis
An in-memory data structure store used as a database, cache, and message broker. Each tenant has its own Redis instance for data storage.

### Replica
A copy of a pod that is running at any given time. Multiple replicas provide redundancy and load balancing.

### Resource Limits
Constraints on the amount of CPU and memory that a container can use. Resource limits can be configured for the NeuralLog server and Redis.

### Resource Requests
The amount of CPU and memory that a container is guaranteed to get. Resource requests can be configured for the NeuralLog server and Redis.

## S

### Service
A Kubernetes resource that defines a logical set of pods and a policy by which to access them. Services are used to expose the NeuralLog server and Redis to other components.

### StatefulSet
A Kubernetes resource for managing stateful applications. Redis is deployed as a StatefulSet for each tenant to ensure stable storage.

### Storage Class
A Kubernetes resource that describes the "classes" of storage offered by the cluster. Used to dynamically provision persistent volumes for Redis.

## T

### Tenant
A logical isolation unit in NeuralLog. Each tenant has its own namespace, server, and Redis instance.

### Tenant Operator
A Kubernetes operator that manages the lifecycle of NeuralLog tenants.

### Tenant Resource
A custom Kubernetes resource that defines a NeuralLog tenant.

## Additional Resources

For more information on Kubernetes concepts, see the [Kubernetes Documentation](https://kubernetes.io/docs/home/).

For more information on NeuralLog, see the [NeuralLog Documentation](https://github.com/NeuralLog/infra/tree/main/docs).
