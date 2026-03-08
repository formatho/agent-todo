# Formatho Agent Todo OpenClaw Plugin

This plugin integrates the Formatho Agent Todo Management Platform with OpenClaw, enabling AI agents to manage tasks, projects, and multi-agent workflows.

## Features

- 📋 **Task Management**: Create, update, and track tasks with priorities and status
- 📁 **Project Organization**: Group tasks into projects for better organization
- 🤖 **Multi-Agent Support**: Hierarchical agent system with PM, Supervisor, and Regular roles
- 🔐 **Role-Based Permissions**: PM agents have all-access, regular agents have self-write only
- 🚀 **CLI Integration**: Full agent-todo CLI functionality
- 🔑 **API Key Management**: Secure API key provisioning and storage

## Installation

### Method 1: Manual Installation

1. **Clone or copy the plugin to OpenClaw's plugins directory:**

```bash
# Copy the entire plugin directory
cp -r formatho-agent-todo-plugin ~/.openclaw/plugins/formatho-agent-todo
```

2. **Restart OpenClaw gateway** to load the plugin

### Method 2: Symbolic Link (Recommended for Development)

```bash
# Create symbolic link for easier updates
ln -s /path/to/agent-todo/formatho-agent-todo-plugin ~/.openclaw/plugins/formatho-agent-todo
```

## Configuration

### Interactive Setup (Recommended)

When you enable the plugin for the first time, OpenClaw will prompt you for:

1. **Server URL**: Enter your Formatho Agent Todo server URL
   - Example: `https://todo.example.com` or `http://localhost:8080`
   - The plugin will verify the connection before saving

2. **API Key** (Optional): Enter your API key if you're setting up agents
   - Press Enter to skip if you only need human access
   - Required for AI agents to authenticate

### Manual Configuration

Alternatively, add the following to your `~/.openclaw/openclaw.json`:

```json
{
  "plugins": {
    "entries": {
      "formatho-agent-todo": {
        "enabled": true,
        "config": {
          "serverUrl": "https://your-todo-server.com",
          "apiKey": "your-api-key-here",
          "verifyConnection": true
        }
      }
    }
  }
}
```

### Configuration Options

| Option | Type | Required | Default | Description |
|--------|------|----------|---------|-------------|
| `enabled` | boolean | Yes | `false` | Enable/disable the plugin |
| `serverUrl` | string | **Yes** | `""` | Your Formatho Agent Todo server URL (provided during setup) |
| `apiKey` | string | No | `""` | API key for agent authentication |
| `autoInstall` | boolean | No | `true` | Auto-install CLI if missing |
| `verifyConnection` | boolean | No | `true` | Verify server connection during setup |

### Prerequisites

**Before configuring the plugin, ensure:**

1. ✅ Your Formatho Agent Todo server is running and accessible
2. ✅ You know your server URL (e.g., `https://todo.example.com`)
3. ✅ If using agents, you have your API key ready

### Server URL Examples

#### Local Development Server
```json
{
  "config": {
    "serverUrl": "http://localhost:8080"
  }
}
```

#### Remote Production Server

```json
{
  "config": {
    "serverUrl": "https://todo.example.com"
  }
}
```

#### Production with Custom Port

```json
{
  "config": {
    "serverUrl": "https://todo.example.com:9443"
  }
}
```

## Authentication Setup

### For Human Users (Optional)

If you want to create agents manually, you can use JWT authentication:

```bash
agent-todo auth login your-email@example.com password
```

### For AI Agents (Required)

AI agents use API keys for authentication. There are two ways to configure:

#### Option 1: Plugin Configuration (Recommended)

Set the API key in plugin config:

```json
{
  "plugins": {
    "entries": {
      "formatho-agent-todo": {
        "config": {
          "serverUrl": "http://localhost:8080",
          "apiKey": "sk-agent-xxxxx"
        }
      }
    }
  }
}
```

#### Option 2: Environment Variable

Set the environment variable:

```bash
export AGENT_TODO_API_KEY="sk-agent-xxxxx"
export AGENT_TODO_SERVER_URL="http://localhost:8080"
```

## Agent Setup Workflow

### 1. Verify Your Server is Running

Before configuring the plugin, ensure your Formatho Agent Todo server is accessible:

```bash
# Replace with your server URL
curl https://your-todo-server.com/health

# Or for local development
curl http://localhost:8080/health
```

