-- Drop the index
DROP INDEX IF EXISTS idx_tasks_created_by_agent_id;

-- Drop the check constraint
ALTER TABLE tasks DROP CONSTRAINT IF EXISTS chk_task_creator;

-- Drop the created_by_agent_id column
ALTER TABLE tasks DROP COLUMN IF EXISTS created_by_agent_id;

-- Make created_by_user_id NOT NULL again (this will fail if there are null values, which is expected)
-- Note: In production, you'd need to handle existing null values first
ALTER TABLE tasks ALTER COLUMN created_by_user_id SET NOT NULL;
