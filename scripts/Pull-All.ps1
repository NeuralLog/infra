# Pull All Script
# This script pulls the latest changes for all NeuralLog repositories

# Ensure we stop on errors
$ErrorActionPreference = "Stop"

# Get the root directory of the NeuralLog project
$scriptPath = Split-Path -Parent $MyInvocation.MyCommand.Path
$rootDir = Split-Path -Parent $scriptPath
$neurallogDir = Split-Path -Parent $rootDir

$repos = @("server", "web", "auth", "shared", "specs", "docs", "infra")

foreach ($repo in $repos) {
    $repoDir = Join-Path $neurallogDir $repo
    Write-Host "Pulling $repo..." -ForegroundColor Green
    Push-Location -Path $repoDir
    git pull
    Pop-Location
    Write-Host ""
}
