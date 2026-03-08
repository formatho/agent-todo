# Multi-Organisation Architecture Requirements

**Task:** 1d763554-bf9a-4933-a7dd-cb901c18c507  
**Date:** 2026-03-08  
**Author:** agent-todo

---

## 1. Executive Summary

This document outlines the architecture requirements for adding multi-organisation support to the Agent Todo Platform. The goal is to enable complete data isolation between organisations while maintaining the current API structure and minimising disruption to existing functionality.

---

## 2. Current State Analysis

### 2.1 Current Data Model

**Existing Models:**
- **User** - Human users (global, not organisation-scoped)
- **Agent** - AI agents (global, not organisation-scoped)
- **Project** - Task grouping (belongs to user)
- **Task** - Todo items (belongs to project, can be assigned to agent)
- **TaskEvent** - Audit log for tasks
- **TaskComment** - Comments on tasks

**Current Relationships:**
```
User → creates → Project
User → creates → Agent
Project → has many → Tasks
Task → belongs to → Project
Task → assigned to → Agent
Task → created by → User or Agent
```

**Current Limitations:**
- No organisation concept
- All users can potentially see all data
- No tenant isolation
- No billing/plan management

---

## 3. Proposed Architecture

### 3.1 Data Model Changes

#### 3.1.1 New Organisation Model

```go
type Organisation struct {
    Base
    Name          string                `gorm:"not null" json:"name"`
    Slug          string                `gorm:"uniqueIndex;not null" json:"slug"`
    Description   string                `json:"description"`
    Plan          OrganisationPlan      `gorm:"not null;default:'free'" json:"plan"`
    Status        OrganisationStatus    `gorm:"not null;default:'active'" json:"status"`
    OwnerID       uuid.UUID             `gorm:"type:uuid;not null" json:"owner_id"`
    Owner         *User                 `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
    Settings      OrganisationSettings  `gorm:"type:jsonb" json:"settings"`
    Members       []OrganisationMember  `gorm:"foreignKey:OrganisationID" json:"members,omitempty"`
    Projects      []Project             `gorm:"foreignKey:OrganisationID" json:"projects,omitempty"`
    Agents        []Agent               `gorm:"foreignKey:OrganisationID" json:"agents,omitempty"`
}

type OrganisationPlan string
const (
    PlanFree      OrganisationPlan = "free"
    PlanPro       OrganisationPlan = "pro"
    PlanEnterprise OrganisationPlan = "enterprise"
)

type OrganisationStatus string
const (
    OrgStatusActive    OrganisationStatus = "active"
    OrgStatusSuspended OrganisationStatus = "suspended"
    OrgStatusDeleted   OrganisationStatus = "deleted"
)

type OrganisationSettings struct {
    MaxAgents       int  `json:"max_agents"`
    MaxProjects     int  `json:"max_projects"`
    MaxTasks        int  `json:"max_tasks"`
    AllowedFeatures []string `json:"allowed_features"`
}
```

#### 3.1.2 Organisation Membership Model

```go
type OrganisationMember struct {
    Base
    OrganisationID uuid.UUID       `gorm:"type:uuid;not null" json:"organisation_id"`
    Organisation   *Organisation   `gorm:"foreignKey:OrganisationID" json:"organisation,omitempty"`
    UserID         uuid.UUID       `gorm:"type:uuid;not null" json:"user_id"`
    User           *User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
    Role           MemberRole      `gorm:"not null;default:'member'" json:"role"`
    JoinedAt       time.Time       `json:"joined_at"`
}

