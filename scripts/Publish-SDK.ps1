# Publish SDK Script
# This script builds and publishes the SDK to the private registry

# Ensure we stop on errors
$ErrorActionPreference = "Stop"

# Get the root directory of the NeuralLog project
$scriptPath = Split-Path -Parent $MyInvocation.MyCommand.Path
$rootDir = Split-Path -Parent $scriptPath
$neurallogDir = Split-Path -Parent $rootDir
$sdkDir = Join-Path $neurallogDir "typescript"

# Start Verdaccio if it's not already running
$verdaccioRunning = docker ps | Select-String "verdaccio"
if (-not $verdaccioRunning) {
    Write-Host "Starting Verdaccio..." -ForegroundColor Yellow
    docker-compose -f $rootDir/docker-compose.web.yml up -d verdaccio
    Start-Sleep -Seconds 5
}

# Configure npm to use the private registry for @neurallog scope
Write-Host "Configuring npm to use private registry for @neurallog scope..." -ForegroundColor Yellow
npm config set @neurallog:registry http://localhost:4873

# Build and publish the SDK
Write-Host "Building and publishing the SDK..." -ForegroundColor Yellow
Push-Location -Path $sdkDir

# Check if user is already logged in to Verdaccio
$npmWhoami = npm whoami --registry http://localhost:4873 2>&1
if ($LASTEXITCODE -ne 0) {
    Write-Host "Logging in to Verdaccio..." -ForegroundColor Yellow
    Write-Host "Use username: admin, password: admin" -ForegroundColor Cyan
    npm adduser --registry http://localhost:4873 --auth-type=legacy
}

# Install dependencies
npm install

# Install the latest shared package
Write-Host "Installing latest @neurallog/shared package..." -ForegroundColor Yellow
npm install @neurallog/shared@latest --registry http://localhost:4873

# Build the SDK
npm run build

# Publish the SDK
npm publish --registry http://localhost:4873

Pop-Location

Write-Host "SDK published successfully!" -ForegroundColor Green
Write-Host "You can now install it in other repositories with:" -ForegroundColor Green
Write-Host "npm install @neurallog/sdk --registry http://localhost:4873" -ForegroundColor Cyan
