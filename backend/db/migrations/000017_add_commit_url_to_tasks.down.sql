-- Remove commit_url field from tasks table
ALTER TABLE tasks DROP COLUMN IF EXISTS commit_url;
