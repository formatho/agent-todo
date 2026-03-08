# Agent Todo Management Platform

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

## License

MIT

## Support

For issues and questions, please open a GitHub issue.
