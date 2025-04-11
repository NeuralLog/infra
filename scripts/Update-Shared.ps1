# Update Shared Package Script
# This script updates the shared package in all repositories

param (
    [Parameter(Mandatory=$false)]
    [string]$Version = "latest"
)

# Ensure we stop on errors
$ErrorActionPreference = "Stop"

# Get the root directory of the NeuralLog project
$scriptPath = Split-Path -Parent $MyInvocation.MyCommand.Path
$rootDir = Split-Path -Parent $scriptPath
$neurallogDir = Split-Path -Parent $rootDir

# Configure npm to use the private registry for @neurallog scope
Write-Host "Configuring npm to use private registry for @neurallog scope..." -ForegroundColor Yellow
npm config set @neurallog:registry http://localhost:4873

# Update the shared package in all repositories
$repos = @("server", "web", "auth")

foreach ($repo in $repos) {
    $repoDir = Join-Path $neurallogDir $repo
    Write-Host "Updating shared package in $repo..." -ForegroundColor Green
    Push-Location -Path $repoDir
    npm install @neurallog/shared@$Version --registry http://localhost:4873
    Pop-Location
    Write-Host ""
}

Write-Host "Shared package updated in all repositories!" -ForegroundColor Green
