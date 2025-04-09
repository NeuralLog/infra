# NeuralLog Docker Configuration Guide

This guide provides detailed information about the Docker configurations used in the NeuralLog infrastructure.

## Table of Contents

- [Overview](#overview)
- [Directory Structure](#directory-structure)
- [Server Docker Configuration](#server-docker-configuration)
  - [Production Dockerfile](#production-dockerfile)
  - [Development Dockerfile](#development-dockerfile)
- [Development Environment](#development-environment)
  - [Docker Compose Configuration](#docker-compose-configuration)
  - [Environment Variables](#environment-variables)
  - [Volumes](#volumes)
  - [Networks](#networks)
- [Building Docker Images](#building-docker-images)
- [Running Docker Containers](#running-docker-containers)
- [Docker Best Practices](#docker-best-practices)
- [Troubleshooting](#troubleshooting)
- [Advanced Configuration](#advanced-configuration)

## Overview

NeuralLog uses Docker for containerization of its components. Docker provides a consistent environment for development, testing, and production.

### Key Features

- **Multi-Stage Builds**: Efficient and secure Docker images
- **Development Environment**: Docker Compose for local development
- **Production Images**: Optimized images for production
- **Security**: Non-root users and minimal images
- **Health Checks**: Built-in health checks for monitoring

## Directory Structure

```
docker/
├── server/               # Server Docker configurations
│   └── Dockerfile        # Production Dockerfile for the server
└── dev/                  # Development Docker configurations
    ├── Dockerfile.dev    # Development Dockerfile for the server
    └── docker-compose.dev.yml  # Docker Compose for development
```

## Server Docker Configuration

### Production Dockerfile

The production Dockerfile for the server uses a multi-stage build for efficiency and security:

```dockerfile
# Build stage
FROM node:18-alpine AS builder

WORKDIR /app

# Copy package files
COPY package.json package-lock.json ./

# Install dependencies
RUN npm ci

# Copy source code
COPY . .

# Build the application
RUN npm run build

# Production stage
FROM node:18-alpine

WORKDIR /app

# Copy package files
COPY package.json package-lock.json ./

# Install production dependencies
RUN npm ci --production

# Copy built files
COPY --from=builder /app/dist ./dist

# Set environment variables
ENV NODE_ENV=production
ENV PORT=3030

# Create a non-root user
RUN addgroup -g 1001 -S nodejs && \
    adduser -S -u 1001 -G nodejs nodejs

# Set ownership
RUN chown -R nodejs:nodejs /app

# Switch to non-root user
USER nodejs

# Expose port
EXPOSE 3030

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
  CMD wget -qO- http://localhost:3030/health || exit 1

# Run the application
CMD ["node", "dist/index.js"]
```

#### Key Features

1. **Multi-Stage Build**: Separates build and runtime environments
2. **Minimal Base Image**: Uses Alpine Linux for a small footprint
3. **Non-Root User**: Runs as a non-root user for security
4. **Health Check**: Includes a health check for monitoring
5. **Environment Variables**: Configurable through environment variables

### Development Dockerfile

The development Dockerfile provides a development environment with hot reloading:

```dockerfile
# Development Dockerfile
FROM node:18-alpine

WORKDIR /app

# Install development dependencies
COPY package.json package-lock.json ./
RUN npm install

# Set environment variables
ENV NODE_ENV=development
ENV PORT=3030

# Expose port
EXPOSE 3030

# Run the application in development mode
CMD ["npm", "run", "dev"]
```

#### Key Features

1. **Development Dependencies**: Includes all dependencies for development
2. **Hot Reloading**: Supports hot reloading for faster development
3. **Development Environment**: Configured for development

## Development Environment

### Docker Compose Configuration

The Docker Compose configuration provides a complete development environment:

```yaml
version: '3.8'

services:
  server:
    build:
      context: ../../../server
      dockerfile: ../infra/docker/dev/Dockerfile.dev
    ports:
      - "3030:3030"
    volumes:
      - ../../../server:/app
      - server_node_modules:/app/node_modules
    environment:
      - NODE_ENV=development
      - PORT=3030
      - REDIS_URL=redis://redis:6379
      - LOG_LEVEL=debug
    depends_on:
      - redis
    command: npm run dev

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis-data:/data
      - ../../redis/conf/redis.conf:/usr/local/etc/redis/redis.conf
    command: ["redis-server", "/usr/local/etc/redis/redis.conf"]

  redis-commander:
    image: rediscommander/redis-commander:latest
    ports:
      - "8081:8081"
    environment:
      - REDIS_HOSTS=local:redis:6379
    depends_on:
      - redis

volumes:
  server_node_modules:
  redis-data:
```

#### Key Features

1. **Multiple Services**: Includes server, Redis, and Redis Commander
2. **Volume Mounts**: Mounts source code for hot reloading
3. **Environment Variables**: Configurable through environment variables
4. **Port Mapping**: Maps container ports to host ports
5. **Dependencies**: Defines service dependencies

### Environment Variables

The development environment uses the following environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `NODE_ENV` | Node.js environment | development |
| `PORT` | Server port | 3030 |
| `REDIS_URL` | Redis connection URL | redis://redis:6379 |
| `LOG_LEVEL` | Logging level | debug |

### Volumes

The development environment uses the following volumes:

| Volume | Description |
|--------|-------------|
| `server_node_modules` | Node.js modules for the server |
| `redis-data` | Redis data for persistence |

### Networks

The development environment uses the default Docker Compose network, which allows services to communicate with each other using service names as hostnames.

## Building Docker Images

### Building the Server Image

```bash
# Navigate to the server directory
cd server

# Build the production image
docker build -t neurallog/server:latest -f ../infra/docker/server/Dockerfile .

# Build the development image
docker build -t neurallog/server:dev -f ../infra/docker/dev/Dockerfile.dev .
```

### Using the Build Script

```bash
# Build the server image
./scripts/build-server-image.sh [tag]

# Build and push the server image
./scripts/build-server-image.sh [tag] -Push
```

## Running Docker Containers

### Running the Server Container

```bash
# Run the production server
docker run -p 3030:3030 -e REDIS_URL=redis://redis:6379 neurallog/server:latest

# Run the development server
docker run -p 3030:3030 -v $(pwd):/app -e REDIS_URL=redis://redis:6379 neurallog/server:dev
```

### Running the Development Environment

```bash
# Start the development environment
cd infra/docker/dev
docker-compose -f docker-compose.dev.yml up -d

# View logs
docker-compose -f docker-compose.dev.yml logs -f

# Stop the development environment
docker-compose -f docker-compose.dev.yml down
```

### Using the Development Scripts

```bash
# Start the development environment
./scripts/start-dev-env.sh

# Stop the development environment
./scripts/stop-dev-env.sh
```

## Docker Best Practices

### Security

1. **Non-Root User**: Run containers as non-root users
2. **Minimal Base Images**: Use minimal base images like Alpine
3. **Multi-Stage Builds**: Use multi-stage builds to minimize image size
4. **Secrets Management**: Use Docker secrets or environment variables for sensitive information

### Performance

1. **Layer Caching**: Optimize Dockerfiles for layer caching
2. **Minimal Dependencies**: Include only necessary dependencies
3. **Resource Limits**: Set resource limits for containers

### Maintainability

1. **Documentation**: Document Dockerfile and Docker Compose configurations
2. **Version Pinning**: Pin versions of base images and dependencies
3. **Environment Variables**: Use environment variables for configuration

## Troubleshooting

### Common Issues

#### Container Not Starting

If a container is not starting, check the logs:

```bash
docker logs <container_id>
```

#### Volume Mount Issues

If volume mounts are not working, check the volume paths:

```bash
docker volume ls
docker volume inspect <volume_name>
```

#### Network Issues

If containers cannot communicate, check the network:

```bash
docker network ls
docker network inspect <network_name>
```

#### Permission Issues

If permission issues occur, check the user and file permissions:

```bash
docker exec -it <container_id> sh
ls -la /app
```

## Advanced Configuration

### Custom Entrypoint

```dockerfile
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
```

### Docker Compose Profiles

```yaml
services:
  server:
    profiles: ["dev", "test"]
    # ...

  redis:
    profiles: ["dev", "test"]
    # ...

  redis-commander:
    profiles: ["dev"]
    # ...
```

### Docker Compose Extensions

```yaml
x-common-variables: &common-variables
  NODE_ENV: development
  LOG_LEVEL: debug

services:
  server:
    environment:
      <<: *common-variables
      PORT: 3030
      REDIS_URL: redis://redis:6379
    # ...
```
