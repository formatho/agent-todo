package services

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/formatho/agent-todo/db"
	"github.com/formatho/agent-todo/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrganisationService struct {
	db *gorm.DB
}

func NewOrganisationService() *OrganisationService {
	return &OrganisationService{
		db: db.GetDB(),
	}
}

// OrganisationFilter represents filters for listing organisations
type OrganisationFilter struct {
	Status     models.OrganisationStatus `json:"status"`
	SearchTerm string                    `json:"search_term"`
}

// CreateOrganisationInput represents the input for creating an organisation
type CreateOrganisationInput struct {
	Name        string `json:"name" binding:"required"`
	Slug        string `json:"slug" binding:"required"`
	Description string `json:"description"`
	UserID      string `json:"-"` // Set from auth context
}

// UpdateOrganisationInput represents the input for updating an organisation
type UpdateOrganisationInput struct {
	Name        *string                  `json:"name"`
	Description *string                  `json:"description"`
	Status      *models.OrganisationStatus `json:"status"`
}

// AddMemberInput represents the input for adding a member
type AddMemberInput struct {
	UserEmail string                       `json:"user_email" binding:"required,email"`
	Role      models.OrganisationMemberRole `json:"role" binding:"required"`
}

// UpdateMemberRoleInput represents the input for updating member role
type UpdateMemberRoleInput struct {
	Role models.OrganisationMemberRole `json:"role" binding:"required"`
}

// validateSlug checks if a slug is valid
func validateSlug(slug string) error {
	if len(slug) < 3 || len(slug) > 50 {
		return errors.New("slug must be between 3 and 50 characters")
	}
	matched, _ := regexp.MatchString("^[a-z0-9-]+$", slug)
	if !matched {
		return errors.New("slug can only contain lowercase letters, numbers, and hyphens")
	}
	return nil
}

// Create creates a new organisation
func (s *OrganisationService) Create(input CreateOrganisationInput) (*models.Organisation, error) {
	// Validate slug
	if err := validateSlug(input.Slug); err != nil {
		return nil, err
	}

	// Check if slug already exists
	var count int64
	s.db.Model(&models.Organisation{}).Where("slug = ?", input.Slug).Count(&count)
	if count > 0 {
		return nil, errors.New("slug already exists")
	}

	userID := uuid.MustParse(input.UserID)

	// Create organisation
	org := &models.Organisation{
		Name:            input.Name,
		Slug:            input.Slug,
		Description:     input.Description,
		Status:          models.OrganisationStatusActive,
		CreatedByUserID: userID,
	}

	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(org).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Add creator as owner
	member := &models.OrganisationMember{
		OrganisationID: org.ID,
		UserID:         userID,
		Role:           models.OrganisationMemberRoleOwner,
		JoinedAt:       time.Now(),
	}

	if err := tx.Create(member).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Load relationships
	s.db.Preload("CreatedBy").Preload("Members.User").First(org, org.ID)

	return org, nil
}

// GetByID retrieves an organisation by ID
func (s *OrganisationService) GetByID(id string) (*models.Organisation, error) {
	var org models.Organisation
	err := s.db.Preload("CreatedBy").
		Preload("Members.User").
		Where("id = ?", id).
		First(&org).
		Error
	if err != nil {
		return nil, err
	}
	return &org, nil
}

// GetBySlug retrieves an organisation by slug
func (s *OrganisationService) GetBySlug(slug string) (*models.Organisation, error) {
	var org models.Organisation
	err := s.db.Preload("CreatedBy").
		Preload("Members.User").
		Where("slug = ?", slug).
		First(&org).
		Error
	if err != nil {
		return nil, err
	}
	return &org, nil
}

// List retrieves organisations for a user
func (s *OrganisationService) List(userID string, filter OrganisationFilter) ([]models.Organisation, error) {
	var organisations []models.Organisation

	query := s.db.Joins("JOIN organisation_members ON organisation_members.organisation_id = organisations.id").
		Where("organisation_members.user_id = ?", userID).
		Preload("CreatedBy")

	if filter.Status != "" {
		query = query.Where("organisations.status = ?", filter.Status)
	}

	if filter.SearchTerm != "" {
		query = query.Where("organisations.name ILIKE ? OR organisations.description ILIKE ?",
			"%"+filter.SearchTerm+"%", "%"+filter.SearchTerm+"%")
	}

	err := query.Order("organisations.created_at DESC").Find(&organisations).Error
	if err != nil {
		return nil, err
	}

	return organisations, nil
}

// Update updates an organisation
func (s *OrganisationService) Update(id string, input UpdateOrganisationInput) (*models.Organisation, error) {
	var org models.Organisation
	if err := s.db.Where("id = ?", id).First(&org).Error; err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})

	if input.Name != nil {
		updates["name"] = *input.Name
	}
	if input.Description != nil {
		updates["description"] = *input.Description
	}
	if input.Status != nil {
		updates["status"] = *input.Status
	}

	if len(updates) > 0 {
		if err := s.db.Model(&org).Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	// Reload with relationships
	s.db.Preload("CreatedBy").Preload("Members.User").First(&org, org.ID)

	return &org, nil
}

