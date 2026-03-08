-- Create default organisation for existing data
INSERT INTO organisations (id, name, slug, description, status, created_by_user_id, created_at, updated_at)
SELECT 
    gen_random_uuid(),
    'Default Organisation',
    'default',
    'Default organisation for existing data',
    'active',
    (SELECT id FROM users LIMIT 1),
    NOW(),
    NOW()
WHERE NOT EXISTS (SELECT 1 FROM organisations WHERE slug = 'default');

-- Get the default organisation ID for use in subsequent statements
DO $$
DECLARE
    default_org_id UUID;
BEGIN
    SELECT id INTO default_org_id FROM organisations WHERE slug = 'default' LIMIT 1;
    
    IF default_org_id IS NOT NULL THEN
        -- Add all existing users as members of the default organisation
        INSERT INTO organisation_members (id, organisation_id, user_id, role, joined_at, created_at, updated_at)
        SELECT 
            gen_random_uuid(),
            default_org_id,
            u.id,
            CASE 
                WHEN u.id = (SELECT created_by_user_id FROM organisations WHERE slug = 'default') THEN 'owner'
                ELSE 'member'
            END,
            NOW(),
            NOW(),
            NOW()
        FROM users u
        WHERE NOT EXISTS (
            SELECT 1 FROM organisation_members om 
            WHERE om.organisation_id = default_org_id AND om.user_id = u.id
        );
        
        -- Assign all existing projects to the default organisation
        UPDATE projects SET organisation_id = default_org_id WHERE organisation_id IS NULL;
        
        -- Assign all existing tasks to the default organisation
        UPDATE tasks SET organisation_id = default_org_id WHERE organisation_id IS NULL;
        
        -- Assign all existing agents to the default organisation
        UPDATE agents SET organisation_id = default_org_id WHERE organisation_id IS NULL;
        
        -- Set current_org_id for all users to the default organisation
        UPDATE users SET current_org_id = default_org_id WHERE current_org_id IS NULL;
    END IF;
END $$;
