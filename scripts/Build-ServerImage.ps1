# PowerShell script to build and push the NeuralLog server Docker image
# This script builds the server Docker image and optionally pushes it to Docker Hub

param (
    [string]$Tag = "latest",
    [switch]$Push = $false
)

# Ensure we stop on errors
$ErrorActionPreference = "Stop"

# Configuration
$ImageName = "neurallog/server"
$FullImageName = "${ImageName}:${Tag}"

Write-Host "Building NeuralLog server Docker image: $FullImageName" -ForegroundColor Green

# Navigate to the server directory
$scriptPath = Split-Path -Parent $MyInvocation.MyCommand.Path
$serverPath = Join-Path (Split-Path -Parent (Split-Path -Parent $scriptPath)) "server"
$dockerfilePath = Join-Path (Split-Path -Parent $scriptPath) "docker\server\Dockerfile"

# Check if server directory exists
if (-not (Test-Path $serverPath)) {
    Write-Host "Error: Server directory not found at $serverPath" -ForegroundColor Red
    Write-Host "Please ensure the NeuralLog/server repository is cloned alongside the infra repository." -ForegroundColor Yellow
    exit 1
}

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

# Build the Docker image
Write-Host "Building Docker image..." -ForegroundColor Yellow
Set-Location $serverPath
try {
    docker build -t $FullImageName -f $dockerfilePath .
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Error: Failed to build Docker image." -ForegroundColor Red
        exit 1
    }
} catch {
    Write-Host "Error: Failed to build Docker image: $_" -ForegroundColor Red
    exit 1
}

Write-Host "Docker image built successfully: $FullImageName" -ForegroundColor Green

# Push the image if requested
if ($Push) {
    Write-Host "Pushing Docker image to Docker Hub..." -ForegroundColor Yellow
    try {
        docker push $FullImageName
        if ($LASTEXITCODE -ne 0) {
            Write-Host "Error: Failed to push Docker image." -ForegroundColor Red
            exit 1
        }
    } catch {
        Write-Host "Error: Failed to push Docker image: $_" -ForegroundColor Red
        exit 1
    }
    Write-Host "Docker image pushed successfully: $FullImageName" -ForegroundColor Green
} else {
    Write-Host "Skipping push to Docker Hub. Use -Push parameter to push the image." -ForegroundColor Yellow
}

Write-Host "Build process completed successfully!" -ForegroundColor Green
