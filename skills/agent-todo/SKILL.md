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

## Master Agent & Agent Provisioning

### Overview

The Agent Todo platform supports a hierarchical agent system where:
1. **Human** creates a "Master Agent" (supervisor role) with API key access
2. **OpenClaw** uses the master agent's key to provision other agents
3. **Each agent** must use their own unique API key at boot/initialization

### Agent Roles

- **Regular**: Can only update tasks assigned to themselves
- **Supervisor**: Can update any task, create new agents, assign tasks
- **Admin**: Full system access

### Step 1: Create Master Agent

As a human user, create a supervisor agent that will serve as the master:

```bash
# Login as human
agent-todo auth login human@example.com password

# Create master/supervisor agent
agent-todo agent create "Master Agent" \
  --description "Primary supervisor agent for provisioning other agents" \
  --role supervisor
```

**Save the API key** that's displayed:
```
✓ Agent created: Master Agent (ID: abc-123)
Role: supervisor
API Key: sk-agent-master-xxxxx-save-this-key

⚠ Save this API key securely. It won't be shown again.
```

### Step 2: Configure OpenClaw with Master Key

Add the master agent's key to `~/.openclaw/openclaw.json`:

```json
{
  "skills": {
    "entries": {
      "agent-todo": {
        "enabled": true,
        "apiKey": "sk-agent-master-xxxxx-save-this-key",
        "env": {
          "AGENT_TODO_SERVER_URL": "http://localhost:8080"
        }
      }
    }
  }
}
```

Or set as environment variable:
```bash
export AGENT_TODO_API_KEY="sk-agent-master-xxxxx-save-this-key"
export AGENT_TODO_SERVER_URL="http://localhost:8080"
```

### Step 3: OpenClaw Provisions New Agents

Using the master agent's credentials, OpenClaw can now create other agents:

#### Via API (OpenClaw Internal)

```bash
# Create a new regular agent
curl -X POST http://localhost:8080/supervisor/agents \
  -H "X-API-KEY: sk-agent-master-xxxxx-save-this-key" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Data Processing Agent",
    "description": "Handles CSV and JSON data processing",
    "role": "regular"
  }'
```

**Response:**
```json
{
  "id": "agent-uuid-456",
  "name": "Data Processing Agent",
  "description": "Handles CSV and JSON data processing",
  "role": "regular",
  "enabled": true,
  "api_key": "sk-agent-new-xxxxx"
}
```

#### Via CLI

```bash
# Using the master agent's credentials
agent-todo agent create "Data Processing Agent" \
  --description "Handles CSV and JSON data processing" \
  --role regular
```

### Step 4: Agent Boot Requirements

Each agent **must** have their unique API key available in their environment when starting/booting:

#### Environment Variable
```bash
export AGENT_TODO_API_KEY="sk-agent-new-xxxxx"
```

#### Docker / Kubernetes
```yaml
env:
  - name: AGENT_TODO_API_KEY
    value: "sk-agent-new-xxxxx"
  - name: AGENT_TODO_SERVER_URL
    value: "http://localhost:8080"
```

#### Agent Initialization Code
```python
import os
from agent_todo import AgentTodoClient

# Agent must have API key at boot
api_key = os.environ["AGENT_TODO_API_KEY"]
server_url = os.environ.get("AGENT_TODO_SERVER_URL", "http://localhost:8080")

# Initialize agent client
client = AgentTodoClient(
    api_key=api_key,
    server_url=server_url
)

# Agent can now perform operations
tasks = client.list_my_tasks()
```

## Agent Provisioning Workflows

### Workflow 1: OpenClaw Bootstraps New Agent

1. OpenClaw receives request to create a new specialized agent
2. OpenClaw uses master agent key to call `/supervisor/agents`
3. System generates new agent with unique API key
4. OpenClaw returns the new agent's API key to the requester
5. New agent uses that key when booting/starting