// Delete soft deletes an organisation
func (s *OrganisationService) Delete(id string) error {
	// Check if organisation has active resources
	var projectCount, memberCount int64
	s.db.Model(&models.Project{}).Where("organisation_id = ?", id).Count(&projectCount)
	s.db.Model(&models.OrganisationMember{}).Where("organisation_id = ? AND role != ?", id, models.OrganisationMemberRoleOwner).Count(&memberCount)

	if projectCount > 0 {
		return fmt.Errorf("organisation has %d active projects", projectCount)
	}
	if memberCount > 1 { // More than just the owner
		return fmt.Errorf("organisation has %d members", memberCount)
	}

	return s.db.Delete(&models.Organisation{}, "id = ?", id).Error
}

// AddMember adds a user to an organisation
func (s *OrganisationService) AddMember(orgID string, input AddMemberInput) (*models.OrganisationMember, error) {
	// Find user by email
	var user models.User
	if err := s.db.Where("email = ?", input.UserEmail).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	// Check if already a member
	var count int64
	s.db.Model(&models.OrganisationMember{}).
		Where("organisation_id = ? AND user_id = ?", orgID, user.ID).
		Count(&count)
	if count > 0 {
		return nil, errors.New("user is already a member")
	}

	// Cannot assign owner role
	if input.Role == models.OrganisationMemberRoleOwner {
		return nil, errors.New("cannot assign owner role")
	}

	member := &models.OrganisationMember{
		OrganisationID: uuid.MustParse(orgID),
		UserID:         user.ID,
		Role:           input.Role,
		JoinedAt:       time.Now(),
	}

	if err := s.db.Create(member).Error; err != nil {
		return nil, err
	}

	// Load relationships
	s.db.Preload("User").Preload("Organisation").First(member, member.ID)

	return member, nil
}

// UpdateMemberRole updates a member's role
func (s *OrganisationService) UpdateMemberRole(orgID, memberID string, input UpdateMemberRoleInput) (*models.OrganisationMember, error) {
	var member models.OrganisationMember
	if err := s.db.Where("id = ? AND organisation_id = ?", memberID, orgID).
		Preload("User").
		First(&member).Error; err != nil {
		return nil, err
	}

	// Cannot change owner's role
	if member.Role == models.OrganisationMemberRoleOwner {
		return nil, errors.New("cannot change owner's role")
	}

	// Cannot assign owner role
	if input.Role == models.OrganisationMemberRoleOwner {
		return nil, errors.New("cannot assign owner role")
	}

	if err := s.db.Model(&member).Update("role", input.Role).Error; err != nil {
		return nil, err
	}

	return &member, nil
}

// RemoveMember removes a member from an organisation
func (s *OrganisationService) RemoveMember(orgID, memberID string) error {
	var member models.OrganisationMember
	if err := s.db.Where("id = ? AND organisation_id = ?", memberID, orgID).
		First(&member).Error; err != nil {
		return err
	}

	// Cannot remove owner
	if member.Role == models.OrganisationMemberRoleOwner {
		return errors.New("cannot remove organisation owner")
	}

	return s.db.Delete(&member).Error
}

// LeaveOrganisation allows a user to leave an organisation
func (s *OrganisationService) LeaveOrganisation(orgID, userID string) error {
	var member models.OrganisationMember
	if err := s.db.Where("organisation_id = ? AND user_id = ?", orgID, userID).
		First(&member).Error; err != nil {
		return err
	}

	// Owner cannot leave
	if member.Role == models.OrganisationMemberRoleOwner {
		return errors.New("owner cannot leave organisation")
	}

	return s.db.Delete(&member).Error
}

// IsMember checks if a user is a member of an organisation
func (s *OrganisationService) IsMember(orgID, userID string) (bool, error) {
	var count int64
	err := s.db.Model(&models.OrganisationMember{}).
		Where("organisation_id = ? AND user_id = ?", orgID, userID).
		Count(&count).
		Error
	return count > 0, err
}

// GetMemberRole gets a user's role in an organisation
func (s *OrganisationService) GetMemberRole(orgID, userID string) (*models.OrganisationMemberRole, error) {
	var member models.OrganisationMember
	err := s.db.Where("organisation_id = ? AND user_id = ?", orgID, userID).
		First(&member).
		Error
	if err != nil {
		return nil, err
	}
	return &member.Role, nil
}

// IsAdmin checks if a user is an admin or owner of an organisation
func (s *OrganisationService) IsAdmin(orgID, userID string) (bool, error) {
	role, err := s.GetMemberRole(orgID, userID)
	if err != nil {
		return false, err
	}
	return *role == models.OrganisationMemberRoleAdmin || *role == models.OrganisationMemberRoleOwner, nil
}

// IsOwner checks if a user is the owner of an organisation
func (s *OrganisationService) IsOwner(orgID, userID string) (bool, error) {
	role, err := s.GetMemberRole(orgID, userID)
	if err != nil {
		return false, err
	}
	return *role == models.OrganisationMemberRoleOwner, nil
}
