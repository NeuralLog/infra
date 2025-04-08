# PowerShell script to initialize the NeuralLog development environment
# This script checks prerequisites and sets up the initial development environment

param (
    [switch]$Force = $false
)

# Ensure we stop on errors
$ErrorActionPreference = "Stop"

Write-Host "Initializing NeuralLog development environment..." -ForegroundColor Green

# Navigate to the infra directory
$scriptPath = Split-Path -Parent $MyInvocation.MyCommand.Path
$infraPath = Split-Path -Parent $scriptPath
$neurallogPath = Split-Path -Parent $infraPath
Set-Location $infraPath

# Check prerequisites
$prerequisites = @{
    "Docker" = {
        try {
            $dockerVersion = docker --version 2>&1
            if ($LASTEXITCODE -ne 0) { return $false }
            return $true
        } catch {
            return $false
        }
    }
    "Docker Compose" = {
        try {
            $dockerComposeVersion = docker-compose --version 2>&1
            if ($LASTEXITCODE -ne 0) { return $false }
            return $true
        } catch {
            return $false
        }
    }
    "Git" = {
        try {
            $gitVersion = git --version 2>&1
            if ($LASTEXITCODE -ne 0) { return $false }
            return $true
        } catch {
            return $false
        }
    }
    "Node.js" = {
        try {
            $nodeVersion = node --version 2>&1
            if ($LASTEXITCODE -ne 0) { return $false }
            return $true
        } catch {
            return $false
        }
    }
}

$allPrerequisitesMet = $true
foreach ($prereq in $prerequisites.Keys) {
    $check = & $prerequisites[$prereq]
    if ($check) {
        Write-Host "✓ $prereq is installed" -ForegroundColor Green
    } else {
        Write-Host "✗ $prereq is not installed" -ForegroundColor Red
        $allPrerequisitesMet = $false
    }
}

if (-not $allPrerequisitesMet) {
    Write-Host "Please install all prerequisites and try again." -ForegroundColor Yellow
    if (-not $Force) {
        exit 1
    } else {
        Write-Host "Continuing anyway due to -Force parameter..." -ForegroundColor Yellow
    }
}

# Check if Docker is running
try {
    $dockerStatus = docker info 2>&1
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Error: Docker is not running. Please start Docker Desktop and try again." -ForegroundColor Red
        if (-not $Force) {
            exit 1
        } else {
            Write-Host "Continuing anyway due to -Force parameter..." -ForegroundColor Yellow
        }
    } else {
        Write-Host "✓ Docker is running" -ForegroundColor Green
    }
} catch {
    Write-Host "Error: Docker is not running. Please start Docker Desktop and try again." -ForegroundColor Red
    if (-not $Force) {
        exit 1
    } else {
        Write-Host "Continuing anyway due to -Force parameter..." -ForegroundColor Yellow
    }
}

# Check if server repository exists
$serverPath = Join-Path $neurallogPath "server"
if (Test-Path $serverPath) {
    Write-Host "✓ NeuralLog/server repository exists" -ForegroundColor Green
} else {
    Write-Host "✗ NeuralLog/server repository not found" -ForegroundColor Red
    Write-Host "Would you like to clone the server repository? (Y/N)" -ForegroundColor Yellow
    $response = Read-Host
    if ($response -eq "Y" -or $response -eq "y") {
        try {
            Set-Location $neurallogPath
            git clone https://github.com/NeuralLog/server.git
            if ($LASTEXITCODE -ne 0) {
                Write-Host "Error: Failed to clone server repository." -ForegroundColor Red
                exit 1
            }
            Write-Host "✓ NeuralLog/server repository cloned successfully" -ForegroundColor Green
        } catch {
            Write-Host "Error: Failed to clone server repository: $_" -ForegroundColor Red
            exit 1
        }
    } else {
        Write-Host "Please clone the server repository manually and try again." -ForegroundColor Yellow
        exit 1
    }
}

# Create Redis data directory if it doesn't exist
$redisDataPath = Join-Path $infraPath "redis\data"
if (-not (Test-Path $redisDataPath)) {
    New-Item -ItemType Directory -Path $redisDataPath -Force | Out-Null
    Write-Host "✓ Created Redis data directory" -ForegroundColor Green
}

# Pull required Docker images
Write-Host "Pulling required Docker images..." -ForegroundColor Yellow
try {
    docker pull redis:7-alpine
    docker pull rediscommander/redis-commander:latest
    docker pull node:18-alpine
    Write-Host "✓ Docker images pulled successfully" -ForegroundColor Green
} catch {
    Write-Host "Warning: Failed to pull some Docker images. They will be pulled when starting the environment." -ForegroundColor Yellow
}

Write-Host "NeuralLog development environment initialized successfully!" -ForegroundColor Green
Write-Host ""
Write-Host "Next steps:" -ForegroundColor Cyan
Write-Host "1. Start the development environment: .\Start-DevEnvironment.ps1" -ForegroundColor Cyan
Write-Host "2. Build the server Docker image: .\Build-ServerImage.ps1" -ForegroundColor Cyan
Write-Host "3. Set up a test Kubernetes cluster: .\Setup-TestCluster.ps1" -ForegroundColor Cyan
