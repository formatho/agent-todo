# Agent Todo Manager

A comprehensive skill for managing projects, tasks, and agents in the Agent Todo Management Platform.

## Overview

This skill allows you to interact with the Agent Todo API through either:
1. **Direct API calls** (for agents with API keys)
2. **CLI tool** (for human users or local operations)

## API Endpoints (For Agents)

### Base URL
```
http://localhost:8080
```

### Authentication
- **Agents**: Use `X-API-KEY` header with your agent API key
- **Users**: Use `Authorization: Bearer <token>` header with JWT token

### Available Endpoints

#### Task Management

**Create Task**
```
POST /tools/tasks/create
Headers: { "X-API-KEY": "<your-api-key>" }
Body: {
  "title": "Task title",
  "description": "Task description (optional)",
  "priority": "high|medium|low",
  "project_id": "uuid (optional)"
}
```

**Update Task Status**
```
POST /tools/tasks/update
Headers: { "X-API-KEY": "<your-api-key>" }
Body: {
  "task_id": "task-uuid",
  "status": "pending|in_progress|completed|cancelled",
  "comment": "Optional comment about the update"
}
```

**List My Tasks**
```
POST /tools/tasks/list
Headers: { "X-API-KEY": "<your-api-key>" }
Body: {
  "status": "pending (optional)",
  "limit": 50
}
```

**Get Task Status**
```
GET /tools/tasks/status/{task_id}
Headers: { "X-API-KEY": "<your-api-key>" }
```

## CLI Tool (For Humans)

### Installation

```bash
cd cli
make build
sudo cp bin/agent-todo /usr/local/bin/
```

### Configuration

Config is stored in `~/.agent-todo/config.yaml`:

```yaml
server_url: http://localhost:8080
token: your-jwt-token
api_key: your-agent-api-key
```

### CLI Commands

#### Authentication

```bash
# Login (for human users)
agent-todo auth login email@example.com password

# Register new account
agent-todo auth register email@example.com password "Display Name"

# Check current user
agent-todo auth whoami

# Logout
agent-todo auth logout
```

#### Project Management

```bash
# Create a project
agent-todo project create "Project Name" \
  --description "Description" \
  --status active

# List all projects
agent-todo project list

# List projects with filters
agent-todo project list --status active
agent-todo project list --search "keyword"

# Get project details
agent-todo project get <project-id>

# Update a project
agent-todo project update <project-id> \
  --name "New Name" \
  --status completed

# Delete a project
agent-todo project delete <project-id>
```

#### Task Management

```bash
# Create a task
agent-todo task create "Task title" \
  --description "Task description" \
  --project <project-id> \
  --priority high

# List all tasks
agent-todo task list

# List tasks with filters
agent-todo task list --status pending
agent-todo task list --priority high
agent-todo task list --project <project-id>
agent-todo task list --agent <agent-id>
agent-todo task list --search "keyword"

# Get task details
agent-todo task get <task-id>

# Update a task
agent-todo task update <task-id> \
  --status in_progress \
  --priority high

# Delete a task
agent-todo task delete <task-id>

# Assign agent to task
agent-todo task assign <task-id> <agent-id>

# Unassign agent from task
agent-todo task unassign <task-id>

# Add comment to task
agent-todo task comment <task-id> "Comment text"

# List task comments
agent-todo task comments <task-id>
```

#### Agent Management

```bash
# Create an agent
agent-todo agent create "Agent Name" \
  --type openai \
  --model gpt-4
# Note: Save the API key that's displayed!

# List all agents
agent-todo agent list

# Get agent details
agent-todo agent get <agent-id>

# Update an agent
agent-todo agent update <agent-id> \
  --name "New Name" \
  --enabled true

# Delete an agent
agent-todo agent delete <agent-id>
```

## Usage Examples

### Example 1: Agent Creates and Updates Tasks

```python
import requests

API_KEY = "sk-agent-123"
BASE_URL = "http://localhost:8080"
headers = {"X-API-KEY": API_KEY}

# Create a task
response = requests.post(
    f"{BASE_URL}/tools/tasks/create",
    headers=headers,
    json={
        "title": "Analyze user data",
        "description": "Process the CSV files and generate insights",
        "priority": "high"
    }
)
task = response.json()
task_id = task["data"]["id"]

# Update task status
requests.post(
    f"{BASE_URL}/tools/tasks/update",
    headers=headers,
    json={
        "task_id": task_id,
        "status": "in_progress",
        "comment": "Started data analysis"
    }
)

# Mark as completed
requests.post(
    f"{BASE_URL}/tools/tasks/update",
    headers=headers,
    json={
        "task_id": task_id,
        "status": "completed",
        "comment": "Analysis complete, generated 5 insights"
    }
)
```

### Example 2: Human Uses CLI

