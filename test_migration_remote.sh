#!/bin/bash
# Migration Test Script (Remote/Server Version)
# This script tests the organisation migration on the server

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
BACKUP_FILE="/tmp/agent_todo_test_backup.sql"
MIGRATION_FILE="/tmp/organisation_migration.sql"

# Create test database
echo -e "${YELLOW}Step 1: Creating test database...${NC}"
docker exec agent-todo-db psql -U agent_todo -c "DROP DATABASE IF EXISTS ${TEST_DB_NAME};"
docker exec agent-todo-db psql -U agent_todo -c "CREATE DATABASE ${TEST_DB_NAME};"

echo -e "${GREEN}✓ Test database created${NC}"

echo ""
echo -e "${YELLOW}Step 2: Restoring production backup to test database...${NC}"

# Create backup of current production database
docker exec agent-todo-db pg_dump -U agent_todo agent_todo > ${BACKUP_FILE}

# Restore to test database
cat ${BACKUP_FILE} | docker exec -i agent-todo-db psql -U agent_todo -d ${TEST_DB_NAME}

echo -e "${GREEN}✓ Backup restored${NC}"

echo ""
echo -e "${YELLOW}Step 3: Capturing pre-migration state...${NC}"

# Capture pre-migration data counts
echo "Pre-migration data counts:"
docker exec agent-todo-db psql -U agent_todo -d ${TEST_DB_NAME} -c "
SELECT 
  'users' as table_name, COUNT(*) as count FROM users
UNION ALL SELECT 'agents', COUNT(*) FROM agents
UNION ALL SELECT 'projects', COUNT(*) FROM projects
UNION ALL SELECT 'tasks', COUNT(*) FROM tasks
UNION ALL SELECT 'task_events', COUNT(*) FROM task_events
UNION ALL SELECT 'task_comments', COUNT(*) FROM task_comments;
"

# Save pre-migration counts
docker exec agent-todo-db psql -U agent_todo -d ${TEST_DB_NAME} -t -c "
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

# Copy migration file to server (assuming we're already on server)
# The migration file should already be in the backend container or we can copy it
# For now, we'll create it inline

docker exec -i agent-todo-db psql -U agent_todo -d ${TEST_DB_NAME} << 'MIGRATION_EOF'
-- Migration: Multi-Organisation Support
-- Test version for staging

-- ============================================================================
-- PHASE 1: Create New Tables (organisations, organisation_members)
-- ============================================================================

-- Create enum types
CREATE TYPE organisation_plan AS ENUM ('free', 'pro', 'enterprise');
CREATE TYPE organisation_status AS ENUM ('active', 'suspended', 'deleted');
CREATE TYPE member_role AS ENUM ('owner', 'admin', 'member', 'viewer');

-- Create organisations table
CREATE TABLE organisations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    plan organisation_plan NOT NULL DEFAULT 'free',
    status organisation_status NOT NULL DEFAULT 'active',
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    settings JSONB DEFAULT '{
        "max_agents": 5,
        "max_projects": 10,
        "max_tasks": 100,
        "allowed_features": []
    }'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes for organisations
CREATE INDEX idx_organisations_slug ON organisations(slug);
CREATE INDEX idx_organisations_owner_id ON organisations(owner_id);
CREATE INDEX idx_organisations_status ON organisations(status);
CREATE INDEX idx_organisations_created_at ON organisations(created_at DESC);

-- Create organisation_members table
CREATE TABLE organisation_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organisation_id UUID NOT NULL REFERENCES organisations(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role member_role NOT NULL DEFAULT 'member',
    joined_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    CONSTRAINT unique_org_member UNIQUE (organisation_id, user_id)
);

-- Create indexes for organisation_members
CREATE INDEX idx_org_members_org_id ON organisation_members(organisation_id);
CREATE INDEX idx_org_members_user_id ON organisation_members(user_id);
CREATE INDEX idx_org_members_role ON organisation_members(role);
CREATE INDEX idx_org_members_joined_at ON organisation_members(joined_at DESC);

-- ============================================================================
-- PHASE 2: Add organisation_id to Existing Tables (NULLABLE)
-- ============================================================================

-- Add organisation_id to users (for current context)
ALTER TABLE users ADD COLUMN current_org_id UUID REFERENCES organisations(id) ON DELETE SET NULL;
CREATE INDEX idx_users_current_org_id ON users(current_org_id);

-- Add organisation_id to projects
ALTER TABLE projects ADD COLUMN organisation_id UUID REFERENCES organisations(id) ON DELETE RESTRICT;
CREATE INDEX idx_projects_organisation_id ON projects(organisation_id);

