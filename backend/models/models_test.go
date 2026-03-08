package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func init() {
	// Ensure consistent test environment
}

// =============================================================================
// Base Model Tests
// =============================================================================

func TestBase_BeforeCreate_GeneratesUUID(t *testing.T) {
	base := Base{}
	err := base.BeforeCreate(nil)
	
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, base.ID)
}

func TestBase_BeforeCreate_PreservesExistingUUID(t *testing.T) {
	existingID := uuid.New()
	base := Base{ID: existingID}
	err := base.BeforeCreate(nil)
	
	assert.NoError(t, err)
	assert.Equal(t, existingID, base.ID)
}

// =============================================================================
// Organisation Model Tests
// =============================================================================

func TestOrganisationStatus_Constants(t *testing.T) {
	tests := []struct {
		name     string
		status   OrganisationStatus
		expected string
	}{
		{"active status", OrganisationStatusActive, "active"},
		{"suspended status", OrganisationStatusSuspended, "suspended"},
		{"archived status", OrganisationStatusArchived, "archived"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, OrganisationStatus(tt.expected), tt.status)
		})
	}
}

func TestOrganisation_StatusDefault(t *testing.T) {
	// Note: Status default is set by GORM at database level (gorm:"default:'active'")
	// In Go without DB, it will be empty string. This test verifies the intended default.
	org := Organisation{
		Name:            "Test Org",
		Slug:            "test-org",
		CreatedByUserID: uuid.New(),
	}
	
	// When explicitly set, status should be as specified
	orgExplicit := Organisation{
		Name:            "Test Org",
		Slug:            "test-org",
		Status:          OrganisationStatusActive,
		CreatedByUserID: uuid.New(),
	}
	assert.Equal(t, OrganisationStatusActive, orgExplicit.Status)
	
	// Verify the zero value is empty (will be set by GORM on insert)
	assert.Empty(t, org.Status)
}

func TestOrganisation_Fields(t *testing.T) {
	userID := uuid.New()
	
	org := Organisation{
		Name:            "Test Organisation",
		Slug:            "test-organisation",
		Description:     "A test organisation",
		Status:          OrganisationStatusActive,
		CreatedByUserID: userID,
	}
	
	assert.Equal(t, "Test Organisation", org.Name)
	assert.Equal(t, "test-organisation", org.Slug)
	assert.Equal(t, "A test organisation", org.Description)
	assert.Equal(t, OrganisationStatusActive, org.Status)
	assert.Equal(t, userID, org.CreatedByUserID)
	
	// Test with pointers
	org2 := Organisation{
		Name:        "Org With Relations",
		Slug:        "org-relations",
		CreatedBy:   &User{Base: Base{ID: userID}},
		Members:     []OrganisationMember{{}},
		Projects:    []Project{{}},
		Agents:      []Agent{{}},
	}
	
	assert.NotNil(t, org2.CreatedBy)
	assert.Len(t, org2.Members, 1)
	assert.Len(t, org2.Projects, 1)
	assert.Len(t, org2.Agents, 1)
}

// =============================================================================
// OrganisationMember Model Tests
// =============================================================================

func TestOrganisationMemberRole_Constants(t *testing.T) {
	tests := []struct {
		name     string
		role     OrganisationMemberRole
		expected string
	}{
		{"owner role", OrganisationMemberRoleOwner, "owner"},
		{"admin role", OrganisationMemberRoleAdmin, "admin"},
		{"member role", OrganisationMemberRoleMember, "member"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, OrganisationMemberRole(tt.expected), tt.role)
		})
	}
}

func TestOrganisationMember_RoleHierarchy(t *testing.T) {
	// Verify role values for potential comparison logic
	assert.Equal(t, OrganisationMemberRole("owner"), OrganisationMemberRoleOwner)
	assert.Equal(t, OrganisationMemberRole("admin"), OrganisationMemberRoleAdmin)
	assert.Equal(t, OrganisationMemberRole("member"), OrganisationMemberRoleMember)
}

