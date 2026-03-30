package handlers

import (
	"net/http"

	"github.com/formatho/agent-todo/db"
	"github.com/formatho/agent-todo/middleware"
	"github.com/formatho/agent-todo/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FeedbackHandler struct{}

func NewFeedbackHandler() *FeedbackHandler {
	return &FeedbackHandler{}
}

// SubmitFeedbackRequest represents the request body for submitting feedback
type SubmitFeedbackRequest struct {
	TesterEmail  string           `json:"tester_email"`
	TesterName   string           `json:"tester_name"`
	FeedbackType models.FeedbackType `json:"feedback_type" binding:"required"`
	Title        string           `json:"title" binding:"required"`
	Description  string           `json:"description" binding:"required"`
	Priority     models.TaskPriority `json:"priority"`
	Page         string           `json:"page"`
	Rating       int              `json:"rating"`
}

// SubmitFeedback godoc
// @Summary Submit beta tester feedback
// @Description Submit feedback from beta testers (public endpoint)
// @Tags feedback
// @Accept json
// @Produce json
// @Param request body SubmitFeedbackRequest true "Feedback data"
// @Success 201 {object} models.BetaFeedback
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /feedback [post]
func (h *FeedbackHandler) SubmitFeedback(c *gin.Context) {
	var req SubmitFeedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate feedback type
	validTypes := map[models.FeedbackType]bool{
		models.FeedbackTypeBug:            true,
		models.FeedbackTypeFeatureRequest: true,
		models.FeedbackTypeGeneral:        true,
		models.FeedbackTypeImprovement:    true,
	}
	if !validTypes[req.FeedbackType] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid feedback type"})
		return
	}

	// Validate rating if provided
	if req.Rating != 0 && (req.Rating < 1 || req.Rating > 5) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Rating must be between 1 and 5"})
		return
	}

	// Get metadata
	userAgent := c.GetHeader("User-Agent")

	feedback := models.BetaFeedback{
		TesterEmail:  req.TesterEmail,
		TesterName:   req.TesterName,
		FeedbackType: req.FeedbackType,
		Title:        req.Title,
		Description:  req.Description,
		Priority:     req.Priority,
		Page:         req.Page,
		UserAgent:    userAgent,
		Status:       models.FeedbackStatusNew,
		Rating:       req.Rating,
	}

	if err := db.DB.Create(&feedback).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit feedback"})
		return
	}

	c.JSON(http.StatusCreated, feedback)
}

