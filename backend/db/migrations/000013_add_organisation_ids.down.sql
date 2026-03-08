-- Remove current_org_id from users
DROP INDEX IF EXISTS idx_users_current_org;
ALTER TABLE users DROP COLUMN IF EXISTS current_org_id;

-- Remove organisation_id from agents
DROP INDEX IF EXISTS idx_agents_org_id;
ALTER TABLE agents DROP COLUMN IF EXISTS organisation_id;

-- Remove organisation_id from tasks
DROP INDEX IF EXISTS idx_tasks_org_id;
ALTER TABLE tasks DROP COLUMN IF EXISTS organisation_id;

-- Remove organisation_id from projects
DROP INDEX IF EXISTS idx_projects_org_id;
ALTER TABLE projects DROP COLUMN IF EXISTS organisation_id;
