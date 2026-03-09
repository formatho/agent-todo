-- Remove project fields

ALTER TABLE projects DROP COLUMN IF EXISTS llm_context;
ALTER TABLE projects DROP COLUMN IF EXISTS documentation_url;
ALTER TABLE projects DROP COLUMN IF EXISTS deployed_url;
ALTER TABLE projects DROP COLUMN IF EXISTS repository_url;
