-- This migration is not easily reversible as it involves data changes
-- The down migration will just remove the default organisation and memberships

-- Remove users from default organisation
DELETE FROM organisation_members WHERE organisation_id IN (SELECT id FROM organisations WHERE slug = 'default');

-- Nullify organisation references for projects
UPDATE projects SET organisation_id = NULL WHERE organisation_id IN (SELECT id FROM organisations WHERE slug = 'default');

-- Nullify organisation references for tasks
UPDATE tasks SET organisation_id = NULL WHERE organisation_id IN (SELECT id FROM organisations WHERE slug = 'default');

-- Nullify organisation references for agents
UPDATE agents SET organisation_id = NULL WHERE organisation_id IN (SELECT id FROM organisations WHERE slug = 'default');

-- Nullify current_org_id for users
UPDATE users SET current_org_id = NULL WHERE current_org_id IN (SELECT id FROM organisations WHERE slug = 'default');

-- Delete default organisation
DELETE FROM organisations WHERE slug = 'default';
