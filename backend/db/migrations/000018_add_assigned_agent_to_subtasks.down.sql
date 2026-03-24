-- Drop foreign key constraint first
ALTER TABLE subtasks DROP CONSTRAINT IF EXISTS fk_subtask_assigned_agent;

-- Drop the index
DROP INDEX IF EXISTS idx_subtasks_assigned_agent;

-- Remove the column
ALTER TABLE subtasks DROP COLUMN IF EXISTS assigned_agent_id;