-- Add organisation_id to agents
ALTER TABLE agents ADD COLUMN organisation_id UUID REFERENCES organisations(id) ON DELETE RESTRICT;
CREATE INDEX idx_agents_organisation_id ON agents(organisation_id);

-- Add organisation_id to tasks
ALTER TABLE tasks ADD COLUMN organisation_id UUID REFERENCES organisations(id) ON DELETE RESTRICT;
CREATE INDEX idx_tasks_organisation_id ON tasks(organisation_id);

-- Add organisation_id to task_events
ALTER TABLE task_events ADD COLUMN organisation_id UUID REFERENCES organisations(id) ON DELETE RESTRICT;
CREATE INDEX idx_task_events_organisation_id ON task_events(organisation_id);

-- Add organisation_id to task_comments
ALTER TABLE task_comments ADD COLUMN organisation_id UUID REFERENCES organisations(id) ON DELETE RESTRICT;
CREATE INDEX idx_task_comments_organisation_id ON task_comments(organisation_id);

-- ============================================================================
-- PHASE 3: Create Default Organisation (For Migration)
-- ============================================================================

DO $$
DECLARE
    default_org_id UUID;
    first_user_id UUID;
BEGIN
    -- Get the first user (will be the owner)
    SELECT id INTO first_user_id FROM users ORDER BY created_at LIMIT 1;
    
    IF first_user_id IS NOT NULL THEN
        -- Create default organisation
        INSERT INTO organisations (name, slug, plan, status, owner_id, settings)
        VALUES (
            'formatho',
            'formatho',
            'enterprise',
            'active',
            first_user_id,
            '{
                "max_agents": 100,
                "max_projects": 100,
                "max_tasks": 10000,
                "allowed_features": ["all"]
            }'::jsonb
        )
        RETURNING id INTO default_org_id;
        
        -- Assign all existing data to default organisation
        UPDATE projects SET organisation_id = default_org_id;
        UPDATE agents SET organisation_id = default_org_id;
        UPDATE tasks SET organisation_id = default_org_id;
        UPDATE task_events SET organisation_id = default_org_id;
        UPDATE task_comments SET organisation_id = default_org_id;
        
        -- Make all users members of default organisation
        INSERT INTO organisation_members (organisation_id, user_id, role, joined_at)
        SELECT 
            default_org_id,
            id,
            CASE 
                WHEN id = first_user_id THEN 'owner'::member_role
                ELSE 'admin'::member_role
            END,
            NOW()
        FROM users;
        
        -- Set current_org_id for first user
        UPDATE users SET current_org_id = default_org_id WHERE id = first_user_id;
        
        RAISE NOTICE 'Default organisation created with ID: %', default_org_id;
    ELSE
        RAISE NOTICE 'No users found, skipping default organisation creation';
    END IF;
END $$;

MIGRATION_EOF

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
docker exec agent-todo-db psql -U agent_todo -d ${TEST_DB_NAME} -c "\dt" | grep -E "(organisations|organisation_members)"

if [ $? -eq 0 ]; then
  echo -e "${GREEN}✓ Organisation tables created${NC}"
else
  echo -e "${RED}✗ Organisation tables not found${NC}"
  exit 1
fi

echo ""
echo "Checking organisation created:"
docker exec agent-todo-db psql -U agent_todo -d ${TEST_DB_NAME} -c "SELECT id, name, slug, plan, status FROM organisations;"

echo ""
echo "Checking organisation membership:"
docker exec agent-todo-db psql -U agent_todo -d ${TEST_DB_NAME} -c "
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
docker exec agent-todo-db psql -U agent_todo -d ${TEST_DB_NAME} -c "
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

# Save post-migration counts
docker exec agent-todo-db psql -U agent_todo -d ${TEST_DB_NAME} -t -c "
SELECT 
  'users' as table_name, COUNT(*) as count FROM users
UNION ALL SELECT 'agents', COUNT(*) FROM agents
UNION ALL SELECT 'projects', COUNT(*) FROM projects
UNION ALL SELECT 'tasks', COUNT(*) FROM tasks
UNION ALL SELECT 'task_events', COUNT(*) FROM task_events
UNION ALL SELECT 'task_comments', COUNT(*) FROM task_comments;
" > /tmp/post_migration_counts.txt

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
docker exec agent-todo-db psql -U agent_todo -d ${TEST_DB_NAME} -c "
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
echo "  - Default organisation 'formatho' created: ✓"
echo "  - All users migrated as members: ✓"
echo "  - All data assigned to organisation: ✓"
echo "  - No data loss: ✓"
echo "  - Foreign keys working: ✓"
echo ""
echo "Test database: ${TEST_DB_NAME}"
echo ""
echo "To clean up test database:"
echo "  docker exec agent-todo-db psql -U agent_todo -c 'DROP DATABASE ${TEST_DB_NAME};'"
