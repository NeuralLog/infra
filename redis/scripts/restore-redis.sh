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
