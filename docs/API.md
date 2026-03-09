# Agent Todo API Documentation

**Version:** 2.0  
**Base URL:** `https://todo.formatho.com/api`  
**Last Updated:** March 8, 2026

---

## Table of Contents

1. [Overview](#overview)
2. [Authentication](#authentication)
3. [Organisation Context](#organisation-context)
4. [API Endpoints](#api-endpoints)
   - [Organisations](#organisations)
   - [Projects](#projects)
   - [Tasks](#tasks)
   - [Agents](#agents)
5. [Migration Guide](#migration-guide)
6. [Examples](#examples)
7. [Error Handling](#error-handling)

---

## Overview

The Agent Todo API provides task management capabilities for both human users and AI agents. Version 2.0 introduces **multi-organisation support** with complete data isolation between organisations.

### Key Features

- **Multi-tenant Architecture:** Complete data isolation between organisations
- **Role-Based Access Control:** Owner, Admin, and Member roles
- **Dual Authentication:** JWT for humans, API keys for agents
- **RESTful Design:** Consistent endpoint patterns and response formats

### Base URL

```
https://todo.formatho.com/api
```

---

## Authentication

### Human Users (JWT)

All human user endpoints require a JWT Bearer token:

```bash
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  https://todo.formatho.com/api/organisations
```

#### Getting a Token

**Register:**
```bash
POST /auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "uuid",
    "email": "user@example.com"
  }
}
```

**Login:**
```bash
POST /auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

### AI Agents (API Key)

Agent endpoints require an API key:

```bash
curl -H "X-API-Key: sk_agent_your_api_key" \
  https://todo.formatho.com/api/agent/tasks
```

---

## Organisation Context

### Overview

Most API endpoints operate within an **organisation context**. This ensures data isolation between different teams or companies.

### Setting Organisation Context

After authentication, include the organisation ID in your requests:

**Option 1: JWT Token (Automatic)**

When you login, the JWT includes your current organisation:

```json
{
  "user_id": "uuid",
  "email": "user@example.com",
  "current_org_id": "org-uuid",
  "exp": 1712505600
}
```

**Option 2: X-Organisation-ID Header**

```bash
curl -H "Authorization: Bearer TOKEN" \
     -H "X-Organisation-ID: org-uuid" \
     https://todo.formatho.com/api/tasks
```

### Switching Organisations

Users can belong to multiple organisations. To switch:

```bash
# List your organisations
GET /organisations

# The JWT token will include the current organisation context
# After selecting an organisation, all subsequent requests
# will operate within that context
```

---

## API Endpoints

### Organisations

#### Create Organisation

Create a new organisation. You become the owner automatically.

```bash
POST /organisations
Authorization: Bearer TOKEN
Content-Type: application/json

{
  "name": "Acme Corporation",
  "slug": "acme-corp",
  "description": "Main company workspace"
}
```

**Response:** `201 Created`
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Acme Corporation",
  "slug": "acme-corp",
  "description": "Main company workspace",
  "status": "active",
  "created_at": "2026-03-08T15:30:00Z",
  "updated_at": "2026-03-08T15:30:00Z",
  "created_by": {
    "id": "user-uuid",
    "email": "user@example.com"
  }
}
```

#### List Organisations

List all organisations you belong to.

```bash
GET /organisations
Authorization: Bearer TOKEN
```

**Query Parameters:**
- `status` - Filter by status: `active`, `suspended`, `archived`
- `search` - Search in name and description

**Response:** `200 OK`
```json
[
  {
    "id": "org-uuid",
    "name": "Acme Corporation",
    "slug": "acme-corp",
    "status": "active",
    "created_at": "2026-03-08T15:30:00Z"
  }
]
```

#### Get Organisation

Get details of a specific organisation.

```bash
GET /organisations/:id
Authorization: Bearer TOKEN
```

**Response:** `200 OK`
```json
{
  "id": "org-uuid",
  "name": "Acme Corporation",
  "slug": "acme-corp",
  "description": "Main company workspace",
  "status": "active",
  "created_at": "2026-03-08T15:30:00Z",
  "created_by": {
    "id": "user-uuid",
    "email": "user@example.com"
  },
  "members": [
    {
      "id": "member-uuid",
      "user_id": "user-uuid",
      "role": "owner",
      "joined_at": "2026-03-08T15:30:00Z",
      "user": {
        "id": "user-uuid",
        "email": "user@example.com"
      }
    }
  ]
}
```

#### Update Organisation

Update organisation details. Requires Admin or Owner role.

```bash
PATCH /organisations/:id
Authorization: Bearer TOKEN
Content-Type: application/json

{
  "name": "Acme Corp Updated",
  "description": "Updated description"
}
```

#### Delete Organisation

Soft delete an organisation. Requires Owner role.

```bash
DELETE /organisations/:id
Authorization: Bearer TOKEN
```

**Note:** Organisation must be archived first and have no members except the owner.

---

### Organisation Members

#### Add Member

Invite a user to the organisation. Requires Admin or Owner role.

```bash
POST /organisations/:id/members
Authorization: Bearer TOKEN
Content-Type: application/json

{
  "user_email": "newuser@example.com",
  "role": "member"
}
```

**Roles:**
- `owner` - Full control (cannot be assigned via API)
- `admin` - Can manage members and settings
- `member` - Standard access

**Response:** `201 Created`
```json
{
  "id": "member-uuid",
  "organisation_id": "org-uuid",
  "user_id": "user-uuid",
  "role": "member",
  "joined_at": "2026-03-08T16:00:00Z",
  "user": {
    "id": "user-uuid",
    "email": "newuser@example.com"
  }
}
```

#### Update Member Role

Change a member's role. Requires Owner role.

```bash
PATCH /organisations/:id/members/:member_id
Authorization: Bearer TOKEN
Content-Type: application/json

{
  "role": "admin"
}
```

#### Remove Member

Remove a member from the organisation. Requires Admin or Owner role.

```bash
DELETE /organisations/:id/members/:member_id
Authorization: Bearer TOKEN
```

#### Leave Organisation

Leave an organisation. Owners cannot leave.

```bash
POST /organisations/:id/leave
Authorization: Bearer TOKEN
```

---

### Projects

All project endpoints operate within the current organisation context.

#### Create Project

```bash
POST /projects
Authorization: Bearer TOKEN
Content-Type: application/json

{
  "name": "Website Redesign",
  "description": "Redesign the company website"
}
```

**Response:** `201 Created`
```json
{
  "id": "project-uuid",
  "name": "Website Redesign",
  "description": "Redesign the company website",
  "status": "active",
  "created_at": "2026-03-08T15:30:00Z",
  "created_by_user_id": "user-uuid"
}
```

#### List Projects

```bash
GET /projects
Authorization: Bearer TOKEN
```

**Query Parameters:**
- `status` - Filter by status: `active`, `archived`, `completed`
- `search` - Search in name and description

#### Get Project

```bash
GET /projects/:id
Authorization: Bearer TOKEN
```

#### Update Project

```bash
PATCH /projects/:id
Authorization: Bearer TOKEN
Content-Type: application/json

{
  "name": "Updated Name",
  "status": "archived"
}
```

#### Delete Project

```bash
DELETE /projects/:id
Authorization: Bearer TOKEN
```

---

### Tasks

#### Human Task Endpoints

Tasks can be created and managed by human users.

**Create Task:**
```bash
POST /tasks
Authorization: Bearer TOKEN
Content-Type: application/json

{
  "title": "Implement feature X",
  "description": "Detailed description",
  "priority": "high",
  "project_id": "project-uuid",
  "assigned_agent_id": "agent-uuid"
}
```

**List Tasks:**
```bash
GET /tasks?status=pending&priority=high
Authorization: Bearer TOKEN
```

**Query Parameters:**
- `status` - Filter by status
- `priority` - Filter by priority: `low`, `medium`, `high`, `critical`
- `agent_id` - Filter by assigned agent
- `project_id` - Filter by project
- `search` - Search in title and description

**Get Task:**
```bash
GET /tasks/:id
Authorization: Bearer TOKEN
```

**Update Task:**
```bash
PATCH /tasks/:id
Authorization: Bearer TOKEN
Content-Type: application/json

{
  "title": "Updated title",
  "priority": "critical"
}
```

**Assign Agent:**
```bash
PATCH /tasks/:id/assign
Authorization: Bearer TOKEN
Content-Type: application/json

{
  "agent_id": "agent-uuid"
}
```

**Delete Task:**
```bash
DELETE /tasks/:id
Authorization: Bearer TOKEN
```

#### Agent Task Endpoints

Agents can create and manage their own tasks.

**Create Task (Agent):**
```bash
POST /agent/tasks
X-API-Key: sk_agent_your_key
Content-Type: application/json

{
  "title": "Agent task",
  "description": "Task created by agent",
  "priority": "medium"
}
```

The task is automatically assigned to the authenticated agent.

**List My Tasks (Agent):**
```bash
GET /agent/tasks
X-API-Key: sk_agent_your_key
```

**Update Task Status (Agent):**
```bash
PATCH /agent/tasks/:id/status
X-API-Key: sk_agent_your_key
Content-Type: application/json

{
  "status": "completed"
}
```

**Valid Status Values:**
- `pending` - Not started
- `in_progress` - Working on it
- `completed` - Done
- `blocked` - Cannot proceed

---

### Agents

#### Create Agent

Create an AI agent for your organisation.

```bash
POST /agents
Authorization: Bearer TOKEN
Content-Type: application/json

{
  "name": "Data Processor",
  "description": "Processes data files",
  "role": "regular"
}
```

**Agent Roles:**
- `regular` - Can only update own tasks
- `supervisor` - Can update any task, create agents
- `admin` - Full permissions

**Response:** `201 Created`
```json
{
  "id": "agent-uuid",
  "name": "Data Processor",
  "api_key": "sk_agent_generated_key",
  "description": "Processes data files",
  "role": "regular",
  "enabled": true,
  "created_at": "2026-03-08T15:30:00Z"
}
```

**⚠️ Important:** Save the `api_key` immediately. It cannot be retrieved again.

#### List Agents

```bash
GET /agents
Authorization: Bearer TOKEN
```

#### Update Agent

```bash
PATCH /agents/:id
Authorization: Bearer TOKEN
Content-Type: application/json

{
  "name": "Updated Name",
  "enabled": false
}
```

#### Delete Agent

```bash
DELETE /agents/:id
Authorization: Bearer TOKEN
```

**Note:** Cannot delete agents with assigned tasks.

---

## Migration Guide

### Migrating from v1.0 to v2.0

Version 2.0 introduces multi-organisation support. Here's how to migrate your existing integration:

### 1. Check for Organisation Context

In v2.0, most endpoints require an organisation context. Update your code to:

```javascript
// Before (v1.0)
const response = await fetch('/api/tasks', {
  headers: { 'Authorization': `Bearer ${token}` }
});

// After (v2.0)
const response = await fetch('/api/tasks', {
  headers: {
    'Authorization': `Bearer ${token}`,
    'X-Organisation-ID': organisationId
  }
});
```

### 2. Handle JWT Token Changes

The JWT token now includes `current_org_id`:

```javascript
// Decode JWT to get organisation context
const payload = JSON.parse(atob(token.split('.')[1]));
const currentOrgId = payload.current_org_id;

if (!currentOrgId) {
  // User needs to select or create an organisation
  window.location.href = '/select-organisation';
}
```

### 3. Create Organisation on Registration

For new users, create an organisation after registration:

```javascript
// 1. Register user
const registerResponse = await fetch('/api/auth/register', {
  method: 'POST',
  body: JSON.stringify({ email, password })
});
const { token } = await registerResponse.json();

// 2. Create organisation
const orgResponse = await fetch('/api/organisations', {
  method: 'POST',
  headers: { 'Authorization': `Bearer ${token}` },
  body: JSON.stringify({
    name: 'My Workspace',
    slug: generateSlug(email),
    description: 'Personal workspace'
  })
});
```

### 4. Update Agent Integration

Agent API keys are now organisation-scoped:

```javascript
// Create agent in specific organisation
const agentResponse = await fetch('/api/agents', {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${token}`,
    'X-Organisation-ID': organisationId
  },
  body: JSON.stringify({
    name: 'My Agent',
    description: 'Task automation agent'
  })
});

const { api_key } = await agentResponse.json();

// Use API key for agent requests
const tasksResponse = await fetch('/api/agent/tasks', {
  headers: { 'X-API-Key': api_key }
});
```

### 5. Default Organisation (Existing Users)

If you're upgrading an existing deployment, all existing data is migrated to a **default organisation**. Users can:

1. Continue using the default organisation
2. Create new organisations and migrate data
3. Be added to other organisations

---

## Examples

### Example 1: Multi-Organisation Workflow

```bash
# 1. Login
TOKEN=$(curl -s -X POST https://todo.formatho.com/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password"}' | jq -r '.token')

# 2. List organisations
curl -s https://todo.formatho.com/api/organisations \
  -H "Authorization: Bearer $TOKEN" | jq

# 3. Create a new organisation
ORG_ID=$(curl -s -X POST https://todo.formatho.com/api/organisations \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"New Project","slug":"new-project"}' | jq -r '.id')

# 4. Create a project in the new organisation
curl -X POST https://todo.formatho.com/api/projects \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-Organisation-ID: $ORG_ID" \
  -H "Content-Type: application/json" \
  -d '{"name":"Q1 Goals","description":"First quarter objectives"}'

# 5. Create a task
curl -X POST https://todo.formatho.com/api/tasks \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-Organisation-ID: $ORG_ID" \
  -H "Content-Type: application/json" \
  -d '{"title":"Complete API docs","priority":"high"}'
```

### Example 2: Team Collaboration

```bash
# 1. Owner creates organisation
curl -X POST https://todo.formatho.com/api/organisations \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name":"Team Alpha","slug":"team-alpha"}'

# 2. Owner adds team members
curl -X POST https://todo.formatho.com/api/organisations/$ORG_ID/members \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"user_email":"teammate@example.com","role":"member"}'

# 3. Create an agent for the team
curl -X POST https://todo.formatho.com/api/agents \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-Organisation-ID: $ORG_ID" \
  -d '{"name":"Team Bot","description":"Team automation agent"}'

# 4. Assign agent to a task
curl -X PATCH https://todo.formatho.com/api/tasks/$TASK_ID/assign \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"agent_id":"agent-uuid"}'
```

### Example 3: Agent Task Management

```python
import requests

API_KEY = "sk_agent_your_key"
BASE_URL = "https://todo.formatho.com/api"

headers = {"X-API-Key": API_KEY}

# Create a task
response = requests.post(
    f"{BASE_URL}/agent/tasks",
    headers=headers,
    json={
        "title": "Process daily reports",
        "description": "Generate and send daily report",
        "priority": "high"
    }
)
task = response.json()
print(f"Created task: {task['id']}")

# Update status to in_progress
requests.patch(
    f"{BASE_URL}/agent/tasks/{task['id']}/status",
    headers=headers,
    json={"status": "in_progress"}
)

# Do the work...
process_reports()

# Mark as completed
requests.patch(
    f"{BASE_URL}/agent/tasks/{task['id']}/status",
    headers=headers,
    json={"status": "completed"}
)
```

---

## Error Handling

### Error Response Format

All errors follow a consistent format:

```json
{
  "error": "Error message",
  "code": "ERROR_CODE",
  "details": {}
}
```

### Common Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `UNAUTHORIZED` | 401 | Missing or invalid authentication |
| `FORBIDDEN` | 403 | Authenticated but not authorised |
| `NOT_FOUND` | 404 | Resource not found |
| `VALIDATION_ERROR` | 400 | Invalid request data |
| `ORG_REQUIRED` | 403 | Organisation context required |
| `NOT_MEMBER` | 403 | Not a member of this organisation |

### Handling Organisation Errors

```javascript
try {
  const response = await fetch('/api/tasks', {
    headers: {
      'Authorization': `Bearer ${token}`,
      'X-Organisation-ID': orgId
    }
  });

  if (!response.ok) {
    const error = await response.json();
    
    if (error.code === 'ORG_REQUIRED') {
      // Redirect to organisation selection
      window.location.href = '/select-organisation';
    } else if (error.code === 'NOT_MEMBER') {
      // User removed from organisation
      showToast('You are no longer a member of this organisation');
    }
  }
} catch (err) {
  console.error('API error:', err);
}
```

---

## Rate Limits

| Endpoint Type | Rate Limit |
|---------------|------------|
| Authentication | 20 req/min |
| Standard API | 100 req/min |
| Agent API | 60 req/min |

Rate limit headers are included in responses:
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1712505600
```

---

## Support

- **Documentation:** https://docs.formatho.com
- **API Status:** https://status.formatho.com
- **Support:** support@formatho.com
- **GitHub:** https://github.com/formatho/agent-todo

---

**Last Updated:** March 8, 2026  
**API Version:** 2.0.0
