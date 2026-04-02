#!/bin/bash

# Database Backup Script for Agent-Todo Platform
# Creates compressed, timestamped backups of the PostgreSQL database

set -e

# Configuration
BACKUP_DIR="/home/deploy/backups"
DB_CONTAINER="agent-todo-db"
DB_NAME="agent_todo"
DB_USER="agent_todo"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="${BACKUP_DIR}/agent-todo-${TIMESTAMP}.sql.gz"
RETENTION_DAYS=30
RETENTION_WEEKS=4
RETENTION_MONTHS=12

# Logging
LOG_FILE="${BACKUP_DIR}/backup.log"
log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1" | tee -a "$LOG_FILE"
}

# Create backup directory if it doesn't exist
mkdir -p "$BACKUP_DIR"

log "Starting database backup..."

# Create backup using pg_dump
if docker exec "$DB_CONTAINER" pg_dump -U "$DB_USER" "$DB_NAME" | gzip > "$BACKUP_FILE"; then
    log "✅ Backup created successfully: $BACKUP_FILE"
    
    # Get backup size
    BACKUP_SIZE=$(du -h "$BACKUP_FILE" | cut -f1)
    log "Backup size: $BACKUP_SIZE"
    
    # Verify backup integrity
    if gzip -t "$BACKUP_FILE" 2>/dev/null; then
        log "✅ Backup integrity verified"
    else
        log "❌ Backup integrity check failed!"
        exit 1
    fi
else
    log "❌ Backup failed!"
    exit 1
fi

# Cleanup old backups (keep last 30 days)
log "Cleaning up old backups..."
find "$BACKUP_DIR" -name "agent-todo-*.sql.gz" -type f -mtime +$RETENTION_DAYS -delete
log "✅ Old backups cleaned up"

# List current backups
log "Current backups:"
ls -lh "$BACKUP_DIR"/agent-todo-*.sql.gz 2>/dev/null | tail -10

log "Backup process completed successfully!"

# TODO: Add off-site upload (S3, Backblaze B2, etc.)
# Example for S3:
# aws s3 cp "$BACKUP_FILE" s3://your-bucket/backups/

# TODO: Add notification (email, Slack, etc.)
# Example for Slack webhook:
# curl -X POST -H 'Content-type: application/json' \
#   --data "{\"text\":\"✅ Agent-Todo backup completed: $BACKUP_SIZE\"}" \
#   YOUR_WEBHOOK_URL

exit 0
