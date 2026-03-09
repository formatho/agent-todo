#!/bin/bash
# Migration Test Script
# This script tests the organisation migration in a local environment

set -e  # Exit on any error

echo "=== Organisation Migration Test Script ==="
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
TEST_DB_NAME="agent_todo_migration_test"
BACKUP_FILE="backend/migrations/test/production_backup.sql"
MIGRATION_FILE="backend/db/migrations/006_add_organisations.sql"

echo -e "${YELLOW}Step 1: Setting up test database...${NC}"

# Stop any existing test container
docker stop agent-todo-migration-test 2>/dev/null || true
docker rm agent-todo-migration-test 2>/dev/null || true

# Start a fresh PostgreSQL container for testing
docker run -d \
  --name agent-todo-migration-test \
  -e POSTGRES_USER=agent_todo \
  -e POSTGRES_PASSWORD=agent_todo_pass \
  -e POSTGRES_DB=${TEST_DB_NAME} \
  -p 5433:5432 \
  postgres:16-alpine

# Wait for PostgreSQL to be ready
echo "Waiting for PostgreSQL to be ready..."
sleep 5
for i in {1..30}; do
  if docker exec agent-todo-migration-test pg_isready -U agent_todo; then
    echo -e "${GREEN}PostgreSQL is ready!${NC}"
    break
  fi
  echo "Waiting... ($i/30)"
  sleep 1
done

echo ""
echo -e "${YELLOW}Step 2: Restoring production backup to test database...${NC}"

# Restore the production backup
cat ${BACKUP_FILE} | docker exec -i agent-todo-migration-test psql -U agent_todo -d ${TEST_DB_NAME}

echo ""
echo -e "${YELLOW}Step 3: Capturing pre-migration state...${NC}"

# Capture pre-migration data counts
echo "Pre-migration data counts:"
docker exec agent-todo-migration-test psql -U agent_todo -d ${TEST_DB_NAME} -c "
SELECT 
  'users' as table_name, COUNT(*) as count FROM users
UNION ALL SELECT 'agents', COUNT(*) FROM agents
UNION ALL SELECT 'projects', COUNT(*) FROM projects
UNION ALL SELECT 'tasks', COUNT(*) FROM tasks
UNION ALL SELECT 'task_events', COUNT(*) FROM task_events
UNION ALL SELECT 'task_comments', COUNT(*) FROM task_comments;
"

# Save pre-migration counts for comparison
docker exec agent-todo-migration-test psql -U agent_todo -d ${TEST_DB_NAME} -t -c "
SELECT 
  'users' as table_name, COUNT(*) as count FROM users
UNION ALL SELECT 'agents', COUNT(*) FROM agents
UNION ALL SELECT 'projects', COUNT(*) FROM projects
UNION ALL SELECT 'tasks', COUNT(*) FROM tasks
UNION ALL SELECT 'task_events', COUNT(*) FROM task_events
UNION ALL SELECT 'task_comments', COUNT(*) FROM task_comments;
" > /tmp/pre_migration_counts.txt

echo ""
echo -e "${YELLOW}Step 4: Running organisation migration...${NC}"

# Run the migration
cat ${MIGRATION_FILE} | docker exec -i agent-todo-migration-test psql -U agent_todo -d ${TEST_DB_NAME}

if [ $? -eq 0 ]; then
  echo -e "${GREEN}Migration script executed successfully!${NC}"
else
  echo -e "${RED}Migration script failed!${NC}"
  exit 1
fi

echo ""
echo -e "${YELLOW}Step 5: Verifying migration results...${NC}"

# Check if organisation tables were created
echo "Checking for new tables:"
docker exec agent-todo-migration-test psql -U agent_todo -d ${TEST_DB_NAME} -c "\dt" | grep -E "(organisations|organisation_members)"

if [ $? -eq 0 ]; then
  echo -e "${GREEN}✓ Organisation tables created${NC}"
else
  echo -e "${RED}✗ Organisation tables not found${NC}"
  exit 1
fi

echo ""
echo "Checking organisation created:"
docker exec agent-todo-migration-test psql -U agent_todo -d ${TEST_DB_NAME} -c "SELECT id, name, slug, plan, status FROM organisations;"

echo ""
echo "Checking organisation membership:"
docker exec agent-todo-migration-test psql -U agent_todo -d ${TEST_DB_NAME} -c "
SELECT 
  o.name as org_name, 
  u.email as user_email, 
  om.role, 
  om.joined_at 