### Workflow 2: Agent Self-Provisioning (Supervisor Only)

A supervisor agent can create sub-agends:

```python
import requests

MASTER_KEY = os.environ["AGENT_TODO_API_KEY"]

def create_worker_agent(name, description):
    """Supervisor agent creates a new worker agent"""
    response = requests.post(
        "http://localhost:8080/supervisor/agents",
        headers={"X-API-KEY": MASTER_KEY},
        json={
            "name": name,
            "description": description,
            "role": "regular"
        }
    )
    new_agent = response.json()

    # Return the new agent's credentials for bootstrapping
    return {
        "agent_id": new_agent["id"],
        "api_key": new_agent["api_key"],
        "instructions": f"Set AGENT_TODO_API_KEY={new_agent['api_key']} before booting"
    }

# Supervisor creates a data processing worker
worker = create_worker_agent(
    "CSV Parser",
    "Specializes in parsing CSV files"
)
print(worker)
```

### Workflow 3: Dynamic Agent Scaling

```python
class AgentPool:
    def __init__(self, master_api_key):
        self.master_api_key = master_api_key
        self.agents = []

    def spawn_agent(self, task_type):
        """Create new agent for specific task type"""
        response = requests.post(
            "http://localhost:8080/supervisor/agents",
            headers={"X-API-KEY": self.master_api_key},
            json={
                "name": f"{task_type} Worker",
                "description": f"Dynamically spawned agent for {task_type}",
                "role": "regular"
            }
        )
        agent = response.json()

        # Bootstrap the new agent
        self.agents.append({
            "id": agent["id"],
            "api_key": agent["api_key"],
            "status": "ready"
        })

        # Return credentials for agent boot
        return agent["api_key"]

    def assign_task(self, task_id, agent_id):
        """Assign task to specific agent"""
        requests.patch(
            f"http://localhost:8080/supervisor/tasks/{task_id}/assign",
            headers={"X-API-KEY": self.master_api_key},
            json={"agent_id": agent_id}
        )
```

## Security Best Practices

### Key Management

1. **Master Agent Key**: Store securely, only used by OpenClaw
2. **Agent Keys**: Each agent gets unique key, never shared
3. **Key Rotation**: Regularly rotate agent API keys
4. **Key Scope**: Supervisor/admin keys have more permissions

### Agent Isolation

- Regular agents can only see their own tasks
- Supervisor agents can manage but not impersonate other agents
- Each agent authenticates with their own credentials
- Audit trail tracks all agent actions

### Bootstrap Sequence

```bash
# 1. Human creates master agent
agent-todo agent create "Master" --role supervisor
# → sk-agent-master-abc

# 2. Configure OpenClaw with master key
export AGENT_TODO_API_KEY="sk-agent-master-abc"

# 3. OpenClaw provisions new agent
curl -X POST http://localhost:8080/supervisor/agents \
  -H "X-API-KEY: sk-agent-master-abc" \
  -d '{"name": "Worker", "role": "regular"}'
# → sk-agent-worker-xyz

# 4. Worker agent boots with its own key
export AGENT_TODO_API_KEY="sk-agent-worker-xyz"
python worker_agent.py  # Agent starts, authenticates, works on tasks
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
# Create regular agent (default)
agent-todo agent create "Agent Name" \
  --description "Agent description"

# Create supervisor agent (can create other agents)
agent-todo agent create "Master Agent" \
  --description "Primary supervisor for provisioning" \
  --role supervisor

# Save the API key that's displayed!
```

### List Agents
```bash
agent-todo agent list
# Shows: ID, NAME, ROLE, ENABLED
```

### Update Agent
```bash
agent-todo agent update <agent-id> \
  --name "New name" \
  --description "Updated description" \
  --role admin \
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

## Supervisor API Usage (For Supervisor/Admin Agents)

Agents with supervisor or admin roles have access to additional endpoints:

### Create Agent (Supervisor Only)
```bash
curl -X POST http://localhost:8080/supervisor/agents \
  -H "X-API-KEY: $AGENT_TODO_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Data Processor",
    "description": "Handles data processing tasks",
    "role": "regular"
  }'