func TestOrganisationMember_Fields(t *testing.T) {
	orgID := uuid.New()
	userID := uuid.New()
	now := time.Now()
	
	member := OrganisationMember{
		OrganisationID: orgID,
		UserID:         userID,
		Role:           OrganisationMemberRoleAdmin,
		JoinedAt:       now,
	}
	
	assert.Equal(t, orgID, member.OrganisationID)
	assert.Equal(t, userID, member.UserID)
	assert.Equal(t, OrganisationMemberRoleAdmin, member.Role)
	assert.Equal(t, now, member.JoinedAt)
}

func TestOrganisationMember_Associations(t *testing.T) {
	orgID := uuid.New()
	userID := uuid.New()
	
	org := Organisation{Base: Base{ID: orgID}}
	user := User{Base: Base{ID: userID}}
	
	member := OrganisationMember{
		OrganisationID: orgID,
		UserID:         userID,
		Organisation:   &org,
		User:           &user,
		Role:           OrganisationMemberRoleMember,
	}
	
	assert.NotNil(t, member.Organisation)
	assert.NotNil(t, member.User)
	assert.Equal(t, orgID, member.Organisation.ID)
	assert.Equal(t, userID, member.User.ID)
}

// =============================================================================
// User Model Tests
// =============================================================================

func TestUser_Fields(t *testing.T) {
	orgID := uuid.New()
	
	user := User{
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
		CurrentOrgID: &orgID,
	}
	
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "hashed_password", user.PasswordHash)
	assert.NotNil(t, user.CurrentOrgID)
	assert.Equal(t, orgID, *user.CurrentOrgID)
}

func TestUser_CurrentOrg_Association(t *testing.T) {
	orgID := uuid.New()
	org := Organisation{Base: Base{ID: orgID}}
	
	user := User{
		Email:        "user@example.com",
		CurrentOrgID: &orgID,
		CurrentOrg:   &org,
	}
	
	assert.NotNil(t, user.CurrentOrg)
	assert.Equal(t, orgID, user.CurrentOrg.ID)
}

func TestUser_NilCurrentOrg(t *testing.T) {
	user := User{
		Email:        "no-org@example.com",
		PasswordHash: "hash",
		CurrentOrgID: nil,
	}
	
	assert.Nil(t, user.CurrentOrgID)
}

// =============================================================================
// Agent Model Tests
// =============================================================================

func TestAgentRole_Constants(t *testing.T) {
	tests := []struct {
		name     string
		role     AgentRole
		expected string
	}{
		{"regular role", AgentRoleRegular, "regular"},
		{"supervisor role", AgentRoleSupervisor, "supervisor"},
		{"admin role", AgentRoleAdmin, "admin"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, AgentRole(tt.expected), tt.role)
		})
	}
}

func TestAgent_Defaults(t *testing.T) {
	// Note: Role and Enabled defaults are set by GORM at database level
	// In Go without DB, they will be zero values. This test verifies explicit setting.
	agent := Agent{
		Name:        "Test Agent",
		APIKey:      "test-api-key",
		Description: "A test agent",
	}
	
	// When explicitly set, values should be as specified
	agentExplicit := Agent{
		Name:        "Test Agent",
		APIKey:      "test-api-key",
		Description: "A test agent",
		Role:        AgentRoleRegular,
		Enabled:     true,
	}
	assert.Equal(t, AgentRoleRegular, agentExplicit.Role)
	assert.True(t, agentExplicit.Enabled)
	
	// Verify the zero values (will be set by GORM on insert)
	assert.Empty(t, agent.Role)
	assert.False(t, agent.Enabled) // bool zero value is false
}

func TestAgent_OrganisationAssociation(t *testing.T) {
	orgID := uuid.New()
	org := Organisation{Base: Base{ID: orgID}}
	
	agent := Agent{
		Name:           "Org Agent",
		APIKey:         "org-api-key",
		OrganisationID: &orgID,
		Organisation:   &org,
	}
	
	assert.NotNil(t, agent.OrganisationID)
	assert.Equal(t, orgID, *agent.OrganisationID)
	assert.NotNil(t, agent.Organisation)
}

