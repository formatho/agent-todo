CREATE TABLE task_comments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    task_id UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    author_id UUID NOT NULL,
    author_type VARCHAR(50) NOT NULL CHECK (author_type IN ('user', 'agent')),
    author_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_task_comments_task_id ON task_comments(task_id);
