# PowerShell script to clean up the test Kubernetes cluster
# This script deletes the kind cluster and cleans up any temporary files

# Ensure we stop on errors
$ErrorActionPreference = "Stop"

# Configuration
$ClusterName = "neurallog-test"

Write-Host "Cleaning up NeuralLog test Kubernetes cluster..." -ForegroundColor Yellow

# Navigate to the infra directory
$scriptPath = Split-Path -Parent $MyInvocation.MyCommand.Path
$infraPath = Split-Path -Parent $scriptPath
Set-Location $infraPath

# Check if kind is installed
try {
    $kindVersion = kind version 2>&1
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Error: kind is not installed. Nothing to clean up." -ForegroundColor Red
        exit 0
    }
} catch {
    Write-Host "Error: kind is not installed. Nothing to clean up." -ForegroundColor Red
    exit 0
}

# Check if the cluster exists
$existingClusters = kind get clusters 2>&1
if ($existingClusters -contains $ClusterName) {
    Write-Host "Deleting kind cluster: $ClusterName" -ForegroundColor Yellow
    try {
        kind delete cluster --name $ClusterName
        if ($LASTEXITCODE -ne 0) {
            Write-Host "Error: Failed to delete kind cluster." -ForegroundColor Red
            exit 1
        }
    } catch {
        Write-Host "Error: Failed to delete kind cluster: $_" -ForegroundColor Red
        exit 1
    }
    Write-Host "Kind cluster $ClusterName deleted successfully" -ForegroundColor Green
} else {
    Write-Host "Kind cluster $ClusterName does not exist. Nothing to clean up." -ForegroundColor Yellow
}

# Remove any temporary files
$kindConfigPath = Join-Path $infraPath "kind-config.yaml"
if (Test-Path $kindConfigPath) {
    Remove-Item $kindConfigPath -Force
    Write-Host "Removed temporary kind configuration file" -ForegroundColor Green
}

Write-Host "Cleanup completed successfully!" -ForegroundColor Green