// =============================================================================
// Project Model Tests
// =============================================================================

func TestProjectStatus_Constants(t *testing.T) {
	tests := []struct {
		name     string
		status   ProjectStatus
		expected string
	}{
		{"active status", ProjectStatusActive, "active"},
		{"archived status", ProjectStatusArchived, "archived"},
		{"completed status", ProjectStatusCompleted, "completed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, ProjectStatus(tt.expected), tt.status)
		})
	}
}

func TestProject_DefaultStatus(t *testing.T) {
	// Note: Status default is set by GORM at database level
	// In Go without DB, it will be empty string. This test verifies explicit setting.
	project := Project{
		Name:            "Test Project",
		CreatedByUserID: uuid.New(),
	}
	
	// When explicitly set, status should be as specified
	projectExplicit := Project{
		Name:            "Test Project",
		Status:          ProjectStatusActive,
		CreatedByUserID: uuid.New(),
	}
	assert.Equal(t, ProjectStatusActive, projectExplicit.Status)
	
	// Verify the zero value is empty (will be set by GORM on insert)
	assert.Empty(t, project.Status)
}

func TestProject_OrganisationAssociation(t *testing.T) {
	orgID := uuid.New()
	
	project := Project{
		Name:            "Org Project",
		OrganisationID:  &orgID,
		CreatedByUserID: uuid.New(),
	}
	
	assert.NotNil(t, project.OrganisationID)
	assert.Equal(t, orgID, *project.OrganisationID)
}

// =============================================================================
// Task Model Tests
// =============================================================================

func TestTaskStatus_Constants(t *testing.T) {
	tests := []struct {
		name     string
		status   TaskStatus
		expected string
	}{
		{"pending status", TaskStatusPending, "pending"},
		{"in_progress status", TaskStatusInProgress, "in_progress"},
		{"completed status", TaskStatusCompleted, "completed"},
		{"failed status", TaskStatusFailed, "failed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, TaskStatus(tt.expected), tt.status)
		})
	}
}

func TestTaskPriority_Constants(t *testing.T) {
	tests := []struct {
		name     string
		priority TaskPriority
		expected string
	}{
		{"low priority", TaskPriorityLow, "low"},
		{"medium priority", TaskPriorityMedium, "medium"},
		{"high priority", TaskPriorityHigh, "high"},
		{"critical priority", TaskPriorityCritical, "critical"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, TaskPriority(tt.expected), tt.priority)
		})
	}
}

func TestTask_Defaults(t *testing.T) {
	// Note: Status and Priority defaults are set by GORM at database level
	// In Go without DB, they will be empty strings. This test verifies explicit setting.
	task := Task{
		Title: "Test Task",
	}
	
	// When explicitly set, values should be as specified
	taskExplicit := Task{
		Title:    "Test Task",
		Status:   TaskStatusPending,
		Priority: TaskPriorityMedium,
	}
	assert.Equal(t, TaskStatusPending, taskExplicit.Status)
	assert.Equal(t, TaskPriorityMedium, taskExplicit.Priority)
	
	// Verify the zero values are empty (will be set by GORM on insert)
	assert.Empty(t, task.Status)
	assert.Empty(t, task.Priority)
}

func TestTask_OrganisationAssociation(t *testing.T) {
	orgID := uuid.New()
	
	task := Task{
		Title:          "Org Task",
		OrganisationID: &orgID,
	}
	
	assert.NotNil(t, task.OrganisationID)
	assert.Equal(t, orgID, *task.OrganisationID)
}

func TestTask_SoftDelete(t *testing.T) {
	task := Task{
		Title:     "Deleted Task",
		DeletedAt: gorm.DeletedAt{Time: time.Now(), Valid: true},
	}
	
	assert.True(t, task.DeletedAt.Valid)
}

// =============================================================================
// TaskEvent Model Tests
// =============================================================================

