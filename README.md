# Agent Todo Management Platform

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Made with Vue](https://img.shields.io/badge/Vue-3.0+-4FC08D?style=flat&logo=vue.js)](https://vuejs.org)

**Task management for AI agents. Because they keep forgetting everything.**

🌐 **Live Demo:** [todo.formatho.com](https://todo.formatho.com)
📚 **Landing Page:** [formatho.com/tools/agent-todo](https://formatho.com/tools/agent-todo)
📖 **Documentation:** [GitHub Wiki](https://github.com/formatho/agent-todo/wiki)

---

## Why Agent-Todo?

Traditional task managers (Todoist, Asana, Trello) are built for **humans** with UIs and per-user pricing. They don't work for AI agents.

**Agent-Todo is different:**

| Feature | Traditional Tools | Agent-Todo |
|---------|------------------|------------|
| **Primary Interface** | UI/Clicks | API-First ✅ |
| **User Model** | 1 human/account | 100s of agents/key ✅ |
| **Context** | None | Agent Ownership ✅ |
| **Integration** | Complex OAuth | Simple API Key ✅ |
| **Pricing** | Per-user | One API, many agents ✅ |
| **Memory** | Session-based | Persistent Storage ✅ |

**Perfect for:**
- 🤖 AI agents and autonomous systems
- 🧠 LLM-powered workflows
- ⚙️ Multi-agent orchestration
- 🔄 Background automation tasks
- 📊 Agent performance tracking

**Quick Example:**
```bash
# Agent creates a task
curl -X POST "https://todo.formatho.com/api/agent/tasks" \
  -H "X-API-Key: YOUR_KEY" \
  -d '{"title":"Process data","priority":"high"}'

# Agent marks it complete
curl -X PATCH "https://todo.formatho.com/api/agent/tasks/{id}/status" \
  -H "X-API-Key: YOUR_KEY" \
  -d '{"status":"completed"}'
```

That's it. Your agents now have persistent memory.

**Free Tier:** 3 agents, 100 tasks/day, no credit card required.

---

## 🎬 Demo

**[📺 Watch 3-minute demo](https://youtube.com/placeholder)** (Coming soon)

Or try it live: [todo.formatho.com](https://todo.formatho.com)

**What you'll see:**
- ✅ Agent creating tasks via API (2 lines of code)
- ✅ Human reviewing and managing agent tasks
- ✅ Real-time status updates
- ✅ Persistent memory across sessions
- ✅ Team collaboration in action

**Want the demo video sooner?** ⭐ Star this repo and we'll prioritize it!

---

## About

A production-grade task management system designed for collaboration between humans and AI agents. Built with Go (Gin), PostgreSQL, Vue.js, and Docker.

## Features

- **Human Task Management**: Full CRUD operations for tasks via web UI
- **Agent Integration**: AI agents can manage tasks via API using API keys
- **Advanced Filtering**: Filter tasks by status, priority, agent, search terms
- **Comments & History**: Task comments and full audit log
- **Priority Levels**: Low, Medium, High, Critical
- **Due Dates**: Optional deadline tracking
- **OpenClaw Compatible**: Tool endpoints for agent frameworks
- **Swagger Documentation**: Auto-generated API docs at `/docs`

## Tech Stack

### Backend
- Go 1.21
- Gin Web Framework
- PostgreSQL 16
- GORM
- JWT Authentication
- golang-migrate

### Frontend
- Vue 3 (Composition API)
- Vite
- TailwindCSS
- Vue Router
- Pinia
- Axios

## Quick Start

### Prerequisites
- Docker and Docker Compose
- Git

### Running with Docker

1. Clone the repository:
```bash
git clone <repository-url>
cd agent-todo
```

2. Start the system:
```bash
docker compose up
```

3. Access the application:
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- Swagger Docs: http://localhost:8080/docs/index.html
- Database: postgresql://agent_todo:agent_todo_pass@localhost:5432/agent_todo

4. Default credentials:
- Email: `admin@example.com`
- Password: `admin123`

## Development Setup

### Backend Development

1. Install Go 1.21+
2. Install PostgreSQL 16
3. Install golang-migrate:
```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

4. Setup environment:
```bash
cd backend
cp .env.example .env
# Edit .env with your settings
```

5. Run migrations:
```bash
migrate -path db/migrations -database "postgres://agent_todo:agent_todo_pass@localhost:5432/agent_todo?sslmode=disable" up
```

6. Install dependencies:
```bash
go mod download
```

7. Run server:
```bash
go run cmd/api/main.go
```

8. Generate Swagger docs:
```bash
swag init -g cmd/api/main.go
```

### Frontend Development

1. Install Node.js 20+
2. Install dependencies:
```bash
cd frontend
npm install
```

3. Run dev server:
```bash
npm run dev
```

4. Build for production:
```bash
npm run build
```

## API Usage

### Human Authentication

Register:
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

Login:
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"admin123"}'
```

Use the returned token in the `Authorization: Bearer TOKEN` header.

### Creating Tasks

```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "title": "Complete project documentation",
    "description": "Write comprehensive docs",
    "priority": "high",
    "due_date": "2024-12-31T23:59:59Z",
    "assigned_agent_id": "AGENT_ID"
  }'
```

### Agent Authentication

Agents authenticate using API keys in the `X-API-KEY` header.

### Agent Task Creation

```bash
curl -X POST http://localhost:8080/agent/tasks \
  -H "Content-Type: application/json" \
  -H "X-API-KEY: sk_agent_example_key_12345" \
  -d '{
    "title": "Process data files",
    "description": "Analyze CSV files",
    "priority": "medium"
  }'
```

### Agent Status Update

```bash
curl -X PATCH http://localhost:8080/agent/tasks/TASK_ID/status \
  -H "Content-Type: application/json" \
  -H "X-API-KEY: sk_agent_example_key_12345" \
  -d '{"status": "in_progress"}'
```

## CLI Usage

The `agent-todo` CLI provides convenient commands for managing tasks:

### Quick Status Commands

```bash
# Mark task as in progress
agent-todo task start <task-id> --comment "Starting work"

# Mark task as completed
agent-todo task complete <task-id> --comment "All done!"
# Alias: task done <task-id>

# Mark task as blocked
agent-todo task block <task-id> --reason "Waiting for API credentials"
```

### Task Management

```bash
# List tasks with filters
agent-todo task list --status pending --priority high

# Create a new task
agent-todo task create "Task title" --description "Details" --priority high

# Update task
agent-todo task update <task-id> --status in_progress

# Add comment
agent-todo task comment <task-id> "Progress update"

# Assign to agent
agent-todo task assign <task-id> <agent-id>
```

### Authentication

The CLI reads credentials from `~/.agent-todo/config.yaml` or uses the `--server` and `--api-key` flags:

```bash
agent-todo --server https://todo.formatho.com --api-key YOUR_KEY task list
```

## OpenClaw Tool Endpoints

The platform provides OpenClaw-compatible endpoints:

```bash
# Create task
curl -X POST http://localhost:8080/tools/tasks/create \
  -H "Content-Type: application/json" \
  -H "X-API-KEY: YOUR_AGENT_KEY" \
  -d '{"title": "Task title", "description": "Description", "priority": "high"}'

# Update task
curl -X POST http://localhost:8080/tools/tasks/update \
  -H "Content-Type: application/json" \
  -H "X-API-KEY: YOUR_AGENT_KEY" \
  -d '{"task_id": "TASK_ID", "status": "completed", "comment": "Done!"}'

# List tasks
curl -X POST http://localhost:8080/tools/tasks/list \
  -H "Content-Type: application/json" \
  -H "X-API-KEY: YOUR_AGENT_KEY" \
  -d '{"status": "pending", "limit": 10}'
```

## Creating Users

### Via Web UI
1. Navigate to http://localhost:3000
2. Click "Create a new account"
3. Enter email and password (min 6 characters)

### Via API
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"newuser@example.com","password":"password123"}'
```

## Creating Agents

### Via Web UI
1. Login to the dashboard
2. Navigate to "Agents" page
3. Click "Create Agent"
4. Enter name and description
5. Copy the generated API key

### Via API
```bash
curl -X POST http://localhost:8080/agents \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"name":"My Agent","description":"An AI assistant"}'
```

The response will contain the API key for the agent.

## Testing Agent Workflows

### Test Agent Creating a Task

```bash
export AGENT_KEY="sk_agent_example_key_12345"

# Agent creates a task
curl -X POST http://localhost:8080/agent/tasks \
  -H "Content-Type: application/json" \
  -H "X-API-KEY: $AGENT_KEY" \
  -d '{
    "title": "Test task from agent",
    "description": "This task was created by an agent",
    "priority": "medium"
  }'
```

### Test Human Editing Agent Task

```bash
export TOKEN="YOUR_JWT_TOKEN"

# View all tasks (including agent-created ones)
curl -X GET http://localhost:8080/tasks \
  -H "Authorization: Bearer $TOKEN"

# Update the agent's task
curl -X PATCH http://localhost:8080/tasks/TASK_ID \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"title": "Updated by human"}'
```

### Test Agent Updating Task Status

```bash
# Agent updates task status
curl -X PATCH http://localhost:8080/agent/tasks/TASK_ID/status \
  -H "Content-Type: application/json" \
  -H "X-API-KEY: $AGENT_KEY" \
  -d '{"status": "completed"}'
```

## Task Status Flow

Tasks can have the following statuses:
- `pending`: Task is not yet started
- `in_progress`: Task is currently being worked on
- `completed`: Task has been completed successfully
- `failed`: Task failed

## Task Priorities

- `low`: Low priority task
- `medium`: Medium priority (default)
- `high`: High priority task
- `critical`: Critical priority task

## Database Schema

The system uses the following tables:
- `users`: Human users
- `agents`: AI agents with API keys
- `tasks`: Tasks with assignments
- `task_events`: Audit log for task changes
- `task_comments`: Comments on tasks

## Architecture

```
┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│   Vue.js    │────▶│  Go Backend  │────▶│ PostgreSQL  │
│  Frontend   │     │   (Gin)      │     │  Database   │
└─────────────┘     └──────────────┘     └─────────────┘
                           │
                           ▼
                    ┌──────────────┐
                    │   Agents     │
                    │ (API Keys)   │
                    └──────────────┘
```

## Environment Variables

### Backend
- `DATABASE_URL`: PostgreSQL connection string
- `PORT`: Server port (default: 8080)
- `GIN_MODE`: gin mode (debug/release)
- `JWT_SECRET`: Secret for JWT tokens
- `FRONTEND_URL`: Frontend URL for CORS

### Frontend
- `VITE_API_URL`: Backend API URL (default: http://localhost:8080)

## Production Deployment

1. Change `JWT_SECRET` to a strong random value
2. Set `GIN_MODE=release`
3. Use strong database passwords
4. Enable HTTPS/TLS
5. Set up proper CORS origins
6. Use environment-specific configurations
7. Set up database backups
8. Monitor logs and metrics

## 🗺️ Roadmap

**Current Status:** Production-ready v1.0

**Q2 2026 (In Progress):**
- ✅ Core task management API
- ✅ Agent authentication & API keys
- ✅ Human-agent collaboration
- ✅ Team workspaces
- ✅ State persistence for agents
- 🔄 Advanced analytics dashboard
- 🔄 Webhook notifications
- 🔄 Integration marketplace

**Q3 2026 (Planned):**
- 📋 Multi-tenant enterprise features
- 📋 Advanced permissions (RBAC)
- 📋 Custom workflows & automation
- 📋 Audit trail exports (SOC 2, GDPR)
- 📋 Performance benchmarking tools

**Q4 2026 (Vision):**
- 🔮 Agent marketplace
- 🔮 Cross-platform SDKs (Python, Node, Go)
- 🔮 Enterprise SSO integration
- 🔮 On-premise deployment option

**Recent Wins:**
- 🚀 Live demo launched: [todo.formatho.com](https://todo.formatho.com)
- 🤝 Partnership discussions with leading agent frameworks
- ⭐ Growing community (join us!)

**Want to influence the roadmap?** [Open a discussion](https://github.com/formatho/agent-todo/discussions) or [vote on features](https://github.com/formatho/agent-todo/issues)!

---

## 🤝 Contributing

We welcome contributions! Agent-Todo is built by the community, for the community.

**Ways to contribute:**
- 🐛 **Report bugs** - [Open an issue](https://github.com/formatho/agent-todo/issues)
- 💡 **Suggest features** - [Start a discussion](https://github.com/formatho/agent-todo/discussions)
- 📖 **Improve docs** - Fix typos, add examples, clarify explanations
- 🔧 **Submit PRs** - Check out [good first issues](https://github.com/formatho/agent-todo/labels/good%20first%20issue)
- 🌟 **Spread the word** - Star the repo, share on social media

**Development setup:**
1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`go test ./...` for backend, `npm test` for frontend)
5. Commit changes (`git commit -m 'Add amazing feature'`)
6. Push to branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

**Contributors get:**
- 📛 Credit in our README
- 🎁 Swag (for significant contributions)
- 🚀 Early access to new features
- 💼 Reference letter (upon request)

See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines.

---

## 📜 License

MIT License - Use it however you want!

---

## 💬 Support

**Need help?** We're here for you:

- 📖 **Documentation:** [GitHub Wiki](https://github.com/formatho/agent-todo/wiki)
- 🐛 **Bug reports:** [Open an issue](https://github.com/formatho/agent-todo/issues)
- 💡 **Feature requests:** [Start a discussion](https://github.com/formatho/agent-todo/discussions)
- 📧 **Direct contact:** team@formatho.com
- 💬 **Community:** Join our Discord (coming soon)

**Enterprise support** available - contact us for SLAs, dedicated support, and custom features.
