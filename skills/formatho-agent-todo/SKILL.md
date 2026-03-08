---
name: formatho-agent-todo
description: Manage tasks, projects, and agents in the Formatho Agent Todo platform via CLI or API
metadata: {"openclaw":{"emoji":"✅","homepage":"https://github.com/formatho/agent-todo","requires":{"bins":["agent-todo"],"anyBins":["curl"]},"primaryEnv":"AGENT_TODO_API_KEY","install":[{"id":"go-build","kind":"go","repo":"github.com/formatho/agent-todo","importPath":"github.com/formatho/agent-todo/cli","bins":["agent-todo"],"label":"Install agent-todo CLI (go)"},{"id":"manual","kind":"download","url":"https://github.com/formatho/agent-todo/releases","label":"Download from GitHub Releases"}]}}
---

# Formatho Agent Todo Management

Interact with the Formatho Agent Todo Management Platform to create and manage tasks, projects, and AI agents.

## Installation

### Prerequisite Check

**Before using this skill, always verify the CLI is installed:**

```bash
# Check if agent-todo CLI is installed
if ! command -v agent-todo &> /dev/null; then
    echo "agent-todo CLI not found. Installing..."
    # Auto-install using Go
    go install github.com/formatho/agent-todo/cli@latest
    # Verify installation
    if command -v agent-todo &> /dev/null; then
        echo "✓ agent-todo CLI installed successfully"
        agent-todo --version
    else
        echo "✗ Installation failed. Please install manually:"
        echo "  go install github.com/formatho/agent-todo/cli@latest"
        exit 1
    fi
else
    echo "✓ agent-todo CLI is already installed"
    agent-todo --version
fi
```

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

### OpenClaw Plugin Configuration

When using the Formatho Agent Todo OpenClaw plugin, configure the server URL in `~/.openclaw/openclaw.json`:

```json
{
  "plugins": {
    "entries": {
      "formatho-agent-todo": {
        "enabled": true,
        "config": {
          "serverUrl": "http://localhost:8080",
          "apiKey": "sk-agent-xxxxx",
          "autoInstall": true
        }
      }
    }
  }
}
```

**Required Configuration:**
- `serverUrl`: Your Formatho Agent Todo server URL (e.g., `http://localhost:8080` or `https://todo.example.com`)
- `apiKey`: API key for agent authentication (optional for humans, required for agents)

**Plugin Configuration Options:**
| Option | Type | Required | Default | Description |
|--------|------|----------|---------|-------------|
| `enabled` | boolean | Yes | `false` | Enable/disable the plugin |
| `serverUrl` | string | **Yes** | `http://localhost:8080` | Formatho Agent Todo server URL |
| `apiKey` | string | No | `""` | API key for agent authentication |
| `autoInstall` | boolean | No | `true` | Auto-install CLI if missing |

For full plugin installation instructions, see `../formatho-agent-todo-plugin/README.md`

## Master Agent & Agent Provisioning

### Overview

The Agent Todo platform supports a hierarchical agent system where:
1. **Human** creates a "Project Manager (PM) Agent" (admin role) with all-access API key
2. **PM agent** is responsible for creating API keys for all subordinate agents
3. **Each agent** must use their own unique API key at boot/initialization, loaded from boot.md

### Agent Roles and Permissions

- **Project Manager (PM)**: Full system access (all-access API key)
  - Can create, update, delete any task, project, or agent
  - Can provision new agents and generate their API keys
  - Can assign tasks to any agent
  - Responsible for managing team of agents

- **Supervisor**: Can manage tasks and create regular agents
  - Can update any task status
  - Can create regular agents (not supervisor/admin)
  - Can assign tasks to agents
  - Responsible for provisioning and managing sub-ordinate agents

- **Regular**: Self-write access only (restricted API key)
  - Can only view and update tasks assigned to themselves
  - Cannot create other agents
  - Cannot view or modify other agents' tasks
  - Must use their own unique API key

### Step 1: Create Project Manager (PM) Agent

As a human user, create a Project Manager agent with full access:

```bash
# Create PM agent with all-access
agent-todo agent create "Project Manager" \
  --description "Primary PM agent for managing all agents and tasks" \
  --role admin
```

**Save the API key** that's displayed:
```
✓ Agent created: Project Manager (ID: pm-abc-123)
Role: admin
API Key: sk-agent-pm-all-access-xxxxx-save-this-key

⚠ Save this API key securely. It won't be shown again.
```

### Step 2: Configure OpenClaw with PM Key

**IMPORTANT:** PM agent must save its API key in OpenClaw's agent boot configuration.

Add the PM agent's key to `~/.openclaw/agents/<pm-agent-id>/boot.md`:

```markdown
# Agent Boot Configuration

## Agent Todo Credentials

- **Agent ID**: pm-abc-123
- **Agent Name**: Project Manager
- **Role**: admin
- **API Key**: sk-agent-pm-all-access-xxxxx-save-this-key
- **Server URL**: http://localhost:8080

## Environment Variables

```bash
export AGENT_TODO_API_KEY="sk-agent-pm-all-access-xxxxx-save-this-key"
export AGENT_TODO_SERVER_URL="http://localhost:8080"
```

## Permissions

- Full system access
- Can create/update/delete any task, project, or agent
- Can provision new agents
- Responsible for managing subordinate agents
```

Alternatively, add to `~/.openclaw/openclaw.json`:

```json
{
  "skills": {
    "entries": {
      "agent-todo": {
        "enabled": true,
        "apiKey": "sk-agent-pm-all-access-xxxxx-save-this-key",
        "env": {
          "AGENT_TODO_SERVER_URL": "http://localhost:8080"
        }
      }
    }
  }
}
```

### Step 3: PM Agent Provisions New Agents

**PM/S supervisor agent responsibility:** The PM agent is responsible for creating API keys for all subordinate agents.

Using the PM agent's all-access credentials, OpenClaw can now create other agents:

#### Via API (OpenClaw Internal)

```bash
# PM agent creates a new regular agent
curl -X POST http://localhost:8080/supervisor/agents \
  -H "X-API-KEY: sk-agent-pm-all-access-xxxxx-save-this-key" \
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
# Using the PM agent's credentials
export AGENT_TODO_API_KEY="sk-agent-pm-all-access-xxxxx-save-this-key"

agent-todo agent create "Data Processing Agent" \
  --description "Handles CSV and JSON data processing" \
  --role regular
```

### Step 4: Save New Agent Key in Boot Configuration

**PM agent responsibility:** After creating a new agent, the PM agent must save the new agent's API key in OpenClaw's agent boot configuration.

Create or update `~/.openclaw/agents/<new-agent-id>/boot.md`:

```markdown
# Agent Boot Configuration

## Agent Todo Credentials

- **Agent ID**: agent-uuid-456
- **Agent Name**: Data Processing Agent
- **Role**: regular
- **API Key**: sk-agent-new-xxxxx
- **Server URL**: http://localhost:8080

## Environment Variables

```bash
export AGENT_TODO_API_KEY="sk-agent-new-xxxxx"
export AGENT_TODO_SERVER_URL="http://localhost:8080"
```

## Permissions

- Self-write access only
- Can only update tasks assigned to this agent
- Cannot create other agents
- Cannot view other agents' tasks

## Boot Instructions

1. Load this boot.md file to get credentials
2. Set AGENT_TODO_API_KEY environment variable
3. Verify authentication with API
4. Start processing assigned tasks
```

### Step 5: Agent Boot Requirements

Each agent **must** load their unique API key from their boot.md file when starting/booting:

#### Environment Variable
```bash
# Load from boot.md
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

# Agent loads API key from boot.md at boot time
api_key = os.environ["AGENT_TODO_API_KEY"]
server_url = os.environ.get("AGENT_TODO_SERVER_URL", "http://localhost:8080")

# Initialize agent client
client = AgentTodoClient(
    api_key=api_key,
    server_url=server_url
)

# Verify authentication and permissions
me = client.get_agent_info()
print(f"Agent {me['name']} (Role: {me['role']}) authenticated successfully")

# Agent can now perform operations within permissions
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

1. **PM Agent Key**: All-access key, stored in PM agent's boot.md, used for provisioning
2. **Supervisor Keys**: Can create regular agents, stored in their boot.md
3. **Regular Agent Keys**: Self-write only, stored in their boot.md, never shared
4. **Key Creation**: PM/supervisor agents are responsible for creating keys for subordinate agents
5. **Key Storage**: All agents must save their keys in `~/.openclaw/agents/<agent-id>/boot.md`
6. **Key Rotation**: PM agent can rotate keys by creating new ones and updating boot.md files
7. **Key Scope**: Permissions are tied to agent role (PM=all, supervisor=manage, regular=self)

### Agent Isolation

- Regular agents can only see their own tasks
- Supervisor agents can manage but not impersonate other agents
- Each agent authenticates with their own credentials
- Audit trail tracks all agent actions

### Bootstrap Sequence

```bash
# 1. Human creates PM agent with all-access
PM_KEY=$(agent-todo agent create "Project Manager" --role admin | grep "API Key" | awk '{print $3}')
# → sk-agent-pm-all-access-abc

# 2. PM saves its key in boot.md
mkdir -p ~/.openclaw/agents/pm-abc-123
cat > ~/.openclaw/agents/pm-abc-123/boot.md << EOF
# Agent Boot Configuration
## Agent Todo Credentials
- **API Key**: $PM_KEY
- **Role**: admin
## Environment Variables
export AGENT_TODO_API_KEY="$PM_KEY"
EOF

