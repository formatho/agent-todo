# OpenClaw Plugin Installation Guide

Quick guide to install and configure the Formatho Agent Todo OpenClaw plugin.

## Prerequisites

Before installing the plugin, ensure you have:

- ✅ A running Formatho Agent Todo server
- ✅ Your server URL (e.g., `https://todo.formatho.com` or `http://localhost:8080`)
- ✅ Access to configure OpenClaw plugins

**Verify your server is running:**
```bash
# Replace with your server URL
curl https://todo.formatho.com/health

# Or for local development
curl http://localhost:8080/health
```

## Quick Start (Interactive)

### 1. Copy Plugin to OpenClaw

```bash
# From the agent-todo repository
cp -r formatho-agent-todo-plugin ~/.openclaw/plugins/formatho-agent-todo
```

### 2. Enable Plugin in OpenClaw

The plugin will prompt you for configuration:

1. **Server URL**: Enter your Formatho Agent Todo server URL
   - Example: `https://todo.formatho.com`
   - Example: `http://localhost:8080`
   - The plugin will verify the connection

2. **API Key** (Optional):
   - Press Enter to skip if you only need human access
   - Enter your API key if setting up AI agents

### 3. Restart OpenClaw

Restart the OpenClaw gateway to load the plugin.

### 4. Verify Installation

In OpenClaw, test:
```
Use agent-todo to list all projects
```

## Manual Configuration

If you prefer manual configuration, edit `~/.openclaw/openclaw.json`:

```json
{
  "plugins": {
    "entries": {
      "formatho-agent-todo": {
        "enabled": true,
        "config": {
          "serverUrl": "https://todo.formatho.com",
          "apiKey": "sk-agent-xxxxx",
          "verifyConnection": true
        }
      }
    }
  }
}
```

**Replace with your actual server URL.**

## Configuration Examples

### Local Development Server

```json
{
  "plugins": {
    "entries": {
      "formatho-agent-todo": {
        "enabled": true,
        "config": {
          "serverUrl": "http://localhost:8080"
        }
      }
    }
  }
}
```

### Remote Production Server

```json
{
  "plugins": {
    "entries": {
      "formatho-agent-todo": {
        "enabled": true,
        "config": {
          "serverUrl": "https://todo.formatho.com"
        }
      }
    }
  }
}
```

### With API Key (For Agents)

```json
{
  "plugins": {
    "entries": {
      "formatho-agent-todo": {
        "enabled": true,
        "config": {
          "serverUrl": "https://todo.formatho.com",
          "apiKey": "sk-agent-pm-all-access-xxxxx"
        }
      }
    }
  }
}
```

## Creating a PM Agent

If you need to create a Project Manager agent:

```bash
# Create PM agent
agent-todo --server https://todo.formatho.com agent create "Project Manager" \
  --description "Primary PM agent for managing all agents and tasks" \
  --role admin

# Save the API key from the output
# Then update your config with that key
```

## Troubleshooting

### Plugin Not Loading

1. Check plugin directory exists:
   ```bash
   ls -la ~/.openclaw/plugins/formatho-agent-todo/
   ```

2. Verify `openclaw.plugin.json` is valid

3. Check OpenClaw logs: `~/.openclaw/logs/gateway.log`

### Connection Failed

1. **Verify your server is running:**
   ```bash
   curl https://todo.formatho.com/health
   ```

2. **Check serverUrl in config:**
   ```json
   {
     "config": {
       "serverUrl": "https://todo.formatho.com"
     }
   }
   ```

3. **Ensure network permissions allow your server domain**

### Wrong Server URL

If you entered the wrong URL, edit `~/.openclaw/openclaw.json`:
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

### CLI Not Found

The plugin will auto-install if enabled. Manual installation:
```bash
go install github.com/formatho/agent-todo/cli@latest
```

## Full Documentation

For complete documentation, see:
- **Plugin README**: `formatho-agent-todo-plugin/README.md`
- **Skill Documentation**: `skills/formatho-agent-todo/SKILL.md`
- **Project README**: `README.md`
- **API Documentation**: Available at your server's `/docs` endpoint

## Support

- **Issues**: https://github.com/formatho/agent-todo/issues
- **Documentation**: https://github.com/formatho/agent-todo

