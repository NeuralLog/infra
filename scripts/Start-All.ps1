# Start All Script
# This script starts all NeuralLog components using Docker Compose

# Ensure we stop on errors
$ErrorActionPreference = "Stop"

# Get the root directory of the NeuralLog project
$scriptPath = Split-Path -Parent $MyInvocation.MyCommand.Path
$rootDir = Split-Path -Parent $scriptPath

Write-Host "Starting NeuralLog components..." -ForegroundColor Green

# Create the network if it doesn't exist
$networkExists = docker network ls | Select-String "neurallog-network"
if (-not $networkExists) {
    Write-Host "Creating neurallog-network..." -ForegroundColor Yellow
    docker network create neurallog-network
}

# Start Verdaccio
Write-Host "Starting Verdaccio..." -ForegroundColor Yellow
docker-compose -f $rootDir/docker-compose.web.yml up -d verdaccio
Start-Sleep -Seconds 5

# Start Redis
Write-Host "Starting Redis..." -ForegroundColor Yellow
docker-compose -f $rootDir/docker-compose.server.yml up -d redis
Start-Sleep -Seconds 2

# Start PostgreSQL and OpenFGA
Write-Host "Starting PostgreSQL and OpenFGA..." -ForegroundColor Yellow
docker-compose -f $rootDir/docker-compose.auth.yml up -d postgres openfga-migrate openfga
Start-Sleep -Seconds 5

# Start Auth Service
Write-Host "Starting Auth Service..." -ForegroundColor Yellow
docker-compose -f $rootDir/docker-compose.auth.yml up -d auth
Start-Sleep -Seconds 5

# Start Logs Server
Write-Host "Starting Logs Server..." -ForegroundColor Yellow
docker-compose -f $rootDir/docker-compose.server.yml up -d server

Write-Host "All components started successfully!" -ForegroundColor Green
Write-Host "You can now run the web application locally with:" -ForegroundColor Green
Write-Host "cd web && npx next dev" -ForegroundColor Cyan
Write-Host "Or start the web container with:" -ForegroundColor Green
Write-Host "docker-compose -f $rootDir/docker-compose.web.yml up -d web" -ForegroundColor Cyan
