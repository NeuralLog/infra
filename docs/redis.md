# NeuralLog Redis Configuration Guide

This guide provides detailed information about the Redis configuration used in the NeuralLog infrastructure.

## Table of Contents

- [Overview](#overview)
- [Directory Structure](#directory-structure)
- [Redis Configuration File](#redis-configuration-file)
  - [Basic Configuration](#basic-configuration)
  - [Memory Management](#memory-management)
  - [Persistence](#persistence)
  - [Logging](#logging)
  - [Security](#security)
  - [Performance Tuning](#performance-tuning)
- [Kubernetes Configuration](#kubernetes-configuration)
  - [StatefulSet](#statefulset)
  - [Service](#service)
  - [ConfigMap](#configmap)
  - [PersistentVolumeClaim](#persistentvolumeclaim)
- [Docker Configuration](#docker-configuration)
- [Backup and Restore](#backup-and-restore)
  - [Backup Script](#backup-script)
  - [Restore Script](#restore-script)
- [Monitoring](#monitoring)
- [Troubleshooting](#troubleshooting)
- [Advanced Configuration](#advanced-configuration)

## Overview

Redis is used as the primary data store for NeuralLog. It stores logs, patterns, rules, and other data required by the system.

### Key Features

- **Persistence**: Configured for data durability
- **Memory Management**: Optimized for memory usage
- **Performance**: Tuned for high performance
- **Security**: Configured for secure operation
- **Monitoring**: Includes monitoring capabilities

## Directory Structure

```
redis/
├── conf/                 # Redis configuration files
│   └── redis.conf        # Main Redis configuration file
└── scripts/              # Redis utility scripts
    ├── backup-redis.sh   # Script for backing up Redis data
    └── restore-redis.sh  # Script for restoring Redis data
```

## Redis Configuration File

The Redis configuration file (`redis.conf`) contains the configuration for the Redis instance:

```
# Redis configuration for NeuralLog

# Basic configuration
port 6379
bind 0.0.0.0
protected-mode yes
daemonize no

# Memory management
maxmemory 256mb
maxmemory-policy allkeys-lru

# Persistence
appendonly yes
appendfsync everysec

# Logging
loglevel notice
logfile ""

# Security
# requirepass should be set in production environments
# requirepass yourpassword

# Performance tuning
tcp-keepalive 300
timeout 0
tcp-backlog 511
databases 16

# Snapshotting
save 900 1
save 300 10
save 60 10000

# Advanced configuration
hash-max-ziplist-entries 512
hash-max-ziplist-value 64
list-max-ziplist-size -2
set-max-intset-entries 512
zset-max-ziplist-entries 128
zset-max-ziplist-value 64
hll-sparse-max-bytes 3000
activerehashing yes
client-output-buffer-limit normal 0 0 0
client-output-buffer-limit slave 256mb 64mb 60
client-output-buffer-limit pubsub 32mb 8mb 60
hz 10
aof-rewrite-incremental-fsync yes
```

### Basic Configuration

| Parameter | Description | Default |
|-----------|-------------|---------|
| `port` | Port to listen on | 6379 |
| `bind` | IP addresses to bind to | 0.0.0.0 |
| `protected-mode` | Enable protected mode | yes |
| `daemonize` | Run as a daemon | no |

### Memory Management

| Parameter | Description | Default |
|-----------|-------------|---------|
| `maxmemory` | Maximum memory usage | 256mb |
| `maxmemory-policy` | Eviction policy | allkeys-lru |

#### Eviction Policies

- **allkeys-lru**: Evict any key using LRU algorithm
- **volatile-lru**: Evict keys with expiration using LRU algorithm
- **allkeys-random**: Evict random keys
- **volatile-random**: Evict keys with expiration randomly
- **volatile-ttl**: Evict keys with expiration using TTL
- **noeviction**: Return errors when memory limit is reached

### Persistence

| Parameter | Description | Default |
|-----------|-------------|---------|
| `appendonly` | Enable append-only file | yes |
| `appendfsync` | Fsync policy | everysec |
| `save` | Save points for RDB persistence | 900 1, 300 10, 60 10000 |

#### Fsync Policies

- **always**: Fsync after every write (slow but safe)
- **everysec**: Fsync every second (good compromise)
- **no**: Let the OS handle fsync (fast but unsafe)

### Logging

| Parameter | Description | Default |
|-----------|-------------|---------|
| `loglevel` | Logging level | notice |
| `logfile` | Log file path | "" (stdout) |

#### Logging Levels

- **debug**: Verbose debugging information
- **verbose**: More information than debug
- **notice**: Informational messages
- **warning**: Warning messages

### Security

| Parameter | Description | Default |
|-----------|-------------|---------|
| `requirepass` | Authentication password | (commented out) |

### Performance Tuning

| Parameter | Description | Default |
|-----------|-------------|---------|
| `tcp-keepalive` | TCP keepalive interval | 300 |
| `timeout` | Client timeout | 0 (disabled) |
| `tcp-backlog` | TCP backlog | 511 |
| `databases` | Number of databases | 16 |
| `activerehashing` | Enable active rehashing | yes |
| `hz` | Redis hz value | 10 |

## Kubernetes Configuration

Redis is deployed in Kubernetes using a StatefulSet, Service, ConfigMap, and PersistentVolumeClaim.

### StatefulSet

The Redis StatefulSet defines the Redis container:

```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: redis
  labels:
    app: redis
spec:
  serviceName: redis
  replicas: 1
  selector:
    matchLabels:
      app: redis
  template:
    metadata:
      labels:
        app: redis
    spec:
      containers:
      - name: redis
        image: redis:7-alpine
        command:
        - redis-server
        - /etc/redis/redis.conf
        ports:
        - containerPort: 6379
          name: redis
        volumeMounts:
        - name: redis-data
          mountPath: /data
        - name: redis-config
          mountPath: /etc/redis
        resources:
          requests:
            cpu: 100m
            memory: 128Mi
          limits:
            cpu: 300m
            memory: 256Mi
        livenessProbe:
          tcpSocket:
            port: redis
          initialDelaySeconds: 15
          periodSeconds: 20
        readinessProbe:
          exec:
            command:
            - redis-cli
            - ping
          initialDelaySeconds: 5
          periodSeconds: 10
      volumes:
      - name: redis-config
        configMap:
          name: redis-config
          items:
          - key: redis.conf
            path: redis.conf
  volumeClaimTemplates:
  - metadata:
      name: redis-data
    spec:
      accessModes: [ "ReadWriteOnce" ]
      resources:
        requests:
          storage: 1Gi
```

### Service

The Redis service exposes the Redis instance:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: redis
  labels:
    app: redis
spec:
  selector:
    app: redis
  ports:
  - port: 6379
    targetPort: redis
    name: redis
  clusterIP: None
```

### ConfigMap

The Redis ConfigMap contains the Redis configuration:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: redis-config
  labels:
    app: redis
data:
  redis.conf: |
    # Redis configuration for NeuralLog
    port 6379
    bind 0.0.0.0
    protected-mode yes
    daemonize no
    
    # Memory management
    maxmemory 256mb
    maxmemory-policy allkeys-lru
    
    # Persistence
    appendonly yes
    appendfsync everysec
    
    # Logging
    loglevel notice
    logfile ""
```

### PersistentVolumeClaim

The Redis PersistentVolumeClaim provides persistent storage for Redis data:

```yaml
volumeClaimTemplates:
- metadata:
    name: redis-data
  spec:
    accessModes: [ "ReadWriteOnce" ]
    resources:
      requests:
        storage: 1Gi
```

## Docker Configuration

Redis is also configured for use with Docker Compose:

```yaml
redis:
  image: redis:7-alpine
  ports:
    - "6379:6379"
  volumes:
    - redis-data:/data
    - ../../redis/conf/redis.conf:/usr/local/etc/redis/redis.conf
  command: ["redis-server", "/usr/local/etc/redis/redis.conf"]
```

## Backup and Restore

Redis backup and restore scripts are provided for data protection.

### Backup Script

The backup script (`backup-redis.sh`) creates Redis backups:

```bash
#!/bin/bash
# Script to backup Redis data

# Configuration
BACKUP_DIR="/backups"
TIMESTAMP=$(date +%Y%m%d%H%M%S)
BACKUP_FILE="$BACKUP_DIR/redis-backup-$TIMESTAMP.rdb"

# Create backup directory if it doesn't exist
mkdir -p $BACKUP_DIR

# Trigger Redis SAVE command
redis-cli SAVE

# Copy the RDB file
cp /data/dump.rdb $BACKUP_FILE

echo "Redis backup created: $BACKUP_FILE"
```

### Restore Script

The restore script (`restore-redis.sh`) restores Redis from backups:

```bash
#!/bin/bash
# Script to restore Redis data from backup

# Check if backup file is provided
if [ -z "$1" ]; then
  echo "Usage: $0 <backup-file>"
  exit 1
fi

BACKUP_FILE="$1"

# Check if backup file exists
if [ ! -f "$BACKUP_FILE" ]; then
  echo "Error: Backup file not found: $BACKUP_FILE"
  exit 1
fi

# Stop Redis server
echo "Stopping Redis server..."
redis-cli SHUTDOWN SAVE

# Wait for Redis to stop
sleep 2

# Backup the current dump.rdb file
if [ -f /data/dump.rdb ]; then
  mv /data/dump.rdb /data/dump.rdb.bak
fi

# Copy the backup file to the Redis data directory
echo "Restoring from backup: $BACKUP_FILE"
cp "$BACKUP_FILE" /data/dump.rdb

# Start Redis server
echo "Starting Redis server..."
redis-server /etc/redis/redis.conf &

echo "Redis restore completed successfully."
```

## Monitoring

Redis can be monitored using various tools:

### Redis Commander

Redis Commander is included in the development environment:

```yaml
redis-commander:
  image: rediscommander/redis-commander:latest
  ports:
    - "8081:8081"
  environment:
    - REDIS_HOSTS=local:redis:6379
  depends_on:
    - redis
```

### Redis CLI

Redis CLI can be used to monitor Redis:

```bash
# Connect to Redis
redis-cli

# Monitor Redis
redis-cli monitor

# Get Redis info
redis-cli info

# Get Redis stats
redis-cli info stats
```

### Prometheus and Grafana

Redis can be monitored using Prometheus and Grafana:

1. Install the Redis Prometheus exporter
2. Configure Prometheus to scrape the exporter
3. Import Redis dashboards into Grafana

## Troubleshooting

### Common Issues

#### Redis Not Starting

If Redis is not starting, check the logs:

```bash
# Kubernetes
kubectl logs -l app=redis -n neurallog

# Docker
docker logs <container_id>
```

#### Redis Connection Issues

If the server cannot connect to Redis, check the Redis URL:

```bash
# Test Redis connection
redis-cli -u redis://redis:6379 ping
```

#### Redis Memory Issues

If Redis is running out of memory, check the memory usage:

```bash
redis-cli info memory
```

#### Redis Persistence Issues

If Redis persistence is not working, check the AOF and RDB files:

```bash
ls -la /data
```

## Advanced Configuration

### Redis Sentinel

Redis Sentinel can be used for high availability:

```
sentinel monitor mymaster redis 6379 2
sentinel down-after-milliseconds mymaster 5000
sentinel failover-timeout mymaster 60000
sentinel parallel-syncs mymaster 1
```

### Redis Cluster

Redis Cluster can be used for horizontal scaling:

```
cluster-enabled yes
cluster-config-file nodes.conf
cluster-node-timeout 5000
```

### Redis ACLs

Redis Access Control Lists (ACLs) can be used for fine-grained access control:

```
user default on >password ~* +@all
user readonly on >password ~* +@read
```
