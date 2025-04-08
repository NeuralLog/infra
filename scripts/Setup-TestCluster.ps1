# PowerShell script to set up a test Kubernetes cluster using kind
# This script creates a kind cluster and deploys NeuralLog to it

# Ensure we stop on errors
$ErrorActionPreference = "Stop"

# Configuration
$ClusterName = "neurallog-test"
$Namespace = "neurallog"

Write-Host "Setting up NeuralLog test Kubernetes cluster..." -ForegroundColor Green

# Navigate to the infra directory
$scriptPath = Split-Path -Parent $MyInvocation.MyCommand.Path
$infraPath = Split-Path -Parent $scriptPath
Set-Location $infraPath

# Check if Docker is running
try {
    $dockerStatus = docker info 2>&1
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Error: Docker is not running. Please start Docker Desktop and try again." -ForegroundColor Red
        exit 1
    }
} catch {
    Write-Host "Error: Docker is not installed or not running. Please install Docker Desktop and try again." -ForegroundColor Red
    exit 1
}

# Check if kind is installed
try {
    $kindVersion = kind version 2>&1
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Error: kind is not installed. Please install kind and try again." -ForegroundColor Red
        Write-Host "You can install kind from: https://kind.sigs.k8s.io/docs/user/quick-start/#installation" -ForegroundColor Yellow
        exit 1
    }
} catch {
    Write-Host "Error: kind is not installed. Please install kind and try again." -ForegroundColor Red
    Write-Host "You can install kind from: https://kind.sigs.k8s.io/docs/user/quick-start/#installation" -ForegroundColor Yellow
    exit 1
}

# Check if kubectl is installed
try {
    $kubectlVersion = kubectl version --client 2>&1
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Error: kubectl is not installed. Please install kubectl and try again." -ForegroundColor Red
        Write-Host "You can install kubectl from: https://kubernetes.io/docs/tasks/tools/install-kubectl/" -ForegroundColor Yellow
        exit 1
    }
} catch {
    Write-Host "Error: kubectl is not installed. Please install kubectl and try again." -ForegroundColor Red
    Write-Host "You can install kubectl from: https://kubernetes.io/docs/tasks/tools/install-kubectl/" -ForegroundColor Yellow
    exit 1
}

# Create the kind cluster if it doesn't exist
$existingClusters = kind get clusters 2>&1
if ($existingClusters -notcontains $ClusterName) {
    Write-Host "Creating kind cluster: $ClusterName" -ForegroundColor Yellow
    
    # Create a config file for the kind cluster
    $kindConfigPath = Join-Path $infraPath "kind-config.yaml"
    @"
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  extraPortMappings:
  - containerPort: 30030
    hostPort: 30030
    protocol: TCP
"@ | Out-File -FilePath $kindConfigPath -Encoding utf8
    
    # Create the cluster with the config
    try {
        kind create cluster --name $ClusterName --config $kindConfigPath
        if ($LASTEXITCODE -ne 0) {
            Write-Host "Error: Failed to create kind cluster." -ForegroundColor Red
            exit 1
        }
    } catch {
        Write-Host "Error: Failed to create kind cluster: $_" -ForegroundColor Red
        exit 1
    }
    
    Write-Host "Kind cluster $ClusterName created successfully" -ForegroundColor Green
} else {
    Write-Host "Kind cluster $ClusterName already exists" -ForegroundColor Yellow
}

# Create the namespace if it doesn't exist
$existingNamespaces = kubectl get namespace 2>&1
if ($existingNamespaces -notmatch $Namespace) {
    Write-Host "Creating namespace: $Namespace" -ForegroundColor Yellow
    try {
        kubectl create namespace $Namespace
        if ($LASTEXITCODE -ne 0) {
            Write-Host "Error: Failed to create namespace." -ForegroundColor Red
            exit 1
        }
    } catch {
        Write-Host "Error: Failed to create namespace: $_" -ForegroundColor Red
        exit 1
    }
} else {
    Write-Host "Namespace $Namespace already exists" -ForegroundColor Yellow
}

# Apply the Kubernetes configurations
Write-Host "Applying Kubernetes configurations..." -ForegroundColor Yellow
try {
    $kubernetesPath = Join-Path $infraPath "kubernetes\overlays\test"
    kubectl apply -k $kubernetesPath
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Error: Failed to apply Kubernetes configurations." -ForegroundColor Red
        exit 1
    }
} catch {
    Write-Host "Error: Failed to apply Kubernetes configurations: $_" -ForegroundColor Red
    exit 1
}

# Wait for deployments to be ready
Write-Host "Waiting for deployments to be ready..." -ForegroundColor Yellow
try {
    kubectl -n $Namespace wait --for=condition=available --timeout=300s deployment/neurallog-server
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Warning: Deployment may not be ready yet. Check status with: kubectl get pods -n $Namespace" -ForegroundColor Yellow
    }
} catch {
    Write-Host "Warning: Deployment may not be ready yet. Check status with: kubectl get pods -n $Namespace" -ForegroundColor Yellow
}

Write-Host "NeuralLog test environment is ready!" -ForegroundColor Green
Write-Host "To access the server, run:" -ForegroundColor Cyan
Write-Host "kubectl port-forward -n $Namespace svc/neurallog-server 3030:3030" -ForegroundColor Cyan
Write-Host "Then access http://localhost:3030" -ForegroundColor Cyan
Write-Host ""
Write-Host "To clean up the test environment, run: .\Cleanup-TestCluster.ps1" -ForegroundColor Yellow
