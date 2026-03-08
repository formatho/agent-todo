-- Add organisation_id to projects table
ALTER TABLE projects ADD COLUMN organisation_id UUID REFERENCES organisations(id) ON DELETE SET NULL;
CREATE INDEX idx_projects_org_id ON projects(organisation_id);

-- Add organisation_id to tasks table
ALTER TABLE tasks ADD COLUMN organisation_id UUID REFERENCES organisations(id) ON DELETE SET NULL;
CREATE INDEX idx_tasks_org_id ON tasks(organisation_id);

-- Add organisation_id to agents table
ALTER TABLE agents ADD COLUMN organisation_id UUID REFERENCES organisations(id) ON DELETE SET NULL;
CREATE INDEX idx_agents_org_id ON agents(organisation_id);

-- Add current_org_id to users table
ALTER TABLE users ADD COLUMN current_org_id UUID REFERENCES organisations(id) ON DELETE SET NULL;
CREATE INDEX idx_users_current_org ON users(current_org_id);
