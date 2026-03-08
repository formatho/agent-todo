-- Create organisations table
CREATE TABLE organisations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    created_by_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create organisation_members table
CREATE TABLE organisation_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organisation_id UUID NOT NULL REFERENCES organisations(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL DEFAULT 'member',
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(organisation_id, user_id)
);

-- Create indexes for performance
CREATE INDEX idx_organisations_slug ON organisations(slug);
CREATE INDEX idx_organisations_status ON organisations(status);
CREATE INDEX idx_organisations_created_by ON organisations(created_by_user_id);
CREATE INDEX idx_org_members_org_id ON organisation_members(organisation_id);
CREATE INDEX idx_org_members_user_id ON organisation_members(user_id);
CREATE INDEX idx_org_members_role ON organisation_members(role);

-- Add check constraint for valid status
ALTER TABLE organisations ADD CONSTRAINT chk_org_status 
    CHECK (status IN ('active', 'suspended', 'archived'));

-- Add check constraint for valid member role
ALTER TABLE organisation_members ADD CONSTRAINT chk_member_role 
    CHECK (role IN ('owner', 'admin', 'member'));
