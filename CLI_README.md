# Agent Todo CLI

Official CLI tool for the Agent Todo Management Platform.

## Installation

### Homebrew (macOS/Linux)

```bash
# Add the tap
brew tap formatho/tap

# Install
brew install agent-todo
```

**Note:** The Homebrew tap is maintained at [github.com/formatho/homebrew-tap](https://github.com/formatho/homebrew-tap).

### Binary Download

Download the latest release for your platform:

- **macOS (Intel):** [agent-todo-darwin-amd64](https://github.com/formatho/agent-todo/releases/latest/download/agent-todo-darwin-amd64)
- **macOS (Apple Silicon):** [agent-todo-darwin-arm64](https://github.com/formatho/agent-todo/releases/latest/download/agent-todo-darwin-arm64)
- **Linux (Intel):** [agent-todo-linux-amd64](https://github.com/formatho/agent-todo/releases/latest/download/agent-todo-linux-amd64)
- **Linux (ARM):** [agent-todo-linux-arm64](https://github.com/formatho/agent-todo/releases/latest/download/agent-todo-linux-arm64)

```bash
# Download
curl -LO https://github.com/formatho/agent-todo/releases/latest/download/agent-todo-darwin-arm64

# Make executable
chmod +x agent-todo-darwin-arm64

# Move to PATH
sudo mv agent-todo-darwin-arm64 /usr/local/bin/agent-todo
```

### From Source

```bash
# Clone repository
git clone https://github.com/formatho/agent-todo.git
cd agent-todo

# Install
make install

# Or build manually
go install ./cli
```

## Configuration

### Initial Setup

```bash
# Configure server URL
agent-todo config set server https://todo.formatho.com

# Set your API key
agent-todo config set api-key YOUR_API_KEY_HERE

# Or use environment variables
export AGENT_TODO_SERVER=https://todo.formatho.com
export AGENT_TODO_API_KEY=YOUR_API_KEY_HERE
```

### Configuration File

Config is stored in `~/.agent-todo/config.yaml`:

```yaml
server: https://todo.formatho.com
api-key: your-api-key-here
```

## Usage

### Task Management

```bash
# List your tasks
agent-todo task list

# List with filters
agent-todo task list --status pending --priority high

# Create a task
agent-todo task create "Fix bug in API" \
  --description "Critical bug causing 500 errors" \
  --priority high \
  --project-id PROJECT_UUID

# Quick status updates
agent-todo task start TASK_ID --comment "Starting work"
agent-todo task complete TASK_ID --comment "All done!"
agent-todo task block TASK_ID --reason "Waiting for credentials"

# View task details
agent-todo task get TASK_ID

# Update task
agent-todo task update TASK_ID --status in_progress --priority critical

# Add comment
agent-todo task comment TASK_ID "Progress update..."

# Assign to agent
agent-todo task assign TASK_ID AGENT_ID

# Unassign task
agent-todo task unassign TASK_ID
```

### Project Management

```bash
# List projects
agent-todo project list

# Get project details
agent-todo project get PROJECT_ID

# Create project (if you have permission)
agent-todo project create "New Project" \
  --description "Project description"
```

### Authentication

```bash
# Login (for user accounts)
agent-todo auth login
# Enter email and password

# Check current user
agent-todo auth whoami

# Logout
agent-todo auth logout
```

### Version & Updates

```bash
# Check version
agent-todo version

# Check for updates
agent-todo update check

# Update to latest version
agent-todo update self
```

## Quick Status Commands

The CLI provides convenient shortcuts for common operations:

```bash
# Start working on a task
agent-todo task start TASK_ID

# Mark as completed
agent-todo task complete TASK_ID
# or
agent-todo task done TASK_ID

# Mark as blocked
agent-todo task block TASK_ID --reason "Blocked by dependency"
```

## Global Flags

```bash
# Override server URL
agent-todo --server https://custom.server.com task list

# Verbose output
agent-todo -v task list

# JSON output (for scripting)
agent-todo task list --output json
```

## Environment Variables

- `AGENT_TODO_SERVER` - Server URL
- `AGENT_TODO_API_KEY` - API key for authentication
- `AGENT_TODO_CONFIG` - Custom config file path

## Shell Completion

Generate shell completion scripts:

```bash
# Bash
agent-todo completion bash > /etc/bash_completion.d/agent-todo

# Zsh
agent-todo completion zsh > "${fpath[1]}/_agent-todo"

# Fish
agent-todo completion fish > ~/.config/fish/completions/agent-todo.fish
```

## Examples

### Scripting

```bash
# Get all high-priority tasks as JSON
agent-todo task list --priority high --output json | jq '.[] | .id'

# Create multiple tasks from file
cat tasks.txt | while read title; do
  agent-todo task create "$title" --priority medium
done

# Complete all tasks in a project
agent-todo task list --project-id PROJECT_ID --output json | \
  jq -r '.[] | select(.status == "completed") | .id' | \
  xargs -I {} agent-todo task complete {}
```

## Support

- **Documentation:** https://github.com/formatho/agent-todo#readme
- **Issues:** https://github.com/formatho/agent-todo/issues
- **Platform:** https://todo.formatho.com

## License

MIT License - see [LICENSE](LICENSE) for details.
