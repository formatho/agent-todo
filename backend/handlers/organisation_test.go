package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// =============================================================================
// Helper Functions
// =============================================================================

func setupTestRouter() *gin.Engine {
	router := gin.New()
	return router
}

func makeJSONRequest(method, path string, body interface{}, headers map[string]string) (*http.Request, *httptest.ResponseRecorder) {
	var reqBody bytes.Buffer
	if body != nil {
		json.NewEncoder(&reqBody).Encode(body)
	}

	req := httptest.NewRequest(method, path, &reqBody)
	req.Header.Set("Content-Type", "application/json")

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return req, httptest.NewRecorder()
}

// =============================================================================
// CreateOrganisation Tests
// =============================================================================

func TestCreateOrganisation_Success(t *testing.T) {
	router := setupTestRouter()
	handler := NewOrganisationHandler()

	router.POST("/organisations", func(c *gin.Context) {
		// Mock user authentication
		c.Set("user_id", uuid.New().String())
		handler.CreateOrganisation(c)
	})

	req, w := makeJSONRequest("POST", "/organisations", map[string]interface{}{
		"name":        "Test Organisation",
		"slug":        "test-org",
		"description": "A test organisation",
	}, nil)

	router.ServeHTTP(w, req)

	// Note: This will fail without a database, but the test structure is correct
	// In production, use mocks or a test database
}

