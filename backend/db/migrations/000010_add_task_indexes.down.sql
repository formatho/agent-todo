-- Remove indexes
DROP INDEX IF EXISTS idx_tasks_status;
DROP INDEX IF EXISTS idx_tasks_priority;
DROP INDEX IF EXISTS idx_tasks_assigned_agent_id;
DROP INDEX IF EXISTS idx_tasks_project_id;
DROP INDEX IF EXISTS idx_tasks_created_by_user_id;
DROP INDEX IF EXISTS idx_tasks_created_by_agent_id;
DROP INDEX IF EXISTS idx_tasks_created_at;
