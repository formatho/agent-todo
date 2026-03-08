DROP INDEX IF EXISTS idx_tasks_due_date;
DROP INDEX IF EXISTS idx_tasks_priority;
DROP INDEX IF EXISTS idx_tasks_created_by_user_id;
DROP INDEX IF EXISTS idx_tasks_assigned_agent_id;
DROP INDEX IF EXISTS idx_tasks_status;
DROP TYPE IF EXISTS task_priority;
DROP TYPE IF EXISTS task_status;
DROP TABLE IF EXISTS tasks;
