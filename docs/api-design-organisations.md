# Organisation API Design

**Version:** 1.0  
**Status:** Draft  
**Created:** March 8, 2026  
**Author:** Agent-Todo

## Overview

This document defines the REST API endpoints for organisation management in the Agent Todo platform. Organisations enable multi-tenant isolation, allowing different teams or companies to use the platform with complete data separation.

## Goals

1. **Multi-tenant Isolation:** Ensure complete data separation between organisations
2. **Backward Compatibility:** Support existing single-organisation setup
3. **Flexible Membership:** Allow users to belong to multiple organisations
4. **Secure Access Control:** Enforce organisation-based permissions at the API level

---

## Data Model

### Organisation Model

```go
type Organisation struct {
    Base
    Name            string              `gorm:"not null" json:"name"`
    Slug            string              `gorm:"uniqueIndex;not null" json:"slug"`
    Description     string              `json:"description"`
    Status          OrganisationStatus  `gorm:"not null;default:'active'" json:"status"`
    CreatedByUserID uuid.UUID           `gorm:"type:uuid;not null" json:"created_by_user_id"`
    CreatedBy       *User               `gorm:"foreignKey:CreatedByUserID" json:"created_by,omitempty"`
    Members         []OrganisationMember `gorm:"foreignKey:OrganisationID" json:"members,omitempty"`
    Projects        []Project           `gorm:"foreignKey:OrganisationID" json:"projects,omitempty"`
    Agents          []Agent             `gorm:"foreignKey:OrganisationID" json:"agents,omitempty"`
}

type OrganisationStatus string

const (
    OrganisationStatusActive    OrganisationStatus = "active"
    OrganisationStatusSuspended OrganisationStatus = "suspended"
    OrganisationStatusArchived  OrganisationStatus = "archived"
)
```

### Organisation Member Model

```go
type OrganisationMember struct {
    Base
    OrganisationID   uuid.UUID             `gorm:"type:uuid;not null" json:"organisation_id"`
    Organisation     *Organisation         `gorm:"foreignKey:OrganisationID" json:"organisation,omitempty"`
    UserID           uuid.UUID             `gorm:"type:uuid;not null" json:"user_id"`
    User             *User                 `gorm:"foreignKey:UserID" json:"user,omitempty"`
    Role             OrganisationMemberRole `gorm:"not null;default:'member'" json:"role"`
    JoinedAt         time.Time             `json:"joined_at"`
}

type OrganisationMemberRole string

const (
    OrganisationMemberRoleOwner  OrganisationMemberRole = "owner"
    OrganisationMemberRoleAdmin  OrganisationMemberRole = "admin"
    OrganisationMemberRoleMember OrganisationMemberRole = "member"
)
```

### Database Migration

Add `organisation_id` to existing models:

```go
// Project
type Project struct {
    // ... existing fields ...
    OrganisationID  *uuid.UUID    `gorm:"type:uuid;index" json:"organisation_id"`
    Organisation    *Organisation `gorm:"foreignKey:OrganisationID" json:"organisation,omitempty"`
}

// Task
type Task struct {
    // ... existing fields ...
    OrganisationID  *uuid.UUID    `gorm:"type:uuid;index" json:"organisation_id"`
    Organisation    *Organisation `gorm:"foreignKey:OrganisationID" json:"organisation,omitempty"`
}

// Agent
type Agent struct {
    // ... existing fields ...
    OrganisationID  *uuid.UUID    `gorm:"type:uuid;index" json:"organisation_id"`
    Organisation    *Organisation `gorm:"foreignKey:OrganisationID" json:"organisation,omitempty"`
}
```

---

## API Endpoints

### Base URL

```
https://todo.formatho.com/api/organisations
```

### Authentication

All organisation endpoints require authentication:
- **Human Users:** JWT Bearer token in `Authorization` header
- **AI Agents:** API key in `X-API-KEY` header (limited read-only access)

---

## Endpoints

### 1. Create Organisation

**POST** `/api/organisations`

