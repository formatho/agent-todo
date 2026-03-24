-- Add assigned_agent_id column to subtasks table
ALTER TABLE subtasks ADD COLUMN assigned_agent_id uuid;

-- Add foreign key constraint
ALTER TABLE subtasks 
ADD CONSTRAINT fk_subtask_assigned_agent 
FOREIGN KEY (assigned_agent_id) REFERENCES agents(id);

-- Create index for faster lookups
CREATE INDEX idx_subtasks_assigned_agent ON subtasks(assigned_agent_id);
