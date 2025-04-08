#!/bin/bash
# Script to stop the development environment

# Navigate to the docker/dev directory
cd ../docker/dev

# Stop the development environment
echo "Stopping development environment..."
docker-compose -f docker-compose.dev.yml down

echo "Development environment stopped successfully"