Creates a new organisation. The creator automatically becomes the owner.

**Request:**

```json
{
  "name": "Formatho Technologies",
  "slug": "formatho",
  "description": "AI-powered task management platform"
}
```

**Validation Rules:**
- `name`: Required, 1-255 characters
- `slug`: Required, 3-50 characters, alphanumeric and hyphens only, must be unique
- `description`: Optional, max 1000 characters

**Response:** `201 Created`

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Formatho Technologies",
  "slug": "formatho",
  "description": "AI-powered task management platform",
  "status": "active",
  "created_at": "2026-03-08T15:30:00Z",
  "updated_at": "2026-03-08T15:30:00Z",
  "created_by": {
    "id": "fd8ada49-6ac2-4c83-83ea-b37cfab47650",
    "email": "admin@example.com"
  }
}
```

**Error Responses:**
- `400 Bad Request` - Invalid input data or slug already exists
- `401 Unauthorized` - Missing or invalid authentication
- `403 Forbidden` - User lacks permission to create organisations

---

### 2. List Organisations

**GET** `/api/organisations`

Lists organisations the authenticated user belongs to.

**Query Parameters:**
- `status` (optional): Filter by status (`active`, `suspended`, `archived`)
- `search` (optional): Search in name and description
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 20, max: 100)

**Response:** `200 OK`

```json
{
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "Formatho Technologies",
      "slug": "formatho",
      "description": "AI-powered task management platform",
      "status": "active",
      "created_at": "2026-03-08T15:30:00Z",
      "updated_at": "2026-03-08T15:30:00Z",
      "member_count": 5,
      "project_count": 12,
      "user_role": "owner"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 1,
    "total_pages": 1
  }
}
```

**Error Responses:**
- `401 Unauthorized` - Missing or invalid authentication

---

### 3. Get Organisation

**GET** `/api/organisations/:id`

Get details of a specific organisation.

**Response:** `200 OK`

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Formatho Technologies",
  "slug": "formatho",
  "description": "AI-powered task management platform",
  "status": "active",
  "created_at": "2026-03-08T15:30:00Z",
  "updated_at": "2026-03-08T15:30:00Z",
  "created_by": {
    "id": "fd8ada49-6ac2-4c83-83ea-b37cfab47650",
    "email": "admin@example.com"
  },
  "members": [
    {
      "id": "660e8400-e29b-41d4-a716-446655440001",
      "user_id": "fd8ada49-6ac2-4c83-83ea-b37cfab47650",
      "role": "owner",
      "joined_at": "2026-03-08T15:30:00Z",
      "user": {
        "id": "fd8ada49-6ac2-4c83-83ea-b37cfab47650",
        "email": "admin@example.com"
      }
    }
  ],
  "statistics": {
    "member_count": 5,
    "project_count": 12,
    "task_count": 48,
    "agent_count": 3
  }
}
```

**Error Responses:**
- `401 Unauthorized` - Missing or invalid authentication
- `403 Forbidden` - User is not a member of this organisation
- `404 Not Found` - Organisation not found

---

### 4. Update Organisation

**PATCH** `/api/organisations/:id`

Update organisation details. Requires `owner` or `admin` role.

**Request:**

```json
{
  "name": "Formatho Tech Inc.",
  "description": "Updated description",
  "status": "active"
}
```

**Validation Rules:**
- `name`: Optional, 1-255 characters
- `description`: Optional, max 1000 characters
- `status`: Optional, one of `active`, `suspended`, `archived`
- Note: `slug` cannot be changed after creation

