package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/formatho/agent-todo/db"
	"github.com/formatho/agent-todo/handlers"
	"github.com/formatho/agent-todo/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestDB() {
	// Use test database
	databaseURL := "postgres://agent_todo:agent_todo_pass@localhost:5432/agent_todo_test?sslmode=disable"
	db.Connect(databaseURL)
}

func setupRouter() *gin.Engine {
	gin.Set(gin.TestMode)
	router := gin.Default()

	// Setup routes
	authHandler := handlers.NewAuthHandler()
	agentHandler := handlers.NewAgentHandler()
	taskHandler := handlers.NewTaskHandler()
	agentTaskHandler := handlers.NewAgentTaskHandler()

	router.POST("/auth/register", authHandler.Register)
	router.POST("/auth/login", authHandler.Login)
	router.POST("/agents", agentHandler.CreateAgent)
	router.GET("/agents", agentHandler.ListAgents)
	router.POST("/tasks", taskHandler.CreateTask)
	router.GET("/tasks", taskHandler.ListTasks)
	router.GET("/tasks/:id", taskHandler.GetTask)
	router.POST("/agent/tasks", agentTaskHandler.CreateTask)
	router.GET("/agent/tasks", agentTaskHandler.ListTasks)
	router.PATCH("/agent/tasks/:id/status", agentTaskHandler.UpdateStatus)

	return router
}

func TestUserRegistrationAndLogin(t *testing.T) {
	router := setupRouter()

	// Test registration
	registerPayload := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	jsonPayload, _ := json.Marshal(registerPayload)

	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response, "token")
	assert.Contains(t, response, "user")

	// Test login
	loginPayload := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	jsonPayload, _ = json.Marshal(loginPayload)

	req, _ = http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var loginResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	assert.Contains(t, loginResponse, "token")
}

func TestAgentCreationAndTaskManagement(t *testing.T) {
	router := setupRouter()

	// Register and login user
	registerPayload := map[string]string{
		"email":    "agentuser@example.com",
		"password": "password123",
	}
	jsonPayload, _ := json.Marshal(registerPayload)

	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var registerResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &registerResponse)
	token := registerResponse["token"].(string)

	// Create agent
	agentPayload := map[string]string{
		"name":        "Test Agent",
		"description": "An agent for testing",
	}
	jsonPayload, _ = json.Marshal(agentPayload)

	req, _ = http.NewRequest("POST", "/agents", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var agentResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &agentResponse)
	agentID := agentResponse["id"].(string)
	apiKey := agentResponse["api_key"].(string)
	assert.NotEmpty(t, apiKey)

	// Agent creates a task
	agentTaskPayload := map[string]interface{}{
		"title":       "Agent Task",
		"description": "Task created by agent",
		"priority":    "high",
	}
	jsonPayload, _ = json.Marshal(agentTaskPayload)

	req, _ = http.NewRequest("POST", "/agent/tasks", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", apiKey)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var taskResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &taskResponse)
	taskID := taskResponse["id"].(string)
	assert.Equal(t, "Agent Task", taskResponse["title"])
	assert.Equal(t, "pending", taskResponse["status"])

	// Agent updates task status
	updatePayload := map[string]string{
		"status": "in_progress",
	}
	jsonPayload, _ = json.Marshal(updatePayload)

	req, _ = http.NewRequest("PATCH", fmt.Sprintf("/agent/tasks/%s/status", taskID), bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", apiKey)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var updatedTaskResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &updatedTaskResponse)
	assert.Equal(t, "in_progress", updatedTaskResponse["status"])

	// List agent tasks
	req, _ = http.NewRequest("GET", "/agent/tasks", nil)
	req.Header.Set("X-API-KEY", apiKey)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var tasksResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &tasksResponse)
	tasks := tasksResponse.([]interface{})
	assert.Greater(t, len(tasks), 0)
}

func TestHumanEditsAgentTask(t *testing.T) {
	router := setupRouter()

	// Register user
	registerPayload := map[string]string{
		"email":    "humanuser@example.com",
		"password": "password123",
	}
	jsonPayload, _ := json.Marshal(registerPayload)

	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var registerResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &registerResponse)
	token := registerResponse["token"].(string)

	// Create agent
	agentPayload := map[string]string{
		"name":        "Test Agent 2",
		"description": "Another test agent",
	}
	jsonPayload, _ = json.Marshal(agentPayload)

	req, _ = http.NewRequest("POST", "/agents", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var agentResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &agentResponse)
	apiKey := agentResponse["api_key"].(string)

	// Agent creates task
	agentTaskPayload := map[string]interface{}{
		"title":       "Original Title",
		"description": "Original description",
		"priority":    "medium",
	}
	jsonPayload, _ = json.Marshal(agentTaskPayload)

	req, _ = http.NewRequest("POST", "/agent/tasks", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", apiKey)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var taskResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &taskResponse)
	taskID := taskResponse["id"].(string)

	// Human edits the task
	newTitle := "Updated by Human"
	newDescription := "Updated description"
	updatePayload := map[string]*string{
		"title":       &newTitle,
		"description": &newDescription,
	}
	jsonPayload, _ = json.Marshal(updatePayload)

	req, _ = http.NewRequest("PATCH", fmt.Sprintf("/tasks/%s", taskID), bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var updatedTaskResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &updatedTaskResponse)
	assert.Equal(t, "Updated by Human", updatedTaskResponse["title"])
	assert.Equal(t, "Updated description", updatedTaskResponse["description"])

	// Agent can still see the updated task
	req, _ = http.NewRequest("GET", fmt.Sprintf("/tasks/%s", taskID), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var fetchedTaskResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &fetchedTaskResponse)
	assert.Equal(t, "Updated by Human", fetchedTaskResponse["title"])
}

func TestTaskFiltering(t *testing.T) {
	router := setupRouter()

	// Register user
	registerPayload := map[string]string{
		"email":    "filteruser@example.com",
		"password": "password123",
	}
	jsonPayload, _ := json.Marshal(registerPayload)

	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var registerResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &registerResponse)
	token := registerResponse["token"].(string)

	// Create multiple tasks
	for i := 0; i < 3; i++ {
		taskPayload := map[string]interface{}{
			"title":       fmt.Sprintf("Task %d", i),
			"description": "Test task",
			"priority":    []string{"low", "medium", "high"}[i],
		}
		jsonPayload, _ = json.Marshal(taskPayload)

		req, _ = http.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}

	// Filter by priority
	req, _ = http.NewRequest("GET", "/tasks?priority=high", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var tasksResponse []interface{}
	json.Unmarshal(w.Body.Bytes(), &tasksResponse)
	assert.Greater(t, len(tasksResponse), 0)
}
