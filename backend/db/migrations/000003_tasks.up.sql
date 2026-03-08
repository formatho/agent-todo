CREATE TYPE task_status AS ENUM ('pending', 'in_progress', 'completed', 'failed');

CREATE TYPE task_priority AS ENUM ('low', 'medium', 'high', 'critical');

CREATE TABLE tasks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(500) NOT NULL,
    description TEXT,
    status task_status NOT NULL DEFAULT 'pending',
    priority task_priority NOT NULL DEFAULT 'medium',
    due_date TIMESTAMP,
    created_by_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    assigned_agent_id UUID REFERENCES agents(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_assigned_agent_id ON tasks(assigned_agent_id);
CREATE INDEX idx_tasks_created_by_user_id ON tasks(created_by_user_id);
CREATE INDEX idx_tasks_priority ON tasks(priority);
CREATE INDEX idx_tasks_due_date ON tasks(due_date);