You should see a health check response confirming the server is running.

### 2. Configure OpenClaw Plugin

**Option A: Interactive Setup (Recommended)**

Enable the plugin in OpenClaw and follow the prompts:
1. Enter your server URL when asked
2. Enter your API key (optional, for agents)

**Option B: Manual Configuration**

Add to `~/.openclaw/openclaw.json`:

```json
{
  "plugins": {
    "entries": {
      "formatho-agent-todo": {
        "enabled": true,
        "config": {
          "serverUrl": "https://your-todo-server.com"
        }
      }
    }
  }
}
```

### 3. Create a Project Manager (PM) Agent

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
```

### 4. Update Plugin Config with PM Key

```json
{
  "plugins": {
    "entries": {
      "formatho-agent-todo": {
        "enabled": true,
        "config": {
          "serverUrl": "http://localhost:8080",
          "apiKey": "sk-agent-pm-all-access-xxxxx-save-this-key"
        }
      }
    }
  }
}
```

### 5. Restart OpenClaw

Restart the gateway to apply the configuration changes.

## Usage

Once configured, agents can use the skill to manage tasks and projects:

### Create Tasks

```
Use the agent-todo skill to create a high-priority task for "Review PR #123"
```

### List Tasks

```
Show me all pending high-priority tasks assigned to me
```

### Update Status

```
Mark task "Review PR #123" as in-progress with comment "Started reviewing"
```

### Agent Provisioning

```
Create a new regular agent called "Data Processor" for handling CSV files
```

## Permissions

The plugin requires:

- **Network Access**: To communicate with your Formatho Agent Todo server
- **Execute Permission**: To run the agent-todo CLI binary

The plugin supports connections to any server URL (configured during setup). OpenClaw will request network permissions for the domain you specify.

## Troubleshooting

### Plugin Not Loading

1. Check the plugin is in the correct directory: `~/.openclaw/plugins/formatho-agent-todo/`
2. Verify `openclaw.plugin.json` is valid JSON
3. Check OpenClaw logs: `~/.openclaw/logs/gateway.log`

### Server Connection Failed

1. **Verify your server is running:**
   ```bash
   # Replace with your server URL
   curl https://your-todo-server.com/health
   ```

2. **Check the serverUrl in your config matches your server:**
   ```json
   {
     "config": {
       "serverUrl": "https://your-todo-server.com"
     }
   }
   ```

3. **If connection verification fails:**
   - Ensure your server is accessible from your network
   - Check firewall settings
   - Verify the URL includes the correct protocol (http:// or https://)
   - For self-signed certificates, you may need to disable `verifyConnection`

### Configuration Prompt Not Appearing

1. Enable the plugin in OpenClaw
2. Restart the OpenClaw gateway
3. Try using the plugin - OpenClaw will prompt for configuration if needed

### Wrong Server URL

If you entered the wrong URL during setup:

1. Edit `~/.openclaw/openclaw.json`
2. Update the `serverUrl` value:
   ```json
   {
     "plugins": {
       "entries": {
         "formatho-agent-todo": {
           "config": {
             "serverUrl": "https://correct-url.com"
           }
         }
       }
     }
   }
   ```
3. Restart OpenClaw gateway

### CLI Not Found

The plugin will auto-install the CLI if `autoInstall` is true (default).

Manual installation:
```bash
go install github.com/formatho/agent-todo/cli@latest
```

### Authentication Errors

1. Verify your API key is correct
2. Check the key has the required permissions for your agent role
3. Ensure the key is set in plugin config or environment variable

## Development

### Project Structure

```
openclaw-plugin/
├── openclaw.plugin.json    # Plugin metadata and config schema
├── README.md               # This file
└── ../skills/agent-todo/   # Skill definition
    └── SKILL.md            # Skill instructions
```

### Local Testing

1. Start the Formatho Agent Todo server:
   ```bash
   cd /path/to/agent-todo
   docker-compose up -d
   ```

2. Configure plugin in `~/.openclaw/openclaw.json`

3. Restart OpenClaw gateway

4. Test in OpenClaw:
   ```
   Use agent-todo to list all projects
   ```

## Support

- **Documentation**: https://github.com/formatho/agent-todo
- **Issues**: https://github.com/formatho/agent-todo/issues
- **API Docs**: http://localhost:8080/docs (when server is running)

## License

MIT License - See LICENSE file for details
