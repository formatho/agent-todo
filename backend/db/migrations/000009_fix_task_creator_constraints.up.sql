-- Migration: Fix task creator foreign key constraints
-- This migration fixes the foreign key constraint issue when agents create tasks

-- Drop the old foreign key constraint on created_by_user_id if it exists
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'fk_tasks_created_by'
    ) THEN
        ALTER TABLE tasks DROP CONSTRAINT fk_tasks_created_by;
    END IF;
END $$;

-- Drop the old foreign key constraint on created_by_agent_id if it exists
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'fk_tasks_created_by_agent'
    ) THEN
        ALTER TABLE tasks DROP CONSTRAINT fk_tasks_created_by_agent;
    END IF;
END $$;

-- Add proper foreign key constraints that allow NULL values
-- The constraints will only be enforced when the value is NOT NULL

ALTER TABLE tasks
    ADD CONSTRAINT fk_tasks_created_by_user
    FOREIGN KEY (created_by_user_id)
    REFERENCES users(id)
    ON DELETE SET NULL
    ON UPDATE CASCADE;

ALTER TABLE tasks
    ADD CONSTRAINT fk_tasks_created_by_agent
    FOREIGN KEY (created_by_agent_id)
    REFERENCES agents(id)
    ON DELETE SET NULL
    ON UPDATE CASCADE;

-- Ensure the check constraint exists (mutual exclusivity)
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'chk_task_creator'
    ) THEN
        ALTER TABLE tasks ADD CONSTRAINT chk_task_creator
        CHECK (
            (created_by_user_id IS NOT NULL AND created_by_agent_id IS NULL) OR
            (created_by_user_id IS NULL AND created_by_agent_id IS NOT NULL)
        );
    END IF;
END $$;

COMMENT ON COLUMN tasks.created_by_user_id IS 'User who created the task (NULL if created by agent)';
COMMENT ON COLUMN tasks.created_by_agent_id IS 'Agent who created the task (NULL if created by user)';
