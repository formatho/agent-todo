#!/bin/bash

# Database Restore Script for Agent-Todo Platform
# Restores the PostgreSQL database from a compressed backup file

set -e

# Configuration
DB_CONTAINER="agent-todo-db"
DB_NAME="agent_todo"
DB_USER="agent_todo"
BACKUP_FILE="$1"

# Logging
log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1"
}

# Validate input
if [ -z "$BACKUP_FILE" ]; then
    echo "❌ Error: No backup file specified"
    echo "Usage: $0 <backup-file.sql.gz>"
    echo ""
    echo "Available backups:"
    ls -lht /home/deploy/backups/agent-todo-*.sql.gz 2>/dev/null | head -10
    exit 1
fi

if [ ! -f "$BACKUP_FILE" ]; then
    echo "❌ Error: Backup file not found: $BACKUP_FILE"
    exit 1
fi

log "⚠️  WARNING: This will REPLACE the current database!"
log "Backup file: $BACKUP_FILE"
read -p "Are you sure you want to continue? (yes/no): " confirm

if [ "$confirm" != "yes" ]; then
    log "Restore cancelled"
    exit 0
fi

log "Starting database restore..."

# Verify backup integrity
log "Verifying backup integrity..."
if gzip -t "$BACKUP_FILE" 2>/dev/null; then
    log "✅ Backup integrity verified"
else
    log "❌ Backup file is corrupted!"
    exit 1
fi

# Create a temporary backup of current database
log "Creating safety backup of current database..."
TEMP_BACKUP="/tmp/pre-restore-$(date +%Y%m%d_%H%M%S).sql.gz"
if docker exec "$DB_CONTAINER" pg_dump -U "$DB_USER" "$DB_NAME" | gzip > "$TEMP_BACKUP"; then
    log "✅ Safety backup created: $TEMP_BACKUP"
else
    log "⚠️  Warning: Could not create safety backup (continuing anyway)"
fi

# Stop backend to prevent connections during restore
log "Stopping backend service..."
cd /home/deploy/todo
docker-compose stop backend || true

# Drop existing connections
log "Terminating existing database connections..."
docker exec "$DB_CONTAINER" psql -U "$DB_USER" -d postgres -c \
    "SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = '$DB_NAME' AND pid <> pg_backend_pid();" || true

# Restore the database
log "Restoring database from backup..."
if gunzip -c "$BACKUP_FILE" | docker exec -i "$DB_CONTAINER" psql -U "$DB_USER" "$DB_NAME"; then
    log "✅ Database restored successfully!"
else
    log "❌ Restore failed! Attempting to recover from safety backup..."
    if [ -f "$TEMP_BACKUP" ]; then
        gunzip -c "$TEMP_BACKUP" | docker exec -i "$DB_CONTAINER" psql -U "$DB_USER" "$DB_NAME"
        log "✅ Recovered from safety backup"
    fi
    exit 1
fi

# Restart backend
log "Starting backend service..."
docker-compose start backend

# Wait for backend to be healthy
log "Waiting for backend to be ready..."
sleep 5
for i in {1..30}; do
    if docker exec agent-todo-backend wget -q --spider http://localhost:18765/health 2>/dev/null; then
        log "✅ Backend is healthy"
        break
    fi
    if [ $i -eq 30 ]; then
        log "⚠️  Warning: Backend health check timeout"
    fi
    sleep 2
done

# Verify the restore
log "Verifying database..."
TABLE_COUNT=$(docker exec "$DB_CONTAINER" psql -U "$DB_USER" -d "$DB_NAME" -t -c \
    "SELECT count(*) FROM information_schema.tables WHERE table_schema = 'public';" | tr -d ' ')
log "Tables in database: $TABLE_COUNT"

TASK_COUNT=$(docker exec "$DB_CONTAINER" psql -U "$DB_USER" -d "$DB_NAME" -t -c \
    "SELECT count(*) FROM tasks;" | tr -d ' ')
log "Tasks in database: $TASK_COUNT"

# Clean up temporary backup
if [ -f "$TEMP_BACKUP" ]; then
    log "Cleaning up safety backup..."
    rm "$TEMP_BACKUP"
fi

log "✅ Restore completed successfully!"
log ""
log "Next steps:"
log "1. Verify application functionality at https://todo.formatho.com"
log "2. Check logs: docker logs agent-todo-backend"
log "3. Test API endpoints"
log "4. Update incident log if this was a disaster recovery"

exit 0
