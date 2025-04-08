#!/bin/bash
# Script to set up a test Kubernetes cluster using kind

# Exit on error
set -e

# Configuration
CLUSTER_NAME="neurallog-test"
NAMESPACE="neurallog"

# Create the kind cluster if it doesn't exist
if ! kind get clusters | grep -q "$CLUSTER_NAME"; then
  echo "Creating kind cluster: $CLUSTER_NAME"
  
  # Create a config file for the kind cluster
  cat <<EOF > kind-config.yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  extraPortMappings:
  - containerPort: 30030
    hostPort: 30030
    protocol: TCP
EOF

  # Create the cluster with the config
  kind create cluster --name "$CLUSTER_NAME" --config kind-config.yaml
  
  echo "Kind cluster $CLUSTER_NAME created successfully"
else
  echo "Kind cluster $CLUSTER_NAME already exists"
fi

# Create the namespace if it doesn't exist
if ! kubectl get namespace | grep -q "$NAMESPACE"; then
  echo "Creating namespace: $NAMESPACE"
  kubectl create namespace "$NAMESPACE"
else
  echo "Namespace $NAMESPACE already exists"
fi

# Apply the Kubernetes configurations
echo "Applying Kubernetes configurations"
kubectl apply -k ../kubernetes/overlays/test

# Wait for deployments to be ready
echo "Waiting for deployments to be ready"
kubectl -n "$NAMESPACE" wait --for=condition=available --timeout=300s deployment/neurallog-server

echo "NeuralLog test environment is ready!"
echo "To access the server, run:"
echo "kubectl port-forward -n $NAMESPACE svc/neurallog-server 3030:3030"
echo "Then access http://localhost:3030"