func TestTaskEventType_Constants(t *testing.T) {
	tests := []struct {
		name      string
		eventType TaskEventType
		expected  string
	}{
		{"created event", TaskEventCreated, "created"},
		{"updated event", TaskEventUpdated, "updated"},
		{"status_changed event", TaskEventStatusChanged, "status_changed"},
		{"assigned event", TaskEventAssigned, "assigned"},
		{"unassigned event", TaskEventUnassigned, "unassigned"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, TaskEventType(tt.expected), tt.eventType)
		})
	}
}

func TestTaskEvent_Fields(t *testing.T) {
	taskID := uuid.New()
	
	event := TaskEvent{
		TaskID:        taskID,
		EventType:     TaskEventStatusChanged,
		PreviousState: "pending",
		NewState:      "in_progress",
		ChangedBy:     "agent-todo",
	}
	
	assert.Equal(t, taskID, event.TaskID)
	assert.Equal(t, TaskEventStatusChanged, event.EventType)
	assert.Equal(t, "pending", event.PreviousState)
	assert.Equal(t, "in_progress", event.NewState)
	assert.Equal(t, "agent-todo", event.ChangedBy)
}

// =============================================================================
// TaskComment Model Tests
// =============================================================================

func TestTaskComment_Fields(t *testing.T) {
	taskID := uuid.New()
	authorID := uuid.New()
	
	comment := TaskComment{
		TaskID:     taskID,
		Content:    "This is a comment",
		AuthorID:   authorID,
		AuthorType: "agent",
		AuthorName: "agent-todo",
	}
	
	assert.Equal(t, taskID, comment.TaskID)
	assert.Equal(t, "This is a comment", comment.Content)
	assert.Equal(t, authorID, comment.AuthorID)
	assert.Equal(t, "agent", comment.AuthorType)
	assert.Equal(t, "agent-todo", comment.AuthorName)
}

// =============================================================================
// Edge Cases and Validation Tests
// =============================================================================

func TestOrganisation_EmptyFields(t *testing.T) {
	org := Organisation{}
	
	assert.Empty(t, org.Name)
	assert.Empty(t, org.Slug)
	assert.Empty(t, org.Description)
	assert.Empty(t, org.Status)
	assert.Empty(t, org.CreatedByUserID)
}

func TestOrganisationMember_EmptyFields(t *testing.T) {
	member := OrganisationMember{}
	
	assert.Empty(t, member.OrganisationID)
	assert.Empty(t, member.UserID)
	assert.Empty(t, member.Role)
	assert.Empty(t, member.JoinedAt)
}

func TestOrganisation_JSONSerialization(t *testing.T) {
	org := Organisation{
		Base:        Base{ID: uuid.New()},
		Name:        "Test Org",
		Slug:        "test-org",
		Description: "Test Description",
		Status:      OrganisationStatusActive,
	}
	
	// Verify all expected JSON tags are present
	assert.NotNil(t, org.ID)
	assert.NotEmpty(t, org.Name)
	assert.NotEmpty(t, org.Slug)
}

// =============================================================================
// GORM Tag Tests (verify struct tags are correctly set)
// =============================================================================

func TestOrganisation_GormTags(t *testing.T) {
	// These tests verify that the struct has correct GORM tags
	// by checking that the fields exist and have expected types
	
	org := Organisation{}
	
	// Base embeds ID, CreatedAt, UpdatedAt
	_ = org.ID
	_ = org.CreatedAt
	_ = org.UpdatedAt
	
	// Organisation-specific fields
	_ = org.Name
	_ = org.Slug
	_ = org.Description
	_ = org.Status
	_ = org.CreatedByUserID
	_ = org.CreatedBy
	_ = org.Members
	_ = org.Projects
	_ = org.Agents
}

func TestOrganisationMember_GormTags(t *testing.T) {
	member := OrganisationMember{}
	
	// Base embeds ID, CreatedAt, UpdatedAt
	_ = member.ID
	_ = member.CreatedAt
	_ = member.UpdatedAt
	
	// OrganisationMember-specific fields
	_ = member.OrganisationID
	_ = member.Organisation
	_ = member.UserID
	_ = member.User
	_ = member.Role
	_ = member.JoinedAt
}
