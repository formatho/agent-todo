---
name: agent-todo
description: Manage tasks, projects, and agents in the Agent Todo platform via CLI or API
metadata:
{
  "openclaw":
  {
    "emoji": "✅",
    "homepage": "https://github.com/formatho/agent-todo",
    "requires":
    {
      "bins": ["agent-todo"],
      "anyBins": ["curl"]
    },
    "primaryEnv": "AGENT_TODO_API_KEY",
    "install":
    [
      {
        "id": "go-build",
        "kind": "go",
        "repo": "github.com/formatho/agent-todo",
        "importPath": "github.com/formatho/agent-todo/cli",
        "bins": ["agent-todo"],
        "label": "Install agent-todo CLI (go)"
      },
      {
        "id": "manual",
        "kind": "download",
        "url": "https://github.com/formatho/agent-todo/releases",
        "label": "Download from GitHub Releases"
      }
    ]
  }
}
---

# Agent Todo Management

Interact with the Agent Todo Management Platform to create and manage tasks, projects, and AI agents.

## Installation

### Quick Install (Go)
```bash
go install github.com/formatho/agent-todo/cli@latest
```

### Build from Source
```bash
git clone https://github.com/formatho/agent-todo.git
cd agent-todo/cli
make build
sudo cp bin/agent-todo /usr/local/bin/
```

### Verify Installation
```bash
agent-todo --version
agent-todo --help
```

## Authentication

### For Human Users (JWT)
```bash
agent-todo auth login email@example.com password
```

### For AI Agents (API Key)
Set the `AGENT_TODO_API_KEY` environment variable:
```bash
export AGENT_TODO_API_KEY="sk-agent-xxxxx"
```

Or in `~/.openclaw/openclaw.json`:
```json
{
  "skills": {
    "entries": {
      "agent-todo": {
        "enabled": true,
        "apiKey": "sk-agent-xxxxx",
        "env": {
          "AGENT_TODO_SERVER_URL": "http://localhost:8080"
        }
      }
    }
  }
}
```

## Configuration

The CLI uses `~/.agent-todo/config.yaml`:
```yaml
server_url: http://localhost:8080
token: eyJhbGciOiJIUzI1NiIs...
api_key: sk-agent-xxxxx
```

Override server URL:
```bash
agent-todo --server https://api.example.com project list
```

## Task Management

### Create a Task
```bash
agent-todo task create "Task title" \
  --description "Detailed description" \
  --priority high \
  --project <project-id>
```

### List Tasks
```bash
# All tasks
agent-todo task list

# Filter by status
agent-todo task list --status pending

# Filter by priority
agent-todo task list --priority high

# Filter by project
agent-todo task list --project <project-id>

# Filter by agent
agent-todo task list --agent <agent-id>

# Search
agent-todo task list --search "keyword"

# Combine filters
agent-todo task list --status pending --priority high --project <project-id>
```

### Update Task Status
```bash
# Update status
agent-todo task update <task-id> --status in_progress

# Update multiple fields
agent-todo task update <task-id> \
  --status completed \
  --priority medium \
  --description "Updated description"
```

### Delete a Task
```bash
agent-todo task delete <task-id>
```

### Assign Agents
```bash
# Assign agent to task
agent-todo task assign <task-id> <agent-id>

# Unassign agent
agent-todo task unassign <task-id>
```

### Task Comments
```bash
# Add comment
agent-todo task comment <task-id> "This is a comment"

# List comments
agent-todo task comments <task-id>
```

## Project Management

### Create Project
```bash
agent-todo project create "Project Name" \
  --description "Project description" \
  --status active
```

### List Projects
```bash
# All projects
agent-todo project list

# Filter by status
agent-todo project list --status active

# Search
agent-todo project list --search "website"
```

### Update Project
```bash
agent-todo project update <project-id> \
  --name "New name" \
  --status completed
```

### Delete Project
```bash
agent-todo project delete <project-id>
```

### Get Project Details
```bash
agent-todo project get <project-id>
```

## Agent Management

### Create Agent
```bash
agent-todo agent create "Agent Name" \
  --type openai \
  --model gpt-4

# Save the API key that's displayed!
```

### List Agents
```bash
agent-todo agent list
```

### Update Agent
```bash
agent-todo agent update <agent-id> \
  --name "New name" \
  --enabled true
```

### Delete Agent
```bash
agent-todo agent delete <agent-id>
```

## API Usage (For Agents)

