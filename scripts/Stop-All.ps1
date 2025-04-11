# Stop All Script
# This script stops all NeuralLog components

# Ensure we stop on errors
$ErrorActionPreference = "Stop"

# Get the root directory of the NeuralLog project
$scriptPath = Split-Path -Parent $MyInvocation.MyCommand.Path
$rootDir = Split-Path -Parent $scriptPath

Write-Host "Stopping NeuralLog components..." -ForegroundColor Yellow

# Stop and remove containers
Write-Host "Stopping and removing containers..." -ForegroundColor Yellow
docker-compose -f $rootDir/docker-compose.web.yml down
docker-compose -f $rootDir/docker-compose.server.yml down
docker-compose -f $rootDir/docker-compose.auth.yml down

Write-Host "All components stopped!" -ForegroundColor Green
