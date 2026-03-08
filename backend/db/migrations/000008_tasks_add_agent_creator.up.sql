-- Add created_by_agent_id column to tasks table
ALTER TABLE tasks ADD COLUMN created_by_agent_id UUID REFERENCES agents(id) ON DELETE CASCADE;

-- Make created_by_user_id nullable to allow agent-created tasks
ALTER TABLE tasks ALTER COLUMN created_by_user_id DROP NOT NULL;

-- Add constraint to ensure at least one creator is set
ALTER TABLE tasks ADD CONSTRAINT chk_task_creator CHECK (
    (created_by_user_id IS NOT NULL AND created_by_agent_id IS NULL) OR
    (created_by_user_id IS NULL AND created_by_agent_id IS NOT NULL)
);

-- Add index for agent-created tasks
CREATE INDEX idx_tasks_created_by_agent_id ON tasks(created_by_agent_id);

-- Add comment to document the constraint
COMMENT ON CONSTRAINT chk_task_creator ON tasks IS 'Ensures exactly one creator (user or agent) is set';