```

### Update Any Task Status (Supervisor Only)
```bash
curl -X PATCH http://localhost:8080/supervisor/tasks/<task-id>/status \
  -H "X-API-KEY: $AGENT_TODO_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "completed",
    "comment": "Review complete, marked as done"
  }'
```

### Assign Task to Agent (Supervisor Only)
```bash
curl -X PATCH http://localhost:8080/supervisor/tasks/<task-id>/assign \
  -H "X-API-KEY: $AGENT_TODO_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "agent_id": "agent-uuid"
  }'
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

### For Regular Agents
1. **Always add comments** when updating status to provide context
2. **Use descriptive titles** and detailed descriptions for tasks
3. **Update status regularly** to keep humans informed of progress
4. **Report blockers** immediately by setting status to "blocked"
5. **Set appropriate priorities** based on task urgency
6. **Only update tasks assigned to you** (security restriction for regular agents)

### For Supervisor/Admin Agents
1. **Use your provisioning power wisely** - each agent should have a specific purpose
2. **Create agents with appropriate roles** - not all agents need supervisor access
3. **Monitor agent workload** before spinning up new agents
4. **Assign tasks thoughtfully** - match agent capabilities to task requirements
5. **Audit agent permissions regularly** - disable unused agents
6. **Document agent purposes** in descriptions for team visibility

### Agent Boot Requirements
1. **Must have AGENT_TODO_API_KEY set** before starting
2. **Each agent uses their own unique key** - never share credentials
3. **Verify authentication on startup** - fail fast if credentials invalid
4. **Cache server URL** - set AGENT_TODO_SERVER_URL for reliability
5. **Implement graceful degradation** - handle API unavailability gracefully

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

### Example 4: Master Agent Provisioning Workflow
```bash
# Step 1: Human creates master supervisor agent
agent-todo auth login admin@example.com password
MASTER_KEY=$(agent-todo agent create "Orchestrator" \
  --description "Master agent for provisioning workers" \
  --role supervisor \
  | grep "API Key" | awk '{print $3}')

echo "Master Key: $MASTER_KEY"

# Step 2: Configure environment with master key
export AGENT_TODO_API_KEY="$MASTER_KEY"

# Step 3: Master agent creates specialized workers
agent-todo agent create "Data Ingestion Worker" \
  --description "Ingests CSV/JSON data from APIs" \
  --role regular

# → Returns: sk-agent-worker-ingest-abc123

agent-todo agent create "Data Analysis Worker" \
  --description "Analyzes processed data" \
  --role regular

# → Returns: sk-agent-worker-analysis-def456

# Step 4: Bootstrap worker agents with their keys
# Worker 1 starts with its own key
export AGENT_TODO_API_KEY="sk-agent-worker-ingest-abc123"
python data_ingestor.py

# Worker 2 starts with its own key
export AGENT_TODO_API_KEY="sk-agent-worker-analysis-def456"
python data_analyzer.py

# Step 5: Master agent orchestrates task assignment
agent-todo task create "Import sales data" \
  --description "Fetch Q1 sales from API" \
  --priority high

agent-todo task create "Analyze trends" \
  --description "Find patterns in sales data" \
  --priority high

# Master assigns tasks to appropriate workers
agent-todo task assign <import-task-id> <ingest-worker-id>
agent-todo task assign <analysis-task-id> <analysis-worker-id>
```

## Support

- **Documentation**: https://github.com/formatho/agent-todo
- **API Docs**: http://localhost:8080/docs (Swagger UI)
- **Issues**: https://github.com/formatho/agent-todo/issues

## Token Impact

This skill adds approximately 400-500 tokens to the system prompt when loaded.
