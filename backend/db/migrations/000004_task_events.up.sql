CREATE TYPE task_event_type AS ENUM ('created', 'updated', 'status_changed', 'assigned', 'unassigned');

CREATE TABLE task_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    task_id UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    event_type task_event_type NOT NULL,
    previous_state TEXT,
    new_state TEXT,
    changed_by VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_task_events_task_id ON task_events(task_id);
CREATE INDEX idx_task_events_event_type ON task_events(event_type);
