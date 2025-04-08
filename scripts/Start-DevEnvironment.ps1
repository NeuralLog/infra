# PowerShell script to start the NeuralLog development environment
# This script sets up and starts the development environment using Docker Compose

# Ensure we stop on errors
$ErrorActionPreference = "Stop"

Write-Host "Setting up NeuralLog development environment..." -ForegroundColor Green

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

# Start the development environment
Write-Host "Starting development environment..." -ForegroundColor Yellow
try {
    docker-compose -f docker-compose.dev.yml up -d
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Error: Failed to start development environment." -ForegroundColor Red
        exit 1
    }
} catch {
    Write-Host "Error: Failed to start development environment: $_" -ForegroundColor Red
    exit 1
}

Write-Host "Development environment started successfully!" -ForegroundColor Green
Write-Host "Server is accessible at: http://localhost:3030" -ForegroundColor Cyan
Write-Host "Redis Commander is accessible at: http://localhost:8081" -ForegroundColor Cyan
Write-Host ""
Write-Host "To view logs, run: docker-compose -f docker-compose.dev.yml logs -f" -ForegroundColor Yellow
Write-Host "To stop the environment, run: .\Stop-DevEnvironment.ps1" -ForegroundColor Yellow
