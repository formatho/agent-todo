-- Add deleted_at column for soft deletes
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP WITH TIME ZONE;

-- Add index for soft delete queries
CREATE INDEX IF NOT EXISTS idx_tasks_deleted_at ON tasks(deleted_at);