# 3. PM agent provisions new regular agent
export AGENT_TODO_API_KEY="$PM_KEY"
WORKER_KEY=$(agent-todo agent create "Worker" --role regular | grep "API Key" | awk '{print $3}')
# → sk-agent-worker-xyz

# 4. PM saves worker's key in worker's boot.md
mkdir -p ~/.openclaw/agents/worker-uuid-456
cat > ~/.openclaw/agents/worker-uuid-456/boot.md << EOF
# Agent Boot Configuration
## Agent Todo Credentials
- **API Key**: $WORKER_KEY
- **Role**: regular
## Environment Variables
export AGENT_TODO_API_KEY="$WORKER_KEY"
EOF

# 5. Worker agent boots with its own key from boot.md
source ~/.openclaw/agents/worker-uuid-456/boot.md
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

### Get Task Details
```bash
agent-todo task get <task-id>
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

### Get Agent Details
```bash
agent-todo agent get <agent-id>
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

### For Project Manager (PM) Agents
1. **Create API keys for all subordinate agents** - this is your primary responsibility
2. **Save each agent's key in their boot.md** - ensure proper boot configuration
3. **Use your all-access power wisely** - create agents only when needed
4. **Monitor agent workload** before spinning up new agents
5. **Assign tasks thoughtfully** - match agent capabilities to task requirements
6. **Audit agent permissions regularly** - disable unused agents, rotate keys if needed
7. **Document agent purposes** in descriptions for team visibility

### For Supervisor Agents
1. **Create regular agents for specific tasks** - you cannot create other supervisors
2. **Save each agent's key in their boot.md** - ensure proper boot configuration
3. **Monitor your sub-ordinate agents** - track their performance and task completion
4. **Assign tasks thoughtfully** - match agent capabilities to task requirements
5. **Document agent purposes** in descriptions for visibility

### For Regular Agents
1. **Always add comments** when updating status to provide context
2. **Use descriptive titles** and detailed descriptions for tasks
3. **Update status regularly** to keep humans informed of progress
4. **Report blockers** immediately by setting status to "blocked"
5. **Set appropriate priorities** based on task urgency
6. **Only update tasks assigned to you** (security restriction for regular agents)
7. **Load your API key from boot.md** at boot time

### Agent Boot Requirements
1. **Must have AGENT_TODO_API_KEY set** before starting (from boot.md)
2. **Each agent uses their own unique key** - never share credentials
3. **Verify authentication on startup** - fail fast if credentials invalid
4. **Cache server URL** - set AGENT_TODO_SERVER_URL for reliability
5. **Implement graceful degradation** - handle API unavailability gracefully
6. **Load credentials from boot.md** - PM/Supervisor agents create this for you

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

### Example 4: PM Agent Provisioning Workflow
```bash
# Step 1: Human creates PM agent with all-access
PM_KEY=$(agent-todo agent create "Project Manager" \
  --description "Primary PM agent for managing all agents and tasks" \
  --role admin \
  | grep "API Key" | awk '{print $3}')

echo "PM Key: $PM_KEY"

# Step 2: PM saves its key in boot.md
mkdir -p ~/.openclaw/agents/pm-abc-123
cat > ~/.openclaw/agents/pm-abc-123/boot.md << EOF
# Agent Boot Configuration
## Agent Todo Credentials
- **API Key**: $PM_KEY
- **Role**: admin
## Environment Variables
export AGENT_TODO_API_KEY="$PM_KEY"
EOF

# Step 3: PM agent provisions new regular agent
export AGENT_TODO_API_KEY="$PM_KEY"
WORKER_KEY=$(agent-todo agent create "Worker" \
  --description "Handles task execution" \
  --role regular \
  | grep "API Key" | awk '{print $3}')

# → Returns: sk-agent-worker-xyz

# Step 4: PM saves worker's key in worker's boot.md
mkdir -p ~/.openclaw/agents/worker-uuid-456
cat > ~/.openclaw/agents/worker-uuid-456/boot.md << EOF
# Agent Boot Configuration
## Agent Todo Credentials
- **API Key**: $WORKER_KEY
- **Role**: regular
## Environment Variables
export AGENT_TODO_API_KEY="$WORKER_KEY"
EOF

# Step 5: Worker agent boots with its own key from boot.md
source ~/.openclaw/agents/worker-uuid-456/boot.md
python worker_agent.py  # Agent starts, authenticates, works on tasks
```

## Support

- **Documentation**: https://github.com/formatho/agent-todo
- **API Docs**: http://localhost:8080/docs (Swagger UI)
- **Issues**: https://github.com/formatho/agent-todo/issues

## Token Impact

This skill adds approximately 400-500 tokens to the system prompt when loaded.
