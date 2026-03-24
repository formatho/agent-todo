# Blocked Status Feature Implementation

## Overview
Added "Blocked" status option to tasks, allowing tasks to be marked as blocked when waiting on external dependencies.

## Changes Made

### 1. Backend Changes

#### Database Migration
- **File:** `backend/db/migrations/000019_add_blocked_status_to_tasks.up.sql`
- **Change:** Added 'blocked' value to task_status enum
- **SQL:**
  ```sql
  ALTER TYPE task_status ADD VALUE 'blocked';
  ```

#### Models
- **File:** `backend/models/models.go`
- **Change:** Added `TaskStatusBlocked` constant
- **Code:**
  ```go
  const (
      TaskStatusPending    TaskStatus = "pending"
      TaskStatusInProgress TaskStatus = "in_progress"
      TaskStatusCompleted  TaskStatus = "completed"
      TaskStatusFailed     TaskStatus = "failed"
      TaskStatusBlocked    TaskStatus = "blocked"  // NEW
  )
  ```

### 2. Frontend Changes

#### Kanban Board
- **File:** `frontend/src/components/KanbanBoard.vue`
- **Changes:**
  - Added "Blocked" column (5th column)
  - Updated `loadTasks()` to filter blocked tasks
  - Added CSS styling for `.column-blocked`

#### Status Colors
- **File:** `frontend/src/utils/agentColors.js`
- **Change:** Added `blocked` status to `STATUS_COLORS` object
- **Color Scheme:**
  - Background: `#EDE9FE` (Light purple)
  - Text: `#5B21B6` (Dark purple)
  - Border: `#8B5CF6` (Purple)
  - Icon: `🚧`

#### Task Card
- **File:** `frontend/src/components/TaskCard.vue`
- **Change:** Added `.task-card.status-blocked` CSS styling
- **Style:** Purple left border to indicate blocked state

#### Task Grid
- **File:** `frontend/src/components/TaskGrid.vue`
- **Changes:**
  - Added "Blocked" option to status filter dropdown
  - Added CSS styling for `.status-badge.blocked`

## Deployment Instructions

### Apply Database Migration
Run the migration to add the 'blocked' value to the database:

**Option 1: Using migrate CLI**
```bash
cd backend
make migrate-up
```

**Option 2: Manual SQL**
```bash
docker exec -i agent-todo-db psql -U agent_todo -d agent_todo -c "ALTER TYPE task_status ADD VALUE 'blocked';"
```

**Option 3: Via psql**
```bash
psql -h localhost -U agent_todo -d agent_todo -c "ALTER TYPE task_status ADD VALUE 'blocked';"
```

### Rebuild and Restart Services
```bash
cd /Users/studio/sandbox/formatho/agent-todo

# Rebuild frontend
cd frontend
npm run build

# Restart backend
cd ../backend
go run cmd/api/main.go
```

## Testing

### 1. Test API
Create a task with blocked status:
```bash
curl -X POST "https://todo.formatho.com/api/agent/tasks" \
  -H "X-API-Key: YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Test Blocked Task",
    "description": "This task is blocked",
    "status": "blocked",
    "priority": "medium"
  }'
```

### 2. Test Frontend
1. Open the Kanban Board
2. Verify the "Blocked" column appears
3. Create a new task or drag an existing task to "Blocked"
4. Verify the status badge displays correctly with the 🚧 icon and purple color
5. Verify the Task Grid filter includes "Blocked" option

### 3. Test Status Transitions
- Test moving tasks between all statuses
- Verify blocked tasks have proper visual styling
- Test filtering by "blocked" status in Task Grid

## Rollback

If you need to rollback this change:

**Note:** PostgreSQL doesn't support removing enum values. To rollback:
1. Drop and recreate the enum type without the 'blocked' value
2. Update all blocked tasks to another status first
3. Migrate data
4. Recreate enum type

See `backend/db/migrations/000019_add_blocked_status_to_tasks.down.sql` for more details.

## Summary

✅ Database migration created
✅ Backend models updated
✅ Frontend components updated (Kanban, Task Card, Task Grid)
✅ Status colors and styling added
⚠️ Migration needs to be applied to live database

The feature is complete and ready for deployment. The database migration is the only remaining step.
