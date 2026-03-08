# Agent Todo Skills

This directory contains AI/agent skills for interacting with the Agent Todo Management Platform.

## Available Skills

### Agent Todo Manager

**File:** `agent-todo-manager.md`

A comprehensive skill for managing projects, tasks, and agents. Supports both:
- **API endpoints** for AI agents with API keys
- **CLI tool** for human users or local operations

#### Capabilities

- Create, read, update, and delete projects
- Create, read, update, and delete tasks
- Assign agents to tasks
- Add comments to tasks
- Filter tasks by status, priority, project, or agent
- List and manage agents

#### Usage with Claude

When this skill is loaded, Claude can:

```
User: Create a high-priority task to analyze Q1 sales data
Claude: [Uses API/CLI] ✓ Created task "Analyze Q1 sales data" (ID: abc-123)

User: What tasks are assigned to the QA agent?
Claude: [Lists tasks] The QA agent has 5 tasks: ...

User: Mark the analysis task as in progress
Claude: [Updates status] ✓ Task status updated to in_progress
```

## Installation

### For Claude Desktop

1. Copy the skill file to Claude's skills directory:
```bash
mkdir -p ~/Library/Application\ Support/Claude/claude-desktop/skills
cp agent-todo-manager.md ~/Library/Application\ Support/Claude/claude-desktop/skills/
```

2. Restart Claude Desktop

3. The skill will be available in all conversations

### For Cursor/VSCode

1. Copy to Cursor skills directory:
```bash
mkdir -p ~/.cursor/skills
cp agent-todo-manager.md ~/.cursor/skills/
```

2. Reload Cursor

## Configuration

The skill requires either:

1. **For API usage**: An agent API key
2. **For CLI usage**: The `agent-todo` CLI installed

### CLI Installation

```bash
cd /path/to/agent-todo/cli
make build
sudo cp bin/agent-todo /usr/local/bin/
```

### Authentication

```bash
# Login as user
agent-todo auth login email@example.com password

# Or configure API key in ~/.agent-todo/config.yaml
```

## API Endpoints

The skill documents these OpenClaw-compatible endpoints:

- `POST /tools/tasks/create` - Create a new task
- `POST /tools/tasks/update` - Update task status
- `POST /tools/tasks/list` - List agent's tasks
- `GET /tools/tasks/status/{id}` - Get task status

All endpoints require `X-API-KEY` header authentication.

## Examples

### Example 1: Task Delegation

```
User: I need to analyze the customer feedback data
Claude: I'll create a task for that and assign it to the data analysis agent.

[Creates task via API]
✓ Created task "Analyze customer feedback data" (priority: high)
✓ Assigned to agent: data-analyzer-001
```

### Example 2: Progress Tracking

```
User: How are the website tasks going?
Claude: Let me check the progress of all website-related tasks.

[Queries API]
Website Redesign Project:
- Design homepage (in_progress) - Assigned to design-bot
- Implement navigation (completed) - Assigned to dev-bot
- Write content (pending) - Assigned to writer-bot

2 of 4 tasks completed (50%)
```

### Example 3: Batch Operations

```
User: Mark all completed testing tasks as done
Claude: I'll update all testing tasks that are marked as complete.

[Updates multiple tasks]
✓ Updated 5 testing tasks to "completed" status
✓ Added completion comments to each task
```

## Agent Capabilities

When this skill is active, AI agents can:

1. **Create tasks** with titles, descriptions, priorities
2. **Update task status** (pending → in_progress → completed)
3. **Add comments** to communicate progress
4. **List their assigned tasks** with filters
5. **Query task details** and status

### Security Restrictions

- Agents can **only update tasks assigned to them**
- Agents cannot modify tasks assigned to other agents
- Agents cannot delete tasks or projects
- All actions are logged with agent identity

## CLI vs API

| Feature | CLI | API |
|---------|-----|-----|
| Target | Human users | AI agents |
| Auth | JWT token | API key |
| Scope | All operations | Agent-specific |
| Access | Full access | Restricted to assigned tasks |
| Use case | Manual management | Automated workflows |

## Troubleshooting

### CLI Not Found

```bash
# Check installation
which agent-todo

# Reinstall if needed
cd cli && make build && sudo cp bin/agent-todo /usr/local/bin/
```

### Authentication Errors

```bash
# Verify credentials
agent-todo auth whoami

# Re-login if needed
agent-todo auth login email@example.com password
```

### API Key Issues

```bash
# Check agent API key
agent-todo agent get <agent-id>

# Create new agent if key lost
agent-todo agent create "Agent Name" --type openai --model gpt-4
```

## Contributing

To add new skills or improve existing ones:

1. Create a new `.md` file in this directory
2. Follow the skill structure format
3. Include examples and usage instructions
4. Test with Claude Desktop/Cursor
5. Submit a pull request

## License

Part of the Agent Todo Management Platform.
See main project LICENSE for details.
