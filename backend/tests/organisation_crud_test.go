package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/formatho/agent-todo/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestOrganisationCRUD tests organisation create, read, update, delete operations
func TestOrganisationCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	router := setupOrganisationRouter()
	testDB := db.GetDB()

	// Clean up test data
	testDB.Exec("DELETE FROM task_events")
	testDB.Exec("DELETE FROM task_comments")
	testDB.Exec("DELETE FROM tasks")
	testDB.Exec("DELETE FROM agents")
	testDB.Exec("DELETE FROM projects")
	testDB.Exec("DELETE FROM organisation_members")
	testDB.Exec("DELETE FROM organisations")
	testDB.Exec("DELETE FROM users")

	t.Run("Create organisation", func(t *testing.T) {
		// Register user
		user := registerTestUser(t, router, "org-create@example.com")
		token := user["token"].(string)

		// Create organisation
		orgPayload := map[string]string{
			"name":        "Test Organisation",
			"slug":        "test-org",
			"description": "A test organisation",
		}
		jsonPayload, _ := json.Marshal(orgPayload)

		req, _ := http.NewRequest("POST", "/organisations", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "Test Organisation", response["name"])
		assert.Equal(t, "test-org", response["slug"])
		assert.NotEmpty(t, response["id"])
	})

	t.Run("Create organisation with duplicate slug fails", func(t *testing.T) {
		// Create first org
		_ = registerTestUser(t, router, "org-dup1@example.com")
		createTestOrganisationWithUser(t, router, "dup-slug", "Org 1", "org-dup1@example.com")

		// Try to create second org with same slug
		user2 := registerTestUser(t, router, "org-dup2@example.com")
		token2 := user2["token"].(string)

		orgPayload := map[string]string{
			"name":        "Duplicate Slug Org",
			"slug":        "dup-slug", // Already exists
			"description": "Should fail",
		}
		jsonPayload, _ := json.Marshal(orgPayload)

		req, _ := http.NewRequest("POST", "/organisations", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token2))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("List organisations", func(t *testing.T) {
		// Create user and organisation
		user := registerTestUser(t, router, "org-list@example.com")
		token := user["token"].(string)

		orgPayload := map[string]string{
			"name":        "List Test Org",
			"slug":        "list-test-org",
			"description": "Test listing",
		}
		jsonPayload, _ := json.Marshal(orgPayload)

		req, _ := http.NewRequest("POST", "/organisations", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// List organisations
		req, _ = http.NewRequest("GET", "/organisations", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.GreaterOrEqual(t, len(response), 1)
	})

	t.Run("Get organisation by ID", func(t *testing.T) {
		// Create organisation
		org := createTestOrganisationWithUser(t, router, "get-org", "Get Test Org", "org-get@example.com")
		token := org["token"].(string)
		orgID := org["organisation_id"].(string)

		// Get organisation
		req, _ := http.NewRequest("GET", fmt.Sprintf("/organisations/%s", orgID), nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, orgID, response["id"])
		assert.Equal(t, "Get Test Org", response["name"])
	})

	t.Run("Update organisation", func(t *testing.T) {
		// Create organisation
		org := createTestOrganisationWithUser(t, router, "update-org", "Update Test Org", "org-update@example.com")
		token := org["token"].(string)
		orgID := org["organisation_id"].(string)

		// Update organisation
		updatePayload := map[string]string{
			"name":        "Updated Organisation",
			"description": "Updated description",
		}
		jsonPayload, _ := json.Marshal(updatePayload)

		req, _ := http.NewRequest("PATCH", fmt.Sprintf("/organisations/%s", orgID), bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "Updated Organisation", response["name"])
		assert.Equal(t, "Updated description", response["description"])
	})

	t.Run("Delete organisation", func(t *testing.T) {
		// Create organisation
		org := createTestOrganisationWithUser(t, router, "delete-org", "Delete Test Org", "org-delete@example.com")
		token := org["token"].(string)
		orgID := org["organisation_id"].(string)

		// Delete organisation
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("/organisations/%s", orgID), nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify deleted
		req, _ = http.NewRequest("GET", fmt.Sprintf("/organisations/%s", orgID), nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Add member to organisation", func(t *testing.T) {
		// Create organisation with owner
		org := createTestOrganisationWithUser(t, router, "add-member-org", "Add Member Org", "owner-add@example.com")
		token := org["token"].(string)
		orgID := org["organisation_id"].(string)

		// Create another user to add as member
		memberUser := registerTestUser(t, router, "new-member@example.com")
		memberID := memberUser["user"].(map[string]interface{})["id"].(string)

		// Add member
		memberPayload := map[string]string{
			"user_id": memberID,
			"role":    "member",
		}
		jsonPayload, _ := json.Marshal(memberPayload)

		req, _ := http.NewRequest("POST", fmt.Sprintf("/organisations/%s/members", orgID), bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Update member role", func(t *testing.T) {
		// Create organisation with owner
		org := createTestOrganisationWithUser(t, router, "role-update-org", "Role Update Org", "owner-role@example.com")
		token := org["token"].(string)
		orgID := org["organisation_id"].(string)

		// Create member
		memberUser := registerTestUser(t, router, "member-role@example.com")
		memberID := memberUser["user"].(map[string]interface{})["id"].(string)

		// Add member
		memberPayload := map[string]string{
			"user_id": memberID,
			"role":    "member",
		}
		jsonPayload, _ := json.Marshal(memberPayload)

		req, _ := http.NewRequest("POST", fmt.Sprintf("/organisations/%s/members", orgID), bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Get member ID from response
		var addMemberResponse map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &addMemberResponse)
		memberRecordID := addMemberResponse["id"].(string)

		// Update role to admin
		rolePayload := map[string]string{
			"role": "admin",
		}
		jsonPayload, _ = json.Marshal(rolePayload)

		req, _ = http.NewRequest("PATCH", fmt.Sprintf("/organisations/%s/members/%s", orgID, memberRecordID), bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Remove member from organisation", func(t *testing.T) {
		// Create organisation with owner
		org := createTestOrganisationWithUser(t, router, "remove-member-org", "Remove Member Org", "owner-remove@example.com")
		token := org["token"].(string)
		orgID := org["organisation_id"].(string)

		// Create and add member
		memberUser := registerTestUser(t, router, "member-remove@example.com")
		memberID := memberUser["user"].(map[string]interface{})["id"].(string)

		memberPayload := map[string]string{
			"user_id": memberID,
			"role":    "member",
		}
		jsonPayload, _ := json.Marshal(memberPayload)

		req, _ := http.NewRequest("POST", fmt.Sprintf("/organisations/%s/members", orgID), bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var addMemberResponse map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &addMemberResponse)
		memberRecordID := addMemberResponse["id"].(string)

		// Remove member
		req, _ = http.NewRequest("DELETE", fmt.Sprintf("/organisations/%s/members/%s", orgID, memberRecordID), nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Non-owner cannot add members", func(t *testing.T) {
		// Create organisation with owner
		org := createTestOrganisationWithUser(t, router, "auth-test-org", "Auth Test Org", "owner-auth@example.com")
		orgID := org["organisation_id"].(string)

		// Create another user (not a member)
		nonMemberUser := registerTestUser(t, router, "non-member-auth@example.com")
		nonMemberToken := nonMemberUser["token"].(string)

		// Create a third user to try to add
		thirdUser := registerTestUser(t, router, "third-user-auth@example.com")
		thirdUserID := thirdUser["user"].(map[string]interface{})["id"].(string)

		// Non-member tries to add member
		memberPayload := map[string]string{
			"user_id": thirdUserID,
			"role":    "member",
		}
		jsonPayload, _ := json.Marshal(memberPayload)

		req, _ := http.NewRequest("POST", fmt.Sprintf("/organisations/%s/members", orgID), bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", nonMemberToken))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}