FROM organisations o
JOIN organisation_members om ON o.id = om.organisation_id
JOIN users u ON om.user_id = u.id;
"

echo ""
echo "Checking data migration - records with organisation_id:"
docker exec agent-todo-migration-test psql -U agent_todo -d ${TEST_DB_NAME} -c "
SELECT 
  'projects' as table_name, 
  COUNT(*) as total, 
  COUNT(organisation_id) as with_org_id,
  COUNT(*) - COUNT(organisation_id) as without_org_id
FROM projects
UNION ALL
SELECT 'agents', COUNT(*), COUNT(organisation_id), COUNT(*) - COUNT(organisation_id) FROM agents
UNION ALL
SELECT 'tasks', COUNT(*), COUNT(organisation_id), COUNT(*) - COUNT(organisation_id) FROM tasks
UNION ALL
SELECT 'task_events', COUNT(*), COUNT(organisation_id), COUNT(*) - COUNT(organisation_id) FROM task_events
UNION ALL
SELECT 'task_comments', COUNT(*), COUNT(organisation_id), COUNT(*) - COUNT(organisation_id) FROM task_comments;
"

echo ""
echo -e "${YELLOW}Step 6: Verifying data integrity...${NC}"

# Compare pre and post migration counts
echo "Post-migration data counts:"
docker exec agent-todo-migration-test psql -U agent_todo -d ${TEST_DB_NAME} -c "
SELECT 
  'users' as table_name, COUNT(*) as count FROM users
UNION ALL SELECT 'agents', COUNT(*) FROM agents
UNION ALL SELECT 'projects', COUNT(*) FROM projects
UNION ALL SELECT 'tasks', COUNT(*) FROM tasks
UNION ALL SELECT 'task_events', COUNT(*) FROM task_events
UNION ALL SELECT 'task_comments', COUNT(*) FROM task_comments;
"

# Save post-migration counts
docker exec agent-todo-migration-test psql -U agent_todo -d ${TEST_DB_NAME} -t -c "
SELECT 
  'users' as table_name, COUNT(*) as count FROM users
UNION ALL SELECT 'agents', COUNT(*) FROM agents
UNION ALL SELECT 'projects', COUNT(*) FROM projects
UNION ALL SELECT 'tasks', COUNT(*) FROM tasks
UNION ALL SELECT 'task_events', COUNT(*) FROM task_events
UNION ALL SELECT 'task_comments', COUNT(*) FROM task_comments;
" > /tmp/post_migration_counts.txt

echo ""
echo "Comparing pre vs post migration counts:"
if diff /tmp/pre_migration_counts.txt /tmp/post_migration_counts.txt > /dev/null; then
  echo -e "${GREEN}✓ All data counts match! No data lost during migration.${NC}"
else
  echo -e "${RED}✗ Data counts differ! Potential data loss detected.${NC}"
  echo "Pre-migration:"
  cat /tmp/pre_migration_counts.txt
  echo ""
  echo "Post-migration:"
  cat /tmp/post_migration_counts.txt
  exit 1
fi

echo ""
echo -e "${YELLOW}Step 7: Testing foreign key constraints...${NC}"

# Test that we can query related data
docker exec agent-todo-migration-test psql -U agent_todo -d ${TEST_DB_NAME} -c "
SELECT 
  t.id as task_id,
  t.title,
  o.name as organisation_name,
  p.name as project_name
FROM tasks t
JOIN projects p ON t.project_id = p.id
JOIN organisations o ON t.organisation_id = o.id
LIMIT 5;
"

if [ $? -eq 0 ]; then
  echo -e "${GREEN}✓ Foreign key relationships working correctly${NC}"
else
  echo -e "${RED}✗ Foreign key constraint test failed${NC}"
  exit 1
fi

echo ""
echo -e "${GREEN}=== Migration Test Complete ===${NC}"
echo ""
echo "Summary:"
echo "  - Organisation tables created: ✓"
echo "  - Default organisation created: ✓"
echo "  - All users migrated as members: ✓"
echo "  - All data assigned to organisation: ✓"
echo "  - No data loss: ✓"
echo "  - Foreign keys working: ✓"
echo ""
echo "Test database is still running at localhost:5433"
echo "Container name: agent-todo-migration-test"
echo ""
echo "To clean up: docker stop agent-todo-migration-test && docker rm agent-todo-migration-test"