```bash
# Login
agent-todo auth login alice@example.com password123

# Create a project
agent-todo project create "Website Redesign" --description "Redesign company website"

# Create tasks
agent-todo task create "Design homepage" --priority high
agent-todo task create "Implement user authentication" --priority high
agent-todo task create "Write documentation" --priority medium

# List pending tasks
agent-todo task list --status pending

# Update task status
agent-todo task update abc-123 --status in_progress

# Add comment
agent-todo task comment abc-123 "Started working on design mockups"
```

### Example 3: View Tasks by Project

```bash
# List tasks in a specific project
agent-todo task list --project <project-id>

# Filter by status within project
agent-todo task list --project <project-id> --status pending

# Search within project
agent-todo task list --project <project-id> --search "bug"
```

## Task Status Values

- `pending` - Task not yet started
- `in_progress` - Task currently being worked on
- `completed` - Task finished successfully
- `cancelled` - Task cancelled
- `blocked` - Task blocked/waiting for something

## Priority Levels

- `low` - Low priority
- `medium` - Medium priority (default)
- `high` - High priority

## Project Status Values

- `active` - Active project
- `archived` - Archived project
- `completed` - Completed project

## Best Practices

### For Agents

1. **Always provide comments** when updating task status
2. **Create descriptive task titles** and detailed descriptions
3. **Update status regularly** to keep humans informed
4. **Only update tasks assigned to you** (security restriction)
5. **Set appropriate priorities** based on urgency

### For Humans

1. **Use projects** to organize related tasks
2. **Set clear priorities** to help agents understand urgency
3. **Use comments** to communicate with agents
4. **Review agent comments** regularly
5. **Assign appropriate agents** based on capabilities

## Common Workflows

### Workflow 1: Delegation to Agent

```bash
# Human creates project
agent-todo project create "Data Analysis Project"

# Human creates tasks
agent-todo task create "Analyze Q1 sales data" \
  --project <project-id> \
  --priority high

# Human assigns task to agent
agent-todo task assign <task-id> <agent-id>

# Agent picks up task and updates status
# Agent: POST /tools/tasks/update { status: "in_progress", comment: "Started analysis" }

# Agent completes task
# Agent: POST /tools/tasks/update { status: "completed", comment: "Analysis complete" }
```

### Workflow 2: Multi-Agent Collaboration

```bash
# Create project with multiple tasks
agent-todo project create "Feature Launch"
agent-todo task create "Design feature" --priority high
agent-todo task create "Implement feature" --priority high
agent-todo task create "Test feature" --priority medium
agent-todo task create "Write documentation" --priority medium

# Assign to different agents
agent-todo task assign <design-task-id> <design-agent-id>
agent-todo task assign <implement-task-id> <dev-agent-id>
agent-todo task assign <test-task-id> <qa-agent-id>
agent-todo task assign <doc-task-id> <writer-agent-id>

# Monitor progress
agent-todo task list --project <project-id>
```

### Workflow 3: Task Review

```bash
# List completed tasks
agent-todo task list --status completed

# Review task comments
agent-todo task comments <task-id>

# Reopen task if needed
agent-todo task update <task-id> --status in_progress
agent-todo task comment <task-id> "Need additional work"
```

## Error Handling

### Common Errors

**401 Unauthorized**
- Cause: Invalid or missing API key/token
- Solution: Check authentication credentials

**403 Forbidden**
- Cause: Agent trying to update task not assigned to them
- Solution: Only update tasks assigned to you

**404 Not Found**
- Cause: Task/project/agent ID doesn't exist
- Solution: Verify the ID is correct

**400 Bad Request**
- Cause: Invalid request body
- Solution: Check request format and required fields

## Tips for AI Agents

1. **Be specific** in task titles and descriptions
2. **Update status frequently** to provide transparency
3. **Add meaningful comments** explaining progress and decisions
4. **Ask for clarification** if task requirements are unclear (use comments)
5. **Report blockers immediately** by setting status to "blocked"
6. **Confirm completion** by updating status to "completed" with summary

## Integration Examples

### Claude Integration

When using this skill with Claude:

```
User: "Create a task to analyze the sales data"
Agent: Creates task via API or CLI
Agent: "✓ Created task 'Analyze sales data' (ID: abc-123) with high priority"

User: "What tasks are pending?"
Agent: Lists tasks using CLI or API
Agent: "You have 3 pending tasks: ..."

User: "Mark the analysis task as in progress"
Agent: Updates task status
Agent: "✓ Task 'Analyze sales data' status updated to in_progress"
```

### Automation Scripts

```bash
#!/bin/bash
# Auto-assign high-priority tasks to available agent

PENDING_TASKS=$(agent-todo task list --status pending --priority high --format json)
AGENT_ID="agent-001"

for task in $PENDING_TASKS; do
  agent-todo task assign $task $AGENT_ID
  agent-todo task comment $task "Auto-assigned to $AGENT_ID"
done
```

## Support

- **API Documentation**: http://localhost:8080/docs (Swagger UI)
- **CLI Help**: `agent-todo --help` or `agent-todo <command> --help`
- **Issues**: https://github.com/formatho/agent-todo/issues
