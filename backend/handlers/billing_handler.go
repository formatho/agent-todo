package handlers

import (
	"net/http"

	"github.com/formatho/agent-todo/backend/models"
	"github.com/formatho/agent-todo/backend/services"
	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v76"
)

// BillingHandler handles billing-related HTTP requests
type BillingHandler struct {
	billingSvc *services.BillingService
}

// NewBillingHandler creates a new billing handler
func NewBillingHandler(billingSvc *services.BillingService) *BillingHandler {
	return &BillingHandler{
		billingSvc: billingSvc,
	}
}

// GetSubscription returns the current user's subscription
// GET /api/billing/subscription
func (h *BillingHandler) GetSubscription(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	sub, err := h.billingSvc.GetUserSubscription(userID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to retrieve subscription",
		})
	}

	return c.JSON(fiber.Map{
		"subscription": sub,
		"limits": fiber.Map{
			"agents":        sub.AgentsLimit,
			"organizations": sub.OrganizationsLimit,
		},
	})
}

// CreateSubscription creates a new subscription for checkout
// POST /api/billing/subscribe
func (h *BillingHandler) CreateSubscription(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var req struct {
		Email string `json:"email" validate:"required,email"`
		Tier  string `json:"tier" validate:"required,oneof=starter pro enterprise"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	sub, err := h.billingSvc.CreateUserSubscription(userID, req.Email, services.SubscriptionTier(req.Tier))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"subscription": sub,
		"message":      "subscription created successfully",
	})
}

// HandleStripeWebhook processes webhooks from Stripe
// POST /api/billing/webhook
func (h *BillingHandler) HandleStripeWebhook(c *fiber.Ctx) error {
	payload := c.Body()
	sigHeader := c.Get("Stripe-Signature")

	event, err := stripe.ConstructEvent(payload, sigHeader, h.billingSvc.config.StripeWebhookSecret)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid webhook signature",
		})
	}

	// Handle the webhook event
	if err := h.billingSvc.HandleWebhook(event); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to process webhook",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
	})
}

// CancelSubscription cancels the current subscription
// POST /api/billing/cancel
func (h *BillingHandler) CancelSubscription(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	if err := h.billingSvc.CancelSubscription(userID); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "subscription canceled successfully",
	})
}

// UpdateSubscription updates the subscription tier
// PUT /api/billing/update
func (h *BillingHandler) UpdateSubscription(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var req struct {
		Tier string `json:"tier" validate:"required,oneof=free starter pro enterprise"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if err := h.billingSvc.UpdateSubscription(userID, services.SubscriptionTier(req.Tier)); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "subscription updated successfully",
	})
}

// CheckLimits checks if user can perform an action based on limits
// GET /api/billing/check-limits/:type
func (h *BillingHandler) CheckLimits(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	limitType := c.Params("type") // "agents" or "organizations"

	var err error
	switch limitType {
	case "agents":
		err = h.billingSvc.CheckAgentLimit(userID)
	case "organizations":
		err = h.billingSvc.CheckOrganizationLimit(userID)
	default:
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid limit type",
		})
	}

	if err != nil {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{
			"allowed": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"allowed": true,
	})
}

// GetPaymentHistory returns payment history for the user
// GET /api/billing/history
func (h *BillingHandler) GetPaymentHistory(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var payments []models.PaymentHistory
	if err := h.billingSvc.db.Where("user_id = ?", userID).Order("created_at DESC").Limit(50).Find(&payments).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to retrieve payment history",
		})
	}

	return c.JSON(fiber.Map{
		"payments": payments,
	})
}

// GetUsageStats returns usage statistics for the user
// GET /api/billing/usage
func (h *BillingHandler) GetUsageStats(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	// Get current subscription
	sub, err := h.billingSvc.GetUserSubscription(userID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to retrieve subscription",
		})
	}

	// Count current usage
	var agentCount int64
	var orgCount int64
	var taskCount int64

	h.billingSvc.db.Model(&models.Agent{}).Where("user_id = ?", userID).Count(&agentCount)
	h.billingSvc.db.Model(&models.Organisation{}).Where("owner_id = ?", userID).Count(&orgCount)
	h.billingSvc.db.Model(&models.Todo{}).Where("user_id = ?", userID).Count(&taskCount)

	return c.JSON(fiber.Map{
		"subscription": sub,
		"usage": fiber.Map{
			"agents":        int(agentCount),
			"organizations": int(orgCount),
			"tasks":         int(taskCount),
		},
		"limits": fiber.Map{
			"agents":        sub.AgentsLimit,
			"organizations": sub.OrganizationsLimit,
		},
		"percent_used": fiber.Map{
			"agents":        calculatePercent(int(agentCount), sub.AgentsLimit),
			"organizations": calculatePercent(int(orgCount), sub.OrganizationsLimit),
		},
	})
}

// calculatePercent calculates percentage used
func calculatePercent(current, limit int) float64 {
	if limit <= 0 {
		return 0 // Unlimited
	}
	percent := float64(current) / float64(limit) * 100
	if percent > 100 {
		return 100
	}
	return percent
}

// GetPricingTiers returns available pricing tiers
// GET /api/billing/pricing
func (h *BillingHandler) GetPricingTiers(c *fiber.Ctx) error {
	tiers := []fiber.Map{
		{
			"id":          "free",
			"name":        "Free",
			"price":       0,
			"price_id":    "",
			"features":    []string{"3 agents", "Basic task management", "1 organization", "Community support"},
			"recommended": false,
		},
		{
			"id":          "starter",
			"name":        "Starter",
			"price":       19,
			"price_id":    "price_starter_monthly",
			"features":    []string{"10 agents", "Advanced tasks", "3 organizations", "Email support", "Basic analytics"},
			"recommended": false,
		},
		{
			"id":          "pro",
			"name":        "Pro",
			"price":       49,
			"price_id":    "price_pro_monthly",
			"features":    []string{"Unlimited agents", "All features", "Unlimited organizations", "Priority support", "Custom integrations", "API access"},
			"recommended": true,
		},
		{
			"id":          "enterprise",
			"name":        "Enterprise",
			"price":       199,
			"price_id":    "price_enterprise_monthly",
			"features":    []string{"Everything in Pro", "SSO/SAML", "Advanced security", "Dedicated support", "White-label options", "SLA guarantee"},
			"recommended": false,
		},
	}

	return c.JSON(fiber.Map{
		"tiers": tiers,
	})
}