When using as an agent with API key authentication, you can also use HTTP endpoints:

### Create Task
```bash
curl -X POST http://localhost:8080/tools/tasks/create \
  -H "X-API-KEY: $AGENT_TODO_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Analyze data",
    "description": "Process CSV files",
    "priority": "high",
    "project_id": "uuid"
  }'
```

### Update Task Status
```bash
curl -X POST http://localhost:8080/tools/tasks/update \
  -H "X-API-KEY: $AGENT_TODO_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "task_id": "task-uuid",
    "status": "in_progress",
    "comment": "Started working on this"
  }'
```

### List My Tasks
```bash
curl -X POST http://localhost:8080/tools/tasks/list \
  -H "X-API-KEY: $AGENT_TODO_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "pending",
    "limit": 50
  }'
```

### Get Task Status
```bash
curl -X GET http://localhost:8080/tools/tasks/status/<task-id> \
  -H "X-API-KEY: $AGENT_TODO_API_KEY"
```

## Status Values

**Task Status:**
- `pending` - Not started
- `in_progress` - Currently working
- `completed` - Finished
- `cancelled` - Cancelled
- `blocked` - Blocked/waiting

**Priority Levels:**
- `low`
- `medium` (default)
- `high`

**Project Status:**
- `active`
- `archived`
- `completed`

## Common Workflows

### 1. Delegation Workflow
```bash
# Create project
agent-todo project create "Q1 Initiatives"

# Create and assign tasks
agent-todo task create "Analyze Q1 sales" \
  --project <project-id> \
  --priority high

agent-todo task assign <task-id> <agent-id>
```

### 2. Progress Tracking
```bash
# Monitor project progress
agent-todo task list --project <project-id>

# Check agent workload
agent-todo task list --agent <agent-id>

# View task details and comments
agent-todo task get <task-id>
agent-todo task comments <task-id>
```

### 3. Status Updates
```bash
# Mark task as in progress
agent-todo task update <task-id> --status in_progress

# Add progress comment
agent-todo task comment <task-id> "Started analysis phase"

# Mark as completed
agent-todo task update <task-id> --status completed
agent-todo task comment <task-id> "Analysis complete, generated insights"
```

## Tips for AI Agents

1. **Always add comments** when updating status to provide context
2. **Use descriptive titles** and detailed descriptions for tasks
3. **Update status regularly** to keep humans informed of progress
4. **Report blockers** immediately by setting status to "blocked"
5. **Set appropriate priorities** based on task urgency
6. **Only update tasks assigned to you** (security restriction)

## Error Handling

### Common Errors

**"Not logged in"**
```bash
# Login first
agent-todo auth login user@example.com password
```

**"Connection refused"**
```bash
# Check server is running
curl http://localhost:8080/health

# Override server URL
agent-todo --server http://localhost:8080 task list
```

**"Task not assigned to this agent"** (API only)
- Agents can only update tasks assigned to them
- Verify the task is assigned to your agent ID

## Examples

### Example 1: Daily Standup Automation
```bash
# Create daily tasks
agent-todo task create "Daily standup preparation" \
  --description "Prepare metrics and blockers" \
  --priority medium

agent-todo task create "Review pull requests" \
  --priority high

agent-todo task create "Update documentation" \
  --priority low
```

### Example 2: Project Tracking
```bash
# Create project
PROJECT_ID=$(agent-todo project create "Website Redesign" | grep -oP 'ID: \K\w+')

# Create multiple tasks
agent-todo task create "Design new homepage" --project $PROJECT_ID --priority high
agent-todo task create "Implement user auth" --project $PROJECT_ID --priority high
agent-todo task create "Write API docs" --project $PROJECT_ID --priority medium

# List project tasks
agent-todo task list --project $PROJECT_ID
```

### Example 3: Agent Collaboration
```bash
# List all tasks for review
agent-todo task list --status pending

# Assign appropriate tasks
agent-todo task assign <design-task> <design-agent>
agent-todo task assign <dev-task> <dev-agent>
agent-todo task assign <qa-task> <qa-agent>

# Monitor progress
agent-todo task list --agent <design-agent>
agent-todo task list --agent <dev-agent>
```

## Support

- **Documentation**: https://github.com/formatho/agent-todo
- **API Docs**: http://localhost:8080/docs (Swagger UI)
- **Issues**: https://github.com/formatho/agent-todo/issues

## Token Impact

This skill adds approximately 200 tokens to the system prompt when loaded.
