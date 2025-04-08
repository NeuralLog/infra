# PowerShell script to stop the NeuralLog development environment
# This script stops the development environment containers

# Ensure we stop on errors
$ErrorActionPreference = "Stop"

Write-Host "Stopping NeuralLog development environment..." -ForegroundColor Yellow

# Navigate to the docker/dev directory
$scriptPath = Split-Path -Parent $MyInvocation.MyCommand.Path
$dockerDevPath = Join-Path (Split-Path -Parent $scriptPath) "docker\dev"
Set-Location $dockerDevPath

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

# Stop the development environment
try {
    docker-compose -f docker-compose.dev.yml down
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Error: Failed to stop development environment." -ForegroundColor Red
        exit 1
    }
} catch {
    Write-Host "Error: Failed to stop development environment: $_" -ForegroundColor Red
    exit 1
}

Write-Host "Development environment stopped successfully!" -ForegroundColor Green
