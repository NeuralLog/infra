#!/bin/bash
# Script to start the development environment

# Exit on error
set -e

# Navigate to the docker/dev directory
cd ../docker/dev

# Start the development environment
echo "Starting development environment..."
docker-compose -f docker-compose.dev.yml up -d

echo "Development environment started successfully"
echo "Server is accessible at: http://localhost:3030"
echo "Redis Commander is accessible at: http://localhost:8081"
echo ""
echo "To view logs, run: docker-compose -f docker-compose.dev.yml logs -f"
echo "To stop the environment, run: docker-compose -f docker-compose.dev.yml down"
