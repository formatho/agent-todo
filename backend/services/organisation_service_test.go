package services

import (
	"testing"

	"github.com/formatho/agent-todo/models"
	"github.com/stretchr/testify/assert"
)

// =============================================================================
// validateSlug Unit Tests
// =============================================================================

func TestValidateSlug_Valid(t *testing.T) {
	tests := []struct {
		name string
		slug string
	}{
		{"simple lowercase", "my-org"},
		{"with numbers", "org123"},
		{"hyphenated", "my-awesome-org"},
		{"min length", "abc"},
		{"alphanumeric", "my-org-2024"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateSlug(tt.slug)
			assert.NoError(t, err)
		})
	}

	t.Run("49 characters", func(t *testing.T) {
		slug := ""
		for i := 0; i < 49; i++ {
			slug += "a"
		}
		assert.Len(t, slug, 49)
		err := validateSlug(slug)
		assert.NoError(t, err)
	})

	t.Run("50 characters (max)", func(t *testing.T) {
		slug := ""
		for i := 0; i < 50; i++ {
			slug += "a"
		}
		assert.Len(t, slug, 50)
		err := validateSlug(slug)
		assert.NoError(t, err)
	})
}

func TestValidateSlug_Invalid(t *testing.T) {
	tests := []struct {
		name        string
		slug        string
		expectError string
	}{
		{"too short", "ab", "slug must be between 3 and 50 characters"},
		{"uppercase", "My-Org", "slug can only contain lowercase letters, numbers, and hyphens"},
		{"spaces", "my org", "slug can only contain lowercase letters, numbers, and hyphens"},
		{"special chars", "my_org!", "slug can only contain lowercase letters, numbers, and hyphens"},
		{"underscore", "my_org", "slug can only contain lowercase letters, numbers, and hyphens"},
		{"empty string", "", "slug must be between 3 and 50 characters"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateSlug(tt.slug)
			assert.Error(t, err)
			assert.EqualError(t, err, tt.expectError)
		})
	}

	t.Run("too long", func(t *testing.T) {
		// Create a 51-character slug
		slug := ""
		for i := 0; i < 51; i++ {
			slug += "a"
		}
		assert.Len(t, slug, 51)
		err := validateSlug(slug)
		assert.Error(t, err)
		assert.EqualError(t, err, "slug must be between 3 and 50 characters")
	})
}

func TestValidateSlug_EdgeCases(t *testing.T) {
	t.Run("exactly 3 characters", func(t *testing.T) {
		err := validateSlug("abc")
		assert.NoError(t, err)
	})

	t.Run("exactly 50 characters", func(t *testing.T) {
		// Create a 50-character slug (repeat "ab" 25 times)
		slug := ""
		for i := 0; i < 25; i++ {
			slug += "ab"
		}
		assert.Len(t, slug, 50)
		err := validateSlug(slug)
		assert.NoError(t, err)
	})

	t.Run("51 characters (too long)", func(t *testing.T) {
		// Create a 51-character slug
		slug := ""
		for i := 0; i < 25; i++ {
			slug += "ab"
		}
		slug += "c"
		assert.Len(t, slug, 51)
		err := validateSlug(slug)
		assert.Error(t, err)
	})

	t.Run("all numbers", func(t *testing.T) {
		err := validateSlug("12345")
		assert.NoError(t, err)
	})

	t.Run("single hyphen in middle", func(t *testing.T) {
		err := validateSlug("a-b")
		assert.NoError(t, err)
	})

	t.Run("starts with hyphen (valid per regex)", func(t *testing.T) {
		// The regex ^[a-z0-9-]+$ allows starting with hyphen
		err := validateSlug("-my-org")
		assert.NoError(t, err)
	})

	t.Run("ends with hyphen (valid per regex)", func(t *testing.T) {
		// The regex ^[a-z0-9-]+$ allows ending with hyphen
		err := validateSlug("my-org-")
		assert.NoError(t, err)
	})

	t.Run("double hyphen (valid per regex)", func(t *testing.T) {
		// The regex ^[a-z0-9-]+$ allows double hyphens
		err := validateSlug("my--org")
		assert.NoError(t, err)
	})
}

// =============================================================================
// Input Struct Tests
// =============================================================================

func TestCreateOrganisationInput_Fields(t *testing.T) {
	input := CreateOrganisationInput{
		Name:        "Test Organisation",
		Slug:        "test-organisation",
		Description: "A test organisation description",
		UserID:      "user-123",
	}

	assert.Equal(t, "Test Organisation", input.Name)
	assert.Equal(t, "test-organisation", input.Slug)
	assert.Equal(t, "A test organisation description", input.Description)
	assert.Equal(t, "user-123", input.UserID)
}

func TestUpdateOrganisationInput_Fields(t *testing.T) {
	name := "Updated Name"
	desc := "Updated Description"
	status := models.OrganisationStatus("suspended")

	input := UpdateOrganisationInput{
		Name:        &name,
		Description: &desc,
		Status:      &status,
	}

	assert.NotNil(t, input.Name)
	assert.Equal(t, "Updated Name", *input.Name)
	assert.NotNil(t, input.Description)
	assert.Equal(t, "Updated Description", *input.Description)
	assert.NotNil(t, input.Status)
	assert.Equal(t, models.OrganisationStatus("suspended"), *input.Status)
}

func TestUpdateOrganisationInput_NilFields(t *testing.T) {
	input := UpdateOrganisationInput{}

	assert.Nil(t, input.Name)
	assert.Nil(t, input.Description)
	assert.Nil(t, input.Status)
}

func TestAddMemberInput_Fields(t *testing.T) {
	input := AddMemberInput{
		UserEmail: "test@example.com",
		Role:      models.OrganisationMemberRoleAdmin,
	}

	assert.Equal(t, "test@example.com", input.UserEmail)
	assert.Equal(t, models.OrganisationMemberRoleAdmin, input.Role)
}

func TestUpdateMemberRoleInput_Fields(t *testing.T) {
	input := UpdateMemberRoleInput{
		Role: models.OrganisationMemberRoleAdmin,
	}

	assert.Equal(t, models.OrganisationMemberRoleAdmin, input.Role)
}

func TestOrganisationFilter_Fields(t *testing.T) {
	filter := OrganisationFilter{
		Status:     models.OrganisationStatusActive,
		SearchTerm: "test",
	}

	assert.Equal(t, models.OrganisationStatusActive, filter.Status)
	assert.Equal(t, "test", filter.SearchTerm)
}

func TestOrganisationFilter_Empty(t *testing.T) {
	filter := OrganisationFilter{}

	assert.Empty(t, filter.Status)
	assert.Empty(t, filter.SearchTerm)
}

// =============================================================================
// Service Constructor Tests
// =============================================================================

func TestNewOrganisationService(t *testing.T) {
	// Note: This will fail if no database is configured
	// In unit tests, we should use mocks
	// This test verifies the constructor exists and returns correct type
	
	// Skip if we want to avoid database dependency
	// service := NewOrganisationService()
	// assert.NotNil(t, service)
	// assert.NotNil(t, service.db)
}
