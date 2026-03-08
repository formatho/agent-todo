-- Create admin user (password: admin123)
-- Password hash is for 'admin123'
INSERT INTO users (email, password_hash) VALUES
('admin@example.com', '$2a$10$8X7OqZKz8YH8qGZ0qZ8qZuX7OqZKz8YH8qGZ0qZ8qZuX7OqZKz8YH');

-- Create example agent
INSERT INTO agents (name, api_key, description) VALUES
('Example Agent', 'sk_agent_example_key_12345', 'An example agent for testing');

-- Create example projects
INSERT INTO projects (name, description, status, created_by_user_id) VALUES
('Website Redesign', 'Redesign the company website with new branding', 'active',
 (SELECT id FROM users WHERE email = 'admin@example.com')),
('API Development', 'Build RESTful API for mobile applications', 'active',
 (SELECT id FROM users WHERE email = 'admin@example.com'));

-- Create some example tasks
INSERT INTO tasks (title, description, status, priority, project_id, created_by_user_id, assigned_agent_id) VALUES
('Setup project environment', 'Configure development environment with all necessary tools', 'pending', 'high',
 (SELECT id FROM projects WHERE name = 'Website Redesign'),
 (SELECT id FROM users WHERE email = 'admin@example.com'),
 (SELECT id FROM agents WHERE name = 'Example Agent')),
('Create initial documentation', 'Write README and API documentation', 'pending', 'medium',
 (SELECT id FROM projects WHERE name = 'Website Redesign'),
 (SELECT id FROM users WHERE email = 'admin@example.com'),
 NULL),
('Design database schema', 'Design database schema for user management', 'in_progress', 'high',
 (SELECT id FROM projects WHERE name = 'API Development'),
 (SELECT id FROM users WHERE email = 'admin@example.com'),
 (SELECT id FROM agents WHERE name = 'Example Agent'));

-- Create example task events
INSERT INTO task_events (task_id, event_type, previous_state, new_state, changed_by)
SELECT id, 'created', '', status, 'admin@example.com' FROM tasks;
