-- Add repository, deployment, and LLM context fields to projects table

ALTER TABLE projects ADD COLUMN IF NOT EXISTS repository_url TEXT;
ALTER TABLE projects ADD COLUMN IF NOT EXISTS deployed_url TEXT;
ALTER TABLE projects ADD COLUMN IF NOT EXISTS documentation_url TEXT;
ALTER TABLE projects ADD COLUMN IF NOT EXISTS llm_context TEXT;

-- Add comment for documentation
COMMENT ON COLUMN projects.repository_url IS 'GitHub/GitLab repository URL';
COMMENT ON COLUMN projects.deployed_url IS 'Production/staging deployment URL';
COMMENT ON COLUMN projects.documentation_url IS 'Documentation URL';
COMMENT ON COLUMN projects.llm_context IS 'LLM context: instructions, guidelines, goals for AI agents';
