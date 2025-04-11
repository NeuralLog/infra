# Push All Script
# This script commits and pushes changes for all NeuralLog repositories

param (
    [Parameter(Mandatory=$true)]
    [string]$CommitMessage
)

# Ensure we stop on errors
$ErrorActionPreference = "Stop"

# Get the root directory of the NeuralLog project
$scriptPath = Split-Path -Parent $MyInvocation.MyCommand.Path
$rootDir = Split-Path -Parent $scriptPath
$neurallogDir = Split-Path -Parent $rootDir

$repos = @("server", "web", "auth", "shared", "specs", "docs", "infra")

foreach ($repo in $repos) {
    $repoDir = Join-Path $neurallogDir $repo
    Write-Host "Pushing $repo..." -ForegroundColor Green
    Push-Location -Path $repoDir
    git add .
    git commit -m $CommitMessage
    git push
    Pop-Location
    Write-Host ""
}
