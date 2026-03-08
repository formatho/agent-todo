-- Migration: Multi-Organisation Support
-- Task: e38a9d26-c23d-464a-9229-c6fdbcc16e4c
-- Date: 2026-03-08

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

-- Add comments
COMMENT ON TABLE organisations IS 'Organisations (tenants) in the multi-tenant system';
COMMENT ON COLUMN organisations.slug IS 'URL-friendly organisation identifier (unique)';
COMMENT ON COLUMN organisations.settings IS 'JSON object containing plan limits and feature flags';

-- ============================================================================
-- Create organisation_members table
-- ============================================================================

CREATE TABLE organisation_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organisation_id UUID NOT NULL REFERENCES organisations(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role member_role NOT NULL DEFAULT 'member',
    joined_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    -- Ensure a user can only be a member once per organisation
    CONSTRAINT unique_org_member UNIQUE (organisation_id, user_id)
);

-- Create indexes for organisation_members
CREATE INDEX idx_org_members_org_id ON organisation_members(organisation_id);
CREATE INDEX idx_org_members_user_id ON organisation_members(user_id);
CREATE INDEX idx_org_members_role ON organisation_members(role);
CREATE INDEX idx_org_members_joined_at ON organisation_members(joined_at DESC);

-- Add comments
COMMENT ON TABLE organisation_members IS 'Membership relationship between users and organisations';
COMMENT ON COLUMN organisation_members.role IS 'User role within the organisation (owner, admin, member, viewer)';

-- ============================================================================
-- PHASE 2: Add organisation_id to Existing Tables (NULLABLE)
-- ============================================================================

-- Add organisation_id to users (for current context)
ALTER TABLE users ADD COLUMN current_org_id UUID REFERENCES organisations(id) ON DELETE SET NULL;
CREATE INDEX idx_users_current_org_id ON users(current_org_id);

COMMENT ON COLUMN users.current_org_id IS 'Currently active organisation context for the user';

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
-- PHASE 3: Create Composite Indexes for Performance
-- ============================================================================

-- Projects: organisation + name unique within org
CREATE UNIQUE INDEX idx_org_project_unique ON projects(organisation_id, name) 
    WHERE organisation_id IS NOT NULL;

-- Agents: organisation + name unique within org
CREATE UNIQUE INDEX idx_org_agent_unique ON agents(organisation_id, name)
    WHERE organisation_id IS NOT NULL;

-- Tasks: organisation + status + priority for filtered queries
CREATE INDEX idx_org_task_status_priority ON tasks(organisation_id, status, priority)
    WHERE organisation_id IS NOT NULL;

-- Tasks: organisation + created_at for pagination
CREATE INDEX idx_org_task_created_at ON tasks(organisation_id, created_at DESC)
    WHERE organisation_id IS NOT NULL;

-- Task events: organisation + created_at for audit logs
CREATE INDEX idx_org_event_created_at ON task_events(organisation_id, created_at DESC)
    WHERE organisation_id IS NOT NULL;

-- ============================================================================
-- PHASE 4: Create Functions and Triggers
-- ============================================================================

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Triggers for updated_at
CREATE TRIGGER update_organisations_updated_at
    BEFORE UPDATE ON organisations
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_org_members_updated_at
    BEFORE UPDATE ON organisation_members
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- ============================================================================
-- PHASE 5: Create Default Organisation (For Migration)
-- ============================================================================

-- This will be executed after deployment
-- Creates a default organisation for existing data

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
            'Default Organisation',
            'default',
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

-- ============================================================================
-- PHASE 6: Make organisation_id NOT NULL (After Migration)
-- ============================================================================

-- Run these AFTER verifying all data has been migrated
-- These statements should be executed separately after confirming success

-- ALTER TABLE projects ALTER COLUMN organisation_id SET NOT NULL;
-- ALTER TABLE agents ALTER COLUMN organisation_id SET NOT NULL;
-- ALTER TABLE tasks ALTER COLUMN organisation_id SET NOT NULL;
-- ALTER TABLE task_events ALTER COLUMN organisation_id SET NOT NULL;
-- ALTER TABLE task_comments ALTER COLUMN organisation_id SET NOT NULL;

-- ============================================================================
-- ROLLBACK SCRIPT (If Needed)
-- ============================================================================

-- To rollback this migration:
/*
-- Drop triggers
DROP TRIGGER IF EXISTS update_organisations_updated_at ON organisations;
DROP TRIGGER IF EXISTS update_org_members_updated_at ON organisation_members;

-- Drop functions
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_organisations_slug;
DROP INDEX IF EXISTS idx_organisations_owner_id;
DROP INDEX IF EXISTS idx_organisations_status;
DROP INDEX IF EXISTS idx_organisations_created_at;
DROP INDEX IF EXISTS idx_org_members_org_id;
DROP INDEX IF EXISTS idx_org_members_user_id;
DROP INDEX IF EXISTS idx_org_members_role;
DROP INDEX IF EXISTS idx_org_members_joined_at;
DROP INDEX IF EXISTS idx_users_current_org_id;
DROP INDEX IF EXISTS idx_projects_organisation_id;
DROP INDEX IF EXISTS idx_agents_organisation_id;
DROP INDEX IF EXISTS idx_tasks_organisation_id;
DROP INDEX IF EXISTS idx_task_events_organisation_id;
DROP INDEX IF EXISTS idx_task_comments_organisation_id;
DROP INDEX IF EXISTS idx_org_project_unique;
DROP INDEX IF EXISTS idx_org_agent_unique;
DROP INDEX IF EXISTS idx_org_task_status_priority;
DROP INDEX IF EXISTS idx_org_task_created_at;
DROP INDEX IF EXISTS idx_org_event_created_at;

-- Drop columns
ALTER TABLE users DROP COLUMN IF EXISTS current_org_id;
ALTER TABLE projects DROP COLUMN IF EXISTS organisation_id;
ALTER TABLE agents DROP COLUMN IF EXISTS organisation_id;
ALTER TABLE tasks DROP COLUMN IF EXISTS organisation_id;
ALTER TABLE task_events DROP COLUMN IF EXISTS organisation_id;
ALTER TABLE task_comments DROP COLUMN IF EXISTS organisation_id;

-- Drop tables
DROP TABLE IF EXISTS organisation_members;
DROP TABLE IF EXISTS organisations;

-- Drop enums
DROP TYPE IF EXISTS member_role;
DROP TYPE IF EXISTS organisation_status;
DROP TYPE IF EXISTS organisation_plan;
*/

-- ============================================================================
-- VERIFICATION QUERIES
-- ============================================================================

-- After migration, run these to verify:
/*
-- Check organisations created
SELECT id, name, slug, plan, status FROM organisations;

-- Check all projects have organisation_id
SELECT COUNT(*) AS projects_without_org FROM projects WHERE organisation_id IS NULL;

-- Check all tasks have organisation_id
SELECT COUNT(*) AS tasks_without_org FROM tasks WHERE organisation_id IS NULL;

-- Check all agents have organisation_id
SELECT COUNT(*) AS agents_without_org FROM agents WHERE organisation_id IS NULL;

-- Check organisation membership
SELECT o.name, COUNT(om.id) AS member_count
FROM organisations o
LEFT JOIN organisation_members om ON o.id = om.organisation_id
GROUP BY o.id, o.name;
*/
