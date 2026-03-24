package handlers

import (
	"net/http"

	"github.com/formatho/agent-todo/db"
	"github.com/formatho/agent-todo/models"
	"github.com/formatho/agent-todo/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SubscriptionHandler handles subscription-related requests
type SubscriptionHandler struct {
	sequenceService *services.EmailSequenceService
}

// NewSubscriptionHandler creates a new subscription handler
func NewSubscriptionHandler() *SubscriptionHandler {
	return &SubscriptionHandler{
		sequenceService: services.NewEmailSequenceService(db.DB),
	}
}

// StartTrialRequest represents a request to start a trial
type StartTrialRequest struct {
	OrganisationID string `json:"organisation_id" binding:"required"`
}

// StartTrial starts a trial for an organisation
// @Summary Start trial
// @Description Start a 14-day trial for an organisation
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param request body StartTrialRequest true "Start trial request"
// @Success 200 {object} models.Subscription
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/trial [post]
// @Security Bearer
func (h *SubscriptionHandler) StartTrial(c *gin.Context) {
	var req StartTrialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orgID, err := uuid.Parse(req.OrganisationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organisation ID"})
		return
	}

	if err := h.sequenceService.StartTrial(orgID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	subscription, err := h.sequenceService.GetSubscription(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subscription)
}

// GetSubscription returns the current subscription for an organisation
// @Summary Get subscription
// @Description Get subscription details for an organisation
// @Tags subscriptions
// @Produce json
// @Param organisation_id path string true "Organisation ID"
// @Success 200 {object} models.Subscription
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /subscriptions/{organisation_id} [get]
// @Security Bearer
func (h *SubscriptionHandler) GetSubscription(c *gin.Context) {
	orgID, err := uuid.Parse(c.Param("organisation_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organisation ID"})
		return
	}

	subscription, err := h.sequenceService.GetSubscription(orgID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
		return
	}

	c.JSON(http.StatusOK, subscription)
}

// CancelTrial cancels an ongoing trial
// @Summary Cancel trial
// @Description Cancel an ongoing trial for an organisation
// @Tags subscriptions
// @Produce json
// @Param organisation_id path string true "Organisation ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{organisation_id}/cancel [post]
// @Security Bearer
func (h *SubscriptionHandler) CancelTrial(c *gin.Context) {
	orgID, err := uuid.Parse(c.Param("organisation_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organisation ID"})
		return
	}

	if err := h.sequenceService.CancelTrial(orgID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Trial cancelled successfully"})
}

// EmailHandler handles email-related requests
type EmailHandler struct {
	sequenceService *services.EmailSequenceService
}

// NewEmailHandler creates a new email handler
func NewEmailHandler() *EmailHandler {
	return &EmailHandler{
		sequenceService: services.NewEmailSequenceService(db.DB),
	}
}

// ProcessQueueRequest represents a request to process the email queue
type ProcessQueueRequest struct {
	BatchSize int `json:"batch_size"`
}

// ProcessQueueResponse represents the response from processing the queue
type ProcessQueueResponse struct {
	Sent   int `json:"sent"`
	Failed int `json:"failed"`
}

// ProcessEmailQueue processes pending emails in the queue
// @Summary Process email queue
// @Description Process pending scheduled emails (for cron jobs)
// @Tags emails
// @Accept json
// @Produce json
// @Param request body ProcessQueueRequest true "Process queue request"
// @Success 200 {object} ProcessQueueResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /emails/process [post]
func (h *EmailHandler) ProcessEmailQueue(c *gin.Context) {
	var req ProcessQueueRequest
	batchSize := 50 // Default
	if err := c.ShouldBindJSON(&req); err == nil && req.BatchSize > 0 {
		batchSize = req.BatchSize
	}

	sent, failed, err := h.sequenceService.ProcessEmailQueue(batchSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ProcessQueueResponse{
		Sent:   sent,
		Failed: failed,
	})
}

// GetEmailQueue returns pending emails in the queue
// @Summary Get email queue
// @Description Get pending scheduled emails
// @Tags emails
// @Produce json
// @Success 200 {array} models.EmailQueue
// @Failure 401 {object} map[string]string
// @Router /emails/queue [get]
// @Security Bearer
func (h *EmailHandler) GetEmailQueue(c *gin.Context) {
	var queue []models.EmailQueue
	if err := db.DB.Limit(100).Order("scheduled_at ASC").Find(&queue).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, queue)
}

// GetEmailLogs returns sent email logs
// @Summary Get email logs
// @Description Get history of sent emails
// @Tags emails
// @Produce json
// @Param limit query int false "Limit" default(50)
// @Success 200 {array} models.EmailLog
// @Failure 401 {object} map[string]string
// @Router /emails/logs [get]
// @Security Bearer
func (h *EmailHandler) GetEmailLogs(c *gin.Context) {
	limit := 50
	if l := c.Query("limit"); l != "" {
		// Parse limit if provided
	}

	var logs []models.EmailLog
	if err := db.DB.Limit(limit).Order("sent_at DESC").Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, logs)
}

// SeedSequenceRequest represents a request to seed email sequences
type SeedSequenceRequest struct {
	Sequence string `json:"sequence"` // "trial_conversion" or empty for all
}

// SeedSeeds seeds the default email sequences and templates
// @Summary Seed email sequences
// @Description Seed default email sequences and templates
// @Tags emails
// @Accept json
// @Produce json
// @Param request body SeedSequenceRequest false "Seed request"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /emails/seed [post]
// @Security Bearer
func (h *EmailHandler) SeedSequences(c *gin.Context) {
	if err := h.sequenceService.SeedTrialConversionSequence(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email sequences seeded successfully"})
}

// GetEmailSequences returns all email sequences
// @Summary Get email sequences
// @Description Get all configured email sequences
// @Tags emails
// @Produce json
// @Success 200 {array} models.EmailSequence
// @Failure 401 {object} map[string]string
// @Router /emails/sequences [get]
// @Security Bearer
func (h *EmailHandler) GetEmailSequences(c *gin.Context) {
	var sequences []models.EmailSequence
	if err := db.DB.Preload("Steps.Template").Find(&sequences).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sequences)
}
