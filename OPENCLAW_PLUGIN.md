# OpenClaw Plugin Installation Guide

Quick guide to install and configure the Formatho Agent Todo OpenClaw plugin.

## Quick Start

### 1. Copy Plugin to OpenClaw

```bash
# From the agent-todo repository
cp -r formatho-agent-todo-plugin ~/.openclaw/plugins/formatho-agent-todo
```

### 2. Configure Server URL

Edit `~/.openclaw/openclaw.json` and add:

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

**IMPORTANT**: Replace `http://localhost:8080` with your Formatho Agent Todo server URL.

### 3. Start Formatho Agent Todo Server

```bash
# If not already running
cd /path/to/agent-todo
docker-compose up -d
```

### 4. Restart OpenClaw

Restart the OpenClaw gateway to load the plugin.

### 5. Verify Installation

In OpenClaw, test:
```
Use agent-todo to list all projects
```

## Configuration Examples

### Local Server (Default)

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

### Remote Server

```json
{
  "plugins": {
    "entries": {
      "formatho-agent-todo": {
        "enabled": true,
        "config": {
          "serverUrl": "https://todo.example.com"
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
          "serverUrl": "http://localhost:8080",
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
agent-todo agent create "Project Manager" \
  --description "Primary PM agent for managing all agents and tasks" \
  --role admin

# Save the API key from the output
# Then update your config with that key
```

## Troubleshooting

### Plugin Not Loading

1. Check plugin directory exists:
   ```bash
   ls -la ~/.openclaw/plugins/agent-todo/
   ```

2. Verify openclaw.plugin.json is valid

3. Check OpenClaw logs: `~/.openclaw/logs/gateway.log`

### Connection Failed

1. Verify server is running:
   ```bash
   curl http://localhost:8080/health
   ```

2. Check serverUrl in config matches your server

### CLI Not Found

The plugin will auto-install if enabled. Manual installation:
```bash
go install github.com/formatho/agent-todo/cli@latest
```

## Full Documentation

For complete documentation, see:
- **Plugin README**: `openclaw-plugin/README.md`
- **Skill Documentation**: `skills/agent-todo/SKILL.md`
- **Project README**: `README.md`
- **API Documentation**: http://localhost:8080/docs (when server is running)

## Support

- **Issues**: https://github.com/formatho/agent-todo/issues
- **Documentation**: https://github.com/formatho/agent-todo