type MemberRole string
const (
    RoleOwner      MemberRole = "owner"      // Full control, billing access
    RoleAdmin      MemberRole = "admin"      // Full control, no billing
    RoleMember     MemberRole = "member"     // Can create/manage own tasks
    RoleViewer     MemberRole = "viewer"     // Read-only access
)
```

#### 3.1.3 Updated Existing Models

**User Model:**
```go
type User struct {
    Base
    Email           string                `gorm:"uniqueIndex;not null" json:"email"`
    PasswordHash    string                `gorm:"not null" json:"-"`
    Memberships     []OrganisationMember  `gorm:"foreignKey:UserID" json:"memberships,omitempty"`
    CurrentOrgID    *uuid.UUID            `gorm:"type:uuid" json:"current_org_id"`
}
```

**Project Model:**
```go
type Project struct {
    Base
    OrganisationID  uuid.UUID       `gorm:"type:uuid;not null;index" json:"organisation_id"` // NEW
    Organisation    *Organisation   `gorm:"foreignKey:OrganisationID" json:"organisation,omitempty"` // NEW
    Name            string          `gorm:"not null;index:idx_org_project,composite:org_project" json:"name"` // Composite index
    Description     string          `json:"description"`
    Status          ProjectStatus   `gorm:"not null;default:'active'" json:"status"`
    CreatedByUserID uuid.UUID       `gorm:"type:uuid;not null" json:"created_by_user_id"`
    CreatedBy       *User           `gorm:"foreignKey:CreatedByUserID" json:"created_by,omitempty"`
    Tasks           []Task          `gorm:"foreignKey:ProjectID" json:"tasks,omitempty"`
}
```

**Agent Model:**
```go
type Agent struct {
    Base
    OrganisationID  uuid.UUID       `gorm:"type:uuid;not null;index" json:"organisation_id"` // NEW
    Organisation    *Organisation   `gorm:"foreignKey:OrganisationID" json:"organisation,omitempty"` // NEW
    Name            string          `gorm:"not null;index:idx_org_agent,composite:org_agent" json:"name"` // Composite index
    APIKey          string          `gorm:"uniqueIndex;not null" json:"api_key"`
    Description     string          `json:"description"`
    Role            AgentRole       `gorm:"not null;default:'regular'" json:"role"`
    Enabled         bool            `gorm:"not null;default:true" json:"enabled"`
}
```

**Task Model:**
```go
type Task struct {
    Base
    OrganisationID  uuid.UUID       `gorm:"type:uuid;not null;index" json:"organisation_id"` // NEW
    Organisation    *Organisation   `gorm:"foreignKey:OrganisationID" json:"organisation,omitempty"` // NEW
    Title           string          `gorm:"not null;index:idx_org_task,composite:org_task" json:"title"`
    Description     string          `json:"description"`
    Status          TaskStatus      `gorm:"not null;default:'pending'" json:"status"`
    Priority        TaskPriority    `gorm:"not null;default:'medium'" json:"priority"`
    DueDate         *time.Time      `json:"due_date"`
    ProjectID       *uuid.UUID      `gorm:"type:uuid;index" json:"project_id"`
    CreatedByUserID *uuid.UUID      `gorm:"type:uuid" json:"created_by_user_id"`
    CreatedByAgentID *uuid.UUID     `gorm:"type:uuid" json:"created_by_agent_id"`
    AssignedAgentID *uuid.UUID      `gorm:"type:uuid" json:"assigned_agent_id"`
    DeletedAt       gorm.DeletedAt  `gorm:"index" json:"deleted_at,omitempty"`
    Project         *Project        `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
    CreatedBy       *User           `gorm:"foreignKey:CreatedByUserID" json:"created_by,omitempty"`
    CreatedByAgent  *Agent          `gorm:"foreignKey:CreatedByAgentID" json:"created_by_agent,omitempty"`
    AssignedAgent   *Agent          `gorm:"foreignKey:AssignedAgentID" json:"assigned_agent,omitempty"`
    Comments        []TaskComment   `gorm:"foreignKey:TaskID" json:"comments,omitempty"`
    Events          []TaskEvent     `gorm:"foreignKey:TaskID" json:"events,omitempty"`
}
```

### 3.2 Database Schema Changes

#### 3.2.1 New Tables

**organisations:**
- id (UUID, PK)
- name (VARCHAR, NOT NULL)
- slug (VARCHAR, UNIQUE, NOT NULL)
- description (TEXT)
- plan (VARCHAR, NOT NULL, DEFAULT 'free')
- status (VARCHAR, NOT NULL, DEFAULT 'active')
- owner_id (UUID, FK → users.id)
- settings (JSONB)
- created_at (TIMESTAMP)
- updated_at (TIMESTAMP)

**organisation_members:**
- id (UUID, PK)
- organisation_id (UUID, FK → organisations.id)
- user_id (UUID, FK → users.id)
- role (VARCHAR, NOT NULL, DEFAULT 'member')
- joined_at (TIMESTAMP)
- created_at (TIMESTAMP)
- updated_at (TIMESTAMP)
- UNIQUE(organisation_id, user_id)

#### 3.2.2 Modified Tables

**users:**
- ADD current_org_id (UUID, FK → organisations.id, NULLABLE)

**projects:**
- ADD organisation_id (UUID, FK → organisations.id, NOT NULL)
- ADD INDEX idx_org_project (organisation_id, name)
- DROP existing unique constraints (now scoped to organisation)

**agents:**
- ADD organisation_id (UUID, FK → organisations.id, NOT NULL)
- ADD INDEX idx_org_agent (organisation_id, name)
- DROP existing unique constraints (now scoped to organisation)

**tasks:**
- ADD organisation_id (UUID, FK → organisations.id, NOT NULL)
- ADD INDEX idx_org_task (organisation_id, title)

**task_events:**
- ADD organisation_id (UUID, FK → organisations.id, NOT NULL)
- ADD INDEX idx_org_event (organisation_id)

**task_comments:**
- ADD organisation_id (UUID, FK → organisations.id, NOT NULL)
- ADD INDEX idx_org_comment (organisation_id)

---

## 4. API Endpoint Changes

### 4.1 New Endpoints

#### Organisation CRUD
```
POST   /api/organisations                 - Create organisation
GET    /api/organisations                 - List user's organisations
GET    /api/organisations/:id             - Get organisation details
PATCH  /api/organisations/:id             - Update organisation
DELETE /api/organisations/:id             - Delete organisation (soft delete)
```

#### Organisation Membership
```
POST   /api/organisations/:id/members     - Invite member
GET    /api/organisations/:id/members     - List members
PATCH  /api/organisations/:id/members/:user_id - Update member role
DELETE /api/organisations/:id/members/:user_id - Remove member
```

#### Organisation Settings
```
GET    /api/organisations/:id/settings    - Get organisation settings
PATCH  /api/organisations/:id/settings    - Update organisation settings
```

### 4.2 Modified Endpoints

All existing endpoints now require organisation context:

**Option 1: Header-based (Recommended)**
```
GET /api/tasks
Headers:
  X-Organisation-ID: <org_id>
```

**Option 2: Path-based**
```
GET /api/organisations/:org_id/tasks
```

**Recommendation:** Use Header-based approach to maintain backward compatibility and cleaner URLs.

### 4.3 Authentication & Authorization Flow

1. User logs in → receives JWT with `user_id`
2. User selects organisation → JWT updated with `current_org_id`
3. All API requests include `X-Organisation-ID` header
4. Middleware validates:
   - User is member of organisation
   - User has required role for action
   - All data queries filter by organisation_id

---

## 5. Access Control Requirements

### 5.1 Role-Based Access Control (RBAC)

**Owner:**
- Full organisation control
- Billing management
- Delete organisation
- Manage all members

**Admin:**
- Manage projects, tasks, agents
- Invite/remove members (except owners)
- Update organisation settings (non-billing)

**Member:**
- Create/manage own projects and tasks
- View all organisation data
- Manage own profile

**Viewer:**
- Read-only access to organisation data
- No create/update/delete permissions

### 5.2 Resource-Level Permissions

| Resource | Owner | Admin | Member | Viewer |
|----------|-------|-------|--------|--------|
| Organisation | CRUD | RU | R | R |
| Projects | CRUD | CRUD | CRD (own) | R |
| Tasks | CRUD | CRUD | CRD (own) | R |
| Agents | CRUD | CRUD | R | R |
| Members | CRUD | CRD | R | R |

### 5.3 Implementation

Use middleware to check permissions:

```go
func RequireOrgRole(role MemberRole) gin.HandlerFunc {
    return func(c *gin.Context) {
        orgID := c.GetHeader("X-Organisation-ID")
        userID := c.GetString("user_id")
        
        member, err := getMember(orgID, userID)
        if err != nil || !hasRequiredRole(member.Role, role) {
            c.AbortWithStatusJSON(403, gin.H{"error": "Forbidden"})
            return
        }
        
        c.Next()
    }
}
```

---

## 6. Migration Strategy

### 6.1 Phase 1: Database Migration (Zero Downtime)

1. **Create new tables** (organisations, organisation_members)
2. **Add nullable columns** (organisation_id to all tables)
3. **Deploy application changes** (organisation-aware code)
4. **Create default organisation** for existing data
5. **Backfill data** with default organisation_id
6. **Make columns NOT NULL** with constraints
7. **Add foreign keys and indexes**

### 6.2 Default Organisation Creation

For existing deployment:
```sql
-- Create default organisation
INSERT INTO organisations (id, name, slug, plan, status, owner_id)
SELECT 
    gen_random_uuid(),
    'Default Organisation',
    'default',
    'enterprise',
    'active',
    (SELECT id FROM users LIMIT 1);

-- Assign all existing data to default organisation
UPDATE projects SET organisation_id = (SELECT id FROM organisations WHERE slug = 'default');
UPDATE agents SET organisation_id = (SELECT id FROM organisations WHERE slug = 'default');
UPDATE tasks SET organisation_id = (SELECT id FROM organisations WHERE slug = 'default');

-- Make all users members of default organisation
INSERT INTO organisation_members (organisation_id, user_id, role, joined_at)
SELECT 
    (SELECT id FROM organisations WHERE slug = 'default'),
    id,
    'owner',
    NOW()
FROM users;
```

### 6.3 Rollback Plan

If migration fails:
1. Remove organisation_id columns (nullable phase)
2. Revert application code
3. Drop new tables
4. System returns to single-tenant mode

---

## 7. Data Isolation Requirements

### 7.1 Database-Level Isolation

**All queries MUST include organisation_id filter:**

```go
// Middleware to inject organisation filter
func OrgFilterMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        orgID := c.GetHeader("X-Organisation-ID")
        c.Set("organisation_id", orgID)
        c.Next()
    }
}

// Query scope for GORM
func OrganisationScope(db *gorm.DB, orgID string) *gorm.DB {
    return db.Where("organisation_id = ?", orgID)
}
```

### 7.2 API Response Filtering

All API responses must only include data from the current organisation:
- Lists filtered by organisation_id
- Single entity lookups verify organisation_id
- No cross-organisation references allowed

### 7.3 Cross-Organisation Prevention

- Users can only see organisations they're members of
- Agents can only access data in their organisation
- API keys are scoped to single organisation
- No global queries permitted

---

## 8. Testing Strategy

### 8.1 Unit Tests

- Test organisation CRUD operations
- Test membership management
- Test role-based permissions
- Test data isolation in queries

### 8.2 Integration Tests

- Test complete user flow across organisations
- Test API endpoint isolation
- Test concurrent organisation operations
- Test data migration

### 8.3 Isolation Tests

```go
func TestOrganisationIsolation(t *testing.T) {
    // Create two organisations
    org1 := createOrganisation("Org 1")
    org2 := createOrganisation("Org 2")
    
    // Create data in org1
    task1 := createTask(org1.ID, "Task 1")
    
    // Try to access from org2
    tasks := listTasks(org2.ID)
    
    // Should be empty
    assert.Empty(t, tasks)
    
    // Try to access task1 from org2
    _, err := getTask(task1.ID, org2.ID)
    assert.Error(t, err) // Should fail
}
```

---

## 9. Performance Considerations

### 9.1 Indexing Strategy

**Composite Indexes:**
```sql
CREATE INDEX idx_org_project ON projects(organisation_id, name);
CREATE INDEX idx_org_agent ON agents(organisation_id, name);
CREATE INDEX idx_org_task ON tasks(organisation_id, status, priority);
CREATE INDEX idx_org_event ON task_events(organisation_id, created_at DESC);
```

### 9.2 Query Optimization

- Always query by organisation_id first (uses index)
- Avoid SELECT * (fetch only needed columns)
- Use pagination for large datasets
- Cache organisation membership in JWT

### 9.3 Connection Pooling

No changes needed - organisation filtering happens at query level, not connection level.

---

## 10. Security Considerations

### 10.1 JWT Token Updates

```json
{
  "user_id": "uuid",
  "email": "user@example.com",
  "current_org_id": "org_uuid",
  "orgs": [
    {"id": "org1", "role": "owner"},
    {"id": "org2", "role": "member"}
  ]
}
```

### 10.2 API Key Scoping

Agent API keys should include organisation_id:
- Current: `sk_agent_86f63599-1044-4800-9853-72b64ad1fb90`
- Future: `sk_agent_<org_id>_<random>`

### 10.3 Audit Logging

Log all cross-organisation access attempts:
- User ID
- Source organisation
- Target organisation
- Action attempted
- Timestamp

---

## 11. Implementation Priority

**Phase 1 (MVP):**
1. Database migrations
2. Organisation CRUD endpoints
3. Membership management
4. Data isolation middleware
5. Update existing endpoints with org context

**Phase 2 (Enhancement):**
1. Role-based access control
2. Organisation settings
3. API key scoping
4. Plan limits enforcement

**Phase 3 (Polish):**
1. Billing integration
2. Usage analytics per org
3. Custom domains
4. SSO integration

---

## 12. Success Criteria

- ✅ All existing data migrated to default organisation
- ✅ All API endpoints organisation-aware
- ✅ Complete data isolation between organisations
- ✅ Role-based access control working
- ✅ Zero data leakage in tests
- ✅ Performance maintained (query time < 100ms)
- ✅ Backward compatibility with existing clients

---

## 13. Open Questions

1. **Default organisation slug:** Use 'default' or generate unique slug?
2. **User registration:** Create personal organisation automatically?
3. **Billing:** Stripe integration? Usage-based pricing?
4. **Data export:** Allow organisation data export?
5. **Deletion:** Soft delete with grace period? Hard delete after X days?

---

## 14. References

- [GORM Composite Indexes](https://gorm.io/docs/indexes.html)
- [Multi-tenant Architecture Patterns](https://www.postgresql.org/docs/current/ddl-schemas.html)
- [RBAC Best Practices](https://auth0.com/docs/manage-users/access-control/rbac)

---

**Document Status:** Draft v1.0  
**Next Steps:** Review with team, create database migration scripts
