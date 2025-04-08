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
