# Quick Start Guide

## 🚀 Start the System

The entire system can be started with a single command:

```bash
docker compose up
```

This will start:
- **PostgreSQL** on port 5432
- **Backend API** on port 8080
- **Frontend UI** on port 3000

## 📍 Access Points

Once running:
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- Swagger Docs: http://localhost:8080/docs/index.html
- Database: postgresql://agent_todo:agent_todo_pass@localhost:5432/agent_todo

## 🔑 Default Credentials

**Admin User:**
- Email: `admin@example.com`
- Password: `admin123`

**Example Agent:**
- Name: `Example Agent`
- API Key: `sk_agent_example_key_12345`

## 📝 Quick Workflow

### 1. Login as User

1. Navigate to http://localhost:3000
2. Login with `admin@example.com` / `admin123`

### 2. View Tasks

You'll see the default tasks created by the seed data.

### 3. Create a Task

1. Click "Create Task" button
2. Fill in title, description, priority
3. Optionally assign to an agent

### 4. Manage Agents

1. Go to "Agents" page
2. View existing agents
3. Create new agents (API keys are auto-generated)

### 5. Test Agent Integration

Using the example agent API key:

```bash
export AGENT_KEY="sk_agent_example_key_12345"

# Agent creates a task
curl -X POST http://localhost:8080/agent/tasks \
  -H "Content-Type: application/json" \
  -H "X-API-KEY: $AGENT_KEY" \
  -d '{
    "title": "Test task from agent",
    "description": "Created via API",
    "priority": "high"
  }'

# Agent lists their tasks
curl http://localhost:8080/agent/tasks \
  -H "X-API-KEY: $AGENT_KEY"

# Agent updates task status
curl -X PATCH http://localhost:8080/agent/tasks/{TASK_ID}/status \
  -H "Content-Type: application/json" \
  -H "X-API-KEY: $AGENT_KEY" \
  -d '{"status": "completed"}'
```

## 🧪 Test Complete Workflow

### Scenario: Human-Agent Collaboration

1. **Human creates task**
```bash
# Login first to get token
export TOKEN="YOUR_JWT_TOKEN"

curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title": "Analyze customer data",
    "description": "Process CSV files and generate report",
    "priority": "high",
    "assigned_agent_id": "ASSIGNED_AGENT_ID"
  }'
```

2. **Agent sees task and starts work**
```bash
curl http://localhost:8080/agent/tasks \
  -H "X-API-KEY: sk_agent_example_key_12345"

# Update status
curl -X PATCH http://localhost:8080/agent/tasks/{TASK_ID}/status \
  -H "Content-Type: application/json" \
  -H "X-API-KEY: sk_agent_example_key_12345" \
  -d '{"status": "in_progress"}'
```

3. **Human monitors progress in UI**
   - Refresh dashboard
   - See status change to "In Progress"

4. **Agent completes task**
```bash
curl -X PATCH http://localhost:8080/agent/tasks/{TASK_ID}/status \
  -H "Content-Type: application/json" \
  -H "X-API-KEY: sk_agent_example_key_12345" \
  -d '{"status": "completed"}'
```

5. **Human views task details**
   - Click on task in dashboard
   - See full history and status changes

## 🛠 OpenClaw Tool Endpoints

Test OpenClaw-compatible endpoints:

```bash
export AGENT_KEY="sk_agent_example_key_12345"

# Create task via tool endpoint
curl -X POST http://localhost:8080/tools/tasks/create \
  -H "Content-Type: application/json" \
  -H "X-API-KEY: $AGENT_KEY" \
  -d '{
    "title": "Tool endpoint test",
    "description": "Created via OpenClaw endpoint",
    "priority": "medium"
  }'

# Update via tool endpoint
curl -X POST http://localhost:8080/tools/tasks/update \
  -H "Content-Type: application/json" \
  -H "X-API-KEY: $AGENT_KEY" \
  -d '{
    "task_id": "TASK_ID",
    "status": "completed",
    "comment": "Finished successfully"
  }'

# List via tool endpoint
curl -X POST http://localhost:8080/tools/tasks/list \
  -H "Content-Type: application/json" \
  -H "X-API-KEY: $AGENT_KEY" \
  -d '{"limit": 5}'
```

## 📊 View API Documentation

Visit http://localhost:8080/docs/index.html for interactive Swagger documentation.

## 🔍 Troubleshooting

### Services not starting
```bash
# Check logs
docker compose logs

# Restart services
docker compose down
docker compose up
```

### Database connection issues
```bash
# Check if PostgreSQL is running
docker compose ps

# View database logs
docker compose logs postgres
```

### Frontend can't reach backend
- Ensure backend is running on port 8080
- Check CORS settings in backend/.env

### Port conflicts
- Edit ports in docker-compose.yml if needed
- Frontend default: 3000
- Backend default: 8080
- Database default: 5432

## 📚 Next Steps

1. **Create your own agents** in the Agents page
2. **Integrate with your AI agents** using the API
3. **Build custom workflows** using the tool endpoints
4. **Extend the system** with additional features

See README.md for complete API documentation and advanced usage.
