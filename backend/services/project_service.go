package services

import (
	"github.com/formatho/agent-todo/db"
	"github.com/formatho/agent-todo/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectService struct {
	db *gorm.DB
}

func NewProjectService() *ProjectService {
	return &ProjectService{
		db: db.GetDB(),
	}
}

// ProjectFilter represents filters for listing projects
type ProjectFilter struct {
	Status          models.ProjectStatus `json:"status"`
	CreatedByUserID string               `json:"created_by_user_id"`
	SearchTerm      string               `json:"search_term"`
	OrganisationID  string               `json:"organisation_id"` // Optional: filter by organisation
}

// Create creates a new project
func (s *ProjectService) Create(name, description string, createdByUserID string) (*models.Project, error) {
	project := &models.Project{
		Name:            name,
		Description:     description,
		Status:          models.ProjectStatusActive,
		CreatedByUserID: uuid.MustParse(createdByUserID),
	}

	if err := s.db.Create(project).Error; err != nil {
		return nil, err
	}

	return project, nil
}

// CreateWithOrganisation creates a new project with organisation context
func (s *ProjectService) CreateWithOrganisation(name, description, createdByUserID, organisationID string) (*models.Project, error) {
	orgUUID := uuid.MustParse(organisationID)
	project := &models.Project{
		Name:            name,
		Description:     description,
		Status:          models.ProjectStatusActive,
		CreatedByUserID: uuid.MustParse(createdByUserID),
		OrganisationID:  &orgUUID,
	}

	if err := s.db.Create(project).Error; err != nil {
		return nil, err
	}

	return project, nil
}

// GetByID retrieves a project with relationships
func (s *ProjectService) GetByID(id string) (*models.Project, error) {
	var project models.Project
	err := s.db.Preload("CreatedBy").Preload("Tasks").Where("id = ?", id).First(&project).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

// List retrieves projects with filters
func (s *ProjectService) List(filter ProjectFilter) ([]models.Project, error) {
	var projects []models.Project
	query := s.db.Preload("CreatedBy")

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	if filter.CreatedByUserID != "" {
		query = query.Where("created_by_user_id = ?", filter.CreatedByUserID)
	}

	if filter.SearchTerm != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?",
			"%"+filter.SearchTerm+"%", "%"+filter.SearchTerm+"%")
	}

	// Filter by organisation if provided
	if filter.OrganisationID != "" {
		query = query.Where("organisation_id = ?", filter.OrganisationID)
	}

	err := query.Order("created_at DESC").Find(&projects).Error
	if err != nil {
		return nil, err
	}

	return projects, nil
}

// Update updates a project
func (s *ProjectService) Update(id string, name, description *string, status *models.ProjectStatus) (*models.Project, error) {
	var project models.Project
	if err := s.db.Where("id = ?", id).First(&project).Error; err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})

	if name != nil {
		updates["name"] = *name
	}
	if description != nil {
		updates["description"] = *description
	}
	if status != nil {
		updates["status"] = *status
	}

	if err := s.db.Model(&project).Updates(updates).Error; err != nil {
		return nil, err
	}

	return s.GetByID(id)
}

// Delete deletes a project
func (s *ProjectService) Delete(id string) error {
	return s.db.Where("id = ?", id).Delete(&models.Project{}).Error
}

// GetTasks retrieves tasks for a project
func (s *ProjectService) GetTasks(projectID string) ([]models.Task, error) {
	var project models.Project
	if err := s.db.Where("id = ?", projectID).First(&project).Error; err != nil {
		return nil, err
	}

	var tasks []models.Task
	err := s.db.Preload("CreatedBy").Preload("AssignedAgent").
		Where("project_id = ?", projectID).
		Order("created_at DESC").
		Find(&tasks).Error

	if err != nil {
		return nil, err
	}

	return tasks, nil
}
