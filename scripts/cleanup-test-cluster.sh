#!/bin/bash
# Script to clean up the test Kubernetes cluster

# Configuration
CLUSTER_NAME="neurallog-test"

# Check if the cluster exists
if kind get clusters | grep -q "$CLUSTER_NAME"; then
  echo "Deleting kind cluster: $CLUSTER_NAME"
  kind delete cluster --name "$CLUSTER_NAME"
  echo "Kind cluster $CLUSTER_NAME deleted successfully"
else
  echo "Kind cluster $CLUSTER_NAME does not exist"
fi

# Remove any temporary files
if [ -f kind-config.yaml ]; then
  rm kind-config.yaml
fi

echo "Cleanup completed"
