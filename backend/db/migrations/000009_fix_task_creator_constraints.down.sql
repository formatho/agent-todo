-- Rollback migration for task creator foreign key constraints

-- Drop the new constraints
ALTER TABLE tasks DROP CONSTRAINT IF EXISTS fk_tasks_created_by_user;
ALTER TABLE tasks DROP CONSTRAINT IF EXISTS fk_tasks_created_by_agent;

-- Drop the check constraint
ALTER TABLE tasks DROP CONSTRAINT IF EXISTS chk_task_creator;

-- Restore the original constraint (for rollback purposes)
-- Note: This may fail if there are agent-created tasks
ALTER TABLE tasks
    ADD CONSTRAINT fk_tasks_created_by
    FOREIGN KEY (created_by_user_id)
    REFERENCES users(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;