// ListFeedbackRequest represents query parameters for listing feedback
type ListFeedbackRequest struct {
	Status       *models.FeedbackStatus   `form:"status"`
	FeedbackType *models.FeedbackType     `form:"feedback_type"`
	Priority     *models.TaskPriority     `form:"priority"`
	Limit        int                      `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset       int                      `form:"offset" binding:"omitempty,min=0"`
}

// ListFeedback godoc
// @Summary List beta tester feedback
// @Description Get all feedback submissions (admin only)
// @Tags feedback
// @Produce json
// @Param status query string false "Filter by status"
// @Param feedback_type query string false "Filter by feedback type"
// @Param priority query string false "Filter by priority"
// @Param limit query int false "Number of results (default: 50)"
// @Param offset query int false "Offset for pagination (default: 0)"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /feedback [get]
func (h *FeedbackHandler) ListFeedback(c *gin.Context) {
	var req ListFeedbackRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set defaults
	if req.Limit == 0 {
		req.Limit = 50
	}

	query := db.DB.Model(&models.BetaFeedback{})

	// Apply filters
	if req.Status != nil {
		query = query.Where("status = ?", req.Status)
	}
	if req.FeedbackType != nil {
		query = query.Where("feedback_type = ?", req.FeedbackType)
	}
	if req.Priority != nil {
		query = query.Where("priority = ?", req.Priority)
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count feedback"})
		return
	}

	// Get feedback items
	var feedback []models.BetaFeedback
	if err := query.Order("created_at DESC").Limit(req.Limit).Offset(req.Offset).Find(&feedback).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list feedback"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"feedback": feedback,
		"total":    total,
		"limit":    req.Limit,
		"offset":   req.Offset,
	})
}

// GetFeedback godoc
// @Summary Get specific feedback
// @Description Get detailed information about a specific feedback submission (admin only)
// @Tags feedback
// @Produce json
// @Param id path string true "Feedback ID"
// @Success 200 {object} models.BetaFeedback
// @Failure 404 {object} map[string]string
// @Router /feedback/{id} [get]
func (h *FeedbackHandler) GetFeedback(c *gin.Context) {
	id := c.Param("id")
	feedbackID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid feedback ID"})
		return
	}

	var feedback models.BetaFeedback
	if err := db.DB.First(&feedback, "id = ?", feedbackID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Feedback not found"})
		return
	}

	c.JSON(http.StatusOK, feedback)
}

// UpdateFeedbackStatusRequest represents the request body for updating feedback status
type UpdateFeedbackStatusRequest struct {
	Status     models.FeedbackStatus `json:"status" binding:"required"`
	AdminNotes string                `json:"admin_notes"`
}

// UpdateFeedbackStatus godoc
// @Summary Update feedback status
// @Description Update the status of a feedback submission (admin only)
// @Tags feedback
// @Accept json
// @Produce json
// @Param id path string true "Feedback ID"
// @Param request body UpdateFeedbackStatusRequest true "Status update data"
// @Success 200 {object} models.BetaFeedback
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /feedback/{id}/status [patch]
func (h *FeedbackHandler) UpdateFeedbackStatus(c *gin.Context) {
	id := c.Param("id")
	feedbackID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid feedback ID"})
		return
	}

	var req UpdateFeedbackStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate status
	validStatuses := map[models.FeedbackStatus]bool{
		models.FeedbackStatusNew:          true,
		models.FeedbackStatusAcknowledged: true,
		models.FeedbackStatusInProgress:   true,
		models.FeedbackStatusResolved:     true,
		models.FeedbackStatusClosed:       true,
	}
	if !validStatuses[req.Status] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	var feedback models.BetaFeedback
	if err := db.DB.First(&feedback, "id = ?", feedbackID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Feedback not found"})
		return
	}

	updates := map[string]interface{}{
		"status": req.Status,
	}
	if req.AdminNotes != "" {
		updates["admin_notes"] = req.AdminNotes
	}

	if err := db.DB.Model(&feedback).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update feedback"})
		return
	}

	// Reload to get updated data
	db.DB.First(&feedback, "id = ?", feedbackID)

	c.JSON(http.StatusOK, feedback)
}

// UpdateFeedbackNotesRequest represents the request body for updating feedback notes
type UpdateFeedbackNotesRequest struct {
	AdminNotes string `json:"admin_notes" binding:"required"`
}

// UpdateFeedbackNotes godoc
// @Summary Update feedback admin notes
// @Description Update admin notes for a feedback submission (admin only)
// @Tags feedback
// @Accept json
// @Produce json
// @Param id path string true "Feedback ID"
// @Param request body UpdateFeedbackNotesRequest true "Admin notes"
// @Success 200 {object} models.BetaFeedback
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /feedback/{id}/notes [patch]
func (h *FeedbackHandler) UpdateFeedbackNotes(c *gin.Context) {
	id := c.Param("id")
	feedbackID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid feedback ID"})
		return
	}

	var req UpdateFeedbackNotesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var feedback models.BetaFeedback
	if err := db.DB.First(&feedback, "id = ?", feedbackID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Feedback not found"})
		return
	}

	if err := db.DB.Model(&feedback).Update("admin_notes", req.AdminNotes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update notes"})
		return
	}

	// Reload to get updated data
	db.DB.First(&feedback, "id = ?", feedbackID)

	c.JSON(http.StatusOK, feedback)
}

// GetFeedbackStats godoc
// @Summary Get feedback statistics
// @Description Get aggregated statistics about feedback submissions (admin only)
// @Tags feedback
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /feedback/stats [get]
func (h *FeedbackHandler) GetFeedbackStats(c *gin.Context) {
	stats := make(map[string]interface{})

	// Total feedback count
	var total int64
	db.DB.Model(&models.BetaFeedback{}).Count(&total)
	stats["total"] = total

	// Count by status
	var statusStats []struct {
		Status models.FeedbackStatus
		Count  int64
	}
	db.DB.Model(&models.BetaFeedback{}).Select("status, count(*) as count").Group("status").Scan(&statusStats)
	stats["by_status"] = statusStats

	// Count by type
	var typeStats []struct {
		FeedbackType models.FeedbackType
		Count        int64
	}
	db.DB.Model(&models.BetaFeedback{}).Select("feedback_type, count(*) as count").Group("feedback_type").Scan(&typeStats)
	stats["by_type"] = typeStats

	// Average rating
	var avgRating float64
	db.DB.Model(&models.BetaFeedback{}).Where("rating > 0").Select("AVG(rating)").Scan(&avgRating)
	stats["average_rating"] = avgRating

	// Recent feedback (last 7 days)
	var recentCount int64
	db.DB.Model(&models.BetaFeedback{}).Where("created_at > NOW() - INTERVAL '7 days'").Count(&recentCount)
	stats["recent_count"] = recentCount

	c.JSON(http.StatusOK, stats)
}
