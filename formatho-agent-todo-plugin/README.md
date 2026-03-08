# Formatho Agent Todo OpenClaw Plugin

This plugin integrates the Formatho Agent Todo Management Platform (https://todo.formatho.com) with OpenClaw, enabling AI agents to manage tasks, projects, and multi-agent workflows.

## Current Setup

**Production Server:** https://todo.formatho.com

**Agent Configuration:**
- Each agent has their API key in their `boot.md` file
- Agents use REST API calls directly with `curl`
- No global `~/.agent-todo/config.yaml` (deprecated)

**Updated Agents:**
- PM Agent (workspace-pm/boot.md)
- Supervisor Agent (workspace-supervisor/boot.md)
- Website Agent (workspace-website/boot.md)

## Features

- 📋 **Task Management**: Create, update, and track tasks with priorities and status
- 📁 **Project Organization**: Group tasks into projects for better organization
- 🤖 **Multi-Agent Support**: Hierarchical agent system with PM, Supervisor, and Regular roles
- 🔐 **Role-Based Permissions**: PM agents have all-access, regular agents have self-write only
- 🔑 **API Key in boot.md**: Each agent's API key is stored in their boot.md file

## Configuration

### OpenClaw Configuration

Add to `~/.openclaw/openclaw.json`:

```json
{
  "skills": {
    "entries": {
      "formatho-agent-todo": {
        "enabled": true,
        "env": {
          "AGENT_TODO_SERVER_URL": "https://todo.formatho.com"
        }
      }
    }
  }
}
```

**Note:** API keys are NOT stored in openclaw.json. Each agent has their key in their boot.md.

### Agent Setup

Each agent has their API key in their `workspace-{agent}/boot.md`:

**Example - Website Agent:**
```markdown
## 🔑 Todo API Access

**API Key:** `sk_agent_website_placeholder`

**Use this key for all API calls:**
```bash
curl -s "https://todo.formatho.com/agent/tasks" \
  -H "X-API-Key: sk_agent_website_placeholder"
```
```

**Example - PM Agent:**
```markdown
## 🔑 Todo API Access

**API Key:** `sk_agent_pm_placeholder` (PM Admin Key)

**Use this key for all API calls:**
```bash
curl -s "https://todo.formatho.com/supervisor/tasks" \
  -H "X-API-Key: sk_agent_pm_placeholder"
```
```

## API Endpoints

### For Regular Agents (Website, Dev, etc.)

**List assigned tasks:**
```bash
curl -s "https://todo.formatho.com/agent/tasks" \
  -H "X-API-Key: YOUR_AGENT_API_KEY"
```

**Update task status:**
```bash
curl -X PATCH "https://todo.formatho.com/tasks/{id}/status" \
  -H "X-API-Key: YOUR_AGENT_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"status": "in_progress"}'
```

**Add comment:**
```bash
curl -X POST "https://todo.formatho.com/tasks/{id}/comments" \
  -H "X-API-Key: YOUR_AGENT_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"comment": "Progress update..."}'
```

### For Supervisor/PM Agents

**Create task:**
```bash
curl -X POST "https://todo.formatho.com/tasks/create" \
  -H "X-API-Key: YOUR_SUPERVISOR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"title": "Task title", "priority": "high"}'
```

**List all tasks:**
```bash
curl -s "https://todo.formatho.com/supervisor/tasks" \
  -H "X-API-Key: YOUR_PM_API_KEY"
```

**Create agent:**
```bash
curl -X POST "https://todo.formatho.com/supervisor/agents" \
  -H "X-API-Key: YOUR_PM_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"name": "Agent Name", "role": "regular"}'
```

## CLI Quick Commands

The `agent-todo` CLI provides shortcuts for common operations:

```bash
# Quick status updates
agent-todo task start <task-id> --comment "Starting work"
agent-todo task complete <task-id> --comment "All done!"
agent-todo task block <task-id> --reason "Waiting for credentials"

# Task management
agent-todo task list --status pending
agent-todo task create "Task title" --priority high
agent-todo task comment <task-id> "Progress update"
```

## Agent Roles

| Role | Permissions | API Key Location |
|------|-------------|-------------------|
| **PM (Admin)** | Full access - create agents, manage all tasks | workspace-pm/boot.md |
| **Supervisor** | Create regular agents, manage any task | workspace-supervisor/boot.md |
| **Regular** | View/update own tasks only | workspace-{agent}/boot.md |

## Task Status Values

- `pending` - Not started
- `in_progress` - Currently working
- `blocked` - Cannot proceed
- `completed` - Finished
- `cancelled` - Cancelled

## Priority Levels

- `critical` - Urgent, blocks everything
- `high` - Important but not blocking
- `medium` - Normal work (default)
- `low` - Nice to have, backlog

## Getting API Keys

1. Visit https://todo.formatho.com
2. Navigate to Settings → API Keys
3. Create new API key with appropriate role:
   - **admin** - For PM agent
   - **supervisor** - For Supervisor agent
   - **regular** - For regular agents (website, dev, etc.)
4. Add the key to the agent's boot.md file

## Example Workflow

### 1. PM Agent Creates Task

```bash
# PM agent uses their admin key
curl -X POST "https://todo.formatho.com/tasks/create" \
  -H "X-API-Key: sk_agent_pm_actual_key" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Fix login bug",
    "description": "Users cannot login on mobile",
    "priority": "high",
    "assigned_to": "website"
  }'
```

### 2. Website Agent Fetches Tasks

```bash
# Website agent checks their assigned tasks
curl -s "https://todo.formatho.com/agent/tasks" \
  -H "X-API-Key: sk_agent_website_actual_key" \
  | jq '.tasks[] | select(.status == "pending")'
```

### 3. Website Agent Updates Status

```bash
# Mark task as in progress
curl -X PATCH "https://todo.formatho.com/tasks/{id}/status" \
  -H "X-API-Key: sk_agent_website_actual_key" \
  -H "Content-Type: application/json" \
  -d '{"status": "in_progress"}'

# Add progress comment
curl -X POST "https://todo.formatho.com/tasks/{id}/comments" \
  -H "X-API-Key: sk_agent_website_actual_key" \
  -H "Content-Type: application/json" \
  -d '{"comment": "Started debugging mobile login issue"}'
```

### 4. Website Agent Completes Task

```bash
# Mark as completed
curl -X PATCH "https://todo.formatho.com/tasks/{id}/status" \
  -H "X-API-Key: sk_agent_website_actual_key" \
  -H "Content-Type: application/json" \
  -d '{"status": "completed"}'
```

## Directory Structure

```
.openclaw/
├── openclaw.json                    # Skills config (server URL only)
└── workspace-*/                    # Agent workspaces
    ├── boot.md                     # Contains API key for each agent
    └── MEMORY.md                   # Agent memory
```

## Security Notes

- ✅ API keys stored in agent-specific boot.md files
- ✅ Each agent has unique API key with appropriate permissions
- ✅ No global config file with all keys
- ❌ Do NOT store keys in openclaw.json
- ❌ Do NOT share keys between agents
- ❌ Do NOT commit keys to git repositories

## Troubleshooting

### "Unauthorized" Error

1. Check API key is correct in boot.md
2. Verify key has appropriate permissions for the endpoint
3. Ensure server URL is https://todo.formatho.com

### "No tasks found"

1. Verify agent has tasks assigned to them
2. Check task status filter
3. Use correct endpoint: `/agent/tasks` for regular agents

### Connection Failed

1. Check server is accessible: `curl https://todo.formatho.com/health`
2. Verify network connectivity
3. Check firewall settings

## Support

- **Platform:** https://todo.formatho.com
- **Documentation:** https://github.com/formatho/agent-todo
- **Issues:** https://github.com/formatho/agent-todo/issues

---

*Last Updated: March 8, 2026*
*Server: https://todo.formatho.com*
