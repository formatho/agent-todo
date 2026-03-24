DELETE FROM task_events WHERE changed_by = 'admin@example.com';
DELETE FROM tasks WHERE created_by_user_id IN (SELECT id FROM users WHERE email = 'admin@example.com');
DELETE FROM agents WHERE name = 'Example Agent';
DELETE FROM users WHERE email = 'admin@example.com';
