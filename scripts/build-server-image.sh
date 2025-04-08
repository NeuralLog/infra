#!/bin/bash
# Script to build and push the server Docker image

# Exit on error
set -e

# Configuration
IMAGE_NAME="neurallog/server"
TAG=${1:-"latest"}

# Navigate to the server directory
cd ../../server

# Build the Docker image
echo "Building Docker image: $IMAGE_NAME:$TAG"
docker build -t "$IMAGE_NAME:$TAG" -f ../infra/docker/server/Dockerfile .

echo "Docker image built successfully: $IMAGE_NAME:$TAG"

# Ask if the image should be pushed
read -p "Do you want to push the image to Docker Hub? (y/n): " PUSH_IMAGE

if [[ "$PUSH_IMAGE" =~ ^[Yy]$ ]]; then
  echo "Pushing Docker image to Docker Hub..."
  docker push "$IMAGE_NAME:$TAG"
  echo "Docker image pushed successfully: $IMAGE_NAME:$TAG"
else
  echo "Skipping push to Docker Hub"
fi

echo "Build process completed"
