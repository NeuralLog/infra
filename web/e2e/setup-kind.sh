#!/bin/bash
set -e

# Create a Kind cluster if it doesn't exist
if ! kind get clusters | grep -q "neurallog-e2e"; then
  echo "Creating Kind cluster 'neurallog-e2e'..."
  kind create cluster --name neurallog-e2e --config - <<EOF
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  extraPortMappings:
  - containerPort: 30000
    hostPort: 30000
    protocol: TCP
EOF
else
  echo "Kind cluster 'neurallog-e2e' already exists."
fi

# Set kubectl context to the Kind cluster
kubectl config use-context kind-neurallog-e2e

# Install the Tenant CRD
echo "Installing Tenant CRD..."
kubectl apply -f ../operator/config/crd/bases/neurallog.io_tenants.yaml

# Create a test namespace
echo "Creating test namespace..."
kubectl create namespace neurallog-test --dry-run=client -o yaml | kubectl apply -f -

# Build and load the operator image into Kind
echo "Building and loading operator image..."
docker build -t neurallog/operator:test ../operator
kind load docker-image neurallog/operator:test --name neurallog-e2e

# Deploy the operator
echo "Deploying operator..."
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: ServiceAccount
metadata:
  name: neurallog-operator
  namespace: neurallog-test
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: neurallog-operator
rules:
- apiGroups: [""]
  resources: ["namespaces", "services", "configmaps", "secrets"]
  verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]
- apiGroups: ["apps"]
  resources: ["deployments", "statefulsets"]
  verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]
- apiGroups: ["networking.k8s.io"]
  resources: ["networkpolicies"]
  verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]
- apiGroups: ["neurallog.io"]
  resources: ["tenants"]
  verbs: ["get", "list", "watch", "create", "update", "patch", "delete", "status"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: neurallog-operator
subjects:
- kind: ServiceAccount
  name: neurallog-operator
  namespace: neurallog-test
roleRef:
  kind: ClusterRole
  name: neurallog-operator
  apiGroup: rbac.authorization.k8s.io
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: neurallog-operator
  namespace: neurallog-test
spec:
  replicas: 1
  selector:
    matchLabels:
      app: neurallog-operator
  template:
    metadata:
      labels:
        app: neurallog-operator
    spec:
      serviceAccountName: neurallog-operator
      containers:
      - name: operator
        image: neurallog/operator:test
        imagePullPolicy: IfNotPresent
        env:
        - name: WATCH_NAMESPACE
          value: ""
        - name: POD_NAME
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
        - name: OPERATOR_NAME
          value: "neurallog-operator"
EOF

# Create a NodePort service for the API
echo "Creating API service..."
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Service
metadata:
  name: neurallog-api
  namespace: neurallog-test
spec:
  type: NodePort
  selector:
    app: neurallog-operator
  ports:
  - port: 8080
    targetPort: 8080
    nodePort: 30000
EOF

echo "Setup complete!"