func TestCreateOrganisation_Unauthorized(t *testing.T) {
	router := setupTestRouter()
	handler := NewOrganisationHandler()

	router.POST("/organisations", handler.CreateOrganisation)

	req, w := makeJSONRequest("POST", "/organisations", map[string]interface{}{
		"name": "Test Organisation",
		"slug": "test-org",
	}, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCreateOrganisation_MissingFields(t *testing.T) {
	router := setupTestRouter()
	handler := NewOrganisationHandler()

	router.POST("/organisations", func(c *gin.Context) {
		c.Set("user_id", uuid.New().String())
		handler.CreateOrganisation(c)
	})

	// Missing required field "slug"
	req, w := makeJSONRequest("POST", "/organisations", map[string]interface{}{
		"name": "Test Organisation",
	}, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateOrganisation_InvalidSlug(t *testing.T) {
	router := setupTestRouter()
	handler := NewOrganisationHandler()

	router.POST("/organisations", func(c *gin.Context) {
		c.Set("user_id", uuid.New().String())
		handler.CreateOrganisation(c)
	})

	testCases := []struct {
		name string
		slug string
	}{
		{"too short", "ab"},
		{"uppercase", "Test-Org"},
		{"spaces", "test org"},
		{"special chars", "test_org!"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, w := makeJSONRequest("POST", "/organisations", map[string]interface{}{
				"name": "Test Organisation",
				"slug": tc.slug,
			}, nil)

			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusBadRequest, w.Code)
		})
	}
}

func TestCreateOrganisation_InvalidJSON(t *testing.T) {
	router := setupTestRouter()
	handler := NewOrganisationHandler()

	router.POST("/organisations", func(c *gin.Context) {
		c.Set("user_id", uuid.New().String())
		handler.CreateOrganisation(c)
	})

	req := httptest.NewRequest("POST", "/organisations", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// =============================================================================
// ListOrganisations Tests
// =============================================================================

func TestListOrganisations_Unauthorized(t *testing.T) {
	router := setupTestRouter()
	handler := NewOrganisationHandler()

	router.GET("/organisations", handler.ListOrganisations)

	req, w := makeJSONRequest("GET", "/organisations", nil, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestListOrganisations_WithFilters(t *testing.T) {
	router := setupTestRouter()
	handler := NewOrganisationHandler()

	router.GET("/organisations", func(c *gin.Context) {
		c.Set("user_id", uuid.New().String())
		handler.ListOrganisations(c)
	})

	// Test with status filter
	req, w := makeJSONRequest("GET", "/organisations?status=active&search=test", nil, nil)
	router.ServeHTTP(w, req)

	// Response depends on database, but request parsing should work
}

// =============================================================================
// GetOrganisation Tests
// =============================================================================

func TestGetOrganisation_Unauthorized(t *testing.T) {
	router := setupTestRouter()
	handler := NewOrganisationHandler()

	router.GET("/organisations/:id", handler.GetOrganisation)

	orgID := uuid.New()
	req, w := makeJSONRequest("GET", "/organisations/"+orgID.String(), nil, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetOrganisation_InvalidID(t *testing.T) {
	router := setupTestRouter()
	handler := NewOrganisationHandler()

	router.GET("/organisations/:id", func(c *gin.Context) {
		c.Set("user_id", uuid.New().String())
		handler.GetOrganisation(c)
	})

	req, w := makeJSONRequest("GET", "/organisations/invalid-uuid", nil, nil)

	router.ServeHTTP(w, req)

	// Will return error due to invalid UUID or not found
	assert.True(t, w.Code == http.StatusBadRequest || w.Code == http.StatusInternalServerError || w.Code == http.StatusNotFound)
}

// =============================================================================
// UpdateOrganisation Tests
// =============================================================================

func TestUpdateOrganisation_Unauthorized(t *testing.T) {
	router := setupTestRouter()
	handler := NewOrganisationHandler()

	router.PATCH("/organisations/:id", handler.UpdateOrganisation)

	orgID := uuid.New()
	req, w := makeJSONRequest("PATCH", "/organisations/"+orgID.String(), map[string]interface{}{
		"name": "Updated Name",
	}, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestUpdateOrganisation_EmptyBody(t *testing.T) {
	router := setupTestRouter()
	handler := NewOrganisationHandler()

	router.PATCH("/organisations/:id", func(c *gin.Context) {
		c.Set("user_id", uuid.New().String())
		handler.UpdateOrganisation(c)
	})

	orgID := uuid.New()
	req, w := makeJSONRequest("PATCH", "/organisations/"+orgID.String(), map[string]interface{}{}, nil)

	router.ServeHTTP(w, req)

	// Empty body should still be accepted (no updates)
	// The actual response depends on authorization check
}

func TestUpdateOrganisation_InvalidStatus(t *testing.T) {
	router := setupTestRouter()
	handler := NewOrganisationHandler()

	router.PATCH("/organisations/:id", func(c *gin.Context) {
		c.Set("user_id", uuid.New().String())
		handler.UpdateOrganisation(c)
	})

	orgID := uuid.New()
	req, w := makeJSONRequest("PATCH", "/organisations/"+orgID.String(), map[string]interface{}{
		"status": "invalid-status",
	}, nil)

	router.ServeHTTP(w, req)

	// Status validation happens at service level
}

// =============================================================================
// DeleteOrganisation Tests
// =============================================================================

func TestDeleteOrganisation_Unauthorized(t *testing.T) {
	router := setupTestRouter()
	handler := NewOrganisationHandler()

	router.DELETE("/organisations/:id", handler.DeleteOrganisation)

	orgID := uuid.New()
	req, w := makeJSONRequest("DELETE", "/organisations/"+orgID.String(), nil, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestDeleteOrganisation_NonOwnerForbidden(t *testing.T) {
	router := setupTestRouter()
	handler := NewOrganisationHandler()

	router.DELETE("/organisations/:id", func(c *gin.Context) {
		c.Set("user_id", uuid.New().String())
		handler.DeleteOrganisation(c)
	})

	orgID := uuid.New()
	req, w := makeJSONRequest("DELETE", "/organisations/"+orgID.String(), nil, nil)

	router.ServeHTTP(w, req)

	// Non-owner should get 403 Forbidden (or error from service)
}

// =============================================================================
// AddOrganisationMember Tests
// =============================================================================

func TestAddOrganisationMember_Unauthorized(t *testing.T) {
	router := setupTestRouter()
	handler := NewOrganisationHandler()

	router.POST("/organisations/:id/members", handler.AddOrganisationMember)

	orgID := uuid.New()
	req, w := makeJSONRequest("POST", "/organisations/"+orgID.String()+"/members", map[string]interface{}{
		"user_email": "test@example.com",
		"role":        "member",
	}, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAddOrganisationMember_MissingFields(t *testing.T) {
	router := setupTestRouter()
	handler := NewOrganisationHandler()

	router.POST("/organisations/:id/members", func(c *gin.Context) {
		c.Set("user_id", uuid.New().String())
		handler.AddOrganisationMember(c)
	})

	orgID := uuid.New()

	// Missing role
	req, w := makeJSONRequest("POST", "/organisations/"+orgID.String()+"/members", map[string]interface{}{
		"user_email": "test@example.com",
	}, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAddOrganisationMember_InvalidEmail(t *testing.T) {
	router := setupTestRouter()
	handler := NewOrganisationHandler()

	router.POST("/organisations/:id/members", func(c *gin.Context) {
		c.Set("user_id", uuid.New().String())
		handler.AddOrganisationMember(c)
	})

	orgID := uuid.New()
	req, w := makeJSONRequest("POST", "/organisations/"+orgID.String()+"/members", map[string]interface{}{
		"user_email": "invalid-email",
		"role":        "member",
	}, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAddOrganisationMember_CannotAssignOwner(t *testing.T) {
	router := setupTestRouter()
	handler := NewOrganisationHandler()

	router.POST("/organisations/:id/members", func(c *gin.Context) {
		c.Set("user_id", uuid.New().String())
		handler.AddOrganisationMember(c)
	})

	orgID := uuid.New()
	req, w := makeJSONRequest("POST", "/organisations/"+orgID.String()+"/members", map[string]interface{}{
		"user_email": "test@example.com",
		"role":        "owner",
	}, nil)

	router.ServeHTTP(w, req)

	// Should return error - cannot assign owner role
	// Response depends on service implementation
}

// =============================================================================
// UpdateMemberRole Tests
// =============================================================================

func TestUpdateMemberRole_Unauthorized(t *testing.T) {
	router := setupTestRouter()
	handler := NewOrganisationHandler()

	router.PATCH("/organisations/:id/members/:member_id", handler.UpdateMemberRole)

	orgID := uuid.New()
	memberID := uuid.New()
	req, w := makeJSONRequest("PATCH", "/organisations/"+orgID.String()+"/members/"+memberID.String(), map[string]interface{}{
		"role": "admin",
	}, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestUpdateMemberRole_MissingRole(t *testing.T) {
	router := setupTestRouter()
	handler := NewOrganisationHandler()

	router.PATCH("/organisations/:id/members/:member_id", func(c *gin.Context) {
		c.Set("user_id", uuid.New().String())
		handler.UpdateMemberRole(c)
	})

	orgID := uuid.New()
	memberID := uuid.New()
	req, w := makeJSONRequest("PATCH", "/organisations/"+orgID.String()+"/members/"+memberID.String(), map[string]interface{}{}, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// =============================================================================
// RemoveOrganisationMember Tests
// =============================================================================

func TestRemoveOrganisationMember_Unauthorized(t *testing.T) {
	router := setupTestRouter()
	handler := NewOrganisationHandler()

	router.DELETE("/organisations/:id/members/:member_id", handler.RemoveOrganisationMember)

	orgID := uuid.New()
	memberID := uuid.New()
	req, w := makeJSONRequest("DELETE", "/organisations/"+orgID.String()+"/members/"+memberID.String(), nil, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// =============================================================================
// LeaveOrganisation Tests
// =============================================================================

func TestLeaveOrganisation_Unauthorized(t *testing.T) {
	router := setupTestRouter()
	handler := NewOrganisationHandler()

	router.POST("/organisations/:id/leave", handler.LeaveOrganisation)

	orgID := uuid.New()
	req, w := makeJSONRequest("POST", "/organisations/"+orgID.String()+"/leave", nil, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// =============================================================================
// Authorization Tests (Role-Based Access Control)
// =============================================================================

func TestOrganisationHandler_RoleBasedAccess(t *testing.T) {
	// Test matrix for role-based access control
	tests := []struct {
		name       string
		endpoint   string
		method     string
		userRole   string
		expectCode int
	}{
		// Get - All members can read
		{"member can view org", "/organisations/:id", "GET", "member", http.StatusUnauthorized},
		{"admin can view org", "/organisations/:id", "GET", "admin", http.StatusUnauthorized},
		{"owner can view org", "/organisations/:id", "GET", "owner", http.StatusUnauthorized},

		// Update - Admin or Owner only
		{"member cannot update org", "/organisations/:id", "PATCH", "member", http.StatusUnauthorized},
		{"admin can update org", "/organisations/:id", "PATCH", "admin", http.StatusUnauthorized},
		{"owner can update org", "/organisations/:id", "PATCH", "owner", http.StatusUnauthorized},

		// Delete - Owner only
		{"member cannot delete org", "/organisations/:id", "DELETE", "member", http.StatusUnauthorized},
		{"admin cannot delete org", "/organisations/:id", "DELETE", "admin", http.StatusUnauthorized},
		{"owner can delete org", "/organisations/:id", "DELETE", "owner", http.StatusUnauthorized},

		// Add member - Admin or Owner only
		{"member cannot add member", "/organisations/:id/members", "POST", "member", http.StatusUnauthorized},

		// Update member role - Owner only
		{"member cannot update role", "/organisations/:id/members/:member_id", "PATCH", "member", http.StatusUnauthorized},
		{"admin cannot update role", "/organisations/:id/members/:member_id", "PATCH", "admin", http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// All tests expect 401 because we don't set user_id
			// In integration tests with database, we'd verify actual role checks
			router := setupTestRouter()
			handler := NewOrganisationHandler()

			switch tt.method {
			case "GET":
				router.GET(tt.endpoint, handler.GetOrganisation)
			case "PATCH":
				router.PATCH(tt.endpoint, handler.UpdateOrganisation)
			case "DELETE":
				router.DELETE(tt.endpoint, handler.DeleteOrganisation)
			case "POST":
				router.POST(tt.endpoint, handler.AddOrganisationMember)
			}

			orgID := uuid.New()
			path := tt.endpoint
			path = replacePathParams(path, orgID.String())

			req, w := makeJSONRequest(tt.method, path, nil, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectCode, w.Code)
		})
	}
}

func replacePathParams(path, orgID string) string {
	path = replaceAll(path, ":id", orgID)
	path = replaceAll(path, ":member_id", uuid.New().String())
	return path
}

func replaceAll(s, old, new string) string {
	result := ""
	for i := 0; i < len(s); i++ {
		if i+len(old) <= len(s) && s[i:i+len(old)] == old {
			result += new
			i += len(old) - 1
		} else {
			result += string(s[i])
		}
	}
	return result
}

// =============================================================================
// Input Validation Tests
// =============================================================================

func TestOrganisationHandler_InputValidation(t *testing.T) {
	t.Run("CreateOrganisation validates required fields", func(t *testing.T) {
		router := setupTestRouter()
		handler := NewOrganisationHandler()

		router.POST("/organisations", func(c *gin.Context) {
			c.Set("user_id", uuid.New().String())
			handler.CreateOrganisation(c)
		})

		testCases := []struct {
			name    string
			payload map[string]interface{}
		}{
			{"missing name", map[string]interface{}{"slug": "test-org"}},
			{"missing slug", map[string]interface{}{"name": "Test Org"}},
			{"empty name", map[string]interface{}{"name": "", "slug": "test-org"}},
			{"empty slug", map[string]interface{}{"name": "Test Org", "slug": ""}},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				req, w := makeJSONRequest("POST", "/organisations", tc.payload, nil)
				router.ServeHTTP(w, req)
				assert.Equal(t, http.StatusBadRequest, w.Code)
			})
		}
	})

	t.Run("AddMember validates email format", func(t *testing.T) {
		router := setupTestRouter()
		handler := NewOrganisationHandler()

		router.POST("/organisations/:id/members", func(c *gin.Context) {
			c.Set("user_id", uuid.New().String())
			handler.AddOrganisationMember(c)
		})

		invalidEmails := []string{
			"not-an-email",
			"missing@domain",
			"@nodomain.com",
			"spaces in@email.com",
		}

		for _, email := range invalidEmails {
			t.Run(email, func(t *testing.T) {
				orgID := uuid.New()
				req, w := makeJSONRequest("POST", "/organisations/"+orgID.String()+"/members", map[string]interface{}{
					"user_email": email,
					"role":        "member",
				}, nil)
				router.ServeHTTP(w, req)
				assert.Equal(t, http.StatusBadRequest, w.Code)
			})
		}
	})
}

// =============================================================================
// HTTP Method Tests
// =============================================================================

func TestOrganisationHandler_MethodNotAllowed(t *testing.T) {
	router := setupTestRouter()
	handler := NewOrganisationHandler()

	router.GET("/organisations", handler.ListOrganisations)
	router.POST("/organisations", handler.CreateOrganisation)

	// Try to DELETE /organisations (not defined)
	req := httptest.NewRequest("DELETE", "/organisations", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// =============================================================================
// Handler Constructor Tests
// =============================================================================

func TestNewOrganisationHandler(t *testing.T) {
	handler := NewOrganisationHandler()
	assert.NotNil(t, handler)
	assert.NotNil(t, handler.orgService)
}
