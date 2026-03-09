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

// ProjectInput represents input for creating/updating a project
type ProjectInput struct {
	Name             *string
	Description      *string
	RepositoryURL    *string
	DeployedURL      *string
	DocumentationURL *string
	LLMContext       *string
	Status           *models.ProjectStatus
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

// CreateProject creates a new project with all fields
func (s *ProjectService) CreateProject(input ProjectInput, createdByUserID string, organisationID *string) (*models.Project, error) {
	project := &models.Project{
		Status:          models.ProjectStatusActive,
		CreatedByUserID: uuid.MustParse(createdByUserID),
	}

	if input.Name != nil {
		project.Name = *input.Name
	}
	if input.Description != nil {
		project.Description = *input.Description
	}
	if input.RepositoryURL != nil {
		project.RepositoryURL = *input.RepositoryURL
	}
	if input.DeployedURL != nil {
		project.DeployedURL = *input.DeployedURL
	}
	if input.DocumentationURL != nil {
		project.DocumentationURL = *input.DocumentationURL
	}
	if input.LLMContext != nil {
		project.LLMContext = *input.LLMContext
	}
	if organisationID != nil {
		orgUUID := uuid.MustParse(*organisationID)
		project.OrganisationID = &orgUUID
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

// GetByIDAndOrganisation retrieves a project by ID, ensuring it belongs to the organisation
func (s *ProjectService) GetByIDAndOrganisation(id, organisationID string) (*models.Project, error) {
	var project models.Project
	err := s.db.Preload("CreatedBy").Preload("Tasks").
		Where("id = ? AND organisation_id = ?", id, organisationID).First(&project).Error
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

// UpdateProject updates a project with all fields
func (s *ProjectService) UpdateProject(id string, input ProjectInput) (*models.Project, error) {
	var project models.Project
	if err := s.db.Where("id = ?", id).First(&project).Error; err != nil {
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
	if input.RepositoryURL != nil {
		updates["repository_url"] = *input.RepositoryURL
	}
	if input.DeployedURL != nil {
		updates["deployed_url"] = *input.DeployedURL
	}
	if input.DocumentationURL != nil {
		updates["documentation_url"] = *input.DocumentationURL
	}
	if input.LLMContext != nil {
		updates["llm_context"] = *input.LLMContext
	}

	if err := s.db.Model(&project).Updates(updates).Error; err != nil {
		return nil, err
	}

	return s.GetByID(id)
}

// UpdateByOrganisation updates a project, verifying it belongs to the organisation
func (s *ProjectService) UpdateByOrganisation(id, organisationID string, name, description *string, status *models.ProjectStatus) (*models.Project, error) {
	var project models.Project
	if err := s.db.Where("id = ? AND organisation_id = ?", id, organisationID).First(&project).Error; err != nil {
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

// UpdateProjectByOrganisation updates a project with all fields, verifying organisation
func (s *ProjectService) UpdateProjectByOrganisation(id, organisationID string, input ProjectInput) (*models.Project, error) {
	var project models.Project
	if err := s.db.Where("id = ? AND organisation_id = ?", id, organisationID).First(&project).Error; err != nil {
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
	if input.RepositoryURL != nil {
		updates["repository_url"] = *input.RepositoryURL
	}
	if input.DeployedURL != nil {
		updates["deployed_url"] = *input.DeployedURL
	}
	if input.DocumentationURL != nil {
		updates["documentation_url"] = *input.DocumentationURL
	}
	if input.LLMContext != nil {
		updates["llm_context"] = *input.LLMContext
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

// DeleteByOrganisation deletes a project, verifying it belongs to the organisation
func (s *ProjectService) DeleteByOrganisation(id, organisationID string) error {
	result := s.db.Where("id = ? AND organisation_id = ?", id, organisationID).Delete(&models.Project{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
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

// GetTasksByOrganisation retrieves tasks for a project, verifying organisation
func (s *ProjectService) GetTasksByOrganisation(projectID, organisationID string) ([]models.Task, error) {
	var project models.Project
	if err := s.db.Where("id = ? AND organisation_id = ?", projectID, organisationID).First(&project).Error; err != nil {
		return nil, err
	}

	var tasks []models.Task
	err := s.db.Preload("CreatedBy").Preload("AssignedAgent").
		Where("project_id = ? AND organisation_id = ?", projectID, organisationID).
		Order("created_at DESC").
		Find(&tasks).Error

	if err != nil {
		return nil, err
	}

	return tasks, nil
}
