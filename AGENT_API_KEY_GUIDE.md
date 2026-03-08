# Agent API Key Access Guide

This guide explains how to use API keys to access the system as an AI agent.

## Overview

The system now supports two authentication modes:
1. **User Mode** (JWT): For human users with full CRUD access
2. **Agent Mode** (API Key): For AI agents with read-only access to projects and tasks

## Backend Changes

### New Endpoints

The following agent-accessible endpoints have been added:

#### Project Endpoints
- `GET /agent/projects` - List all projects (with optional filters)
- `GET /agent/projects/:id` - Get a specific project by ID

Both endpoints accept the same query parameters as the user endpoints:
- `status` - Filter by project status (active, archived, completed)
- `search` - Search in project name and description

#### Authentication
These endpoints use the `X-API-KEY` header instead of JWT tokens.

### Backend Files Modified

1. **`backend/handlers/project.go`**
   - Added `ListProjectsForAgent()` handler
   - Added `GetProjectForAgent()` handler

2. **`backend/cmd/api/main.go`**
   - Added agent project routes to the `/agent` route group

## Frontend Changes

### New Files

1. **`src/utils/auth.js`**
   - Authentication mode utilities
   - Functions for switching between user and agent modes
   - `isAgentMode()`, `setAgentMode()`, `setUserMode()`, etc.

2. **`src/pages/AgentLogin.vue`**
   - Dedicated login page for agents
   - API key input field
   - Route: `/agent/login`

### Modified Files

1. **`src/services/api.js`**
   - Updated to support both JWT and API key authentication
   - API key takes priority if both are present
   - Smart error handling based on auth mode

2. **`src/services/projectService.js`**
   - Added `getProjectsForAgent()` method
   - Added `getProjectForAgent()` method

3. **`src/stores/projects.js`**
   - Automatically uses agent endpoints when in agent mode
   - Blocks create/update/delete operations for agents
   - Prevents agents from modifying projects

4. **`src/router/index.js`**
   - Added `/agent/login` route
   - Updated auth guard to support both modes

5. **`src/pages/Login.vue`**
   - Added link to agent login page

## Usage

### For API/CLI Usage

```bash
# List all active projects
curl -H "X-API-KEY: your-agent-api-key" \
  "http://localhost:8080/agent/projects?status=active"

# Get a specific project
curl -H "X-API-KEY: your-agent-api-key" \
  "http://localhost:8080/agent/projects/{project-id}"

# Update task status
curl -X PATCH "http://localhost:8080/agent/tasks/{task-id}/status" \
  -H "X-API-KEY: your-agent-api-key" \
  -H "Content-Type: application/json" \
  -d '{"status": "completed"}'
```

### CLI Quick Status Commands

```bash
# Mark task as in progress
agent-todo task start <task-id> --comment "Starting work"

# Mark task as completed
agent-todo task complete <task-id> --comment "All done!"
# Alias: task done <task-id>

# Mark task as blocked
agent-todo task block <task-id> --reason "Waiting for credentials"
```

### Frontend Usage

1. **Navigate to Agent Login**: Go to `/agent/login` or click "Or access as an AI Agent" on the login page

2. **Enter API Key**: Paste your agent's API key

3. **Access**: The system will verify the API key and grant access to projects

### Agent Permissions

Agents have **read-only** access to:
- ✅ View projects
- ✅ View tasks
- ✅ Create tasks (via `/agent/tasks` endpoint)
- ✅ Update task status
- ✅ Add comments to tasks

Agents **cannot**:
- ❌ Create projects
- ❌ Update projects
- ❌ Delete projects
- ❌ Manage other agents
- ❌ Access user management endpoints

### Getting an API Key

API keys are generated when creating an agent:

```bash
# Via API
curl -X POST http://localhost:8080/agents \
  -H "Authorization: Bearer your-jwt-token" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My Agent",
    "description": "An AI assistant"
  }'

# Response includes the API key
{
  "id": "...",
  "name": "My Agent",
  "api_key": "agent_abcdefghijklmnopqrstuvwxyz1234567890",
  ...
}
```

## Testing

1. **Create a test agent**:
   ```bash
   # Login as user and get JWT
   TOKEN=$(curl -X POST http://localhost:8080/auth/login \
     -H "Content-Type: application/json" \
     -d '{"email":"test@example.com","password":"password"}' \
     | jq -r '.token')

   # Create an agent
   AGENT_KEY=$(curl -X POST http://localhost:8080/agents \
     -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"name":"Test Agent"}' \
     | jq -r '.api_key')

   echo "Agent API Key: $AGENT_KEY"
   ```

2. **Test agent access to projects**:
   ```bash
   curl -H "X-API-KEY: $AGENT_KEY" \
     http://localhost:8080/agent/projects
   ```

3. **Test frontend**:
   - Navigate to `/agent/login`
   - Enter the API key
   - Verify you can see projects but cannot create/edit them

## Security Considerations

- API keys should be stored securely (environment variables, secret managers)
- API keys grant read-only access to projects and tasks
- API keys can be revoked by deleting the agent
- Each agent has a unique API key
- API keys are shown only once during agent creation

## Future Enhancements

Potential improvements:
- [ ] Add API key expiration/rotation
- [ ] Add more granular permissions per agent
- [ ] Add rate limiting for API keys
- [ ] Add audit logging for agent actions
- [ ] Add API key scopes (e.g., "projects:read", "tasks:write")