**Response:** `200 OK`

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Formatho Tech Inc.",
  "slug": "formatho",
  "description": "Updated description",
  "status": "active",
  "updated_at": "2026-03-08T16:00:00Z"
}
```

**Error Responses:**
- `400 Bad Request` - Invalid input data
- `401 Unauthorized` - Missing or invalid authentication
- `403 Forbidden` - User lacks admin/owner permissions
- `404 Not Found` - Organisation not found

---

### 5. Delete Organisation

**DELETE** `/api/organisations/:id`

Soft delete an organisation. Requires `owner` role.

**Prerequisites:**
- Organisation must be in `archived` status
- All members except owner must be removed
- All projects must be deleted or transferred

**Response:** `204 No Content`

**Error Responses:**
- `401 Unauthorized` - Missing or invalid authentication
- `403 Forbidden` - User is not the organisation owner
- `404 Not Found` - Organisation not found
- `409 Conflict` - Organisation still has active resources

---

### 6. Add Organisation Member

**POST** `/api/organisations/:id/members`

Add a user to the organisation. Requires `owner` or `admin` role.

**Request:**

```json
{
  "user_email": "newuser@example.com",
  "role": "member"
}
```

**Validation Rules:**
- `user_email`: Required, valid email, user must exist
- `role`: Required, one of `admin`, `member` (cannot assign `owner`)

**Response:** `201 Created`

```json
{
  "id": "770e8400-e29b-41d4-a716-446655440002",
  "organisation_id": "550e8400-e29b-41d4-a716-446655440000",
  "user_id": "880e8400-e29b-41d4-a716-446655440003",
  "role": "member",
  "joined_at": "2026-03-08T16:30:00Z",
  "user": {
    "id": "880e8400-e29b-41d4-a716-446655440003",
    "email": "newuser@example.com"
  }
}
```

**Error Responses:**
- `400 Bad Request` - Invalid input or user already a member
- `401 Unauthorized` - Missing or invalid authentication
- `403 Forbidden` - User lacks admin/owner permissions
- `404 Not Found` - Organisation or user not found

---

### 7. Update Member Role

**PATCH** `/api/organisations/:id/members/:member_id`

Update a member's role. Requires `owner` role.

**Request:**

```json
{
  "role": "admin"
}
```

**Validation Rules:**
- `role`: Required, one of `admin`, `member`
- Cannot change owner's role
- Cannot demote yourself

**Response:** `200 OK`

```json
{
  "id": "770e8400-e29b-41d4-a716-446655440002",
  "role": "admin",
  "updated_at": "2026-03-08T17:00:00Z"
}
```

**Error Responses:**
- `400 Bad Request` - Invalid role or attempting to change owner
- `401 Unauthorized` - Missing or invalid authentication
- `403 Forbidden` - User is not the organisation owner
- `404 Not Found` - Organisation or member not found

---

### 8. Remove Organisation Member

**DELETE** `/api/organisations/:id/members/:member_id`

Remove a member from the organisation. Requires `owner` or `admin` role.

**Restrictions:**
- Cannot remove the owner
- Admins cannot remove other admins
- Cannot remove yourself (use leave endpoint)

**Response:** `204 No Content`

**Error Responses:**
- `401 Unauthorized` - Missing or invalid authentication
- `403 Forbidden` - User lacks permissions or attempting restricted action
- `404 Not Found` - Organisation or member not found

---

### 9. Leave Organisation

**POST** `/api/organisations/:id/leave`

Leave an organisation. Not available for owners.

**Response:** `204 No Content`

**Error Responses:**
- `400 Bad Request` - User is the owner (must transfer ownership first)
- `401 Unauthorized` - Missing or invalid authentication
- `404 Not Found` - Organisation not found or user is not a member

---

## Authentication & Access Control

### JWT Token Enhancement

Add organisation context to JWT claims:

```json
{
  "user_id": "fd8ada49-6ac2-4c83-83ea-b37cfab47650",
  "email": "admin@example.com",
  "current_organisation_id": "550e8400-e29b-41d4-a716-446655440000",
  "exp": 1712505600
}
```

### Access Control Middleware

```go
func OrganisationAccessMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get organisation ID from context or query param
        orgID := getOrganisationID(c)
        
        // Verify user is a member
        userID := middleware.GetUserID(c)
        if !isOrganisationMember(userID, orgID) {
            c.JSON(403, gin.H{"error": "Access denied"})
            c.Abort()
            return
        }
        
        // Set organisation context
        c.Set("organisation_id", orgID)
        c.Next()
    }
}
```

### Query Filtering

All data queries must include `organisation_id` filter:

```go
// In service layer
func (s *TaskService) List(userID, orgID uuid.UUID) ([]Task, error) {
    var tasks []Task
    err := s.db.Where("organisation_id = ?", orgID).Find(&tasks).Error
    return tasks, err
}
```

---

## Migration Strategy

### Phase 1: Data Model & Endpoints
1. Add Organisation and OrganisationMember models
2. Implement CRUD endpoints
3. Add organisation_id to existing models (nullable)

### Phase 2: Data Migration
1. Create default organisation ("formatho")
2. Migrate all existing data to default organisation
3. Add all existing users as members

### Phase 3: Access Control
1. Add organisation middleware
2. Update all queries to filter by organisation_id
3. Make organisation_id required for new records

### Phase 4: Frontend Integration
1. Add organisation switcher UI
2. Update all API calls to include organisation context
3. Test multi-organisation scenarios

---

## Backward Compatibility

### Default Organisation

- Create a "formatho" organisation during migration
- All existing data is migrated to this organisation
- Existing users become members with appropriate roles
- API continues to work without organisation_id for backward compatibility
- If no organisation context is provided, use default organisation

### API Versioning

Consider versioned endpoints:
- `/api/v1/organisations` - New organisation endpoints
- `/api/v1/projects` - Updated with organisation context
- Keep existing `/api/projects` for backward compatibility

---

## Security Considerations

### Data Isolation
- **Row-Level Security (RLS):** Consider PostgreSQL RLS policies
- **Middleware Enforcement:** Every request validates organisation membership
- **Query Scoping:** All database queries include organisation filter

### Permission Levels
- **Owner:** Full control, can delete organisation
- **Admin:** Can manage members, update organisation
- **Member:** Can view and use organisation resources

### Rate Limiting
- Organisation creation: 5 per hour per user
- Member invites: 50 per hour per organisation
- Standard API limits: 100 requests per minute

---

## Testing Requirements

### Unit Tests
- Organisation model validation
- Slug generation and uniqueness
- Role-based permissions

### Integration Tests
- Multi-organisation isolation
- Member management flows
- Access control enforcement

### API Tests
- All CRUD endpoints
- Error handling
- Authentication/authorization
- Pagination and filtering

---

## Performance Considerations

### Database Indexing

```sql
CREATE INDEX idx_projects_org_id ON projects(organisation_id);
CREATE INDEX idx_tasks_org_id ON tasks(organisation_id);
CREATE INDEX idx_agents_org_id ON agents(organisation_id);
CREATE INDEX idx_org_members_user_org ON organisation_members(user_id, organisation_id);
```

### Caching Strategy
- Cache user's organisation membership list
- Cache organisation context in JWT
- Invalidate on membership changes

---

## Open Questions

1. **Organisation Slugs:**
   - Should slugs be globally unique or unique per subdomain?
   - Allow custom domains?

2. **Ownership Transfer:**
   - How to handle ownership transfer process?
   - Require confirmation from new owner?

3. **Organisation Quotas:**
   - Limit number of organisations per user?
   - Limit members per organisation?

4. **Agent Access:**
   - Should agents be organisation-specific or global?
   - How to handle agent permissions across organisations?

---

## Success Criteria

- ✅ Complete data isolation between organisations
- ✅ All existing data migrated successfully
- ✅ No performance degradation
- ✅ 100% test coverage for organisation endpoints
- ✅ Documentation complete and accurate
- ✅ Backward compatibility maintained

---

## References

- Existing API patterns: `/backend/handlers/project.go`
- Authentication: `/backend/middleware/auth.go`
- Data models: `/backend/models/models.go`
- Database: `/backend/db/database.go`

---

**Document Status:** Draft - Ready for Review  
**Next Steps:** Implement models and migrations
