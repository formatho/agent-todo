-- First, add a default project for existing tasks
INSERT INTO projects (name, description, status, created_by_user_id)
SELECT 'Default Project', 'Default project for existing tasks', 'active', created_by_user_id
FROM tasks
LIMIT 1;

-- Add project_id column as nullable first
ALTER TABLE tasks ADD COLUMN project_id UUID REFERENCES projects(id) ON DELETE CASCADE;

-- Update existing tasks to use the default project
UPDATE tasks
SET project_id = (SELECT id FROM projects WHERE name = 'Default Project' LIMIT 1)
WHERE project_id IS NULL;

-- Now make it NOT NULL
ALTER TABLE tasks ALTER COLUMN project_id SET NOT NULL;

CREATE INDEX idx_tasks_project_id ON tasks(project_id);
