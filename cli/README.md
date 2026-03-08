# Agent Todo CLI

A comprehensive CLI tool for managing projects, tasks, and AI agents in the Agent Todo Management Platform.

## Table of Contents

- [Installation](#installation)
- [Configuration](#configuration)
- [Quick Start](#quick-start)
- [Commands Reference](#commands-reference)
  - [Authentication](#authentication)
  - [Projects](#projects)
  - [Tasks](#tasks)
  - [Agents](#agents)
- [Shell Completion](#shell-completion)
- [Development](#development)
- [Examples](#examples)
- [Troubleshooting](#troubleshooting)

---

## Installation

### From Source

```bash
cd cli
make install
```

This builds and installs the CLI to your `$GOPATH/bin` or `$HOME/go/bin` directory. Make sure this directory is in your `PATH`.

### Pre-built Binaries

Download the appropriate binary for your platform:

```bash
# Linux (amd64)
wget https://releases.example.com/agent-todo-linux-amd64 -O agent-todo
chmod +x agent-todo
sudo mv agent-todo /usr/local/bin/

# macOS (Intel)
wget https://releases.example.com/agent-todo-darwin-amd64 -O agent-todo
chmod +x agent-todo
sudo mv agent-todo /usr/local/bin/

# macOS (Apple Silicon)
wget https://releases.example.com/agent-todo-darwin-arm64 -O agent-todo
chmod +x agent-todo
sudo mv agent-todo /usr/local/bin/

# Windows (amd64)
wget https://releases.example.com/agent-todo-windows-amd64.exe -O agent-todo.exe
# Add to PATH
```

### Build from Source

```bash
# Build for current platform
make build

# Build for all platforms
make build-all
```

Binaries are created in the `bin/` directory.

### Verify Installation

```bash
agent-todo --version
agent-todo --help
```

---

## Configuration

The CLI stores its configuration in `~/.agent-todo/config.yaml`. This file is automatically created on first use.

### Configuration File Structure

```yaml
server_url: http://localhost:8080
token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
api_key: ""
insecure: false
```

### Configuration Options

| Option | Type | Description | Default |
|--------|------|-------------|---------|
| `server_url` | string | The API server URL | `http://localhost:8080` |
| `token` | string | JWT authentication token (auto-set on login) | - |
| `api_key` | string | API key for agent authentication | - |
| `insecure` | boolean | Skip SSL verification | `false` |

### Environment Variables

You can also use environment variables:

- `AGENT_TODO_SERVER_URL` - Override server URL
- `AGENT_TODO_TOKEN` - Override auth token
- `AGENT_TODO_API_KEY` - Override API key
- `AGENT_TODO_INSECURE` - Skip SSL verification

### Command-Line Flags

Override configuration with command-line flags:

```bash
# Override server URL
agent-todo -s https://api.example.com project list

# Verbose output
agent-todo -v task list
```

| Flag | Short | Description |
|------|-------|-------------|
| `--server` | `-s` | Server URL (overrides config file) |
| `--verbose` | `-v` | Enable verbose output |
| `--help` | `-h` | Show help for command |

---

## Quick Start

### 1. Login to the Platform

```bash
agent-todo auth login user@example.com password123
```

### 2. Create Your First Project

```bash
agent-todo project create "My First Project" \
  --description "Getting started with Agent Todo"
```

### 3. Create an AI Agent

```bash
agent-todo agent create "Developer Bot" \
  --type openai \
  --model gpt-4
```

Save the API key that's displayed!

### 4. Create Tasks

```bash
# Create a high-priority task
agent-todo task create "Fix critical bug" \
  --description "Login form not working" \
  --priority high

# Create a task in a project
agent-todo task create "Design homepage" \
  --project <project-id>
```

### 5. Assign Tasks to Agents

```bash
agent-todo task assign <task-id> <agent-id>
```

### 6. Track Progress

```bash
# List all tasks
agent-todo task list

# List tasks in a project
agent-todo task list --project <project-id>

# Update task status
agent-todo task update <task-id> --status "in_progress"
```

---

## Commands Reference

### Authentication

Commands for user authentication and session management.

#### `auth login`

Login to the platform and save authentication token.

```bash
agent-todo auth login <email> <password>
```

**Example:**
```bash
agent-todo auth login alice@example.com secretpassword
# Output: ✓ Logged in successfully as alice@example.com (Alice)
```

#### `auth register`

Register a new user account.

```bash
agent-todo auth register <email> <password> [display-name]
```

**Example:**
```bash
agent-todo auth register bob@example.com password123 "Bob Smith"
# Output: ✓ Registered and logged in successfully as bob@example.com
```

#### `auth whoami`

Display current user information.

```bash
agent-todo auth whoami
```

**Example:**
```bash
$ agent-todo auth whoami
Email: alice@example.com
Name: Alice Johnson
ID: 550e8400-e29b-41d4-a716-446655440000
Joined: 2024-01-15T10:30:00Z
```

#### `auth logout`

Logout and clear stored credentials.

```bash
agent-todo auth logout
```

---

### Projects

Commands for managing projects.

#### `project create`

Create a new project.

```bash
agent-todo project create <name> [flags]
```

**Flags:**
- `-d, --description string` - Project description
- `-s, --status string` - Project status (e.g., "active", "archived", "completed")

**Examples:**
```bash
# Basic project
agent-todo project create "Website Redesign"

# Full project with details
agent-todo project create "Mobile App" \
  --description "Build iOS and Android app" \
  --status active
```

#### `project list`

List all projects with optional filtering.

```bash
agent-todo project list [flags]
```

**Flags:**
- `-s, --status string` - Filter by status
- `-q, --search string` - Search by name or description

**Examples:**
```bash
# List all projects
agent-todo project list

# List only active projects
agent-todo project list --status active

# Search projects
agent-todo project list --search "website"

# Output format:
# ID                 NAME                STATUS    DESCRIPTION
# abc123-def456      Website Redesign    active    Redesign company website
# 789ghi-012jkl      Mobile App          active    Build mobile app
#
# Total: 2 project(s)
```

#### `project get`

Get detailed information about a specific project.

```bash
agent-todo project get <id>
```

**Example:**
```bash
$ agent-todo project get abc123-def456
ID:          abc123-def456
Name:        Website Redesign
Description: Redesign company website
Status:      active
Created:     2024-01-15T10:30:00Z
Updated:     2024-01-20T15:45:00Z
```

#### `project update`

Update an existing project.

```bash
agent-todo project update <id> [flags]
```

**Flags:**
- `-n, --name string` - New project name
- `-d, --description string` - New project description
- `-s, --status string` - New project status

**Examples:**
```bash
# Update status
agent-todo project update abc123-def456 --status completed

# Update name and description
agent-todo project update abc123-def456 \
  --name "Website Redesign Phase 2" \
  --description "Second phase of redesign"
```

#### `project delete`

Delete a project.

```bash
agent-todo project delete <id>
```

**Warning:** This will also delete all tasks associated with the project.

---

### Tasks

Commands for managing tasks.

#### `task create`

Create a new task.

```bash
agent-todo task create <title> [flags]
```

**Flags:**
- `-d, --description string` - Task description
- `-p, --project string` - Project ID to add task to
- `-P, --priority string` - Task priority (low, medium, high)

**Examples:**
```bash
# Simple task
agent-todo task create "Review pull request"

# Full task with all options
agent-todo task create "Fix authentication bug" \
  --description "Users cannot login with SSO" \
  --project abc123-def456 \
  --priority high
```

#### `task list`

List tasks with filtering options.

```bash
agent-todo task list [flags]
```

**Flags:**
- `-s, --status string` - Filter by status (pending, in_progress, completed, etc.)
- `-P, --priority string` - Filter by priority (low, medium, high)
- `-p, --project string` - Filter by project ID
- `-a, --agent string` - Filter by assigned agent ID
- `-q, --search string` - Search by title or description

**Examples:**
```bash
# List all tasks
agent-todo task list

# List high-priority tasks
agent-todo task list --priority high

# List tasks in a project
agent-todo task list --project abc123-def456

# List pending tasks
agent-todo task list --status pending

# Search tasks
agent-todo task list --search "authentication"

# Complex filter
agent-todo task list \
  --status pending \
  --priority high \
  --project abc123-def456

# Output format:
# ID                 TITLE                   STATUS      PRIORITY    PROJECT
# xyz789-uvw012      Fix auth bug            pending     high        abc123-def456
# def456-ghi789      Update documentation    in_prog    medium      abc123-def456
#
# Total: 2 task(s)
```

#### `task get`

Get detailed information about a specific task.

```bash
agent-todo task get <id>
```

**Example:**
```bash
$ agent-todo task get xyz789-uvw012
ID:          xyz789-uvw012
Title:       Fix authentication bug
Description: Users cannot login with SSO
Status:      pending
Priority:    high
Project ID:  abc123-def456
Agent ID:    agent-001
Created:     2024-01-15T10:30:00Z
Updated:     2024-01-20T15:45:00Z
```

#### `task update`

Update an existing task.

```bash
agent-todo task update <id> [flags]
```

**Flags:**
- `-t, --title string` - New task title
- `-d, --description string` - New task description
- `-s, --status string` - New task status
- `-P, --priority string` - New task priority

**Examples:**
```bash
# Update status
agent-todo task update xyz789-uvw012 --status in_progress

# Update multiple fields
agent-todo task update xyz789-uvw012 \
  --title "Fix SSO authentication bug" \
  --priority high \
  --status in_progress
```

#### `task delete`

Delete a task.

```bash
agent-todo task delete <id>
```

#### `task assign`

Assign an agent to a task.

```bash
agent-todo task assign <task-id> <agent-id>
```

**Example:**
```bash
agent-todo task assign xyz789-uvw012 agent-001
# Output: ✓ Agent assigned to task
```

#### `task unassign`

Unassign an agent from a task.

```bash
agent-todo task unassign <task-id>
```

**Example:**
```bash
agent-todo task unassign xyz789-uvw012
# Output: ✓ Agent unassigned from task
```

#### `task comment`

Add a comment to a task.

```bash
agent-todo task comment <task-id> <content>
```

**Example:**
```bash
agent-todo task comment xyz789-uvw012 "Started investigating the issue"
# Output: ✓ Comment added
```

#### `task comments`

List all comments on a task.

```bash
agent-todo task comments <task-id>
```

**Example:**
```bash
$ agent-todo task comments xyz789-uvw012
[2024-01-15T10:30:00Z] Started investigating the issue
[2024-01-15T11:45:00Z] Found the bug in the SSO handler
[2024-01-15T14:20:00Z] Fix deployed to staging

Total: 3 comment(s)
```

---

### Agents

Commands for managing AI agents.

#### `agent create`

Create a new AI agent.

```bash
agent-todo agent create <name> [flags]
```

**Flags:**
- `-t, --type string` - Agent type (e.g., "openai", "anthropic", "custom")
- `-m, --model string` - Model identifier (e.g., "gpt-4", "claude-3")

**Examples:**
```bash
# Create OpenAI agent
agent-todo agent create "Developer Bot" \
  --type openai \
  --model gpt-4

# Output:
# ✓ Agent created: Developer Bot (ID: agent-001)
# API Key: sk-agent-abc123xyz789...
#
# ⚠ Save this API key securely. It won't be shown again.
```

**Important:** Save the API key securely! It won't be displayed again.

#### `agent list`

List all agents.

```bash
agent-todo agent list
```

**Example:**
```bash
$ agent-todo agent list
ID                 NAME                TYPE        MODEL        ENABLED
agent-001          Developer Bot       openai      gpt-4        ✓
agent-002          Writer Bot          anthropic   claude-3     ✓
agent-003          Tester Bot          custom      test-model   ✗

Total: 3 agent(s)
```

#### `agent get`

Get detailed information about an agent.

```bash
agent-todo agent get <id>
```

**Example:**
```bash
$ agent-todo agent get agent-001
ID:        agent-001
Name:      Developer Bot
Type:      openai
Model:     gpt-4
Enabled:   true
Created:   2024-01-15T10:30:00Z
Updated:   2024-01-20T15:45:00Z
```

#### `agent update`

Update an existing agent.

```bash
agent-todo agent update <id> [flags]
```

**Flags:**
- `-n, --name string` - New agent name
- `-t, --type string` - New agent type
- `-m, --model string` - New model identifier
- `-e, --enabled bool` - Enable or disable agent

**Examples:**
```bash
# Update name
agent-todo agent update agent-001 --name "Senior Developer Bot"

# Update model
agent-todo agent update agent-001 --model gpt-4-turbo

# Disable agent
agent-todo agent update agent-001 --enabled false

# Enable agent
agent-todo agent update agent-001 --enabled true
```

#### `agent delete`

Delete an agent.

```bash
agent-todo agent delete <id>
```

**Warning:** This will also unassign the agent from all tasks.

---

## Shell Completion

The CLI supports shell auto-completion for bash, zsh, fish, and PowerShell.

### Installation

#### Bash

```bash
# Add to ~/.bashrc or ~/.bash_profile
source <(agent-todo completion bash)

# Or generate completion file
agent-todo completion bash > /etc/bash_completion.d/agent-todo
```

#### Zsh

```bash
# Add to ~/.zshrc
source <(agent-todo completion zsh)

# Or for oh-my-zsh users
agent-todo completion zsh > ~/.zsh/completion/_agent-todo
```

#### Fish

```bash
# Add to ~/.config/fish/completions/agent-todo.fish
agent-todo completion fish > ~/.config/fish/completions/agent-todo.fish
```

#### PowerShell

```powershell
# Add to PowerShell profile
agent-todo completion powershell | Out-String | Invoke-Expression

# Or save to profile
agent-todo completion powershell > agent-todo.ps1
# Then add: . agent-todo.ps1 to your profile
```

### Usage

Once installed, use tab completion:

```bash
agent-todo <TAB>
# Shows: agent, auth, project, task

agent-todo task <TAB>
# Shows: assign, comment, comments, create, delete, get, list, unassign, update

agent-todo task update <TAB>
# Shows flags: --description, --priority, --status, --title
```

---

## Development

### Prerequisites

- Go 1.21 or later
- make (optional, for using Makefile)

### Building

```bash
# Build for current platform
make build
# or
go build -o bin/agent-todo .

# Build for all platforms
make build-all

# Development build with race detection
make dev
```

### Testing

```bash
# Run tests
make test

# Run tests with coverage
make test-coverage
```

### Code Quality

```bash
# Format code
make fmt

# Run linter
make lint

# Run both
make check
```

### Project Structure

```
cli/
├── main.go              # Entry point
├── cmd/                 # Command implementations
│   ├── root.go         # Root command
│   ├── auth.go         # Authentication commands
│   ├── project.go      # Project commands
│   ├── task.go         # Task commands
│   └── agent.go        # Agent commands
├── client/             # API client
│   └── client.go       # HTTP client wrapper
├── config/             # Configuration management
│   └── config.go       # Config file handling
├── Makefile            # Build commands
├── go.mod              # Go module definition
└── README.md           # This file
```

---

## Examples

### Complete Workflow Example

```bash
# 1. Login
agent-todo auth login alice@example.com password123

# 2. Create a project
PROJECT_ID=$(agent-todo project create "Website Redesign" \
  --description "Redesign company website and CMS" \
  --status active | grep -oP 'ID: \K\w+(?=-\w+)')

echo "Project ID: $PROJECT_ID"

# 3. Create an AI agent
agent-todo agent create "Developer Bot" \
  --type openai \
  --model gpt-4

# Save the agent ID from output
AGENT_ID="agent-001"

# 4. Create multiple tasks
agent-todo task create "Design new homepage" \
  --project $PROJECT_ID \
  --priority high

agent-todo task create "Implement user authentication" \
  --project $PROJECT_ID \
  --priority high

agent-todo task create "Write API documentation" \
  --project $PROJECT_ID \
  --priority medium

# 5. List tasks in project
agent-todo task list --project $PROJECT_ID

# 6. Assign agent to first task
TASK_ID="xyz789-uvw012"
agent-todo task assign $TASK_ID $AGENT_ID

# 7. Add a comment
agent-todo task comment $TASK_ID "Agent assigned to work on this task"

# 8. Update task status
agent-todo task update $TASK_ID --status in_progress

# 9. View task details
agent-todo task get $TASK_ID

# 10. Mark task as completed
agent-todo task update $TASK_ID --status completed

# 11. View task comments
agent-todo task comments $TASK_ID

# 12. Logout when done
agent-todo auth logout
```

### Managing Multiple Projects

```bash
# Create projects
agent-todo project create "Website Project" --status active
agent-todo project create "Mobile App Project" --status active

# List all projects
agent-todo project list

# Create tasks in different projects
agent-todo task create "Fix homepage bug" --project <website-id>
agent-todo task create "Add login screen" --project <mobile-id>

# View tasks by project
agent-todo task list --project <website-id>
agent-todo task list --project <mobile-id>
```

### Agent Workflow

```bash
# Create specialized agents
agent-todo agent create "Frontend Dev" --type openai --model gpt-4
agent-todo agent create "Backend Dev" --type anthropic --model claude-3
agent-todo agent create "QA Tester" --type openai --model gpt-4

# List agents
agent-todo agent list

# Create tasks
FRONTEND_TASK=$(agent-todo task create "Design React component" --priority high | grep ID)
BACKEND_TASK=$(agent-todo task create "Build API endpoint" --priority high | grep ID)

# Assign to appropriate agents
agent-todo task assign $FRONTEND_TASK <frontend-agent-id>
agent-todo task assign $BACKEND_TASK <backend-agent-id>

# Monitor progress
agent-todo task list --agent <frontend-agent-id>
agent-todo task list --agent <backend-agent-id>
```

### Filtering and Search

```bash
# Find all high-priority pending tasks
agent-todo task list --status pending --priority high

# Search for specific tasks
agent-todo task list --search "authentication"

# Find tasks assigned to a specific agent
agent-todo task list --agent agent-001

# View completed tasks in a project
agent-todo task list --project <project-id> --status completed

# Find all tasks across all projects
agent-todo task list
```

### Task Comments Collaboration

```bash
# Create a task
TASK_ID=$(agent-todo task create "Investigate performance issue" | grep ID)

# Add initial comment
agent-todo task comment $TASK_ID "Issue reported by user: slow load times"

# Team member adds comment
agent-todo task comment $TASK_ID "Found that database queries are not optimized"

# Agent adds comment
agent-todo task comment $TASK_ID "Created indexes on frequently queried columns"

# View conversation
agent-todo task comments $TASK_ID

# Mark as resolved
agent-todo task update $TASK_ID --status completed
agent-todo task comment $TASK_ID "Performance improved by 80%"
```

---

## Troubleshooting

### "Not logged in" Error

**Problem:** You're trying to execute a command but get a "not logged in" error.

**Solution:**
```bash
# Login first
agent-todo auth login your@email.com password

# Verify you're logged in
agent-todo auth whoami
```

### Connection Refused

**Problem:** Cannot connect to the API server.

**Solutions:**

1. Check if the server is running:
```bash
curl http://localhost:8080/health
```

2. Verify the server URL in your config:
```bash
cat ~/.agent-todo/config.yaml
```

3. Override the server URL:
```bash
agent-todo -s http://localhost:8080 project list
```

4. Check your network connection and firewall settings.

### Invalid Credentials

**Problem:** Login fails with "invalid credentials".

**Solutions:**

1. Verify email and password:
```bash
agent-todo auth login your@email.com password
```

2. Register a new account if needed:
```bash
agent-todo auth register your@email.com password "Your Name"
```

3. Check if the server is running and accessible.

### Permission Denied

**Problem:** You get a permission error when trying to access/modify a resource.

**Solutions:**

1. Verify you're logged in:
```bash
agent-todo auth whoami
```

2. Check if you have the necessary permissions for the resource.

3. Ensure the resource belongs to your account or project.

### Config File Issues

**Problem:** CLI is not using the correct configuration.

**Solutions:**

1. View current config:
```bash
cat ~/.agent-todo/config.yaml
```

2. Reset config:
```bash
rm ~/.agent-todo/config.yaml
agent-todo auth login your@email.com password
```

3. Use command-line flags to override:
```bash
agent-todo -s https://correct-server.com project list
```

### Shell Completion Not Working

**Problem:** Tab completion is not working.

**Solutions:**

**Bash:**
```bash
# Ensure completion is loaded
echo 'source <(agent-todo completion bash)' >> ~/.bashrc
source ~/.bashrc
```

**Zsh:**
```bash
# Ensure completion is loaded
echo 'source <(agent-todo completion zsh)' >> ~/.zshrc
source ~/.zshrc
```

**Fish:**
```bash
# Ensure completion file exists
agent-todo completion fish > ~/.config/fish/completions/agent-todo.fish
```

### Agent API Key Lost

**Problem:** You lost an agent's API key.

**Solution:**
Unfortunately, API keys are only shown once when creating an agent. You'll need to:
1. Delete the agent: `agent-todo agent delete <agent-id>`
2. Create a new agent: `agent-todo agent create "<name>" --type <type> --model <model>`
3. Save the new API key securely

### Verbose Debugging

Enable verbose output to see what's happening:

```bash
agent-todo -v task list
```

This will show detailed HTTP requests and responses.

### Getting Help

- Use `--help` flag: `agent-todo task --help`
- Read API documentation: Visit your API's `/docs` endpoint (Swagger UI)
- Check logs: Application logs may have more details
- Report issues: Open an issue on GitHub

---

## CLI vs API Feature Coverage

The CLI currently supports all major API endpoints:

| Feature | API | CLI | Notes |
|---------|-----|-----|-------|
| User Authentication | ✅ | ✅ | Login, register, logout, whoami |
| Project CRUD | ✅ | ✅ | Full CRUD operations |
| Task CRUD | ✅ | ✅ | Full CRUD operations |
| Agent CRUD | ✅ | ✅ | Full CRUD operations |
| Task Assignment | ✅ | ✅ | Assign/unassign agents |
| Task Comments | ✅ | ✅ | Add and view comments |
| Agent Task API | ✅ | ❌ | Agent-specific endpoints (via API key) |
| OpenClaw Tools | ✅ | ❌ | Tool integration endpoints |
| Search/Filter | ✅ | ✅ | Basic filtering |
| Bulk Operations | ❌ | ❌ | Not yet implemented |
| Subtasks | ❌ | ❌ | Not yet implemented |
| Due Dates | ❌ | ❌ | Not yet implemented |
| Tags/Labels | ❌ | ❌ | Not yet implemented |
| Dependencies | ❌ | ❌ | Not yet implemented |
| Notifications | ❌ | ❌ | Not yet implemented |

### Planned CLI Features

Future versions of the CLI will include:

1. **Interactive Mode** - Menu-driven interface
2. **Batch Operations** - Process multiple commands from a file
3. **Output Formatting** - JSON, table, custom formats
4. **Task Templates** - Create tasks from templates
5. **Time Tracking** - Track time spent on tasks
6. **Interactive Filtering** - Select tasks interactively
7. **Task Dependencies** - Manage task dependencies
8. **Webhooks** - Configure webhook notifications
9. **Import/Export** - Import/export tasks and projects
10. **Configuration Profiles** - Multiple server configurations

---

## Best Practices

### Security

1. **Never share API keys** - Agent API keys are sensitive
2. **Use environment variables** - For sensitive configuration
3. **Enable SSL verification** - Only use `--insecure` for local development
4. **Rotate credentials** - Regularly update passwords and API keys
5. **Use read-only tokens** - When possible, use tokens with limited permissions

### Productivity

1. **Use shell completion** - Speed up command entry
2. **Create aliases** - For frequently used commands
3. **Use projects** - Organize tasks logically
4. **Set priorities** - Use priority levels effectively
5. **Add descriptions** - Provide context for tasks
6. **Use comments** - Document task progress
7. **Regular updates** - Keep task status current

### Workflow

1. **Plan projects first** - Create projects before tasks
2. **Break down tasks** - Use descriptive titles and descriptions
3. **Assign appropriately** - Match agent capabilities to tasks
4. **Track progress** - Update status regularly
5. **Communicate** - Use comments for collaboration
6. **Review regularly** - Check task lists and project status

### Example Aliases

Add these to your `~/.bashrc` or `~/.zshrc`:

```bash
# Quick aliases
alias at='agent-todo'
alias atl='agent-todo task list'
alias atp='agent-todo project list'
alias ata='agent-todo agent list'
alias atc='agent-todo task create'
alias att='agent-todo task update'

# Custom functions
at-todo() {
  agent-todo task create "$*" --priority medium
}

at-urgent() {
  agent-todo task create "$*" --priority high
}

at-done() {
  agent-todo task update "$1" --status completed
}

at-todo "Fix the login bug"        # Create medium priority task
at-urgent "Server is down"         # Create high priority task
at-done <task-id>                  # Mark task as done
```

---

## Contributing

We welcome contributions! Please feel free to submit issues and pull requests.

### Development Setup

```bash
# Fork the repository
git clone https://github.com/your-username/agent-todo.git
cd agent-todo/cli

# Install dependencies
go mod download

# Make changes
# ...

# Test
make test

# Build
make build
```

### Submitting Changes

1. Create a feature branch
2. Make your changes
3. Add tests
4. Ensure all tests pass
5. Submit a pull request

---

## License

This CLI is part of the Agent Todo Management Platform. See the main project LICENSE file for details.

---

## Support

- **Documentation**: [Full API Docs](http://localhost:8080/docs)
- **Issues**: [GitHub Issues](https://github.com/formatho/agent-todo/issues)
- **Discussions**: [GitHub Discussions](https://github.com/formatho/agent-todo/discussions)

---

**Version**: 1.0.0
**Last Updated**: 2024-01-20
