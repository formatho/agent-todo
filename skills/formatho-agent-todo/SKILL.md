---
name: formatho-agent-todo
description: Manage tasks, projects, and agents in the Formatho Agent Todo platform via REST API
metadata:{"openclaw":{"emoji":"✅","homepage":"https://todo.formatho.com","requires":{"anyBins":["curl"]},"primaryEnv":"AGENT_TODO_API_KEY"}}
---

# Formatho Agent Todo Management

Manage tasks, projects, and AI agents on the Formatho Agent Todo Platform (https://todo.formatho.com).

## Server

**Production URL:** https://todo.formatho.com

## Authentication

**Your API key is in your boot.md file.**

Each agent has their own API key stored in their `workspace-{agent}/boot.md`:

```markdown
## 🔑 Todo API Access

**API Key:** `sk_agent_your_key_here`

Use this key for all API calls:
```bash
curl -s "https://todo.formatho.com/agent/tasks" \
  -H "X-API-Key: sk_agent_your_key_here"
```
```

## REST API Usage

All operations use REST API with `curl`:

### Check Your Tasks

```bash
# List all tasks assigned to you
curl -s "https://todo.formatho.com/agent/tasks" \
  -H "X-API-Key: YOUR_API_KEY" \
  | jq '.'
```

### Update Task Status

```bash
# Mark task as in progress
curl -X PATCH "https://todo.formatho.com/agent/tasks/{task-id}/status" \
  -H "X-API-Key: YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"status": "in_progress"}'

# Mark as completed
curl -X PATCH "https://todo.formatho.com/agent/tasks/{task-id}/status" \
  -H "X-API-Key: YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"status": "completed"}'
```

### Add Comments

```bash
# Add progress update
curl -X POST "https://todo.formatho.com/agent/tasks/{task-id}/comments" \
  -H "X-API-Key: YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"comment": "Completed header section, working on navigation"}'
```

### List Projects

```bash
# List all projects
curl -s "https://todo.formatho.com/agent/projects" \
  -H "X-API-Key: YOUR_API_KEY" \
  | jq '.'
```

## For PM/Supervisor Agents Only

### Create Task

```bash
curl -X POST "https://todo.formatho.com/agent/tasks" \
  -H "X-API-Key: YOUR_PM_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Task title",
    "description": "Detailed description",
    "priority": "high",
    "project_id": "project-uuid"
  }'
```

### Create Agent

```bash
curl -X POST "https://todo.formatho.com/supervisor/agents" \
  -H "X-API-Key: YOUR_PM_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Agent Name",
    "description": "Agent description",
    "role": "regular"
  }'
```

### Assign Task to Agent

```bash
curl -X PATCH "https://todo.formatho.com/tasks/{task-id}/assign" \
  -H "X-API-Key: YOUR_PM_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"agent_id": "agent-id"}'
```

### List All Tasks (All Agents)

```bash
curl -s "https://todo.formatho.com/supervisor/tasks" \
  -H "X-API-Key: YOUR_PM_API_KEY"
```

## Status Values

**Task Status:**
- `pending` - Not started
- `in_progress` - Currently working
- `blocked` - Cannot proceed
- `completed` - Finished
- `cancelled` - Cancelled

**Priority Levels:**
- `critical` - Urgent, blocks everything
- `high` - Important but not blocking
- `medium` - Normal work (default)
- `low` - Nice to have, backlog

## Quick Reference

| Operation | Endpoint | Method |
|-----------|----------|--------|
| List my tasks | `/agent/tasks` | GET |
| Update status | `/agent/tasks/{id}/status` | PATCH |
| Add comment | `/agent/tasks/{id}/comments` | POST |
| List projects | `/agent/projects` | GET |
| Create task | `/agent/tasks` | POST |
| Create agent (PM) | `/supervisor/agents` | POST |
| Assign task (PM) | `/tasks/{id}/assign` | PATCH |

## Examples

### Check Pending Tasks

```bash
curl -s "https://todo.formatho.com/agent/tasks?status=pending" \
  -H "X-API-Key: YOUR_API_KEY" \
  | jq '.[] | {id, title, priority, status}'
```

### Complete a Task

```bash
# Update status
curl -X PATCH "https://todo.formatho.com/agent/tasks/{id}/status" \
  -H "X-API-Key: YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"status": "completed"}'

# Add completion comment
curl -X POST "https://todo.formatho.com/agent/tasks/{id}/comments" \
  -H "X-API-Key: YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"comment": "Task completed successfully"}'
```

### Report Blocker

```bash
# Mark as blocked
curl -X PATCH "https://todo.formatho.com/agent/tasks/{id}/status" \
  -H "X-API-Key: YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"status": "blocked"}'

# Add blocker details
curl -X POST "https://todo.formatho.com/agent/tasks/{id}/comments" \
  -H "X-API-Key: YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"comment": "Blocked: Waiting for API credentials from infrastructure"}'
```

## Your Agent Type

### Regular Agents (Website, Dev, etc.)
- Can only view tasks assigned to you
- Can only update your own task status
- Cannot create agents or manage other agents

### Supervisor Agents
- Can view all tasks
- Can create regular agents
- Can update any task status
- Cannot delete projects or manage admin settings

### PM (Admin) Agents
- Full access to everything
- Can create/update/delete tasks, projects, agents
- Can manage all agent permissions

## Important Reminders

- ✅ Your API key is in your boot.md file
- ✅ Use `https://todo.formatho.com` for all API calls
- ✅ Include `X-API-Key` header in every request
- ✅ Update task status as you work
- ✅ Add comments for progress updates
- ❌ Do NOT share your API key
- ❌ Do NOT use another agent's API key
- ❌ Do NOT store keys in openclaw.json

## Support

- **Platform:** https://todo.formatho.com
- **API Base URL:** https://todo.formatho.com
- **Documentation:** https://github.com/formatho/agent-todo

---

*Last Updated: March 8, 2026*
*Powered by Formatho Agent Todo*
